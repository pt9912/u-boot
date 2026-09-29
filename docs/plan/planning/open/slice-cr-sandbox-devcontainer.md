# Slice CR: Sandbox-Profil für Devcontainer (autonome Agenten)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`done/slice-cr-adr-format-madr`](../done/slice-cr-adr-format-madr.md)).
> Anlass: autonomer Agenten-Einsatz (deaktivierte Rückfrage) auf
> u-boot-erzeugten Devcontainern. Der CR ist zugleich das **Trigger-
> Ereignis**, das [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md)
> in seinem Re-Evaluierungs-Trigger nennt („Agenten-Sandbox-Use-Case
> auf u-boot-erzeugten Devcontainern").

**Lifecycle:** Zustand = Verzeichnis (`open/` → `next/` → `in-progress/`
→ `done/`), Wechsel nur per `git mv`.

**Welle:** ohne Welle (CR am Vertrags-Stratum).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md) bis [`LH-FA-DEV-009`](../../../../spec/lastenheft.md) (neu),
[`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte)
(Ergänzung),
[`LH-NFA-SEC-003`](../../../../spec/lastenheft.md#lh-nfa-sec-003--sichere-defaults),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md).

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel

Lastenheft-Änderung **0.2.0 → 0.3.0** und anschließende Umsetzung eines
opt-in-Sandbox-Profils für Devcontainer. Der CR wurde gegen den
Repo-Stand geprüft (Review 2026-09-29): Versionsfeld steht auf 0.2.0,
die IDs -006 bis -009 sind frei (§4.3 endet bei
[`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports)),
die Motivation deckt sich mit den tatsächlichen Templates (Build aus
Repo-Kontext, `remoteUser: vscode`, kein Netzwerkschutz).

### CR-Kern (vorgeschlagene Anforderungen)

