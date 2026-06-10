# Slice Later: Lokale User-Templates (`u-boot init --template ./pfad`, [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates))

> **Status: DONE** — `u-boot init --template ./pfad` ausgeliefert
> ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates), `Priorität: Later`; auf Nutzer-Wunsch vor
> regulärem Trigger umgesetzt). DoD-Hashes: T1 `66c347d`, T2
> `87a8704`, T3 `5031b5f`, T4 `adaafbe`, Pre-T5-Review `7d63532`,
> T5-Closure `4fdf1a6`. Funktional, getestet (Unit + FS-Resolver +
> End-to-End-Exit-Code-Matrix + `--json`-Envelope + CLI-Acceptance),
> dokumentiert (README EN+DE, CHANGELOG, `cli-json-output.md` §6.4,
> [ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte ✅). Der [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)-Carveout ist
> mit dieser Closure aufgelöst → von §Temporäre Carveouts in
> [`carveouts.md`](../in-progress/carveouts.md) §Carveout-Auflösungs-
> Slices (historisch) gewandert. `--var`-Variablen bleiben
> out-of-scope (eigener Folge-Slice).
>
> **Post-Closure-Review #1 (`d03d73d`):** Zweite Review-Runde
> (Konvergenz-Test) fand einen MEDIUM, den Runde 1 übersah —
> `ClassifyTemplateRef` stufte bare `.`/`..` als Katalog-Namen ein
> (statt FS-Pfad), wodurch `init --template .` auf die Katalog-Wurzel
> und `..` via `path.Join` auf die embed-FS-Wurzel auflöste. Gefixt:
> `.`/`..` → `TemplateRefPath` **plus** Defense-in-Depth-Guard in
> `catalog.Open` (pfad-förmige Namen → `ErrTemplateNotFound`). Tests
> erweitert. `make gates` grün.

## Auslöser

`spec/lastenheft.md` [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates) (Later) fordert lokale
User-Templates:

```bash
u-boot init --template ./my-template
```

[ADR-0009](../../adr/0009-template-format-yaml-files.md) §Entscheidung hat den Pfad bereits vorgezeichnet:

> *„Lokale User-Templates ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates), Later): `--template
> ./mein-template` löst gegen das Dateisystem statt `embed.FS` auf.
> Same Schema, same Engine."*

Die zwei Voraussetzungs-Slices sind `done`:

- [`slice-v1-template-list`](../done/slice-v1-template-list.md) —
  `template.yaml`-Schema (`apiVersion`-Gate, `KnownFields`),
  `domain.TemplateMetadata` + `Catalog`-Adapter.
- [`slice-v1-template-init`](../done/slice-v1-template-init.md) —
  Render-Loop (`TemplateInitService`), `driven.TemplateFiles`-Port
  (`Open(ctx, name) (iofs.FS, error)`), `domain.TemplatePath`-
  Eskalations-Validator, drei Driving-Sentinels
  (`ErrTemplateNotFound`/`ErrInvalidTemplatePath`/`ErrTemplateRender`).

Der Render-Pfad ist also vollständig FS-agnostisch: `TemplateInitService`
hängt nur an der `driven.TemplateFiles`-Abstraktion. Dieser Slice
liefert eine **zweite `TemplateFiles`-Implementierung**, die gegen
das reale Dateisystem statt `embed.FS` auflöst, plus einen
**Composite-`TemplateFiles`-Adapter**, der via reiner
`domain`-Klassifikation zwischen Katalog-Name und lokalem Pfad
unterscheidet (T0-(a)/(a2)). Die CLI bleibt unverändert und reicht
den `--template`-String roh durch — kein Business-Dispatch im CLI.

## Trigger

Konkrete Nutzeranfrage nach projekteigenen/lokalen Templates (z. B.
ein Team, das ein firmeninternes Service-Template versionieren und
per `u-boot init --template ./templates/our-service` ausrollen will).
Bis dahin bleibt der eingebaute Katalog (`basic`, …) die einzige
Template-Quelle.

