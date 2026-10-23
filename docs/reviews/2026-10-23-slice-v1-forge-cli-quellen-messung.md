# Review-Report: slice-v1-forge-cli-quellen-messung — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, Entscheidungen und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** Commit-Range `origin/main..HEAD` (`23d3f28f` bis `c2b6d99b`, 7 Commits, 91 Dateien; Slice [`slice-v1-forge-cli-quellen-messung`](../plan/planning/done/slice-v1-forge-cli-quellen-messung.md), Welle [`welle-forge-cli-gh-glab`](../plan/planning/welle-forge-cli-gh-glab.md), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) `Proposed`, Messprotokoll und Wortlaut-Entwurf unter [`docs/messungen/slice-v1-forge-cli-quellen-messung/`](../messungen/slice-v1-forge-cli-quellen-messung/README.md), Skripte unter `tools/messung/forge-cli/`).

**Eingang gelesen:** Slice-Plan (alle Sektionen), Messprotokoll (vollständig), [ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (vollständig), Wortlaut-Entwurf (vollständig), `tools/messung/forge-cli/lib.sh` und `05-bau.sh`, `08-aufraeumen.sh`, beide `block-*tgz.dockerfile`, Rohdaten stichprobenartig (`rohdaten-03/05/06/07`), [`AGENTS.md`](../../AGENTS.md) §3.1, §3.4 bis §3.10, §5 Regel 1 und 6, [`harness/conventions.md`](../../harness/conventions.md) (Ablage, Modus), MR-008, im Code `Dockerfile.tmpl`, `devcontainerTemplateData`, `egressBaseHosts`, `agentCredentialEnv`.

## Ergebnis

Keine HIGH-Findings. Ein MEDIUM, fünf LOW, drei INFO. Die Rohdaten sind frei von Zugangsdaten (eigene Suche, siehe Schwerpunkt 2); [ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) ist durchgängig `Proposed` und die Freigabe-Punkte a bis g sind als offene Wahl des Projektinhabers dargestellt. Das MEDIUM betrifft die Wiederverwendung der Skripte für die Nachmessung mit echtem Token (Punkt f der ADR), nicht die vorliegenden Rohdaten. `make gates` auf `c2b6d99b` (vor diesem Report): Exit 0.

## Findings

### F-1 — MEDIUM — `redact()` schwärzt nur einen Ausschnitt der Formen, in denen ein Token in einer Ausgabe steht

- `kategorie`: MEDIUM · `quelle`: [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Slice §6 Risiko 3 ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md), 2×) · `klasse`: rohdaten-schwaerzung-lueckenhaft
- `pfad`: `tools/messung/forge-cli/lib.sh` · „Schwaerzt vorsignierte URLs, Autorisierungs-Header, bekannte Token-Praefixe und den Home-Pfad“ (Funktion `redact`)
- `befund`: Ich habe die acht `sed`-Ausdrücke von `redact` gegen Beispielzeilen laufen lassen. Geschwärzt werden `Authorization: …` (Großschreibung) und `Bearer`/`Basic`-Werte sowie Präfixe (`ghp_`, `glpat-` …). Nicht geschwärzt bleiben: `authorization: token abc123` (Kleinschreibung, wie HTTP/2-Ausgaben sie führen), `PRIVATE-TOKEN: abc123def`, `GITLAB_TOKEN=abc123def`, `{"token":"abc123def"}`, `https://user:abc123@gitlab.com/x`, `x-access-token:abc123@github.com`. Die letzte Berichtigung (`c2b6d99b`, Wert hinter `Bearer`) schloss genau eine dieser Lücken.
- `risiko`: In den vorliegenden Rohdaten entsteht kein Schaden (kein Token verwendet, Suche leer). [ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkt f schlägt aber eine Nachmessung mit echtem, kurzlebigem Token vor, und die Skripte (`inner-07.sh`, `07-laufzeit.sh`) sind dafür die Grundlage. Mit echtem Token landet ein Wert in einer der ungeschwärzten Formen (etwa `GITLAB_TOKEN=…` in einer Fehlermeldung, ein `https://user:token@`-Remote, ein Kleinschreib-Header einer `-v`-Ausgabe) in einer Rohdatei und im Verlauf: dritter Beleg für die Beobachtung, und der Abschnitt „Token-Sicherheit“ des Protokolls nennt `redact` als Schutz, ohne die Grenze zu nennen. Der nächste Änderer, der `redact` für den Lauf mit Token wiederverwendet, folgt dieser Zusage.
- `verifizierbar`: Review-only (Beispielzeilen durch dieselben Ausdrücke; ein Test-Skript existiert nicht).

### F-2 — LOW — Die Spalte „Läufe“ folgt ihrer eigenen Definition nicht

- `kategorie`: LOW · `quelle`: Slice §2 DoD 1 („je Aussage … Zahl der Läufe“), [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) · `klasse`: laufzahl-uneinheitlich
- `pfad`: [`docs/messungen/slice-v1-forge-cli-quellen-messung/README.md`](../messungen/slice-v1-forge-cli-quellen-messung/README.md) · „*Läufe*: Zahl der Läufe, deren Rohdaten im Verzeichnis liegen; **gewertet ist der letzte Lauf**“
- `befund`: Die Definition zählt Läufe mit vorliegenden Rohdaten, also 1 je Rohdatei. Die Tabellen führen für A1 bis A9 „2“ (eine Rohdatei `rohdaten-02.txt`, der Text nennt einen früheren Lauf), für B1 „3 / 3“, für A5 „2 / 1“, aber für A10 bis A14 „1“, obwohl für Schritt 03 vier frühere Läufe genannt sind. Beim Bau `ghk` (A5, „1“) nennt der Text drei frühere Läufe, B2 führt für dieselben Bauten „3“. Mal zählt die Spalte alle Läufe, mal nur den gewerteten.
- `risiko`: Wer aus dem Protokoll die Belastbarkeit einer Aussage abliest (der CR-Autor über die Liste „nicht gemessen“), liest nicht dieselbe Größe von Zeile zu Zeile. Die Zahlen D „1 gewertet“ sind konsistent.
- `verifizierbar`: Review-only.

### F-3 — LOW — Einschränkungen des Protokolls fallen in ADR, Skript-Kommentar und Alternativentabelle weg

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 („auch eine Aussage über Verhalten nennt nur, was gemessen wurde“) · `klasse`: zusage-weiter-als-messung
- `pfad`: `docs/plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md` · „Kein Feature trägt Signatur oder Attestierung am OCI-Artefakt (A3)“; ebenso „für `glab` zusätzlich, dass das Projekt keine Signatur veröffentlicht“ (Entscheidung 6); `tools/messung/forge-cli/block-glabtgz.dockerfile` · „Das Projekt veroeffentlicht keine Signatur und keine Attestierung fuer die Archive“; ADR-Tabelle Messung · „Alle drei bauen (C4 bis C6)“
- `befund`: A3 lautet „in diesen beiden Mechanismen (Referrers, Tags); andere nicht geprüft“, A12 „Suchmuster, nicht abschließend“, C4 bis C6 sind Handvarianten (`k`, Schlüssel ohne Schema), während derselbe ADR-Abschnitt sagt, der Weg über `source:` baue nicht. Die ADR und der Kommentar im Block geben die unbedingte Fassung wieder. Die Wortlaut-Tabelle (S3) führt den Vorbehalt dagegen richtig („ob `glab` außerhalb des Suchmusters signiert“).
- `risiko`: Die Entscheidungsgrundlage, die der Projektinhaber freigibt, ist an diesen Stellen stärker als die Messung; der Kommentar im Block wandert in den Code-Slice (Muster der Handvarianten) und wird dort als Zusage über `glab` gelesen.
- `verifizierbar`: Review-only.

### F-4 — LOW — Der Wortlaut-Entwurf trägt Messwerte im Lastenheft-Text, den seine eigene Konvention ausschließt

- `kategorie`: LOW · `quelle`: Maintainability (eigene Konvention des Entwurfs, Muster [`slice-cr-devcontainer-ollama`](../plan/planning/done/slice-cr-devcontainer-ollama.md)), [`AGENTS.md`](../../AGENTS.md) §3.10 (Bedingung ohne Test) · `klasse`: wortlaut-konvention-selbstbruch
- `pfad`: [`docs/messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md`](../messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md) · „Gemessen ist das auf Docker 29.9.2, Linux `amd64`, `gh` 2.103.0 und `glab` 1.124.0;“ (Satz S6)
- `befund`: Die Konvention oben im Entwurf sagt, Aussagen über das Verhalten des Generators seien Zielbild und „tragen keinen Messwert“. S6 setzt Docker-Version und Werkzeugversionen in den Lastenheft-Absatz; S4 nennt Host-Namen, die die Stelle 3 selbst der Spezifikation zuordnet (der Entwurf bietet das Streichen an). Die Sätze S1, S3 und S5 („kein Token gespeichert“) sind Bedingungen über Verhalten des Codes; die Fitness-Function der ADR nennt Tests dafür, der Entwurf verweist nicht darauf.
- `risiko`: Der CR-Autor übernimmt den Wortlaut; ein Messstand im Lastenheft altert mit dem nächsten Release der Werkzeuge. Kein Verweis nach unten (Kennungsfreiheit eingehalten: der Wortlaut nennt nur `LH-*`).
- `verifizierbar`: Review-only.

### F-5 — LOW — Skript-Kommentar verweist auf einen Abschnitt, den das Protokoll nicht hat

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 (Rang-Zeiger muss auflösen), §3.9 (b) · `klasse`: kommentar-zeiger-ins-leere
- `pfad`: `tools/messung/forge-cli/lib.sh` · „Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Abschnitt Wiederholung).“
- `befund`: Das Protokoll hat keinen Abschnitt „Wiederholung“; die Zeilen `AUSSAGE …` erklärt der Abschnitt *Lesehilfen*.
- `risiko`: Gering; der nächste Änderer sucht einen Abschnitt, den es nicht gibt. Sonst sind die Kommentare der Skripte im Indikativ und tragen Zusage, Kopplung oder Grenze (geprüft, siehe Negativbefund).
- `verifizierbar`: Review-only.

