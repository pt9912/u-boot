# Clone-Host: Eine lokale Pfad-URL wird zum „Host“ der Egress-Liste

**Sub-Area:** `internal/hexagon/application` (Kürzel APP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Berechnung der Egress-Liste ([`LH-FA-DEV-008`](../../../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) liegt in der Anwendungsschicht.

`cloneSourceHost` (`internal/hexagon/application/devcontainer_egress.go`) leitet den Host der Clone-Quelle ab: Ohne `://` und ohne `@` schneidet die Funktion die URL an `:` und gibt den Rest klein geschrieben zurück. Eine lokale Pfad-URL wie `/srv/git/repo.git` besteht die Prüfung der zulässigen Zeichen und ergibt den „Host“ `/srv/git/repo.git`; dieser Name steht dann in der `hosts=`-Zeile des Startscripts und in `u-boot egress list` mit der Herkunft `clone`. Das Verhalten ist vorbestehend; `u-boot egress list` macht es sichtbar. Gelesen am Quelltext, nicht gefahren.
