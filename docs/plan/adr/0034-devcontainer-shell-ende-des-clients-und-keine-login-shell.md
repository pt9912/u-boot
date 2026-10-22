# ADR 0034: `u-boot devcontainer shell` — die Shell im Container endet nicht in jedem Fall mit dem Ende des Clients, das Ende von `u-boot` beendet den Client nicht in jedem Fall, die Shell ist keine Login-Shell; die Prämissen aus ADR-0031 werden berichtigt, die Entscheidungen von ADR-0031 gelten

**Status:** Accepted

**Datum:** 2026-10-22

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [ADR-0031](0031-devcontainer-shell-ohne-editor.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (der neue Unterabschnitt für `shell`, soweit er Signale, Terminal, die Lebensdauer der Shell und die Art der Shell beschreibt; die Spezifikation führt ihn nach dem Code nach). Das Lastenheft sagt dazu nichts ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gelesen: „interaktive Shell“, Exit-Code `0` nach Beginn der Shell; keine Aussage über das Ende der Shell nach dem Ende von `u-boot`, keine über eine Login-Shell); diese ADR ändert es nicht.

**Supersedes (teilweise):** vier Stellen von [ADR-0031](0031-devcontainer-shell-ohne-editor.md), und nur diese: (1) in Entscheidung 1, Spiegelstrich *Shell*, der erste Satz („**Shell:** eine Login-Shell, `bash`, wenn sie im Container liegt, sonst `sh`.“); (2) in Entscheidung 4 der Satz „Endet `u-boot` durch ein Signal, endet der Client mit ihm.“; (3) in Entscheidung 4, in der Aufzählung hinter „**Nicht belegt, Prüfaufträge des Umsetzungs-Slice:**“, die ersten zwei Glieder („dass die Shell im Container mit dem Ende des Clients endet (mit und ohne Terminal); das Verhalten von Strg-C ohne Terminal;“); (4) in den akzeptierten Negativen der Punkt „(2) `userEnvProbe` wird nicht nachgebildet, die Login-Shell liest ihre Profile selbst;“. Für diese vier Stellen ist diese ADR die geltende Fassung (Wortlaut in Entscheidung 1; dieselbe Form wie [ADR-0033](0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md) gegenüber ADR-0030 und ADR-0031). Das dritte Glied der Aufzählung (der Terminal-Pfad selbst in der CI) und alles Übrige von ADR-0031 gilt unverändert; ADR-0031 bleibt `Accepted` und wird nicht geändert.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0031](0031-devcontainer-shell-ohne-editor.md) Entscheidung 4 sagt ohne Vorbehalt „Endet `u-boot` durch ein Signal, endet der Client mit ihm“ und führt zwei Aussagen als nicht belegt: dass die Shell im Container mit dem Ende des Clients endet (mit und ohne Terminal), und das Verhalten von Strg-C ohne Terminal. Entscheidung 1 nennt die Shell eine „Login-Shell“, und der akzeptierte Nachteil (2) stützt sich darauf; dieselbe Entscheidung und die Fitness Function nennen die Argumentliste von `docker exec` ohne `-l`. Der Umsetzungs-Slice hat die Prüfaufträge gemessen (Docker-Tests im Adapter-Paket, dazu Läufe im Test-Image) und legt den Befund vor, weil sein Plan bei einer von einem Docker-Test widerlegten Annahme einen Befund an den Architect verlangt. **Eine neue Produktentscheidung steht nicht an;** die ADR prüft, was an den Aussagen hängt, berichtigt sie dort, wo eine `Accepted`-ADR sie trägt, und legt den Widerspruch um das Wort „Login-Shell“ in Richtung des Aufrufs aus, den die Argumentliste, die Fitness Function und der Handbuch-Weg tragen und den der Adapter ruft (Entscheidung 3).

**Gelesen am 2026-10-22** (jede Aussage über Code, Spezifikation oder Verhalten, `AGENTS.md` §3.10):