### F-6 — LOW — Prozessabweichungen des Implementers: Heredoc-Commits und `python3 --version` auf dem Host

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only), Kommentar in `.claude/hooks/pretooluse-command-guard.sh` („Abhilfe: Commit-Messages via `git commit -F <datei>`“) · `klasse`: host-werkzeug-und-heredoc
- `pfad`: zwei Commits der Range (nach Meldung des Implementers; welche, ist aus dem Diff nicht ablesbar) und der gemeldete Aufruf `python3 --version` (kein Artefakt im Diff)
- `befund`: Der Implementer meldet zwei Commits mit Heredoc statt Write-Datei und einen Host-Aufruf von `python3 --version`. Die Messages selbst sind inhaltlich in Ordnung (Slice-Kennung, ADR, `LH-*`; Regel 1 erfüllt; kein Heredoc-Rest im Text). Der Guard nennt `-F <datei>` ausdrücklich als Weg, weil er quote-blind ist. Der Python-Aufruf ist ein Versionsprobe ohne Nachweis-Funktion: kein Skript im Diff verwendet `python` (Suche leer), der Guard prüft Interpreter bewusst nicht, aber §3.1 zielt auf Host-Toolchains nur als Nachweis.
- `risiko`: Gering. Heredoc-Messages bergen das Guard-False-Positive und Quoting-Fehler; der Host-Aufruf ist harmlos, solange kein Ergebnis in einen Handoff als Beleg eingeht (hier: nein). Beides ist Abweichung, keine Entlastung, und gehört in den Handoff als Abweichung.
- `verifizierbar`: Review-only (`git log --format=%B`).

