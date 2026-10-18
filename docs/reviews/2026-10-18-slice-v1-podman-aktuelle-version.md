# Review-Report: slice-v1-podman-aktuelle-version — 2026-10-18

**Review-Art:** Code — Code-Review gegen Plan, ADR und Hard Rules (Schwerpunkt Lieferkette und Sicherheit der Dockerfile-Vorlage).

**Gegenstand:** Commit-Range `314b1e9..HEAD` (8 Commits, 11 Dateien) für den Slice `slice-v1-podman-aktuelle-version`. <!-- d-check:ignore (Kennungen als Zitat im Report) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `6d55655` · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-podman-aktuelle-version` (in-progress), §3 Plan-Tabelle, §5 Messbefunde <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- ADR-0023 (Accepted 2026-10-18) samt Anhang; `ADR-0014`, `ADR-0019`, `ADR-0021` (Kontext, unverändert) <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- `LH-FA-DEV-007`, `LH-FA-DEV-008`, `LH-NFA-SEC-004`; Spezifikation `LH-FA-DEV-007.b`, `SPEC-019`, `SPEC-045` <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- `AGENTS.md` §3.1 bis §3.10, §5

**Eigene Läufe (Docker-only, Image `u-boot:test-docker-tools`):**

- Archive `podman-linux-arm64.tar.gz` und `podman-linux-amd64.tar.gz` (v6.1.4) aus dem Netz geladen: beide SHA-256 stimmen mit den `ARG`s im Template und mit ADR-0023 überein. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- `tar -tvz` beider Archive: 67 bzw. 69 Einträge, nur Pfade unter `usr/local/` und `etc/containers/` sowie `README.md`; keine absoluten Pfade, kein `..`; drei Symlinks, alle relativ und innerhalb `usr/local` (`pasta` → `passt`, zwei `quadlet`-Generatoren); keine Hardlinks; kein setuid/setgid/sticky-Bit; alle Einträge `runner/runner` (daher `--no-same-owner` wirksam nötig).
- `gpg --import` des Schlüssels von `keys.openpgp.org` und `gpg --status-fd 1 --verify`: `VALIDSIG` trägt den Primär-Fingerprint `0CCF102C…B503` in **Feld 12** (`$12` ist korrekt); Feld 3 ist der Unterschlüssel `6E8EF2F5…0BA4`.
- Release Notes `containers/podman` v6.0.0 gelesen: „Support for running on cgroups v1 systems has been removed“ — die Handbuch-Aussage trägt.
- `make doc-immutable RANGE="314b1e9..HEAD"`: 8 Commits, OK. `make docs-check` siehe Verdikt.
- Unit-Mutationen (`go test -run 'SandboxPodman|SandboxEgressWithPodman'`), je danach `git checkout -- <Datei>`: `--no-same-owner` entfernt (gefangen), `registries.conf` mit zweiter Registry (gefangen), `sha256sum -c -` durch `true` ersetzt (gefangen), `mount_program` aus dem Drop-in entfernt (gefangen), stiller Architektur-Fallback statt `exit 1` (**nicht gefangen**, F-2).
- E2E-Mutationen (`make test-docker`-Rezeptur, Einzeltests): `VALIDSIG`-Vergleich durch `true` ersetzt — `…_PinMismatchBreaksBuild/signing_subkey_as_the_pin` schlägt fehl (gefangen); `amd64`-/`arm64`-Zweig der Prüfsumme vertauscht — `…_NestedPodman_Build` bricht am Image-Build mit `FAILED` (gefangen). Unmutiert: `…_NestedPodman` grün (33,6 s).
- Host danach: Arbeitsbaum sauber, keine `sbxe2e-*`-Container oder -Volumes, keine Probe-Images des Reviews, `alpine:3` gecacht.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Plan-Tabelle §3 führt `docs/plan/planning/open/README.md` als „update“, der Diff berührt die Datei nicht; die Zeile „`README.md` Zeile 224 die Host-Mindestversion“ trifft im Baum Zeile 225 (`Docker Compose ≥ 2.20.0`). | `AGENTS.md` §3.10 | `slice-v1-podman-aktuelle-version.md` · „`docs/plan/planning/open/README.md` \| update \| Bestandszeile (dieser Plan)“ | ja — `git diff --stat 314b1e9..HEAD` gegen §3 | Plan-Tabelle gegen Diff nicht abgeglichen | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-2 | LOW | Der Zweig `*) echo "podman-static has no archive for $arch" >&2; exit 1` ist durch keinen Test gepinnt: ersetzt durch `*) sha256="${PODMAN_SHA256_ARM64}" ;;` bleiben Unit- und E2E-Tests grün. Auf anderen Architekturen als `amd64`/`arm64` (etwa `armhf`) wäre der Fehler damit nicht abgesichert. | ADR-0023 (Entscheidung 2, Fail-closed der Prüfkette) | `Dockerfile.tmpl` · `has no archive for $arch`; `generate_sandbox_test.go` · `assertPodmanStaticDockerfile` | ja — `make test` mit der Mutation (eigener Lauf: nicht gefangen) | Fail-closed-Zweig ohne Pin-Test | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-3 | LOW | Handbuch und Spezifikation nennen, dass Archiv und Schlüssel beim Image-Build geladen werden, aber nicht die dafür nötigen Hosts (`github.com`, `objects.githubusercontent.com`, `keys.openpgp.org`); nur ADR-0023 (Entscheidung 6) nennt sie. Ein Build hinter einer Firewall oder ohne Netz bricht damit ohne Hinweis im Nutzerhandbuch. | ADR-0023 (Entscheidung 6); `AGENTS.md` §5 Zeile 5 | `docs/user/benutzerhandbuch.md` · „Archiv und Schlüssel lädt der Bau des Images“ | nein — Review-only | Build-Zeit-Hosts nur in der ADR | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-4 | INFO | Nach Ablauf des Schlüssels (2027-03-08) meldet `gpg --verify` laut ADR-0023 `EXPKEYSIG` mit Exit-Code 0 und unverändertem `VALIDSIG`; die Vorlage erzwingt weder Ablauf noch Widerruf, die Integrität trägt die gepinnte SHA-256. Handbuch, Spezifikation und CHANGELOG nennen den Ablauf nicht (nur die ADR als Negativfolge). Für den Verifier/Planer: Pin-Hebung vor 2027-03-08 ist eine Wiedervorlage-Pflicht. | ADR-0023 (Negativ, Re-Evaluierungs-Trigger) | `Dockerfile.tmpl` · `test "$(echo "$status" \| awk '$2 == "VALIDSIG"` | nein | Signatur-Ablauf nicht erzwungen, nur in ADR | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-5 | INFO | Die vier Pins sind `ARG`s und per `--build-arg` ersetzbar (der Negativ-Test nutzt das bewusst). Der Pin ist damit ein Default, kein erzwungener Wert; wer die Build-Argumente kontrolliert, kontrolliert auch das Dockerfile. Das Handbuch sagt „Die Version ist festgelegt“, ohne die Überschreibbarkeit zu erwähnen. | ADR-0023 (Entscheidung 2) | `Dockerfile.tmpl` · `ARG PODMAN_SHA256_AMD64=` | nein | Pin als ARG überschreibbar, nicht dokumentiert | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-6 | INFO | `tar -xzf … -C / --strip-components=1` legt neben `usr/local` und `etc/containers` auch `README.md` des Archivs nach `/README.md` und überschreibt `etc/containers/{policy.json,seccomp.json,containers.conf,storage.conf}`; die Spezifikation nennt nur `/usr/local` und „die Konfiguration des Archivs unter `/etc/containers/`“. Kein Pfad außerhalb dieser Wurzeln, kein setuid, keine Traversal (Archiv geprüft). | `LH-FA-DEV-007.b` (Spezifikation) | `Dockerfile.tmpl` · `tar -xzf "podman-linux-${arch}.tar.gz" -C /` | nein | Archivinhalt außerhalb `/usr/local` nicht benannt | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-7 | INFO | Das Startscript schreibt `~/.config/containers/containers.conf` vollständig neu (eine vorhandene Benutzerdatei geht verloren); das vfs-Drop-in bleibt bei einem späteren Start mit `/dev/fuse` liegen, die Statusdatei meldet dann `ok`, der Treiber bleibt `vfs`. Das Verhalten stammt in dieser Form aus der Vorfassung (`storage.conf` wurde ebenso nur geschrieben, nie entfernt). | `SPEC-045` | `sandbox-init.sh.tmpl` · `> "$conf_dir/containers.conf"` | nein | Fallback-Konfiguration ohne Rückbau | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-8 | INFO | Von den drei Negativproben prüft nur „signing subkey as the pin“ den `VALIDSIG`-Vergleich (belegt durch die eigene Mutation); „unknown fingerprint“ scheitert bereits am Schlüsselabruf (HTTP 404, `curl -f`). Beide Proben fordern `.tar.gz: OK` in der Ausgabe und sind damit bei Netzausfall rot statt falsch-grün. | ADR-0023 (Negativprobe) | `sandbox_devcontainer_docker_test.go` · `"unknown fingerprint as the pin"` | ja — `make test-docker` | Negativprobe: Wirkstufe je Fall verschieden | <!-- d-check:ignore (Kennungen als Zitat im Report) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Dockerfile-Block `NestedPodman`: Reihenfolge Prüfsumme → Schlüsselabruf → Import → `--verify` → `VALIDSIG` → `tar`; `&&`-Kette, jeder Fehler bricht den `RUN` (`status="$(…)"` trägt den gpg-Exit-Code, leeres `VALIDSIG` bricht per `test`) | geprüft, ohne Befund (Mutationen `VALIDSIG` und Prüfsumme gefangen) |
| `awk '$12'` des `VALIDSIG`-Status | geprüft, ohne Befund (Feld 12 ist der Primär-Fingerprint, eigener Lauf) |
| Schlüsselserver kompromittiert / leere Antwort / HTTP-Fehler | geprüft, ohne Befund (`curl -f` bricht; ein fremder Schlüssel liefert einen anderen Primär-Fingerprint in Feld 12; mehrere Schlüssel in `key.asc` ändern das nicht) |
| Temporäres GPG-Home und Arbeitsverzeichnis | geprüft, ohne Befund (`$GNUPGHOME` liegt unter `$work`, `rm -rf "$work"` nach dem Entpacken) |
| Archivinhalt (Pfade, Links, Rechte, Eigentümer) | geprüft, ohne Befund außer F-6 |
| Architekturwahl `dpkg --print-architecture` | geprüft; Fehler statt Fallback im Code, Pin-Lücke als F-2 |
| Prüfsummen, Version, Fingerprint: eine Stelle, deckungsgleich mit ADR-0023 und Slice | geprüft, ohne Befund (Test zählt jede Pin-Zeichenkette genau einmal; Archive gegen das Netz geprüft) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `registries.conf`, `docker`-Symlink, subuid/subgid-Zeilen | geprüft, ohne Befund (subuid/subgid-Zeilen unverändert zur Vorfassung) |
| Startscript: Pfade, Idempotenz bei Wiederholung (`mkdir -p`, Überschreiben), Exit-Codes, Statusdatei `vfs` | geprüft, ohne Befund außer F-7; Exit-Codes und Stderr unverändert |
| Alter Wortlaut im Baum (`storage.conf`, `podman 5.4`, `5.4.2`, `Debian-Paket`, `passt`, `fuse-overlayfs`) | geprüft, ohne Befund (Treffer nur in ADR-0014/-0018/-0019, Changelog-Historie, Negativ-Pins der Tests, `done/`) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| Tests: `assertPodmanStaticDockerfile` nicht tautologisch (vier Mutationen gefangen); `NestedPodman`, `_Build`-Pfad, vfs-Test; Skip-Ehrlichkeit (`skipIfNestedUsernsRestricted` prüft Exit 11 vor dem Skip); Cleanup (`rmi`, `volume rm`); Zeitlimits 15 Minuten | geprüft, ohne Befund außer F-2, F-8 |
| Spezifikation `LH-FA-DEV-007.b` gegen Dockerfile/Startscript; `SPEC-019` (kein neuer Host); Lastenheft unverändert | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| Handbuch (Herkunft/Version, Kurznamen, vfs, Gemessen, Grenzen), Aussage „Podman 6 unterstützt nur cgroups v2“ | geprüft, belegt durch Release Notes v6.0.0; F-3 |
| CHANGELOG `[Unreleased]` (Changed, nichts Ungeliefertes, Messgrenzen genannt), `harness/replay.md`, ADR-Index | geprüft, ohne Befund |
| `ADR-0014`/`-0021` unverändert; ADR-0023 nur Status und Geschichte geändert | geprüft, ohne Befund (`make doc-immutable` OK) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| §3.7/§3.9 (a)(b)(c) im Diff (Kommentare der Vorlagen und Tests, Platzierung der Helfer) | geprüft, ohne Befund |
| Sprachvertrag (Skript-Ausgaben Englisch, Spec/Planung Deutsch) | geprüft, ohne Befund |
| Plan-Tabelle §3 gegen `git diff --stat` | geprüft; F-1 |
| Verschachtelter Lauf auf `amd64`; Hosts mit cgroups v1 | nicht geprüft (nur arm64-Host; Download, Prüfsumme und Signatur von `amd64` im Review per Archivlauf geprüft, kein Nested-Lauf) |
| `make gates`, `make test-docker` (Gesamtlauf), `make verify-depguard` | nicht gelaufen: Gesamtlauf Verifier-Aufgabe; keine Importregel berührt (`verify-depguard` entfällt), Einzeltests des Docker-Pfads gelaufen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Plan-Tabelle gegen Diff nicht abgeglichen · Fail-closed-Zweig ohne Pin-Test · Build-Zeit-Hosts nur in der ADR · Signatur-Ablauf nicht erzwungen, nur in ADR · Pin als ARG überschreibbar, nicht dokumentiert · Archivinhalt außerhalb `/usr/local` nicht benannt · Fallback-Konfiguration ohne Rückbau · Negativprobe: Wirkstufe je Fall verschieden

## Verdikt

**Merge-blockierend:** nein — kein HIGH und kein MEDIUM; die Lieferkettenprüfung ist fail-closed (zwei Mutationen der Prüfkette gefangen, Archiv und Prüfsummen gegen das Netz bestätigt). Die drei LOW sind Nachzieher.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen zusätzlich in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation — DoD und Gesamtläufe (`make gates`, `make test-docker`) prüft der Verifier.
