package shellcomplete

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testMarker = "# t completion (t install-completions)"

func TestApplyRCLine(t *testing.T) {
	line := "source x   " + testMarker
	next, action := applyRCLine("a\n", line, testMarker)
	if action != rcAdded || next != "a\n\n"+line+"\n" {
		t.Fatalf("add: %d %q", action, next)
	}
	if _, action := applyRCLine(next, line, testMarker); action != rcSame {
		t.Fatalf("same: %d", action)
	}
	other := "source y   " + testMarker
	next2, action := applyRCLine(next, other, testMarker)
	if action != rcReplaced || next2 != "a\n\n"+other+"\n" {
		t.Fatalf("replace: %d %q", action, next2)
	}
	if next, action := applyRCLine("", line, testMarker); action != rcAdded || next != "\n"+line+"\n" {
		t.Fatalf("empty: %q", next)
	}
}

func TestStripMarked(t *testing.T) {
	text := "a\n\nsource x   " + testMarker + "\nb\n"
	next, removed := stripMarked(text, testMarker)
	if !removed || next != "a\nb\n" {
		t.Fatalf("got %v %q", removed, next)
	}
	if _, removed := stripMarked("a\n", testMarker); removed {
		t.Fatal("removed nothing")
	}
}

func TestShellFromEnv(t *testing.T) {
	cases := map[string]string{"/bin/bash": "bash", "/usr/bin/zsh": "zsh", "/usr/bin/fish": "all", "": "all"}
	for in, want := range cases {
		if got := ShellFromEnv(in); got != want {
			t.Errorf("%q: got %q", in, got)
		}
	}
}

func TestInstallUninstall(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_DATA_HOME", filepath.Join(home, "data"))
	s := Spec{Bin: "t"}
	var out bytes.Buffer

	if err := s.Install("all", true, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".bashrc")); err == nil {
		t.Fatal("dry run wrote .bashrc")
	}

	if err := s.Install("all", false, &out); err != nil {
		t.Fatal(err)
	}
	bashFile := filepath.Join(home, "data", "bash-completion", "completions", "t")
	zshFile := filepath.Join(home, "data", "zsh", "site-functions", "_t")
	for _, f := range []string{bashFile, zshFile} {
		if _, err := os.Stat(f); err != nil {
			t.Fatal(err)
		}
	}
	rc, _ := os.ReadFile(filepath.Join(home, ".bashrc"))
	if !strings.Contains(string(rc), bashFile) || !strings.Contains(string(rc), testMarker) {
		t.Fatalf("bashrc: %q", rc)
	}
	if err := s.Install("bash", false, &out); err != nil {
		t.Fatal(err)
	}
	rc2, _ := os.ReadFile(filepath.Join(home, ".bashrc"))
	if string(rc2) != string(rc) {
		t.Fatal("second install changed .bashrc")
	}

	if err := s.Uninstall("all", false, &out); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{bashFile, zshFile} {
		if _, err := os.Stat(f); err == nil {
			t.Fatalf("%s left", f)
		}
	}
	for _, name := range []string{".bashrc", ".zshrc"} {
		data, _ := os.ReadFile(filepath.Join(home, name))
		if strings.Contains(string(data), testMarker) {
			t.Fatalf("%s still marked: %q", name, data)
		}
	}
	if err := s.Install("fish", false, &out); err == nil {
		t.Fatal("fish accepted")
	}
}
