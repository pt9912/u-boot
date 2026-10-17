# Zusage ohne CI-Evidenz: der nested-Podman-Build-Test endet auf dem Runner womöglich im Skip

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Test liegt dort, der Workflow `integration.yml` führt ihn aus.

Die Tests `TestE2E_SandboxDevcontainer_NestedPodman_Build` und `…_BuildEgress` überspringen (`t.Skip`), wenn `kernel.apparmor_restrict_unprivileged_userns` auf dem Docker-Host 1 ist. Ob das auf `ubuntu-latest` (Workflow `integration.yml`, `make test-docker`) der Fall ist, ist ungemessen; ist es so, hat die Zusage zu `podman build` ([`LH-FA-DEV-007`](../../../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Lastenheft 0.12.0) auf amd64 und in CI keine Maschinenevidenz, und ein SKIP belegt nichts.

**Trigger:** Nach dem Push zeigt der CI-Lauf, ob die Build-Tests auf `ubuntu-latest` überspringen. Die Entscheidung des Projektinhabers (2026-10-17) lautet: keine Workflow-Änderung (eine Änderung wäre eine Gate-Änderung nach [`AGENTS.md`](../../../../../../AGENTS.md) §3.6); die Beobachtung wird geführt, bis ein Lauf auf amd64 vorliegt.
