package auth

import (
	proto_oidcuser "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	services_common_claimsprincipal "github.com/fluffy-bunny/fluffycore/services/common/claimsprincipal"
)

var writeEndpoints = []string{
	proto_oidcuser.RageUserService_CreateRageUser_FullMethodName,
	proto_oidcuser.RageUserService_UpdateRageUser_FullMethodName,
	proto_oidcuser.RageUserService_LinkRageUser_FullMethodName,
	proto_oidcuser.RageUserService_UnlinkRageUser_FullMethodName,
}
var readEndpoints = []string{
	proto_oidcuser.RageUserService_GetRageUser_FullMethodName,
}

var noAuthEndpoints = []string{
	"/grpc.health.v1.Health/Check",
}

func BuildGrpcEntrypointPermissionsClaimsMap() map[string]contracts_common.IEntryPointConfig {
	entryPointClaimsBuilder := services_common_claimsprincipal.NewEntryPointClaimsBuilder()
	for _, endpoint := range noAuthEndpoints {
		entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(endpoint)
	}

	for _, endpoint := range writeEndpoints {
		entrypointConfig := &services_common_claimsprincipal.EntryPointConfig{
			FullMethodName: endpoint,
			ClaimsAST: &services_common_claimsprincipal.ClaimsAST{
				Or: []contracts_common.IClaimsValidator{
					&services_common_claimsprincipal.ClaimsAST{
						ClaimFacts: []contracts_common.IClaimFact{
							services_common_claimsprincipal.NewClaimFact(contracts_common.Claim{
								Type:  "permission",
								Value: "User.Write",
							}),
							services_common_claimsprincipal.NewClaimFact(contracts_common.Claim{
								Type:  "permission",
								Value: "User.Write.All",
							}),
						},
					},
				},
			},
		}
		entryPointClaimsBuilder.EntrypointClaimsMap[endpoint] = entrypointConfig
	}
	for _, endpoint := range readEndpoints {
		entrypointConfig := &services_common_claimsprincipal.EntryPointConfig{
			FullMethodName: endpoint,
			ClaimsAST: &services_common_claimsprincipal.ClaimsAST{
				Or: []contracts_common.IClaimsValidator{
					&services_common_claimsprincipal.ClaimsAST{
						ClaimFacts: []contracts_common.IClaimFact{
							services_common_claimsprincipal.NewClaimFact(contracts_common.Claim{
								Type:  "permission",
								Value: "User.ReadWrite.All",
							}),
							services_common_claimsprincipal.NewClaimFact(contracts_common.Claim{
								Type:  "permission",
								Value: "User.Read.All",
							}),
						},
					},
				},
			},
		}
		entryPointClaimsBuilder.EntrypointClaimsMap[endpoint] = entrypointConfig
	}
	entryPointClaimsBuilder.DumpExpressions()
	return entryPointClaimsBuilder.GetEntryPointClaimsMap()
}
