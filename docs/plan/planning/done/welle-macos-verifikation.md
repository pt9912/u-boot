# Welle welle-macos-verifikation: macOS-Verifikation (Homebrew-Formel)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld. **Geplante
Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter *Nächste
Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-16.

---

## 1. Welle-Ziel

Die Zusage, die u-boot für macOS macht, ist auf einem echten Mac belegt: `brew install` der Formel aus
dem Tap liefert ein lauffähiges `u-boot` auf Apple Silicon und Intel. Danach sagen README und Handbuch den
macOS-Weg zu, weil er gemessen ist, nicht weil er vorgesehen ist.

Die Docker-Desktop-Messung für `nestedRuntime: docker` gehört nicht mehr zur Welle: Das Lastenheft 0.7.0
sagt Docker Desktop nicht mehr zu ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)),
und es gibt keinen Docker-Desktop-Host.

## 2. Trigger (Welle startet)

- Das Tap-Repo `pt9912/homebrew-u-boot` besteht und das Secret `HOMEBREW_TAP_GITHUB_TOKEN` besteht
  `tap-check` (belegt durch den Tap-Nachzug der Releases `v0.7.0` bis `v0.9.0`).
- Der Release `v0.9.0` ist released ([`slice-v1-release-cut-v0.9.0`](slice-v1-release-cut-v0.9.0.md)).

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber der Slice-DoD: Der Smoke im Tap und die Erprobung auf einem echten Mac decken
beide Architekturen ab, und erst danach ändert sich die Aussage der Nutzerdoku.

- Der Slice der Welle liegt in `done/`.
- Der macOS-Smoke im Tap zur aktuellen Formel ist auf Apple Silicon **und** Intel grün.
- `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot` und `brew install u-boot` sind auf einem
  echten Mac erprobt (`u-boot --version`, `init`, `doctor` ohne Absturz).
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md) | Homebrew-Formel für u-boot | [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: nichts mehr; Smoke und Mac-Erprobung liegen vor (siehe den Slice).

## 6. Out-of-Scope für diese Welle

- Homebrew-Core-Submission: eigener Antragsprozess beim Homebrew-Maintainerteam, erst bei stabilem Tap
  und nachweisbarer Nutzerbasis.
- Linux-Brew: bietet keinen Mehrwert über das direkte Binary oder das GHCR-Image.
- Docker Desktop, Podman als Host-Engine und VS Code/Codespaces:
  [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
  sagt sie für `docker` nicht zu. Die Docker-Desktop-Messung bleibt als
  [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md)
  in `open/`, ohne Welle, bis ein Docker-Desktop-Host oder ein Nutzerbericht vorliegt.

## 7. Closure-Notiz

Detailquelle: [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md) §7 (die Welle trägt nur diesen Slice; einen Release-Cut-Slice gibt es
nicht, weil sie kein Release ist).

- **Geliefert:** Die Homebrew-Formel ist auf macOS belegt: Smoke im Tap auf `macos-14` (arm64) und `macos-15-intel` (amd64) grün
  (Lauf `38900477533`, Tap-Commit `a8498fd`), `brew tap`/`trust`/`install` auf Apple Silicon durch den Projektinhaber erprobt (Zusage).
- **Was hat funktioniert:** Die Closure-Trigger aus §3 ließen sich einzeln belegen; der Verifier prüfte den Smoke-Lauf am Log.
- **Was ging anders als geplant:** Die Welle sollte auch die Docker-Desktop-Messung tragen. Es gibt keinen Docker-Desktop-Host; das
  Lastenheft 0.7.0 nahm die Zusage für Docker Desktop zurück ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)),
  und der Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) verließ die Welle.
  Die Intel-Abbrüche des Smokes kamen vom abgekündigten Label `macos-13`, nicht von knappen Runnern.
- **Trigger-Audit:** Carveout: [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) in [`carveouts.md`](../in-progress/carveouts.md) aufgelöst.
  Bootstrap-aware Gate: `coverage-gate` unverändert, 90,90 % gegen 90 %; 0 offen. ADR: [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md)
  und [`ADR-0018`](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) tragen keinen eingetretenen Re-Evaluierungs-Trigger (die Geschichte der zweiten
  vermerkt die Lastenheft-Eingrenzung); 0 offen.
- **Steering-Loop-Einträge:** Kein Eintrag über der Schwelle (das Beobachtungs-Register ist leer). Gezählt, nicht verkörpert: ein Smoke gegen ein
  abgekündigtes Runner-Label bleibt „queued" statt rot zu werden.
- **Beobachtungs-Register:** keine Beobachtung angelegt.
- **Folge-Slices:** [`slice-v2-homebrew-smoke-haertung`](../open/slice-v2-homebrew-smoke-haertung.md) (Smoke-Härtung, Formel-`desc`, Workflow-Kommentar);
  [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) (ohne Welle, wartet auf einen
  Docker-Desktop-Host).
- **Verifikation:** `make gates` grün (Exit 0); Verification Evidence
  [`2026-10-16-verification-slice-v2-homebrew-formula`](../../../reviews/2026-10-16-verification-slice-v2-homebrew-formula.md), Review-Report
  [`2026-10-16-welle-macos-verifikation`](../../../reviews/2026-10-16-welle-macos-verifikation.md).
