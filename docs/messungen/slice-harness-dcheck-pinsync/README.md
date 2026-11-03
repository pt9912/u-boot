# Messprotokoll: Pin-Gleichheit im Gatter und Mismatch-Gegenproben zu [`slice-harness-dcheck-pinsync`](../../plan/planning/in-progress/slice-harness-dcheck-pinsync.md)

Beleg-Läufe des Umsetzungs-Diffs (Pin-Inventar, [`LH-FA-DEV-011`](../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)): der positive Gatter-Lauf, je Pin-Gruppe eine Mismatch-Gegenprobe mit Rückführung, die Negativprobe des GNU-make-Guards, die Byte-Stabilität der `go.mod`-Normierung und die Digest-Verifikation des d-check-Bumps. Rohdaten in diesem Verzeichnis.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-11-03 |
| Host | Linux `x86_64`, lokale Docker-Engine, GNU make 4.4.1 (Host und Devcontainer-Image `vsc-dcsdmlh0zxpm8ik-…:latest` je `make --version` gemessen) |
| d-check | `ghcr.io/pt9912/d-check@sha256:5d2f650d102c1d8224bb26930ec96d237506c07612fb8e0dbcdc42291702c57f` (Tag `v0.91.0`) |
| Gegenstand | Arbeitsbaum des Umsetzungs-Diffs; die Gegenproben laufen in einer Kopie dieses Baums (Verzeichnis außerhalb des Repos), das veränderte Byte und die Rückführung je Lauf in den Rohdaten |

## Reihenfolge und Aufrufe

Positive und Gegen-Läufe sind je ein `docker run --rm --network none` gegen die d-check-Referenz an der Repo-Wurzel der Kopie (dieselbe Zeile wie der `docs-check` der Gatter, ohne Argumente — die Modulmenge kommt aus `.d-check.yml`):

```bash
docker run --rm --network none -v "$PWD":/repo:ro \
  ghcr.io/pt9912/d-check@sha256:5d2f650d102c1d8224bb26930ec96d237506c07612fb8e0dbcdc42291702c57f
```

Gegenprobe: ein Pin-Wert in der Kopie bewusst verändert, Lauf, Muster der Befund-Zeile notiert, Byte zurück, Gegenprobe des Grünlufs. Die Byte-Stabilität der tidy-Normierung vergleicht SHA-256 vor und nach einem weiteren Aufruf des Modul-Werkzeugs am Docker-Stage (`-u "$(id -u):$(id -g)"`, Caches im Container). Die Negativprobe des make-Guards ruft `make MAKE_VERSION=4.2 -n help` (Kommandozeilen-Setzung sticht den eingezogenen Wert), die Positivprobe `make -n help` ohne Setzung. Die Digest-Verifikation vergleicht die Digest-Zeile der Release-Notes (GitHub-API-Lesung des Tags) und die `RepoDigests`-Ausgabe des gezogenen Images.

## Ergebnisse

**1. Positiver Gatterlauf** ([`roh-01-positive.txt`](roh-01-positive.txt)): drei Pins, elf Quellen gelesen, kein Befund, Lauf grün (`pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`). Die Pin-Gruppen `go`, `golangci-lint` und `devcontainer-cli` tragen je Quelle genau einen Treffer; der Vergleich ist byte-gleich.

**2. Mismatch der Go-Direktive** ([`roh-02-mismatch-go.txt`](roh-02-mismatch-go.txt)): die `go`-Zeile der `go.mod` der Kopie auf `1.26.0` zurückgesetzt — ein Befund `pinsync-mismatch` auf der Zeile der Direktive, die Meldung nennt Pin und alle Quellen mit ihrem Wert; der Lauf endet rot. Nach der Rückführung auf `1.27.0` ist der Lauf byte-identisch grün (zweite Hälfte der Datei).

**3. Mismatch der CLI-Pin-Stelle** ([`roh-03-mismatch-cli.txt`](roh-03-mismatch-cli.txt)): die Abhängigkeits-Version in `package.json` der Kopie auf `0.90.0` geändert — ein Befund `pinsync-mismatch` auf der Zeile der `package.json`-Quelle; nach der Rückführung auf `0.89.0` grün.

**4. GNU-make-Guard** ([`roh-04-make-guard.txt`](roh-04-make-guard.txt)): mit `MAKE_VERSION=4.2` bricht der Lauf vor jeder Logik mit `$(error)` ab (`GNU make 4.3 or newer is required …`), ohne Targets zu starten; ohne die Setzung (ist 4.4.1) parst derselbe Aufruf normal.

**5. Byte-Stabilität der tidy-Normierung** ([`roh-05-tidy-byte-stabil.txt`](roh-05-tidy-byte-stabil.txt)): der Aufruf am Docker-Stage über der gehobenen Direktive `1.27.0` ändert die `go.mod` nicht (SHA-256 vor und nach dem Lauf gleich); der erste Aufruf des Umsetzungs-Diffs normalisierte das Require-Layout (Cobra-Abhängigkeiten aus `// indirect` in das direkte `require`-Block, Indirect-Aufräumzeile getrennt) — der Stand des Umsetzungs-Diffs ist dieser normierte.

**6. Digest-Verifikation** ([`roh-06-digest-verifikation.txt`](roh-06-digest-verifikation.txt)): die Release-Notes des Tags `v0.91.0` tragen die Digest-Zeile `ghcr.io/pt9912/d-check@sha256:5d2f65…c57f`; die `RepoDigests`-Ausgabe des gezogenen Images nennt denselben Digest, `DCHECK_DIGEST` in `d-check.mk` trägt diesen Wert.

## Grenzen der Messung

Der Wächter prüft Gleichheit, keine Aktualität. Die Fanggruppe der `go`-Gruppe liest den Major.Minor-Teil (die Direktive steht auf dem Pin-Minor, die Patch-Stelle der Toolchain bleibt in den `GO_VERSION`-Stellen) — eine Patch-Hebung der `GO_VERSION`-Stellen allein bleibt im Gatter grün und ist die Routine der `GO_VERSION`-Stellen je [`SPEC-054`](../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv). Eine Pin-Hebung über Minor oder Wert trifft jede Quellen-Zeile einzeln: eine vergessene Zeile ist ein rotierter Lauf (`pinsync-mismatch`), nicht eine stille Drift. Die Gegenproben laufen je einmal je Pin-Gruppe (keine Streuung bestimmt); das JDK ist nicht gruppiert (eine einzige Einstellen-Stelle, Template-`ARG`).