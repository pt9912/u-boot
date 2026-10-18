# Slice CR: Anzeige der wirksamen Egress-Liste mit Herkunft je Host (`LH-FA-DEV-008`, `LH-FA-CLI-005`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — zweiter Slice der Welle; ihm geht [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) voraus (Ergebnis zur Subdomain-Aussage), ihm folgt [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) (Umsetzung).

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; Erweiterung), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Abfragebefehle unter `--quiet`; Änderung der Aufzählung). Verwandt: [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md).
Anlass: Projektinhaber, 2026-10-18, Vorschlag (a) („Ok deine Vorschläge für die Verbesserungen“).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md)).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusage „die wirksame Egress-Allowlist samt Herkunft jedes Hosts ist abfragbar, offline und ohne laufenden Container“ steht im Lastenheft (Änderung **0.13.0 → 0.14.0**), in der vom Projektinhaber freigegebenen Fassung. Der Code-Slice beginnt erst danach.

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) sagt heute nur: Erlaubte Ziele sind `devcontainer.sandbox.egress.allow` und eine „dokumentierte“ Default-Allowlist. Ein Befehl, der die wirksame Liste zeigt, ist eine neue Produktzusage (neuer Befehl, neue Ausgabe); sie gehört ins Lastenheft, nicht in die Spezifikation. Ein neuer Abfragebefehl ändert außerdem die Aufzählung der Abfragebefehle in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).

### Gefundener Stand (gelesen am 2026-10-18)

