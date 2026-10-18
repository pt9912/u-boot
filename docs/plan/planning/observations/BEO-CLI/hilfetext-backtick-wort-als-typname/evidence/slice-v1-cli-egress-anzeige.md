**Vorgang:** [`slice-v1-cli-egress-anzeige`](../../../../done/slice-v1-cli-egress-anzeige.md)
**Fund:** Bericht des Implementers zu diesem Slice: Die Hilfe zeigte `--quiet up` statt `--quiet`, weil Cobra das Backtick-Wort `up` im Usage-Text der Option als Typnamen des Wertes las. Der Planner hat den Befund aus dem Bericht übernommen; er ist weder im Review-Report noch im Verification-Report ein Finding, und die Stelle im Code (Usage-Text der Option `--quiet`) hat der Planner nicht gelesen.

**Beleg-Qualität:** schwach — nur die Angabe des Implementers; kein Test im Repo ist dem Planner als Pin für die Hilfe-Zeile bekannt. Ein Nachfahren (`u-boot --help` im Docker-Image) steht aus.
