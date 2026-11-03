# CR-Vorschlag: SPEC-031 — der einzig-gültige Bau-Pfad (`bau.sh` im Devcontainer; `BUILD_MODE` entfällt) und die Träger-Ebenen je Gatter

**Status:** freigegeben wie vorgelegt (Projektinhaber 2026-11-03).
**Mess-Basis:** [`docs/messungen/slice-v1-bau-nur-devcontainer/README.md`](../../../messungen/slice-v1-bau-nur-devcontainer/README.md) samt sechs Rohdaten — M1 (Image/Werkzeug-Lage, Kontext-Konvention, Cache-Warm-Grenze), M2 (DIND trägt), M3 (gates im Container Exit 0, 33–62 s, 92,20 %), M4 (test-docker am verschachtelten Daemon 12 FAIL/1758 PASS — Hybrid-Grund), M5 (Flagge honoriert die erzeugende Engine, Sockelform der Ebene 1). **Entscheidungen des Projektinhabers:** drei bejaht (2026-11-03), test-docker-Platzierung Hybrid je Messbasis (2026-11-03).

**Stelle (a) — [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Technik):**
> **Der Bau-Pfad.** Das Repo trägt einen Bau-Pfad: die Rezepte von `scripts/bau.sh` laufen in der Devcontainer-Umgebung (`BUILD_MODE=local` ist der einzige Modus-Wert, eine zweite Modus-Entscheidung entfällt). Die Gatter tragen ihre Ebenen je Messbasis: `lint`/`test`/`coverage`/`docs-check`/`baseline-verify`/`ablage-pruefung` laufen im Devcontainer-Container — dessen verschachtelter Daemon trägt die Container-Züge der Dokumentation (M2/M3) — und `test-docker` läuft an dem Docker, der den Devcontainer-Container erzeugt (Sockelform; der verschachtelte Daemon trägt die docker-getaggten Tests nicht — Messung M4). `modus-pruefung` entfällt; ein Aufruf mit dem entfallenen Schalter (`BUILD_MODE=docker`) ist ein unbekanntes Verhalten und endet mit dem make-Fehler der ungeöffneten Variable.

**Stelle (b) — `AGENTS.md` §3.1:**
> Der Host braucht Docker und GNU make (Bestand); die Bau-Nutzung läuft über den Devcontainer des Repos: die Rezepte von `scripts/bau.sh` sind der einzige Pfad; `make test-docker` bindet an den Docker-Daemon außerhalb des Devcontainer-Containers (Sockelform), die übrigen Gatter laufen im Container.

**Stelle (c) — CI:** die Workflows laufen `gates`/`ci`/`fullbuild` im Devcontainer-Container (`docker run` auf dem gebauten Image — der M5-Nachweis: unprivilegierte Ebene 1 mit Socket) und `test-docker` am Runner-Docker in der Sockelform (M4-Grund); `integration-docker` bleibt der Nachweis-Träger.

**Stelle (d) — ADR:** eine Folge-ADR (`Supersedes` [ADR-0037](../../../../docs/plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md)) trägt die Entscheidung (Index-Zeile je §5-Regel).

**Grenzen je Messbasis:** der Netzbau des Images ist nicht belegt (cache-warm, Grenze M1); die Ursachen von zwei der sechs Fehl-Klassen des M4 sind nur bis zur Meldungs-Stelle bestimmt (der Hybrid-Weg umgeht die Frage, keine Berichtigung zugesagt); die `modus-pruefung`-Entfernung ist je Messung M2 getragen (die Prüfung lief im verbleibenden Stand OK — sie entfällt, weil die Frage, die sie trug, entfällt).