# Devcontainer Features — u-boot

| Dokument    | Devcontainer-Features-User-Guide |
| ----------- | -------------------------------- |
| Projektname | `u-boot` |
| Bezug       | `LH-FA-DEV-003` in [`spec/lastenheft.md`](../../spec/lastenheft.md) §692-721 + §1340-1353 + §2394 |
| Slice       | [`slice-v1-devcontainer-features`](../plan/planning/done/slice-v1-devcontainer-features.md) |
| ADR         | [`docs/plan/adr/0008-plugin-system-statisch.md`](../plan/adr/0008-plugin-system-statisch.md) §78 |
| Status      | v0.4.0 |

## Zweck

Erklärt, wie u-boot die [Dev Containers
Features](https://containers.dev/implementors/features/) im
`.devcontainer/devcontainer.json` verwaltet: welche Features
out-of-the-box im Built-in-Katalog stehen, wie der User externe
Quellen freigeben muss, und welche Doctor- bzw.
Generator-Verhaltensregeln daraus folgen.

Die Pflichtaussagen leben im Lastenheft (`LH-FA-DEV-003`); dieses
Dokument ist die Bedien-Doku.

---

## 1. Built-in-Katalog (8 Features)

u-boot bringt einen kuratierten Default-Satz devcontainer-Features
mit. Sie sind ohne Allowlist-Eintrag aktivierbar (Spec §711). Alle
zeigen aktuell auf `ghcr.io/devcontainers/features/<name>:1` —
Versions-Pin per Default `1`, override pro Feature siehe §3.

| Catalogue-Key       | Source                                                     | Inhalt                          |
| ------------------- | ---------------------------------------------------------- | ------------------------------- |
| `git`               | `ghcr.io/devcontainers/features/git`                       | Git CLI                         |
| `docker-cli`        | `ghcr.io/devcontainers/features/docker-outside-of-docker`  | Docker CLI (outside-of-docker)  |
| `node`              | `ghcr.io/devcontainers/features/node`                      | Node.js                         |
| `java`              | `ghcr.io/devcontainers/features/java`                      | Java + SDKMAN                   |
| `go`                | `ghcr.io/devcontainers/features/go`                        | Go toolchain                    |
| `cpp`               | `ghcr.io/devcontainers/features/cpp`                       | C++ toolchain                   |
| `kubectl-helm`      | `ghcr.io/devcontainers/features/kubectl-helm-minikube`     | kubectl + helm + minikube       |
| `postgres-client`   | `ghcr.io/devcontainers/features/postgresql-client`         | PostgreSQL client               |

Externe Features (jede andere Source) brauchen den Allowlist-Weg
aus §4.

---

## 2. Catalogued Feature aktivieren

Drei-Schritt-Workflow für built-in Features:

```bash
# 1. Projekt mit Devcontainer initialisieren
u-boot init --devcontainer

# 2. Feature aktivieren (config set akzeptiert true/false/1/0)
u-boot config set devcontainer.features.node.enabled true

# 3. devcontainer.json regenerieren
u-boot generate devcontainer
```

Resultat in `.devcontainer/devcontainer.json` innerhalb des
managed-Blocks:

```jsonc
// BEGIN U-BOOT MANAGED BLOCK: init
{
  "name": "demo",
  "build": { … },
  "features": {
    "ghcr.io/devcontainers/features/node:1": {}
  },
  "remoteUser": "vscode"
}
// END U-BOOT MANAGED BLOCK: init
```

`u-boot generate devcontainer` ist idempotent (zweiter Aufruf
= No-Op). Mehrere Features werden alphabetisch nach Source
sortiert, damit das JSON byte-deterministisch bleibt.

---

## 3. Version-Override pro Feature

Standardmäßig zieht jedes Catalogued Feature die Version `1`. Für
einen abweichenden Pin (z. B. Java 21 statt der default-major):

```bash
u-boot config set devcontainer.features.java.enabled true
u-boot config set devcontainer.features.java.version 21
u-boot generate devcontainer
```

`u-boot.yaml` zeigt dann:

```yaml
devcontainer:
  enabled: true
  features:
    java:
      enabled: true
      version: "21"
```

Render-Output: `"ghcr.io/devcontainers/features/java:21": {}`.

---

## 4. Externe Feature-Quellen (Allowlist)

Spec §710-721 verlangt, dass externe Features nur mit expliziter
Freigabe aktiviert werden können. `--yes` reicht **nicht** als
Substitut (LH-NFA-SEC-004).

### Allowlist befüllen

Drei gleichwertige Pfade — alle persistieren die URL in
`devcontainer.featureSources.allow`:

```bash
# Pfad A: bei init
u-boot init --devcontainer \
  --allow-external-feature-sources https://ghcr.io/orgX/features/custom-rust

# Pfad B: vor jeder generate (kann mehrfach kumulieren)
u-boot generate devcontainer \
  --allow-external-feature-sources https://ghcr.io/orgX/features/custom-rust

# Pfad C: direkt via config set
u-boot config set devcontainer.featureSources.allow \
  https://ghcr.io/orgX/features/custom-rust
```

Die drei Befehle sind die einzigen, auf denen
`--allow-external-feature-sources` gültig ist (Spec §714-717);
andernorts (z. B. `generate readme`) wird das Flag abgelehnt.

Mehrere URLs gehen via Komma:

```bash
u-boot config set devcontainer.featureSources.allow \
  https://a.test/x,https://b.test/y
```

Der Flag-Vorkommen ist additiv (kein last-wins) und das Schreiben
ist silent-dedupe (Spec §1352): `https://a.test/x` taucht in der
Allowlist nur einmal auf, auch wenn du es mehrfach passierst.

### Externes Feature aktivieren

Nachdem die URL in der Allowlist steht:

```bash
u-boot config set devcontainer.features.custom-rust.source \
  https://ghcr.io/orgX/features/custom-rust
u-boot config set devcontainer.features.custom-rust.enabled true
u-boot generate devcontainer
```

Resultat:

```yaml
devcontainer:
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
```

### Allowlist-Vergleich — was muss matchen

Die Source-Override wird **byte-equal** gegen die Allowlist-
Einträge geprüft. Das heißt:

- Trailing-Slashes sind signifikant: `https://x/y` und
  `https://x/y/` sind verschiedene Einträge.
- Host-Case ist signifikant: `https://X.io/y` und `https://x.io/y`
  matchen nicht.

Der Hint in der Doctor- / Config-Fehlermeldung weist auf beide
Fallen hin.

### Was passiert ohne Allowlist-Eintrag

```bash
u-boot config set devcontainer.features.custom-rust.source \
  https://ghcr.io/orgX/features/custom-rust
# → Exit-Code 10 (LH-FA-DEV-003 / LH-NFA-SEC-004):
#   external source "..." is not in devcontainer.featureSources.allow
```

Repair-Hint im Fehler nennt den genauen `config set ...allow`-
Aufruf.

---

## 5. Doctor-Verhalten

`u-boot doctor` enthält den Check
`devcontainer.features.allowlist` mit drei Klassifikationen:

| Schweregrad | Trigger                                                                                       |
| ----------- | --------------------------------------------------------------------------------------------- |
| **Error**   | `source:`-Override gesetzt, aber URL nicht in `featureSources.allow` (LH-FA-DEV-003 §720).    |
| **Warn**    | Orphan-Activation: `source:` leer und Name nicht im Built-in-Katalog (Renderer skippt den Eintrag still). |
| **Warn**    | `enabled:` fehlt für einen Feature-Eintrag (LH-FA-ADD-005 §893-Analog).                       |

Die Worst-Severity gewinnt: wenn Allowlist-Violation UND
Orphan-Activation gleichzeitig vorliegen, surface Doctor die
Error-Meldung.

Doctor schweigt (OK) wenn:

- `u-boot.yaml` fehlt / nicht parsbar (gehört zum primären
  `uboot.yaml.valid`-Check).
- `cfg.Devcontainer` fehlt oder `features:`-Map leer (Spec §2394
  negative pin: kein Error ohne legitimen Anlass).

### Drift-Check `devcontainer.features.drift`

Der zweite, ergänzende Check vergleicht `u-boot.yaml` gegen die
Keys in `.devcontainer/devcontainer.json`'s `features:`-Map und
erkennt drei Drift-Situationen (jeweils Severity **Warn** mit
Repair-Hint „`u-boot generate devcontainer`"):

| Case | Trigger | Repair |
| ---- | ------- | ------ |
| **1** | Feature ist `enabled: true` in u-boot.yaml, fehlt aber im JSON (oder JSON-Datei fehlt ganz). | `generate devcontainer` ausführen. |
| **2a** | User hat `enabled: false` (oder unset) gesetzt, der JSON-Key steht aber noch drin. | `generate devcontainer` ausführen. |
| **2b** | JSON enthält einen Feature-Key, für den u-boot.yaml *keinen* Eintrag hat (Hand-Edit oder Drift aus früherem u-boot-Stand). | Eintrag in u-boot.yaml ergänzen oder Key aus JSON entfernen. |

Der Drift-Check skippt (OK), wenn weder u-boot.yaml-Features noch
JSON-Features konfiguriert sind, oder wenn das JSON nicht parsbar
ist (dafür ist `devcontainer.json.valid` zuständig).

`nil` (kein `devcontainer.features:`-Block) und explizit leere
Map (`features: {}`) werden unterschieden: bei expliziter leerer
Map feuert Case 2b weiterhin, wenn das JSON Keys enthält.

---

## 6. Renderer-Verhalten

`u-boot generate devcontainer` projiziert
`cfg.Devcontainer.Features` per Catalogue-Lookup oder
Source-Override und schreibt den `"features": {…}`-Block innerhalb
des managed-Blocks. Render-Disziplin:

- Nur Einträge mit `enabled: true` landen im JSON.
- Reihenfolge: alphabetisch nach Source (deterministisches Output-
  Byte-Layout).
- Unbekannte Feature-Namen ohne `source:`-Override werden still
  übergangen — Doctor surface das (siehe §5).
- Idempotenz: zwei aufeinanderfolgende `generate devcontainer`-
  Aufrufe schreiben kein zweites Mal (Byte-Equal-NoOp).

---

## 7. Beispiel: vollständiger Workflow

```bash
# Setup
mkdir myproj && cd myproj
u-boot init --devcontainer

# Catalogued Features
u-boot config set devcontainer.features.git.enabled true
u-boot config set devcontainer.features.node.enabled true
u-boot config set devcontainer.features.java.enabled true
u-boot config set devcontainer.features.java.version 21

# Externes Feature
u-boot config set devcontainer.featureSources.allow \
  https://ghcr.io/orgX/features/custom-rust
u-boot config set devcontainer.features.custom-rust.source \
  https://ghcr.io/orgX/features/custom-rust
u-boot config set devcontainer.features.custom-rust.enabled true

# Render + Sanity-Check
u-boot generate devcontainer
u-boot doctor

# Ergebnis: devcontainer.json mit 4 features, doctor OK.
```

---

## 8. Out of Scope

- **Eigene/lokale Features** (Custom-Feature im Repo-Pfad statt
  externer Quelle): geht über `LH-FA-DEV-003` hinaus; eigener
  Folge-Slice mit eigenem Trigger.
- **Sprach-spezifische Build-Aktionen** (Gradle-Wrapper anlegen,
  Go-Module-Init, …): das ist Template-Job (`LH-FA-TPL-*`), nicht
  Devcontainer-Feature-Job.
- **Feature-Version-Updates** (Renovate/Dependabot-Hook): nice-to-
  have; eigener Folge-Slice nach Erstauslieferung.
