# Slice V1: Hilfetexte und Meldungen des Binarys ohne Anforderungs-, Spezifikations- und Slice-Kennungen — der Benutzer liest Worte statt `LH-FA-CLI-006`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung offen (Frage 2 in §6). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er ist der Ausgang der Beobachtung [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md): Die Produktfrage dort („Kennungen aus Ausgabe und Hilfe entfernen oder als Produktmerkmal behalten“) hat der Projektinhaber am 2026-11-01 entschieden: entfernen („damit kann kein Benutzer etwas anfangen“).

**Bezug:** [`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache) (CLI-Ausgaben und Fehlermeldungen Englisch), [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) (Hilfeausgabe), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes; der Inhalt der Hilfe-Zeile „Exit codes“ bleibt, nur die Kennung geht), [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6 (Nutzer-Doku trägt keine Kennungen; hier sinngemäß für die Ausgabe des Binarys), [`AGENTS.md`](../../../../AGENTS.md) §3.7 (Kommentare im Code dürfen Kennungen tragen).

**Berührte Spec-Stellen:** `spec/spezifikation.md` Zeilen 887 und 916 (gelesen am 2026-11-01): Beide legen die Textform der Warnzeile fest, `Warning (<Code>): <Text>`. Streicht der Slice den Code aus der Warnzeile (Frage 1 in §6), ändern sich diese zwei Stellen mit; sonst bleibt die Spezifikation unverändert. Der Verweis zeigt aufwärts, die Spezifikation nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Implementer, Entwurf im Auftrag des Projektinhabers). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Hilfetexte aller Befehle (`Short`, `Long`, Beispiele, Usage-Text jeder Option) **und** die für den Benutzer sichtbaren Meldungen des Binarys (Fehlertexte, Warnungen, Hinweise, Diagnose-Texte von `doctor`) tragen keine Anforderungs-, Spezifikations-, ADR-, Slice- oder `§`-Verweise und keine deutschen Wörter. Sie sagen in Worten, was die Kennung nur hieß. Der Inhalt der Zusagen bleibt; es kommt keine neue Zusage hinzu. Ein Wächter-Test pro Fläche hält den Stand (Hilfebaum, Meldungs-Literale).

**Begründung (Projektinhaber, 2026-11-01):** „in den help messages werden spec Kennungen referenziert — damit kann kein Benutzer etwas anfangen“, und dasselbe bei den Meldungen. Die Nutzer-Doku trägt nach [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6 keine Kennungen; für Hilfe und Meldungen des Binarys gab es keine Regel (Beobachtung [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md), vier Belege, Schwelle erreicht).

### Gefundener Stand (gelesen am 2026-11-01 am Binary `u-boot version 0.17.0` mit `u-boot <befehl> --help`, rekursiv über alle 26 Befehle, und im Quelltext; Zeilenangaben nennen den Stand vor dem Diff)

Befehlsbaum (gelesen): `add`, `completion` (`bash`, `fish`, `powershell`, `zsh`), `config` (`get`, `list`, `set`), `devcontainer` (`down`, `shell`, `up`), `doctor`, `down`, `egress` (`list`), `generate`, `help`, `init`, `logs`, `remove`, `template` (`list`), `up`. `completion`, `help`, `egress` (Wurzel), `devcontainer` (Wurzel), `devcontainer shell` und die Hilfe des Wurzelbefehls tragen **keine** Kennung in der Ausgabe von `--help` (Suche nach `LH-`, `SPEC-`, `ADR-`, `ARC-`, `slice`, `§` in der gesammelten Ausgabe: Treffer nur in den unten genannten Befehlen).

**A. Hilfetexte** (Datei `internal/adapter/driving/cli/<befehl>.go`; Zahl der Zeilen mit Treffer, ohne Zeilenkommentare und ohne das Feld `Code:`; vorläufig, der Implementer zählt vor dem Diff neu):

| Befehl | Zeilen | Beispiele (Fundstelle) |
|---|---|---|
| `add` | 9 | `add.go:59` „supports only postgres (LH-FA-ADD-002)“, `:69` „Flag combinations (LH-FA-CLI-007/008):“, `:74` und `:66` „voll-schema“ (deutsches Wort in englischer Hilfe), `:98`, `:100`, `:102` Option-Texte mit Kennung | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `config` (Wurzel, `set`) | 10 | `config.go:168` „Writable paths (LH-FA-CONF-001 / §D1 of slice-m8-config.md):“ (Verweis auf eine Slice-Datei und einen Paragraphen), `:178` „(slice §D3)“, `:183` „LH-FA-CLI-007/008 envelope“, `:185` „Exit codes (LH-FA-CLI-006):“, `:274`, `:276`, `:278` Option-Texte; „fachlich“ (deutsch) im Text der Exit-Codes von `config` und `egress list`, `remove.go` (Zeile im Exit-Code-Block) | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `doctor` | 4 | `doctor.go:67` „per LH-FA-DIAG-003“, `:70` „Default checks (LH-FA-DIAG-002):“, `:92` | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `down` | 5 | `down.go:82`, `:84` „Destructive confirmation gate (LH-FA-CLI-005A):“, `:95` „LH-FA-CLI-006 exit codes:“, `:101` „LH-NFA-PERF-002 progress:“, `:120` | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `generate` | 12 | `generate.go:59` bis `:62`, `:65`, `:76`, `:107`, `:109`, `:111`, `:113` |
| `init` | 18 | `init.go:112` bis `:124`, `:146` „# LH-AK-005 devcontainer flow“, `:165` bis `:183`; `:179` nennt außerdem „slice-v1-template-init T4“ | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `logs` | 5 | `logs.go:139`, `:143`, `:155`, `:160` („per T0-(b)“, Verweis auf eine Tranche), `:174` |
| `remove` | 9 | `remove.go:101`, `:107`, `:114`, `:121`, `:126` („voll-schema“), `:128`, `:155`, `:157`, `:159`; Go-Sentinel-Namen (`ErrServiceUnsupported` u. a.) im Text der Exit-Codes sind für den Benutzer ebenso unlesbar |
| `template` (Wurzel, `list`) | 7 | `template.go:53`, `:58` „(slice-v1-template-init)“ und `:60` „(slice-later-local-templates)“ unter „Coming in later slices:“, `:92` (`Short`), `:96`, `:97`, `:100` |
| `up` | 6 | `up.go:122`, `:124`, `:128`, `:130`, `:136`, `:146` („LH-FA-UP-001/LH-FA-UP-001“, die Kennung doppelt) | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `devcontainer up`, `devcontainer down` | 2 | `devcontainer.go:174`, `:230` „Exit codes (LH-FA-CLI-006):“ | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
| `egress list` | 1 | `egress.go:110` „Exit codes (LH-FA-CLI-006):“ | <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->

Zusammen **88 Zeilen** in 12 Dateien; vier davon sind Meldungen, keine Hilfe (`config.go:445`, `:659`, `generate.go:181`, `:188`, siehe B).

**B. Meldungen** (Fehlertexte, Warnungen, Hinweise, Diagnose-Texte; vorläufig gezählt, ohne Zeilenkommentare und ohne `Code:`-Felder und ohne die Argumente `warn("LH-…", …)`/`relax("LH-…", …)`, die ein Code sind):

- **Treiber `internal/adapter/driving/cli`:** `config.go:445` und `:659` („… (LH-FA-DEV-003)“ im Fehlertext), `generate.go:181` („(LH-FA-DEV-003)“), `:188` („(LH-FA-DEV-006)“). <!-- d-check:ignore (wörtliches Zitat der Meldung) -->
- **Anwendung `internal/hexagon/application`** (Zeilen mit Kennung im Meldungstext, 33): `initproject.go:1281`, `:1285`; `config.go:364` („the LH-FA-ADD-005 state machine“), `:577`, `:813`; `generate.go:346` („LH-SA-FILE-002“); `doctor.go:276` („see slice-v0.1.1-doctor-container-awareness“ im Text `hostHintSkippedInContainer`), `:422`, `:430`, `:581`, `:600`, `:608`, `:622`, `:676`, `:942`, `:1268`, `:1294`; `doctor_sandbox.go:92`, `:97`, `:109`, `:114`, `:125`, `:142`, `:152`, `:158`, `:237`, `:243`; `devcontainer_sandbox.go:188`, dazu die Texte der Aufrufe `warn(…)`/`relax(…)` in `:317` bis `:339` (erste Zeile des Aufrufs nennt den Code, ob der Text selbst eine Kennung trägt, ist nicht gelesen). <!-- d-check:ignore (wörtliches Zitat der Meldung) -->
- **Domain:** `internal/hexagon/domain/projectname.go:45` („… does not match … (LH-FA-INIT-006)“). <!-- d-check:ignore (wörtliches Zitat der Meldung) -->
- **Adapter `internal/adapter/driven`:** keine Meldung mit Kennung (Suche über alle Nicht-Test-Dateien: Treffer nur in Kommentaren).
- **Textform der Warnzeile:** vier Stellen schreiben `Warning (<Code>): <Text>` auf die Ausgabe: `up.go:215`, `devcontainer.go:555`, `generate.go:232`, `init.go:443`. Der Benutzer sieht dort die Kennung als Präfix (Beispiel: `Warning (LH-FA-DEV-006): no git remote …`). Die Zeile ist **in der Spezifikation festgelegt** (siehe *Berührte Spec-Stellen*) und in Tests gepinnt (`internal/adapter/driving/cli/sandbox_test.go:55`, `:89`; `updown_cluster_test.go:64`; `internal/e2e/devcontainer_shell_docker_test.go:380`). Das ist eine Entscheidung des Projektinhabers, kein stilles Mitnehmen (Frage 1). <!-- d-check:ignore (wörtliches Zitat der Meldung) -->
- **Erzeugte Dateien:** Kommentare in den Vorlagen (`internal/hexagon/application/templates/**/*.tmpl`, `internal/adapter/driven/externaltemplates/templates/basic/*.tmpl`) tragen Kennungen (Beispiele: `Dockerfile.tmpl:5` „(LH-FA-DEV-004)“, `egress-init.sh.tmpl:2` „(LH-FA-DEV-008, ADR-0012)“, `gitignore.tmpl:6` „LH-SA-FILE-002 §611“). Das sind keine Meldungen des Binarys, aber Text, den der Benutzer in seinem Projekt liest. **Nicht gelesen im Hinblick auf die Wirkung** (Golden-Dateien, Replay-Pflicht [`harness/replay.md`](../../../../harness/replay.md)); Frage 3. <!-- d-check:ignore (wörtliches Zitat der Vorlagen-Kommentare) -->

**Nicht im Umfang, mit Grund:**

- **Das JSON-Feld `code` der Diagnostics** (`add.go:259` und die gleichartigen `Code: "LH-…"`-Stellen in allen Befehlen, `composesentinel.go:35`, `devcontainer_start.go:26`, `devcontainer_fallback.go:49`, `devcontainer_features.go:595`): maschinenlesbarer Vertrag ([`docs/maintainer/cli-json-contract.md`](../../../maintainer/cli-json-contract.md) §2: `diagnostics.code` folgt der Konvention „LH-Kennung der verursachenden Anforderung“). Der Slice ändert es nicht.
- **Kommentare im Code** ([`AGENTS.md`](../../../../AGENTS.md) §3.7: sie dürfen Kennungen tragen, wenn die Kennung die Herkunft nennt).
- **Test-Dateien und Hilfspakete für Tests** (`internal/adapter/driving/cli/jsontestutil/`: Meldungen von `t.Errorf` und die Beschreibungen in `coderegistry.go` sind nicht für den Benutzer; `coderegistry.go:44` bis `:51` enthalten außerdem deutsche Wörter, Folge-Kandidat).
- **`--quiet up` in der Hilfe** (Backtick-Wort im Usage-Text der Option `--quiet`): eigene Beobachtung [`BEO-CLI/hilfetext-backtick-wort-als-typname`](../observations/BEO-CLI/hilfetext-backtick-wort-als-typname/observation.md), gelesen im Binary am 2026-11-01 (`--quiet up`, `--no-git git init`, `--devcontainer .devcontainer/devcontainer.json`, `--template u-boot template list`, `--purge docker volume rm`, `--sandbox generate devcontainer`, `--allow-external-feature-sources generate devcontainer`); derselbe Wortlaut-Eingriff am selben Text. Ob der Slice diese Stellen mitnimmt, ist Frage 4.
- **Inhaltliche Berichtigungen:** `template.go:56` bis `:60` nennt „Coming in later slices“ unter anderem `init --template <name>`, obwohl `init` die Option `--template` trägt (gelesen im Binary: `init --help`). Der Slice streicht die Slice-Namen; ob der Satz dabei neu geschnitten wird (Option vorhanden, lokaler Pfad geplant), entscheidet Frage 5. `remove.go` sagt „In v0.3.0 …“ (Stand einer alten Version); nicht Gegenstand, Folge-Kandidat.

## 2. Definition of Done

Bis zu drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [ ] **Hilfe:** Alle Hilfetexte (`Short`, `Long`, Beispiele, Usage-Text jeder Option, rekursiv über den ganzen Befehlsbaum, auch die Ausgabe von `--help`) tragen keines der Muster `LH-<Familie>-…-<Nummer>`, `LH-AK-<Nummer>`, `SPEC-<Nummer>`, `ADR-<Nummer>`, `ARC-<Nummer>`, `slice-<name>` als Dateiverweis, `§<Zeichen>`, `T<Nummer>-(…)` und kein deutsches Wort („voll-schema“, „fachlich“); jede entfernte Kennung ist durch die Aussage in Worten ersetzt (Beispiel: „Exit codes:“ statt „Exit codes (LH-FA-CLI-006):“); der Inhalt der Zusagen ist derselbe. **Wächter:** ein Test in `internal/adapter/driving/cli` durchläuft den Cobra-Befehlsbaum rekursiv und meldet Befehlspfad und Fundstelle (`Short`, `Long`, `Example`, Option `--<name>`). Rot färbende Mutation: im `Long` von `add` die Kennung „(LH-FA-ADD-002)“ wieder einsetzen; zweite: im Usage-Text einer Option `--dry-run` ein `(LH-FA-CLI-007)` einsetzen. Der Implementer sieht beide rot. <!-- d-check:ignore (wörtliches Zitat der Mutation) -->
- [ ] **Meldungen:** Die Texte der Fehler, Warnungen, Hinweise und Diagnose-Texte in `internal/adapter/driving/cli`, `internal/hexagon/{application,port,domain}` und `internal/adapter/driven` tragen keines der Muster; das Feld `code` bleibt unverändert. Tests, die den alten Wortlaut pinnen, folgen dem neuen. **Wächter:** ein Test nach dem Muster von `internal/hexagon/application/config_filename_literal_test.go` (`go/parser`, Pfad ab Wurzel über `runtime.Caller`) durchläuft die Nicht-Test-Quellen der genannten Pakete und meldet String-Literale in `errors.New`, `fmt.Errorf`, `Message:`/`Hint:`/`MissingHint:`-Feldern und Fehlertext-Konstanten mit Kennungsmuster, ausgenommen die Felder `Code:`/`code:` und die Argumente von `warn(`/`relax(`. Rot färbende Mutation: in `projectname.go:45` „(LH-FA-INIT-006)“ wieder einsetzen. Ob die Textform `Warning (<Code>)` mitgeht, steht in Frage 1. <!-- d-check:ignore (wörtliches Zitat der Mutation) -->
- [ ] **Doku und Eintrag:** Wörtliche Zitate der geänderten Texte in `README.md`, `README.de.md`, `docs/user/` sind nachgezogen (Suche nach den alten Wortlauten; `grep` nach „Exit codes (“, „voll-schema“, „Warning (LH-“ über diese Dateien fand am 2026-11-01 keinen Treffer in `docs/user/`, `README.md`, `README.de.md`); `CHANGELOG.md` `[Unreleased]` `### Changed`: „Help texts and messages of all commands no longer carry requirement, specification or slice identifiers.“; Beispiele in [`docs/maintainer/cli-json-contract.md`](../../../maintainer/cli-json-contract.md) bleiben, solange sie nur das Feld `code` zeigen (Zeile 1269 zeigt eine Meldung ohne Kennung).
- [ ] `make gates` grün. Ein volles `make test-docker` läuft nicht in der Umsetzung (die Docker-Tests lesen den Wortlaut der Warnzeile, `devcontainer_shell_docker_test.go:380`); sie laufen in der CI, der Implementer nennt die Tests, die den Wortlaut tragen.
- [ ] Review (`.harness/skills/reviewer.md`), Verification, Closure-Notiz mit Lerneintrag, Beobachtungs-Register fortgeschrieben (`BEO-CLI/kennung-in-cli-ausgabe` auf `geplant`/`verkörpert` mit der Kennung dieses Slice), jedes Risiko aus §6 mit Ausgang.

## 3. Plan (vor Code)

Die Tabelle nennt jede Datei, die der Diff voraussichtlich berührt ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Implementer gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab. Alle Fundstellen sind gelesen (Quelltext oder Binary, 2026-11-01), soweit nicht anders vermerkt.

| Datei / Komponente | Änderungs-Art | Begründung / Fundstelle |
|---|---|---|
| `internal/adapter/driving/cli/add.go`, `config.go`, `devcontainer.go`, `doctor.go`, `down.go`, `egress.go`, `generate.go`, `init.go`, `logs.go`, `remove.go`, `template.go`, `up.go` | update (Hilfetexte, Fehlertexte in `config.go:445`, `:659`, `generate.go:181`, `:188`) | Fundstellen in §1 A und B |
| `internal/adapter/driving/cli/up.go`, `devcontainer.go`, `generate.go`, `init.go` (Warnzeile) | update, nur bei Antwort „ja“ auf Frage 1 | Zeilen `up.go:215`, `devcontainer.go:555`, `generate.go:232`, `init.go:443` |
| `internal/adapter/driving/cli/help_identifiers_test.go` (neu) | neu | Wächter A; Muster der Tests und Test-Doubles aus `cli_test.go` und `export_test.go` (nicht im Detail gelesen; der Implementer liest sie vor dem Schreiben) |
| `internal/hexagon/application/message_identifiers_test.go` (neu, oder im Paket, das alle genannten Pakete lesen kann) | neu | Wächter B nach `config_filename_literal_test.go` (gelesen: `go/parser`, `runtime.Caller`, Pfad ab `go.mod`) |
| `internal/hexagon/application/initproject.go`, `config.go`, `doctor.go`, `doctor_sandbox.go`, `devcontainer_sandbox.go`, `generate.go` | update (Meldungstexte) | Fundstellen in §1 B |
| `internal/hexagon/domain/projectname.go` | update (Zeile 45) | §1 B |
| `internal/adapter/driving/cli/sandbox_test.go`, `updown_cluster_test.go`, `doctor_json_test.go` und Tests, die einen Wortlaut mit Kennung pinnen; `internal/e2e/devcontainer_shell_docker_test.go` | update (nur bei Frage 1 / bei gepinntem Wortlaut) | gelesen: `sandbox_test.go:55`, `:89`, `updown_cluster_test.go:64`, `devcontainer_shell_docker_test.go:380`, `doctor_json_test.go:214` (Testeingabe); weitere Treffer sucht der Implementer per `grep` über `*_test.go` und `internal/e2e` |
| `spec/spezifikation.md` | update, nur bei Antwort „ja“ auf Frage 1 | Zeilen 887, 916 (Warnzeile) |
| `CHANGELOG.md` | update (`[Unreleased]` `### Changed`) | ein englischer Satz, siehe DoD |
| `README.md`, `README.de.md`, `docs/user/` | unverändert (belegt) | Suche am 2026-11-01 nach `LH-`, `slice-` und `Warning (LH` in `docs/user`, `README.md`, `README.de.md`: kein Treffer; wörtliche Hilfe-Zitate nicht gefunden |
| `docs/maintainer/cli-json-contract.md` | unverändert | Zeilen mit Kennung betreffen `diagnostics.code` und Verweise im Maintainer-Dokument (kein Nutzertext); Beispiel Zeile 1269 trägt keine Kennung in `message` |
| `internal/hexagon/application/templates/**`, `internal/adapter/driven/externaltemplates/templates/**` | unverändert, bis Frage 3 entschieden | erzeugte Dateien, Golden-Pflicht |
| `docs/plan/planning/open/slice-v1-hilfetexte-ohne-kennungen.md`, `docs/plan/planning/open/README.md` | neu / update | Lifecycle; Zeile im Bestand |
| `docs/plan/planning/observations/BEO-CLI/kennung-in-cli-ausgabe/` | update bei der Closure (`state.md`, `evidence/slice-v1-hilfetexte-ohne-kennungen.md`) | Ausgang der Beobachtung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Der Projektinhaber startet die Umsetzung in einer eigenen Devcontainer-Instanz mit `u-boot devcontainer up --instance` (Entscheidung des Projektinhabers, 2026-11-01); kein Termin. Keine Abhängigkeit von einem anderen Slice. WIP-Limit 1 in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): wenn die Zählung der Meldungen deutlich über die vorläufigen 33 + 1 + 4 Stellen steigt oder Frage 1 mit „ja“ die Spezifikation und die Docker-Tests mitzieht; dann in zwei Slices schneiden (Hilfe · Meldungen), der Wächter A gehört zum ersten.
- `in-progress` → `open` (blockiert): solange Frage 1 oder Frage 3 unbeantwortet ist und die Umsetzung davon abhängt.

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Review und Verification in eigenen Kontexten, kein Self-Review; Closure-Notiz in §7 mit Lerneintrag; `git mv` nach `done/` als eigener Commit nach dem Inhalt ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2).

Lerneintrag: Mögliche Antwort ist, dass die Beobachtung `BEO-CLI/kennung-in-cli-ausgabe` mit der Kennung dieses Slice auf `geplant`/`verkörpert` gesetzt wird und beide Wächter die Regel tragen; der Implementer nennt, ob ein Hilfetext vor dem Wächter Kennungen trug, die das Muster nicht kennt (Lücke des Musters).

## 6. Risiken und offene Punkte

**Fragen an den Projektinhaber** (Entwurf; keine ist entschieden):

1. **Warnzeile `Warning (<Code>): <Text>`:** Der Code steht als Präfix in der Ausgabe von `up`, `devcontainer up`/`shell`, `generate`, `init` und ist in `spec/spezifikation.md` (887, 916) so festgelegt. Weg A: Präfix streichen (`Warning: <Text>`), Spezifikation und Tests (auch ein Docker-Test) nachziehen; der Code bleibt im JSON. Weg B: Präfix lassen, weil die Zeile ein Vertrag ist, und nur die Texte ändern. Empfehlung: A, sonst bleibt die häufigste Kennung in der Ausgabe stehen.
2. **Release:** Patch `v0.17.1` (Einordnung *Changed*) oder mit dem nächsten Release (Vorschlag `v0.18.0` im Entwurf `slice-v1-katalog-feature-jdk-gepinnt`).
3. **Vorlagen der erzeugten Dateien:** Kommentare in den Vorlagen tragen Kennungen (§1 B). Gleicher Slice (Golden-Dateien und Replay-Evidence nötig) oder Folge-Slice.
4. **Backtick-Wörter im Usage-Text** (`--quiet up` u. a.): im selben Zug mitnehmen (selbe Texte) oder getrennt lassen ([`BEO-CLI/hilfetext-backtick-wort-als-typname`](../observations/BEO-CLI/hilfetext-backtick-wort-als-typname/observation.md)).
5. **`template`-Hilfe:** „Coming in later slices“ neu schneiden (vorhandene Option `init --template <name>` · geplanter lokaler Pfad) oder nur die Slice-Namen streichen.

**Risiken:**

- **1. Wächter A sieht nur, was der Cobra-Baum trägt:** Texte, die zur Laufzeit zusammengesetzt werden (Fehlertexte), liegen außerhalb; dafür Wächter B. Wächter B sieht nur Literale in den genannten Aufrufen; ein Text, der über eine Hilfsfunktion oder ein `%s`-Argument eine Kennung bekommt (zum Beispiel ein `code`, der in `Message` eingesetzt wird), entkommt. — **Ausgang:** offen bis zur Umsetzung; der Implementer nennt im Bericht die Grenze der Wächter.
- **2. Tests pinnen den Wortlaut, auch außerhalb des Pakets:** Docker-Tests (`internal/e2e`) laufen nicht in `make gates`; ein dort gepinnter Wortlaut bricht erst in der CI. — **Ausgang:** offen; der Implementer sucht `grep` über `internal/e2e` und nennt jeden Treffer im Plan.
- **3. Zusage-Verlust beim Umschreiben:** „Exit codes (LH-FA-CLI-006):“ → „Exit codes:“ ist verlustfrei; Sätze wie „per LH-FA-UP-001“ tragen die Kennung als Begründung, die ein Wort braucht. — **Ausgang:** der Reviewer liest je Befehl Alt und Neu nebeneinander (Pflicht in §7). <!-- d-check:ignore (wörtliches Zitat des Hilfetexts) -->
- **4. Wächter-Test und Lastenheft:** Die Hilfe darf Wörter wie „ADR“ oder „SPEC“ nur als Kennung-Muster nicht tragen; ein Muster ohne Ziffern wäre zu breit. — **Ausgang:** entfällt mit dem Muster *Buchstaben + Bindestrich + Ziffer*; der Test pinnt es.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/adapter/driving/cli` (CLI, Greenfield), `internal/hexagon/application` (APP, Greenfield), `internal/hexagon/domain` (DOM, Greenfield), `docs/plan/` (DOK, Greenfield), bei Frage 1 `spec/` (DOK) und `internal/e2e` (E2E, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle ≥ 2 von 3 Achsen ist dort entschieden. Alle berührten Sub-Areas Greenfield (kein Modus-Begründungsblock nötig).

**Vorgelagert — offene Beobachtungen sichten** (Register `docs/plan/planning/observations/` gelesen, nicht verändert):

- [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) — offen, vier Belege; dieser Slice ist der Folge-Slice der Entscheidung vom 2026-11-01. Das Register wird bei der Closure fortgeschrieben, nicht im Entwurf.
- [`BEO-CLI/hilfetext-backtick-wort-als-typname`](../observations/BEO-CLI/hilfetext-backtick-wort-als-typname/observation.md) — offen, ein Beleg; berührt dieselben Usage-Texte (Frage 4).
- [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) — verkörpert; Prüflinse für die neuen Kommentare der Wächter-Tests ([`AGENTS.md`](../../../../AGENTS.md) §3.7, §3.9).
- [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) — verkörpert; Prüflinse für diese Plan-Tabelle: die Zeilenzahlen in §1 sind vorläufig und tragen den Vermerk „der Implementer zählt neu“.

Evidenz-/Diskrepanz-Risiko: mittel (Zeilenzahlen vorläufig, Frage 1 ändert den Umfang).
