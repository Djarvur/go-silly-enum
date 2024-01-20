package flags

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/Djarvur/go-silly-enum/internal/generator"
)

const (
	GenerateErrorCode = 2

	buildTagsFlag    = "buildTags"
	excludeTestsFlag = "excludeTests"
)

var Generate = &cobra.Command{
	Use:   "generate",
	Short: "read sources and generate the code",
	Run:   generateRun,
}

func init() {
	Generate.Flags().StringArray(buildTagsFlag, nil, "build tags to be used for sources parsing")
	Generate.Flags().Bool(excludeTestsFlag, false, "do not process test files")
}

func generateRun(cmd *cobra.Command, args []string) {
	log := slogNew(must(cmd.Flags().GetBool("verbose")))

	unprocessed, err := generator.Generate(
		args,
		must(cmd.Flags().GetStringArray(buildTagsFlag)),
		!must(cmd.Flags().GetBool(excludeTestsFlag)),
		log,
	)
	if err != nil {
		exitWithLog(GenerateErrorCode, log, "generate", "error", err)
	}

	slog.Debug("generate", "unprocessed", unprocessed)
}
