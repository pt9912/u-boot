# Review-Report: slice-v1-release-cut-v0.18.0 — 2026-11-02

**Review-Art:** Code (Release-Cut: Cut-Commit gegen Plan, Hard Rules und
[`releasing.md`](../maintainer/releasing.md); Lead des Blocks `[0.18.0]` gegen
die drei Lastenheft-Stände und die ausgelieferten Texte)

**Gegenstand:** [`slice-v1-release-cut-v0.18.0`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md),
Tag `v0.18.0` (= HEAD `e6c87379`, gepusht). Der Cut selbst ist der Diff des
Commits `e6c87379` gegen `dcc3697c` (7 Dateien); die Träger-Slices
[`slice-v1-devcontainer-host-git-uebernehmen`](../plan/planning/done/slice-v1-devcontainer-host-git-uebernehmen.md)
(Lastenheft `0.28.0`),
[`slice-v1-devcontainer-ressourcenlimits`](../plan/planning/done/slice-v1-devcontainer-ressourcenlimits.md)
(`0.27.0`) und
[`slice-v1-devcontainer-container-name`](../plan/planning/done/slice-v1-devcontainer-container-name.md)
(`0.29.0`) sowie [`slice-v1-hilfetexte-ohne-kennungen`](../plan/planning/done/slice-v1-hilfetexte-ohne-kennungen.md)
(Lastenheft-Beleg der Hygiene-Einträge) tragen je ein eigenes Review und sind
hier nur am Produkttext und als Eintrags-Belege gelesen; hinzu kamen der
MR-003-Nachzug (`827601be`) und die Move-/Marker-Commits (`b160d4ac`,
`77e88143`, `44c6b644`, `dcc3697c`).

