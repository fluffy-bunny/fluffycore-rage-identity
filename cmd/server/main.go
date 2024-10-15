package main

import (
	"context"
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	_ "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/server/docs" // docs is generated by Swag CLI, you have to import it.
	services_AuthorizationCodeClaimsAugmentor "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/AuthorizationCodeClaimsAugmentor"
	services_EmailTemplateData "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/EmailTemplateData"
	"github.com/rs/xid"

	services_EventSink "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/EventSink"
	services_handlers_account_about "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/about"
	services_handlers_account_api_api_user_profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/api/api_user_profile"
	services_handlers_account_callback "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/callback"

	services_handlers_account_home "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/home"
	services_handlers_account_login "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/login"
	services_handlers_account_logout "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/logout"
	services_handlers_account_passkey_management "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/passkey_management"
	services_handlers_account_personal_information "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/personal_information"
	services_handlers_account_profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/profile"
	services_handlers_account_totp_management "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/totp_management"
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/oidcflowstore"
	services_user_id_generator "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/user_id_generator"
	services_oidcuser_inmemory "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/userstore/inmemory"
	services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"

	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	pkg_runtime "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/runtime"
	pkg_version "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/version"
	fluffycore_cobracore_cmd "github.com/fluffy-bunny/fluffycore/cobracore/cmd"
)

// https://github.com/swaggo/echo-swagger

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9044
// @BasePath /
// @securityDefinitions.basic BasicAuth

func main() {
	/*
		processDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}

	*/
	startup := pkg_runtime.NewStartup(
		pkg_runtime.WithConfigureServices(MyConfigServices),
	)
	fmt.Printf("Version: %s\n", pkg_version.Version())
	fluffycore_cobracore_cmd.SetVersion(pkg_version.Version())
	fluffycore_cobracore_cmd.Execute(startup)
}

func MyConfigServices(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) {
	// this extension point is called by the runtime at the end of the startup process
	// it allows you to swap out external services like the user store
	fmt.Println("MyConfigServices")
	services_user_id_generator.AddSingletonIUserIdGenerator(builder)
	services_oidcuser_inmemory.AddSingletonIFluffyCoreUserServiceServer(builder)
	services_oidcflowstore.AddSingletonAuthorizationRequestStateStoreServer(builder)
	services_AuthorizationCodeClaimsAugmentor.AddSingletonIClaimsAugmentor(builder)
	services_EventSink.AddSingletonIEventSink(builder)
	services_EmailTemplateData.AddSingletonIEmailTemplateData(builder)
	// Account Handlers
	//--------------------------------------------------------
	services_handlers_account_about.AddScopedIHandler(builder)
	services_handlers_account_callback.AddScopedIHandler(builder)
	services_handlers_account_api_api_user_profile.AddScopedIHandler(builder)
	services_handlers_account_home.AddScopedIHandler(builder)
	services_handlers_account_login.AddScopedIHandler(builder)
	services_handlers_account_logout.AddScopedIHandler(builder)
	services_handlers_account_personal_information.AddScopedIHandler(builder)
	services_handlers_account_passkey_management.AddScopedIHandler(builder)
	services_handlers_account_profile.AddScopedIHandler(builder)
	services_handlers_account_totp_management.AddScopedIHandler(builder)
	guid := xid.New().String()
	if pkg_version.Version() != "dev-build" {
		guid = pkg_version.Version()
	}
	cacheBustingHTMLConfig := &contracts_config.CacheBustingHTMLConfig{
		FilePath:   "./static/blazor/management/wwwroot/index_template.html",
		StaticPath: "./static/blazor/management/wwwroot",
		EchoPath:   "/management/*",
		RootPath:   "/management/",
		ReplaceParams: []*contracts_config.KeyValuePair{
			{
				Key:   "{title}",
				Value: config.ApplicationName,
			},
			{
				Key:   "{version}",
				Value: guid,
			},
		},
	}
	services_handlers_cache_busting_static_html.AddScopedIHandler(builder, cacheBustingHTMLConfig)

}
