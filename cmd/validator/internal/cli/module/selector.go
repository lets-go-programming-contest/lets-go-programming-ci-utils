package module

import (
	"fmt"

	"github.com/lets-go-programming-contest/lets-go-programming-ci-utils/internal/config"
	"github.com/spf13/cobra"
)

type (
	runEFunc   func(cmd *cobra.Command, args []string) error
	funcMapper map[config.Mode]runEFunc
	funcGetModule func(cfg config.Config) config.Mode
)

func selectorRun(
	runsMapper funcMapper,
	modeGetterFunc funcGetModule,
) runEFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := config.ReadConfig(getConfigPath(cmd.Flags()))
		if err != nil {
			return fmt.Errorf("read stage config: %w", err)
		}

		runMode := modeGetterFunc(cfg)

		runFunc, ok := runsMapper[runMode]
		if !ok {
			return fmt.Errorf("unsupported mode %q for current stage", runMode)
		}

		return runFunc(cmd, args)
	}
}
