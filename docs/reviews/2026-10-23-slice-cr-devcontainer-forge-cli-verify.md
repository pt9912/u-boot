# Verify-Report: slice-cr-devcontainer-forge-cli — 2026-10-23

**Rolle:** Verifier (`harness/verification.md`), frischer Kontext. **Modell:** claude-sonnet-5-5.

**Gegenstand:** Slice [`slice-cr-devcontainer-forge-cli`](../plan/planning/done/slice-cr-devcontainer-forge-cli.md), Welle `welle-forge-cli-gh-glab`, Diff `dd4d5e1a..ffcefebd`; Review-Report [`2026-10-23-slice-cr-devcontainer-forge-cli.md`](2026-10-23-slice-cr-devcontainer-forge-cli.md).

**Urteil: schließbar** (offen bleiben nur Closure-Schritte der Implementer-Rolle, siehe Punkt "Offen für die Closure").

## Sensoren (selbst gefahren)

| Lauf | Exit-Code |
|---|---|
| `make gates` | 0 (Log im Scratchpad, letzte Zeile `GATES EXIT=0`) |
| `make doc-trace` | 0 |
| `make doc-planning` | 0 |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Code-, Docker- oder Importregel-Pfad im Diff (nur `spec/lastenheft.md` und Planung) |

## DoD Punkt für Punkt

