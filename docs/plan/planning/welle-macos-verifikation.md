# Welle welle-macos-verifikation: macOS-Verifikation (Homebrew-Formel und Docker Desktop)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld. **Geplante
Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter *Nächste
Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-16.

---

## 1. Welle-Ziel

Die beiden Zusagen, die u-boot für macOS macht, sind auf einem echten Mac belegt: `brew install`
der Formel aus dem Tap liefert ein lauffähiges `u-boot`, und der Sandbox-Devcontainer mit
`nestedRuntime: docker` läuft unter Docker Desktop. Danach sagen README und Handbuch die macOS-Wege
zu, weil sie gemessen sind, nicht weil sie vorgesehen sind.

## 2. Trigger (Welle startet)

- Das Tap-Repo `pt9912/homebrew-u-boot` besteht und das Secret `HOMEBREW_TAP_GITHUB_TOKEN` besteht
  `tap-check` (belegt durch den Tap-Nachzug der Releases `v0.7.0` bis `v0.9.0`).
- Der Release `v0.9.0` ist released ([`slice-v1-release-cut-v0.9.0`](done/slice-v1-release-cut-v0.9.0.md)).

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den beiden Slice-DoDs: beide Prüfungen laufen in **einer** Mac-Sitzung, und
erst danach ändert sich die Aussage der Nutzerdoku.

- Beide Slices der Welle liegen in `done/`.
- Der macOS-Smoke im Tap zur aktuellen Formel ist auf Apple Silicon **und** Intel grün.
- `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot` und `brew install u-boot` sind auf einem
  echten Mac erprobt (`u-boot --version`, `init`, `doctor` ohne Absturz).
- Die Messreihe zu `nestedRuntime: docker` liegt für Docker Desktop (macOS) vor, und die Doku sagt
  die Lauffähigkeit erst danach zu.
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-v2-homebrew-formula`](in-progress/slice-v2-homebrew-formula.md) | Homebrew-Formel für u-boot | [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung) |
| [`slice-v1-sandbox-docker-nachholmessung-desktop`](open/slice-v1-sandbox-docker-nachholmessung-desktop.md) | Docker-Desktop-Messung für `nestedRuntime: docker` | [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: Zugriff auf einen Mac mit Docker Desktop (außerhalb des Repos); ohne ihn
  bleiben beide Slices stehen.

## 6. Out-of-Scope für diese Welle

- Homebrew-Core-Submission: eigener Antragsprozess beim Homebrew-Maintainerteam, erst bei stabilem Tap
  und nachweisbarer Nutzerbasis.
- Linux-Brew: bietet keinen Mehrwert über das direkte Binary oder das GHCR-Image.
- Podman als Host-Engine und VS Code/Codespaces:
  [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
  sagt sie für `docker` nicht zu.
- Änderungen am Lockerungsprofil der Sandbox: ein Messergebnis, das `--privileged` nicht trägt,
  führt über den Re-Evaluierungs-Trigger von [`ADR-0018`](../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
  zu einem neuen CR.

## 7. Closure-Notiz

Offen.
