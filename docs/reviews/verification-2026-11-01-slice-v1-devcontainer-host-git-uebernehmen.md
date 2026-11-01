# Verification-Report: [`slice-v1-devcontainer-host-git-uebernehmen`](../plan/planning/in-progress/slice-v1-devcontainer-host-git-uebernehmen.md) — 2026-11-01

**Verification-Art:** DoD-Abgleich im frischen Kontext (Baseline-Regelwerk
`modul-11-verification.md` — DoD-Konformität gegen Code und Tests je Fundstelle, die
Behauptung der Implementer-Übergabe nicht als Beleg genommen, nachgetragene Rot-Belege
je Wächter als Mutationstest am Test), gegen den Erst-Report
[`review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen.md`](review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen.md)
(HIGH 2 · MEDIUM 2 · LOW 4 · INFO 3, alle umgesetzt) und den Nachlauf
[`review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen-nachlauf.md`](review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen-nachlauf.md)
(HIGH 0 · MEDIUM 0). Review-Findings sind hier nicht Gegenstand; die Fixes wurden nur
soweit gelesen, wie der DoD sie berührt.

**Gegenstand:** `git diff 82682508..4ab02fe7` — 20 Commits des Slices ab `82682508`, 43 Dateien, 5208 Einfügungen,
67 Löschungen; letzter den Slice berührender Commit `4ab02fe7dc9c9664e21b21df487808e5068cde28`. Während der Verification landeten zwei Commits eines anderen Entwurfs (`39ba7b07`, `69f88dc8`, Slice [slice-v1-test-docker-wechselschutz](../plan/planning/open/slice-v1-test-docker-wechselschutz.md), nur `docs/plan/planning/open/`) — kein Slice-Code berührt; dieser Report belegt den Slice bis `4ab02fe7`.

**Verifier:** Agentenlauf (Rolle Verifier), 2026-11-01 · Sensor-Läufe in Docker-Modus
(Host Linux 7.0.0-42-generic `amd64`, Engine 29.9.2, lokale Engine — der Modus der
umgebenden Standard-`make gates`-Umgebung).

---

## Scope

- Slice: [slice-v1-devcontainer-host-git-uebernehmen](../plan/planning/in-progress/slice-v1-devcontainer-host-git-uebernehmen.md) (ohne Welle; Minor/`Added`)
- IDs: [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) (schärft [ADR-0028](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md))
- Artefakte: `internal/hexagon/port/driven/hostgitsource.go`, `internal/adapter/driven/git/hostsource.go`, `internal/hexagon/application/devcontainer_start.go`, `internal/adapter/driven/docker/devcontainer_launcher.go` + `devcontainer_shell.go`, `internal/adapter/driving/cli/devcontainer.go`, `cmd/uboot/main.go`, `internal/hexagon/domain/configpath.go`, `internal/e2e/devcontainer_hostgit_docker_test.go`, Tests der benannten Dateien, `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `docs/plan/adr/0040-…` + ADR-Index, Messprotokoll und Rohdaten, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md`, `docs/reviews/` (zwei Berichte)

---

## DoD-Abgleich (Slice §2)

### 1. Messung und Freigabe — **erfüllt**

