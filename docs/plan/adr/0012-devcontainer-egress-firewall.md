# ADR 0012: Devcontainer-Egress-Firewall im Sandbox-Profil

**Status:** Accepted

**Datum:** 2026-09-30

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0014](0014-nested-podman-sandbox-devcontainer.md)

**Schärft:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) — Mechanismus, Default-Allowlist und Verhalten bei nicht gewährbarer Capability.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Ein autonomer Agent im Sandbox-Devcontainer ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) soll nicht beliebige Hosts erreichen (Exfiltration, unerwünschte Downloads). Verbreitetes Muster („network-hardened devcontainer“): ein Startscript mit `--cap-add=NET_ADMIN`, das ausgehenden Verkehr per `iptables`/`ipset` auf eine Allowlist beschränkt. Die Egress-Restriktion ist das **Runtime-Pendant** zur Build-Time-Allowlist für Feature-Quellen ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)).

**Grenze:** Die Restriktion ist ein **Guardrail, keine Sandbox-Grenze**. Ein Prozess im Container, der die Capability hat (und das `sudo` des Base-Images nutzen kann), kann die Regeln aufheben; der Verkehr verschachtelter Podman-Container läuft über `FORWARD`, nicht `OUTPUT`, und wird nicht erfasst.

**Messung (2026-09-30, Docker 29.8.1, `devcontainers/base:debian`, `nft` und `dnsmasq` 2.91).** Mit `--cap-add=NET_ADMIN`:

| Variante | Ergebnis |
|---|---|
| `iptables`-Regeln pro einmal aufgelöster IPv4-Adresse | `github.com` liefert pro DNS-Abfrage **wechselnde Adressen** (140.82.112.3, .113.3, .121.3, .121.4); die Verbindung scheitert, sobald `curl` eine andere Adresse als das Script nutzt (im Integrationstest in rund jedem zweiten Lauf Timeout trotz Allowlist) — **nicht tragfähig** |
| lokaler `dnsmasq` (nur erlaubte Namen, `nftset=` füllt ein nftables-Set mit den Antworten) plus nftables-Policy `drop` | sechs von sechs Abrufen von `github.com` erfolgreich; `example.com` ist **nicht auflösbar** und nicht erreichbar; das Set enthält nur die tatsächlich aufgelöste Adresse (Timeout 1 h) |
| ohne `NET_ADMIN` | `nft list tables` bricht mit Fehler ab: der Zustand ist im Startscript erkennbar |

Kein `ipset`-Kernelmodul nötig (nftables-Sets). Zusammen mit dem Nested-Podman-Profil ([ADR-0014](0014-nested-podman-sandbox-devcontainer.md)) bleibt `NET_ADMIN` zusätzlich gewährbar.

## Entscheidung

