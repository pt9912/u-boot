# Test grün nur mit lokalem Cache

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Tests laufen lokal mit gewärmtem Image-Cache grün und werden auf einem Runner ohne Cache rot; der Fehler liegt im Test, nicht im Produkt.

Ein e2e-Test, der ein Registry-Image ohne vorherigen Pull startet und dessen Ausgabe auswertet, hängt vom Cache des Hosts ab: Lokal schreibt `docker run -d` nur die Container-ID, auf einem frischen Runner zusätzlich den Pull-Text. Kein Gate fängt das lokal, weil der Cache lokal immer gefüllt ist.

## Benannt, nicht gezählt

- **Flake in `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` (30,05 s, ein Lauf von vier ohne Meldung):** Ursache unbelegt, nicht reproduzierbar (28 Läufe des Tests ohne Fehlschlag, davon 12 nach `docker rmi` des Images, 4 unter gleichzeitiger CPU-Last). Die 30,05 s liegen nahe am 30-s-Zeitlimit des Registry-Clients im Daemon; eine kurzzeitige Registry-Störung beim neuen `docker pull -q` des Tests ist die plausibelste Vermutung, aber ohne Meldung nicht belegt. Beim nächsten Auftreten die volle Ausgabe sichern (`-v`, Datei): `docker pull alpine:3: …` bedeutet Netz oder Registry, `seed files`, `FindContainer` oder `ReadFile` einen Fehler in Test oder Adapter. Der Flake gilt nicht als behoben und zählt nicht als Vorgang.
