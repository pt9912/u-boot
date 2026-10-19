# Slice v1: Egress — Wiederholungslauf von `egress-init.sh` ohne Lücke (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nach dem Release-Cut `v0.12.0` (Empfehlung des Verifiers, Abschnitt 5 seines Berichts; die Entscheidung über die Reihenfolge liegt beim Projektinhaber). Der Slice hält den Cut nicht auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht aus dem Befund B-1 von [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) hervor.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; eine Leitplanke, keine Sandbox-Grenze), [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Messung M2 mit atomarer Variante, Entscheidungen 1 bis 3, Konsequenz „das Startscript ist bei einem zweiten Lauf nicht atomar“), [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set).
Anlass: Befund B-1 der `amd64`-Messung: Der Wiederholungslauf von `egress-init.sh` löscht die Tabellen und baut sie danach neu; dazwischen hielt die Restriktion einen Verbindungsversuch in 27 von 30 Läufen nicht auf (Spanne bis 48 ms; zweiter Lauf 28 von 30, bis 50 ms), und die Namensauflösung des Containers fiel in 30 und 27 Läufen kurz aus. Messdaten: [`rohdaten-08.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-08.txt) und [`rohdaten-lauf2-08.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-lauf2-08.txt), Protokoll [`README.md`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md); `arm64`: [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) Messung M2 (Fenster 24, 30, 30, 34 und 45 ms). Das Ergebnis ist kein neuer Befund: `amd64` bestätigt M2.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Startscripte, Statusdateien) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Allowlist); [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) (Kopf, Feld `Schärft`) hält fest, dass `spec/spezifikation.md` keinen Satz zum zweiten Lauf des Startscripts trägt (gelesen am 2026-10-19). Dieser Slice ändert die Spezifikation nicht; ob sie eine Zusage samt Grenze aufnimmt, entscheidet [`slice-v2-spec-upstream-dns-egress-init`](../open/slice-v2-spec-upstream-dns-egress-init.md). Der Verweis zeigt aufwärts: Die Spezifikation nennt diesen Slice nie.

