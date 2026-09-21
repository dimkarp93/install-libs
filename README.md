# install-libs

The shared part of the [CONVENTIONS.md](https://github.com/dimkarp93/install) requirements —
build version and origin — packaged as a Go library.

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
