// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package tally

import (
	"context"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/storj"
	"storj.io/common/sync2"
	"storj.io/common/uuid"
	"storj.io/storj/satellite/accounting"
	"storj.io/storj/satellite/metainfo"
	"storj.io/storj/satellite/metainfo/metabase"
)

// Error is a standard error class for this package.
var (
	Error = errs.Class("tally error")
	mon   = monkit.Package()
)

// Config contains configurable values for the tally service.
type Config struct {
	Interval time.Duration `help:"how frequently the tally service should run" releaseDefault:"1h" devDefault:"30s"`
}

// Service is the tally service for data stored on each storage node
//
// architecture: Chore
type Service struct {
	log  *zap.Logger
	Loop *sync2.Cycle

	metainfoLoop            *metainfo.Loop
	liveAccounting          accounting.Cache
	storagenodeAccountingDB accounting.StoragenodeAccounting
	projectAccountingDB     accounting.ProjectAccounting
	nowFn                   func() time.Time
}

// New creates a new tally Service.
func New(log *zap.Logger, sdb accounting.StoragenodeAccounting, pdb accounting.ProjectAccounting, liveAccounting accounting.Cache, metainfoLoop *metainfo.Loop, interval time.Duration) *Service {
	return &Service{
		log:  log,
		Loop: sync2.NewCycle(interval),

		metainfoLoop:            metainfoLoop,
		liveAccounting:          liveAccounting,
		storagenodeAccountingDB: sdb,
		projectAccountingDB:     pdb,
		nowFn:                   time.Now,
	}
}

// Run the tally service loop.
func (service *Service) Run(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	return service.Loop.Run(ctx, func(ctx context.Context) error {
		err := service.Tally(ctx)
		if err != nil {
			service.log.Error("tally failed", zap.Error(err))
		}
		return nil
	})
}

// Close stops the service and releases any resources.
func (service *Service) Close() error {
	service.Loop.Close()
	return nil
}

// SetNow allows tests to have the Service act as if the current time is whatever
// they want. This avoids races and sleeping, making tests more reliable and efficient.
func (service *Service) SetNow(now func() time.Time) {
	service.nowFn = now
}

// Tally calculates data-at-rest usage once.
func (service *Service) Tally(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	initialLiveTotals, err := service.liveAccounting.GetAllProjectTotals(ctx)
	if err != nil {
		return Error.Wrap(err)
	}
	// Fetch when the last tally happened so we can roughly calculate the byte-hours.
	lastTime, err := service.storagenodeAccountingDB.LastTimestamp(ctx, accounting.LastAtRestTally)
	if err != nil {
		return Error.Wrap(err)
	}
	if lastTime.IsZero() {
		lastTime = service.nowFn()
	}

	// add up all nodes and buckets
	observer := NewObserver(service.log.Named("observer"), service.nowFn())
	err = service.metainfoLoop.Join(ctx, observer)
	if err != nil {
		return Error.Wrap(err)
	}
	finishTime := service.nowFn()

	// calculate byte hours, not just bytes
	hours := time.Since(lastTime).Hours()
	for id := range observer.Node {
		observer.Node[id] *= hours
	}

	// save the new results
	var errAtRest, errBucketInfo error
	if len(observer.Node) > 0 {
		err = service.storagenodeAccountingDB.SaveTallies(ctx, finishTime, observer.Node)
		if err != nil {
			errAtRest = errs.New("StorageNodeAccounting.SaveTallies failed: %v", err)
		}
	}

	if len(observer.Bucket) > 0 {
		// record bucket tallies to DB
		err = service.projectAccountingDB.SaveTallies(ctx, finishTime, observer.Bucket)
		if err != nil {
			errAtRest = errs.New("ProjectAccounting.SaveTallies failed: %v", err)
		}

		// update live accounting totals
		tallyProjectTotals := projectTotalsFromBuckets(observer.Bucket)
		latestLiveTotals, err := service.liveAccounting.GetAllProjectTotals(ctx)
		if err != nil {
			return Error.Wrap(err)
		}

		// empty projects are not returned by the metainfo observer. If a project exists
		// in live accounting, but not in tally projects, we would not update it in live accounting.
		// Thus, we add them and set the total to 0.
		for projectID := range latestLiveTotals {
			if _, ok := tallyProjectTotals[projectID]; !ok {
				tallyProjectTotals[projectID] = 0
			}
		}

		for projectID, tallyTotal := range tallyProjectTotals {
			delta := latestLiveTotals[projectID] - initialLiveTotals[projectID]
			if delta < 0 {
				delta = 0
			}
			err = service.liveAccounting.AddProjectStorageUsage(ctx, projectID, -latestLiveTotals[projectID]+tallyTotal+(delta/2))
			if err != nil {
				return Error.Wrap(err)
			}
		}
	}

	// report bucket metrics
	if len(observer.Bucket) > 0 {
		var total accounting.BucketTally
		for _, bucket := range observer.Bucket {
			monAccounting.IntVal("bucket.objects").Observe(bucket.ObjectCount)

			monAccounting.IntVal("bucket.segments").Observe(bucket.Segments())
			monAccounting.IntVal("bucket.inline_segments").Observe(bucket.InlineSegments)
			monAccounting.IntVal("bucket.remote_segments").Observe(bucket.RemoteSegments)

			monAccounting.IntVal("bucket.bytes").Observe(bucket.Bytes())
			monAccounting.IntVal("bucket.inline_bytes").Observe(bucket.InlineBytes)
			monAccounting.IntVal("bucket.remote_bytes").Observe(bucket.RemoteBytes)
			total.Combine(bucket)
		}
		monAccounting.IntVal("total.objects").Observe(total.ObjectCount) //mon:locked

		monAccounting.IntVal("total.segments").Observe(total.Segments())            //mon:locked
		monAccounting.IntVal("total.inline_segments").Observe(total.InlineSegments) //mon:locked
		monAccounting.IntVal("total.remote_segments").Observe(total.RemoteSegments) //mon:locked

		monAccounting.IntVal("total.bytes").Observe(total.Bytes())            //mon:locked
		monAccounting.IntVal("total.inline_bytes").Observe(total.InlineBytes) //mon:locked
		monAccounting.IntVal("total.remote_bytes").Observe(total.RemoteBytes) //mon:locked
	}

	// return errors if something went wrong.
	return errs.Combine(errAtRest, errBucketInfo)
}

