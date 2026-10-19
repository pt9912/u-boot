# Review-Report: slice-v1-e2e-pin-mismatch-zwischencontainer — 2026-10-19

**Review-Art:** Code-Review (Docker-Test-Hygiene: Aufräumen des Zwischencontainers der Negativprobe zur Ollama-Prüfsumme) gegen Plan, Entscheidungen und Hard Rules; Test, Makefile-Fragment, Dockerfile-Stage und Workflow selbst gelesen, nicht aus dem Bericht des Implementers übernommen.

**Gegenstand:** `git diff -M 5eab741..HEAD` (Commits `1689c68` Code, `be66a0a` Slice-Stand): 2 Dateien, 75 Einfügungen, 6 Löschungen; keine neue Datei.

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

**Eingangs-Kontext:**

- Slice [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../plan/planning/in-progress/slice-v1-e2e-pin-mismatch-zwischencontainer.md) (§1, §2 DoD, §3 Plan, §4, §6, Abschnitt „Umsetzung (Implementer)“).
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md); `AGENTS.md` §3.1, §3.2, §3.7, §3.9, §3.10; `harness/review.md`; `spec/architecture.md` Abschnitt Tests.
- Selbst gelesen: `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`, `containerIDs`, `newContainersWithCommand`, `dockerOutput` und `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` (`internal/e2e/sandbox_devcontainer_docker_test.go`); Ziel `test-docker` (`repo.mk`) mit `docker run … go test -tags docker`; Stage `test-docker-tools` (`Dockerfile`); `.github/workflows/integration.yml`; Fundstellen von `t.Parallel` im Paket (nur `sandbox_status_docker_test.go`).
- Sensor-Läufe des Reviewers (nur `make`): `make lint` (Exit 0), `make test` (Exit 0), `make docs-check` (515 Dateien, 0 Befunde), `make doc-trace` (83 Anforderungen, 0 Waisen), `make doc-planning` (0 Befunde), `make ablage-pruefung` (26 Fälle, 0 Fehlschläge). `make test-docker` nicht gefahren (Verifier-Aufgabe, kein Befund verlangt es).

---

## Findings

Keine HIGH, keine MEDIUM; der Konflikt-Pfad ist nicht ausgelöst. Zwei LOW, vier INFO.

LOW `internal/e2e/sandbox_devcontainer_docker_test.go` · "newContainersWithCommand(context.Background(), t, before, pin)" - Zuordnung „neu und Null-Pin“ trennt nicht den eigenen vom gleichartigen fremden Lauf
Klasse: Zuordnung nicht eigen-eindeutig
Quelle: Slice §2 DoD („Der Eingriff berührt keinen fremden Container“), Slice §6 Risiko 1
Befund: Der Pin (`OLLAMA_SHA256_AMD64=` plus 64 Nullen) ist für jeden Lauf dieses Tests derselbe. Läuft derselbe Test ein zweites Mal gleichzeitig gegen denselben Daemon (zweiter Checkout, zweiter Entwickler, selbst gehosteter Runner), ist dessen Zwischencontainer in der Vorher-Liste des ersten Laufs nicht enthalten, trägt den Pin und wird vom `t.Cleanup` des ersten Laufs mit `docker rm -f` entfernt, auch mitten im Schritt des anderen Baus. Innerhalb eines Laufs ist die Zuordnung eindeutig: Der Test ist nicht parallel (`t.Parallel` steht nur in `sandbox_status_docker_test.go`, dort ohne Null-Pin), und der Podman-Test des Pakets nutzt andere Variablen.
Risiko: Der andere Lauf verliert seinen Container, sein Ausgang (Text, ob `FAILED` erscheint) wird nicht mehr durch seinen eigenen Fehlschlag bestimmt; die DoD-Zusage „kein fremder Container“ gilt dann nur für fremde Objekte anderer Art. Die Wahrscheinlichkeit ist klein (ein geteilter Daemon und ein gleichzeitiger Lauf desselben Pakets); die Zusage in DoD und §6 ist trotzdem weiter als der Code.
Verifizierbar: Review-only (Lesen von `newContainersWithCommand` und der Pin-Konstruktion)

