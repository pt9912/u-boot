# Review-Report (Nachlauf): [`slice-v1-devcontainer-host-git-uebernehmen`](../plan/planning/done/slice-v1-devcontainer-host-git-uebernehmen.md) — 2026-11-01

**Review-Art:** Nachlauf der Review-Fix-Runde (2. Lauf, frischer Kontext) gegen den Erst-Report [`review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen.md`](review-2026-11-01-slice-v1-devcontainer-host-git-uebernehmen.md).

**Gegenstand:** `git diff 08633e7f..HEAD` (Commit `e5a8cb63`, „Review-Fix-Runde“). 10 Dateien, 241 Einfügungen, 19 Löschungen; alles vor `08633e7f` ist der gelesene Stand des Erst-Laufs.

**Skill:** `.harness/skills/reviewer.md` @ `e5a8cb63` · **Datum:** 2026-11-01

---

## Stand der Vorgänger-Findings

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (HIGH, ADR-Verweis abwärts) | behoben | [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md): alle vier Prosa-Verweise entfernt — „Anlass, gelesen am 2026-11-01 gegen den Code“, „wurden vor dem Code abgewogen (Messung: Protokoll in `docs/messungen/`)“, „Der Projektinhaber hat am 2026-11-01 entschieden“ ohne Slice-Nennung, „Entscheidung des Projektinhabers 2026-1, deren Messprotokoll ein Ablage-Ort und keine Referenz auf einen Slice ist |
| F-2 (HIGH, Kommentarzahl) | behoben | `internal/e2e/devcontainer_hostgit_docker_test.go`: „(33 + the newline of printenv = 34) is the expected count in the session.“ — 33 Zeichen Platzhalter, `hostGitSessionLen = "34"`; Kommentar und Konstante stimmen |
| F-3 (MEDIUM, Hilfetext) | behoben | `internal/adapter/driving/cli/devcontainer.go`: der Hilfetext führt die Auflösung nach (Sandbox-Profil, Schalter, `GIT_TOKEN` nicht gesetzt, Quelle `gh auth token` genannt, Wert nie, fehlende Quelle ohne Abbruch); englisch, ohne Kennung, ohne Chronik |
| F-4 (MEDIUM, vier gepinnte Bedingungen) | behoben | vier neue Tests: `TestDevcontainerUp_HostGitToken_SwitchWithoutReferenceWarns` (pins Warnungstext „references no ${localEnv:GIT_TOKEN}“ und: nichts aufgelöst, Quelle leer), `TestDevcontainerUp_HostGitToken_IdentityReadFailureWarns` (pins Warnung mit Ursache, keine Identität übergeben), `TestLauncher_Up_MasksResolvedValue` (pins Maskierung im Fortschritt und im Fehlertext), `TestDevcontainerUp_HostGitToken_InstanceStillResolves` (pins Auflösung unter `--instance inst-1`) |
| F-5 (LOW, Slice §7 Zustand) | behoben | „Noch nicht geschlossen; der Slice liegt in `in-progress/` (Review-Fix-Runde).“ |
| F-6 (LOW, Handbuch-Leitsatz) | behoben | „… oder, mit dem Opt-in-Schalter des nächsten Abschnitts, die er auf dem Host ermittelt hat —, und zwar als Name …“ |
| F-7 (LOW, Architektur-Sicht) | behoben | `spec/architecture.md` §2.4 trägt `HostGitSource` (`GitToken`/`ConfigValue`) mit Grenze (nur `up`/`shell`, keine eigene Konfigurationsdatei, optionale Verdrahtung); Verweise nur auf `LH-*` und [[`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) — kein ADR-, Slice- oder Commit-Bezug (`AGENTS.md` §3.4 erfüllt) |
| F-8 (LOW, Port-Doc Fehlerpfad) | behoben | Der Doc von `HostGitSource.GitToken` trägt jetzt die Fehlerbedeutung als Invariante („The error carries only a failed environment read; the adapters read the environment of their own process.“) statt der unerreichten Produktaussage |
| F-9 (INFO, open/README-Bestand) | offen, für Closure akzeptiert | Nachzug mit der Closure |
| F-10 (INFO, Belege) | Teilbehoben | Rot-Beleg ohne `/tmp`-Verweis („Meldung oben wörtlich; die Mutation ist zurückgenommen“); die CI-Evidenz des Docker-Tests und die Mutation-Wiederholung bleiben Verifier-Aufgabe |
| F-11 (INFO, Restcontainer) | offen, für Verifier akzeptiert | Messprotokoll-Regie (`aufraeumen`) |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-N1 | LOW | Der Abgleich der Plan-Tabelle (Slice §3) wurde nach der Fix-Runde nicht wiederholt, obwohl §3.10 ihn nach „jeder späteren Änderung“ verlangt; die Fix-Runde berührt zwei Dateien ohne Plan-Eintrag: `spec/architecture.md` (F-7-Fix) und `docs/reviews/…` (Review-Berichte). Die Test-, CLI-, Handbuch- und ADR-Dateien sind in der Tabelle geführt. | `AGENTS.md` §3.10 (Abgleich nach jeder späteren Änderung); Klasse [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`slice-v1-devcontainer-host-git-uebernehmen.md`](../plan/planning/done/slice-v1-devcontainer-host-git-uebernehmen.md) · „wiederholt nach jeder Fix-Runde“ | Review-only | Abgleich nach der Fix-Runde stehen geblieben |
| F-N2 | INFO | Der Inline-Kommentar im neuen Test „t.Setenv verbietet t.Parallel in dieser Datei: deterministisch ohne Umgebungslesung des Ports“ begründet das Setzen von `NEUTRAL=x` mit einer Umwelt-Unabhängigkeit, die das `Setenv` selbst nicht herstellt — die Unabhängigkeit trägt die Fake-Quelle, die `Setenv` wirkt nur über das Parallel-Verbot. Harmlos, aber der nächste Änderer kopiert die Begründung mit. | `AGENTS.md` §3.7 (Kommentar nennt die geltende Zusage) | `internal/hexagon/application/devcontainer_start_test.go` · „deterministisch ohne Umgebungslesung des Ports“ | Review-only | Kommentar begründet mehr, als die Zeile wirkt |

## Negativbefunde

Geprüft, ohne Befund:

- [`ADR-0040`](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) gegen `AGENTS.md` §3.4/§3.5 und die [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)-Matrix (kein abwärts-Verweis mehr, Inhalt sonst unverändert, Index nicht betroffen)
- Hilfetext von `up` gegen §3.7 (Indikativ, Zustandsform, kein Chronik-Muster) und [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) (englisch)
- Die vier neuen Tests gegen die Plan-Wächter-Tabelle (je Wächter, für den die Fix-Runde eine Lücke schloss, jetzt ein pinnender Fall) und gegen §3.9 (c) (keine Helfer-Einfügung zwischen Doc-Kommentar und Funktion)
- `spec/architecture.md` §2.4-Neuabschnitt gegen §3.4 (kein ADR-/Slice-/Commit-Bezug, kein Wellenbezug) und gegen die Portliste des Use-Cases
- Handbuch-Leitsatz (F-6) gegen den neu formulierten Zugangsdaten-Absatz — konsistent
- Port-Doc `GitToken` (F-8) gegen den Adapter (`internal/adapter/driven/git/hostsource.go`) — Invariante jetzt korrekt gelesen
- Slice §7 und Rot-Beleg-Zeile gegen §3.7 (Zustand, keine Chronik, keine Wegwerf-Pfade `/tmp`)
- §3.9 (a)/(b) über den Fix-Diff: kein Chronik-Kommentar, kein Kommentar neben geändertem Zahlenwert (F-2 behebt gerade diesen)

Nicht geprüft:

- Ausführung der Test-Suiten — `make test` und `make lint` sind laut Implementer grün (Host-Lauf); der Nachweis für einen Handoff bleibt `make gates` (Stempel), Verifier
- `make test-docker` — der geänderte E2E-Kommentar ist produkt-inaktiv; die vier HostGit-E2E-Tests sind seit `08633e7f` unverändert, der Gesamtlauf des Handoffs bleibt Beleg
- Die Mutation hinter `TestLauncher_Up_MasksResolvedValue` („aufgelöste Werte aus `passedValues` weglassen“ färbt ihn rot) — Verifier führt die Rot-Belege, wie im Slice §3 vereinbart

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Verdikt:** **Handoff-reif für den Verifier: ja.** Keine HIGH- und kein MEDIUM-Befund des Erst-Laufs bleibt offen; alle vier Wächter, die F-4 benannte, tragen jetzt je einen pinnenden Test. F-N1 (Abgleich nachholen, zwei Dateien in die Plan-Tabelle) ist vor der Closure nachzuziehen; F-N2 ist ein Kommentar-Feinschliff ohne Blockierwirkung. Offen für Verifier/Closure: CI-Evidenz des Docker-Tests und die Mutation-Rot-Belege (F-10-Rest), F-9 (open/README-Bestand bei der Closure), F-11 (Restcontainer-Inventur), Stempel `make gates` auf dem Abschluss-Commit.