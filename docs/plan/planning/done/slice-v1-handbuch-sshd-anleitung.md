# Slice v1: Handbuch — Schritt-für-Schritt-Anleitung für den `sshd` im Sandbox-Container

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `v0.11.0` (Release-Zuordnung des Projektinhabers vom 2026-10-18; Release-Versionen sind Wellen, [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)). Der Slice läuft vor dem Cut [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md); der Cut nimmt den Stand des Handbuchs mit.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `sshd`),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sicherheitsaussage: das Profil ist Schadensbegrenzung, keine Isolationsgrenze),
[`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md). Der Slice ändert **keine Zusage**: Er ordnet vorhandene, belegte Aussagen des Handbuchs zu einer Anleitung.

**Berührte Spec-Stellen:** — (das Handbuch ist ein öffentlicher Vertrag, kein Spec-Stratum; Lastenheft und Spezifikation bleiben unverändert).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Wer den `sshd` im Sandbox-Container nutzen will, findet im Handbuch eine Anleitung in der Reihenfolge, in der er sie ausführt: Feature aktivieren, Zugang einrichten, verbinden, Sicherheitshinweis.

### Entscheidung des Projektinhabers (2026-10-18)

„Das Handbuch können wir dazu verbessern, sonst nichts.“ Folge: Der Slice schreibt Handbuch-Text und sonst nichts — keine neue Zusage im Lastenheft, kein Feature, keine Änderung am Code; die Nutzer-Doku trägt keine Anforderungs-, ADR- oder Spec-Kennungen und keine Chronik, sie nennt den Stand ([`AGENTS.md`](../../../../AGENTS.md) §5, Zeile 6).

### Gefundener Stand (Ist-Stand, gelesen am 2026-10-18)

- Das Handbuch hat den Abschnitt „Login-Dienst `sshd` (optional)“ als vierte Überschriftsebene in 4.11 (`docs/user/benutzerhandbuch.md`, Überschrift `#### Login-Dienst `sshd` (optional)`): zwei Befehle (`u-boot config set devcontainer.features.sshd.enabled true`, `u-boot generate devcontainer`), danach vier Aufzählungspunkte (Passwort setzen oder Schlüssel eintragen · Port 2222 nicht veröffentlicht · keine Lockerung des Profils · Grenzen des Dienstes). Es fehlt die Reihenfolge der Handlungen: Container neu bauen, Zugang einrichten, verbinden.
- Der Katalog-Eintrag: `featureCatalogue()["sshd"]` in `internal/hexagon/application/devcontainer_features.go` (Quelle `ghcr.io/devcontainers/features/sshd`, Version `1`, Kurzbeschreibung „SSH server (login service, port 2222)“); die Warnung beim Erzeugen: `sandboxRelaxationWarnings` in `internal/hexagon/application/devcontainer_sandbox.go` (Text: „login service sshd is active in the container (port 2222; password and key login, root login permitted; the accounts stay locked until you set a password with `sudo passwd <user>` or add a key); the profile limits damage and is not an isolation boundary“), nur im Sandbox-Profil.
- Der Port wird nicht weitergeleitet: `forwardPorts` entsteht nur aus den Ports aktiver Dienste (`collectDevcontainerForwardPorts` in `internal/hexagon/application/generate.go`) und steht im erzeugten `devcontainer.json` nur, wenn die Liste nicht leer ist (`{{- if .ForwardPorts}}` in `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`); die `runArgs` des Sandbox-Profils enthalten keine Port-Veröffentlichung (`runArgs` entstehen in `devcontainerTemplateData`, `devcontainer_sandbox.go`, nur aus den Lockerungen für Podman, Docker und Egress).
- Im Sandbox-Profil ist der Projektordner ein benanntes Volume, kein Bind-Mount (`"workspaceMount": ""` und `mounts` mit `type=volume` in derselben Vorlage); Host-Dateien wie `~/.ssh` sind im Container nicht eingebunden ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)). Ein Schlüssel kommt also nicht über einen Mount in den Container.
- Der Container-Benutzer ist `vscode` (`"remoteUser": "vscode"` in der Vorlage).

