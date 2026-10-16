# Slice CR: `--quiet` zeigt Fehler und Warnungen (`LH-FA-CLI-005`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md)).
> Anlass: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) sagt `--quiet` – nur Fehler; der Code
> hebt den Pegel nur auf Warn. Der Projektinhaber hat entschieden: der Code
> behält recht, das Lastenheft wird angepasst.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — Spec-Wartung am Vertrags-Stratum; es gibt keine
Closure-Bedingung außer der DoD dieses Slice (Wellen sind Release-Versionen,
[`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
Verwandt: [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Code:
„letzte Option gewinnt"; unabhängig von diesem CR umsetzbar — dieser CR klärt
die Zusage von `--quiet`, jener die Wahl zwischen mehreren Optionen).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16
(Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der
CR nicht; sie beschreiben das Verhalten erst nach
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) und diesem CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.7.0 → 0.8.0**: `--quiet` heißt „nur Fehler
und Warnungen". Der Code (`internal/adapter/driving/cli/root.go`) setzt den
Logger bei `--quiet` auf `Warn`; das Lastenheft sagt heute „nur Fehler".

**Warum das Lastenheft und nicht der Code:** Warnungen können
sicherheitsrelevant sein, etwa die Ausweispflicht von Sicherheitslockerungen im
Sandbox-Profil ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): `--privileged` wird in
der Befehlsausgabe ausgewiesen). Ein `--quiet`, das sie verschluckt, würde die
Ausweispflicht aushebeln.

### Offener Punkt — vor der Lastenheft-Änderung zu klären

Das Sicherheitsargument trägt nur, wenn die Lockerungs-Hinweise bei `--quiet`
tatsächlich sichtbar bleiben. **Prüfauftrag:** laufen sie über den Logger (dann
trägt `--quiet` auf Warn sie weiter) oder direkt auf stdout/stderr (dann ist
das Argument schwächer und die Entscheidung neu zu bewerten)?

**Gefundener Stand im Code (nicht entschieden):**

- Die Lockerungs-Hinweise entstehen in
  `internal/hexagon/application/devcontainer_sandbox.go`
  (`sandboxRelaxationWarnings`, Codes [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)/[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) als
  `driving.WarningEntry` in der Use-Case-Antwort (`generate.go` hängt sie an
  `warnings`). Sie laufen **nicht** über den Logger (`slog`).
- Der CLI-Adapter schreibt sie **direkt** auf die Ausgabe:
  `fmt.Fprintf(out, "Warning (%s): %s\n", …)` in `cli/generate.go` und
  `cli/init.go`; in `cli/up.go` ebenfalls, dort aber nach dem Rücksprung bei
  `flags.Quiet` (`up --quiet` unterdrückt Tabelle, Diagnosen und Warnungen).
  Eine `--quiet`-Prüfung in `generate.go`/`init.go` auf diesem Pfad wurde nicht
  gefunden. Im JSON-Modus werden sie über `mapWarningsToDiagnostics` zu
  `diagnostics[]` mit `level: "warn"`.
- Folge: Das Sicherheitsargument stützt sich bei `init`/`generate` auf
  Direktausgabe, nicht auf den Logger-Pegel; die Aussage „`--quiet` auf Warn
  trägt sie weiter" gilt dort nur, solange die Direktausgabe nicht an `quiet`
  gekoppelt wird. Bei `up --quiet` werden Warnungen heute unterdrückt.
  Zu bewerten, ob der CR dadurch an Tragkraft verliert oder die
  Ausnahme `up` mitzuregeln ist.

### CR-Kern (vorgeschlagene Fassung von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging))

> ### [LH-FA-CLI-005](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) – Verbosity und Logging
>
> *Priorität: MVP.* Das Produkt muss ein konfigurierbares Ausgabeverbose-Level (Verbosity) unterstützen.
>
> Mindestens müssen folgende Stufen unterstützt werden:
>
> - `--quiet` – nur Fehler und Warnungen
> - Standard – Statusmeldungen und Fehler
> - `--verbose` – zusätzlich Detailinformationen
> - `--debug` – zusätzlich interne Diagnoseausgaben
>
> Werden mehrere Verbosity-Optionen gleichzeitig angegeben (z. B. `--quiet --verbose`), gewinnt die zuletzt auf der Kommandozeile angegebene Option. Eine Validierungsabweisung wegen Mehrfachangabe erfolgt nicht.

**Unverändert, mit Begründung:**

- Standard, `--verbose`, `--debug`: der Code entspricht dem Lastenheft; kein
  Änderungsgrund.
- Der Satz zu mehreren Optionen: bleibt; der Code folgt ihm erst mit
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung im Code** — der Code ist die Bezugsgröße dieses CR; nötige
  Anpassungen (z. B. bei `up --quiet`, falls der Prüfauftrag das verlangt) sind
  Folge-Slices.
- **Die Regel „letzte Option gewinnt"** — Code-Anpassung in
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md), nicht Teil des CR.
- **Eine Spec-Beschreibung in `spezifikation.md`** — erst nach beiden Vorgängen,
  sonst beschreibt sie einen Zwischenstand.

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change):**

1. Ein Test belegt, dass `--quiet` Warnungen zeigt (mindestens die
   Lockerungs-Warnung des Sandbox-Profils aus [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) bei `init`/`generate`
   mit `--quiet`) und keine Statusmeldungen.
2. `docs/user/benutzerhandbuch.md` nennt `--quiet` als „nur Fehler und
   Warnungen" (Optionstabellen, Zeilen um 299 und 1280).
3. `CHANGELOG.md` trägt den Eintrag.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.8.0 | (Annahmedatum) | LH-FA-CLI-005: `--quiet` zeigt Fehler und
Warnungen (statt nur Fehler), weil Warnungen sicherheitsrelevant sein können
(Ausweispflicht von Lockerungen im Sandbox-Profil). | Vereinbarung mit dem
Projektinhaber |
```

## 2. Definition of Done

- [ ] **Lastenheft 0.8.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) in der Fassung oben (nach Klärung des
      offenen Punkts), Versionsfeld, Historie-Zeile.
- [ ] **Umsetzungsbedarf** geklärt: Abnahmekriterien 1–3 sind einem
      Umsetzungs-Slice mit Kennung zugeordnet oder als bereits erfüllt belegt.
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
| `docs/user/benutzerhandbuch.md` | update (Umsetzung) | Optionstabellen nennen „Warnungen und Fehler bleiben" bereits; Abgleich mit der neuen Fassung |

Reihenfolge: Offenen Punkt klären (Prüfauftrag §1) → Lastenheft ändern →
Umsetzungsbedarf (Abnahmekriterium 1) bewerten.

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber gibt den CR frei; der
Prüfauftrag zu den Warn-Pfaden (§1) ist abgeschlossen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Prüfauftrag zeigt
  weitere Warn-Pfade mit abweichendem `--quiet`-Verhalten — dann je Pfad ein
  Umsetzungs-Slice.
- `in-progress` → `open` (blockiert): die Prüfung zeigt, dass die Warnungen
  nicht über den Logger laufen und das Sicherheitsargument entfällt — die
  Entscheidung ist dann mit dem Projektinhaber neu zu bewerten.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.8.0 ist angenommen;
etwaige Umsetzungs-Slices liegen in `open/` oder weiter.

## 6. Risiken und offene Punkte

- **Warn-Pfad ohne Logger:** Die Lockerungs-Hinweise laufen heute über
  `WarningEntry` und Direktausgabe, nicht über den Logger; `up --quiet`
  unterdrückt sie. Das Sicherheitsargument des CR ist damit nicht ohne
  Weiteres durch „Logger auf Warn" belegt. — **Ausgang:** weiter offen:
  Prüfauftrag in §1, vor der Lastenheft-Änderung zu klären.
- **Zwischenstand der Spec:** Bis Slice
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) und CR beide geschlossen sind,
  darf die Spezifikation das Verhalten nicht beschreiben. — **Ausgang:** weiter
  offen: bis zur Closure beider.

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