1. **Freigabe vor der Änderung:** erfüllt. §7 *Freigabe* trägt Stelle 1 (Frage 1 Ja, Frage 2 Nein), Stelle 2 wie Entwurf, Stelle 3 keine Änderung, Stelle 4 Version 0.22.0. Commit-Reihenfolge `3b2179c7` (Freigabe) vor `5b55250f` (Lastenheft). Die Freigabe des Projektinhabers selbst ist im Repo nicht belegbar, nur durch die Zusage des Auftraggebers dieses Laufs (2026-10-23); das ist die Grenze der Verifikation.
2. **Wortlaut im Lastenheft, wortgleich mit der Freigabe:** erfüllt. Zeile 413 gegen den in §7 als freigegeben ausgewiesenen Absatz, maschinell verglichen (Linkziele normalisiert, da §7 relativ zum Slice verlinkt): Text identisch. Halbsatz "in beiden Profilen (Standard und Sandbox)" vorhanden; keine Hostnamen (`api.github.com`, `github.com`, `gitlab.com` fehlen im Diff); Spiegelstriche `GitHub CLI (gh)`, `GitLab CLI (glab)`. Stelle 2 ([`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), drei Passagen) entspricht dem Entwurf-*Neu* (Einleitungssatz um Forge-Werkzeug, Dokumentation mit `GH_TOKEN`/`GITLAB_TOKEN`, `doctor` bei aktivem `gh`/`glab`). Beleg-Tabelle B1 bis B11 in §7 vorhanden, je Bedingung Messaussage oder Test-Auftrag.
3. **Versionsfeld und Historie:** erfüllt. Versionsfeld 0.22.0; Historie-Zeile 0.22.0 (2026-10-23) nennt nur Stelle 1 und 2, trägt "in beiden Profilen", keine Hostnamen, nur Lastenheft-Anker.
4. **`make gates`, `docs-check`, `doc-trace`:** erfüllt (Exit-Codes oben).
5. Review/Doku/Closure-Notiz/Register/Risikoausgänge/Paarungen (offene Haken in §2): Aufgabe der Closure, nicht Gegenstand dieser Verifikation.

## Nachlauf-Prüfung Klammer (Commit ffcefebd)

- `git diff dd4d5e1a..HEAD -- spec/lastenheft.md` zeigt in Zeile 413 die Klammer "(GitHub Enterprise, GitLab Self-Managed)" bei "selbst gehostete Instanzen"; Commit `ffcefebd` ändert im Lastenheft genau diese Zeile (1 Einfügung/1 Löschung, Slice §7 ebenfalls eine Zeile).
- Der gelieferte Satz S6 lautet im Lastenheft wörtlich: "… Telemetrie-/Update-Hosts einer interaktiven Sitzung, selbst gehostete Instanzen (GitHub Enterprise, GitLab Self-Managed), andere Versionen als die gepinnten; außerdem `glab duo`, `glab orbit` und das Nachladen weiterer Programme durch `glab`." Das ist der vom Auftraggeber genannte freigegebene Vorschau-Wortlaut; auch der Entwurf (`wortlaut-entwurf.md`) trägt dieselbe Klammer. Slice §7 nennt sie jetzt als Teil des freigegebenen Absatzes (Vorschau hatte sie). Keine Zusageverschiebung: es ist eine Nicht-Zusage, die Beispiele konkretisieren. Messbarkeit/Decken-Regel unberührt (keine Kennung nach unten).
- Der Review-Report liegt vor dieser Änderung; ein eigenes Nachlauf-Review (DoD-Punkt "Review", `BEO-DOK/vertragsaenderung-nach-review`) liegt nicht als Datei vor. Diese Prüfung liest genau den Diff und ersetzt es inhaltlich; als Datei ist sie dieser Bericht. Der Planner kann das als Nachlauf-Review führen oder eigens nachziehen.

## Weitere Prüfungen

- **Decken-Regel (AGENTS.md §3.4):** Die hinzugefügten Zeilen im Lastenheft enthalten keine Treffer auf ADR, Slice, Welle, Roadmap, SPEC, ARC; Links nur auf Lastenheft-Anker (existieren; `docs-check` grün).
- **Test-Aufträge je Bedingung (Welle §5):** Die Welle trägt den Auftrag "Aufträge des Lastenheft-CR an den Code-Slice" mit Verweis auf die Beleg-Tabelle, ausdrücklich "beide Profile" (B1) sowie B5 (Ausweisung beider Archive und `gh`-Attestierung), B6 (Hosts von `gh` in der erzeugten Default-Allowlist) und B9 (Token weder im Dockerfile noch in erzeugter Datei noch in Mount). Damit sind die Review-Findings F-2 und F-1 (Ausweisung der Abweichungen: verdichteter Wortlaut, drei Punkte in §7) behoben. B2, B3, B4, B7, B8, B10, B11 sind in der Slice-Tabelle belegt. Der Code-Slice selbst existiert noch nicht; die Aufträge sind Aufträge, keine Tests.
- **Messstände-CR (F-3):** Zeilenangaben selbst gegengelesen mit `grep -n` auf `spec/lastenheft.md`: 411 (Ollama-Absatz), 467 und 469 (`podman`), 1583/1588/1589/1591 (Historie 0.7.0/0.12.0/0.13.0/0.15.0), 1096 ([`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort)), 428 (UID-Beispiel), 1504 und 1579 ([`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung)/0.3.5) stimmen. Versionsangaben "nach 0.22.0" / "zuletzt 0.22.0" stimmen mit Versionsfeld. Zeile 413 des Slice-Vorabsuche-Hinweises stimmt.
- **Nutzer-Doku:** `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `spec/spezifikation.md`, `spec/architecture.md` nicht im Diff (`git diff --stat`). Passt zu Plan ("unverändert") und AGENTS.md §5 Regel 6.
- **Plan-vs-Diff (§3.10):** `git diff --stat dd4d5e1a..HEAD`: roadmap.md, Slice (Move plus Inhalt), open/README.md, Messstände-CR, Welle, Review-Report, lastenheft.md. Jede Datei hat Plan-Eintrag (Review-Report ist laut Plan Reviewer-Werkstück); kein Eintrag ohne Diff außer `observations/` (laut Plan erst bei Closure).
- Lifecycle: Slice liegt nur in `in-progress/`, Roadmap nennt ihn; `doc-planning` grün.

## Findings

Keine blockierenden. Hinweise:

- **V-1 LOW (Prozess):** Zwischen Review und Nachlauf-Prüfung wurde Lastenheft-Text (Klammer) geändert; ein separates Nachlauf-Review-Artefakt fehlt, dieser Bericht deckt den Diff. Closure sollte das in §7 vermerken.
- **V-2 INFO:** Die Freigabe des Projektinhabers ist nur über den Auftraggeber belegt (nicht im Repo auditierbar).

## Offen für die Closure (Implementer/Planner)

Was hat funktioniert / Was ging anders, Steering-Loop-Eintrag, Register-Fortschreibung (`BEO-DOK/vertragsaenderung-nach-review`, `BEO-DOK/zusage-weiter-als-messung`), Risikoausgänge §6, DoD-Haken, danach `git mv` nach `done/` und Roadmap-Zeile zurück auf "Nichts in Arbeit".
