# Slice v1: Egress — Wiederholungslauf von `egress-init.sh` ohne Lücke (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nach dem Release-Cut `v0.12.0` (Empfehlung des Verifiers, Abschnitt 5 seines Berichts; die Entscheidung über die Reihenfolge liegt beim Projektinhaber). Der Slice hält den Cut nicht auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht aus dem Befund B-1 von [`slice-v1-devcontainer-ollama-amd64-messung`](../in-progress/slice-v1-devcontainer-ollama-amd64-messung.md) hervor.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; eine Leitplanke, keine Sandbox-Grenze), [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Messung M2 mit atomarer Variante, Entscheidungen 1 bis 3, Konsequenz „das Startscript ist bei einem zweiten Lauf nicht atomar“), [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set).
Anlass: Befund B-1 der `amd64`-Messung: Der Wiederholungslauf von `egress-init.sh` löscht die Tabellen und baut sie danach neu; dazwischen hielt die Restriktion einen Verbindungsversuch in 27 von 30 Läufen nicht auf (Spanne bis 48 ms; zweiter Lauf 28 von 30, bis 50 ms), und die Namensauflösung des Containers fiel in 30 und 27 Läufen kurz aus. Messdaten: [`rohdaten-08.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-08.txt) und [`rohdaten-lauf2-08.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-lauf2-08.txt), Protokoll [`README.md`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md); `arm64`: [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) Messung M2 (Fenster 24, 30, 30, 34 und 45 ms). Das Ergebnis ist kein neuer Befund: `amd64` bestätigt M2.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Startscripte, Statusdateien) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Allowlist); der Architect liest beide und benennt, ob ein Satz über den zweiten Lauf des Startscripts dort steht (in diesem Plan nicht gelesen; `grep` nach „zweiten Lauf“ und „idempot“ in `spec/spezifikation.md` fand dazu keinen Treffer). Der Verweis zeigt aufwärts: Die Spezifikation nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein erneuter Lauf von `egress-init.sh` in einem Container mit aktiver Restriktion lässt keinen Zeitraum ohne Restriktion und keinen Ausfall der Namensauflösung eines erlaubten Namens entstehen — oder die Messung zeigt, dass und warum das mit dem Mechanismus von [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) nicht geht, und die Grenze steht ausgewiesen im Handbuch.

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Neuaufbau im Template:** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`, Zeilen 40 bis 43: der alte `dnsmasq` bekommt `kill`, danach `nft delete table ip u_boot_egress` und `nft delete table ip6 u_boot_egress6`; Zeilen 57 bis 72: `nft add table`, `set`, `chain` und jede `rule` einzeln, danach `dnsmasq --conf-file` und `echo "nameserver 127.0.0.1" > /etc/resolv.conf`. Zwischen dem `delete` (Zeile 42) und den ersten Regeln (Zeilen 57 bis 62) gilt keine Tabelle, zwischen dem `kill` (Zeile 41) und dem Start (Zeile 71) läuft kein lokaler Resolver. Die Kommentarzeile 40 sagt „Re-run safe“.
- **Gemessen:** [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) M2 (`arm64`, Colima): Fenster ohne Filter-Regelwerk 24 bis 45 ms, Ausfall des lokalen `dnsmasq` 28 bis 55 ms; die **atomare Variante** (Messskript, kein Produktcode: Wechsel in einer nft-Transaktion, `dnsmasq` erst nach dem Ende des alten) zeigte in fünf Läufen keinen Zustand `offen` und einen Ausfall von höchstens 24 ms. `amd64` (Messfrage A7, 30 Läufe je Messung): Befund B-1 oben. Die Sonde ist ein Verbindungsversuch zum Gateway des Containers auf Port 9, abgelehnt statt verworfen; die Zählung ist eine Untergrenze.
- **Wirkung:** Der Wiederholungslauf braucht `sudo` (Rootrechte für `nft`), mit dem sich die Restriktion nach [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) Entscheidung 2 ohnehin aufheben lässt. Die Sicherheitsrelevanz ist begrenzt: kein neuer Weg für einen Prozess ohne `sudo`, ein Fenster nur für einen gleichzeitigen Prozess, während jemand mit Rootrechten das Skript erneut ausführt. Nicht gemessen und nicht gesagt: ob irgendein normaler Betriebspfad das Skript im laufenden Container erneut ausführt (`postStartCommand` führt es beim Start aus; der Docker-Test mit `egressSubdomainProbeScript` führt es nach dem Umschreiben der Liste im laufenden Container erneut aus, `internal/e2e/sandbox_devcontainer_docker_test.go`, Zeilen 180 bis 185).
- **Was das Handbuch sagt:** Der Abschnitt *Internet-Beschränkung (Egress)* nennt die Bedingung („in einem laufenden Container … erneut ausführen“), beide Zählungen, die Metrik, „Untergrenze“ und „Leitplanke und keine Sicherheitsgrenze“ (`docs/user/benutzerhandbuch.md`, Zeilen 1355 bis 1366).

### Vorbehalt: Entscheidung des Projektinhabers vor dem Start

[`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) ist `Accepted` und bleibt unverändert (`AGENTS.md` §3.5). Entscheidung 1 schließt einen Befehl aus, der ein Programm im Container ausführt oder eine Liste hineinlegt (Neu-Laden), und Entscheidung 3 hält R1 bis R4 und S offen; Folge-Slices zu R1 und R2 entstehen nicht. Ein atomares Neuschreiben der **eigenen** Tabelle im Startscript ist kein Neu-Laden-Weg, nähert aber den Entscheidungsraum von R1 an (die atomare Variante ist dort als Zelle der Option R1 geführt). Ob die Änderung unter den bestehenden Entscheidungen trägt oder eine **Folge-ADR** braucht (die [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) nicht umschreibt, sondern ergänzt: `Supersedes`/`Schärft`), entscheidet der Projektinhaber vor dem Start; bis dahin ist der Slice **Rückführung/Entscheidung** (`open/`). Der Planner nimmt diese Entscheidung nicht vorweg.

