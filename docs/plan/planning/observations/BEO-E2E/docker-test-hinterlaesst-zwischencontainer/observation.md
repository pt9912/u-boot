# Docker-Test hinterlässt den Zwischencontainer eines absichtlich fehlgeschlagenen Baus

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Fehler liegt im Test, nicht im Produkt; er trifft jeden Lauf von `make test-docker` auf einem geteilten Docker-Host.

Ein Docker-Test, der ein Image absichtlich fehlschlagend baut (Negativprobe zu einem Pin), hinterlässt mit dem klassischen Builder des Test-Images den Container des fehlgeschlagenen Schritts, solange der Bau ohne `--force-rm` läuft und kein Cleanup ihn entfernt. Der Container trägt von Docker vergebenen Namen und kein Label; nichts ordnet ihn dem Test zu. Kein Gate fängt das, weil der Test selbst grün bleibt.

## Benannt, nicht gezählt

- **Zweites Vorkommen, gezählt:** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` hinterließ in drei Teilfällen je einen Container, gleicher Mechanismus, Pin im Env; gemessen und geräumt im Folge-Slice [`slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../../../done/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (Evidence-Datei). Der Test trägt zusätzlich einen Wächter (`failedStepRemoved`), der unter BuildKit rot statt vakuum-grün wird.
- **Ollama-Wächterlücke:** `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `-q --force-rm` und trägt keinen Wächter; unter BuildKit (Wechsel des Test-Images) bliebe er grün, ohne dass ein Test es sagt. Offen, Ausgang beim Projektinhaber (Folge-Slice oder nicht).
- **CI-Runner nicht gemessen:** ob der Daemon des Runners denselben Container hinterlässt, ist offen.
- **Grenzen der Eingriffe:** zwei gleichzeitige Läufe gegen einen Daemon sind nicht ausgeschlossen; der Wächter hängt am Wortlaut des klassischen Builders.

**Trigger:** ein drittes Vorkommen in einem weiteren Docker-Test, die Entscheidung des Projektinhabers zur Ollama-Wächterlücke oder eine Messung des CI-Runners.
