# Review-Report: spec-prozesszustaende-nachlauf — 2026-10-31

**Review-Art:** Design (Doku-Diff gegen Entscheidungen des Projektinhabers, Hard Rules und Lastenheft)

**Gegenstand:** uncommittete Änderungen (`git diff HEAD`) an `spec/spezifikation.md` (3 Stellen) und `docs/maintainer/zusagen-und-messstaende.md` (Abschnitt 3.15). Sie setzen die drei Entscheidungen des Projektinhabers vom 2026-10-31 um, die der Report `2026-10-31-spec-prozesszustaende.md` (F-2, F-3, F-6) offen ließ.

**Skill:** `.harness/skills/reviewer.md` (Stand des Arbeitsbaums) · Kopie der Vorlage per `cp`, `diff -q` gegen `.harness/baseline/v6.18.1/templates/docs/reviews/review-report.template.md` vor dem Füllen: gleich.
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-31. Der Reviewer hatte kein Edit-Werkzeug; der Report ist mit Write geschrieben.

**Eingangs-Kontext:**

- Report `2026-10-31-spec-prozesszustaende.md` (F-2, F-3, F-6)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge: Docker auf Linux-Host mit lokaler Engine, `amd64`; kein Basis-Image, keine Engine-Versionen)
- [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
- Planungsdokument [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) (Zeile „Messung vor dem Bau“)
- `AGENTS.md` §3.4, §3.7, §3.9, §3.10

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| N-1 | MEDIUM | Die Reichweite „zugesagt für Container mit dem Basis-Image `…base:debian`“ ist weiter als die Tests, die sie tragen: Die Shell-Tests (Signale, Login-Shell) laufen gegen einfache `sleep 3600`-Container dieses Images (`internal/adapter/driven/docker/devcontainer_shell_docker_test.go`, `newShellContainer`), nicht gegen einen von `devcontainer up` erzeugten Container (anderer Entrypoint, Init, Feature-Schicht). Der Vorbehalt „nicht ein von `devcontainer up` erzeugter“ fehlt jetzt; die Zusage-Menge schließt genau den Container des Befehls (erzeugtes Dockerfile mit diesem Basis-Image nach [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)) ein, den kein Test pinnt. Nächster Änderer liest Signal- und Shell-Aussagen als belegt für den `up`-Container. | `AGENTS.md` §3.10 (Bedingung ohne Test) | `spec/spezifikation.md` · „Sie sind zugesagt für Container mit dem Basis-Image“ | Review-only (Test gegen einen `up`-Container nicht gefunden) | Zusage-Menge weiter als Test-Menge |
| N-2 | LOW | Die Reichweite lässt „Engines“ unbestimmt: „Andere Engines, Hosts und Terminals sind nicht zugesagt“ nennt keine zugesagte Engine-Menge; die Engine-Versionen (29.9.2, 28.0.4) stehen nur beim Punkt der Erkennung weiter unten. Für Signale und Shell-Art ist „andere“ gegen nichts bestimmt. Der Satz stand wortgleich schon in HEAD; neu ist, dass die Angabe „Docker 29.9.2“ davor fehlt. | `AGENTS.md` §3.10 | `spec/spezifikation.md` · „Andere Engines, Hosts und Terminals sind nicht zugesagt.“ | Review-only | Zusage-Reichweite unbestimmt |
| N-3 | LOW | Im Lebenszyklus-Absatz steht „andere Hosts und Engines“ zweimal als nicht zugesagt: im neuen Satz („andere Hosts und Engines sind nicht zugesagt“) und in der folgenden Aufzählung („… und andere Hosts und Engines“). Die Liste „Nicht zugesagt sind …“ ist der alte Ort; der neue Satz dupliziert ihn. | Maintainability | `spec/spezifikation.md` · „andere Hosts und Engines sind nicht zugesagt. Nicht zugesagt sind“ | Review-only | Doppelung einer Grenzaussage |
| N-4 | LOW | Die gestrichene Zeile hinterließ die Leerzeile vor „Alles andere gilt als Sitzung: …“ nicht mehr (Diff: `-` Leerzeile). Der eingerückte Absatz folgt jetzt unmittelbar dem verschachtelten Aufzählungspunkt „unable to find user“ und wird im Markdown als dessen Fortsetzung gelesen, nicht als eigener Absatz des übergeordneten Punkts; die Aussage „Alles andere gilt als Sitzung“ hängt damit am letzten Unterpunkt. | Maintainability | `spec/spezifikation.md` · „der Doppelpunkt ist die Namensgrenze, `vscode2` ist nicht `vscode`.“ (Zeilenende, danach „Alles andere gilt als Sitzung“) | Review-only (Rendering) | Zufälliger Strukturbruch durch Streichung |
| N-5 | LOW | Der Einleitungsabsatz von 3.15 sagt weiter, was in der Spezifikation „ausdrücklich als ‚bleibt in der Spec bis zur Entscheidung‘ steht, steht dort noch“. Nach dem Diff steht dort keine solche Stelle mehr (Suche in `spec/spezifikation.md`: keine Treffer); die Aussage über die Spec ist abgelaufen (`AGENTS.md` §3.10 Zustandsmarker mit Ablaufpunkt). | `AGENTS.md` §3.10 | `docs/maintainer/zusagen-und-messstaende.md` · „bleibt in der Spec bis zur Entscheidung“ | Review-only | Abgelaufener Zustandsmarker |
| N-6 | INFO | Die Fundstelle des Basis-Images im Maintainer-Dokument ist ehrlich begrenzt, belegt aber nicht, was sie scheinbar belegt: Die Zeile „Messung vor dem Bau“ des Planungsdokuments nennt `…base:debian` und `busybox:1.36` für die Messung der Engine 28.0.4 (Fall `user_removed`), nicht für Sitzungsende, Signale, Art der Shell; das Maintainer-Dokument sagt das („welche Aussage auf welchem Image lief, nennt die Zeile nicht“). Das Basis-Image als solches stimmt mit dem Planungsdokument und mit [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) überein. | Maintainability | `docs/maintainer/zusagen-und-messstaende.md` · „Das Basis-Image steht im Planungsdokument“ | Review-only | Fundstelle trägt nur das Image |

Zu den Prüfpunkten:

- **(1) F-2.** Das Lastenheft nennt kein Basis-Image und sagt nichts zu Signalen oder Art der Shell; die Reichweite verengt also nichts, was das Lastenheft zusagt, und widerspricht keinem Lastenhefttext (Spec darf präzisieren, eine Verengung einer nicht im Lastenheft stehenden Aussage ist zulässig). Das Image `mcr.microsoft.com/devcontainers/base:debian` stimmt mit dem Planungsdokument und [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) überein. Sie steht als Zusage-Menge, nicht als Test-Pin; die Tests decken nur den einfachen Container (N-1). Der Satz „Andere Engines, Hosts und Terminals sind nicht zugesagt“ bleibt stimmig zu den Hosts (Zusage-Menge des Lastenhefts), für „Engines“ unbestimmt (N-2).
- **(2) F-3.** Keine andere Stelle der Spezifikation oder des Lastenhefts hängt am gestrichenen Satz (Suche „anderer Container“, „innerhalb der Sitzung“: kein Treffer; die Fehlerklasse `12` für Fehler der Sitzung kommt in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nur für Fehler vor Beginn der Shell vor, danach `0`). Die Streichung ist im Maintainer-Dokument wörtlich festgehalten (Vergleich mit `git show HEAD:spec/spezifikation.md`: Satz vollständig, auch „(aus der Form abgeleitet, nicht gefahren)“) und als unbelegte Beobachtung ausgewiesen. Folge: Die Spec ist zur Mehrdeutigkeit (Zeile des fehlenden Benutzers aus einer Sitzung heraus) jetzt still; das ist der Gehalt der Entscheidung, kein Verlust.
- **(3) F-6.** „unter Linux `amd64` mit lokaler Engine“ ist die Zusage-Menge des Lastenhefts („unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`)“); wörtlich gleich in der Substanz. „Gepinnte Version der CLI“ und „Befehle in der Form einer Zeichenkette, die gelingen“ stehen nicht im Lastenheft; sie sind Verengungen (nicht Erweiterungen) und standen unverändert vorher in der Spec. Die Engine-Version „Docker 29.9.2“ entfällt; das Lastenheft kennt keine, die Aussage ist damit für jede lokale Engine zugesagt, gepinnt durch einen Test nur für die Engine des Test-Images (Einzelfall); Reichweite dem Lastenheft gemäß, vom Test enger.
- **(4) Maintainer 3.15.** Wortlaut vorher (HEAD) steht für alle drei Stellen vollständig dort, das Abschnittsnummernschema bleibt; die Verweise (Planungsdokument, Messprotokoll, Test, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)) lösen auf (`make docs-check` grün). Die Fundstellen stimmen mit dem Fix des letzten Reviews. Doppelung: N-3 (Spec), N-5 (Maintainer-Einleitung).
- **(5) Decken-Regel.** Im Diff von `spec/spezifikation.md` kein neuer Verweis auf ADR, Slice, Carveout, Roadmap; die neuen Sätze sind indikativ. Geprüft, ohne Befund.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/spezifikation.md` Decken-Regel (§3.4) | drei geänderte Stellen gelesen, keine Verweise nach unten, geprüft, ohne Befund |
| `spec/lastenheft.md` | `git diff HEAD` leer, geprüft, ohne Befund |
| `AGENTS.md` §3.7/§3.9 Muster (a) im Spec-Diff | kein „früher/former/no longer/ohne X wäre“ neu in der Spec; Chronik steht im Maintainer-Dokument, wo sie hingehört; geprüft, ohne Befund |
| Streichung „ein Lauf je Aufruf“ / „Docker 29.9.2“ im Lebenszyklus-Absatz | im Maintainer-Dokument im Wortlaut vorher enthalten, geprüft, ohne Befund |
| `make docs-check` | 942 Dateien, 0 Befunde (vor Anlegen dieses Reports) |
| `make ablage-pruefung` | 26 Fälle, 0 Fehlschläge |

Nicht geprüft:

- `make gates`, `make test-docker`, `make verify-depguard` — Gates sind Sache des Verifiers; kein Docker-Pfad und keine Importregel im Diff berührt, keine Container gestartet.
- Ob ein Test die Signal- und Shell-Aussagen gegen einen von `devcontainer up` erzeugten Container fährt — kein solcher Test gefunden, nicht Zeile für Zeile aller Docker-Tests gelesen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zusage-Menge weiter als Test-Menge · Zusage-Reichweite unbestimmt · Doppelung einer Grenzaussage · Zufälliger Strukturbruch durch Streichung · Abgelaufener Zustandsmarker · Fundstelle trägt nur das Image

## Verdikt

**Merge-blockierend:** nein für HIGH (keine); N-1 (MEDIUM) normalerweise vor dem Merge klären. Keine der drei Stellen macht die Bedeutung weiter als das Lastenheft: Alle Verengungen (Basis-Image, gepinnte CLI, Zeichenkettenform) liegen innerhalb der Zusage-Menge von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten). Weiter als der Test ist die Reichweite des Basis-Images (N-1) und die Engine-Menge (N-2, Lebenszyklus-Absatz).

**Übergabe:** Findings gehen an den Implementer. Der Report ersetzt keine Verifikation.
