package main

import (
	"fmt"
	"os"

	"github.com/oraichain/orai/oraivisor"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

// Run is the main loop, but returns an error
func Run(args []string) error {
	cfg, err := oraivisor.GetConfigFromEnv()
	if err != nil {
		return err
	}

	doUpgrade, err := oraivisor.LaunchProcess(cfg, args, os.Stdout, os.Stderr)
	// if RestartAfterUpgrade, we launch after a successful upgrade (only condition LaunchProcess returns nil)
	for cfg.RestartAfterUpgrade && err == nil && doUpgrade {
		doUpgrade, err = oraivisor.LaunchProcess(cfg, args, os.Stdout, os.Stderr)
	}
	return err
}
