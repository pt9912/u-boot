# `make verify-depguard` — beweist, dass jede `depguard`-Schichtregel auf einen echten verbotenen Import anschlägt

## Vertrag

Kein Gate, sondern ein Nachweis über das Gate `make lint`: Das Skript fügt je
Regel einen absichtlich verbotenen Import ein, lässt `make lint` laufen, prüft
die erwartete Regel-Beschreibung in der Ausgabe und nimmt den Import wieder
heraus. Es urteilt nicht über den Zustand des Repos, sondern darüber, ob die
Regeln aus [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
überhaupt feuern können. Es steht nicht in `make gates`, weil jede Regel einen
vollen `make lint`-Lauf kostet.

Wo kein vorhandenes, von der Regel verbotenes Paket ohne Rückweg zum Zielpaket
existiert (Regel `domain-isoliert`: jedes Paket unter `application`, `port` und
`adapter` importiert `domain`, ein Import daraus bildet einen Zyklus, und der
Typecheck bricht vor `depguard` ab), legt das Skript ein Wegwerf-Paket ohne
Importe an (`internal/hexagon/port/verifydepguardprobe`, fünftes Feld der Zeile
in `CASES`), importiert es und entfernt es nach dem Lauf, auch bei Abbruch. Das
Paket trägt einen Paketkommentar, damit `revive` (`package-comments`) nicht
anschlägt: Im Lauf dieses Falls ist der `depguard`-Treffer der einzige Befund,
und der rote Exit von `make lint` belegt die Regel, nicht einen Begleitbefund.

## Grenze — was das Grün nicht abdeckt

1. **Je Regel ein Fall** — geprüft wird ein verbotener Import je Regel, nicht
   jede Import-Kombination; permanent.
2. **Nur die Regeln der Tabelle `CASES`** in `scripts/verify-depguard.sh` — eine
   neue Regel in `.golangci.yml` ohne neuen Fall bleibt ungeprüft; heilbar durch
   den Eintrag in `CASES`.
3. **Nicht im Gate** — nach einer Änderung der Schichten oder der
   `depguard`-Konfiguration muss der Lauf von Hand angestoßen werden.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`bash scripts/verify-depguard.sh` listet je Regel eine Zeile.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | alle Regeln haben mit der erwarteten Beschreibung angeschlagen |
| 1 | mindestens eine Regel hat nicht angeschlagen (Logs im ausgegebenen Verzeichnis) |
| 2 | Vorbedingung verletzt: Arbeitsbaum in einem Schicht-Verzeichnis nicht sauber, oder ein Probe-Verzeichnis aus `CASES` existiert schon |

## Sperren

- Schmutziger Arbeitsbaum — `git status --porcelain` über den Schicht-Verzeichnissen
  muss leer sein, sonst wären Einfügung und Aufräumen nicht von eigenen
  Änderungen unterscheidbar → committen oder stashen.

## Bindung

[`ADR-0003`](../../docs/plan/adr/0003-solid-nahes-lint-profil.md);
[`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln).
