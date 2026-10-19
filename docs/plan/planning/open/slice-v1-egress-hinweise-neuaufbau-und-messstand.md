# Slice v1: Egress-Texte — Hinweis auf den Neuaufbau bei der Listenänderung und Messstand der Subdomain-Aussage (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Der Slice ist ein eigener Code-Slice **nach** dem Release der Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) (Übergabe der Closure dieser Welle, Lücken L-2 und L-3 des Validators); er hat keine Closure-Bedingung, die von seiner DoD verschieden wäre.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion: Anzeige, Subdomain-Aussage), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (was `--quiet` zeigt); [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Neuaufbau statt Neu-Laden, `Accepted`).
Anlass: Validator-Bericht [`2026-10-19-validation-welle-sandbox-egress-ollama`](../../../reviews/2026-10-19-validation-welle-sandbox-egress-ollama.md), Lücken **L-2** (Listenänderung wirkt erst nach Neuaufbau, `generate devcontainer` sagt es nicht) und **L-3** (Subdomain-Texte der Ausgabe ohne Messstand), und die Beobachtung [`BEO-CLI/subdomain-text-ohne-vorbehalt`](../observations/BEO-CLI/subdomain-text-ohne-vorbehalt/observation.md) (Verification V-2 des Spezifikations-Slice).

**Berührte Spec-Stellen:** [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) (Egress-Anzeige; trägt die `Notes:`-Zeilen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Subdomain-Aussage und Messstand), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Verhalten unter `--quiet`). Der Slice **liest** sie; der Architect liest ihren Inhalt (hier sind nur die Titel gelesen) und benennt, welche Stelle eine Nachführung braucht. Der Verweis zeigt aufwärts: Die Spec nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Wer die Egress-Liste ändert und `generate devcontainer` ausführt, liest an dieser Stelle, dass ein bestehender Container die Liste seines Images behält und neu aufgebaut werden muss (L-2); und die Texte, die die Subdomain-Aussage ohne Vorbehalt tragen, nennen ihren Messstand oder sagen nur, was gemessen ist (L-3, V-2).

### Gefundener Stand (gelesen am 2026-10-19)