**Skill:** `.harness/skills/reviewer.md` @ `v0.18.0`
**Modell:** glm-5.3-flash · **Datum:** 2026-11-02

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-release-cut-v0.18.0`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md),
  Neufassung [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)
  (Release-Cuts wellenlos, 2026-11-02) und [`releasing.md`](../maintainer/releasing.md) §§3–6
- Lastenheft, Stand
  [`0.27.0`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration),
  [`0.28.0`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer),
  [`0.29.0`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
  (Volltexte der Absätze §§495–511, §835 und die Historie-Zeilen 1613–1616)
- [`CHANGELOG.md`](../../CHANGELOG.md) Block `[0.18.0]` samt der drei
  Träger-Einträge, Benutzerhandbuch (Abschnitte *Der Name des Containers*,
  *Auflösung der Token-Quelle*, *Identität*, *Ressourcen*, *Instanzen*, *shell*)
- [`ADR-0040`](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md)
  (Transport der Werte), vorheriger Cut
  [`2026-11-01-slice-v1-release-cut-v0.17.0.md`](2026-11-01-slice-v1-release-cut-v0.17.0.md)
  (Negativbefund-Klassen als Suchbild),
  Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
  (Risiko 1 des Plans)
- `AGENTS.md` §§3.7, 3.9, 3.10, 5 Regel 6; [`AGENTS.md`](../../AGENTS.md) §3.10
  Abgleich Plan-Tabelle gegen `git show --stat e6c87379`

**Nicht ausgeführte Sensors** (Auftrag Review-only — keine Gate-Läufe, keine
Aktion mit Außenwirkung): `make gates`, `make ci`, `make fullbuild` (beim
Hauptagenten grün auf dem Cut-Stand, Cut-Commit-Message; CI auf dem Tip grün
`41403124364`/`41403124438` — laut Auftragskontext, nicht eigenständig
nachgelesen), `make test-docker` (Auftrag), `make verify-depguard` (keine
Importregel berührt), `make doc-immutable RANGE="v0.17.0..HEAD"` (grün, 96
Commits, 0 Befunde — laut Auftragskontext), `make docs-check` (Reviewer
schreibt seinen Report vor dem Gate-Lauf des Verifiers; Link-Pfade manuell
gegen die Muster des v0.17.0-Reports geprüft). Verifikation der DoD und der
CI-Evidence: Verifier, nicht dieser Lauf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Lead-Satz „…; `down`, `shell`, `doctor`, `init` and `generate` work as before.“ widerspricht der eigenen Lieferung in jeder Lesart. In der Ausgabe-Lesart: dieselbe Release-Fassung trägt den Bullet „Comments in the files that `init` and `generate devcontainer` write no longer carry requirement identifiers“ — die geschriebenen Dateien beider Befehle haben sich in diesem Release geändert (Vorlagen-Diff `v0.17.0..HEAD`, 7 Dateien). In der Exit-/Verhalten-Lesart: die Shell setzt den Namen zusammen ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): „Für den Start gelten die Fehlerklassen von `up`“), nennt ihn auf der Standardfehlerausgabe vor der Sitzung — die Startzeile ist neu in diesem Release (`devcontainer_shell_test.go:455`, kam erst mit dem Container-Name-Slice; die `v0.17.0`-Fassung trug sie nicht) — und endet bei Namenskonflikt mit Exit-Code `10` (Code `devcontainer_start.go:525` ordnet `ErrDevcontainerContainerNameConflict`). | [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Plan-Risiko 1) | `CHANGELOG.md` · “`down`, `shell`, `doctor`, `init` and `generate` work as before.” (Zeile 17) | nein — kein Gate liest Lead gegen Einträge oder Code; `make docs-check` prüft nur Links | Lead widerspricht der Lieferung (ausgelieferte Stellen gegeneinander gelesen) |
| F-2 | MEDIUM | Die Not-promised-Liste des Lead ist kürzer als die Vertragsliste: es fehlt „the behaviour of a host `gh` with several accounts“ aus dem Not-promised-Teil des [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)-Eintrags (0.28.0). Die übrigen Grenzen (echter Push, andere Hosts/Engines, Forge-Token, `git credential fill`, zu kleines Limit, Summe über Instanzen, verschachteltes Docker) stimmen mit den drei Einträgen überein. Die Liste ist nicht als Auswahl gekennzeichnet. Die Klasse wiederkehrt aus dem v0.17.0-Cut (dort F-2, LOW) — nach der Skill-Tabelle („Ein LOW-Muster wiederholt sich und wird zum Drift-Signal“) ist der zweite Befund MEDIUM. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), `AGENTS.md` §3.10, Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | `CHANGELOG.md` · “**Not promised:** a real push with a real token, other hosts and engines, …” (Zeile 21) | nein — Review-only | Lead-Grenzliste kürzer als Vertragsliste (2. Befund, Drift-Signal) |
| F-3 | LOW | Die Lead-Formulierung „the value travels by name on the command line and in the environment of the `docker` client process“ ist lesbar als Wert auf der Kommandozeile. Der Vertrag sagt das Gegenteil: der Wert steht in keiner Kommandozeile ([`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)); der Eintrag löst die Mehrdeutigkeit („by name on the command line, the value in the environment …“), der Lead behält sie. Die stärkste Aussage der Zusage (Transport des Tokens) verschwimmt im gelesensten Text des Release. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0040`](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) (Transport), Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | `CHANGELOG.md` · “the value travels by name on the command line” (Zeile 17) | nein — Review-only | Lead-Formulierung gegen die Vertragsaussage lesbar |
| F-4 | LOW | Slice §7 trägt „Noch nicht geschlossen; der Slice liegt in `open/` (Plan).“ Der Slice liegt seit `77e88143`/`44c6b644` in `in-progress/`; der Marker ist abgelaufen. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt) | [`slice-v1-release-cut-v0.18.0.md`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md) §7 · “der Slice liegt in `open/` (Plan)” | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut |
| F-5 | LOW | Die Roadmap trägt „**Letzte Änderung:** 2026-11-01“, wurde aber zuletzt am 2026-11-02 geändert (`44c6b644`: Freigabe-Zeile und Ruhe-Marker-Berichtigung) — der Stand-Wert liegt einen Tag zurück. | [`AGENTS.md`](../../AGENTS.md) §3.7 (Zustandsfeld nennt Zustand und Beleg), §3.10 | `docs/plan/planning/in-progress/roadmap.md` · “**Letzte Änderung:** 2026-11-01” | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut |
| F-6 | LOW | Der DoD-Punkt „`releasing.md`-Zeiger ‚zuletzt‘“ nennt Arbeit, die es nicht mehr gibt: der Zeiger auf das letzte Cut-Slice wurde im Vorzug (`8cd0976b`, Plan-Tabelle-Zeile „im Cut kein Diff“) gestrichen. Der DoD-Eintrag blieb beim alten Stand stehen und wird beim Abhaken gegen ein nicht existierendes Ziel geprüft. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Aussage über Code/Doku) | [`slice-v1-release-cut-v0.18.0.md`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md) §DoD · “`releasing.md`-Zeiger „zuletzt“” | nein — Review-only | Plan-Aussage über Code veraltet (`BEO-DOK/plan-aussage-ueber-code`) |
| F-7 | LOW | `releasing.md` trägt im Kopf des Abschnitts 3 die neue Setzung „seine **Form** kommt aus der Skelett-Vorlage …, der **Ablauf** steht vollständig in diesem Dokument“, aber Punkt 5 derselben Liste sagt weiterhin „Release-Cut-Slice nach dem Muster der früheren Cuts anlegen“ — der Rest des vor dem Cut gezogenen Schnitts. Ein Leser, der Punkt 5 folgt, greift zum letzten Cut-Slice statt zur Vorlage, genau das, was die Korrektur abstellen wollte. | [`releasing.md`](../maintainer/releasing.md) §3, [`AGENTS.md`](../../AGENTS.md) §3.10 | `docs/maintainer/releasing.md` · “nach dem Muster der früheren Cuts anlegen” (Zeile 78) | ja — `make docs-check` prüft Links, nicht die Aussage; sonst Review-only | Doku-Stellen gegeneinander (Residuum einer Berichtigung) |
| F-8 | LOW | Die Roadmap führt unter *Nächste Wellen* die Zeile „`v0.18.0` — Feature-Katalog: Quellen, Optionen und Basis-Image (geplant)“ mit Trigger „`v0.17.0` geschlossen …“ — der Tag `v0.18.0` ist jetzt gesetzt und trägt anderen Inhalt; die geplante Version ist vergeben. | [`AGENTS.md`](../../AGENTS.md) §3.10, `releasing.md` §4 („unter ‚Offene Wellen‘ führen, bis der Tag gesetzt ist“) | `docs/plan/planning/in-progress/roadmap.md` · “`v0.18.0` — Feature-Katalog: Quellen, Optionen und Basis-Image (geplant)” (Zeile 32) | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut (Vorschlags-Version vergeben) |
| F-9 | LOW | Die zwei CR-Entwürfe tragen den abgelaufenen Vorschlag `v0.18.0`: JDK-Slice-Feld „Vorschlag Release `v0.18.0`“, CR-Entwurf „Vorschlag `v0.18.0`“ und in §5 „Release-Zuordnung (`v0.18.0` Katalog-Welle …)“, dazu die Bestand-Zeilen in `open/README.md` („ohne Welle (Vorschlag `v0.18.0`)“ je Stelle). Mit dem Tag ist der Vorschlagswert belegt; die nächsten Zuordnung ist `v0.19.0`. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker) | [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/open/slice-v1-katalog-gepinnte-toolchain-entwurf.md) · „Vorschlag `v0.18.0`“; gleicher Stand in [`slice-v1-katalog-feature-jdk-gepinnt.md`](../plan/planning/open/slice-v1-katalog-feature-jdk-gepinnt.md) und `docs/plan/planning/open/README.md` (Zeilen 86–87) | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut (Vorschlags-Version vergeben) |
| F-10 | LOW | Die Kopfliste des Slice „Planungs-Commits ohne Release-Wirkung“ nennt nur den CR-Entwurf und die Messungs-Commits des Container-Name-Slices. Im Range liegen daneben: die Messungs-Commits des Host-git-Slices (30 Dateien unter `tools/messung/` und `docs/messungen/`) und zwei `tools/mk`-Commits (`82682508`/`7bbb8185`: `make dev-up`/`dev-shell` GIT_TOKEN-/GH_TOKEN-Bridge und `INSTANCE`-Weiterleitung in `harness/README.md`/`tools/mk/devcontainer.mk`) — keine Release-Beiträge, aber auch nicht als solche geführt; der Abgleich „was steckt in `v0.18.0`“ gegen `git diff --stat` findet sie unerklärt. Die Klasse wiederkehrt aus dem v0.17.0-Cut (dort F-6). | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`slice-v1-release-cut-v0.18.0.md`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md) Kopf · “Planungs-Commits ohne Release-Wirkung” | nein — Review-only | Inhaltsliste des Release unvollständig (2. Befund) |
| F-11 | INFO | Handbuch-Version `2.1` unverändert, obwohl drei Träger-Slices Handbuchinhalt brachten: [`releasing.md`](../maintainer/releasing.md) §§3 und 5 verlangt nur den Abgleich von `Stand` gegen das Tag-Datum (2026-11-02, stimmt) und keine Anhebung der Handbuch-Version; die Inhaltsänderungen kamen je im Träger-Review. Keine Regel verletzt. | `releasing.md` §5 | `docs/user/benutzerhandbuch.md` · “Handbuch-Version: 2.1” | nein — Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `CHANGELOG.md` Kopf `[0.18.0]`: Datum 2026-11-02 = Tag-Datum (Tag liegt auf `e6c87379`, 2026-11-02), Zählung „Eighteenth release.“ (v0.1.0–v0.17.0 = 17, fortlaufend; `0.17.0` „Seventeenth“) | geprüft, ohne Befund |
| Lead, Promised-Teil gegen die drei Stände: „the three features on a Linux host (amd64) with a local engine, as far as the tests … show“, Grenze „only on an engine that is not nested“, Testbasis (Unit, fake sources, Docker-Tests mit Build-Tag `docker`) — deckungsgleich mit 0.27.0/0.28.0/0.29.0 | geprüft, ohne Befund (Grenzliste: F-2, Formulierung: F-3, Shell-Satz: F-1) |
| Drei `[0.18.0]`-Einträge gegen die Lastenheft-Stände: Host-Token-Quelle (Reihenfolge Host-Variable → `gh auth token`, Quelle genannt nie Wert, Warnung ohne Abbruch), Identität (`user.name`/`user.email`, bestehender Wert bleibt, beide Profile), Ressourcen-Schlüssel (`--memory`/`--cpus` in `runArgs` beider Profile, invalid → Exit `10`), Containername (Namensform, Suche über Labels unverändert, Halter nicht entfernt/umbenannt, Zufallsname bis Neuanlegen, Startzeilen von `up`/`shell`, Exit `10` vor dem Hilfs-Container) | geprüft, ohne Befund |
| Hygiene-Einträge mit Beleg im Range: Hilfetexte/Meldungen/Vorlagen (`fcbafc2e`, [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)); Vorlagen-Diff bestätigt; keine `LH-`-Treffer mehr in `internal/hexagon/application/templates/` | geprüft, ohne Befund |
| `CHANGELOG.md` Fuß: `[Unreleased]: …/compare/v0.18.0...HEAD`, `[0.18.0]: …/compare/v0.17.0...v0.18.0`, `[0.17.0]` unverändert; `[Unreleased]` leer | geprüft, ohne Befund |
| Versionszeilen `cmd/uboot/main.go:43`, `repo.mk:61`, `Dockerfile:45` samt Kommentar-Nennungen (`Dockerfile:41,43`, `repo.mk:57,60`); kein weiterer `0.17.0`-Treffer in Code, READMEs, Handbuch, `repo.mk`, `Dockerfile`, `.github/` | geprüft, ohne Befund; Reste nur in CHANGELOG-/Lastenheft-Historie, Test-Fixtures (`help_identifiers_test.go`, eigene Version-Eingabe) und Bestandsaufnahmen ([`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)) |
| `README.md`/`README.de.md`: Status-Zeile und `VERSION=`-Beispiel; keine Anforderungs-/ADR-/Spec-Kennung, kein Planverweis, kein „released“ (`AGENTS.md` §5 Regel 6; verlinktes Lastenheft ohne Anker zulässig) | geprüft, ohne Befund |
| Handbuch: Kopf `Software-Version` v0.18.0 und `Stand` 2026-11-02 = Tag-Datum, drei Versionsbeispiele (`0.18.0`, `download/v0.18.0`, `VERSION=0.18.0`), Gültigkeitsbereich; Kennungen-Anzahl 0 (`grep`-Abgleich) | geprüft, ohne Befund (Handbuch-Version: F-11) |
| Handbuch-Abschnitte der Träger gegeneinander und gegen den Lead gelesen (*Der Name des Containers*, *Zugangsdaten*, *Auflösung der Token-Quelle*, *Identität*, Ressourcen-Hinweis, *Instanzen*, *shell*-Abschnitt): Transport („unter dem Namen auf die Kommandozeile, der Wert in die Umgebung“), Namenszusammensetzung, Exit-Klassen konsistent | geprüft, ohne Befund (F-1 betreffen nur Lead gegen Einträge/Code, nicht die Handbuchstellen) |
| Plan-Tabelle (§3) gegen `git show --stat e6c87379`: 7 Dateien (`CHANGELOG.md`, `Dockerfile`, beide READMEs, `main.go`, Handbuch, `repo.mk`), jede mit Plan-Zeile; „kein Diff“ für `releasing.md`, `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` im Cut-Commit bestätigt | geprüft, ohne Befund |
| MR-003-Nachzug (`827601be`): [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md) Neufassung (Release-Cuts wellenlos, vier welle-lose Cuts `v0.13.1`/`v0.15.1`/`v0.16.1`/`v0.17.0` als Grundlage), Index-Zeile, Roadmap-Format-Regel, Planning-README, `plan-welle`/`close-welle`/`planner`; Link-Tiefen der Nachzüge geprüft (`../../harness/…`, `../../../harness/…` lösen korrekt) | geprüft, ohne Befund |
| Move-/Marker-Commits `44c6b644`/`dcc3697c`: Ruheort-Regel (Roadmap-Zeile, Bestand `open/README.md`, Verweis-Tiefen der done-/in-progress-Fassungen), gelöste Marker „noch offen: Review/Verification/Closure“ und „(Review läuft)“ | geprüft, ohne Befund |
| `git mv`-Disziplin (§3.3): `b160d4ac` (open→next), `77e88143` (next→in-progress), `2667010d` (Container-Name→done) als reine Moves, Verweis-Nachzug je getrennt | geprüft, ohne Befund |
| §3.7/§3.9 über den Cut-Diff: nur Versionswerte und Kommentar-Nennungen auf den neuen Stand; kein Kommentar über Verworfenenes, kein Verweis auf den alten Wert, kein Helfer zwischen Doc-Kommentar und Funktion | geprüft, ohne Befund |
| Reviewer-HIGH-Liste (Dual-Classifier, Referenzmodell, Sprachvertrag, Architekturgrenze, Gate-Lockerung, Kommentar-Klassen, Carveout, Immutabilität, Lifecycle): der Cut ändert keinen Produktcode, kein Gate, kein ADR, keine Spec; `doc-immutable` über den Range grün (Auftragskontext) | geprüft, ohne Befund |
| Roadmap-Zeile *Offene Wellen* für den Cut: „liegt in `in-progress/` … Freigabe … erteilt 2026-11-02“, drei Träger mit Lastenheft-Ständen verlinkt; `open/README.md` führt den Cut nicht mehr als Bestand | geprüft, ohne Befund (Marker: F-5, Nächste-Wellen-Zeile: F-8) |
| CI-/Gate-Evidence (CI-Läufe des Tips, `gates`/`fullbuild`, `doc-immutable`, Trockenlauf des Tap-Pfads, `tap-check`) | nicht geprüft: Review-only; Nachweis beim Verifier, IDs im Auftragskontext |
| Produktcode der drei Träger (Sentinel-Abdeckung, `devcontainer_shell_test.go`/`devcontainer_start.go` nur fundstellenweise gelesen) | nicht geprüft im Ganzen: die Träger haben je ein Review; hier nur die für F-1 nötigen Stellen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (F-1, F-2) |
| LOW | 8 (F-3 bis F-10) |
| INFO | 1 (F-11) |

