# ADR 0012: Devcontainer-Egress-Firewall (network-hardened devcontainer)

## Status

Proposed

> **Entwurf — noch nicht ratifiziert.** Festgehalten, damit die Idee
> nicht verloren geht (Roadmap-AP `slice-vN-devcontainer-egress-firewall`).
> Die §Entscheidung unten ist ein **Vorschlag**; die §Offenen Fragen
> müssen vor `Accepted` beantwortet werden. Kein Code, bis ratifiziert
> + Spec-Erweiterung + Slice-Plan stehen. Gleiche Klasse wie
> [ADR-0011](0011-agent-harness-scaffolding.md) (prospektives Feature,
> Produkt-Scope-Entscheidung).

## Datum

2026-06-09

## Kontext

`u-boot init --devcontainer` / `generate devcontainer` erzeugt heute
einen **minimalen** Devcontainer: `devcontainer.json` (`name`, `build`,
`forwardPorts`, `features`, `remoteUser: vscode`) + ein `Dockerfile`
(`FROM mcr.microsoft.com/devcontainers/base:debian`, non-root
`USER vscode`). **Keine Netzwerk-Restriktion** — kein
`postCreate`/`initializeCommand`, kein `runArgs`, kein
`iptables`/`ipset`. In `spec/lastenheft.md` kommt Firewall/Egress
nirgends vor (`LH-FA-DOC-003` „Netzwerk" meint das gemeinsame
*Compose*-Netzwerk, nicht Egress-Kontrolle).

