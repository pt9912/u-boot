# docs/plan/planning/next

Für die nächste Iteration eingeplante Slice- und Tranchen-Pläne. Inhalte
sind ausgereift genug, um in `in-progress/` übernommen zu werden.

Übergang nach `in-progress/` per `git mv`, sobald die Arbeit beginnt.
Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-harness-dcheck-bump-reviewer-skill`](../done/slice-harness-dcheck-bump-reviewer-skill.md) | Auffrischung des Reviewer-Skills (Baseline v6.18.1); d-check-Bump liegt bei der ai-harness-init-Adoption |
| [`slice-harness-internal-readmes-festlegungen-in-spec`](../done/slice-harness-internal-readmes-festlegungen-in-spec.md) | Festlegungen aus den `internal/`-READMEs in die Spezifikation, READMEs auf Zweck und Verweis kürzen |
| [`slice-v1-sandbox-fallback-host-ausweisung`](slice-v1-sandbox-fallback-host-ausweisung.md) | Daemon-, Egress- und `vfs`-Fallback host-seitig ermitteln (Statusdateien im Container, `docker cp`, [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md)) und in `u-boot doctor` sowie in `generate devcontainer` und `init` ausweisen ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Lastenheft unverÃ¤ndert); startet nach [`slice-cr-sandbox-fuse-startfehler`](../open/slice-cr-sandbox-fuse-startfehler.md) |

Zuvor: Der Konsolidierungs-Slice
[`slice-v1-cli-json-envelope-consolidation`](../done/slice-v1-cli-json-envelope-consolidation.md) ist am 2026-06-08
vollständig nach `done/` ausgeliefert (T0→T4). Der JSON-CLI-Cluster
[`slice-v1-cli-json-dry-run`](../done/slice-v1-cli-json-dry-run.md) ist ebenfalls vollständig in `done/`
(9/9 Folge-Slices + T_close).
