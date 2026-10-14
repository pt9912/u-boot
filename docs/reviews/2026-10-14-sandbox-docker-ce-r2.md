# Review-Report: slice-v1-sandbox-docker-ce (r2) — 2026-10-14

**Review-Art:** Code — Folgelauf zur Behebung der Findings F-1 bis F-6 des Erstreports.

**Gegenstand:** Diff `eea88cc..3cfa710` (Commit `3cfa710`: Fingerprint-Prüfung nur über den Primärschlüssel, stärkere Golden Cases und E2E).

**Skill:** `.harness/skills/reviewer.md` @ `3cfa710`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- [`slice-v1-sandbox-docker-ce`](../plan/planning/in-progress/slice-v1-sandbox-docker-ce.md)
- [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)
- Vorgänger-Report `2026-10-14-sandbox-docker-ce.md` (Findings F-1 bis F-7)
- `AGENTS.md` (Hard Rules, §3.7), `harness/review.md`, `harness/replay.md`

Nicht ausgeführt: Docker-Builds und `make test-docker` (Reviewer-Rolle, keine privilegierten Container). Das awk-Programm wurde lokal mit synthetischer `--with-colons`-Ausgabe geprüft (ein Primärschlüssel mit Subkey, zwei Primärschlüssel).

---

## Status der Findings des Erstreports

