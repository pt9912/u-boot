# Slice v1: Egress — Subdomains und Wildcard-Eintrag prüfen: messen, belegen, entscheiden (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — erster Slice der Welle (Analyse; hängt von keinem anderen Slice ab). Ihm folgen
[`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) und [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md), deren Anzeige die hier belegte Subdomain-Aussage trägt.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; `egress.allow` ist eine „Liste von Hostnamen“),
[`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set; Entscheidung 2 nennt „samt Subdomains“).
Anlass: Projektinhaber, 2026-10-18 („in sandbox modus kann man egress nicht erweitern/konfigurieren — geht das besser?“, Vorschlag (c) Wildcards `*.example.org`).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Zeile *Nutzer*: „kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard“) — der Slice **liest** sie und ändert die Spezifikation nicht;
eine Änderung nimmt [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) vor. Das Lastenheft ändert der Slice nicht.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine Messung belegt, ob und wie Subdomains eines erlaubten Namens im Sandbox-Container aufgelöst und erreicht werden, und beantwortet die Frage des Projektinhabers: **Gibt es einen Fall, der ohne Wildcard-Eintrag nicht geht?**
Das Ergebnis ist entweder *entfallen* (kein solcher Fall; „entfallen“ ist ein gültiges Ergebnis) oder ein benannter Fall mit Entscheidungsvorlage. Die Aussage „Subdomains eines erlaubten Namens sind mitgemeint“ steht danach auf einer Messung und einem Test, oder sie wird gestrichen.

### Gefundener Stand im Code und in der Doku (gelesen am 2026-10-18)

