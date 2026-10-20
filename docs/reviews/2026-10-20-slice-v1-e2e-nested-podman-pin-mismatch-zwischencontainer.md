# Review-Report: slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer — 2026-10-20

**Review-Art:** Code-Review eines Docker-Tests (Test-Hygiene, kein Produktcode) mit Prüfung der Slice-Aussagen im Abschnitt Umsetzung; Mutationen und Rohdaten selbst nachvollzogen, nichts vom Implementer übernommen.

**Gegenstand:** `git show 034e124` (zwei Dateien: `internal/e2e/sandbox_devcontainer_docker_test.go`, [`docs/plan/planning/done/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md`](../plan/planning/done/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md)); Vorgänger [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../plan/planning/done/slice-v1-e2e-pin-mismatch-zwischencontainer.md).

**Skill:** `.harness/skills/reviewer.md` @ HEAD `034e124` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

**Eingangs-Kontext:**

- Slice (Ziel, DoD, §3 Plan-Tabelle, §6 Risiken, Abschnitt Umsetzung); [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); [`ADR-0023`](../plan/adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md); `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §3.9, §3.10; `harness/review.md`.
- Selbst gelesen: `NestedPodman_PinMismatchBreaksBuild` und `Ollama_PinMismatchBreaksBuild` im Wortlaut, die Helfer `containerIDs`, `newContainersWithCommand`, `dockerOutput`, `Dockerfile.tmpl` (ARG-Pins), `repo.mk` (`test-docker`).
- Sensor-Läufe des Reviewers: `make gates` (Exit 0); der berührte Test gezielt im Image `u-boot:test-docker-tools` (Aufruf wie in `make test-docker`, mit `-run`), Mutationen in Kopien des Baums unter dem Scratchpad (kein Repo-Eingriff). `make verify-depguard` nicht gelaufen: keine Importregel berührt. Voller `make test-docker` nicht wiederholt (siehe Nicht geprüft).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Kommentar vor `failedStepRemoved.MatchString` beginnt mit „Without the classic builder no intermediate container exists and the leftover check above finds nothing whatever the flag does.“ Die Aussage ist im Indikativ und beschreibt eine geltende Kopplung (Rest-Prüfung hängt am klassischen Builder), hat aber die Form „Ohne X …“, die `AGENTS.md` §3.9 (a) als Suchmuster nennt. Die zweite Hälfte („earlier steps print the same words when the layer cache is cold“) ist eine Abgrenzung und stimmt (Messung unten, 15 Treffer „Removed intermediate container“ bei `--no-cache`). | AGENTS.md §3.7, §3.9 (a) | `internal/e2e/sandbox_devcontainer_docker_test.go` · „Without the classic builder no intermediate container exists“ | Review-only | kommentar-ohne-x-form |
| F-2 | INFO | Das Geschwister `Ollama_PinMismatchBreaksBuild` baut weiter mit `-q` und trägt keinen Wächter gegen den vakuum-grünen Lauf. Wechselt das Test-Image auf BuildKit, bleibt dieser Test grün, der Nested-Podman-Test wird rot. Der Slice schließt das Ollama-Gegenstück aus (§1) und hält das Risiko in §6 offen; das Ausgangs-Feld des Risikos verlangt für den Slice eine Prüfung oder eine Beobachtung im Register — für den Nested-Test ist die Prüfung geliefert, für die Ollama-Hälfte steht der Ausgang noch aus. Zuständig: Closure/Verifier. | `BEO-E2E/docker-test-hinterlaesst-zwischencontainer`, Slice §6 | `internal/e2e/sandbox_devcontainer_docker_test.go` · „\"build\", \"-q\", \"--force-rm\", \"-t\", tag,“ | Review-only | wächter-nur-an-einer-hälfte |
| F-3 | INFO | Zwei gleichzeitige Läufe gegen denselben Daemon: Der Eingriff sieht den laufenden Build-Container des anderen Laufs als „neu mit Pin“ und entfernt ihn mit `rm -f`; der andere Bau scheitert dann an einem Fremdgrund, nicht an der Pin-Prüfung. Slice §1/§6, Doc-Kommentar von `newContainersWithCommand` und „Nicht gemessen“ nennen die Grenze. Innerhalb des Pakets läuft der Test nicht neben anderen: die `t.Parallel`-Tests (`sandbox_status_docker_test.go`) starten erst nach den seriellen Tests des Pakets. Die Grenze gilt nur für zwei Läufe gegen einen Daemon. | Slice §6 | `internal/e2e/sandbox_devcontainer_docker_test.go` · „Two runs against the same daemon at the same time are not told apart“ | Review-only | — |
| F-4 | INFO | Slice §1 „Gefundener Stand (gelesen am 2026-10-19)“ nennt „ab Zeile 572“. Am 2026-10-19 stimmte das (Stand bei Anlage des Slice, Commit `8c30937`); im Baum vor dem Diff steht die Funktion bei Zeile 865 (Commits des Egress-Slice dazwischen). Die Angabe ist datiert und damit kein Verstoß gegen §3.10; bei der Closure kann der Marker entfallen. | AGENTS.md §3.10 (Zustandsmarker mit Ablaufpunkt) | Slice · „(ab Zeile 572)“ | Review-only | — |
| F-5 | INFO | Das Wächter-Muster hängt am Wortlaut des klassischen Builders. Dessen Ausgabe trägt „DEPRECATED: The legacy builder is deprecated and will be removed in a future release“. Verschwindet der Builder oder ändert sich der Wortlaut, wird der Test rot mit einer Meldung, die beide Ursachen nennt (Builder, `--force-rm`); das ist die gewollte Richtung, aber ein Wortlaut-Wechsel ist von „Mechanismus nicht mehr belegt“ am Ausgang nicht zu unterscheiden. | Maintainability | `internal/e2e/sandbox_devcontainer_docker_test.go` · „builder is not the classic one, or --force-rm is missing“ | `make test-docker` | — |
| F-6 | INFO | Der Cleanup-Block (Vorher-Liste, `newContainersWithCommand`, `rm -f`, `t.Errorf`) steht jetzt zweimal (Ollama, Nested Podman) fast wortgleich. Der Slice verschiebt eine gemeinsame Aufräum-Hilfe ausdrücklich auf ein drittes Vorkommen (§1, Nicht-Ziele); das ist die erklärte Zählung des Beobachtungs-Registers, kein Mangel. | Slice §1 | `internal/e2e/sandbox_devcontainer_docker_test.go` · „the failed build left container %.12s with the pin %s“ | Review-only | duplikat-aufräumblock |

Keine HIGH-, keine MEDIUM-Findings.

## Prüfung der Schwerpunkte

**(1) Wächter `failedStepRemoved`.** Selbst gemessen im Image `u-boot:test-docker-tools` (Docker 29.9.1, klassischer Builder, Hinweis „Install the buildx component“ in der Ausgabe):

| Lauf | Ergebnis |
|---|---|
| geliefert | PASS, 3 Teilfälle (30,9 s) |
| `--force-rm` entfernt | FAIL in allen drei Teilfällen, beide Prüfungen feuern („does not show the failed step's container removed“ und „left container … cleanup removed it“); Ausgabe dort: „sha256sum: WARNING …“ direkt vor „The command '…“, **keine** Zeile „Removed intermediate container“ dazwischen |
| `--force-rm` entfernt, `--no-cache` (kalter Cache) | FAIL in allen drei Teilfällen; die Ausgabe enthält 15 Zeilen „Removed intermediate container“ (erfolgreiche Zwischenschritte wie `USER root`), das Muster mit Nachbarzeile bleibt trotzdem rot — die Nachbarzeile trägt die Aussage |
| `--force-rm` und `--no-cache` | PASS (kalter Cache stört den Wächter nicht) |
| Muster verfälscht („The command ZZZ“), `--force-rm` bleibt | FAIL in allen drei Teilfällen, die Rest-Prüfung schweigt (unabhängig) |
| BuildKit | `docker build --force-rm --progress=plain` auf dem Host (dort buildx vorhanden) mit einem scheiternden `RUN`: 0 Zeilen „Removed intermediate“; der Wächter fände keinen Treffer, der Test würde rot, nicht vakuum-grün. Das Test-Image selbst wurde nicht auf BuildKit umgestellt (Gate-Änderung, Slice-Nicht-Ziel). |

Warmer Cache: Der fehlgeschlagene Schritt wird nie aus dem Cache bedient, der Container entsteht in jedem Lauf (Step 7/12 lief in allen Läufen als „Running in …“). Parallelität im Build: der Wortlaut „Removed … \n The command“ ist je Bau eine zusammenhängende Strecke eines Prozesses (`CombinedOutput`, eine Pipe, Schreibreihenfolge erhalten). Die Aussage „klassischer Builder legte einen Container des fehlgeschlagenen Schritts an, `--force-rm` entfernte ihn“ trägt das Muster damit; es belegt nicht, *welcher* Container entfernt wurde — das tut die Rest-Prüfung mit dem Pin (die in der Mutation ohne `--force-rm` nachweislich die Container fand).

**(2) Cleanup und Zuordnung.** Entfernt wird nur, was nach der Vorher-Liste neu entstand und den Pin des Teilfalls in `Path`/`Args`/`Env` trägt; die Nullen-Pins entstehen sonst nirgends im Paket (Suche nach `PODMAN_SHA256_AMD64=`/`PODMAN_SIGNING_KEY=` in `internal/`: nur Template und Template-Test, der Test setzt keine Werte aus `docker build`). Teilfälle laufen seriell, jeder hat eigene Vorher-Liste und eigenes `t.Cleanup` (am Subtest-`t`). Reihenfolge LIFO: der `rmi`-Cleanup des unerwartet gelungenen Baus (später registriert) läuft vor der Rest-Prüfung. Fehlerpfade: Docker nicht erreichbar → `requireDocker` und `t.Fatalf` in `containerIDs` vor Registrierung des Cleanups; schlägt `rm -f` fehl, meldet es `t.Errorf` mit Ausgabe. Mutationslauf: nach den Läufen war keiner der von den Mutationen hinterlassenen Container mehr vorhanden (per ID geprüft), also hat der Cleanup sie entfernt. Fremde Container wurden nicht angefasst; die Container-Liste des Hosts änderte sich während des Reviews durch Vorgänge außerhalb dieser Läufe, deren Objekte hier nicht aufgeführt werden.

**(3) Fehlendes `-q`.** Die Ausgabe wird nur bei Fehler in die Testmeldung übernommen (laut `-v`-Lauf: 1,3 KB Log im Erfolgsfall; je fehlgeschlagenem Teilfall rund 8 KB bei Fehler, Gesamtlog 25 KB). Inhalt: apt-Ausgabe, Pins (öffentliche Prüfsummen und Fingerabdrücke), kein Token — `sandboxCanaryToken` geht in andere Tests, nicht in diesen Bau. Kein Befund.

**(4) Mutationen** siehe Tabelle oben; alle in Kopien unter dem Scratchpad, Repo unverändert. Docker-Hygiene: nur Container der Testläufe entstanden und wurden vom Test entfernt; erzeugte Images trugen das Test-Präfix bzw. wurden nicht getaggt (Sonde ohne Tag, scheiterte, hinterließ nichts außer Build-Cache).

**(6) Slice-Aussagen, gelesen/gemessen.** „Alle drei Teilfälle scheitern im selben `RUN`-Schritt (Step 7/12)“: bestätigt, dasselbe `RUN apt-get …`. Exit-Codes 1/1/22: bestätigt (`returned a non-zero code: 1` zweimal, `22` einmal mit `curl: (22) … 404`). „Mit `--force-rm` bleibt keiner zurück“: bestätigt durch den Rest-Check in der gelieferten Fassung. „Wirkung bei BuildKit aus der Mutation abgeleitet, nicht gesehen“: der Slice sagt das ehrlich; ich habe die Wortlaut-Aussage auf dem Host gegengeprüft (0 Treffer). Plan-Tabelle gegen `git diff --stat 034e124~1 034e124`: zwei Dateien (Test, Slice) — beide in §3 (Zeile 1 und Sammelzeile); `quality.md`, `docs/messungen/`, `tools/messung/`, Produktcode, `Makefile`, `repo.mk`, Workflows, `spec/` ohne Diff, wie ausgewiesen. Die DoD-Punkte 1 und 2 sind gehakt und durch die Mutationen gedeckt; Punkt 3 steht offen (Review, Closure) — das ist Sache von Verifier und Closure, nicht dieses Reports.

## Negativbefunde

Geprüft, ohne Befund:

- Hard Rules: Docker-only (alle Läufe über das Test-Image, `make gates` Exit 0), kein `//nolint`, keine Lockerung von Gate, Schwelle, Importregel (§3.1, §3.2, §3.6); kein Produktcode, kein Accepted-ADR berührt, kein Carveout.
- Kommentare §3.7/§3.9 außer F-1: (a) kein Kommentar beschreibt Entferntes oder Früheres; (b) die Suche nach `-q`, `force-rm`, `newContainersWithCommand` im Baum findet nur Test, Slice und Vorgänger-Dokumente (datierte Lesung); der Doc-Kommentar „Two runs against the same daemon …“ passt zum Ist-Stand; (c) `failedStepRemoved` steht nach der Probe-Konstante und vor dem Doc-Kommentar der Testfunktion, nicht zwischen Kommentar und Funktion.
- Helfer `containerIDs`, `newContainersWithCommand`: unverändert (laut Diff), Signatur genügt dem zweiten Aufrufer.
- Sprachvertrag: Test-Meldungen Englisch, Slice Deutsch.
- Traceability: Commit-Nachricht nennt Slice-Kennung und [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).

Nicht geprüft:

- Der Docker-Daemon des CI-Runners und Docker Desktop (Slice nennt es als nicht gemessen).
- Zwei gleichzeitige Läufe gegen denselben Daemon (F-3), nicht erzeugt.
- Der Test-Image-Wechsel auf BuildKit selbst (Gate-Änderung); nur der Builder-Wortlaut wurde auf dem Host geprüft.
- Voller `make test-docker` (655,7 s laut Slice): nicht wiederholt, die Angabe des Implementers „Exit 0, Liste der Container-IDs gleich“ ist hier nicht belegt; Aufgabe des Verifiers.
- `make verify-depguard`: keine Importregel berührt.
