# Review-Report: slice-v1-handbuch-sshd-anleitung — 2026-10-18

**Review-Art:** Code (Dokumentations-Diff) — geprüft gegen Plan, Hard Rules und Messquelle.

**Gegenstand:** Commit-Range `0bfc46f..HEAD` (`cda63fa`, `7dc8521`), Slice `slice-v1-handbuch-sshd-anleitung`. <!-- d-check:ignore (Slice-Kennung im Titel; Report friert ein) -->

**Skill:** `.harness/skills/reviewer.md` @ `7dc8521` · <!-- d-check:ignore (Skill-Pfad, existiert im Repo) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-handbuch-sshd-anleitung` (in `in-progress/`), Messquelle `slice-v1-devcontainer-feature-sshd` (`done/`) <!-- d-check:ignore (Slice-Kennungen; Report friert ein) -->
- `AGENTS.md` §3.7, §3.9, §3.10 (beide Schärfungen), §5 Zeile 6
- Code-Fundstellen: `generate.go` (`collectDevcontainerForwardPorts`), `devcontainer_sandbox.go` (`sandboxRelaxationWarnings`, `loginServiceActive`), Vorlage `devcontainer.json.tmpl`
- `harness/review.md`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Schritt 3 rät als Weg für den Port „eine Port-Option der Container-Engine“ ohne den Hinweis, die Veröffentlichung an die Loopback-Adresse des Hosts zu binden. Beim Passwort-Login und erlaubtem Root-Login auf `0.0.0.0` im Container kann eine unqualifizierte Veröffentlichung den Login-Dienst auf allen Host-Schnittstellen öffnen; der Beispielbefehl nennt nur `localhost` auf der Client-Seite. | Maintainability (Sicherheitsaussage, Slice §6 „Sicherheitsaussage verwässert“) | `docs/user/benutzerhandbuch.md` · „oder eine Port-Option der Container-Engine (nicht gemessen, wie ein Editor oder die CLI das tut)“ | nein — Review-only | Anleitung-Rat ohne Sicherheitsgrenze |
| F-2 | MEDIUM | Der Slice schreibt vor: eine Fremdwerkzeug-Aussage ohne gelesene Quelle „steht nicht im Text“. Der Text führt `ssh-keygen -t ed25519`, die Rechte `700`/`600`, `docker exec`, `ssh -i … -p 2222 vscode@localhost` und die Weiterleitung durch Editor/Engine; die Plan-Tabelle räumt „Allgemeinwissen, nicht an einer Dokumentation gelesen“ ein. Die Kennzeichnung „nicht gemessen“ ist ehrlich und im Text sichtbar, ersetzt aber die verlangte Quelle nicht. | Slice-Plan §1 „Was die Anleitung nur sagen darf“ (b); `AGENTS.md` §3.10 | `docs/user/benutzerhandbuch.md` · „Die Rechte `700` und `600` sind eine Annahme über die Prüfung des“ | nein — Review-only | Fremdwerkzeug-Aussage ohne Quelle |
| F-3 | LOW | Schritt 2 stellt Schlüssel- und Passwort-Login gleichrangig („Wählen Sie einen der beiden Wege“); eine Präferenz oder ein Hinweis auf die Passwortstärke fehlt, obwohl Root-Login erlaubt und der Port auf allen Schnittstellen offen ist. Der Sicherheitshinweis steht am Ende; der Querverweis davor ist vorhanden. | Slice-Plan §6 (Risiko „Sicherheitsaussage verwässert“) | `docs/user/benutzerhandbuch.md` · „Wählen Sie einen der beiden Wege.“ | nein — Review-only | Sicherheitsaussage gleichrangig neben bequemerem Weg |
| F-4 | LOW | Schritt 1 gibt „Neubau statt Neustart“ als Anweisung („Bauen Sie den Container deshalb … neu auf“); gemessen oder gelesen ist nur „Installationsskript läuft beim Image-Build“, die Folgerung ist nicht belegt (Slice-Plan: „schreibt das Wort, das belegt ist“). Die Egress-Beschreibung im selben Abschnitt sagt für eine Änderung: „starten den Container neu“. | Slice-Plan §1 offener Punkt „neu starten oder neu bauen“ | `docs/user/benutzerhandbuch.md` · „statt ihn nur neu zu“ | nein — Review-only | Folgerung als Anweisung, Begriff uneinheitlich im Abschnitt |
| F-5 | INFO | Die Klammer „`ssh-keygen` liegt im Basis-Image …, gemessen mit `command -v ssh-keygen`“ steht in Schritt 1 des Host-Weges; die Plan-Tabelle führt die Messung als „vom Implementer nicht wiederholt“. Der Reviewer hat sie wiederholt: `/usr/bin/ssh-keygen` im Image, Aussage wahr. Für das Erzeugen auf dem Host ist sie ohne Bezug. | Maintainability | `docs/user/benutzerhandbuch.md` · „gemessen mit `command -v ssh-keygen`“ | ja — wiederholt (siehe Negativbefunde) | Messaussage nicht vom Autor wiederholt |
| F-6 | INFO | `docs/plan/planning/open/README.md` führt den Slice weiter mit Zeiger auf `in-progress/`; die Plan-Tabelle führt die Datei als „unverändert“ mit der Begründung „Bestandszeile entstand beim Anlegen“, nicht mit einer gelesenen Aussage zum jetzigen Inhalt (zweite Schärfung von §3.10). Kein Gate-Befund. | `AGENTS.md` §3.10 (zweite Schärfung) | `docs/plan/planning/in-progress/slice-v1-handbuch-sshd-anleitung.md` · „Bestandszeile entstand beim Anlegen des Slice“ | ja — `make docs-check` (grün) | „unverändert“-Zeile ohne gelesene Aussage | <!-- d-check:ignore (Slice-Kennung; Report friert ein) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Port 2222, `sshd -T` (Passwort-/Schlüssel-Login, `PermitRootLogin yes`), gesperrte Konten `root`/`vscode`, `~/.ssh` fehlt | geprüft gegen Messtabelle `slice-v1-devcontainer-feature-sshd` §3, ohne Befund; Messung (CLI 0.89.0, Basis-Image) im Text genannt | <!-- d-check:ignore (Slice-Kennung; Report friert ein) -->
| `remoteUser vscode`, Volume statt Bind-Mount, `workspaceMount: ""`, `forwardPorts` nur bei nicht leerer Liste, `runArgs` ohne Port-Veröffentlichung | geprüft an Vorlage `devcontainer.json.tmpl` und `collectDevcontainerForwardPorts` (`generate.go`), ohne Befund |
| Warntext und Sichtbarkeit unter `--quiet`, nur im Sandbox-Profil | geprüft an `sandboxRelaxationWarnings` (`relax(…)`, `sandbox && loginServiceActive`) und Handbuch-Zeile `--quiet` Ausnahme (1), ohne Befund |
| Ehrlichkeit der Kennzeichnung „nicht gemessen“ | geprüft: SSH-Login, `ssh-keygen`/`authorized_keys`/Rechte, Weiterleitung, Neubau sind im Fließtext gekennzeichnet; der Absatz vor Schritt 1 nennt, dass kein SSH-Login gemessen ist; Befund nur F-2/F-4 (Quelle, nicht Kennzeichnung) |
| Rechte-/Befehls-Korrektheit (praktisch) | Wegwerf-Container `mcr.microsoft.com/devcontainers/base:debian` mit Debian-`openssh-server` (Port 2222, Root-Login `yes`, Passwort-Login `yes`, `UsePAM yes`), `vscode` und `root` per Sperre `L`, `ssh-keygen -t ed25519`, `~/.ssh` 700, `authorized_keys` 600: `ssh -i … -p 2222 vscode@localhost` meldete Login ohne Passwort; `sudo` und `ssh-keygen` im Image vorhanden. Der Aufbau ist eine Näherung (Debian-Paket statt Feature-Installationsskript), ohne Portweiterleitung des Editors/der Engine und ohne Dev-Container-CLI |
| `AGENTS.md` §5 Zeile 6 | geprüft: keine `LH-`/`ADR-`/`SPEC-`/Slice-Kennung im neuen Abschnitt (Suche), keine Verweise auf Planungs-/Entwicklerdokumente, keine Chronik; „(gemessen mit …)“ nennt Stand |
| Querverweise | `#login-dienst-sshd-optional` (Überschrift unverändert, Feature-Tabelle Zeile 797) und `#sicherheitshinweis-sshd` (expliziter Anker) lösen auf |
| Andere `sshd`-Stellen (Feature-Tabelle, Egress-Tabelle „keine zusätzlichen Hosts“, `--quiet`-Zeile) | geprüft, ohne Widerspruch zur neuen Fassung |
| `README.md`, `README.de.md`, `CHANGELOG.md` | geprüft: `sshd` nur als Katalog-Schlüssel und in der Releases-Tabelle, ohne Nutzungsaussage, ohne Widerspruch |
| Sprachvertrag und Stil | Handbuch Deutsch, Befehle und Dateinamen englisch/literal; Stil des Bestands (Fettdruck-Schrittüberschriften, Bindestrich-Listen) gewahrt, ohne Befund |
| Plan-Tabelle gegen `git diff --stat` | drei geänderte Dateien (`benutzerhandbuch.md`, `roadmap.md`, Slice) stehen alle in der Tabelle; Fundstellen-Tabelle ist breit, aber ein getrennter Zusatz neben der Datei-Tabelle und deckt jede Aussage mit Fundstelle; `CHANGELOG.md`, README-Dateien mit gelesener Aussage belegt (zweite Schärfung, ohne Befund außer F-6) |
| Roadmap-Marker | „Nichts in Arbeit.“ entfernt; Abschnitt „Offene Wellen“ trägt den Derivativ-Satz, ohne Befund |
| Nicht geprüft | echter SSH-Login mit dem Feature-Installationsskript und der Dev-Container-CLI (kein Node im Repo-Werkzeug; Slice schließt die Messung aus); `make test-docker`, `make verify-depguard` (kein Docker-/Importpfad berührt, kein Code im Diff); Portweiterleitung eines Editors oder der Engine |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Anleitung-Rat ohne Sicherheitsgrenze · Fremdwerkzeug-Aussage ohne Quelle · Sicherheitsaussage gleichrangig neben bequemerem Weg · Folgerung als Anweisung, Begriff uneinheitlich im Abschnitt · Messaussage nicht vom Autor wiederholt · „unverändert“-Zeile ohne gelesene Aussage

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) sind offen; F-1 betrifft die Sicherheitsaussage der Anleitung, F-2 eine ausdrückliche Vorgabe des Plans. Keine HIGH-Anker getroffen.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation — DoD-Konformität prüft der Verifier.