### F-7 — INFO — `glab`-Konfiguration führt Schlüssel für automatisches Nachladen weiterer Binaries

- `kategorie`: INFO (Hinweis an Verifier/Validator) · `quelle`: [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) · `klasse`: ungeprueft-beobachteter-schluessel
- `pfad`: [`docs/messungen/slice-v1-forge-cli-quellen-messung/rohdaten-07-glabsk-discover.txt`](../messungen/slice-v1-forge-cli-quellen-messung/rohdaten-07-glabsk-discover.txt) · „duo_cli_auto_run duo_cli_auto_download duo_cli_binary_path“ (Schlüsselliste der `config.yml`, auch `orbit_cli_auto_download`, `notify_skill_updates`)
- `befund`: Die Konfiguration von `glab` 1.124.0 führt Schlüssel, die nach automatischem Download weiterer Programme aussehen. Das Protokoll (E4, D5) nennt `git_protocol`, `telemetry`, `check_update`, aber diese nicht; D5 bleibt für die gemessenen Aufrufe richtig („nicht gemessen“ in interaktiver Sitzung).
- `risiko`: Wenn ein Aufruf in einer Sitzung ein weiteres Binary von einem weiteren Host lädt, wäre die Host-Liste (`gitlab.com`) und die Zusage nach [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) zu knapp. Aus den Rohdaten nicht zu entscheiden.
- `verifizierbar`: Review-only; Nachmessung (Punkt f) oder Nennung unter „nicht gemessen“.