- `internal/adapter/driven/docker/devcontainer_shell.go`, `shellArgs`: `exec -i [-t] -u <Benutzer> -w <Workspace> [-e <NAME> …] <Kennung> <Shell>`, nach der Shell kein Argument. `probeShell`: `bash -c true`, sonst `sh -c true`. `runSession`: mit Terminal läuft der Client unter `context.WithoutCancel(ctx)` und ohne `cmd.Cancel`; ohne Terminal ist `cmd.Cancel` ein `os.Interrupt` an den Client.
- `cmd/uboot/main.go`, `main`: `signal.NotifyContext` mit `os.Interrupt` und `SIGTERM`; beide Signale brechen den Kontext ab, den Prozess beenden sie nicht. Andere Signale (etwa `SIGHUP`) sind nicht angemeldet.
- `internal/adapter/driven/docker/devcontainer_shell_test.go`: `TestShellArgs` (die Liste, nach der Shell kein Argument), `TestLauncher_Shell_ProbesThenOpensTheSession` (zwei `docker exec`: die Probe und die Sitzung), `TestLauncher_Shell_TerminalSessionIgnoresTheContext` und `TestLauncher_Shell_CancelWithoutTerminalEndsWithTheContextError`.
- `internal/adapter/driven/docker/devcontainer_shell_docker_test.go`: `TestShell_Docker_ClientEnd` (vier Fälle) und `TestShell_Docker_Terminal_UnderScript` (Strg-C; `shopt -q login_shell` wird gemessen und mit `t.Logf` ausgegeben, nicht behauptet). Der Container der Tests ist `docker run -d … mcr.microsoft.com/devcontainers/base:debian sleep 3600` (`newShellContainer`, `devBaseImage`).
- `docs/user/benutzerhandbuch.md`, 4.9 *Im Container arbeiten*: der Weg ohne den Befehl ist `docker exec -it -u <Benutzer> -w <Workspace> <Container> bash`, ohne `-l`. Die Vorlagen unter `internal/hexagon/application/templates/devcontainer/` setzen kein `userEnvProbe` (Suche in `internal`: kein Treffer).

**Messungen** (Docker 29.9.2 für Engine und Client, der Client ist der des Test-Images; Linux `amd64`, Kernel 6.8; lokale Engine; ein Lauf je Zeile; Rohausgabe der Läufe im Test-Image nicht im Repo; von mir nicht wiederholt, die Tests habe ich gelesen):

| # | Fall | Ergebnis | Beleg |
|---|---|---|---|
| 1 | ohne Terminal, Shell im Leerlauf (liest die Standardeingabe), Abbruch des Kontexts: der Adapter sendet SIGINT an den Client | die Shell endet | `TestShell_Docker_ClientEnd` |
| 2 | dasselbe, SIGKILL an den Client | die Shell endet | Lauf im Test-Image, nicht als Test |
| 3 | ohne Terminal, die Shell führt `sleep 300` aus, Abbruch des Kontexts | `bash` und `sleep` bleiben | `TestShell_Docker_ClientEnd` |
| 4 | mit Terminal (unter `script`), Shell im Leerlauf, SIGKILL an den Client | `bash` bleibt | `TestShell_Docker_ClientEnd` |
| 5 | mit Terminal, die Shell führt `sleep 300` aus, SIGKILL an den Client | `bash` bleibt | `TestShell_Docker_ClientEnd` |
| 6 | mit Terminal, SIGTERM, SIGINT oder SIGHUP an den Client | `bash` bleibt | Lauf im Test-Image, nicht als Test |
| 7 | der Client bei SIGINT und SIGTERM | Exit-Code 0 | Lauf im Test-Image, nicht als Test |
| 8 | mit Terminal, SIGKILL an den Elternprozess des Clients (`sh` an Stelle von `u-boot`) | der Client lebt weiter und liest und schreibt das Terminal | ein Lauf im Test-Image, nicht als Test (der Ersatz ist kein `u-boot`) |
| 9 | mit Terminal, Strg-C (0x03) während `sleep 300` | `sleep` endet, die Shell nimmt den nächsten Befehl an | `TestShell_Docker_Terminal_UnderScript` |
| 10 | `shopt -q login_shell` in der Sitzung mit Terminal | nein (`$0` ist `bash`) | `TestShell_Docker_Terminal_UnderScript`, ausgegeben, nicht behauptet |