- **`generate devcontainer` meldet die Änderung ohne Hinweis:** `printGenerateSummary` (`internal/adapter/driving/cli/generate.go`, Zeilen 335 bis 368) schreibt für `GenerateActionUpdatedBlock` „Updated %s managed block (%s).“ mit Artefaktname und Dateiliste, unabhängig davon, welche Datei sich geändert hat; die Funktion kennt den Grund der Änderung nicht (sie bekommt `driving.GenerateResponse` und `dryRun`). Der Validator sah beim Lauf nach einer Listenänderung „Updated devcontainer managed block (.devcontainer/egress-init.sh)“.
- **Die Subdomain-Texte stehen ohne Vorbehalt in vier Stellen:** die Langhilfe und die Zeile `Notes:` von `u-boot egress list` („A name in the list also allows its subdomains.“, `internal/adapter/driving/cli/egress.go`, Zeilen 99 und 213; gepinnt durch `TestExecute_EgressList_Notes` in `egress_test.go`), die Fehlermeldung zu einem Eintrag mit `*` („subdomains of an allowed name are already included“, `internal/hexagon/domain/sandbox.go`, `NewEgressHost`, Zeilen 134 und 137; gepinnt durch `TestConfigSet_EgressAllowWildcard_Exit10WithHint` in `config_egress_wildcard_test.go`) und der Kommentar „an entry never opens its parent name“ in `internal/e2e/sandbox_devcontainer_docker_test.go` (Zeile 197).
- **Die Ausgabe von `egress list` nennt den Neuaufbau nicht:** Der Block `Notes:` (`internal/adapter/driving/cli/egress.go`, Zeilen 212 bis 220) sagt, dass ein Container die Liste des Images behält, aus dem er gebaut wurde, und nennt `devcontainer up --remove-existing-container`; die Wörter „rebuild“ und „Neuaufbau“ stehen nicht im Text. `TestExecute_EgressList_Notes` pinnt den Behalt und den Satz zum `--remove-existing-container`. Handbuch (Punkt „Nach einer Änderung führen Sie …“) und CHANGELOG (`[Unreleased]`, Abschnitt Fixed, Eintrag zum Neuaufbau nach einer Änderung der Liste) tragen den Neuaufbau als Ist-Zustand. Die Spezifikation ([`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), *Abgrenzung zum laufenden Container*) sagt über die Ausgabe nur den Behalt und den Aufruf; die Zusage „eine geänderte Liste braucht einen Neuaufbau“ ist dort gestrichen, weil kein Test sie pinnt (Übergabe aus `slice-v2-spec-zielbild`, Verification V-1).
- **Was Spezifikation und Handbuch dazu sagen:** den Messstand und die Grenzen (Eintrag `www.example.org` ungemessen, Verbindung über eine geteilte Adresse nicht geprüft). Diese Zeilen sind hier nicht neu gelesen; der Implementer liest sie vor der Textwahl.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Neu-Laden der Liste im laufenden Container** — entschieden gegen (R0, [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)); der Hinweis nennt den Neuaufbau, er ändert ihn nicht.
- **Ein Satz zu Cloud-Modellen im Handbuch (Lücke L-1)** — Handbuch-Text ohne Codeänderung; er steht als Auftrag im Release-Cut der Welle und in der Folge-Welle `welle-ollama-cloud-modelle`, nicht hier.
- **Messungen** zur Subdomain-Tiefe (`a.b.example.org`) oder zur geteilten Adresse — eigener Messvorgang; hier wird gekürzt oder der vorhandene Messstand genannt, nichts Neues zugesagt (`AGENTS.md` §3.10).
- **Anforderungs-Kennungen als Präfix der Warnungen (Lücke L-9)** — Vorbestand außerhalb der Welle, eigener Vorgang.

## 2. Definition of Done

- [ ] **Hinweis zum Neuaufbau:** `generate devcontainer` sagt bei einer Änderung von `.devcontainer/egress-init.sh`, dass ein bestehender Container die Liste seines Images behält und neu aufgebaut werden muss; ein Test pinnt den Fall (Änderung der Liste) und den Gegenfall (Änderung ohne `egress-init.sh`, NoOp, `--dry-run`); die Behandlung unter `--quiet` ist vor dem Code entschieden und getestet ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): ohne Ausnahme entfällt ein Hinweis unter `--quiet`).
- [ ] **Neuaufbau in der Ausgabe von `egress list`:** Die `Notes:` nennen den Neuaufbau mit dem Wort (Neuaufbau beziehungsweise rebuild) neben dem Behalt der Liste; `TestExecute_EgressList_Notes` pinnt den Wortlaut. Danach darf [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) die Hälfte „eine geänderte Liste braucht einen Neuaufbau“ wieder zusagen; der Slice führt sie in der Spezifikation nach (ein Vertrag, vier Stellen).
- [ ] **Subdomain-Texte:** die vier Stellen aus §1 nennen den Messstand oder sagen nur das Gemessene; die Tests (`TestExecute_EgressList_Notes`, `TestConfigSet_EgressAllowWildcard_Exit10WithHint`) pinnen den neuen Wortlaut; Spezifikation, Handbuch und CHANGELOG sagen über die Texte dasselbe wie der Code.
- [ ] `make gates` grün; Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review; Doku-Update für Handbuch und CHANGELOG, soweit der öffentliche Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/generate.go` (`printGenerateSummary`) und, falls der Grund der Änderung in der Antwort fehlt, `internal/hexagon/port/driving/generate.go` samt Anwendungsschicht | update | der Hinweis braucht die Information, dass `egress-init.sh` geändert wurde; ob sie aus `resp.Changed` abgeleitet oder in der Antwort getragen wird, entscheidet der Architect |
| `internal/adapter/driving/cli/generate_acceptance_test.go` (Block „Human-Mode“) | update | Fall Listenänderung, Gegenfall, `--dry-run`, `--quiet` |
| `internal/adapter/driving/cli/egress.go`, `internal/adapter/driving/cli/egress_test.go` (`TestExecute_EgressList_Notes`) | update | Langhilfe und `Notes:` mit Messstand beziehungsweise gekürzt; `Notes:` nennen den Neuaufbau (Pin) |
| `internal/hexagon/domain/sandbox.go` (`NewEgressHost`), `internal/adapter/driving/cli/config_egress_wildcard_test.go` | update | Wortlaut des Wildcard-Fehlers |
| `internal/e2e/sandbox_devcontainer_docker_test.go` (Kommentar, Zeile 197) | update | Kommentar auf das Gemessene |
| `spec/spezifikation.md` ([`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) und die vom Architect benannten Stellen), `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Ein Vertrag, vier Stellen (Welle-Kriterium der Vorgänger-Welle: Text im Code gleich Text in Spec, Handbuch, CHANGELOG) |
| `harness/replay.md` und Golden-Dateien | update, falls die Ausgabe von `generate` oder `egress list` in einem Golden Case steht | Replay-Evidence nach `harness/replay.md` |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; der Release der Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) ist geschnitten (der Slice ändert Code und Texte, die der Release-Cut mit der Welle ausliefert); der Projektinhaber hat vor dem Start entschieden, (a) ob die Subdomain-Texte einen Vorbehalt bekommen oder gekürzt werden, und (b) ob der Hinweis zum Neuaufbau unter `--quiet` entfällt. Kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Grund der Änderung ist in der Antwort von `generate` nicht vorhanden und braucht eine Port-Änderung durch mehrere Schichten; dann trennt der Architect den Hinweis (Schnittstelle) von den Texten (Kürzung).
- `in-progress` → `open` (blockiert): Die Entscheidung des Projektinhabers zu (a) oder (b) fehlt.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der Hinweis erscheint bei jeder Änderung von `egress-init.sh`, auch wenn kein Container existiert.** — **Ausgang:** weiter offen bis zum Architect-Entwurf; der Wortlaut sagt „ein bestehender Container“ und lässt offen, ob einer existiert.
- **Die Antwort von `generate` trägt den Grund der Änderung nicht.** — **Ausgang:** weiter offen; Rückführung `in-progress` → `next` in §4.
- **Ein Vorbehalt im Wortlaut erscheint dem Nutzer als Unsicherheit der Funktion, die gemessen funktioniert.** — **Ausgang:** weiter offen; die Entscheidung (a) in §4 liegt beim Projektinhaber.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** voraussichtlich berührt werden `internal/adapter/driving/cli` (CLI, Greenfield), `internal/hexagon/domain` (DOM, Greenfield), `internal/hexagon/application` und `internal/hexagon/port` (APP und PORT, Greenfield, nur falls die Antwort erweitert wird), `internal/e2e` (E2E, Greenfield; ein Kommentar) und `spec/`, `docs/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle aus zwei von drei Achsen laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-CLI/subdomain-text-ohne-vorbehalt`](../observations/BEO-CLI/subdomain-text-ohne-vorbehalt/observation.md) (offen, 1×; dieser Slice ist die Antwort), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (offen, 2×; jeder Vorbehalts-Satz wird gegen das Messprotokoll gelesen, nicht gegen die Erinnerung), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen, zwei Belege; Text im Code, Spezifikation, Handbuch und CHANGELOG sind zusammen nachzuziehen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide verkörpert in `AGENTS.md` §3.10). Evidenz-/Diskrepanz-Risiko: niedrig bis mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