- **[`LH-FA-DEV-006`](../../../../spec/lastenheft.md) — Sandbox-Profil (V1):** opt-in über
  `u-boot generate devcontainer --sandbox` bzw. `devcontainer.profile: sandbox`.
  Nicht-root-Benutzer ([`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte)),
  kein Bind-Mount des Host-Arbeitsverzeichnisses (Workspace in benanntem
  Volume, Clone im Container), kein Socket-Mount, kein `--privileged`,
  keine Host-Dateien mit Secrets.
- **[`LH-FA-DEV-007`](../../../../spec/lastenheft.md) — Container-Runtime im Sandbox-Devcontainer (V1):**
  optional rootless Podman im Container
  (`devcontainer.sandbox.nestedRuntime: podman | none`, Default `none`),
  docker-Kompatibilitäts-Alias, `/dev/fuse`, subuid/subgid, Storage-Volume.
  Engine-neutral; Sicherheitslockerungen (Seccomp/AppArmor) benennen und
  in der Befehlsausgabe ausweisen.
- **[`LH-FA-DEV-008`](../../../../spec/lastenheft.md) — Egress-Restriktion (V2):** Allowlist für ausgehenden
  Verkehr; Guardrail, keine Sandbox-Grenze; Doctor-Warnung, wenn die
  Capability nicht gewährbar ist (kein Abbruch).
- **[`LH-FA-DEV-009`](../../../../spec/lastenheft.md) — Git-Zugangsdaten (V1):** Zugangsdaten weder im Image,
  im Volume noch in `u-boot.yaml`; Laufzeit-Übergabe per Env oder
  read-only-Secret-Mount; Doku zu kurzlebigen, repo-begrenzten Tokens
  (Private Key nie im Container) inkl. Branch-Protection-Forderung;
  Doctor-Check auf Token-Quelle (warn).
- **Ergänzung [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte):** UID des Container-Benutzers an den Host
  anpassbar (z. B. macOS `501` unter Colima).

**Nicht-Ziele:** u-boot startet keinen Agenten und setzt keinen
Berechtigungsmodus; keine GitHub-Apps/Tokens/Host-Dienste; keine
Host-Colima-/Podman-Konfiguration; Schadensbegrenzung statt harter
Isolation.

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change):**

1. `--sandbox` erzeugt Dateien ohne `--privileged`, ohne Socket-Mount,
   ohne Bind-Mount des Repos, ohne Host-Secrets.
2. Mit `nestedRuntime: podman` läuft `docker build` im Container ohne
   Docker-Socket.
3. Dieselbe Konfiguration startet unter Docker (inkl. Colima) und unter
   Podman (Integrationstest, `//go:build docker`).
4. Zugangsdaten tauchen in keiner erzeugten Datei auf.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.3.0 | (Annahmedatum) | Neue Anforderungen LH-FA-DEV-006 bis -009
(Sandbox-Profil, Container-Runtime, Egress-Restriktion,
Git-Zugangsdaten); Ergänzung von LH-FA-DEV-004. | Vereinbarung mit dem
Projektinhaber |
```

## 2. Definition of Done

- [ ] **CR-Charakter dokumentiert** (dieser Slice): Trigger, betroffenes
  Vertrags-Stratum, Resolutions der offenen CR-Fragen. Kein stiller
  Spec-Edit.
- [ ] **[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md)-Kohärenz** (Review M1): Der CR-Bezug zu
  [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) ist im ADR
  (Status/Re-Evaluierungs-Trigger/Geschichte) festgehalten; `-008` trägt
  Priorität V2 und hängt an der Ratifizierung des ADR. Die ADR-offenen
  Fragen (Default-Allowlist, Degradations-Politik, Aktivierungs-Surface,
  iptables vs. nft, Verhältnis zur
  [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Allowlist)
  sind im Zuge des Spec-Changes beantwortet oder explizit an den
  Umsetzungs-ADR delegiert.
- [ ] **Spec-Änderung 0.3.0** (drei Spuren nach §16): Versions-Feld,
  Historie-Zeile, neue Anforderungen in §4.3 inkl. Prioritäts-Markierung
  (V1/V2) und Ergänzung von
  [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte).
- [ ] **Exit-Code-Vertrag** (Review M2): Fehlerpfade der neuen Oberfläche
  (`--sandbox`, Config-Keys, Doctor-Checks) sind auf die
  [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)-Kategorien
  abgebildet (u. a. fachlicher Fehler → Exit 10, Umgebungsproblem → Exit
  11); Test-Pins für relevante Sentinels im Umsetzungs-Slice.
- [ ] **Degradationsverhalten definiert** (Review M3): definiertes
  Verhalten für (a) fehlendes `/dev/fuse`, (b) nicht gewährbare
  Egress-Capability (NET_ADMIN), (c) durch Seccomp/AppArmor blockierte
  nested User-Namespaces — je Pfad warn vs. hart abbrechen, im
  Anforderungstext oder Umsetzungs-ADR.
- [ ] **Traceability-Matrix** (Review L1): [`PH-DEV-006`](../../../../spec/lastenheft.md)..[`PH-DEV-009`](../../../../spec/lastenheft.md) und
  [`TC-DEV-006`](../../../../spec/lastenheft.md)..[`TC-DEV-009`](../../../../spec/lastenheft.md)-Zeilen in §13.
- [ ] **Quellen-Gating-Verhältnis** (Review L2): explizit geklärt, ob/wie
  das Sandbox-Profil Devcontainer-Features ergänzt und dass dann
  [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)/
  [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  greifen.
- [ ] **`-004`-Ergänzung in Vertrags-Form** (Review L3): konkreter
  Mechanismus der UID-Anpassbarkeit (Config-Key und/oder Build-Arg),
  nicht nur die Absicht.
- [ ] **Technische Verifikation** (CR-Fragen 2 und 3): NET_ADMIN vs.
  rootless User-Namespaces und Seccomp/AppArmor-Lockerungen unter Docker
  sind getestet; Befund ist im Umsetzungs-ADR festgehalten (auch
  Negativbefunde).
- [ ] `make docs-check` grün (nächstgelegener Sensor für den
  Spec-Change; Umsetzung des V1-Pakets erfolgt als Folge-Slice mit
  `make gates`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | §4.3: Anforderungen -006..-009 (Prioritäten V1/V1/V2/V1), Ergänzung -004; §13: PH/TC-Zeilen; Versions-Feld 0.3.0; §16: Historie-Zeile |
| `docs/plan/adr/0012-devcontainer-egress-firewall.md` | update | CR als Re-Evaluierungs-Trigger-Ereignis eintragen; Ratifizierungspfad für -008 klären (ADR bleibt `Proposed`, bis seine offenen Fragen beantwortet sind) |
| Umsetzungs-ADR (neu) | neu | Base-Image für Podman-im-Container (CR-Frage 1) + nested-Runtime-Design; erst nach Spec-Annahme |
| `docs/user/`, `README`, `CHANGELOG` | update | Erst bei Umsetzung des V1-Pakets (öffentliche Verträge) |

## 4. Trigger

- **`open` → `next`:** Projektinhaber priorisiert den CR (Einordnung in
  eine Iteration).
- **`next` → `in-progress`:** Beginn des Spec-Changes; [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md)-Kohärenz
  ist Voraussetzung (M1), da -008 im selben Versionspunkt landet.
- **Umsetzung:** das V1-Paket (-006, -007, -009 plus -004-Ergänzung)
  bekommt einen eigenen Umsetzungs-Slice; -008 folgt nach
  [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md)-Ratifizierung als eigener Slice (V2).
- **Rückführungen:** `in-progress` → `next` bei Zerlegungsbedarf (Spec-
  Change vs. Umsetzung trennen); `in-progress` → `open` bei Blockade
  (z. B. wenn die Seccomp-Verifikation das Default-Profil als nicht
  vertretbar erweist — dann CR-Teilrevision).

## 5. Closure-Trigger

DoD vollständig + Spec-Change committet (drei Spuren) +
[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md)-Bezug
nachgetragen + `make docs-check` grün. Umsetzung ist bewusst **kein**
Closure-Kriterium dieses Slices (Prioritäts-Lage: Spec zuerst,
V1-V2-Pakete als Folge-Slices).

## 6. Risiken und offene Punkte

- **NET_ADMIN vs. rootless Podman im Container** (CR-Frage 2,
  unverifiziert): nested rootless Runtime braucht User-Namespaces,
  Egress-Firewall will NET_ADMIN. Realistischer Fallback: Egress-Präferenz
  auf Netz- bzw. DNS-Ebene statt Container-iptables — Entscheidung liegt
  beim Umsetzungs-ADR bzw. [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md).
- **Seccomp/AppArmor unter Docker** (CR-Frage 3, unverifiziert): falls
  nested User-Namespaces ohne Lockerungen nicht klappt, stellt sich die
  Frage nach der Vertretbarkeit für ein Default-Profil — dann
  CR-Teilrevision (siehe Rückführung).
- **Base-Image** (CR-Frage 1): bewusst zurückgestellt auf
  Umsetzungs-ADR; der Spec-Text bleibt image-agnostisch.
- **Spec-Wachstum:** vier neue Anforderungen in einem Verschnitt; Gegen
  eine Aufblähung hält der CR den Vertragscharakter (Verhalten, nicht
  Implementierungsdetails) — im Spec-Change durchhalten.
- **Carveout-Potential:** falls der Spec-Change Teile zurückstellt (z. B.
  Doctor-Checks), Carveout-Eintrag in
  [`carveouts.md`](../in-progress/carveouts.md) mit Plan-Anker ergänzen.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Alle berührten Sub-Areas GF (siehe Kurs Modul 5 §Worked
Mini-Example).
