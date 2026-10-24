# Slice CR: Cloud-Modelle von Ollama und Anschluss an einen Ollama-Dienst des Hosts — Lastenheft-CR ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Zusage und Grenzen, [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Netzanschluss, [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-ollama-cloud-modelle`](../welle-ollama-cloud-modelle.md) (Release `v0.16.0`). Vierter von sieben Zeilen der Tabelle in Welle §4 (Reihenfolge Messung → Messung mit Konto → ADR → **CR** → Code → Spezifikation → Release-Cut; **der CR läuft vor dem Code-Slice**, [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) *Folgepflicht*). Muster: [`slice-cr-devcontainer-forge-cli`](../done/slice-cr-devcontainer-forge-cli.md) (Freigabe als eigener Schritt vor der Lastenheft-Änderung) und [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Absatz `ollama`: Cloud-Modelle und Anschluss an einen Dienst des Hosts), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Ausweisung eines Netzanschlusses an den Host), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1 unter `--quiet`); [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (`Accepted`, 2026-10-24; Entscheidungen 2 bis 4, 8, 10, O1 bis O3, Zusage-Menge Z1 bis Z13). Unverändert und nur gelesen: [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Begründung §1, *Stelle 8*).

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), §4.1 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (das tut `slice-v2-spec-ollama-cloud` nach dem Code). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** Implementer (Start 2026-10-24); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Lastenheft sagt zu, was [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) mit der Zusage-Menge Z1 bis Z13 entschieden hat, und nicht mehr: Cloud-Modelle von Ollama sind im Container nutzbar (Anmeldung im Container durch den Nutzer, `ollama.com` durch den Nutzer in `devcontainer.sandbox.egress.allow`), und `u-boot` erzeugt auf Wunsch des Nutzers einen Anschluss des Containers an einen Ollama-Dienst des Hosts, der im Sandbox-Profil ausgewiesen wird und unter `--quiet` sichtbar bleibt. Der Satz „nicht zugesagt: … ein Anschluss an einen Ollama-Dienst des Hosts“ in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) entfällt und wird durch diese Zusage ersetzt. Es ist das **Zielbild**: kein Messwert, kein Messstand, keine Hostnamen von Messumgebungen, kein Prozesszustand, keine Slice-, Wellen- oder ADR-Kennung (Vorgabe des Projektinhabers vom 2026-10-19). Ob der Code es trägt, entscheiden die Tests des Code-Slice, nicht dieser CR.

**Reihenfolge zu [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) (Planänderung, Entscheidung des Standards, 2026-10-24).** Dieser CR läuft **vor** jenem Slice: Er gehört zu einer laufenden Welle (Release `v0.16.0`), jener hat keine Welle und hält keinen Release auf. Die Version vergibt der Implementer **erst beim Start** dieses CR als nächste freie Minor-Version nach dem Stand von `spec/lastenheft.md` (heute `0.22.0`, erwartet `0.23.0`; [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Freigabe-Punkt g). [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) passt sich danach an: Er nimmt die folgende Version, liest seine Fundstellen (Zeilen 411, 467, 469, 1583 bis 1591) am Stand nach diesem CR neu und streicht dann die Messstände, die dieser CR stehen lässt (*Stelle 2*, Frage F3). Die Anpassung der Angaben jenes Slice ist Teil dieser Planung (Plan-Tabelle §3) und wird bei Start und Closure dieses CR nachgezogen. Der Projektinhaber kann die Reihenfolge im Gespräch ändern; dann gilt der Wortlaut dieser Zeile nicht mehr.

### Gefundener Stand (gelesen am 2026-10-24 in `spec/lastenheft.md`, Lastenheft `0.22.0`; Zeilen mit `grep -n` bestimmt, die Absätze gelesen)

| Fundstelle | Aussage (gelesen) | Folge |
|---|---|---|
| Kopf, Zeile 9 | `| Version          | 0.22.0 |` | Versionsfeld auf die beim Start vergebene Version (*Stelle 7*) |
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beispiele, Zeilen 399 bis 409 | Liste endet mit „- Ollama (lokale Sprachmodelle, nur CPU)“ (407), „- GitHub CLI (`gh`)“ (408), „- GitLab CLI (`glab`)“ (409) | ein Spiegelstrich für den Anschluss an den Host (*Stelle 1*) |
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Absatz `ollama`, Zeile 411 | Satz 1 „Ollama ist als eingebautes Feature `ollama` wählbar … und wirkt im Sandbox-Profil“; letzter Satz „Gemessen ist das auf Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6; nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen und Container-Engines als gemessen, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.“ | der Ausschluss „ein Anschluss an einen Ollama-Dienst des Hosts“ entfällt (*Stelle 2*); der Halbsatz „Gemessen ist das auf …“ ist Messstand der Klasse A und bleibt in diesem CR stehen (Frage F3) |
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Absatz `gh`/`glab`, Zeile 413 | Absatz in Form „… sind in beiden Profilen … wählbar …; Nicht zugesagt: …“ (Muster für Stelle 3 und 4: je ein eigener Absatz, Zusage und „Nicht zugesagt“ am Ende) | unverändert; Vorbild der Form |
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Zeile 415 | „Für optionale externe Feature-Quellen gilt:“ | die neuen Absätze stehen zwischen Zeile 411 und Zeile 413 (Cloud-Modelle) bzw. zwischen Zeile 413 und Zeile 415 (Anschluss an den Host) |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Zeile 455 | „- jede Sicherheitslockerung (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) in der Befehlsausgabe ausweisen;“; Zeile 456 „- einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.“ | Netzanschluss an den Host in die Klammer der Zeile 455 (*Stelle 5*) |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 1, Zeile 191 | Aufzählung „(zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up` und `u-boot devcontainer shell`) und eines im Container aktivierten Login-Dienstes …“; Verweisliste „([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten))“; der Satz davor „Vier Ausnahmen bleiben unter `--quiet` sichtbar“ steht in Zeile 189 | Aufzählung und Verweisliste ergänzen (*Stelle 6*); die Zahl „Vier Ausnahmen“ bleibt |
| [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), Zeilen 473 bis 482 | `devcontainer.sandbox.egress.allow` (Liste von Hostnamen) ist der Ort für Nutzer-Einträge; „Ein Name in der Liste schließt seine Subdomains ein“ | der Nutzer-Eintrag `ollama.com` braucht keinen neuen Satz; unverändert (*Stelle 8*) |
| [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Zeile 486 | Zusage auf Zugangsdaten eines Agenten oder Forge-Werkzeugs: weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer erzeugten Datei | unverändert (*Stelle 8*): der Geräteschlüssel im Modell-Volume wird nicht als Zugangsdatum gewertet ([`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Entscheidung 5) |
| [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 502 | „Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`); alle übrigen Hosts und Engines, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt.“ | Vorbild der Zusage-Menge als Wortlaut ohne Messstand; die Stellen 3 und 4 übernehmen die Form |
| §16, Zeilen 1556 bis 1598; Datei 1605 Zeilen | letzte Historie-Zeile `0.22.0` (Zeile 1598); Zeile 1591 (`0.15.0`) nennt „nicht zugesagt: GPU, Anschluss an einen Ollama des Hosts, Modellverwaltung“ | neue Zeile nach 1598 (*Stelle 7*); Zeile 1591 ist Chronik (Klasse D) und bleibt |
| Suche über `spec/lastenheft.md` (`grep -n -F`) | `-cloud`, `host.docker.internal`, `Relay`, `ollama.com`, `ollama-host`, `signin`, `Geräteschlüssel`, `OLLAMA_`, `ANTHROPIC`, `add-host`, `Netzanschluss`, `Host-Anschluss`: **keine Treffer**; `Cloud` nur Zeile 1394 (Abgrenzung „Cloud-Provisioning“, [`LH-ABG-001`](../../../../spec/lastenheft.md#lh-abg-001--kein-vollständiges-deployment-system)); `Ollama` nur Zeilen 407, 411, 1591 | die Zusage ist neu; keine zweite Stelle widerspricht |

### Wortlaut-Entwurf (so, wie er ins Lastenheft geht — nicht verdichtet)

Der Entwurf ist bis zur Freigabe kein Vertragstext. Er trägt keinen Messstand (keine Engine-Version, keinen Host einer Messung, keine Modellnamen), keine Kennung von ADR, Slice oder Welle und keinen Hostnamen außer `ollama.com` (Frage F1). Die Verweise auf Lastenheft-Anforderungen sind hier als Datei-Links gesetzt; im Lastenheft sind es Anker innerhalb der Datei, wie in den Nachbarabsätzen. Die Vorlage an den Projektinhaber zeigt **diesen Text ungekürzt**, je Stelle mit *Alt* und *Neu* ([`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md)).