### Was die Anleitung nur sagen darf

Jede Aussage der Anleitung ist eines von drei: (a) im Handbuch, im Code oder in einer Messung **bereits belegt** (Fundstelle im Plan-Eintrag des Implementers, §3.10 von [`AGENTS.md`](../../../../AGENTS.md)); (b) eine Aussage über ein Fremdwerkzeug (`ssh-keygen`, `ssh`, `ssh-copy-id`, Port-Weiterleitung des Editors), die durch eine Quelle oder einen Praxistest belegt ist (der Beleg steht im Plan-Eintrag des Implementers) und, soweit sie nicht in u-boot gemessen ist, als **nicht gemessen** gekennzeichnet wird; (c) **gestrichen**. Eine Aussage der Klasse (b) ohne Beleg steht nicht im Text.

Die Reihenfolge, die der Projektinhaber vorgibt: Feature aktivieren → `generate devcontainer` → Container neu starten; Zugang einrichten (`ssh-keygen` auf dem Host, Schlüssel in die `authorized_keys` des Container-Benutzers, alternativ ein Passwort setzen, weil die Konten bis dahin gesperrt sind); verbinden über Port 2222 und wie der Port auf den Host kommt; Sicherheitshinweis (Root-Login erlaubt, Passwort- und Schlüssel-Login; das Profil begrenzt Schaden, ist keine Isolationsgrenze).

**Offene Punkte, die der Implementer **vor** dem Schreiben löst (lesen oder messen, sonst streichen):**

- **„neu starten“ oder „neu bauen“:** Die Reihenfolge des Projektinhabers sagt „neu starten“; das Handbuch sagt für die Drittanbieter-Features, ihr Installationsskript laufe beim Image-Build, und für die Egress-Beschränkung „Nach einer Änderung führen Sie `u-boot generate devcontainer` aus und starten den Container neu“. Ob eine Änderung der Feature-Liste einen Neubau des Images braucht, sagt das Handbuch für `sshd` nicht; das Handbuch schreibt nur, was belegt ist („Neu bauen nötig“, weil das Installationsskript beim Image-Build läuft), und kennzeichnet als „nicht gemessen“, ob ein Neustart genügt.
- **Wie der Schlüssel in den Container kommt:** kein Mount (siehe oben); der Weg über ein Terminal im Container (Editor-Terminal oder die Engine) ist eine Aussage der Klasse (b).
- **Wie der Port auf den Host kommt:** `forwardPorts` bleibt leer; das Handbuch sagt „Erreichbar ist der Port über das Port-Forwarding Ihres Dev-Container-Werkzeugs“. Gemessen ist, dass der Dienst auf Port 2222 lauscht (Dev-Container-CLI 0.89.0, Basis-Image `mcr.microsoft.com/devcontainers/base:debian`); gemessen ist **nicht**, wie ein Editor oder die CLI den Port auf den Host bringt. Was darüber hinausgeht, steht als „nicht gemessen“ oder nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Zusage oder ein Feature** (etwa ein `forwardPorts`-Eintrag für `sshd`, ein Befehl zum Eintragen eines Schlüssels, ein gesetztes Passwort): wäre eine Änderung am Vertrag oder am Code; der Projektinhaber hat „das Handbuch verbessern, sonst nichts“ entschieden.
- **Kennungen und Chronik in der Nutzer-Doku:** Nutzer können damit nichts anfangen ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 6); der Text nennt den Stand.
- **Eine Änderung der Sicherheitsaussagen:** Sie stehen schon im Abschnitt und in der Warnung des Codes; der Slice ordnet sie um, ändert sie nicht.
- **Eine Messung mit der Dev-Container-CLI:** die CLI ist nicht Teil des Repo-Werkzeugs (kein Node im Image `u-boot:test-docker-tools`); Aussagen, die nur sie belegen würde, stehen als „nicht gemessen“.

