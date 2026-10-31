# Review-Report: spec-prozesszustaende — 2026-10-31

**Review-Art:** Design (Doku-Diff gegen Auftrag, Hard Rules und Lastenheft)

**Gegenstand:** uncommittete Änderungen (`git diff HEAD`) an `spec/architecture.md`, `spec/spezifikation.md` und `docs/maintainer/zusagen-und-messstaende.md` (neuer Abschnitt 3.15). Auftrag des Projektinhabers vom 2026-10-31: Spec als Zielbild ohne Prozesszustände, Messstände und Prüfumgebungen; Bedeutung normativer Aussagen unverändert; Messdetails verlustfrei ins Maintainer-Dokument.

**Skill:** `.harness/skills/reviewer.md` (Stand des Arbeitsbaums) · Kopie der Vorlage per `cp`, `diff -q` gegen `.harness/baseline/v6.18.1/templates/docs/reviews/review-report.template.md` vor dem Füllen: gleich.
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-31. Der Reviewer hatte kein Edit-Werkzeug; der Report ist mit Write geschrieben.

**Eingangs-Kontext:**

- Auftrag des Projektinhabers (siehe Gegenstand)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge: Docker auf Linux-Host mit lokaler Engine, `amd64`; keine Engine-Versionen)
- [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- `AGENTS.md` §3.4, §3.5, §3.7, §3.9, §3.10

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Fundstelle des neuen Abschnitts 3.15 stimmt nicht. Das Protokoll [`slice-v1-devcontainer-instanz-messung`](../plan/planning/done/slice-v1-devcontainer-instanz-messung.md) enthält keines der Stichworte `sleep 3600`, `Pseudo-Terminal`, `script`, `postAttach`, `28.0.4`, Login-Shell; es misst Namen, Abbruch und Haenger (Rohdaten `namen`, `stufe2` bis `stufe4`). Die Aussagen zu Sitzung, Signalen, Shell und Engine 28.0.4 stammen aus dem Planungsdokument [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) (Docker 29.9.2, Basis-Image `devcontainers/base:debian`, Test unter `script`, Runner-Engine 28.0.4) und die Lebenszyklus-Aussage aus dem Test `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` (`internal/e2e/devcontainer_start_docker_test.go`). Es gibt zu den Shell-Aussagen kein Messprotokoll unter `docs/messungen/`. Der Zusatz „Zuordnung der einzelnen Aussage ist nicht geprüft“ benennt die Lücke, löst sie nicht. | Auftrag (Fundstellen stimmen); `AGENTS.md` §3.10 | `docs/maintainer/zusagen-und-messstaende.md` · „**Fundstelle:** [[`slice-v1-devcontainer-instanz-messung`](../plan/planning/done/slice-v1-devcontainer-instanz-messung.md)]“ | ja, Review-only (Stichwortsuche im Protokoll leer) | Fundstelle trägt die Aussage nicht |
| F-2 | MEDIUM | Der Punkt „Messbedingungen“ in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) war ein Vorbehalt: die Aussagen zu Ende der Sitzung, Signalen und Art der Shell gelten für einen `sleep 3600`-Container des Basis-Images, ausdrücklich nicht für einen von `devcontainer up` erzeugten. Der Ersatz „Grenze der Zusage … Andere Engines, Hosts und Terminals sind nicht zugesagt“ trägt diesen Vorbehalt nicht mehr und nennt kein Bezugsobjekt für „andere“ (Engine-Version? Host? Terminal wovon?). Die Signal- und Shell-Aussagen lesen sich damit für jeden Container des Befehls, die Menge „Engines“ ist gegen die Messmenge nicht mehr bestimmbar. Der Satz „Andere Engines, Hosts und Terminals sind nicht zugesagt“ stand wörtlich schon vorher dort, ist also nicht neu; neu ist sein Wegfall des Bezugs. | `AGENTS.md` §3.10 (Bedingung ohne Test); Auftrag (Bedeutung bewahrt) | `spec/spezifikation.md` · „**Grenze der Zusage** für die Aussagen zum Ende der Sitzung“ | Review-only | Zusage weiter als Beleg (Vorbehalt verloren) |
| F-3 | MEDIUM | „(aus der Form abgeleitet, nicht gefahren)“ ist gestrichen. Er begrenzte die Aussage, ein `docker exec -u <Benutzer> <anderer Container>` innerhalb der Sitzung ende mit `12` statt `0`. Der Satz bleibt als Verhaltensaussage stehen. Gefunden ist ein Einheitstest mit der Zeile eines anderen Benutzers (`devcontainer_shell_test.go`, „status 126, the missing-user line of another user“), kein Test des Falls „innerhalb der Sitzung gegen anderen Container“; das Dokument [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) und das Beobachtungs-Register führen die Aussage als abgeleitet. Der Vorbehalt war hier keine Messsprache, sondern der Beleg-Status der Aussage. | `AGENTS.md` §3.10 (Aussage über Verhalten ohne Test); Auftrag (Bedeutung bewahrt) | `spec/spezifikation.md` · „der Befehl endete dann mit `12` statt `0`“ | Review-only (Suche nach einem pinnenden Test: keiner gefunden) | Zusage weiter als Beleg (Vorbehalt verloren) |
| F-4 | LOW | Verlustfreiheit unvollständig: Abschnitt 3.15 übernimmt nur die vier Sätze aus [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes). Die in der Historie gestrichenen Messaussagen („gemessen in einem Lauf“ zu [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), „als gemessen“/„als nicht geprüft ausgewiesen“ zu [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), „Stabilität des Blob-Hosts ungemessen“ zu [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), „gemessen in einem Lauf“ zu [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) stehen nicht wörtlich dort. Inhaltlich liegen sie teils an anderer Stelle des Dokuments (Zeilen zu „alte Liste gelten“ und zur ersetzbaren Statusdatei); für [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (übergeordneter Name) und [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (Blob-Host) habe ich keine gleichwertige Stelle belegt. Der Abschnittstext behauptet „wörtlich übernommen“ nur für die [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)-Sätze, nennt die Historie aber nicht. | Auftrag (verlustfrei) | `docs/maintainer/zusagen-und-messstaende.md` · „### 3.15 Messbedingungen der Spezifikation“ | Review-only | Streichung ohne Ziel |
| F-5 | LOW | Zeile „[`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) führt die Herkunft je Quelle …“ der Historie ist nach der Streichung sprachlich mehrdeutig: „gepinnte Fälle und der übergeordnete Name bei geteilter IP-Adresse getrennt“ benennt eine andere Trennung als vorher (gepinnte gegen gemessene Fälle). Die Historie bleibt damit eine Änderungsliste, aber mit verschobener Aussage zu dieser Änderung. | Auftrag (Bedeutung bewahrt) | `spec/spezifikation.md` · „gepinnte Fälle und der übergeordnete Name bei geteilter IP-Adresse getrennt“ | Review-only | Bedeutungsverschiebung in Historie |
| F-6 | LOW | Zeile zum Lebenszyklus: „Gemessen ist das mit … (Docker 29.9.2), ein Lauf je Aufruf“ wird zu „Das gilt für … unter Linux `amd64` mit lokaler Engine“. Die Engine-Version entfällt; „mit lokaler Engine“ ist jetzt Zusage-Reichweite jeder Engine-Version, während derselbe Absatz „andere Hosts und Engines“ als nicht zugesagt führt. Das Lastenheft nennt keine Versionen, die Spec verengt also nicht gegen den Vertrag, der Absatz widerspricht sich aber in der Reichweite. Der Test pinnt Docker 29.9.2 (Einzelfall). | `AGENTS.md` §3.10 | `spec/spezifikation.md` · „Das gilt für die gepinnte Version der CLI“ | Review-only | Zusage-Reichweite unbestimmt |
| F-7 | INFO | Zeile zu den Engines: „Zugesagt sind die Engines 29.9.2 und 28.0.4 mit dem Client 29.9.2“ macht aus der Messmenge eine positive Zusage-Menge. Vorher stand die Menge als „Messbedingungen … andere Engine-Versionen sind nicht zugesagt“, was dieselbe Menge implizit zusagte; die Bedeutung ist praktisch gleich. Gegen das Lastenheft ist das eine Verengung ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) kennt keine Engine-Versionen), keine Erweiterung; Tests pinnen beide Engines (Einheitstests und CI-Runner 28.0.4). Zur Beachtung für den Projektinhaber: die Zusage-Menge steht erstmals als solche in der Spec. | Maintainability | `spec/spezifikation.md` · „Zugesagt sind die Engines 29.9.2 und 28.0.4“ | Review-only | Messmenge als Zusage-Menge |
| F-8 | INFO | Offene Messsprache in der Spec, vom Auftrag nicht erfasst, bewertet nicht geändert: `spec/spezifikation.md` Z. 1478 (Historie, beschreibt die Entfernung selbst; unproblematisch); Historie Z. 1486 bis 1492 nicht auf Messwörter geprüft. `spec/lastenheft.md` Z. 416 „Messung“, Z. 1101 „gemessen ohne Docker-Kommunikation“ (Latenzbedingung, Anforderungstext, bleibt), Z. 1588, 1593, 1594, 1596 („gemessene Hosts“, Messwerte in der Historie) und 1607, [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) (Z. ~1497) sowie die „Nicht zugesagt“-Sprache (Z. 474, 507) sind Vertragstext bzw. Historie des Lastenhefts; der Auftrag lässt das Lastenheft unberührt, Änderung des Vertragstexts ist Sache des Projektinhabers. `spec/architecture.md` und `spec/spezifikation.md` enthalten sonst keine Messstände, Colima-/Docker-Versionen als Testumgebung oder Prozesszustände (Treffer nur „geplant“ im Sinn des Produktverhaltens, „Lauf“ als Befehlsausführung, Z. 743 „Geplante Erweiterungen … keine“). | Maintainability | siehe Befund | Review-only | Restbestand Messsprache |

Zu den Prüfpunkten im Einzelnen:

- **(a)/(F-2):** Der Satz ist nicht neu, sein Bezug ist weg. Die Zusage-Menge des Lastenhefts (Linux, `amd64`, lokale Engine; Colima, Docker Desktop, rootless, Podman, entfernte Engine, `arm64`, Windows nicht zugesagt) deckt „andere Hosts“; „andere Engines“ bleibt gegen das Lastenheft ohne Version bestimmt.
- **(b)/(F-6):** Weder enger noch weiter im Vertrag; weiter gegenüber der Messung (Version entfällt).
- **(c)/(F-7):** Praktisch gleiche Menge; aus Messaussage wird Zusage-Wortlaut.
- **(d)/(F-3):** Ja, ein Vorbehalt geht verloren.
- **(e)/(F-5):** Nur Messsprache gestrichen, kein Änderungsinhalt entfernt; die Historie bleibt vollständig, eine Zeile ([`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) ist aber mehrdeutig.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` | `git diff HEAD` leer, geprüft, ohne Befund |
| `spec/architecture.md` Decken-Regel (§3.4) | gelesen: Slice-, Welle-, ADR-, Roadmap-, Commit-Bezüge stehen nur in der Hard Rule (Z. 9 bis 14); die einzige Änderung („damit andere Dokumente auf eine Komponente verweisen können“) ist indikativ, ohne Bezug nach unten; geprüft, ohne Befund |
| `spec/spezifikation.md` Verweise nach unten | gelesen: Kopf Z. 6 bis 9 verbietet sie, im Text kein Link auf ADR, Slice, Carveout, Roadmap; geprüft, ohne Befund |
| §3.5 ADR-Immutabilität | keine ADR berührt, geprüft, ohne Befund |
| §3.7/§3.9 Kommentarmuster im Diff | Muster (a) „früher/former/no longer“ nicht neu eingeführt; der Absatz 3.15 sagt „standen bis zum 2026-10-31 im Abschnitt“ (Chronik, in einem Maintainer-Dokument, dessen Zweck sie ist): geprüft, ohne Befund |
| §3.10 „unter Linux amd64 mit lokaler Engine“ als Eigenschaft des Codes | Zusage-Menge des Lastenhefts ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)), keine Code-Eigenschaft: im Adapter fand ich keine GOARCH-Prüfung; der Wortlaut sagt „gilt für“, nicht „nur wenn“; für das Verhalten selbst gibt es den Test `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters`; geprüft, ohne Befund |
| `make docs-check` | 941 Dateien, 0 Befunde (vor Anlegen dieses Reports) |
| `make ablage-pruefung` | 26 Fälle, 0 Fehlschläge |

Nicht geprüft:

- `make gates`, `make test-docker`, `make verify-depguard` — Gates sind Sache des Verifiers; kein Docker-Pfad und keine Importregel im Diff berührt.
- Wörtlicher Satz-für-Satz-Abgleich der Historie 1486 bis 1492 — außerhalb des Diffs.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Fundstelle trägt die Aussage nicht · Zusage weiter als Beleg (Vorbehalt verloren) · Streichung ohne Ziel · Bedeutungsverschiebung in Historie · Zusage-Reichweite unbestimmt · Messmenge als Zusage-Menge · Restbestand Messsprache

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM (F-1 bis F-3) normalerweise vor dem Merge klären; keine HIGH. Der Diff erfüllt die Decken-Regel und lässt das Lastenheft unberührt; die Bedeutungsverschiebungen betreffen die Vorbehalte „Container ist kein `devcontainer up`-Container“ und „abgeleitet, nicht gefahren“ sowie eine falsche Fundstelle.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen zusätzlich in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
