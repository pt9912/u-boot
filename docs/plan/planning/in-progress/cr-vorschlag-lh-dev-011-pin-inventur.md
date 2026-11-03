# CR-Vorschlag: Pin-Inventur (`pinsync`) — die Werkzeugketten der `LH-FA-DEV-011`-Kopplung samt fünftem Ort `go.mod`

**Status:** vorgelegt (2026-11-03); wartet auf die Freigabe des Projektinhabers.
**Grundlage:** Projektinhaber 2026-11-03 („die Nummern müssen zusammenpassen"; „… auch alle anderen Tools, die für Entwicklung an diese Repo nötig sind: bash, sed, go, docker, …"); d-check v0.91.0 Modul `pinsync` (`DC-FA-PSYNC-001`, byte-gleicher Wert an ≥2 Dateistellen).

**Die vorgeschlagenen Pin-Gruppen** (je eine benannte pinsync-Regel, jede Quellen-Zeile ein RE2-Ausdruck mit einer Fanggruppe):

| Gruppe | Quellen-Stellen | Hinweis |
|---|---|---|
| `go` | Root-`Dockerfile` `ARG GO_VERSION=`, `repo.mk` `GO_VERSION ?=`, Katalog-Template `ARG GO_VERSION=`, `.devcontainer/Dockerfile` (erzeugter Block), **`go.mod` `go <dir>`** (fünfter Ort) | `go.mod` auf `1.27.0` gehoben (Minor = Pin `1.27.3`); jede Pin-Hebung zieht die Direktive fortan nach |
| `golangci-lint` | dieselben Stellen minus `go.mod`, `GOLANGCI_LINT_VERSION` | `v`-Präfix im Ausdruck |
| `cli` | `package.json` Version-Feld ⇔ `DCHECK`/CLI-Image-Tag-Nenung | das gepinnte CLI-Image |
| `jdk` | Template-`ARG JDK_VERSION=` als Einstellen-Stelle | **offen:** eine zweite echte Quelle existiert nicht (die e2e-Konstante ist ein Test, kein Pin-Ort) — pinsync braucht ≥2 Stellen; Vorschlag: JDK **nicht** gruppen, sondern als Grenze im Inventar-Absatz, dessen Wortlaut bei der [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Stelle frei gegeben wird |
| Basis-Image-Gruppe | `golang:*`-Tags in Root-`Dockerfile`/templates | Grenze: `bash`/`sed`/`coreutils`/`curl` tragen keine Einzel-Pins — sie kommen aus dem Basis-Image; die Inventur benennt das ausdrücklich, keine Einzel-Pin-Zusage |

**Der freizugebende Wortlaut (Stelle der `LH-FA-DEV-011`-Absatz):**

> **Pin-Inventur.** Die Version der zwei Go-Werkzeugketten-Schlüssel steht an den vier Kopplungsstellen und zusätzlich an der `go`-Direktive der `go.mod` dieses Repos; die Direktive trägt denselben Pin-Wert, und jede Hebung des Pins zieht sie fortan nach. `bash`, `sed` und weitere Shell-Werkzeuge des Entwicklungs- und Bau-Verlaufs tragen keine Einzel-Pins — ihre Versionen kommen mit dem Basis-Image des Pins. Ein `pinsync`-Wächter prüft die Gleichheit der Quellen-Zeilen je benannten Pin-Gruppe beim Gatter-Lauf; der Wächter prüft Gleichheit, keine Aktualität.

**Offen zur Freigabe:** die JDK-Gruppe (nicht gruppen, Grenze statt), die `go.mod`-Hebung auf `1.27.0` (Minor-Gleichheit je „zusammenpassen"), die Basis-Image-Grenze (keine Einzel-Pins für Shell-Werkzeuge).