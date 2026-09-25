package xdgpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigHome(t *testing.T) {
	t.Setenv("HOME", "/h")
	t.Setenv("XDG_CONFIG_HOME", "")
	if got := ConfigHome(); got != "/h/.config" {
		t.Fatalf("got %q", got)
	}
	t.Setenv("XDG_CONFIG_HOME", "/x")
	if got := ConfigHome(); got != "/x" {
		t.Fatalf("got %q", got)
	}
	if got := ConfigDir("app"); got != "/x/app" {
		t.Fatalf("got %q", got)
	}
}

func TestDataHome(t *testing.T) {
	t.Setenv("HOME", "/h")
	t.Setenv("XDG_DATA_HOME", " ")
	if got := DataHome(); got != "/h/.local/share" {
		t.Fatalf("got %q", got)
	}
	t.Setenv("XDG_DATA_HOME", "/d")
	if got := DataHome(); got != "/d" {
		t.Fatalf("got %q", got)
	}
}

func TestExpandHome(t *testing.T) {
	t.Setenv("HOME", "/h")
	cases := map[string]string{
		"~":      "/h",
		"~/a/b":  "/h/a/b",
		"/abs":   "/abs",
		"rel":    "rel",
		"~other": "~other",
	}
	for in, want := range cases {
		if got := ExpandHome(in); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}

func TestResolve(t *testing.T) {
	t.Setenv("HOME", "/h")
	if got := Resolve("", "/def"); got != "/def" {
		t.Fatalf("got %q", got)
	}
	if got := Resolve("~/c", "/def"); got != "/h/c" {
		t.Fatalf("got %q", got)
	}
}

func TestWithLegacy(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "new")
	old := filepath.Join(dir, "old")
	if got := WithLegacy(primary, old); got != primary {
		t.Fatalf("nothing exists: got %q", got)
	}
	if err := os.WriteFile(old, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if got := WithLegacy(primary, "", old); got != old {
		t.Fatalf("legacy exists: got %q", got)
	}
	if err := os.WriteFile(primary, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if got := WithLegacy(primary, old); got != primary {
		t.Fatalf("both exist: got %q", got)
	}
}