1. **Opt-in im Sandbox-Profil** über `devcontainer.sandbox.egress.enabled: true` (Default `false`); keine Flag-Variante (wie im Spec), daher keine Interaktion mit `--yes`/`--no-interactive`. Ohne Sandbox-Profil hat der Schlüssel keine Wirkung (Doctor-Warnung).
2. **Mechanismus:** DNS-gesteuerte Allowlist. Das Startscript `.devcontainer/egress-init.sh` startet einen lokalen `dnsmasq` (Paket `dnsmasq-base`), der **nur die erlaubten Namen** (samt Subdomains) an die ursprünglichen Resolver weiterleitet und die Antworten per `nftset=` in ein nftables-Set einträgt; `/etc/resolv.conf` zeigt auf `127.0.0.1`. Eine nftables-Kette `output` mit Policy `drop` lässt Loopback, bestehende Verbindungen, DNS zu den Upstream-Resolvern und Ziele im Set zu; IPv6 ist vollständig gesperrt. Die Allowlist ist im Script eingebacken. Ausführung per `postStartCommand` mit `sudo` (die Regeln überleben keinen Container-Neustart); die Clone-Phase läuft vorher im `postCreateCommand`.
3. **Capability:** `--cap-add=NET_ADMIN` in `runArgs`, in der Befehlsausgabe als Lockerung ausgewiesen (Code [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
4. **Default-Allowlist** (dokumentiert in `docs/user/devcontainer-sandbox.md`): gemeinsame Basis `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org`; bei `nestedRuntime: podman` zusätzlich `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io`; je aktiviertem Feature: `node` → `registry.npmjs.org`; `go` → `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com`; `java` → `repo.maven.apache.org`, `repo1.maven.org`. Der Host der Clone-Quelle kommt automatisch dazu. Hosts des Agenten selbst (z. B. dessen API) kennt u-boot nicht: sie stehen in `devcontainer.sandbox.egress.allow` (Liste gültiger Hostnamen, Schreibweise klein).
5. **Getrennt von `devcontainer.featureSources.allow`:** Build-Quelle und Laufzeit-Ziel sind eigene Schlüssel mit eigener Semantik (Entscheidung zu Frage 5 des Entwurfs).
6. **Degradation** nach der Tabelle in [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): ist `nft` nicht nutzbar (Capability nicht gewährt), entfällt die Restriktion mit Warnung (`onUnavailable: warn`) oder das Startscript endet mit Exit 11 (`fail`). `u-boot doctor` prüft die Konfiguration statisch (Schlüssel ohne Sandbox-Profil → `warn`); die Capability ist vom Host aus nicht zuverlässig bestimmbar und wird beim Containerstart geprüft.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Keine Egress-Kontrolle | keine Capability, einfach | Exfiltration und beliebige Downloads möglich |
| B — `iptables` + `ipset` im Container (verbreitetes Muster) | dynamische Sets | einmalige Auflösung verliert rotierende Adressen (gemessen: nicht tragfähig); `ip_set`-Kernelmodul auf Docker Desktop/Colima nicht garantiert |
| C — `iptables`-Regeln pro einmal aufgelöster IP | keine weiteren Pakete | rotierende Adressen (GitHub) brechen die Verbindung (gemessen: nicht tragfähig) |
| **C2 — `dnsmasq` + nftables-Set (DNS-gesteuert)** | deterministisch auch bei rotierenden Adressen, nicht erlaubte Namen sind nicht auflösbar, kein `ipset`-Modul | zwei zusätzliche Pakete (`dnsmasq-base`, `nftables`), `/etc/resolv.conf` wird im Container umgeleitet |
| D — Restriktion außerhalb des Containers (Netzwerk-/DNS-Ebene, Proxy) | kein `NET_ADMIN`, nicht vom Container aus aufhebbar | braucht Host-/Netzwerk-Konfiguration außerhalb von u-boots Scope (Nicht-Ziel: keine Host-Konfiguration) |

## Konsequenzen

- Positiv: Laufzeit-Pendant zur Build-Allowlist; additiv und opt-in; Template, Config und Doctor sind etablierte u-boot-Muster.
- Negativ: `NET_ADMIN` ist eine weitere Lockerung; Guardrail, nicht Grenze (siehe Kontext); die Allowlist braucht Pflege (zu eng bricht Builds, zu weit ist wertlos); IPv6 ist vollständig gesperrt; `/etc/resolv.conf` wird umgeleitet; nested Podman-Verkehr wird nicht erfasst; Doctor kann die Capability nicht vom Host aus prüfen.
- Folgepflicht: Doku (Grenzen, Default-Liste), Golden Cases und ein Docker-Integrationstest (erlaubter und gesperrter Host, Degradation).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Golden Case | Egress aus: keine Egress-Ausgabe; Egress an: `NET_ADMIN` in `runArgs`, Script mit Default- plus Nutzer-Hosts, `postStartCommand` | `make test` |
| Integrationstest (`//go:build docker`) | erlaubter Host erreichbar, gesperrter Host nicht, ohne Capability Exit 11/Warnung | `make test-docker` |

## Re-Evaluierungs-Trigger

Nutzerbericht über blockierte legitime Ziele (dann Proxy-Variante D), eine Umgebung, in der `nft` im Container nicht nutzbar ist, oder eine Entscheidung, Restriktionen auf Host-/Netzwerk-Ebene zu unterstützen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-06-09 | Entwurf (`Proposed`) | — |
| 2026-09-30 | Trigger „Agenten-Sandbox-Use-Case“ eingetreten; offene Fragen 2, 3 und 5 beantwortet | Lastenheft 0.3.0 |
| 2026-09-30 | Fragen 1 und 4 entschieden, auf MADR-Form gebracht, `Accepted`; Mechanismus gemessen (erste Variante mit einmaliger IP-Auflösung verworfen, `dnsmasq` + nftables-Set gewählt) | Vereinbarung mit dem Projektinhaber („alles fertig machen“) |
