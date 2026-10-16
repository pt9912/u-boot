# `make coverage-gate` — Coverage-Schwelle über `./internal/...`, bootstrap-aware

## Vertrag

Rot, wenn die Gesamt-Coverage der Produktionspakete unter `THRESHOLD` liegt
(Vorgabe 90 %), oder wenn die Messung nichts liefert, obwohl Produktionscode
existiert. Was als Coverage-Bootstrap gilt, legt
[`LH-FA-BUILD-008`](../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap) fest.

## Grenze — was das Grün nicht abdeckt

1. **Nur `./internal/...`** — Pakete außerhalb (`cmd/`) fließen nicht in die
   Zahl ein; permanent, weil sie nur verdrahten.
2. **Nur Default-Tests** — Tests mit dem Build-Tag `docker` (`make test-docker`)
   laufen nicht im Coverage-Stage und tragen nichts bei; permanent.
3. **Eine Zahl für das Ganze** — ein einzelnes Paket kann weit unter der
   Schwelle liegen, solange der Durchschnitt sie hält; nicht heilbar ohne
   paketweise Schwellen.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`make coverage-gate` gibt `coverage-gate: OK — coverage <p>% meets threshold <t>%` aus.

## Ausgabe und Ausgänge

Die Auswertung läuft in `scripts/coverage-gate.sh` über die Ausgabe von
`go tool cover -func`.

| Exit | Bedeutung |
|---|---|
| 0 | Schwelle erreicht, oder Bootstrap-Modus ohne Produktionspakete |
| 1 | Coverage unter der Schwelle (`FAIL — coverage <p>% below threshold <t>%`) |
| 2 | Eingabe fehlt oder ist leer ohne Bootstrap, keine `total:`-Zeile, Prozentwert nicht lesbar — die Messung selbst ist gescheitert, nicht die Schwelle |

Exit 1 und Exit 2 unterscheiden den geprüften Baum (zu wenig Tests) von der
Umgebung (`go test` ist vor der Messung gescheitert); die Zeile mit
`hint:` nennt dann die wahrscheinliche Ursache.

## Sperren

- `COVERAGE_BOOTSTRAP` — nur der Coverage-Stage des `Dockerfile` setzt `1`, und
  nur solange `./internal/...` keine Produktionspakete trägt; eine leere Eingabe
  ist sonst Exit 2.

## Bindung

Schwelle 90 %; [`LH-FA-BUILD-008`](../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap).
