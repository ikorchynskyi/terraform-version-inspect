package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/ikorchynskyi/terraform-version-inspect/internal"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terraform-version-inspect",
	Short: "Inspects terraform project to determine the required terraform version",
	Long: `A CLI application to determine the required terraform version.

Does the shallow terraform project parsing to provide the required version.
The list of available versions is taken from https://releases.hashicorp.com/.
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
			PartsOrder: []string{
				zerolog.LevelFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			},
		})
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// SilenceErrors is an option to quiet errors down stream.
		cmd.SilenceErrors = true

		module, err := internal.GetModule(dir)
		if err != nil {
			return err
		}
		constraints, err := internal.GetConstraints(module)
		if err != nil {
			return err
		}
		versions, err := internal.GetVersions()
		if err != nil {
			return err
		}
		latestRequired, err := internal.GetLatestRequired(constraints, versions, registry)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, latestRequired)
		return nil
	},
}

var debug bool
var dir string
var registry string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// DisableDefaultCmd prevents Cobra from creating a default 'completion' command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// SilenceUsage is an option to silence usage when an error occurs.
	rootCmd.SilenceUsage = true

	// Persistent flags which will be global for the application.
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "turn on debug logging")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", ".", "path that contains terraform configuration files")
	rootCmd.PersistentFlags().StringVar(&registry, "registry", "", "ensure the terraform image being available in the specified registry")
}
