# Slice v1: Degradations-Fallbacks (Daemon, Egress) host-seitig ausweisen — oder `LH-FA-DEV-007` präzisieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — kein Release-Bezug, keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.
Herkunft: [`welle-cli-verbosity`](../done/welle-cli-verbosity.md), Slice
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (zweiter Review F-4, Verifier-Auflage 3).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Absatz Degradation: jeder Fallback wird ausgewiesen),
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 2),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Wegfall der Egress-Restriktion als Fallback).
Review:
[`2026-10-16-slice-v1-cli-quiet-nur-fehler-2`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md) (F-4) und
[Verification Evidence](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-nur-fehler.md) (Bewertung Entscheidung 3).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) — nur im CR-Zweig
(Lösung B), dann mit Kennung `slice-cr-sandbox-fallback-ausweisung`; im Code-Zweig (Lösung A) wird das Lastenheft nicht geändert.
`spezifikation.md` und `architecture.md` folgen der gewählten Lösung nach der Closure.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Lücke zwischen [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) („jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen“) und der Umsetzung ist
geschlossen — entweder im Code oder im Vertrag.

**Gefundener Stand im Code (Ausgangslage):**

- Produktiv erreichbar ist nur der `vfs`-Fall: `doctor` weist ihn über das Item `devcontainer.sandbox.runtime` aus
  (`internal/hexagon/application/doctor_sandbox.go`, wenn `/dev/fuse` fehlt und die Policy `warn` gilt).
- „Container ohne Docker-Daemon“ und „Wegfall der Egress-Restriktion“ entstehen im Startscript im Container
  (`internal/hexagon/application/templates/devcontainer/`); der Host kennt sie nicht, `doctor` hat kein Item dafür.
- Für `init`, `generate` und `up` erzeugt keine produktive Quelle einen Eintrag mit `DegradationFallback`; die Pins laufen über Fake-Use-Cases.
- Die Lücke besteht unabhängig von `--quiet`: Die Option verbirgt nichts, was der Host ermittelt (Ausnahme 2 von
  [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) hält ausgewiesene Fallbacks sichtbar).

**Entscheidung des Projektinhabers (2026-10-17): Lösung B.** Zwei Lösungen standen zur Wahl; der Projektinhaber hat „B: Lastenheft präzisieren“ entschieden.

- **Lösung A — Code (verworfen).** Der Fallback wird host-seitig ermittelbar, etwa meldet das Startscript ihn an den Host (Statusdatei oder Volume), oder
  `doctor` liest den Container-Zustand; danach erzeugen `init`/`generate`/`up`/`doctor` die Einträge.
- **Lösung B — CR (gewählt).** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) wird auf das host-seitig Ermittelbare präzisiert; der Gegenstand liegt im CR
  [`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md) (dieser Slice geht mit `Gegenstand:` dorthin über).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung am Verhalten von `--quiet`** — die Option verursacht die Lücke nicht; Ausnahme 2 und ihre Pins stehen im Slice
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md).
- **Die Ausgabe der Startscripte selbst** — sie laufen im Container und kennen weder `--quiet` noch den Host-Prozess; Gegenstand ist die Ausweisung auf dem Host.
- **Neue Fallback-Klassen oder Sicherheitslockerungen** — eine neue Zusage; [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) zählt drei Fallbacks.
- **Die Vorschau von `--dry-run` unter `--quiet`** — anderer Befund
  ([`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md)).

## 2. Definition of Done

- [ ] Die Entscheidung des Projektinhabers (Code oder CR) ist in §1 festgehalten; danach: im Code-Zweig weist jeder der drei Fallbacks in Befehlsausgabe und
      `doctor` aus, soweit der gewählte Weg ihn ermittelt, mit Test je Fallback ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer));
      im CR-Zweig steht die präzisierte Fassung im Lastenheft.
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` und CHANGELOG nennen den geltenden Stand der Ausweisung.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft, im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Entfällt: Der Slice ist ohne Lieferung geschlossen (§7 `Gegenstand:`); der Plan des CR steht in
[`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md) §3. Die Tabelle bleibt als Stand der verworfenen Lösung A:

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/templates/devcontainer/*.tmpl` | update | Startscript meldet den Fallback an den Host (Weg offen) |
| `internal/hexagon/application/doctor_sandbox.go` | update | `doctor` liest die Meldung und erzeugt Items mit `DegradationFallback` |
| `internal/hexagon/application/` (Init-/Generate-/Up-Services) | update | Einträge mit `DegradationFallback` in der Befehlsausgabe |
| Tests der genannten Dateien | neu / update | Happy/Boundary/Negative je Fallback — nach [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |

## 4. Trigger

**Start:** entfällt; der Slice ist ohne Lieferung geschlossen (§7 `Gegenstand:`). Der Projektinhaber hat am 2026-10-17 den CR gewählt; ein späterer Code-Weg (Lösung A) beginnt als neuer Slice.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Meldung vom Container an den Host verlangt einen eigenen Mechanismus (Volume, Statusdatei) — dann ein Slice je Fallback.
- `in-progress` → `open` (blockiert): Die Entscheidung des Projektinhabers fehlt oder wechselt von Code zu CR.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Lösung A braucht einen Kanal Container → Host:** Eine Statusdatei oder ein Volume ist eine neue Schnittstelle und berührt das Sandbox-Profil
  ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): kein Host-Socket). — **Ausgang:** entfallen: Lösung A ist verworfen (Entscheidung des Projektinhabers 2026-10-17); ein Kanal Container → Host entsteht nicht.
- **Lösung B schwächt den Vertrag:** Die Präzisierung auf das host-seitig Ermittelbare nimmt eine Zusage zurück. — **Ausgang:** entfallen: der Projektinhaber entscheidet sie ausdrücklich als CR
  (`slice-cr-sandbox-fallback-ausweisung`); die Risiken der Änderung führt der CR in seinem §6.

## 7. Closure-Notiz

- **Was hat funktioniert:** — (nichts geliefert; der Gegenstand liegt im CR.) Die Ausgangslage in §1 hielt der Prüfung am Code stand: Der Planungslauf des CR las Startscripte und `doctor` erneut an der Stelle und fand dieselbe Lücke.
- **Was ging anders als geplant:** Der Slice war für zwei Wege angelegt (Code oder CR) und lag bis zur Entscheidung in `open/`; am 2026-10-17 fiel sie auf den CR.
- **Gegenstand:** übernommen von [`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md): Lösung B (Lastenheft präzisieren); die Liefer-Punkte der DoD liegen dort.
- **Steering-Loop-Eintrag:** Ein Slice mit zwei Lösungswegen, dessen Entscheidung beim Projektinhaber liegt, steht bis dahin in `open/` und geht bei der Wahl des CR-Wegs mit `Gegenstand:` an den CR über, statt zwei Pläne parallel zu halten. Gezählt, nicht verkörpert; erstes Auftreten.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt.
- **Folge-Slices:** [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md) — ist eine Datei in `open/`; der zwischenzeitliche CR [`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md) entfiel am 2026-10-17 (Zielbild, Code zum Ziel).
- **Risiken aus §6:** jedes trägt einen Ausgang — siehe §6 (beide entfallen).
- **Drei Paarungen:** Anker: keine verkörperte Regel, daher kein Anker; Folge-Slice: der CR liegt in `open/`; Register: keine Beobachtung angelegt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `internal/hexagon/application`
(Kürzel APP) — Modus Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), erfüllt die Schwelle
≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Im CR-Zweig kommt `spec/` (DOK, GF) hinzu.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
