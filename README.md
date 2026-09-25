# install-libs

The shared part of the [CONVENTIONS.md](https://github.com/dimkarp93/install) requirements —
build version and origin, config paths, flag/environment options and shell completion —
packaged as a Go library.

| Package | What it does |
|---|---|
| `buildinfo` | `--version`, `--origin`, `--buildinfo` |
| `xdgpath` | config and data paths under the XDG base directories |
| `envflag` | options set by a flag or an environment variable, with the `--envs` report |
| `shellcomplete` | bash/zsh completion: `completion`, `install-completions`, `__complete` |

Module: `github.com/dimkarp93/install-libs`.

## Why

The conventions require identical behaviour from every installable program:

- `--version` prints a bare semver without the `v` prefix;
- `--origin` prints the canonical URL of the repository the binary was built from;
- `--buildinfo` prints `key=value` in the fixed order
  `origin`, `upstream`, `version`, `commit`, `channel`;
- when built through `go install` (without `-ldflags`), the version and the origin are taken
  from the module metadata.

The `buildinfo` package implements this once, instead of copying the same code into every
tool.

The list of projects that depend on it is in [using.md](using.md).

## Usage

The values are still baked into `package main` — that is what the conventions and the release
workflows require (`-X main.version=...` and so on):

```go
package main

import "github.com/dimkarp93/install-libs/buildinfo"

var (
	version  string
	origin   string
	upstream string
	commit   string
	channel  string
)

func build() buildinfo.Info {
	return buildinfo.Info{
		Version:  version,
		Origin:   origin,
		Upstream: upstream,
		Commit:   commit,
		Channel:  channel,
	}
}

func main() {
	if build().Handle(os.Args[1:]) {
		return
	}
	// the program's normal work
}
```

`Handle` parses `--version`, `-v`, `--origin`, `--buildinfo` and the `version` subcommand
(only as the first argument), prints the result to stdout and returns `true` if the argument
was handled. Parsing stops at the `--` separator and at the first argument that does not look
like a flag, so the flags of a child command (`kdbx-env run -- app --version`, `mono-vcs do git
--version`) are not intercepted.

When you need individual values rather than the whole parsing:

```go
info := build()
info.VersionString()   // "0.4.0", or the module version under go install, or "dev"
info.OriginString()    // the canonical https://host/owner/name
info.UpstreamString()  // upstream.txt, otherwise origin
info.Lines()           // []string for --buildinfo
info.Fprint(w)         // the same lines into an arbitrary io.Writer
```

`buildinfo.NormalizeURL` brings a URL to its canonical form (strips the user info, `.git` and
the trailing slash, turns the ssh form into https) — `OriginString` does the same, so a token
from `git remote get-url origin` does not leak into a release binary.

Channel constants: `ChannelLocal`, `ChannelGoInstall`, `ChannelGithubRelease`,
`ChannelGiteaRelease`.

## xdgpath

Config paths follow the XDG base directories. The base from `mono-vcs` and `ExpandHome` from
`kdbx-cli` were merged into one package.

```go
xdgpath.ConfigHome()              // $XDG_CONFIG_HOME, otherwise ~/.config
xdgpath.DataHome()                // $XDG_DATA_HOME, otherwise ~/.local/share
xdgpath.ConfigDir("kdbx-cli")     // ConfigHome()/kdbx-cli
xdgpath.ExpandHome("~/x")         // ~ and ~/... in user input
xdgpath.Resolve(flagValue, def)   // the --config override when set, otherwise def
xdgpath.WithLegacy(path, old)     // path if it exists, otherwise the first existing old path
```

`WithLegacy` keeps an existing config readable when a tool moves from a hardcoded
`~/.config/<app>` to `ConfigDir`: if `XDG_CONFIG_HOME` points elsewhere and the new file does
not exist yet, the old one is used.

## envflag

Options that can be set by a flag or by an environment variable, with the priority
`flag > env > default` and the source of every value. Generalised from `net-install`.

```go
var retries int
opts := envflag.New("NET_",
	envflag.Option{Name: "retries", Default: "5", Parse: func(v string) error {
		n, err := strconv.Atoi(v)
		if err == nil {
			retries = n
		}
		return err
	}},
	envflag.Option{Name: "retry-all-errors", Env: "NET_RETRY_ALL", Default: "0"},
)
if err := opts.LoadEnv(os.LookupEnv); err != nil { ... }
if opts.HandleEnvs(os.Stdout, os.Args[1:]) {
	return
}
opts.Set("retries", "7", envflag.SourceFlag)   // from the tool's own flag parser
```

- The variable name defaults to the prefix plus the upper-cased option name with `-` turned
  into `_`; `Option.Env` overrides it.
- `Parse` converts and validates the string and stores the typed value in the tool; a failed
  `Parse` leaves the previous value and source in place.
- `--envs` (only as the first argument) prints `Report()`: one `NAME=value (default|env|flag)`
  line per option, in declaration order.
- `envflag.ParseBool` accepts `1/true/yes/on` and `0/false/no/off`.

## shellcomplete

Shell completion for bash and zsh without third-party packages. Taken from
`git-repos` (`internal/app/autocomplete.go`); what differs is that the command tree is
described declaratively, so tools with their own argument parsers can use it too.

```go
spec := shellcomplete.Spec{
	Bin:   "mono-vcs",
	Flags: []shellcomplete.Flag{{Name: "--version", Bool: true}},
	Commands: []shellcomplete.Command{
		{Name: "pull", Flags: shellcomplete.FromFlagSet(pullFlags())},
		{Name: "do", Args: shellcomplete.Anything},
		{Name: "alias", Args: shellcomplete.Positional("add", "rm", "list"), Intermixed: true},
	},
}
if code, ok := spec.Handle(os.Stdout, os.Stderr, os.Args[1:]); ok {
	os.Exit(code)
}
```

`Handle` serves the first argument:

- `completion bash|zsh` prints the script, for `source <(tool completion bash)`;
- `install-completions [--dry-run] [bash|zsh|all]` writes the script to
  `$XDG_DATA_HOME/bash-completion/completions/<bin>` and `$XDG_DATA_HOME/zsh/site-functions/_<bin>`
  and adds one marked line to `~/.bashrc` / `~/.zshrc`; running it again changes nothing;
  without an argument the shell comes from `$SHELL`;
- `uninstall-completions` removes the files and the marked line;
- `__complete <words...>` is what the script calls on every TAB: it prints the candidates one
  per line, or one of the directives `:none` (nothing), `:dirs` (directories), `:default`
  (file names).

Describing the tree:

- `Flag.Name` is the spelling the tool accepts (`--config`, `-jobs`, `-y`); a single-dash long
  flag is offered as `--name` when the user has typed `--`, as the `flag` package accepts both.
- `Flag.Bool` flags take no value. For the rest, the value comes from `Values`, or `Dirs`, or
  `Files`, otherwise nothing is suggested.
- `FromFlagSet` builds the list from a `*flag.FlagSet`, `With` adds values to some of its
  flags.
- `Command.Args(pos)` suggests the positional argument number `pos`: `Positional(...)` for the
  first one, `Anything` for file names.
- Flags stop at the first positional argument, as in the `flag` package, unless `Intermixed`
  is set; after `--` only positional arguments are completed.
- A tool without subcommands leaves `Commands` empty and sets `Spec.Args`.

## Build

```sh
make check        # go vet + go test
make bump-patch   # 1.2.3 -> 1.2.4
```

The version lives in `versions.txt`; the release is created by the workflow from a semver tag
and is only needed for `go get` (the library ships no archives).

`make pack` builds a file-based Go proxy `dist/install-libs-<version>-proxy.tar.gz` — a way to
wire the module into consumers before it is published.

## License

[MIT](LICENSE)
