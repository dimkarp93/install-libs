package buildinfo

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionString(t *testing.T) {
	if got := (Info{Version: "0.4.0"}).VersionString(); got != "0.4.0" {
		t.Fatalf("got %q", got)
	}
	if got := (Info{Version: "v1.2.3"}).VersionString(); got != "1.2.3" {
		t.Fatalf("got %q", got)
	}
	if got := (Info{}).VersionString(); got != "dev" {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeURL(t *testing.T) {
	cases := map[string]string{
		"":                                       "local",
		"https://github.com/dimkarp93/mytool":    "https://github.com/dimkarp93/mytool",
		"https://user:token@host/owner/repo.git": "https://host/owner/repo",
		"git@github.com:dimkarp93/mytool.git":    "https://github.com/dimkarp93/mytool",
		"ssh://git@host:2222/o/r":                "https://host:2222/o/r",
		"https://host/owner/repo/":               "https://host/owner/repo",
		"/some/local/path":                       "local",
	}
	for in, want := range cases {
		if got := NormalizeURL(in); got != want {
			t.Errorf("NormalizeURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUpstreamFallsBackToOrigin(t *testing.T) {
	i := Info{Origin: "https://gitea.example.org/dima/mytool"}
	if got := i.UpstreamString(); got != i.OriginString() {
		t.Fatalf("got %q, want %q", got, i.OriginString())
	}
}

func TestLinesOrderAndOmissions(t *testing.T) {
	i := Info{
		Version:  "0.4.0",
		Origin:   "https://gitea.example.org/dima/mytool",
		Upstream: "https://github.com/dimkarp93/mytool",
		Commit:   "431b60b",
		Channel:  ChannelGiteaRelease,
	}
	want := []string{
		"origin=https://gitea.example.org/dima/mytool",
		"upstream=https://github.com/dimkarp93/mytool",
		"version=0.4.0",
		"commit=431b60b",
		"channel=gitea-release",
	}
	got := i.Lines()
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for n := range want {
		if got[n] != want[n] {
			t.Errorf("line %d: got %q, want %q", n, got[n], want[n])
		}
	}

	short := Info{Version: "0.1.0", Origin: "https://host/o/r"}.Lines()
	if len(short) != 3 {
		t.Fatalf("got %v", short)
	}
}

func TestHandle(t *testing.T) {
	i := Info{Version: "0.4.0", Origin: "https://host/o/r", Commit: "abc1234", Channel: ChannelLocal}

	var buf bytes.Buffer
	if !i.FHandle(&buf, []string{"--version"}) {
		t.Fatal("--version not handled")
	}
	if buf.String() != "0.4.0\n" {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if !i.FHandle(&buf, []string{"-v"}) || buf.String() != "0.4.0\n" {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if !i.FHandle(&buf, []string{"--origin"}) || buf.String() != "https://host/o/r\n" {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if !i.FHandle(&buf, []string{"--buildinfo"}) {
		t.Fatal("--buildinfo not handled")
	}
	if n := len(strings.Split(strings.TrimSpace(buf.String()), "\n")); n != 5 {
		t.Fatalf("got %d lines: %q", n, buf.String())
	}

	buf.Reset()
	if !i.FHandle(&buf, []string{"version"}) || buf.String() != "0.4.0\n" {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if i.FHandle(&buf, []string{"run", "--", "--version"}) {
		t.Fatal("args after -- must not be handled")
	}
	if buf.Len() != 0 {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if i.FHandle(&buf, []string{"do", "git", "--version"}) {
		t.Fatal("flags of a subcommand must not be handled")
	}
	if buf.Len() != 0 {
		t.Fatalf("got %q", buf.String())
	}

	buf.Reset()
	if i.FHandle(&buf, []string{"show", "SECRET"}) {
		t.Fatal("unrelated args must not be handled")
	}

	buf.Reset()
	if i.FHandle(&buf, nil) {
		t.Fatal("empty args must not be handled")
	}

	buf.Reset()
	if !i.FHandle(&buf, []string{"--no-color", "--buildinfo"}) {
		t.Fatal("flag after another flag must be handled")
	}
}
