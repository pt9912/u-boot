# u-boot Roadmap

**Status:** Aktiv. **Letzte Änderung:** 2026-07-24.

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**, keine
Reihenfolge von Terminen (Kurs Modul 6, vendored unter
`.harness/baseline/v3.5.1/regelwerk/modul-06-roadmap.md`). u-boots Wellen
entsprechen den **Release-Versionen**; Termine erscheinen nur als *Konsequenz*
einer abgeschlossenen Welle (Release-Datum), nicht als Treiber. Diese Datei
liegt dauerhaft in `in-progress/` und bleibt bewusst knapp: der per-Slice-
Audit-Trail lebt in den `done/`-Slices, die Release-Historie in
[`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md).

## Aktuelle Welle

**Welle-ID:** welle-harness-adoption-v3.5.1
**Start:** 2026-07-24
**Geplantes Ende:** offen (Wartungs-Welle; trigger-getriebene Folge-Slices)

**Closure-Trigger:** DoD des Adoptions-Slice
[`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md)
erfüllt (`make gates` grün, Review + Verification) **und** die Folge-Slices
FS-1..FS-4 (Reviewer-Skill, `docs/reviews/`, Sub-Area-Modus-Audit,
Freshness-Audit) sind entweder ausgeliefert oder als eigene `open/`-Slices
vertagt.

Die Produktarbeit wartet nach v0.4.0 (released 2026-06-08) auf externe Trigger
(siehe Nächste Wellen); die aktive Welle ist die Harness-/Baseline-Wartung.

## Nächste Wellen

| Welle | Trigger | Wichtigste Slices | Aufwand |
|---|---|---|---|
| macOS-Distribution | konkrete Homebrew-Nutzeranfrage | [`slice-v2-homebrew-formula`](../open/slice-v2-homebrew-formula.md) | S |
| Linux-Pakete | konkrete Debian-/RPM-Anfrage | [`slice-v2-distro-pakete`](../open/slice-v2-distro-pakete.md) ([ADR-0007](../../adr/0007-distributionswege-ghcr.md)) | M |
| Devcontainer-Robustheit | Real-World-Half-State-Beschwerde oder Schema-Erweiterung | [`slice-v2-generate-devcontainer-rollback-aware-write`](../open/slice-v2-generate-devcontainer-rollback-aware-write.md) | M |
| CI-Stabilität | belastbare Keycloak-Flake-Logs (Quay-/Mirror-Befund) | [`slice-v1-keycloak-ci-flake`](../open/slice-v1-keycloak-ci-flake.md) | S |
| Harness-Scaffolding (Produkt) | [ADR-0011](../../adr/0011-agent-harness-scaffolding.md) accepted + Spec-Erweiterung + Lizenz-Check | `slice-vN-harness-bootstrap-scaffold` (noch kein Plan) | L |
| Devcontainer-Egress-Firewall | [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) accepted + Spec-Erweiterung | `slice-vN-devcontainer-egress-firewall` (noch kein Plan) | M |
| Podman-first / Migration / Custom-Sources | Konkretisierung je Thema offen | `slice-vN-podman-formal`; `slice-later-migration` ([`LH-FA-CONF-006`](../../../../spec/lastenheft.md#lh-fa-conf-006--konfiguration-migrieren)); `slice-later-custom-data-sources` ([`LH-DA-004`](../../../../spec/lastenheft.md#lh-da-004--schema-migration)) | L |

## Meilensteine

<!-- Externe Versprechen / interne Trigger-Punkte; Datum als Konsequenz, nicht Treiber. -->

| Meilenstein | Welle(n) | Trigger / Fokus | Status |
|---|---|---|---|
| MVP-CLI + Release-Pipeline | v0.1.0 | Core-Subcommands, GHCR, CI/Gates | erreicht (2026-05-31) |
| Container-aware `doctor` + Binaries | v0.2.0 | Plattform-Binaries, Template-Katalog | erreicht (2026-06-01) |
| Add-on Catalogue Expansion | v0.3.0 | `remove`, `--with-deps`, Keycloak, OTel | erreicht (2026-06-01) |
| Maschinenlesbare CLI | v0.4.0 | `--json`/`--dry-run`/`--diff` alle 10 Subcommands, `logs`, Devcontainer-Features | erreicht (2026-06-08) |
| Nächstes Produkt-Release | v0.5.0 | trigger-getrieben (siehe Nächste Wellen) | offen |

## Abhängigkeitsgraph

```mermaid
flowchart LR
    V1[v0.1.0 MVP]
    V2[v0.2.0 doctor+binaries]
    V3[v0.3.0 add-ons]
    V4[v0.4.0 JSON-CLI]
    H[welle-harness-adoption-v3.5.1]
    N[Naechste Wellen: trigger-getrieben]

    V1 --> V2 --> V3 --> V4 --> H --> N
```

## Abgeschlossene Wellen

| Welle | Abschluss | Closure-Notiz (Detailquelle) |
|---|---|---|
| v0.1.0 (MVP M1..M8) | 2026-05-31 | [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md); [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md) |
| v0.2.0 | 2026-06-01 | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | 2026-06-01 | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| v0.4.0 (maschinenlesbare CLI, `logs`, Devcontainer-Features) | 2026-06-08 | [`slice-v1-release-cut-v0.4.0`](../done/slice-v1-release-cut-v0.4.0.md); JSON-CLI-Cluster [`slice-v1-cli-json-dry-run`](../done/slice-v1-cli-json-dry-run.md) (9/9 Folge-Slices + `T_close`, `3a35d58`) |
| Lokale FS-Templates | 2026-06 | [`slice-later-local-templates`](../done/slice-later-local-templates.md) ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates)) |

> **Closure-Form (Abweichung, MR-003).** u-boot führt **keine**
> `welle-NN-results.md`; die Welle-Closure lebt vollständig im jeweiligen
> `done/`-Release-Cut-Slice (Detailquelle inkl. Tranchen-Hashes).

## Historische Trigger-Verschiebungen

<!-- Umplanungen: Datum, Änderung, Grund. Steering-Loop-relevant. -->

| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-06-08 | JSON-CLI als 9-teiliger Cluster vor v0.4.0 gezogen | [ADR-0010](../../adr/0010-kein-http-driving-adapter.md) Re-Eval-Trigger 2; V1-pünktlich statt HTTP-Adapter |
| laufend | Debian/RPM + Homebrew auf konkrete Anfrage vertagt | Packaging-Overhead ohne belastbare Nachfrage ([ADR-0007](../../adr/0007-distributionswege-ghcr.md)) |

## Verwandte Dokumente

- [`carveouts.md`](carveouts.md) — Master-Inventar aller temporären und
  permanenten Carveouts ([`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)),
  plus Audit-Trail der Slices, die offene Carveouts geschlossen haben.
- [`README.md`](../README.md) — Slice-/Tranche-Konventionen für Dateinamen in
  `docs/plan/planning/` ([`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)).
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
- Diese Datei ist die einzige zulässige Ausnahme von der
  `slice-`/`tranche-`-Namenskonvention in `docs/plan/planning/` (siehe
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  und [`../README.md`](../README.md)).
