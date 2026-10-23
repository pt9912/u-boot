# Validation-Bericht: welle-forge-cli-gh-glab — 2026-10-23 <!-- d-check:ignore (Wellen-Kennung, Welle liegt flach unter planning/ und wandert bei Closure) -->

**Rolle:** Validator (Modul 8) — Frage „Bauen wir das Richtige?“ gegen den realen Nutzerbedarf, nicht gegen die DoD. Kein Review, keine Verifikation, keine Reparatur. Die grünen Verifier-Berichte der Slices sind Eingang, nicht Ergebnis. Muster: [`2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung`](2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung.md).

**Gegenstand:** die Welle [`welle-forge-cli-gh-glab`](../plan/planning/done/welle-forge-cli-gh-glab.md) für `v0.15.0`, vor dem Release-Cut: die Features `gh` und `glab` im Devcontainer ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)); Handbuch Abschnitt „Forge-CLIs (optional)“ (4.11), Katalogzeilen in 4.10, Egress-Tabelle, „Zugangsdaten“; `CHANGELOG.md` `[Unreleased]`. Geprüft wurde der Nutzerpfad nach Welle-Closure-Punkt 4, am gebauten Binary.

**Bedarf:** Ein Nutzer mit nur Docker auf dem Host schaltet `gh` und/oder `glab` per `u-boot config set` ein, bekommt einen Container mit funktionierendem Werkzeug, erfährt vorher, was an der Lieferkette ungeprüft bleibt, und kann sich anmelden, ohne dass `u-boot` ein Token speichert. Im Sandbox-Profil hinter der Egress-Liste muss der Aufruf gehen. Wo er mit der Zusage nicht weiterkommt, muss das Handbuch es ihm sagen.

**Urteil: GO für den Release-Cut (release-reif: ja), mit Auflagen als Text.** Was die Welle zusagt, trägt am Binary: Installation, Version, Ausweisung, Egress-Hosts, Verhalten ohne Variable, `doctor`. Zwei MEDIUM-Befunde betreffen, was ein Nutzer darüber hinaus braucht (Standardprofil: kein gangbarer Weg im Handbuch genannt; `config set` auf `version` sperrt das Projekt). Kein Fall „Verifikation grün, Validation rot“, keine Prozess-Drift. **Nicht geprüft wurde ein gültiges Token** (siehe unten).

## Validation Evidence

### Messumgebung und Verfahren

- **Host:** Linux x86_64 (`amd64`), native Docker Engine 29.9.2, Netz vorhanden. Kein Host-Go, kein Host-Python.
- **Binary:** `make build-binaries IMAGE=u-boot-validator PLATFORMS=linux/amd64 VERSION=0.15.0-validator` aus HEAD `68584349` (Docker-only).
- **Wegwerf-Projekte** im Scratchpad: `std` (`init --devcontainer`), `sbx` (`init --devcontainer --sandbox`), `eg` (Kopie von `sbx` mit `devcontainer.sandbox.egress.enabled true`), `dr` (Kopie von `sbx` ohne die beiden Token-Referenzen). Eingabe wie ein Nutzer: `u-boot config set devcontainer.features.gh.enabled true`, ebenso `glab`, dann `u-boot generate devcontainer`. Container mit `u-boot devcontainer up`, Befehle über `u-boot devcontainer shell` (Eingabe per Standardeingabe), `docker exec` nur für den Vergleich im Standardprofil.
- **Token:** Es gibt keins; es wurde keines verwendet oder abgelegt. Die Aufrufe mit Variable liefen mit den **Platzhaltern** `ghp_PLATZHALTER` und `glpat-PLATZHALTER`.
- **Aufgeräumt:** alle Devcontainer der Läufe (Label `devcontainer.local_folder` unter dem Scratchpad), das Workspace-Volume `eg-…`, alle Images `vsc-…` und `vval-…`, `u-boot-validator:binary-archive`. Fremde Container unberührt.

### Ergebnis je Prüfpunkt

