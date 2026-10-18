# Welle welle-sandbox-egress-ollama: Egress im Sandbox-Profil komfortabler und Ollama als Katalog-Feature

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld.
**Geplante Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter
*Nächste Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug. Die Zuordnung zu einem Release-Cut (nach v0.11.0) liegt beim Projektinhaber und ist offen; die Welle ist keine Release-Version.

**Verantwortlich:** pt9912. **Datum:** 2026-10-18.

---

## 1. Welle-Ziel

Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) lässt sich die Egress-Allowlist ([`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) besser überblicken und pflegen, und der Katalog bietet Ollama als wählbaren Schlüssel `ollama`:

1. **Anzeige:** Ein Befehl zeigt die wirksame Allowlist mit der Herkunft jedes Hosts (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste), offline und ohne laufenden Container; dieselbe Berechnung füllt das Startscript.
2. **Subdomains und Wildcards:** Die Aussage „Subdomains eines erlaubten Namens sind mitgemeint“ steht auf einer Messung und einem Test; ob ein Wildcard-Eintrag (`*.example.org`) einen Fall löst, den der Name allein nicht löst, ist beantwortet — „entfallen“ ist ein gültiges Ergebnis.
3. **Neu-Laden:** Optionen, ihre Sicherheitsfolgen und eine Messung auf Colima liegen dem Projektinhaber vor; er hat R0 gewählt (2026-10-18): kein Neu-Laden, der Weg „Liste ändern“ ist der Neuaufbau des Containers; die Vorlage ([`ADR-0025`](../adr/0025-egress-liste-im-laufenden-container.md)) sagt auch, was ein Prozess im Container mit und ohne `sudo` an der Beschränkung ändern kann.
4. **Ollama:** `ollama` ist ein Katalog-Feature für den Sandbox-Devcontainer: CPU-only, gepinnte und geprüfte Installation (kein ungeprüftes `curl | sh`), Modelle im Container, Hosts des Modell-Pull in der Default-Allowlist — jeweils nur, soweit gemessen und im Lastenheft zugesagt.

Anlass (Projektinhaber, 2026-10-18): „in sandbox modus kann man egress nicht erweitern/konfigurieren — geht das besser?“ und, auf die vier Vorschläge, „Ok deine Vorschläge für die Verbesserungen“ (a Anzeige, b Neu-Laden, c Wildcards, d mehr Agenten) sowie „zu 4. ollama“ mit der Antwort „Feature: Ollama im Container“. Vorschlag (d) „mehr Agenten im Katalog“ entsteht **nicht** als Slice: Der Projektinhaber hat außer Ollama keinen Agenten genannt.

## 2. Trigger (Welle startet)

- Der Projektinhaber hat die Welle am 2026-10-18 beauftragt (eingetreten): Egress-Komfort (Anzeige, Neu-Laden, Wildcard-Prüfung) und Ollama als Katalog-Feature, als eigene Welle.
- Kein Slice liegt in `in-progress/` (WIP-Limit 1; beobachtbar an der Verzeichnis-Position und am Vermerk *Nichts in Arbeit* der Roadmap, eingetreten). Der erste Slice der Welle wird erst danach nach `in-progress/` bewegt.
- Die Zuordnung zu einer Release-Version ist **kein** Trigger (offen, Projektinhaber).

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den Slice-DoDs: Die Egress-Verbesserungen und `ollama` tragen **zusammen**, und die Stellen, an denen der Vertrag steht, sagen dasselbe.

- Alle Slices der Welle liegen in `done/` (Tabelle §4; der Entwurfs-Slice hat keine Folge-Slices, weil R0 gewählt ist); die CR-Slices liegen dort mit ihren Lastenheft-Versionen.
- **Eine Aussage über die Allowlist:** Lastenheft, Spezifikation ([`SPEC-019`](../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) mit der Subdomain-Aussage), Handbuch (Abschnitt *Internet-Beschränkung (Egress)*) und CHANGELOG sagen über die Anzeige, die Subdomain-Semantik, den Weg „Liste ändern“ (Neubau des Containers, gemessen: ein Neustart genügt nicht; Neu-Laden entfällt nach der Entscheidung R0 vom 2026-10-18) und die Hosts von `ollama` dasselbe wie der Code (`make doc-trace` und `make docs-check` ohne Befund).
- **Anzeige gleich Startscript:** Ein Test zeigt für Standard-, `podman`-, `docker`-, Feature-, Clone- und Nutzer-Fälle, dass die Liste des Anzeige-Befehls gleich der `hosts`-Zeile von `egress-init.sh` ist; ein Golden Case mit `ollama` (Sandbox-Profil, Egress, auch zusammen mit `claude-code`) belegt, dass die Hosts von `ollama` in beiden erscheinen, `devcontainer.json`, Dockerfile und `egress-init.sh` stabil sind und der zweite Lauf ein NoOp ist.
- **Lauffähigkeit belegt, nicht angenommen:** Messungen sind benannt und ihr Messstand steht in den Evidences (Subdomain-Verhalten, Neuaufbau statt Neu-Laden, Installation und Lauf von Ollama hinter der Egress-Restriktion; Docker auf Colima, arm64); ein Docker-Test (`//go:build docker`) belegt, was der CR für `ollama` zusagt, oder die Closure nennt den Grund, warum er nicht lief (kein stilles Rot). Was nicht wiederholbar ist (Messprotokoll ohne Skript), ist als Grenze ausgewiesen.
- **Zusage nur für Gemessenes:** Jede Bedingung in Lastenheft, Handbuch und CHANGELOG zu `ollama` und zur Subdomain-Aussage hat einen Test oder steht dort nicht (`AGENTS.md` §3.10).
- Die Spezifikation trägt die Einträge zu Anzeige und Katalog und ist **nach** dem Abschluss aller Code- und CR-Slices entstanden und gegen den Code gelesen (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise; `make docs-check` meldet 0 Befunde; `make doc-immutable` meldet für die Range der Welle 0 Befunde (Accepted-ADRs unverändert; Korrekturen stehen in Folge-ADRs); `make gates` ist grün.
- Closure-Notiz in §7 mit Lerneintrag; die Entscheidungen des Projektinhabers (Neu-Laden, Lieferkette von Ollama, Release-Zuordnung) sind im Gespräch belegt und in §7 vermerkt.

## 4. Slices in dieser Welle

Die Reihenfolge steht in §5; die Verzeichnis-Position nennt den Zustand, keine Status-Spalte.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-v1-egress-wildcard-pruefung`](done/slice-v1-egress-wildcard-pruefung.md) | Analyse und Messung: Subdomains eines erlaubten Namens, Wildcard-Eintrag — gibt es einen Fall, der ohne ihn nicht geht (Ergebnis E1: Aussage bestätigt und durch einen Docker-Test gepinnt, Wildcard-Syntax entfallen, erklärender Fehlertext) | [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| [`slice-cr-sandbox-egress-anzeige`](open/slice-cr-sandbox-egress-anzeige.md) | CR: Anzeige der wirksamen Egress-Liste mit Herkunft je Host (Lastenheft 0.14.0) | [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| [`slice-v1-cli-egress-anzeige`](open/slice-v1-cli-egress-anzeige.md) | Code: Befehl für die wirksame Egress-Liste (offline, `--json`), eine Berechnung für Anzeige und Startscript | [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| [`slice-v1-sandbox-egress-neuladen-entwurf`](in-progress/slice-v1-sandbox-egress-neuladen-entwurf.md) | Entwurf und Messung: Neu-Laden im laufenden Container, Sicherheits-Optionen (mit und ohne `sudo`), Folge-ADR (Ergebnis: R0, Guardrail-Status quo, kein Neu-Laden; ein Neustart genügt nach einer Listenänderung nicht, Neuaufbau nötig; Folge-ADR [`ADR-0025`](../adr/0025-egress-liste-im-laufenden-container.md) `Accepted`; keine Folge-Slices) | [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) |
| [`slice-v1-devcontainer-ollama-messung`](done/slice-v1-devcontainer-ollama-messung.md) | Messung und Quellenwahl: Ollama (CPU-only) im Sandbox-Devcontainer, Lieferkette, Hosts, Modellspeicher, Start; Folge-ADR (Ergebnis: Quelle B, Release-Archiv mit SHA-256-Pin, Folge-ADR [`ADR-0024`](../adr/0024-ollama-im-sandbox-devcontainer.md) Accepted) | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| [`slice-cr-devcontainer-ollama`](open/slice-cr-devcontainer-ollama.md) | CR: Katalog erweitert um `ollama`, Zusage nur für Gemessenes (nächste freie Lastenheft-Version) | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) |
| [`slice-v1-devcontainer-feature-ollama`](open/slice-v1-devcontainer-feature-ollama.md) | Code: Katalog-Schlüssel `ollama`, Installation nach ADR, Hosts, Modell-Volume, Start, Golden Case, Docker-Test, Doku | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) |
| [`slice-v2-spec-sandbox-egress-ollama`](open/slice-v2-spec-sandbox-egress-ollama.md) | Spezifikation: Egress-Anzeige, Subdomain-Aussage, Katalog und Hosts von `ollama`, Nachführung von [`SPEC-043`](../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Sicht (zuletzt, nach dem Code) | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |

Folge-Slices zum Neu-Laden gibt es nicht: Der Projektinhaber hat am 2026-10-18 R0 gewählt; ein Lastenheft-CR und ein Code-Slice zum Neu-Laden (R1/R2) entstehen nicht.

**Nicht Teil der Welle:** [`slice-v1-egress-dns-upstream-schliessen`](open/slice-v1-egress-dns-upstream-schliessen.md) (Pfad zum Upstream-Resolver für Prozesse ohne `sudo`, Befund M1 a des Entwurfs-Slice) ist ein eigener Slice nach dieser Welle (Entscheidung des Projektinhabers, 2026-10-18); er hält die Closure der Welle und das Release nicht auf, und kein Closure-Trigger in §3 hängt an ihm.

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: keiner offenen Welle; der Start verlangt, dass kein Slice in `in-progress/` liegt (§2).
- Release-Zuordnung: offen, Entscheidung des Projektinhabers; kein Slice der Welle hängt an ihr.
- **Reihenfolge in der Welle:**
  - [`slice-v1-egress-wildcard-pruefung`](done/slice-v1-egress-wildcard-pruefung.md) zuerst: Er hängt von nichts ab, ist klein, und sein Ergebnis (Subdomain-Aussage belegt oder korrigiert) steht im Wortlaut des Anzeige-CR.
  - [`slice-cr-sandbox-egress-anzeige`](open/slice-cr-sandbox-egress-anzeige.md) vor [`slice-v1-cli-egress-anzeige`](open/slice-v1-cli-egress-anzeige.md): Der Code setzt den Befehlsnamen des freigegebenen Lastenhefts um (Lastenheft ist Zielbild, Code folgt).
  - [`slice-v1-sandbox-egress-neuladen-entwurf`](in-progress/slice-v1-sandbox-egress-neuladen-entwurf.md) nach der Anzeige: Der Slice hängt nicht an der Anzeige und ändert keinen Code; ein Neu-Laden-Befehl entsteht nach R0 nicht.
  - [`slice-v1-devcontainer-ollama-messung`](done/slice-v1-devcontainer-ollama-messung.md) nach den Egress-Slices, damit ein Subdomain-Fall aus den Hosts des Modell-Pull an einen geschlossenen Wildcard-Slice **zurückgereicht** wird: Zeigt die Messung einen Fall, den der Name allein nicht löst, öffnet das den Gegenstand des Wildcard-Slice neu — als Eintrag hier, nicht still.
  - [`slice-cr-devcontainer-ollama`](open/slice-cr-devcontainer-ollama.md) nach der Messung (Zusage nur für Gemessenes), [`slice-v1-devcontainer-feature-ollama`](open/slice-v1-devcontainer-feature-ollama.md) nach dem CR und der ADR (`Accepted`).
  - [`slice-v2-spec-sandbox-egress-ollama`](open/slice-v2-spec-sandbox-egress-ollama.md) zuletzt: Die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und Code es tragen.
- **Lastenheft-Versionen:** jeder CR nimmt die nächste freie Minor-Version zum Zeitpunkt des Lastenheft-Commits (Stand 2026-10-18: `0.13.0`; der Anzeige-CR `0.14.0`, danach der Ollama-CR; ein Neu-Laden-CR entsteht nach R0 nicht). Der Projektinhaber gibt jeden Wortlaut vor dem Lastenheft-Commit frei.
- **Zwei Slices berühren dieselben Dateien** (`devcontainer_egress.go`, Handbuch-Abschnitt Egress, `harness/replay.md`): WIP-Limit 1 macht die Reihenfolge der Tabelle verbindlich.
- **Mögliche Rückführung:** Ein Rückführungs-Ausgang des Entwurfs- oder Mess-Slice kann einen weiteren Slice in die Welle ziehen; die Zahl der Slices ändert sich dann durch einen Eintrag in §4, nicht stillschweigend.

## 6. Out-of-Scope für diese Welle

- **Anschluss an den Ollama des Hosts** (Host-Dienst, `host.docker.internal`, Socket): vom Projektinhaber ausgeschlossen; widerspräche [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (keine Host-Verbindung im Profil).
- **Weitere Agenten im Katalog** (Codex, Gemini, andere): Der Projektinhaber hat außer Ollama keinen genannt; jede Quelle braucht eine eigene Entscheidung (Vorschlag (d) entsteht nicht als Slice).
- **GPU-Beschleunigung für Ollama**: Entscheidung des Projektinhabers (CPU-only).
- **Auswahl, Vorbelegung oder Verwaltung von Modellen durch `u-boot`** und **Feature-Optionen** (`pull`): der Katalog kennt keine Optionen ([`ADR-0020`](../adr/0020-drittanbieter-features-im-katalog.md), Re-Evaluierungs-Trigger); der Nutzer holt Modelle im Container.
- **Ollama als Compose-Add-on-Service**: gewählt ist „Feature: Ollama im Container“.
- **Eine Wildcard-Syntax in `egress.allow`** ohne belegten Fall: Sie entsteht nur durch eine Entscheidung des Projektinhabers (Wildcard-Slice, Ergebnis E3).
- **Eine Restriktion außerhalb des Containers** (Proxy, Netz-/DNS-Ebene) und **das Entziehen von `sudo`** für den Container-Benutzer: eigene Vorgänge ([`ADR-0012`](../adr/0012-devcontainer-egress-firewall.md) Alternative D; Option S des Entwurfs-Slice).
- **Eine Änderung der Lockerungen des Sandbox-Profils**, des Basis-Images oder der `nestedRuntime`-Werte: Sicherheitsentscheidung des Projektinhabers; verlangt die Messung eine, hält der betroffene Slice an (Rückführung nach `open/`).
- **Docker Desktop, Linux-Host und amd64 als Messumgebung**: nicht gemessen, nicht zugesagt; jede Aussage nennt ihren Messstand (Colima, arm64).
- **Der Release-Cut und die Zuordnung der Welle zu ihm**: Sache des Projektinhabers, eigener Slice nach dieser Welle, wenn er ihn wünscht.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt — nach dem letzten Slice. Die Pfade dieses Abschnitts lösen vom Ruheort `done/` auf (Ruheort-Regel). Inhalt: Geliefert, Closure-Trigger (§3) je Punkt, Was hat funktioniert, Was ging anders als geplant, Trigger-Audit (Carveout, Gate, ADR, Hard Rule), Steering-Loop-Einträge, Lücken zwischen Bedarf und Zusage (Validator), Bestand, Beobachtungs-Register, Folge-Slices, Auftrag an den Release-Cut, offene Entscheidungen beim Projektinhaber, Verifikation, drei Paarungen. Vorlage: [`welle-devcontainer-sshd-claude-code`](done/welle-devcontainer-sshd-claude-code.md) §7.)*
