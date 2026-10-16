# `make verify-depguard` — beweist, dass jede `depguard`-Schichtregel auf einen echten verbotenen Import anschlägt

## Vertrag

Kein Gate, sondern ein Nachweis über das Gate `make lint`: Das Skript fügt je
Regel einen absichtlich verbotenen Import ein, lässt `make lint` laufen, prüft
die erwartete Regel-Beschreibung in der Ausgabe und nimmt den Import wieder
heraus. Es urteilt nicht über den Zustand des Repos, sondern darüber, ob die
Regeln aus [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
überhaupt feuern können. Es steht nicht in `make gates`, weil jede Regel einen
vollen `make lint`-Lauf kostet.

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
| 2 | Vorbedingung verletzt: Arbeitsbaum in einem Schicht-Verzeichnis nicht sauber |

## Sperren

- Schmutziger Arbeitsbaum — `git status --porcelain` über den Schicht-Verzeichnissen
  muss leer sein, sonst wären Einfügung und Aufräumen nicht von eigenen
  Änderungen unterscheidbar → committen oder stashen.

## Bindung

[`ADR-0003`](../../docs/plan/adr/0003-solid-nahes-lint-profil.md);
[`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln).
