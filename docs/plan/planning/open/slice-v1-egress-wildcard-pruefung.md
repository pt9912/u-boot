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
[`slice-cr-sandbox-egress-anzeige`](slice-cr-sandbox-egress-anzeige.md) und [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md), deren Anzeige die hier belegte Subdomain-Aussage trägt.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; `egress.allow` ist eine „Liste von Hostnamen“),
[`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set; Entscheidung 2 nennt „samt Subdomains“).
Anlass: Projektinhaber, 2026-10-18 („in sandbox modus kann man egress nicht erweitern/konfigurieren — geht das besser?“, Vorschlag (c) Wildcards `*.example.org`).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Zeile *Nutzer*: „kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard“) — der Slice **liest** sie und ändert die Spezifikation nicht;
eine Änderung nimmt [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md) vor. Das Lastenheft ändert der Slice nicht.

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

Die Hosts der Ollama-Messung ([`slice-v1-devcontainer-ollama-messung`](slice-v1-devcontainer-ollama-messung.md)) können einen Fall zeigen (zum Beispiel Modell-Blobs unter wechselnden Subdomains eines Speicher-Dienstes); zeigt sie einen, reicht der Ollama-Slice ihn an diesen Gegenstand zurück (Welle §5), er ändert die Wildcard-Entscheidung nicht selbst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Wildcard-Syntax im Code** (`*.example.org` annehmen) — eine neue Eingabe-Zusage ist eine Entscheidung des Projektinhabers (E3), keine Folge der Analyse; ohne belegten Fall entsteht sie nicht.
- **Eine geänderte Allowlist-Mechanik** (Proxy, andere DNS-Logik) — [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative D (Restriktion außerhalb des Containers) ist ein eigener Vorgang.
- **Die Anzeige der wirksamen Liste** — [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md); sie nutzt das Ergebnis dieses Slice für den Hinweis „Subdomains eingeschlossen“.
- **Das Neu-Laden im laufenden Container** — [`slice-v1-sandbox-egress-neuladen-entwurf`](slice-v1-sandbox-egress-neuladen-entwurf.md).
- **Lastenheft und Spezifikation** — der Slice ändert keine von beiden (Lastenheft ist Zielbild; die Spezifikation folgt im Spezifikations-Slice der Welle).

## 2. Definition of Done

- [ ] **Messung:** Die Fälle a bis d sind gemessen und mit Messstand (Docker-Version, Colima-Version, Architektur, `dnsmasq`-Version aus dem Image) in der Evidence und in §3 festgehalten; das Ergebnis ist E1, E2 oder E3 (§1), jeweils mit der gemessenen Begründung.
- [ ] **Belegte Aussage:** Bei E1 und E2 steht die Subdomain-Aussage in Handbuch (`docs/user/benutzerhandbuch.md`, Abschnitt Egress) im Wortlaut der Messung, und ein Test (`//go:build docker`) pinnt sie oder sie ist gestrichen; bei E1 trägt der Fehlertext von `NewEgressHost` den Hinweis, mit Test.
- [ ] **Ergebnis für die Welle:** §7 trägt `Gegenstand:` (entfallen mit Begründung, oder die benannte Vorlage an den Projektinhaber), und die Welle ([`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §5) ist nachgeführt, falls E3 einen Slice in die Welle zieht.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für das Handbuch, falls ein öffentlicher Vertrag berührt ist; `CHANGELOG.md` (`[Unreleased]`) nur bei geändertem Fehlertext.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — Beleg in [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), falls die Subdomain-Aussage ohne Beleg stand (Befund steht dann in der Evidence); sonst „keine Beobachtung angefallen“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | neuer Test neben `TestE2E_SandboxDevcontainer_Egress` (Build-Tag `docker`): Fälle a und b, nach [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion); Fall c nur, wenn die Messung einen reproduzierbaren Namen liefert |
| `internal/hexagon/domain/sandbox.go` | update (bedingt, E1) | Fehlertext von `NewEgressHost` für Einträge mit `*` |
| `internal/hexagon/domain/sandbox_test.go` | update (bedingt, E1) | Test des neuen Fehlertextes (Negativfall `*.example.org`) |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt *Internet-Beschränkung (Egress)*, Punkt *Wichtig* (Subdomain-Aussage mit Messstand) |
| `CHANGELOG.md` | update (bedingt, E1) | `[Unreleased]`: geänderter Fehlertext bei Wildcard-Eintrag |
| `docs/plan/adr/` (Folge-ADR, Index `docs/plan/adr/README.md`) | neu / update (bedingt, E2) | nur wenn die Messung die Aussage von [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) widerlegt |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/` | neu (bedingt) | Beleg, falls die Aussage ohne Test stand und das Ergebnis es bestätigt |
| diese Slice-Datei | update | Messung (§3), Closure (§7) |

Der Abgleich der Tabelle gegen `git diff --stat` läuft nach der **letzten** Fix-Runde (`AGENTS.md` §3.10); jede Zeile mit „bedingt“ wird dabei als eingetreten oder entfallen markiert.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zieht eine Änderung des Mechanismus nach sich (E3 mit Umsetzung); die Umsetzung ist dann ein eigener Slice.
- `in-progress` → `open` (blockiert): die Messumgebung (Docker auf Colima mit `NET_ADMIN`) steht nicht zur Verfügung; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); das Ergebnis (E1, E2 oder E3) steht in §7 und in der Welle.

## 6. Risiken und offene Punkte

- **Die Messung unter Colima gilt nur dort** (Docker auf Colima, arm64; `dnsmasq`-Version des Images): Linux-Hosts und Docker Desktop sind nicht gemessen. — **Ausgang:** offen bis zur Closure; erwartet: weiter offen, Grenze im Handbuch benannt (kein Register-Eintrag, solange kein zweiter Fall auftritt).
- **Ein `CNAME`-Fall (Frage c) lässt sich ohne öffentliches Messziel nicht stabil testen**: ein Test hinge an fremdem DNS. — **Ausgang:** offen bis zur Messung; erwartet: nur Messung in der Evidence, kein Test (Grenze ausweisen).
- **Die Messung findet, dass die Aussage in [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) nicht stimmt** (E2): [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) ist `Accepted` und unveränderlich. — **Ausgang:** offen bis zur Messung; eingetreten wäre: Folge-ADR (`0024`), nicht Umschreiben.
- **Die Ollama-Messung zeigt später einen Fall** ([`slice-v1-devcontainer-ollama-messung`](slice-v1-devcontainer-ollama-messung.md)), nachdem dieser Slice „entfallen“ gemeldet hat. — **Ausgang:** offen; der Rückweg steht in §1 und in der Welle (§5), kein stilles Wiederöffnen.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Messergebnis a–d, E1/E2/E3, `Gegenstand:` bei „entfallen“, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/domain` (DOM, Greenfield; Fehlertext von `NewEgressHost`, bedingt), `internal/e2e` (E2E, Greenfield; ein Test) und `docs/user/`, `docs/plan/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle aus zwei von drei Achsen (eigenes Verzeichnis, eigene Konvention, eigene Inventur-Linie) laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Stand `verkörpert` in `AGENTS.md` §3.10; die Subdomain-Aussage ist ein möglicher weiterer Vorgang) und [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (offen, ein Beleg; der neue Test baut ein Image und ist ein weiterer e2e-Test mit Image-Cache-Abhängigkeit). Evidenz-/Diskrepanz-Risiko: mittel, weil die Aussage ohne Beleg im Handbuch steht.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
