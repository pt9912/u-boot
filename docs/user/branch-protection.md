# Branch Protection für `main`

[LH-QA-003](../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) verlangt, dass die drei CI-Jobs aus `.github/workflows/ci.yml`
für jeden Pull-Request **blockierend** sind. Die exakten Required-
Status-Check-Namen sind die Workflow-`name:`-Felder
(`gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`,
`image-scan (trivy HIGH+CRITICAL)`), nicht die kürzeren
`jobs.<key>`-Identifier — siehe Schritt 4 „Require status checks to
pass before merging" → Bullet „Status checks that are required" unten.
Der GitHub-Actions-Workflow allein reicht für die Aktivierung nicht —
er muss zusätzlich in den Repository-Settings als **Required Status
Check** eingetragen werden. Diese Aktivierung lebt im GitHub-UI, nicht
im Repo-Code; um sie reproduzierbar zu halten, dokumentiert diese
Datei die Schritte.

Der zugehörige Carveout ([ADR-0004](../plan/adr/0004-ci-system.md) Folgepunkt „Branch-Protection nicht
versioniert") wird mit dem Vorhandensein dieser Checkliste aufgelöst.
Spätestens vor dem ersten externen PR müssen die Schritte einmalig im UI
ausgeführt sein.

## Einmalige UI-Aktivierung

1. **GitHub → Settings → Branches → Add branch protection rule.**
2. **Branch name pattern:** `main`.
3. **Require a pull request before merging:** aktivieren.
   - **Required approvals:** für ein Solo-Projekt bewusst `0`. Sobald ein
     zweiter Contributor dazukommt, auf `1` anheben.
   - **Dismiss stale pull request approvals when new commits are pushed:**
     aktivieren (sinnvoll auch im Solo-Setup für die Zukunft).
4. **Require status checks to pass before merging:** aktivieren.
   - **Require branches to be up to date before merging:** aktivieren.
   - **Status checks that are required:** mindestens — die exakten
     Display-Namen sind die `name:`-Felder aus `.github/workflows/ci.yml`,
     nicht die kürzeren `jobs.<key>`-Identifier. Bei einer Hebung der
     `name:`-Werte muss diese Checkliste mitgezogen werden.
     - `gates (lint + test + coverage-gate)` ([LH-QA-003](../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions))
     - `security-gates (govulncheck)` ([LH-QA-003](../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions))
     - `image-scan (trivy HIGH+CRITICAL)` ([LH-QA-003](../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions), geliefert mit
       [`slice-v1-release-pipeline`](../plan/planning/done/slice-v1-release-pipeline.md)
       T3, siehe
       [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md))
5. **Require conversation resolution before merging:** aktivieren (Review-
   Kommentare müssen aufgelöst sein, bevor gemerged werden darf).
6. **Restrict who can push to matching branches:** für Solo-Projekt
   irrelevant; bei mehreren Contributors einschränken.
7. **Do not allow bypassing the above settings:** **aktivieren** —
   das ist der Toggle in der GitHub-UI; aktivieren bedeutet „Schutz-
   regeln gelten auch für Admins / Repo-Owner". Ohne diesen Toggle
   kann der Repo-Owner die Regeln umgehen, was den Zweck aushebelt.
8. **Allow force pushes:** **deaktivieren** (Default).
9. **Allow deletions:** **deaktivieren**.
10. Optional: **Require linear history:** aktivieren, wenn das Projekt
    rebase- statt merge-orientiert arbeiten soll.

## Verifikation

Nach der Aktivierung:

- Ein PR mit absichtlich fehlschlagendem `gates` muss als „Merge blocked"
  angezeigt werden.
- Ein direkter Push auf `main` (z. B. `git push origin main`) muss vom
  Server abgewiesen werden, auch für den Repo-Owner.
- `git push --force origin main` muss vom Server abgewiesen werden.
- `git push origin :main` (Branch-Delete) muss abgewiesen werden.

## Optional: Repository-Ruleset-Export

GitHub bietet ab den Repository-Rulesets die Möglichkeit, die obigen
Regeln als JSON zu exportieren und im UI wieder zu importieren. Sobald
das Repo öffentlich oder team-geteilt wird, lohnt sich ein Export nach
`docs/user/branch-protection-ruleset.json` als zusätzliche Quelle der
Wahrheit. Für das Solo-Bootstrap reicht diese Markdown-Checkliste.

## Bezug

- Auslösende Spec: [LH-QA-003](../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) (drei Jobs `gates` / `security-gates` /
  `image-scan`, alle PR-blockierend; image-scan via
  [slice-v1-release-pipeline](../plan/planning/done/slice-v1-release-pipeline.md) T3 in den Spec-Pflicht-Block geschrieben);
  [ADR-0004](../plan/adr/0004-ci-system.md) Folgepunkt „Branch-Protection nicht versioniert".
- Slice: [`slice-v1-release-pipeline`](../plan/planning/done/slice-v1-release-pipeline.md)
  (Teilabschluss Branch-Protection 2026-05-27; T2 `publish.yml` +
  T3 `image-scan` 2026-05-31; [LH-OPEN-002](../../spec/lastenheft.md#lh-open-002--paketierung)-Restwege bleiben mit
  Trigger-Slices vertagt — siehe
  [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md)).
- Carveout-Inventar: [`carveouts.md`](../plan/planning/in-progress/carveouts.md).
