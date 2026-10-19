# Messprotokoll: Egress-Allowlist, Subdomains und Wildcard

Rohdaten der Messung des Slice
[`slice-v1-egress-wildcard-pruefung`](../../plan/planning/done/slice-v1-egress-wildcard-pruefung.md).
Messstand, Auswertung und Grenzen stehen in der Evidence-Datei
[`slice-v1-egress-wildcard-pruefung`](../../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md);
diese Datei wiederholt sie nicht. Skripte liegen unter `tools/messung/egress-wildcard-pruefung/`
(Ablage: [`docs/messungen/README.md`](../README.md)).

## Aufruf

`run.sh <repo-wurzel> [fall...]` im Verzeichnis `tools/messung/egress-wildcard-pruefung/`, mit Fall aus `a a2 b c d d2 d3 cd-c cd-d`
(ohne Angabe: `a a2 b c d2 d3 cd-c cd-d`). Voraussetzung: Docker mit Zugriff auf `debian:stable` und `--cap-add=NET_ADMIN`.
`run.sh` ruft `inner.sh` (Startscript `egress-init.sh.tmpl` mit ersetzter `hosts=`-Zeile) und `inner-cd.sh` (Messungen ohne Allowlist)
in Wegwerf-Containern. Die Ausgabe wird in die Rohdatei in diesem Verzeichnis umgeleitet.

## Rohdaten

| Datei | Inhalt (laut Evidence-Datei) |
|---|---|
| `rohdaten-1.txt` | Fälle a, a2, b |
| `rohdaten-2.txt` | CNAME-Kandidaten ohne Allowlist, Muster `*.` und `.` ohne Startscript |
| `rohdaten-3.txt` | Fälle c, d2, d3 |
