# Slice V2: Egress-Restriktion im Sandbox-Devcontainer

**Lifecycle:** Zustand = Verzeichnis (`open/` → `next/` → `in-progress/`
→ `done/`), Wechsel nur per `git mv`.

**Welle:** ohne Welle (Folge der Sandbox-Umsetzung).

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md).

**Autor:** pt9912. **Datum:** 2026-09-30.

**Status:** **abgeschlossen** (2026-09-30, **Delivery-Hash: `9b1c414`**).

---

## 1. Ziel und Abgrenzung

**Ziel:** Opt-in Egress-Restriktion (Allowlist) für das Sandbox-Profil als
Guardrail, inklusive Ratifizierung von [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md).

**Ausdrücklich NICHT in diesem Slice:**

- Restriktion auf Host-/Netzwerk-Ebene oder per Proxy — Nicht-Ziel des
  Lastenhefts (keine Host-Konfiguration); Alternative D im ADR.
- Erfassung des Verkehrs verschachtelter Podman-Container (`FORWARD`) — im ADR als Grenze benannt.

## 2. Definition of Done

- [x] **[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) ratifiziert** (MADR-Form, Alternativen, Mechanismus gemessen, offene Fragen 1 und 4 entschieden).
- [x] **Spec:** Lastenheft 0.3.3 (Prüfung der Capability beim Containerstart statt im `doctor`).
- [x] **Config:** `devcontainer.sandbox.egress.enabled` und `.allow` (Liste, Hostnamen-Validierung, Exit 10, Append/Dedupe).
- [x] **Generator:** `NET_ADMIN`-Lockerung ausgewiesen, `egress-init.sh`, `postStartCommand`, nftables/dnsmasq im Dockerfile, Default-Allowlist je Stack.
- [x] **Degradation:** ohne `NET_ADMIN` Warnung oder Exit 11 (`onUnavailable`).
- [x] **Doctor:** `devcontainer.sandbox.egress` (statisch; 16 Checks).
- [x] **Tests:** Golden Cases (Egress an/aus/ohne Sandbox/mit Podman), Config-Tests, Doctor-Test, Docker-Integrationstest (erlaubter Host erreichbar, fremder Host nicht auflösbar/erreichbar, Degradation warn/fail).
- [x] **Doku:** `docs/user/devcontainer-sandbox.md` §7, CHANGELOG, README, `harness/replay.md`.

## 3. Plan (vor Code)

Umgesetzt wie im ADR beschrieben; Dateien: `devcontainer_egress.go`,
`devcontainer_sandbox_config.go`, Templates (`egress-init.sh.tmpl`,
`devcontainer.json.tmpl`, `Dockerfile.tmpl`), `doctor_sandbox.go`,
`internal/e2e/sandbox_devcontainer_docker_test.go`.

## 4. Trigger

Erledigt: `open` → `next` → `in-progress` am 2026-09-30 auf Anweisung des
Projektinhabers („alles fertig machen“); [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) war Voraussetzung.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Docker-Integrationstest ausgeführt.

## 6. Risiken und offene Punkte

- **Erster Ansatz verworfen:** `iptables` mit einmal aufgelösten IPs scheiterte
  im Integrationstest an den rotierenden Adressen von `github.com`
  (Timeout trotz Allowlist); ersetzt durch `dnsmasq` + nftables-Set.
- **Guardrail:** Mit `NET_ADMIN` und `sudo` aufhebbar; nested Podman-Verkehr nicht erfasst.
- **Doctor** kann die Capability nicht vom Host aus prüfen (Spec 0.3.3).
- **Hosts des Agenten** kennt u-boot nicht; Nutzer tragen sie in `egress.allow` ein.
- Nicht geprüft: Colima/macOS, Podman-Host, VS Code/Codespaces.

## 7. Closure-Notiz (nach `done/`)

- **Sensoren:** `make gates` grün (Coverage 91,7 %); `go test -tags docker
  ./internal/e2e -run SandboxDevcontainer` grün (Egress dreimal in Folge
  stabil); `make ci` und voller `make test-docker` nicht ausgeführt.
- **Folge-Slices:** keine.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas GF (Generator, Templates, Doku), siehe Kurs Modul 5
§Worked Mini-Example.