Verbreiteter Pattern („network-hardened devcontainer"): ein
`init-firewall.sh`, das mit `--cap-add=NET_ADMIN` läuft und per
`iptables`+`ipset` nur eine **Allowlist ausgehender Ziele**
(z. B. GitHub, Paket-Registry, ggf. ein API-Endpoint) zulässt und den
Rest `DROP`t. Zweck: einen autonomen Agenten oder fremden Code im
Container eindämmen (keine Exfiltration, keine beliebigen Hosts).

**Fit mit u-boot:**

- Die Sicherheits-Philosophie existiert bereits: `LH-NFA-SEC-004`
  (keine verdeckte Fremd-Code-Ausführung) + `LH-FA-DEV-003`
  (`--allow-external-feature-sources`, explizite Allowlist für
  Feature-*Quellen*). Eine Egress-Firewall ist das **Runtime-Pendant**
  zur bestehenden **Build-Time/Supply-Chain-Allowlist** —
  komplementär, nicht doppelt.
- Deterministisch/template-bar: ein `.devcontainer/init-firewall.sh` +
  `runArgs` + `postCreateCommand` + ein Config-Key in `u-boot.yaml` —
  u-boots Wheelhouse (Template + Managed-Block + Config + Defaults).

**Zentrale Grenze:** Eine Egress-Firewall im Container ist ein
**Guardrail, kein Sandbox.** Ein In-Container-Prozess *mit* der
Capability kann die Regeln rückgängig machen; die Firewall begrenzt
nur den Egress, nicht den In-Container-Root. Das muss ehrlich als
Defense-in-Depth dokumentiert werden, nicht als harte
Sicherheitsgrenze.

## Entscheidung

**(Vorschlag — Status `Proposed`, noch nicht ratifiziert.)**

u-boot erzeugt eine Devcontainer-Egress-Firewall als **opt-in**:

1. **Opt-in, nicht Default.** Aktivierung über ein Flag
   (`generate devcontainer --firewall` o. ä.) oder einen Config-Key
   `devcontainer.firewall.enabled: true` — nicht Default, weil
   `NET_ADMIN` nicht überall verfügbar ist (siehe §Konsequenzen).
2. **Artefakte:** `.devcontainer/init-firewall.sh` (iptables+ipset,
   `default DROP` + Allowlist) als Managed-Block-Datei; in
   `devcontainer.json` `runArgs: ["--cap-add=NET_ADMIN"]` +
   `postCreateCommand`/`initializeCommand`, der das Script ausführt.
3. **Allowlist als Config:** `devcontainer.firewall.allow: [<host>…]`
   in `u-boot.yaml`, mit **sinnvollen Defaults je Ökosystem**
   (z. B. GitHub + die Registry des gewählten Service-/Sprach-Stacks).
4. **doctor-Check + graceful degradation:** `u-boot doctor` prüft, ob
   `NET_ADMIN` gewährbar ist; ist es das nicht, `warn` (nicht `error`)
   mit klarem Hinweis statt eines `up`-Abbruchs. Kein hartes Scheitern
   auf Umgebungen ohne die Capability.
5. **Engine/Format wie heute** (`text/template` + Managed-Block), kein
   neuer Stack.

## Konsequenzen

Positiv:

- Schließt die Runtime-Lücke neben der bestehenden Build-Time-Allowlist
  (`LH-FA-DEV-003`); konsistente Sicherheits-Story.
- Rein additiv, opt-in — kein Bruch für bestehende Devcontainer.
- Template + Config + doctor sind etablierte u-boot-Muster.

Negativ / Risiken:

- **`NET_ADMIN`-Abhängigkeit / Portabilität:** rootless Docker,
  manche CI-Runner, Docker-Desktop-Eigenheiten gewähren die Capability
  nicht → ohne Degradation bricht `up`. Der doctor-Check (Punkt 4) ist
  load-bearing.
- **Guardrail ≠ Sandbox** (siehe §Kontext) — Erwartungs-Management in
  der Doku Pflicht.
- **Allowlist-Pflege:** zu eng → Builds brechen (npm/pip/go proxy
  fehlt); zu weit → Schutz wertlos. Defaults je Ökosystem müssen
  sorgfältig kuratiert und dokumentiert werden.
- **Verifikation schwergewichtig:** „blockt die Firewall Host X" ist
  ein Integrationstest (`//go:build docker` + `NET_ADMIN`), kein
  Unit-Test — passt zur e2e-Harness, kostet aber mehr.
- **Spec-Wachstum:** neue `LH-FA-DEV-*` (und ggf. `LH-NFA-SEC-*`).

## Offene Fragen (vor `Accepted` zu beantworten)

1. **Default-Allowlist je Ökosystem:** welche Hosts pro Service-/
   Sprach-Stack (postgres/keycloak/otel; Node/Python/Go/Java)? Eine
   gemeinsame Basis (GitHub, ggf. Distro-Mirror) + stack-spezifische
   Ergänzungen?
2. **Degradations-Politik:** `warn` + trotzdem starten (Firewall
   inaktiv) vs. opt-in-`--require-firewall`, das ohne `NET_ADMIN`
   hart abbricht (Exit 11)?
3. **Aktivierungs-Surface:** Flag, Config-Key, oder beides; Interaktion
   mit `--no-interactive`/`--yes` (`LH-FA-CLI-005A`).
4. **iptables vs. nftables**, und Verhältnis zur Distro-Basis des
   Devcontainer-Image (`debian` → iptables-legacy/nft?).
5. **Verhältnis zu `LH-FA-DEV-003`:** geteilte Allowlist-Semantik/
   Config-Form oder bewusst getrennt (Build-Source vs. Runtime-Egress)?

## Folgepunkte

Dieses ADR liefert nur die Entscheidungs-Rahmung. Vor Implementierung:

- Ratifizierung (`Proposed` → `Accepted`) nach Klärung der §Offenen
  Fragen.
- Spec-Erweiterung: neue `LH-FA-DEV-*`-Anforderungen (+ ggf.
  `LH-NFA-SEC-*` für die Guardrail-Semantik).
- Slice-Plan `slice-vN-devcontainer-egress-firewall` in `open/` (heute
  nur als Roadmap-AP geführt, ohne Plan); doctor-Check als eigene
  Tranche (analog `slice-followup-devcontainer-features-drift-doctor`).

Re-Evaluation-Trigger: konkrete Nutzer-/Team-Nachfrage nach
egress-restringierten Devcontainern, oder ein Agenten-Sandbox-Use-Case
auf u-boot-erzeugten Devcontainern.
