// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package console

import (
	"context"
	"math"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/memory"
	"storj.io/common/storj"
	"storj.io/private/version"
	"storj.io/storj/private/date"
	"storj.io/storj/private/version/checker"
	"storj.io/storj/storagenode/bandwidth"
	"storj.io/storj/storagenode/contact"
	"storj.io/storj/storagenode/operator"
	"storj.io/storj/storagenode/payouts/estimatedpayouts"
	"storj.io/storj/storagenode/pieces"
	"storj.io/storj/storagenode/pricing"
	"storj.io/storj/storagenode/reputation"
	"storj.io/storj/storagenode/satellites"
	"storj.io/storj/storagenode/storageusage"
	"storj.io/storj/storagenode/trust"
)

var (
	// SNOServiceErr defines sno service error.
	SNOServiceErr = errs.Class("storage node dashboard service error")

	mon = monkit.Package()
)

// Service is handling storage node operator related logic.
//
// architecture: Service
type Service struct {
	log            *zap.Logger
	trust          *trust.Pool
	usageCache     *pieces.BlobsUsageCache
	bandwidthDB    bandwidth.DB
	reputationDB   reputation.DB
	storageUsageDB storageusage.DB
	pricingDB      pricing.DB
	satelliteDB    satellites.DB
	pieceStore     *pieces.Store
	contact        *contact.Service

	estimation *estimatedpayouts.Service
	version    *checker.Service
	pingStats  *contact.PingStats

	allocatedDiskSpace memory.Size

	walletAddress  string
	walletFeatures operator.WalletFeatures
	startedAt      time.Time
	versionInfo    version.Info
}

// NewService returns new instance of Service.
func NewService(log *zap.Logger, bandwidth bandwidth.DB, pieceStore *pieces.Store, version *checker.Service,
	allocatedDiskSpace memory.Size, walletAddress string, versionInfo version.Info, trust *trust.Pool,
	reputationDB reputation.DB, storageUsageDB storageusage.DB, pricingDB pricing.DB, satelliteDB satellites.DB,
	pingStats *contact.PingStats, contact *contact.Service, estimation *estimatedpayouts.Service, usageCache *pieces.BlobsUsageCache, walletFeatures operator.WalletFeatures) (*Service, error) {
	if log == nil {
		return nil, errs.New("log can't be nil")
	}

	if usageCache == nil {
		return nil, errs.New("usage cache can't be nil")
	}

	if bandwidth == nil {
		return nil, errs.New("bandwidth can't be nil")
	}

	if pieceStore == nil {
		return nil, errs.New("pieceStore can't be nil")
	}

	if version == nil {
		return nil, errs.New("version can't be nil")
	}

	if pingStats == nil {
		return nil, errs.New("pingStats can't be nil")
	}

	if contact == nil {
		return nil, errs.New("contact service can't be nil")
	}

	if estimation == nil {
		return nil, errs.New("estimation service can't be nil")
	}

	return &Service{
		log:                log,
		trust:              trust,
		usageCache:         usageCache,
		bandwidthDB:        bandwidth,
		reputationDB:       reputationDB,
		storageUsageDB:     storageUsageDB,
		pricingDB:          pricingDB,
		satelliteDB:        satelliteDB,
		pieceStore:         pieceStore,
		version:            version,
		pingStats:          pingStats,
		allocatedDiskSpace: allocatedDiskSpace,
		contact:            contact,
		estimation:         estimation,
		walletAddress:      walletAddress,
		startedAt:          time.Now(),
		versionInfo:        versionInfo,
		walletFeatures:     walletFeatures,
	}, nil
}

// SatelliteInfo encapsulates satellite ID and disqualification.
type SatelliteInfo struct {
	ID                 storj.NodeID `json:"id"`
	URL                string       `json:"url"`
	Disqualified       *time.Time   `json:"disqualified"`
	Suspended          *time.Time   `json:"suspended"`
	CurrentStorageUsed int64        `json:"currentStorageUsed"`
}