## 2. Definition of Done

- [x] **Anleitung im Handbuch:** Der Abschnitt „Login-Dienst `sshd` (optional)“ führt Schritt für Schritt von „Feature aktivieren“ über „Zugang einrichten“ (Schlüssel oder Passwort) zu „Verbinden“ (Port 2222, Weg des Ports auf den Host) und endet im Sicherheitshinweis; jede Aussage ist belegt, als „nicht gemessen“ gekennzeichnet oder gestrichen.
- [x] `make gates` grün (`docs-check` ohne Befund); `docs/user/benutzerhandbuch.md` trägt keine `LH-*`-, `ADR-*`-, `SPEC-*`- oder Slice-Kennung und keine Chronik.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (öffentlicher Vertrag; der einzige berührte).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Jede Aussage über Code steht mit Fundstelle (gelesen am 2026-10-18). Der Implementer erweitert die Tabelle um jede Aussage der Anleitung samt Fundstelle und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10); eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitt `#### Login-Dienst `sshd` (optional)` in 4.11 (gelesen: zwei Befehle, vier Aufzählungspunkte) wird zur Anleitung umgebaut; die Aussagen zu Port 2222, gesperrten Konten, `sudo passwd`, `~/.ssh/authorized_keys`, `PermitRootLogin` und der Grenze „Profil begrenzt Schaden“ stehen schon dort und wandern in die Schritte; die Katalog-Zeile `sshd` in 4.10 und die Zeile „Feature `sshd` — keine zusätzlichen Hosts“ in der Egress-Tabelle bleiben, solange sie zur neuen Fassung stimmen (Querverweis `#login-dienst-sshd-optional` bleibt gültig, wenn die Überschrift bleibt) |
| `CHANGELOG.md` | unverändert | gelesen: der Block `[Unreleased]` führt Produktverhalten (Added, Changed); eine Handbuch-Verbesserung ohne neue Zusage ist kein Eintrag; der Release-Cut entscheidet, ob der Lead sie nennt |
| `README.md`, `README.de.md` | unverändert | gelesen: beide nennen `sshd` nur als Katalog-Schlüssel (`README.md:11`, `:153`, `README.de.md:12`, `:155`) und in der Releases-Tabelle; keine Aussage zur Nutzung |
| `docs/plan/planning/open/README.md` | update | gelesen: die Zeile zum Slice trägt den Zeiger `../done/slice-v1-handbuch-sshd-anleitung.md` (Zeile 15) und beschreibt den Gegenstand zutreffend; `make slice-mv` nach `done/` zieht den Zeiger nach (reiner Verweis, kein Inhalt) |
| `docs/reviews/2026-10-18-slice-v1-handbuch-sshd-anleitung.md` | neu | Review-Report des Reviewers; zusammen mit der Behebung der Findings committet |
| `docs/plan/planning/in-progress/roadmap.md` | update | gelesen: Abschnitt „Offene Wellen“ trug „Nichts in Arbeit.“; während der Arbeit entfernt (d-check `planning` bemängelt den Marker bei einem Slice in `in-progress/`), nach dem Move nach `done/` wieder gesetzt |
| `docs/plan/planning/done/slice-v1-handbuch-sshd-anleitung.md` | update | diese Tabelle um die Fundstellen der Aussagen erweitert; §6, §7, DoD bei der Closure gefüllt; Ruheort `done/` |
| `docs/reviews/2026-10-18-verification-slice-v1-handbuch-sshd-anleitung.md` | neu | Verification-Report des Verifiers (Auflagen A-1 bis A-5, Diskrepanz D-1); mit der Closure committet |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-v1-handbuch-sshd-anleitung.md` | neu | Beleg Nr. 11 der Beobachtung (Klasse „unverändert“ ohne gelesene Aussage, Review F-6) |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/state.md` | update | gelesen: nennt zehn Vorgänge, zwei nach der Verkörperung; auf elf und drei nachgezogen, Entscheidung über eine weitere Verkörperung bleibt offen |

