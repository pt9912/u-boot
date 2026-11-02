# Review-Report: slice-v1-katalog-gepinnte-toolchain-entwurf — 2026-11-02

**Review-Art:** Plan (Entwurf-Slice gegen Plan, Entscheidungen des
Projektinhabers und Hard Rules; Messprotokoll, CR-Vorschlag, ADR und
Folge-Slices als Berührflächen; **nicht** gegen die DoD — Verifier-Ebene)

**Gegenstand:** [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md),
alle Commits seit dem Move (`8353b5a1`/`4910bf9d`/`e04e3329`) bis HEAD
`6dd2a14d`, geprüft als Diff `e04e3329^..HEAD` (59 Dateien, +2673/−17):
Messung M1–M3 und Zusatzlauf `m3z` (`f76b7099`, `5e08fd50`), vier
§4-Entscheidungen (`63e16263`, `b89943ea`), CR-Vorschlag (`3cb04e77`),
[`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (`56ce6c4b`), Folge-Slices + Freigabe-Block (`dd3763a4`),
F3-Klarstellung (`6dd2a14d`).

**Skill:** `.harness/skills/reviewer.md` @ `v6.18.1` ·
**Modell:** glm-5.3-flash · **Datum:** 2026-11-02

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md)
  (Stand `6dd2a14d`), CR-Vorschlag [`cr-vorschlag-gepinnte-werkzeugkette`](../plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md)
  (Freigabe 2026-11-02, F3-Klarstellung `6dd2a14d`)
- [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
  (`Accepted`, eigener Commit `56ce6c4b`), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md),
  [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md)
- Messprotokoll [`docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md`](../messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md)
  mit 40 Rohdaten-Dateien; Messskripte `tools/messung/katalog-gepinnte-toolchain/`
  (5 Dateien)
- Folge-Slices [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md)
  und [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen),
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)
  (Grabstein), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe),
  [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)
- `AGENTS.md` (Hard Rules, §§3.4, 3.7, 3.9, 3.10),
  [`harness/conventions.md`](../../harness/conventions.md)
  (`MR-008` ADR-Form, Ablage-Regel), Beobachtungs-Register
  ([`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md),
  [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md),
  [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md),
  [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../plan/planning/observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md))

---

## Findings