// Dashboard encapsulates dashboard stale data.
type Dashboard struct {
	NodeID         storj.NodeID `json:"nodeID"`
	Wallet         string       `json:"wallet"`
	WalletFeatures []string     `json:"walletFeatures"`

	Satellites []SatelliteInfo `json:"satellites"`

	DiskSpace DiskSpaceInfo `json:"diskSpace"`
	Bandwidth BandwidthInfo `json:"bandwidth"`

	LastPinged time.Time `json:"lastPinged"`

	Version        version.SemVer `json:"version"`
	AllowedVersion version.SemVer `json:"allowedVersion"`
	UpToDate       bool           `json:"upToDate"`

	StartedAt time.Time `json:"startedAt"`
}

// GetDashboardData returns stale dashboard data.
func (s *Service) GetDashboardData(ctx context.Context) (_ *Dashboard, err error) {
	defer mon.Task()(&ctx)(&err)
	data := new(Dashboard)

	data.NodeID = s.contact.Local().ID
	data.Wallet = s.walletAddress
	data.WalletFeatures = s.walletFeatures
	data.Version = s.versionInfo.Version
	data.StartedAt = s.startedAt

	data.LastPinged = s.pingStats.WhenLastPinged()
	data.AllowedVersion, data.UpToDate = s.version.IsAllowed(ctx)

	stats, err := s.reputationDB.All(ctx)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	for _, rep := range stats {
		url, err := s.trust.GetNodeURL(ctx, rep.SatelliteID)
		if err != nil {
			s.log.Warn("unable to get Satellite URL", zap.String("Satellite ID", rep.SatelliteID.String()),
				zap.Error(SNOServiceErr.Wrap(err)))
			continue
		}
		_, currentStorageUsed, err := s.usageCache.SpaceUsedBySatellite(ctx, rep.SatelliteID)
		if err != nil {
			s.log.Warn("unable to get Satellite Current Storage Used", zap.String("Satellite ID", rep.SatelliteID.String()),
				zap.Error(SNOServiceErr.Wrap(err)))
			continue
		}

		data.Satellites = append(data.Satellites,
			SatelliteInfo{
				ID:                 rep.SatelliteID,
				Disqualified:       rep.DisqualifiedAt,
				Suspended:          rep.SuspendedAt,
				URL:                url.Address,
				CurrentStorageUsed: currentStorageUsed,
			},
		)
	}

	pieceTotal, _, err := s.pieceStore.SpaceUsedForPieces(ctx)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	trash, err := s.pieceStore.SpaceUsedForTrash(ctx)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	bandwidthUsage, err := s.bandwidthDB.MonthSummary(ctx, time.Now())
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	data.DiskSpace = DiskSpaceInfo{
		Used:      pieceTotal,
		Available: s.allocatedDiskSpace.Int64(),
		Trash:     trash,
	}

	overused := s.allocatedDiskSpace.Int64() - pieceTotal - trash
	if overused < 0 {
		data.DiskSpace.Overused = int64(math.Abs(float64(overused)))
	}

	data.Bandwidth = BandwidthInfo{
		Used: bandwidthUsage,
	}

	return data, nil
}

// PriceModel is a satellite prices for storagenode usage TB/H.
type PriceModel struct {
	EgressBandwidth int64
	RepairBandwidth int64
	AuditBandwidth  int64
	DiskSpace       int64
}