LOW `docs/plan/planning/in-progress/slice-v1-e2e-pin-mismatch-zwischencontainer.md` · "dessen `Path`/`Args` den Pin `OLLAMA_SHA256_AMD64` tragen" - Plan-Aussage über den Container bleibt nach der Messung stehen <!-- d-check:ignore (Dateipfad des Slice als Fundstelle) -->
Klasse: Plan-Aussage über Code nicht berichtigt
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10 (Abgleich nach der letzten Fix-Runde; Aussage nur, wenn gelesen)
Befund: §1 („Gefundener Stand“), die DoD (erste Zeile: „mit dem Pin im `Path`/`Args`“) und §6 Risiko 1 („mit dem Pin im Befehl“) führen Path/Args als Träger des Pins. Der Abschnitt „Umsetzung (Implementer)“ stellt fest, dass Path/Args nur den Variablennamen tragen und der Wert in `Config.Env` steht; die Prüfung liest deshalb das Env. Die genannten Stellen bleiben unberichtigt. Das als Beleg zitierte Rohdatum `rohdaten-09.txt` enthält nach einer Textsuche keine Zeile mit `Path`, `Args` oder `Env`; die Fundstelle der Aussage ist dort nicht auflösbar.
Risiko: Der nächste Änderer (etwa beim Folge-Vorgang für den zweiten Test) liest die DoD-Zeile, prüft Path/Args und erzeugt die zwei grünen Zwischenläufe, die der Implementer verworfen hat; kein Gate widerlegt die Plan-Aussage.
Verifizierbar: Review-only

INFO `docs/plan/planning/in-progress/slice-v1-e2e-pin-mismatch-zwischencontainer.md` · "Zweites Vorkommen (nur gemeldet, nicht behoben)" - Auslegung der DoD-Bedingung und Empfehlung <!-- d-check:ignore (Dateipfad des Slice als Fundstelle) -->
Klasse: Folge-Vorgang ohne Anker
Quelle: Slice §1 NICHT-Punkt 4 („erst, wenn die Prüfung in der DoD ein zweites Vorkommen zeigt“), §4 Rückführung `in-progress` → `next`
Befund: Gelesen: `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` ruft in drei Teilfällen `docker build -q` ohne `--force-rm` und ohne Cleanup auf und erwartet Exit-Code ungleich 0; derselbe klassische Builder, also derselbe Mechanismus (nicht gemessen, vom Implementer ebenso ausgewiesen). Die Bedingung des Plans verlangt, dass die Prüfung der DoD das zweite Vorkommen zeigt; die Prüfung (Pin des Eingriffs) ist auf den Ollama-Test zugeschnitten und zeigt keines, das zweite Vorkommen entsteht aus Lesen. Das ist vom Wortlaut gedeckt, wenn „zeigt“ als „gemessen“ gelesen wird, und nicht gedeckt, wenn Lesen genügt. Der Plan enthält die Rückführung für den Fall, dass „das Aufräumen eine Änderung an mehreren Docker-Tests verlangt“, mit Schnitt durch den Architect; das spricht für die Auslegung „Folge-Vorgang“. Empfehlung: Folge-Vorgang, nicht jetzt. Gründe: (a) drei Teilfälle haben verschiedene Pins (`PODMAN_SHA256_AMD64`, `PODMAN_SIGNING_KEY` mit zwei Werten) und scheitern an verschiedenen Schritten, der Pin der Zuordnung ist damit je Teilfall zu wählen und der Mechanismus je Teilfall zu messen; (b) der Slice hat seine Rückführung dafür vorab benannt; (c) `newContainersWithCommand` nimmt das Suchwort bereits als Parameter, die spätere Verallgemeinerung ist klein.
Risiko: Der Folge-Vorgang hat im Diff keinen auflösbaren Anker (kein Register-Eintrag, kein Slice in `open/`); die DoD verlangt vor der Closure das Fortschreiben des Beobachtungs-Registers oder die Notiz in §7, und §7 ist noch leer. Ohne Anker bleibt die Lücke eine Aussage im Fließtext des Slice.
Verifizierbar: Review-only; Zuständigkeit Verifier (DoD-Punkt „Register fortgeschrieben“) und Architect (Schnitt des Folge-Vorgangs)