- **Eingabe:** `domain.NewEgressHost` (`internal/hexagon/domain/sandbox.go`) lässt nur Namen zu, die `egressHostPattern` entsprechen (kleingeschrieben, Labels aus `a-z0-9-`, **mindestens ein Punkt**); der Fehlertext nennt „no scheme, port, path or wildcard“. `validateEgressHosts` und `normaliseEgressHosts` (`internal/hexagon/application/devcontainer_egress.go`) rufen sie je Eintrag auf; `setEgressAllow` (dieselbe Datei) hängt neue Hosts an die bestehende Liste an.
- **Mechanismus:** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` schreibt je Host zwei Zeilen in `/etc/u-boot-egress-dnsmasq.conf`: `server=/$h/$ns` (je Upstream-Resolver) und `nftset=/$h/4#ip#u_boot_egress#allowed`; andere Namensregeln gibt es nicht. Der Namensabgleich ist damit Sache von `dnsmasq`.
- **Die Aussage über Subdomains** steht an zwei Stellen: `docs/plan/adr/0012-devcontainer-egress-firewall.md` (Entscheidung 2: „nur die erlaubten Namen (samt Subdomains)“) und `docs/user/benutzerhandbuch.md` (Abschnitt *Internet-Beschränkung (Egress)*, Punkt *Wichtig*: „Subdomains eines erlaubten Namens sind mitgemeint“).
- **Beleg dafür: keiner gefunden.** `TestE2E_SandboxDevcontainer_Egress` (`internal/e2e/sandbox_devcontainer_docker_test.go`) prüft einen erlaubten Host (`github.com`) und einen gesperrten (`example.com`); die Suche nach `subdomain` unter `internal/` (Quelltext und Tests) hat keinen Treffer; die Messung von [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) nennt Subdomains nicht. Das ist die Klasse von [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (`AGENTS.md` §3.10: eine Bedingung in Spec oder Handbuch hat einen Test oder wird gestrichen).
- **Hypothese des Planners, nicht belegt** (die Messung prüft sie, sie steht nicht als Tatsache im Handbuch): Der Abgleich von `server=/<name>/` und `nftset=/<name>/` in `dnsmasq` gilt für den Namen und alle Subdomains; dann wäre `*.example.org` inhaltlich `example.org` (der Name selbst eingeschlossen), und ein Wildcard-Eintrag fügte keine Ausdrucksstärke hinzu, höchstens eine andere Schreibweise.

### Messfragen (Arbeitsplan; Docker auf Colima, Wegwerf-Container aus dem erzeugten Dockerfile mit Egress, wie die Messreihe von [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md))

| # | Frage | Warum |
|---|---|---|
| a | Eintrag `example.org`: lösen und erreichen `example.org`, `www.example.org` und eine tiefe Subdomain (`a.b.example.org`)? | belegt oder widerlegt die Aussage in ADR und Handbuch |
| b | Eintrag nur `api.github.com` (ohne Eltern-Namen): bleibt `github.com` gesperrt, und lösen Subdomains von `api.github.com` auf? | die Gegenrichtung: ein Eintrag öffnet nie den Eltern-Namen |
| c | Ein erlaubter Name, dessen DNS-Antwort per `CNAME` auf einen Namen **außerhalb** der Allowlist zeigt: wird die Adresse am Ende der Kette in das nftables-Set aufgenommen und der Verkehr erreicht sie? | der wahrscheinlichste Fall, der ohne Eintrag eines weiteren Namens nicht ginge; entscheidet, ob ein Wildcard-Eintrag etwas lösen könnte, was ein einfacher Eintrag nicht kann |
| d | Lässt sich „nur Subdomains, nicht der Name selbst“ in der `dnsmasq`-Konfiguration überhaupt ausdrücken? | wäre der einzige Inhalt, den ein `*.`-Eintrag über den Eintrag des Namens hinaus hätte |

### Verzweigung nach der Messung (vorab benannt)

- **E1 — bestätigt, kein Fall:** Der Slice endet ohne neue Zusage. Die Aussage bekommt einen Test (Fall a, b), das Handbuch nennt den Messstand, und der Fehlertext von `NewEgressHost` sagt bei einem Eintrag mit `*` ausdrücklich, dass Subdomains eines erlaubten Namens enthalten sind und der Name ohne `*.` einzutragen ist (nur ein anderer Text zu demselben Exit-Code `10`, kein neuer Vertrag). Wildcard-Syntax: *entfallen* (Begründung in §7).
- **E2 — Aussage falsch oder nur teilweise wahr** (etwa Subdomains tiefer als eine Ebene): Handbuch und der Folge-Eintrag in der Spezifikation werden korrigiert; [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) bleibt unverändert `Accepted` (`AGENTS.md` §3.5), die Korrektur steht in einer Folge-ADR (nächste freie Nummer; geprüft 2026-10-18: `0024`) oder, wenn die Messung nichts Entscheidungsbedürftiges findet, in der Evidence. Ob dadurch ein Wildcard-Eintrag nötig wird, entscheidet die Messung d.
- **E3 — echter Fall** (zum Beispiel Frage c oder d zeigt etwas, das nur ein Wildcard-Eintrag löst): Der Slice liefert eine **Entscheidungsvorlage** an den Projektinhaber (Optionen: Wildcard-Eintrag als Alias auf den Namen, Wildcard mit eigener Semantik, Eintrag des genauen Namens als Anleitung) und **keinen Code**; die Umsetzung wäre ein eigener Slice nach Lastenheft-Entscheidung (Welle §5, Eintrag).

Die Hosts der Ollama-Messung ([`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md)) können einen Fall zeigen (zum Beispiel Modell-Blobs unter wechselnden Subdomains eines Speicher-Dienstes); zeigt sie einen, reicht der Ollama-Slice ihn an diesen Gegenstand zurück (Welle §5), er ändert die Wildcard-Entscheidung nicht selbst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Wildcard-Syntax im Code** (`*.example.org` annehmen) — eine neue Eingabe-Zusage ist eine Entscheidung des Projektinhabers (E3), keine Folge der Analyse; ohne belegten Fall entsteht sie nicht.
- **Eine geänderte Allowlist-Mechanik** (Proxy, andere DNS-Logik) — [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative D (Restriktion außerhalb des Containers) ist ein eigener Vorgang.
- **Die Anzeige der wirksamen Liste** — [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md); sie nutzt das Ergebnis dieses Slice für den Hinweis „Subdomains eingeschlossen“.
- **Das Neu-Laden im laufenden Container** — [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md).
- **Lastenheft und Spezifikation** — der Slice ändert keine von beiden (Lastenheft ist Zielbild; die Spezifikation folgt im Spezifikations-Slice der Welle).

## 2. Definition of Done

- [x] **Messung:** Die Fälle a bis d sind gemessen und mit Messstand (Docker-Version, Colima-Version, Architektur, `dnsmasq`-Version aus dem Image) in der Evidence und in §3 festgehalten; das Ergebnis ist E1, E2 oder E3 (§1), jeweils mit der gemessenen Begründung.
- [x] **Belegte Aussage:** Bei E1 und E2 steht die Subdomain-Aussage in Handbuch (`docs/user/benutzerhandbuch.md`, Abschnitt Egress) im Wortlaut der Messung, und ein Test (`//go:build docker`) pinnt sie oder sie ist gestrichen; bei E1 trägt der Fehlertext von `NewEgressHost` den Hinweis, mit Test.
- [x] **Ergebnis für die Welle:** §7 trägt `Gegenstand:` (entfallen mit Begründung, oder die benannte Vorlage an den Projektinhaber), und die Welle ([`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §5) ist nachgeführt, falls E3 einen Slice in die Welle zieht.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für das Handbuch, falls ein öffentlicher Vertrag berührt ist; `CHANGELOG.md` (`[Unreleased]`) nur bei geändertem Fehlertext.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — Beleg in [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), falls die Subdomain-Aussage ohne Beleg stand (Befund steht dann in der Evidence); sonst „keine Beobachtung angefallen“ in §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | neuer Test neben `TestE2E_SandboxDevcontainer_Egress` (Build-Tag `docker`): Fälle a und b (`TestE2E_SandboxDevcontainer_EgressSubdomains`, eingetreten), nach [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion); Fall c nicht getestet (Messung ohne stabiles Ziel, nur Evidence) |
| `internal/hexagon/domain/sandbox.go` | update (eingetreten, E1) | Fehlertext von `NewEgressHost` für Einträge mit `*` |
| `internal/hexagon/domain/sandbox_test.go` | update (eingetreten, E1; `TestNewEgressHost`, `TestNewEgressHost_WildcardLengthAndRawQuote`) | Test des neuen Fehlertextes (Negativfall `*.example.org`); die Grenze von 253 Zeichen im Wildcard-Zweig und das Zitat der rohen (ungetrimmten) Eingabe |
| `internal/adapter/driving/cli/config_egress_wildcard_test.go` | neu (eingetreten) | `config set` mit Wildcard-Eintrag über die Produktionskette (`ConfigService`, echtes Dateisystem und YAML-Codec): Exit-Code `10`, Hinweistext, `u-boot.yaml` unverändert; ein Paar (`setEgressAllow`) und mehrere Paare (`planEgressAllow`) |
| `docs/user/benutzerhandbuch.md` | update (eingetreten) | Abschnitt *Internet-Beschränkung (Egress)*: Punkte *Subdomains*, *Namen hinter einem CNAME* (auf die drei gemessenen Fälle begrenzt) und *Messstand* |
| `CHANGELOG.md` | update (eingetreten, E1) | `[Unreleased]`: geänderter Fehlertext bei Wildcard-Eintrag |
| `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` | unverändert | gelesen: `if [ -f "$pidfile" ]; then kill "$(cat "$pidfile")" 2>/dev/null \|\| true; fi` beendet den alten `dnsmasq` ohne Warten (Startrennen); Entscheidung in §6 |
| `docs/plan/adr/` (Folge-ADR, Index `docs/plan/adr/README.md`) | entfallen (E1, Aussage bestätigt) | nur wenn die Messung die Aussage von [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) widerlegt |
| `docs/reviews/2026-10-18-slice-v1-egress-wildcard-pruefung.md` | neu (eingetreten) | Review-Report nach `harness/review.md`; Verification Evidence folgt mit der Closure (`harness/verification.md`) |
| `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md` | neu (eingetreten) | Beleg, dass die Aussage ohne Test stand, mit Messung, Auswertung der sauberen Läufe und Erklärung des „Address already in use“-Artefakts |
| `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/messung-slice-v1-egress-wildcard-pruefung/` (`run.sh`, `inner.sh`, `inner-cd.sh`, `rohdaten-1.txt` bis `-3.txt`) | neu (eingetreten) | Messskripte und Rohdaten im Repo, neben der Evidence-Datei (Fälle a bis d reproduzierbar) |
| `docs/plan/planning/welle-sandbox-egress-ollama.md`, `docs/plan/planning/open/README.md` | update (Verweis-Nachzug durch `make slice-mv`; Closure) | Pfad des Slice nach `in-progress/` und bei der Closure nach `done/`; Welle-Tabelle (§4) nennt das Ergebnis E1 statt „erwartet“, `open/README.md` führt den Slice als geliefert |
| `docs/plan/planning/done/slice-cr-sandbox-egress-anzeige.md`, `…/slice-v1-cli-egress-anzeige.md`, `…/slice-v1-devcontainer-feature-ollama.md`, `…/slice-v1-devcontainer-ollama-messung.md`, `…/slice-v1-sandbox-egress-neuladen-entwurf.md` | update (Verweis-Nachzug durch `make slice-mv`) | Pfad des Slice nach `in-progress/`; ihre Aussagen zur Subdomain-Messung sind gelesen (am 2026-10-18): Vorbehalte und Verweise auf dieses Ergebnis (E1), keine Behauptung über das Handbuch; der Satz mit Bedingung in `slice-cr-sandbox-egress-anzeige` („gilt nur, wenn … E1 bestätigt“) und das zugehörige Risiko (§6) stehen bei der Closure auf dem Ist-Stand (E1 bestätigt, Wortlaut unverändert) |
| `docs/plan/planning/open/slice-v2-spec-sandbox-egress-ollama.md` | update (Verweis-Nachzug und berichtigte Aussage) | Pfad nach `in-progress/`; die Aussage, der Subdomain-Satz stehe nur im Handbuch, ist gegen `docs/user/benutzerhandbuch.md` gelesen und berichtigt |
| `docs/reviews/2026-10-18-verification-slice-v1-egress-wildcard-pruefung.md` | neu (eingetreten) | Verification-Report (GO) mit Auflagen A-1 bis A-4 für die Closure |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/` (`state.md`, `evidence/slice-v1-egress-wildcard-pruefung.md`) | update, neu (eingetreten) | fünfzehnter Vorgang (Review F-4 bis F-6), Zählbasis im Stand |
| `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/state.md` | update (eingetreten) | vierter Vorgang in der Zählbasis (Evidence-Datei besteht seit der Messung) |
| `docs/plan/planning/observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/` (`observation.md`, `state.md`, `evidence/slice-v1-egress-wildcard-pruefung.md`) | neu (eingetreten) | neue Beobachtung aus Review F-2, Zähler 1× |
| diese Slice-Datei | update; Move nach `done/` | Messung (§3), Risiken (§6), Marker (§8), Closure (§7); reiner Move bei der Closure |
| `docs/plan/planning/in-progress/roadmap.md` | update | Arbeitszeile statt Ruhe-Marker solange der Slice in `in-progress/` liegt; Ruhe-Marker bei der Closure wiederherzustellen |

### Messung (2026-10-18)

**Messstand:** Docker 29.5.2 (Server, arm64) auf Colima 0.10.3 (macOS Virtualization.Framework, aarch64); Wegwerf-Container `debian:stable` mit `--cap-add=NET_ADMIN`, `dnsmasq` 2.91 (Paket `dnsmasq-base` aus Debian stable), `nftables`; ausgeführt wurde das Startscript `egress-init.sh.tmpl` mit ersetzter `hosts=`-Zeile, sonst unverändert. Skripte, Rohdaten und Auswertung der sauberen Läufe (samt Erklärung des „Address already in use“-Artefakts in der Rohdatei) liegen im Repo, neben der Evidence: [`messung-slice-v1-egress-wildcard-pruefung/`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/messung-slice-v1-egress-wildcard-pruefung/run.sh) und [Evidence-Datei](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md).

| # | Eintrag | Befund | Bewertung |
|---|---|---|---|
| a | `example.org` | `example.org` und `www.example.org`: DNS `NOERROR`, HTTPS 200; `a.b.example.org`: DNS `NXDOMAIN` (an den Upstream weitergereicht, Name existiert nicht; ein nicht passender Name wie `example.com` ergibt `REFUSED`) | Subdomain erlaubt. Tiefe Subdomain mit echtem Namen: Eintrag `amazonaws.com` → `s3.amazonaws.com` (307), `s3.us-east-1.amazonaws.com` (307, zwei Ebenen), `sts.us-east-1.amazonaws.com` (302) erreichbar, `example.com` `REFUSED` |
| b | `api.github.com` | `api.github.com` HTTPS 200; `github.com` und `www.github.com`: `REFUSED`; `uploads.api.github.com`: `NXDOMAIN` (also weitergereicht, Regel greift) | Eintrag öffnet nie den Eltern-Namen; Subdomains des Eintrags sind eingeschlossen |
| c | `www.mozilla.org` (CNAME auf `www-mozilla.fastly-edge.com`), `files.pythonhosted.org` (CNAME auf `dualstack.python.map.fastly.net`), `deb.debian.org` (CNAME auf `debian.map.fastlydns.net`) | alle drei erreichbar (HTTPS 200, 404 für den Wurzelpfad von `files.pythonhosted.org`); die Zielnamen selbst: `REFUSED` | Alle drei Ketten enden beim selben CDN-Anbieter (Fastly), ein Lauf: Die Adresse am Ende der Kette gelangt in das nftables-Set. Der Zielname braucht keinen Eintrag und löst direkt nicht auf. Kein Fall für einen Wildcard-Eintrag |
| d | `*.example.org` und `.example.org` | `dnsmasq --test` nimmt beide an. `server=/*.example.org/…`: `example.org` `REFUSED`, `www.example.org` `NOERROR`. Mit dem vollen Startscript (`nftset=/*.example.org/…`) löst `www.example.org` auf, die Verbindung läuft in den Timeout (kein Eintrag im Set); `.example.org` verhält sich wie `example.org` | „Nur Subdomains“ ist für die Auflösung ausdrückbar, **nicht** für die Freigabe der Adresse: der `nftset`-Treffer fehlt. Mit `dnsmasq` 2.91 ergibt sich daraus keine nutzbare Form |

**Ergebnis: E1** (bestätigt, kein Fall). Ein Wildcard-Eintrag fügte weder für a bis c noch für d (mit diesem `dnsmasq`) Ausdrucksstärke hinzu. Nicht gemessen: andere `dnsmasq`-Versionen, Linux-Hosts ohne VM, Docker Desktop, amd64, IPv6 (das Script sperrt IPv6 ganz).

Die Tabelle ist nach der letzten Änderung (Closure: Move nach `done/`, Roadmap, `open/README.md`, Register) gegen `git diff --stat 1cec55f` abgeglichen (`AGENTS.md` §3.10; 32 Dateien, der Move erscheint als Neuanlage in `done/` und Löschung in `open/`): jede berührte Datei hat eine Zeile, jede Zeile einen Diff, ausgenommen die Einträge „unverändert“ und „entfallen“; es steht keine Zeile mit „bedingt“ mehr in ihr. Eine spätere Änderung wiederholt den Abgleich.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zieht eine Änderung des Mechanismus nach sich (E3 mit Umsetzung); die Umsetzung ist dann ein eigener Slice.
- `in-progress` → `open` (blockiert): die Messumgebung (Docker auf Colima mit `NET_ADMIN`) steht nicht zur Verfügung; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); das Ergebnis (E1, E2 oder E3) steht in §7 und in der Welle.

## 6. Risiken und offene Punkte

- **Die Messung unter Colima gilt nur dort** (Docker auf Colima, arm64; `dnsmasq`-Version des Images): Linux-Hosts und Docker Desktop sind nicht gemessen. — **Ausgang:** weiter offen; die Grenze steht im Handbuch (Punkt *Messstand*), kein Register-Eintrag, solange kein zweiter Fall auftritt.
- **Ein `CNAME`-Fall (Frage c) lässt sich ohne öffentliches Messziel nicht stabil testen**: ein Test hinge an fremdem DNS. — **Ausgang:** eingetreten; der Fall ist nur gemessen (drei Hosts hinter demselben CDN, ein Lauf, Evidence), nicht getestet. Das Handbuch begrenzt die Aussage auf die drei gemessenen Fälle und nennt, dass es keinen Test gibt (Punkt *Namen hinter einem CNAME*). Ein Test mit lokalem `dnsmasq`-Upstream und `cname=`-Eintrag wurde nicht gebaut: der Test-Container hat keinen zweiten Resolver, die Kette des Startscripts nimmt die Upstreams aus `/etc/resolv.conf`, und der Aufbau ginge über die Messfrage hinaus; ein Folgeslice kann ihn nachholen, falls die Zusage gebraucht wird.
- **Die Messung findet, dass die Aussage in [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) nicht stimmt** (E2): [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) ist `Accepted` und unveränderlich. — **Ausgang:** entfallen; die Messung bestätigt die Aussage (E1), keine Folge-ADR.
- **Die Ollama-Messung zeigt später einen Fall** ([`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md)), nachdem dieser Slice „entfallen“ gemeldet hat. — **Ausgang:** weiter offen; der Rückweg steht in §1 und in der Welle (§5), kein stilles Wiederöffnen.
- **Startrennen im Startscript** (Review F-9): `egress-init.sh.tmpl` beendet beim zweiten Lauf den alten `dnsmasq` mit `kill` und startet den neuen sofort (`bind-interfaces`); ein verlorenes Rennen färbt `B example.org` in `TestE2E_SandboxDevcontainer_EgressSubdomains` rot (die alte Instanz antwortet noch), nie falsch-grün. — **Ausgang:** weiter offen, bewusst nicht behoben: Das Script ist Produktcode im Generator-Output (Replay-/Golden-Fälle und ein 13-minütiger Docker-Lauf stünden für eine Änderung an), das Neu-Laden im Container ist Gegenstand von `slice-v1-sandbox-egress-neuladen-entwurf` (dort gehört der geordnete Neustart hin). Im einzigen Lauf des Tests (`make test-docker`, Stand der Umsetzung, grün) trat das Rennen nicht auf; ein Lauf belegt keine Flakefreiheit. Tritt der Test flakig rot, ist die Ursache hier benannt.
- **Der Verifier und das Handbuch ohne Spec-Anker** (Review F-10): Die Handbuch-Zusagen zu Subdomains, `CNAME` und Wildcard-Abweisung stehen bis zum Folgeslice `slice-v2-spec-sandbox-egress-ollama` ohne Anker in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion). — **Ausgang:** weiter offen; der Spec-Anker kommt im Folgeslice (Nicht-Ziel dieses Slice: Spezifikation unverändert).
- **Wächter-Grenze des e2e-Tests** (Review F-8): Der Status `000` unterscheidet Allowlist nicht von sonstigem Netzfehler und der Test hängt an fremdem DNS und am Image-Cache. — **Ausgang:** weiter offen; falsch-grün ist ausgeschlossen (die Positivfälle im selben Lauf belegen das Netz), falsch-rot möglich; Register `BEO-E2E/test-gruen-nur-mit-lokalem-cache`.

## 7. Closure-Notiz

**Ergebnis: E1** (Aussage bestätigt, kein Fall). Messstand: Docker 29.5.2 (arm64) auf Colima 0.10.3, `dnsmasq` 2.91, Wegwerf-Container mit dem unveränderten Startscript; Skripte und Rohdaten im Repo (Messprotokoll unter [`evidence/messung-slice-v1-egress-wildcard-pruefung/`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/messung-slice-v1-egress-wildcard-pruefung/run.sh)). Die Fälle stehen in §3 *Messung*; in Kürze: (a) ein Eintrag schließt `www.` und tiefere Subdomains ein (`amazonaws.com` bis `s3.us-east-1.amazonaws.com`, zwei Ebenen); (b) ein Eintrag `api.github.com` öffnet `github.com` nicht für die Namensauflösung; (c) drei Namen mit `CNAME` auf einen Namen außerhalb der Liste sind erreichbar, der Zielname braucht keinen Eintrag; (d) „nur Subdomains“ ist für die Auflösung ausdrückbar, nicht für die Freigabe der Adresse (der `nftset`-Treffer fehlt), eine nutzbare Form ergibt sich mit `dnsmasq` 2.91 nicht.

- **Was hat funktioniert:** Die Messung lief gegen das echte Startscript statt gegen eine Nachbildung, und die Hypothese des Planners (Abgleich von `server=/<name>/` und `nftset=/<name>/` gilt für den Namen samt Subdomains) stand als Hypothese im Plan, nicht im Handbuch; die Aussage kam erst nach der Messung dorthin. Der Review fand vor der Closure zwei MEDIUM (F-1 CNAME-Aussage breiter als der Beleg, F-2 Messprotokoll außerhalb des Repos); beide sind behoben, bevor der Verifier die Sensoren selbst fuhr (`make gates` 91,20 %, `make doc-trace` 0 Waisen, `make test-docker` grün).
- **Was ging anders als geplant:** Die erste Fassung des Handbuchs sagte zu `CNAME` mehr, als drei Hosts hinter einem CDN in einem Lauf tragen (F-1), und „gesperrt“ meinte nur die Namensauflösung (F-3); beides steht jetzt begrenzt auf das Gemessene im Handbuch. Der Projektinhaber entschied nach der Messung: kein Wildcard-Alias, stattdessen ein erklärender Fehlertext (`NewEgressHost`, derselbe Exit-Code `10`). Die Plan-Tabelle führte zuerst die Link-Nachzüge von `make slice-mv` nicht, und ein Satz in `slice-v2-spec-sandbox-egress-ollama` über das Handbuch war nach dem Diff veraltet (F-4, F-5, F-6).
- **Gegenstand:** entfallen: die Wildcard-Syntax (`*.example.org` als Eingabe) — kein Fall in a bis d: (a) und (b) belegen, dass der Name allein Subdomains einschließt und den Eltern-Namen nicht öffnet, (c) braucht keinen weiteren Eintrag, (d) ergibt mit `dnsmasq` 2.91 keine Form, die der Eintrag des Namens nicht schon hätte; ein Alias wäre nur eine andere Schreibweise, und der Projektinhaber hat ihn nicht beauftragt. Kommt aus der Ollama-Messung ([`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md)) ein Subdomain-Fall, geht er laut Welle §5 an diesen Gegenstand zurück, nicht still in die Default-Allowlist.
- **Steering-Loop-Eintrag:** gezählt, nicht verkörpert; mit diesem Slice wurde keine Regel, kein Sensor und keine Spec-Lücke verkörpert. Zwei Einträge im Register: (1) `BEO-DOK/spec-bedingung-ohne-code-entsprechung`, vierter Vorgang: Die Aussage „Subdomains eines erlaubten Namens sind mitgemeint“ stand seit dem Egress-Slice (2026-09-30) im Handbuch und in [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) ohne Test und war wahr; die Beleg-Lücke ist damit Altbestand vor der Bedingungs-Klausel in `AGENTS.md` §3.10, kein Rückfall nach ihr. (2) `BEO-DOK/plan-aussage-ueber-code`, fünfzehnter Vorgang und Rückfall nach der Verkörperung (F-4 Plan-Tabelle ohne die Link-Nachzüge, F-5 Zustandsmarker „bedingt“ nach Ablauf, F-6 Aussage über das Handbuch in einem anderen Slice veraltet). Ob eine weitere Schärfung oder ein Sensor folgt, ist nicht entschieden (Entscheidung des Projektinhabers, siehe Bericht).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-v1-egress-wildcard-pruefung.md` in `BEO-DOK/spec-bedingung-ohne-code-entsprechung/` (vorhanden) und in `BEO-DOK/plan-aussage-ueber-code/` ergänzt, `state.md` beider Einträge nachgeführt; neu angelegt `BEO-DOK/messprotokoll-ausserhalb-des-repos/` (F-2: Messdaten nur außerhalb des Repos, ein Beleg, Zähler 1×, offen). `BEO-E2E/test-gruen-nur-mit-lokalem-cache` bekommt keine Evidence: Die Wächter-Grenze des neuen e2e-Tests (F-8) ist kein lokal-grün/CI-rot-Fall.
- **Folge-Slices:** keine neuen. Vorhanden und betroffen: [`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) (der Subdomain-Satz des Wortlauts trägt jetzt die Messung), [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) (trägt den Spec-Anker für Subdomains, `CNAME` und Wildcard-Abweisung), [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md) (Startrennen im Startscript).
- **Grenzen (nur gemessen, was hier steht):** Colima/arm64 und `dnsmasq` 2.91; Linux-Hosts ohne VM, Docker Desktop, amd64 und IPv6 (das Script sperrt IPv6 ganz) nicht gemessen; ein Lauf je Fall gegen fremdes DNS. Fall c (`CNAME`) ist nur gemessen (drei Hosts, ein CDN-Anbieter), nicht getestet. **N-2:** Die Aussage, dass bei geteilter IP-Adresse der Eltern-Name per Adresse offen ist, ist nur an identischen Adressen in den Rohdaten (`rohdaten-1.txt`, `example.org` und `www.example.org`) belegt; kein Verbindungstest über die Adresse, sie steht im Handbuch als Bedingung („teilen sich eine Adresse“) und nicht als Zusage über Namen. Der e2e-Test unterscheidet mit dem Status `000` die Allowlist nicht von sonstigem Netzfehler (falsch-grün ausgeschlossen, falsch-rot möglich); das Startrennen im Startscript ist in zwei grünen Läufen nicht aufgetreten, das belegt keine Flakefreiheit. Mutationen an `NewEgressHost` und am e2e-Test sind nicht gefahren (Review und Verifier lasen sie); `make verify-depguard`, `make doc-immutable` und `make ci` liefen nicht (kein Import-, ADR- oder Abhängigkeits-Diff). Die Handbuch-Zusagen stehen bis [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) ohne Anker in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion).
- **Verifikation:** [Verification-Report](../../../reviews/2026-10-18-verification-slice-v1-egress-wildcard-pruefung.md) (GO), [Review-Report](../../../reviews/2026-10-18-slice-v1-egress-wildcard-pruefung.md) (F-1 bis F-10 behoben oder begründet).
- **Risiken aus §6:** jedes trägt in §6 genau einen Ausgang: eingetreten (CNAME-Fall), entfallen (Aussage der ADR falsch), weiter offen (Messung gilt nur unter Colima, Ollama-Rückgabe, Startrennen, Handbuch ohne Spec-Anker, Wächter-Grenze).
- **Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: die drei genannten Dateien liegen in `open/`. Register: `BEO-DOK/spec-bedingung-ohne-code-entsprechung`, `BEO-DOK/plan-aussage-ueber-code`, `BEO-DOK/messprotokoll-ausserhalb-des-repos` und `BEO-E2E/test-gruen-nur-mit-lokalem-cache` sind Verzeichnisse im Register.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/domain` (DOM, Greenfield; Fehlertext von `NewEgressHost`), `internal/e2e` (E2E, Greenfield; ein Test) und `docs/user/`, `docs/plan/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle aus zwei von drei Achsen (eigenes Verzeichnis, eigene Konvention, eigene Inventur-Linie) laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Stand `verkörpert` in `AGENTS.md` §3.10; die Subdomain-Aussage ist ein möglicher weiterer Vorgang) und [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (offen, ein Beleg; der neue Test baut ein Image und ist ein weiterer e2e-Test mit Image-Cache-Abhängigkeit). Evidenz-/Diskrepanz-Risiko: mittel, weil die Aussage ohne Beleg im Handbuch steht.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
