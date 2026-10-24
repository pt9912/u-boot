# Review-Report: slice-v1-ollama-cloud-messung

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff e60956c3..HEAD` (HEAD `23639d60`), Slice
  [`slice-v1-ollama-cloud-messung`](../plan/planning/in-progress/slice-v1-ollama-cloud-messung.md)
  der Welle `welle-ollama-cloud-modelle`: Skripte `tools/messung/ollama-cloud/`, Messprotokoll mit 22 Rohdaten-Dateien unter `docs/messungen/slice-v1-ollama-cloud-messung/`
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), `AGENTS.md` §3.7, §3.9, §3.10
* Nicht Gegenstand: DoD-Prüfung (Verifier)

## Urteil

Ein HIGH (abgebrochener Kommentarsatz in `Dockerfile.cli`), ein MEDIUM (DoD-Wortlaut vom Implementer umformuliert), vier LOW, vier INFO.
Zugangsdaten: kein Fund. Messaussagen: bis auf eine falsche Zählung (LOW-3) gegen die Rohdaten nachgezählt, ohne Abweichung.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit-Code 0 (`EXIT=0` im Log; `ablage-pruefung`: 26 Fälle, 0 Fehlschläge) |
| `redact-probe.sh` | Exit 0: „24 Formen geschwaerzt, neutrale Zeilen unveraendert“ |
| `redact-probe.sh`, rote Gegenprobe | nicht selbst gefahren (Skript müsste verändert werden, kein sed/Host-Edit erlaubt); stattdessen unabhängige Suche (unten) |
| `06-suche.sh` | auf Messordner plus `tools/messung/ollama-cloud/`: Exit 1, aber nur Treffer in `README.md` des Protokolls (siehe LOW-2); auf ganz `docs/messungen/` und `tools/messung/`: nur ältere Protokolle (siehe MEDIUM-1) |
| eigene Suche | Präfixe `ghp_ gho_ github_pat glpat sk-ant sk-… ollama_…`, `Bearer`, `Authorization`, `BEGIN … KEY`, `ssh-ed25519`, `ssh-rsa`, `AAAA…`, `connect?`, `signin?`, `key=`, `nonce`, `eyJ…`, Wörter ≥ 40 Zeichen über Messordner und Skripte: ohne Treffer ausser Namen, `[Query entfernt]`-Marken und Docker-Kennungen |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Produktcode, keine Importregel berührt; Auftrag (Speicher) |
| Aufräumen | `docker ps -a`, `docker images -a`, `docker volume ls`: 0 Objekte mit Präfix `u-boot-messung-ocloud`/`vsc-…`, 0 Container |

## Findings

### HIGH-1

* `kategorie`: HIGH
* `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen), Skill-Anker „Kommentar bricht mitten im Satz ab“
* `pfad`: `tools/messung/ollama-cloud/Dockerfile.cli` · „(npm loest sie beim Bau auf) und die Debian-Pakete des Node-Images.“
* `klasse`: Kommentar-Rest nach Teilersetzung
* `befund`: Der Kommentarblock endet mit einem Satzrest, der keinen Satz fortsetzt; die Pin-Aussage steht davor schon vollständig („Nicht gepinnt: die npm-Abhaengigkeiten … und die Debian-Pakete“).
* `risiko`: Wer das Hilfs-Image ändert (Pin heben), liest zwei Aussagen zur Pin-Grenze, von denen eine unvollständig ist, und kann den Rest für eine zusätzliche Bedingung halten; das Muster ist die ausdrücklich benannte Klasse, kein Gate fängt sie.
* `verifizierbar`: Review-only

### MEDIUM-1

