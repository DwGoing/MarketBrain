package api

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Treasury struct {
	UnimplementedTreasuryServer
}

func (Self *Treasury) createChargeOrder() (string, error) {
	return "", nil
}

func (Self *Treasury) CreateChargeOrder(ctx context.Context, request *emptypb.Empty) (*CreateChargeOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChargeOrder not implemented")
}
