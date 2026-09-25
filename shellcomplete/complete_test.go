package shellcomplete

import (
	"flag"
	"slices"
	"testing"
)

func testSpec() Spec {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.Int("jobs", 1, "")
	fs.Bool("dry-run", false, "")
	fs.String("C", "", "")
	fs.String("protocol", "", "")
	return Spec{
		Bin:   "tool",
		Flags: []Flag{{Name: "--version", Bool: true}, {Name: "--config", Files: true}},
		Commands: []Command{
			{Name: "sync", Flags: With(FromFlagSet(fs), Flag{Name: "-C", Dirs: true}, Flag{Name: "-protocol", Values: Static("ssh", "https")})},
			{Name: "do", Flags: []Flag{{Name: "-all", Bool: true}}, Args: Anything},
			{Name: "alias", Args: Positional("add", "rm", "list"), Intermixed: true, Flags: []Flag{{Name: "-x", Bool: true}}},
			{Name: "secret", Hidden: true},
			{Name: "help"},
		},
	}
}

func TestComplete(t *testing.T) {
	s := testSpec()
	cases := []struct {
		words []string
		want  []string
	}{
		{nil, []string{"sync", "do", "alias", "help", "completion", "install-completions", "uninstall-completions"}},
		{[]string{"s"}, []string{"sync"}},
		{[]string{"i"}, []string{"install-completions"}},
		{[]string{"--v"}, []string{"--version"}},
		{[]string{"--config", ""}, []string{DirDefault}},
		{[]string{"--config", "f", "s"}, []string{"sync"}},
		{[]string{"nope", ""}, []string{DirNone}},
		{[]string{"secret", ""}, []string{DirNone}},
		{[]string{"sync", "-j"}, []string{"-jobs"}},
		{[]string{"sync", "--j"}, []string{"--jobs"}},
		{[]string{"sync", "-d"}, []string{"-dry-run"}},
		{[]string{"sync", "-z"}, []string{DirNone}},
		{[]string{"sync", "-C", ""}, []string{DirDirs}},
		{[]string{"sync", "--protocol", "h"}, []string{"https"}},
		{[]string{"sync", "-protocol", ""}, []string{"ssh", "https"}},
		{[]string{"sync", "--protocol=s"}, []string{"--protocol=ssh"}},
		{[]string{"sync", "--protocol", "=", "h"}, []string{"https"}},
		{[]string{"sync", "--protocol", "="}, []string{"ssh", "https"}},
		{[]string{"sync", "-jobs", ""}, []string{DirNone}},
		{[]string{"sync", "-dry-run", ""}, []string{DirNone}},
		{[]string{"do", ""}, []string{DirDefault}},
		{[]string{"do", "-all", "git", "-"}, []string{DirDefault}},
		{[]string{"do", "-a"}, []string{"-all"}},
		{[]string{"alias", "a"}, []string{"add"}},
		{[]string{"alias", "add", "-x", ""}, []string{DirNone}},
		{[]string{"alias", "add", "-"}, []string{"-x"}},
		{[]string{"completion", ""}, []string{"bash", "zsh"}},
		{[]string{"completion", "bash", ""}, []string{DirNone}},
		{[]string{"install-completions", "--dry-run", "a"}, []string{"all"}},
		{[]string{"install-completions", "zsh", "--d"}, []string{"--dry-run"}},
	}
	for _, c := range cases {
		if got := s.Complete(c.words); !slices.Equal(got, c.want) {
			t.Errorf("%q: got %q, want %q", c.words, got, c.want)
		}
	}
}

func TestCompleteNoCommands(t *testing.T) {
	s := Spec{
		Bin:   "md",
		Flags: []Flag{{Name: "-in", Files: true}, {Name: "-out", Files: true}, {Name: "-toc", Bool: true}},
		Args:  Anything,
	}
	cases := []struct {
		words []string
		want  []string
	}{
		{[]string{""}, []string{DirDefault}},
		{[]string{"-"}, []string{"-in", "-out", "-toc"}},
		{[]string{"--o"}, []string{"--out"}},
		{[]string{"-in", ""}, []string{DirDefault}},
		{[]string{"-toc", ""}, []string{DirDefault}},
		{[]string{"completion", "z"}, []string{"zsh"}},
		{[]string{"file.md", "completion", ""}, []string{DirDefault}},
	}
	for _, c := range cases {
		if got := s.Complete(c.words); !slices.Equal(got, c.want) {
			t.Errorf("%q: got %q, want %q", c.words, got, c.want)
		}
	}
}

func TestCompleteAfterDoubleDash(t *testing.T) {
	s := Spec{
		Bin:      "kdbx",
		Flags:    []Flag{{Name: "--secrets"}, {Name: "--dry-run", Bool: true}},
		Args:     Anything,
		Commands: []Command{{Name: "show"}},
	}
	if got := s.Complete([]string{"--secrets", "a", "--", "show", ""}); !slices.Equal(got, []string{DirDefault}) {
		t.Fatalf("got %q", got)
	}
	if got := s.Complete([]string{"--dry-run", "sh"}); !slices.Equal(got, []string{"show"}) {
		t.Fatalf("got %q", got)
	}
}

func TestFromFlagSet(t *testing.T) {
	fs := flag.NewFlagSet("x", flag.ContinueOnError)
	fs.Bool("b", false, "")
	fs.String("s", "", "")
	got := FromFlagSet(fs)
	want := []Flag{{Name: "-b", Bool: true}, {Name: "-s"}}
	if len(got) != 2 || got[0].Name != want[0].Name || !got[0].Bool || got[1].Name != want[1].Name || got[1].Bool {
		t.Fatalf("got %+v", got)
	}
}

func TestWith(t *testing.T) {
	got := With([]Flag{{Name: "-a"}, {Name: "-b", Bool: true}}, Flag{Name: "--a", Dirs: true}, Flag{Name: "-c", Files: true})
	if len(got) != 3 || !got[0].Dirs || got[0].Name != "-a" || !got[1].Bool || got[2].Name != "-c" {
		t.Fatalf("got %+v", got)
	}
}
