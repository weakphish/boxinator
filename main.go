package main

import (
	"context"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is the CLI version printed by fang when requested.
//
// It exists as a package variable so release builds can replace it with
// ldflags without adding a version file yet.
var Version = "dev"

// main builds the root command and lets fang handle CLI execution.
//
// Fang owns the terminal-facing parts of command execution so the command
// definitions stay boring and easy to replace once Boxinator has real work.
func main() {
	if err := fang.Execute(context.Background(), newRootCommand()); err != nil {
		os.Exit(1)
	}
}

// newRootCommand returns the top-level Boxinator command.
//
// The command only wires configuration and logging for now; behavior belongs in
// subcommands once the first real workflow is known.
func newRootCommand() *cobra.Command {
	config := viper.New()
	config.SetEnvPrefix("BOXINATOR")
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	config.AutomaticEnv()

	var verbose bool

	cmd := &cobra.Command{
		Use:     "boxinator",
		Short:   "Run coding agents inside Apple container sandboxes",
		Version: Version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := config.BindPFlag("verbose", cmd.Root().PersistentFlags().Lookup("verbose")); err != nil {
				return err
			}
			configureLogger(config.GetBool("verbose"))
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")
	cmd.AddCommand(newRunCommand())

	return cmd
}

// newRunCommand returns the stub for the future sandbox launcher.
//
// Keeping it as a stub makes the intended CLI shape visible without pretending
// container policy exists before it has been designed.
func newRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run PROJECT",
		Short: "Start a sandboxed agent session",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			log.Info("run is not implemented yet", "project", args[0])
		},
	}
}

// configureLogger applies process-wide logging defaults.
//
// It centralizes the debug toggle so future commands do not each invent their
// own logging setup.
func configureLogger(verbose bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
		return
	}
	log.SetLevel(log.InfoLevel)
}
