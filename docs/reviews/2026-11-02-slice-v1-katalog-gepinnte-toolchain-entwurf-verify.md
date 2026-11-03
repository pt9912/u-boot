# Verification: slice-v1-katalog-gepinnte-toolchain-entwurf — 2026-11-02

**Rolle:** Verifier (`harness/roles.md`, frischer Kontext). **Frage:** „Bauen wir es richtig?" — die drei Liefer-Zeilen der DoD (§2) gegen die Belege.

**Gegenstand:** [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md), alle Berührungen der DoD-Zeilen samt Fix-Runde (`2d53e63d`, Report-Berichtigung `bd71715f`) am Stand HEAD `bd71715f`.

**Eingangs-Kontext:** Entwurf-Slice (Stand `bd71715f`), CR-Vorschlag [`cr-vorschlag-gepinnte-werkzeugkette`](../plan/planning/done/cr-vorschlag-gepinnte-werkzeugkette.md) (Freigabe 2026-11-02, `3cb04e77`, Klarstellung `6dd2a14d`, Fix `2d53e63d`), Messprotokoll [`docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md`](../messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md) (44 Rohdaten-Dateien, Skripte `tools/messung/katalog-gepinnte-toolchain/`, 5 Dateien), [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (`Accepted`, eigener Commit `56ce6c4b`), zwei Folge-Slices in `open/`, Review-Report [`2026-11-02-slice-v1-katalog-gepinnte-toolchain-entwurf.md`](2026-11-02-slice-v1-katalog-gepinnte-toolchain-entwurf.md) (0 HIGH/4 MEDIUM/1 LOW/2 INFO).

**Fundstellen am Code selbst gelesen** ([`AGENTS.md`](../../AGENTS.md) §3.10): `.devcontainer/Dockerfile` Zeilen 75 bis 120 — `ARG`-Block 87–92, `RUN`-Install-Block mit `sha256sum -c` 93–108, Cache-/ENV-Kommentar 109–112, `ENV`-Block 113–118, Prüfsummen-Kommentar 81–83 (`go.dev/dl`, „golangci-lint checksums file"), Hebung-Kommentar 84–85. Lastenheft-Version unverändert **0.29.0** (Fundstelle `spec/lastenheft.md` Tabelle „Version", Zeile 8); [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) als Grabstein unverändert (Zeilen 884–886, „*Entfallen (0.6.0).*").

---

## DoD-Abgleich (§2 des Entwurfs, drei Liefer-Zeilen)

- [x] **Messung zur ENV- und Cache-Mechanik** — erfüllt. Protokoll [`README.md` der Messung](../messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md) trägt **M1** (m1a/m1b/m1c/m1d, je zwei Läufe: `${localWorkspaceFolder}` löst zum Host-Pfad — in der Sandbox-Form nicht aufs Volume; `${containerWorkspaceFolder}` löst zum Container-Pfad, Fundstelle m1c; `remoteEnv` erreicht den `docker exec`-Weg des Adapters nicht, m1d), **M2** (m2a/m2b/m2c, je zwei Läufe: Block-`ENV` greift; `containerEnv` übersteuert das Block-`ENV` bei gleichem Namen), **M3** (Dockerfile-ENV-Form, zwei Läufe: Marker + 64-KiB-Pad vollständig, gleiches Volume, Eigentumsrechte `vscode:vscode`, Container-ID gewechselt) und **M3-Zusatz** (`m3z`-gocache/`m3z`-gomod, je zwei Läufe, dieselbe Prüfstellen-Tabelle wie M3 in der `containerEnv`-Form). Messstand datiert (2026-11-02, 07:59–08:10 UTC), `linux-host amd64`, lokale Engine 29.9.2, gepinnte CLI 0.89.0 — die DoD-Umgebung. Verworfene Vorläufe (M2-Reihe, M3-Vorlauf) sind im Protokoll vermerkt, nicht untergeschlagen. Beleg-Commits: `f76b7099` (M1–M3, Skripte, Protokoll, Bestand `docs/messungen/README.md`), `5e08fd50` (Zusatzlauf `m3z` zur §4-Frage-3-Entscheidung `b89943ea`). Skripte liegen laut Plan-Tabelle der Ablage-Regel nach `tools/messung/katalog-gepinnte-toolchain/` (Commit-Dateiliste bestätigt: `lib.sh`, `m1-…`, `m2-…`, `m3-…`, `m3z-…`).
- [x] **CR-Wortlaut und Träger** — erfüllt, **mit einem Restmarker (F-5-Rest, unten)**. [`cr-vorschlag-gepinnte-werkzeugkette`](../plan/planning/done/cr-vorschlag-gepinnte-werkzeugkette.md) trägt Alt/Neu je Stelle für **alle acht Stellen** (§§2.1–2.8) samt Regelrahmen (Decken-Regel, Grabstein-Regel — [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) zitiert, nicht belebt; neue Kennung [`LH-FA-DEV-011`](../../spec/lastenheft.md) als Entwurfs-Kennung mit `d-check:ignore`) und ADR-Gegenprüfung gegen das Katalog-Regelwerk (Stelle 1, Satz S6; [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) §Verhältnis). **Freigabe festgetragen:** Freigabe-Block oben im CR-Vorschlag (2026-11-02) — acht Stellen je wie vorgelegt, F1 (Kennung [`LH-FA-DEV-011`](../../spec/lastenheft.md) in §4.3), F2 (Abweisung `go` + `go-toolchain`, Exit-Code `10`), F3 (**`10` je Fehlerfall nach der Klarstellung „deine Empfehlung"** gegen die vorgelegte 13/10-Empfehlung) entschieden; F4–F6 ohne eigenen Aussatz bzw. als offener Punkt an den Umsetzungs-Slice (Stelle 8). Träger der Umsetzung: der CR-Slice [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) in `open/` (Lastenheft 0.29.0 → 0.30.0, wortgleich, CR-Disziplin: schreiben, nicht neu diskutieren). [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) ist ein eigener Commit des Entwurfs-Umfangs (`56ce6c4b`, `Accepted`, ADR-Index nachgezogen — Fundstelle `docs/plan/adr/README.md` Zeile 69).
- [x] **Folge-Slice-Plan in `open/` mit vollständiger Plan-Tabelle** — erfüllt. [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md): Plan-Tabelle §3 mit **13 Zeilen + Ansatzliste** — alle Punkte der DoD-Zeile 3 benannt (Template-Block Go + `golangci-lint`, Katalog-Regeln mit `version:`-Übersteuerung und den Abweisungen, `containerEnv`-Render-Stelle, Golden Cases + Regenerate der eigenen `.devcontainer/` ohne Hand-written-Teil, Docker-Test mit Pins + Cache-Recreate + Negativprobe der Prüfsumme, Pin-Stellen-Abgleichstest der vier Stellen, Hebung-Abgleich über Root-`Dockerfile`/`repo.mk`/Template-`ARG`), `Verantwortlich:` **pt9912** gesetzt, Start-Bedingung „nach dem CR-Slice in `done/`". Der CR-Slice trägt **seine** vollständige Plan-Tabelle (8 Zeilen, Fundstellen der acht Lastenheft-Stellen, [`CHANGELOG.md`](../../CHANGELOG.md) `[Unreleased]`, Handbuch ohne Kennungen). Bestand beider Slices in `open/README.md` (Zeilen 89–90) und Roadmap-Zeile `v0.19.0` (Nennung beider, ADR und Freigabe belegt) vorhanden.

**Plan-vs-Code (§3 des Entwurfs, Stand nach `2d53e63d`)** gegen `git diff --stat b160d4ac^..HEAD` (75 Dateien, +3251/−49): Der Range mischt **zwei** Slices. Attributionsgetrennt:

- **Entwurf-attributiert** (`2481485d`, `f76b7099`, `63e16263`, `b89943ea`, `5e08fd50`, `56ce6c4b`, `3cb04e77`, `dd3763a4`, `6dd2a14d`, `2d53e63d`, `bd71715f`, Moves `05cb1aa9`/`8353b5a1`/`4910bf9d`/`e04e3329`): jede Diff-Datei trägt eine §3-Zeile — Messverzeichnis + Skripte (Zeile 1), `docs/messungen/README.md` (Zeile 2), `spec/lastenheft.md` als ausgewiesener Nicht-Eintrag (Zeile 3, im Range nicht geändert — bestätigt), diese Datei/`open/README`/beide Folge-Slices (Zeile 4), [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) + ADR-Index (Zeile 5), CR-Vorschlag (Zeile 6), Roadmap (Zeile 7), `docs/reviews/` (Zeile 8). Kein Entwurf-attributierter Diff-Pfad ohne Plan-Zeile — **abgedeckt**.
- **Nicht Entwurf-attribuert** (Release-Cut [`slice-v1-release-cut-v0.18.0`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md), `e6c87379`, `9a7db2f1`, Moves/Doku-Nachzug): [`CHANGELOG.md`](../../CHANGELOG.md), beide READMEs, `cmd/uboot/main.go`, `docs/maintainer/releasing.md`, `repo.mk`, `Dockerfile`, `docs/user/benutzerhandbuch.md`, `harness/conventions/MR-003-roadmap-wellen.md`, zwei Review-/Verification-Reports, Closure-Zügel am Container-Name-Slice, Beobachtungs-Evidence — diese Dateien haben ihren Plan im (bereits geschlossenen) Cut-Slice, nicht hier. Lifecycle-mechanische Berührung des Entwurfs-Perimeters: die `slice-mv` <!-- d-check:ignore (Make-Target, keine Slice-Kennung) ---> Verweis-Nachzüge am JDK-Slice [`slice-v1-katalog-feature-jdk-gepinnt`](../plan/planning/done/slice-v1-katalog-feature-jdk-gepinnt.md) (Moves `4910bf9d`, `e04e3329`; Cut-Fix `9a7db2f1`) — mechanisch, keine Entwurfs-Aussage.
- `internal/hexagon/**`, `tests/**`, e2e sind — wie in §3 Zeile 9 versprochen — **ungeöffnet** (Diff bestätigt: kein Produktcode im Entwurf-Umfang).

---

## Review-Findings F-1 bis F-7 am Fix-Stand

Quelle: [`2026-11-02-slice-v1-katalog-gepinnte-toolchain-entwurf.md`](2026-11-02-slice-v1-katalog-gepinnte-toolchain-entwurf.md); Fix-Runde `2d53e63d`, Report-Berichtigung `bd71715f`. Jedes Finding am Diff `6dd2a14d..2d53e63d` und am Ist-Stand geprüft.

| ID | Fix-Stand-Prüfung | Ergebnis |
|---|---|---|
| **F-1** (MEDIUM, veralteter Exit-Code in Freigabe-Notiz und Satz-Beleg) | Diff `2d53e63d` am CR-Vorschlag: Freigabe-Block jetzt „(F1 und F2 wie empfohlen; F3 nach Klarstellung mit Exit-Code `10` je Fehlerfall)"; Beleg S4/S5 jetzt „nach Klarstellung: Exit-Code `10` je Fehlerfall"; F3-Tabellenzeile §4 trägt in der Entscheidungs-Spalte „**`10` je Fehlerfall** (Klarstellung … der 13/10-Vorschlag des Satzes S5 ist überholt)". Der Fragetext der F3-Zeile (Frage- und Empfehlungs-Spalte mit `13`/`10`) bleibt als Vorlagen-Historie der gestellten Frage stehen — von der Fix-Runde ausgewiesen, kein Widerspruch zur Entscheidung. | **behoben** |
| **F-2** (MEDIUM, neun Reststellen in den zwei Folge-Slices) | Eigene Suche nach `` `13` `` in den vier Planungsdateien (CR-Vorschlag, Entwurf, CR-Slice, Umsetzungs-Slice): **0 Reste** außer der F3-Frage-Zeile (Vorlagen-Historie, s. F-1). CR-Slice: Bezug- und Spec-Zeilen „endend je mit Exit-Code `10` (Sätze S5, S6, S8)"; Umsetzungs-Slice: DoD-Zeile 1 „je Exit-Code `10`", Plan-Tabelle „Fehler-Codes je `10` (S5 nach der F3-Klarstellung, Muster der Exit-Klassifikation [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes))", Test-Zeile „Fixture je Lesefall `10`", Abweisungs- und Kombinationsfälle je `10`. Bestand (`open/README.md` Zeilen 89–90) trägt je „Exit-Code `10`" samt Klarstellungs-Nennung. | **behoben** |
| **F-3** (MEDIUM, §3-Plan-Tabelle ohne Diff-Dateien) | `2d53e63d` ergänzte die neun Nachzugs-Zeilen (sichtbar im Diff des Slices: Messungen-Index, [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) + Index, CR-Vorschlag, CR-Slice, Roadmap, `docs/reviews/`); Abgleich gegen den Diff-Range oben — alle Entwurf-attributierten Dateien abgedeckt. | **behoben** |
| **F-4** (MEDIUM, Fundstellen-Drift) | Diff `2d53e63d` am Entwurf: Befund-Fundstelle jetzt „Zeilen 77 bis 118 … `ARG`-Zeilen 87 bis 92 … (`ENV`-Zeilen 113 bis 118)". Umsetzungs-Slice §3: „ENV-Zeilen 113 bis 118, Cache-Pfad-Kommentar Zeilen 109 bis 112". Am Code nachgelesen (`.devcontainer/Dockerfile` 75–120, Fundstellen oben in der Kopfzeile): exakt deckungsgleich. | **behoben** |
| **F-5** (LOW, Zustandsmarker auf Vor-Klarstellungs-Stand) | Der Fix der Fix-Runde („Freigabe-Präfix trennt F1/F2-wie-empfohlen von F3-nach-Klarstellung") traf nur den **Freigabe-Block des CR-Vorschlags**. Der benannte Fundstand — **DoD-Zeile 2 des Entwurfs** — trägt weiter „alle acht Stellen, wie vorgelegt; **F1 bis F3 entschieden, je wie empfohlen**" ([`slice-…-entwurf.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) §2, Beleg-Satz der Zeile). F3 wurde **nicht** je wie empfohlen gestellt (`10` nach Klarstellung gegen die Vorlage 13/10). Der Verweis auf den Freigabe-Block (der den Ist-Stand trägt) mildert, hebt den falschen Marker aber nicht — nach [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker) vor der Closure des Entwurfs zu berichtigen. | **offen am Fundstand** |
| **F-6** (INFO, Zusage „in beiden Profilen" gegen Messstand) | Bewusst getragen, nicht beruhigt: die Stelle-3-Zusage ruht auf der Golden-Case-/Docker-Test-Pflicht des Umsetzungs-Slices; die Messung deckt das Standardprofil für `containerEnv` ausweislich der Protokoll-Grenzen („ebensowenig `containerEnv` mit `${containerWorkspaceFolder}` im Standardprofil") nicht. Die Pflicht steht in der Umsetzungs-DoD (Beleg-Tabelle, Docker-Test), aber die Plan-Zeile der Golden-Case-Daten nennt das Standardprofil nicht ausdrücklich — als Rest-Risiko an die Umsetzungs-Closure übergeben (unten). | **offen als Rest-Risiko** (nach Auftragslage) |
| **F-7** (INFO, Sensoren nicht ausgeführt im Review) | Entfällt für diesen Lauf: der Verifier hat die Sensoren selbst gefahren bzw. den Stempel geprüft (Tabelle unten). | **entkräftet in diesem Lauf** |

---

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| Gate-Stempel | pass (abgeprüft, nicht neu gefahren) | `.harness/state/gates-passed.head` = `bd71715fb965…` = `git rev-parse HEAD` — der Stempel deckt den Prüfstand, ein erneuter `make gates` war nach Auftragslage nicht fällig. |
| CI der Mess-Tips | pass (4 Läufe) | `f76b7099`: `ci` 41410945794 **success** + `integration-docker` 41410946065 **success** (17m30s); `5e08fd50`: `ci` 41411456556 **success** + `integration-docker` 41411456578 **success** (16m27s) — per `gh run list` bestätigt. |
| CI auf dem Tip `2d53e63d`/`bd71715f` | läuft noch | Push 09:29/09:40 UTC: `ci` 41421125014 **success**, `integration-docker` 41421125007 in_progress; Tip `bd71715f` (41422617775/41422617683) beim Prüfmoment beide in_progress — kein roter Befund, aber der grüne Integrations-Nachweis für den finalen Tip steht noch aus. |
| `make gates` | nicht neu gefahren | Grund: Stempel deckt HEAD (oben). Der Stempel trägt `docs-check` + `ablage-pruefung` + Baseline-Verifikation als Teile des Gates. |
| `make test-docker` | nicht ausgeführt | Grund: kein Produkt-Pfad berührt (Diff des Slices: nur `docs/`, `tools/messung/`, Planning, ADR — `internal/**` ungeöffnet, §3-Zeile 9); die Docker-Pipeline ist über die CI-`integration-docker`-Läufe der zwei Mess-Tips belegt. |
| `make verify-depguard` | nicht ausgeführt | Grund: keine Import-/Architekturregel berührt (kein Go-Diff im Slice-Range). |

**Nicht ausgeführt (Pflicht nach `harness/verification.md`):** siehe Tabelle — alle mit Grund.

---

## Traceability

| ID | Beleg |
| --- | --- |
| [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) | Messprotokoll (Bezug, Fundstellen-Tabelle); Commits `f76b7099`/`5e08fd50` nennen sie; Grabstein-Stand unverändert am Code gelesene Fundstelle (`spec/lastenheft.md` 884–886) |
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) | CR-Vorschlag Stelle 3 (`containerEnv`-Form), Messzellen m1c/m2c/m3z |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | CR-Vorschlag Stelle 1 (Eintragstyp, S1–S11), Fundstellen am Katalog-Code; [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) `Schärft:` |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md) (Entwurfs-Kennung) | CR-Vorschlag Stelle 2 (Pin-Politik, vier Stellen, Abgleich per Test), F1 entschieden; Anker entsteht mit dem CR-Lastenheft-Commit |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | F2/F3-Entscheidungen (Abweisung je `10`), S5/S6/S8 |
| [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) | eigener Commit `56ce6c4b` (`Accepted`, ADR-Index Zeile 69), vom CR-Vorschlag und beiden Folge-Slices verlinkt |

Commit-Messages aller Entwurfs-Commits tragen mindestens eine `LH-*`- oder `ADR-`-Kennung bzw. Slice-ID (Commit-Hook-Vertrag), die Moves als `slice-mv` <!-- d-check:ignore (Make-Target, keine Slice-Kennung) --> rein.

---

## Carveouts

- Neu: **keine**. Gelöst: **keine**. Unverändert: [`carveouts.md`](../plan/carveouts/README.md) berührt der Range nicht. Keine Gate-Schwelle gelockert, Lastenheft unverändert (0.29.0, am Kopf gelesen).

## Rest-Risiken und offene Punkte (Übergabe an den Planner)

1. **F-5-Rest (klein, vor Closure zu berichtigen):** DoD-Zeile 2 des Entwurfs trägt „F1 bis F3 entschieden, je wie empfohlen" — F3 ist mit `10` gegen die Empfehlung gestellt. Eine Zeilenberichtigung im Entwurf-Slice; kein DoD-Gegenstand selbst.
2. **F-6 des Reviews — Golden-Case-Pflicht fürs Standard-Profil:** die Stelle-3-Zusage „in beiden Profilen" braucht, dass der Golden-/Docker-Test des Umsetzungs-Slices das `containerEnv`-Rendering auch im Standardprofil deckt; die aktuelle Plan-Zeile der Golden-Case-Daten nennt das Profil nicht ausdrücklich. Verifier-Hinweis an den Implementer des Umsetzungs-Slices, die Profilkante im Golden-Satz zu führen (Prüflinse [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
3. **F6 des CR-Vorschlags (ENV-Rest `BUILD_MODE`/`GOFLAGS`/`GOTOOLCHAIN`):** als offener Punkt an den §7 der Umsetzungs-Slice-Closure gebunden (Umsetzungs-Slice §1, letzter Ausdrücklich-NICHT-Punkt) — kein Defekt, ein benannter Rest.
4. **CI auf dem finalen Tip** (`bd71715f`) beim Prüfmoment noch in_progress; kein roter Befund auf keiner Vorstufe — bei der Closure der Tip-Stand gegen den Stempel erneut abgleichen, falls bis dahin etwas nachgeschoben wird.
5. **M7 / Offenes nichts weiter:** kein Release- oder Nutzerpfad im Slice; `open/README`-Bestandszeile des Entwurfs trägt den Stand; kein weiterer offener Lieferpunkt.

---

## Gesamturteil

**Closure-fähig** ([`slice-…-entwurf.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) §5), mit zwei Bedingungen:

- Vor dem Closure-Commit die **F-5-DoD-Zeile** berichtigen (Zustandsmarker auf Klarstellungsstand: F1/F2 wie empfohlen, F3 `10` entgegen der Vorlage) — §3.10-Abgleich ist sonst erfüllt (Plan-Tabelle deckt den Diff, F-3/F-4 behoben).
- Alle drei **DoD-Liefer-Zeilen sind belegt** (Messung mit Rohdaten je Lauf und Zusatzlauf `m3z`; CR-Wortlaut achttellig freigegeben mit F1–F3 entschieden und Fix-Runde nachgezogen — Suchbild `13` leer; zwei Folge-Slices in `open/` mit vollständigen Plan-Tabellen, Bestand und Roadmap getragen); Gate-Stempel deckt HEAD; CI der Mess-Tips grün; Rest-Risiken oben stehen, keines blockiert die Closure des Entwurfs — sie liegen bei den Folge-Slices bzw. der F-5-Berichtigung.

**Keine Änderung außer diesem Report** (Auftrag); keine Commits, keine Spec-Berührung.

---

## Commit / Artefakt

- Prüfstand: HEAD `bd71715fb9651a1efc98f2f0afab5ea8a4b37aea`; Diff-Range `b160d4ac^..HEAD` (75 Dateien, +3251/−49); Fix-Runde `2d53e63d`; Report-Berichtigung `bd71715f`.
- Beleg-Commits: `f76b7099` (M1–M3), `5e08fd50` (m3z), `b89943ea` (§4-Frage-3-Entscheidung), `3cb04e77` (CR-Vorschlag), `56ce6c4b` ([`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)), `dd3763a4` (Freigabe + zwei Folge-Slices), `6dd2a14d` (F3-Klarstellung), `2d53e63d` (F1–F5-Fix-Runde).