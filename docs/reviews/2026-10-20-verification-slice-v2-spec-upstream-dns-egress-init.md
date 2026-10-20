# Verification-Report: slice-v2-spec-upstream-dns-egress-init — 2026-10-20

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v2-spec-upstream-dns-egress-init`](../plan/planning/done/slice-v2-spec-upstream-dns-egress-init.md) (`in-progress/`), Commit-Range `058bbf8..8bb5d59` (vier Commits: `b6706f1` Umsetzung, `eae98a6` Review-Report, `d3a7d06` Fix-Runde, `8bb5d59` Historienzeile). Entscheidungen: [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) (Entscheidung 2: Abbruch auch unter `warn`; Entscheidung 5: Spezifikationsfrage geht an diesen Slice), [`ADR-0027`](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) (Entscheidung 3: Zusage nicht weiter als Beleg; Entscheidung 4: fehlgeschlagene Transaktion endet mit Exit 11). Anforderungen: [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer). Spezifikation: [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-008.a`](../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability).

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle samt Ist-Stand, Fundstellen, Befunde des Gegenlesens im Slice), Review [`2026-10-20-slice-v2-spec-upstream-dns-egress-init`](2026-10-20-slice-v2-spec-upstream-dns-egress-init.md) (F-1 bis F-12), Fix-Commits `d3a7d06`, `8bb5d59`.

**Gesamturteil: GO mit Auflagen.** Die Abbruchzweige, Exit-Codes, Richtlinien-Unabhängigkeit und die Struktur des Wiederholungslaufs in der Spezifikation stimmen mit dem Template und mit gelesenen Tests überein; die Zielbild-Regel ist in den hinzugefügten Spec-Zeilen eingehalten; alle Review-Findings sind behoben oder im Ist-Stand erklärt, mit einer Ausnahme (V-4). Offen sind drei Spec-Bedingungen ohne pinnenden Test (zwei davon erst in der Fix-Runde eingeführt), eine Stärke-Differenz zwischen Handbuch und Spezifikation unter einer abgehakten DoD-Zeile und die Closure-Voraussetzungen.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass | eigener Lauf auf `8bb5d59`, sauberer Arbeitsbaum, Exit 0; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; Stempel durch `record-gates` |
| `make docs-check` | pass | `d-check: 537 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-trace` | pass | `83 Anforderung(en), 0 Waise(n)` |
| `make test-docker` | nicht gelaufen | Der Diff enthält keinen Code, kein Template, keinen Test (`git diff --stat`: `CHANGELOG.md`, Slice, Review-Report, Handbuch, `spec/spezifikation.md`). Die Aussagen sind gegen den **Wortlaut** der Docker-Tests gelesen, nicht gegen ein neues Ergebnis; die Tests selbst sind unverändert seit dem Verifier-Lauf des Code-Slice. |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt |
| `git diff --stat 058bbf8..HEAD -- internal Makefile repo.mk .github spec/lastenheft.md spec/architecture.md docs/plan/adr` | leer | Code, Gates, Lastenheft, Sicht und ADRs unverändert |

## 2. Aussagen gegen Template und Tests (Stichprobe: alle Bedingungen von [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf))

Template `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` selbst gelesen: fünf `exit 11` (Zeilen 29, 44, 100, 123, 130), jeder nach `write_status failed`; Zeile 29 liegt in `if [ "$on_unavailable" = "fail" ]` (Zeile 27), die anderen vier nicht. Alle Template-Zeilenangaben der Fundstellen-Tabelle im Slice stimmen.

