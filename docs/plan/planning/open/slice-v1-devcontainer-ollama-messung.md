# Slice v1: Ollama im Sandbox-Devcontainer — Quelle wählen, messen, belegen (Katalog-Schlüssel `ollama`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — fünfter Slice der Welle (Messung, Quellenwahl, Folge-ADR; ändert weder Lastenheft noch Produktcode). Ihm folgen [`slice-cr-devcontainer-ollama`](slice-cr-devcontainer-ollama.md) (Lastenheft), [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md) (Umsetzung) und [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (keine verdeckte Ausführung fremder Skripte; bindet die Wahl der Quelle), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts des Modell-Pull); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Drittanbieter-Features im Katalog; Muster Messreihe mit der Dev-Container-CLI), [`ADR-0019`](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) und [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) (Muster gepinnte, geprüfte Lieferkette), [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Egress).
Anlass: Projektinhaber, 2026-10-18 („zu 4. ollama“; Rückfrage beantwortet mit „Feature: Ollama im Container“).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Default-Allowlist) — der Slice **liest** sie und ändert keine; die Nachführung folgt in [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md). Das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Mit gemessenen Belegen steht fest, **wie** Ollama (CPU-only) in den Sandbox-Devcontainer kommt, mit welcher Lieferkette, welchen Hosts im laufenden Container, wohin die Modelle gehen und ob der Server automatisch startet. Ergebnis: eine Folge-ADR (Quelle und Version gepinnt **nach** der Messung) und der mit Messwerten gefüllte Entwurf für den Lastenheft-CR. Zugesagt wird später nur, was hier gemessen ist.

### Gefundener Stand (gelesen am 2026-10-18)

**Im Repo:**

- `grep -ril ollama` über den Baum (ohne `.git`, `.harness`) trifft nur `docs/plan/planning/done/slice-v1-release-cut-v0.11.0.md` (Zeile: „Eine Welle zu Egress und Ollama; `buildah`.“ unter *Offen beim Projektinhaber*); Code, Spezifikation, Handbuch und ADRs kennen Ollama nicht.
- **Katalog:** `featureCatalogue()` in `internal/hexagon/application/devcontainer_features.go` ist eine Karte Schlüssel → `featureCatalogueEntry{source, defaultVersion, shortDesc, publisher}`; `source` ist eine OCI-Referenz ohne Schema. `collectDevcontainerFeatures` (gleiche Datei) rendert die aktivierten Einträge als `"<source>:<version>": {}` unter `features` in `devcontainer.json`; ein Name ohne Katalogeintrag und ohne `source:` wird dort übersprungen und an anderer Stelle mit Exit-Code `10` gemeldet (Kommentar zu `collectDevcontainerFeatures`). **Ein Eintrag, der kein OCI-Feature ist (Dockerfile-Block wie bei Podman und Docker), hat im Katalog heute keinen Platz.**
- **Drittanbieter:** `featureCatalogueEntry.publisher` kennzeichnet `homebrew` (`meaningful-ooo`) und `claude-code` (`anthropics`); `thirdPartyFeatureWarnings` (`internal/hexagon/application/devcontainer_features.go`) erzeugt daraus eine Warnung mit Code [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) ohne `SecurityRelaxation`, die unter `--quiet` entfällt ([`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)). Der Katalog verlangt für Einträge außerhalb des offiziellen Namensraums den Herausgeber als Registry-Namensraum.
- **Lieferkette Podman (Muster):** `Dockerfile.tmpl`, Block `NestedPodman`, pinnt Version, SHA-256 je Architektur und den Fingerprint des Signaturschlüssels als `ARG`, bricht vor dem Entpacken bei Abweichung ab und lädt Archiv und Signatur beim Image-Bau; [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) begründet das; [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) hält die Pin-Alterung als offene Beobachtung (Trigger: Pin-Hebung vor 2027-03-08).
- **Egress zur Laufzeit, Image-Bau unbeschränkt:** Die Restriktion greift erst im laufenden Container (`egress-init.sh` im `postStartCommand`); der Image-Bau läuft auf dem Container-Host ohne sie ([`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), Absatz zu Podman: „der Build läuft auf dem Container-Host ohne die Egress-Restriktion, braucht diese Hosts also offen, und im laufenden Container ist keiner davon nötig“; gleiche Aussage im Slice [`slice-v1-devcontainer-feature-claude-code`](../done/slice-v1-devcontainer-feature-claude-code.md) §3). **Folge für den Wunsch des Projektinhabers „Hosts für Installation und Modell-Pull in die Default-Allowlist“:** Hosts der **Installation** gehören nicht in die Laufzeit-Allowlist, sofern die Installation beim Image-Bau läuft; nur die Hosts des **Modell-Pull** (und was der Server zur Laufzeit sonst braucht) tun es. Die Messung belegt, welche das sind.
- **Loopback ist offen:** `egress-init.sh.tmpl` lässt `oifname lo` zu (Regel `nft add rule ip u_boot_egress output oifname lo accept`); ein Server auf `127.0.0.1` im selben Container ist von der Restriktion nicht betroffen.

