# docs/reviews

Ablage der **Review-Reports** — das Übergabe-Artefakt der Reviewer-Rolle an die
Implementation. Ein Report hält fest, *wogegen* geprüft wurde, *was* gefunden
wurde und *was ohne Befund* blieb.

Kanonische Review-Regeln stehen in
[`harness/review.md`](../../harness/review.md); die repo-spezifischen
Prüf-Anker im Reviewer-Skill unter `.harness/skills/reviewer.md`. Dieses
Verzeichnis ist nur die Ablage.

## Konvention

- **Ein Report pro Lauf.** Folgeläufe erzeugen eine **neue Datei**, kein
  Überschreiben — sonst geht die Historie der Findings verloren und der
  Steering-Loop („dasselbe Finding zum dritten Mal") ist nicht belegbar.
- **Dateiname:** `<YYYY-MM-DD>-<slice-oder-diff-ref>.md`. Mehrere Läufe am
  selben Tag zum selben Gegenstand bekommen ein Suffix (`-r2`, `-r3`).
- **Kopiervorlage:**
  `.harness/baseline/v3.5.1/templates/docs/reviews/review-report.template.md`
  (vendorte Referenz-Form). Es liegt bewusst **keine** zweite Kopie der Vorlage
  in diesem Verzeichnis — eine Kopie wäre eine Drift-Quelle gegenüber der
  Baseline.
- **Negativbefunde sind Pflicht.** Jeder Report nennt die geprüften Bereiche
  ohne Befund und die *nicht* geprüften Bereiche mit Grund. Ohne diesen Block
  ist „keine Findings" nicht von „nicht geprüft" unterscheidbar.
- **Reports sind keine Verification-Evidence.** Sie können Evidence auslösen;
  der DoD-Abgleich bleibt Verifier-Aufgabe nach
  [`harness/verification.md`](../../harness/verification.md).

## Abgrenzung

| Artefakt | Frage | Ort |
|---|---|---|
| Review-Report | Welche Risiken/Vertragsbrüche enthält der Diff? | dieses Verzeichnis |
| Verification Evidence | Wurde der Slice gemessen an DoD/Spec richtig gebaut? | im jeweiligen Slice unter [`docs/plan/planning/done/`](../plan/planning/done/) |
| Closure-Notiz | Was ist das Lernsignal für den Steering-Loop? | im jeweiligen Slice |

Die Ortswahl (`docs/reviews/` neben `.harness/skills/`) ist im Adaptions-Ledger
[`harness/conventions.md`](../../harness/conventions.md) als `MR-009`
dokumentiert.