- **Berechnung heute:** `egressHosts(dc, cloneURL)` (`internal/hexagon/application/devcontainer_egress.go`) bildet die Vereinigung aus `egressBaseHosts()`, `egressRegistryHosts()` (bei `nestedRuntime: podman` oder `docker`), `egressDockerAptHost()` (bei `docker`), `egressFeatureHosts()` je aktiviertem Feature, `cloneSourceHost(cloneURL)` und `dc.Sandbox.Egress.Allow`; das Ergebnis ist sortiert und ohne Duplikate. **Die Herkunft geht dabei verloren:** die Funktion gibt nur `[]string` zurück; `ghcr.io` stammt aus zwei Quellen (`egressRegistryHosts()` und `egressFeatureHosts()["homebrew"]`), `registry.npmjs.org` aus zwei (`node`, `claude-code`).
- **Verbraucher:** `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`) setzt `data.EgressHosts = egressHosts(dc, cloneURL)`; `egress-init.sh.tmpl` schreibt die Liste als `hosts="…"` in das Startscript. Die Clone-Quelle liefert `effectiveCloneSource` (dieselbe Datei): `devcontainer.sandbox.repository`, sonst die `origin`-URL aus `.git/config` (`readGitOriginURL`, liest die Datei über den `FileSystem`-Port; kein `git`-Aufruf, kein Netz).
- **Heute sichtbar:** die Liste steht nur im erzeugten Startscript (`generate devcontainer --dry-run` oder `--diff` zeigt es als Änderungsplan) und als Tabelle der Quellen im Handbuch (Abschnitt *Internet-Beschränkung (Egress)*) und in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion); ein Befehl, der sie samt Herkunft zeigt, existiert nicht (`internal/adapter/driving/cli/` kennt `add`, `config`, `doctor`, `down`, `generate`, `init`, `logs`, `remove`, `template`, `up`).
- **Aufzählung der Abfragebefehle:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) nennt „Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt.“; dieselbe Aufzählung steht in [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (der Spezifikations-Slice zieht sie nach).

**Befehlsname (freigegeben):** Der Projektinhaber hat Befehlsname und Wortlaut am 2026-10-18 freigegeben (§7); `u-boot egress list` steht im Lastenheft 0.14.0. Nicht gewählt wurden `u-boot config egress` (Untergruppe von `config`) und eine Option an `u-boot doctor`. [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) hält den Namen ein; ein anderer Name wäre ein CR-Nachlauf (§6).

### Wortlaut (alt/neu; am 2026-10-18 vom Projektinhaber freigegeben, im Lastenheft 0.14.0 eingetragen)

Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen die Verweise als Anker `#…` der Datei.

**Stelle 1 — [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), neuer Spiegelstrich** (nach dem Spiegelstrich „Erlaubte Ziele“; die übrigen Spiegelstriche bleiben unverändert).

Alt (Spiegelstrich 2, bleibt):

> - Erlaubte Ziele: `devcontainer.sandbox.egress.allow` (Liste von Hostnamen); die Default-Allowlist (gemeinsame Basis plus Ergänzungen je gewähltem Stack) ist dokumentiert.

Neu (zusätzlich, direkt danach):

> - Die wirksame Allowlist und die Herkunft jedes Hosts (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste) sind mit dem Befehl `u-boot egress list` abfragbar. Der Befehl liest nur die Projektkonfiguration (`u-boot.yaml`) und die Git-Konfiguration des Projekts, braucht weder Netz noch Docker noch einen laufenden Container, ändert nichts und zeigt dieselbe Liste, die `u-boot generate devcontainer` in das Startscript schreibt. Er weist aus, ob die Restriktion für die Konfiguration aktiv ist (Sandbox-Profil und `devcontainer.sandbox.egress.enabled`); ist sie es nicht, zeigt er die Liste, die bei Aktivierung gälte, und nennt das. Ein Name in der Liste schließt seine Subdomains ein.

Der letzte Satz („Ein Name … Subdomains ein“) trägt die Messung von [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md): Ergebnis E1 (bestätigt; Messstand Colima/arm64, `dnsmasq` 2.91, im Handbuch), gepinnt durch `TestE2E_SandboxDevcontainer_EgressSubdomains`; der Wortlaut bleibt unverändert.

**Stelle 2 — [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Aufzählung der Abfragebefehle.**

Alt:

> Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt.

Neu:

> Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `egress list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt.

**Stelle 3 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.13.0 |`

Neu: `| Version          | 0.14.0 |`

Historie-Zeile (Datum 2026-10-18; freigegebener Wortlaut, die Zeile im Lastenheft trägt zusätzlich die Klausel zu Stelle 4):

> | 0.14.0 | 2026-10-18 | [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion): Die wirksame Egress-Allowlist samt Herkunft jedes Hosts ist abfragbar (`u-boot egress list`; offline, ohne laufenden Container, ändert nichts); [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): der Abfragebefehl behält unter `--quiet` seine Ausgabe. | Vereinbarung mit dem Projektinhaber |

**Stelle 4 — [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Befehlstabelle (aus dem freigegebenen Befehl abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend).** Die Tabelle führt jeden Befehl des Vertrags; mit `u-boot egress list` in [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) führt sie ihn ebenfalls. Der Satz ist nicht wortgleich vorgelegt worden.

Alt: Tabelle ohne Zeile für den Befehl (letzte Zeilen: `u-boot config migrate` | Konfigurationsschema migrieren | Later; `u-boot template list` | Templates anzeigen | V1).

Neu (eine Zeile, zwischen `u-boot config migrate` und `u-boot template list`; Priorität V2 wie die Anforderung [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)):

> | `u-boot egress list`         | Egress-Liste mit Herkunft anzeigen | V2        |

Die Historie-Zeile 0.14.0 nennt die Befehlsübersicht und den Zustand (abgeleitet, Bestätigung ausstehend). Weitere Stellen im Lastenheft, die Befehle nennen (gelesen am 2026-10-18, Docker-`grep` über `spec/lastenheft.md`; nicht geändert): die Wertemenge des Feldes `command` im Absatz zu `--dry-run --json` in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (gilt für dateiverändernde Befehle; der neue Befehl lehnt `--dry-run` ab, siehe [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) §2), die Beispiele gruppierter Befehle (`template`, `config`) in [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) und in der JSON-Feldliste von [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) sowie das Laufzeit-Beispiel (`config get`) in [`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort); sie zählen Beispiele auf, keine Vollständigkeitsliste.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Spezifikation, Handbuch, CHANGELOG** — [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) und [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md); das Lastenheft ist Zielbild, der Code folgt.
- **Das Neu-Laden im laufenden Container** — [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md); eine eigene Zusage mit Sicherheitsentscheidung.
- **Eine Änderung der Default-Allowlist oder der Eingabe-Syntax** (Wildcard) — [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md).
- **Die Anzeige des Zustands des laufenden Containers** (ob die Restriktion beim letzten Start griff) — das leistet `u-boot doctor` (`devcontainer.sandbox.fallback.egress`, [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)); der neue Befehl bleibt offline und trifft keine Aussage über einen Container.

## 2. Definition of Done

- [x] **Lastenheft 0.14.0:** Stellen 1 bis 3 in der vom Projektinhaber freigegebenen Fassung (Wortlaut, Befehlsname und, falls geändert, der Subdomain-Satz), Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“); die Freigabe ist im Gespräch belegt und in §7 vermerkt (Datei-Beleg besteht nicht).
- [x] **Vorab-Suche und Nachzug:** Die Suche gegen ältere Lastenheft-Zusagen zum Befehlsumfang (`config`, Abfragebefehle, `--quiet`, Egress) ist in §3 festgehalten; Stellen, die dadurch mitgeändert werden müssen, sind genannt.
- [x] `make gates` grün; `make doc-trace` ohne neue Waise (die neue Anforderungsaussage hat im Code-Slice ihren Beleg; bis dahin steht sie in `ADR`-/Slice-Bezug). Läufe und Exit-Codes in §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `CHANGELOG.md` (`[Unreleased]`) nennt die Lastenheft-Änderung nicht als Produktänderung (kein Code); die öffentlichen Verträge ziehen die Code-Slices nach.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft (Stand der Prüfung bei der Slice-Closure in §7).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 bis 3 (§3 [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); Versionsfeld; §16 Historie); Wortlaut am 2026-10-18 freigegeben (§7). Stelle 4: Zeile in der Befehlstabelle [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und Klausel in der Historie-Zeile 0.14.0 (abgeleitet, Bestätigung ausstehend, §7) |
| `docs/plan/planning/done/slice-cr-sandbox-egress-anzeige.md` (diese Datei) | update | Freigabe, Vorab-Suche, Stelle 4, Abweichungen, Closure (§7) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ durch die Arbeitszeile dieses Slice ersetzt; die Closure setzt ihn zurück und zieht den Pfad der Welle-Zeile nach `done/` nach |
| `docs/plan/planning/in-progress/slice-v1-cli-egress-anzeige.md`, `docs/plan/planning/open/slice-v2-spec-sandbox-egress-ollama.md` | update | Verweis-Nachzug durch `make slice-mv` (beim Start und beim Wechsel nach `done/`); dazu Übergabe der Abgrenzung und der Auslegungsfragen aus dem Review (§7), je als gelesene Aussage mit Fundstelle |
| `docs/plan/planning/done/slice-v1-egress-wildcard-pruefung.md`, `docs/plan/planning/done/slice-v1-sandbox-egress-neuladen-entwurf.md`, `docs/reviews/2026-10-18-verification-slice-v1-sandbox-egress-neuladen-entwurf.md` | update | Verweis-Nachzug durch `make slice-mv` (Pfad von `open/` nach `in-progress/`, bei der Closure nach `done/`), ohne inhaltliche Änderung |
| `docs/plan/planning/open/README.md` | update | Verweis-Nachzug durch `make slice-mv`; die Closure setzt den Zustand der Zeile dieses Slice auf „geliefert“ (Stelle 4 offen) |
| `docs/plan/planning/welle-sandbox-egress-ollama.md` | update | Verweis-Nachzug durch `make slice-mv`; die Closure führt die Tabellenzeile dieses Slice nach und trägt das Release-Gate zu Stelle 4 in §3 ein |
| `docs/plan/planning/observations/` (`BEO-DOK/plan-aussage-ueber-code`, `BEO-WZG/host-werkzeug-ausserhalb-docker`: je `state.md` und eine Evidence-Datei; neu `BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht` und `BEO-DOK/vertragsaenderung-nach-review`) | update, neu | Register bei der Closure (§7): je Vorgang eine Evidence-Datei mit dem Dateinamen dieses Slice |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `CHANGELOG.md` | unverändert | der CR ändert kein Produktverhalten; ein Eintrag entsteht im Code-Slice. Gelesen: `[Unreleased]` führt nur Code-Änderungen |
| `docs/maintainer/releasing.md` | unverändert | gelesen: §3 (Vorbereitung, fünf Punkte) und §4 (Vorab-Prüfung) führen keine Gate-Liste für Vertragsbestätigungen; das Release-Gate zu Stelle 4 steht in der Welle-Datei §3 und fällt dem Release-Cut-Slice zu (§7) |

Die Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md)): `grep -i` über `spec/lastenheft.md` nach `config list`, `Abfragebefehl` und `allowlist` (am 2026-10-18 gelesen): Treffer sind der Absatz mit der Aufzählung in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Stelle 2), die Historie-Zeile 0.9.0 (Verlaufstext, bleibt unverändert, denn Historie wird nicht umgeschrieben), der Spiegelstrich zu `podman build` in [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) („Default-Allowlist“; unberührt) und die Spiegelstriche in [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Stelle 1). Der Implementer hat die Suche vor dem Lastenheft-Commit am 2026-10-18 in einem Docker-Container (`alpine:3`, `grep -n -i`) mit denselben drei Mustern wiederholt: dieselben Treffer wie oben (Zeilen der Aufzählung in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), `podman build` in [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), die Spiegelstriche in [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), Historie 0.9.0); weitere Stellen, die mitgeändert werden müssen, gibt es nicht. Die Zusage in [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) („`u-boot config` zeigt die gesamte Konfiguration“) bleibt: die wirksame Egress-Liste ist abgeleitet, keine Konfiguration.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) liegt in `done/` (Ergebnis zur Subdomain-Aussage ist Teil des Wortlauts); der Projektinhaber hat den Wortlaut und den Befehlsnamen freigegeben (Gespräch).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Wortlaut verlangt mehr als zwei Anforderungen.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts; das Lastenheft bleibt bei 0.13.0, kein Code-Slice startet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

- **Der Befehlsname ändert sich zwischen Freigabe und Architect-Entscheidung** ([`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md)). — **Ausgang:** weiter offen → [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) (Datei in `open/`): das Lastenheft nennt `u-boot egress list`, der Architect hält den Namen ein; sonst CR-Nachlauf (Review des Diffs der Änderung, nicht der ganzen Datei).
- **Der Subdomain-Satz steht im Lastenheft, ohne dass die Messung ihn trägt** (E2 oder E3 des Wildcard-Slice). — **Ausgang:** entfallen; Ergebnis E1, der Satz bleibt im Wortlaut.
- **Ein Review läuft vor der Freigabe und deckt die spätere Fassung nicht ab** (Lehre der Vorgänger-Welle: Gezählt, nicht verkörpert, [`welle-devcontainer-sshd-claude-code`](../done/welle-devcontainer-sshd-claude-code.md) §7). — **Ausgang:** entfallen: der Review lief über die freigegebene Fassung (Inhalts-Commit `21df190`, nach der Freigabe). Die Fix-Runde `36fb552` änderte das Lastenheft danach noch einmal (Stelle 4); sie ist nicht erneut reviewt, der Verifier hat sie gelesen (Verification-Report V-2) — dieselbe Klasse, im Register als [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) geführt.
- **Die Stelle 4 (Befehlstabelle, Historie-Klausel) wird vom Projektinhaber nicht bestätigt oder in der Priorität geändert.** — **Ausgang:** weiter offen → Release-Gate (§7, Welle-Datei [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §3): Lehnt er die Zeile ab oder ändert die Priorität, ist das ein CR-Nachlauf (Review des Diffs der Änderung, nicht der ganzen Datei).

## 7. Closure-Notiz

**Freigabe des Wortlauts:** Der Projektinhaber hat Stellen 1 bis 3 und den Befehlsnamen `u-boot egress list` am 2026-10-18 im Gespräch ausdrücklich freigegeben („ok“ auf die wörtlich vorgelegten Texte); der Wortlaut der Stellen 1 und 2 steht unverändert im Lastenheft 0.14.0, die Historie-Zeile trägt das Datum 2026-10-18.

**Nicht freigegebener Satz (Stelle 4):** Die Zeile für `u-boot egress list` in der Befehlstabelle [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und die zugehörige Klausel der Historie-Zeile 0.14.0 sind aus dem freigegebenen Befehl abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend. Priorität V2 wie die Anforderung [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).

**Ist-Stand des Starts:** Der Move nach `in-progress/` (Commit `8a7a827`) erfolgte mit `Verantwortlich: —`; die Start-Bedingung „`Verantwortlich:` gesetzt“ aus §4 ist nicht erfüllt, das Feld bleibt `—`.

**Abgrenzung der Anzeige (Bedarf der Folge-Slices):** Das Lastenheft nennt die Liste „die wirksame Allowlist“ und sagt, der Befehl zeige „dieselbe Liste, die `u-boot generate devcontainer` in das Startscript schreibt“; es sagt nicht, dass ein laufender Container die Liste seines Image-Standes trägt. Gemessen im Vorgänger-Slice [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md): Stoppen und `devcontainer up` behält die alte Liste, erst ein Neuaufbau übernimmt die neue. Der freigegebene Wortlaut bleibt unverändert; die Abgrenzung steht als Bedarf in [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) (Hinweis in der Ausgabe) und [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) (Spezifikation) §1.

**Übergabe der Auslegungsfragen und der Folge-Stellen:** Auslegung von „Restriktion nicht aktiv“ unter `--quiet` und Reichweite der Clone-Quelle: [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) §1; weitere Stellen mit der Aufzählung der Abfragebefehle ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Kommentar in `root.go`, Handbuch): dort und in [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) §1.

**Stelle 4: Bestätigung offen, Gate vor Release.** Der Projektinhaber hat am 2026-10-18 angeordnet, bis zum Release ohne Rückfragen weiterzuarbeiten; der Slice schließt deshalb mit unbestätigter Stelle 4, und deren Bestätigung ist ein verbindliches Gate **vor dem Release** (nicht vor `done/`). Das weicht vom Urteil des Verifiers ab („Der Slice ist bis dahin nicht schließbar“, Verification-Report, Gesamturteil); es ist die Entscheidung des Projektinhabers. Zustand im Lastenheft (unverändert gelassen): die Zeile `u-boot egress list` in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Zeile 1121) und die Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“ in der Historie-Zeile 0.14.0 (Zeile 1567). Bei der Bestätigung setzt eine Änderung genau dieser Klausel sie auf den Ist-Stand oder streicht sie (Verification-Report V-1) und dieser Abschnitt trägt den Vermerk mit Datum; lehnt der Projektinhaber die Zeile ab oder ändert die Priorität, ist das ein CR-Nachlauf (Review des Diffs der Änderung). Das Gate steht in den Closure-Kriterien der Welle [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §3 (Release-Gate); [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) bleibt unverändert (Plan-Tabelle §3), der Release-Cut-Slice der Welle nimmt das Gate in seine Vorbereitung auf.

**`Verantwortlich:`** bleibt `—`, wie bei den Vorgänger-Slices der Welle: der Projektinhaber hat keine Person gesetzt, der Autor `pt9912` steht in der Rolle Planner. Die Start-Bedingung aus §4 ist damit nicht erfüllt (Ist-Stand des Starts oben).

**Belege:** Review-Report [`2026-10-18-slice-cr-sandbox-egress-anzeige`](../../../reviews/2026-10-18-slice-cr-sandbox-egress-anzeige.md) (0 HIGH, 1 MEDIUM, 4 LOW, 4 INFO; F-1 bis F-7 in `36fb552` behoben, begründet oder übergeben, F-8 ohne Handlungsbedarf, F-9 ausgewiesen) und Verification-Report [`2026-10-18-verification-slice-cr-sandbox-egress-anzeige`](../../../reviews/2026-10-18-verification-slice-cr-sandbox-egress-anzeige.md) (GO mit Auflage; V-1 bis V-6). Sensoren des Verifiers (nach `36fb552`): `make gates` Exit 0, `make doc-trace` 83 Anforderungen und 0 Waisen, `make docs-check` 0 Befunde.

**Vorab-Suche:** Die Suche in §3 (Muster `config list`, `Abfragebefehl`, `allowlist`; Docker-`grep` in `alpine:3`) fand die Befehlsübersicht [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) nicht; der Review fand sie (F-1), Stelle 4 schließt die Lücke. Der Verifier hat die übrigen Stellen mit Befehlsnennungen ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort)) als Beispiel-Aufzählungen gelesen und unverändert gelassen.

**Was hat funktioniert:** Der Review wiederholte die Vorab-Suche und fand die Lücke, die sie hatte; der Verifier verglich die Stellen 1 bis 3 Zeile für Zeile mit dem freigegebenen Wortlaut (wortgleich) und wies Stelle 4 als nicht freigegeben aus, statt sie als freigegeben durchgehen zu lassen. Der Subdomain-Satz trägt die Messung E1 und den Docker-Test (Review F-8).

**Was ging anders als geplant:** Die Befehlsübersicht fehlte im Wortlaut-Vorschlag (F-1); die Fix-Runde fügte Stelle 4 ohne Freigabe des Wortlauts hinzu und ließ sie als abgeleitet ausweisen. Plan-Tabelle ohne die Verweis-Nachzüge und Zustandsmarker „Vorschlag“ nach der Freigabe (F-3, F-4) standen nach dem Review noch im Plan; beides ist behoben. Der Slice schließt mit offenem Gate (siehe oben).

**Steering-Loop-Eintrag:** gezählt, nicht verkörpert; mit diesem Slice wurde keine Regel und kein Sensor verkörpert. Lerneintrag als **geschärfte Regel für die Vorab-Suche** (Kandidat, erstes Auftreten): Eine Vertragsänderung, die einen Befehl einführt, sucht in der Vorab-Suche zusätzlich nach dem Namen eines benachbarten Befehls (hier `template list`) und der Befehlsübersicht, nicht nur nach Stichworten der Anforderung — Register [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), 1×. Zweites Auftreten der Klasse „Vertrag ändert sich nach dem Review, ohne Nachlauf-Review über die Änderung“ (erstes: [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) §7): Register [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), 2×, unter der Schwelle von drei. Der Rest läuft im Register weiter (unten).

**Beobachtungs-Register (`../observations/`):** [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): `evidence/slice-cr-sandbox-egress-anzeige.md` ergänzt (achtzehnter Vorgang, Review F-3 und F-4: Plan-Tabelle ohne Verweis-Nachzug, Zustandsmarker nach der Freigabe nicht auf den Ist-Stand; Rückfall nach der Verkörperung; Zählung nach der Evidence-Regel „ein Vorgang zählt einmal“), `state.md` nachgeführt, eine weitere Schärfung nicht entschieden. [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md): `evidence/slice-cr-sandbox-egress-anzeige.md` ergänzt (vierter Vorgang: lesende Host-Aufrufe im Review-Lauf, F-9, im Verifier-Lauf, V-5, und im Closure-Lauf; Verkörperung weiter nicht entschieden, Entscheidung beim Projektinhaber). Neu: [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (1×) und [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (2×). Kein Eintrag erreicht mit diesem Vorgang erstmals die Schwelle von drei; die zwei Einträge, die darüber liegen, tragen ihre Stände.

**Offener Hinweis (V-3):** Der Tabellenkopf der Befehlsübersicht („Priorität: MVP/V1 gemischt (siehe Spalte)“) nennt weder `Later` (bestand schon) noch das neue `V2`; die Spalte ist maßgeblich. Nicht geändert (das Lastenheft bleibt, siehe Gate); beim Bestätigen von Stelle 4 mit zu entscheiden.

**Folge-Slices:** [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) (Code, Befehlsname, Hinweis zur Abgrenzung, Auslegungsfragen, Stellen mit der Aufzählung) und [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) (Spezifikation) — Dateien in `open/`; neue Folge-Slices entstehen nicht.

**Risiken aus §6:** jedes trägt in §6 genau einen Ausgang: entfallen (Subdomain-Satz, Review vor Freigabe), weiter offen (Befehlsname → Code-Slice, Stelle 4 → Release-Gate).

**Grenzen:** (1) Die Freigabe des Wortlauts ist nur im Gespräch belegt (Verification-Report Abschnitt 8). (2) Die Fix-Runde `36fb552` (drei Lastenheft-Zeilen) ist nicht erneut reviewt; der Verifier hat sie gelesen (V-2). (3) „Die wirksame Allowlist“ nennt im Lastenheft die Grenze zum Stand eines laufenden Containers nicht (F-2, V-4); die Grenze lebt in den Plänen der Folge-Slices, im Code-Slice als Pflichtpunkt (Hinweis in der Ausgabe, mit Test). (4) `make test-docker`, `make verify-depguard` und `make doc-immutable` liefen nicht (kein Produktcode, keine ADR im Diff); der Docker-Test zum Subdomain-Satz lief nicht neu (Messstand im Vorgänger-Slice).

**Verifikation der Closure:** Plan-Tabelle §3 gegen `git diff --stat 1b833ea..HEAD` nach der letzten Änderung: siehe Abgleich unten; Sensor-Läufe der Closure (`make gates` Exit-Code von `make` selbst, `make doc-trace`, `make docs-check`) unten.

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: [`slice-v1-cli-egress-anzeige`](../in-progress/slice-v1-cli-egress-anzeige.md) und [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) liegen in `open/`. Register: die vier genannten Kennungen sind Verzeichnisse im Register und tragen je ein nicht leeres `evidence/`.

**Abgleich der Plan-Tabelle und Sensor-Läufe der Closure:** `git diff --stat 1b833ea..HEAD` nach dem letzten Inhalts-Commit vor diesem Eintrag (`6772aba`): 24 Einträge — Lastenheft, dieser Slice (Move mit Inhalt, im Diff als Löschung unter `open/` und Neuanlage unter `done/`), Roadmap, `open/README.md`, zwei Folge-Slices in `open/`, Welle-Datei, zwei Slices in `done/`, Review-Report, Verification-Report, Verification-Report des Vorgänger-Slice und elf Dateien des Registers (zwei `observation.md`, vier `state.md`, fünf Evidence-Dateien). Jeder Eintrag hat eine Zeile in der Plan-Tabelle §3, jede Zeile mit Datei-Eintrag einen Diff; `CHANGELOG.md` und `docs/maintainer/releasing.md` stehen nicht im Diff (gelesen, §3). Läufe nach `6772aba` (Exit-Code des `make`-Prozesses selbst): `make gates` Exit 0, `make doc-trace` Exit 0 (83 Anforderungen, 0 Waisen), `make docs-check` Exit 0 (431 Dateien, 0 Befunde). Dieser Eintrag ändert nur Text in dieser Datei; ein erneuter Lauf nach ihm steht im Handoff der Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer für Lastenheft-Änderungen unter den offenen Beobachtungen (`BEO-APP/podman-static-pin-alterung`, `BEO-E2E/*`, `BEO-WZG/kopplung-nicht-nachgezogen` betreffen Code und Werkzeug). [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) gilt (verkörpert in `AGENTS.md` §3.10): jede Aussage dieses Plans über Code ist an Datei und Symbol gelesen.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