**Im Netz (gelesen am 2026-10-18; Lesen, nicht Messen — die Messung unten prüft es):**

- Das Verzeichnis `https://containers.dev/static/devcontainer-index.json` führt **genau ein** Feature mit Namen `ollama`: `ghcr.io/prulloac/devcontainer-features/ollama`, Version `1.0.1`, Option `pull` („comma separated list of models to pull“). Quelle dort: `src/ollama/install.sh` und `devcontainer-feature.json` im Repository `prulloac/devcontainer-features` (Branch `main`): die `install.sh` ruft `curl -fsSL https://ollama.com/install.sh | sh` auf und lädt bei gesetzter Option Modelle (`ollama serve &`, `sleep 3`, `ollama pull`); die Metadaten tragen `"entrypoint": "ollama serve"` und `installsAfter: …/common-utils`. **Das ist ein ungeprüftes `curl | sh`** und ein `entrypoint` im Feature, beides gegen die Vorgaben des Projektinhabers; ob das Feature auf Debian 13 (`mcr.microsoft.com/devcontainers/base:debian`) baut, ist ungemessen.
- `https://ollama.com/install.sh` (455 Zeilen): lädt `ollama-linux-<arch>` von `https://ollama.com/download`, verlangt `zstd`, richtet `systemd` oder einen Dienst ein und enthält Zweige, die bei erkannter NVIDIA-GPU CUDA-Repositories und -Treiber installieren; eine Prüfung von SHA-256 oder Signatur habe ich im Skript nicht gefunden (Suche nach `sha256`, `gpg`, `cosign`, `.asc`: keine Treffer).
- Upstream-Release `v0.40.6` (GitHub-API `repos/ollama/ollama/releases/latest`, veröffentlicht 2026-10-16) bietet unter anderem `ollama-linux-amd64.tar.zst`, `ollama-linux-arm64.tar.zst` (daneben `-rocm`, `-mlx`, `-jetpack`), `install.sh` und `sha256sum.txt`; **eine Signaturdatei (`.asc`, `.sig`) steht nicht in der Asset-Liste.** Ob das Release eine Build-Provenance trägt (zum Beispiel `gh attestation verify`), ist ungemessen.

### Quellen-Alternativen (der Architect misst und entscheidet; Einordnung des Planners, nicht entscheidend)

| Opt. | Quelle | Pro | Contra |
|---|---|---|---|
| A | Community-Feature `ghcr.io/prulloac/devcontainer-features/ollama:1` als Katalogeintrag (Muster `homebrew`, [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)) | kein Eingriff in den Katalog-Typ; Warnung und Herausgeber-Regel vorhanden | `curl \| sh` ungeprüft (Projektinhaber: „kein ungeprüftes `curl \| sh`“); `entrypoint` im Feature (startet `ollama serve` bei Containerstart, am Start-Script der Egress-Kette vorbei); ein Einzelmaintainer; Version `1.0.1` schwimmt |
| B | Eigener Block in `Dockerfile.tmpl` mit dem Release-Archiv (`ollama-linux-<arch>.tar.zst`), Version und SHA-256 je Architektur gepinnt, `zstd` aus Debian; aktiviert durch den Schlüssel `ollama` (Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)) | gepinnte, geprüfte Lieferkette ohne Fremdskript; Upstream-Quelle; kein Drittanbieter-Feature | **keine Signatur gefunden** (nur SHA-256 aus demselben Release; [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) und [ADR-0023](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) verlangten zusätzlich Schlüssel-Pin bzw. Signaturprüfung); der Katalog bekommt einen zweiten Eintragstyp (kein OCI-Feature); Pin-Alterung wie bei Podman |
| C | `COPY --from=ollama/ollama@sha256:…` aus dem Docker-Hub-Image (Mehrstufen-Build) | Digest-Pin | Binärkompatibilität des Image-Inhalts auf Debian 13, Größe, Basis und Host des Image-Baus ungemessen; Registry-Zugriff beim Bau |
| D | Debian-Paket `ollama` | `apt` als Update-Pfad | Existenz in Debian 13 ungemessen (`apt-cache policy ollama` im Basis-Image) |
| E | Ein von `u-boot` erzeugtes lokales Feature im Projekt (`.devcontainer/features/ollama/`, [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): „lokal hinterlegte … Features erlaubt“) | volle Kontrolle über Installation und Pin; kein Fremdskript | neuer Generator-Zweig und eine zweite Wahrheit neben dem Katalog; Aufwand |

