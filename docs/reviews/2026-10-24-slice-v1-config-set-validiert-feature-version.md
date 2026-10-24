# Review: slice-v1-config-set-validiert-feature-version

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff a67acfc0..HEAD` (HEAD `bb6a401f`), Slice
  [`slice-v1-config-set-validiert-feature-version`](../plan/planning/in-progress/slice-v1-config-set-validiert-feature-version.md)
* Sensors: `make gates` selbst gefahren, Exit-Code 0 (sichtbar im Log, `GATES EXIT=0`).
  `make test-docker` nicht gelaufen: kein Docker-Pfad berührt. `make verify-depguard` nicht gelaufen:
  keine Importregel berührt (nur `internal/hexagon/application` ohne neue Importe).
* Gezielte Läufe: `docker run … u-boot:test-docker-tools go test ./internal/hexagon/application/ -run 'TestConfigSet|TestConfigSetMany|TestValidateDevcontainer|ForgeCLI'`.

## Findings

Kein HIGH, kein MEDIUM.

### LOW-1

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Plan-Tabelle gegen `git diff --stat`)
* `pfad`: Plan-Tabelle des Slice (§3) · „| `docs/reviews/2026-10-23-…` | neu |“
* `klasse`: plan-tabelle-abgleich
* `befund`: Der Diff berührt drei Dateien unter `done/` (zwei Slices, eine Welle) und zwei Evidence-Dateien des Beobachtungs-Registers (Verweis-Nachzug des Verschiebe-Commits `open/` → `in-progress/`); keine Zeile der Plan-Tabelle nennt sie. Der Review-Report-Name der Tabelle trägt das Datum 2026-10-23, der Lauf ist vom 2026-10-24.
* `risiko`: Die Tabelle gibt den Diff nicht vollständig wieder; der Abgleich nach der letzten Fix-Runde ist offen.
* `verifizierbar`: Review-only (`git diff --stat` gegen die Tabelle)

### INFO-1

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `spec/spezifikation.md` · „Eine `version` ohne `source` ist bei `gh` und `glab` in beiden Profilen“
* `klasse`: geltungsbereich-wortlaut
* `befund`: [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) nennt für `generate` nur `gh` und `glab`, [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), Handbuch und CHANGELOG nennen `ollama`, `gh`, `glab`. Der Code (`featureVersionNeedsSource`, `dockerfileBlock: true` für alle drei) trägt beide Wortlaute; der Abschnitt zu `gh`/`glab` ist durch seinen Titel begrenzt, kein Widerspruch.
* `risiko`: keines; ein Leser, der beide Abschnitte nebeneinander legt, sieht zwei Mengen.
* `verifizierbar`: Review-only

### INFO-2

* `kategorie`: INFO (Zuständigkeit Verifier)
* `quelle`: [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
* `pfad`: `internal/hexagon/application/config_features_test.go` · `TestConfigSet_DockerfileBlockVersionWithoutSource_Rejected`
* `klasse`: beleg-ebene
* `befund`: Exit-Code `10` und `--dry-run`/`--diff` sind auf Service-Ebene (`PreviewDryRun`, `PreviewAndApply`) und über die beiden Sentinel-Pins (`cli_test.go` „ErrConfigValueInvalid (validation)“ → 10, `config_internal_test.go` „ErrConfigValueInvalid“ → Diagnose-Code [`LH-FA-CONF-001`](../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)) belegt; ein Test durch den CLI-Befehl mit `--dry-run` existiert im Diff nicht.
* `risiko`: gering; Kette Service → Sentinel → Pin ist geschlossen.
* `verifizierbar`: `make test`

## Prüfung der Schwerpunkte

1. **Code.** `featureVersionNeedsSource` ist die einzige Aussage der Regel; `validateDevcontainerFeatures` fragt sie mit unverändertem Text und Sentinel (`ErrInvalidFeatureSource`), `revalidateFeatureEntry` im Fall `ConfigDevcontainerFeatureVersion`. `Set` (`config.go` Stufe 5) und `SetMany` (`applySetItem` → `revalidateConfigDomain`) laufen über denselben Aufruf, keine Doppelung. Fehlerklasse `ErrConfigValueInvalid` (nicht `ErrConfigSchemaInvalid`), keine neue Sentinel, beide Klassifikatoren unverändert und gepinnt (selbst gelesen). Meldung Englisch, nennt Pfad und Ausweg „source first“ (Sprachvertrag erfüllt). Reihenfolge bei `SetMany`: Prüfung nach jedem Paar, `version` vor `source` wird abgelehnt, atomar (Test prüft Bytes und Schreibzahl). Datei bei Ablehnung byte-identisch. Hexagon-/depguard: keine neuen Importe, Änderung nur in `application`.
   Randfälle: (a) `featureSources.allow` kann nur angehängt, `source` nicht geleert werden (leer wird abgewiesen) – ein Weg, `source` nachträglich zu entfernen und die Datei zu sperren, ist über `config set` nicht offen. (b) Ein No-Op-Set (gleicher Wert) umgeht die Revalidierung in `applySetItem`, trifft aber nur Dateien, die der Lesepfad schon ablehnt. (c) Ein bereits vorhandener Eintrag `version` ohne `source` bleibt gesperrt (Handbuch und CHANGELOG sagen es, Test `TestConfigSetMany_RejectsFileWithForgeVersionWithoutSource`). (d) Der Meldungstext „source first“ nennt die Freigabe in `featureSources.allow` nicht; der folgende `source`-Schritt liefert dafür den eigenen Hinweis. Kein Finding.
2. **Tests, Mutationen** (Rücknahme per `git apply -R` bzw. Gegenpatch; Baum vorher und nachher sauber, `git status --short` leer):
   * M1: gesamter Code-Diff in `config.go` zurückgenommen → `TestConfigSet_DockerfileBlockVersionWithoutSource_Rejected` und `TestConfigSetMany_DockerfileBlockVersionNeedsSourceFirst` rot.
   * M2: Sentinel `ErrConfigValueInvalid` → `ErrConfigSchemaInvalid` → beide rot.
   * M3: Helfer ohne `dockerfileBlock`-Bedingung → sieben Tests rot (u. a. `TestConfigSet_VersionWithSourceOrWithoutDockerfileBlock_Succeeds`, `TestValidateDevcontainerFeatures`).
   Die Umdrehung des Charakterisierungstests ist vollzogen; sein alter Name kommt nur noch in der Historie (`done/`, Reviews) vor.
3. **Doku.** [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Tabelle, Stufe 6) und [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz `version` und `source`): jede Bedingung hat einen Test (B1 bis B6 des Plans abgedeckt). §3.4: die Spec verweist nicht abwärts. Handbuch kennungsfrei (§5 Regel 6); der Teilsatz mit „Zurzeit“ ist entfernt. CHANGELOG `[Unreleased]` → `Fixed`, sachlich passend zum Verhalten; die Abschnitte ab `0.15.0` unberührt.
4. **Lastenheft** unverändert; [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) und [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) tragen die Zusage.
5. **Kommentare §3.7/§3.9.** Die neuen Doc-Kommentare (`featureVersionNeedsSource`, `revalidateFeatureEntry`) tragen Zusage bzw. Kopplung, keine Chronik. Der Wortlaut „zurzeit“/„Charakterisierung“ steht im Baum nur noch in `done/`-Dokumenten und Reviews als Historie, nicht in Handbuch, Spec oder Code. Plan-Tabelle: siehe LOW-1.
6. **Tag** `v0.15.0` zeigt weiter auf `a895ffd2`.

## Negativbefund

* Geprüft, ohne Befund: `internal/hexagon/application/config.go`, `devcontainer_features.go` (Regel einmal, Lesepfad unverändert), Dual-Classifier (kein neuer Sentinel, Pins vorhanden), Sprachvertrag der Meldung, Hexagon-/depguard-Grenzen, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `spec/lastenheft.md`, Tag `v0.15.0`.
* Nicht geprüft: `make test-docker` (kein Docker-Pfad), `make verify-depguard` (keine Importregel), `govulncheck`/`image-scan` (Abhängigkeiten unberührt), DoD und Closure (Verifier).
