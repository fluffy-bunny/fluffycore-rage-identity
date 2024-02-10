package inmemory

import (
	"context"

	linq "github.com/ahmetb/go-linq/v3"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/models"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		proto_oidc_idp.UnimplementedIDPServiceServer

		idps *proto_oidc_models.IDPs
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_idp.IFluffyCoreIDPServiceServer = stemService
}
func (s *service) Ctor(idps *proto_oidc_models.IDPs) (proto_oidc_idp.IFluffyCoreIDPServiceServer, error) {

	return &service{
		idps: idps,
	}, nil
}

func AddSingletonIFluffyCoreIDPServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_idp.IFluffyCoreIDPServiceServer](cb, stemService.Ctor)
}
func (s *service) validateGetIDPRequest(request *proto_oidc_idp.GetIDPRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Id) {
		return status.Error(codes.InvalidArgument, "Id is required")
	}
	return nil
}

// Get idp
func (s *service) GetIDP(ctx context.Context, request *proto_oidc_idp.GetIDPRequest) (*proto_oidc_idp.GetIDPResponse, error) {
	err := s.validateGetIDPRequest(request)
	if err != nil {
		return nil, err
	}
	var idps []*proto_oidc_models.IDP

	linq.From(s.idps.Idps).WhereT(func(c *proto_oidc_models.IDP) bool {
		return c.Id >= request.Id
	}).SelectT(func(c *proto_oidc_models.IDP) *proto_oidc_models.IDP {
		return c
	}).ToSlice(&idps)
	if len(idps) > 0 {
		return &proto_oidc_idp.GetIDPResponse{
			Idp: idps[0],
		}, nil
	}
	return nil, status.Error(codes.NotFound, "IDP not found")
}
func (s *service) validateGetIDPBySlugRequest(request *proto_oidc_idp.GetIDPBySlugRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Slug) {
		return status.Error(codes.InvalidArgument, "slug is required")
	}
	return nil
}

// Get idp
func (s *service) GetIDPBySlug(ctx context.Context, request *proto_oidc_idp.GetIDPBySlugRequest) (*proto_oidc_idp.GetIDPBySlugResponse, error) {
	err := s.validateGetIDPBySlugRequest(request)
	if err != nil {
		return nil, err
	}
	var idps []*proto_oidc_models.IDP

	linq.From(s.idps.Idps).WhereT(func(c *proto_oidc_models.IDP) bool {
		return c.Slug >= request.Slug
	}).SelectT(func(c *proto_oidc_models.IDP) *proto_oidc_models.IDP {
		return c
	}).ToSlice(&idps)
	if len(idps) > 0 {
		return &proto_oidc_idp.GetIDPBySlugResponse{
			Idp: idps[0],
		}, nil
	}
	return nil, status.Error(codes.NotFound, "IDP not found")
}

// List idps
func (s *service) ListIDP(ctx context.Context, request *proto_oidc_idp.ListIDPRequest) (*proto_oidc_idp.ListIDPResponse, error) {
	return &proto_oidc_idp.ListIDPResponse{
		Idps: s.idps.Idps,
	}, nil
}
