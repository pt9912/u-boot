# Slice V1: Homebrew-Pakete im Devcontainer — ein Dockerfile-Block installiert Homebrew und eine frei wählbare Formel-Liste beim Image-Bau, ohne dass jedes Werkzeug einen Katalogeintrag braucht

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen (eine Zusage-Erweiterung im
Devcontainer-Bereich; keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist).

**Bezug:** [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (zwei neue Schlüssel in `u-boot.yaml`), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Scope: Katalog-Eintragstyp Dockerfile-Block und Feature-Schema), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Installer-Skript und Tap als Quellen), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Abgrenzung: Pakete sind ungepinnt), [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md) (Bau-Pfad, nur als Rahmen).

**Berührte Spec-Stellen:** [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (beide laut Architect-Bericht zu [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md); die Abschnitte liest der Implementer beim Start), [`LH-FA-DEV-003.a`](../../../../spec/spezifikation.md#lh-fa-dev-003a--mechanik-der-freigabe-externer-feature-quellen) (Freigabe externer Feature-Quellen), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Katalog-Features im Sandbox-Profil), [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) (Egress-Anzeige), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (Startscript der Egress-Restriktion), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Tabelle der Warnungen zu `homebrew` und `claude-code`); `architecture.md` — keine.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers; Fassung 2 auf den Dockerfile-Weg, 2026-11-05). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Nutzer trägt Homebrew-Formeln (auch `<user>/<tap>/<formel>`) als Liste in `u-boot.yaml` ein, und das erzeugte Devcontainer-Dockerfile installiert Homebrew und diese Formeln beim Image-Bau, ohne dass `u-boot` das Werkzeug kennen muss.

Ausgangslage (gelesen):

- Der Katalog steht im Code (`internal/hexagon/application/devcontainer_features.go`). `ociFeatureCatalogue` führt `homebrew` als OCI-Feature (`ghcr.io/meaningful-ooo/devcontainer-features/homebrew`, Version `2`). `blockAndOptionCatalogue` führt Einträge ohne OCI-Quelle mit `dockerfileBlock: true` (`ollama`, `gh`, `glab`, `go-toolchain`, `golangci-lint`, `java-toolchain`).
- Das Dockerfile-Template (`internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl`) geht von `mcr.microsoft.com/devcontainers/base:debian` aus. Die Blöcke laufen unter `USER root` (die Bedingung `if or .UID .Sandbox .Gh .Glab .GoToolchain .GolangciLint .JdkToolchain` setzt es), am Ende steht `USER vscode`.
- Das Schema eines Feature-Eintrags kennt nur `enabled`, `source`, `version` (`ubootYAMLDevcontainerFeature` in `internal/hexagon/application/initproject.go`). Optionen oder Paketlisten gibt es nicht.

**Warum der Dockerfile-Weg:** Ein OCI-Feature wird nach dem Dockerfile-Build in das Image gelegt (Annahme aus der Dev-Container-Spezifikation, im Repo nicht gelesen); ein `RUN brew install` im Dockerfile fände dort noch kein `brew`. Installiert der Block Homebrew selbst, steht `brew` im selben Dockerfile bereit. Die Pakete stecken im Image und überstehen das Neuanlegen des Containers. Die Egress-Restriktion greift nach der Doku über ein Startscript beim Containerstart ([`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)), nicht beim Bau; ob der Bau dadurch ohne Egress-Einträge auskommt, ist zu messen (DoD 3).

Beispiel der Ziel-Form (Schlüssel `brew` ist ein Arbeitsname, siehe §6 Entscheidung 1):

```yaml
devcontainer:
  profile: sandbox
  features:
    brew:
      enabled: true
      packages:
        - ripgrep
        - pt9912/ai-harness/ai-harness
```

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Versions-Pins je Paket** (`ripgrep@14`) — Homebrew pinnt Formeln nicht verlässlich; ein Pin-Versprechen wäre eine Zusage, die `u-boot` nicht halten kann. Die Liste ist ausdrücklich ungepinnt und wird so ausgewiesen.
- **Automatische Egress-Einträge aus dem Paketnamen** — die Hosts einer Formel (Download-URL, Abhängigkeiten) sind aus dem Namen nicht ableitbar; Egress-Einträge für den Betrieb im laufenden Container (`brew install`, `brew update`) bleiben Nutzerentscheidung (Bestand, `devcontainer.sandbox.egress.allow`).
- **Casks, `brew services`, `Brewfile`-Import** — andere Vorgänge mit eigenem Sicherheits- und Zustandsmodell; erst nach Nutzung der Formel-Liste zu beurteilen.
- **Ein allgemeines `postCreateCommand` in `u-boot.yaml`** — wäre die größere Erweiterung (beliebige Befehle) und entfällt hier, weil der Dockerfile-Weg keinen Lebenszyklus-Befehl braucht.
- **Änderung am OCI-Feature `homebrew`** (`meaningful-ooo`) — Drittanbieter-Feature, nicht Gegenstand von `u-boot`; sein Bestand bleibt (§6 Entscheidung 1 klärt das Verhältnis).

## 2. Definition of Done

- [ ] **Schema und Block:** `devcontainer.features.<Schlüssel>.packages` (Liste von Strings) wird gelesen und validiert (Zeichensatz ohne Leerzeichen und Shell-Zeichen; `packages` nur am Schlüssel `brew` ohne `source:`, sonst Exit-Code `10`; Eintrag `<formel>` oder `<user>/<tap>/<formel>` nach [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md)); das erzeugte Dockerfile trägt einen Block, der Homebrew aus einer gepinnten Quelle ohne Installer-Skript einrichtet und die Formeln als Benutzer `vscode` installiert. `config set`/`config get` kennen das Feld (Listenpfad hängt an); die Tap-Freigabe `devcontainer.brewTaps.allow` prüft `generate devcontainer` vor jedem Schreiben und `doctor` im Befund der Feature-Freigabe (Form und Orte nach der ADR). Test referenziert [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features).
- [ ] **Ausweisung und Konflikt:** Die Befehlsausgabe weist den Block als Drittanbieter-Code mit ungepinnter Formel-Liste aus; ein Eintrag der Form `<user>/<tap>/<formel>` wird als eigene Quelle ausgewiesen und wird nur ausgewiesen und installiert, wenn der Tap in der Freigabeliste steht (sonst fachlicher Validierungsfehler, Exit-Code `10`, nichts geschrieben; Entscheidung 2 aus §6); der Schlüssel zusammen mit dem OCI-Feature `homebrew` ist ein fachlicher Validierungsfehler (Exit-Code `10`, nichts geschrieben), wie die Regel für `go` zusammen mit `go-toolchain`. Test referenziert [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte).
- [ ] **Messung im Build:** Ein Test mit Build-Tag `docker` baut das Image mit `ripgrep` und einer Tap-Formel, in beiden Profilen und, wo die Plattform es zulässt, auf `amd64` und `arm64`; das Messprotokoll unter `docs/messungen/` nennt, ob der Bau ohne Egress-Einträge gelingt, welche Hosts hinter der Egress-Restriktion im laufenden Container ein `brew install` braucht, und ob `brew` in Login- und Nicht-Login-Shells im `PATH` steht. Test referenziert [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).

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
| `internal/hexagon/application/devcontainer_features.go` (`blockAndOptionCatalogue`, `validateDevcontainerFeatures`) | update | neuer Katalogeintrag `brew` mit `dockerfileBlock: true`, Validierung (Zeichensatz, Konflikt mit dem OCI-Feature `homebrew`); die Ausweisung kommt aus einer neuen Warnfunktion nach dem Muster `forgeCLIWarnings` (`thirdPartyFeatureWarnings` hängt am `publisher` einer OCI-Quelle und trägt den Block nicht; `projectFeatureEntry` bleibt unverändert, Befund des Architects zu [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md)) |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | neuer Block (Build-Abhängigkeiten per `apt-get`, Homebrew aus gepinnter Quelle, `brew install` als `vscode`, `ENV`-Zeilen für den `PATH`); die Bedingung `if or .UID .Sandbox …` am Kopf bekommt das neue Template-Datum |
| `internal/hexagon/application/generate.go` (Template-Daten) | update (zu lesen) | das neue Template-Datum und die Liste müssen dort in die Daten des Dockerfile-Templates; die Fundstelle ist nicht gelesen |
| `internal/hexagon/domain/configpath.go` und `internal/hexagon/application/config.go` | update | neue Pfad-Art neben `ConfigDevcontainerFeatureEnabled`/`Source`/`Version`/`SourcesAllow`; Listenform wie bei `SourcesAllow` (die Parser-Stelle in `config.go` ist nicht gelesen) |
| `internal/hexagon/application/doctor.go` | update | die Tap-Regel im Befund der Feature-Freigabe (`classifyFeatureEntries`); eine Drift-Prüfung des Blocks entfällt, weil `projectFeatureSourceVersion` für Einträge ohne Quelle nichts liefert und `checkDevcontainerDockerfile` nur `FROM` liest (Befund des Architects) |
| `internal/hexagon/application/devcontainer_features_test.go` (+ neue Testdatei) | update / neu | Happy/Boundary/Negative nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); Build-Tag-`docker`-Test für die Messung |
| `spec/spezifikation.md` | update | Zielbild des Blocks, Schlüsselnamen und Pfade, Verfahren der Tap-Freigabe, Konflikt-Regel mit dem Katalog-Feature `homebrew` (Exit-Code `10`), Plattform-, Bau- und Messgrenzen der Zusage; jede Bedingung bekommt einen pinnenden Test oder entfällt (`AGENTS.md` §3.10). Das Lastenheft ändert nur der CR, nicht dieser Slice |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Beispiel-Block, ungepinnt, Egress im laufenden Container |
| `docs/messungen/slice-v1-devcontainer-homebrew-pakete/` und `tools/messung/…` | neu | Messprotokoll und Skript getrennt nach der Ablage-Regel (`harness/conventions.md`) |

Ansatz: Die ADR liegt vor ([`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), `Proposed`); sie legt Schlüssel `brew`, die Freigabeliste `devcontainer.brewTaps.allow` und die Pin-Stufe (voller Commit von `Homebrew/brew`, keine Übersteuerung, die Vier-Stellen-Politik aus [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) gilt nicht) fest. Das Lastenheft ändert der Projektinhaber im vorgeschalteten CR [`slice-cr-devcontainer-brew-pakete`](../in-progress/slice-cr-devcontainer-brew-pakete.md); weder ADR noch Slice ändern `LH-*`.

## 4. Trigger

**Start** (`next` → `in-progress`): [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) ist `Accepted` (Bestätigung von Schlüsselname, Listenpfad und Pin-Stufe durch den Projektinhaber); [`slice-cr-devcontainer-brew-pakete`](../in-progress/slice-cr-devcontainer-brew-pakete.md) liegt in `done/` (Lastenheft-Wortlaut freigegeben und übernommen); `in-progress/` ist leer (WIP-Limit 1). Die zwei Entscheidungen aus §6 und die Sichtung des Beobachtungs-Registers (§8) sind erfolgt.

**Rückführungen:**

- `in-progress` → `next` (zu groß): wenn Tap-Unterstützung und Formel-Liste getrennte Zusagen brauchen — dann Tap-Unterstützung als Folge-Slice abspalten.
- `in-progress` → `open` (blockiert): wenn die Messung zeigt, dass Homebrew auf einer der zugesagten Plattformen im Block nicht sauber läuft.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün im Devcontainer, Review-Report und Messprotokoll liegen vor, Closure-Notiz mit Lerneintrag geschrieben, jedes Risiko aus §6 hat einen Ausgang.

## 6. Risiken und offene Punkte

**Entscheidungen des Projektinhabers** (getroffen 2026-11-05)

1. **Verhältnis zum OCI-Feature `homebrew`: neuer Schlüssel daneben.** Der Block steht als eigener Katalogschlüssel neben dem bestehenden Feature `homebrew`, dessen Bestand unverändert bleibt; beide zusammen aktiv ist ein Validierungsfehler. Der Name des Schlüssels (hier `brew`) ist ein Arbeitsname, den der Architect vor dem Code festlegt.
2. **Tap als externe Quelle.** [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nimmt Paketmanager-Installationen im Image-Bau aus, nennt aber Skripte „nicht ausdrücklich in `u-boot.yaml` freigegeben“. Ein Tap ist Drittanbieter-Code mit eigener Vertrauensfrage. **Entscheidung: eigene Freigabeliste** analog `featureSources.allow`; `--yes` ersetzt die Freigabe nicht, wie bei externen Features. Name und Ort der Liste legt der Architect vor dem Code fest.

**Risiken**

3. **Installer-Skript.** Der übliche Homebrew-Installer ist ein Skript aus dem Netz; [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nennt ad hoc geladene Shell-Skripte ausdrücklich. Vorschlag: `git clone` von `Homebrew/brew` auf eine gepinnte Revision, ohne Skript (die Pin-Stufe — Tag oder Commit — und ihre Prüfbarkeit sind offen; eine Prüfsumme wie bei den Release-Archiven gibt es nicht). Ob dieser Weg ohne den Installer vollständig läuft, ist nicht gelesen und zu messen. — **Ausgang:** offen bis Messung.
4. **Rechte und Build-Abhängigkeiten.** Homebrew verweigert Root; der Block müsste `brew` als `vscode` ausführen, obwohl die Blöcke unter `USER root` stehen. Welche Debian-Pakete zusätzlich nötig sind (Werkzeuge zum Bauen und zum Laden), ist nicht gelesen. — **Ausgang:** offen bis Messung.
5. **Architektur.** Ob Homebrew auf Linux-`arm64` im Block sauber läuft, ist nicht gelesen. Die Zusage gilt nur für die gemessenen Plattformen (DoD 3), Linux-`amd64` mit lokaler Engine als Mindestmaß nach dem Muster der übrigen Blöcke. — **Ausgang:** offen bis Messung.
6. **Image-Größe und Bauzeit** wachsen mit Homebrew und den Formeln. — **Ausgang:** offen; das Messprotokoll nennt beide Werte, `u-boot` begrenzt sie nicht.
7. **Ungepinnte Pakete widersprechen dem Geist von [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)** (Pin-Politik gilt nur für die gepinnte Werkzeugkette). — **Ausgang:** offen; Abgrenzung in der Zusage festhalten (Ausweisung statt Pin).
8. **Egress im laufenden Container** (`brew install`, `brew update`, Bottles auf `ghcr.io`) ist aus dem Paketnamen nicht vorhersagbar. — **Ausgang:** entfallen als Zusage-Risiko, sobald die Doku nur die gemessene Beispielliste nennt und keine Vollständigkeit behauptet.

9. **Kombination mit den übrigen Blöcken im selben Dockerfile** (`ollama`, `gh`, `glab`, Werkzeugketten, `nestedRuntime`) wird nicht gebaut; der Golden Case pinnt nur die Reihenfolge. — **Ausgang:** eingetreten: weiterer Beleg in `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` bei der Closure (der Trigger der Beobachtung); kein Lauf in dieser DoD.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `internal/hexagon/application` (APP), `internal/hexagon/domain` (DOM), `spec/` und `docs/plan/` (DOK) sowie `tools/` und `docs/messungen/` (WZG); alle stehen in `harness/conventions.md` §Modus-Deklaration und sind nicht zu grob. Das Dockerfile-Template liegt unter `internal/hexagon/application/` und gehört zu APP.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stichworte Dockerfile-Block, Katalog, Feature, Validierung); vier Treffer, alle offen:

- [`BEO-APP/block-schritt-ohne-golden-pin`](../observations/BEO-APP/block-schritt-ohne-golden-pin/observation.md) (1 Beleg): der Golden Case des neuen Blocks pinnt auch, was im Block fehlen muss (kein Installer-Skript, kein `curl | sh`, keine Anmeldung); Auswirkung auf DoD 1 und die Testzeile in §3.
- [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) (2 Belege): jede neue Regel für `packages` (Zeichensatz, Konflikt mit dem OCI-Feature `homebrew`, Freigabe des Taps) gilt auch in der Schreibpipeline von `config set`, nicht nur beim Lesen; Auswirkung auf DoD 1.
- [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md) (5 Evidence-Dateien; der Ausgang ist laut `state.md` nicht zugewiesen und liegt beim Projektinhaber): der Trigger „weiterer Katalogeintrag mit Dockerfile-Block“ tritt mit diesem Slice ein; der Golden Case pinnt die Reihenfolge der Blöcke, ein Lauf der Kombination mit `ollama`, `gh`, `glab` und `nestedRuntime` ist nicht Teil der DoD (siehe Risiko 9).
- [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) (4 Evidence-Dateien): ein gepinnter Stand altert ohne rotes Gate; für die Homebrew-Revision gilt dasselbe (Risiko 3).

**Modus-Begründungsblock:** Alle berührten Sub-Areas sind Greenfield (Spec-geführt: Zusage zuerst in Lastenheft/Spezifikation, dann Code). Kein BF-/Hybrid-Block nötig.
