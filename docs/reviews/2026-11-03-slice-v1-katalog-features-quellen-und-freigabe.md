# Review-Report: slice-v1-katalog-features-quellen-und-freigabe — 2026-11-03

**Review-Art:** Code-Review — gegen Plan + Vertrags-Carrier + Hard Rules
(kein DoD-Abgleich; der ist Verifier-Frage).

**Gegenstand:** Commit `06e81ac4` (23 Dateien, +597/−113) des Vorgangs
[`slice-v1-katalog-features-quellen-und-freigabe`](../plan/planning/done/slice-v1-katalog-features-quellen-und-freigabe.md);
Vertrags-Range `d3c24eb6..HEAD`, Vertrags-Norm
[`slice-cr-devcontainer-katalog-quellen`](../plan/planning/done/slice-cr-devcontainer-katalog-quellen.md)
(Übernahme `151a7976`, Fix-Runde `b122f697`).

**Skill:** `.harness/skills/reviewer.md` @ `06e81ac4` · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** glm-5.3-flash · **Datum:** 2026-11-03

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-v1-katalog-features-quellen-und-freigabe`](../plan/planning/done/slice-v1-katalog-features-quellen-und-freigabe.md) (Plan samt DoD, §3-Tabelle, neuer §3.10-Abgleich)
- [`slice-cr-devcontainer-katalog-quellen`](../plan/planning/done/slice-cr-devcontainer-katalog-quellen.md) samt [`cr-vorschlag-katalog-quellen.md`](../plan/planning/done/cr-vorschlag-katalog-quellen.md) (Stellen (a)–(d), Mess-Basis M1/M2)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Lastenheft 0.33.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
- [Messprotokoll der N1/N2-Messung](../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md) (M1 Token-Verweigerung, M2 Schlüsselform-Probe, Tarball offen)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules §3), [`harness/conventions.md`](../../harness/conventions.md)
- vorherige Findings am Bereich: [CR-Review-Report](2026-11-03-slice-cr-devcontainer-katalog-quellen.md) (F-1..F-4 et al.)

> **Zitier-Form** *(dieser Block bleibt stehen).** Kennung, nicht Adresse;
> Baseline-Stellen als **Tag + Pfad in Inline-Code**.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der unveränderte Erstsatz von [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) („externe Feature-Quellen **müssen** als URL mit Schema `http`/`https`/`oci` angegeben werden") widerspricht dem nachgeschobenen gemessenen Satz („die Registry-Referenz verarbeitet, `oci`/`https` unverarbeitbar") und der neuen Beleg-Zeile, die die Registry-Referenz als **die** verarbeitete Form festhält; die Tests pinnen das Gegenteil des Erstsatzes. Der CR-Abgleich §7 trägt den Ausgleich nicht (Stellen (a)–(d) benennen ihn nicht; `151a7976`/`b122f697` ließen den Erstsatz stehen) — Antwort auf die Übergabe-Frage: **nicht getragen**. Szenario für den nächsten Änderer: er liest den ersten Satz als Norm und schreibt Dokumentation oder Fixtures auf Schema-URLs zurück. | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Zusage-Seite), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) | `spec/spezifikation.md` · „Externe Feature-Quellen müssen als URL mit Schema" | review-only (d-check prüft keine Inhaltskohärenz) | spec-erstsatz-im-widerspruch-zur-form-regel |
| F-2 | MEDIUM | Das [Messprotokoll](../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md) widerspricht sich über seine eigene Messung: Umgebung-Satz („die Dev-Container-CLI trat nicht auf … in diesem Protokoll nicht geführt") und Grenze 3 („hier nicht geführt", „keine Aussage über Form-Schwäche der externen Quellen") behaupten, M2 sei nicht geführt, während §M2 (seit `d3c24eb6`) die Schlüsselform-Probe samt Rohdaten trägt. Alle 0.33.0-Carrier dieses Commits (Beleg-Tabellenkopf „Messbasis 2026-11-03, Registry- und CLI-Lesung", CHANGELOG, Historie) zitieren genau dieses Protokoll. Szenario: der Verifier liest die Grenze, stuft die N2-Messbasis zum Closure-Zeitpunkt als fehlend ein; oder ein späterer Änderer folgt „keine Aussage über Form-Schwäche" und öffnet Schema-Formen wieder. | [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md); [`AGENTS.md`](../../AGENTS.md) §3.10 | `docs/messungen/`[`slice-v1-katalog-features-quellen-und-freigabe`](../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md)`/README.md` · „hier nicht geführt" | review-only | messprotokoll-grenze-im-widerspruch-zur-messung |
| F-3 | MEDIUM | Der Slice-Plan trägt beim Kopf-Update weiterhin den abgelaufenen Zustandsmarker, der CR liege „in `in-progress/`" — der Verweis daneben zeigt (richtig) nach `done/`, das Verzeichnis ist die Wahrheit ([`roadmap.md`](../plan/planning/in-progress/roadmap.md) trägt die CR-Closure). Dritter Kontext derselben Zustandsmarker-Familie am WIP-Wechsel (nach F-1/F-2 des CR-Reviews und der CR-Bemerkung „dritter Kontext in Folge" in dessen §7) — der recurring-LOW-Hebel macht es zum MEDIUM. Kein Gate fängt Zelleninhalte. | [`AGENTS.md`](../../AGENTS.md) §3.10; [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | `docs/plan/planning/in-progress/`[`slice-v1-katalog-features-quellen-und-freigabe`](../plan/planning/done/slice-v1-katalog-features-quellen-und-freigabe.md)` · „liegt in `in-progress/`" | review-only | abgelaufener-zustandsmarker-am-lifecycle-wechsel |
| F-4 | MEDIUM | Die Konfigurationstabelle des [Handbuchs](../user/benutzerhandbuch.md) nennt als Wert der `devcontainer.features.<name>.source` „**URL** aus der Freigabeliste" — genau die Form, die der Umsetzungs-Diff an beiden Aufrufstellen abweist; die Nachbarzeile der Allowlist wurde im selben Diff auf „Registry-Referenzen (`host/pfad`)" berichtigt, diese Zeile nicht. Szenario: der Nutzer liest die Referenztabelle, trägt eine URL als Quelle ein und bekommt Exit-Code `10`. | [LH-FA-DEV-003](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | `docs/user/benutzerhandbuch.md` · „URL aus der Freigabeliste" | review-only (d-check prüft Zellinhalte nicht) | doku-tabelle-driftet-gegen-quell-form |
| F-5 | LOW | Die neue Schritt-6-Zeile der [SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Beleg-Tabelle nennt „am Nutzerpfad" `TestRun_Generate_UnprocessableFeatureSourceFormIsRejected` — dieser Test pinnt die Erzeugungs-Abweisung, nicht `config set …source` am CLI-Pfad; ein Nutzerpfad-Test der Schritt-6-Abweisung existiert nicht (Repo-Muster: `config set`-Regeln pinnen am Dienst-Pfad). Szenario: eine spätere CLI-Änderung des `config`-Fehlerwegs verliert laut Tabelle scheinbar einen Nutzerpfad-Pin, der keiner ist. | [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) | `spec/spezifikation.md` · „und am Nutzerpfad `TestRun_Generate_UnprocessableFeatureSourceFormIsRejected`" | review-only | beleg-zelle-nennt-test-anderer-stelle |
| F-6 | LOW | Antwort auf die Übergabe-Frage zur form-permissiven Freigabeliste: der einzige Prosa-Ausweis („die Freigabeliste selbst bleibt form-permissiv") steht in einer Historie-Zelle der Spezifikation — Chronik, nicht normativer Träger; keine Spec- oder Handbuch-Stelle trägt das Verhalten (`config set devcontainer.featureSources.allow` nimmt Schema-Formen an; nur beiläufig gepinnt durch den erwarteten `wantNew` von `TestConfigSet_FeatureSourcesAllow_Append`). Solange keine Zusage der Strenge besteht, bricht nichts — aber der Ausweis veraltet mit der nächsten Historie-Rotation. | [`AGENTS.md`](../../AGENTS.md) §3.7 (Zusagen tragen auflösbare Anker) | `spec/spezifikation.md` §7 · „die Freigabeliste selbst bleibt form-permissiv" | review-only | verhaltensausweis-nur-in-historie-zelle |
| F-7 | INFO | Übergabe-Frage (3), alte Plan-DoD gegen CR-Muster: der Docker-Test eines Baus aus einer externen Quelle, die Tarball-Probe und der doctor-Drift-Vergleich „derselben Form" aus der unveränderten DoD-Zeile 1 des Plans sind nicht geliefert; die Plan-§3-Zeile erklärt sie zum „Reste in §6/§7 des Verifier" und der CR §7 (*Folge-Slice*) trägt den schmaleren Umsetzungs-Scope. Diese Auslegung ist die des Implementers — ob Grenze oder offene Pflicht, entscheidet der Verifier je Closure-Trigger (der DoD-Haken bleibt sein Instrument, nicht die §3-Prosa). Zuständigkeit: Verifier. | [`harness/verification.md`](../../harness/verification.md) (Rolle) | `docs/plan/planning/in-progress/`[slice-v1-katalog-features-quellen-und-freigabe](../plan/planning/done/slice-v1-katalog-features-quellen-und-freigabe.md)` §3 · „Reste in §6/§7 des Verifier" | review-only | dod-rest-von-cr-muster-abgegrenzt |
| F-8 | INFO | Übergabe-Frage (2), Schärfe der Abweisung: das Prädikat weist **jede** Form mit `://` ab, gemessen sind nur zwei (`oci://…:<tag>`, `https://…:<tag>`); `http://` und ungemessene Schema-Formen fallen darunter ohne eigene Messaussage — ausgewiesen in CHANGELOG (Not-promised) und Historie-Zeile („URL-Schema-Formen außerhalb der zwei tragen keine eigene Messaussage"). Verträglich mit dem verallgemeinernden [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Bullet 0.33.0; der Umsetzungsausweis trägt dieselbe Regel — kein Bruch, Auslegungsrest beim Verifier. | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | `internal/hexagon/application/devcontainer_features.go` · `strings.Contains(source, "://")` | `make test` (pinnt die zwei gemessenen Formen) | abweisung-schärfer-als-messung |
| F-9 | INFO | Die neue Registry-Referenz-Grammatik nimmt degenerierte Formen an, die die alte Grammatik abwies (Slash > 0 genügt: `ghcr.io/` mit leerem Pfad, nackte `host/pfad`-Kurzformen); TestValidateFeatureSource pinnt nur die Abweisung des reinen Hostworts. Nicht gemessen, nichts versprochen — Beobachtung zur Stage-1-Grammatik ohne Änderungspflicht. | Maintainability | `internal/hexagon/application/devcontainer_features.go` · `slash <= 0` | `make test` | stage-1-grammatik-toleranter-als-vorher |

---

## Übergabe-Fragen des Auftraggebers — Antwort-Zuordnung

1. **[`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)-Erstsatz-Spannung** — nicht über die CR-Träger getragen; F-1 (MEDIUM).
2. **Form-permissive Freigabeliste** — Ausweis existiert, aber als Historie-Chronik (F-6, LOW); die Nutzertabelle sagt strikter als das Verhalten (F-4).
3. **Alte Plan-DoD vs. CR-Muster** — vom Implementer als Grenze abgeschrieben, ist aber Verifier-Frage (F-7, INFO).
4. **Fehlgeschlagene Zwischen-Mutation** — der Endstand ist vollständig: kein `cpp`/`postgres-client`-Rest in `internal/` (nur die unberührten Namens-Grammatik-Beispiele in `internal/hexagon/domain/featurename_test.go`), keine Reste in `docs/user/` und beiden READMEs (18-Feature-Zählung stimmt), Reste in `spec/` nur als Chronik/Evidence-Zeilen; das geklebte Historie-Fragment (V-1) ist im Endstand beseitigt. Kein Finding; die Entfernungs-Diffs (`ociFeatureCatalogue`, Doc-Kommentar, Katalog-Tabellentest) sind vollständig.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `devcontainer_features.go` — Entfernungs-Diff (zwei Katalog-Einträge samt Doc-Kommentar zur „unread"-Aussage) | geprüft, ohne Befund: vollständig, keine verwaiste Aussage |
| `generate.go` — Extraktion `rejectFeatureSourceViolations` | geprüft, ohne Befund: Verhalten der Freigabe-Prüfung unverändert (Wrap-Form und Meldung identisch), nur Reihenfolge Form-vor-Freigabe neu, vor jedem plan and write |
| `config.go` — `revalidateFeatureEntry` | geprüft, ohne Befund: Form-Prüfung vor Listen-Zugehörigkeit, `--yes` kann die Prüfung nicht umgehen (die Fachprüfung liest keinen Yes-Zweig), `u-boot.yaml` bleibt byte-identisch |
| Dual-Classifier-Regel | geprüft, ohne Befund: kein neuer Driving-Sentinel; `ErrConfigValueInvalid` ist zentral klassifiziert, Nutzerpfad-Pin `--json` Code [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) vorhanden (`TestRun_Generate_UnprocessableFeatureSourceFormIsRejected`) |
| Sprachvertrag [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | geprüft, ohne Befund: alle neuen Meldungs- und Kommentar-Texte englisch |
| §3.7/§3.9-Kommentar-Selbstcheck über den Diff | geprüft, ohne Befund: keine „früher"-Formulierungen, keine Symptom-Chronik; die zwei Grep-Treffer „formerly"/„no longer" liegen außerhalb des Diffs |
| Dateisicherheit ([`LH-NFA-REL-001`](../../spec/lastenheft.md#lh-nfa-rel-001--kein-stilles-überschreiben)) | geprüft, ohne Befund: nur Vorab-Abweisungen, kein neuer Schreibpfad, Write-los-Pinntests je Aufrufstelle |
| Architektur/Importregeln ([`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)) | geprüft, ohne Befund: keine neuen Imports, keine Port-Änderung, keine depguard-Berührung |
| 11 Fixture-Dateien (Set-/Generate-Pfade) | je Hunk geprüft, ohne Befund: Berichtigung nur zur Registry-Referenz, kein Vergleichsverlust (Sandbox-/Freigabe-/Sorted-/Host-Case-Assertionen blieben bei gleichem Verhalten) |
| Neue Test-Dateien (`generate_sourceform_test.go`, cmd-Nutzerpfad-Test) | geprüft, ohne Befund: alle im Auftrag genannten Test-Namen existieren; Abhilfe-Kette (Exit 0, Schlüssel gerendert) am Nutzerpfad gepinnt |
| [`spec/lastenheft.md`](../../spec/lastenheft.md) im Vertrags-Range | geprüft, ohne Befund: Umsetzungs-Commit rührt es nicht; Range-Änderungen stammen nur aus den CR-Commits (`b122f697`, `151a7976`) — Beleg-Auftrag, nicht Umsetzungs-Edit |
| Beleg-Tabellen-Neuaufnahme ([`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)) | geprüft, ohne Befund: alle zitierten Test-Namen existieren (F-5 betrifft nur die Zuständigkeits-Zelle der Schritt-6-Zeile) |
| [`CHANGELOG.md`](../../CHANGELOG.md) (zwei Changed-Einträge) | geprüft, ohne Befund: Promised/Not-promised deckt die Messgrenze (CLI 0.89.0, drei Formen über existierenden Tag, Tarball offen) — F-2 betrifft das Protokoll, nicht den Eintrag |
| [`harness/replay.md`](../../harness/replay.md) Golden-Zeile | geprüft, ohne Befund: Inhalt deckt Form-Abweisung, Freigabe-Reihenfolge, `config set`-Gleichlauf; keine bestehende Zeile berührt |
| Lifecycle/Lifecycle-Moves, ADR-Immutabilität, Baseline-Dateien | geprüft, ohne Befund: kein Move im Commit, keine ADR-Berührung, keine Gate- oder Schwellen-Lockerung (Katalog-Tabellentest geringert aus Vertrags-Grund) |
| `doctor.go` und die drei als „unverändert" geführten Testdateien | geprüft, ohne Befund: Dateien sind im Commit nicht enthalten (Diff-Abgleich), Plan-Aussage stimmt zur Stelle |

**Nicht geprüft** (mit Grund):

- `make gates` / `make test-docker` nicht erneut gelaufen — Voll-Läufe sind Verifier-Instrument; Implementer-Lauf ist belegt (`record-gates`-Stempel `.harness/state/gates-passed.head` auf `06e81ac4`). Gelaufen: `make docs-check` (vom Report mitgetragen, Ergebnis im Lauf-Status).
- `tools/messung/katalog-quellen/registry-quellen.sh` nur als Beleg-Kontext gelesen, nicht inhaltlich auditiert — Messobjekt außerhalb des Commits.
- Inhaltliche Kohärenz-Prüfung der Spec-BODY-Änderungen der CR-Commits (`151a7976`, `b122f697`) — Gegenstand des CR-Reviews (liegt vor); hier nur soweit gelesen, wie die Umsetzungs-Belege auf sie zeigen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** spec-erstsatz-im-widerspruch-zur-form-regel · messprotokoll-grenze-im-widerspruch-zur-messung · abgelaufener-zustandsmarker-am-lifecycle-wechsel · doku-tabelle-driftet-gegen-quell-form · beleg-zelle-nennt-test-anderer-stelle · verhaltensausweis-nur-in-historie-zelle · dod-rest-von-cr-muster-abgegrenzt · abweisung-schärfer-als-messung · stage-1-grammatik-toleranter-als-vorher

## Verdikt

**Merge-blockierend:** ja (0 HIGH / 4 MEDIUM). Kein Code-Defekt gefunden — die produktiven Pfade, Tests und Doku-Abszüge tragen den
0.33.0-Schnitt; alle vier MEDIUM sind Träger-Reconciliation (Spec-Erstsatz, Messprotokoll-Grenze, Plan-Marker, Handbuch-Zelle) und
in einer Fix-Runde berichtigbar. Kein Finding wurde wegen Implementer-Widerspruchs herabgestuft; die Entfernungs- und
Fixture-Berichtigungen sind gegen die CR-Träger nachgelesen, nicht übernommen.

**Übergabe:** Findings F-1..F-6 gehen an den Implementer (Fix-Runde); F-7..F-9 an den Verifier.
Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler
([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
und [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
gewinnen je einen neuen Fund; die Zustandsmarker-Familie tritt zum dritten Kontext in Folge auf).
Dieser Report ist ein Lauf-Beleg; DoD-Closure und Verification-Evidence sind Verifier-Aufgabe.