## Aufhebungsbedingung

`u-boot init --template ./my-template my-service` in einem leeren
Verzeichnis rendert ein lokales Template-Verzeichnis (mit gültiger
`template.yaml` + `*.tmpl`/Plain-Files) byte-identisch zum
erwarteten Output ins Projekt; ein Pfad ohne `template.yaml` und ein
nicht existierender Pfad werden mit Exit 10 abgewiesen; ein
malformed `template.yaml` wird mit Exit 10 abgewiesen (T0-(d)); ein
Template mit einem **Symlink** im Datei-Baum, das aus dem Root zeigt,
wird deterministisch mit Exit 10 abgewiesen (T0-(e), primärer
Escape-Vektor). `generatedFiles` wird dabei **nicht** als
Render-Vertrag geprüft (T0-(g)).

## Architektur-Anker (aus [ADR-0009](../../adr/0009-template-format-yaml-files.md) + done-Slices)

| Baustein | Status heute | Dieser Slice |
| --- | --- | --- |
| `driven.TemplateFiles.Open(ctx, name)` | Port existiert, eine Impl (`externaltemplates.Catalog`, `embed.FS`); Docs sprechen noch von Katalog-Name | **zweite Impl** (FS-Resolver) + **Composite**, der via `domain`-Klassifikation an Katalog oder FS delegiert; Port-/Request-Docs werden auf rohen Template-Ref (`name` oder Pfad) präzisiert |
| `application.TemplateInitService` | mappt alle `Open`-Fehler ≠ `ErrTemplateNotFound` auf `ErrTemplateRender` (Exit 14); Walk-Callback mappt alle Walk-Fehler auf `ErrTemplateRender` (Exit 14) | **zwei** umrissene Änderungen: (1) `Open`-Mapping-Branch `ErrTemplateInvalid` → Exit 10 (T0-(d)); (2) Symlink-Guard im Walk-Callback → `ErrInvalidTemplatePath` Exit 10 (T0-(e)). Sonst unverändert |
| `domain.TemplatePath` | validiert Output-Pfade gerenderter Dateien (greift im `renderOne`) | **unverändert** (defense-in-depth, T0-(g)) |
| `domain` Name-vs-Pfad-Klassifikator | existiert nicht | **neu** als reine Funktion (T0-(a)) |
| `readTemplate` (apiVersion-Gate, `KnownFields`, `Validate`) | privat in `externaltemplates/catalog.go:194` | **teilen/extrahieren** in ein kleines Driven-Helper-Paket (z. B. `internal/adapter/driven/templateyaml`), damit `externaltemplates` und `localtemplates` denselben Parser nutzen ohne Adapter-zu-Adapter-Import |
| `driven`-Sentinels | `ErrTemplateNotFound` | **neu** `ErrTemplateInvalid` (malformed `template.yaml`, Exit 10) |
| CLI `--template`-Flag | reicht Wert als `TemplateName` durch | **unverändert** (kein Business-Dispatch im CLI) |
| `main.go`-Wiring | ein `templateCatalogAdapter` | **ein Composite-Resolver** (Katalog + FS) verdrahtet |

## T0 — Decisions (vor Code festzurren)

