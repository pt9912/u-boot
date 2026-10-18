# Review-Report: slice-cr-sandbox-fuse-startfehler — 2026-10-18

**Review-Art:** Code (Plan-Konformität des Diffs gegen Slice, Lastenheft, Spezifikation, `AGENTS.md`; CR am Vertrags-Stratum)

**Gegenstand:** Commit-Range `17e4b90..HEAD` (`c1ef60e`, `4c16901`, `a28a9f5`, `e4e5d0d`, `b0919ad`, `702335a`) des Slice `slice-cr-sandbox-fuse-startfehler` <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `702335a` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice `slice-cr-sandbox-fuse-startfehler` (§1 Wortlaut Stelle 1/2/3, §3 Plan-Tabelle, §6 Risiken) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Messung M13/M14)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman)
- `AGENTS.md` §3.5–§3.10, §5; `harness/review.md`
- Freigaben des Projektinhabers (2026-10-17/18) laut Auftrag: Wortlaut Stelle 1/2/3, Version 0.13.0, `doctor` immer Fehler, `vfs`-Fallback bleibt

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Spezifikation bindet die Fehler-Einstufung an einen „Linux-Host mit lokaler Engine“; `checkSandboxRuntime` prüft die Engine-Lage nicht (kein `DOCKER_HOST`-Bezug), jeder Linux-Host mit fehlendem `/dev/fuse` liefert den Fehler. Slice §6 (Falsch-positiv bei entfernter Engine) führt das Risiko als „offen“. | `AGENTS.md` §3.10 (Aussage über Code lesen), Maintainability | `spec/spezifikation.md` · „als Fehler aus (Umgebungsproblem, Exit-Code `11`, unabhängig von `onUnavailable`“; `internal/hexagon/application/doctor_sandbox.go` · „if fuse, _ := s.fs.Exists("/dev/fuse"); fuse“ | nein — Review-only | Spec-Bedingung ohne Code-Entsprechung |
| F-2 | LOW | Die Statuszeile in `README.md` und `README.de.md` nennt als Voraussetzung auf dem Container-Host nur `/dev/net/tun`; ab 0.13.0 gehören beide Devices zur Start-Voraussetzung. Die Plan-Tabelle §3 hat die READMEs über die Suche nach `fuse`/`vfs` als „unverändert“ geführt, nicht über `/dev/net/tun`. Die Aussage ist als „neu für bestehende Projekte“ lesbar und nicht falsch, aber unvollständig. | `AGENTS.md` §5 Zeile 5 (öffentliche Verträge nachziehen) | `README.md` · „projects need `/dev/net/tun` on the container host“; `README.de.md` · „brauchen `/dev/net/tun` auf dem Container-Host“ | nein — Review-only | Öffentliche Voraussetzungs-Aussage nicht auf beide Devices nachgezogen |
| F-3 | LOW | Der Hinweistext des Fehler-Items (`Hint`: `modprobe fuse` / `nestedRuntime` auf `none`) ist durch keinen Test gepinnt; die Mutation, die den `nestedRuntime`-Teil streicht, überlebt `make test`. Meldungstext, Severity, Fallback-Flag und Nicht-Linux-Zweig sind gepinnt. | `harness/review.md` §Prüflinsen (Test belegt Verhalten), Maintainability | `internal/hexagon/application/doctor_sandbox_test.go` · `TestDoctor_SandboxRuntime` (prüft nur `msgPart`) | ja — `make test` (Mutation) | Hinweistext einer Prüfung ungepinnt |
| F-4 | INFO | Die Zeile `{"runtime: fuse missing, fail policy is an unflagged error",sandboxYAMLPodman + …` trägt kein Leerzeichen nach dem Komma (nicht `gofmt`-konform). Kein Gate fängt das (`make lint` läuft grün; `gofmt` steht nicht in `.golangci.yml`). | Maintainability (ohne Konventions-Anker, daher INFO) | `internal/hexagon/application/doctor_sandbox_test.go` · `unflagged error",sandboxYAMLPodman` | nein — `make lint` meldet es nicht | Formatierung ohne Gate |
| F-5 | INFO | `SPEC-043` Ausnahme 2 und `LH-FA-CLI-005` Ausnahme 2 zählen `vfs`-Storage weiter unter den Fallbacks, die `doctor` „soweit vom Host aus ermittelbar“ ausweist; Host-seitig ist `vfs` jetzt nicht mehr ermittelbar. Das ist durch Slice §3 (Umsetzungs-Slice zieht `SPEC-043` nach) und die Freigabe „`vfs`-Fallback bleibt“ gedeckt, kein Widerspruch im Vertrag. | Slice §1/§3, Hinweis an den Umsetzungs-Slice | `spec/spezifikation.md` · „(2) die Ausweisung der Degradations-Fallbacks (Container ohne Daemon, `vfs`-Storage“ | nein — Review-only | Bewusst offener Nachzug (Folge-Slice) | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-6 | INFO | Der Auftrag verweist auf `AGENTS.md` §5 „Zeile 6“; die Tabelle hat fünf Zeilen. Geprüft wurde die dem Sinn nach zutreffende Zeile 5 (Verträge nachziehen) und die Handbuch-Regel (keine Kennungen/Chronik): Der Handbuch-Diff trägt keine `LH-*`/`ADR-*`/Slice-Kennung und keine Chronik. | `AGENTS.md` §5 | `AGENTS.md` · „### 5. Dokumentations-Regeln“ | nein — Review-only | Zeilenverweis im Auftrag |
| F-7 | INFO | Der e2e-Test `TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart` ist nicht gelaufen (`make test-docker` nicht ausgeführt: baut ein Podman-Image und startet Container mit `SYS_ADMIN`, im Review-Lauf nicht vorgesehen). Gelesen: Kontrolllauf mit den echten `runArgs` vor dem Fehlerlauf (nicht tautologisch), `--rm`, Image-Aufräumen per `t.Cleanup`, 15-Minuten-Kontext, fehlendes `docker` ist `Fatalf` wie bei den Nachbartests (kein stilles Skip). Die Meldungsprüfung „error gathering device information“ ist Docker-spezifisch (Nachbar-Muster). | Slice §2 (`make test-docker` auf dem Entwicklungs-Host) | `internal/e2e/sandbox_devcontainer_docker_test.go` · „func TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart“ | ja — `make test-docker` (Verifier) | Sensor nicht gelaufen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` Diff: Stelle 1, Stelle 2, Versionsfeld `0.13.0`, §16-Zeile (Datum 2026-10-18, Messstand, „Vereinbarung mit dem Projektinhaber“) | geprüft, ohne Befund — Wortlaut von Stelle 1 und 2 stimmt mit Slice §1 überein; die §16-Zeile stimmt mit dem Vorschlag überein; keine weitere Zusage verändert (Diff: 7 Zeilen); Ausnahme 2 in `LH-FA-CLI-005` unverändert wie im Slice begründet | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Referenzmodell: Lastenheft und Spezifikation verweisen nicht abwärts | geprüft, ohne Befund — Diff nennt weder Slice noch ADR noch Carveout |
| Lastenheft ↔ Spezifikation ↔ Code ↔ Handbuch ↔ CHANGELOG (Suche `/dev/fuse`, `vfs`, `/dev/net/tun`, `Degradation`, `onUnavailable` im Baum) | geprüft, ohne Widerspruch außer F-1/F-2/F-5; Fixtures in `internal/adapter/driving/cli/*_test.go` mit „falls back to vfs storage“ sind Vertreter eines `DegradationFallback`-Items, keine Prüfung-Texte (wie in Slice §3 geplant) |
| `checkSandboxRuntime` Zweige: ohne Sandbox/Runtime, Profil-Mismatch, `docker`, Nicht-Linux, AppArmor-Sysctl, `/dev/fuse` vorhanden, fehlend | geprüft, ohne Befund — fehlendes Device ist `SeverityError`, unabhängig von `onUnavailable`, ohne `DegradationFallback`; Meldung englisch und im Muster der Nachbarn („(LH-FA-DEV-007, exit 11)“ wie beim AppArmor-Zweig) | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Wirkung auf `doctor --quiet`/`--strict`/`--json`/Exit 11 | geprüft, ohne Befund — ein Fehler-Item wird über die bestehende Severity-Abbildung ausgegeben (`--quiet` zeigt Fehler, Exit 11 über die zentrale Klassifikation); es entsteht kein neuer Sentinel, damit entfällt der Dual-Classifier-Anker |
| Mutationen (`make test`, je danach `git checkout -- <Datei>`): Nicht-Linux-Meldung gekürzt · Severity `Warn` · `DegradationFallback: true` · Meldung „falls back to vfs“ · Nicht-Linux-Zweig `Error` · `/dev/fuse`-Prüfung invertiert | geprüft, alle sechs getötet; siebte Mutation (Hinweistext) überlebt, siehe F-3 |
| Tests nicht tautologisch (`doctor_sandbox_test.go`) | geprüft, ohne Befund — beide Policies, Nicht-Linux mit und ohne Policy, Flag-Test mit `fallback=false`; Test-Kommentare tragen Zusage/Abgrenzung |
| Kommentare (`AGENTS.md` §3.7/§3.9 (a)/(b)/(c)) in Go-Diff, Tests, e2e | geprüft, ohne Befund — keine Chronik- oder Konjunktiv-Kommentare; kein Helfer zwischen Doc-Kommentar und Funktion; Verweise auf geänderte Symbole/Werte (Fixtures, Startscript-Kommentar) im Baum gegengesucht und stimmig |
| Handbuch (`docs/user/benutzerhandbuch.md`) | geprüft, ohne Befund — keine Kennungen, keine Chronik, Tabelle/Absätze stimmen mit Spec und Code überein |
| CHANGELOG | geprüft, ohne Befund — Regression für Skripte ohne `--strict` benannt, Messstand genannt, „generated files do not change“ stimmt mit dem Diff (kein Template-/Generator-Eintrag im Diff); nichts Ungeliefertes |
| Plan-Tabelle §3 gegen `git diff --stat` (`AGENTS.md` §3.10) | geprüft, ohne Befund — alle neun berührten Dateien stehen in der Tabelle, jede „update“-Zeile hat einen Diff; Ausnahme siehe F-2 (README als „unverändert“ geführt) |
| Sprachvertrag (`LH-LESE-002`) | geprüft, ohne Befund — Meldungen und Hinweis englisch | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Doc-Immutable | geprüft, ohne Befund — `git diff 17e4b90..HEAD -- docs/plan/adr` leer; ADR-0022 unberührt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Hard Rules §3.1/§3.2/§3.6 (Docker-only, kein `//nolint`, keine Gate-Lockerung) | geprüft, ohne Befund — kein `nolint`, `.golangci.yml`/Makefile unberührt; `make lint` und `make test` laufen grün |
| Roadmap-Diff (Ruhe-Marker entfernt) | geprüft, ohne Befund |
| `make test-docker` | nicht geprüft — nicht gelaufen (siehe F-7); Zuständigkeit Verifier |
| `make verify-depguard` | nicht geprüft — Importregeln nicht berührt |
| Replay/Golden (`harness/replay.md`) | nicht geprüft — kein Generator-Output geändert (`runArgs`/Templates unberührt) |
| DoD-Konformität (Closure-Notiz, Register, Risiko-Ausgänge) | nicht geprüft — Verifier-Aufgabe |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Spec-Bedingung ohne Code-Entsprechung · Öffentliche Voraussetzungs-Aussage nicht auf beide Devices nachgezogen · Hinweistext einer Prüfung ungepinnt · Formatierung ohne Gate · Bewusst offener Nachzug (Folge-Slice) · Zeilenverweis im Auftrag · Sensor nicht gelaufen

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die drei LOW-Findings sind Doku- und Test-Präzision ohne Vertragsbruch; F-2 betrifft einen öffentlichen Vertrag und sollte vor der Closure entschieden werden (Nachzug oder begründet unverändert).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. F-1 hängt am offenen Risiko „Falsch-positiv bei entfernter Engine“ (Slice §6), dessen Ausgang bei der Closure zu setzen ist. Dieser Report ersetzt keine Verifikation; `make test-docker` und die DoD-Prüfung liegen beim Verifier.
