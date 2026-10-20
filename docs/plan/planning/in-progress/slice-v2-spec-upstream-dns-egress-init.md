# Slice v2: Spezifikation — Abbruch des Egress-Startscripts auch unter `warn`, Kette `egress && docker`, Zusage zum Upstream-DNS (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nächster Release nach `v0.12.0` (Zuordnung des Planners; die Entscheidung über die Reihenfolge liegt beim Projektinhaber). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er folgt dem Code von [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md) (Reihenfolge „Spezifikation nach dem Code“, Muster [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md)); der Start verlangt, dass jener Slice in `done/` liegt (§4).

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation und Strenge der Sandbox-Fähigkeiten, `onUnavailable`); [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) (`Accepted`; Entscheidung 2: Abbruch auch bei `warn`; Entscheidung 5: die Spezifikationsfrage geht an diesen Slice), [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md).
Anlass: Auflage 1 des Verification-Reports [`2026-10-19-verification-slice-v1-egress-dns-upstream-schliessen`](../../../reviews/2026-10-19-verification-slice-v1-egress-dns-upstream-schliessen.md): Die Spezifikation trägt drei Stellen, die der Code nach [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) überholt hat, und für die Frage aus Entscheidung 5 gab es keinen Anker.

**Berührte Spec-Stellen:** `spezifikation.md`: [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz zu `result`, Absatz *Stand*), die Degradationstabelle von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Zeile der Egress-Capability), [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) (der Architect liest, ob der Satz zur Capability-Prüfung mit dem Abbruch aus anderem Grund vereinbar ist), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Allowlist; kein Satz zur DNS-Regel, der Architect liest, ob einer dorthin gehört), gegebenenfalls ein neuer Eintrag (nächste freie Kennung nach [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen); geprüft 2026-10-19: die höchste vergebene), §7 Historie. Das Lastenheft ändert dieser Slice **nicht**. Die Spezifikation nennt diesen Slice nie (Verweise zeigen aufwärts).

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-20)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation beschreibt den Egress-Abbruch und die Regel zum Upstream-DNS so, wie der Code sie liefert — als Zielbild, ohne Messstände, jede Bedingung nur mit einem Test, der sie pinnt — und nennt die Zusage samt Grenze an genau einer Stelle, falls der Architect eine eigene Stelle für nötig hält.

### Gefundener Stand (vom Verifier gelesen am 2026-10-19; beim Schreiben gegen Spezifikation und Code neu zu lesen)

