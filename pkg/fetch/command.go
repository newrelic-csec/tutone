package fetch

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/newrelic/tutone/internal/util"
)

const (
	DefaultAPIKeyEnv       = "TUTONE_API_KEY"
	DefaultSchemaCacheFile = "schema.json"
)

var refetch bool

var Command = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch GraphQL Schema",
	Long: `Fetch GraphQL Schema

Query the GraphQL server for schema and write it to a file.
`,
	Example: "tutone fetch --config configs/tutone.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		Fetch(
			viper.GetString("endpoint"),
			viper.GetBool("auth.disable"),
			viper.GetString("auth.header"),
			viper.GetString("auth.api_key_env_var"),
			viper.GetString("cache.schema_file"),
			refetch,
		)
	},
}

func Fetch(
	endpoint string,
	disableAuth bool,
	authHeader string,
	authEnvVariableName string,
	schemaFile string,
	refetch bool,
) {
	e := NewEndpoint()
	e.URL = endpoint
	e.Auth.Disable = disableAuth
	e.Auth.Header = authHeader
	e.Auth.APIKey = os.Getenv(authEnvVariableName)

	_, err := os.Stat(schemaFile)

	if os.IsNotExist(err) || refetch {
		schema, err := e.Fetch()
		if err != nil {
			log.Fatal(err)
		}

		if schemaFile != "" {
			util.LogIfError(log.ErrorLevel, schema.Save(schemaFile))
		}

		log.WithFields(log.Fields{
			"endpoint":    endpoint,
			"schema_file": schemaFile,
		}).Info("successfully fetched schema")
	}
}

func init() {
	Command.Flags().StringP("endpoint", "e", "", "GraphQL Endpoint")
	util.LogIfError(log.ErrorLevel, viper.BindPFlag("endpoint", Command.Flags().Lookup("endpoint")))

	Command.Flags().String("header", DefaultAuthHeader, "Header name set for Authentication")
	util.LogIfError(log.ErrorLevel, viper.BindPFlag("auth.header", Command.Flags().Lookup("header")))

	Command.Flags().String("api-key-env", DefaultAPIKeyEnv, "Environment variable to read API key from")
	util.LogIfError(log.ErrorLevel, viper.BindPFlag("auth.api_key_env_var", Command.Flags().Lookup("api-key-env")))

	Command.Flags().StringP("schema", "s", DefaultSchemaCacheFile, "Output file for the schema")
	util.LogIfError(log.ErrorLevel, viper.BindPFlag("cache.schema_file", Command.Flags().Lookup("schema")))

	Command.Flags().BoolVar(&refetch, "refetch", false, "Force a refetch of your GraphQL schema to ensure the generated types are up to date.")
}