**Messgrenzen:** der Container ist ein `sleep 3600` des Basis-Images, nicht ein von `devcontainer up` erzeugter; das Terminal ist ein Pseudo-Terminal von `script` aus dem Test-Image; der CI-Lauf ist nicht gemessen. **Nicht gemessen, nicht zugesagt:** das Schließen des Terminals (gemessen ist nur SIGHUP am Client, Zeile 6); Strg-C an der Tastatur, wenn die Standardeingabe von `u-boot` kein Terminal ist, `u-boot` aber ein steuerndes Terminal hat; die Wirkung eines Signals an `u-boot` als Lauf (die Wege über `main` und den Kontext sind gelesen und im Unit-Test des Adapters gepinnt, nicht als Prozess gestartet); andere Hosts als Linux `amd64`.

**Befund.** (1) Das Ende des Clients lässt die Shell im Container bestehen, außer bei einer untätigen Shell ohne Terminal, deren Standardeingabe endet (Zeilen 1 und 2); mit einem laufenden Befehl (Zeile 3) und mit Terminal (Zeilen 4 bis 6) bleibt sie, samt dem Befehl. Der Mechanismus ist nicht Gegenstand der Messung; gemessen ist, dass das Ende des Clients, wie es auch eintritt (SIGKILL ist der Fall, in dem der Client nichts mehr tun kann), die Shell nicht beendet. (2) Der Satz „Endet `u-boot` durch ein Signal, endet der Client mit ihm“ trägt ohne Terminal für SIGINT und SIGTERM an `u-boot` (der Kontext bricht, der Adapter sendet SIGINT an den Client); mit Terminal trägt er nicht: Der Adapter beobachtet den Kontext nach Beginn der Sitzung nicht, und ein SIGKILL an `u-boot` lässt den Client weiterlaufen (Zeile 8). Er ist falsch, soweit er alle Signale und beide Modi meint. (3) Die Shell ist keine Login-Shell (Zeile 10, dazu `shellArgs`); „Login-Shell“ in Entscheidung 1 und der Satzteil „die Login-Shell liest ihre Profile selbst“ im Nachteil (2) stehen im Widerspruch zu der Argumentliste derselben Entscheidung.

## Entscheidung

Jede Entscheidung nennt die Gegenoption, die verworfen ist.

1. **Die Prämissen lauten, wie gemessen.** Die vier Stellen lauten künftig:

   | Stelle in ADR-0031 | geltende Fassung |
   |---|---|
   | Entscheidung 1, Spiegelstrich *Shell*, erster Satz | „**Shell:** `bash`, wenn sie im Container liegt, sonst `sh`, jeweils ohne Argument (keine Login-Shell).“ |
   | Entscheidung 4, Satz „Endet `u-boot` durch ein Signal, endet der Client mit ihm.“ | „Das Ende von `u-boot` beendet den Client nicht in jedem Fall: Ohne Terminal unterbricht ein Abbruch des Kontexts (SIGINT, SIGTERM) den Client mit SIGINT; mit Terminal beobachtet der Adapter den Kontext nach Beginn der Sitzung nicht, und ein SIGKILL an `u-boot` lässt den Client weiterlaufen (Ersatzlauf, Messgrenzen im Kontext).“ |
   | Entscheidung 4, die ersten zwei Glieder der Aufzählung „Nicht belegt, Prüfaufträge“ | „Die Shell im Container endet nicht in jedem Fall mit dem Ende des Clients: Ohne Terminal endet eine untätige Shell (ihre Standardeingabe endet); eine Shell, die einen Befehl ausführt, und jede Shell mit Terminal bleiben bestehen, samt ihrem Befehl (Messtabelle, Zeilen 1 bis 6). Strg-C ohne Terminal: Der Adapter sendet bei Abbruch des Kontexts SIGINT an den Client; Strg-C an der Tastatur ist dort nicht gemessen.“ |
   | akzeptierte Negative, Punkt (2) | „(2) `userEnvProbe` wird nicht nachgebildet; die Shell ist keine Login-Shell und bekommt von `u-boot` nur die Umgebung aus Entscheidung 1;“ |

   **Gegenoption:** die Sätze stehen lassen und den Befund nur im Plan vermerken — der Plan hat im Rang weniger Gewicht als die ADR, und die Pläne für Handbuch und Spezifikation leiten ihre Sätze aus ADR-0031 ab; ein Vermerk neben einem `Accepted`-Satz ließe den Satz gelten.

