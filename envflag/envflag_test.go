package envflag

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func newSet(retries *int) *Set {
	return New("NET_",
		Option{Name: "retries", Default: "5", Parse: func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			*retries = n
			return nil
		}},
		Option{Name: "retry-all-errors", Env: "NET_RETRY_ALL", Default: "0"},
	)
}

func lookup(env map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := env[k]
		return v, ok
	}
}

func TestEnvNames(t *testing.T) {
	var n int
	s := newSet(&n)
	if got := s.Env("retries"); got != "NET_RETRIES" {
		t.Fatalf("got %q", got)
	}
	if got := s.Env("retry-all-errors"); got != "NET_RETRY_ALL" {
		t.Fatalf("got %q", got)
	}
	if s.Has("nope") || s.Env("nope") != "" {
		t.Fatal("unknown option reported")
	}
}

func TestPriority(t *testing.T) {
	var n int
	s := newSet(&n)
	if s.Source("retries") != SourceDefault || s.Value("retries") != "5" {
		t.Fatalf("default: %s %s", s.Value("retries"), s.Source("retries"))
	}
	if err := s.LoadEnv(lookup(map[string]string{"NET_RETRIES": "3", "NET_RETRY_ALL": ""})); err != nil {
		t.Fatal(err)
	}
	if n != 3 || s.Source("retries") != SourceEnv {
		t.Fatalf("env: %d %s", n, s.Source("retries"))
	}
	if s.Source("retry-all-errors") != SourceDefault {
		t.Fatal("empty env must be ignored")
	}
	if err := s.Set("retries", "7", SourceFlag); err != nil {
		t.Fatal(err)
	}
	if n != 7 || s.Source("retries") != SourceFlag {
		t.Fatalf("flag: %d %s", n, s.Source("retries"))
	}
}

func TestLoadEnvError(t *testing.T) {
	var n int
	s := newSet(&n)
	err := s.LoadEnv(lookup(map[string]string{"NET_RETRIES": "x"}))
	if err == nil || !strings.HasPrefix(err.Error(), "NET_RETRIES: ") {
		t.Fatalf("got %v", err)
	}
	if s.Source("retries") != SourceDefault || s.Value("retries") != "5" {
		t.Fatal("failed parse must not change the value")
	}
	if err := s.Set("nope", "1", SourceFlag); err == nil {
		t.Fatal("unknown option accepted")
	}
}

func TestReportAndHandle(t *testing.T) {
	var n int
	s := newSet(&n)
	_ = s.LoadEnv(lookup(map[string]string{"NET_RETRIES": "3"}))
	want := "NET_RETRIES=3 (env)\nNET_RETRY_ALL=0 (default)\n"
	var buf bytes.Buffer
	if !s.HandleEnvs(&buf, []string{"--envs"}) || buf.String() != want {
		t.Fatalf("got %q", buf.String())
	}
	buf.Reset()
	if s.HandleEnvs(&buf, []string{"fetch", "--envs"}) || s.HandleEnvs(&buf, nil) || buf.Len() != 0 {
		t.Fatal("--envs handled outside the first argument")
	}
}

func TestParseBool(t *testing.T) {
	for _, v := range []string{"1", "true", "YES", " on "} {
		if b, err := ParseBool(v); err != nil || !b {
			t.Errorf("%q: %v %v", v, b, err)
		}
	}
	for _, v := range []string{"0", "false", "no", "off"} {
		if b, err := ParseBool(v); err != nil || b {
			t.Errorf("%q: %v %v", v, b, err)
		}
	}
	if _, err := ParseBool("maybe"); err == nil {
		t.Error("maybe accepted")
	}
}

func TestDuplicatePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("no panic")
		}
	}()
	New("", Option{Name: "a"}, Option{Name: "a"})
}
