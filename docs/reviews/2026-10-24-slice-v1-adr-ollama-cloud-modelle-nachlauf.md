# Review-Report: slice-v1-adr-ollama-cloud-modelle (Nachlauf) — 2026-10-24

**Review-Art:** Design, Nachlauf — [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Status `Proposed`) nach den Fixes aus Commit `d5a030cf` gegen den Vorreport [`2026-10-24-slice-v1-adr-ollama-cloud-modelle.md`](2026-10-24-slice-v1-adr-ollama-cloud-modelle.md) (F-1 bis F-16); Maintainability, nicht gegen die DoD.

**Gegenstand:** Commit `d5a030cf` (Fixes) auf `e9f94698` und `5034770b`; Datei `docs/plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `d5a030cf` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24 · frischer Kontext, kein Übernehmen von Einschätzungen des Architect.

**Eingangs-Kontext:** `harness/review.md`, Reviewer-Skill, `AGENTS.md` (§3.7, §3.9, §3.10), die ADR (ganz), der Vorreport, die Messprotokolle [`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md) und [`slice-v1-ollama-cloud-messung-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md) (Stichproben in Rohdaten: `rohdaten-00-vorlauf.txt`, `rohdaten-01.txt`, `rohdaten-05-std-1-111134.txt`; Protokoll-Zeilen zu Abschnitt F, Handvarianten, Zählung 28/14). Code an der Fundstelle gelesen: `featureRelaxationWarnings`, `agentCredentialEnv`, `devcontainerTemplateData` (Quellen von `RunArgs`), `featureCatalogueEntry.relaxations` und Eintrag `go`, `egressFeatureHosts` samt Doc-Kommentar, `templates/devcontainer/devcontainer.json.tmpl` (`runArgs` im Zweig `.Sandbox`), `cli/warnings.go` (Sichtbarkeit von `SecurityRelaxation`). Spec: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1, [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) Sandbox-Regel, [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidungen 2, 4, 5, 7, 8.

**Sensors:** `make doc-immutable RANGE="e9f94698^..HEAD"` — 4 Commits, 0 Befunde. `git diff e9f94698^ HEAD` über die Datei von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) — leer. `make docs-check` — 841 Dateien, 0 Befunde. Nicht gelaufen: `make gates`, `make test-docker`, `make verify-depguard` — der Diff trägt nur Dokumente, keinen Go-Code, keinen Docker-Pfad, keine Importregel.

---

## 1. Prüfung der 16 Behebungen am Text

| Vor-ID | Urteil | Beleg im Text |
|---|---|---|
| F-1 (HIGH, Welle verlinkt/zitiert) | behoben | Kein Link und kein Zitat auf Welle, Roadmap oder Slice-Datei mehr (Suche über die Datei); `docs-check` 0 Befunde. Rest als N-4 (LOW). |
| F-2 (Handvariante gegen eigene Regel) | behoben | Regel im Abschnitt *Messung* neu gefasst (Handvariante sagt nur den Aufbau der Spalte *Pfad* zu); Z1, Z3, Z4, Z7 bis Z11 tragen „Handvariante“; Z10 behauptet kein Sandbox-Profil der Auslieferung mehr, *Nicht gemessen* führt das erzeugte Sandbox-Profil ohne Egress und begrenzt Z1 bis Z3, Z7, Z8, Z10 auf den `discover`-Container. |
| F-3 (Pfad Z2/Z6) | teilweise behoben | Z2 jetzt „Handvariante (`discover`-Container, Anmeldung von Hand)“, Z1 bis Z4 einheitlich. Z6 ist aufgeteilt, nennt aber für den Ohne-Konto-Teil „Auslieferung“, obwohl dieser Lauf laut Protokoll ohne Konto (Handvarianten-Tabelle) `devcontainer up` der Dev-Container-CLI nutzt: N-3 (LOW). |
| F-4 (`relaxations` als Muster) | behoben | Kontext und O1 stimmen mit dem Code: `relaxations` speist nur `featureRelaxationWarnings`; `RunArgs` kommt allein aus `devcontainerTemplateData` (Zeilen mit `podmanRelaxations`, `dockerRelaxation`, `egressRelaxation`); `runArgs` steht im Template nur im Zweig `.Sandbox`; „neuer Eintragstyp mit neuem Rendering-Pfad“ und die Konsequenzen („größer als ein Eintrag, eine Zeile“) sind angeglichen. |
| F-5 (`--quiet` gegen Ausweisung) | behoben | Entscheidung 8 trennt: Sandbox = `SecurityRelaxation`, sichtbar unter `--quiet`; Standardprofil = gewöhnliche Warnung, entfällt. Code stimmt (`featureRelaxationWarnings` nur mit `sandbox`; `ollamaWarnings` ohne `SecurityRelaxation`; `cli/warnings.go` zeigt `SecurityRelaxation`). Fitness Function angeglichen. Rest als N-1 (MEDIUM). |
| F-6 (Herkunft) | behoben | Definition der Spalte *Herkunft* an [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) angeglichen; e bis i „Empfehlung, nicht gefragt“, i zusätzlich „wird vor `Accepted` gefragt“; Entscheidung 10 nennt sich „Empfehlung, kein Beschluss“ (Punkt i). Punkt g und h tragen dieselbe Herkunft. |
| F-7 (Allaussage) | behoben | „trägt nirgends“ ersetzt durch „In dieser Lage trägt die erzeugte Zeile allein nicht“ mit Messlage (Linux, lokale Engine, Bindung `127.0.0.1`); Meldungstext in Entscheidung 8 neutral; *Verbleibendes Risiko*, Konsequenzen und Alternative C auf „in der gemessenen Lage“ gebracht; Re-Evaluierungs-Trigger nennt andere Engine. |
| F-8 (Supersedes-Spanne) | behoben | Kopf begrenzt auf die Vorgabe im Kontext von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) (Zeile 19 dort, wörtlich vorhanden), sagt, wo der Wortlaut „Anschluss an einen Ollama-Dienst des Hosts“ tatsächlich steht, und legt die Nummerierung fest („ohne Zusatz … dieser ADR“). Die genannten Bezüge stimmen mit [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md): 2 Wirkungsbereich (Z. 57), 5 Bindung, 7 Hosts, 8 Volume. |
| F-9 (Beleg-Zitat) | behoben | `rohdaten-05-std-1-111134.txt` steht in Entscheidung 4 und Z9; Datei ist Lauf 1, Standardprofil, `claude:2`. |
| F-10 (Messumgebung) | behoben | „Der Messlauf war bereinigt … gemessen ist also nur der Platzhalter“; `rohdaten-00-vorlauf.txt` („nach der Bereinigung im selben Prozess“) und `rohdaten-01.txt` (`umgebung-pruefung exit=0`) stützen genau das. |
| F-11 (28 Prompts) | behoben | „28 Aufrufe mit Exit 0 und Ausgabe, davon 14 `claude -p` mit Kurzantwort; bei `ollama run`/`launch` ist die erbetene Antwort nicht belegt“; Protokoll Zeile „Exit 0 mit Ausgabe 28 (Weg a 16, Weg b 12)“ und „14 / 14“ stimmen. |
| F-12 (Rückfallebene B) | behoben | B streicht Z1 bis Z4, Z7, Gültig-Teil Z8, Z9, Z10, Prompt-Teil Z12 und den Mit-Konto-Satz in Entscheidung 3. Gegengeprüft: Z5, Z6 (Ohne-Konto-Teil), Z11, `curl`-Teil Z12 hängen nicht am Konto. |
| F-13 (abgeleitet) | behoben | Z10 und Entscheidung 4: „Herkunft der `launch`-Antwort abgeleitet“ / „abgeleitet … nicht an der Umgebung von `claude` gelesen“, deckungsgleich mit dem Protokoll. |
| F-14 (Z13) | behoben | Z13 aus der Tabelle genommen, als „Produktversprechen, keine Messaussage“ mit dem noch zu liefernden Test geführt. Rest als N-5 (LOW). |
| F-15 (Folgepflicht-Lücken) | teilweise behoben | Code-Slice nennt den Doc-Kommentar von `egressFeatureHosts` („`ollama.com` is not needed“, im Code gelesen: Zeilen 78 bis 80); [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) *Sandbox-Regel* ist genannt. Der Satz, der durch `ollama-host` falsch wird, ist die Klammer „(`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion)“: N-2 (LOW). |
| F-16 (Reihenfolge, Bezug auf die Egress-ADR) | behoben | Reihenfolge Lastenheft-CR vor Code-Slice steht in der Folgepflicht; [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) wird in Entscheidung 3 im Text genutzt. |

