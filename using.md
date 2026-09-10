# Кто использует install-libs

Проекты, которые импортируют `github.com/dimkarp93/install-libs/buildinfo` — через него
реализованы флаги `--version`, `--origin` и `--buildinfo` из
[CONVENTIONS.md](https://github.com/dimkarp93/install/blob/master/CONVENTIONS.md).

| Проект | Репозиторий | Модуль | Где подключено |
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

## Способы подключения

- **`Info.Handle(os.Args[1:])`** — тулзы с ручным разбором аргументов: git-remote-gz, mvpy-env,
  mono-vcs, repos, duty, kdbx-env (внутри `cmd.Execute`).
- **Отдельные флаги `flag.Bool`** + `VersionString()` / `OriginString()` / `Print()` — тулзы на
  пакете `flag`: md-pdf, md-docx.
- **cobra** — `cmd.Version = info.VersionString()` с `SetVersionTemplate("{{.Version}}\n")`
  плюс флаги `--origin` / `--buildinfo`: rcopy.

## Примечания

- **duty** — репозиторий пока только локальный, remote не настроен; ссылка указывает на путь
  из `go.mod`, куда проект и должен уехать.
- **git-remote-gz** и **kdbx-env** живут ещё и в Gitea-зеркале, поэтому у них есть
  `upstream.txt`: в релизах оттуда `origin` указывает на зеркало, а `upstream` — на GitHub.
- **[dimkarp93/remote](https://github.com/dimkarp93/remote)** соответствует тем же конвенциям,
  но библиотеку не подключает: это POSIX-shell тулза, значения подставляются в скрипт при
  сборке.

## Обновление списка

Найти всех потребителей в рабочих копиях:

```sh
grep -rl 'install-libs/buildinfo' ~/tools ~/program --include='*.go'
```
