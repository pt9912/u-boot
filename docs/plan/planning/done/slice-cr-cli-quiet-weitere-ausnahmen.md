# Slice CR: `--quiet` zeigt nur Fehler, mit vier benannten Ausnahmen (`LH-FA-CLI-005`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md)).
> Anlass: Der Review des Code-Slices
> [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
> ([Report](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler.md), F-5 und F-6) fand zwei
> Spannungen zwischen [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und anderen Zusagen des Lastenhefts. Der Projektinhaber hat am
> 2026-10-16 entschieden: Die Degradations-Fallbacks und der Compose-Fortschritt bei `up`
> bleiben unter `--quiet` sichtbar; beides wird als Ausnahme in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) benannt.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — zweiter Slice der Welle; ihm geht
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) voraus, ihm folgen
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) und
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation, Ausweisung jedes Fallbacks),
[`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker) (Fortschrittsanzeige bei `up`).
Verwandt: [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (Code: setzt die vier Ausnahmen um).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16
(Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der
CR nicht; sie beschreiben das Verhalten erst nach allen Slices der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.8.0 → 0.9.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) behält
`--quiet` – nur Fehler und benennt **vier** Ausnahmen statt einer: die Ausweisung von Sicherheitslockerungen
(bereits in 0.8.0), die Ausweisung der Degradations-Fallbacks
([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), den Compose-Fortschritt von `u-boot up`
([`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)) und die Zusammenfassung der betroffenen Pfade
von `--force` ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung),
[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).

**Warum Ausnahmen statt Auslegung:** Die drei neuen Ausnahmen folgen aus eigenen Zusagen des Lastenhefts —
„jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen“, „der Fortschritt einzelner Services
(Pull, Create, Start, Healthcheck) muss sichtbar sein“ und „vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden“. Ohne Benennung widersprächen sich zwei Sätze, und der Code
müsste den Vorrang still entscheiden; die Ausnahme macht ihn im Vertrag ausdrücklich.

**Wirkung von `--quiet` (Festlegung des Projektinhabers, 2026-10-16):** `--quiet` unterdrückt die Erfolgs- und
Zusammenfassungsmeldung **aller** Subkommandos, die eine ausgeben (`init`, `generate`, `add`, `remove`, `up`, `down`,
`doctor`: nur Fehler); `remove` bekommt dafür die Option `--quiet`. Das ist der heutige Wortlaut „nur Fehler“ und braucht
**keinen** Lastenheft-Wechsel; die Umsetzung liegt im Code-Slice
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md).

### Gefundener Stand im Code

- Die Lockerungs-Hinweise (Ausnahme 1) sind umgesetzt: `driving.WarningEntry` trägt ein Kennzeichen, `cli/warnings.go`
  filtert danach (siehe Review-Report F-5).
- Die Degradations-Fallbacks (`vfs`-Storage, Container ohne Daemon, Wegfall der Egress-Restriktion) entstehen nicht beim
  Rendern, sondern beim Containerstart (Skripte unter `internal/hexagon/application/templates/devcontainer/`) und als
  Items von `u-boot doctor` (`internal/hexagon/application/doctor_sandbox.go`). `sandboxRelaxationWarnings` erzeugt sie
  nicht; `doctor --quiet` blendet Items unterhalb des Fehler-Niveaus aus, also auch Fallback-Items der Stufe Warn.
