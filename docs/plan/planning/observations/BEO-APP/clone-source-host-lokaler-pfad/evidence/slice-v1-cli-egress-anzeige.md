**Vorgang:** [`slice-v1-cli-egress-anzeige`](../../../../in-progress/slice-v1-cli-egress-anzeige.md)
**Fund:** Review-Report [`2026-10-18-slice-v1-cli-egress-anzeige`](../../../../../../reviews/2026-10-18-slice-v1-cli-egress-anzeige.md) F-11 (INFO) und Verification-Report [`2026-10-18-verification-slice-v1-cli-egress-anzeige`](../../../../../../reviews/2026-10-18-verification-slice-v1-cli-egress-anzeige.md) F-11/Auflage 4 (a): `cloneSourceHost` macht aus `/srv/git/repo.git` den „Host“ `/srv/git/repo.git`.

**Beleg-Qualität:** gelesen, nicht gefahren — der Reviewer, der Verifier und der Planner lasen den Quelltext von `cloneSourceHost` (Schnitt an `:` ohne `://` und ohne `@`); kein Test pinnt das Verhalten, kein Lauf zeigt die erzeugte `hosts=`-Zeile für eine lokale Pfad-URL.
