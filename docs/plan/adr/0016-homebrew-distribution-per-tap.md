# ADR 0016: Homebrew-Distribution über einen eigenen Tap mit Formel als Release-Asset

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912

**Bezug:** [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung), [ADR-0007](0007-distributionswege-ghcr.md)

**Schärft:** [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung) — Homebrew ist von „vertagt mit Trigger“ auf „gewählt“ gesetzt; Mechanismus und Tap-Form.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0007](0007-distributionswege-ghcr.md) wählte GHCR und das Einzel-Binary und vertagte Homebrew mit dem Trigger „konkrete macOS-Nachfrage“. Die Binaries für sechs Plattformen hängen seit v0.2.0 als Release-Assets am Tag. Der Projektinhaber betreibt für andere Werkzeuge bereits Homebrew-Taps (`pt9912/homebrew-ai-harness-init`, `pt9912/homebrew-d-migrate`) und will dieselbe Distribution für u-boot. Beide Referenzen zeigen zwei Muster: `ai-harness-init` füllt im Release-Workflow ein Formel-Skeleton aus der `SHA256SUMS` desselben Tags, hängt die Formel als Release-Asset an und zieht sie per Job ins Tap nach; `d-migrate` nutzt die Drittanbieter-Action `Justintime50/homebrew-releaser` (nötig dort wegen der JVM-Distribution als Tarball).

Der Accepted-ADR-0007 bleibt unverändert; diese Folgeentscheidung ersetzt die Zeile „Homebrew: vertagt“ (ADR-Disziplin: Accepted ADRs werden nicht umgeschrieben).

## Entscheidung

1. **Eigener Tap pro Werkzeug:** `pt9912/homebrew-u-boot`, Installation per `brew install pt9912/u-boot/u-boot`.
2. **Formel als Release-Asset:** Der `publish`-Workflow erzeugt für jeden Tag eine `SHA256SUMS` der sechs Binaries; für **stabile** Tags füllt `scripts/homebrew-formula-fill.sh` das Skeleton `scripts/homebrew-formula.rb.tmpl` aus genau diesen Digests (fail-closed bei fehlendem oder ungültigem Digest) und hängt `u-boot.rb` als Asset an. Vier Plattformen (macOS/Linux × amd64/arm64; Homebrew trägt kein Windows).
3. **Nachzug per Job:** Der Job `tap` (nur stabile Tags, `needs: publish`) ruft `scripts/tap-nachzug.sh`, das das Asset lädt, die Version prüft und `Formula/u-boot.rb` ins Tap pusht (idempotent). Er läuft nur mit dem Secret `HOMEBREW_TAP_GITHUB_TOKEN`; ohne Secret endet er mit einer Warnung, damit der Release nie vom Tap abhängt.
4. **Prereleases** (`vX.Y.Z-…`) landen nicht im Tap (Homebrew trackt nur Stable).
5. **Smoke im Tap:** `packaging/homebrew-tap/.github/workflows/smoke.yml` (Vorlage für das Tap-Repo) installiert die Formel auf macOS und prüft `--version` und `init`; `doctor` ist dort wegen fehlendem Docker nur auf „kein Absturz“ geprüft.
6. **Nicht Teil dieser Entscheidung:** Homebrew-Core-Einreichung (eigener Antragsprozess, erst bei stabiler Nutzerbasis) und Linux-Brew als beworbener Weg.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Homebrew-Core-Einreichung | maximale Reichweite (`brew install u-boot`) | Antragsprozess, zusätzliche Qualitätsregeln, Kontrolle beim Homebrew-Team |
| B — Drittanbieter-Action (`homebrew-releaser`) | wenig eigener Code | zusätzliche Fremdaktion mit Token-Zugriff (Supply-Chain), Digest nicht an die Release-Assets gekoppelt, für Einzel-Binaries mehr Konfiguration |
| **C — Skeleton + Fill-Skript + Nachzug-Job (Muster `ai-harness-init`)** | Digest stammt aus den Release-Assets desselben Tags, Formel reist als Asset (prüfbar), kein Fremdcode, direkt testbar | eigener Code in `scripts/`, PAT-Secret nötig |
| D — nichts tun | kein Aufwand | macOS-Nutzer ohne Paketweg |

## Konsequenzen

- Positiv: Ein Tag erzeugt Binaries, Digests, Formel und Tap-Commit in einem Vorgang; das Fill-Skript ist per Go-Test (`cmd/uboot/homebrew_scripts_test.go`) gedeckt.
- Negativ: Zwei Repos zu pflegen (Tap, Hauptrepo); das Tap-Repo und das Secret müssen einmalig von Hand angelegt werden; die Installation auf echten macOS-Runnern ist erst mit dem ersten Tag nach Einrichtung belegt.
- Folgepflicht: Tap-Repo anlegen (Vorlagen in `packaging/homebrew-tap/`), Secret `HOMEBREW_TAP_GITHUB_TOKEN` setzen, ersten stabilen Tag abwarten und `brew install` auf macOS verifizieren.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test | Fill-Skript befüllt alle Platzhalter, bricht ohne gültigen Digest vor dem Schreiben ab | `make test` |

## Re-Evaluierungs-Trigger

Eine Nutzerbasis, die die Homebrew-Core-Einreichung rechtfertigt, oder ein Wechsel des Release-Mechanismus (z. B. neue Plattformen).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Entschieden und `Accepted` (Referenzen: `homebrew-ai-harness-init`, `homebrew-d-migrate`) | Vereinbarung mit dem Projektinhaber |
