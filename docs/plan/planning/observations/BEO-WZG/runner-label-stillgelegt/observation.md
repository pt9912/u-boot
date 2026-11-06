# Runner-Label stillgelegt: ein Label der Smoke-Matrix bekommt keinen Runner mehr

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die nächste Zeile der Modus-Deklaration für die Release- und Tap-Werkzeuge; die berührte Datei `packaging/homebrew-tap/.github/workflows/smoke.yml` liegt in keiner Sub-Area (Aspirantin, Plan-Abschnitt 8 des Vorgangs [`slice-v1-tap-smoke-runner-macos-15`](../../../done/slice-v1-tap-smoke-runner-macos-15.md)).

Die Matrix des Smoke-Workflows für den Homebrew-Tap nennt ein Runner-Label, das GitHub stilllegt: der Job bekommt keinen Runner (`runner_name` leer), die Datei bleibt unverändert. Der Workflow läuft nur bei einem Push auf `Formula/u-boot.rb` und per `workflow_dispatch` (Datei, Zeilen 10 bis 14, kein `schedule`); die Stilllegung fällt deshalb erst mit dem nächsten Formel-Push auf, nicht beim Stilllegungsdatum.

## Benannt, nicht gezählt

Keine.
