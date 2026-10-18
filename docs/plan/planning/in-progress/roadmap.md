# u-boot Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-10-18.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**, keine
Reihenfolge von Terminen (Kurs Modul 6, vendored unter
`.harness/baseline/v6.18.1/regelwerk/modul-06-roadmap.md`). u-boots Wellen
entsprechen den **Release-Versionen**; Termine erscheinen nur als *Konsequenz*
einer abgeschlossenen Welle (Release-Datum), nicht als Treiber. Diese Datei
liegt dauerhaft in `in-progress/` und bleibt bewusst knapp: der per-Slice-
Audit-Trail lebt in den `done/`-Slices, die Release-Historie in
[`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md).

## Offene Wellen

Derivativ: Der Zustand sind die flachen Welle-Dateien unter `docs/plan/planning/`; woran gearbeitet
wird, sagt das `Welle:`-Feld der Slices in `in-progress/`. Ziel, Trigger und Closure-Kriterien
stehen in der Welle-Datei, nicht hier.

In Arbeit: [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) (Ollama im Sandbox-Devcontainer messen und Quelle wählen, fünfter Slice von `welle-sandbox-egress-ollama`).

## Nächste Wellen

| Welle | Trigger | Wichtigste Slices | Aufwand |
|---|---|---|---|
| [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) (Egress im Sandbox-Profil komfortabler, Ollama als Katalog-Feature; Release-Zuordnung offen) | Beauftragung durch den Projektinhaber (2026-10-18, eingetreten); der erste Slice startet, wenn kein Slice in `in-progress/` liegt | [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md); [`slice-cr-sandbox-egress-anzeige`](../open/slice-cr-sandbox-egress-anzeige.md) und [`slice-v1-cli-egress-anzeige`](../open/slice-v1-cli-egress-anzeige.md); [`slice-v1-sandbox-egress-neuladen-entwurf`](../open/slice-v1-sandbox-egress-neuladen-entwurf.md); [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md), [`slice-cr-devcontainer-ollama`](../open/slice-cr-devcontainer-ollama.md), [`slice-v1-devcontainer-feature-ollama`](../open/slice-v1-devcontainer-feature-ollama.md); [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) | L |
| Harness-Scaffolding (Produkt) | [ADR-0011](../../adr/0011-agent-harness-scaffolding.md) accepted + Spec-Erweiterung + Lizenz-Check | `slice-vN-harness-bootstrap-scaffold` (noch kein Plan) | L |
| Podman-first / Migration / Custom-Sources | Konkretisierung je Thema offen | `slice-vN-podman-formal`; `slice-later-migration` ([`LH-FA-CONF-006`](../../../../spec/lastenheft.md#lh-fa-conf-006--konfiguration-migrieren)); `slice-later-custom-data-sources` ([`LH-DA-004`](../../../../spec/lastenheft.md#lh-da-004--schema-migration)) | L |

## Meilensteine

<!-- Externe Versprechen / interne Trigger-Punkte; Datum als Konsequenz, nicht Treiber. -->

| Meilenstein | Welle(n) | Trigger / Fokus | Status |
|---|---|---|---|
| MVP-CLI + Release-Pipeline | v0.1.0 | Core-Subcommands, GHCR, CI/Gates | erreicht (2026-05-31) |
| Container-aware `doctor` + Binaries | v0.2.0 | Plattform-Binaries, Template-Katalog | erreicht (2026-06-01) |
| Add-on Catalogue Expansion | v0.3.0 | `remove`, `--with-deps`, Keycloak, OTel | erreicht (2026-06-01) |
| Maschinenlesbare CLI | v0.4.0 | `--json`/`--dry-run`/`--diff` alle 10 Subcommands, `logs`, Devcontainer-Features | erreicht (2026-06-08) |
| Nächstes Produkt-Release | v0.5.0 | trigger-getrieben (siehe Nächste Wellen) | erreicht (2026-07-25) |

## Abhängigkeitsgraph

```mermaid
flowchart LR
    V1[v0.1.0 MVP]
    V2[v0.2.0 doctor+binaries]
    V3[v0.3.0 add-ons]
    V4[v0.4.0 JSON-CLI]
    H[welle-harness-adoption-v3.5.1]
    K[welle-harness-konformitaet-nachlauf]
    G[welle-gate-ausbau-v0.51]
    V5[v0.5.0 security+local-templates]
    N[Naechste Wellen: trigger-getrieben]

    V1 --> V2 --> V3 --> V4 --> H --> K --> G --> V5 --> N
```

## Abgeschlossene Wellen

| Welle | Abschluss | Closure-Notiz (Detailquelle) |
|---|---|---|
| v0.1.0 (MVP M1..M8) | 2026-05-31 | [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md); [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md) |
| v0.2.0 | 2026-06-01 | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | 2026-06-01 | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| v0.4.0 (maschinenlesbare CLI, `logs`, Devcontainer-Features) | 2026-06-08 | [`slice-v1-release-cut-v0.4.0`](../done/slice-v1-release-cut-v0.4.0.md); JSON-CLI-Cluster [`slice-v1-cli-json-dry-run`](../done/slice-v1-cli-json-dry-run.md) (9/9 Folge-Slices + `T_close`, `3a35d58`) |
| Lokale FS-Templates | 2026-06 | [`slice-later-local-templates`](../done/slice-later-local-templates.md) ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates)) |
| welle-harness-adoption-v3.5.1 (Regelwerk v3.5.1 + ADR-Format-CR auf MADR) | 2026-07-24 | [`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md) (`8d897d3`); [`slice-cr-adr-format-madr`](../done/slice-cr-adr-format-madr.md) (`7633cbb`) |
| welle-harness-konformitaet-nachlauf (FS-1..FS-4 + `architecture.md`-Template-Konformität) | 2026-07-25 | [`slice-harness-reviewer-skills-und-review-ablage`](../done/slice-harness-reviewer-skills-und-review-ablage.md) (`d5d896a`); [`slice-harness-sub-area-modus-audit`](../done/slice-harness-sub-area-modus-audit.md) (`be3d33c`); [`slice-harness-baseline-freshness-audit`](../done/slice-harness-baseline-freshness-audit.md) (`c7fe437`); [`slice-harness-architecture-template-konformitaet`](../done/slice-harness-architecture-template-konformitaet.md) (`c35249d`); Review-Findings (`4fed84f`), Report in [`docs/reviews/`](../../../reviews/README.md) |
| welle-spec-technik-stratum (Technik-Stratum `spezifikation.md`, Lastenheft 0.4.0, Architektur-Angleichung) | 2026-10-02 | [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) (`fdf4c92`); [`slice-v2-spec-schnitt-cli-init-devcontainer`](../done/slice-v2-spec-schnitt-cli-init-devcontainer.md) (`8119daa`); [`slice-v2-spec-schnitt-addons-betrieb`](../done/slice-v2-spec-schnitt-addons-betrieb.md) (`5d4716e`); [`slice-v2-spec-schnitt-projektkapitel`](../done/slice-v2-spec-schnitt-projektkapitel.md) (`d4abd2b`); [`slice-v2-spec-schnitt-qualitaet-akzeptanz-abschluss`](../done/slice-v2-spec-schnitt-qualitaet-akzeptanz-abschluss.md) (`3b98b30`); [`slice-v2-spec-architektur-angleichung`](../done/slice-v2-spec-architektur-angleichung.md) (`5661be3`); Review-Report `docs/reviews/2026-10-02-welle-spec-technik-stratum.md`; [`slice-v2-spec-lastenheft-coverage-default`](../done/slice-v2-spec-lastenheft-coverage-default.md) (durch die Kürzung erledigt); [`slice-v2-spec-lastenheft-verschlankung-runde-3`](../done/slice-v2-spec-lastenheft-verschlankung-runde-3.md) (`05212da`) |
| welle-lastenheft-verschlankung (Verhaltensdetails, Beispiele und Regelwerk-Dopplungen aus dem Lastenheft, 2984 → 1626 Zeilen) | 2026-10-02 | [`slice-v2-spec-lastenheft-verschlankung`](../done/slice-v2-spec-lastenheft-verschlankung.md) (`4cd66a5`); [`slice-v2-spec-lastenheft-verschlankung-runde-2`](../done/slice-v2-spec-lastenheft-verschlankung-runde-2.md) (`d024a1b`) |
| welle-gate-ausbau-v0.51 (d-check-Bump + vier Sensoren) | 2026-07-25 | [`slice-gate-print-mk-einbindung`](../done/slice-gate-print-mk-einbindung.md) (`e246bbf`); [`slice-gate-planning-targets-module`](../done/slice-gate-planning-targets-module.md) (`83efe46`); [`slice-gate-ids-link-policy-always`](../done/slice-gate-ids-link-policy-always.md) (`0d896fc`); [`slice-gate-immutabilitaets-sensor`](../done/slice-gate-immutabilitaets-sensor.md) (`7d9c4c2`); Vorlauf: Image-Bump (`225627c`) |
| v0.5.0 (Security-Release + lokale Templates) | 2026-07-25 | [`slice-v1-release-cut-v0.5.0`](../done/slice-v1-release-cut-v0.5.0.md) — Tag `v0.5.0` auf `1bc071f` |
| Spec-Wartung: Traceability-Matrix konfiguriert | 2026-09-29 | [`slice-gate-rtm-traceability`](../done/slice-gate-rtm-traceability.md) — RTM 79 Anforderungen, 0 Waisen |
| Spec-Wartung: Sandbox-Profil (Lastenheft 0.3.0) | 2026-09-30 | [`slice-cr-sandbox-devcontainer`](../done/slice-cr-sandbox-devcontainer.md) (`712af41`) |
| Sandbox-Profil Umsetzung V1 (Lastenheft 0.3.2) | 2026-09-30 | [`slice-v1-sandbox-devcontainer-umsetzung`](../done/slice-v1-sandbox-devcontainer-umsetzung.md) (`50dbd84`) |
| Sandbox: Egress-Restriktion (Lastenheft 0.3.3) | 2026-09-30 | [`slice-v2-sandbox-egress`](../done/slice-v2-sandbox-egress.md) (`9b1c414`) |
| Config-Erweiterungen (V1-Cleanup) | 2026-10-01 | [`slice-v1-config-structured-hint`](../done/slice-v1-config-structured-hint.md), [`slice-v1-config-multi-path-get`](../done/slice-v1-config-multi-path-get.md), [`slice-v1-config-multi-path-set`](../done/slice-v1-config-multi-path-set.md), [`slice-v1-config-list-subcommand`](../done/slice-v1-config-list-subcommand.md) (`acc0233`) |
| Logs-Erweiterungen (V1-Cleanup) | 2026-10-01 | [`slice-v1-logs-multi-service-filter`](../done/slice-v1-logs-multi-service-filter.md), [`slice-v1-logs-format-flags`](../done/slice-v1-logs-format-flags.md), [`slice-v1-logs-time-range-filter`](../done/slice-v1-logs-time-range-filter.md) (`4db9f90`) |
| Up/Down-Erweiterungen (V1-Cleanup) | 2026-10-01 | [`slice-v1-recreate-detection`](../done/slice-v1-recreate-detection.md), [`slice-v1-multi-port-services`](../done/slice-v1-multi-port-services.md), [`slice-v1-up-partial-snapshot-on-failure`](../done/slice-v1-up-partial-snapshot-on-failure.md), [`slice-v1-down-volumes-named-list`](../done/slice-v1-down-volumes-named-list.md), [`slice-v1-volume-auto-removal`](../done/slice-v1-volume-auto-removal.md) (`c4c9b09`) |
| Generate-Devcontainer-Rollback (V2-Hardening) | 2026-10-01 | [`slice-v2-generate-devcontainer-rollback-aware-write`](../done/slice-v2-generate-devcontainer-rollback-aware-write.md) (`2f12ec6`) |
| v0.6.0 (Sandbox-Devcontainer + V1-Cleanup + Homebrew) | 2026-10-01 | [`slice-v1-release-cut-v0.6.0`](../done/slice-v1-release-cut-v0.6.0.md) — Tag `v0.6.0` auf `76a8200` |
| Spec-Wartung und Umsetzung: Docker als verschachtelte Runtime (Lastenheft 0.5.0, `nestedRuntime: docker`) | 2026-10-14 | [`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md) (`6d726e0`); [`slice-v1-sandbox-docker-generator`](../done/slice-v1-sandbox-docker-generator.md); [`slice-v1-sandbox-docker-degradation-doctor`](../done/slice-v1-sandbox-docker-degradation-doctor.md); [`slice-v1-sandbox-docker-integrationstest`](../done/slice-v1-sandbox-docker-integrationstest.md) (letzte Code-Änderung `e1c2c3a`); Review-Reports `docs/reviews/2026-10-14-sandbox-docker-runtime.md`, `-r2` |
| v0.8.0 (Docker im Sandbox-Devcontainer, `nestedRuntime: docker`) | 2026-10-14 | [`slice-v1-release-cut-v0.8.0`](../done/slice-v1-release-cut-v0.8.0.md) — Tag `v0.8.0` auf `1f2b41b` |
| Docker aus dem Docker-APT-Repo (`nestedRuntime: docker`, [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md)) | 2026-10-14 | [`slice-v1-sandbox-docker-ce`](../done/slice-v1-sandbox-docker-ce.md) (`90f9122`); Review-Reports `docs/reviews/2026-10-14-sandbox-docker-ce.md`, `-r2` |
| v0.8.1 (aktueller Docker im Sandbox-Devcontainer) | 2026-10-14 | [`slice-v1-release-cut-v0.8.1`](../done/slice-v1-release-cut-v0.8.1.md) — Tag `v0.8.1` auf `5aadaf3` |
| Devcontainer-Features `homebrew` und `sdkman` ([ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md)) | 2026-10-15 | [`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md) (`ac59b83`); Review-Report `docs/reviews/2026-10-15-devcontainer-features-homebrew-sdkman.md` |
| Devcontainer-Feature `python` | 2026-10-15 | [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md); Review-Report `docs/reviews/2026-10-15-devcontainer-feature-python.md` |
| v0.9.0 (Devcontainer-Features Homebrew, SDKMAN, Python) | 2026-10-15 | [`slice-v1-release-cut-v0.9.0`](../done/slice-v1-release-cut-v0.9.0.md) — Tag `v0.9.0` auf `7bfd8d6` |
| welle-macos-verifikation (Homebrew-Formel auf macOS belegt; Docker-Desktop-Zusage aus dem Lastenheft 0.7.0 genommen) | 2026-10-16 | [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) (Detailquelle); Welle-Datei [`welle-macos-verifikation`](../done/welle-macos-verifikation.md) §7 |
| welle-harness-wartung (Reviewer-Skill auf Baseline-Schema; `internal/`-README-Festlegungen in der Spezifikation) | 2026-10-16 | [`slice-harness-internal-readmes-festlegungen-in-spec`](../done/slice-harness-internal-readmes-festlegungen-in-spec.md) (Detailquelle); [`slice-harness-dcheck-bump-reviewer-skill`](../done/slice-harness-dcheck-bump-reviewer-skill.md); Welle-Datei [`welle-harness-wartung`](../done/welle-harness-wartung.md) §7 |
| welle-cli-verbosity (Verbosity-Vertrag: `--quiet` nur Fehler mit vier Ausnahmen, Vorschau bleibt, letzte Option gewinnt; Lastenheft 0.8.0 bis 0.10.0, [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)) | 2026-10-16 | [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Detailquelle); [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md); [`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md); [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md); [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md); [`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md); Welle-Datei [`welle-cli-verbosity`](../done/welle-cli-verbosity.md) §7 |
| welle-devcontainer-sshd-claude-code (Devcontainer-Features `sshd` und Claude Code im Katalog, Sicherheitswarnungen unter `--quiet`, `podman build` im Sandbox-Devcontainer; Lastenheft 0.11.0 und 0.12.0, [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)) | 2026-10-17 | [`slice-v1-devcontainer-podman-build`](../done/slice-v1-devcontainer-podman-build.md) (Detailquelle); [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md); [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md); [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md); [`slice-v1-devcontainer-feature-claude-code`](../done/slice-v1-devcontainer-feature-claude-code.md); [`slice-v2-spec-devcontainer-sshd-claude-code`](../done/slice-v2-spec-devcontainer-sshd-claude-code.md); [`slice-v1-cli-egress-ohne-profil-sichtbar`](../done/slice-v1-cli-egress-ohne-profil-sichtbar.md); Validator-Bericht [`2026-10-17-validation-welle-devcontainer-sshd-claude-code`](../../../reviews/2026-10-17-validation-welle-devcontainer-sshd-claude-code.md); Welle-Datei [`welle-devcontainer-sshd-claude-code`](../done/welle-devcontainer-sshd-claude-code.md) §7 |
| v0.10.0 (Devcontainer-Features `sshd` und Claude Code, `podman build`, Verbosity-Vertrag; trägt [`welle-devcontainer-sshd-claude-code`](../done/welle-devcontainer-sshd-claude-code.md), [`welle-cli-verbosity`](../done/welle-cli-verbosity.md), [`welle-macos-verifikation`](../done/welle-macos-verifikation.md) und [`welle-harness-wartung`](../done/welle-harness-wartung.md)) | 2026-10-17 | [`slice-v1-release-cut-v0.10.0`](../done/slice-v1-release-cut-v0.10.0.md) — Tag `v0.10.0` auf `f50690c` |
| v0.11.0 (Fallback-Ausweisung vom Host, fehlendes `/dev/fuse` als Startfehler, `doctor` im Container, Dry-Run „Would“, Podman 6.1.4, `sshd`-Anleitung, e2e-Helfer; trägt die Slices seit `v0.10.0`) | 2026-10-18 | [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) — Tag `v0.11.0` auf `6926ad6` |

