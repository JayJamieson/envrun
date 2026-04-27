// Package cli handles argument parsing and command dispatch for envrun.
package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JayJamieson/envrun/internal/exec"
)

const usage = `envrun - run a command inside a controlled environment

USAGE:
  envrun [OPTIONS] <command> [args...]

OPTIONS:
  --env=KEY=VALUE    Set an environment variable (repeatable)
  --clean            Start with an empty environment (no inherited vars)
  --help             Show this help

EXAMPLES:
  # Run bash with a proxy environemnt variable set
  envrun --env=HTTPS_PROXY=https://proxy:8080 /bin/bash

  # Run curl with forced proxy and no inherited env
  envrun --clean --env=HTTPS_PROXY=https://proxy:8080 --env=HOME=/tmp /usr/bin/curl https://example.com
`

// Run is the main entry point for the CLI.
func Run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usage)
		return nil
	}

	if args[0] == "--help" || args[0] == "-h" {
		fmt.Print(usage)
		return nil
	}

	return runExec(args)
}

// runExec parses flags and executes the target command.
func runExec(args []string) error {
	var (
		envPairs []string
		cleanEnv bool
		cmdArgs  []string
	)

	i := 0
	for i < len(args) {
		arg := args[i]

		switch {
		case arg == "--":
			// Everything after -- is the command.
			cmdArgs = args[i+1:]
			i = len(args)

		case strings.HasPrefix(arg, "--env="):
			pair := strings.TrimPrefix(arg, "--env=")
			if !strings.Contains(pair, "=") {
				return fmt.Errorf("--env requires KEY=VALUE format, got: %q", pair)
			}
			envPairs = append(envPairs, pair)
			i++

		case arg == "--clean":
			cleanEnv = true
			i++

		case strings.HasPrefix(arg, "--"):
			return fmt.Errorf("unknown flag: %q", arg)

		default:
			// First non-flag argument starts the command.
			cmdArgs = args[i:]
			i = len(args)
		}
	}

	if len(cmdArgs) == 0 {
		return errors.New("no command specified; use --help for usage")
	}

	opts := exec.Options{
		EnvPairs: envPairs,
		CleanEnv: cleanEnv,
		Command:  cmdArgs[0],
		Args:     cmdArgs[1:],
	}

	return exec.Run(opts)
}