* `kategorie`: MEDIUM
* `quelle`: `AGENTS.md` §3.6 (Gates/Akzeptanz nicht still ändern), `harness/roles.md` (Rollenvertrag: Plan ändert der Planner/Architect)
* `pfad`: `docs/plan/planning/in-progress/slice-v1-ollama-cloud-messung.md` · „**Keine Zugangsdaten:**“ (DoD-Punkt 3)
* `klasse`: DoD-Wortlaut vom Implementer angepasst
* `befund`: Die Zeile verlangte „die Suche … über `docs/messungen/` und `tools/messung/` ist leer“ und wurde im selben Diff, der sie abhakt, auf „leer über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud/`“ plus „älteres von Hand als keine Zugangsdaten gelesen“ umformuliert. Die Ausweitung trägt: ich habe die älteren Treffer selbst gelesen (öffentliche `ssh-ed25519`-Zeilen und „Your new public key is:“ in `slice-v1-devcontainer-ollama-messung`, GPG-Schlüssel-IDs, `DUMMY`/`platzhalter-…`, Testnamen `…_Ollama_…`): kein privater Teil, kein Zugangsdatum.
* `risiko`: Die ursprüngliche Zeile war über den Bestand nicht erfüllbar (Planungsfehler); die Korrektur ist sachlich richtig, aber der Rollenwechsel „Plan ändern“ ist nicht als Artefakt belegt, und der Verifier prüft sonst gegen den vom Implementer selbst weichgestellten Maßstab.
* `verifizierbar`: Review-only (Verifier bewertet, ob er die geänderte DoD-Zeile gelten lässt)

### LOW-1

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.9 (b), §3.7
* `pfad`: `tools/messung/ollama-cloud/lib.sh` · „Aufraeumen in 07-aufraeumen.sh arbeitet nur darauf“ und „Varianten: std (Standardprofil, C7), dis (Sandbox ohne Egress: Modus discover),“
* `klasse`: Kommentar verweist auf nicht vorhandenes Symbol
* `befund`: Das Aufräum-Skript heißt `05-aufraeumen.sh`; die Variante `dis` gibt es nicht (`04-lauf.sh` startet `discover` als zweiten Container aus dem Image von `res`).
* `risiko`: Der nächste Änderer sucht ein Skript und eine Variante, die nicht existieren.
* `verifizierbar`: Review-only

### LOW-2

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Aussage nennt nur, was gelesen oder gemessen wurde)
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/README.md` · „über dieses Verzeichnis und `tools/messung/ollama-cloud/` (…): 0 Treffer“
* `klasse`: Messaussage nicht reproduzierbar
* `befund`: `06-suche.sh` auf genau diese beiden Orte endet mit Exit 1: zwei Treffer in `README.md` (Zeile mit der `ollama signin`-Aussage und der Satz mit `GH_TOKEN=ghp_…`), Muster „Zuweisung“ und „Anmelde-URL“. Die Suche lief vor dem Schreiben dieser Zeilen oder zählt sie nicht; der Text sagt es nicht.
* `risiko`: Ein Wiederholungslauf (Verifier, Folge-Slice mit Konto) beginnt mit rotem Ergebnis und muss selbst klären, ob es ein Zugangsdatum ist.
* `verifizierbar`: ja, `tools/messung/ollama-cloud/06-suche.sh <Messordner> tools/messung/ollama-cloud` (Exit 1)

### LOW-3

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/README.md` · „in den sechs Dateipaaren unterscheidet sich **nur** der Exit-Code“
* `klasse`: Zählung falsch
* `befund`: Es gibt drei Paare (`discover`, `restricted`, `allow`, je Lauf 1 und 2), sechs Dateien; der Vergleich der `AUSSAGE`-Zeilen bestätigt die Sachaussage (nur `C3-launch-claude-tty-exit` 124/137).
* `risiko`: Gering; eine falsche Zahl in der Lesehilfe senkt das Vertrauen in die übrigen Zahlen.
* `verifizierbar`: Review-only

### LOW-4

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Aussage nur nach Lesen/Messen), `harness/verification.md` (Evidence)
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/README.md` · „Vorgänge gegen `ollama.com`/`registry.ollama.ai` im Lauf 0: Pull, Run und Launch (ohne und mit Pseudo-Terminal) je einmal in `discover`“
* `klasse`: Aussage ohne erhaltene Rohdaten
* `befund`: Lauf 0 ist ehrlich ausgewiesen (abgebrochen, Rohdaten beim Neuaufsetzen überschrieben, Ursache des Hängens nicht untersucht), nennt aber Einzelheiten des verlorenen Laufs (welche Aufrufe, je einmal), die kein Beleg trägt.
* `risiko`: Die Zahl der Vorgänge gegen die Cloud ist für ein späteres Kontingent relevant und steht auf Erinnerung.
* `verifizierbar`: Review-only