INFO `internal/e2e/sandbox_devcontainer_docker_test.go` · "out, code := dockerOutput(ctx, t, \"build\", \"-q\", \"--force-rm\", \"-t\", tag," - Mechanismus: Ursache plausibel, drei Teile ungeprüft
Klasse: Mechanismus-Aussage teils ungemessen
Quelle: Slice §6 Risiko 2, Slice „Nicht gemessen“
Befund: Bestätigt durch Lesen: Die Stage `test-docker-tools` installiert nur `docker-ce-cli` und `docker-compose-plugin`, kein buildx-Plugin; `repo.mk` führt den Test mit `docker run … $(IMAGE):test-docker-tools go test -tags docker` aus, also mit dem Client des Images. Ein Docker-CLI ab 23 fällt ohne buildx auf den klassischen Builder zurück; dieser behält den Container eines fehlgeschlagenen Schritts, solange `--force-rm` fehlt (Docker-Verhalten, so weit ohne Messlauf entscheidbar). Die Messung des Implementers (Wegwerf-Probe, Container-Zahl, Befehl `docker buildx version`) habe ich nicht wiederholt. Die CI setzt `DOCKER_BUILDKIT: '1'` auf Job-Ebene (`integration.yml`), aber `docker run` im Ziel reicht keine Umgebung in den Test-Container durch (kein `-e`, kein `--env`); der Test-Client ist dort derselbe wie lokal, der Daemon des Runners ist nicht gemessen. Ungeprüft bleibt, wie `--force-rm` auf einem Client mit buildx oder einem Daemon mit BuildKit-Standard wirkt (Flag wirkungslos, ignoriert mit Hinweis oder abgelehnt): wirkungslos oder ignoriert lässt den Test grün, weil dann kein Container entsteht; abgelehnt würde den Test laut rot machen (Exit ungleich 0, aber ohne `FAILED`), nicht still. Die Erwartung des Tests (Exit ungleich 0, `FAILED` in der Ausgabe) ändert `--force-rm` im klassischen Builder nicht: Es entfernt nur den Container nach dem Schritt, die Ausgabe des Schritts erscheint vorher.
Risiko: Gering; ein Wechsel des Test-Images auf buildx würde die Prüfung vakuum-grün machen (kein Container entsteht), ohne dass ein Test es sagt.
Verifizierbar: `make test-docker` (Verifier); der Wechsel-Fall ist Review-only

INFO `docs/plan/planning/in-progress/slice-v1-e2e-pin-mismatch-zwischencontainer.md` · "Ein fremder Lauf (Image" - Beleg-Tabelle ehrlich, ein Objektname eines fremden Laufs in der Datei <!-- d-check:ignore (Dateipfad des Slice als Fundstelle) -->
Klasse: Beleg-Lesung
Quelle: Slice §2 DoD (Pin „einmal rot gesehen“); Auftragsvorgabe „keine Namen fremder Docker-Objekte in Dateien“
Befund: Die Tabelle liest sich ehrlich: Erste Zeile (Eingriff entfernt): 3 → 3, weil der Cleanup den einen Container nach der Meldung „left 1 container(s)“ entfernt hat, der Test ist rot; zweite Zeile (geliefert): 3 → 4, `PASS`. Die Differenz der Gesamtzahl in der zweiten Zeile wird mit einem fremden Lauf erklärt und ausdrücklich nicht als Beleg gewertet; der Beleg ist die Prüfung nach Pin und neuer ID (`PASS` heißt: kein neuer Container mit dem Null-Pin). Die Aussage „kein Container im Status `exited` auf dem Host“ in derselben Zeile ist eine Hostaussage und trägt den Test nicht. Die Mutation ist als gemessen ausgewiesen; ich habe sie nicht wiederholt. In der Datei steht das Basis-Image des fremden Laufs mit Namen und Tag; es ist ein generisches öffentliches Basis-Image, kein Objekt des Projekts, aber ein Name aus einem fremden Lauf.
Risiko: Keines für den Diff; die Namensnennung ist durch die Auftragsvorgabe berührt.
Verifizierbar: Review-only

