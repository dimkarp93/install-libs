# Who uses install-libs

Projects that import `github.com/dimkarp93/install-libs/buildinfo` — it is how they implement
the `--version`, `--origin` and `--buildinfo` flags from
[CONVENTIONS.md](https://github.com/dimkarp93/install/blob/master/CONVENTIONS.md).

| Project | Repository | Module | Where it is wired in |
|---|---|---|---|
| md-pdf | [dimkarp93/md-pdf](https://github.com/dimkarp93/md-pdf) | `github.com/dimkarp93/md-pdf` | `cmd/md-pdf/main.go` |
| md-docx | [dimkarp93/md-docx](https://github.com/dimkarp93/md-docx) | `github.com/dimkarp93/md-docx` | `cmd/md-docx/main.go` |
| git-remote-gz | [dimkarp93/git-remote-gz](https://github.com/dimkarp93/git-remote-gz) | `github.com/dimkarp93/git-remote-gz` | `main.go` |
| kdbx-env | [dimkarp93/kdbx-env](https://github.com/dimkarp93/kdbx-env) | `github.com/dimkarp93/kdbx-env` | `main.go`, `internal/cmd/execute.go` |
| mvpy-env | [dimkarp93/mvpy-env](https://github.com/dimkarp93/mvpy-env) | `github.com/dimkarp93/mvpy-env` | `main.go` |
| mono-vcs | [dimkarp93/mono-vcs](https://github.com/dimkarp93/mono-vcs) | `github.com/dimkarp93/mono-vcs` | `main.go` |
| repos | [dimkarp93/repos](https://github.com/dimkarp93/repos) | `github.com/dimkarp93/repos` | `main.go` |
| duty | [dimkarp93/duty](https://github.com/dimkarp93/duty) | `github.com/dimkarp93/duty` | `main.go` |
| rcopy | [rcopy/rcopy](https://github.com/rcopy/rcopy) | `github.com/rcopy/rcopy` | `cmd/rcopy/main.go`, `internal/cli/root.go` |

## Ways of wiring it in

- **`Info.Handle(os.Args[1:])`** — tools with manual argument parsing: git-remote-gz, mvpy-env,
  mono-vcs, repos, duty, kdbx-env (inside `cmd.Execute`).
- **Separate `flag.Bool` flags** + `VersionString()` / `OriginString()` / `Print()` — tools
  built on the `flag` package: md-pdf, md-docx.
- **cobra** — `cmd.Version = info.VersionString()` with `SetVersionTemplate("{{.Version}}\n")`
  plus the `--origin` / `--buildinfo` flags: rcopy.

## Notes

- **duty** — the repository is still local only, no remote is configured; the link points at
  the path from `go.mod`, which is where the project is meant to end up.
- **git-remote-gz** and **kdbx-env** also live in a Gitea mirror, which is why they have an
  `upstream.txt`: in releases made there `origin` points at the mirror and `upstream` at GitHub.
- **[dimkarp93/remote](https://github.com/dimkarp93/remote)** follows the same conventions but
  does not use the library: it is a POSIX shell tool, and the values are substituted into the
  script at build time.

## Updating the list

Find every consumer in the local working copies:

```sh
grep -rl 'install-libs/buildinfo' ~/tools ~/program --include='*.go'
```
