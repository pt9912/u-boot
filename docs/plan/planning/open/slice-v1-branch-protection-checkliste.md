# Slice V1: Branch-Protection-Checkliste

## Auslöser

ADR-0004 erwähnt: *„Branch-Protection im UI ist nicht im Repo
versioniert."* Required-Status-Checks für `gates` und `security-gates`
müssen im GitHub-UI manuell aktiviert werden, sonst sind die beiden
CI-Jobs zwar grün, aber **nicht PR-blockierend** — was die zentrale
Pflicht von `LH-QA-003` aushebelt (`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Eine versionierte Checkliste in `docs/user/` (z. B.
`docs/user/branch-protection.md`), die einmalig beim Bootstrap des
GitHub-Repos abgearbeitet wird. Optional ergänzt durch einen
GitHub-Repository-Ruleset-JSON-Export für reproduzierbares Setup.

## Akzeptanzkriterien

- `docs/user/branch-protection.md` beschreibt Schritt-für-Schritt:
  - Repository → Settings → Branches → Add branch protection rule für
    `main`.
  - Required status checks: `gates` und `security-gates` (beide
    einzeln required).
  - Require PR before merging: ja, mindestens 1 Approval (oder als
    Solo-Projekt: bewusst null, dokumentiert).
  - Block force pushes auf `main`.
  - Block branch deletion.
  - Optional: linear history erzwingen.
- Optional `docs/user/branch-protection-ruleset.json` als
  GitHub-Repository-Ruleset-Export (importierbar via UI oder API).
- README (de/en) Section „Setup" verweist auf die Checkliste.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf den Aufhebungs-Commit als gelöst markieren.

## Out of Scope

- DCO-Bot-Aktivierung (separater ADR-0004-Folgepunkt; lebt im
  GitHub-Marketplace, kein Repo-Artefakt).
- CODEOWNERS-Datei (eigener Slice, wenn Teilautoren dazukommen).

## Bezug

- Auslösende ADR: `0004-ci-system.md` Folgepunkte.
- Auslösende Spec: `LH-QA-003` „beide Jobs PR-blockierend".
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Branch-Protection nicht versioniert.
- Hängt von: erstem PR-Workflow (vermutlich M3-PR).
