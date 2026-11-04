# Docker-getaggte Tests pinnen einen Vertragswert, den `make gates` nicht prüft

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — ein Test mit Build-Tag `docker` pinnt einen Wert des erzeugten Vertrags (hier die Liste `runArgs`); `make gates` übersetzt und fährt solche Tests nicht, ein Vertragswechsel lässt sie unbemerkt rot, bis die volle Suite von `make test-docker` läuft.

Der Implementer eines Vertragswechsels fährt `make gates` und gezielt die Docker-Tests, die er kennt; ein weiterer Docker-Test, der denselben Wert gleich-läufig prüft, fällt erst im Volllauf (rund 27 Minuten) auf.