**Offen für den Projektinhaber** (Sicherheitsentscheidung, nicht Messung; hier dargelegt, nicht gefragt): ob ein SHA-256-Pin aus derselben Release-Quelle ohne Signatur genügt (B), oder ob die Lieferkette erst mit einer Signatur oder Build-Provenance tragen muss; und ob ein Drittanbieter-`curl | sh` unter irgendeiner Bedingung zulässig ist (A). Die ADR legt die Optionen mit der Messung dar; die Entscheidung trägt der Projektinhaber.

### Messfragen (Arbeitsplan; Dev-Container-CLI im Wegwerf-Container auf Colima wie bei `sshd` und `claude-code`: [`slice-v1-devcontainer-feature-claude-code`](../done/slice-v1-devcontainer-feature-claude-code.md) §3; Basis-Image `mcr.microsoft.com/devcontainers/base:debian`)

| # | Frage | Warum |
|---|---|---|
| O1 | Welche der Optionen A bis E bauen auf dem Basis-Image (arm64), mit welcher Größe des Images, welchem Eigentümer und welchen Rechten der Binärdatei? Für A: Feature-Metadaten (`entrypoint`, `mounts`, `capAdd`, `privileged`) im OCI-Manifest, Einfluss des `entrypoint` auf den Container-Start | Quellenwahl; das Profil darf nicht unterlaufen werden (wie bei `claude-code`) |
| O2 | Gibt es zu den Release-Archiven eine Signatur oder Build-Provenance (`gh attestation verify`, Sigstore)? Stimmt `sha256sum.txt` mit den Archiven überein? | trägt oder verwirft die Prüfung in B |
| O3 | Laufzeit unter der erzeugten Egress-Restriktion: `ollama serve`, `ollama pull <kleines Modell nach Wahl des Messenden>`, `ollama run` mit einem Prompt, CPU-only; Dauer, RAM, Plattenbedarf; Bindung des Servers (Loopback oder alle Adressen), Standard-Port | Zusage nur für Gemessenes; entscheidet, ob eine Ausweisung eines Netzdienstes nötig wird (Vergleich `sshd`) |
| O4 | DNS-Mitschnitt (protokollierender Resolver, Verfahren aus dem `claude-code`-Slice) für Pull und Run: **welche Hosts** braucht der Server zur Laufzeit, auch wechselnde Subdomains eines Dienstes (Blob-Speicher)? | Default-Allowlist ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)); ein Subdomain-Fall wird an [`slice-v1-egress-wildcard-pruefung`](../in-progress/slice-v1-egress-wildcard-pruefung.md) zurückgereicht |
| O5 | Modellspeicher: Standardpfad des Servers für den Benutzer `vscode`, Größe typischer Modelle, Verhalten ohne Volume (Container-Schicht), mit benanntem Volume (Eigentümer, Muster `…-containers`-Volume bei Podman: `devcontainer.json.tmpl`) | Volume-Zusage und Platzbedarf |
| O6 | Automatischer Start: (i) Feature-`entrypoint`, (ii) `postStartCommand` hinter dem Egress-Script (`sudo sh …/u-boot-egress-init && …`, `devcontainerTemplateData`), (iii) kein Auto-Start, Nutzer startet selbst; Verhalten bei `onUnavailable: fail` (Exit 11) und bei Neustart des Containers | Entscheidung zu „Auto-Start offen“ |
| O7 | Version und Pflege: welche Version wird gepinnt, wie wird die Pin-Hebung sichtbar (Muster [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md): Trigger mit Datum) | Re-Evaluierungs-Trigger der ADR |