Schema nach Skill §Output-Schema (`v6.18.1` · `regelwerk/modul-10-review-harness.md`
§Ziel-Form: Reviewer-Skill).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die F3-Klarstellung (`6dd2a14d`: alle Übersteuerungs-Fehlerfälle Exit `10`) wurde nur in Satz S5 und der F3-Tabellenzeile nachgezogen; die Freigabe-Notiz oben im CR-Vorschlag trägt weiter „F3 — die Exit-Codes … wie im Satz S5: `13` (Quelle nicht lesbar), `10` (Version nicht überdeckt)“, und der Satz-Beleg S4/S5 nennt „F3 … `13`/`10` wie im Satz“. Beide Stellen widersprechen dem freigegebenen Wortlaut desselben Dokuments. | [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); `AGENTS.md` §3.7, §3.9(b) | `docs/plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md` · „**F3** — die Exit-Codes der Übersteuerungs-Fehlerfälle wie im Satz S5: `13` (Quelle nicht lesbar)“ | Review-only | Veralteter Wert nach Entscheidungs-Nachzug |
| F-2 | MEDIUM | Der F3-Ausgang ist in beiden Folge-Slices nicht nachgezogen: sie ordnen dem Fehlerfall „Quelle nicht lesbar“ weiterhin Exit-Code `13` zu (CR-Slice: Bezugs-Zeile „Exit-Codes `10`/`13`“, Spec-Stellen-Zeile „S5, S6, S8“, Plan-Tabelle „Abweisungs-Fälle S5 (`13`/`10`)“, CHANGELOG-Zeile „Fehler-Codes `13`/`10`“; Umsetzungs-Slice: Bezug, Ziel, DoD, Plan-Tabelle „Fehler-Codes `13`/`10` (S5 …)“, Test-Zeile „Fixture je Lesefall `13`/`10`“) — neun Stellen gesamt. Die Slices sind vor der Klarstellung (`dd3763a4`) geschnitten und von `6dd2a14d` nicht erreicht worden. | [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); `AGENTS.md` §3.7 | [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) · „Abweisungs-Fälle mit Exit-Code `10`/`13` (Sätze S5, S6, S8)“; [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) · „Fehler-Codes `13`/`10` (S5, Muster der Exit-Klassifikation“ | Review-only | Veralteter Wert nach Entscheidungs-Nachzug |
| F-3 | MEDIUM | Der §3.10-Abgleich ist offen: der Diff `e04e3329^..HEAD` berührt Dateien mit keiner Plan-Zeile in §3 des Entwurfs — [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) samt ADR-Index, die CR-Vorschlag-Datei, der CR-Slice (nur der Umsetzungs-Slice ist in der Zeile genannt), die Roadmap-Zeile `v0.19.0` und der Messungen-Index `docs/messungen/README.md`; dazu die Verweis-Nachzüge an `done/`-Slice und Review-Reports (lifecycle-mechanisch). | `AGENTS.md` §3.10; [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) §3 · „diese Datei, `docs/plan/planning/open/README.md`, Umsetzungs-Slice \| neu / update“ | Review-only | Plan-Tabelle berücksichtigt Diff-Datei nicht |
| F-4 | MEDIUM | Fundstellen-Drift zum Hand-written-Teil: der Entwurf-Befund nennt „`ARG`-Zeilen 83 bis 88“ und „`ENV`-Zeilen 97 bis 107“, der Plan des Umsetzungs-Slices führt dies fort („ENV-Zeilen 97 bis 107, Cache-Pfad-Kommentar Zeilen 117 bis 118“). Am Code trägt der `ARG`-Block die Zeilen 87–92, der `ENV`-Block 113–118 und der Cache-Kommentar 109–112 (`.devcontainer/Dockerfile` gelesen am Stand `6dd2a14d`). Das Messprotokoll des selben Vorgangs und der CR-Vorschlag („ab Zeile 113“, „Zeilen 84 bis 85“, „93 bis 108“) korrigieren die Zahlung teilweise — die Plan-Fassungen widersprechen den Beleg-Dokumenten des selben Vorhabens. | `AGENTS.md` §3.10; [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) · „(`ENV`-Zeilen 97 bis 107)“; [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §3 · „ENV-Zeilen 97 bis 107, Cache-Pfad-Kommentar Zeilen 117 bis 118“ | Review-only | Fundstellen-Drift zu gelesener Stelle |
| F-5 | LOW | Der Zustandsmarker der DoD-Zeile 2 des Entwurfs („Freigabe … F1 bis F3 entschieden, je wie empfohlen — Freigabe-Block oben im CR-Vorschlag“) steht auf dem Vor-Klarstellungsstand: F3 wurde gegen die vorgelegte Empfehlung (13/10) mit „`10` je Fehlerfall“ gestellt — die F3-Zeile des CR-Vorschlags trägt die Klarstellung, „je wie empfohlen“ beschreibt sie nicht mehr. | `AGENTS.md` §3.7, §3.9 | [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) §2 · „F1 bis F3 entschieden, je wie empfohlen“ | Review-only | Zustandsmarker auf Vor-Entscheidungs-Stand |
| F-6 | INFO | Stelle 3 des CR-Vorschlags sagt „Zugesagt ist die Erzeugung dieser Werte in der erzeugten Datei in beiden Profilen“ zu; die Messung (m1c/m3z) ist im Sandbox-Profil gemessen, das Standardprofil trägt kein Protokoll-Stand der Erzeugung. Die Pin-Pflicht läuft über die Golden-Case-Zeile des CR-/Umsetzungs-Slices — Hinweis an Verifier und Umsetzungs-Implementer, dass der Golden Case das `containerEnv`-Rendering auch im Standardprofil decken muss, sonst steht die Zusage ohne Beleg ([`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). | [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | `docs/plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md` §2, Stelle 3 · „in beiden Profilen“ | Review-only | Zusage mit Beleg-Pflicht an späteren Test gebunden |
| F-7 | INFO | Kein Sensor wurde in diesem Review-Lauf ausgeführt (Review-only-Auftrag, keine Gate-Läufe, keine Commits); die Angabe „`make gates` grün (Stempel `6dd2a14d` = HEAD), `docs-check` 0 Befunde, CI-Läufe der Tips `f76b7099` und `5e08fd50` je success“ ist Auftragskontext des Hauptagenten und nicht eigenständig nachgelesener Befund. | Workflow (Reviewer-Rolle) | — | ja — `make gates` beim Hauptagenten | Nicht ausgeführte Sensors |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Messprotokoll (DoD-Form: je Zelle zwei Läufe, datiert 2026-11-02, Umgebung `linux-host amd64`, Engine 29.9.2, CLI 0.89.0; verworfene Vorläufe m2-Reihe und m3-Vorlauf vermerkt; Inventar-Abschnitt je Rohdaten-Log; Projekt-Ordner, Feature-Images und Mess-Image geräumt) | geprüft, ohne Befund |
| Rohdaten auf Zugangsdaten (Sweep `token`/`secret`/`password`/`ghp_`/`AKIA`/`Bearer` über alle 41 Dateien von Protokoll und Skripten) | geprüft, ohne Befund ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)) |
| Rohdaten-Stichprobe (`rohdaten-m3-lauf1.log`/`.txt`: Zeitstempel, Container-ID-Wechsel, Volume, Teardown-Inventar leer) und Ablage nach `harness/conventions.md` §Ablage (`docs/messungen/` vs `tools/messung/`) | geprüft, ohne Befund |
| CR-Vorschlag: keine Zusage weiter als [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) oder Messung in den Grenz-Sätzen (S10/S11, Stelle-3-Grenz-Sätze, Stelle 4 „Was der Docker-Test belegt“); Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) zitiert, nicht belebt (neue Kennung `LH-FA-DEV-011` <!-- d-check:ignore (Kennung des Entwurfs; Anker entsteht mit dem Lastenheft-Commit) --> als Entwurf); `d-check:ignore`-Marker als Übergangsstand benannt | geprüft, ohne Befund (Grenz-INFO F-6) |
| [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md): MADR-Form nach `MR-008`, `Schärft:`-Träger konsistent, keine Abwärtsecken zu Slice/Roadmap/Carveout (§3.4), keine Duplication des CR-Wortlauts, ADR-Index nachgezogen, eigener Commit | geprüft, ohne Befund |
| `spec/lastenheft.md` im Prüf-Range unverändert (Version 0.29.0, Grabstein-Zeilen 884–886 unverändert, `LH-FA-DEV-011` <!-- d-check:ignore (Kennung des Entwurfs; Anker entsteht mit dem Lastenheft-Commit) --> frei) | geprüft, ohne Befund |
| Fundstellen-Stichproben des CR-Vorschlags und der Folge-Slices (Lastenheft 227/231, 414/418, 377–390/383/385, 1614–1616; [`SPEC-053`](../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) als letzte Spec-Kennung; `featureVersionNeedsSource` Zeilen 798/731, `ollamaWarnings`/`ollamaHostWarnings` 589/616, `TestE2E_…_PinMismatchBreaksBuild`) | geprüft, ohne Befund (die gestrichenen `.devcontainer`-Fundstellen stehen in F-4) |
| Folge-Slices: Start-Bedingungen konsistent (CR-Slice vor Umsetzungs-Slice, WIP-Limit 1 mit laufendem Entwurf, Freigabe-Beleg 2026-11-02, `Verantwortlich` im Umsetzungs-Slice gesetzt), `open/README`-Bestand und Roadmap-Zeile gezogen, Träger-Kennungen verlinkt | geprüft, ohne Befund (Exit-Codes: F-2) |
| Sprach-Grenzen: keine Handbuch-/README-/CHANGELOG-Berührung im Prüf-Range; die CR-Slice-Plan-Tabelle ordnet diese Stellen korrekt dem CR-Slice zu, CHANGELOG-Muster Englisch, Handbuch ohne Kennungen | geprüft, ohne Befund |
| Messskripte (`lib.sh`, `m1`–`m3z`): Kommentare tragen Zusage/Kopplung, kein §3.7-Muster (Sweep nach Chronik-Formulierungen leer), Namespace-Präfix und `leftovers after teardown`-Inventar im Skript verankert | geprüft, ohne Befund |
| Slice-Edits §4: die fünf Entscheidungen nennen jeweils Datum, Träger und Fundstelle (m1c/m2c/m3z, `devcontainer.json.tmpl` 12–13/33–41, [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)-Stellenliste) — keine Aussage ohne Fundstelle gefunden | geprüft, ohne Befund |

## Nicht geprüft

| Bereich | Grund |
|---|---|
| `make gates`, `make docs-check`, `make test-docker`, `make modus-pruefung` | Review-only-Auftrag (keine Gate-Läufe, keine Commits); Sandbox-Umgebung. Auftragskontext: grün, Stempel `6dd2a14d` = HEAD — eigenständig nicht verifiziert (F-7) |
| DoD-Häkchen, Lieferpunkte, Verification-/Closure-Bedingungen, Ausgang je §6-Risiko | Verifier-Rolle (dieser Lauf prüft nicht gegen die DoD) |
| Validator-Ebene (reale Nutzerbedarf-Prüfung der CR-Zusagen, etwa Stelle 3 `containerEnv`-Nutzbarkeit) | Validator-Rolle |
| Sämtliche 40 Rohdaten-Dateien einzeln im Volltext | Umfang; geprüft per Protokoll-Tabellen gegen die Datei-Bestände, Volltext-Stichprobe `m3-lauf1` + Credential-Sweep über alle Dateien |
| Andere `open/`-Slices (etwa [`slice-v1-katalog-feature-jdk-gepinnt`](../plan/planning/open/slice-v1-katalog-feature-jdk-gepinnt.md)) außer den zwei Folge-Slices | außerhalb des Auftrags |
| Interpretation der F3-Klarstellung (welche drei Fehlerfälle „alle drei“ zählt) | Entscheidung des Projektinhabers; der Reviewer prüft den Nachzug, nicht die Entscheidung |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Veralteter Wert nach Entscheidungs-Nachzug ·
Plan-Tabelle berücksichtigt Diff-Datei nicht · Fundstellen-Drift zu gelesener
Stelle · Zustandsmarker auf Vor-Entscheidungs-Stand · Zusage mit
Beleg-Pflicht an späteren Test gebunden · Nicht ausgeführte Sensors

## Verdikt

**Merge-blockierend:** ja — die MEDIUM-Findings blockieren den Start des
CR-Slices, nicht dessen Plan: F-1 und F-2 sind derselbe Wurzelbefund (`6dd2a14d`
hat den F3-Ausgang nur teilweise nachgezogen; seine Commit-Message behauptet,
„Freigabe-Zeile und F3-Tabelle tragen die Entscheidung“, während die Diff die
Freigabe-Notiz in Zeile 5 und den Satz-Beleg S4/S5 unberührt ließ). Der
CR-Slice ist an angewiesen, den Wortlaut **wortgleich** zu schreiben und
[`SPEC-054`](../../spec/spezifikation.md) samt CHANGELOG nach der Plan-Tabelle
zu füllen — folgt er den Plänen, zieht Exit-Code `13` gegen den freigegebenen
Satz S5 (`10`) in einen Produktvertrag, kein Gate fängt das semantisch. F-3 und
F-4 fallen unter [`§3.10`](../../AGENTS.md) und sind vor der Closure des
Entwurfs zu berichtigen (Plan-Tabelle gegen `git diff --stat`, Fundstellen der
`.devcontainer`-Zeilen).

**Wichtigste Abhilfe-Reihenfolge für die Closure:**
1. F-1/F-2 in einem Berichtigungslauf: Freigabe-Notiz und Satz-Beleg im
   CR-Vorschlag auf „Exit `10` je Fehlerfall“ stellen und die neun
   `13`-Stellen in den beiden Folge-Slices nachziehen (Suchbild: `` `13` ``
   in `docs/plan/planning/`, nur die Lastenheft-Definition des Codes `13`
   selbst und die F3-Frage-Zeile bleiben als Zitat).
2. F-3: §3-Plan-Tabelle des Entwurfs um [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (+Index), CR-Vorschlag,
   CR-Slice, Roadmap und Messungen-Index ergänzen oder als ausgewiesene
   Nicht-Einträge führen — Abgleich gegen `git diff --stat` wiederholen.
3. F-4: Fundstellen des Hand-written-Teils auf `ARG` 87–92, `ENV` 113–118,
   Kommentar 109–112 berichtigen (Entwurf §1 Befund und
   Umsetzungs-Slice §3-Vorspann).
4. F-5: DoD-Zeile 2 des Entwurfs auf den Klarstellungsstand (F3: `10`
   entgegen der Vorlage) stellen.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review → Plan bei
Plan-Defekt, hier: F-3/F-4); die Finding-Klassen gehen in die Slice-Closure §7
und von dort in den Zähler. Das Befund-Bild stützt die Closure des Entwurfs
nicht negativ — Messung, CR-Wortlaut-Umfang und ADR sind gegen die
Entscheidungen geprüft und ohne contract-überschreitende Zusage.