- **[`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), Absatz zu `result`:** nennt `failed` nur bei der Richtlinie `fail`. Nach dem Code (`egress-init.sh.tmpl`, Abbruchzweige mit `write_status failed` und Exit 11) endet das Startscript unabhängig von `onUnavailable` mit Exit 11 und `failed`, wenn der Benutzer `nobody` fehlt, der `dnsmasq` nicht startet oder unter einer anderen UID läuft.
- **Degradationstabelle von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), Zeile der Egress-Capability:** `warn` führt zu „Egress-Restriktion entfällt“, nur `fail` zu Exit 11. Das gilt für die nicht gewährbare Capability; die Abbrüche aus dem Absatz oben stehen in der Tabelle nicht.
- **[`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), Absatz *Stand*:** „Die Kette `egress && docker` lässt `docker-init.sh` aus, wenn `egress-init.sh` mit Exit-Code `11` endet“ — nennt den Anlass nicht und gilt nach dem Code auch unter `warn`.
- **Offene Frage aus [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5:** Die Spezifikation führt keinen Satz zur Regel „Upstream-DNS nur für den Benutzer des `dnsmasq`“. Ob die Zusage samt Grenze (Grenzen: `sudo`, Adressen erlaubter Namen, Unternamen erlaubter Namen) eine eigene Spezifikationsstelle braucht, entscheidet dieser Slice.
- **Pinnende Tests (Namen gelesen):** `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` (Unterfälle `direct_path_closed`, `abort_branches`; `internal/e2e/sandbox_devcontainer_docker_test.go`), `TestE2E_SandboxDevcontainer_StatusFiles_Egress` (`internal/e2e/sandbox_status_docker_test.go`), `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser`, `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` (vier `failed`-Zweige) und `TestGenerateDevcontainer_SandboxEgressWithDocker` (die Verkettung in `postStartCommand`; das Verhalten der Shell-Kette ist durch keinen Test ausgeführt) in `internal/hexagon/application/generate_sandbox_test.go`, `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` in `internal/hexagon/application/devcontainer_fallback_test.go`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Templates** — der Code-Slice hat sie geliefert; stellt dieser Slice eine Abweichung fest, ist das ein Befund und ein neuer Slice, kein stilles Umschreiben.
- **Das Lastenheft** — Zielbild; [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) sagt zum DNS-Pfad nichts und wird nicht abgeschwächt ([`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5).
- **Messstände in der Spezifikation** (36 s Wartezeit, 43 gegen 7 s, `--dns`, Obergrenze der Anfragen je Sekunde, Messumgebung) — sie stehen im Handbuch und im Messprotokoll; die Spezifikation trägt Zusagen nur dort, wo ein Test sie pinnt.
- **Eine Zusage zu `nestedRuntime: docker`, benutzerdefiniertem Netz, Colima, `arm64` und CI** — ungemessen; was nicht gemessen ist, steht nicht in der Spezifikation.
- **Die Entscheidung, ob ein fehlender Benutzer `nobody` im Erstlauf offen (ohne Restriktion) bleiben darf** — Entscheidung des Projektinhabers; die Spezifikation beschreibt den Ist-Stand und übernimmt keine andere Entscheidung.
- **Handbuch und `CHANGELOG.md`** — der Code-Slice hat sie auf den Ist-Stand gesetzt; dieser Slice liest sie gegen die Spezifikation (Auflage „ein Vertrag, vier Stellen“) und ändert nur bei einer festgestellten Abweichung.

## 2. Definition of Done

- [x] **Drei überholte Stellen nachgeführt:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz zu `result` und Absatz *Stand*) und die Degradationstabelle von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) nennen den Abbruch mit Exit-Code `11` und `failed` auch unter `warn` (Benutzer fehlt, Resolver startet nicht, Resolver läuft unter anderer UID) und die Folge für die Kette `egress && docker`; jede Bedingung trägt einen Test aus dem Abschnitt *Pinnende Tests* oder wird gestrichen; eine **Fundstellen-Tabelle** (Datei und Symbol, gelesen) steht in §3.
- [x] **Frage der Spezifikationsstelle beantwortet:** Der Architect entscheidet, ob die Zusage („Upstream-DNS nur für den Benutzer des lokalen Resolvers“) samt Grenze einen eigenen Eintrag braucht (neue Kennung) oder in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) beziehungsweise [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) steht oder ganz entfällt; jede Zusage hat einen pinnenden Test, die Grenze (`sudo`, Adressen erlaubter Namen, Unternamen) steht als Grenze, nicht als Messstand; §7 Historie trägt je Änderung eine Zeile; `make doc-trace` meldet keine Waise, `make docs-check` keine Befunde.
- [x] **Ein Vertrag, vier Stellen:** Lastenheft, Spezifikation, Handbuch und `CHANGELOG.md` sagen über den Abbruch unter `warn` und die Regel dasselbe wie der Code (gelesen, Stichprobe des Verifiers).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb in der Closure dieses Slice geprüft.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz zu `result`, Absatz *Stand*), Degradationstabelle von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), gegebenenfalls [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) oder ein neuer Eintrag zur Regel, §7 Historie |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | unverändert, Lesepflicht | Gegenlesen gegen die Spezifikation; Änderung nur bei festgestellter Abweichung (der Code-Slice hat sie auf den Ist-Stand gesetzt) |
| `spec/lastenheft.md`, `spec/architecture.md` | unverändert, Lesepflicht | Lastenheft bleibt (Zielbild); der Architect liest, ob die Sicht den Abbruch nennt, und trägt die Lesestelle ein |
| `internal/**`, `Makefile`, `repo.mk`, `.github/workflows/` | unverändert | kein Code, keine Gate-Änderung |
| diese Slice-Datei, `open/README.md`, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |

### Ist-Stand der Plan-Tabelle (Implementer, gelesen gegen `git diff --stat`)

| Datei | Ist |
|---|---|
| `spec/spezifikation.md` | geändert: [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz zu `result`, Absatz *Stand*), Degradationstabelle von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (neue Zeile), [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) (ein Satz), neu [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), §7 Historie (eine Zeile). [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) unverändert |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | unverändert, gelesen (Befunde siehe unten) |
| `spec/lastenheft.md` | unverändert, gelesen: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) Absatz *Degradation* (Zeile 467) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Zeilen 469 bis 478) nennen den Abbruch aus anderem Grund nicht und widersprechen ihm nicht |
| `spec/architecture.md` | unverändert, gelesen: die Treffer zu „egress“ (Zeilen 120 bis 122, 161, 187, 242) betreffen `EgressListService` und den Befehl `egress list`; die Sicht nennt weder Startscripte noch ihren Abbruch |
| `internal/**`, `Makefile`, `repo.mk`, `.github/workflows/` | unverändert |
| diese Datei | geändert (Verantwortlich, DoD, dieser Abschnitt) |

### Entscheidung zur Stelle der Zusagen

**Eigener Eintrag [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)** (nächste freie Kennung: [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) war die höchste vergebene, geprüft per Suche im Dokument am 2026-10-20). Begründung: Drei Stellen brauchen dieselbe Aussage (Statusdatei in [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), Degradationstabelle in [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), Capability-Satz in [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability)); trüge [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) sie, stünden Upstream-DNS-Regel und Wiederholungslauf in einem Eintrag über Statusdateien, trüge [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) sie, in einem Eintrag über Hosts. Die Abbruchzweige, die Regel zum Upstream-DNS samt Grenze und die Struktur des Wiederholungslaufs sind eine Eigenschaft des Startscripts; sie stehen einmal in [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), die drei anderen Stellen verweisen darauf. [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) bekommt keinen Satz: es führt Hosts, und die Subdomain-Aussage dort ist durch `TestE2E_SandboxDevcontainer_EgressSubdomains` gepinnt, die Unternamen-Grenze von [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) ist als „nicht zugesagt“ formuliert und verweist dorthin.

**Abbruchzweige am Template gezählt** (`internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`): fünf Stellen mit `exit 11`. Vier gelten unter jeder Richtlinie (Benutzer fehlt, `nft`-Transaktion scheitert, Resolver startet nicht, Resolver unter anderer UID); die fünfte (Zeile 29, `NET_ADMIN` fehlt) gilt nur bei `fail`. Weitere Abbruchzweige gibt es nicht; die Spezifikation führt die vier und grenzt die fünfte auf die Degradationstabelle ab. Der Test `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` zählt fünf `write_status failed` vor `exit 11` (Fall `egress`).

### Fundstellen (Datei und Symbol, gelesen am 2026-10-20)

| Aussage der Spezifikation | Code | Pinnender Test |
|---|---|---|
| Benutzer einmal festgelegt: `user=` in der Resolver-Konfiguration, UID zur Laufzeit aus demselben Wert | Template Zeilen 39 bis 40 (`dns_user`, `dns_uid`), 65 (`user=$dns_user`) | `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser` (`internal/hexagon/application/generate_sandbox_test.go`), beide Richtlinien |
| Port 53 zu den Upstream-Resolvern nur für diese UID, jede Regel mit Benutzerbedingung | Template Zeilen 87 bis 90 | derselbe Test (zwei Regeln, jede Zeile mit `dport 53 accept` trägt `meta skuid`); `TestE2E_SandboxDevcontainer_EgressUpstreamDNS`, Unterfall `direct_path_closed` (`DIRECT udp/tcp example.com=none`, `RULES skuid`, `UID dnsmasq` gleich Regel) in `internal/e2e/sandbox_devcontainer_docker_test.go` |
| erlaubte Namen lösen auf, Verbindung kommt zustande | Template Zeilen 67 bis 70 | `…EgressUpstreamDNS/direct_path_closed` (`LOCAL api.github.com=1`, `LOCAL example.com=0`, `HTTPS api.github.com=200`, `STATUS result=ok`) |
| Grenze: als Benutzer des Resolvers kommt eine Antwort | — | `…EgressUpstreamDNS/direct_path_closed` (`LIMIT udp/tcp sudo-as-resolver-user=answer`) |
| Grenze: Adressen erlaubter Namen, Unternamen — nicht zugesagt | Template Zeile 86 (`@allowed`), 68 (`server=/…`) | kein Test pinnt das Offensein; deshalb nur als „nicht zugesagt“ formuliert |
| Abbruch 1, Benutzer fehlt, vor den Regeln, Exit 11 und `failed` | Template Zeilen 41 bis 45 | Unit: Verzweigungstext in `…UpstreamDNSOnlyForTheResolverUser`; `…EgressUpstreamDNS/abort_branches` (`ABORT user-missing exit=11 result=failed`, `tables=0`, `resolver-in-resolv.conf=0`, Meldung auf `stderr`) |
| Abbruch 2, `nft`-Eingabe nicht angenommen | Template Zeilen 97 bis 101 | Unit: letzter Check in `TestGenerateDevcontainer_SandboxEgress_RerunIsOneTransaction`; `TestE2E_SandboxDevcontainer_EgressRerun/failed_transaction` (erster Lauf: `tables=0 resolv=0`; Wiederholung: `chain-unchanged=1 same-resolver=1 resolv=1`; ungültige Regelzeile) |
| Abbruch 3, Resolver startet nicht; Regeln `drop` stehen, `resolv.conf` unverändert (erster Lauf) | Template Zeilen 120 bis 124 | Unit: Verzweigungstext; `…abort_branches` (`ABORT not-started`, `tables=2 drop=1 resolver-in-resolv.conf=0`) |
| Abbruch 4, Resolver unter anderer UID | Template Zeilen 125 bis 131 | Unit: `running_uid`-Prüfung und Verzweigungstext; `…abort_branches` (`ABORT uid-differs`, `tables=2 drop=1 resolver-in-resolv.conf=0`) |
| unabhängig von `onUnavailable` | Zeilen 44, 100, 123, 130 liegen außerhalb von `if [ "$on_unavailable" = "fail" ]` (Zeile 27) | Unit: beide Tests laufen mit `warn` und `fail`; E2E-Läufe mit `generateEgressProject(ctx, t, "warn")` |
| `NET_ADMIN` fehlt: Exit 11 nur bei `fail` (Abgrenzung) | Template Zeilen 25 bis 33 | `TestE2E_SandboxDevcontainer_StatusFiles_Egress` (`internal/e2e/sandbox_status_docker_test.go`); Unit: `failed: 5` im Fall `egress` von `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` |
| `doctor`: `failed` ist Fehler, Meldung und Hinweis nennen die vier Ursachen | — | `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` (`internal/hexagon/application/devcontainer_fallback_test.go`); `…StatusFiles_Egress` (`expectItem … SeverityError`) |
| Wiederholungslauf: eine `nft -f`-Eingabe, `add`/`delete`/`add` je Tabelle, kein anderer `nft`-Aufruf, geordneter Resolver-Wechsel, `resolv.conf` zuletzt | Template Zeilen 73 bis 101, 103 bis 132 | `…RerunIsOneTransaction`; `TestE2E_SandboxDevcontainer_EgressRerun/repeated_runs` (elf Läufe: Exit 0, `ok`, ein Resolver, zwei Tabellen, UID gleich Regel, `SET before=1 after=0`) |
| Kette `egress && docker` im `postStartCommand` | — | `TestGenerateDevcontainer_SandboxEgressWithDocker` (`generate_sandbox_test.go`, Zeile 1014): nur der Text der Verkettung; **kein ausgeführter Test**, deshalb steht in der Spezifikation „gepinnt ist die Verkettung als Text“ und die Folge nur als Semantik von `&&` |

### Befunde des Gegenlesens (nicht geändert)

- `CHANGELOG.md` (`[Unreleased]`, Absatz zum Upstream-DNS): nennt als Abbruchgründe nur den fehlenden Benutzer und den Resolver unter anderem Benutzer; „Resolver startet nicht“ fehlt dort (der nftables-Fall steht im Absatz darüber). Unvollständig, nicht widersprüchlich.
- Handbuch, Abschnitt *Fallbacks sehen* (Zeile 1143; Satz „endet das Egress-Script mit Exit-Code 11, läuft das Daemon-Script nicht“): stärker formuliert als die Spezifikation und als der Satz im Abschnitt zur Egress-Beschränkung („gelesen, kein Test mit Abbruch“); gleiche Aussage, unterschiedliche Stärke.
- [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz *Nicht ermittelbar*, Zeile 654): „weil ein früheres Script der Startkette mit Exit-Code `11` endete“ trägt dieselbe Verkettungs-Annahme ohne ausgeführten Test; an dieser Stelle nicht geändert (Folgefrage für den Reviewer: streichen oder wie im Absatz *Stand* auf den gepinnten Text begrenzen).
- Lastenheft: keine Abweichung; es schweigt zum Abbruch aus anderem Grund als der fehlenden Capability.

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1); die Reihenfolge entscheidet der Projektinhaber.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Frage der Spezifikationsstelle verlangt einen eigenen Eintrag mit mehreren Bedingungen, deren Tests erst entstehen müssen; dann trennt der Architect Nachführung und neuen Eintrag.
- `in-progress` → `open` (blockiert): eine Entscheidung des Projektinhabers zum fehlenden Benutzer im Erstlauf ändert den Code; dann folgt die Spezifikation dem geänderten Code.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Eine Spec-Bedingung ohne pinnenden Test** (Klasse [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)) — **Ausgang:** weiter offen bis zur Closure; jede Bedingung trägt einen Test aus §1 oder entfällt.
- **Die Verkettung `egress && docker` ist nur als Text gepinnt, das Verhalten der Shell-Kette nicht ausgeführt** — **Ausgang:** weiter offen bis zur Closure; die Spezifikation sagt „die Kette ist so verkettet“, nicht „das Docker-Script läuft nicht“, solange kein Test den Abbruch ausführt, oder der Architect benennt den Test, der es tut.
- **Der Code ändert sich vor dem Start** (Entscheidung des Projektinhabers zum fehlenden Benutzer im Erstlauf) — **Ausgang:** weiter offen; Rückführung `in-progress` → `open` (§4).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); sie erfüllt die Schwelle aus zwei von drei Achsen laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert in `AGENTS.md` §3.10; jede Bedingung der Spezifikation bekommt einen Test oder entfällt), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Schwelle erreicht, Ausgang beim Projektinhaber; die Spezifikation trägt keinen Messstand), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert; jede Aussage über den Code mit Fundstelle). Evidenz-/Diskrepanz-Risiko: mittel (drei überholte Stellen vom Verifier gelesen, nicht vom Planner).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