**Stelle 1 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beispiele (nach Zeile 409).** *Alt:* die Liste endet mit „- GitLab CLI (`glab`)“. *Neu:* ein zusätzlicher Spiegelstrich nach dieser Zeile:

> - Anschluss an einen Ollama-Dienst des Hosts (`ollama-host`)

**Stelle 2 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Absatz `ollama` (Zeile 411), letzter Satz.** *Alt:*

> Gemessen ist das auf Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6; nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen und Container-Engines als gemessen, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.

*Neu* (der übrige Absatz bleibt Wort für Wort; geändert ist nur die Aufzählung nach „nicht zugesagt sind“ — der Halbsatz „Gemessen ist das auf …“ bleibt, Frage F3):

> Gemessen ist das auf Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6; nicht zugesagt sind GPU-Beschleunigung, andere Architekturen und Container-Engines als gemessen, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.

**Stelle 3 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), neuer Absatz Cloud-Modelle (zwischen Zeile 411 und Zeile 413).** *Neu:*

> Mit dem Feature `ollama` sind Cloud-Modelle von Ollama (Modellname mit der Endung `-cloud`) im Sandbox-Profil nutzbar. Der Nutzer startet den Server im Container selbst, meldet sich im Container mit `ollama signin` an und ruft danach `ollama run`, `claude -p` oder `ollama launch claude` auf; ein vorheriges `ollama pull` ist nicht nötig. `u-boot` führt `ollama signin` nicht aus, speichert keine Anmeldung und reicht keine Umgebungsvariable des Hosts für Ollama in den Container. Der Server legt beim ersten Start einen Geräteschlüssel im Volume der Modelle an, nicht im Workspace-Volume; ein neu angelegter Container mit demselben Volume antwortet ohne neue Anmeldung. Hinter der Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) scheitert der Aufruf eines Cloud-Modells an der Namensauflösung von `ollama.com`, solange der Name nicht in `devcontainer.sandbox.egress.allow` steht; mit dem Eintrag fordert der Aufruf ohne Anmeldung zur Anmeldung auf. Die Default-Allowlist enthält den Namen nicht; der Nutzer trägt ihn ein. Zugesagt ist das auf einem Linux-Host mit lokaler Docker-Engine (`amd64`); alle übrigen Hosts und Engines, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt. Nicht zugesagt: der Lauf eines bestimmten Cloud-Modells, Kontingent und Kosten, der interaktive Ablauf von Claude Code und von `ollama launch claude`, ein Aufruf bei gesetztem `ANTHROPIC_API_KEY`, ein gültiger `OLLAMA_API_KEY`, die Wirkung von `ollama signout` im Konto, das Verhalten bei abgelaufener Anmeldung, der Lauf mit Anmeldung hinter der Egress-Restriktion ohne den Eintrag von `ollama.com`, andere Versionen von Ollama als die gepinnte sowie eine bestimmte Version von Claude Code.

**Stelle 4 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), neuer Absatz Anschluss an einen Ollama-Dienst des Hosts (zwischen Zeile 413 und Zeile 415).** *Neu:*

