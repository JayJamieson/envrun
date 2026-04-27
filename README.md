# envrun

Run any binary with a preconfigured set of environment variables where that
binary becomes the main process (via `execve`/`syscall.Exec`).

```sh
envrun --env=HTTPS_PROXY=https://proxy:8080 /bin/bash
```

`/bin/bash` (or any binary you name) becomes the process — `envrun` is gone from
the process table the moment exec succeeds.

## Usage reference

```txt
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
```
