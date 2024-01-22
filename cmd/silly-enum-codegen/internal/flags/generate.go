package flags

//nolint:gci
import (
	"github.com/spf13/cobra"

	"github.com/Djarvur/go-silly-enum/internal/generator"
)

const (
	generateErrorCode = 2

	buildTagsFlag    = "buildTags"
	excludeTestsFlag = "excludeTests"
)

// Generate is a generate CLI command.
var Generate = func() *cobra.Command { //nolint:gochecknoglobals
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "generate",
		Short: "read sources and generate the code",
		Run:   generateRun,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().StringArray(buildTagsFlag, nil, "build tags to be used for sources parsing")
	cmd.Flags().Bool(excludeTestsFlag, false, "do not process test files")

	return cmd
}()

func generateRun(cmd *cobra.Command, args []string) {
	log := slogNew(must(cmd.Flags().GetBool("verbose")))

	err := generator.Generate(
		args[0],
		must(cmd.Flags().GetStringArray(buildTagsFlag)),
		!must(cmd.Flags().GetBool(excludeTestsFlag)),
		log,
	)
	if err != nil {
		exitWithLog(generateErrorCode, log, "generate", "error", err)
	}
}