2. **Die Entscheidungen von ADR-0031 gelten unverändert; der Adapter bleibt, wie er ist; am Code ist nichts zu ändern.** Geprüft wurde jede Entscheidung gegen ihren Anker: Die Entscheidungen 1 (bis auf das Wort „Login-Shell“, Entscheidung 3 dieser ADR), 2, 3, 5 und 6 hängen am `docker exec` des Hosts, an der Probe, der Ausgabe, den Exit-Codes, dem Port und der Zusage-Menge, nicht an der Lebensdauer der Shell; die Grenze in Entscheidung 3 („was in der Sitzung geschieht, ändert den Exit-Code nicht“) trägt auch für eine Shell, die das Ende des Clients überlebt. Entscheidung 4 gilt, soweit sie belegt ist: In der Sitzung bricht `u-boot` den `docker exec` mit Terminal nicht über seinen Kontext ab, und Strg-C erreicht die Shell (gepinnt). Berichtigt ist nur die Begründung „endet der Client mit ihm“.
   - **Gegenoption C — den Client an `u-boot` koppeln** (Signal-Weiterleitung auch mit Terminal, eigene Prozessgruppe, Elternbindung des Prozesses): ändert die Lebensdauer der Shell im Container nicht (das Ende des Clients lässt sie bestehen, SIGKILL am Client ist der stärkste Fall und gemessen); verschiebt den Terminal-Pfad, den der Test unter `script` pinnt (eine eigene Prozessgruppe ist nicht die des Vordergrunds, die Folgen für das Terminal sind nicht gemessen); ein SIGKILL an `u-boot` bliebe unabfangbar. Nicht gewählt.
   - **Gegenoption D — die Shell im Container an die Verbindung koppeln** (ein Wächter in der Shell, ein Zwischenprozess, eine Zeitgrenze der Shell): verändert die Shell, ihre Umgebung und ihren Prozessbaum und damit die Argumentliste von Entscheidung 1; eine Produktentscheidung ohne Bedarfsmeldung, das Lastenheft schweigt. Nicht gewählt.

3. **Die Shell ist keine Login-Shell; der Adapter ruft `bash` oder `sh` ohne Argument.** Die Argumentliste in Entscheidung 1 und die Fitness Function von ADR-0031 nennen kein `-l`, der Handbuch-Weg, den der Befehl ersetzt, ruft `bash` ohne `-l`, das Lastenheft sagt „interaktive Shell“ und stellt keine Anforderung an Profil-Dateien, und `TestShellArgs` pinnt die Liste. Die Auslegung ändert das Verhalten des Adapters nicht und verlangt keine neue Entscheidung des Projektinhabers: Die Entscheidungen a bis e in ADR-0031 enthalten keine Login-Shell. Die Zugangsdaten und die Variablen der erzeugten Konfiguration kommen über `-e` (Entscheidung 1), nicht über Profil-Dateien. **Gegenoption E — `-l`:** die Shell läse zusätzlich die Profil-Dateien; ihre Wirkung auf Umgebung und `PATH` in den Images, die `u-boot` erzeugt oder die Nutzer wählen, ist nicht gemessen; sie wich vom Handbuch-Weg und von der gepinnten Argumentliste ab und änderte nichts daran, dass `userEnvProbe` nicht nachgebildet ist (ADR-0031 Entscheidung 6). Nicht gewählt; ein Nutzerbericht über fehlende Variablen oder Pfade ist der Trigger, sie neu zu bewerten (eine Folge-ADR, ein Argument und ein Eintrag der Golden-Liste).

