// Package exec handles constructing the child environment and handing
// execution over to the target binary via syscall.Exec (not os/exec).
//
// syscall.Exec replaces the current process image entirely — the envrun
// process IS the target process after the call.  There is no parent to
// wait on, no signal forwarding to worry about, and the kernel's OOM
// killer sees only the target binary.
package exec

import (
	"fmt"
	"os"
	gosys "os/exec"
	"slices"
	"strings"
	"syscall"
)

// Options controls how the target command is executed.
type Options struct {
	// EnvPairs are KEY=VALUE strings added to (or replacing) the environment.
	EnvPairs []string

	// CleanEnv, when true, starts from an empty environment (only EnvPairs).
	CleanEnv bool

	// UnsetEnv, when true, selectively unsets environment variables.
	UnsetEnv bool

	// UnsetEnvKeys is a list of environment variables to unset
	UnsetEnvKeys []string

	// Command is the path (or name, resolved via PATH) of the binary.
	Command string

	// Args are the arguments passed to the binary (not including argv[0]).
	Args []string
}

// Run builds the environment and exec's the target binary.
// It never returns on success.
func Run(opts Options) error {
	binaryPath, err := gosys.LookPath(opts.Command)
	if err != nil {
		return fmt.Errorf("resolving %q: %w", opts.Command, err)
	}

	env := buildEnv(opts)

	argv := append([]string{binaryPath}, opts.Args...)

	return syscall.Exec(binaryPath, argv, env)
}

// buildEnv constructs the environment slice.
//
//	Base: either the current process environment or empty (--clean).
//	Then: all --env=KEY=VALUE pairs are applied, overwriting any existing KEY.
func buildEnv(opts Options) []string {
	var base []string
	if !opts.CleanEnv {
		base = os.Environ()
	}

	envIndex := make(map[string]int, len(base))
	unsetIndex := make(map[string]int, len(opts.UnsetEnvKeys))
	for i, key := range base {
		envIndex[key] = i
	}

	for i, key := range opts.UnsetEnvKeys {
		unsetIndex[key] = i
	}

	if opts.UnsetEnv {
		base = slices.DeleteFunc(base, func(e string) bool {
			key := envKey(e)
			_, ok := unsetIndex[key]

			return ok
		})
	}

	for _, pair := range opts.EnvPairs {
		key := envKey(pair)
		if i, exists := envIndex[key]; exists {
			base[i] = pair // overwrite
		} else {
			envIndex[key] = len(base)
			base = append(base, pair) // append
		}
	}

	return base
}

// envKey returns the KEY portion of a KEY=VALUE string.
func envKey(pair string) string {
	if before, _, ok := strings.Cut(pair, "="); ok {
		return before
	}
	return pair
}
