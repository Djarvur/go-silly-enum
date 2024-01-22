// Package main is an app itself
package main

//nolint:gci
import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/Djarvur/go-silly-enum/cmd/silly-enum-codegen/internal/flags"
)

func main() {
	app := buildApp()

	err := app.Execute()
	if err != nil {
		slog.Error("run", "error", err)
		os.Exit(1)
	}
}

func buildApp() *cobra.Command {
	app := &cobra.Command{ //nolint:exhaustruct
		Use:   "silly-enum-codegen",
		Short: "Generates some silly but useful methods for Go enum (sort of) types",
	}

	app.PersistentFlags().Bool("verbose", false, "verbose logging")

	app.AddCommand(
		flags.Generate,
	)

	return app
}
