# Slice V1: Mehrere Sandbox-Devcontainer-Instanzen desselben Projekts

**Lifecycle:** Zustand = Verzeichnis (`open/` → `next/` → `in-progress/`
→ `done/`), Wechsel nur per `git mv`.

**Welle:** ohne Welle (kein Release-Ziel, Folge der Sandbox-Umsetzung).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)
(benanntes Volume als Workspace),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
(Podman-Storage-Volume). Folge-Slice von
[`slice-v1-sandbox-devcontainer-umsetzung`](../done/slice-v1-sandbox-devcontainer-umsetzung.md)
(Abschnitt „Mehrere Instanzen“ in
[`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](../../../user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten)).

**Autor:** pt9912. **Datum:** 2026-09-30.

---

## 1. Ziel und Abgrenzung

**Ziel:** Mehrere Sandbox-Devcontainer **desselben Projekts** (z. B.
mehrere autonome Agenten parallel) laufen gleichzeitig, ohne Workspace-
und Podman-Storage-Volumes zu teilen.

**Auslöser:** Die erzeugten Volume-Namen `<projekt>-workspace` und
`<projekt>-containers` hängen nur am Projektnamen. Zwei Container
desselben Projekts — auch aus verschiedenen Ordnern mit gleichem
`project.name` — teilen Workspace und Podman-Storage. Der Podman-Storage
ist nicht für gleichzeitigen Zugriff mehrerer Container ausgelegt.

**Ausdrücklich NICHT in diesem Slice:**

- Compose-Service-Ports (`u-boot up`) mehrerer Instanzen — die Kollision
  der veröffentlichten Host-Ports liegt im Compose-Stack, nicht im
  Devcontainer; eigener Vorgang mit eigener Entscheidung (Port-Offset,
  Projektname pro Stack).
- Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
  V2) — eigener Slice nach Ratifizierung von
  [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md).
- Mehrere Instanzen im **Default-Profil** — dort ist der Projektordner
  per Bind-Mount eingebunden, die Dev-Containers-Werkzeuge trennen
  Instanzen bereits pro Ordner.

## 2. Definition of Done

- [x] **Verifikation vorab (Risiko 1):** Ob Dev Containers
  `${devcontainerId}` bzw. `${localWorkspaceFolderBasename}` in
  `workspaceMount` und `mounts` auflöst, ist mit VS Code **und** der
  `devcontainer`-CLI geprüft; Befund (auch Negativbefund) steht im
  Umsetzungs-ADR.
- [x] **Entscheidung** (Umsetzungs-ADR, ≥ 3 Alternativen): eindeutiger
  Instanz-Anteil im Volume-Namen über (a) `${devcontainerId}`, (b)
  `${localWorkspaceFolderBasename}` oder (c) einen Konfigurationsschlüssel
  (z. B. `devcontainer.sandbox.instance`) — inklusive Verhalten bei
  Wechsel der Namensregel für bestehende Projekte (Volumes der alten
  Benennung, Migration oder bewusster Bruch).
- [x] **Spec-Prüfung:** Falls die Namensregel vertragsrelevant ist
  (Konfigurationsschlüssel, Verhalten bestehender Volumes), Spec-Anhebung
  nach §16 mit Historie-Zeile; sonst ausdrücklich „keine Spec-Änderung“
  begründet.
- [x] **Generator:** Workspace- und Containers-Volume tragen den
  Instanz-Anteil; Default-Profil-Ausgabe bleibt byte-identisch; Golden
  Cases (Fresh, NoOp, Wechsel alte → neue Benennung) nach
  `harness/replay.md`.
- [x] **Integrationstest** (`//go:build docker`): zwei Instanzen mit
  unterschiedlichem Instanz-Anteil laufen gleichzeitig, schreiben in ihren
  Workspace und führen `podman run` aus, ohne dass sich Volumes oder
  Storage berühren.
- [x] **Doku:** `docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten` §„Mehrere Instanzen“
  auf den neuen Stand, CHANGELOG-Eintrag, `make gates` grün.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Umsetzungs-ADR (neu) | neu | Namensregel, Verifikationsbefund, Migrationsverhalten |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`, `Dockerfile.tmpl` | update | Volume-Quellen mit Instanz-Anteil (Mount-Ziele bleiben) |
| `internal/hexagon/application/devcontainer_sandbox.go`, `templates.go` | update | Volume-Namen / ggf. Config-Key |
| `internal/e2e/` | neu | Zwei-Instanzen-Test |
| `spec/lastenheft.md` | ggf. update | nur falls vertragsrelevant (siehe DoD) |
| `docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`, `CHANGELOG.md` | update | öffentliche Verträge |

## 4. Trigger

- **`open` → `next`:** Projektinhaber priorisiert (konkreter Bedarf an
  parallelen Sandbox-Instanzen).
- **`next` → `in-progress`:** Beginn mit der Vorab-Verifikation; fällt sie
  negativ aus (Variablen nicht auflösbar), bleibt nur die
  Konfigurationsschlüssel-Variante (c).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Zwei-Instanzen-Integrationstest
ausgeführt, Delivery-Hash in der Closure-Notiz.

## 6. Risiken und offene Punkte

- **Variablen-Auflösung unverifiziert:** Ob `${devcontainerId}` in
  `workspaceMount` funktioniert, ist nicht geprüft (Werkzeug-abhängig:
  VS Code, `devcontainer`-CLI, GitHub Codespaces).
- **Bestehende Volumes:** Eine geänderte Namensregel verwaist die Volumes
  bereits erzeugter Sandbox-Projekte (Workspace-Inhalt im alten Volume).
- **`devcontainerId` ist pro Ordner stabil**, nicht pro gewünschter
  Instanz: Zwei Instanzen aus **demselben** Ordner bleiben damit eine
  Instanz; parallele Instanzen brauchen getrennte Ordner (z. B. Git-
  Worktrees) oder den Konfigurationsschlüssel.

## 7. Closure-Notiz (nach `done/`)

**Status:** **abgeschlossen** (2026-09-30, **Delivery-Hash: `0723eb1`**).

- **Verifikation:** Mit `devcontainer` CLI 0.80.2 gemessen: `${devcontainerId}`
  wird in `mounts` aufgelöst, in `workspaceMount` **nicht**; `"workspaceMount": ""`
  plus Workspace-Volume in `mounts` startet ohne Host-Bind-Mount
  ([ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md)).
- **Entscheidung:** Variante D (`mounts` mit `${devcontainerId}`); keine
  Spec-Änderung, kein neuer Konfigurationsschlüssel (Spec-Prüfung: das
  „benannte Volume“ aus [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) bleibt erfüllt, der Name ist kein Vertrag).
- **Sensoren:** `make gates` grün (Coverage 91,7 %); `go test -tags docker
  ./internal/e2e -run SandboxDevcontainer` (Zwei-Instanzen-Isolation, nested
  Podman, Degradation) grün; `make ci` und voller `make test-docker` nicht
  ausgeführt.
- **Offen:** VS Code und Codespaces ungeprüft ([ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) Folgepflicht);
  Compose-Service-Ports mehrerer Instanzen; zwei Instanzen aus demselben
  Ordner bleiben eine Instanz.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas GF (Generator, Templates, Doku), siehe Kurs Modul 5
§Worked Mini-Example.
