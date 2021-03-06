// Defines all the CLI command definitions and execution against internal frameworks
package commands

import (
	"fmt"
	"github.com/gondor/depcon/cliconfig"
	"github.com/gondor/depcon/commands/marathon"
	"github.com/gondor/depcon/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	FlagVerbose = "verbose"
	FlagEnv     = "env"
	ViperEnv    = "env_name"
	EnvHelp     = `Specifies the Environment name to use (eg. test | prod | etc). This can be omitted if only a single environment has been defined`
	DepConHelp  = `
DEPCON (Deploy Containers)

Provides management and deployment aid across user-defined clusters such as
  - Mesos/Marathon
  - Kubernetes
  - Amazon ECS (EC2 Container Service)
    `
)

var (
	configFile *cliconfig.ConfigFile

	// Root command for CLI command hierarchy
	rootCmd = &cobra.Command{
		Use:              "depcon",
		Short:            "Manage container clusters and deployments",
		Long:             DepConHelp,
		PersistentPreRun: configureLogging,
	}

	// Default logging levels
	logLevels = map[string]logger.LogLevel{
		"depcon":             logger.WARNING,
		"client":             logger.WARNING,
		"depcon.deploy.wait": logger.INFO,
		"depcon.marathon":    logger.WARNING,
		"depcon.marshal":     logger.WARNING,
	}
)

func init() {
	logger.InitWithDefaultLogger("depcon")
	rootCmd.PersistentFlags().StringP(FlagEnv, "e", "", EnvHelp)
	rootCmd.PersistentFlags().Bool(FlagVerbose, false, "Enables debug/verbose logging")
	viper.BindPFlag(FlagEnv, rootCmd.PersistentFlags().Lookup(FlagEnv))
}

// Main Entry point called by main - responsible for detecting if this is a first run without a config
// to force initial setup
func Execute() {
	file, found := cliconfig.HasExistingConfig()
	if found {
		configFile = file
		executeWithExistingConfig()
	} else {
		logger.Logger().Error("%s file not found.  Generating initial configuration", file.Filename())
		configFile = cliconfig.CreateNewConfigFromUserInput()
	}
}

func executeWithExistingConfig() {
	envName := findEnvNameFromArgs()

	if envName == "" {
		if _, single := configFile.DetermineIfServiceIsRooted(); single {
			envName = configFile.GetEnvironments()[0]
		} else {
			if configFile.DefaultEnv != "" {
				envName = configFile.DefaultEnv
			} else {
				rootCmd.Execute()
				logger.Logger().Error("Multiple environments are defined in config.  You must execute with -e envname.")
				printValidEnvironments()
				return
			}
		}
	}

	if _, err := configFile.GetEnvironment(envName); err != nil {
		logger.Logger().Error("'%s' environment could not be found in config (%s)\n\n", envName, configFile.Filename())
		printValidEnvironments()
		os.Exit(1)
	} else {
		viper.Set(ViperEnv, envName)
		if configFile.RootService {
			marathon.AddJailedMarathonToCmd(rootCmd, configFile)
		} else {
			marathon.AddMarathonToCmd(rootCmd, configFile)
		}
	}
	rootCmd.AddCommand(configCmd)
	rootCmd.Execute()
}

// Profiles the user with a list of current environments found within the config.json based on
// a user error or invalid flags
func printValidEnvironments() {
	envs := configFile.GetEnvironments()
	fmt.Println("Valid Environments:")
	for _, env := range envs {
		fmt.Printf("-  %s\n", env)
	}
	fmt.Println("")
}

func findEnvNameFromArgs() string {
	if len(os.Args) < 2 {
		return ""
	}
	f := os.Args[1]
	if f == "-e" && len(os.Args) > 2 {
		return os.Args[2]
	}
	if strings.HasPrefix(os.Args[1], "--env=") {
		split := strings.Split(os.Args[1], "=")
		return split[1]
	}
	return ""
}

// Configures the logging levels based on the logLevels map.  If --verbose is flagged
// then all categories defined in the map become DEBUG
func configureLogging(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool(FlagVerbose)

	for category, level := range logLevels {
		if verbose {
			logger.SetLevel(logger.DEBUG, category)
		} else {
			logger.SetLevel(level, category)
		}
	}
}
