# MR-012 — `ai-harness-init` schreibt die Werkzeug-Dateien; u-boots Targets in `repo.mk`

- **Datum:** 2026-10-16
- **Geltungsbereich:** [`Makefile`](../../Makefile), [`repo.mk`](../../repo.mk),
  [`d-check.mk`](../../d-check.mk), `harness/mk/`, `tools/harness/` (außer
  `fetch-baseline-cache.sh`, `ablage-pruefung.sh` und `ablage-pruefung-test.sh`), `.claude/settings.json`, `.claude/hooks/`,
  [`harness/mk/ai-harness-init.md`](../mk/ai-harness-init.md),
  [`.d-check.yml`](../../.d-check.yml) (`targets`), [`harness/README.md`](../README.md) §Sensors.
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt](../../.harness/baseline/v6.18.1/regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt)
- **Adaption:** `ai-harness-init` ist das Werkzeug dieses Repos. Es schreibt bei jedem
  Lauf `Makefile`, `d-check.mk`, die Fragmente unter `harness/mk/`, die Skripte
  unter `tools/harness/` und die Hooks unter `.claude/hooks/` kanonisch neu; von Hand
  geänderte Fassungen gehen dabei verloren. Der Gate-Index hat zwei Teile
  ([`harness/README.md`](../README.md) für das Repo,
  [`harness/mk/ai-harness-init.md`](../mk/ai-harness-init.md) für das Werkzeug), und
  `.d-check.yml` misst gegen beide (`authority-disjoint`). Dazu gelten fünf
  Festlegungen des Repos:
  1. **Eigene Targets in `repo.mk`.** Die Go-/Docker-Targets (`lint`, `test`,
     `coverage-gate`, `build`, `ci`, `fullbuild` u. a.) und `GATE_CHECKS += lint test
     coverage-gate` stehen dort; das `Makefile` trägt keinen u-boot-Inhalt.
  2. **Drei Werkzeug-Targets stehen in der Repo-Tabelle.** `record-gates`, `slice-mv` <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
     und `hooks-install` führt [`harness/README.md`](../README.md) §Sensors mit der
     Bindung an diese Adaption; das Werkzeug lässt Targets aus, die der Index schon
     als Zeile führt.
  3. **`fetch-baseline-cache.sh` bleibt.** Es trägt Vendor-Lauf und Freshness-Audit
     ([`MR-004`](MR-004-regelwerk-vendored.md)); das Gate `baseline-verify` des Werkzeugs
     verifiziert den Bestand netzlos und ergänzt `make gates`. Ebenso repo-geführt
     sind `ablage-pruefung.sh` und sein Wächter-Test `ablage-pruefung-test.sh`
     (Target `ablage-pruefung` in `repo.mk`, Zeile in [`harness/README.md`](../README.md) §Sensors,
     Regel in [`conventions.md`](../conventions.md#ablage-von-messprotokollen-und-skripten)):
     Das Werkzeug schreibt nur die Dateien seiner eigenen Liste und lässt eine fremde
     Datei in `tools/harness/` bei einem erneuten Lauf stehen (an einem Testklon gelesen:
     ein zweiter Lauf entfernt eine fremde `*.sh` dort nicht).
  4. **Repo-geführte Dateien bleiben.** Skills, `implement-slice`, die Agenten und
     `.githooks/commit-msg` schreibt das Werkzeug nur an einem freien Pfad. `planner`,
     `plan-welle`, `close-welle` und das Beobachtungs-Register legt es an; sie sind auf u-boot
     gefasst ([`MR-003`](MR-003-roadmap-wellen.md)). `archivierung.mk` und `make archive-welle`
     bleiben ungenutzt: u-boot archiviert die Zeitdokumente geschlossener Wellen nicht.
  5. **Bekannte Grenze des Werkzeugs.** `tools/harness/working-tree-hash.sh` sortiert
     ohne feste Locale; der Gate-Nachweis ist locale-abhängig, wenn `make gates` und der
     Stop-Hook mit verschiedenen Locales laufen (der Hook blockt dann einmal).
- **Begründung:** Die Werkzeug-Dateien sind eine Quelle (das Programm); Handpflege
  daneben driftet beim nächsten Lauf. Eigene Targets gehören in `repo.mk`, weil das
  Werkzeug das `Makefile` ohne Meldung ersetzt.
- **Auflösungs-Trigger:** u-boot gibt das Werkzeug auf, oder u-boot führt die Archivierung der Zeitdokumente geschlossener Wellen ein (Punkt 4).
- **Löst auf:** [`MR-011`](../conventions.md#mr-011)
- **Ausgelöst durch Baseline-Stand:** v6.18.1
