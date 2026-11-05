# Slice V1: Homebrew-Pakete im Devcontainer — `devcontainer.features.homebrew.packages` installiert eine frei wählbare Formel-Liste, ohne dass jedes Werkzeug einen Katalogeintrag braucht

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen (eine Zusage-Erweiterung im
Devcontainer-Bereich; keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Scope: Feature-Schema), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Tap als Quelle), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Abgrenzung: Pakete sind ungepinnt), [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md) (Bau-Pfad, nur als Rahmen).

**Berührte Spec-Stellen:** [`LH-FA-DEV-003.a`](../../../../spec/spezifikation.md#lh-fa-dev-003a--mechanik-der-freigabe-externer-feature-quellen) (Freigabe externer Feature-Quellen), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Katalog-Features im Sandbox-Profil), [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) (Egress-Anzeige), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (Startscript der Egress-Restriktion), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Tabelle der Warnungen zu `homebrew` und `claude-code`); `architecture.md` — keine.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Nutzer trägt Homebrew-Formeln (auch `<user>/<tap>/<formel>`) als Liste in `u-boot.yaml` ein, und der erzeugte Devcontainer installiert sie, ohne dass `u-boot` das Werkzeug kennen muss.

Ausgangslage (gelesen): Der Katalog steht im Code (`internal/hexagon/application/devcontainer_features.go`, `ociFeatureCatalogue`: Eintrag `homebrew` mit `source`, `defaultVersion`, `publisher`). Das Schema eines Feature-Eintrags kennt nur `enabled`, `source`, `version` (`ubootYAMLDevcontainerFeature` in `internal/hexagon/application/initproject.go`). Es gibt keinen Weg, einem Feature Optionen oder dem Container Befehle mitzugeben.

Beispiel der Ziel-Form:

```yaml
devcontainer:
  profile: sandbox
  sandbox:
    egress:
      enabled: true
      allow: [github.com, release-assets.githubusercontent.com]
  features:
    homebrew:
      enabled: true
      packages:
        - ripgrep
        - pt9912/ai-harness/ai-harness
```

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Versions-Pins je Paket** (`ripgrep@14`) — Homebrew pinnt Formeln nicht verlässlich; ein Pin-Versprechen wäre eine Zusage, die `u-boot` nicht halten kann. Die Liste ist ausdrücklich ungepinnt und wird so ausgewiesen.
- **Automatische Egress-Einträge aus dem Paketnamen** — die Hosts einer Formel (Download-URL, Abhängigkeiten) sind aus dem Namen nicht ableitbar; Egress-Einträge bleiben bewusste Nutzerentscheidung (Bestand, `devcontainer.sandbox.egress.allow`).
- **Casks, `brew services`, `Brewfile`-Import** — andere Vorgänge mit eigenem Sicherheits- und Zustandsmodell; erst nach Nutzung der Formel-Liste zu beurteilen.
- **Allgemeines `postCreateCommand` in `u-boot.yaml`** — wäre die größere Erweiterung (beliebige Befehle); hier nur der enge Fall Formel-Liste. Falls der Bedarf bleibt: eigener Slice.
- **Änderung am Feature `homebrew` selbst** (Version `2`, `meaningful-ooo`) — Drittanbieter-Feature, nicht Gegenstand von `u-boot`.

## 2. Definition of Done

- [ ] **Schema und Erzeugung:** `devcontainer.features.homebrew.packages` (Liste von Strings) wird gelesen, validiert (Zeichensatz, kein Leerzeichen/Shell-Zeichen, `homebrew` muss `enabled` sein, sonst Exit-Code `10`) und erzeugt einen Installationsschritt im Devcontainer; `config set`/`config get` und `doctor`-Drift kennen das Feld. Test referenziert [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features).
- [ ] **Ausweisung und Quelle:** Die Befehlsausgabe weist die Liste als ungepinnt aus; ein Eintrag der Form `<user>/<tap>/<formel>` wird als Drittanbieter-Quelle ausgewiesen, und `brew trust <user>/<tap>` läuft vor dem Install-Schritt (Entscheidung zur Freigabe: §6, Risiko 2). Test referenziert [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte).
- [ ] **Egress-Nachweis:** Eine Messung (Container bauen, `docker`-Build-Tag-Test) belegt, welche Hosts `ripgrep` und eine Tap-Formel hinter der Egress-Restriktion brauchen; das Ergebnis steht als Messprotokoll unter `docs/messungen/` und im Handbuch als Beispiel, nicht als automatische Default-Allowlist. Test referenziert [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Doku nachgezogen: Lastenheft/Spezifikation (Zusage und Grenzen), Benutzerhandbuch (ohne Kennungen, `AGENTS.md` §5 Regel 6), CHANGELOG.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Die Tabelle nennt nur Dateien, die ich gelesen habe; der Implementer erweitert sie in seinem ersten Lauf und gleicht sie vor Handoff gegen `git diff --stat` ab (`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/initproject.go` (`ubootYAMLDevcontainerFeature`) | update | Feld `Packages []string \`yaml:"packages,omitempty"\`` neben `Enabled`/`Source`/`Version` |
| `internal/hexagon/application/devcontainer_features.go` (`validateDevcontainerFeatures`, `thirdPartyFeatureWarnings`, `projectFeatureEntry`) | update | Validierung (Zeichensatz, nur mit aktivem `homebrew`), Ausweisung ungepinnt/Tap |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` (Felder `PostCreate`, `PostStart`) und `sandbox-init.sh.tmpl` | update | Installationsschritt; `devcontainer.json.tmpl` trägt `postCreateCommand` und `postStartCommand` schon, was `PostCreate` heute enthält, ist nicht gelesen (§6 Risiko 1). `Dockerfile.tmpl` kommt nicht in Frage, siehe Risiko 1 |
| `internal/hexagon/domain/configpath.go` und `internal/hexagon/application/config.go` | update | neue Pfad-Art neben `ConfigDevcontainerFeatureEnabled`/`Source`/`Version`/`SourcesAllow`; Verarbeitung in `application/config.go` (Listenform wie bei `SourcesAllow`; die Parser-Stelle ist nicht gelesen) |
| `internal/hexagon/application/doctor.go` | update (zu prüfen) | Drift-Prüfung der Features (`driftJSONFeatureKeys` steht in `devcontainer_features.go`); ob die Liste dort sichtbar wird, ist nicht gelesen |
| `internal/hexagon/application/devcontainer_features_test.go` (+ neue Testdatei) | update / neu | Happy/Boundary/Negative nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); Build-Tag-`docker`-Test für die Messung |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update | Zusage mit Grenzen; jede Bedingung bekommt einen pinnenden Test oder entfällt (`AGENTS.md` §3.10) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Beispiel-Block, Egress-Hinweis, ungepinnt |
| `docs/messungen/slice-v1-devcontainer-homebrew-pakete/` und `tools/messung/…` | neu | Messprotokoll und Skript getrennt nach der Ablage-Regel (`harness/conventions.md`) |

Ansatz: Falls `u-boot` für die neue Zusage eine ADR braucht (Install-Zeitpunkt, Tap-Freigabe), legt der Architekt sie vor dem Code an (neue ADR, Index aktualisieren, `AGENTS.md` §3.5).

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber hat die zwei Entscheidungen aus §6 (Install-Zeitpunkt, Tap-Freigabe) getroffen.

**Rückführungen:**

- `in-progress` → `next` (zu groß): wenn Tap-Freigabe und Formel-Liste getrennte Zusagen brauchen — dann Tap-Unterstützung als Folge-Slice abspalten.
- `in-progress` → `open` (blockiert): wenn die Messung zeigt, dass kein Install-Zeitpunkt hinter der Egress-Restriktion funktioniert.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün im Devcontainer, Review-Report und Messprotokoll liegen vor, Closure-Notiz mit Lerneintrag geschrieben, jedes Risiko aus §6 hat einen Ausgang.

## 6. Risiken und offene Punkte

1. **Install-Zeitpunkt (Entscheidung des Projektinhabers).** Dev-Container-Features werden nach dem Dockerfile-Build in das Image gelegt; ein `RUN brew install` im `.devcontainer/Dockerfile` hätte daher kein `brew` — das ist eine Annahme aus der Funktionsweise der Dev-Container-Spezifikation, **nicht im Repo gelesen und zu messen**. Kandidat ist ein `postCreateCommand`. Offen: Läuft er vor oder nach dem Egress-Startscript? Kollidiert er mit einem vorhandenen `postStartCommand`/`postCreateCommand` der Sandbox-Erzeugung (Handbuch erwähnt beide, Code nicht gelesen)? Ein `postCreateCommand` läuft nach Spezifikation (§Lebenszyklus-Befehle) nicht erneut gegen einen laufenden Container. — **Ausgang:** offen bis Messung.
2. **Tap als externe Quelle (Entscheidung des Projektinhabers).** [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nimmt „Pakete über Paketmanager im Image-Build" aus, nennt aber Skripte „nicht ausdrücklich in `u-boot.yaml` freigegeben". Ein Tap ist Drittanbieter-Code mit eigener Vertrauensfrage. Optionen: (a) Eintrag in `packages` gilt als Freigabe, wird nur ausgewiesen; (b) eigene Freigabeliste analog `featureSources.allow`. Empfehlung: (b), weil die Linie des Repos „`--yes` ersetzt die Freigabe bewusst nicht" ist. — **Ausgang:** offen bis Entscheidung.
3. **`brew trust` je Tap** (Homebrew 5, laut Benutzerhandbuch im Abschnitt Homebrew) — nur in der Doku gelesen, Verhalten im Container ungemessen. — **Ausgang:** offen bis Messung.
4. **Ungepinnte Pakete widersprechen dem Geist von [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)** (Pin-Politik gilt nur für die gepinnte Werkzeugkette). — **Ausgang:** offen; Abgrenzung in der Zusage festhalten (Ausweisung statt Pin).
5. **Egress-Hosts nicht vorhersagbar** (Formel-Downloads, Bottles auf `ghcr.io`). — **Ausgang:** entfallen als Zusage-Risiko, sobald die Doku nur die gemessene Beispielliste nennt und keine Vollständigkeit behauptet.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `internal/hexagon/application` (APP), `internal/hexagon/domain` (DOM), `internal/adapter/driving/cli` (CLI), `spec/`, `docs/plan/` (DOK) und `tools/`/`docs/messungen/` (WZG); alle stehen in `harness/conventions.md` §Modus-Deklaration und sind nicht zu grob.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** Alle berührten Sub-Areas sind Greenfield (Spec-geführt: Zusage zuerst in Lastenheft/Spezifikation, dann Code). Kein BF-/Hybrid-Block nötig.