### F-8 — INFO — DoD-Haken „Umgang mit Tokens gemessen“ bei offenem Erfolgspfad

- `kategorie`: INFO (Zuständigkeit Verifier) · `quelle`: Slice §2 DoD 2 · `klasse`: dod-haken-vs-nicht-gemessen
- `pfad`: Slice §2 · „**Umgang mit Tokens gemessen**“ (angehakt) gegen Protokoll F „Jeder Aufruf mit einem gültigen Token“
- `befund`: Die DoD fragt, „ob ein Aufruf mit gesetzter Variable ohne `auth login` gelingt“; gemessen ist, dass er versucht wird und `401` liefert (Platzhalter), nicht dass er gelingt. Das Protokoll sagt das offen (E2, F), der Auftrag verbot ein echtes Token. Ob der Haken so gilt, urteilt der Verifier; für den Reviewer ist der Haken ehrlich begründet, aber die DoD-Formulierung trägt den Fall nicht ausdrücklich.
- `risiko`: gering; Zuordnung zur Verifier-Frage.
- `verifizierbar`: Review-only.

### F-9 — INFO — „15 bis 21 s“ in D3 ist gröber als die Rohdaten

- `kategorie`: INFO · `quelle`: Maintainability · `klasse`: zahl-ungenau
- `pfad`: Protokoll D3 · „Gesamtdauer des Containerlaufs 127 bis 131 s statt 15 bis 21 s“
- `befund`: `gesamtdauer` im Modus `discover` liegt für die `glab`-Läufe bei 18,2 bis 19,8 s (`rohdaten-07-glab*-discover.txt`, Zeile 233), im Modus `restricted` bei 127,3 bis 130,5 s. Die Spanne „15 bis 21“ umfasst beides, ist aber nicht die gemessene.
- `risiko`: keines für die Aussage (Faktor etwa 7 bleibt).
- `verifizierbar`: Review-only.

## Prüfung nach Schwerpunkten

