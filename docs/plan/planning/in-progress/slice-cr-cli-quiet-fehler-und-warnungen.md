# Slice CR: `--quiet` zeigt nur Fehler, mit Ausnahme der Sicherheitslockerungen (`LH-FA-CLI-005`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md)).
> Anlass: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) sagt `--quiet` – nur Fehler; der Code
> hebt den Pegel nur auf Warn. Der Projektinhaber hat am 2026-10-16 entschieden:
> die Stufenbeschreibung des Lastenhefts bleibt, der Code folgt ihr
> ([`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md)); das Lastenheft
> bekommt eine Ausnahme für die Pflicht-Ausweisung von Sicherheitslockerungen.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — erster Slice der Welle; ihm folgen
[`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md) und
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
Verwandt: [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md) (Code: Fehler-Pegel und Ausnahme)
und [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Code: Wahl zwischen
mehreren Optionen; dieser CR ändert die Mehrfachangabe-Regel nicht).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16
(Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der
CR nicht; sie beschreiben das Verhalten erst nach allen drei Slices der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.7.0 → 0.8.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) behält
`--quiet` – nur Fehler und nimmt die Pflicht-Ausweisung von Sicherheitslockerungen
([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): „jede Sicherheitslockerung in der Befehlsausgabe ausweisen";
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): `--privileged` bei `docker`) als einzige Ausnahme auf.

**Warum eine Ausnahme:** Die Ausweispflicht von Sicherheitslockerungen ist eine
eigene Zusage. Ein `--quiet`, das sie verschluckt, würde sie aushebeln; die
Ausnahme macht den Vorrang im Vertrag ausdrücklich, statt ihn dem Code zu
überlassen.

**Wirkung von `--quiet` (Festlegung des Projektinhabers):** `--quiet` unterdrückt
Status- und Erfolgsmeldungen (Statustabelle von `up`, Erfolgsmeldung von `down`,
OK-Items von `doctor`) sowie Warn- und Info-Logs; `doctor --quiet` zeigt nur
Fehler-Items, der Exit-Code trägt das Ergebnis. Die Ausnahme sind die
Sandbox-Lockerungs-Hinweise.

### Gefundener Stand im Code

