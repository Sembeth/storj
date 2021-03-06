// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.20
// source: multinode.proto

package multinodepb

import (
	bytes "bytes"
	context "context"
	errors "errors"

	jsonpb "github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"

	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_multinode_proto struct{}

func (drpcEncoding_File_multinode_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_multinode_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_multinode_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	var buf bytes.Buffer
	err := new(jsonpb.Marshaler).Marshal(&buf, msg.(proto.Message))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (drpcEncoding_File_multinode_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(buf), msg.(proto.Message))
}

type DRPCStorageClient interface {
	DRPCConn() drpc.Conn

	DiskSpace(ctx context.Context, in *DiskSpaceRequest) (*DiskSpaceResponse, error)
}

type drpcStorageClient struct {
	cc drpc.Conn
}

func NewDRPCStorageClient(cc drpc.Conn) DRPCStorageClient {
	return &drpcStorageClient{cc}
}

func (c *drpcStorageClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcStorageClient) DiskSpace(ctx context.Context, in *DiskSpaceRequest) (*DiskSpaceResponse, error) {
	out := new(DiskSpaceResponse)
	err := c.cc.Invoke(ctx, "/multinode.Storage/DiskSpace", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCStorageServer interface {
	DiskSpace(context.Context, *DiskSpaceRequest) (*DiskSpaceResponse, error)
}

type DRPCStorageUnimplementedServer struct{}

func (s *DRPCStorageUnimplementedServer) DiskSpace(context.Context, *DiskSpaceRequest) (*DiskSpaceResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

type DRPCStorageDescription struct{}

func (DRPCStorageDescription) NumMethods() int { return 1 }

func (DRPCStorageDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/multinode.Storage/DiskSpace", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCStorageServer).
					DiskSpace(
						ctx,
						in1.(*DiskSpaceRequest),
					)
			}, DRPCStorageServer.DiskSpace, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterStorage(mux drpc.Mux, impl DRPCStorageServer) error {
	return mux.Register(impl, DRPCStorageDescription{})
}

type DRPCStorage_DiskSpaceStream interface {
	drpc.Stream
	SendAndClose(*DiskSpaceResponse) error
}

type drpcStorage_DiskSpaceStream struct {
	drpc.Stream
}

func (x *drpcStorage_DiskSpaceStream) SendAndClose(m *DiskSpaceResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCBandwidthClient interface {
	DRPCConn() drpc.Conn

	MonthSummary(ctx context.Context, in *BandwidthMonthSummaryRequest) (*BandwidthMonthSummaryResponse, error)
}

type drpcBandwidthClient struct {
	cc drpc.Conn
}

func NewDRPCBandwidthClient(cc drpc.Conn) DRPCBandwidthClient {
	return &drpcBandwidthClient{cc}
}

func (c *drpcBandwidthClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcBandwidthClient) MonthSummary(ctx context.Context, in *BandwidthMonthSummaryRequest) (*BandwidthMonthSummaryResponse, error) {
	out := new(BandwidthMonthSummaryResponse)
	err := c.cc.Invoke(ctx, "/multinode.Bandwidth/MonthSummary", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCBandwidthServer interface {
	MonthSummary(context.Context, *BandwidthMonthSummaryRequest) (*BandwidthMonthSummaryResponse, error)
}

type DRPCBandwidthUnimplementedServer struct{}

func (s *DRPCBandwidthUnimplementedServer) MonthSummary(context.Context, *BandwidthMonthSummaryRequest) (*BandwidthMonthSummaryResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

type DRPCBandwidthDescription struct{}

func (DRPCBandwidthDescription) NumMethods() int { return 1 }

func (DRPCBandwidthDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/multinode.Bandwidth/MonthSummary", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCBandwidthServer).
					MonthSummary(
						ctx,
						in1.(*BandwidthMonthSummaryRequest),
					)
			}, DRPCBandwidthServer.MonthSummary, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterBandwidth(mux drpc.Mux, impl DRPCBandwidthServer) error {
	return mux.Register(impl, DRPCBandwidthDescription{})
}

type DRPCBandwidth_MonthSummaryStream interface {
	drpc.Stream
	SendAndClose(*BandwidthMonthSummaryResponse) error
}

type drpcBandwidth_MonthSummaryStream struct {
	drpc.Stream
}

func (x *drpcBandwidth_MonthSummaryStream) SendAndClose(m *BandwidthMonthSummaryResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCNodeClient interface {
	DRPCConn() drpc.Conn

	Version(ctx context.Context, in *VersionRequest) (*VersionResponse, error)
	LastContact(ctx context.Context, in *LastContactRequest) (*LastContactResponse, error)
	Reputation(ctx context.Context, in *ReputationRequest) (*ReputationResponse, error)
	TrustedSatellites(ctx context.Context, in *TrustedSatellitesRequest) (*TrustedSatellitesResponse, error)
	Operator(ctx context.Context, in *OperatorRequest) (*OperatorResponse, error)
}

type drpcNodeClient struct {
	cc drpc.Conn
}

func NewDRPCNodeClient(cc drpc.Conn) DRPCNodeClient {
	return &drpcNodeClient{cc}
}

func (c *drpcNodeClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcNodeClient) Version(ctx context.Context, in *VersionRequest) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/multinode.Node/Version", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcNodeClient) LastContact(ctx context.Context, in *LastContactRequest) (*LastContactResponse, error) {
	out := new(LastContactResponse)
	err := c.cc.Invoke(ctx, "/multinode.Node/LastContact", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcNodeClient) Reputation(ctx context.Context, in *ReputationRequest) (*ReputationResponse, error) {
	out := new(ReputationResponse)
	err := c.cc.Invoke(ctx, "/multinode.Node/Reputation", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcNodeClient) TrustedSatellites(ctx context.Context, in *TrustedSatellitesRequest) (*TrustedSatellitesResponse, error) {
	out := new(TrustedSatellitesResponse)
	err := c.cc.Invoke(ctx, "/multinode.Node/TrustedSatellites", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcNodeClient) Operator(ctx context.Context, in *OperatorRequest) (*OperatorResponse, error) {
	out := new(OperatorResponse)
	err := c.cc.Invoke(ctx, "/multinode.Node/Operator", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCNodeServer interface {
	Version(context.Context, *VersionRequest) (*VersionResponse, error)
	LastContact(context.Context, *LastContactRequest) (*LastContactResponse, error)
	Reputation(context.Context, *ReputationRequest) (*ReputationResponse, error)
	TrustedSatellites(context.Context, *TrustedSatellitesRequest) (*TrustedSatellitesResponse, error)
	Operator(context.Context, *OperatorRequest) (*OperatorResponse, error)
}

type DRPCNodeUnimplementedServer struct{}

func (s *DRPCNodeUnimplementedServer) Version(context.Context, *VersionRequest) (*VersionResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCNodeUnimplementedServer) LastContact(context.Context, *LastContactRequest) (*LastContactResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCNodeUnimplementedServer) Reputation(context.Context, *ReputationRequest) (*ReputationResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCNodeUnimplementedServer) TrustedSatellites(context.Context, *TrustedSatellitesRequest) (*TrustedSatellitesResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCNodeUnimplementedServer) Operator(context.Context, *OperatorRequest) (*OperatorResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

type DRPCNodeDescription struct{}

func (DRPCNodeDescription) NumMethods() int { return 5 }

func (DRPCNodeDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/multinode.Node/Version", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCNodeServer).
					Version(
						ctx,
						in1.(*VersionRequest),
					)
			}, DRPCNodeServer.Version, true
	case 1:
		return "/multinode.Node/LastContact", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCNodeServer).
					LastContact(
						ctx,
						in1.(*LastContactRequest),
					)
			}, DRPCNodeServer.LastContact, true
	case 2:
		return "/multinode.Node/Reputation", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCNodeServer).
					Reputation(
						ctx,
						in1.(*ReputationRequest),
					)
			}, DRPCNodeServer.Reputation, true
	case 3:
		return "/multinode.Node/TrustedSatellites", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCNodeServer).
					TrustedSatellites(
						ctx,
						in1.(*TrustedSatellitesRequest),
					)
			}, DRPCNodeServer.TrustedSatellites, true
	case 4:
		return "/multinode.Node/Operator", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCNodeServer).
					Operator(
						ctx,
						in1.(*OperatorRequest),
					)
			}, DRPCNodeServer.Operator, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterNode(mux drpc.Mux, impl DRPCNodeServer) error {
	return mux.Register(impl, DRPCNodeDescription{})
}

type DRPCNode_VersionStream interface {
	drpc.Stream
	SendAndClose(*VersionResponse) error
}

type drpcNode_VersionStream struct {
	drpc.Stream
}

func (x *drpcNode_VersionStream) SendAndClose(m *VersionResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCNode_LastContactStream interface {
	drpc.Stream
	SendAndClose(*LastContactResponse) error
}

type drpcNode_LastContactStream struct {
	drpc.Stream
}

func (x *drpcNode_LastContactStream) SendAndClose(m *LastContactResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCNode_ReputationStream interface {
	drpc.Stream
	SendAndClose(*ReputationResponse) error
}

type drpcNode_ReputationStream struct {
	drpc.Stream
}

func (x *drpcNode_ReputationStream) SendAndClose(m *ReputationResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCNode_TrustedSatellitesStream interface {
	drpc.Stream
	SendAndClose(*TrustedSatellitesResponse) error
}

type drpcNode_TrustedSatellitesStream struct {
	drpc.Stream
}

func (x *drpcNode_TrustedSatellitesStream) SendAndClose(m *TrustedSatellitesResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCNode_OperatorStream interface {
	drpc.Stream
	SendAndClose(*OperatorResponse) error
}

type drpcNode_OperatorStream struct {
	drpc.Stream
}

func (x *drpcNode_OperatorStream) SendAndClose(m *OperatorResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayoutsClient interface {
	DRPCConn() drpc.Conn

	Summary(ctx context.Context, in *SummaryRequest) (*SummaryResponse, error)
	SummaryPeriod(ctx context.Context, in *SummaryPeriodRequest) (*SummaryPeriodResponse, error)
	SummarySatellite(ctx context.Context, in *SummarySatelliteRequest) (*SummarySatelliteResponse, error)
	SummarySatellitePeriod(ctx context.Context, in *SummarySatellitePeriodRequest) (*SummarySatellitePeriodResponse, error)
	Earned(ctx context.Context, in *EarnedRequest) (*EarnedResponse, error)
	EarnedSatellite(ctx context.Context, in *EarnedSatelliteRequest) (*EarnedSatelliteResponse, error)
	EstimatedPayoutSatellite(ctx context.Context, in *EstimatedPayoutSatelliteRequest) (*EstimatedPayoutSatelliteResponse, error)
	EstimatedPayout(ctx context.Context, in *EstimatedPayoutRequest) (*EstimatedPayoutResponse, error)
	Undistributed(ctx context.Context, in *UndistributedRequest) (*UndistributedResponse, error)
	PaystubSatellite(ctx context.Context, in *PaystubSatelliteRequest) (*PaystubSatelliteResponse, error)
	Paystub(ctx context.Context, in *PaystubRequest) (*PaystubResponse, error)
	PaystubPeriod(ctx context.Context, in *PaystubPeriodRequest) (*PaystubPeriodResponse, error)
	PaystubSatellitePeriod(ctx context.Context, in *PaystubSatellitePeriodRequest) (*PaystubSatellitePeriodResponse, error)
	HeldAmountHistory(ctx context.Context, in *HeldAmountHistoryRequest) (*HeldAmountHistoryResponse, error)
}

type drpcPayoutsClient struct {
	cc drpc.Conn
}

func NewDRPCPayoutsClient(cc drpc.Conn) DRPCPayoutsClient {
	return &drpcPayoutsClient{cc}
}

func (c *drpcPayoutsClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcPayoutsClient) Summary(ctx context.Context, in *SummaryRequest) (*SummaryResponse, error) {
	out := new(SummaryResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/Summary", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) SummaryPeriod(ctx context.Context, in *SummaryPeriodRequest) (*SummaryPeriodResponse, error) {
	out := new(SummaryPeriodResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/SummaryPeriod", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) SummarySatellite(ctx context.Context, in *SummarySatelliteRequest) (*SummarySatelliteResponse, error) {
	out := new(SummarySatelliteResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/SummarySatellite", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) SummarySatellitePeriod(ctx context.Context, in *SummarySatellitePeriodRequest) (*SummarySatellitePeriodResponse, error) {
	out := new(SummarySatellitePeriodResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/SummarySatellitePeriod", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) Earned(ctx context.Context, in *EarnedRequest) (*EarnedResponse, error) {
	out := new(EarnedResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/Earned", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) EarnedSatellite(ctx context.Context, in *EarnedSatelliteRequest) (*EarnedSatelliteResponse, error) {
	out := new(EarnedSatelliteResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/EarnedSatellite", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) EstimatedPayoutSatellite(ctx context.Context, in *EstimatedPayoutSatelliteRequest) (*EstimatedPayoutSatelliteResponse, error) {
	out := new(EstimatedPayoutSatelliteResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/EstimatedPayoutSatellite", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) EstimatedPayout(ctx context.Context, in *EstimatedPayoutRequest) (*EstimatedPayoutResponse, error) {
	out := new(EstimatedPayoutResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/EstimatedPayout", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) Undistributed(ctx context.Context, in *UndistributedRequest) (*UndistributedResponse, error) {
	out := new(UndistributedResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/Undistributed", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) PaystubSatellite(ctx context.Context, in *PaystubSatelliteRequest) (*PaystubSatelliteResponse, error) {
	out := new(PaystubSatelliteResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/PaystubSatellite", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) Paystub(ctx context.Context, in *PaystubRequest) (*PaystubResponse, error) {
	out := new(PaystubResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/Paystub", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) PaystubPeriod(ctx context.Context, in *PaystubPeriodRequest) (*PaystubPeriodResponse, error) {
	out := new(PaystubPeriodResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/PaystubPeriod", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) PaystubSatellitePeriod(ctx context.Context, in *PaystubSatellitePeriodRequest) (*PaystubSatellitePeriodResponse, error) {
	out := new(PaystubSatellitePeriodResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/PaystubSatellitePeriod", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPayoutsClient) HeldAmountHistory(ctx context.Context, in *HeldAmountHistoryRequest) (*HeldAmountHistoryResponse, error) {
	out := new(HeldAmountHistoryResponse)
	err := c.cc.Invoke(ctx, "/multinode.Payouts/HeldAmountHistory", drpcEncoding_File_multinode_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCPayoutsServer interface {
	Summary(context.Context, *SummaryRequest) (*SummaryResponse, error)
	SummaryPeriod(context.Context, *SummaryPeriodRequest) (*SummaryPeriodResponse, error)
	SummarySatellite(context.Context, *SummarySatelliteRequest) (*SummarySatelliteResponse, error)
	SummarySatellitePeriod(context.Context, *SummarySatellitePeriodRequest) (*SummarySatellitePeriodResponse, error)
	Earned(context.Context, *EarnedRequest) (*EarnedResponse, error)
	EarnedSatellite(context.Context, *EarnedSatelliteRequest) (*EarnedSatelliteResponse, error)
	EstimatedPayoutSatellite(context.Context, *EstimatedPayoutSatelliteRequest) (*EstimatedPayoutSatelliteResponse, error)
	EstimatedPayout(context.Context, *EstimatedPayoutRequest) (*EstimatedPayoutResponse, error)
	Undistributed(context.Context, *UndistributedRequest) (*UndistributedResponse, error)
	PaystubSatellite(context.Context, *PaystubSatelliteRequest) (*PaystubSatelliteResponse, error)
	Paystub(context.Context, *PaystubRequest) (*PaystubResponse, error)
	PaystubPeriod(context.Context, *PaystubPeriodRequest) (*PaystubPeriodResponse, error)
	PaystubSatellitePeriod(context.Context, *PaystubSatellitePeriodRequest) (*PaystubSatellitePeriodResponse, error)
	HeldAmountHistory(context.Context, *HeldAmountHistoryRequest) (*HeldAmountHistoryResponse, error)
}

type DRPCPayoutsUnimplementedServer struct{}

func (s *DRPCPayoutsUnimplementedServer) Summary(context.Context, *SummaryRequest) (*SummaryResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) SummaryPeriod(context.Context, *SummaryPeriodRequest) (*SummaryPeriodResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) SummarySatellite(context.Context, *SummarySatelliteRequest) (*SummarySatelliteResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) SummarySatellitePeriod(context.Context, *SummarySatellitePeriodRequest) (*SummarySatellitePeriodResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) Earned(context.Context, *EarnedRequest) (*EarnedResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) EarnedSatellite(context.Context, *EarnedSatelliteRequest) (*EarnedSatelliteResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) EstimatedPayoutSatellite(context.Context, *EstimatedPayoutSatelliteRequest) (*EstimatedPayoutSatelliteResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) EstimatedPayout(context.Context, *EstimatedPayoutRequest) (*EstimatedPayoutResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) Undistributed(context.Context, *UndistributedRequest) (*UndistributedResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) PaystubSatellite(context.Context, *PaystubSatelliteRequest) (*PaystubSatelliteResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) Paystub(context.Context, *PaystubRequest) (*PaystubResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) PaystubPeriod(context.Context, *PaystubPeriodRequest) (*PaystubPeriodResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) PaystubSatellitePeriod(context.Context, *PaystubSatellitePeriodRequest) (*PaystubSatellitePeriodResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

func (s *DRPCPayoutsUnimplementedServer) HeldAmountHistory(context.Context, *HeldAmountHistoryRequest) (*HeldAmountHistoryResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), 12)
}

type DRPCPayoutsDescription struct{}

func (DRPCPayoutsDescription) NumMethods() int { return 14 }

func (DRPCPayoutsDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/multinode.Payouts/Summary", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					Summary(
						ctx,
						in1.(*SummaryRequest),
					)
			}, DRPCPayoutsServer.Summary, true
	case 1:
		return "/multinode.Payouts/SummaryPeriod", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					SummaryPeriod(
						ctx,
						in1.(*SummaryPeriodRequest),
					)
			}, DRPCPayoutsServer.SummaryPeriod, true
	case 2:
		return "/multinode.Payouts/SummarySatellite", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					SummarySatellite(
						ctx,
						in1.(*SummarySatelliteRequest),
					)
			}, DRPCPayoutsServer.SummarySatellite, true
	case 3:
		return "/multinode.Payouts/SummarySatellitePeriod", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					SummarySatellitePeriod(
						ctx,
						in1.(*SummarySatellitePeriodRequest),
					)
			}, DRPCPayoutsServer.SummarySatellitePeriod, true
	case 4:
		return "/multinode.Payouts/Earned", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					Earned(
						ctx,
						in1.(*EarnedRequest),
					)
			}, DRPCPayoutsServer.Earned, true
	case 5:
		return "/multinode.Payouts/EarnedSatellite", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					EarnedSatellite(
						ctx,
						in1.(*EarnedSatelliteRequest),
					)
			}, DRPCPayoutsServer.EarnedSatellite, true
	case 6:
		return "/multinode.Payouts/EstimatedPayoutSatellite", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					EstimatedPayoutSatellite(
						ctx,
						in1.(*EstimatedPayoutSatelliteRequest),
					)
			}, DRPCPayoutsServer.EstimatedPayoutSatellite, true
	case 7:
		return "/multinode.Payouts/EstimatedPayout", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					EstimatedPayout(
						ctx,
						in1.(*EstimatedPayoutRequest),
					)
			}, DRPCPayoutsServer.EstimatedPayout, true
	case 8:
		return "/multinode.Payouts/Undistributed", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					Undistributed(
						ctx,
						in1.(*UndistributedRequest),
					)
			}, DRPCPayoutsServer.Undistributed, true
	case 9:
		return "/multinode.Payouts/PaystubSatellite", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					PaystubSatellite(
						ctx,
						in1.(*PaystubSatelliteRequest),
					)
			}, DRPCPayoutsServer.PaystubSatellite, true
	case 10:
		return "/multinode.Payouts/Paystub", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					Paystub(
						ctx,
						in1.(*PaystubRequest),
					)
			}, DRPCPayoutsServer.Paystub, true
	case 11:
		return "/multinode.Payouts/PaystubPeriod", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					PaystubPeriod(
						ctx,
						in1.(*PaystubPeriodRequest),
					)
			}, DRPCPayoutsServer.PaystubPeriod, true
	case 12:
		return "/multinode.Payouts/PaystubSatellitePeriod", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					PaystubSatellitePeriod(
						ctx,
						in1.(*PaystubSatellitePeriodRequest),
					)
			}, DRPCPayoutsServer.PaystubSatellitePeriod, true
	case 13:
		return "/multinode.Payouts/HeldAmountHistory", drpcEncoding_File_multinode_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPayoutsServer).
					HeldAmountHistory(
						ctx,
						in1.(*HeldAmountHistoryRequest),
					)
			}, DRPCPayoutsServer.HeldAmountHistory, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterPayouts(mux drpc.Mux, impl DRPCPayoutsServer) error {
	return mux.Register(impl, DRPCPayoutsDescription{})
}

type DRPCPayouts_SummaryStream interface {
	drpc.Stream
	SendAndClose(*SummaryResponse) error
}

type drpcPayouts_SummaryStream struct {
	drpc.Stream
}

func (x *drpcPayouts_SummaryStream) SendAndClose(m *SummaryResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_SummaryPeriodStream interface {
	drpc.Stream
	SendAndClose(*SummaryPeriodResponse) error
}

type drpcPayouts_SummaryPeriodStream struct {
	drpc.Stream
}

func (x *drpcPayouts_SummaryPeriodStream) SendAndClose(m *SummaryPeriodResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_SummarySatelliteStream interface {
	drpc.Stream
	SendAndClose(*SummarySatelliteResponse) error
}

type drpcPayouts_SummarySatelliteStream struct {
	drpc.Stream
}

func (x *drpcPayouts_SummarySatelliteStream) SendAndClose(m *SummarySatelliteResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_SummarySatellitePeriodStream interface {
	drpc.Stream
	SendAndClose(*SummarySatellitePeriodResponse) error
}

type drpcPayouts_SummarySatellitePeriodStream struct {
	drpc.Stream
}

func (x *drpcPayouts_SummarySatellitePeriodStream) SendAndClose(m *SummarySatellitePeriodResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_EarnedStream interface {
	drpc.Stream
	SendAndClose(*EarnedResponse) error
}

type drpcPayouts_EarnedStream struct {
	drpc.Stream
}

func (x *drpcPayouts_EarnedStream) SendAndClose(m *EarnedResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_EarnedSatelliteStream interface {
	drpc.Stream
	SendAndClose(*EarnedSatelliteResponse) error
}

type drpcPayouts_EarnedSatelliteStream struct {
	drpc.Stream
}

func (x *drpcPayouts_EarnedSatelliteStream) SendAndClose(m *EarnedSatelliteResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_EstimatedPayoutSatelliteStream interface {
	drpc.Stream
	SendAndClose(*EstimatedPayoutSatelliteResponse) error
}

type drpcPayouts_EstimatedPayoutSatelliteStream struct {
	drpc.Stream
}

func (x *drpcPayouts_EstimatedPayoutSatelliteStream) SendAndClose(m *EstimatedPayoutSatelliteResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_EstimatedPayoutStream interface {
	drpc.Stream
	SendAndClose(*EstimatedPayoutResponse) error
}

type drpcPayouts_EstimatedPayoutStream struct {
	drpc.Stream
}

func (x *drpcPayouts_EstimatedPayoutStream) SendAndClose(m *EstimatedPayoutResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_UndistributedStream interface {
	drpc.Stream
	SendAndClose(*UndistributedResponse) error
}

type drpcPayouts_UndistributedStream struct {
	drpc.Stream
}

func (x *drpcPayouts_UndistributedStream) SendAndClose(m *UndistributedResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_PaystubSatelliteStream interface {
	drpc.Stream
	SendAndClose(*PaystubSatelliteResponse) error
}

type drpcPayouts_PaystubSatelliteStream struct {
	drpc.Stream
}

func (x *drpcPayouts_PaystubSatelliteStream) SendAndClose(m *PaystubSatelliteResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_PaystubStream interface {
	drpc.Stream
	SendAndClose(*PaystubResponse) error
}

type drpcPayouts_PaystubStream struct {
	drpc.Stream
}

func (x *drpcPayouts_PaystubStream) SendAndClose(m *PaystubResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_PaystubPeriodStream interface {
	drpc.Stream
	SendAndClose(*PaystubPeriodResponse) error
}

type drpcPayouts_PaystubPeriodStream struct {
	drpc.Stream
}

func (x *drpcPayouts_PaystubPeriodStream) SendAndClose(m *PaystubPeriodResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_PaystubSatellitePeriodStream interface {
	drpc.Stream
	SendAndClose(*PaystubSatellitePeriodResponse) error
}

type drpcPayouts_PaystubSatellitePeriodStream struct {
	drpc.Stream
}

func (x *drpcPayouts_PaystubSatellitePeriodStream) SendAndClose(m *PaystubSatellitePeriodResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPayouts_HeldAmountHistoryStream interface {
	drpc.Stream
	SendAndClose(*HeldAmountHistoryResponse) error
}

type drpcPayouts_HeldAmountHistoryStream struct {
	drpc.Stream
}

func (x *drpcPayouts_HeldAmountHistoryStream) SendAndClose(m *HeldAmountHistoryResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_multinode_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
