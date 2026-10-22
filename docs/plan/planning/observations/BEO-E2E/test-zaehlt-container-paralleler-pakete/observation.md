# Ein Docker-Test zählt nach einem Namenspräfix und sieht die Container eines parallel laufenden Pakets

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Fehler liegt im Test, nicht im Produkt; er trifft jeden Lauf von `make test-docker`, weil `go test ./...` die Pakete gleichzeitig gegen eine Engine laufen lässt.

Ein Docker-Test prüft eine Aussage über „alle Container mit diesem Namenspräfix“ (hier: während der Sitzung besteht kein Hilfs-Container). Das Präfix ist das des Produkts, und die Adapter-Tests im Paket `internal/adapter/driven/docker` vergeben dasselbe Präfix, verlängert um eine Testkennung. Läuft ein solcher Container des anderen Pakets zur Zeit der Zählung, zählt der Test ihn mit und färbt sich rot, obwohl das Produkt nichts hinterlassen hat. Lokal fällt es nicht auf, wenn die Fenster der Pakete nicht überlappen; kein Gate prüft, dass eine Suche nur die Container des eigenen Tests sieht.

## Benannt, nicht gezählt

- **Parallelgrenze der Hilfsfunktion `newContainersWithCommand`** (`internal/e2e/sandbox_devcontainer_docker_test.go`): Sie listet alle Container und sucht einen Pin in Pfad, Argumenten und Umgebung; ihr Kommentar benennt, dass zwei Läufe gegen denselben Daemon nicht unterschieden werden. Ein Container eines parallelen Pakets mit demselben Pin ist denkbar, nicht belegt (Nachlauf-Review `slice-v1-cli-devcontainer-shell-ci-fix`, I-4). Kein Vorgang.
- **Verwandt, nicht gezählt:** der Race zwischen Paketen um das Hilfs-Image, das ein Fundament-Test entfernte, während Tests des Pakets `internal/e2e` es benutzten (Lese-Schritt in [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../test-gruen-nur-mit-lokalem-cache/state.md), Vorgang `slice-v1-cli-devcontainer-up-down`). Gleicher Mechanismus (Pakete teilen eine Engine), anderes Objekt (Image statt Container); als zweiter Beleg wird er erst gezählt, wenn jemand die Evidence des Vorgangs gelesen und das Objekt benannt hat.

**Trigger:** ein zweiter Vorgang, in dem ein Docker-Test wegen eines Objekts rot oder vakuum-grün wird, das ein parallel laufendes Paket auf derselben Engine angelegt oder entfernt hat.
