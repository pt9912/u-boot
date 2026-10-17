# Review-Report: slice-cr-devcontainer-agent-features-sandbox — 2026-10-17

**Review-Art:** Plan/Design — Diff gegen den CR-Kern des Slice, die Lastenheft-Zusagen [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), die angrenzenden älteren Zusagen und Hard Rules (Maintainability; nicht gegen die DoD).

**Gegenstand:** Commit `5de388c` für den Slice [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/done/slice-cr-devcontainer-agent-features-sandbox.md) (Welle `welle-devcontainer-sshd-claude-code`): `spec/lastenheft.md` (Version 0.11.0, §3 `LH-FA-DEV-006`, `LH-FA-DEV-009`, `LH-FA-CLI-005` Ausnahme 1, §16-Zeile), `CHANGELOG.md` `[Unreleased]`. <!-- d-check:ignore (Kennung und Pfade als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` (aktueller Stand im Arbeitsbaum), HEAD `5de388c` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/done/slice-cr-devcontainer-agent-features-sandbox.md) (§1 CR-Kern und Prüfungstabelle, §2 DoD, §6 Risiken) und Welle-Plan `welle-devcontainer-sshd-claude-code` <!-- d-check:ignore (Kennung als Zitat) -->
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-SEC-001`](../../spec/lastenheft.md#lh-nfa-sec-001--keine-secrets-einchecken), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)
- [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md)
- Vorherige Reports derselben Art: `2026-10-16-slice-cr-cli-quiet-vorschau.md`, `2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md` <!-- d-check:ignore (Report-Dateinamen als Zitat) -->
- `AGENTS.md` (Hard Rules §3.4, §3.7, §3.9)

**Systematische Suche gegen ältere Zusagen** (Lastenheft; Begriffe `Zugangsdaten`, `Token`, `Secret`, `Credential`, `Passwort`, `Geheimnis`, `Login`, `ssh`, `sshd`, `Agent`, `Sicherheitslockerung`, `ausgewiesen`, `--quiet`, `Umgebungsvariable`, `Secret-Mount`): Treffer im Normtext nur in `LH-FA-CLI-005` (Ausnahme 1), `LH-FA-DEV-003`/`-004`-Querverweis auf `LH-NFA-SEC-004`, `LH-FA-DEV-006` (Clone-URL ohne Zugangsdaten, keine Host-Dateien mit Geheimnissen, Ausweisung, Nicht-Ziel Agent), `LH-FA-DEV-009`, `LH-NFA-SEC-001`, `LH-NFA-SEC-002`/Admin-Beispielwerte (`.env.example`), Nicht-Ziel „Secret-Management für Produktion“. `sshd`, `Claude` und `Anthropic` kommen im Lastenheft-Normtext nicht vor. Kein Widerspruch zu einer älteren Zusage gefunden (Einzelergebnisse unten in „Prüfpunkte“). <!-- d-check:ignore (Suchbegriffe und Kennungen als Zitat) -->

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Eintrag sagt im Indikativ, Agenten-Zugangsdaten „reach the container at run time through an environment variable, and `u-boot doctor` checks the source“. Im Code reicht die Sandbox-Vorlage nur `GIT_TOKEN` als `remoteEnv` durch und `doctor` prüft nur die Git-Token-Quelle; ein Agenten-Schlüssel wird nirgends übergeben oder geprüft. Der Schlusssatz „The `sshd` and Claude Code features themselves follow in a later release“ nimmt nur die Features aus, nicht die Übergabe und die `doctor`-Prüfung. | Docs/Release-Linse (`harness/review.md`: Doku driftet gegen implementiertes Verhalten); [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `CHANGELOG.md` · "they reach the\n  container at run time through an environment variable, and `u-boot doctor` checks the source" (Zeilenumbruch nach „the“; eindeutig: „checks the source“) | nein — Review-only; nach den Code-Slices der Welle per `make test` | Changelog-Zusage vor Code-Stand |
| F-2 | INFO | „bei aktivem Agenten-Feature die Umgebungsvariable des Agenten“: Das Lastenheft nennt weder, welche Features „Agenten-Features“ sind, noch den Variablennamen; es gibt nur das Beispiel „Zugangsschlüssel eines Coding-Agenten“. Wie `doctor` die Variable kennt (Eigenschaft des Katalog-Eintrags oder der Feature-Quelle), entscheidet der Vertrag nicht. Das ist Sache der Spezifikation (Katalog-Metadaten), kein Vertragsloch im Lastenheft; der Slice trägt es als Risiko 2 und der Folge-Slice zur Spezifikation ist zuständig. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `spec/lastenheft.md` · "bei aktivem Agenten-Feature die Umgebungsvariable des Agenten" | nein — Review-only (Zuständigkeit Spezifikation / Verifier) | Begriff „Agenten-Feature“ im Vertrag ohne Abgrenzung |
| F-3 | INFO | „einen im Container aktivierten Login-Dienst (Passwort-Login)“: Die Klammer lässt offen, ob sie den Dienst erklärt (jeder Login-Dienst) oder begrenzt (nur Passwort-Login). Ein `sshd` mit ausschließlich Schlüssel-Login ist damit nicht eindeutig geregelt; der Slice nennt Schlüssel-Login ausdrücklich als Nicht-Ziel („Schlüssel-Login, Passwort-Setzen und Agent-Konfiguration durch `u-boot`“). Der Folge-Slice `sshd` legt den Fall in der Praxis fest (Passwort-Login nach `sudo passwd`). | [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `spec/lastenheft.md` · "einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen." | nein — Review-only | Klammerzusatz erklärt oder begrenzt unklar |
| F-4 | INFO | Die Überschrift „LH-FA-DEV-009 – Git-Zugangsdaten im Sandbox-Devcontainer“ trägt „Git“, der Text deckt jetzt auch Agenten-Zugangsdaten. Die Entscheidung ist im Slice benannt (Anker-Stabilität; die Anker werden in README, Spezifikation, Benutzerhandbuch und Code-Tests verlinkt), die Zusage selbst ist im ersten Absatz eindeutig. Kein Verstoß gegen eine Konvention; ein Leser, der nur das Inhaltsverzeichnis liest, erwartet keine Agenten-Zusage. | Maintainability (kein Konventions-Anker, daher nicht LOW) | `spec/lastenheft.md` · "### LH-FA-DEV-009 – Git-Zugangsdaten im Sandbox-Devcontainer" | ja — `make docs-check` (Anker lösen auf, 0 Befunde) | Überschrift schmaler als Zusage | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-5 | INFO | Der Slice-Text ist nach dem Lastenheft-Commit nicht nachgezogen: §1 nennt die Teilfrage zu Entscheidung 1 „offen“ und „bei Teilfrage ‚ja‘“, der Entwurf der §16-Zeile nennt `LH-FA-CLI-005` nicht (die gesetzte Zeile tut es), und der Welle-Plan sagt „heute gilt es nur für Git-Zugangsdaten“. DoD-Zeile 1 verlangt, dass die Entscheidung „in §1 festgehalten“ ist. Die Dokumente sind veränderlich; der Nachzug ist Sache des Implementers vor der Closure, die Prüfung Sache des Verifiers. | DoD des Slice (Verifier-Zuständigkeit) | `docs/plan/planning/done/slice-cr-devcontainer-agent-features-sandbox.md` · "**Offene Teilfrage zu Entscheidung 1 (der Projektinhaber entscheidet beim Review des CR):**" | nein — Review-only (Verifier) | Slice-Entwurf folgt der Entscheidung nicht nach | <!-- d-check:ignore (Kennungen als Zitat) -->

## Prüfpunkte aus dem Auftrag

1. **Wörtlichkeit.** Der Text von `LH-FA-DEV-009` im Lastenheft ist Wort für Wort der Zitat-Block „Vorschlag für den erweiterten Text“ in Slice §1 (Satz für Satz verglichen: Einleitung mit Klammerbeispiel, drei Aufzählungspunkte; die Datei-Links des Zitats stehen im Lastenheft als `#…`-Anker). Der Zusatz in `LH-FA-DEV-006` ist wörtlich der Vorschlag („- einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.“); der vorherige Punkt endet jetzt mit „;“ statt „.“ und der neue Punkt mit „.“ — Aufzählungsinterpunktion korrekt. Die Variante „ja“ der Teilfrage steht in `LH-FA-CLI-005` Ausnahme 1 und ist im Slice als Vorschlag beschrieben („Ausnahme 1 nennt den Login-Dienst“); der Wortlaut der Ausnahme ist im Slice nicht vorformuliert, also nicht vergleichbar. Überschrift: siehe F-4. <!-- d-check:ignore (Kennungen als Zitat) -->
2. **Eindeutigkeit und Widersprüche.**
   - „Zugangsdaten eines Agenten“: abgegrenzt durch das Beispiel „Zugangsschlüssel eines Coding-Agenten, den ein aktiviertes Devcontainer-Feature nutzt“; Failure-Szenario „Implementer erfasst zu wenig/zu viel“: Der Vertrag zählt als Orte Image, Workspace-Volume, `u-boot.yaml` und „eine andere erzeugte Datei“; das schließt `devcontainer.json`, `compose.yaml` und `.env.example` ein. Wer den Wert dort ablegt, bricht die Zusage; wer nur die Referenz ablegt, erfüllt sie (zweiter Satz des ersten Punkts). Keine Lücke, die eine Implementierung falsch öffnen lässt.
   - `doctor` und die Variable des Agenten: F-2 (INFO, Sache der Spezifikation).
   - „Übergabe zur Laufzeit per Umgebungsvariable“ gegen `LH-FA-DEV-006` „keine Host-Dateien mit Geheimnissen … einbinden“: verträglich, weil eine Umgebungsvariable keine Host-Datei ist; der Secret-Mount bleibt ausdrücklich auf Git-Zugangsdaten beschränkt (Agenten-Zugangsdaten nur per Variable), und `LH-FA-DEV-006` verbietet Host-Dateien mit Geheimnissen — der bisherige Secret-Mount wäre dieselbe Frage, sie ist unverändert und gehört nicht zum Diff. <!-- d-check:ignore (Kennungen als Zitat) -->
   - `LH-NFA-SEC-001` („keine echten Secrets in erzeugte Dateien“): verträglich; der Text verlangt „nur die Referenz auf die Umgebungsvariable des Hosts, nie den Wert“. <!-- d-check:ignore (Kennungen als Zitat) -->
   - „Login-Dienst“ in `LH-FA-DEV-006` und `LH-FA-CLI-005`: gleiche Formulierung an beiden Stellen; Eindeutigkeit des Klammerzusatzes siehe F-3. <!-- d-check:ignore (Kennungen als Zitat) -->
   - Satzstruktur Ausnahme 1: „die Ausweisung von Sicherheitslockerungen (…) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe“ — grammatisch korrekt (zweites Genitiv-Objekt zu „Ausweisung“); die Aufzählung der vier Ausnahmen ist nicht verschoben (Nummern 1–4 unverändert, Einleitung „Vier Ausnahmen“ passt). Ausnahme 1 bleibt die Sicherheitslockerungs-Ausnahme; dass der Login-Dienst in `LH-FA-DEV-006` als eigener Punkt neben „jede Sicherheitslockerung“ steht und in Ausnahme 1 mit ihr geführt wird, ist konsistent mit der Entscheidung des Projektinhabers (Variante „ja“). <!-- d-check:ignore (Kennungen als Zitat) -->
3. **Versionssprung und §16.** 0.10.0 → 0.11.0 folgt dem Muster der Zeilen 0.5.0/0.7.0–0.10.0 (Minor bei neuer oder erweiterter Zusage mit Code-Folge). Kopftabelle: Spaltenbreite erhalten. Die Zeile 0.11.0 trägt Datum 2026-10-17, drei Kennungen als Link, Quelle „Vereinbarung mit dem Projektinhaber“. Abweichung vom CR-Entwurf: die gesetzte Zeile nennt zusätzlich `LH-FA-CLI-005` (der Entwurf tat es nicht, weil die Teilfrage offen war) und sagt „Ausweisung eines im Container aktivierten Login-Dienstes (Passwort-Login)“ statt „eines Login-Dienstes“ — eine Folge der entschiedenen Teilfrage, kein Fehler; Nachzug des Entwurfs siehe F-5. <!-- d-check:ignore (Kennungen als Zitat) -->
4. **Decken-Regel (`AGENTS.md` §3.4).** Der Diff fügt im Lastenheft keinen Verweis auf Slice, Welle, Carveout, Roadmap oder ADR ein; `sshd`, Claude Code und Anthropic werden im Lastenheft-Text nicht genannt (der Text bleibt generisch: „Login-Dienst“, „Coding-Agent“, „Devcontainer-Feature“). Das ist vereinbar mit der Aussage „Beispiele sind keine Zusage“: Die Beispielliste von `LH-FA-DEV-003` bleibt unverändert, das Klammerbeispiel in `LH-FA-DEV-009` benennt keinen Produktnamen. Neu verlinkte Anker (`#lh-fa-dev-006--sandbox-profil`, `#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer`, `#lh-fa-cli-005--verbosity-und-logging`) existieren; `make docs-check` meldet 0 Befunde. <!-- d-check:ignore (Anker als Zitat) -->
5. **CHANGELOG.** Versionsangabe „Lastenheft 0.11.0“ und der Inhalt (Agenten-Zugangsdaten nie in Image/Volume/`u-boot.yaml`/erzeugter Datei, Übergabe per Umgebungsvariable, `doctor` prüft die Quelle, Ausweisung des Login-Dienstes bleibt unter `--quiet` sichtbar) stimmen mit dem Lastenheft überein. Der Schlusssatz nennt die Features als folgend, nicht aber die Übergabe und die `doctor`-Prüfung; Release-Cut-Szenario: F-1. Verweise `LH-FA-DEV-009`, `-006`, `CLI-005` lösen auf. <!-- d-check:ignore (Kennungen als Zitat) -->
6. **Sprachvertrag, `AGENTS.md` §3.7 und §3.9.** Lastenheft und Slice deutsch, CHANGELOG englisch: vertragskonform. Die neuen Sätze beschreiben den Zustand im Indikativ; keine Chronik, keine verworfene Alternative, kein „früher/former/no longer“ im Diff (§3.9 a). §3.9 b: kein geänderter oder umbenannter Symbolname im Diff; der Prüfbefehl nach alten Namen entfällt, weil der Diff keinen Code berührt. §3.9 c: kein Doc-Kommentar-Block betroffen. Der Slice trägt in §1 eine Chronik-Aussage „Ausgangslage“; das ist Planungs-Prosa, kein Kommentar. <!-- d-check:ignore (Kennungen als Zitat) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` Kopftabelle, Wortlaut `LH-FA-DEV-009` und `LH-FA-DEV-006` gegen CR-Kern | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-CLI-005` Ausnahme 1, Satzstruktur, Nummerierung der vier Ausnahmen | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat) -->
| §16-Historie 0.11.0 (Muster, Quelle, Verweise) | geprüft, ohne Befund (Abweichung vom Entwurf: F-5) |
| Ältere Zusagen `LH-FA-DEV-003`, `-007`, `-008`, `LH-NFA-SEC-001`, `-004`, `LH-FA-INIT-005` | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat) -->
| Decken-Regel und Anker im Lastenheft (`make docs-check`: 282 Dateien, 0 Befunde) | geprüft, ohne Befund |
| `AGENTS.md` §3.7 und §3.9 in Lastenheft, CHANGELOG und Slice | geprüft, ohne Befund |
| Sprachvertrag Lastenheft/CHANGELOG | geprüft, ohne Befund |
| README, `README.de.md`, Spezifikation: Verweise auf `LH-FA-DEV-009` (Sandbox-Hinweis, `devcontainer.sandbox.credentials`) | geprüft, ohne Befund (nennen nur Git-Zugangsdaten bzw. die Kennung; kein Widerspruch zur Erweiterung, Nachzug liegt im Folge-Slice zur Spezifikation) | <!-- d-check:ignore (Kennungen als Zitat) -->

Nicht geprüft:

- Code (`devcontainer.json.tmpl`, `doctor_sandbox.go`, `devcontainer_sandbox.go`) — der Diff berührt keinen Code; gelesen nur als Ist-Stand-Beleg für F-1.
- `make gates`, `make test-docker`, `make verify-depguard` — kein Code-, Docker- oder Importregel-Pfad berührt; nur `make docs-check` gelaufen.
- DoD-Zeilen des Slice, Closure-Notiz — Verifier-Aufgabe.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Changelog-Zusage vor Code-Stand · Begriff „Agenten-Feature“ im Vertrag ohne Abgrenzung · Klammerzusatz erklärt oder begrenzt unklar · Überschrift schmaler als Zusage · Slice-Entwurf folgt der Entscheidung nicht nach

## Verdikt

**Merge-blockierend:** nein — kein HIGH; F-1 (MEDIUM) betrifft einen unveröffentlichten Changelog-Eintrag und ist vor dem nächsten Release-Cut zu klären, nicht vor dem Lastenheft-Change.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