### Messfragen (Arbeitsplan; der Architect schneidet)

| # | Frage | Warum |
|---|---|---|
| E1 | Trägt eine Variante (neue Tabelle in einer einzigen `nft -f`-Transaktion, die alte erst danach ersetzt; `dnsmasq` erst beendet, wenn der neue bereit ist; `resolv.conf` zuletzt) auf `amd64` und `arm64` ein Fenster von null (A7-Messung mit derselben Sonde, 30 Läufe) und keinen Ausfall der Auflösung eines erlaubten Namens? | Die Messung belegt die Variante auf dem Zielpfad, nicht nur auf `arm64` (M2) |
| E2 | Bleiben die Start-Eigenschaften (Exit-Codes `0` und `11`, Statusdatei `ok`/`fallback`/`failed`, IPv6-Regeln, Neustart des Containers) gleich? | keine Regression im Erststart |
| E3 | IPv6: Die Tabelle `ip6 u_boot_egress6` wird beim Wiederholungslauf ebenfalls gelöscht und neu angelegt; Messung des Fensters für IPv6 (in B-1 nicht gemessen) | die Aussage im Handbuch führt IPv6 als nicht gemessen |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Neu-Laden der Liste im laufenden Container vom Host aus (R1, R2)** — entschieden gegen (Entscheidung 1 von [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)); der Slice ändert den Wiederholungslauf des Startscripts, nicht den Weg „Liste ändern“.
- **Der offene Pfad zum Upstream-Resolver ohne `sudo`** — Gegenstand von [`slice-v1-egress-dns-upstream-schliessen`](slice-v1-egress-dns-upstream-schliessen.md).
- **Entziehen von `sudo`, Restriktion außerhalb des Containers (S, R4)** — eigene Vorgänge, in [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) offen gehalten.
- **Der Verkehr verschachtelter Container (B-2 der Messung)** — anderer Mechanismus (Beobachtung [`BEO-APP/nested-docker-umgeht-egress-podman-nicht`](../observations/BEO-APP/nested-docker-umgeht-egress-podman-nicht/observation.md)).
- **Änderung des Lastenhefts** — nicht erwartet: Die Restriktion bleibt eine Leitplanke ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).

## 2. Definition of Done