// Satellite encapsulates satellite related data.
type Satellite struct {
	ID                 storj.NodeID            `json:"id"`
	StorageDaily       []storageusage.Stamp    `json:"storageDaily"`
	BandwidthDaily     []bandwidth.UsageRollup `json:"bandwidthDaily"`
	StorageSummary     float64                 `json:"storageSummary"`
	BandwidthSummary   int64                   `json:"bandwidthSummary"`
	EgressSummary      int64                   `json:"egressSummary"`
	IngressSummary     int64                   `json:"ingressSummary"`
	CurrentStorageUsed int64                   `json:"currentStorageUsed"`
	Audits             Audits                  `json:"audits"`
	AuditHistory       reputation.AuditHistory `json:"auditHistory"`
	PriceModel         PriceModel              `json:"priceModel"`
	NodeJoinedAt       time.Time               `json:"nodeJoinedAt"`
}

// GetSatelliteData returns satellite related data.
func (s *Service) GetSatelliteData(ctx context.Context, satelliteID storj.NodeID) (_ *Satellite, err error) {
	defer mon.Task()(&ctx)(&err)
	from, to := date.MonthBoundary(time.Now().UTC())

	bandwidthDaily, err := s.bandwidthDB.GetDailySatelliteRollups(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	storageDaily, err := s.storageUsageDB.GetDaily(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	bandwidthSummary, err := s.bandwidthDB.SatelliteSummary(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	egressSummary, err := s.bandwidthDB.SatelliteEgressSummary(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	ingressSummary, err := s.bandwidthDB.SatelliteIngressSummary(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	storageSummary, err := s.storageUsageDB.SatelliteSummary(ctx, satelliteID, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	_, currentStorageUsed, err := s.usageCache.SpaceUsedBySatellite(ctx, satelliteID)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	rep, err := s.reputationDB.Get(ctx, satelliteID)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	pricingModel, err := s.pricingDB.Get(ctx, satelliteID)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	satellitePricing := PriceModel{
		EgressBandwidth: pricingModel.EgressBandwidth,
		RepairBandwidth: pricingModel.RepairBandwidth,
		AuditBandwidth:  pricingModel.AuditBandwidth,
		DiskSpace:       pricingModel.DiskSpace,
	}

	url, err := s.trust.GetNodeURL(ctx, satelliteID)
	if err != nil {
		s.log.Warn("unable to get Satellite URL", zap.String("Satellite ID", satelliteID.String()),
			zap.Error(SNOServiceErr.Wrap(err)))
	}

	return &Satellite{
		ID:                 satelliteID,
		StorageDaily:       storageDaily,
		BandwidthDaily:     bandwidthDaily,
		StorageSummary:     storageSummary,
		BandwidthSummary:   bandwidthSummary.Total(),
		CurrentStorageUsed: currentStorageUsed,
		EgressSummary:      egressSummary.Total(),
		IngressSummary:     ingressSummary.Total(),
		Audits: Audits{
			AuditScore:      rep.Audit.Score,
			SuspensionScore: rep.Audit.UnknownScore,
			OnlineScore:     rep.OnlineScore,
			SatelliteName:   url.Address,
		},
		AuditHistory: reputation.GetAuditHistoryFromPB(rep.AuditHistory),
		PriceModel:   satellitePricing,
		NodeJoinedAt: rep.JoinedAt,
	}, nil
}

// Satellites represents consolidated data across all satellites.
type Satellites struct {
	StorageDaily     []storageusage.Stamp    `json:"storageDaily"`
	BandwidthDaily   []bandwidth.UsageRollup `json:"bandwidthDaily"`
	StorageSummary   float64                 `json:"storageSummary"`
	BandwidthSummary int64                   `json:"bandwidthSummary"`
	EgressSummary    int64                   `json:"egressSummary"`
	IngressSummary   int64                   `json:"ingressSummary"`
	EarliestJoinedAt time.Time               `json:"earliestJoinedAt"`
	Audits           []Audits                `json:"audits"`
}

// Audits represents audit, suspension and online scores of SNO across all satellites.
type Audits struct {
	AuditScore      float64 `json:"auditScore"`
	SuspensionScore float64 `json:"suspensionScore"`
	OnlineScore     float64 `json:"onlineScore"`
	SatelliteName   string  `json:"satelliteName"`
}

// GetAllSatellitesData returns bandwidth and storage daily usage consolidate
// among all satellites from the node's trust pool.
func (s *Service) GetAllSatellitesData(ctx context.Context) (_ *Satellites, err error) {
	defer mon.Task()(&ctx)(nil)
	from, to := date.MonthBoundary(time.Now().UTC())

	var audits []Audits

	bandwidthDaily, err := s.bandwidthDB.GetDailyRollups(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	storageDaily, err := s.storageUsageDB.GetDailyTotal(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	bandwidthSummary, err := s.bandwidthDB.Summary(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	egressSummary, err := s.bandwidthDB.EgressSummary(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	ingressSummary, err := s.bandwidthDB.IngressSummary(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	storageSummary, err := s.storageUsageDB.Summary(ctx, from, to)
	if err != nil {
		return nil, SNOServiceErr.Wrap(err)
	}

	satellitesIDs := s.trust.GetSatellites(ctx)
	joinedAt := time.Now().UTC()

	for i := 0; i < len(satellitesIDs); i++ {
		stats, err := s.reputationDB.Get(ctx, satellitesIDs[i])
		if err != nil {
			return nil, SNOServiceErr.Wrap(err)
		}

		url, err := s.trust.GetNodeURL(ctx, satellitesIDs[i])
		if err != nil {
			s.log.Warn("unable to get Satellite URL", zap.String("Satellite ID", satellitesIDs[i].String()),
				zap.Error(SNOServiceErr.Wrap(err)))
			continue
		}

		audits = append(audits, Audits{
			AuditScore:      stats.Audit.Score,
			SuspensionScore: stats.Audit.UnknownScore,
			OnlineScore:     stats.OnlineScore,
			SatelliteName:   url.Address,
		})
		if !stats.JoinedAt.IsZero() && stats.JoinedAt.Before(joinedAt) {
			joinedAt = stats.JoinedAt
		}
	}

	return &Satellites{
		StorageDaily:     storageDaily,
		BandwidthDaily:   bandwidthDaily,
		StorageSummary:   storageSummary,
		BandwidthSummary: bandwidthSummary.Total(),
		EgressSummary:    egressSummary.Total(),
		IngressSummary:   ingressSummary.Total(),
		EarliestJoinedAt: joinedAt,
		Audits:           audits,
	}, nil
}

// GetSatelliteEstimatedPayout returns estimated payouts for current and previous months for selected satellite.
func (s *Service) GetSatelliteEstimatedPayout(ctx context.Context, satelliteID storj.NodeID, now time.Time) (estimatedPayout estimatedpayouts.EstimatedPayout, err error) {
	estimatedPayout, err = s.estimation.GetSatelliteEstimatedPayout(ctx, satelliteID, now)
	if err != nil {
		return estimatedpayouts.EstimatedPayout{}, SNOServiceErr.Wrap(err)
	}

	return estimatedPayout, nil
}

// GetAllSatellitesEstimatedPayout returns estimated payouts for current and previous months for all satellites.
func (s *Service) GetAllSatellitesEstimatedPayout(ctx context.Context, now time.Time) (estimatedPayout estimatedpayouts.EstimatedPayout, err error) {
	estimatedPayout, err = s.estimation.GetAllSatellitesEstimatedPayout(ctx, now)
	if err != nil {
		return estimatedpayouts.EstimatedPayout{}, SNOServiceErr.Wrap(err)
	}

	return estimatedPayout, nil
}

// VerifySatelliteID verifies if the satellite belongs to the trust pool.
func (s *Service) VerifySatelliteID(ctx context.Context, satelliteID storj.NodeID) (err error) {
	defer mon.Task()(&ctx)(&err)

	err = s.trust.VerifySatelliteID(ctx, satelliteID)
	if err != nil {
		return SNOServiceErr.Wrap(err)
	}

	return nil
}
