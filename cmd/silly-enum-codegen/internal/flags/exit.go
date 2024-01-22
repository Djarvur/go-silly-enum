package flags

import (
	"os"

	"golang.org/x/exp/slog"
)

func exitWithLog(code int, log *slog.Logger, msg string, logArgs ...any) {
	log.Error(msg, logArgs...)
	os.Exit(code)
}
