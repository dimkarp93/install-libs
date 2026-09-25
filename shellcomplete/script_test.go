package shellcomplete

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScripts(t *testing.T) {
	s := Spec{Bin: "my-tool"}
	for _, sh := range Shells {
		script, err := s.Script(sh)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"my-tool " + CompleteCommand, DirNone, DirDirs, DirDefault} {
			if !strings.Contains(script, want) {
				t.Errorf("%s: no %q", sh, want)
			}
		}
		if strings.Contains(script, "{{") {
			t.Errorf("%s: template not filled", sh)
		}
		if _, err := exec.LookPath(sh); err != nil {
			continue
		}
		cmd := exec.Command(sh, "-n")
		cmd.Stdin = strings.NewReader(script)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("%s -n: %v\n%s", sh, err, out)
		}
	}
	if !strings.Contains(mustScript(t, s, "bash"), "complete -F _my_tool my-tool") {
		t.Error("bash function name")
	}
	if _, err := s.Script("fish"); err == nil {
		t.Error("fish accepted")
	}
}

func mustScript(t *testing.T, s Spec, sh string) string {
	t.Helper()
	out, err := s.Script(sh)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestBashScriptRuns(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash")
	}
	dir := t.TempDir()
	fake := "#!/bin/sh\n[ \"$1\" = __complete ] || exit 1\nshift\necho \"$#:$*\"\necho second\n"
	if err := os.WriteFile(filepath.Join(dir, "my-tool"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	script := mustScript(t, Spec{Bin: "my-tool"}, "bash")
	run := script + `
COMP_WORDS=(my-tool sync "")
COMP_CWORD=2
_my_tool
printf '%s|' "${COMPREPLY[@]}"
`
	cmd := exec.Command("bash", "-c", run)
	cmd.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := string(out); got != "2:sync |second|" {
		t.Fatalf("got %q", got)
	}
}

func TestHandle(t *testing.T) {
	s := testSpec()
	var out, errb bytes.Buffer
	if code, ok := s.Handle(&out, &errb, []string{CompleteCommand, "sy"}); !ok || code != 0 || out.String() != "sync\n" {
		t.Fatalf("__complete: %d %v %q", code, ok, out.String())
	}
	out.Reset()
	if code, ok := s.Handle(&out, &errb, []string{"completion", "zsh"}); !ok || code != 0 || !strings.HasPrefix(out.String(), "#compdef tool") {
		t.Fatalf("completion: %d %v", code, ok)
	}
	if code, ok := s.Handle(&out, &errb, []string{"completion"}); !ok || code != 2 {
		t.Fatalf("completion without shell: %d %v", code, ok)
	}
	if code, ok := s.Handle(&out, &errb, []string{"completion", "fish"}); !ok || code != 2 {
		t.Fatalf("completion fish: %d %v", code, ok)
	}
	if code, ok := s.Handle(&out, &errb, []string{"install-completions", "a", "b"}); !ok || code != 2 {
		t.Fatalf("extra args: %d %v", code, ok)
	}
	if _, ok := s.Handle(&out, &errb, []string{"sync"}); ok {
		t.Fatal("sync handled")
	}
	if _, ok := s.Handle(&out, &errb, nil); ok {
		t.Fatal("empty handled")
	}
}