| # | Frage | Empfehlung |
| - | ----- | ---------- |
| (a) | **Name-vs-Pfad-Disambiguierung — Regel:** wann ist `--template X` ein Katalog-Name, wann ein lokaler Pfad? | Pfad, wenn `X` exakt `.`/`..` ist (Post-Closure-Review #1), mit `./`, `../` oder `/` beginnt, exakt `~` ist, mit `~/` beginnt, einen Slash `/` oder Backslash `\` enthält, oder wie ein Windows-Drive-Pfad (`C:\...`/`C:/...`) aussieht; sonst Katalog-Name. Die Regel ist **plattformunabhängig deterministisch** (nicht `filepath.Separator`-abhängig), damit Linux/macOS/Windows-Binaries gleich klassifizieren. Kein FS-Stat zur Klassifikation (vermeidet TOCTOU + überraschende Reklassifikation bei gleichnamigem Verzeichnis). Die Regel selbst ist als **reine `domain`-Funktion** (z. B. `domain.ClassifyTemplateRef(s)`) zu implementieren — pure, FS-frei, depguard-konform, gezielt unit-testbar. `~`-Expansion ist keine Domain-Aufgabe: der lokale FS-Resolver expandiert nur `~`/`~/...` via `os.UserHomeDir`; `~user` bleibt unsupported und wird nicht als Home-Alias behandelt. |
| (a2) | **Name-vs-Pfad-Disambiguierung — Ort:** wo wird Katalog- vs. FS-Resolver gewählt? | **Composite-`TemplateFiles`-Adapter** in der Driven-Schicht (die Adapter-Schicht darf beide Resolver kennen, [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)). Der Composite ruft die `domain`-Klassifikation (a) und delegiert an Katalog- oder FS-Resolver. **Kein** Business-Dispatch im CLI (CLI darf keinen Driven-Adapter importieren) und **kein** Branch im Wiring: `main.go` verdrahtet genau einen Composite. Die CLI reicht den rohen `--template`-String unverändert als `TemplateName` durch (heute schon so). |
| (b) | **`--var key=value`-Auflösung:** in-scope? | **Out of scope** — analog [`slice-v1-template-init`](slice-v1-template-init.md) (nur `templateData{Name}`-Projektion). Eigener Folge-Slice, sobald ein variable-bedürftiges Template real auftaucht. Hier nur dokumentieren, dass `variables:` im lokalen `template.yaml` aktuell ignoriert/gewarnt wird. |
| (c) | **Metadaten-Validierung lokaler Templates:** wie? | `readTemplate`-Logik (apiVersion-Gate + `KnownFields` + `Validate()`) in einen teilbaren Driven-Helfer extrahieren (T1; z. B. `adapter/driven/templateyaml.Read`) und im FS-Resolver wiederverwenden — ein fehlendes/ungültiges/unsupported `template.yaml` muss am Root genauso hart abweisen wie im Katalog (fail-fast). Kein Export aus `externaltemplates`, damit `localtemplates` nicht an einen konkreten Geschwister-Adapter koppelt. Fehlerklasse siehe (d). |
| (d) | **Fehlerklassen-Split — not-found vs. invalid:** welcher Exit-Code? | **Zwei getrennte Klassen, beide Exit 10:** (1) Root existiert nicht / ist kein Verzeichnis / `template.yaml` fehlt → `driven.ErrTemplateNotFound` (bestehend, mappt via `driving.ErrTemplateNotFound` auf Exit 10). (2) `template.yaml` vorhanden aber malformed/unsupported-apiVersion/Validate-Fail → **neuer** `driven.ErrTemplateInvalid`-Sentinel. **Wichtig:** der heutige `TemplateInitService.Init` mappt *alle* `Open`-Fehler außer `driven.ErrTemplateNotFound` auf `driving.ErrTemplateRender` (Exit **14**). Ohne Anpassung würde ein malformed lokales Template fälschlich als technischer Render-Fehler (14) statt user-fixbarer Validierungsfehler (10) klassifiziert. Deshalb: `Init` bekommt **eine** zusätzliche `errors.Is(err, driven.ErrTemplateInvalid)`-Branch → neuer `driving.ErrTemplateInvalid` (Exit 10). Dies ist die `Open`-Error-Mapping-Änderung; die zweite, separate Service-Änderung ist der Symlink-Guard im Walk-Callback (T0-(e)). Mehr ändert sich am Service nicht. |
| (e) | **Symlinks im Template-Baum:** Verhalten + Sentinel + **Ort** + Exit-Code? | **Hart ablehnen, nicht still überspringen.** Stilles Skippen würde unvollständigen Output erzeugen, der trotzdem grün wirkt — kein deterministischer Produktvertrag. **Ort:** der Check sitzt im **Application-Render-Loop** (`planRender`-`WalkDir`-Callback in `TemplateInitService`), **nicht** im Driven-Resolver — nur dort ist der `iofs.DirEntry`-Typ sichtbar *und* der `driving`-Sentinel importierbar (ein Driven-Adapter darf `driving` nicht importieren, [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)). Sobald der Walk-Callback einen Symlink-Eintrag (`d.Type()&fs.ModeSymlink != 0`) sieht, bricht der **gesamte** Render ab mit `driving.ErrInvalidTemplatePath` (Exit **10**, semantisch korrekt: Pfad-Safety-Verletzung, bestehender Sentinel — kein neuer nötig). **Wichtig:** heute mappt der Walk-Callback alle Walk-Layer-Fehler auf `driving.ErrTemplateRender` (Exit 14) — der Symlink-Guard ist deshalb eine **bewusste Render-Loop-Änderung** (eine zusätzliche, früh greifende Type-Prüfung vor dem `ErrTemplateRender`-Fallback), kein „unverändert". Die Prüfung greift im gemeinsamen Loop für **beide** Quellen; `embed.FS` liefert nie Symlinks → für den Katalog-Pfad ein harmloser No-Op (defense-in-depth). Kein Folgen des Targets, keine Target-Außerhalb-Root-Prüfung (pauschale Ablehnung ist die einfachere, sichere Variante). Table-Test pinnt: Symlink im Baum → `driving.ErrInvalidTemplatePath`/Exit 10, kein Teil-Output (Two-Phase-Render: Reject in Phase 1 vor jedem `WriteFile`). |
| (f) | **Root-Pfad selbst:** absolute Pfade / `..` erlaubt? | Ja — der Root ist vom User explizit gewählt (anders als die Datei-Liste *im* Template). `domain.NewTemplatePath` gilt weiterhin nur für die **Output-Pfade** der gerenderten Dateien (Schutz der Projekt-Base-Dir bleibt erhalten). |
| (g) | **`generatedFiles`-Semantik:** validierter Render-Vertrag oder Anzeige-Metadatum? | **Nur Anzeige-Metadatum** (für `template list`-Surface), **kein** validierter Render-Vertrag. Der Render-Pfad läuft über `iofs.WalkDir` (nicht über `generatedFiles`); `WalkDir` liefert per `io/fs`-Kontrakt ausschließlich *clean, relative* Pfade — ein `..`/absoluter Pfad kann darüber gar nicht entstehen. Die echte Output-Pfad-Safety kommt aus: (1) Symlink-Policy (e), (2) `domain.NewTemplatePath` auf jedem gerenderten Output-Pfad (defense-in-depth, existiert bereits in `renderOne`). `generatedFiles`-Inhalt wird **nicht** gegen das Dateisystem geprüft. |

## Tranchen (vorgeschlagen, wird beim Trigger ausgearbeitet)

| T | Inhalt (Skizze) |
| - | --------------- |
| T0 | Decisions (a)–(g) festzurren; Out-of-Scope-Carveout für `--var` als Stub in `open/` falls nötig. |
| T1 ✅ | **Fundament (pure/infra + Service) — done (`66c347d`):** (1) `templateyaml`-Paket: `Read`/apiVersion-Gate/`KnownFields`/`Validate` aus `externaltemplates/catalog.go` extrahiert, Katalog-Pfad unverändert grün; (2) `domain.ClassifyTemplateRef` reine Funktion + Tabelle (`.`, `..`, `./`, `../`, `/abs`, `~`, `~/x`, `foo/bar`, `foo\bar`, `C:\x`, `c:tpl`, `basic`, `~user`, …; `.`/`..` per Post-Closure-Review #1 ergänzt), 100% Coverage; (3) Port-/Request-Docs auf rohen Template-Ref präzisiert (Feldname `TemplateName` beibehalten); (4) Sentinels `driven.ErrTemplateInvalid` + `driving.ErrTemplateInvalid`; (5) `TemplateInitService` zwei umrissene Änderungen: `Open`-Mapping-Branch `ErrTemplateInvalid`→Exit 10 (d) **und** Symlink-Guard im `planRender`-Walk-Callback → `driving.ErrInvalidTemplatePath` Exit 10 (e), Two-Phase-No-Write-Pin; Dual-Klassifikator (`isTemplateInitValidationError` + `mapInitErrorToDiagnostic` [`LH-FA-TPL-002`](../../../../spec/lastenheft.md#lh-fa-tpl-002-template-metadaten)) + ExitCode-Pin-Test. `make gates` grün (lint 0, coverage 91.40%). |
| T2 ✅ | **FS-Resolver + Composite (Driven-Adapter) — done (`87a8704`):** `internal/adapter/driven/localtemplates/` mit `Resolver.Open(ctx, path)` über stdlib `os`/`io/fs` (`os.Stat` + `os.DirFS`), Root-Existenz/Verzeichnis-Check, `~`/`~/…`-Expansion via `os.UserHomeDir`, `template.yaml`-Gate via `templateyaml.Read` mit `ErrNotExist`-Split (→ `ErrTemplateNotFound` bzw. `ErrTemplateInvalid`). Kein Import von `externaltemplates`/`driven/fs`. **Symlink-Policy nicht hier** — Resolver liefert den gerooteten `iofs.FS` unverfolgt; Walk-Guard sitzt in T1. `Composite.Open` dispatcht via `domain.ClassifyTemplateRef` (raw-passthrough). Test-Seam-frei: TempDir + `t.Setenv("HOME")`-Fixtures (valid/notexist/notdir/missing/malformed/`~`/`~/…`/HOME-unresolvable/ctx-cancel) + Composite-Dispatch-Tabelle. `make gates` grün (coverage 91.40%). |
| T3 ✅ | **Wiring + End-to-End — done (`5031b5f`):** `main.go` verdrahtet `localtemplates.NewComposite(externaltemplates.New(), localtemplates.New())` als einzigen `TemplateFiles` für `NewTemplateInitService`; `templateCatalogAdapter` bleibt für die `template list`-`TemplateCatalog`-Rolle. CLI unverändert (kein Business-Dispatch, `--template` roh durchgereicht). End-to-End-Test (`cli_test`, echte Adapter-Kette): Exit-Code-Matrix not-found 10 / invalid-metadata 10 / symlink 10 / render-IO 14 + local-valid byte-identisch + Katalog-`basic`-Regressionspin. Mutex-Regeln unverändert (erblich aus [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)). **Feature ab hier funktional**; Docs/Acceptance in T4. `make gates` grün. |
| T4 ✅ | **E2E + Docs — done (`adaafbe`):** CLI-Command-Acceptance (`cli_test`, voll verdrahteter `init`-Pfad via `getwd`-Seam): lokales Fixture → Projekt byte-identisch + INIT-003-Strukturdirs, Missing-Path → Exit 10. README EN+DE um `--template ./pfad`-Snippet + `<name\|pfad>`-Referenz ergänzt; CHANGELOG `[Unreleased] ### Added`; [ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte auf ✅. `make gates` grün. |
| T5 ✅ | **Closure — done (`4fdf1a6`):** Slice `git mv` `open/` → `done/`; Status-Block + DoD-Hash-Tabelle finalisiert; roadmap.md-Zeile auf done ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates) ausgeliefert); `carveouts.md` [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates)-Zeile von §Temporäre Carveouts → §Carveout-Auflösungs-Slices (historisch); `open/README.md`-Eintrag entfernt; [ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte-Link auf `done/`-Pfad. `make gates` grün. |

## Akzeptanzkriterien

- ✅ `u-boot init --template ./fixture my-svc` rendert ein lokales
  Template byte-identisch (`.tmpl` gerendert, Plain-Files 1:1
  kopiert, `template.yaml` nicht ins Projekt geleakt).
- ✅ Katalog-Pfad (`--template basic`) bleibt unverändert grün
  (keine Regression durch die Dispatch-Einführung).
- ✅ Lokaler Pfad nicht existent / kein Verzeichnis / `template.yaml`
  fehlt → Exit-Code 10 via `driving.ErrTemplateNotFound`.
- ✅ Lokales `template.yaml` malformed / unsupported-`apiVersion` /
  `Validate`-Fail → Exit-Code **10** via `driving.ErrTemplateInvalid`
  (Regressions-Pin gegen die heutige Exit-14-Fehlklassifikation,
  T0-(d)). Test: `u-boot init --template ./bad-metadata my-svc`.
- ✅ Symlink im lokalen Template-Baum → gesamter Render mit Exit 10
  (`driving.ErrInvalidTemplatePath`) abgewiesen, **kein** Teil-Output
  geschrieben (T0-(e), primärer Escape-Vektor) — durch Table-Test
  gepinnt.
- ✅ Defense-in-depth: ein gerenderter Output-Pfad, der nach
  `.tmpl`-Strip `domain.NewTemplatePath` verletzt → `ErrInvalidTemplatePath`
  (Exit 10). `generatedFiles`-Metadatum wird **nicht** als
  Render-Vertrag geprüft (T0-(g)).
- ✅ `make gates` grün; depguard-Regeln (`domain-isoliert`,
  `application-no-adapter`) eingehalten.

## Out of Scope

- **`--var key=value`-Variable-Resolution + Prompt-UI** — eigener
  Folge-Slice (T0-(b)); hier nur `templateData{Name}`.
- **Lokale Templates in `u-boot template list`** — Listing bleibt
  Katalog-only; lokale Templates sind pfad-adressiert, nicht
  enumeriert.
- **Remote-/Git-/OCI-Templates** — [ADR-0009](../../adr/0009-template-format-yaml-files.md) verwirft OCI explizit
  als prospektiv ohne Trigger; Remote-Fetch wäre ein eigener Slice
  mit [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte)-Sandbox-Betrachtung.
- **Cookiecutter-Kompatibilität** — durch [ADR-0009](../../adr/0009-template-format-yaml-files.md) verworfen.
- **`generatedFiles` als validierter Render-Vertrag** — bleibt
  Anzeige-Metadatum (T0-(g)); ein Abgleich „deklarierte vs.
  tatsächlich gerenderte Files" wäre ein eigener Slice.

## Bezug

- Spec: [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates) (`spec/lastenheft.md`), `Priorität: Later`.
- ADR: [ADR-0009 §Entscheidung „Lokale User-Templates" + §Folgepunkte](../../adr/0009-template-format-yaml-files.md)
  — verbindlicher Architektur-Anker.
- Carveout: [`carveouts.md`](../in-progress/carveouts.md) §Carveout-
  Auflösungs-Slices (historisch), [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates)-Zeile — der
  [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)-Carveout wurde mit dieser Closure aufgelöst
  und aus §Temporäre Carveouts entfernt.
- Voraussetzungs-Slices:
  [`slice-v1-template-list`](../done/slice-v1-template-list.md),
  [`slice-v1-template-init`](../done/slice-v1-template-init.md).
- Security: [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte) (keine verdeckte Fremd-Code-Ausführung)
  — durch `text/template` ohne Pre-/Post-Hooks trivial erfüllt,
  bleibt auch für lokale Templates gewahrt.
- Mutex-Regeln: [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
  T0-(i) (`--template` + `--dry-run`/`--diff`).
- Roadmap: [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0+ Backlog,
  Zeile [`slice-later-local-templates`](slice-later-local-templates.md).
- Phase: Later (nach v0.4.0, Trigger-getrieben).
