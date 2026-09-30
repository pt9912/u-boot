# ADR 0015: Sandbox-Volumes pro Instanz über `${devcontainerId}` in `mounts`

**Status:** Accepted

**Datum:** 2026-09-30

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0014](0014-nested-podman-sandbox-devcontainer.md)

**Schärft:** [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) — Benennung des Workspace-Volumes im Sandbox-Profil (der Spec verlangt „benanntes Volume“, nicht den Namen).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Bis zur Umsetzung von ADR-0014 hießen die Sandbox-Volumes `<projekt>-workspace` und `<projekt>-containers`. Zwei Container desselben Projekts (etwa mehrere autonome Agenten parallel, jeweils aus einem eigenen Ordner oder Git-Worktree) teilten damit Workspace und Podman-Storage; der Podman-Storage ist nicht für gleichzeitigen Zugriff mehrerer Container ausgelegt.

**Messung (2026-09-30, `devcontainer` CLI 0.80.2, Docker 29.8.1).** Zwei Ordner `a` und `b` mit identischer `devcontainer.json`:

| Variante | Ergebnis |
|---|---|
| `${devcontainerId}` in `workspaceMount` | **nicht aufgelöst** (`source=demo-ws-${devcontainerId}` bleibt wörtlich, `docker run` scheitert) |
| `${devcontainerId}` in `mounts` | aufgelöst, **pro Ordner verschieden** (52 Zeichen, stabil pro Ordner) |
| `"workspaceMount": ""` + Workspace-Volume in `mounts` mit `${devcontainerId}` | Container startet; die Mount-Liste enthält **nur die beiden Volumes**, kein Bind-Mount des Host-Ordners; `a` und `b` bekommen verschiedene Volumes |

Nicht gemessen: VS Code selbst (nutzt dieselbe Spezifikation, aber eine eigene Implementierung) und GitHub Codespaces.

## Entscheidung

Im Sandbox-Profil setzt u-boot `"workspaceMount": ""` (schaltet den Standard-Bind-Mount des Projektordners ab) und bindet **beide** Volumes über `mounts` ein, mit dem Instanz-Anteil `-${devcontainerId}` im Namen:

```text
source=<projekt>-workspace-${devcontainerId},target=/workspaces/<projekt>,type=volume
source=<projekt>-containers-${devcontainerId},target=/home/vscode/.local/share/containers,type=volume   # nur mit nestedRuntime: podman
```

Eine Instanz ist ein Ordner (`devcontainerId` ist pro Ordner stabil). Parallele Instanzen desselben Projekts laufen aus getrennten Ordnern (Klone, Git-Worktrees). Es gibt keinen neuen Konfigurationsschlüssel und keine Spec-Änderung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Status quo (`<projekt>-workspace`) | keine Änderung | parallele Instanzen teilen Workspace und Storage (Kern des Problems) |
| B — `${localWorkspaceFolderBasename}` im Namen | Variable ist in `workspaceMount` verfügbar | Ordner gleichen Namens (`a/demo`, `b/demo`) kollidieren weiter |
| C — Konfigurationsschlüssel `devcontainer.sandbox.instance` | explizit, unabhängig von Werkzeug-Variablen | zusätzliche Vertragsfläche (Spec-Anhebung), Nutzer muss pro Ordner einen Wert pflegen |
| **D — `workspaceMount: ""` + `mounts` mit `${devcontainerId}`** | gemessen lauffähig, automatisch eindeutig pro Ordner, keine Spec-Änderung, kein Host-Bind-Mount | zwei Instanzen aus **demselben** Ordner bleiben eine Instanz; Verhalten in VS Code/Codespaces ungeprüft; lange Volume-Namen |

## Konsequenzen

- Positiv: Parallele Sandbox-Instanzen aus getrennten Ordnern sind isoliert; der Bind-Mount-Ausschluss aus [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) bleibt erhalten.
- Negativ: Die alte Benennung entfällt. Das Sandbox-Profil ist noch nicht released; lokal angelegte Volumes `<projekt>-workspace` / `<projekt>-containers` bleiben verwaist und können mit `docker volume rm` entfernt werden. Der Workspace-Inhalt eines neuen Volumes entsteht durch den Clone-Schritt neu (nicht gepushte Arbeit im alten Volume geht nicht automatisch über).
- Folgepflicht: Nachholmessung in VS Code (Dev Containers) und Codespaces; Ergebnis als Zeile in der Geschichte. Compose-Service-Ports mehrerer Instanzen bleiben ein getrennter Vorgang.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Golden Case | Sandbox-`devcontainer.json` hat `"workspaceMount": ""` und Volumes mit `${devcontainerId}` in `mounts`, nie `type=bind` | `make test` |
| Integrationstest (`//go:build docker`) | zwei Container mit verschiedener ID teilen weder Workspace noch Storage | `make test-docker` |

## Re-Evaluierungs-Trigger

Ein Befund, dass VS Code oder Codespaces `workspaceMount: ""` anders behandeln (dann Variante C oder eine werkzeugspezifische Ausgabe), oder eine Spezifikationsänderung, die `${devcontainerId}` in `workspaceMount` erlaubt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-30 | Entschieden und `Accepted` (Messung mit `devcontainer` CLI 0.80.2) | Vereinbarung mit dem Projektinhaber („alles fertig machen“) |