| Erstreport | Stand | Beleg |
|---|---|---|
| F-1 (Substring-Pin) | behoben, Rest siehe R2-3 | Reihenfolge Prüfung, Repo-Eintrag, Install und `signed-by` gepinnt: generate_sandbox_test.go:675 |
| F-2 (Version/Herkunft) | behoben mit Schwächen, siehe R2-1 und R2-2 | internal/e2e/sandbox_devcontainer_docker_test.go:460 |
| F-3 (Egress plus apt-get update) | Test vorhanden, Aussagekraft eingeschränkt, siehe R2-4 | internal/e2e/sandbox_devcontainer_docker_test.go:565 |
| F-4 (Replay-Zeile) | behoben | harness/replay.md:52 |
| F-5 (nur erster Schlüssel) | behoben | Dockerfile.tmpl:32 |
| F-6 (Kommentar egressHosts) | behoben | devcontainer_egress.go:98 |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R2-1 | MEDIUM | Der Egress-Test bildet die Eigenschaft nicht belastbar ab. Das Skript lautet `sudo apt-get update 2>&1 \| tail -5 && echo APT_UPDATE_OK; …`. Ohne `pipefail` ist der Exit-Status der Pipeline der von `tail`, `APT_UPDATE_OK` erscheint also auch bei fehlgeschlagenem `apt-get update`. Der Container-Exit-Code stammt nur vom letzten Kommando (`! curl … && echo BLOCKED_OK`). Übrig bleibt der String-Abgleich auf `Err:` und `Failed to fetch` im Ausschnitt der letzten fünf Zeilen. Die `Err:`-Zeile steht bei einem nicht erreichbaren Host regulär weiter oben und fällt aus dem Ausschnitt; `W: Failed to fetch …` hängt davon ab, dass es unter die letzten fünf Zeilen fällt (mehrzeilige Fortsetzungszeilen, weitere Quellen mit Warnungen verschieben es). Fehlerszenario: `download.docker.com` fehlt in der Allowlist, die Warnzeile liegt außerhalb des Ausschnitts, der Test bleibt grün. Ich konnte das Szenario nicht ausführen; die Fragilität folgt aus dem gelesenen Skript. | [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) Entscheidung 3 | internal/e2e/sandbox_devcontainer_docker_test.go:583 | ja — `make test-docker` mit entferntem Host (Mutation von egressDockerAptHost) | Negativ-Mutation des Tests nicht rot gesichert |
| R2-2 | LOW | Die Versionsprüfung vergleicht die Hauptversion als String: `major < "28"`. Dreistellige Hauptversionen (zum Beispiel `100`) gelten als kleiner als `28` und erzeugen einen Fehlalarm; einstellige sind über `len(major) < 2` korrekt abgefangen. Ist die letzte Ausgabezeile keine Versionszahl (zum Beispiel `bash: docker: command not found`, Hauptversion-Rest `bash: docker: command not found` beginnt mit einem Buchstaben, der größer als `2` ist), gilt sie als nicht kleiner und der Test bleibt grün. | Maintainability, [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) §Fitness Function | internal/e2e/sandbox_devcontainer_docker_test.go:473 | ja — `make test-docker`; Ausdruck als Zahl parsen wäre die Gegenprobe | Versionsvergleich auf Strings statt Zahlen |
| R2-3 | LOW | Die Paketprüfung verwendet `strings.Contains(out, pkg)` auf kombinierter Ausgabe (`CombinedOutput`) von `dpkg-query -W …; docker version`. Bei einem nicht installierten Paket nennt die Fehlermeldung von dpkg-query den Paketnamen („no packages found matching docker-ce“); die Prüfung auf Vorhandensein ist dann trotzdem erfüllt, und das `;` verdeckt den Exit-Code von dpkg-query. Gerettet wird der Test im Szenario Debian-Daemon nur durch die Versionsprüfung (26.1 < 28). Die Herkunft ist damit indirekt, nicht direkt belegt. Außerdem liest `fields[len(fields)-1]` ohne Längenprüfung; bei leerer Ausgabe wäre das ein Index-Panik statt eines Testfehlers. | [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) §Fitness Function | internal/e2e/sandbox_devcontainer_docker_test.go:463 | ja — `make test-docker` | Positiv-Prüfung auf Fehlerausgabe statt auf Zustand |
| R2-4 | LOW | Der Golden Case pinnt Reihenfolge, `signed-by` und Pin-Wert, aber weder das awk-Programm (Primärschlüssel-Auswahl) noch die `&&`-Verkettung zwischen Prüfung und Repo-Eintrag. Die Änderung aus F-5, also die eigentliche Korrektur, kann auf die frühere `exit`-Variante zurückgesetzt werden, ohne dass der Test rot wird; ebenso ein `;` statt `&&` nach der Prüfung. | [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md), `harness/replay.md` §Replay Rules | internal/hexagon/application/generate_sandbox_test.go:676 | ja — `make test` mit mutierter Vorlage | Sicherheitsrelevante Eigenschaft nur teilweise gepinnt |
| R2-5 | INFO | Der Doc-Kommentar von egressHosts nennt den Host nun, die Zeile ist mit über 100 Zeichen deutlich länger als die Folgezeile. Der Inhalt entspricht [AGENTS.md](../../AGENTS.md) §3.7 (Indikativ, Zusammensetzung, Herkunft-Anker [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)); kein Linter für Zeilenlänge ist in `.golangci.yml` gesetzt. Keine Änderungspflicht. | `AGENTS.md` §3.7 | internal/hexagon/application/devcontainer_egress.go:97 | nein — Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| awk-Ausdruck (`pub`/`sub`/`fpr`) | geprüft, ohne Befund: lokal mit synthetischer Ausgabe eines Primärschlüssels mit Subkey ergibt genau den Primär-Fingerprint, ohne Subkey-Fingerprint; bei zwei Primärschlüsseln entstehen zwei Zeilen, der Vergleich `= "9DC8…"` schlägt fehl und die `&&`-Kette bricht ab; leere Ausgabe ergibt den leeren String und bricht ebenfalls ab. Ein vorangestellter gpg-Fehler (Pipeline-Status) ändert daran nichts, da nur die Ausgabe verglichen wird |
| Quoting und `$10` | geprüft, ohne Befund: das awk-Programm steht in einfachen Anführungszeichen innerhalb von `$( … )` in doppelten; `$10` wird nicht von der Shell expandiert; Dockerfile-RUN expandiert nichts selbst. `text/template` sieht keine `{{`-Folge (nur einzelne Klammern `{p=1}`, `{print $10}`) |
| Reihenfolge im Golden Case | geprüft, ohne Befund: im Docker-Dockerfile ohne Egress ist der erste `apt-get install` der Docker-Install (Podman-Block schließt sich aus, Egress-Block folgt später), die Indexvergleiche sind damit zulässig |
| generateSandboxProjectWith | geprüft, ohne Befund: Verhalten für die bestehenden Aufrufer unverändert (gleiche Map); die Iterationsreihenfolge der Map war schon zuvor ungeordnet |
| Egress-Test, übrige Teile | geprüft, ohne Befund: `docker build` und `rmi` per Cleanup, `generatedRunArgs`, Blockade von `example.com` über `! curl` ist nicht durch Netzausfall verdeckt, da `apt-get update` im selben Lauf Zugriff braucht (siehe R2-1 für die Grenze) |
| Replay-Zeile | geprüft, ohne Befund: nennt Paketquelle, Fingerprint-Prüfung vor Repo-Eintrag und Install sowie `download.docker.com` nur bei `docker`; stimmt mit Code und Tests überein |
| Hard Rules | geprüft, ohne Befund: kein `//nolint`, keine Gate-Lockerung, keine ADR-Änderung im Diff, keine Spec-Änderung (Decken-Regel unberührt) |
| Sprachvertrag, Exit-Codes, Architektur | geprüft, ohne Befund: keine neue CLI-Ausgabe, keine neuen Sentinels, keine Importänderung |

Nicht geprüft:

- Tatsächlicher Image-Build, Daemon-Start, `make test-docker`, `make gates` — Reviewer führt keine privilegierten Container aus; nur `make docs-check` im Rahmen dieses Reports.
- Reales Verhalten von `apt-get update` bei fehlendem Host (Ausgabeformat, Zeilenzahl) — nicht ausgeführt; R2-1 ist eine Aussage über die Prüflogik, nicht über eine gemessene Ausgabe.
- Live-Gültigkeit des Fingerprints — kein Netzzugriff auf `download.docker.com`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Negativ-Mutation des Tests nicht rot gesichert · Positiv-Prüfung auf Fehlerausgabe statt auf Zustand · Sicherheitsrelevante Eigenschaft nur teilweise gepinnt

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings; die Korrektur des Produktcodes (Fingerprint-Prüfung) ist ohne Befund. R2-1 (MEDIUM, Testaussagekraft) ist vor Slice-Closure zu klären oder in der Closure-Notiz mit Ausgang zu führen.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure. Dieser Report ersetzt keine Verifikation.