| # | Prüfpunkt | Urteil | Beleg |
| --- | --- | --- | --- |
| 1 | `generate devcontainer` mit `gh` und `glab`, Standard- und Sandbox-Profil | **GO** | Exit 0. Je Feature eine Warnung auf der Fehlerausgabe, in beiden Profilen: `gh` „checked by its SHA-256 only, no signature was found for the archive (by search pattern, not exhaustive), and the GitHub attestation the project publishes is not verified by u-boot, so the check does not protect against a release that was harmful from the start“; `glab` ohne den Attestierungs-Teil. Das Dockerfile trägt `USER root` im Standardprofil, `ARG`-Version und -Prüfsummen an einer Stelle, `sha256sum -c` vor `tar`. `devcontainer.json` des Standardprofils bleibt ohne `remoteEnv`; im Sandbox-Profil stehen `GH_TOKEN` und `GITLAB_TOKEN` als `${localEnv:…}` neben `GIT_TOKEN`. Zweiter Lauf: „already up to date“, die Warnungen erscheinen wieder; `--quiet`: still; `--json`: Status `warn`, beide Meldungen als Diagnostics |
| 2 | Bau und Version im Container | **GO mit Grenze** | Beide Profile: `gh version 2.103.0`, `glab 1.124.0`, `/usr/local/bin`, Benutzer `vscode`. **Grenze:** Die Layer lagen im Docker-Cache (Bau in 1 bis 2 s); der Download und die Prüfsumme wurden in diesem Lauf nicht frisch gefahren (Kalt-Bau nicht gemessen) |
| 3 | Sandbox hinter der Egress-Liste | **GO** | Das Sandbox-Profil schränkt **nicht von selbst** ein: `egress.enabled` ist Voraussetzung (`u-boot egress list` meldet „not active“; das Handbuch sagt es). Mit `egress.enabled true`: Liste enthält `gitlab.com` (Herkunft `feature:glab`) und `api.github.com`, `github.com` (`base, feature:gh`), wie das Handbuch erklärt. Im Container: `api.github.com` 200, `gitlab.com` 301, `example.org` und `bitbucket.org` blockiert. `glab api projects/gitlab-org%2Fcli` gelingt (ohne Variable), `gh api` erreicht den Server (Platzhalter: `Bad credentials`). Host aus der Liste genommen (`egress-init.sh` von Hand): `glab` scheitert sofort mit `dial tcp: lookup gitlab.com on 127.0.0.1:53: server misbehaving` — lesbar, aber ohne Hinweis auf die Egress-Liste (Handbuch nennt `lookup gitlab.com`) |
| 4 | Token-Variable gesetzt und ungesetzt | **GO** | Ohne Variable (Sandbox, Host-Variable ungesetzt): im Container `unset`; `gh api`: Exit 4 mit „please run: gh auth login … populate the GH_TOKEN“; `gh auth status`: Exit 1; `glab api` auf öffentliches Projekt: Exit 0; `glab auth status`: 401. Das deckt auch den im Handbuch als „nicht gemessen“ genannten Fall der **auf dem Host ungesetzten** Variable ab: sie kommt nicht als leere Zeichenkette an, sondern fehlt. Mit Platzhalter: Variable im Container gesetzt; `gh api`: `Bad credentials`, Exit 1; `gh auth status`: „Failed to log in … using token (GH_TOKEN)“; `glab api`: `401 Unauthorized`; `glab auth status` nennt den Token aus `GITLAB_TOKEN`. Das stimmt mit den Handbuch-Zeilen überein |
| 4b | `u-boot doctor` | **GO** | Sandbox mit Referenzen: `devcontainer.sandbox.credentials` ok („passed through GH_TOKEN, GITLAB_TOKEN (reference only)“). Referenzen aus `devcontainer.json` entfernt: Warnung „an active agent or forge-CLI feature needs its credential at runtime but devcontainer.json has no reference for: GH_TOKEN, GITLAB_TOKEN“, Abhilfe `u-boot generate devcontainer`, Exit 0 (`warn`), auch `--json`. Standardprofil: Prüfung übersprungen („Sandbox profile not active“) |
| 5 | Standardprofil ohne Token-Referenz | **Handbuch ehrlich, Nutzerpfad dünn** (M-1) | siehe unten |
| 6 | Handbuch, Fehlerfall, Stolpersteine | **GO mit Befunden** | siehe M-2, L-1 bis L-3 |

### Punkt 5 im Einzelnen

Standardprofil (`std`), `GH_TOKEN` und `GITLAB_TOKEN` auf dem Host gesetzt (Platzhalter):

