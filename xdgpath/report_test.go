package xdgpath

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPathsReport(t *testing.T) {
	dir := t.TempDir()
	present := filepath.Join(dir, "config.json")
	if err := os.WriteFile(present, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "state.json")

	p := NewPaths(
		Entry{Name: "config", Path: present},
		Entry{Name: "state", Path: missing},
	)
	got := p.Report()
	want := []string{
		"config: " + present + " (found)",
		"state: " + missing + " (not found)",
	}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPathsHandlePath(t *testing.T) {
	p := NewPaths(Entry{Name: "config", Path: "/x"})
	var buf bytes.Buffer
	if !p.HandlePath(&buf, []string{"--path"}) {
		t.Fatal("--path not handled")
	}
	if buf.String() != "config: /x (not found)\n" {
		t.Fatalf("got %q", buf.String())
	}
	buf.Reset()
	if p.HandlePath(&buf, []string{"other", "--path"}) || p.HandlePath(&buf, nil) || buf.Len() != 0 {
		t.Fatal("--path handled outside the first argument")
	}
}
