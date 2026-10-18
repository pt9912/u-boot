# Hilfetext: Ein Backtick-Wort im Usage-Text einer Option erscheint als Typname

**Sub-Area:** `internal/adapter/driving/cli` (Kürzel CLI, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Hilfe-Ausgabe ([`LH-FA-CLI-002`](../../../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe)) ist dort Konvention; die Usage-Texte der Optionen stehen im Adapter.

Cobra liest in dem Usage-Text einer Option das erste in Backticks gesetzte Wort als Namen des Wertes und zeigt es in der Hilfe als Typnamen (`--quiet up` statt `--quiet`), obwohl die Option keinen Wert nimmt. Der Befund stammt aus dem Bericht des Implementers von [`slice-v1-cli-egress-anzeige`](../../../done/slice-v1-cli-egress-anzeige.md); der Planner hat ihn weder in der Hilfe-Ausgabe noch im Code nachgefahren.
