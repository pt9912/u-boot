package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func entry(t *testing.T, path, value string) driving.ConfigPathValue {
	t.Helper()
	return driving.ConfigPathValue{Path: mustCfgPath(t, path), Value: value}
}

func runConfig(t *testing.T, uc *fakeConfigUseCase, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := newAppWithConfigStub(uc).Execute(context.Background(), args, &stdout, &stderr)
	return stdout.String(), err
}

// config get with several paths: one value per line in argument order.
func TestConfigGet_MultiPath_Human(t *testing.T) {
	uc := &fakeConfigUseCase{getManyResp: driving.ConfigGetManyResponse{Entries: []driving.ConfigPathValue{
		entry(t, "project.name", "demo"), entry(t, "devcontainer.enabled", "true")}}}
	out, err := runConfig(t, uc, "config", "get", "project.name", "devcontainer.enabled")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out != "demo\ntrue\n" {
		t.Errorf("stdout = %q", out)
	}
	if len(uc.getManyReq.Paths) != 2 {
		t.Errorf("request paths = %v", uc.getManyReq.Paths)
	}
}

// --json: data.entries[]; --json-array forces it for one path; a single
// path without the flag keeps the legacy data shape.
func TestConfigGet_MultiPath_JSON(t *testing.T) {
	uc := &fakeConfigUseCase{
		getResp: driving.ConfigGetResponse{Path: mustCfgPath(t, "project.name"), Value: "demo"},
		getManyResp: driving.ConfigGetManyResponse{Entries: []driving.ConfigPathValue{
			entry(t, "project.name", "demo"), entry(t, "devcontainer.enabled", "true")}},
	}
	type env struct {
		Data struct {
			Path    string                         `json:"path"`
			Value   string                         `json:"value"`
			Entries []struct{ Path, Value string } `json:"entries"`
		} `json:"data"`
	}
	parse := func(args ...string) env {
		t.Helper()
		out, err := runConfig(t, uc, args...)
		if err != nil {
			t.Fatalf("Execute: %v", err)
		}
		var e env
		if jerr := json.Unmarshal([]byte(out), &e); jerr != nil {
			t.Fatalf("envelope: %v\n%s", jerr, out)
		}
		return e
	}
	if e := parse("--json", "config", "get", "project.name", "devcontainer.enabled"); len(e.Data.Entries) != 2 || e.Data.Entries[1].Path != "devcontainer.enabled" {
		t.Errorf("multi: %+v", e.Data)
	}
	if e := parse("--json", "config", "get", "--json-array", "project.name"); len(e.Data.Entries) != 2 {
		t.Errorf("--json-array must use the entries shape: %+v", e.Data)
	}
	if e := parse("--json", "config", "get", "project.name"); e.Data.Path != "project.name" || e.Data.Entries != nil {
		t.Errorf("single path must keep the legacy shape: %+v", e.Data)
	}
}

// config set with several pairs → SetMany; odd argument count → exit 2.
func TestConfigSet_MultiPair(t *testing.T) {
	uc := &fakeConfigUseCase{setManyResp: driving.ConfigSetManyResponse{Entries: []driving.ConfigSetEntry{
		{Path: mustCfgPath(t, "project.name"), OldValue: "a", NewValue: "b"},
		{Path: mustCfgPath(t, "devcontainer.enabled"), OldValue: "true", NewValue: "true"}}}}
	out, err := runConfig(t, uc, "config", "set", "project.name", "b", "devcontainer.enabled", "true")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(uc.setManyReq.Items) != 2 || uc.setManyReq.Items[0].Value != "b" {
		t.Errorf("request = %+v", uc.setManyReq)
	}
	if !strings.Contains(out, "config: project.name a → b.") || !strings.Contains(out, "already true") {
		t.Errorf("stdout = %q", out)
	}

	jsonOut, err := runConfig(t, uc, "--json", "config", "set", "project.name", "b", "devcontainer.enabled", "true")
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	var env struct {
		Data struct {
			NoOp    bool `json:"noOp"`
			Entries []struct {
				Path string `json:"path"`
				NoOp bool   `json:"noOp"`
			} `json:"entries"`
		} `json:"data"`
	}
	if jerr := json.Unmarshal([]byte(jsonOut), &env); jerr != nil || len(env.Data.Entries) != 2 || env.Data.NoOp || !env.Data.Entries[1].NoOp {
		t.Errorf("json = %s (%v)", jsonOut, jerr)
	}

	for _, args := range [][]string{{"config", "set", "project.name"}, {"config", "set", "a", "b", "c"}} {
		_, err := runConfig(t, &fakeConfigUseCase{}, args...)
		if got := cli.ExitCode(err); got != 2 {
			t.Errorf("%v: ExitCode = %d, want 2 (err=%v)", args, got, err)
		}
	}
}

// config list: human path=value lines; JSON data.entries[]; preview
// flags are rejected like on the other read-only forms.
func TestConfigList(t *testing.T) {
	uc := &fakeConfigUseCase{listResp: driving.ConfigListResponse{Entries: []driving.ConfigPathValue{
		entry(t, "devcontainer.enabled", "true"), entry(t, "project.name", "demo")}}}
	out, err := runConfig(t, uc, "config", "list")
	if err != nil || out != "devcontainer.enabled=true\nproject.name=demo\n" {
		t.Errorf("human: err=%v out=%q", err, out)
	}
	jsonOut, err := runConfig(t, uc, "--json", "config", "list")
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	var env struct {
		Subcommand string `json:"subcommand"`
		Data       struct {
			Entries []struct{ Path, Value string } `json:"entries"`
		} `json:"data"`
	}
	if jerr := json.Unmarshal([]byte(jsonOut), &env); jerr != nil || env.Subcommand != "list" || len(env.Data.Entries) != 2 {
		t.Errorf("json = %s (%v)", jsonOut, jerr)
	}
	empty, err := runConfig(t, &fakeConfigUseCase{}, "--json", "config", "list")
	if err != nil || !strings.Contains(empty, `"entries":[]`) && !strings.Contains(empty, `"entries": []`) {
		t.Errorf("empty list must be entries:[]: err=%v %s", err, empty)
	}
	_, err = runConfig(t, uc, "config", "list", "--dry-run")
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("--dry-run on list: ExitCode = %d, want 2", got)
	}
}