var _ metainfo.Observer = (*Observer)(nil)

// Observer observes metainfo and adds up tallies for nodes and buckets.
type Observer struct {
	Now    time.Time
	Log    *zap.Logger
	Node   map[storj.NodeID]float64
	Bucket map[metabase.BucketLocation]*accounting.BucketTally
}

// NewObserver returns an metainfo loop observer that adds up totals for buckets and nodes.
// The now argument controls when the observer considers pointers to be expired.
func NewObserver(log *zap.Logger, now time.Time) *Observer {
	return &Observer{
		Now:    now,
		Log:    log,
		Node:   make(map[storj.NodeID]float64),
		Bucket: make(map[metabase.BucketLocation]*accounting.BucketTally),
	}
}

// ensureBucket returns bucket corresponding to the passed in path.
func (observer *Observer) ensureBucket(ctx context.Context, location metabase.ObjectLocation) *accounting.BucketTally {
	bucketLocation := location.Bucket()
	bucket, exists := observer.Bucket[bucketLocation]
	if !exists {
		bucket = &accounting.BucketTally{}
		bucket.BucketLocation = bucketLocation
		observer.Bucket[bucketLocation] = bucket
	}

	return bucket
}

// Object is called for each object once.
func (observer *Observer) Object(ctx context.Context, object *metainfo.Object) (err error) {
	if object.Expired(observer.Now) {
		return nil
	}

	bucket := observer.ensureBucket(ctx, object.Location)
	bucket.ObjectCount++

	return nil
}

// InlineSegment is called for each inline segment.
func (observer *Observer) InlineSegment(ctx context.Context, segment *metainfo.Segment) (err error) {
	if segment.Expired(observer.Now) {
		return nil
	}

	bucket := observer.ensureBucket(ctx, segment.Location.Object())
	bucket.InlineSegments++
	bucket.InlineBytes += int64(segment.DataSize)
	bucket.MetadataSize += int64(segment.MetadataSize)

	return nil
}

// RemoteSegment is called for each remote segment.
func (observer *Observer) RemoteSegment(ctx context.Context, segment *metainfo.Segment) (err error) {
	if segment.Expired(observer.Now) {
		return nil
	}

	bucket := observer.ensureBucket(ctx, segment.Location.Object())
	bucket.RemoteSegments++
	bucket.RemoteBytes += int64(segment.DataSize)
	bucket.MetadataSize += int64(segment.MetadataSize)

	// add node info
	minimumRequired := segment.Redundancy.RequiredShares

	if minimumRequired <= 0 {
		observer.Log.Error("failed sanity check", zap.ByteString("key", segment.Location.Encode()))
		return nil
	}

	pieceSize := float64(segment.DataSize / int(minimumRequired)) // TODO: Add this as a method to RedundancyScheme

	for _, piece := range segment.Pieces {
		observer.Node[piece.StorageNode] += pieceSize
	}

	return nil
}

func projectTotalsFromBuckets(buckets map[metabase.BucketLocation]*accounting.BucketTally) map[uuid.UUID]int64 {
	projectTallyTotals := make(map[uuid.UUID]int64)
	for _, bucket := range buckets {
		projectTallyTotals[bucket.ProjectID] += (bucket.InlineBytes + bucket.RemoteBytes)
	}
	return projectTallyTotals
}

// using custom name to avoid breaking monitoring.
var monAccounting = monkit.ScopeNamed("storj.io/storj/satellite/accounting")
