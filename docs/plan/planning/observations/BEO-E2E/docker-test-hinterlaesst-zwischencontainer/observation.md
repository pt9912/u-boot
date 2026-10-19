# Docker-Test hinterlässt den Zwischencontainer eines absichtlich fehlgeschlagenen Baus

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Fehler liegt im Test, nicht im Produkt; er trifft jeden Lauf von `make test-docker` auf einem geteilten Docker-Host.

Ein Docker-Test, der ein Image absichtlich fehlschlagend baut (Negativprobe zu einem Pin), hinterlässt mit dem klassischen Builder des Test-Images den Container des fehlgeschlagenen Schritts, solange der Bau ohne `--force-rm` läuft und kein Cleanup ihn entfernt. Der Container trägt von Docker vergebenen Namen und kein Label; nichts ordnet ihn dem Test zu. Kein Gate fängt das, weil der Test selbst grün bleibt.

## Benannt, nicht gezählt

- **Zweites Vorkommen, gelesen, nicht gemessen:** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` baut in drei Teilfällen ohne `--force-rm` und ohne Cleanup. Gleicher Mechanismus ist vermutet; die Messung steht im Folge-Slice [`slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../../../open/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md).
- **CI-Runner nicht gemessen:** ob der Daemon des Runners denselben Container hinterlässt, ist offen.

**Trigger:** Der Folge-Slice misst die drei Teilfälle; ein Teilfall mit Zwischencontainer ist der zweite Beleg (eigene Evidence-Datei), ein Befund ohne Container (oder BuildKit-Verhalten) ist die Antwort, die der Stand festhält.