- Die Lockerungs-Hinweise entstehen in
  `internal/hexagon/application/devcontainer_sandbox.go`
  (`sandboxRelaxationWarnings`, Codes [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)/[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) als
  `driving.WarningEntry` in der Use-Case-Antwort (`generate.go` hängt sie an
  `warnings`). Sie laufen **nicht** über den Logger (`slog`).
- Der CLI-Adapter schreibt sie **direkt** auf die Ausgabe:
  `fmt.Fprintf(out, "Warning (%s): %s\n", …)` in `cli/generate.go` und
  `cli/init.go`, ohne `--quiet`-Prüfung. Im JSON-Modus werden sie über
  `mapWarningsToDiagnostics` zu `diagnostics[]` mit `level: "warn"`.
- `up` gibt die Sandbox-Hinweise **nicht** aus: der Up-Service
  (`internal/hexagon/application/upservice.go`) ruft `sandboxRelaxationWarnings`
  nicht auf; seine `Warnings` sind die Recreate-Warnungen
  ([`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen)). `cli/up.go` schreibt diese Warnungen nach einem Rücksprung bei
  `flags.Quiet`, also nur ohne `--quiet`. Der Rücksprung gilt für alles hinter der
  Statustabelle; eine Sandbox-Lockerung hätte dort keine Sichtbarkeit unter
  `--quiet`.
- Der Code ruft in Produktionscode nirgends `Logger.Warn` auf
  (`internal/adapter/driven/logger/logger.go` definiert die Methode); der Pegel
  `Warn` bei `--quiet` (`cli/root.go`, `PersistentPreRunE`) blendet heute nur
  `Info`-Logs aus.
- `--quiet` wirkt im CLI-Adapter über den Logger-Pegel hinaus: `up` unterdrückt
  Statustabelle, Diagnosen und Warnungen (`cli/up.go`); `down` die Erfolgsmeldung
  und die Liste entfernter Volumes (`cli/down.go`, `renderDownSuccess` in
  `cli/statusview.go`); `doctor` blendet OK-Items aus (`cli/doctor.go`,
  `writeDoctorReport`). `--quiet --json` ist ein No-op.

### CR-Kern (vorgeschlagene Fassung von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging))

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
> Eine Ausnahme gilt für Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices): Ihre Ausweisung in der Befehlsausgabe ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) bleibt unter `--quiet` sichtbar. Alle anderen Hinweise entfallen unter `--quiet`.
>
> Werden mehrere Verbosity-Optionen gleichzeitig angegeben (z. B. `--quiet --verbose`), gewinnt die zuletzt auf der Kommandozeile angegebene Option. Eine Validierungsabweisung wegen Mehrfachangabe erfolgt nicht.

**Unverändert, mit Begründung:**

- Die vier Stufenbeschreibungen: `--quiet` – nur Fehler ist der heutige Wortlaut;
  der Code folgt ihm im Umsetzungs-Slice.
- Der Satz zu mehreren Optionen: bleibt; der Code folgt ihm erst mit
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung im Code** — der Adapter ändert sich in
  [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md); dieser CR ändert nur das Vertrags-Stratum.
- **Die Regel „letzte Option gewinnt"** — Code-Anpassung in
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md), nicht Teil des CR.
- **Eine Spec-Beschreibung in `spezifikation.md`** — erst nach allen drei Slices
  der Welle, sonst beschreibt sie einen Zwischenstand.

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change)** — sie liegen
in [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md):

1. Ein Test belegt je Subkommando, dass `--quiet` keine Statusmeldungen zeigt
   und die Lockerungs-Warnung des Sandbox-Profils ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) bei `init`,
   `generate` und `up` sichtbar bleibt.
2. `docs/user/benutzerhandbuch.md` nennt `--quiet` als „nur Fehler" mit der
   Ausnahme (Optionstabellen, Zeilen um 299 und 1280).
3. `CHANGELOG.md` trägt den Eintrag.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.8.0 | (Annahmedatum) | LH-FA-CLI-005: `--quiet` bleibt „nur Fehler"; die
Ausweisung von Sicherheitslockerungen im Sandbox-Profil (LH-FA-DEV-006,
LH-FA-DEV-007) bleibt unter `--quiet` sichtbar. | Vereinbarung mit dem
Projektinhaber |
```

## 2. Definition of Done

- [ ] **Lastenheft 0.8.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) in der Fassung oben, Versionsfeld,
      Historie-Zeile (Verweis „Vereinbarung mit dem Projektinhaber").
- [ ] **Umsetzungsbedarf** einem Slice zugeordnet: Abnahmekriterien 1–3 liegen in
      [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md).
- [ ] Doku-Update: CHANGELOG (öffentlicher Vertrag berührt).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Versionsfeld 0.8.0, §16-Zeile |
| `CHANGELOG.md` | update | öffentlicher Vertrag berührt |

Reihenfolge: Lastenheft ändern → Umsetzung im Slice
[`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md) (die Ausnahme steht erst nach 0.8.0 im Vertrag).

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber gibt den CR frei. Die
Entscheidung vom 2026-10-16 (Pegel nur Fehler, Ausnahme Sicherheitslockerungen)
liegt vor.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Fassung der
  Ausnahme verlangt, weitere Warn-Klassen (etwa
  [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)) mitzuregeln — dann je Klasse ein Umsetzungs-Slice.
- `in-progress` → `open` (blockiert): der Projektinhaber ändert die
  Entscheidung zur Ausnahme.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.8.0 ist angenommen;
der Umsetzungs-Slice [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md) liegt in `open/` oder weiter.

## 6. Risiken und offene Punkte

- **Reichweite der Ausnahme:** „Sicherheitslockerung" deckt die Sandbox-Hinweise
  ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)); andere Warnungen (Recreate-Warnung
  [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen), Drittanbieter-Features) fallen unter `--quiet` weg. — **Ausgang:**
  entfallen: der Lastenheft-Satz nennt die Arten der Lockerung und schließt alle anderen
  Hinweise aus; der Umsetzungs-Slice [`slice-v1-cli-quiet-nur-fehler`](../open/slice-v1-cli-quiet-nur-fehler.md) bildet das
  im Code ab (die Kennung eines Hinweises genügt dafür nicht, weil [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) auch „has no effect" trägt).
- **Zwischenstand der Spec:** Bis alle drei Slices der Welle geschlossen sind,
  darf die Spezifikation das Verhalten nicht beschreiben. — **Ausgang:**
  entfallen: die Welle [`welle-cli-verbosity`](../welle-cli-verbosity.md) schließt erst mit dem Spezifikations-Eintrag im
  letzten Slice.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/`, `harness/`,
`docs/plan/` (eine Sub-Area laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Abschnitt Modus-Deklaration,
Modus Greenfield). `docs/user/` und `CHANGELOG.md` sind öffentliche Verträge,
keine Sub-Area; der Code (`internal/adapter/driving/cli`) wird von diesem CR
nicht geändert.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
