# ADR 0006: revive Custom-Rules-Profil

## Status

Accepted

## Datum

2026-05-27

## Kontext

ADR-0003 (SOLID-nahes Lint-Profil) hat als offenen Folgepunkt:

> *„Erweiterung um `revive`-Custom-Rules in einem Folge-ADR, falls
> die default-Konfiguration zu schwach wird."*

Bis zu diesem ADR lief `revive` ohne expliziten `rules:`-Block und
nutzte damit die in `golangci-lint` per Default eingestellten 24
revive-Regeln (`blank-imports`, `context-as-argument`, ...,
`var-naming`). Diese Defaults haben sich bewährt: die einzigen
fachlichen revive-Findings im Verlauf der Bootstrap-Phase waren bei
absichtlich provozierten Verstößen (z. B. `verify-depguard.sh`-Stub-
Dateien ohne Package-Doc — siehe M3-T5).

Auslöser für diesen ADR: in der M3-Anker-Triage-Sitzung 2026-05-27
wurde der `slice-v2-revive-custom-rules.md`-Carveout vorgezogen, obwohl
keiner der ursprünglichen Trigger (wiederholte Reviewer-Findings, neuer
Style-Beschluss) eingetreten ist. Begründung des Vorziehens: die
Code-Basis stabilisiert sich gerade auf der MVP-Schwelle, und ein
expliziter Regel-Block macht zukünftige Änderungen (jede Rule-Hebung
wird zur expliziten Policy-Entscheidung) transparenter als das
implizite „läuft halt mit den Defaults".

`golangci-lint`-Mechanik: sobald `linters.settings.revive.rules`
gesetzt ist, **ersetzt** das Slice die Defaults vollständig — es gibt
keinen „defaults + extras"-Modus. Folge: das Profil muss alle
Default-Regeln explizit aufzählen, sonst werden sie deaktiviert.

## Entscheidung

`linters.settings.revive.rules` in `.golangci.yml` enthält:

1. **Alle 24 Default-Regeln** explizit aufgezählt — `blank-imports`,
   `context-as-argument`, `context-keys-type`, `dot-imports`,
   `empty-block`, `error-naming`, `error-return`, `error-strings`,
   `errorf`, `exported`, `if-return`, `increment-decrement`,
   `indent-error-flow`, `package-comments`, `range`,
   `receiver-naming`, `redefines-builtin-id`, `superfluous-else`,
   `time-naming`, `unexported-return`, `unreachable-code`,
   `unused-parameter`, `var-declaration`, `var-naming`.

2. **Eine projekt-spezifische Erweiterung: `unused-receiver`**.

   Begründung: u-boot hat eine hexagonale Architektur mit vielen
   Service- und Adapter-Strukturen. Methoden, die ihren Receiver
   nicht referenzieren, signalisieren entweder Refactoring-Potenzial
   (Free-Function statt Methode) oder schlampige Interface-
   Implementierungen. Der Check hilft, die Struktur intentional zu
   halten.

   Test-Files sind ausgenommen, weil Fakes oft stateless Methoden
   für Interface-Erfüllung implementieren (gleicher Grund wie der
   bestehende `unused-parameter`-Test-Exclude).

3. **Bestehende `unused-parameter`-Test-Exclusion** bleibt; neue
   `unused-receiver`-Test-Exclusion analog hinzugefügt.

Future-Rules-Hinzufügungen: jede neue revive-Regel braucht einen
Eintrag in dieser ADR-Folgesektion plus einen `Why:` im
`.golangci.yml`-Kommentarblock.

## Konsequenzen

- **Refactoring-Beifang in diesem Slice:** `resolveProjectName` in
  `internal/hexagon/application/initproject.go` war eine Methode auf
  `InitProjectService`, ohne den Receiver zu nutzen. Refactored zu
  einer Free-Function (gleicher Code, einfachere Aufruf-Signatur).
- **CI-Stabilität:** die expliziten 24 Default-Regeln machen unseren
  Lint-Stand robust gegen golangci-lint-Default-Drift. Sollte ein
  zukünftiges golangci-lint-Release seine revive-Default-Liste ändern,
  bleibt unser Profil unverändert.
- **Maintenance:** jede revive-Rule-Hebung muss explizit gemacht
  werden; das ist gewollt — kein „silent default-bump".
- **Verifikations-Pflicht beim Pin-Bump:** die Liste der 24 Defaults
  in `.golangci.yml` wurde gegen revive's Default-Set zum Zeitpunkt
  dieses ADR (golangci-lint v2.12.2, revive bundled) festgelegt. Bei
  jedem Hebung des `GOLANGCI_LINT_VERSION`-Pins muss gegen
  `revive --config <empty>` oder die golangci-lint-Release-Notes
  geprüft werden, ob neue Default-Regeln dazugekommen sind, die wir
  haben oder explizit verwerfen wollen. Sonst kann uns ein
  Default-Set-Wachstum entgehen.

## Verworfen

- **`enable-all-rules: true` + Selective Disable:** schaltet auch
  Rules wie `cognitive-complexity` und `cyclomatic` an, die mit den
  bereits aktiven `gocognit`/`cyclop`-Lintern überlappen oder
  konflikten würden.
- **Komplettes Abschalten von `revive`:** würde Lint-Coverage
  ersatzlos verlieren (24 Default-Checks fallen weg).
- **Zusätzliche Custom-Rules wie `early-return`, `var-naming`-
  Whitelisting:** ohne konkreten Trigger schwer zu rechtfertigen; bei
  nächstem Bedarf in einer ADR-Folgesektion ergänzen.

## Bezug

- Auslösende ADR: `0003-solid-nahes-lint-profil.md` Folgepunkte.
- Slice: [`slice-v2-revive-custom-rules.md`](../planning/done/slice-v2-revive-custom-rules.md)
  (vorgezogen 2026-05-27).
- LH-Verweise: `LH-QA-004` (SOLID-Profil), `LH-FA-PROJDOCS-005`
  (Carveout-Disziplin).
