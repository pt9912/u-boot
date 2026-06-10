# ADR 0008: Add-on-System bleibt statisch, kein Plugin-Loader

## Status

Accepted

## Datum

2026-05-31

## Kontext

[`LH-OPEN-003`](../../../spec/lastenheft.md#lh-open-003-plugin-system-entschieden) (Plugin-System) ist in `spec/lastenheft.md` §14 offen:

> *„Es ist zu klären, ob Add-ons langfristig fest eingebaut oder als
> Plugins nachladbar sein sollen."*

`spec/architecture.md` §7 nennt das Plugin-System prospektiv als
„geplante Erweiterung" (Driven-Port `PluginRegistry`). MVP-Stand
(`e0d6c87`): `postgres` ist das einzige ausgelieferte Add-on
([`LH-FA-ADD-001`](../../../spec/lastenheft.md#lh-fa-add-001-add-on-befehl)..[`LH-FA-ADD-002`](../../../spec/lastenheft.md#lh-fa-add-002-postgresql-hinzufügen), [`LH-FA-ADD-005`](../../../spec/lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern)); Keycloak ([`LH-FA-ADD-003`](../../../spec/lastenheft.md#lh-fa-add-003-keycloak-hinzufügen),
[`LH-AK-003`](../../../spec/lastenheft.md#lh-ak-003-keycloak-flow)) und OpenTelemetry ([`LH-FA-ADD-004`](../../../spec/lastenheft.md#lh-fa-add-004-opentelemetry-hinzufügen), [`LH-AK-004`](../../../spec/lastenheft.md#lh-ak-004-opentelemetry-flow)) sind
für V1 statisch geplant. Es gibt heute keinen externen Wunsch nach
einem vierten, nicht-Kern-Add-on.

Die Entscheidungsfindung betrachtet drei Optionen:

1. **Statisch eingebaute Add-ons** — neue Services im u-boot-Binary;
   neue Add-ons brauchen u-boot-Release.
2. **Plugin-System über Driven-Port** — `PluginRegistry`-Port lädt
   externe Plugin-Binaries oder OCI-Bundles zur Laufzeit.
3. **Hybrid** — Kern-Add-ons statisch, exotische via Plugin.

Sicherheits-Rahmen: [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte) (MVP) verbietet die verdeckte
Ausführung externer Skripte ohne ausdrückliche Nutzer-Zustimmung.
Jedes Plugin-Loader-Modell muss dieses Pflicht-Setting respektieren —
mit explizitem Allowlist-Mechanismus, Signature-Verifikation und
ggf. Sandboxing.

Vergleichbare Tools:

- `kubectl`: statische Kommandos + `krew`-Plugin-Manager (Go-Binaries
  mit Signature-Verifikation, eigener Index). Reife Lösung, aber das
  Plugin-System wurde erst lange nach v1 ergänzt.
- `helm`: ähnliches Modell; Plugins sind separate Repositories, der
  Loader prüft `plugin.yaml`-Manifeste.
- `gh`: GitHub-CLI mit Extensions (`gh extension install`). Auch
  hier Plugin-System nach mehreren Jahren stabiler Kern-CLI ergänzt.

In allen drei Beispielen folgte die Plugin-Architektur nach
einer stabilen Kern-Phase mit konkretem externem Bedarf, nicht
prospektiv.

## Entscheidung

**Statisch** (Option 1). Add-ons bleiben im u-boot-Binary fest
eingebaut. Neue Services werden im u-boot-Repository als
Add-on-Implementierung gegen [`LH-FA-ADD-001`](../../../spec/lastenheft.md#lh-fa-add-001-add-on-befehl)..[`LH-FA-ADD-005`](../../../spec/lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern) ergänzt und mit
einem regulären u-boot-Release distribuiert.

Konkrete Setzungen:

- **Kein `PluginRegistry`-Driven-Port.** `spec/architecture.md` §7
  wird auf den Stand „kein Plugin-Loader, statisches Add-on-Modell"
  umgeschrieben; der prospektive Hinweis verschwindet, der HTTP-
  Driving-Adapter bleibt als geplante Erweiterung erhalten.
- **Add-on-Erweiterung über den Planning-Lifecycle.** Jedes neue Add-on (`add
  keycloak`, `add otel`, später ggf. `add redis`, `add minio`) wird
  als eigenes Planungsartefakt angelegt. Damit bleibt der Add-on-Pfad
  reviewbar und folgt der Planning-Disziplin ([`LH-FA-PROJDOCS-005`](../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)).
- **[`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte) automatisch erfüllt.** Ohne Plugin-Loader gibt
  es keinen Pfad, über den u-boot fremden Code aus nicht-
  freigegebenen Quellen lädt. Die einzigen externen Quellen sind
  Docker-Images (`compose.yaml`/Dockerfile, ohnehin von
  [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte) ausgenommen) und devcontainer-Features (über
  `devcontainer.featureSources.allow`, [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)).
- **Re-Evaluation-Trigger explizit dokumentiert** (siehe
  §Folgepunkte). Sobald einer der genannten Trigger eintritt, wird eine
  neue Entscheidung vorbereitet, die dieses ADR superseded.

Hybrid (Option 3) und volles Plugin-System (Option 2) werden
verworfen:

- **Hybrid:** doppelter Mechanismus (Kern + Plugin) ist Wartungs-
  Overhead ohne erkennbaren Nutzen, solange kein Plugin-Use-Case
  existiert.
- **Plugin-System:** komplexes Sicherheits-Modell (Signing,
  Sandboxing, ABI-Versionierung) für einen heute nicht vorhandenen
  Anwendungsfall. Prospektive Architektur widerspricht der
  Planning-Disziplin („Plan-Loch ist nicht erlaubt"; ohne Trigger ist
  jeder temporäre Carveout nur Vertagung mit Trigger).

## Konsequenzen

Positiv:

- **Minimale Surface.** Kein neuer Driven-Port, keine Plugin-Loader-
  Logik, keine zusätzliche Test-Strecke für Plugin-Lifecycle.
- **[`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004-keine-verdeckte-ausführung-fremder-skripte) trivial erfüllt.** Kein Diskussionsbedarf zu
  Signature-Verifikation, Sandbox-Boundaries, Plugin-Allowlist-
  Format.
- **Add-on-Distribution = u-boot-Distribution.** Genau eine GHCR-
  Pipeline ([ADR-0007](0007-distributionswege-ghcr.md)), genau ein Release-Schnitt; keine separate
  Plugin-Registry / kein zweiter Distributions-Weg.
- **Konsistenz mit der Distributionsentscheidung**:
  GHCR als Distributionsweg für `v0.1.0` ([ADR-0007](0007-distributionswege-ghcr.md)) +
  statische Add-ons heißt: ein Pull = vollständige u-boot-
  Funktionalität.

Negativ / Trade-offs:

- **Add-on-Releases sind an u-boot-Release-Zyklus gekoppelt.** Ein
  Bugfix in einem Add-on (z. B. `postgres`-Compose-Template) zwingt
  einen u-boot-Patch-Release. Bei niedriger Add-on-Frequenz
  (heute drei MVP/V1-Add-ons) akzeptabel.
- **Drittanbieter können keine Add-ons hinzufügen, ohne den
  u-boot-Code zu forken oder einen PR einzureichen.** Solange das
  Projekt klein bleibt, ist PR-getrieben der gewünschte Modus
  (Review-Punkt, Planning-Disziplin); wenn das später eng wird,
  greift einer der Re-Eval-Trigger.
- **`PluginRegistry`-Port-Skizze in `architecture.md` §7 fällt
  weg.** Wer den prospektiven Hinweis in der Vergangenheit gelesen
  hat, muss dieses ADR konsultieren; deshalb wird §7 auf einen
  expliziten ADR-Verweis umgeschrieben (nicht nur stumm gelöscht).

Alternativen (verworfen):

- **Plugin-System ab `v0.1.0` (Option 2):** prospektive
  Architektur für einen nicht vorhandenen Use-Case. Plugin-Loader
  bedeutet ein vollständiges Security-Modell (Plugin-Manifest-
  Schema, Signature-Verifikation, Allowlist in `u-boot.yaml`,
  optional Sandbox). Aktuell ohne erkennbaren Nutzen,
  Wartungslast permanent.
- **Hybrid (Option 3):** Kern statisch, „exotisch" via Plugin —
  zwei parallele Add-on-Pfade ohne klare Trennlinie. Erschwert
  Reviews und führt mit hoher Wahrscheinlichkeit dazu, dass jeder
  neue Service nach „Kern oder Plugin?" gefragt wird; der
  Planning-Pfad wäre derselbe Diskussionsoverhead bei einem
  rein statischen Modell, ohne Plugin-Komplexität.

## Folgepunkte (Re-Evaluation-Trigger)

Dieses ADR ist **revertierbar**: sobald einer der folgenden
Trigger eintritt, wird eine neue ADR vorbereitet, die dieses ADR
superseded.

1. **Drittes externes Add-on-Anfrage.** Konkret: ein Issue oder
   PR fordert einen Service, der nicht zum Kern-Katalog gehört
   und auch nicht plausibel in den Kern aufgenommen werden soll.
2. **Add-on-Release-Frequenz übersteigt u-boot-Release-Frequenz.**
   Konkret: ein Add-on braucht zwischen u-boot-Releases einen
   Hotfix, und der Patch-Release-Aufwand ist deutlich höher als
   ein Plugin-Update wäre.
3. **Externer Contributor möchte ein Add-on pflegen, das nicht in
   den u-boot-Mainline soll.** Konkret: jemand betreibt einen
   produktiven `add minio`-Fork und bittet um einen Pfad, der
   ohne Mainline-Commit auskommt.
4. **Sicherheits-Anforderung („Sandboxed Sub-Tooling") aus dem
   produktiven Einsatz.** Wenn u-boot in einer Umgebung läuft, die
   pro Service eine Isolation verlangt (z. B. kein Filesystem-
   Zugriff für Add-on-Code), wäre ein Plugin-Boundary natürlicher
   als reines Compose-File-Output.

Solange keiner dieser Trigger eintritt, bleibt das Add-on-System
statisch.