4. **Zusage-Menge für Handbuch und Spezifikation.** Zur Lebensdauer der Shell und zu ihrer Art sagen Handbuch und Spezifikation nur, was ein Test trägt (`AGENTS.md` §3.10, ADR-0031 Entscheidung 6).
   - *Zugesagt* (Test): `u-boot` ruft im Container nur die Probe und die Sitzung auf (`TestLauncher_Shell_ProbesThenOpensTheSession`); die Sitzung endet mit dem Ende der Shell, Exit-Code `0` (`TestShell_Docker_Terminal_UnderScript`, `exit`); eine Shell, die einen Befehl ausführt, bleibt samt Befehl bestehen, wenn die Sitzung von außen endet (`TestShell_Docker_ClientEnd`: ohne Terminal bei Abbruch des Kontexts, mit Terminal bei SIGKILL am Client); Strg-C erreicht mit Terminal die Shell (`TestShell_Docker_Terminal_UnderScript`). Eine zulässige Aussage für das Handbuch: „Beenden Sie die Sitzung mit `exit`. Beendet sich `u-boot` auf anderem Weg (etwa mit `kill -9`), kann die Shell im Container samt ihren Befehlen weiterlaufen.“ (der Satz sagt „kann“: lebt der Client weiter, lebt die Sitzung weiter; endet er, bleibt die Shell mit Terminal und mit einem laufenden Befehl.)
   - *Nicht zugesagt:* dass die Shell mit `u-boot` endet; dass ein Signal an `u-boot` die Sitzung beendet (mit Terminal beendet es sie nicht); das Verhalten beim Schließen des Terminals; Strg-C ohne Terminal; dass die Shell eine Login-Shell ist; dass Einstellungen ankommen, die nur Login-Shells lesen.
   - „Keine Login-Shell“ sagt das Handbuch nur nach einem Test, der es behauptet (die Messung in Zeile 10 wird bisher ausgegeben, nicht behauptet).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: Vermerk im Plan, ADR-0031 bleibt | kein neues Dokument | eine `Accepted`-ADR trägt einen falschen Satz und ein widersprüchliches Wort mit höherem Rang als der Vermerk; die Pläne leiten daraus ab |
| **B — diese ADR: Prämissen berichtigt (`Supersedes (teilweise)`), Entscheidungen bestätigt, kein Code** | die Datei von ADR-0031 bleibt unverändert; die Berichtigung steht dort, wo man den Satz findet; Handbuch und Spezifikation bekommen eine Zusage-Menge | ein Dokument mehr |
| C — den Client an `u-boot` koppeln | kein verwaister Client nach einem Signal | ändert die Shell nicht; verschiebt den gepinnten Terminal-Pfad; SIGKILL bleibt |
| D — die Shell an die Verbindung koppeln | keine verwaiste Shell | verändert, was der Nutzer als Shell bekommt; Produktentscheidung ohne Bedarf |
| E — `-l` | liest die Profil-Dateien | Wirkung nicht gemessen; weicht von Handbuch-Weg und Argumentliste ab; ersetzt `userEnvProbe` nicht |
| F — ADR-0031 ablösen und neu entscheiden | eine geschlossene Fassung | keine Entscheidung hängt an den berichtigten Sätzen; Neuentscheidung ohne Anlass |

## Konsequenzen