Fundstellen der Aussagen der Anleitung (gelesen am 2026-10-18):

| Aussage im Handbuch | Fundstelle |
|---|---|
| Port 2222, `sshd -T`: Passwort-, Schlüssel-Login, Root-Login `yes`; Konten `root`, `vscode` gesperrt; `~/.ssh` fehlt | `done/slice-v1-devcontainer-feature-sshd.md` §3 Messtabelle („`ss -ltn`: `0.0.0.0:2222`“, „`passwordauthentication yes`, `pubkeyauthentication yes`, `permitrootlogin yes`“, „`passwd -S`: `root` und `vscode` gesperrt (`L`)“) |
| kein SSH-Login gemessen | selbe Datei, „Grenzen der Messung“ (1): „Ein tatsächlicher SSH-Login mit gesetztem Passwort wurde nicht ausgeführt“ |
| Warnung nennt Passwort-/Schlüssel-Login, Root-Login, `sudo passwd <user>`; nur im Sandbox-Profil; unter `--quiet` sichtbar | `devcontainer_sandbox.go`: `if sandbox && loginServiceActive(dc) { relax(…, "login service sshd is active in the container (port 2222; password and key login, root login permitted; …` und Handbuch `--quiet`-Zeile („Sicherheitslockerungen des Sandbox-Profils und ein im Container aktivierter Login-Dienst (Feature `sshd`)“) |
| `forwardPorts` nur aus den Diensten, ohne Dienst-Ports fehlt der Schlüssel | `generate.go`: `collectDevcontainerForwardPorts` („derives the container-side ports of every active service“, `len(services) == 0` → `nil`); Vorlage `devcontainer.json.tmpl`: `{{- if .ForwardPorts}}` |
| `runArgs` ohne Port-Veröffentlichung | `devcontainer_sandbox.go`, `devcontainerTemplateData`: `data.RunArgs` wird nur mit `podmanRelaxations()`, `dockerRelaxation()`, `egressRelaxation()` belegt |
| Projektordner ist ein Volume, kein Host-Mount | Vorlage: `"workspaceMount": ""` und `"mounts": ["source={{.WorkspaceVolume}}-${devcontainerId},target={{.WorkspaceFolder}},type=volume"]` |
| Container-Benutzer `vscode` | Vorlage: `"remoteUser": "vscode"` |
| Installationsskript läuft beim Image-Build; Neubau statt Neustart | Handbuch (Abschnitt Claude Code): „sein Installationsskript läuft beim Image-Build“; der Text sagt „neu aufbauen“ und nennt als „nicht gemessen“, ob ein Neustart genügt; kein Editor-Befehl mehr im Text |
| `ssh-keygen -t ed25519`, `~/.ssh` `700`, `authorized_keys` `600`, `ssh -i … -p 2222 vscode@localhost` ohne Passwort | Praxistest des Reviewers vom 2026-10-18 (Review-Report, F-2): Wegwerf-Container aus `mcr.microsoft.com/devcontainers/base:debian` mit dem Debian-Paket `openssh-server` auf Port 2222; der Login funktionierte. **Näherung:** Debian-Paket statt Feature-Skript, ohne Editor- oder Engine-Portweiterleitung, ohne Dev-Container-CLI; der Text nennt das. Es wurde keine Dokumentation gelesen |
| `docker exec`, Port-Weiterleitung von Editor und Engine, SSH-Login mit dem Feature selbst | weder gelesen noch gemessen; im Text als „nicht gemessen“ gekennzeichnet |
| Loopback-Bindung `127.0.0.1:2222:2222` statt `2222:2222` | Sicherheitsempfehlung des Reviews (F-1): Passwort-Login und erlaubter Root-Login (Messung `sshd -T`, siehe oben) öffnen den Dienst sonst auf allen Host-Schnittstellen; die Bindung selbst ist nicht gemessen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Slice [`slice-v1-sandbox-fallback-host-ausweisung`](../done/slice-v1-sandbox-fallback-host-ausweisung.md) liegt in `done/` (erfüllt, 2026-10-18); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Eine Aussage der Anleitung ist weder belegt noch als „nicht gemessen“ vertretbar und der Projektinhaber will sie behalten; dann braucht sie eine Messung (eigener Vorgang) oder entfällt.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Cut [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md) startet erst, wenn dieser Slice in `done/` liegt.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Aussage ohne Beleg:** Eine Anleitung liest sich wie eine Zusage; ein Schritt, der nur auf dem Papier funktioniert (Weg des Ports, Weg des Schlüssels), wäre eine Doku, die gegen das Verhalten driftet. — **Ausgang:** weiter offen: entfallen für die gemessenen Aussagen (Praxistest des Reviewers, vom Verifier wiederholt; Aussagen zu Code gegen Quelle gelesen); offen bleiben die im Text als „nicht gemessen“ gekennzeichneten (SSH-Login mit dem Feature, Weiterleitung durch Editor und Engine, `docker exec`, Neustart statt Neubau).
- **Zusage durch die Hintertür:** Ein Schritt wie „der Port ist erreichbar“ könnte als Zusage gelesen werden, die das Lastenheft nicht trägt. — **Ausgang:** entfallen: der Text nennt den Stand („u-boot veröffentlicht … nicht“), keine Garantie; Review und Verifier fanden keine Zusage.
- **Sicherheitsaussage verwässert:** Ein Schritt-für-Schritt-Text kann den Hinweis „Root-Login erlaubt, Passwort-Login möglich“ an den Rand drängen. — **Ausgang:** entfallen: Querverweis auf den Sicherheitshinweis steht vor Schritt 1, die Loopback-Bindung in Schritt 3 und im Hinweis, der Schlüssel ist bevorzugt (Review F-1, F-3 behoben).

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-18-slice-v1-handbuch-sshd-anleitung`](../../../reviews/2026-10-18-slice-v1-handbuch-sshd-anleitung.md) (2 MEDIUM, 2 LOW, 2 INFO; F-1 bis F-6 behoben in `541ad35`) und Verification-Report [`2026-10-18-verification-slice-v1-handbuch-sshd-anleitung`](../../../reviews/2026-10-18-verification-slice-v1-handbuch-sshd-anleitung.md) (sachliche DoD erfüllt, Diskrepanzen D-1 und D-2, Auflagen A-1 bis A-5; in dieser Closure erfüllt). Sensoren der Verification am HEAD `541ad35`: `make gates` grün, `make doc-trace` 83 Anforderungen ohne Waise, `docs-check` ohne Befund; `make test-docker` und `make verify-depguard` nicht gefahren (der Diff enthält ausschließlich Markdown). Praxistest des Reviewers, vom Verifier unabhängig wiederholt (Näherung, siehe Grenzen). Commits: `cda63fa` (Anleitung), `7dc8521` (Plan-Tabelle ohne unverlinkte Kennung), `541ad35` (Review-Fixes F-1 bis F-6, Review-Report) und der Closure-Commit (D-1, §6, §7, Register). Nach D-1 ist der Nachsatz zu `ssh-keygen` im Basis-Image aus Schritt 2 des Handbuchs gestrichen (`docs/user/benutzerhandbuch.md`, Schritt „Schlüssel (bevorzugt)“, Nr. 1): er belegte das Werkzeug im Container, der Befehl läuft auf dem Host; die Zeile der Plan-Tabelle dazu entfiel mit ihm.

**Freigabe des Projektinhabers (2026-10-18):** „Das Handbuch können wir dazu verbessern, sonst nichts.“ Der Wortlaut ist die Entscheidung dieses Slice (§1): Handbuch-Text, keine Zusage, kein Feature, keine Änderung am Code; der Diff hält sich daran (`git diff --stat 0bfc46f..HEAD` vor der Closure: Handbuch, Roadmap, Slice, Review-Report). Das Kriterium „Handbuch statt Feature“ hat der Projektinhaber gesetzt, nicht dieser Plan. <!-- d-check:ignore (Wortlaut der Freigabe als Zitat) -->

**Fundstellen je Aussage (`AGENTS.md` §3.10):** die Aussagen der Anleitung stehen mit Fundstelle in der zweiten Tabelle von §3 (Messtabelle des `sshd`-Features, `generate.go`, `devcontainer_sandbox.go`, Vorlage `devcontainer.json.tmpl`, Praxistest); der Verifier hat die Aussagen zu Code gegen die Quelle gelesen (`devcontainer.json.tmpl:43`, `generate.go:934`, `devcontainer_sandbox.go:321`, Handbuch Zeile 310).

- **Was hat funktioniert:** Der Praxistest (Wegwerf-Container aus dem Basis-Image, Debian-`openssh-server` auf Port 2222, Veröffentlichung `127.0.0.1:2222:2222`, Schlüssel erzeugt, Rechte `700`/`600`, `ssh -i … -p 2222 vscode@localhost` mit Exit 0 und `uid=1000(vscode)`) machte aus den Fremdwerkzeug-Aussagen belegte Aussagen; der Reviewer fand die fehlende Loopback-Bindung (F-1), die der Plan nicht vorsah und die ein Gate nie fängt. Der Vorspann „Gemessen / Nicht gemessen“ im Text hielt die Grenzen sichtbar, sodass der Verifier sie nur noch gegen die Plan-Tabelle lesen musste.
- **Was ging anders als geplant:** (1) Der Plan verlangte für Aussagen zu Fremdwerkzeugen eine Quelle oder einen Praxistest; die erste Fassung stützte sich auf „Allgemeinwissen“ und kennzeichnete nur „nicht gemessen“ (F-2), die Kennzeichnung ersetzt den Beleg nicht. (2) §6 stand bis zur Closure auf „offen“ (Auflage A-2 der Verification), wie schon im Vorgänger-Slice. (3) Ein Nachsatz zum Werkzeug im Basis-Image blieb nach den Review-Fixes stehen und belegte das falsche Ziel (D-1); er ist in der Closure gestrichen.
- **Finding-Klassen des Reviews:** Anleitung-Rat ohne Sicherheitsgrenze (F-1, MEDIUM: Loopback-Bindung fehlte; die Aussage steht jetzt in Schritt 3 und im Sicherheitshinweis, die Bindung selbst ist nicht gemessen) · Fremdwerkzeug-Aussage ohne Quelle (F-2, MEDIUM; Plan-Vorgabe gegen Text) · Sicherheitsaussage gleichrangig neben bequemerem Weg (F-3, LOW; der Schlüssel ist jetzt bevorzugt) · Folgerung als Anweisung, Begriff uneinheitlich im Abschnitt (F-4, LOW; „neu aufbauen“ ist als Folgerung aus belegter Prämisse gekennzeichnet, „Nicht gemessen“, ob ein Neustart genügt) · Messaussage nicht vom Autor wiederholt (F-5, INFO; die Klammer ist gestrichen, der Rest D-1) · „unverändert“-Zeile ohne gelesene Aussage (F-6, INFO; behoben, nach Wortsuche auf „unverändert“ trägt keine Zeile der Plan-Tabelle mehr nur eine Begründung ohne Lesebeleg).
- **Steering-Loop-Eintrag (nichts verkörpert, gezählt):** Rückfall der Klasse „unverändert“ ohne gelesene Aussage (F-6) in `BEO-DOK/plan-aussage-ueber-code` ([`observation.md`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)), der dritte Vorgang nach der Verkörperung in `AGENTS.md §3.10` und der erste nach der zweiten Schärfung vom selben Tag. Die Beobachtung steht jetzt bei elf Vorgängen. **Entscheidung nicht getroffen:** ob eine weitere Schärfung oder ein Sensor folgt, ist offen; diese Closure nennt nur den Befund. Der Retirement-Check verlangt nichts, denn die Regel wird weder entfernt noch gelockert. Kein `liegt in`: es wird mit diesem Slice nichts verkörpert.
- **Gezählt, nicht verkörpert:** (1) Fremdwerkzeug-Aussage ohne Quelle (F-2, einmal): sie passt nicht zu `plan-aussage-ueber-code` (dort geht es um Plan-Aussagen über Code, hier um eine Handbuch-Aussage gegen die Plan-Vorgabe), ein eigenes Register-Verzeichnis lohnt bei einem Vorgang nicht. (2) Anleitung-Rat ohne Sicherheitsgrenze (F-1, einmal). (3) Messaussage nicht vom Autor wiederholt (F-5, einmal). Weitere Klassen (F-3, F-4) einmal, ohne Register-Eintrag.
- **Beobachtungs-Register (`../observations/`):** Beleg [`evidence/slice-v1-handbuch-sshd-anleitung.md`](../observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-v1-handbuch-sshd-anleitung.md) in `BEO-DOK/plan-aussage-ueber-code/` ergänzt (elf Belege), `state.md` nachgezogen; `BEO-DOK/spec-bedingung-ohne-code-entsprechung` ohne Wiederauftreten (zwei Vorgänge, unverändert). Kein neues Verzeichnis.
- **Folge-Slices:** keine aus diesem Slice. Der Release-Cut [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md) startet jetzt; ob er die Handbuch-Verbesserung nennt, entscheidet er (CHANGELOG bleibt unverändert, Auflage A-5).
- **Risiken aus §6:** alle drei tragen einen Ausgang: weiter offen (Aussage ohne Beleg, nur für die „nicht gemessen“ gekennzeichneten), entfallen (Zusage durch die Hintertür; Sicherheitsaussage verwässert).
- **Grenzen:** (1) Der Feature-Login selbst ist nicht gemessen: Der Praxistest benutzt das Debian-Paket `openssh-server` statt des Installationsskripts des Features, ohne Dev-Container-CLI. (2) Weiterleitung des Ports durch Editor oder Engine und `docker exec` sind im Text „nicht gemessen“ (der Verifier hat `docker exec` benutzt, es lief; das Handbuch bleibt konservativ). (3) Die Loopback-Bindung `127.0.0.1:2222:2222` ist in der Näherung gemessen (`docker port` zeigt `127.0.0.1:2222`), nicht mit dem Feature. (4) Ob ein Neustart statt eines Neubaus genügt, ist nicht gemessen. (5) Der Passwort-Login wurde nicht ausgeführt.
- **Roadmap:** Der Ruhe-Marker „Nichts in Arbeit.“ steht nach dem Move wieder unter „Offene Wellen“ (`in-progress/` trägt keinen Slice; d-check `planning`).
- **Drei Paarungen:** Anker (kein `liegt in`-Feld, nichts verkörpert) · Folge-Slice (keiner) · Register (ein Pfad, `evidence/` nicht leer) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `docs/plan/` (DOK, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006)); `docs/user/` ist ein öffentlicher Vertrag, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — `BEO-DOK/plan-aussage-ueber-code` (verkörpert, `AGENTS.md` §3.10) und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` (verkörpert) betreffen diesen Slice: jede Aussage über Verhalten steht mit Fundstelle, eine Bedingung im Handbuch, die eine Eigenschaft des Codes behauptet („der Port wird nicht weitergeleitet“), bekommt einen Beleg oder entfällt.

Alle berührten Sub-Areas GF.