INFO `internal/e2e/sandbox_devcontainer_docker_test.go` · "left := newContainersWithCommand(" - Altbestand und Fehlerpfade des Cleanups
Klasse: Grenze der Zusage
Quelle: Maintainability
Befund: Die Vorher-Liste schließt Container früherer Läufe aus; ein Zwischencontainer, den ein älterer Lauf (vor dieser Änderung) hinterlassen hat, bleibt auf dem Host und wird nicht entfernt (gewollt; das Aufräum-Skript des Messprotokolls ist ausdrücklich nicht Gegenstand). Der Rückgabewert von `docker rm -f` wird verworfen; scheitert das Entfernen, meldet der Test dennoch denselben Text wie bei gelungenem Entfernen („left %d container(s)“). `containerIDs` ruft `t.Fatalf` aus dem Cleanup, wenn `docker ps` scheitert; das ist erlaubt und lärmt, nicht still.
Risiko: Gering.
Verifizierbar: Review-only

---

## Negativbefunde

Geprüft, ohne Befund:

- **Zuordnung (False Positive/Negative innerhalb eines Laufs):** Entfernt wird nur, was nicht in der Vorher-Liste steht **und** im Format von `{{.Path}} {{.Args}} {{.Config.Env}}` den Needle `OLLAMA_SHA256_AMD64=` plus 64 Nullen enthält; kein Status, kein Zeitfenster, kein Name. Fremde Container mit anderem Inhalt und Container des Podman-Tests (andere Variablen) werden nicht erkannt; erfolgreiche Zwischencontainer des klassischen Builders entfernt der Builder selbst. Unabhängige Pakete laufen im Test-Aufruf (`./...`) nebeneinander, enthalten aber keine Builds mit diesem Pin. Das Rest-Risiko steht im ersten LOW.
- **Test-Reihenfolge und Cleanup-Reihenfolge:** Der Cleanup wird vor der Prüfung des Exit-Codes registriert, läuft also auch im `t.Fatalf`-Pfad; der `rmi`-Cleanup im Erfolgspfad ist später registriert und läuft vorher (LIFO). `context.Background()` im Cleanup ist gewollt, `ctx` ist dann abgelaufen oder abgebrochen.
- **Erwartungswert des Tests:** Exit ungleich 0 und die Ausgabe `FAILED` hängen nicht an `--force-rm`; im klassischen Builder erscheint die Ausgabe des Schritts vor dem Entfernen des Containers.
- **Test-Hygiene:** Build-Tag `//go:build docker` steht in der Datei; keine neue Abhängigkeit (nur `strings`, bereits importiert); keine Uhr im Spiel, die Fake-Clock-Konvention (`spec/architecture.md`, Abschnitt Tests) ist nicht berührt; kein `//nolint`; `make lint` grün. Fehlermeldungen des Tests sind Englisch. Das Test-Doc nennt die ADR-Kennung der Negativprobe ([`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md)); die Prüfsumme ist über [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) im Slice verankert.
- **Kommentare (`AGENTS.md` §3.7, §3.9):** Die drei neuen Kommentare sind im Indikativ und tragen Zusage oder Kopplung („the `test-docker-tools` image has no buildx“ koppelt an die Stage, die dies bestimmt); keine Chronik, kein „früher“; der Symbolverweis `--force-rm` und der Stage-Name stimmen mit dem Baum überein (Suche im ganzen Baum: `--force-rm` nur im Test und im Slice). Die Helfer `containerIDs` und `newContainersWithCommand` stehen **nach** der Testfunktion, deren Doc-Kommentar davor steht (kein Block-Bruch). Das Doc des Tests nennt das Aufräumen nicht; das ist keine falsche Aussage.
- **Plan-Tabelle gegen `git diff --stat -M 5eab741..HEAD` (`AGENTS.md` §3.10):** Der Diff berührt genau zwei Dateien, den Test (Zeile 1 der Tabelle) und die Slice-Datei (Sammelzeile). Unverändert und im Diff nicht enthalten: Produktcode (`internal/hexagon/**`, `internal/adapter/**`), `Makefile`, `repo.mk`, `.github/workflows/`, `spec/`, Templates, ADRs, `docs/maintainer/quality.md` (`git diff --stat` für `repo.mk`, `.github`, `spec`, `Makefile` leer). Die als „entfällt“ geführten Messungs-Pfade (`docs/messungen/<slice>/`, `tools/messung/<name>/`) fehlen im Diff und sind als Text im Slice belegt. Die Zeile zu `quality.md` nennt eine Suche und Zeilen der Treffer; die Aussage ist hier nicht neu gelesen.
- **Gate-Lockerung, Suppression, Docker-only:** keine Änderung an Gates, Schwellen, `.golangci.yml`, Workflows; der Docker-Pfad bleibt `make test-docker`; die Messungen des Implementers laufen im Image des Ziels.
- **Dual-Classifier, Sprache, Spec-Straten, Dateisicherheit, Architektur:** nicht berührt (kein Produktcode, keine Spec, kein CLI-Pfad).
- **Doku/Traceability:** Die Slice-Änderungen tragen verlinkte Kennungen; `make docs-check`, `make doc-trace`, `make doc-planning` grün; der Slice-Stand ist Zustand mit Beleg (DoD: zwei Punkte `[x]`, der dritte offen).

Nicht geprüft:

- `make test-docker` und die Wirkung von `--force-rm` im Daemon des CI-Runners sowie auf einem Client mit buildx: Verifier-Aufgabe beziehungsweise vom Slice als nicht gemessen ausgewiesen (§6 Risiko 2).
- Die Wiederholung der roten Mutation und der Wegwerf-Probe des Implementers: gelesen, nicht gefahren.
- `make gates`, `make verify-depguard`: Importregeln und Produktcode unberührt; die Einzelsensoren `lint`, `test`, `docs-check` liefen grün.
- DoD-Closure, Closure-Notiz, Register-Fortschreibung, Ausgang der zwei Risiken in §6: Verifier-Aufgabe, der dritte DoD-Punkt ist offen.

## Zusammenfassung

Keine HIGH, keine MEDIUM; zwei LOW, vier INFO. Der Eingriff ist klein, die Ursache (klassischer Builder im Image ohne buildx, ohne `--force-rm` bleibt der Container des fehlgeschlagenen Schritts) ist nach dem Lesen von `repo.mk`, Dockerfile-Stage und Workflow plausibel; die Zuordnung nach neuer ID plus Null-Pin ist innerhalb eines Laufs eindeutig und berührt keine fremden Objekte. Offen für den Implementer: die Zusage „kein fremder Container“ gilt nicht für einen gleichzeitigen Lauf desselben Tests auf demselben Daemon (erstes LOW), und die Plan-Aussagen zu Path/Args sind nach der Messung zu berichtigen (zweites LOW). Zum zweiten Vorkommen: Empfehlung Folge-Vorgang, mit einem Anker (Register-Eintrag oder Slice in `open/`) vor der Closure.
