# Verification Evidence: slice-v1-devcontainer-feature-sshd — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v1-devcontainer-feature-sshd`](../plan/planning/in-progress/slice-v1-devcontainer-feature-sshd.md) (Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/welle-devcontainer-sshd-claude-code.md), dritter Slice), Commits `ff85bbd` (Code, Tests), `f829fc5` (Doku, Replay-Zeile, Messung) und `0574231` (Review-Fixes F-1/F-2/F-3/F-4). Maßstab: §1 (Ziel, Nicht-Ziele), §2 (DoD), §3 (Plan, Umsetzungsstand, Messung), §6 (Risiken), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-005`](../../spec/lastenheft.md#lh-fa-dev-005--ports), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Lastenheft 0.11.0), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md), `harness/replay.md`. Review-Report: [`2026-10-17-slice-v1-devcontainer-feature-sshd`](2026-10-17-slice-v1-devcontainer-feature-sshd.md) (0 HIGH, 2 MEDIUM, 4 INFO; lief gegen `f829fc5`, die Fixes in `0574231`). <!-- d-check:ignore (Pfad als Zitat) -->

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-feature-sshd`](../plan/planning/in-progress/slice-v1-devcontainer-feature-sshd.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-005`](../../spec/lastenheft.md#lh-fa-dev-005--ports), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md)
- Artefakte: `devcontainer_features.go`, `devcontainer_sandbox.go`, die Tests (`generate_features_test.go`, `generate_sandbox_test.go`, `devcontainer_features_test.go`, `doctor_features_drift_test.go`, `cli/quiet_test.go`), `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (§2 des Slice; die offenen Häkchen sind Closure-Arbeit):
- [x] **Katalog und Generator** — Eintrag `sshd` (Quelle `ghcr.io/devcontainers/features/sshd`, Version `1`, kein `publisher`); `egressFeatureHosts()` unverändert (keine Zeile, weil die Messung keinen Host fand; Test `…SandboxEgress_SshdAddsNoHost` vergleicht das Skript mit und ohne `sshd` byte-gleich, `…_FeatureHosts/sshd` erwartet nur `deb.debian.org`); Sandbox-Verhalten über `loginServiceActive` und `relax(...)` mit `SecurityRelaxation`; Render, double-run → NoOp, Happy/Negative/Boundary und CLI-Sichtbarkeit sind getestet (siehe (a), (b), (d)); Mutation „ohne `SecurityRelaxation`“ rot (M-1). **Formulierung „bei `init`, `generate` und `up`“:** trägt den Ist-Stand, weil die DoD-Zeile den Erzeuger ausdrücklich benennt (`generate devcontainer`) und den Test auf den Filter mit eingespeistem Eintrag begrenzt; §1 sagt es ebenso. Die Begründung steht im Code: `init` übergibt `nil`-Features, `up` erzeugt nur Neuanlege-Warnungen, `sandboxRelaxationWarnings` hat genau einen Aufrufer (`generateDevcontainer`). Eine Restunschärfe steht im Test-Doc-Kommentar („keeps it in init, generate and up“), der den Filter meint (V-3, INFO).
- [x] **Messung** — erfüllt, mit benannten Grenzen; Bewertung unter (c).
- [x] **Doku** — Handbuch (Katalogzeile, Egress-Zeile, Abschnitt „Login-Dienst `sshd` (optional)“, beide `--quiet`-Zeilen), `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` nennen den Schlüssel; gegen den Code gelesen (siehe (e) bis (g)). Der Satz „follow in a later release“ ist angepasst.
- [x] `make gates` grün — selbst gefahren, Exit 0.
- [~] Review durchgeführt, Report liegt vor — Report liegt vor (Häkchen steht noch aus); Bewertung der Findings unter (h). Der Review lief gegen `f829fc5`, `0574231` ist nicht erneut gereviewt: die Änderung ist klein (ein Satz im Warntext, Doku, ein Test), ich habe sie selbst gegen Code und Messung gelesen (siehe (h)); kein Nachreview verlangt.
- [ ] Doku-Update bei öffentlichem Vertrag — inhaltlich erledigt (siehe DoD „Doku“); Häkchen Closure.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Wird bei der Closure gefüllt“).
- [ ] Beobachtungs-Register fortgeschrieben — offen, Closure; Beleg-Kandidaten unter (j).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — Ausgang steht an allen vier, Form bei Risiko 2 nicht kanonisch (V-2).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Closure.

### (a) Erkennung und Warnung gegen den Code (`devcontainer_sandbox.go`)

- `loginServiceActive(dc)`: `nil` → falsch; für jedes Feature mit `Enabled == true` (nil und `false` zählen nicht) löst `projectFeatureSourceVersion(name, f.Source, f.Version)` die Quelle auf und vergleicht sie mit `featureCatalogue()["sshd"].source`; ein `version`-Override ändert die Quelle nicht, die Warnung bleibt. Aufruf nur unter `sandbox` (`sandbox && loginServiceActive(dc)`), und `sandbox` ist `req.Sandbox || Profil sandbox`; im Default-Profil keine Warnung, wie das Lastenheft nur das Sandbox-Profil bindet. Erzeugt über `relax(...)`: `Level: "warn"`, `SecurityRelaxation: true`, Code `LH-FA-DEV-006`. <!-- d-check:ignore (Kennung als Zitat) -->
- Einziger Aufrufer ist `generateDevcontainer` (Suche im Baum); `init` rendert ohne Feature-Konfiguration, `up` rendert nichts, `config set` ändert nur `u-boot.yaml`. Ein Container mit `sshd` entsteht also erst nach `generate devcontainer`, und jeder Lauf (auch NoOp) weist ihn aus.
- Grenze (im Slice §3 genannt): ein anderes SSH-Feature über `source:` bleibt unausgewiesen; der Zweig `f.Source` ist ungepinnt, weil ein Override ein URL-Schema braucht (`validateFeatureSource`), die Katalogquelle keines hat. Ich bestätige die Begründung; kein erreichbarer Fall im Rahmen von [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil).

**Warntext gegen Messung und Lastenheft.** Wortlaut: „login service sshd is active in the container (port 2222; password and key login, root login permitted; the accounts stay locked until you set a password with `sudo passwd <user>` or add a key); the profile limits damage and is not an isolation boundary“.

| Aussage im Text | Beleg |
|---|---|
| Port 2222 | `ss -ltn` `0.0.0.0:2222`, `[::]:2222` (Slice §3) |
| password and key login | `sshd -T`: `passwordauthentication yes`, `pubkeyauthentication yes` |
| root login permitted | `sshd -T`: `permitrootlogin yes` |
| accounts stay locked until password set or key added | `passwd -S`: `root`, `vscode` `L`; `~/.ssh` fehlt |
| `sudo passwd <user>` | setzt das Passwort des genannten Nutzers; `<user>` verhindert das Root-Passwort aus dem Wortlaut „ohne Argument“ |
| limits damage, not an isolation boundary | Zielsetzung [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Schadensbegrenzung, keine Isolation) |

Der Text ist so weit wie der Befund und nicht weiter (der Schlüssel-Login ist aus der Konfiguration abgeleitet, nicht aus einem Login: Grenze 1 der Messung; der Text sagt „login“ der Konfiguration, nicht „bestätigt“). Gegen das Lastenheft: die Ausweisung verlangt „einen im Container aktivierten Login-Dienst (Passwort-Login)“; „(Passwort-Login)“ ist Klammerzusatz und kein Ausschluss, der breitere Text erfüllt die Zusage. Das Lastenheft untertreibt den gemessenen Zustand; eine Präzisierung ist eine Entscheidung des Projektinhabers (Planner), keine Auflage dieses Slice (V-1, INFO). Sprachvertrag [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache): Meldung Englisch.

### (b) Katalog, Egress, forwardPorts, Drift

- Katalog: `featureCatalogue()` hat 12 Schlüssel (`git`, `docker-cli`, `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, `python`, `sdkman`, `sshd`, `homebrew`); `TestFeatureCatalogue_KeysCoverSpecExamples` führt `sshd`. README und README.de nennen „12“ und beide Aufzählungen mit `sshd` — stimmt; Suche nach „11 built-in“ in README, Handbuch, Wartungs-Doku und Spezifikation: keine Reste.
- Egress: kein Host (Messung), `egressFeatureHosts()` ohne Zeile; Handbuch-Tabelle „Feature `sshd` | keine zusätzlichen Hosts“ stimmt. Mutation M-3 (Host `ghcr.io` für `sshd` ergänzt) rot.
- `forwardPorts`: `TestGenerateDevcontainer_SshdLoginWarning` prüft je Fall, dass `2222` nicht in `devcontainer.json` steht ([`LH-FA-DEV-005`](../../spec/lastenheft.md#lh-fa-dev-005--ports): Quelle nur Services).
- Drift: `TestDoctor_FeaturesDrift_Sshd` (in sync / fehlender Render-Schlüssel) steht im Diff.

### (c) Messung mit der Dev-Container-CLI — Bewertung

Der Messlauf (CLI 0.89.0 im Wegwerf-Image mit Docker-Socket, Colima, arm64, Basis Debian trixie, `sshd:1` löste auf 1.1.0 auf) steht als Tabelle in Slice §3. Ich habe ihn **nicht wiederholt**: im Repo ist er nicht ausführbar (kein Make-Target, kein Skript, Docker-Socket-Image ad hoc). Die Bewertung stützt sich auf die Formulierung, nicht auf einen Nachlauf.

- **DoD-Punkte abgedeckt:** Profile `default`, `sandbox`, `sandbox` mit Egress (Zeile „Container“ und „Port“ nennen alle drei); Port 2222; Metadaten `mounts`, `capAdd`, `privileged`, `entrypoint`, `installsAfter`, Optionen (Zeile „Feature-Metadaten“); Laufzeit-Hosts (Zeile „Netz zur Laufzeit“). Der Befund steht im Slice (die DoD sagt „Evidence“; der Slice §3 ist der Ort, den die Vorgänger-Slices ebenso nutzen).
- **Vorsichtig und belegt:** ja. Die vier Grenzen sind genannt (kein echter Login; Hosts aus Build-Protokoll, nicht DNS-Mitschnitt; nur arm64/Colima und Major-Pin; OCI-Artefakt statt Repository). Einzelaussagen tragen ihre Belegart („`sshd -T`“, „`passwd -S`“, „Typ nicht geprüft“ beim `default`-Mount). Die Abweichung vom CR („nur Passwort-Login“ ist enger als der Befund) wird offen als Folgerung benannt statt geglättet.
- **Weiter als der Beleg** (zwei Stellen, nicht blockierend): „Risiko ‚Feature-Metadaten‘ entfällt“ gilt für `sshd` 1.1.0, der Pin `1` schwimmt; „Schlüssel-Anmeldung … wirkt“ folgt aus der Konfiguration, nicht aus einem Login. Beide Einschränkungen stehen in den Grenzen 1 und 3; die Closure sollte sie im Ausgang von Risiko 2 mitnehmen (V-2).
- **Nicht wiederholbar: Grenze oder Mangel?** Grenze. [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) trägt dieselbe Art Messung mit derselben Aussage („im Repo nicht ausgeführt, durch die Messung belegt“); der Slice benennt die Nicht-Wiederholbarkeit nicht ausdrücklich, aber alle seine Grenzen sind Einzelheiten dieser einen. **Auflage (klein):** die Closure-Notiz nennt ausdrücklich, dass der Messlauf im Repo nicht wiederholbar ist und was ein Nachlauf braucht (CLI-Version, Feature-Version 1.1.0, Basis-Image-Tag, die drei Profile); dann ist die Grenze benannt und der Punkt nicht still. Ein Mangel wäre es erst, wenn eine Folge-Entscheidung (z. B. Wegfall der Egress-Zeile) an dem Lauf hinge, ohne dass er benannt ist; das ist hier nicht der Fall, weil der Test `…SshdAddsNoHost` die Folge (kein Host) festhält und eine abweichende spätere Messung als Test-Änderung sichtbar würde.
- **Rückführungs-Trigger (§4) nicht eingetreten:** keine Lockerung des Profils gefunden (`Privileged=false`, `CapAdd` leer außer `NET_ADMIN` von `u-boot`, keine Devices, keine Port-Bindungen); der Befund „Schlüssel-/Root-Login“ ist keine Profil-Lockerung, sondern die Konfiguration des Dienstes, die die Warnung jetzt nennt.

### (d) Tests und eigene Mutationen (`make test`, danach `git checkout -- <Datei>`, Baum sauber)

Vorhanden: Happy (`…SshdLoginWarning`: Profil, `--sandbox`, Version-Override), Negative (Default-Profil, ohne `sshd`, `enabled: false`, `sshd: {}`), Boundary (`…_ThirdPartyStaysUnflagged`: Drittanbieter-Warnung ohne Kennzeichen daneben; CLI-Test mit zwei Einträgen unter `--quiet` in `init`, `generate`, `up`), Pins (Port, `sudo passwd`, Grenze; `key login`, `root login`, kein „only“).

| # | Mutation | Ergebnis |
|---|---|---|
| M-1 | `devcontainer_sandbox.go`: `relax(...)` der sshd-Warnung → `warn(...)` (ohne `SecurityRelaxation`) | rot: `…SshdLoginWarning` (`sandbox_profile_`, `sandbox_flag_`, `sandbox_with_sshd_version`), `…_ThirdPartyStaysUnflagged` |
| M-2 | `devcontainer_sandbox.go`: Warntext „password and key login“ → „password login only“ | rot: `…SshdLoginWarning` (drei Sandbox-Fälle) |
| M-3 | `devcontainer_egress.go`: Zeile `"sshd": {"ghcr.io"}` ergänzt | rot: `…SandboxEgress_SshdAddsNoHost` |

Die Mutationen des Reviews (Bedingung ohne `sandbox`, Quellenvergleich, `Enabled`-Prüfung) habe ich nicht wiederholt; M-1 bis M-3 sind andere Stellen. `git status --short` nach den Läufen: leer.

### (e) CHANGELOG `[Unreleased]`

`### Added`: Eintrag `sshd` (offizielles Feature, Port 2222, Passwort- und Schlüssel-Login, Root-Login erlaubt, Konten gesperrt bis Passwort oder Schlüssel, keine Optionen, kein Egress-Host; im Sandbox-Profil erlaubt; `generate` weist den Login-Dienst aus, sichtbar unter `--quiet`; Port nicht weitergeleitet). Gegen (a) und (b) gelesen: stimmt. Im `### Changed`-Absatz ist „(password and key login, root login permitted)“ nachgezogen. Der Restsatz nennt nur noch „the run-time handover of agent credentials with its new `doctor` check and the Claude Code feature“: stimmt mit dem Stand (Login-Dienst-Ausweisung und `sshd` sind geliefert; Zugangsdaten-Übergabe, `doctor`-Check und Claude Code stehen aus; Lastenheft 0.11.0 und Slices der Welle decken es).

Release-Cut-Szenario (Cut nach diesem Slice, ohne `claude-code`): der Eintrag ist im Indikativ, die `sshd`-Zeile ist vollständig lieferbar, und der Restsatz sagt, was fehlt. Die Aussage „`generate` discloses“ ist präzise (nicht „init, generate and up“). Der Cut braucht keine Glättung. Der letzte landende Code-Slice der Welle (`claude-code`) entfernt den Restsatz.

### (f) Handbuch gegen den Code

- Katalogzeile `sshd`: Port 2222, Passwort- und Schlüssel-Login, Root-Login, keine Optionen — stimmt.
- Abschnitt „Login-Dienst `sshd` (optional)“: `config set … true` und `generate devcontainer`; „im Sandbox-Profil erlaubt, Warnung auch unter `--quiet` (anders als die Drittanbieter-Warnung)“, „ohne Sandbox-Profil erscheint sie nicht“ — stimmt mit `sandbox && …` und dem Boundary-Test. „`sudo passwd` ohne Argument setzt das Passwort von `root`“ ist die Begründung für `<user>` im Text. „Port nicht veröffentlicht und nicht weitergeleitet, `forwardPorts` nur aus Services“ — stimmt. „Keine Lockerung des Profils“ — stimmt mit der Messung (nur Entrypoint). Die Aussage „gemessen mit der Dev-Container-CLI 0.89.0“ ist mit Grenze 3 vereinbar (Version genannt, Umgebung nicht; INFO).
- `--quiet`-Zeilen (Optionstabellen an zwei Stellen, textgleich): „die `sshd`-Warnung erzeugt `generate devcontainer`; die übrigen Lockerungen je nach Befehl `init`, `generate`, `up`“ — stimmt mit dem einzigen Aufrufer; die zusammengesetzte Klammer ist lang, aber richtig. Die Formulierung „Vier Ausnahmen“ bleibt (die `sshd`-Warnung ist Teil von Ausnahme 1).
- Grenzen-Absatz: Passwort-, Schlüssel-Login, `PermitRootLogin yes`, Port auf allen Schnittstellen im Container — deckt sich mit `sshd -T` und `ss -ltn`.

### (g) Welle und `open/README.md`

- Welle §4/§5: Reihenfolge CR (`done/`) → `slice-v1-cli-quiet-sicherheitswarnungen` (`done/`) → dieser Slice (`in-progress/`) → `claude-code` → Spec (beide `open/`) stimmt mit dem Verzeichnisstand überein (`ls`: `in-progress/` nur dieser Slice neben den Master-Dokumenten; `open/` vier Slices, davon zwei der Welle, zwei fremde). <!-- d-check:ignore (Kennung als Zitat) -->
- `open/README.md`: die Zeile des `sshd`-Slice steht weiter im Bestand von `open/`, mit korrektem Pfad nach `in-progress/`; gängige Praxis bis zur Closure (V-3 der Vorgänger), kein Sensor meldet es. Die Zählung („dritter“, „vierter“, „letzter Slice“) stimmt.
- **Welle-Datei gegen den Befund (Planner, Welle-Closure):** (1) §1 „sshd mit einer Warnung zum Passwort-Login“ und §6 „das Feature kennt nur Passwort-Login“ stehen im Widerspruch zur Messung (Schlüssel-Login ist in `sshd -T` an); (2) §3 Closure-Trigger „`--quiet` zeigt bei `init`, `generate` und `up` die Warnung zum Login-Dienst“ trifft nur für `generate devcontainer` als Erzeuger zu (bei `init` und `up` bleibt allein der Filter bewiesen). Beide sind Plan-Aussagen, die der Slice korrigiert, die Welle nicht (V-4, INFO, Pflicht vor Welle-Closure, nicht blockierend für diesen Slice).

### (h) Review-Findings

| Finding | Bewertung | Begründung |
|---|---|---|
| F-1 MEDIUM (Warntext „password login only“) | **behoben, nicht blockierend** | `0574231`: Text nennt Passwort- und Schlüssel-Login, Root-Login, `sudo passwd <user>`; Test pinnt „key login“, „root login“, kein „only“; M-2 rot. Text gegen Messung geprüft (siehe (a)). |
| F-2 MEDIUM (Handbuch, CHANGELOG zu eng) | **behoben, nicht blockierend** | `0574231`: Katalogzeile, `--quiet`-Zeilen, Abschnitt und CHANGELOG sagen dasselbe wie Code und Messung; Suche nach „nur Passwort“ im Handbuch und „login only“ im CHANGELOG ohne Rest. |
| F-3 INFO (Erzeuger nur `generate`) | **behoben, nicht blockierend** | Handbuch-Klammer, Slice-§1/§3 und DoD-Zeile benennen den Erzeuger; die Welle-Datei bleibt (V-4). |
| F-4 INFO (`f.Source`-Zweig ungepinnt) | **nicht blockierend** | Grenze in Slice §3 mit Begründung; kein erreichbarer Fall; Katalog bleibt die einzige Quelle der Erkennung. |
| F-5 INFO (Folgerung weiter als Beleg) | **nicht blockierend** | siehe (c); Auflage für den Ausgang von Risiko 2. |
| F-6 INFO (Spezifikation ohne `sshd`) | **nicht blockierend** | Spezifikation zieht [`slice-v2-spec-devcontainer-sshd-claude-code`](../plan/planning/open/slice-v2-spec-devcontainer-sshd-claude-code.md) nach (`spec/spezifikation.md` führt die Katalogzeile und die Egress-Tabelle weiter ohne `sshd`, laut Welle gewollt); der Spec-Slice trägt es. |

### (i) Plan-vs-Code-Diff (§3 gegen `git diff --stat 3d7501c..HEAD`, 14 Dateien)

| Datei im Diff | Plan-Zeile §3 | Befund |
|---|---|---|
| `devcontainer_features.go` | ja | konform (Eintrag `sshd`) |
| `devcontainer_sandbox.go` | ja (`devcontainer_sandbox.go`, `generate.go`, bedingt) | konform; `loginServiceActive` neu, Doc-Kommentar von `sandboxRelaxationWarnings` nennt die vierte Klasse |
| `devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `cli/quiet_test.go` | ja | konform |
| `doctor_features_drift_test.go` | ja | konform |
| `benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | ja | konform |
| Slice-Datei, Review-Report | nicht in §3 | Prozess-Artefakte (Plan-Nachzug, Review-Ausgang); kein Mangel |
| `devcontainer_egress.go`, `generate.go` | ja (bedingt) | im Diff **nicht** berührt; Begründung im Umsetzungsstand (Messung ohne Host; Aufrufer existiert) — bestätigt |
| `AGENTS.md` §3.10 | in §6 genannt, nicht in §3 | nicht im Diff; die Regel steht bereits (Stand-Datei der Beobachtung); keine Änderung nötig |

Abweichung der Plan-Aussagen (Umsetzungsstand dokumentiert): Warn-Erzeuger nur `generate devcontainer`; Messung widerlegt „nur Passwort-Login“ aus Planungsgrundlage und CR; `harness/replay.md` zusätzlich (in der Tabelle geführt). Nicht-Ziele §1 eingehalten: keine Optionen, kein Schlüssel-Login-Code, kein Eintrag in `forwardPorts`, Spezifikation unberührt (`spec/` nicht im Diff), Lastenheft unberührt.

### (j) Beobachtungs-Zähler (Belege für die Closure)

- **Plan-Aussage über Code ([`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md), §3.10): ja, trat auf**, zweimal in diesem Vorgang: (1) Planungsgrundlage und CR sagen „nur Passwort-Login, kein dokumentierter Schlüssel-Login“; `sshd -T` der Messung zeigt `pubkeyauthentication yes` und `permitrootlogin yes` (vom Implementer im Slice §3 notiert, vom Review als F-1/F-2 in Warntext und Doku nachgewiesen). (2) Plan, Welle und Doku-Entwurf nennen `init`, `generate` und `up` als Erzeuger der Warnung; der Code hat einen (`generate devcontainer`; vom Implementer im Umsetzungsstand notiert, Review F-3). Die Regel steht schon in `AGENTS.md` §3.10 (verkörpert); die Beobachtung zählt weiter, die Evidence-Datei im Register ist Closure-Arbeit.
- **Kommentar beschreibt Abwesendes ([`BEO-CLI/kommentar-beschreibt-abwesendes`](../plan/planning/observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md), §3.9): nein.** `git diff` gelesen auf (a) Chronik-Formulierungen, (b) Verweise auf geänderte Werte, (c) Helfer zwischen Doc-Kommentar und Funktion: keiner (`loginServiceActive` steht nach `sandboxRelaxationWarnings`). Grenzfall (V-3, INFO): der Doc-Kommentar des CLI-Tests sagt „keeps it in init, generate and up“, wo die Funktion nur den Filter prüft; das ist Zusage-Umfang, kein §3.9-Muster, und der Test-Kommentar im Slice §3 stellt es richtig. Kein Zähler-Eintrag.
- Review-Finding-Klassen dieses Laufs (Review, F-1 bis F-6) gehören in die Closure §7: Warntext enger als gemessenes Verhalten · Doku widerspricht Messbefund · Warnungs-Erzeuger enger als Befehlsliste · Erkennungszweig ohne Wächter · Folgerung weiter als Messbereich · Spec-Katalogliste hinter dem Code.

### (k) Risiken §6 — Form der Ausgänge

| Risiko | Ausgang im Slice | Befund |
|---|---|---|
| 1 Passwort-Login | `entfallen: DoD 1 trägt die Fassung des CR; bei der Closure bestätigen` | Form kanonisch; Bestätigung offen; Begriff „Passwort-Login“ ist enger als der Befund → in der Closure „Login-Dienst (Passwort- und Schlüssel-Login, Root-Login)“ schreiben |
| 2 Feature-Metadaten | `weiter offen: die Messung (DoD 2) entscheidet …` | **nicht kanonisch** (V-2): die Messung hat entschieden → `entfallen: gemessen (§3): kein mounts/capAdd/privileged/securityOpt, nur entrypoint/installsAfter, sshd 1.1.0, Major-Pin schwimmt` |
| 3 Port 2222 nur über Forwarding | `entfallen: …` | Form kanonisch; stimmt mit Code, Handbuch und Test |
| 4 Plan-Aussage über Code | `weiter offen: → BEO-DOK/plan-aussage-ueber-code` | kanonische Form; Eintritt belegt (j); darf als `eingetreten`-Beleg in die Closure, die Beobachtung bleibt im Register |

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (`lint`, `test`, `coverage-gate: OK — coverage 91.10% meets threshold 90%`, `docs-check`, `baseline-verify`) |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make test` (3 Mutationen) | rot wie erwartet | M-1 bis M-3, jeweils per `git checkout -- <Datei>` zurückgesetzt |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Evidence erneut gefahren (siehe Übergabe) |
| `make test-docker` | not run | kein Docker-Laufzeitpfad berührt: Katalog-Eintrag und Warnung im Generator, keine Compose-/Engine-Aufrufe |
| `make verify-depguard` | not run | keine Importregel berührt, keine neuen Imports |
| Messung mit der Dev-Container-CLI | nicht wiederholt | im Repo nicht ausführbar; Bewertung nach Formulierung, siehe (c) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Katalog-Eintrag, `TestFeatureCatalogue_KeysCoverSpecExamples`, Render-Fall `sshd`, Drift-Test |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | `loginServiceActive`, `TestGenerateDevcontainer_SshdLoginWarning*`, M-1, M-2 |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `…SandboxEgress_SshdAddsNoHost`, `…_FeatureHosts/sshd`, M-3 |
| [`LH-FA-DEV-005`](../../spec/lastenheft.md#lh-fa-dev-005--ports) | Port `2222` nicht in `devcontainer.json` (Test je Fall) |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `TestExecute_Quiet_KeepsSshdLoginWarningDropsThirdParty` (Filter in `init`, `generate`, `up`), M-1 |
| [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) | Katalog ohne `publisher`, keine Optionen, Messgrenze wie dort |

Carveouts:
- Neu: none
- Gelöst: none

## Verdikt

**Verifiziert: ja, ohne blockierende Abweichung.** Plan, DoD (die vier Liefer-Punkte) und Verträge sind getragen; die Behebungen F-1/F-2/F-3 stehen, F-4 bis F-6 sind Grenzen oder Folge-Arbeit. Die Messung ist belegt und vorsichtig formuliert; ihre Nicht-Wiederholbarkeit ist eine Grenze, solange die Closure sie nennt.

**Auflagen für die Closure (keine blockiert den Übergang, alle gehören in §7 oder die Welle):**
1. Risiko 2 in die kanonische Form bringen (V-2), Risiko 1 bestätigen und den Begriff „Passwort-Login“ an den Befund anpassen.
2. Closure-Notiz nennt: Messlauf im Repo nicht wiederholbar (CLI 0.89.0, `sshd` 1.1.0, Basis `…/base:debian`, drei Profile), Major-Pin schwimmt.
3. Beobachtungs-Register: `BEO-DOK/plan-aussage-ueber-code` bekommt die Evidence-Datei dieses Slice mit den zwei Befunden aus (j); `BEO-CLI/kommentar-beschreibt-abwesendes` bekommt keinen Eintrag („keine Beobachtung“ dort in §7 notieren).
4. Planner, vor der Welle-Closure: Welle-Datei §1/§6 („Passwort-Login“, „kennt nur Passwort-Login“) und §3 („bei `init`, `generate` und `up`“) an den Befund anpassen (V-4); optional Lastenheft „(Passwort-Login)“ bei der nächsten Berührung präzisieren (V-1; Entscheidung Projektinhaber); `claude-code` entfernt den CHANGELOG-Restsatz.
5. Häkchen der offenen DoD-Punkte erst nach Closure-Notiz, Register und Risiko-Ausgängen setzen; DoD „Doku-Update bei öffentlichem Vertrag“ ist inhaltlich erfüllt.

Restpunkte: V-1 Lastenheft untertreibt („Passwort-Login“); V-2 Form Risiko 2; V-3 Test-Doc-Kommentar „init, generate and up“ meint den Filter; V-4 Welle-Datei und `open/README.md`-Praxis (Zeile bleibt bis Closure in `open/`).