- Positiv: Wer ADR-0031 liest, findet die berichtigten Sätze; Handbuch und Spezifikation haben eine Zusage-Menge, die an Tests hängt; der Adapter und seine Tests bleiben, wie sie sind.
- Negativ: Eine Shell und ihre Befehle können nach dem gewaltsamen Ende von `u-boot` oder der Sitzung im Container weiterlaufen; der Befehl räumt sie nicht auf.
- Folgepflicht (eigene Schritte, nicht dieser ADR): Spezifikation und Sicht nach dem Code ([`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)): die Sätze der Entscheidung 4 mit den Messgrenzen, der Satz „Endet `u-boot` durch ein Signal, endet der Client mit ihm“ wird nicht übernommen; Handbuch mit dem Befehl: höchstens der Satz der Entscheidung 4. Keine eigene Folge-ADR, keine Änderung an Code oder Tests.
- **Akzeptierte Negative, mit Grund:** (i) Eine verwaiste Shell nach dem Ende des Clients ist das Verhalten des Clients und der Engine bei diesem Aufruf; der Befehl ruft den Client wie der Handbuch-Weg (Liste aus `TestShellArgs`) und beendet sonst nichts im Container; sie zu verhindern verlangte Option C oder D, beide ohne Bedarfsmeldung. (ii) Ein SIGKILL an `u-boot` lässt den Client weiterlaufen (ein Lauf mit Ersatzprozess): kein Programm fängt SIGKILL ab, auch Option C nicht. (iii) Mit Terminal beendet ein SIGINT oder SIGTERM an `u-boot` die Sitzung nicht (Entscheidung 4 von ADR-0031: der Kontext wird nach Beginn der Sitzung nicht beobachtet; gelesen und im Unit-Test gepinnt, nicht als Prozess gemessen); der Nutzer beendet die Shell mit `exit`. (iv) Strg-C an der Tastatur ohne Terminal an der Standardeingabe ist nicht gemessen: Ein Aufruf mit Pipe oder Datei auf der Standardeingabe ist ein Randfall, `u-boot` und Client stehen dann in derselben Prozessgruppe des Vordergrunds, und schlimmstenfalls bleibt die Shell wie unter (i). (v) Ob ein Image Einstellungen nur für Login-Shells setzt, ist nicht gemessen; die Shell liest, was eine interaktive Shell ohne `-l` liest.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Integrationstest (`//go:build docker`) | `TestShell_Docker_ClientEnd`: ohne Terminal endet eine untätige Shell, eine Shell mit Befehl bleibt samt Befehl; mit Terminal bleibt die Shell nach SIGKILL am Client, untätig und mit Befehl; der Test färbt sich rot, wenn Engine oder Client es anders enden lassen — dann werden diese ADR, Handbuch und Spezifikation gelesen | `make test-docker` |
| Integrationstest (`//go:build docker`) | `TestShell_Docker_Terminal_UnderScript`: Strg-C erreicht die Shell, die Sitzung endet mit `exit` und ergibt `ok` | `make test-docker` |
| Unit | `TestShellArgs`: nach der Shell kein Argument (kein `-l`); `TestLauncher_Shell_ProbesThenOpensTheSession`: zwei `docker exec`; `TestLauncher_Shell_TerminalSessionIgnoresTheContext` und `TestLauncher_Shell_CancelWithoutTerminalEndsWithTheContextError` pinnen den Umgang mit dem Kontext | `make test` |

Die Wortwahl von Handbuch und Spezifikation („nur die Sätze der Entscheidung 4“) ist nicht maschinell prüfbar; `AGENTS.md` §3.10 (Bedingung mit Test) und die Prüfung im Review decken sie.

## Re-Evaluierungs-Trigger

`TestShell_Docker_ClientEnd` färbt sich rot (Engine oder Client verhalten sich anders) oder der Docker-Client im Test-Image wird gehoben: Messtabelle, Entscheidung 1 und 4 neu lesen; ein Nutzerbericht, dass eine verwaiste Shell oder ein verwaister Befehl im Container stört (Option D neu bewerten); ein Nutzerbericht, dass `u-boot` in einer Terminal-Sitzung auf SIGTERM nicht endet oder ein Client nach dem Ende von `u-boot` das Terminal hält (Option C neu bewerten); ein Nutzerbericht, dass in der Shell eine Variable oder ein Pfad fehlt, die der Editor oder eine Login-Shell setzt (Option E neu bewerten, zusammen mit `userEnvProbe` aus ADR-0031); eine Messung in einem von `devcontainer up` erzeugten Container oder auf einem anderen Host als Linux `amd64`, die abweicht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-22 | Entwurf und `Accepted` durch den Architect im Auftrag des Projektinhabers, im Gespräch („ohne nachzufragen“); die ADR berichtigt Tatsachen-Prämissen und legt zwei in sich widersprüchliche Stellen von ADR-0031 nach der Argumentliste aus, die der Adapter ruft; sie trifft keine neue Produktentscheidung | Messungen des Umsetzungs-Slice zu ADR-0031 (`internal/adapter/driven/docker/devcontainer_shell_docker_test.go`) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-0034`.