**Finding-Klassen dieses Laufs:** Lead widerspricht der Lieferung (ausgelieferte
Stellen gegeneinander gelesen) · Lead-Grenzliste kürzer als Vertragsliste (2.
Befund, Drift-Signal) · Lead-Formulierung gegen die Vertragsaussage lesbar ·
Zustandsmarker nach Ablauf nicht abgebaut (F-4, F-5, F-6-teilweise, F-8, F-9) ·
Plan-Aussage über Code veraltet (`BEO-DOK/plan-aussage-ueber-code`) ·
Doku-Stellen gegeneinander (Residuum einer Berichtigung) · Inhaltsliste des
Release unvollständig (2. Befund)

**Steering-Loop-Vermerk:** Die Klasse „Zustandsmarker nach Ablauf nicht
abgebaut“ erreicht mit diesem Lauf die Schwelle (v0.17.0-Cut F-3/F-4, dieser
Lauf F-4/F-5/F-8/F-9 trotz der Berichtigungsrunde in `827601be`). Sie gehört in
den Lerneintrag der Closure (§7) und zur Prüfung des Architect-Moduls, ob ein
computational Sensor möglich ist (etwa `make docs-check`-Regel gegen
„Vorschlag Release <tag-aktuell>“ in `open/` nach gesetztem Tag). Die Klassen
„Lead-Grenzliste kürzer als Vertragsliste“ und „Inhaltsliste des Release
unvollständig“ stehen bei zwei Vorgängen je.

