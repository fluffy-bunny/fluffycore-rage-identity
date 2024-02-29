package main

import (
	"context"
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	_ "github.com/fluffy-bunny/fluffycore-rage-identity/cmd/server/docs" // docs is generated by Swag CLI, you have to import it.
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/oidcflowstore"
	services_oidcuser_inmemory "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/oidcuser/inmemory"
	services_user_id_generator "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/user_id_generator"

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
	fluffycore_cobracore_cmd.SetVersion(pkg_version.Version())
	fluffycore_cobracore_cmd.Execute(startup)
}

func MyConfigServices(ctx context.Context, builder di.ContainerBuilder) {
	// this extension point is called by the runtime at the end of the startup process
	// it allows you to swap out external services like the user store
	fmt.Println("MyConfigServices")
	services_user_id_generator.AddSingletonIUserIdGenerator(builder)
	services_oidcuser_inmemory.AddSingletonIFluffyCoreUserServiceServer(builder)
	services_oidcflowstore.AddSingletonIOIDCFlowStore(builder)
}