1. **Aussagen gegen Messung (Stichproben, gegengezeichnet).** B1: sechs Bauten (`gh`, `ghx`, `glabs`, `glabp`, `glabz`, `ghn`) Exit 1 nach 0,8 bis 0,9 s mit der genannten Meldung (`rohdaten-05.txt`); die Handvarianten `*k` Exit 0. C2: 1 479 255 931 minus 1 419 218 247 = 60 037 684 ✓; C7: 1 472 224 607 minus 1 412 193 634 = 60 030 973 ✓. C9/C10/C11 (`rohdaten-06.txt`: 19,3/19,0/16,0/17,2 s; Differenzen 48 809 142, 40 798 568, 56 100 573, 67 439 177; Negativprobe Exit 1) ✓. A10/A12: `immutable: true`, `checksums.txt` gleich Datei, `arm64`-Prüfsummen in ADR und Blöcken gleich denen in `rohdaten-03.txt` (gh `7f322df3…`), Signierschritte nur macOS/Windows im Auszug ✓. D3 `127 bis 131 s` ✓ (F-9 für die Gegenspanne), E8 (`gh auth status` 2.46.0 Exit 0, 2.103.0 Exit 1) ✓, E4 (`device-id` 36 Byte; `glab` `config.yml`, `aliases.yml`, `config.lock`; `platzhalter-treffer=0`) ✓. Aussagen über Code in ADR und Plan (`Dockerfile.tmpl` `USER root` unter `or .UID .Sandbox`, `devcontainerTemplateData` Rückkehr bei `!sandbox` vor `AgentEnv`, `egressBaseHosts` sieben Namen, Credential-Helper, `agentCredentialEnv`) an den Stellen gelesen: stimmen. Abweichungen: F-2, F-3, F-9.
2. **Zugangsdaten.** Eigene Wiederholung der Suche über [`docs/messungen/slice-v1-forge-cli-quellen-messung/`](../messungen/slice-v1-forge-cli-quellen-messung/README.md), `tools/messung/forge-cli/` und `docs/plan/`: Token-Präfixe (`ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`, `github_pat_`, `glpat-`, `glptt-`, `gldt-`, `glrt-`, `sk-ant-`, je mindestens 8 Zeichen) nur ein Treffer in einem fremden, unveränderten Slice-Text (Erkennungsausdruck), keiner in Rohdaten oder Skripten; `Bearer`, `Authorization`, `X-Amz`, `Signature=`, `sig=`, `PRIVATE-TOKEN` in den Rohdaten: 0 Treffer; Home-Pfad des Benutzers in Rohdaten und Skripten: 0 (die Treffer auf „db“ sind `man-db`); Schlüssel `oauth_token`/`token:` nur als Namen von Aufruf-Etiketten; der Platzhalter steht nur im Protokoll und in den Skripten. Die Rohdaten sind sauber. Schwäche der Schutzfunktion selbst: F-1. Der Platzhalter ist als solcher ausgewiesen; die Zusage „kein Token angelegt oder verwendet“ steht in den Rohdaten nirgends im Widerspruch.
3. **Messvariante gegen Auslieferung.** Das Protokoll trennt die Handvarianten (`k`, `*deb`, `*tgz`, `basd`) von den erzeugten Projekten in einer eigenen Tabelle und sagt, dass der Katalogeintrag und der Dockerfile-Block nicht von `u-boot` erzeugt sind; die ADR führt „Bau eines vom Generator erzeugten Dockerfiles mit einem Block“ unter „nicht gemessen“ und legt ihn in die Fitness-Function (`make test-docker`). Das Risiko 2 des Slice ist damit sichtbar, nicht verdeckt. Einschränkung nur F-3 (Wortlaut „bauen“ ohne „Handvariante“).
4. **Kommentare in Skripten (§3.7/§3.9).** Alle `#`-Kommentare von `lib.sh`, `00` bis `08`, `inner-*`, Dockerfile-Blöcke gelesen/durchsucht auf Konjunktiv über Verworfenes, „früher/former/no longer/ohne X wäre“, Namen geänderter Symbole: keiner gefunden; sie tragen Zusage, Kopplung, Abgrenzung (z. B. „HAND-Variante … nicht von u-boot erzeugt“) oder Grenze. Die Berichtigungsnotiz zu `redact` steht im Protokoll (dort als Zustand: „wurde berichtigt“), nicht im Skript. Befund F-5. Das Muster (c) (Helfer zwischen Doc-Kommentar und Funktion) kommt nicht vor.
5. **[ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Status und Freigabe-Punkte.** Status `Proposed` in Kopf, Index-Zeile und Geschichte; Entscheidung ist ausdrücklich „Vorschlag (gilt erst mit `Accepted`)“; a bis g als Tabelle mit „Empfehlung des Implementers (keine Entscheidung)“; Quellenwahl, die vom Muster abweicht, liegt beim Projektinhaber (Slice §2 DoD 3 bleibt offen, richtig). MADR-Form (MR-008): Kontext, Entscheidung, Verglichene Alternativen mit Gegenoption je Entscheidung, Konsequenzen, Fitness Function, Re-Evaluierungs-Trigger, Geschichte; `Schärft:` zeigt nach oben (`LH-*`, `SPEC-*`); keine Verweise auf Slice, Welle, Roadmap oder Carveout (Wörter „Slice“ nur generisch, keine Links); Ablage-Link auf das Protokoll ist ein Doku-Zeiger wie in [ADR-0024](../plan/adr/0024-ollama-im-sandbox-devcontainer.md). ADR-Index ergänzt. Keine Immutabilitätsverletzung (neue ADR, `Proposed`; [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md)/0024/0032 unverändert). Offene Punkte der Darstellung: F-3.
6. **Wortlaut-Entwurf.** Kein Verweis nach unten im vorgeschlagenen Lastenheft-Text (nur `LH-*`); Bedingungen und Nicht-Zugesagtes je Satz mit Spalte „Nicht gemessen“; F-4. Der Entwurf ändert das Lastenheft nicht (`spec/` ist im Diff nicht berührt, geprüft).
7. **Plan-Tabelle gegen `git diff --stat`.** Die 91 Dateien des Diffs sind: Skripte und Hilfsdateien unter `tools/messung/forge-cli/` (27 Dateien: `00` bis `08`, `lib.sh`, `inner-02`, `inner-02-index`, `inner-02-sonde`, `inner-03`, `inner-03-debian`, `inner-07`, `Dockerfile.cli`, vier Blöcke), `docs/messungen/…/` (Protokoll, Wortlaut-Entwurf, Rohdaten), ADR und ADR-Index, `docs/messungen/README.md`, Slice, Welle-Datei, `open/README.md`, Roadmap. Jede Gruppe hat einen Eintrag in §3, jeder Eintrag einen Diff; `internal/`, `cmd/`, `spec/`, `docs/user/`, README, CHANGELOG sind nicht berührt (`git diff --stat` ohne Treffer), wie die Zeile „unverändert“ sagt, und die Plan-Zeile belegt das mit gelesenen Stellen. Roadmap-Zustandszeile nennt den Zustand mit Zeiger, keine Chronik. Ablage: nur Skripte unter `tools/messung/`, nur `.md`/`.txt`/`.log` unter `docs/messungen/` (`make ablage-pruefung` 26 Fälle, 0 Fehlschläge). Commit-Messages tragen Slice-Kennung und `LH-*`/[`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (§5 Regel 1); keine `SPEC-*`/`ARC-*`. Nutzer-Doku unberührt (§5 Regel 6 nicht berührt).
8. **Sensoren.** `make gates`: Exit 0 (`lint`, `test`, `coverage-gate`, `docs-check`, `baseline-verify`, `ablage-pruefung`), von mir auf `c2b6d99b` gelaufen. Nicht gelaufen: `make test-docker` (kein Docker-Pfad des Produkts berührt, nur Skripte unter `tools/messung/`, die nicht Teil von `make`-Gates sind), `make verify-depguard` (keine Importregel berührt), `make ci`/`fullbuild` (nicht Teil der Review-Frage). Die Messskripte selbst habe ich nicht erneut ausgeführt (Netz, Zeit, Zwischen-Images); die Zusagen beruhen auf den Rohdaten, die ich nachgelesen habe.

## Negativbefund (geprüft, ohne Befund)

- Rohdaten auf Zugangsdaten, Home-Pfad, vorsignierte URLs (Schwerpunkt 2) · Skript-Kommentare auf §3.7/§3.9 (bis auf F-5) · [ADR-0035](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md): Status, Index, MADR-Form, `Schärft:`-Richtung, Freigabe-Punkte · Referenzmodell (kein Verweis von Rang 1 bis 4 nach unten) · Ablage-Regel · Plan-Tabelle gegen Diff · Roadmap-Zeiger, Welle-Datei-Zeiger, `open/README.md` · Lifecycle (`333355e3` reiner Move, `2e986785` Verweis-Nachzug getrennt, [`AGENTS.md`](../../AGENTS.md) §3.3) · Sprachvertrag (kein CLI-Output im Diff) · Lastenheft und Spezifikation unverändert.

## Nicht geprüft

- Wiederholung der Messung (Netz, Zeit); die Rohdaten der früheren Läufe sind überschrieben (Protokoll sagt es offen) und nicht prüfbar.
- Inhalt der `arm64`-Pins über die Prüfsummen-Datei hinaus (nicht gebaut, vom Protokoll als solches geführt).
- Die Richtigkeit von Aussagen über Fremdquellen (Installationsskripte der Features, Upstream-Workflows) über das Nachlesen der Auszüge in `rohdaten-02/03.txt` hinaus.
- DoD-Erfüllung und Slice-Closure: Verifier-Aufgabe.
