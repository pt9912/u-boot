# homebrew-u-boot

Homebrew-Tap für [`u-boot`](https://github.com/pt9912/u-boot) (ADR 0016 im
Hauptrepo). Die Formel `Formula/u-boot.rb` wird **nicht von Hand gepflegt**: der
`publish`-Workflow des Hauptrepos hängt sie als Release-Asset `u-boot.rb` an jeden
stabilen Tag, und der Job `tap` zieht sie hierher nach (`scripts/tap-nachzug.sh`).

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot     # Homebrew 5 verlangt, Drittanbieter-Taps ausdrücklich zu vertrauen
brew install u-boot
u-boot --version
```

`brew tap pt9912/u-boot` bindet diesen Tap (`pt9912/homebrew-u-boot`) ein; Homebrew 5
installiert aus Drittanbieter-Taps erst nach `brew trust pt9912/u-boot`. Vier Plattformen: macOS und Linux, je `amd64` und `arm64`
(Homebrew trägt kein Windows). Prereleases (`vX.Y.Z-rc.1`) landen nicht im Tap.

Dieses Verzeichnis (`packaging/homebrew-tap/` im Hauptrepo) enthält die Dateien,
die beim Anlegen des Tap-Repos hineinkopiert werden: diese README und
`.github/workflows/smoke.yml`.