> Auf Wunsch des Nutzers erzeugt `u-boot` einen Anschluss des Containers an einen Ollama-Dienst des Hosts: Das eingebaute Feature `ollama-host` (`u-boot config set devcontainer.features.ollama-host.enabled true`) wirkt in beiden Profilen (Standard und Sandbox), unabhängig vom Feature `ollama`. Es trägt genau einen Eintrag in die Container-Optionen der erzeugten Devcontainer-Konfiguration ein, der den Host aus dem Container adressierbar macht, und erzeugt sonst nichts zu diesem Anschluss: keinen Relay, keine Umgebungsvariable (`OLLAMA_HOST`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`) und keinen Eintrag der Egress-Allowlist. `u-boot` startet keinen Relay und ändert weder die Bindung noch eine Datei des Dienstes auf dem Host. Ob der Dienst aus dem Container erreichbar ist, hängt von seiner Bindung ab: Lauscht er nur auf der Loopback-Adresse des Hosts, ist er erst mit einem Relay oder einer geänderten Bindung erreichbar, die der Nutzer einrichtet. Mit einem Relay des Nutzers und einer vom Nutzer gesetzten Basis-URL für Claude Code antwortet `claude -p` im Container über den Dienst des Hosts. `u-boot` weist den Anschluss beim Erzeugen aus; im Sandbox-Profil ist er eine ausgewiesene Sicherheitslockerung ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)). Hinter der Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) trägt der Anschluss nicht: `u-boot` trägt keine Adresse in die Allowlist ein und warnt beim Erzeugen, ohne abzubrechen. Zugesagt ist das auf einem Linux-Host mit lokaler Docker-Engine (`amd64`) und einem Dienst, der auf der Loopback-Adresse des Hosts lauscht; alle übrigen Hosts, Engines und Bindungen, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt. Nicht zugesagt: die Erreichbarkeit des Dienstes für andere Container, `ollama run` und `ollama launch` mit dem Dienst des Hosts hinter der Egress-Restriktion, die Hosts, die der Dienst des Hosts für den Cloud-Dienst braucht.

**Stelle 5 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Aufzählung des erzeugten Ergebnisses (Zeile 455).** *Alt:*

> - jede Sicherheitslockerung (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) in der Befehlsausgabe ausweisen;

*Neu:*

> - jede Sicherheitslockerung (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices, ein auf Wunsch des Nutzers erzeugter Netzanschluss an den Host nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) in der Befehlsausgabe ausweisen;

**Stelle 6 — [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 1 (Zeile 191).** *Alt* (die geänderte Stelle, der Rest der Zeile bleibt Wort für Wort):

> 1. die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up` und `u-boot devcontainer shell`) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)) …

*Neu:*