## 2. Findings (neu durch die Fixes oder nicht vollständig behoben)

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| N-1 | MEDIUM | Entscheidung 8 stützt die Sichtbarkeit unter `--quiet` auf [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1. Deren Aufzählung der ausgewiesenen Sicherheitslockerungen nennt Capabilities, `--privileged`, Seccomp/AppArmor, Devices, Runtime-Socket des Hilfs-Containers und Login-Dienst; ein Host-Anschluss per `--add-host` steht nicht darin, und „alle anderen Hinweise und Warnungen entfallen“. Die Folgepflicht führt [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), nicht [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging). Szenario: Der Lastenheft-CR folgt der Folgepflicht, Rang 1 bleibt bei „entfallen“, der Code-Slice (Ausweisung als `SecurityRelaxation`, wie das Muster `featureRelaxationWarnings` es tut) steht gegen den Vertrag, oder der Test „sichtbar unter `--quiet`“ wird gegen das Lastenheft gelesen. Der Code selbst stimmt (`SecurityRelaxation` bleibt sichtbar). | Referenzmodell (Rang 1 gewinnt); `AGENTS.md` §3.10 (Bedingung ohne Entsprechung im Vertrag) | `0036` · „(Muster `featureRelaxationWarnings`, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1; Entscheidung O3)“ gegen „[`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): kein neuer Satz aus dieser ADR“ | Review-only (Lastenheft Zeile zu Ausnahme 1 und `cli/warnings.go` gelesen) | Folgepflicht-Liste ohne die Stelle, auf die sich die Aussage stützt |
| N-2 | LOW | Die Folgepflicht zu [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) schützt den Satz „keines von beiden ergänzt `runArgs` …“, der schon auf `sshd` und `claude-code` begrenzt ist. Falsch wird durch `ollama-host` die Klammer im selben Absatz „(`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion)“; sie wird nicht genannt. | `AGENTS.md` §3.9 (b), §3.10 | `0036` · „darf `ollama-host` nicht mitmeinen“; `spec/spezifikation.md` · „(`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion)“ | Review-only | Folgepflicht trifft nicht den Satz, der falsch wird |
| N-3 | LOW | Z6 führt für den Ohne-Konto-Teil „Auslieferung (Volume der CLI, Start von Hand)“. Der Lauf nutzt `devcontainer up` der Dev-Container-CLI statt `u-boot devcontainer up` (Handvarianten-Tabelle des Protokolls ohne Konto); Z5 weist denselben Unterschied aus („Auslieferung bis auf … `devcontainer up` der CLI“), Z6 nicht. | `BEO-DOK/messvariante-belegt-nicht-die-auslieferung`; Regel im Abschnitt *Messung* der ADR | `0036` · „Auslieferung (Volume der CLI, Start von Hand)“ | Review-only | Pfad-Spalte uneinheitlich für dieselbe Messlage |
| N-4 | LOW | Ohne Link, aber in der Sache steht noch Plan-Inhalt in der ADR: die Pflicht des Planners, „Closure-Kriterien des Vorgangs“ anzugleichen und „Zeilen der Folge-Slices und Roadmap nachzuziehen“, die Reihenfolge Lastenheft-CR vor Code-Slice, „das Kriterium des Vorgangs … der Planner gleicht an“. Die Referenzmatrix von [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md) führt ADR gegen Roadmap/Slice nach unten nicht als Kante. Kein Gate fängt die ungelinkte Nennung; Szenario: Der nächste Änderer liest die Closure-Anweisung als Teil der Entscheidung. Das HIGH aus F-1 (Link, Zitat als Entscheidungs-Eingabe) ist nicht mehr gegeben. | [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md); Reviewer-Skill (Referenzmodell) | `0036` · „Pflicht des Planners: die Closure-Kriterien des Vorgangs“ | Review-only (`docs-check` 0 Befunde) | ADR trägt Plan-Pflichten ohne Link |
| N-5 | LOW | Entscheidung 10 und die Folgepflicht sprechen von der „Zusage-Menge Z1 bis Z12“ und davon, dass der Lastenheft-Wortlaut keinen Satz ohne Zeile in der Tabelle *Zugesagt* enthalte. Der Satz, den [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) für den Host-Anschluss braucht (erzeugte Zeile in `runArgs`), hängt an Z13, die nicht in der Tabelle steht. | Entscheidung 10 der ADR (interne Regel) | `0036` · „Der Lastenheft-Wortlaut enthält keinen Satz, der in der Tabelle *Zugesagt* keine Zeile hat“ gegen „*Produktversprechen, keine Messaussage:* **Z13**“ | Review-only | Zusage-Menge und Tabelle decken sich nicht |
| N-6 | INFO | Zwei Zustandsmarker mit Ablaufpunkt: „wird vor `Accepted` gefragt“ (Punkt i) und „`Accepted` nach eigener Freigabe“ (Geschichte). Beim Übergang nach `Accepted` ist die Herkunft von i auf den Ist-Stand zu setzen; danach ist die Datei unveränderlich. | `AGENTS.md` §3.10 (Zustandsmarker) | `0036` · „**wird vor `Accepted` gefragt**“ | Review-only | Marker mit Ablauf |

## 3. Prüfung der weiteren Linsen

| Linse | Ergebnis |
|---|---|
| Referenzmodell: Link, Zitat oder Anker abwärts auf Welle, Roadmap, Slice | geprüft, ohne Befund (kein Link; ungelinkte Nennungen: N-4) |
| Zusage-Zeilen Z1 bis Z12 nicht weiter als gemessen, Pfad je Zeile | geprüft; Z1 bis Z5, Z7 bis Z12 stimmig mit Protokoll und Handvarianten; Z6: N-3 |
| Widersprüche Kontext / Tabelle / *Nicht gemessen* / Rückfall B | geprüft, ohne Befund (Z4 „von Hand“ gegen *Nicht gemessen* „über `config set`“ ist unterschieden, nicht widersprüchlich; Z10 gegen `discover`-Begrenzung stimmig; B vollständig, F-12) |
| Aussagen über Code am Bestand | geprüft, ohne Befund: `featureRelaxationWarnings` (nur Sandbox, `SecurityRelaxation`), `agentCredentialEnv` (drei Variablen), `egressFeatureHosts` (zwei Namen, Kommentar), `runArgs` nur im Zweig `.Sandbox` und nur aus drei Quellen in `devcontainerTemplateData`, `relaxations` ohne Fluss in `runArgs`; `ollamaWarnings` ohne `SecurityRelaxation`. Grenze: N-1 (Vertrag), N-2 (Spec) |
| Herkunftsspalte | geprüft, ohne Befund: i = „Empfehlung, nicht gefragt; wird vor `Accepted` gefragt“; e bis i „Empfehlung, nicht gefragt“; a bis d, j bis p „Projektinhaber, 2026-10-24“ (Gesprächsverlauf nicht im Repo, siehe unten) |
| [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) unverändert | geprüft, ohne Befund (`git diff` leer, `make doc-immutable` 0 Befunde) |
| ADR-Index | geprüft, ohne Befund (Zeile vorhanden, `Proposed`) |
| Spec-Straten / Nutzerdoku | geprüft, ohne Befund (Diff trägt nur ADR, Index, Reports) |

Nicht geprüft:
- Gesprächsverlauf mit dem Projektinhaber (Herkunft der Punkte a bis d, j bis p) — Quelle nicht im Repo.
- Rohdaten über die Stichproben hinaus (Zählungen je Datei, Antwortinhalte) — Gegenstand der Messreviews.
- Welle-Datei, Roadmap, Slice — nicht im Diff; ihr Nachzug ist Folgepflicht.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 1 |

Vorreport: 14 von 16 Behebungen stimmen im Text, 2 teilweise (F-3, F-15; Reste N-3, N-2). Keine Behebung ist nur behauptet.

**Finding-Klassen dieses Laufs:** Folgepflicht-Liste ohne die Stelle, auf die sich die Aussage stützt · Folgepflicht trifft nicht den Satz, der falsch wird · Pfad-Spalte uneinheitlich für dieselbe Messlage · ADR trägt Plan-Pflichten ohne Link · Zusage-Menge und Tabelle decken sich nicht · Marker mit Ablauf

## Verdikt

**Merge-blockierend:** N-1 (MEDIUM) ist vor `Accepted` zu klären (Folgepflicht um [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1 oder die Aussage anders stützen); LOW und INFO blockieren nicht. Keine HIGH-Findings.

**Übergabe:** an den Architect (Autor der ADR); Klassen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation; DoD-Konformität prüft der Verifier.