| Aussage (Spezifikation) | Beleg | Urteil |
|---|---|---|
| vier Abbruchfälle unter jeder Richtlinie, `NET_ADMIN` fehlt nur bei `fail` ([`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) `result`, Degradationstabelle, [`LH-FA-DEV-008.a`](../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability)) | Template (s. o.); `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` (`failed: 5`, `exit11Next: 5`); `TestE2E_SandboxDevcontainer_EgressUpstreamDNS/abort_branches` (`ABORT … exit=11 result=failed` unter `generateEgressProject(…, "warn")` für drei Fälle); `…StatusFiles_Egress` | ok |
| Fall 1 (Benutzer fehlt): vor den Regeln, erster Lauf ohne Tabelle, `resolv.conf` unverändert, Meldung auf Standardfehlerausgabe | `abort_branches`: `STATE user-missing tables=0 drop=0 resolver-in-resolv.conf=0`, `ERR user-missing: …` | ok |
| Fall 2 (`nft`-Eingabe nicht angenommen), Wiederholung: Kette `output` der `ip`-Tabelle, Resolver, `resolv.conf`-Eintrag unverändert; Konfigurationsdatei bereits neu | `…EgressRerun/failed_transaction` (`chain-unchanged=1 same-resolver=1 resolv=1`); Reihenfolge `} > "$conf"` vor `nft -f` im Golden Case `…RerunIsOneTransaction` | ok; gepinnt nur für eine ungültige Regelzeile der `ip`-Tabelle, so formuliert |
| Fall 3 und 4: Tabellen stehen, Kette `drop`, `resolv.conf` unverändert | `abort_branches`: `STATE not-started` / `STATE uid-differs … tables=2 drop=1 resolver-in-resolv.conf=0` | ok mit V-3, V-6 |
| Fall 4 „oder keine auslesbare Pid-Datei“ | nur Template Zeile 126 (`running_uid` leer); kein Mutant, kein Golden-Text für den Leerfall | **V-3** |
| Upstream-DNS: `user=`/UID aus einem Wert; jede Port-53-Regel trägt `meta skuid`; `DIRECT … =none` UDP/TCP; `LIMIT … =answer`; erlaubter Name löst auf, Verbindung kommt zustande | `…UpstreamDNSOnlyForTheResolverUser` (Zählung genau 1, jede Zeile mit `dport 53 accept` trägt `meta skuid`, beide Richtlinien); `…EgressUpstreamDNS/direct_path_closed` (`DIRECT`, `LIMIT`, `LOCAL`, `HTTPS api.github.com=200`, `UID dnsmasq` gleich `rule`, `RULES skuid` gleich 2 je Upstream) | ok; Test gelesen: er prüft die Bedingung, das `LIMIT`-Kontrollexperiment gibt dem `none` Bedeutung |
| Wiederholungslauf: eine `nft -f`-Eingabe, `add`/`delete`/`add` je Tabelle, kein anderes `nft add`/`delete`, Reihenfolge, begrenztes Warten, `resolv.conf` zuletzt | `…RerunIsOneTransaction` (Zählung, Sequenz je Tabelle, `order`-Marker, `[ "$waited" -lt 50 ]`) | ok |
| Lauf ohne Abbruch endet mit Exit 0, `ok`, einem Resolver unter der Regel-UID, beiden Tabellen; Set danach leer und füllt sich neu | `…EgressRerun/repeated_runs` (elf Prüfzeilen per Regex; `SET before=1 after=0`; `HTTPS api.github.com=200` nach dem zwölften Lauf) | ok |
| „bestehende Verbindungen laufen über `established,related` weiter“ | Template Zeile 85; **kein Test** (Suche über `internal/**/*_test.go` nach `established,related`: kein Treffer; die Probe misst keine bestehende Verbindung) | **V-1** |
| „Upstream-Regeln betreffen IPv4-Resolver; Set `ipv4_addr`“ | Template Zeilen 49 (`$2 !~ /:/`), 82; Treffer in Tests: nur der Filter der **Probe** (`UPS=…`), nicht das Script; kein Golden-Text für `ipv4_addr` | **V-2** |
| `ip6`-Tabelle lässt nur Loopback zu | Golden Case pinnt `add chain ip6 … policy drop` und `oifname lo accept` genau einmal; Probe zählt `ip6drop=1 ip6lo=1` je Lauf | ok |
| `doctor`: `failed` ist Fehler; Meldung nennt vier Ursachen; Hinweis nennt Capability, Benutzer, Regeln und den Log-Verweis | `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` (vier Teile der Meldung, vier Teile des Hinweises, die genau diese Aussage tragen); `…StatusFiles_Egress` (`SeverityError`) | ok (F-9 behoben: die Spec behauptet für den Hinweis nur noch, was der Test prüft) |
| Kette `egress && docker`: nur der Text gepinnt | `TestGenerateDevcontainer_SandboxEgressWithDocker` vergleicht den String des `postStartCommand`; die Spec sagt „Gepinnt ist der Text der Verkettung“ und die Folge nur als Semantik von `&&` | inhaltlich ehrlich; Form siehe V-4 |
| Degradationstabellen-Zeile und Satz in [`LH-FA-DEV-008.a`](../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) | verweisen auf [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf); Tabellenzeile zählt die Ursachen ein zweites Mal auf (Wartungskopplung, INFO aus dem Review F-10) | ok |

## 3. Review-Findings F-1 bis F-12 gegen den Ist-Stand (nicht dem Report geglaubt)

| Finding | Stand | Beleg (selbst gelesen) |
|---|---|---|
| F-1 (MEDIUM) Handbuch stellt Verkettungsfolge als Tatsache dar | behoben | Handbuch, Absatz „Die Datei beschreibt den letzten Start“: Klammer zur Folge von Exit 11 gestrichen (`git diff`) |
| F-2 (LOW) CHANGELOG ohne „Resolver startet nicht“ | behoben | CHANGELOG `[Unreleased]`: „…does not exist in the image, the resolver does not start, or the resolver does not run under that user“ |
| F-3 (LOW) Klausel zur Startkette in *Nicht ermittelbar* | behoben | Klausel gestrichen (`git diff`) |
| F-4 (LOW) Teststand-Satz im Zielbild („…ein Lauf der Kette mit Abbruch ist es nicht“) | teilweise: der Teil „ein Lauf … ist es nicht“ ist weg, `Gepinnt ist der Text der Verkettung im postStartCommand.` steht weiter in *Stand* von [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) | **V-4**; der Slice nennt keinen Grund für das Belassen |
| F-5 (LOW) „leert die Menge der erreichbaren Adressen“ | Wortlaut berichtigt („das Set der erlaubten Adressen ist leer und füllt sich neu“), aber um eine neue, ungepinnte Aussage ergänzt | **V-1** |
| F-6 (LOW) Allquantor „Jeder Lauf“ | behoben | „Ein Lauf, der keinen der Abbrüche trifft (mit der Capability, einem startenden Resolver und einem erreichbaren Upstream-Resolver)“ |
| F-7 (LOW) Fall 2 in der Allgemeinheit | behoben | „gepinnt für eine ungültige Regelzeile der `ip`-Tabelle“; Konfigurationsdatei „bereits neu geschrieben“ genannt |
| F-8 (LOW) „keine Namensauflösung, aber keine offene Verbindung“; IPv4-Grenze | Spec behoben (Satz gestrichen, IPv4-Grenze genannt); das Handbuch trägt den Satz weiter | **V-5**, **V-2** |
| F-9 (LOW) Hinweis nennt nicht vier Ursachen | behoben | Spec-Wortlaut jetzt „Hinweis nennt die Capability, den Benutzer und die Regeln“, deckt der Test |
| F-10 bis F-11 (INFO) | F-11 teilweise berücksichtigt („etwa nach einem Wechsel auf diesen Benutzer, der root-Rechte voraussetzt“, „oder keine auslesbare Pid-Datei“); „gilt nicht diese Zeile“ unverändert, folgenlos | ok |
| F-12 (INFO) Zustandsmarker, Historie, „elf Läufe“ | behoben | Kopf nennt [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) als geschrieben; Historie trägt fünf Zeilen für fünf Änderungen; Fundstellen-Tabelle: „elf Läufe mit Prüfzeile; ein zwölfter“ |

## 4. DoD-Abgleich

- [x] **Drei überholte Stellen nachgeführt:** erfüllt. [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (`result`, *Stand*, *Nicht ermittelbar*), Degradationstabelle von [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und [`LH-FA-DEV-008.a`](../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) nennen Exit 11 und `failed` unter `warn`; Fundstellen-Tabelle in §3 vorhanden und stichprobenweise (alle Zeilenangaben, alle benannten Tests) bestätigt. Vorbehalt: V-1 bis V-3 (Bedingungen ohne Test, `AGENTS.md` §3.10 letzter Punkt).
- [x] **Frage der Spezifikationsstelle beantwortet:** eigener Eintrag [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) ([`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5), Begründung in §3 des Slice; Grenze (`sudo`, Adressen erlaubter Namen, Unternamen) steht als Grenze, nicht als Messstand; Historie trägt je Änderung eine Zeile; `make doc-trace` 0 Waisen, `make docs-check` 0 Befunde. Vorbehalt wie oben.
- [x] **Ein Vertrag, vier Stellen:** **nicht belegt, solange V-5 offen ist.** Lastenheft (Zeilen 467 und 469 bis 478 gelesen) schweigt zum Abbruch aus anderem Grund und widerspricht nicht; CHANGELOG und Handbuch nennen die drei Abbruchgründe unter `warn` wie die Spezifikation (Tabelle Z. 1078, `failed`-Zeile Z. 1130, Absätze Z. 1258 und 2185); das Handbuch sagt zu Fall 3 und 4 aber „In beiden Fällen gibt es keine Namensauflösung, aber auch keine offene Verbindung“, was die Spezifikation nach dem Review gestrichen hat und kein Test pinnt.
- [ ] `make gates` grün: Sensor erfüllt (Abschnitt 1); Haken setzt der Implementer in der Closure.
- [ ] Review durchgeführt: Report liegt vor; Haken in der Closure.
- [ ] Doku-Update, Closure-Notiz, Register, §6-Ausgänge, Paarungen: offen (Abschnitt 7).