- Messprotokoll [`docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/README.md`](../messungen/slice-v1-devcontainer-host-git-uebernehmen/README.md): Messstand datiert (2026-11-01, 14:17–14:33 UTC), Host/Umgebung (Linux `amd64`, Engine 29.9.2, git 2.53.0, gh 2.46.0, Dev-Container-CLI 0.89.0), alle acht Messungen (M1–M8) mit Ergebnis-Tabellen.
- Rohdaten im selben Verzeichnis: `rohdaten-m135-leicht-lauf1..3` (drei volle Läufe), `rohdaten-m2-lauf1..2` zuzüglich Vorlauf-Sonde, `rohdaten-m6-lauf1` (je Zelle **zwei Wiederholungen** in der Datei, „Wiederholungen je Zelle: 2" im Kopf), `rohdaten-m8-lauf1..2`. Die DoD-Vorgabe „mindestens zwei Läufe je Zelle" ist je Zelle belegt; bei M6 steckt die Wiederholung („Wiederholungen je Zelle: 2") in der einen Rohdaten-Datei, bei M8 tragen die zwei Läufe je eine Datei und belegen eine Zelle (Protokoll Grenzen, ehrlich benannt).
- Skripte `tools/messung/devcontainer-host-git-uebernehmen/messen-leicht.sh` und `messen-schwer.sh` (Ablage nach `harness/conventions.md` §Ablage; `make ablage-pruefung` grün).
- Keine Zugangsdaten, keine echten Namen/Adressen: eigener Musterscan über [docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/](../messungen/slice-v1-devcontainer-host-git-uebernehmen/README.md) und `tools/messung/devcontainer-host-git-uebernehmen/` (`github_pat_`, `ghp_`): kein Treffer; der echte Token steht in M3 nur als Länge (93) und SHA-256-Präfix (Rohdaten gelesen).
- Freigabe: Entscheidungen (Frage 2 Opt-in, Frage 4a Quellen 1+2 ohne Quelle 3, Frage 4b D2, Frage 6 Sandbox-Profil für Token/Identität in beiden Profilen, Fragen 1/5/9/10) im Slice §6 „Entschieden am 2026-11-01" festgehalten; der Lastenheft-CR 0.28.0 trägt in der Historie „Entscheidung des Projektinhabers 2026-11-01" (`spec/lastenheft.md`, Historien-Zeile 0.28.0). Der Wortlaut der Freigabe (Auflösung der Klammer „[freigegeben: eingeschaltet]" → „und hat der Benutzer die Auflösung in der Projektkonfiguration eingeschaltet") steht wörtlich im CR-Text von [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).

### 2. Übernahme im Container — **erfüllt im geprüften Umfang; der DoD-Nachweis „CI-Lauf" steht aus (s. untere Offene Punkte)**

**(a) Token (Opt-in, Quelle genannt nie der Wert, kein Abbruch):** je Aussage die Fundstelle.

| Zusage | Code/Beleg |
|---|---|
| Nur bei Sandbox-Profil + Schalter `devcontainer.sandbox.hostGitToken` + Referenz `${localEnv:GIT_TOKEN}` | `devcontainer_start.go`: `resolveHostGit` → `hostGitSwitch` (liest `u-boot.yaml`; `resolveToken := sandboxProfile && switchOn`), `resolveTokenSource` verlangt `containsString(names, gitTokenEnvName)` |
| Rangfolge: Host-Variable zuerst, dann `gh auth token`; `gh` fehlend/nicht angemeldet = Quelle skipped, kein Fehler | Port-Doc `hostgitsource.go` (feste Reihenfolge); Adapter `hostsource.go` `GitToken` (`os.Getenv` zuerst, `gh auth token` mit Exit ≠ 0 oder leerer Ausgabe als „keine Quelle"); Tests `TestHostSource_GitToken_HostVariableWins` / `_FallsBackToGh` / `_GhMissingReadsAsNoSource` / `_GhFailingReadsAsNoSource` / `_GhEmptyOutputReadsAsNoSource` (`hostsource_test.go`); Service-Seite `TestDevcontainerUp_HostGitToken_HostVariableWins` (`devcontainer_start_test.go`); Launcher-Seite `TestLauncher_Up_HostValueWinsOverResolvedEntry` |
| Wert nur im Prozess-Umfeld, nie Kommandozeile/Datei/Ausgabe | `devcontainer_launcher.go`: `passedEnv` (gesamte Namen, Werte als `cmd.Env = append(os.Environ(), sortedEnvPairs(envValues)…)` in `runUp` und `runSession`); Tests `TestLauncher_Up_ResolvedEnvValueTravelsInTheClientEnvironment`, `TestLauncher_Shell_ResolvedEnvValueTravelsInTheClientEnvironment` (Environment des Fake-`docker`-Clients gelesen über `/proc/$$/environ`, Wert in keiner Argumentliste), Rot-Belege unten |
| Quelle genannt, nie der Wert | Warnungen `tokenSourceHostEnvWarning` / `tokenSourceGhWarning` (`devcontainer_start.go`); `gitTokenSource` im `--json`-Umschlag (`devcontainer.go` CLI, `devcontainerUpData`); Tests `ResolvedFromGhReachesLauncher`, `WarningsCarryNoValue`; E2E `TestE2E_HostGitToken_OptInReachesTheSession` (Envelope-Suche nach dem Platzhalter) |
| Warnung ohne Abbruch bei fehlender Quelle | `tokenMissingSourceWarning`; Tests `MissingSourceWarnsWithoutAbort` (`err == nil` gepinnt), `IdentityReadFailureWarns`, `TestLauncher_Up_GitIdentityFailedExecReportsNotApplied`; E2E `TestE2E_HostGitToken_NoSourceWarnsAndCarriesNoValue` (Exit 0 bei Fake-`gh` mit Exit 1) |
| Schalter an, Verweis fehlt → Hinweis, nichts aufgelöst | `tokenSwitchWithoutRefWarning`; Test `TestDevcontainerUp_HostGitToken_SwitchWithoutReferenceWarns` (Fix-Runde, F-4a) |
| Opt-in bindet; Standardprofil löst nichts auf | Tests `SwitchOffResolvesNothing`, `StandardProfileResolvesNothing`; E2E `TestE2E_HostGitToken_SwitchOffDoesNotResolve` (Fake-`gh` wird nicht aufgerufen — `asked`-Datei fehlt) |
| Derselbe Wert für Hilfs-Container und Sitzung | `start()`-Ergebnis (`startResolution.envValues`) fließt in `Shell` (`EnvValues`); E2E OptIn-Test: Sitzung liest den Platzhalter (`printenv | wc -c` → 34) und der Credential-Helper des Images liefert ihn (`git credential fill` → `password=platzh…`, „fill-ok") |
| Für jede `--instance` | `TestDevcontainerUp_HostGitToken_InstanceStillResolves` (Fix-Runde, F-4d); M6-Zelle E (Rohdaten) und `up --instance`-E2E-Pfad derselbe Startcode |
| Token nur im Sandbox-Profil (Frage 6) | kein Referenz-Hilfscontainer-Weg im Standardprofil: `StandardProfileResolvesNothing`; M6-Zelle D |

**(b) Identität (D2, beide Profile, Container-Wert bleibt, Warnung ohne Abbruch):**

| Zusage | Code/Beleg |
|---|---|
| `git config` des Hosts im Projektordner (wirksamer Wert) | Port `ConfigValue` (`hostgitsource.go`), Adapter `ConfigValue` (`git -C <dir> config <key>`, Exit 1 = fehlender Schlüssel); Messung M5; Tests `TestHostSource_ConfigValue_*` |
| Setzen mit `git config --global` im Container-Heim, nur wo der Container keinen nichtleeren Wert trägt | `devcontainer_launcher.go` `applyGitIdentity` / `gitIdentityArgs` (`docker exec -e UB_GIT_AUTHOR_NAME … sh -c 'current=$(git config --global …); [ -n "$current" ] || git config --global …'`); Tests `TestLauncher_Up_GitIdentityExec` (Reserve gepinnt), `GitIdentityEmptySetsNothing`, E2E `TestE2E_HostGitIdentity_KeepsTheContainerValue` (zweiter Start behält den ersten Wert; `git var GIT_COMMITTER_IDENT` commit-bereit) |
| Beide Profile (Standardprofil im Test belegt) | Unit `TestDevcontainerUp_HostGitToken_IdentityFromHostConfig` (kein `u-boot.yaml` = Standardprofil) und E2E `TestE2E_HostGitIdentity_KeepsTheContainerValue` (`hostGitPrepare(ctx, t, false, …)` = Standardprofil) |
| Ohne Werte oder ohne `git` → Warnung ohne Abbruch | `gitIdentityMissingWarningT`-Familie (`resolveIdentity`); Tests `IdentityMissingWarnsWithoutAbort`, `IdentityReadFailureWarns` (F-4b, Fix-Runde), `IdentityStepFailureWarns` + Launcher `GitIdentityFailedExecReportsNotApplied` (best effort, Start bleibt) |
| Wert nie in Kommandozeile oder Ausgabe | derselbe Wert-Kanal wie der Token (`cmd.Env` des `docker exec`-Clients); `TestLauncher_Up_GitIdentityExec` prüft „h-name" nicht in `flatCalls` |

**Docker-Test mit Fake-Quelle hinter dem Build-Tag `docker`:**
`internal/e2e/devcontainer_hostgit_docker_test.go` — vier E2E-Tests (OptIn mit AnkommenToken/Hilfs-Container-Umgebung/Metadaten-Freiheit, NoSource warnt ohne Abbruch, SwitchOff löst nicht auf, Identität mit Reserve). Evidenz-Lage siehe Sensors-Abschnitt: der Gesamtlauf des Implementers vor dem heutigen Stand (`08633e7f`) lief grün inklusive der vier gefilterten HostGit-E2E-Läufe; seitdem ändert sich die Datei nur um einen Kommentar (F-2-Fix: 33/34 statt 37 — produkt-inaktiv); der DoD-Nachweis ist ausdrücklich der **CI-Lauf** und steht aus. Ein lokaler Gesamtlauf des Verifiers wurde auf Direktive des Projektinhabers nach zwei Packages abgebrochen (Regie unten); die DoD-Zusage schneidet ihn nicht als Beleg vor.

### 3. Doku und Zusage-Menge — **erfüllt** (jede Aussage über Verhalten an der Fundstelle gelesen, `AGENTS.md` §3.10)

- **Lastenheft CR 0.28.0** (Version-Zeile, Historien-Zeile): Deckungsgleichheit je Aussage:
  - „Ist die Umgebungsvariable des Hosts nicht gesetzt **und** hat der Benutzer die Auflösung … eingeschaltet … mit `gh auth token` … Quelle … nie den Wert … ohne Abbruch … in keiner Kommandozeile … nur auf einem Linux-Host (`amd64`) mit lokaler Engine und nur, soweit ein Test mit Build-Tag `docker` es belegt" ← [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ← `resolveHostGit`/`resolveTokenSource`/`hostGitSwitch` (Sandbox-Profil + Schalter + Referenz), `GitToken`-Rangfolge im Adapter, Warnungen, `cmd.Env`-Kanal; E2E-Bindung hinter `//go:build docker`.
  - [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Zusatz „Ist die Umgebungsvariable des Hosts nicht gesetzt und hat der Benutzer die Auflösung eingeschaltet, ermittelt der Befehl den Token … nennt die benutzte Quelle … weist er ohne Abbruch darauf hin" ← derselbe Code; neue Zeile „`user.name` und `user.email` ermittelt der Befehl auf dem Host mit `git config` im Projektordner und setzt sie mit `git config --global` im Container … ein im Container bereits gesetzter Wert bleibt … Wert nie in einer Kommandozeile oder einer Ausgabe, und `u-boot` öffnet keine Git-Konfigurationsdatei selbst … beide Profile" ← `resolveIdentity`, `applyGitIdentity`/`gitIdentityArgs`, Adapter liest nur `git`/`gh`-Aufrufe (keine Datei-Öffnung außer `os.Getenv`).
  - [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration): optionales Feld `devcontainer.sandbox.hostGitToken` (Schalter, ohne Sandbox-Profil ohne Wirkung) ← `configpath.go` (`ConfigDevcontainerSandboxHostGitToken`, `NewConfigPath`-Fehlerliste, `scalarConfigPaths`), Coercer und YAML-Pfad in `devcontainer_sandbox_config.go`, `HostGitToken *bool` in `initproject.go`; ohne Wirkung im Standardprofil ← `hostGitSwitch`.
- **[`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes):**
  - *Auflösung der Token-Quelle* — Schalter + Verweis + Variable nicht gesetzt; Quelle `gh auth token`; Fehler-Ende von gh „gilt wie fehlend" (Adapter: `err != nil → skip`); Übergabe wie Host-Variable (Name auf der Kommandozeile, Wert im Prozess-Umfeld); Warntexte wörtlich zitiert — die beiden Texte stehen genau so in `devcontainer_start.go` (Konstanten); Schalter ohne Verweis weist hin (`tokenSwitchWithoutRefWarning`); derselbe Wert für Start und Sitzung (`startResolution`).
  - *Wozu `u-boot` das `gh` und `git` des Hosts befragt* — nur `up` und `shell` (`start()` ist die einzige Injektionsstelle `WithHostGitSource`, Verdrahtung nur `cmd/uboot/main.go`), `doctor` liest nur Projektdateien (`doctor_sandbox.go` im Diff unverändert, kein `HostGitSource`-Griff), keine anderen Schlüssel, keine Konfigurationsdatei selbst (Adapter: nur `os/exec` und `os.Getenv`).
  - *Identität* — beide Profile, nach gelungenem Start, effective Werte mit bedingten Einfügungen, leerer Wert wie fehlend, Container-Wert bleibt (nichtleer), best effort mit Warnung, Wert im `docker exec`-Client-Umfeld; deckungsgleich mit `applyGitIdentity` (inklusive der Überschreibung eines **leeren** Container-Werts: `[ -n "$current" ] ||`).
  - *Wo der Wert nicht steht* / *Grenzen* — `docker run`/`docker exec` tragen nur den Namen; `Config.Env` nicht enthalten, Label `devcontainer.metadata` trägt die Referenz (E2E prüft beides; Messung M8 belegt die Stellen).
  - *Ausgabe von `up`* — Warnungen der Auflösung nach den Fallback-Warnungen (`mapWarningsToDiagnostics(append(raised, resp.Warnings…))` im CLI), `--quiet` unterdrückt sie (`visibleWarnings`, [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) Entscheidung 4), `--json` trägt `gitTokenSource` (CLI-Diff).
  - *shell*-Umgebung — Sitzung bekommt Namen, die auf dem Host gesetzt sind **oder die der Start ermittelt hat** (`Shell` reicht `resolution.envValues`).
- **[`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)** Pfad-Tabelle: Zeile `devcontainer.sandbox.hostGitToken` (lesbar + schreibbar, Wahrheitswert, ohne Sandbox-Profil ohne Wirkung) ← `configpath.go` + `devcontainer_sandbox_config.go` (Coercer `ParseBool`, yaml-Pfad `devcontainer/sandbox/hostGitToken`, Lesung) + `configpath_test.go` (Eintrag gepinnt).
- **[SPEC-001](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) / [SPEC-024](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen)**: beide nennen `devcontainer up` in ihren Aufzählungen ([SPEC-001](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) Enum `command`, [SPEC-024](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) Beispielaufrufe); eine Feldliste der `up`-`data` führen beide nicht — das neue Feld `gitTokenSource` ist allein in [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz `--json`) beschrieben, keine Änderung fällig; Code: das Feld existiert nur im `up`-Umschlag (`devcontainerUpData`), nicht in `down`/`egress`.
- [**ADR-0040**](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md): Status Accepted, Datum, Autor mit der Rolle Architect, `Schärft:`-Feld auf [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)/[SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), `Supersedes (teilweise)` auf die drei ADRs (die selbst unverändert bleiben — `git diff` berührt [ADR-0028](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md)/0031/0035 nicht); Entscheidungen 1–5 deckungsgleich mit Code und CR (Kanal, Rangfolge, Identität, Ausgabe, Umfang der Befragung); die vier Prosa-Slice-Verweise aus F-1 sind entfernt; ADR-Index-Eintrag nachgezogen ([`docs/plan/adr/README.md`](../plan/adr/README.md), Zeile ADR 0040).
- **Handbuch** ([`docs/user/benutzerhandbuch.md`](../user/benutzerhandbuch.md)): 4.9 Leitsatz trägt jetzt die Ausnahme („oder, mit dem Opt-in-Schalter des nächsten Abschnitts, die er auf dem Host ermittelt hat" — F-6-Fix); Abschnitt „Auflösung der Token-Quelle (Opt-in)" und „Identität" beschreiben nur das Gemessene (Quellen-Reihenfolge, Wert-Kanal, Quelle genannt nie der Wert, Warnung ohne Abbruch, Grenze auf `up`/`shell`, `doctor` sieht sie nicht, Wert gilt für Start und Sitzung); Abschnitt „Zugangsdaten (Git und Agent)" trägt Schalter, Befragungs-Grenze, Identitäts-Bullet und die erweiterte shell-Weiterleitung. Keine Kennungen im Handbuch-Neutext (§5 Regel 6 — geprüft: kein `LH-`/`SPEC-`/`ADR-`/Plan-Verweis), keine Chronik.
- **CHANGELOG**: `Added`-Eintrag mit Promised/Not-Promised scharf (Auflösung, Quelle, Identität, Transport belegt durch Docker-Tests mit Fake-Quellen; echter Push nicht gemessen; `GH_TOKEN`/`GITLAB_TOKEN` und `git credential fill` nicht zugesagt; Host-`gh` mit mehreren Konten nicht zugesagt) mit Kennungen — Muster der vorigen Einträge gewahrt.
- **Bedingung mit pinnendem Test oder gestrichen** (`AGENTS.md` §3.10, letzter Punkt): die vier aus F-4 nachgetragenen Bedingungen tragen je einen Test (`SwitchWithoutReferenceWarns`, `IdentityReadFailureWarns`, `TestLauncher_Up_MasksResolvedValue`, `InstanceStillResolves`); die übrigen Bedingungen der Spec und des Handbuchs pinnt die Tabelle oben (Rangfolge, Wert-Kanal, Quelle genannt, Opt-in bindet, Reserve, best effort, Standardprofil, Instanz, Metadaten-Grenze, `doctor`-Grenze strukturell: kein `HostGitSource` in `doctor`-Verdrahtung). Der Port-Doc-Rest aus F-8 ist zur invarianten Form behoben („The error carries only a failed environment read…").

### 4. `make gates` grün; Review; Verification; Abgleich der Plan-Tabelle — **erfüllt** (Verification = dieser Report)

- `make gates` vom Verifier **frisch gefahren** auf `4ab02fe7` (nicht nur Stempel geprüft): grün — lint-Stage, test-Stage, coverage-Stage mit `COVERAGE_THRESHOLD=90`, `d-check: 999 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK — 54 Dateien`, `ablage-pruefung` OK, kein `FAIL` im Log; `record-gates` hat die Stempel erneuert (`.harness/state/gates-passed.head` = `4ab02fe7…`); nach den zwei späteren Entwurfs-Commits ist der Stempel nicht mehr die Spitze von main — der Gates-Lauf für die Closure wiederholt sich am Closure-Commit.
- Review: Erst-Report (HIGH 2 · MEDIUM 2 · LOW 4 · INFO 3) und Nachlauf (HIGH 0 · MEDIUM 0) liegen unter `docs/reviews/`; alle Fix-Commit-Zusammenhänge gelesen (`e5a8cb63` Fix-Runde, `4ab02fe7` Nachlauf-Fixes); F-N1 (Abgleich der Plan-Tabelle nach der Fix-Runde) ist in `4ab02fe7` nachgeholt — die Abgleichs-Notiz in Slice §3 trägt jetzt `spec/architecture.md` und die beiden Review-Berichte, und sie stimmt mit `git diff 82682508..4ab02fe7 --stat` überein (43 Dateien: jede berührte Datei hat einen Eintrag; Vorlagen und Golden-Cases unberührt — keine `internal/hexagon/application/templates/`-Zeile im Diff, Replay-Pflicht nicht berührt).

---

## Nachgetragene Rot-Belege je Wächter (Mutationstests des Verifiers)

Verfahren je Baseline-Regelwerk `modul-11-verification.md` §Bewusstes Brechen: Mutation
kurz angewandt, mit der beim Wächter behaupteten Meldung rot gekommen, zurückgenommen;
Baum nach der Rücknahme identisch mit dem Eingangsstand, `git status` leer, HEAD zum
Zeitpunkt der Mutation `4ab02fe7`. Läufe: **zwei Läufe von `make test`** (Docker-Modus,
Test-Stage; Unit-Suite) — der erste Lauf kam für drei Anwendungs-Mutationen nicht zur
Ausführung (Compile-Fehler der Mutation „Schalter ignoriert": `switchOn` unbenutzt), der
zweite Lauf nach Korrektur (`_ = switchOn`, gleichwertige Mutation) deckte alle acht.

| # | Wächter (Plan §3) | Mutation | Test rot mit behaupteter Meldung |
|---|---|---|---|
| 1 | Werte nie in der Argumentliste (`runUp`) | `runUp` hängt die aufgelösten Paare (`sortedEnvPairs`) an die `docker run`-Argumentliste an | `TestLauncher_Up_ResolvedEnvValueTravelsInTheClientEnvironment` → „the resolved value is in a command line" |
| 2 | Werte nie in der Argumentliste (Sitzung) | `runSession` hängt die Paare an die `docker exec`-Argumentliste an | `TestLauncher_Shell_ResolvedEnvValueTravelsInTheClientEnvironment` → „the resolved value is in a command line" |
| 3 | Maskierung der aufgelösten Werte | die Aufnahme der Werte in die Maskierliste von `runUp` entfernt | `TestLauncher_Up_MasksResolvedValue` (beide Subfälle) → „the resolved value is in the progress" bzw. „… in the error" und „shows no placeholder" |
| 4 | Ein im Container gesetzter Wert bleibt (D2) | die Reserve `[ -n "$current" ] \|\|` aus `gitIdentityArgs` entfernt | `TestLauncher_Up_GitIdentityExec` → „the exec command does not keep the container value" |
| 5 | Rangfolge: gesetzte Host-Variable hat Vorrang | die Umgebungsquelle im Adapter umgangen (`value != "" && value == "zz-impossible-zz"`), damit die gh-Quelle bei gesetzter Variable gewinnt | `TestHostSource_GitToken_HostVariableWins` → `source = "gh auth token", want the environment source` |
| 6 | Quelle genannt, Wert nicht | die Quellen-Warnung um den Wert erweitert | `TestDevcontainerUp_HostGitToken_ResolvedFromGhReachesLauncher` → „warning … carries the value"; `TestDevcontainerUp_HostGitToken_WarningsCarryNoValue` → „warning carries the value" |
| 7 | Opt-in-Schalter bindet | der Schalter ignoriert (`resolveToken := sandboxProfile`) | `TestDevcontainerUp_HostGitToken_SwitchOffResolvesNothing` → „switch off resolved something: EnvValues=map[GIT_TOKEN:gh-token-placeh] GitTokenSource=\"gh auth token\"" |
| 8 | Warnung ohne Abbruch (Warn-Klasse) | `warningNoAbort` erzeugt `Level: "error"` statt `"warn"` | `TestDevcontainerUp_HostGitToken_MissingSourceWarnsWithoutAbort` → „want the missing-source warning" (die Einträge tragen sichtbar `Level:error`); `TestDevcontainerUp_HostGitToken_IdentityMissingWarnsWithoutAbort` → „want the identity warning" |

Dazu der im Slice §3 geführte Rot-Beleg des Implementers (mutation „Werte an die
Argumentliste anhängen" → `TestLauncher_Up_PassesEnvByNameOnly`, Meldung „a value
reached a command line") — vom Verifier inhaltlich übernommen; die Mutationen 1 und 2
oben decken dieselbe Verstoßklasse an denselben Stellen frisch nach.

**Grenzen der nachgetragenen Belege** (ehrlich benannt):

- Mutation 8 zeigt die Eskalation der Warn-Klasse; dass die Warnung die Ausführung nicht beendet, pinnt zusätzlich das `err == nil` derselben Tests und Exit 0 im E2E-Fall (`TestE2E_HostGitToken_NoSourceWarnsAndCarriesNoValue`). Eine Mutation, die die Warnung in einen echten Fehler-Abbruch des `Up`-Aufrufs umwandelt, wurde nicht gebaut — der Schadenpfad wäre derselbe Test.
- Nicht mutiert, mit Grund:
  - **Docker-Wächter als Mutation** („die Auflösung auslassen" → E2E rot): nicht gefahren — auf Direktive des Projektinhabers läuft der Verifier heute keine Docker-Suite; die vier E2E-Tests messen das Ankommen direkt (Sitzungslänge, Credential-Helper, `git var`), und die Auflösungs-Strecke hinter ihnen ist dieselbe, deren Unit-Wächter (Mutationen 1, 6, 7) rot kamen. Der Beweis in der DoD-eigenen Form bleibt der CI-Lauf.
  - **Wächter „Fehler von `gh` als Abbruch durchreichen" / „Exit 1 von `git config` als Fehler"**: Korrektheit der Fehlerklassen, nicht sicherheitskritisch; die Antwortformen sind je durch `TestHostSource_GitToken_GhFailing…`/`_GhMissing…` und `TestHostSource_ConfigValue_OtherFailureIsAnError` gepinnt.
  - **Wächter „Schritt im Standardprofil/bei Instanz überspringen"**: die Rot-Richtung (Auflösung außerhalb der Freigabe) ist von Mutation 7 abgedeckt; die Positiv-Fälle sind gepinnt (Standardprofil/Instanz-Tests).

---

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | **pass** (frisch vom Verifier auf `4ab02fe7`) | Log: lint/test/coverage-Stages grün (`COVERAGE_THRESHOLD=90`), `d-check: 999 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK`, `ablage-pruefung: OK`, `MAKE_EXIT=0`; Stempel `.harness/state/gates-passed.head` = `4ab02fe7` erneuert |
| `make test` (mutiert, 2 Läufe) | **rot je Wächter** wie oben | beide Läufe des Verifiers im Docker-Test-Stage; die Meldungen stehen wörtlich in der Tabelle oben, damit die Evidenz im Repo trägt; die Roh-Logs der Läufe sind Session-Log und kein Beleg (keine `/tmp`-Referenz, wie im Erst-Review F-10 kritisiert); Mutationen zurückgenommen, `git status` leer |
| `make test`(unmutierter Inhalt) | **pass** | der unmutierte Lauf ist Teil des grünen `make gates` (Test-Stage mit `--no-cache-filter test`) |
| `make test-docker` | **nicht als Verifier-Frischbeleg geführt** | Implementer-Gesamtlauf `08633e7f` grün inklusive der vier gefilterten HostGit-E2E-Läufe gegen die Engine (Handoff-Evidence); die E2E-Datei ändert sich seitdem nur um den produkt-inaktiven Kommentar (F-2); der Verifier-Lauf wurde auf Direktive des Projektinhabers nach zwei Packages abgebrochen (kein Beweis, keine Reste — s. Regie); **der DoD-Nachweis ist der CI-Lauf**, Workflow `.github/workflows/integration.yml` (`make test-docker` mit echtem Daemon), steht aus |
| `make docs-check` | **pass** | innerhalb `make gates` (d-check 0 Befunde) |
| `make verify-depguard` | not run | keine Importregel berührt: `.golangci.yml` unverändert, Paketgrenzen unverändert (neuer Port im `internal/hexagon/port/driven`, Adapter in `internal/adapter/driven/git`, Anwendung kennt nur den Port, Verdrahtung nur `cmd/uboot/main.go`) — Grenzen gelesen statt gefahren, wie im Erst-Review |
| `make modus-pruefung` | not run | kein Makefile-/Dockerfile-Griff im Diff; der Modus der Läufe ist der Docker-Standardmodus (ohne `BUILD_MODE`), wie die Logs zeigen (`docker build … --target test`) |

**Carveouts:** unverändert — kein neuer Eintrag; `docs/plan/planning/in-progress/carveouts.md` im Diff nicht berührt.

**Traceability:**

| ID | Beleg |
|---|---|
| [LH-FA-DEV-009](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Adapter/Service-Tests + E2E-Docker-Tests (OptIn/NoSource/SwitchOff), CR 0.28.0, [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md), Mutations-Rot-Belege 1/2/3/5/6/7 |
| [LH-FA-DEV-010](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Identitäts-Tests (Reserve, best effort, beide Profile, E2E), CR 0.28.0, [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md), Mutations-Rot-Beleg 4 |
| [LH-FA-CONF-002](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) | `configpath.go`/`configpath_test.go`, `devcontainer_sandbox_config.go`, `initproject.go`, [SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeile, `u-boot config`-Pfad im Fehlertext von `NewConfigPath` |
| [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) | Absatz-Abgleich oben (je Fundstelle), Hist-Zeile in der Spezifikation |
| [SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) | Pfad-Tabellen-Zeile ↔ Code |
| [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) | Deckungsgleichheit Entscheidungen 1–5 mit Code, Index-Eintrag, keine abwärts-Verweise |

**Nicht ausgeführt:**

- `make test-docker` (Verifier) — Direktive des Projektinhabers: kein weiterer Gesamtlauf; Evidenz siehe oben.
- `make verify-depguard` — Grund oben.
- Restliche Wächter-Mutationen — Gründe im Abschnitt Rot-Belege.

**Commit / Artefakt:** `4ab02fe7dc9c9664e21b21df487808e5068cde28` (20 Commits ab `82682508`); Gates-Stempel frisch auf diesem HEAD erneuert — die Spitze von main liegt seither bei den zwei Entwurfs-Commits (`69f88dc8`), die den Slice nicht berühren; der Closure-Lauf stampt den Closure-Commit.

---

## Regie des abgebrochenen lokalen `test-docker`-Laufs (Projektinhaber-Entscheidung)

- Der lokale Gesamtlauf des Verifiers wurde nach dem Start (Build der `test-docker-tools`-Stage, erster Package-Teil) auf Direktive abgebrochen; `make` und der Werkzeug-Container (`u-boot:test-docker-tools`, `keen_wilson`) sind entfernt. Der Lauf hat die E2E-Phase nicht erreicht: kein einziger Restcontainer mit `devcontainer.local_folder`-Label aus diesem Lauf (geprüft per Label-Liste — kein Projektordner unter `/tmp/Test*`). Die Engine trägt durch den Abbruch keine Wegwerf-Objekte dieses Laufs.
- Die Messungs-Restbestände aus F-11 (vorheriges Inventur-Ergebnis des Verifiers: 42 Container mit `devcontainer.local_folder`-Label, 16 `ubmess`/`mubmess`-Volumes, Hilfs-Container) sind vom Verifier im Skript-Scope aufgeräumt: Container mit nicht mehr existierenden Projektordnern (41), alle `ubmess`/`mubmess`-Volumes, `u-boot-devcontainer-up-*`-Hilfs-Container. Schluss-Bestand: **0** Container mit dem Label, **0** Messungs-Volumes, **0** Hilfs-Container. F-11 ist damit erledigt.
- Der verbleibende Docker-Nachweis: **CI-Lauf** `make test-docker` (`.github/workflows/integration.yml`) auf dem gepushten Stand — ausstehend, Zuständigkeit Projektinhaber (Slice §5 Closure-Trigger: „Docker-Test im CI-Lauf grün").

---

## Offene Punkte (mit Zuständigkeit)

1. **CI-Lauf `make test-docker` auf dem gepushten HEAD** — der DoD-eigene Nachweis für „Token/Identität im Container" (und Closure-Trigger). Zuständig: Projektinhaber (Push und CI-Lauf). Nicht verweigert, ausstehend.
2. **F-9** (Bestands-Zeile in `docs/plan/planning/open/README.md` beschreibt den Slice als Entwurf) — Nachzug bei der Closure; Zuständig: Implementer (`git mv` nach `done/` zieht die Zeile mit).
3. **Echter Push mit dem echten Token** — vom CR ausdrücklich nicht zugesagt; Nachweis-Zuständig: Validator/Projektinhaber (Messprotokoll §Grenzen, ehrlich benannt).
4. **Slice §7 Closure-Notiz mit Steering-Loop-Lerneintrag und drei Paarungen** — Zuständig: Implementer bei Closure; nicht Teil dieser Prüfung.

---

## DoD-Fazit

| DoD-Punkt | Stand |
|---|---|
| 1. Messung und Freigabe | **erfüllt** |
| 2. Übernahme im Container | **erfüllt** — Code-Schicht vollständig belegt (Tests je Bedingung, acht Rot-Belege am Ort); die DoD-eigene Nachweisform „CI-Lauf" steht aus (Offener Punkt 1) |
| 3. Doku und Zusage-Menge | **erfüllt** (jede Bedingung gepinnt oder gestrichen; Handbuch ohne Kennungen) |
| 4. `make gates` grün | **erfüllt** (frischer Lauf des Verifiers auf `4ab02fe7`); Review umgesetzt; Verification = dieser Report |

**Schluss:** Closure-reif — **ja, mit einer ausstehenden Bedingung**: Der Slice ist
closure-reif, sobald der CI-Lauf `make test-docker` auf dem gepushten Stand grün ist
(zusammen mit den Closure-Nachzügen F-9 und der Closure-Notiz). Alle übrigen DoD-Punkte
sind an Ort und Stelle belegt; ein DoD-Verstoß wurde im frischen Kontext nicht gefunden.