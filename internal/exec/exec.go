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
	"strings"
	"syscall"
)

// Options controls how the target command is executed.
type Options struct {
	// EnvPairs are KEY=VALUE strings added to (or replacing) the environment.
	EnvPairs []string

	// CleanEnv, when true, starts from an empty environment (only EnvPairs).
	CleanEnv bool

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

	index := make(map[string]int, len(base))
	for i, pair := range base {
		key := envKey(pair)
		index[key] = i
	}

	for _, pair := range opts.EnvPairs {
		key := envKey(pair)
		if i, exists := index[key]; exists {
			base[i] = pair // overwrite
		} else {
			index[key] = len(base)
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