> 1. die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices, eines auf Wunsch des Nutzers erzeugten Netzanschlusses an den Host im Sandbox-Profil, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up` und `u-boot devcontainer shell`) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)) …

**Stelle 7 — Versionsfeld (Zeile 9) und Historie (§16, nach Zeile 1598).** *Neu:* im Versionsfeld die beim Start vergebene Version (erwartet `0.23.0`) und eine Historie-Zeile (Datum des Lastenheft-Commits):

> | 0.23.0 | JJJJ-MM-TT | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Cloud-Modelle von Ollama im Container (Anmeldung im Container durch den Nutzer, `ollama.com` durch den Nutzer in `devcontainer.sandbox.egress.allow`) und optionaler Anschluss an einen Ollama-Dienst des Hosts über das Katalog-Feature `ollama-host` (ein erzeugter Eintrag, Relay beim Nutzer, hinter der Egress-Restriktion nicht tragend, ausgewiesen); der Anschluss an einen Ollama-Dienst des Hosts ist nicht mehr als „nicht zugesagt“ geführt; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Ausweisung des Netzanschlusses an den Host als Sicherheitslockerung; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Ausnahme 1 unter `--quiet` nennt den Netzanschluss an den Host im Sandbox-Profil. | Vereinbarung mit dem Projektinhaber |

**Stelle 8 — keine Änderung** an [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (der Eintrag des Nutzers in `devcontainer.sandbox.egress.allow` ist dort vorgesehen; die Default-Allowlist bekommt keinen neuen Namen), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (der Geräteschlüssel liegt im Modell-Volume, nicht im Workspace-Volume, nicht im Image und in keiner erzeugten Datei; keine Wertung als Zugangsdatum) und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); die Aufzählung der optionalen Felder der Konfiguration bleibt unverändert (kein neuer Konfigurationspfad, [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) O1). Die Begründung gehört zur Freigabe; sie entscheidet der Projektinhaber.

### Herkunft jedes Satzes (Satz → Zusage-Zeile oder Test)

Die ADR verlangt, dass der Wortlaut keinen Satz trägt, der in der Zusage-Menge keine Zeile hat. Die Spalte *Z / Test* nennt die Zeile aus [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (§ *Zugesagt (bedingungsgleich) und nicht zugesagt*) oder den Test der Fitness Function, der den Satz im Code-Slice pinnt. Die ausführliche Beleg-Tabelle (Bedingung → Messaussage oder Test) entsteht in §7 beim Start.

| Satz (Stelle) | Z / Test |
|---|---|
| 3.1 Cloud-Modelle mit dem Feature `ollama` im Sandbox-Profil nutzbar | Z1 (Pfad: Handvariante `discover`-Container) |
| 3.2 Server selbst starten, `ollama signin`, danach `ollama run`, `claude -p`, `ollama launch claude`; `ollama pull` nicht nötig | Z1, Z2, Z3 |
| 3.3 `u-boot` führt `ollama signin` nicht aus, speichert nichts, reicht keine Variable des Hosts durch | Test des Code-Slice: `agentCredentialEnv` ohne Ollama-Namen, keine der Namen in einer erzeugten Datei (Fitness Function, Zeilen 1 und 2) |
| 3.4 Geräteschlüssel im Modell-Volume, nicht im Workspace-Volume | Z6 |
| 3.5 neuer Container mit demselben Volume ohne neue Anmeldung | Z4 (Pfad: Handvariante) |
| 3.6 hinter der Egress-Restriktion Lookup-Fehler ohne Eintrag, Anmeldeaufforderung mit Eintrag; Default-Allowlist ohne den Namen | Z5; Test: Zeile `ollama` in `egressFeatureHosts` hat genau zwei Namen (Fitness Function) |
| 3.7 Zusage-Menge Linux-Host `amd64`, lokale Docker-Engine | Bedingung jeder Zeile (Messlage); Form wie [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) |
| 3.8 Nicht zugesagt | Abschnitt *Nicht gemessen, nicht zugesagt* der ADR; keine Zeile nötig |
| 4.1 `ollama-host` in beiden Profilen, unabhängig von `ollama` | Z13 (Produktversprechen); Test je Profil (Fitness Function) |
| 4.2 genau ein Eintrag, sonst nichts zu diesem Anschluss | Z13; Test: keine der vier Variablen und kein Egress-Eintrag in einer erzeugten Datei |
| 4.3 kein Relay, keine Änderung am Dienst | Entscheidung 4 (Relay beim Nutzer); Nicht-Handeln, kein Test nötig |
| 4.4 Erreichbarkeit hängt von der Bindung ab; Loopback-Bindung nur mit Relay oder geänderter Bindung | Z11 |
| 4.5 mit Relay und Basis-URL antwortet `claude -p` | Z9 (Pfad: Handvariante) |
| 4.6 Ausweisung beim Erzeugen; im Sandbox-Profil Sicherheitslockerung | Z13 trägt den Satz für [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); Test: `SecurityRelaxation` im Sandbox-Profil (Fitness Function) — **keine eigene Z-Zeile, Frage F6** |
| 4.7 hinter der Egress-Restriktion trägt der Anschluss nicht; keine Adresse in der Allowlist; Warnung ohne Abbruch | Z12; Entscheidung O2; Test: Warnung, keine Adresse im Egress-Skript (Fitness Function) — Warnung **ohne Z-Zeile, Frage F6** |
| 4.8 Zusage-Menge Linux-Host `amd64`, lokale Docker-Engine, Loopback-Bindung | Bedingung jeder Zeile; Z11 |
| 4.9 Nicht zugesagt | Abschnitt *Nicht gemessen, nicht zugesagt* der ADR |
| 5 / 6 Ausweisung als Sicherheitslockerung, sichtbar unter `--quiet` im Sandbox-Profil | Entscheidung 8, O3; Test: im Sandbox-Profil `SecurityRelaxation` und sichtbar unter `--quiet`, im Standardprofil entfällt die Meldung unter `--quiet` (Fitness Function) |

### Fragen an den Projektinhaber für die Freigabe (Planner, aus dem Lesen von Lastenheft und ADR; keine Änderung des Entwurfs)

**Standard** heißt: der Entwurf oben trägt diese Antwort; der Planner schlägt sie vor, entschieden ist sie erst mit der Freigabe.

- **F1 — Namen im Wortlaut.** Der Entwurf nennt `ollama.com` (der Nutzer muss den Namen in `devcontainer.sandbox.egress.allow` eintragen; ohne den Namen im Vertrag ist die Zusage nicht ausführbar) und nennt `host.docker.internal` **nicht** (Eintrag der Container-Optionen als Spezifikationsdetail). Der Forge-CR hat Hostnamen aus dem Lastenheft herausgehalten (Antwort Nein, [`slice-cr-devcontainer-forge-cli`](../done/slice-cr-devcontainer-forge-cli.md) §7 *Freigabe*). *Standard:* `ollama.com` genannt, `host.docker.internal` nicht. *Alternative:* `ollama.com` umschrieben („der Name des Cloud-Dienstes von Ollama“); dann steht der Name im Handbuch und in der Spezifikation.
- **F2 — Zusage-Menge des Cloud-Wegs.** Der Entwurf übernimmt die Menge der Messlage (Linux-Host `amd64`, lokale Docker-Engine) in der Form von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); Colima und `arm64` sind für diese zwei Absätze **nicht zugesagt**. Der Absatz zur Installation von `ollama` nennt dagegen weiter Colima und `arm64` als Messstand (Frage F3): im selben Feature stehen zwei Mengen nebeneinander, bis [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) die Menge der Installation entscheidet. *Standard:* so, mit der Spannung ausgewiesen.
- **F3 — Messstand im Absatz `ollama`.** [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (*Folgepflicht*) und Welle §4 sagen: „der Messstand im Absatz entfällt“. Ihn zu streichen verlangt die Menge der Installation im Wortlaut (Frage 1 von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md): `arm64` unter Colima, `amd64` auf Linux oder beides) — eine Entscheidung, die nicht Gegenstand der ADR ist. *Standard (Planänderung gegenüber dem Wortlaut der Folgepflicht, ausgewiesen):* Dieser CR ändert den Absatz `ollama` nur um den Ausschluss „ein Anschluss an einen Ollama-Dienst des Hosts“ (Stelle 2); der Halbsatz „Gemessen ist das auf …“ bleibt und wird von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) gestrichen, der danach läuft. Die neuen Absätze tragen keinen Messstand. *Alternative:* der Halbsatz entfällt schon hier; dann entscheidet die Freigabe zugleich die Zusage-Menge der Installation (Frage 1 des Messstände-CR).
- **F4 — Modelle.** Der Entwurf sagt „Cloud-Modelle“ allgemein und nimmt in „Nicht zugesagt“ den Lauf eines bestimmten Modells aus; die Läufe der Messung betreffen zwei Modelle. Modellnamen im Lastenheft wären Messstand und altern beim Anbieter. *Standard:* allgemein, ohne Namen, wie im Entwurf. *Alternative:* die zwei Modelle als Beispiele nennen.
- **F5 — Sichtbarkeit unter `--quiet`.** Der Entwurf stützt sich auf [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Entscheidung 8: Im Sandbox-Profil bleibt die Ausweisung des Anschlusses unter `--quiet` sichtbar (Stelle 6), im Standardprofil entfällt die Meldung (die Regel „Alle anderen Hinweise und Warnungen entfallen unter `--quiet`“ bleibt, Zeile 196). Die Zusatzwarnung „hinter der Egress-Restriktion trägt der Anschluss nicht“ gehört nach der ADR zur selben Meldung; weil die Egress-Restriktion das Sandbox-Profil voraussetzt, ist sie damit unter `--quiet` sichtbar. Stelle 6 nennt sie nicht gesondert. *Standard:* so; *Alternative:* Stelle 6 nennt die Zusatzwarnung ausdrücklich.
- **F6 — Sätze ohne eigene Zusage-Zeile.** Die ADR verlangt, dass der Wortlaut keinen Satz ohne Zeile in Z1 bis Z13 trägt. Z13 deckt nur die erzeugte Zeile; die Sätze „weist den Anschluss beim Erzeugen aus“, „warnt beim Erzeugen, ohne abzubrechen“ (Sätze 4.6, 4.7) und die Sichtbarkeit unter `--quiet` (Stellen 5 und 6) stehen als Entscheidung 8 und O2 mit einem Test der Fitness Function, aber ohne Messaussage. Der Planner liest das als Produktversprechen derselben Art wie Z13. *Standard:* die Sätze bleiben, ihr Beleg ist der Test des Code-Slice (Auftrag in der Beleg-Tabelle §7). *Alternative:* sie entfallen aus dem Lastenheft, und die ADR bleibt die einzige Quelle; dann fehlt [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) der Satz, den Entscheidung O3 verlangt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Katalogeintrag `ollama-host`, Template, Meldungen, Golden Cases, Tests** — der CR schreibt den Vertrag, der Code folgt ihm (`slice-v1-devcontainer-ollama-cloud`, Welle §4; [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) *Folgepflicht*: der Code beginnt erst nach der Freigabe des CR).
- **`spec/spezifikation.md` und `spec/architecture.md`** — nach dem Code in `slice-v2-spec-ollama-cloud`; [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (gelesen Zeile 739: „**Nicht zugesagt.** GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, …“), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) und [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) präzisieren das Lastenheft und erweitern es nie ([`AGENTS.md`](../../../../AGENTS.md) §2). Bis zu jenem Slice widerspricht Zeile 739 dem Lastenheft; das ist die vorgesehene Folge der Reihenfolge CR → Code → Spezifikation (§6).
- **Handbuch, READMEs, `CHANGELOG.md`, `harness/replay.md`** — es ändert sich kein Verhalten; die Code- und Spezifikations-Slices bringen die Dokumente mit (Handbuch Zeile 2023 „Cloud-Modelle sind nicht enthalten“ und Zeile 2043, `README.de.md` Zeile 14 „Cloud-Modelle von Ollama sind nicht enthalten“, `harness/replay.md` Zeile 58 ohne Zeile zu `ollama-host`).
- **Änderungen an [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)** — die ADR führt aus ihnen keinen neuen Satz (*Folgepflicht*); die Begründung je Anforderung gehört zur Freigabe (Stelle 8).
- **Die übrigen Messstände in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)** (Halbsatz „Gemessen ist das auf …“, Zeilen 467 und 469, Historie-Zeilen) — eigener CR ([`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md)); dieser CR fügt keinen hinzu.
- **Eine Zusage über das Gemessene hinaus** — der Wortlaut bleibt in Z1 bis Z13; ein Nachmessen ist vom Projektinhaber verworfen ([`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Entscheidung 11, Freigabe-Punkt d); fällt ein späteres Urteil „Läufe mit Konto genügen nicht“, ist Ausgang B die Rückfallebene (Folge-ADR), nicht eine Änderung dieses CR.
- **Die Release-Nummer** — Lastenheft-Version ist nicht Release-Version; `v0.16.0` schneidet der Release-Cut der Welle.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Freigabe des Projektinhabers vor der Änderung:** Der Wortlaut der Stellen 1 bis 7 ist dem Projektinhaber **ungekürzt, wie er ins Lastenheft geht**, je Stelle mit *Alt* und *Neu* vorgelegt (nicht verdichtet), die Fragen F1 bis F6 aus §1 sind beantwortet, und die Entscheidung je Stelle (Ja / Ja mit Änderung / Nein) steht in §7 *Freigabe* mit Datum; ein geänderter Satz steht dort im freigegebenen Wortlaut. *Beleg:* §7 *Freigabe*; kein Commit von `spec/lastenheft.md` vor diesem Eintrag (Commit-Reihenfolge in `git log`).
- [ ] **Wortlaut im Lastenheft** nach der Freigabe, wortgleich mit ihr; Zielbild ohne Messwert, ohne Messstand, ohne Hostnamen von Messumgebungen, ohne Prozesszustand und ohne Verweis auf ADR, Slice oder Roadmap (Decken-Regel, [`AGENTS.md`](../../../../AGENTS.md) §3.4); der Ausschluss „ein Anschluss an einen Ollama-Dienst des Hosts“ steht im Absatz `ollama` nicht mehr, die Historie-Zeile `0.15.0` bleibt. Jede Bedingung im Text ist durch eine Zusage-Zeile der ADR oder durch ein Zielbild des Code-Slice mit benanntem Test belegt (Beleg-Tabelle in §7, [`AGENTS.md`](../../../../AGENTS.md) §3.10). *Beleg:* Vergleich Satz für Satz im Verification-Report, Fundstellen in §7.
- [ ] **Versionsfeld und Historie:** die beim Start vergebene Version (nächste freie Minor-Version; heute `0.22.0`, erwartet `0.23.0`) im Versionsfeld und eine Historie-Zeile in §16, die nur die freigegebenen Stellen nennt und keine Messangabe trägt. *Beleg:* Versionsfeld, Historie-Zeile.
- [ ] `make gates` grün, `make docs-check` grün. `make doc-trace` ist für [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) ein Beleg (Familie `LH-FA-*`, `.d-check.yml` Abschnitt `trace`). *Beleg:* Läufe und Exit-Codes in §7.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). **Nachlauf-Review:** Ändert eine Fix-Runde Lastenheft-Text, liest ein Nachlauf-Review genau diesen Diff als eigene Datei, **bevor** der Verifier startet ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)); ändert sie ihn nicht, steht das mit dem Diff-Beleg in §7.
- [ ] **Verifier** (eigener Kontext, `harness/verification.md`): Verification-Report unter `docs/reviews/` vergleicht das Lastenheft Satz für Satz mit der Freigabe in §7 und mit der Zusage-Menge der ADR; die Plan-Tabelle ist gegen `git diff --stat` abgeglichen (nach der letzten Fix-Runde).
- [ ] Doku-Update: keines — der CR ändert kein Verhalten; Roadmap, Welle-Datei und `open/README.md` zeigen den Zustand.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer Sensor · benannte Spec-Lücke). *Beleg:* §7.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — u-boot führt Wellen; die Welle-Closure prüft sie.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Ablauf** (die Reihenfolge ist Teil des Plans):

1. **Start und Version.** Der Implementer setzt `Verantwortlich:`, liest `spec/lastenheft.md` neu (Versionsfeld, letzte Historie-Zeile, Zeilen der Stellen 1 bis 6), vergibt die nächste freie Minor-Version (erwartet `0.23.0`) und trägt sie in §7 ein. Er zieht die Angaben von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) und seiner Zeile in `open/README.md` auf die folgende Version nach (Muster [`slice-cr-devcontainer-forge-cli`](../done/slice-cr-devcontainer-forge-cli.md) §6, *Versionsfolge*).
2. **Vorlage.** Der Implementer legt dem Projektinhaber den Wortlaut der Stellen 1 bis 7 aus §1 vor, **Satz für Satz, ungekürzt, so wie er ins Lastenheft geht** (kein Auszug, keine Zusammenfassung; *Alt* und *Neu* je Stelle), mit der Tabelle *Herkunft jedes Satzes* und den Fragen F1 bis F6. Rückfrage im Gespräch; kein Lastenheft-Commit.
3. **Freigabe.** Die Antwort je Stelle (Ja / Ja mit Änderung / Nein) steht in §7 *Freigabe*; ein geänderter Satz steht dort im freigegebenen Wortlaut, ein verdichteter Satz der Vorlage wird als Abweichung gegen diesen Entwurf ausgewiesen ([`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md)). Der Entwurf in §1 wird nicht umgeschrieben; die Abweichung trägt §7.
4. **Vorab-Suche** vor Schritt 5 mit Fundstellen in §7 (Muster [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)): `Ollama`, `Cloud`, `-cloud`, `Host`, `Anschluss`, `Relay`, `host.docker.internal`, `ollama.com`, `ollama-host`, `signin`, `Geräteschlüssel`, `OLLAMA_`, `ANTHROPIC`, `add-host`, `Sicherheitslockerung`, `Netzanschluss`, `Loopback`, `Katalog`, `Beispiele`, `Zusage-Menge`, `Befehlsübersicht`, `Features` (Aufzählungen aller Feature-Namen, etwa in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) oder der Befehlsübersicht). Gelesen am 2026-10-24 (Planner): Treffer siehe Tabelle *Gefundener Stand* in §1 (`Ollama`: Zeilen 407, 411, 1591; `Sicherheitslockerung`: Zeilen 191, 455, 499, 1584, 1585; `Loopback`: Zeile 411; die übrigen Suchwörter ohne Treffer); die Suche des Implementers wiederholt sie und ergänzt `Katalog`, `Beispiele`, `Befehlsübersicht` und die Feature-Aufzählungen.
5. **Lastenheft.** `spec/lastenheft.md` ändert sich genau im freigegebenen Wortlaut (ein Commit), danach Versionsfeld und Historie-Zeile.
6. Review, gegebenenfalls Fix-Runde mit Nachlauf-Review, Verification, Closure.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat — es gibt keine.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 bis 7 im freigegebenen Wortlaut. Gelesen am 2026-10-24: Versionsfeld `0.22.0` (Zeile 9); Beispiele Zeilen 399 bis 409; Absatz `ollama` Zeile 411; Absatz `gh`/`glab` Zeile 413; „Für optionale externe Feature-Quellen gilt:“ Zeile 415; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Zeilen 441 bis 458, Punkt Sicherheitslockerung Zeile 455; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeile 189 („Vier Ausnahmen“) und Ausnahme 1 Zeile 191; §16 beginnt Zeile 1556, letzte Historie-Zeile `0.22.0` Zeile 1598 (Datei 1605 Zeilen) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „Offen:“ (gelesen: „vier weitere benannt ohne Datei (Lastenheft-CR vor dem Code-Slice, Code, Spezifikation, Release-Cut zuletzt)“) nennt den CR als Datei (Planner, mit der Anlage dieses Plans); Zeile *Historische Trigger-Verschiebungen* zur Reihenfolge (Planner); „Nichts in Arbeit.“ (gelesen) bei Start auf diesen Slice und bei der Closure zurück (`.d-check.yml`, Abschnitt `planning`; `make doc-planning`) |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile dieses Slice (Anlage durch den Planner; bei der Closure „geliefert (liegt in `done/`)“); Zeile von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) (gelesen Zeile 33: „auf die nächste freie Version (nach 0.22.0)“) um die Reihenfolge ergänzt (Planner) und beim Start dieses CR auf die Version nach diesem CR gesetzt (Implementer) |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Reihenfolge und Zitate (Planner, mit der Anlage dieses Plans): gelesen §1 Absatz *Zusage-Grenzen bleiben* und Tabellenzeile 411 zitieren den Ausschluss „ein Anschluss an einen Ollama-Dienst des Hosts“ als geltenden Wortlaut; beides steht nach diesem CR nicht mehr so. Beim Start dieses CR (Implementer): Versionsangaben (gelesen §2 „zuletzt `0.22.0`“) und Zeilenangaben auf den Stand nach diesem CR |
| `docs/plan/planning/welle-ollama-cloud-modelle.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | §4 Tabelle: Zeile des CR (Kennung `slice-cr-lastenheft-ollama-cloud-modelle`, Datei besteht) und der Satz „Die ersten drei Zeilen haben eine Datei“ (gelesen, Zeile 84); §5 Frage 3 auf Ist-Stand und die Abgrenzung zum Messstände-CR (gelesen Zeile 93); Verweis-Nachzug durch `make slice-mv` bei den Moves; die vergebene Lastenheft-Version bei der Closure |
| diese Datei | update | `Verantwortlich:`, Freigabe in §7, Plan-Abgleich, Closure; Moves per `make slice-mv` (eigene Commits, [`AGENTS.md`](../../../../AGENTS.md) §3.3) |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review-Report, Nachlauf-Review(s) und Verification-Report, von Reviewer bzw. Verifier geschrieben |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence-Dateien für angefallene Beobachtungen (§7) |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/`, `harness/replay.md`, `docs/plan/adr/`, `docs/messungen/`, `internal/`, `cmd/`, `tools/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten. Gelesen als Stand, den die Folge-Slices ändern oder der Beleg bleibt: `spec/spezifikation.md` Zeile 739 ([`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), *Nicht zugesagt* mit dem Anschluss an einen Ollama-Dienst des Hosts); `docs/user/benutzerhandbuch.md` Zeile 2023 („**Cloud-Modelle sind nicht enthalten:**“) und Zeile 2043 („… ein Anschluss an einen Ollama-Dienst des Hosts …“); `README.de.md` Zeile 14 („Cloud-Modelle von Ollama sind nicht enthalten“); `CHANGELOG.md` Zeile 12 (`[Unreleased]` leer); `harness/replay.md` Zeile 58 (Zeile zu `ollama` ohne `ollama-host`); [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) `Accepted`, bleibt unverändert (`make doc-immutable RANGE="<Commit der ADR>..HEAD"`); der Katalog und die Egress-Tabelle im Code (`featureCatalogue`, `egressFeatureHosts`) kennen `ollama-host` nicht — Beleg ist die ADR (*Stand im Code und im Vertrag*, am Symbol gelesen vom Architect am 2026-10-24) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) ist `Accepted` (eingetreten: Commit `d9dd171e`) und [`slice-v1-adr-ollama-cloud-modelle`](../done/slice-v1-adr-ollama-cloud-modelle.md) liegt in `done/` (eingetreten); kein anderer Slice in `in-progress/` (WIP-Limit 1; die Slices ohne Welle in `open/` teilen sich das Limit, gelesen am 2026-10-24: `in-progress/` trägt nur `carveouts.md`, `README.md`, `roadmap.md`). Die Freigabe des Wortlauts ist **kein** Start-Trigger, sondern Schritt 3 des Ablaufs (§3) und Bedingung für den Lastenheft-Commit. Die Version wird beim Start vergeben (Schritt 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Freigabe verlangt eine Zusage, die eine neue Entscheidung oder Messung braucht (etwa ein Anschluss hinter der Egress-Restriktion, ein Relay durch `u-boot`, eine Zusage für `arm64` oder Colima, ein Lauf mit gültigem `OLLAMA_API_KEY`); dann zuerst eine Folge-ADR.
- `in-progress` → `open` (blockiert): Die Freigabe des Wortlauts liegt nicht vor, oder der Projektinhaber verlangt einen anderen Wortlaut als den Entwurf, der Rückfragen nach sich zieht, oder er stuft die Läufe mit Konto als nicht tragend ein (Ausgang B, Folge-ADR streicht Z1 bis Z4, Z7 bis Z10 teilweise); §7 nennt die offene Stelle.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; der Wortlaut im Lastenheft ist mit der Freigabe Satz für Satz verglichen (Verification-Report), der Review-Report (und ein Nachlauf-Review, falls eine Fix-Runde Lastenheft-Text änderte) liegt vor, und die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Wortlaut sagt mehr zu, als die Zusage-Menge trägt.** Die Läufe mit Konto laufen über Handvarianten (Pfad je Zeile Z1 bis Z4, Z7 bis Z10), die Menge ist ein Lauf für Weg (b) und zwei für Weg (a); „Cloud-Modelle“ allgemein (Frage F4) und „Zugesagt ist das auf einem Linux-Host …“ (Sätze 3.7 und 4.8) dehnen die Aussage auf Modelle und Läufe aus, die nicht einzeln gemessen sind ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 24 Belege; [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)). — **Ausgang:** weiter offen bis zur Freigabe; Review und Verifier vergleichen jeden Satz gegen die Tabelle *Herkunft jedes Satzes* und gegen den Abschnitt *Nicht gemessen, nicht zugesagt* der ADR.
- **Zwei Zusage-Mengen im selben Feature** (Installation von `ollama` mit Messstand Colima und `arm64`, Cloud-Weg mit Linux `amd64`), solange [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) nicht gelaufen ist (Fragen F2 und F3). — **Ausgang:** weiter offen; die Spannung ist im Wortlaut-Vorschlag ausgewiesen, und [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) löst sie mit der Entscheidung zu seiner Frage 1.
- **Das Lastenheft sagt zu, was Code und Spezifikation noch nicht tragen** (Zwischenzustand nach dem CR: Lastenheft sagt `ollama-host` zu, der Katalog kennt ihn nicht; Zeile 739 der Spezifikation nennt den Anschluss weiter als „nicht zugesagt“). — **Ausgang:** weiter offen bis zur Welle-Closure; vorgesehene Reihenfolge ([`ADR-0036`](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) *Folgepflicht*), der Release-Cut liegt nach Code und Spezifikation, der Zwischenzustand erreicht keinen Release.
- **Eine Bedingung des Wortlauts ohne Test** (Sätze 3.3, 3.6, 4.1, 4.2, 4.6, 4.7, Stellen 5 und 6; [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), 24 Belege, [`AGENTS.md`](../../../../AGENTS.md) §3.10). — **Ausgang:** weiter offen; die Beleg-Tabelle in §7 nennt je Bedingung die Zusage-Zeile oder den Auftrag an den Code-Slice, die Aufträge stehen in der Welle-Datei §5.
- **Der Wortlaut ändert sich nach dem Review.** Jede Fix-Runde mit Lastenheft-Diff bekommt ein Nachlauf-Review über genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), 15 Belege). — **Ausgang:** weiter offen; Nachlauf-Review als eigene Datei vor dem Verifier (§2).
- **Die Freigabe-Vorlage verdichtet den Entwurf** ([`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md), 1 Beleg). — **Ausgang:** weiter offen; der Ablauf (§3 Schritte 2 und 3) verlangt die Vorlage ungekürzt, die Verifikation vergleicht das Lastenheft gegen §7 **und** gegen den Entwurf in §1.
- **Die Freigabe fällt anders aus als der Entwurf** (etwa `ollama.com` umschrieben, der Halbsatz „Gemessen ist das auf …“ gestrichen, Modelle genannt). Die Freigabe in §7 trägt den geänderten Wortlaut, die Beleg-Tabelle folgt ihm; trägt der geänderte Wortlaut eine Zusage über die Zusage-Menge hinaus, ist das ein Rückführungsgrund (§4). — **Ausgang:** weiter offen bis zur Freigabe.
- **Versionsfolge mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md).** Beide CR ändern [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); dieser nimmt die nächste freie Version beim Start, jener die folgende (Entscheidung des Standards, §1). — **Ausgang:** weiter offen bis zum Start; Angaben jenes Slice und seiner Zeile in `open/README.md` werden beim Start nachgezogen (§3 Schritt 1).
- **Eine Aussage über den Code oder das Lastenheft im Plan stimmt nicht** (Zeilennummern, Symbole, Zustand der Welle). Abgleich der Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), 54 Belege, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10). — **Ausgang:** weiter offen; die Zeilenangaben sind am Stand `0.22.0` gelesen und werden beim Start neu bestimmt (§3 Schritt 1).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Reihenfolge dieser Sektion: Version und Start · Freigabe je Stelle (Rückfrage beim Projektinhaber) · Vorab-Suche mit Fundstellen · Beleg-Tabelle (Bedingung → Zusage-Zeile oder Test) · Verifikation mit Exit-Codes · Closure.

### Version und Start

Start am 2026-10-24 (`date +%F`, Ortszeit). Versionsfeld gelesen: `0.22.0` (Zeile 9), letzte Historie-Zeile `0.22.0`. Vergeben: **`0.23.0`**. `in-progress/` trug vorher nur `carveouts.md`, `README.md`, `roadmap.md`. Angaben von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) (Version, Zeilen 412, 472, 474, 1588, 1593, 1594, 1596) und seiner Zeile in `open/README.md` sind nachgezogen.

### Freigabe je Stelle

Projektinhaber, 2026-10-24, je Stelle eigene Rückfrage, Antwort „Freigeben wie vorgelegt“ (Vorlage ungekürzt mit *Alt* und *Neu*):

| Stelle | Entscheidung | Wortlaut |
|---|---|---|
| 1, 2, 3, 4, 5, 6 | Ja | wie in §1, Zeichen für Zeichen; Links als Anker innerhalb der Datei |
| 7 | Ja | Version beim Start vergeben (`0.23.0`), Historie-Zeile mit Datum 2026-10-24, sonst Wortlaut §1 |
| 8 | Ja (keine Änderung) | wie in §1 |

Fragen: F1 `ollama.com` genannt, `host.docker.internal` nicht; F2 zwei Zusage-Mengen bis zum Messstände-CR; F3 Halbsatz „Gemessen ist das auf …“ bleibt; F4 „Cloud-Modelle“ ohne Namen; F5 Zusatzwarnung nicht gesondert unter `--quiet`; F6 Sätze ohne eigene Z-Zeile bleiben, Beleg ist der Test des Code-Slice (alle Standard). **Abweichung gegen den Entwurf in §1:** keine.

### Vorab-Suche (Implementer, 2026-10-24)

`grep -n -i -F` über `spec/lastenheft.md` auf `Anschluss an einen Ollama`, `ollama-host`: Treffer nur die neuen Stellen (Zeilen 410, 418, 1604) und die Chronik-Zeile `0.15.0` (Zeile 1596, Klasse D, bleibt). Dieselben Suchwörter über `harness/`, `docs/user/`, `docs/maintainer/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `spec/spezifikation.md`, `spec/architecture.md`: Treffer außerhalb des Slice, bewusst bis zu den Folge-Slices stehen geblieben: `spec/spezifikation.md` Zeile 739 ([`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), „Nicht zugesagt. GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, …“, widerspricht dem Lastenheft bis `slice-v2-spec-ollama-cloud`), `docs/user/benutzerhandbuch.md` Zeile 2043 (gleiche Aufzählung) und Zeile 2023 („Cloud-Modelle sind nicht enthalten“), `README.de.md` Zeile 14 („Cloud-Modelle von Ollama sind nicht enthalten“). Weitere Spiegelstellen im Lastenheft (Befehlsübersicht, Feature-Aufzählungen, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)): keine Aussage über den Anschluss oder Cloud-Modelle gefunden (Suche `Cloud`, `Ollama`, `Anschluss` im Lastenheft: nur die genannten Stellen und Zeile 1394 „Cloud-Provisioning“).

### Ist-Stand der Änderung

Geänderte Stellen in `spec/lastenheft.md` (Zeilen am Stand `0.23.0`): Stelle 1 Zeile 410; Stelle 2 Zeile 412 (Aufzählung); Stelle 3 Zeile 414 (neuer Absatz); Stelle 4 Zeile 418 (neuer Absatz); Stelle 5 Zeile 460 ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)); Stelle 6 Zeile 191 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1); Stelle 7 Zeile 9 und Zeile 1604.

*(Beleg-Tabelle, Verifikation und Closure werden von Reviewer, Verifier und bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (Lastenheft, Kürzel DOK) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area), dazu `docs/plan/` (DOK); beide Greenfield. Kein Code, keine zu grobe Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler = Dateien in `evidence/`, gelesen am 2026-10-24). Treffer in DOK: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (24, Schwelle erreicht, Ausgang beim Projektinhaber; §6 Risiko 1), [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (15; §2 Nachlauf-Review, §6), [`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md) (1, unter der Schwelle; §3 Schritte 2 und 3, §6), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (Pfad je Zusage-Zeile; §6), [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (4; die Stellen 3 und 4 benennen den Nutzerpfad in beiden Profilen: `ollama` nur Sandbox, `ollama-host` beide), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (24) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (54), beide in [`AGENTS.md`](../../../../AGENTS.md) §3.10 verkörpert; dieser Plan nennt die Fundstellen (§1 Tabelle, §3 Tabelle). Evidenz-/Diskrepanz-Risiko: mittel (Zusage-Menge, Spannung der zwei Mengen im Feature `ollama`).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