- Der Compose-Fortschritt von `up` läuft live auf stderr und wird von `--quiet` nicht berührt; ein Kommentar in
  `internal/adapter/driving/cli/up.go` nennt [`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker) als Grund. Das Lastenheft nennt die Ausnahme bisher nicht
  (Review-Report F-6).

### CR-Kern (angenommene Fassung von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging))

> ### [LH-FA-CLI-005](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) – Verbosity und Logging
>
> *Priorität: MVP.* Das Produkt muss ein konfigurierbares Ausgabeverbose-Level (Verbosity) unterstützen.
>
> Mindestens müssen folgende Stufen unterstützt werden:
>
> - `--quiet` – nur Fehler
> - Standard – Statusmeldungen und Fehler
> - `--verbose` – zusätzlich Detailinformationen
> - `--debug` – zusätzlich interne Diagnoseausgaben
>
> Unter `--quiet` geben die Subkommandos, die etwas ändern (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`), nur Fehler aus; ihre Erfolgs- und Zusammenfassungsmeldungen entfallen. `u-boot doctor` zeigt nur Fehler-Einträge. Vier Ausnahmen bleiben unter `--quiet` sichtbar:
>
> 1. die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) in der Befehlsausgabe ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion));
> 2. die Ausweisung der Degradations-Fallbacks (Container ohne Docker-Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in der Befehlsausgabe und in `u-boot doctor` ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer));
> 3. der Fortschritt einzelner Services bei `u-boot up` (Pull, Create, Start, Healthcheck; [`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker));
> 4. die Zusammenfassung der betroffenen Pfade, die `--force` vor dem Überschreiben ausgibt ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).
>
> Alle anderen Hinweise und Warnungen entfallen unter `--quiet`. Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt.
>
> Werden mehrere Verbosity-Optionen gleichzeitig angegeben (z. B. `--quiet --verbose`), gewinnt die zuletzt auf der Kommandozeile angegebene Option. Eine Validierungsabweisung wegen Mehrfachangabe erfolgt nicht.

Der Satz zu den Abfragebefehlen nennt das Kriterium (der Befehl ändert nichts) und die Befehle (`config get`, `config list`,
`logs`, `template list`); `config set` zählt zu den ändernden Subkommandos. Die vierte Ausnahme folgt aus dem Review
(F-2): [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) verlangen die Zusammenfassung vor dem
Überschreiben; sie ist eine Ausweisung vor einer zerstörenden Aktion wie die Sicherheitslockerungen.

**Unverändert, mit Begründung:**

- Die vier Stufenbeschreibungen: `--quiet` – nur Fehler bleibt der Wortlaut.
- Der Satz zu mehreren Optionen: bleibt; der Code folgt ihm mit
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung im Code** — Adapter und Use-Cases ändern sich in
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md); dieser CR ändert nur das Vertrags-Stratum.
- **Die Reichweite „alle Subkommandos“** — steht im Wortlaut „nur Fehler“, braucht keinen Lastenheft-Wechsel.
- **Die Regel „letzte Option gewinnt“** — Code in
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).
- **Eine Spec-Beschreibung in `spezifikation.md`** — erst nach allen Slices der Welle, sonst beschreibt sie einen Zwischenstand.

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change)** — sie liegen in
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md):

1. Je Subkommando belegt ein Test, dass `--quiet` Erfolgs- und Zusammenfassungsmeldungen unterdrückt (`remove` mit neuer Option).
2. Ein Test pinnt je Ausnahme, dass sie sichtbar bleibt: Lockerungs-Hinweis, Degradations-Fallback (Befehlsausgabe und `doctor`), Compose-Fortschritt bei `up`, `--force`-Zusammenfassung.
3. `docs/user/benutzerhandbuch.md` und `CHANGELOG.md` nennen `--quiet` mit den vier Ausnahmen.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.9.0 | 2026-10-16 | LH-FA-CLI-005: Vier Ausnahmen unter `--quiet`
(Sicherheitslockerungen, Ausweisung der Degradations-Fallbacks nach
LH-FA-DEV-007, Fortschritt von `u-boot up` nach LH-NFA-PERF-002,
`--force`-Zusammenfassung nach LH-FA-CLI-005A und LH-FA-INIT-005); die
ändernden Subkommandos geben unter `--quiet` nur Fehler aus, `u-boot
doctor` nur Fehler-Einträge, die Abfragebefehle (`config get`, `config
list`, `logs`, `template list`) behalten ihre Ausgabe; die Stufen und die
Mehrfachangabe-Regel bleiben unverändert. | Vereinbarung mit dem
Projektinhaber |
```

## 2. Definition of Done

- [x] **Lastenheft 0.9.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) in der Fassung oben, Versionsfeld,
      Historie-Zeile (Verweis „Vereinbarung mit dem Projektinhaber").
- [x] **Umsetzungsbedarf** einem Slice zugeordnet: Abnahmekriterien 1–3 liegen in
      [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md).
- [x] Doku-Update: CHANGELOG (öffentlicher Vertrag berührt).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Versionsfeld 0.9.0, §16-Zeile |
| `CHANGELOG.md` | update | öffentlicher Vertrag berührt |

Reihenfolge: Lastenheft ändern → Umsetzung im Slice
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (die Ausnahmen stehen erst nach 0.9.0 im Vertrag).

## 4. Trigger

**Start** (`next` → `in-progress`): sofort; die Entscheidung des Projektinhabers vom 2026-10-16 (Reichweite,
Compose-Fortschritt, Degradations-Fallbacks) liegt vor.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Fassung verlangt, weitere Hinweis-Klassen mitzuregeln
  (etwa [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)) — dann je Klasse ein Umsetzungs-Slice.
- `in-progress` → `open` (blockiert): der Projektinhaber ändert die Entscheidung zu einer der vier Ausnahmen.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.9.0 ist angenommen;
der Umsetzungs-Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) liegt in `open/` oder weiter.

## 6. Risiken und offene Punkte

- **Abgrenzung „Degradations-Fallback“:** Der Begriff muss sich im Code an `doctor`-Items und an Skript-Ausgaben beim
  Containerstart festmachen lassen; ein Fallback-Item der Stufe Warn könnte sonst mit gewöhnlichen Warn-Items verwechselt
  werden. — **Ausgang:** entfallen: der Lastenheft-Satz zählt die drei Fallbacks aus
  [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) einzeln auf; der Code-Slice pinnt jeden einzeln.
- **Containerstart-Ausgabe liegt außerhalb von u-boots Prozess:** Die Fallback-Meldungen beim Containerstart stammen aus
  generierten Skripten, nicht aus dem CLI-Lauf von `--quiet`. — **Ausgang:** entfallen: das Lastenheft verlangt die
  Ausweisung „in der Befehlsausgabe und in `u-boot doctor`“; was u-boot selbst ausgibt, regelt der Code-Slice, die Skripte kennen
  `--quiet` nicht.
- **Zwischenstand der Spec:** Bis alle Slices der Welle geschlossen sind, darf die Spezifikation das Verhalten nicht
  beschreiben. — **Ausgang:** entfallen: die Welle [`welle-cli-verbosity`](../welle-cli-verbosity.md) schließt erst mit dem
  Spezifikations-Eintrag im letzten Slice.

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen`](../../../reviews/2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md)
(0 HIGH, 3 MEDIUM, 1 LOW, 3 INFO; alle behoben in `591acb7`); Verification Evidence
[`2026-10-16-verification-slice-cr-cli-quiet-weitere-ausnahmen`](../../../reviews/2026-10-16-verification-slice-cr-cli-quiet-weitere-ausnahmen.md)
(bestanden, fünf Doku-Befunde zur Zählung in den Planungsdateien nachgezogen); `make gates` grün, `make doc-trace`: 83 Anforderungen, 0 Waisen.

- **Was hat funktioniert:** Der Review las das Lastenheft gegen sich selbst:
  [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) verlangen die
  Zusammenfassung vor dem Überschreiben, die der erste Entwurf unter `--quiet` entfallen ließ. Daraus wurde die vierte Ausnahme.
- **Was ging anders als geplant:** Der Satz zu den Ergebnis-Ausgaben war als Vorschlag in den CR gekommen und im Lastenheft gelandet, ohne
  ein Kriterium zu nennen („Zweck eines Befehls" trägt auch `config set`). Er nennt jetzt das Kriterium (der Befehl ändert nichts) und die Befehle.
  Die Zählung („drei", dann „vier" Ausnahmen) war an mehreren Planungsdateien nachzuziehen.
- **Steering-Loop-Eintrag:** Eine Ausnahme-Regel im Lastenheft wird vor der Änderung gegen alle älteren Zusagen geprüft, die dasselbe Verhalten berühren
  (Suche nach `--quiet`, „Zusammenfassung", „ausgewiesen"); der Review fand die Kollision erst nach dem Commit. Gezählt, nicht verkörpert
  (zweites Auftreten nach dem ersten CR, aber ein anderes Fehlermuster).
- **Grenze:** Das CHANGELOG sagt zukünftig, was `--quiet` leistet; bis die Slices `slice-v1-cli-quiet-nur-fehler` und
  `slice-v1-cli-verbosity-letzte-option-gewinnt` geliefert sind, kein Release-Cut. Die Purge- und Abhängigkeits-Warnungen von `remove` sind keine Ausnahme; der
  Code-Slice pinnt das.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt.
- **Folge-Slices:** [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (trägt die Abnahmekriterien 1–3),
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).
- **Risiken aus §6:** alle drei tragen ihren Ausgang (siehe §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/`, `harness/`,
`docs/plan/` (eine Sub-Area laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Abschnitt Modus-Deklaration,
Modus Greenfield). `CHANGELOG.md` ist ein öffentlicher Vertrag, keine Sub-Area; der Code
(`internal/adapter/driving/cli`) wird von diesem CR nicht geändert.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