### INFO-1

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/rohdaten-01.txt` · „== Messstand, 2026-10-24T08:38:04.111Z“
* `klasse`: Messstand nach dem Lauf aufgenommen
* `befund`: `rohdaten-01.txt` und `rohdaten-01-vorher.txt` tragen 08:38, nach den Läufen (04: 08:18 bis 08:34) und nach dem Aufräumen (05: 08:34); die Tabelle „Reihenfolge“ führt 01 als ersten Schritt. Die Umgebungsprüfung der Läufe selbst steht in `rohdaten-04-*` (`umgebung-container-verbotene-namen=0`) und in `04-lauf.sh` (`umgebung_pruefen`); der Zeitpunkt der Datei „vorher“ ist im Protokoll nicht genannt.
* `risiko`: Der Leser nimmt „vorher“ als Stand vor der Messung.
* `verifizierbar`: Review-only

### INFO-2

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/README.md` · „`sed`, `awk`, `curl` laufen nur in Containern“
* `klasse`: Verfahrensabweichung nicht im Protokoll
* `befund`: Der Auftrag nennt zweimal Host-`sed` während der Arbeit (Verfahrensabweichung); weder README noch Slice erwähnen sie. Die Skripte selbst nutzen Host-`sed` nicht (`redact` und `tool()` laufen in Containern); ich kann die zwei Aufrufe nicht aus dem Diff erkennen, also keinen Einfluss auf Rohdaten feststellen.
* `risiko`: Die Aussage des READMEs gilt für die Skripte, nicht für den gesamten Arbeitsvorgang; ob das für die Messung zählt, entscheidet der Verifier.
* `verifizierbar`: Review-only

### INFO-3

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `docs/messungen/slice-v1-ollama-cloud-messung/rohdaten-04-1.txt` · „[Query entfernt: Anmelde-URL] entfernt: Anmelde-URL]“; `docs/messungen/slice-v1-ollama-cloud-messung/rohdaten-04-allow-1.txt` · „== Modus restricted“
* `klasse`: Rohdaten-Artefakte
* `befund`: In den Sammeldateien `rohdaten-04-<n>.txt` ist die Schwärzung doppelt angewandt und nicht idempotent (zerstückelte Marke, nur Kosmetik, kein Wert sichtbar); `rohdaten-04-allow-<n>.txt` trägt die Überschrift „Modus restricted“, weil `inner.sh` für `allow` den Modus `restricted` übergibt (im README erklärt). `rohdaten-02-c7.txt` wird nicht durch `redact()` geschrieben (Ausgabe von `generate`, kein Zugang, Hand gelesen).
* `risiko`: Kosmetik; ein Leser, der die Dateien maschinell vergleicht, stolpert.
* `verifizierbar`: Review-only

