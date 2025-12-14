package runtime

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/contracts/config"

	rage_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	zerolog "github.com/rs/zerolog"
)

// onLoadMyAppConfig will load a file and merge it over the default config
// you can still use ENV variables to replace as well.  i.e. for secrets that only come in that way.
// ---------------------------------------------------------------------
func onLoadMyAppConfig(ctx context.Context, myAppConfigFilePath string) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadIDPs").Logger()

	// If no config file path is specified, skip loading
	if myAppConfigFilePath == "" {
		log.Info().Msg("myAppConfigFilePath is empty, skipping config file load")
		return nil
	}

	log.Info().Str("myAppConfigFilePath", myAppConfigFilePath).Msg("loading my app config")
	fileContent, err := os.ReadFile(myAppConfigFilePath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read myAppConfigFilePath")
		return err
	}
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")
	overlay := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(fixedFileContent)).Decode(&overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal myAppConfigFilePath")
		return err
	}
	src := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(string(contracts_config.ConfigDefaultJSON))).Decode(&src)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
		return err
	}
	err = rage_utils.ReplaceMergeMap(overlay, src)
	if err != nil {
		log.Error().Err(err).Msg("failed to ReplaceMergeMap")
		return err
	}
	bb, err := json.Marshal(overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal overlay")
		return err
	}
	contracts_config.ConfigDefaultJSON = bb

	return nil

}
