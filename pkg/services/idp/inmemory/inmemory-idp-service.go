package inmemory

import (
	"context"
	"strings"

	linq "github.com/ahmetb/go-linq/v3"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
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

	for _, idp := range idps.Idps {
		if idp.ClaimedDomains == nil {
			idp.ClaimedDomains = make([]string, 0)
		}
		tolowerDomains := make([]string, 0)
		for _, v := range idp.ClaimedDomains {
			tolowerDomains = append(tolowerDomains, strings.ToLower(v))
		}
		idp.ClaimedDomains = tolowerDomains

	}
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
		return c.Id == request.Id
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
		return c.Slug == request.Slug
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

func (s *service) validateListIDPRequest(request *proto_oidc_idp.ListIDPRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if request.Filter != nil {
		if request.Filter.Slug != nil {
			request.Filter.Slug.Eq = strings.ToLower(request.Filter.Slug.Eq)
			if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Slug.In) {
				for i, v := range request.Filter.Slug.In {
					request.Filter.Slug.In[i] = strings.ToLower(v)
				}
			}
		}
		if request.Filter.ClaimedDomain != nil {
			request.Filter.ClaimedDomain.Eq = strings.ToLower(request.Filter.ClaimedDomain.Eq)
			if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomain.In) {
				for i, v := range request.Filter.ClaimedDomain.In {
					request.Filter.ClaimedDomain.In[i] = strings.ToLower(v)
				}
			}
		}
	}
	return nil

}

// List idps
func (s *service) ListIDP(ctx context.Context, request *proto_oidc_idp.ListIDPRequest) (*proto_oidc_idp.ListIDPResponse, error) {
	err := s.validateListIDPRequest(request)
	if err != nil {
		return nil, err
	}
	var idps []*proto_oidc_models.IDP

	linq.From(s.idps.Idps).WhereT(func(c *proto_oidc_models.IDP) bool {
		if request.Filter != nil {
			if request.Filter.Enabled != nil {
				if request.Filter.Enabled.Eq != c.Enabled {
					return false
				}
			}
			if request.Filter.Slug != nil {
				if request.Filter.Slug.Eq != c.Slug {
					return false
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Slug.In) {
					gotHit := false
					for _, v := range request.Filter.Slug.In {
						if v == c.Slug {
							gotHit = true
							break
						}
					}
					if !gotHit {
						return false
					}
				}
			}
			if request.Filter.EmailVerificationRequired != nil {
				if request.Filter.EmailVerificationRequired.Eq != c.EmailVerificationRequired {
					return false
				}
			}
			if request.Filter.AutoCreate != nil {
				if request.Filter.AutoCreate.Eq != c.AutoCreate {
					return false
				}
			}
			if request.Filter.Hidden != nil {
				if request.Filter.Hidden.Eq != c.Hidden {
					return false
				}
			}
			if request.Filter.Metadata != nil {
				metadataValue, ok := c.Metadata[request.Filter.Metadata.Key]
				if !ok {
					return false
				}
				if metadataValue != request.Filter.Metadata.Value.Eq {
					return false
				}
			}
			if request.Filter.ClaimedDomain != nil {
				claimedDomainsMap := make(map[string]bool)
				for _, v := range c.ClaimedDomains {
					claimedDomainsMap[v] = true
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomain.Eq) {
					_, ok := claimedDomainsMap[request.Filter.ClaimedDomain.Eq]
					if !ok {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomain.In) {
					gotHit := false
					for _, v := range request.Filter.ClaimedDomain.In {
						_, ok := claimedDomainsMap[v]
						if ok {
							gotHit = true
							break
						}
					}
					if !gotHit {
						return false
					}
				}
			}
			return true
		} else {
			return true
		}

	}).SelectT(func(c *proto_oidc_models.IDP) *proto_oidc_models.IDP {
		return c
	}).ToSlice(&idps)

	return &proto_oidc_idp.ListIDPResponse{
		Idps: idps,
	}, nil
}
