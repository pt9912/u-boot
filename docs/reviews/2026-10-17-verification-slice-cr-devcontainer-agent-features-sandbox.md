# Verification Evidence: slice-cr-devcontainer-agent-features-sandbox — 2026-10-17

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/in-progress/slice-cr-devcontainer-agent-features-sandbox.md) (Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/welle-devcontainer-sshd-claude-code.md), erster Slice), Commits `5de388c` (Lastenheft 0.11.0, CHANGELOG), `4df2b54` (Ausnahme 1 um die `doctor`-Sicherheits-Einträge erweitert, CHANGELOG sagt nur den Vertrag, Review-Report) und `48b2537` (Planungsdateien: fünf Slices, Entscheidung 5). Maßstab: §1 (CR-Kern wörtlich, Entscheidungen 1 bis 5, Historie-Zeile), §2 (DoD), §3 (Plan), §6 (Risiken), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging). Review-Report: [`2026-10-17-slice-cr-devcontainer-agent-features-sandbox.md`](2026-10-17-slice-cr-devcontainer-agent-features-sandbox.md) (lief gegen `5de388c`, also vor der Erweiterung in `4df2b54`).

## Verification Evidence

Scope:
- Slice: [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/in-progress/slice-cr-devcontainer-agent-features-sandbox.md)
- IDs: [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md)
- Artefakte: `spec/lastenheft.md` (Version, §3 `LH-FA-DEV-006`, `LH-FA-DEV-009`, `LH-FA-CLI-005` Ausnahme 1, §16), `CHANGELOG.md` `[Unreleased]`, die Planungsdateien der Welle <!-- d-check:ignore (Pfade und Kennungen als Zitat) -->

DoD-Abgleich (§2 des Slice):
- [x] Entscheidungen 1 bis 5 stehen in §1; Prüfungstabelle vorhanden und vom Review bestätigt — Evidence: Slice §1; die Tabelle ist um die Erweiterung (Ausnahme 1, doctor) in `48b2537` nachgezogen. Die Erweiterung ist **nicht** im Review-Report geprüft (Auflage A-1); der Verifier hat sie selbst geprüft (siehe unten).
- [~] Lastenheft **0.11.0** — Zusagen in `LH-FA-DEV-009`, `LH-FA-DEV-006`, `LH-FA-CLI-005` Ausnahme 1 stehen, Versionsfeld und §16-Zeile stehen; **Abweichung**: die Historie-Zeile, die DoD und §1 zitieren („Ausnahme 1 unter `--quiet` nennt auch die Sicherheits-Einträge von `u-boot doctor`“), ist nicht wörtlich die geltende §16-Zeile (diese sagt „nennt auch diesen Login-Dienst und die Sicherheits-Einträge von `u-boot doctor`“). Das Lastenheft ist vollständiger als der Slice-Text; der Slice ist anzupassen (V-1, Auflage A-2). Die Suche belegt keinen Widerspruch zu älteren Zusagen (eigene, erweiterte Suche unten). <!-- d-check:ignore (Kennungen als Zitat) -->
- [x] `make gates` grün — selbst gefahren, Exit 0 (siehe Sensors).
- [x] Review durchgeführt, Report liegt vor (0 HIGH, 1 MEDIUM, 4 INFO); Bewertung der Findings unten. Einschränkung: der Review deckt `4df2b54` nicht ab (A-1).
- [x] Doku-Update, falls öffentlicher Vertrag berührt — `CHANGELOG.md` trägt den Vertrag (F-1 behoben); Benutzerhandbuch und Spezifikation folgen laut Plan in den Code- und Spec-Slices; kein Widerspruch zum jetzigen Code-Stand, solange die Zusage dort als „folgt“ steht.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Wird bei der Closure gefüllt“); gehört in die Closure.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 — offen, Closure.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — Ausgang steht an allen vier Risiken, aber drei in nicht-kanonischer Form (V-3, Auflage A-3).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Closure; die Folge-Slices (vier, in `open/`) sind als Dateien und in der Welle §4 getragen (V-2).

### Wörtlichkeit (a): CR-Kern gegen Lastenheft

Satz für Satz verglichen (Slice §1 gegen `git show 4df2b54:spec/lastenheft.md`):