> **Closure-Form (Abweichung, MR-003).** u-boot führt **keine**
> `welle-<Kennung>-results.md`. Die Closure bisheriger Wellen lebt vollständig im jeweiligen
> `done/`-Release-Cut-Slice (Detailquelle inkl. Tranchen-Hashes); bei neuen Wellen steht die
> Closure-Notiz in §7 der Welle-Datei und verweist auf diesen Slice.

## Historische Trigger-Verschiebungen

<!-- Umplanungen: Datum, Änderung, Grund. Steering-Loop-relevant. -->

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-08 | JSON-CLI als 9-teiliger Cluster vor v0.4.0 gezogen | [ADR-0010](../../adr/0010-kein-http-driving-adapter.md) Re-Eval-Trigger 2; V1-pünktlich statt HTTP-Adapter |
| laufend | Debian/RPM + Homebrew auf konkrete Anfrage vertagt | Packaging-Overhead ohne belastbare Nachfrage ([ADR-0007](../../adr/0007-distributionswege-ghcr.md)) |

## Verwandte Dokumente

- [`carveouts.md`](carveouts.md) — Master-Inventar aller temporären und
  permanenten Carveouts (Baseline-Regelwerk `modul-07-carveouts.md`),
  plus Audit-Trail der Slices, die offene Carveouts geschlossen haben.
- [`README.md`](../README.md) — Slice-/Tranche-Konventionen für Dateinamen in
  `docs/plan/planning/` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).
- [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md)
  — ausgelagerte Release-Historie.

## Pflege-Regeln

- Diese Roadmap beschreibt Wellen-Steuerung (Reihenfolge, Trigger,
  Abhängigkeiten), nicht jede historische Tranche.
- Release-Details, lange Commit-Listen und retrospektive Tabellen gehören in
  `done/`-Slices oder nach `docs/archive/`.
- Neue Wellen brauchen einen beobachtbaren Trigger und die wichtigsten Slices.
  Ohne ausgearbeiteten Plan bleibt eine Welle als benannter Eintrag hier, nicht
  als halbfertiger Slice.
- Diese Datei ist ein Master-Dokument und folgt keiner der Namenskonventionen in
  `docs/plan/planning/` (siehe [`../README.md`](../README.md)).
- Eine Welle bekommt ihre Datei (`welle-<kebab-slug>.md`) bei der Eröffnung; geplante
  Wellen stehen nur unter *Nächste Wellen*. Schließt eine Welle, wandert die Datei nach `done/`,
  die Zeile unter *Abgeschlossene Wellen* zeigt auf den Release-Cut-Slice als Detailquelle.