**Grenzen vorab:** Nur arm64 unter Colima messbar; amd64 erst über CI (siehe [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md): ein CI-Lauf ohne `-v` belegt keinen Pass). GPU-Beschleunigung wird nicht gemessen (CPU-only, Entscheidung des Projektinhabers). Modell- und Rechenbedarf hängt an der Größe der Colima-VM.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Tests, Spezifikation, Handbuch, Lastenheft** — [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md), [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md) und [`slice-cr-devcontainer-ollama`](slice-cr-devcontainer-ollama.md); ein Slice, der misst und entscheidet, setzt nicht um (Rollentrennung, `harness/roles.md`).
- **Anschluss an den Ollama des Hosts** (Host-Dienst, `host.docker.internal`, Socket-Mount) — vom Projektinhaber ausgeschlossen; widerspräche zudem [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (keine Host-Verbindung im Profil).
- **GPU-Beschleunigung und Grafik-Treiber** — Entscheidung des Projektinhabers: CPU-only; die `install.sh`-Zweige für NVIDIA bleiben unbenutzt.
- **Auswahl, Vorbelegung oder Verwaltung von Modellen durch `u-boot`** (Option `pull` des Community-Features) — der Katalog kennt keine Optionen ([`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md), Re-Evaluierungs-Trigger Optionsmodell); der Nutzer holt Modelle im Container.
- **Weitere Agenten im Katalog** — der Projektinhaber nannte außer Ollama keinen; ohne Nennung kein Slice.
- **Ollama als Compose-Add-on-Service** — der Projektinhaber wählte „Feature: Ollama im Container“.

## 2. Definition of Done

- [ ] **Messung:** O1 bis O7 sind gemessen und mit Messstand (Docker-, Colima-, Dev-Container-CLI-Version, Architektur, Basis-Image-Digest, Ollama-Version) und Grenzen in §3 und in der Evidence festgehalten; die Kommandos stehen wörtlich dabei (Messprotokoll **im Repo ablegen**, nicht nur das Ergebnis — das Skript bleibt außerhalb, die Grenze ist ausgewiesen).
- [ ] **Folge-ADR:** Eine ADR (nächste freie Nummer; Index `docs/plan/adr/README.md`; `Proposed` bis zur Entscheidung des Projektinhabers zur Lieferkette, danach `Accepted`) entscheidet Quelle (A bis E), Pin von Version und Prüfsumme(n) **nach der Messung**, Prüfung (SHA-256, Signatur oder Provenance, soweit vorhanden), Hosts der Default-Allowlist (gemessen), Modellspeicher, Auto-Start, Kennzeichnung (Warnung, nicht `SecurityRelaxation`, falls ein Drittanbieter-Feature gewählt wird) und Re-Evaluierungs-Trigger; sie trägt den **Wortlaut-Entwurf** für den CR, mit Messwerten gefüllt.
- [ ] **Übergabe:** [`slice-cr-devcontainer-ollama`](slice-cr-devcontainer-ollama.md) und [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md) sind gegen das Ergebnis geprüft und, wo nötig, in ihren Plan-Tabellen berichtigt (Variante A oder B/E); eine Beobachtung zur Pin-Alterung ist im Register angelegt oder in §7 begründet nicht.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: ADR-Index; kein Handbuch, kein CHANGELOG (kein Produktverhalten).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00NN-ollama-im-sandbox-devcontainer.md` (nächste freie Nummer) und `docs/plan/adr/README.md` | neu / update | Folge-ADR: Quelle, Pin, Prüfung, Hosts, Volume, Auto-Start, Trigger, Wortlaut-Entwurf |
| `docs/plan/planning/observations/BEO-APP/…` | neu (bedingt) | Beobachtung zur Pin-Alterung des Ollama-Pins, falls B oder C gewählt wird (Muster `podman-static-pin-alterung`) |
| `docs/plan/planning/open/slice-cr-devcontainer-ollama.md`, `docs/plan/planning/open/slice-v1-devcontainer-feature-ollama.md` | update | Plan-Tabellen und Wortlaut gegen das Messergebnis berichtigt |
| `docs/reviews/` (Review-Report, Verification Evidence, Messprotokoll) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md`; Messprotokoll mit Kommandos |
| diese Slice-Datei | update | Messung (§3), Closure (§7) |

Kein Produktcode und kein Test in diesem Slice (`internal/**`, `spec/**` unverändert). Abgleich der Tabelle gegen `git diff --stat` nach der **letzten** Fix-Runde (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1). Er hängt von keinem Egress-Slice der Welle ab; die Reihenfolge der Welle stellt ihn hinter die Egress-Slices, damit ein Subdomain-Fall aus O4 an einen geschlossenen Wildcard-Slice zurückgereicht wird (Welle §5).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung braucht je Quelle eine eigene Messreihe mit eigenem Entscheidungsbedarf; dann ein Slice je Quelle.
- `in-progress` → `open` (blockiert): die Messumgebung steht nicht zur Verfügung (Colima, Dev-Container-CLI, Netz zu `ollama.com`); oder die Quelle verlangt eine Lockerung des Profils (zusätzliche Capability, Device, Mount), die der Projektinhaber entscheiden muss — dann zuerst eine Sicherheitsentscheidung, wie bei [`slice-v1-devcontainer-podman-build`](../done/slice-v1-devcontainer-podman-build.md).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die ADR trägt ihren Status, die Entscheidung des Projektinhabers zur Lieferkette ist im Gespräch belegt und in §7 vermerkt.

## 6. Risiken und offene Punkte

- **Keine Quelle erfüllt „gepinnt, geprüft, kein `curl | sh`“**, weil das Release keine Signatur trägt (O2). — **Ausgang:** offen bis zur Messung; eingetreten wäre: die ADR legt die Auslegung (SHA-256 ohne Signatur) dem Projektinhaber vor; ohne Entscheidung bleibt `ollama` ungebaut (Rückführung nach `open/`).
- **Ollama ist ein Netzdienst im Container** (Port, Bindung) und braucht eine Ausweisung wie `sshd`. — **Ausgang:** offen bis O3; gebunden an Loopback wäre das entfallen, sonst Wortlaut im CR.
- **Modell-Blobs liegen unter wechselnden Subdomains**; die Allowlist kennt dann nur den Eltern-Namen (breit) oder scheitert. — **Ausgang:** offen bis O4; eingetreten: Rückgabe an [`slice-v1-egress-wildcard-pruefung`](../in-progress/slice-v1-egress-wildcard-pruefung.md) (Welle §5), keine stille Erweiterung der Default-Allowlist.
- **CPU-only ist langsam und speicherhungrig**; ein Modell, das in die Colima-VM nicht passt, lässt die Messung scheitern. — **Ausgang:** offen; die Messung nennt das gewählte Modell und die VM-Größe, die Zusage ist darauf begrenzt.
- **Versions-Pin altert** (Ollama veröffentlicht häufig; `v0.40.6` ist der Stand vom 2026-10-16). — **Ausgang:** weiter offen: Register-Eintrag und Trigger in der ADR (Muster `podman-static`).
- **Die Messung ist im Repo nicht wiederholbar** (kein Skript; Grenze wie bei `sshd`, `claude-code` und `podman build`). — **Ausgang:** weiter offen als bekannte Grenze; das Messprotokoll mit Kommandos liegt in `docs/reviews/`; ob ein Make-Target gewollt ist, bleibt eine Entscheidung des Projektinhabers ([`welle-devcontainer-sshd-claude-code`](../done/welle-devcontainer-sshd-claude-code.md) §7, offene Entscheidung (c)).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Messergebnis O1–O7 mit Messstand, gewählte Quelle, ADR-Nummer und Status, Entscheidung des Projektinhabers, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/plan/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); gemessen wird in Wegwerf-Containern, nicht in einer Sub-Area. Die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) (offen, ein Beleg; ein zweiter gepinnter Binär-Download im Dockerfile, falls B gewählt wird: gleiche Klasse, Zähler wäre zwei, wenn der Beleg dort ergänzt wird), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen: CI-Beleg für die Messung fehlt), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert; Aussagen über Netz-Quellen oben sind als „gelesen, nicht gemessen“ markiert). Evidenz-/Diskrepanz-Risiko: mittel (Lieferkette und Hosts sind ungemessen).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
