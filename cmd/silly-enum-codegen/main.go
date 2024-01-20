package main

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/Djarvur/go-silly-enum/cmd/silly-enum-codegen/internal/flags"
)

func main() {
	err := app.Execute()
	if err != nil {
		slog.Error("run", "error", err)
		os.Exit(1)
	}
}

var app = &cobra.Command{
	Use:   "silly-enum-codegen",
	Short: "Generates some silly but useful methods for Go enum (sort of) types",
}

func init() {
	app.PersistentFlags().Bool("verbose", false, "verbose logging")

	app.AddCommand(
		flags.Generate,
	)
}
