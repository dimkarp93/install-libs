# install-libs

Общая часть требований [CONVENTIONS.md](https://github.com/dimkarp93/install) — версия и
происхождение сборки — в виде Go-библиотеки.

Модуль: `github.com/dimkarp93/install-libs`.

## Зачем

Конвенции требуют от каждой installable-программы одинакового поведения:

- `--version` печатает голый semver без префикса `v`;
- `--origin` печатает канонический URL репозитория, из которого собран бинарь;
- `--buildinfo` печатает `key=value` в фиксированном порядке
  `origin`, `upstream`, `version`, `commit`, `channel`;
- при сборке через `go install` (без `-ldflags`) версия и origin берутся из метаданных модуля.

Пакет `buildinfo` реализует это один раз, вместо копирования одного и того же кода
в каждую тулзу.

Список проектов, которые её подключают, — в [using.md](using.md).

## Использование

Значения по-прежнему зашиваются в `package main` — так требуют конвенции и релизные
workflow (`-X main.version=...` и далее):

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
	// обычная работа программы
}
```

`Handle` разбирает `--version`, `-v`, `--origin`, `--buildinfo` и подкоманду `version`
(только первым аргументом), печатает результат в stdout и возвращает `true`, если аргумент был
обработан. Разбор останавливается на разделителе `--` и на первом аргументе, не похожем на
флаг, поэтому флаги дочерней команды (`kdbx-env run -- app --version`, `mono-vcs do git
--version`) не перехватываются.

Когда нужен не весь разбор, а отдельные значения:

```go
info := build()
info.VersionString()   // "0.4.0", либо версия модуля при go install, либо "dev"
info.OriginString()    // канонический https://host/owner/name
info.UpstreamString()  // upstream.txt, иначе origin
info.Lines()           // []string для --buildinfo
info.Fprint(w)         // те же строки в произвольный io.Writer
```

`buildinfo.NormalizeURL` приводит URL к канонической форме (срезает user info, `.git`,
конечный слэш, превращает ssh-форму в https) — это же делает и `OriginString`, так что
токен из `git remote get-url origin` не утечёт в релизный бинарь.

Константы каналов: `ChannelLocal`, `ChannelGoInstall`, `ChannelGithubRelease`,
`ChannelGiteaRelease`.

## Сборка

```sh
make check        # go vet + go test
make bump-patch   # 1.2.3 -> 1.2.4
```

Версия — в `versions.txt`; релиз создаётся workflow'ом по semver-тегу и нужен только для
`go get` (архивов у библиотеки нет).

`make pack` собирает file-based Go-прокси `dist/install-libs-<version>-proxy.tar.gz` —
способ подключить модуль в потребителях до того, как он опубликован.
