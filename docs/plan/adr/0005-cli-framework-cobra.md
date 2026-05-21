# ADR 0005: CLI-Framework Cobra

## Status

Accepted

## Datum

2026-05-21

## Kontext

ADR-0001 (Implementierungssprache Go) hat als offenen Folgepunkt:

> *„Wahl des CLI-Frameworks (`flag` aus stdlib reicht für MVP-Stub;
> Cobra wird mit `add`/`generate`/`config`-Subkommandos
> wahrscheinlich nötig)."*

M1 hat `flag` aus der stdlib für den `--help`/`--version`-Stub
benutzt; mit M3-T3 (dem ersten Subkommando `init`) wird das zu eng.
Die offene Entscheidung muss formal getroffen und dokumentiert werden
(`LH-FA-PROJDOCS-005`, slice
[`slice-m3-cli-framework-adr.md`](../planning/done/slice-m3-cli-framework-adr.md)).

Vorlagen / Markt:

- `spf13/cobra` — De-facto-Standard im Go-Ökosystem (kubectl, hugo,
  helm, gh, docker-cli). Subkommandos, automatisches Help-Layout,
  Bash-/Zsh-/Fish-Completion, durchdachter Test-Pfad
  (`SetArgs`/`SetOut`).
- `urfave/cli` v3 — Alternative mit ähnlichem Funktionsumfang,
  etwas schlanker, kleinere Verbreitung.
- `stdlib flag` — keine Subkommandos, kein Help-Layout für komplexe
  CLIs.

## Entscheidung

**Cobra** als CLI-Framework, mit `github.com/spf13/cobra` als
Modul-Dependency. Pin im aktuellen Stand: `v1.10.2` (transitive
Dep `github.com/spf13/pflag v1.0.9`, `github.com/inconshreveable/mousetrap v1.1.0`).

Konkrete Setzungen:

- `internal/adapter/driving/cli/` als Cobra-Adapter (LH-FA-ARCH-002).
- App-Konstruktor `cli.New(version, useCase, opts...)` nimmt die
  Driving-Ports als Konstruktor-Args; Wiring erfolgt in `cmd/uboot/`.
- Funktionale Options (`cli.WithGetwd(...)`) für Test-Seams.
- Exit-Code-Mapping (`cli.ExitCode(err)`) bündelt die
  `LH-FA-CLI-006`-Logik im Adapter:
  - `nil` → 0
  - `driving.ErrProjectExists` → 10
  - Cobra-Usage-Errors (unbekannte Subkommandos/Flags, falsche
    Argumentzahl) → 2
  - sonst → 1
- `RunE`-Closure ruft `cmd.Context()` und reicht ihn an die
  `runInit`-Funktion durch, die Context als ersten Parameter
  explizit annimmt — so passt es zu `contextcheck` auf der Tiefe
  unter dem Cobra-Closure. Cobras `RunE`-Signatur selbst kennt
  keinen Context-Parameter; das ist ein permanenter `contextcheck`-
  Carveout für `internal/adapter/driving/cli/`
  (`carveouts.md`-Permanent-Tabelle).

## Konsequenzen

Positiv:

- **Subkommando-Routing** automatisch (`init` heute, `add`/`up`/
  `down`/`doctor`/`generate`/`config`/`template` folgen).
- **Help-Layout** ohne Boilerplate; `--help` und `<command> --help`
  identisch in jeder Tiefe.
- **Shell-Completion** out-of-the-box; in einem späteren Slice für
  alle Subkommandos aktiviert.
- **Idiomatisches Test-Pattern** (`SetArgs`/`SetOut`/`SetErr` für
  In-Memory-Buffers).
- **Etabliert in der Domäne** — neue Beitragende kennen Cobra aus
  kubectl/gh/docker-cli; keine Lernkurve.

Negativ / Trade-offs:

- **Modul-Dependencies wachsen** um Cobra + pflag + mousetrap
  (~3 zusätzliche transitive Module, alle stabil und gepflegt).
- **`contextcheck`-Carveout** für den CLI-Adapter — Cobras
  `RunE`-Signatur hat keinen Context-Parameter; das Pattern
  „`cmd.Context()` im Closure extrahieren" ist unvermeidbar
  (siehe Permanent-Carveout in carveouts.md).
- **Cobra-Versions-Major-Bumps** können Help-/Error-Messaging
  ändern; der `cli.ExitCode`-Test pinnt die heutigen Cobra-Usage-
  Error-Präfixe, ein Bump erfordert ggf. eine Aktualisierung.

Alternativen (verworfen):

- **stdlib `flag` weiterbehalten:** kein Subkommando-Routing,
  jedes Subkommando bräuchte einen eigenen Sub-`flag.FlagSet`
  plus manuelles Dispatchen. Skaliert nicht über zwei, drei
  Subkommandos hinaus.
- **`urfave/cli` v3:** vergleichbarer Funktionsumfang, etwas
  geringerer Verbreitungsgrad. Bei vergleichbarer Funktionalität
  gewinnt Verbreitung.

## Folgepunkte

- Shell-Completion-Generation (`cobra completion bash|zsh|fish`)
  in einem späteren Slice exposed.
- Bei Bedarf: `cobra/doc` für automatische Manpage-/Markdown-
  Generation der Subkommandos.
- ADR-0001 „Offene Folgepunkte" wird mit M3-Closure auf
  „CLI-Framework geschlossen via ADR-0005" verkürzt.