### INFO-4

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `docs/plan/planning/in-progress/slice-v1-ollama-cloud-messung.md` · „bereinigter Umgebung (`env -u` der genannten Variablen“
* `klasse`: Plan-Text und Umsetzung weichen ab
* `befund`: Der Plan nennt `env -u`, umgesetzt ist `env -i` (strenger: nur sechs Variablen); die Umsetzung erfüllt die Absicht, der Plan-Satz ist veraltet.
* `risiko`: gering
* `verifizierbar`: Review-only

## Prüfung der Schwerpunkte

1. **Zugangsdaten.** Keine Token-Präfixe, Header, Schlüsselblöcke oder Anmelde-Query in Rohdaten und Skripten; Anmelde-URLs stehen als `https://ollama.com/connect?[Query entfernt: Anmelde-URL]`; `id_ed25519` und `.pub` erscheinen nur als Pfad und Größe (387/81 Byte); „Your new public key is“ kommt in den neuen Rohdaten nicht vor (nur in alten Protokollen, öffentlich, siehe MEDIUM-1); der Serverlog-Auszug (`server config env=map[…`) bricht bei 220 Zeichen ab und enthält keine Werte. `redact-probe.sh` grün; `redact()` deckt die Formen der Probe, die README bindet die Zusage ausdrücklich an diese Formen. Die erzeugten Dateien tragen nur `${localEnv:…}`-Verweise; der Hilfs-Container der CLI bekommt keine Host-Variable (`docker run` ohne `-e`).
2. **Messaussagen nachgezählt (C0–C7).** Version `0.40.6`, `launch`/`claude` in der Hilfe (`c0-launch-claude-in-help=4`), `claude 2.1.310`; Serverstart: `ollama.com` 2 Anfragen (`discover`/`allow`) gegen 4 (`restricted`); Pull: nur `registry.ollama.ai` (2 bis 4); `run`: „You need to be signed in …“ Exit 0 (`discover`/`allow`) gegen „lookup ollama.com … server misbehaving“ Exit 1 (`restricted`); `launch claude`: „requires sign in“ Exit 1 gegen echtes `claude` mit Fehler „Input must be provided …“, `api.anthropic.com` 2, `ollama.com` 12 abgelehnt; Stub: sieben `CLAUDE_CODE_*`; `OLLAMA_NO_CLOUD` (Run, Launch, Server-Log „Ollama cloud disabled: true“); Dateien (`model-recommendations.json` 2370 nur mit erreichbarem `ollama.com`, `config.json` 104 nur `restricted`, `.claude.json` 852 auf dem Wurzel-Dateisystem); Egress-Liste 11 Namen ohne `ollama.com`; C7 Warnungen und Zähler. Alle Einträge der Spalte Läufe stimmen; eine Zählung falsch (LOW-3). „Fail-open“ wird nicht behauptet: das README beschreibt nur, dass in `restricted` das installierte `claude` startet und warum der Ablauf davon abhängt ausdrücklich als nicht gemessen.
3. **Messvariante und Auslieferung.** Alle Handvarianten stehen in einer Tabelle (CLI statt `u-boot devcontainer up`, Zweitcontainer für `discover`, `docker exec` ohne `remoteEnv`, `dnsmasq` mit `log-queries`, Platzhalter für `claude`, Aufrufe ohne Eingabe, Lauf 2 mit vorhandenem Volume, ein Modellname). „Nicht gemessen“ ist vollständig, jede Zeile hat Aufruf und Bedingung; keine Zusage geht über das Gemessene hinaus („ohne Konto gemessen, Erfolgspfad nicht gemessen“ in der Kopfbedingung).
4. **DoD-3 und Lauf 0.** Siehe MEDIUM-1 und LOW-4; Lauf 0 ist offen ausgewiesen, der Verlust der Rohdaten wird nicht verschwiegen.
5. **Ablage, Kommentare, Plan.** `make ablage-pruefung` grün (Teil von `make gates`); `git diff --stat` gegen die Plan-Tabelle (§3): alle Pfade sind erfasst (Skripte, Protokoll, `docs/messungen/README.md`, Slice, Welle-Datei, `open/README.md`, Roadmap); Reviews, Register-Fortschreibung und Verify folgen zur Closure. Keine Kennung nach unten: das Protokoll (Rang 7) verweist aufwärts auf Lastenheft und seitwärts auf Slice/Messprotokolle. Skript-Kommentare bis auf HIGH-1 und LOW-1 nach §3.7 tragend.
6. **Aufräumen.** Keine Objekte mit Präfix; die drei „fremden“ Container in `rohdaten-05-nachher.txt` sind nicht mehr vorhanden und tragen kein Präfix.
7. **Hinweis `settings-json`.** Keine Datei im Diff enthält Anweisungen an Agenten oder Einstellungsänderungen (kein Pfad unter `.claude/`, `harness/`, `AGENTS.md`, `.harness/`); der Treffer ist `tools/messung/ollama-cloud/claude-stub.sh` (Platzhalter, gibt nur Namen aus) und Variablennamen `ANTHROPIC_*`/`CLAUDE_CODE_*` als Messdaten.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `lib.sh` `redact()`/`umgebung_pruefen()`/`bestand()`/`praefix_entfernen()` | geprüft, ohne Befund (Funktion); Kommentar siehe LOW-1 |
| `redact-probe.sh`, `06-suche.sh` | geprüft, ohne Befund |
| `01-messstand.sh`, `02-projekte.sh`, `03-lauf-vorbereitung.sh`, `04-lauf.sh`, `05-aufraeumen.sh`, `00-alles.sh`, `claude-stub.sh`, `inner.sh` | geprüft, ohne Befund |
| Rohdaten `rohdaten-02-erzeugt-*`, `rohdaten-04-*` (alle), `-up-*`, `-05*`, `-03` | gelesen/durchsucht, ohne Zugangsdaten |
| `docs/messungen/README.md`, `docs/plan/planning/` (Roadmap, Welle, `open/README.md`) | geprüft, ohne Befund (reine Verweis- und Zeigerpflege) |
| Produktcode, Spec, Lastenheft, ADRs | nicht berührt, nicht geprüft |
| rote Gegenprobe von `redact-probe.sh` | nicht gefahren (Grund oben) |
| Zusage „Messung wiederholt gleich“ bei neuem Lauf | nicht geprüft: Messung nicht wiederholt (Auftrag, Speicher) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Kommentar-Rest nach Teilersetzung · DoD-Wortlaut vom Implementer angepasst · Kommentar verweist auf nicht vorhandenes Symbol · Messaussage nicht reproduzierbar · Zählung falsch · Aussage ohne erhaltene Rohdaten

## Verdikt

**Merge-blockierend:** ja (HIGH-1, MEDIUM-1). Beide sind billig zu beheben (ein Kommentar; die DoD-Zeile vom Planner/Architect bestätigen lassen oder dem Verifier ausdrücklich vorlegen). Zugangsdaten und Messaussagen tragen.

**Übergabe:** Findings an den Implementer; Finding-Klassen in die Slice-Closure §7. Die DoD-Prüfung bleibt beim Verifier.