- [ ] **Entscheidung und Messung:** Der Projektinhaber hat vor dem Start entschieden, ob die Änderung unter [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) trägt oder eine Folge-ADR braucht (Trägerin ist dann diese ADR, `Accepted` vor dem Code); E1 bis E3 sind gemessen, mit Skripten und Rohdaten im Repo (Skripte unter `tools/messung/<name>/`, Rohdaten unter `docs/messungen/<slice>/`; das Protokoll [`slice-v1-devcontainer-ollama-amd64-messung`](../in-progress/slice-v1-devcontainer-ollama-amd64-messung.md) liefert Messfrage A7 als Vorlage).
- [ ] **Template, Golden Case und Docker-Test:** `egress-init.sh.tmpl` baut die Tabellen ohne Zeitraum ohne Restriktion neu auf (Variante nach der Messung); Golden Case nach [`harness/replay.md`](../../../../harness/replay.md) (Fresh-State und zweiter Lauf des Startscripts); ein Docker-Test (`//go:build docker`) führt das Startscript wiederholt aus und weist ein Fenster ab; ob die Sonde deterministisch genug für ein Gate ist, ist ein Risiko (§6).
- [ ] **Handbuch und CHANGELOG** auf den Ist-Stand (Abschnitt *Internet-Beschränkung (Egress)*: die Zählungen und die Aussage „Untergrenze“ gelten für die alte Fassung; Messstand der neuen Fassung); ein Befund aus E1 bis E3, der nicht schließt, steht ausgewiesen dort.
- [ ] `make gates` grün; `make test-docker` grün (Code geändert); Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` | update | Neuaufbau ohne Lücke (Zeilen 40 bis 72), Variante nach E1 |
| `internal/hexagon/application/generate_sandbox_test.go` und gegebenenfalls ein Golden-Literal | update | Golden Case für den zweiten Lauf des Startscripts ([`harness/replay.md`](../../../../harness/replay.md)); Testnamen ergänzt der Architect |
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | Docker-Test für den Wiederholungslauf (Fenster, Auflösung eines erlaubten Namens) |
| `tools/messung/<name>/` und `docs/messungen/<slice>/` | neu | Skripte und Rohdaten für E1 bis E3 (Vorlage: A7 in `tools/messung/ollama-amd64/08-wiederholungslauf.sh` und `inner-08.sh`) |
| `docs/plan/adr/` (Folge-ADR) | neu, falls der Projektinhaber es entscheidet | Ergänzung zu [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md); `Accepted` bleibt unverändert |
| `spec/spezifikation.md` ([`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und vom Architect benannte Stellen) | update, falls die Spezifikation das Verhalten des zweiten Laufs zusagt | Spezifikations-Nachführung als eigener Schritt nach dem Code |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Ist-Stand und Messstand der neuen Fassung |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |
| `spec/lastenheft.md`, `.github/workflows/`, `Makefile`, `repo.mk` | unverändert | keine Zusage-Änderung, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; der Projektinhaber hat die Frage nach der Folge-ADR entschieden (Vorbehalt in §1); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Messung (E1 bis E3), Folge-ADR und Code tragen zusammen mehr als drei Liefer-Punkte; dann trennt der Architect die Messung mit ADR vom Code.
- `in-progress` → `open` (blockiert): Die Entscheidung zur Folge-ADR fehlt, oder die Messung zeigt, dass die atomare Variante die Start-Eigenschaften ändert (E2); Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Die Änderung ist unter den bestehenden Entscheidungen von [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) nicht gedeckt.** — **Ausgang:** weiter offen bis zur Entscheidung des Projektinhabers vor dem Start (Vorbehalt in §1); Folge-ADR als Träger.
- **Ein Docker-Test mit Sonde im Takt von Millisekunden ist nicht deterministisch** (A7 zählte eine Untergrenze). — **Ausgang:** weiter offen bis zum Architect-Schnitt; der Golden Case pinnt die Struktur des Scripts, der Docker-Test misst mit Wiederholungen und fällt nur bei einem Treffer.
- **Die atomare Variante ändert das Verhalten beim Erststart oder bei IPv6.** — **Ausgang:** weiter offen bis E2/E3.
- **Die Messung auf dem geteilten Host des Projektinhabers verändert dessen Docker-Bestand.** — **Ausgang:** weiter offen; Präfix-Regel und Aufräum-Schritt des Vorgängers übernehmen (nur Objekte mit dem Präfix der Messung).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (Template, APP), `internal/e2e` (E2E), `tools/` (WZG) und `docs/` (DOK) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-APP/nested-docker-pull-timeout-nach-egress-start`](../observations/BEO-APP/nested-docker-pull-timeout-nach-egress-start/observation.md) (offen, 1×; berührt die Reihenfolge von Beschränkung und Daemon-Start), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (3×, Ausgang beim Projektinhaber; die Zählungen im Handbuch werden je Lauf gegen die Rohdaten gelesen), [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Schwelle erreicht; Messskripte nutzen nur `docker`, `make` und die Shell). Evidenz-/Diskrepanz-Risiko: mittel (Messung auf Untergrenze).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
