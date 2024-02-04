package flags

//nolint:gci
import (
	"github.com/spf13/cobra"

	"github.com/Djarvur/go-silly-enum/internal/generator"
)

const (
	generateErrorCode = 2

	buildTagsFlag    = "buildTags"
	envVarsFlag      = "envVars"
	excludeTestsFlag = "excludeTests"
	enumNameFlag     = "enumName"
)

type generateCustom struct {
	enumName *Regexp
}

// Generate is a generate CLI command.
var Generate = func() *cobra.Command { //nolint:gochecknoglobals
	enumName := must(NewRegexp("^.+Enum$"))

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:                   "generate package [package...]",
		DisableFlagsInUseLine: true,
		Short:                 "read sources and generate the code",
		Run:                   func(cmd *cobra.Command, args []string) { generateRun(cmd, args, enumName) },
		Args:                  cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringArray(buildTagsFlag, nil, "build tags to be used for sources parsing")
	cmd.Flags().StringArray(envVarsFlag, nil, "environment variables to be used for sources parsing")
	cmd.Flags().Bool(excludeTestsFlag, false, "do not process test files")
	cmd.Flags().Var(enumName, enumNameFlag, "regexp to find the Enum type(s)")

	return cmd
}()

func generateRun(
	cmd *cobra.Command,
	args []string,
	enumName *Regexp,
) {
	log := slogNew(must(cmd.Flags().GetBool("verbose")))

	err := generator.Generate(
		args,
		must(cmd.Flags().GetStringArray(buildTagsFlag)),
		must(cmd.Flags().GetStringArray(envVarsFlag)),
		!must(cmd.Flags().GetBool(excludeTestsFlag)),
		enumName,
		log,
	)
	if err != nil {
		exitWithLog(generateErrorCode, log, "generate", "error", err)
	}
}
