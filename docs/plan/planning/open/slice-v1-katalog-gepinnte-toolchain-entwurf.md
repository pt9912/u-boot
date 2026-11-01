# Slice V1: Katalog-Feature für eine gepinnte Werkzeugkette — Entwurf (Katalog-Eintragstyp, CR-Wortlaut, Messung zur ENV- und Cache-Mechanik), damit der Hand-written-Teil der eigenen `.devcontainer/Dockerfile` entfällt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle; Vorschlag `v0.18.0` (die Feature-Katalog-Welle der Roadmap *Nächste Wellen*, in der [`slice-v1-katalog-feature-jdk-gepinnt`](slice-v1-katalog-feature-jdk-gepinnt.md) ebenfalls vorgeschlagen ist). Anlass: Auftrag des Projektinhabers 2026-11-01 — "Ja ein CR-Slice aufsetzen (so dass wir den ganzen 'Hand-written part' ersetzen können)".

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog, Dockerfile-Blöcke mit Pin), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Konfiguration), [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Pin-Politik, Hebung dreier Stellen), [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog und Drittanbieter), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Muster Release-Archiv mit Prüfsumme), [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (lokaler Bau-Modus des eigenen Devcontainers).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (prüfen); gelesen beim Start, sie stehen hier nur als Zeiger. Der Verweis zeigt **aufwärts**.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neuer Katalog-Eintragstyp "gepinnte Werkzeugkette" — ein Dockerfile-Block, der Version und Release-Archiv samt SHA-256 je Architektur fest im Template trägt (`sha256sum -c` vor dem Entpacken, Muster der Einträge `ollama`, `gh`, `glab`) — trägt die Werkzeuge des lokalen Bau-Modus, und das eigene Devcontainer-Projekt von `u-boot` nutzt ihn, **sodass der handgeschriebene Teil der eigenen `.devcontainer/Dockerfile` (unter dem Managed Block) vollständig entfällt**: der Bau des Devcontainer-Images ist dann ganz erzeugtes Ergebnis, und eine Hebung von `GO_VERSION`/`GOLANGCI_LINT_VERSION` berührt die Template-Stellen des Katalogs statt eines Hand-written-Teils.

**Befund** (gelesen am 2026-11-01; was nicht gelesen oder gemessen ist, steht je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Der Hand-written-Teil.** `.devcontainer/Dockerfile` Zeilen 75 bis 110 (`# ---- Hand-written part`) installiert nach dem Managed Block: Go `1.27.3` und `golangci-lint` `v2.14.0` als Release-Archive mit `sha256sum -c` (je `amd64`/`arm64`, `ARG`-Zeilen 83 bis 88), und setzt die Umgebung des lokalen Bau-Modus (`ENV`-Zeilen 97 bis 107): `BUILD_MODE=local`, `GOFLAGS="-mod=readonly -buildvcs=false"`, `GOTOOLCHAIN=local`, und `GOCACHE`/`GOMODCACHE` unter `/workspaces/u-boot/.cache/…`, damit die Caches das `--recreate` des Sandbox-Profils überstehen. Der Kommentar nennt die Kopplung: "Hebung betrifft drei Stellen: Dockerfile (ARG), repo.mk and this file; the versions must stay equal."
- **Die Pin-Politik dreht schon um drei Stellen.** Root-`Dockerfile` Zeilen 23 bis 27 (`# Pin policy (LH-FA-BUILD-003)` … "Hebung betrifft drei Stellen: this file, repo.mk and .devcontainer/Dockerfile (hand-written part behind the managed block)") und `repo.mk` Zeilen 47 bis 50 (`GO_VERSION ?= 1.27.3`, `GOLANGCI_LINT_VERSION ?= v2.14.0`). Der Hand-written-Teil ist die dritte, handgetragene Kopie derselben Pins; der Eintragstyp verschiebt die dritte Stelle in die erzeugte Form — die Kopplung bleibt und wird in §4 als Frage aufgerufen, wo ihre eine Quelle liegen soll.
- **Das Muster der Dockerfile-Blöcke existiert.** `featureCatalogueEntry.dockerfileBlock` (`internal/hexagon/application/devcontainer_features.go`, Zeilen 165 bis 171): Eintrag ohne OCI-Quelle, "the pin stays in the template"; drei Einträge (`ollama` Sandbox nur, `gh` und `glab` beide Profile, Zeilen 305 bis 316). Version und Prüfsummen stehen als `ARG` im Dockerfile-Template; die genauen `ARG`- und Block-Zeilen von `templates/devcontainer/Dockerfile.tmpl` sind **nicht gelesen** (Vorbild-Blöcke liest der Implementer im Umsetzungs-Slice).
- **Der Katalog-Eintrag `go` ist heute ein externes OCI-Feature** mit der Standardversion `1` (Zeilen 257 bis 263, Quelle `ghcr.io/devcontainers/features/go`); es installiert nur Go, ohne Pin des Releases und ohne `golangci-lint`. Der Hand-written-Teil wird von ihm nicht abgedeckt. Ob und wie der Eintrag mit dem neuen Typ zusammengeht (ersetzen, ergänzen, umbenennen), entscheidet der Projektinhaber (§4 Frage 2).
- **Das eigene Projekt ist ein Sandbox-Projekt.** `u-boot.yaml` der Wurzel: `profile: sandbox`, `nestedRuntime: docker`, Features `gh`, `node` (extern, `2.1.1`), `claude-code` (`1.0.5`), `ollama`, `homebrew` (gelesen); die erzeugte `.devcontainer/devcontainer.json` trägt `workspaceFolder: /workspaces/u-boot` und die drei Volumes mit `${devcontainerId}`. Die ENV-Zeilen des Hand-written-Teils setzen deshalb auf diese Pfadkonvention — für einen generischen Katalog-Eintrag sind die Cache-Pfade **nicht** statisch formulierbar; über `containerEnv` der erzeugten `devcontainer.json` mit `${localWorkspaceFolder}` wäre es denkbar, ob die Dev-Container-CLI diese Form dort auflöst, ist **nicht gemessen** (Messung 1).
- **Was der Typ ersetzen muss, um "ganz" zu ersetzen:** (a) der Dockerfile-Block für Go und `golangci-lint` (Muster vorhanden), (b) die ENV-Zeilen (`BUILD_MODE`, `GOFLAGS`, `GOTOOLCHAIN`) — im Dockerfile-Block renderbar, aber `BUILD_MODE=local` gilt für dieses Projekt, nicht für jedes Projekt mit dem Feature — und (c) die Cache-Pfade. Der Umfang von (b) und (c) in einem generischen Eintrag ist die offene Entwurfsfrage: Ein Eintrag nur für (a) ersetzt den Hand-written-Teil schon zu sieben von zehn Zeilen, die ENV bleibt dann als kleiner Rest oder wandert in `containerEnv`/`remoteEnv` der erzeugten Konfiguration.

**Ausdrücklich NICHT in diesem Slice (Entwurf, kein Produktcode) — je Punkt mit Begründung:**

- **Der Katalog-Eintrag JDK** — eigenes Vorhaben mit eigenem Hersteller-Auswahl-Entscheidungsraum ([`slice-v1-katalog-feature-jdk-gepinnt`](slice-v1-katalog-feature-jdk-gepinnt.md), §4 Wege A bis C); der Eintragstyp hier ist seine Voraussetzung, nicht seine Umsetzung. Der JDK-Slice wird auf diesen Entwurf verweisen und ihn nicht duplizieren.
- **Hebung der Versionen von Go und `golangci-lint`** — bleibt die Routine-Hebung der Pin-Politik [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik); der Slice ändert keine Version.
- **Der Umsetzungs-Code** (Template-Block, Renderer, Tests, das Ersetzen im eigenen `.devcontainer/` und der Regenerate-Vergleich) — eigener Folge-Slice nach der Entscheidung (§2, Liefer-Punkt 3).
- **Das externe Projekt oder eine andere Engine** — der Nachweis folgt dem üblichem Muster: Linux-Host `amd64`, lokale Engine, belegt bis zum Docker-Test.
- **Das Hand-manuelle Aufräumen alter Hilfs-Images oder Caches** — kein Gegenstand des Katalogs.

## 2. Definition of Done

- [ ] **Messung zur ENV- und Cache-Mechanik** (nicht gemessen, nicht angenommen; Protokoll unter `docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/`, Skripte unter `tools/messung/`, linux-host `amd64`, lokale Engine, je Zelle zwei Läufe): (M1) löst die Dev-Container-CLI `${localWorkspaceFolder}` (oder eine Form des Workspace-Pfads) in `containerEnv` bzw. `remoteEnv` auf, sodass `GOCACHE`/`GOMODCACHE` je Projekt unter dem Workspace-Volume liegen; (M2) greift ein `ENV` des Dockerfile-Blocks im Container, und wie verhalten sich `ENV` aus dem Block und `containerEnv` aus der Konfiguration zu einander; (M3) übersteht der Cache unter dem Workspace-Volume ein `up --recreate`.
- [ ] **CR-Wortlaut und Träger** (Alt/Neu je Stelle, nicht geschrieben bis zur Freigabe): neuer Katalog-Eintragstyp (Name) in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) mit der Pin-Zusage des Musters und der Grenze "nur soweit ein Docker-Test ihn belegt"; die ENV-/Cache-Form nach Messung 1; Freigabe durch den Projektinhaber je Stelle; falls der Eintragstyp die Katalog-Regeln verlässt (etwa `version:` wählbar), ADR-Gegenprüfung gegen [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md).
- [ ] **Folge-Slice-Plan** in `docs/plan/planning/open/` mit der vollständigen Plan-Tabelle: Template-Block (Go + `golangci-lint` je Weg), Renderer/Daten, Katalog-Regeln, Golden Cases, Docker-Test (Bau läuft, `go version` und `golangci-lint version` melden die Pins, Cache liegt im Workspace-Volume), Regenerate-Vergleich der eigenen `.devcontainer/` (der Hand-written-Teil entfällt, Blöcke und Konfiguration erzeugt), Hebung-Abgleich über die verbleibenden Stellen. `Verantwortlich:` gesetzt; Start nach den §4-Entscheidungen.

## 3. Plan (vor Code)

Dieser Slice schreibt kein Produktcode. Gelesene Stellen am 2026-11-01.

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/`, `tools/messung/<name>/` | neu | M1 bis M3; Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) |
| `spec/lastenheft.md` | **in diesem Slice nicht** | der CR-Wortlaut entsteht hier als Alt/Neu-Vorschlag (§2, Liefer-Punkt 2); geschrieben wird er nach der Freigabe — im Umsetzungs-Slice oder in einem eigenen CR-Slice, wie das Muster [`slice-cr-lastenheft-messstaende`](../done/slice-cr-lastenheft-messstaende.md) es trägt |
| diese Datei, `docs/plan/planning/open/README.md`, Umsetzungs-Slice | neu / update | Zeile im Bestand; der Umsetzungs-Slice trägt die Folge-Plan-Tabelle |
| `internal/hexagon/application/**`, `tests/**`, e2e | **ungeöffnet** | bewusst: der Entwurf liest die Fundstellen (befundet), ändert sie nicht |

## 4. Trigger

**Entscheidungen des Projektinhabers vor dem Start des Umsetzungs-Slices** (dieser Entwurf bereitet sie vor, der Umsetzungs-Slice braucht sie):

1. **Eintragstyp:** **A:** fest gepinnte Version im Template (wie `ollama`, `gh`, `glab`; Hebung ändert das Template); **B:** Version aus einer gepinnten Tabelle wählbar (`version:` bei einem Dockerfile-Block verlässt die heutige Abweisungsregel, [`slice-v1-katalog-feature-jdk-gepinnt`](slice-v1-katalog-feature-jdk-gepinnt.md) §4 Weg B trägt dieselbe Frage). Empfehlung des Entwurfs: **A**, weil die Kopplung an die CI-Pins ([`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)) die Version ohnehin zentral wählt.
2. **Schlüssel-Namen und Verhältnis zum Eintrag `go`:** zwei neue Schlüssel (`go-toolchain`, `golangci-lint`) oder einer (Werkzeugkette in einem Block); bleibt das externe `go`-Feature bestehen (für Projects, die kein gepinntes Release brauchen)? Empfehlung: zwei Schlüssel, eins je Werkzeug, das externe `go` bleibt.
3. **Umfang der ENV-/Cache-Form nach der Messung:** Dockerfile-`ENV` im Block, `containerEnv`/`remoteEnv` der erzeugten Konfiguration, oder ein Rest bleibt im eigenen Projekt handgeschrieben (sichtbar in den Rohdaten von M1/M2; "ganz ersetzen" heißt hier nicht um jeden Preis).
4. **Wo die zwei Versionen künftig ihre eine Quelle haben:** die Hebung berührt heute Root-`Dockerfile`, `repo.mk` und `.devcontainer/Dockerfile`; mit dem Typ tritt eine vierte, erzeugte Stelle an die dritte (Template-ARG). Zu entscheiden, ob die Pin-Politik in [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) die Stellenliste benennt und wie der Abgleich geprüft wird (Test oder Gate-Prüfung der drei/vier Stellen).
5. **Release-Zuordnung** (`v0.18.0` Katalog-Welle, gemeinsam mit dem JDK-Slice) und **ADR-Bedarf** (der Eintragstyp ist dem Muster von [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) treu; ein ADR steht an, wenn Frage 1/2 die Katalog-Regeln ändern).

**Start des Entwurfs** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; Messung M1 bis M3 geplant; kein anderer Slice in `in-progress/` (WIP-Limit 1). **Start des Umsetzungs-Slices:** die §4-Fragen entschieden, der CR-Wortlaut je Stelle freigegeben, der CR in `done/` falls gesondert.

**Rückführungen — vorab benannt:** `in-progress` → `open`: der Projektinhaber entscheidet die §4-Fragen nicht. `in-progress` → `next`: die Messung zeigt, dass keine ENV-/Cache-Form der drei trägt (M1/M2 beide nicht auflösend); der Umsetzungs-Plan geht dann mit einem benannten Hand-written-Rest zum Projektinhaber.

## 5. Closure-Trigger

DoD vollständig; `make gates` grün (`docs-check` 0 Befunde); Messprotokoll mit Rohdaten, je Zelle zwei Läufe, datiert; CR-Wortlaut je Stelle vorgelegt und freigegeben; Umsetzungs-Slice in `open/` mit Plan-Tabelle; Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 mit Ausgang; die drei Paarungen getragen.

## 6. Risiken und offene Punkte

1. `${localWorkspaceFolder}` in `containerEnv`/`remoteEnv` wird von der CLI nicht aufgelöst (M1) — dann trägt die Cache-Form nur als Dockerfile-`ENV` mit festen Pfaden, was projektegens ist, oder der Rest bleibt Hand-written. — **Ausgang:** bei Closure.
2. Es entsteht eine **zweite Definition** der Pins und driftet zur CI-Pin-Politik (drei Stellen werden vier). — **Ausgang:** bei Closure (Frage 4: Quelle und Prüfung benennen).
3. Das externe `go`-Feature und ein neuer `go-toolchain`-Block widersprechen einander, wenn ein Projekt beide einschaltet. — **Ausgang:** bei Closure (Frage 2; Verhältnis im CR benennen, Abweisung oder Warnung im Umsetzungs-Slice).
4. Ein ENV-Block im Dockerfile gilt für **jede** Konfiguration, die den Katalog-Eintrag aktiviert — ein Nutzer-Projekt erbt `GOFLAGS`/`GOTOOLCHAIN`-Werte, die sein Build verändern. — **Ausgang:** bei Closure (Wirkung je Form in M2 messen; CR formuliert nur das Gemessene).
5. Ein Messwert wird in den CR geschrieben, den der Docker-Test nicht belegt. — **Ausgang:** bei Closure ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md); Beleg-Tabelle im Umsetzungs-Slice).
6. Cache unter dem Workspace-Volume verliert Rechte oder Größe beim Regenerate (M3). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden in diesem Entwurf nur `docs/messungen/`, `docs/plan/` und `tools/` (DOK, WZG) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area). Der Umsetzungs-Slice berührt dann `internal/hexagon/application` (APP), `internal/e2e` (E2E), `spec/`, `harness/` (DOK); Greenfield derselben Deklaration.

**Vorgelagert — offene Beobachtungen sichten** (Register 2026-11-01 durchgesehen, Zähler nicht neu erhoben): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 5), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Beleg-Tabelle des Umsetzungs-Slices), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Prüflinse des Befunds, Zeilen genannt). Evidenz-/Diskrepanz-Risiko: mittel (M1/M2 entscheiden über die Form des CR). Alle berührten Sub-Areas Greenfield — kein Modus-Begründungsblock nötig.