| Stelle | Ergebnis |
|---|---|
| `LH-FA-DEV-009` (Einleitung, drei Punkte) | wörtlich gleich (Datei-Links im Slice, `#…`-Anker im Lastenheft) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-DEV-006` Zusatz „einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.“ | wörtlich gleich; Vorgänger-Punkt endet mit „;“ | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-CLI-005` Ausnahme 1 | wörtlich gleich dem Block „Wortlaut der Ausnahme 1“ (Aufzählung der Lockerungen, Login-Dienst, Warn-Einträge von `u-boot doctor` mit den drei Klammerbezügen); Ausnahme 2 bis 4 und der übrige Text unverändert (Diff zeigt nur Zeile 1) | <!-- d-check:ignore (Kennungen als Zitat) -->
| §16 Zeile 0.11.0 | Inhalt trägt alle drei Anforderungen, ist aber **nicht** wörtlich die im Slice zitierte Fassung (V-1) |
| Versionsfeld 0.10.0 → 0.11.0 | konform |

### Satzstruktur und Zahl „vier“ (b)

Ausnahme 1 koordiniert „die Ausweisung von … und eines … Login-Dienstes … in der Befehlsausgabe (…) **sowie** die Warn-Einträge von `u-boot doctor`, die …“; der Satz ist grammatisch eindeutig, die Klammer nach „Befehlsausgabe“ trägt `LH-FA-DEV-006`/`-007`/`-008`, die zweite `LH-FA-DEV-009`/`-008`/`-007`. Die Nummerierung 1 bis 4 und die Einleitung „Vier Ausnahmen“ sind unverändert. Die Zahl „vier“ stimmt an allen Stellen: `spec/lastenheft.md` (Zeile 189, §16 0.9.0), `CHANGELOG.md` (0.10.0-Absatz), `docs/user/benutzerhandbuch.md` (beide `--quiet`-Zeilen), `spec/spezifikation.md` `SPEC-043`. Der **Inhalt** von Handbuch und `SPEC-043` nennt die neue Ausnahme-1-Fassung noch nicht (Login-Dienst, `doctor`-Einträge); das ist geplant: Handbuch im Code-Slice, `SPEC-043` im Spec-Slice. Bis dahin gilt Lastenheft vor Spezifikation (Source Precedence), kein Widerspruch der Zahl. <!-- d-check:ignore (Kennungen als Zitat) -->

### Kollisionen mit älteren Zusagen (c)

Eigene Suche im Lastenheft 0.11.0 (`doctor`, `--quiet`, `sichtbar`, `Zugangsdaten`, `Egress`, `Warn`, dazu `strict`): <!-- d-check:ignore (Suchbegriffe als Zitat) -->