- `u-boot devcontainer shell`: im Container **beide Variablen ungesetzt**; `gh api` verlangt `gh auth login`. Das ist konsequent (keine Referenz, keine Durchreichung) und im Handbuch-Spiegelstrich „Standardprofil, keine Token-Referenz“ beschrieben.
- `docker exec -e GH_TOKEN -u vscode <Container> …`: Variable gesetzt, `gh api` erreicht GitHub (`Bad credentials`). Ohne `-e`: ungesetzt. Dieser Weg ist im Abschnitt „Zugangsdaten“ und in 4.9 für `GIT_TOKEN` beschrieben, im Abschnitt „Forge-CLIs“ aber **nicht** für `GH_TOKEN` und `GITLAB_TOKEN` genannt; dort steht „ein bestimmter Weg ist weder zugesagt noch gemessen“.
- `gh auth login` im Container (interaktiv) wurde nicht gefahren; ein Token landet dann in `~/.config/gh` im Container und geht mit dem Container verloren. Weder die Eignung noch diese Folge steht im Handbuch.

## Lücken zwischen Bedarf und Zusage

Kein Vertragsfehler: die Zusage der Welle stimmt in Verhalten, Hilfe, Meldungen und Handbuch überein.

- **M-1 (MEDIUM): Das Standardprofil hat im Handbuch keinen Weg zum Token.** Ein Nutzer im Standardprofil aktiviert `gh`, setzt `GH_TOKEN` auf dem Host, startet `u-boot devcontainer shell` und bekommt `gh auth login` verlangt — ohne Meldung, ohne `doctor`-Hinweis. Der Satz im Spiegelstrich „Token-Variablen“ („`u-boot devcontainer shell` öffnet die Shell mit den beiden Variablen, soweit sie auf dem Host gesetzt sind“) und der gleichlautende Satz in „Zugangsdaten“ lesen sich im Standardprofil falsch; der ehrliche Spiegelstrich steht erst später. Der gemessene Weg `docker exec -e GH_TOKEN` fehlt im Forge-Abschnitt. Der Befund erfüllt das Muster [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../plan/planning/observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md): der Vertrag sagt „weder zugesagt noch gemessen“, der Nutzer braucht einen Satz, was er tun kann. Vorschlag (Text): im Spiegelstrich „Standardprofil“ `docker exec -e GH_TOKEN` als gemessen mit Platzhalter nennen (Token gültig: nicht gemessen), den `shell`-Satz mit „im Sandbox-Profil“ einschränken, und sagen, dass ein Token aus `gh auth login` im Container bleibt.
- **M-2 (MEDIUM): `config set` auf `devcontainer.features.gh.version` sperrt das Projekt.** `u-boot config set devcontainer.features.gh.version 2.0.0` meldet Erfolg und schreibt die Datei. Danach lehnen `generate devcontainer` (Exit 10), `doctor` (`uboot.yaml.valid` Fehler) **und jedes weitere `config set`** ab („schema validation failed … a version override needs a source: override“); ein `config unset` gibt es nicht. Die Meldung von `generate` beginnt mit „managed block missing or malformed:“, was nicht stimmt. Abhilfe ist nur das Bearbeiten von `u-boot.yaml` von Hand (hier geprüft; die Meldung sagt es nicht). Das Handbuch sagt „`version` ohne `source` weist `u-boot` mit Exit-Code `10` ab“ — für `generate` richtig, für `config set` nicht. Ein Nutzer, der die Version ändern will (naheliegend: „andere Version“), läuft in diesen Pfad. Folge-Slice (`config set` prüft wie `generate`) oder mindestens ein Handbuch-Satz mit der Abhilfe.
- **L-1 (LOW): Auffindbarkeit und Lesbarkeit.** Der Abschnitt „Forge-CLIs (optional)“ liegt in Kapitel 4.11 (Sandbox-Devcontainer) obwohl die Features in beiden Profilen wirken; ein Standardprofil-Nutzer findet ihn über die Link in der Katalogzeile von 4.10. Der Abschnitt umfasst rund 110 Zeilen Spiegelstriche mit Messvorbehalten; die Antwort auf „wie melde ich mich an?“ ist über vier Spiegelstriche verteilt. Eine Zwei-Satz-Einleitung mit dem Rezept je Profil würde das schließen.
- **L-2 (LOW, Wortlaut): Egress-Fehler ohne Hinweis.** Fehlt `gitlab.com` in der Liste, meldet `glab` `lookup gitlab.com … server misbehaving`. Das ist das Verhalten des Werkzeugs (Handbuch nennt es); `u-boot` kann es nicht ändern, aber ein Satz „bei `lookup … misbehaving` den Host in `egress.allow` eintragen“ stünde dem Nutzer bei.
- **L-3 (LOW, Wortlaut): die Warnungen.** Zwei Sätze von rund 400 Zeichen mit Präfix „Warning ([`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte))“ (Kennung in der Ausgabe, wie bei der Socket-Ausweisung, Bestand) und sie erscheinen bei jedem `generate`, auch ohne Änderung. Das ist so zugesagt (Handbuch: „je Feature eine Meldung in jedem Profil“), nur laut.
- **I-1 (INFO): Sandbox ohne Egress.** Wer „Sandbox-Profil“ liest, erwartet eine Beschränkung; sie ist nur mit `egress.enabled`. Handbuch und `egress list` sagen es.
- **I-2 (INFO): Grenzen der Nutzer-Prüfung.** Siehe unten.

## Was ohne gültiges Token ungeprüft blieb

Alle Aufrufe mit Variable liefen mit Platzhaltern und endeten in 401 / `Bad credentials`. **Nicht gemessen:** ein erfolgreicher authentifizierter Aufruf von `gh` oder `glab`; die Hosts, die nach einer Anmeldung folgen (Handbuch führt sie als nicht aufgenommen); `gh pr create`, `glab mr create`, `glab repo clone` (Standardprotokoll `ssh`, im Container kein Schlüssel); `gh auth setup-git` und die Wirkung auf `git`; `gh auth login` interaktiv; GitHub Enterprise und GitLab Self-Managed; ein Terminal eines Menschen; `arm64`; ein Kalt-Bau mit frischem Download (Cache-Treffer); Docker Desktop, Colima, macOS. Die Zusage des Handbuchs („Ein Aufruf mit einem gültigen Token ist nicht gemessen und nicht zugesagt“) deckt das ab; sie ist ehrlich.

## Sensors

`make gates`: **Exit 0** (alle inneren Gates, `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge; Bericht danach angelegt, Gates vor dem Commit des Berichts gelaufen). Nicht gelaufen: `make test-docker` und `make verify-depguard` (Verifier-Aufgabe der Slices, hier kein Code geändert); die Docker-Läufe dieses Berichts sind die Nutzerpfad-Messung, nicht die Integrationstests.

## Release-Reife aus Nutzersicht

- **Was trägt:** Installation und Version in beiden Profilen, ehrliche Ausweisung der Lieferkette, Egress-Hosts stimmen mit der Erklärung, Verhalten mit und ohne Variable wie beschrieben, `doctor` warnt und nennt die Abhilfe, Handbuch ist in den gemessenen Aussagen belegt und nennt seine Grenzen.
- **Auflagen an den Cut (Text, kein Code):**
  - **A-1:** `CHANGELOG.md` / Lead nennen die Zusage-Menge (Linux `amd64`, lokale Engine) und dass **kein authentifizierter Aufruf** gemessen ist.
  - **A-2:** M-1: Handbuch-Satz zum Standardprofil (`docker exec -e`, Einschränkung des `shell`-Satzes). Kann im Cut-Slice oder als kleiner Handbuch-Slice laufen.
  - **A-3:** M-2: mindestens der Handbuch-Satz „bei einer ungültigen `version` ohne `source` sperrt auch `config set`; Abhilfe: `u-boot.yaml` von Hand berichtigen“, sonst Folge-Slice mit Ausgang im Register.
- **Folge, nicht Teil des Cuts:** `config set`-Validierung (M-2), Gliederung des Abschnitts (L-1), Wortlaute (L-2, L-3).

## Bewertung

Verifikation grün, Validation grün mit Auflagen: kein Fall „perfekt das Falsche gebaut“. Das Gebaute ist das Richtige für den Nutzer, der die Features im Sandbox-Profil mit durchgereichter Variable nutzt; für den Nutzer im Standardprofil trägt die Welle das Werkzeug, aber nicht den Weg zur Anmeldung (M-1). Keine Prozess-Drift.