## Verdikt

**Merge-blockierend:** nein — der Tag `v0.18.0` ist gesetzt und veröffentlicht
und wird nicht zurückgenommen ([`releasing.md`](../maintainer/releasing.md)
*Zweck*). F-1 und F-2 betreffen den ausgelieferten `CHANGELOG.md`-Text und
führen zu Nachzugs-Commits auf `main`; der Slice geht nicht nach `done/`, bevor
der Hauptagent F-1 und F-2 entschieden hat. Die LOW-Findings gehen an den
Implementer (F-4, F-6 in den Slice selbst; F-5, F-8, F-9 in die Planning-Doku;
F-7 in `releasing.md`; F-3 je nach F-1-Runde in den Nachzug). Die Findings
werden nicht herabgestuft, wenn der Implementer widerspricht; HIGH gibt es
nicht, ein Konfliktpfad über den Architect entfällt.

**DoD-Vermerk:** Nicht Bestandteil dieses Laufs (Verifier). Festgehalten: die
offenen DoD-Punkte (Roadmap-Zeile nach *Abgeschlossene Wellen* mit Tag-Hash,
Closure-Notiz, `git mv` nach `done/`) sind im Slice ehrlich als offen
ausgewiesen; der DoD-Punkt zur „`releasing.md`-Zeiger“-Doku ist durch F-6 gegen
den Ist-Stand zu berichtigen, bevor er abgehakt wird.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in
die Closure-Notiz §7 des Slice und von dort in den Zähler. Dieser Report ist
ein Lauf-Beleg und ersetzt keine Verifikation (Verifier, anderer Kontext). Der
Report wurde vor dem Gate-Lauf des Verifiers geschrieben; seine Links sind
manuell gegen die Fundstellen geprüft, nicht durch `make docs-check`.