**Verantwortlich:** pt9912 (Rolle Implementer)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein erneuter Lauf von `egress-init.sh` in einem Container mit aktiver Restriktion ersetzt beide Tabellen in einer einzigen `nft -f`-Transaktion und wechselt den Resolver geordnet ([`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 1 und 4). Zugesagt wird die Struktur, die ein Golden Case pinnt; ein Fenster von null und ein Ausfall der Namensauflösung von null sind **nicht** zugesagt, und was ungemessen ist, steht ausgewiesen im Handbuch (Entscheidung 3).

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Neuaufbau im Template:** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`, Abschnitt „Re-run safe“ (Kommentarzeile): der alte `dnsmasq` bekommt `kill`, danach `nft delete table ip u_boot_egress` und `nft delete table ip6 u_boot_egress6`; dahinter `nft add table`, `set`, `chain` und jede `rule` einzeln, danach der Start von `dnsmasq` (`dnsmasq --conf-file`, mit Prüfung von Rückgabewert und UID) und `echo "nameserver 127.0.0.1" > /etc/resolv.conf`. Zwischen den `nft delete table` und den ersten `nft add rule` gilt keine Tabelle, zwischen dem `kill` und dem Start von `dnsmasq` läuft kein lokaler Resolver. Der Kommentar sagt „Re-run safe“ (gelesen im Template nach der Umsetzung von [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md); Zeilenbezüge stehen nicht mehr, weil die Umsetzung das Template verlängert hat).
- **Gemessen:** [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) M2 (`arm64`, Colima): Fenster ohne Filter-Regelwerk 24 bis 45 ms, Ausfall des lokalen `dnsmasq` 28 bis 55 ms; die **atomare Variante** (Messskript, kein Produktcode: Wechsel in einer nft-Transaktion, `dnsmasq` erst nach dem Ende des alten) zeigte in fünf Läufen keinen Zustand `offen` und einen Ausfall von höchstens 24 ms. `amd64` (Messfrage A7, 30 Läufe je Messung): Befund B-1 oben. Die Sonde ist ein Verbindungsversuch zum Gateway des Containers auf Port 9, abgelehnt statt verworfen; die Zählung ist eine Untergrenze.
- **Wirkung:** Der Wiederholungslauf braucht `sudo` (Rootrechte für `nft`), mit dem sich die Restriktion nach [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) Entscheidung 2 ohnehin aufheben lässt. Die Sicherheitsrelevanz ist begrenzt: kein neuer Weg für einen Prozess ohne `sudo`, ein Fenster nur für einen gleichzeitigen Prozess, während jemand mit Rootrechten das Skript erneut ausführt. Nicht gemessen und nicht gesagt: ob irgendein normaler Betriebspfad das Skript im laufenden Container erneut ausführt (`postStartCommand` führt es beim Start aus; der Docker-Test mit `egressSubdomainProbeScript` führt es nach dem Umschreiben der Liste im laufenden Container erneut aus, `internal/e2e/sandbox_devcontainer_docker_test.go`, Zeilen 180 bis 185).
- **Was das Handbuch sagt:** Der Abschnitt *Internet-Beschränkung (Egress)* nennt die Bedingung („in einem laufenden Container … erneut ausführen“), beide Zählungen, die Metrik, „Untergrenze“ und „Leitplanke und keine Sicherheitsgrenze“ (`docs/user/benutzerhandbuch.md`, Zeilen 1355 bis 1366).

### Entscheidung und Messung (Stand der Umsetzung)

Der Projektinhaber entschied am 2026-10-19: Die Änderung trägt eine Folge-ADR, [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) (`Accepted`), die [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) ergänzt und nicht umschreibt. **Die Messung E1 bis E3 entfällt** (abgebrochen; nichts davon liegt im Repo): kein Messskript, kein Eintrag unter `docs/messungen/`. Belege sind M2 (`arm64`, atomare Variante als Messskript) und B-1 (`amd64`, heutige Fassung). Daraus folgt der Zuschnitt der Zusage in [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 3: Struktur zugesagt; `amd64`-Wirkung der neuen Fassung, IPv6, Last, `nestedRuntime` und die Zeiten des Produkt-Scripts ungemessen und ausgewiesen. Die Fragen E1 bis E3 des Arbeitsplans sind damit nicht beantwortet, sondern entfallen; E2 (Exit-Codes, Statusdatei, Neustart) trägt der Docker-Test als Regression, E1 und E3 bleiben ungemessen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Neu-Laden der Liste im laufenden Container vom Host aus (R1, R2)** — entschieden gegen (Entscheidung 1 von [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)); der Slice ändert den Wiederholungslauf des Startscripts, nicht den Weg „Liste ändern“.
- **Der offene Pfad zum Upstream-Resolver ohne `sudo`** — Gegenstand von [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md).
- **Entziehen von `sudo`, Restriktion außerhalb des Containers (S, R4)** — eigene Vorgänge, in [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) offen gehalten.
- **Der Verkehr verschachtelter Container (B-2 der Messung)** — anderer Mechanismus (Beobachtung [`BEO-APP/nested-docker-umgeht-egress-podman-nicht`](../observations/BEO-APP/nested-docker-umgeht-egress-podman-nicht/observation.md)).
- **Änderung des Lastenhefts** — nicht erwartet: Die Restriktion bleibt eine Leitplanke ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).

## 2. Definition of Done

- [x] **Entscheidung:** Trägerin der Umsetzung ist [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md), `Accepted` vor dem Code (2026-10-19). Die Messung E1 bis E3 entfällt auf Entscheidung des Projektinhabers; es gibt kein Messskript und kein Messprotokoll, die Lücken stehen in [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 3 und im Handbuch.
- [ ] **Template, Golden Case und Docker-Test:** `egress-init.sh.tmpl` ersetzt beide Tabellen in einer `nft -f`-Transaktion und wechselt den Resolver nach [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 1 und 4; Golden Case nach [`harness/replay.md`](../../../../harness/replay.md) pinnt die Struktur (Fresh-State und zweiter Lauf des Startscripts); ein Docker-Test (`//go:build docker`) führt das Startscript wiederholt aus, weist ein Fenster nur als Treffer ab (Untergrenze) und prüft den Fehlerzweig der Transaktion (§6).
- [ ] **Handbuch und CHANGELOG** auf den Ist-Stand (Abschnitt *Internet-Beschränkung (Egress)*: die Zählungen und die Aussage „Untergrenze“ gelten für die alte Fassung; für die neue Fassung steht die Struktur und was ungemessen ist, [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 3).
- [ ] `make gates` grün; `make test-docker` grün (Code geändert); Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` | update | Neuaufbau ohne Lücke (Abschnitt von `kill` und `nft delete table` bis zum Schreiben von `resolv.conf`), Variante nach E1 |
| `internal/hexagon/application/generate_sandbox_test.go` und gegebenenfalls ein Golden-Literal | update | Golden Case für den zweiten Lauf des Startscripts ([`harness/replay.md`](../../../../harness/replay.md)); Testnamen ergänzt der Architect |
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | Docker-Test für den Wiederholungslauf (Fenster, Auflösung eines erlaubten Namens) |
| `tools/messung/<name>/` und `docs/messungen/<slice>/` | entfällt | Die Messung E1 bis E3 entfällt (Entscheidung des Projektinhabers, 2026-10-19); es entsteht nichts |
| `docs/plan/adr/0027-egress-init-wiederholungslauf-atomar.md` und `docs/plan/adr/README.md` | neu (vor der Umsetzung, `Accepted`) | Trägerin der Umsetzung; ergänzt [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md), das `Accepted` bleibt unverändert |
| `spec/spezifikation.md` | unverändert | trägt keinen Satz zum zweiten Lauf ([`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md), Feld `Schärft`); eine Nachführung entscheidet [`slice-v2-spec-upstream-dns-egress-init`](../open/slice-v2-spec-upstream-dns-egress-init.md) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Ist-Stand und Messstand der neuen Fassung |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |
| `spec/lastenheft.md`, `.github/workflows/`, `Makefile`, `repo.mk` | unverändert | keine Zusage-Änderung, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; der Projektinhaber hat die Frage nach der Folge-ADR entschieden ([`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md), `Accepted`); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Code, Tests und Handbuch tragen zusammen mehr als drei Liefer-Punkte; dann trennt der Architect den Code von der Doku.
- `in-progress` → `open` (blockiert): Ein Test zeigt, dass die Transaktion die Start-Eigenschaften ändert (Exit-Codes, Statusdatei, `skuid`-Regel) und die Abweichung von [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) nicht trägt; Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Die Änderung ist unter den bestehenden Entscheidungen von [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) nicht gedeckt.** — **Ausgang:** entschieden: [`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) (`Accepted`) ist die Trägerin.
- **Ein Docker-Test mit Sonde im Takt von Millisekunden ist nicht deterministisch** (A7 zählte eine Untergrenze). — **Ausgang:** der Golden Case pinnt die Struktur; der Docker-Test zählt Treffer einer Sonde (Verbindungsversuch zum Gateway, Port 9) während zehn Wiederholungsläufen und fällt nur bei einem Treffer; eine Nullzählung trägt keine Zusage ([`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 3). Antwortet das Gateway ohne Restriktion nicht, bleibt die Sonde auf dem Host ohne Aussage und der Test sagt es.
- **Die Transaktion ändert das Verhalten beim Erststart oder bei IPv6.** — **Ausgang:** E2 trägt der Docker-Test (Exit `0`, Status `ok`, beide Tabellen, `skuid`-Regel, ein Resolver unter der Regel-UID, erlaubter Name löst auf); die IPv6-Tabelle ist nur auf ihr Vorhandensein und ihre Loopback-Regel geprüft, ihr Fenster ist ungemessen und ausgewiesen. Neu: Die Tabellen werden jetzt auf ihr Ergebnis geprüft; ein Host, auf dem `nft` die Familie `ip6` nicht anlegen kann, endet mit Exit `11` ([`ADR-0027`](../../adr/0027-egress-init-wiederholungslauf-atomar.md) Entscheidung 4).
- **Die Messung auf dem geteilten Host des Projektinhabers verändert dessen Docker-Bestand.** — **Ausgang:** entfällt mit der Messung; die Docker-Tests legen nur Images und Container mit dem Präfix `uboot-e2e-` an und räumen sie auf.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (Template, APP), `internal/e2e` (E2E), `tools/` (WZG) und `docs/` (DOK) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-APP/nested-docker-pull-timeout-nach-egress-start`](../observations/BEO-APP/nested-docker-pull-timeout-nach-egress-start/observation.md) (offen, 1×; berührt die Reihenfolge von Beschränkung und Daemon-Start), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (3×, Ausgang beim Projektinhaber; die Zählungen im Handbuch werden je Lauf gegen die Rohdaten gelesen), [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Schwelle erreicht; Messskripte nutzen nur `docker`, `make` und die Shell). Evidenz-/Diskrepanz-Risiko: mittel (Messung auf Untergrenze).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
