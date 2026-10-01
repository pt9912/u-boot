package docker

import "testing"

// slice-v1-recreate-detection: only `Container <name> Recreate` lines
// become plan actions; duplicates and other actions are ignored.
func TestParseComposePlan(t *testing.T) {
	t.Parallel()
	out := ` DRY-RUN MODE -  Image nginx:1.26-alpine Pulled
 DRY-RUN MODE -  Container demo-web-1  Recreate
 DRY-RUN MODE -  Container demo-web-1  Recreate
 DRY-RUN MODE -  Container demo-db-1  Running
 Container demo-api-1 Recreate
 Container demo-api-1 Error response from daemon: boom
`
	plan := parseComposePlan(out)
	if len(plan) != 2 || plan[0].Container != "demo-web-1" || plan[1].Container != "demo-api-1" || plan[0].Action != "Recreate" {
		t.Errorf("plan = %+v, want demo-web-1 and demo-api-1 once each", plan)
	}
	if got := parseComposePlan("nothing relevant\n"); len(got) != 0 {
		t.Errorf("plan = %+v, want empty", got)
	}
}