## 5. Konformität mit den Entscheidungen

| Entscheidung | Befund |
|---|---|
| [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 (Abbruch auch unter `warn`: Benutzer fehlt, Resolver unter anderer UID) | konform; die Spezifikation führt zusätzlich die zwei Abbruchfälle, die [`ADR-0027`](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) einführte (Transaktion, Resolver startet nicht); beide sind im Template und in Tests belegt. |
| [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5 (Spezifikationsfrage; Lastenheft nicht abgeschwächt) | konform; `spec/lastenheft.md` im Diff leer, Grenze der Zusage steht als Grenze. |
| [`ADR-0027`](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 3 (Zusage nicht weiter als Beleg; „Fenster von null“ und unterbrechungsfreie Auflösung nicht zugesagt) | im Kern konform („Zugesagt ist diese Struktur; nicht zugesagt sind …“); verletzt in den Einzelaussagen V-1 und V-2 (Wortlaut stärker als der Beleg) und im Handbuch (V-5). |
| [`ADR-0027`](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 4 (fehlgeschlagene Transaktion: Exit 11, `failed`, auch `warn`) | konform; Fall 2 mit dem gepinnten Fehlerfall und ausgewiesener Klasse. |
| Zielbild-Regel `AGENTS.md` §3.4 | konform für die hinzugefügten Zeilen: Suche in [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) und im Diff nach `ADR-`, `slice`, `welle`, `roadmap`, `carveout`, Zeiten und „gemessen“: kein Treffer; die neuen Verweise zeigen aufwärts oder lateral. Einzige Form-Abweichung: V-4. |
| `AGENTS.md` §5 Regel 6 (Handbuch ohne Kennungen) | konform; Suche nach `LH-`, `ADR-`, `SPEC-`, `slice-`, `welle-` in `docs/user/benutzerhandbuch.md`: kein Treffer. |
| Accepted-ADRs unverändert, kein Gate gelockert, kein `//nolint` | konform (Diff berührt keine ADR, kein Gate, keinen Code). |

## 6. Plan-Tabelle gegen `git diff --stat 058bbf8..HEAD`

Fünf Dateien im Diff: `spec/spezifikation.md`, `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, die Slice-Datei, der Review-Report. Jede steht im Ist-Stand-Abschnitt von §3 des Slice; „unverändert“-Zeilen (`internal/**`, `Makefile`, `repo.mk`, Workflows, Lastenheft, Sicht) bestätigt durch leeren `git diff`, die Lesestellen für Lastenheft (Z. 467, 469 bis 478) und Sicht („egress“-Treffer) stimmen mit meiner Lesung überein. Roadmap und `open/README.md` stehen im Plan als „bei Start und Closure“ und sind im Range nicht berührt (Roadmap-Zeile „In Arbeit“ aus `058bbf8`, Stand korrekt für die Phase). Kein Eintrag ohne Diff, keine Diff-Zeile ohne Eintrag.

## 7. Befunde

| ID | Schwere | Befund | Fundstelle |
|---|---|---|---|
| V-1 | MEDIUM | [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (*Wiederholungslauf*) sagt „bestehende Verbindungen laufen über `established,related` weiter“. Die Aussage ist in der Fix-Runde zu F-5 entstanden, hat nur das Template (Zeile 85) als Beleg und keinen Test: kein Treffer für `established,related` in `internal/**/*_test.go`, die Rerun-Probe misst keine bestehende Verbindung. Das Review F-5 verlangte, die Aussage nicht weiter zu fassen als der Beleg; die Korrektur hat sie durch eine ungepinnte Zusage ersetzt. `AGENTS.md` §3.10 letzter Punkt: Test oder streichen. | `spec/spezifikation.md` · „bestehende Verbindungen laufen über `established,related` weiter“ |
| V-2 | LOW | „Die Upstream-Regeln betreffen IPv4-Resolver; … das Set der erlaubten Adressen ist `ipv4_addr`“: Das Template trägt es (Zeilen 49, 82), die Tests nicht (der Treffer `!~ /:/` im E2E-Test ist der Filter der Probe, nicht das Script). Als Grenze formuliert und plausibel, aber eine Bedingung ohne Test. | `spec/spezifikation.md` · „betreffen IPv4-Resolver“ |
| V-3 | LOW | Fall 4 „oder keine auslesbare Pid-Datei“ ist nur durch Template Zeile 126 belegt; `abort_branches` erzeugt die UID-Abweichung (`dns_uid=1234`), nicht den Leerfall. | `spec/spezifikation.md` · „oder keine auslesbare Pid-Datei“ |
| V-4 | LOW | Review F-4 empfahl, den Teststand-Satz aus dem Zielbild zu nehmen. Geblieben ist „Gepinnt ist der Text der Verkettung im `postStartCommand`.“ in *Stand* von [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung): eine Aussage über den Teststand, wie sie die Bereinigung vom 2026-10-19 (Historie) aus [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)/045/046/047 entfernt hat. Weder Slice noch Fix-Commit nennen einen Grund für das Belassen. | `spec/spezifikation.md` · [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand* |
| V-5 | MEDIUM | DoD „Ein Vertrag, vier Stellen“ ist `[x]`, aber das Handbuch sagt zu Fall 3 und 4 weiterhin „In beiden Fällen gibt es keine Namensauflösung, aber auch keine offene Verbindung“. Die Spezifikation hat denselben Satz nach F-8 gestrichen, weil kein Test ihn pinnt (`abort_branches` zählt Tabellen, Kette und `resolv.conf`; keine Auflösung, kein Verbindungsversuch). Zwei Stärken derselben Aussage, in der Nutzer-Doku die ungeprüfte (`AGENTS.md` §3.10: auch Handbuch-Bedingungen). Dazu das Handbuch: „(ein Docker-Test prüft sie)“ für die drei Zustände, wo der Test nur die dort genannten Zähler prüft. | `docs/user/benutzerhandbuch.md` · Abschnitt Egress-Beschränkung, Spiegelpunkt „DNS-Server startet nicht oder läuft unter einem anderen Benutzer“ |
| V-6 | INFO | „Bei 3 und 4 stehen **im ersten Lauf** beide Tabellen …“: In der Abbruch-Probe läuft `uid-differs` als dritter Lauf im selben Container nach `not-started` (Tabellen bestehen schon); der gemessene Zustand (`tables=2 drop=1 resolver-in-resolv.conf=0`) gilt, die Eigenschaft „erster Lauf“ ist nur für `not-started` frisch. Inhaltlich folgenlos, weil `resolv.conf` in beiden Fällen nie geschrieben wurde. | `spec/spezifikation.md` · [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) *Abbruch*; `egressAbortProbeScript` |
| V-7 | INFO | Zustandsmarker für die Closure: Der Slice-Kopf nennt „Verantwortlich: pt9912 (Rolle Implementer, Start 2026-10-20)“ (stimmt); die Roadmap-Zeile „In Arbeit“ und die Zeile des Slice in `open/README.md` müssen bei der Closure auf den Ist-Stand gesetzt oder entfernt werden (`AGENTS.md` §3.10, Zustandsmarker mit Ablaufpunkt). | `docs/plan/planning/in-progress/roadmap.md` Z. 20; `docs/plan/planning/open/README.md` Z. 42 |

## 8. Auflagen vor Closure

1. V-1: Satz zu `established,related` streichen, oder einen Test schreiben, der eine bestehende Verbindung über den Lauf hält; V-2 und V-3 ebenso (streichen oder pinnen; ein Golden-Text für `ipv4_addr` und `$2 !~ /:/` genügt für V-2; für V-3 ein Mutant mit leerer Pid-Datei, oder die Klausel streichen).
2. V-5: Handbuch auf die Stärke der Spezifikation bringen (Satz „keine Namensauflösung, aber auch keine offene Verbindung“ streichen oder durch einen Test pinnen) und „(ein Docker-Test prüft sie)“ auf das tatsächlich Geprüfte begrenzen; danach die DoD-Zeile „Ein Vertrag, vier Stellen“ erneut bestätigen. Nach einer Handbuch- oder Spec-Änderung die Plan-Tabelle (§3 des Slice) erneut gegen `git diff --stat` abgleichen (`AGENTS.md` §3.10, Stand nach der letzten Fix-Runde).
3. V-4: [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand*: den Satz „Gepinnt ist der Text der Verkettung im `postStartCommand`.“ streichen, oder den Grund des Belassens im Slice (§7) festhalten.
4. §6 des Slice: jedes der drei Risiken braucht einen Ausgang (eingetreten / entfallen / weiter offen) mit Beleg; das erste („eine Spec-Bedingung ohne pinnenden Test“) ist durch V-1 bis V-3 **eingetreten**, solange diese offen sind; das zweite (Verkettung nur als Text gepinnt) ist in [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) begrenzt formuliert und kann als „entfallen: Spezifikation sagt nur den gepinnten Text“ geschlossen werden (nach V-4).
5. §7 des Slice: Closure-Notiz mit Steering-Loop-Lerneintrag schreiben; das Register (`../observations/`) fortschreiben oder „keine Beobachtung angefallen“ belegen. Kandidaten: V-1 (Korrektur erzeugt eine neue ungepinnte Bedingung) und V-5 gehören in die Klasse `BEO-DOK/spec-bedingung-ohne-code-entsprechung` beziehungsweise `BEO-DOK/zusage-weiter-als-messung`; die Zählung gegen deren `state.md` prüfen, nicht dem Zählstand dieses Reports glauben.
6. Paarungen (Anker, Folge-Slice, Register) in der Closure prüfen; Roadmap-Zeile „In Arbeit“ und Zeile in `open/README.md` entfernen oder auf den Ist-Stand setzen (V-7).
7. Closure-Reihenfolge nach `AGENTS.md` §3.3 Fall 2: erst Inhalt (DoD-Haken, Closure-Notiz mit Commit-Hash), dann der reine `git mv` nach `done/`; die Haken „`make gates` grün“, „Review“, „Doku-Update“ erst nach 1 bis 5.

## 9. Nicht geprüft

- `make test-docker` (Begründung in Abschnitt 1): die Docker-Tests `…EgressUpstreamDNS`, `…EgressRerun`, `…StatusFiles_Egress` nicht neu gefahren; ihre Wortlaute sind gelesen, ihre Ergebnisse stammen aus dem Verification-Report des Code-Slice.
- Mutationen (kein Produktcode; V-1 bis V-3 sind deshalb als „kein Test vorhanden“ belegt, nicht durch eine fehlgeschlagene Mutation).
- Verhalten auf Hosts mit Kernel-Ablehnung einer `nft`-Transaktion, mit IPv6-Upstream-Resolvern, `arm64`, Docker Desktop, CI, `nestedRuntime: docker`: nicht erzeugbar auf diesem Host, in der Spezifikation nicht zugesagt.
- Handbuch und CHANGELOG außerhalb der Absätze zu Egress-Abbruch, Upstream-DNS und Statusdateien: Stichprobe, nicht vollständig gegen die Spezifikation gelesen.
- Beobachtungs-Register (`observations/`) und die drei Paarungen: nicht gelesen; Closure-Aufgabe.

## 10. Carveouts und Artefakt

Carveouts: neu none, gelöst none, unverändert (Inventar `carveouts.md` im Diff leer).

Verifizierter Stand: `8bb5d59`. Dieser Report liegt unter `docs/reviews/` (Verification-Report des Slice).
