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
sagt Docker Desktop nicht mehr zu ([`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)),
und es gibt keinen Docker-Desktop-Host.

## 2. Trigger (Welle startet)

- Das Tap-Repo `pt9912/homebrew-u-boot` besteht und das Secret `HOMEBREW_TAP_GITHUB_TOKEN` besteht
  `tap-check` (belegt durch den Tap-Nachzug der Releases `v0.7.0` bis `v0.9.0`).
- Der Release `v0.9.0` ist released ([`slice-v1-release-cut-v0.9.0`](done/slice-v1-release-cut-v0.9.0.md)).

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
| [`slice-v2-homebrew-formula`](in-progress/slice-v2-homebrew-formula.md) | Homebrew-Formel für u-boot | [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: nichts mehr; Smoke und Mac-Erprobung liegen vor (siehe den Slice).

## 6. Out-of-Scope für diese Welle

- Homebrew-Core-Submission: eigener Antragsprozess beim Homebrew-Maintainerteam, erst bei stabilem Tap
  und nachweisbarer Nutzerbasis.
- Linux-Brew: bietet keinen Mehrwert über das direkte Binary oder das GHCR-Image.
- Docker Desktop, Podman als Host-Engine und VS Code/Codespaces:
  [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
  sagt sie für `docker` nicht zu. Die Docker-Desktop-Messung bleibt als
  [`slice-v1-sandbox-docker-nachholmessung-desktop`](open/slice-v1-sandbox-docker-nachholmessung-desktop.md)
  in `open/`, ohne Welle, bis ein Docker-Desktop-Host oder ein Nutzerbericht vorliegt.

## 7. Closure-Notiz

Offen.