- `LH-FA-CLI-005` Zeile 189: „`u-boot doctor` zeigt nur Fehler-Einträge. Vier Ausnahmen …“ — verträglich, die Ausnahmen sind genau der Rahmen dazu; Ausnahme 2 nennt `u-boot doctor` bereits (Degradations-Fallbacks), die Erweiterung in Ausnahme 1 steht daneben. Überschneidung benigne: „Wegfall der Egress-Restriktion“ (Ausnahme 2) und „Egress-Restriktion“ (Ausnahme 1, Warn-Einträge) sind beide sichtbar, keiner schaltet den anderen stumm. <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-FA-DEV-008` (Egress): „`u-boot doctor` prüft die Konfiguration“ — verträglich; die Klammer nennt die Egress-Warn-Einträge als sichtbar. <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-FA-DEV-007` (Degradation, „jeder Fallback … in `u-boot doctor` ausgewiesen“) — verträglich. <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-FA-DEV-009` — Erweiterung, Git-Zusagen wörtlich erhalten; die `doctor`-Warnung zur Quelle ist die in Ausnahme 1 genannte. <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-FA-DIAG-001` (`warn` ohne Eskalation, `--strict` macht jedes `warn` zu Exit ungleich 0) — verträglich: Sichtbarkeit unter `--quiet` ändert nicht die Zählung; die ausgeblendeten Warnungen zählen weiter für `--strict`. <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-NFA-SEC-001`, `LH-NFA-SEC-004`, `LH-FA-DEV-003` — wie im Review (verträglich). <!-- d-check:ignore (Kennungen als Zitat) -->
- Kein weiterer Treffer zu `Zugangsdaten`/`Warn` im Normtext, der eine Zusage über die Sichtbarkeit trifft.

Ergebnis: kein Widerspruch zu einer älteren Zusage.

### Wortlaut-Lücke „Sicherheitslockerungen des Docker-Daemons“ (d)

Befund aus dem Code (`internal/hexagon/application/doctor_sandbox.go`): `checkSandboxRuntime` gibt für `nestedRuntime: docker` die Stufe **OK** zurück (Meldung nennt `--privileged`); die einzige Warnung mit `docker` ist in `checkSandboxEgress` (Egress zusammen mit `docker`); der Fall `/dev/fuse` fehlt ist ein Degradations-Fallback (Ausnahme 2), `nestedRuntime` ohne Sandbox-Profil ist „has no effect“. Es gibt also heute **keinen** `doctor`-Warn-Eintrag, der allein eine „Sicherheitslockerung des Docker-Daemons“ ausweist. <!-- d-check:ignore (Pfad und Bezeichner als Zitat) -->

Bewertung: **leere Menge, kein Vertragsloch.** Ausnahme 1 definiert eine Klasse („die Warn-Einträge …, die eine Sicherheitseigenschaft … betreffen“) und nennt Beispielkategorien; ein Relativsatz verspricht keinen Eintrag, den es nicht gibt. Die Zusagen, die tatsächlich verlangt werden, sind erfüllbar und decken die Sicherheitslage: (1) `--privileged` des Docker-Daemons ist eine Sicherheitslockerung und steht im ersten Satzteil von Ausnahme 1 unter `--quiet` in der Befehlsausgabe von `init`, `generate` und `up`; (2) die Kombination Egress und `docker` ist als `doctor`-Warnung sichtbar. Failure-Szenario geprüft: ein Nutzer mit `nestedRuntime: docker` ohne Egress ruft `doctor --quiet` auf und sieht nichts. Das bricht keine Zusage (weder `LH-FA-DEV-007` noch `LH-FA-DIAG-001` verlangen dort ein `warn`), und die Lockerung ist in der Befehlsausgabe schon ausgewiesen. Verworfen als Lücke. Restrisiko Leserirrtum: der Wortlaut könnte als Zusage eines Warn-Eintrags gelesen werden. Der Slice [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/open/slice-v1-cli-quiet-sicherheitswarnungen.md) §1/§6 benennt den Fall bereits; der Spec-Slice trägt ihn in `SPEC-043` ein (Auflage A-5, keine Vertragsänderung). <!-- d-check:ignore (Kennungen als Zitat) -->

### Einstufung im Quiet-Slice (e)

Tabelle in [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/open/slice-v1-cli-quiet-sicherheitswarnungen.md) §1 gegen den Code (alle `SeverityWarn` in `doctor_sandbox.go`): vollständig und korrekt abgebildet (credentials ×2, egress ×2, runtime ×2; Zeilen 80, 127, 137, 187, 192 und der `DegradationFallback`-Fall). Die übrigen `SeverityWarn` in `doctor.go` (Docker-Version, Pflichtdateien, `compose.yaml`, `enabled`-Schlüssel, `forwardPorts`, Feature-Allowlist, Feature-Drift) betreffen nicht das Sandbox-Profil; die Tabelle nennt sie nicht, ein Satz „übrige Warn-Einträge, darunter Feature-Allowlist und -Drift (Drittanbieter-Herkunft, Entscheidung 5), nicht eingestuft“ fehlt (INFO, V-6).

Gegen den Wortlaut „eine Sicherheitseigenschaft des Sandbox-Profils betreffen“:
- Klartext-Zugangsschlüssel, https-Klonquelle ohne Token-Quelle, Egress mit `docker`: eindeutig gedeckt (die Klammer nennt Quelle der Zugangsdaten und Egress-Restriktion).
- `runtime` „has no effect“ (`nestedRuntime` ohne Sandbox): keine Sicherheitseigenschaft, Konfigurationsmeldung — Einstufung trägt.
- **`egress` „has no effect“ (`egress.enabled` ohne Sandbox-Profil): nicht eindeutig.** Die Klammer in Ausnahme 1 nennt „Egress-Restriktion nach `LH-FA-DEV-008`“ ohne Einschränkung; der Eintrag betrifft die Egress-Restriktion, nur eben ihr Nichtwirken. Failure-Szenario: Nutzer setzt `egress.enabled: true`, das Profil ist (noch) nicht `sandbox`, die CI ruft `doctor --quiet`, sieht nichts und hält die Restriktion für aktiv. Die Auslegung des Quiet-Slice (nicht sichtbar) ist vertretbar („des Sandbox-Profils“: das Profil ist nicht aktiv), aber der Wortlaut entscheidet es nicht. Der Slice führt es als Risiko („weiter offen bis zum Review“); die Entscheidung gehört dem Projektinhaber **vor** dem Code-Slice (Auflage A-4), nicht in die Closure des Quiet-Slice. <!-- d-check:ignore (Kennungen als Zitat) -->

### Risiko-Form im CR §6 (zulässig?)

Kanonische Form (`slice.template.md`, `modul-05-planning-harness.md` §Offene Risiken): *eingetreten* → Carveout oder Folge-Slice mit ID · *entfallen* → Grund · *weiter offen* → **Beobachtungs-Register** (`BEO-<KUERZEL>/<slug>`). Im CR tragen drei Risiken „weiter offen“ mit Verweis auf einen Folge-Slice und eine Rückführung, **ohne** `BEO-`-Verweis (Risiko 1 Wortlaut/Code, Risiko 2 Anmeldeweg, Risiko 3 Feature-Metadaten); das vierte ist „entfallen“ mit Grund. Das Register enthält nur `BEO-CLI/kommentar-beschreibt-abwesendes` (nicht einschlägig). Die Form ist **nicht kanonisch**. Repo-Präzedenz (mehrere `done/`-Slices schreiben „weiter offen: Re-Evaluierungs-Trigger der ADR“) bindet die Regel nicht: dort stand der Trigger in einer ADR, hier steht er in einem anderen, noch offenen Slice. Zudem sind die drei Risiken **bedingte Ereignisse der Folge-Slices**, keine Dauerzustände: bei der Closure des CR ist „weiter offen“ nur wahr, wenn sie an einen Zähler hängen. Zulässig sind bei der Closure: (i) Risiko 1: **entfallen** mit der Begründung dieser Evidence (Wortlaut und Einstufung gegen den Code gelesen, leere Menge belegt) — oder Register; (ii) Risiko 2 und 3: Register-Eintrag `BEO-DOK/…` mit dem Folge-Slice als Prüfer, oder **eingetreten**-frei als „weiter offen → BEO“. Für den Quiet-Slice gilt dasselbe (vier „weiter offen“ ohne `BEO-`); er ist noch `open/`, die Auflage gilt bei dessen Closure, der Planer kann sie früher beheben. Nicht blockierend für den Lastenheft-Stand, **blockierend für die Closure des CR** (Auflage A-3).

### Bewertung der Review-Findings

| Finding | Bewertung | Begründung |
|---|---|---|
| F-1 MEDIUM (CHANGELOG sagte Umsetzendes im Indikativ) | **behoben, nicht blockierend** | `4df2b54` ergänzt „This entry states the contract only: the run-time handover of agent credentials, the new `doctor` checks and their visibility under `--quiet`, the login-service disclosure and the `sshd` and Claude Code features themselves follow in a later release.“ Der Einleitungsabsatz bleibt im Indikativ, der Schlusssatz nimmt jetzt alle Umsetzungsteile aus. Restpflicht: beim Landen der Code-Slices den Satz „follow in a later release“ entfernen bzw. anpassen (A-6). |
| F-2 INFO („Agenten-Feature“, Variablenname) | nicht blockierend | Sache der Spezifikation (Katalog-Metadaten), Risiko 2 trägt es, Spec-Slice zuständig. |
| F-3 INFO (Klammer „Passwort-Login“ erklärt oder begrenzt) | nicht blockierend | die Klammer ist in `LH-FA-DEV-006` und Ausnahme 1 gleich; ein reiner Schlüssel-Login ist ausdrücklich Nicht-Ziel; der sshd-Slice legt den Fall fest. Beim Spec-Slice klären. | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 INFO (Überschrift „Git-Zugangsdaten“) | nicht blockierend | Anker-Stabilität, Entscheidung im Slice benannt, `make docs-check` löst die Anker. |
| F-5 INFO (Slice nicht nachgezogen) | **behoben bis auf V-1/V-2** | Kein „Teilfrage“, „offen“ zur Teilfrage oder „bedingt“ zur sshd-Entscheidung mehr im CR und in den Slice-Texten (Suche über Welle, vier Folge-Slices, `open/README.md`, Roadmap); verbleibende „bedingt“ betreffen Messergebnisse in `slice-v1-devcontainer-feature-sshd` und `slice-v1-devcontainer-feature-claude-code` (Tabellenzeilen, kein Vertragspunkt). Offen: V-1 (Historie-Zitat), V-2 (Tippfehler im Link-Text). | <!-- d-check:ignore (Kennungen als Zitat) -->

## Planungsdateien der Welle

Geprüft: Welle §1 bis §5, `open/README.md`, Roadmap-Zeile, die vier Folge-Slices.
- „fünf Slices“ und Entscheidungen 1 bis 5 durchgängig: Welle §1/§3/§4/§5 („fünf Slices“, „fünf Entscheidungen“), Roadmap-Zeile („fünf Slices“), `open/README.md` (vier Folge-Slices als Zeilen 15 bis 18, CR-Zeile nennt Sicherheitswarnungen unter `--quiet`), Slice `slice-v2-spec-devcontainer-sshd-claude-code` („die vier Slices“ vor dem Spec-Slice, korrekt). <!-- d-check:ignore (Kennungen als Zitat) -->
- Keine Aussage „Teilfrage offen“ oder „bedingt“ zu den Entscheidungen.
- Der neue Slice [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/open/slice-v1-cli-quiet-sicherheitswarnungen.md) trägt die Abnahmekriterien des CR: Kennzeichen am Eintrag, `--quiet` zeigt Klartext-Schlüssel, fehlende Token-Quelle und Egress mit `docker`; „has no effect“ und die Drittanbieter-Warnung bleiben aus; JSON, Stufen und Exit-Codes unverändert; `--strict` unverändert; Reihenfolge (vor `claude-code` zwingend, vor `sshd` aus Wahl) mit Begründung; Mutationen und Boundary-Test je Check-Kennung.
- Mängel: (V-2) Welle §4, Zeile 2: Link-Text `slice-v1-cli-quiet-sicherheswarnungen` (Tippfehler „sicherheswarnungen“; das Ziel löst, der Text nicht); (V-7) Welle §1 „heute gilt es nur für Git-Zugangsdaten“ ist nach dem Lastenheft-Commit schon veraltet und bei der Closure der Welle falsch. <!-- d-check:ignore (Kennungen als Zitat) -->

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `0 issues.` (golangci-lint), `coverage-gate: OK — coverage 91.10% meets threshold 90%`, `d-check: 284 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK` |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Evidence erneut gefahren (siehe Übergabe) |
| `make test-docker` | not run | kein Docker-Laufzeitpfad berührt (Diff: Lastenheft, CHANGELOG, Planungsdokumente) |
| `make verify-depguard` | not run | keine Importregel berührt |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Lastenheft-Text wörtlich gleich dem CR-Kern; Code-Beleg folgt im Code-Slice für Claude Code |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | Zusatz wörtlich; Code-Beleg folgt im sshd-Slice |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Ausnahme 1 wörtlich; Einstufungstabelle gegen `doctor_sandbox.go` gelesen; Code-Beleg folgt im Quiet-Slice |
| [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) | im Diff unverändert (Accepted-immutable); Entscheidung 3 folgt ihr |

Plan-vs-Code-Diff (§3 gegen `5de388c^..HEAD`):
| Plan-Zeile | Diff | Befund |
|---|---|---|
| `spec/lastenheft.md` (§3 drei Anforderungen, Versionsfeld, §16) | geändert | konform; Historie-Zitat im Slice weicht ab (V-1) |
| dieser Slice (§1 Entscheidungsdokumentation, Prüfungstabelle, §7) | §1 und Prüfungstabelle geändert, §7 leer | konform bis auf §7 (Closure) |
| nicht im Plan: `CHANGELOG.md` `[Unreleased]` | geändert (`5de388c`, `4df2b54`) | Abweichung vom Plan, vertretbar: DoD „Doku-Update, falls ein öffentlicher Vertrag berührt“; §3 trägt es nicht, beim Closure nachtragen |
| nicht im Plan: Review-Report, Welle, `open/README.md`, Roadmap, vier Folge-Slices | geändert bzw. neu (`4df2b54`, `48b2537`) | durch Entscheidung 5 und den neuen Slice begründet; Welle und Roadmap tragen es |
| Nicht-Ziele §1 (Code, Katalog, Egress-Hosts, Drittanbieter-Warnung unter `--quiet`, Spezifikation, Architektur, ADR) | nicht berührt; `spezifikation.md`, `architecture.md`, `internal/`, `templates/` im Diff unverändert | eingehalten |
| Decken-Regel `AGENTS.md` §3.4 | Lastenheft nennt weder Slice, Welle, ADR noch Produktnamen | eingehalten |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — nicht berührt (kein Code, keine Importregel).
- Kein Nachlauf-Review für `4df2b54` — nicht Aufgabe des Verifiers; siehe Auflage A-1.

Commit / Artefakt:
- `5de388c`, `4df2b54`, `48b2537` (HEAD zum Zeitpunkt der Verifikation `48b2537`)

## Befunde des Verifiers

| ID | Stufe | Befund | Auflage |
|---|---|---|---|
| V-1 | LOW | Historie-Zeile in Slice §1 und DoD ist nicht wörtlich die geltende §16-Zeile (fehlt „diesen Login-Dienst und“). | A-2 |
| V-2 | LOW | Welle §4: Link-Text `slice-v1-cli-quiet-sicherheswarnungen` (Tippfehler). | A-2 | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-3 | MEDIUM (Closure) | CR §6: drei Risiken „weiter offen“ ohne `BEO-`-Verweis; Quiet-Slice §6: vier. | A-3 |
| V-4 | MEDIUM (Entscheidung) | Egress „has no effect“ gedeckt oder nicht: Wortlaut entscheidet es nicht. | A-4 |
| V-5 | INFO | „Sicherheitslockerungen des Docker-Daemons“ ist heute eine leere Menge; kein Vertragsloch. | A-5 |
| V-6 | INFO | Tabelle im Quiet-Slice nennt die übrigen `doctor`-Warnungen (Feature-Allowlist, -Drift) nicht als ausgenommen. | A-5 |
| V-7 | INFO | Welle §1 „heute gilt es nur für Git-Zugangsdaten“ bei Closure falsch. | A-2 |

## Verdikt

**Erfüllt mit Auflagen.** Der Vertragsinhalt (Lastenheft 0.11.0, CHANGELOG) entspricht dem CR-Kern wörtlich, die Satzstruktur und die Zahl „vier“ stimmen, es gibt keinen Widerspruch zu älteren Zusagen, und `make gates` sowie `make doc-trace` sind grün. Die Wortlaut-Lücke aus (d) ist eine leere Menge, kein Vertragsloch. Der Slice darf **nicht** nach `done/`, bevor die Auflagen A-1 bis A-4 erledigt sind; die Folge-Slices dürfen starten, sobald A-2 und A-4 geklärt sind (A-4 vor dem Code des Quiet-Slice).

Auflagen für die Closure:
- **A-1:** Ein Nachlauf-Review (kurz) des Diffs `5de388c..4df2b54` auf die Erweiterung von Ausnahme 1; der vorhandene Review lief davor.
- **A-2:** Slice §1 und DoD auf die geltende §16-Zeile anpassen (V-1); Link-Text in Welle §4 berichtigen (V-2); Welle §1 „heute gilt es nur …“ im Präsens des Lastenhefts 0.11.0 formulieren oder streichen (V-7).
- **A-3:** Ausgang der Risiken in kanonischer Form: Risiko 1 „entfallen“ mit Begründung (Wortlaut und Einstufung gegen den Code gelesen, leere Menge belegt, diese Evidence) oder Register-Eintrag; Risiko 2 und 3 als `BEO-DOK/<slug>` im Register mit Verweis auf den prüfenden Folge-Slice. Dasselbe für die Risiken des Quiet-Slice bei dessen Closure.
- **A-4:** Entscheidung des Projektinhabers zur Einstufung des Egress-Eintrags „has no effect“ (sichtbar oder nicht) **vor** dem Code des Quiet-Slice; bei „sichtbar“ ändert sich Ausnahme 1 oder nur die Tabelle (Klammer „Egress-Restriktion“ deckt es dann ausdrücklich).
- **A-5:** Im Quiet-Slice den Fall `nestedRuntime: docker` (OK-Stufe, leere Menge) und die ausgenommenen Feature-Warnungen in der Tabelle benennen; `SPEC-043` im Spec-Slice trägt es. <!-- d-check:ignore (Kennungen als Zitat) -->
- **A-6:** Beim Landen der Code-Slices den CHANGELOG-Satz „follow in a later release“ anpassen; CR-Closure §7 trägt den Steering-Loop-Lerneintrag (Klassen aus dem Review plus V-1 bis V-7).

## Nachverifikation abd906f

Gegenstand: Commit `abd906f` (Behebung F-1 des Nachlauf-Reviews [`2026-10-17-slice-cr-devcontainer-agent-features-sandbox-2`](2026-10-17-slice-cr-devcontainer-agent-features-sandbox-2.md)); Zusagen [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).

1. **Wortlaut-Gleichheit:** Der Ausnahme-1-Text (Lastenheft Zeile 191, Slice §1 Zeile 72) ist nach Normalisierung der Link-Pfade zeichengleich, einschließlich des Zusatzes „dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift; die Meldung, dass die Container-Runtime ohne Sandbox-Profil keine Wirkung hat, ist eine Konfigurationsmeldung und gehört nicht dazu“. Die §16-Zeile 0.11.0 (Lastenheft Zeile 1562) ist zeichengleich mit der Historie-Zeile in Slice §1 (Zeile 77) und der DoD (Zeile 114): „die Warn-Einträge von `u-boot doctor` zu Sicherheitseigenschaften des Sandbox-Profils, darunter die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift.“ Befund: ok.
2. **F-1 beseitigt:** ja. Egress „no effect“ ist ausdrücklich eingeschlossen, Runtime „no effect“ ausdrücklich als Konfigurationsmeldung ausgeschlossen; beide Richtungen stehen im Vertrag. Failure-Szenarien: (a) Leser deutet „Sicherheitseigenschaft des Sandbox-Profils“ weit und zieht Runtime ein: verworfen, der Ausschluss ist wörtlich. (b) Egress-Meldung „nicht greift“ gegen die Programm-Meldung „has no effect“: verworfen, der Vertrag nennt die Wirkung (Restriktion greift nicht), nicht den Meldungstext; die Einstufungstabelle ordnet beide Einträge je Check-Kennung am Eintrag zu. (c) Gleiche Check-Kennung `devcontainer.sandbox.runtime` mit Fallback `vfs`: bleibt über Ausnahme 2 sichtbar, nicht berührt.
3. **Planung:** [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/open/slice-v1-cli-quiet-sicherheitswarnungen.md) §1 stuft Egress „has no effect“ als sichtbar mit `SecurityRelaxation` ein und `nestedRuntime` „has no effect“ als unterdrückt (kein Kennzeichen); DoD und Dateiliste führen dieselben vier gekennzeichneten Einträge. Deckungsgleich mit dem Lastenheft-Satz.
4. **Zählung:** „Vier Ausnahmen“ unverändert (Lastenheft Zeile 189 und §16-Zeile 0.9.0; Slice §1 „die Zahl der Ausnahmen bleibt vier“); der Diff von `spec/lastenheft.md` fügt keine Zahländerung hinzu.
5. **Kollisionssuche** (`Egress`, `Konfigurationsmeldung`, `keine Wirkung`, `no effect`) im Lastenheft: Ausnahme 2 nennt den Wegfall der Egress-Restriktion als Degradations-Fallback (anderer Fall: Capability nicht gewährbar, nicht fehlendes Profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und Konfigurationsschlüssel-Text widersprechen nicht. `Konfigurationsmeldung` kommt nur im neuen Satz vor. Keine Kollision.

Sensor: `make gates` Exit-Code 0 (selbst gefahren). `make docs-check` nach dem Anhang: siehe Zusammenfassung der Nachverifikation.

**Verdikt: bestanden.** F-1 behoben, 0 offene Befunde.
