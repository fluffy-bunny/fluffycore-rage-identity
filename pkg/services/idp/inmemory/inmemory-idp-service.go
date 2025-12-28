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
		proto_oidc_idp.UnimplementedSingletonIDPServiceServer

		idps *proto_oidc_models.IDPs
	}
)

var stemService = (*service)(nil)
var _ proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer = stemService

func (s *service) Ctor(idps *proto_oidc_models.IDPs) (proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer, error) {

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

func AddSingletonIFluffyCoreSingletonIDPServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer](cb, stemService.Ctor)
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
		if request.Filter.ClaimedDomains != nil {
			request.Filter.ClaimedDomains.Eq = strings.ToLower(request.Filter.ClaimedDomains.Eq)
			if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomains.In) {
				for i, v := range request.Filter.ClaimedDomains.In {
					request.Filter.ClaimedDomains.In[i] = strings.ToLower(v)
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
			if request.Filter.Hidden != nil {
				if request.Filter.Hidden.Eq != c.Hidden {
					return false
				}
			}
			if request.Filter.ClaimedDomains != nil {
				claimedDomainsMap := make(map[string]bool)
				for _, v := range c.ClaimedDomains {
					claimedDomainsMap[v] = true
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomains.Eq) {
					_, ok := claimedDomainsMap[request.Filter.ClaimedDomains.Eq]
					if !ok {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.ClaimedDomains.In) {
					gotHit := false
					for _, v := range request.Filter.ClaimedDomains.In {
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
		IDPs: idps,
	}, nil
}
