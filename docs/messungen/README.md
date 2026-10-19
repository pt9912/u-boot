# docs/messungen

Messprotokolle: das Dokument einer Messung und ihre Rohdaten. Eine Messung, die ein Slice oder ein ADR als Beleg nennt,
liegt hier im Repo, nicht im Home-Verzeichnis des Messenden.

## Ablage

Die Regel steht in [`harness/conventions.md`](../../harness/conventions.md#ablage-von-messprotokollen-und-skripten); hier ihre Anwendung auf dieses Verzeichnis.

- **Protokoll hier:** ein Verzeichnis je Vorgang unter `docs/messungen/` mit `README.md` (Messstand, Reihenfolge und Aufruf, Lesehilfen)
  und den Rohdaten als `.txt` oder `.log`. Der Verzeichnisname ist die Kennung des Vorgangs, der die Messung auslöste.
- **Skripte und Hilfs-Dockerfiles unter `tools/messung/<name>/`**, nicht hier: Ein Skript und seine Hilfsdateien
  (`Dockerfile.*`, `inner-*.sh`, Sonden) liegen im selben Verzeichnis, weil das Skript sie über sein eigenes Verzeichnis findet.
  Das `README.md` des Protokolls nennt den Skriptpfad und den Aufruf, der die Ausgabe in dieses Verzeichnis umleitet.
- **Im Beobachtungs-Register** bleibt je Beobachtung nur die Form `observation.md`, `state.md`, `evidence/<vorgangs-id>.md`;
  eine Evidence-Datei verweist auf das Protokoll hier.

## Bestand

| Verzeichnis | Skripte |
|---|---|
| [`slice-v1-devcontainer-ollama-messung`](slice-v1-devcontainer-ollama-messung/README.md) | `tools/messung/devcontainer-ollama/` |
| [`slice-v1-sandbox-egress-neuladen-entwurf`](slice-v1-sandbox-egress-neuladen-entwurf/README.md) | `tools/messung/sandbox-egress-neuladen-entwurf/` |
| [`slice-v1-egress-wildcard-pruefung`](slice-v1-egress-wildcard-pruefung/README.md) | `tools/messung/egress-wildcard-pruefung/` |
| [`slice-v1-devcontainer-ollama-amd64-messung`](slice-v1-devcontainer-ollama-amd64-messung/README.md) | `tools/messung/ollama-amd64/` |
| [`slice-v1-egress-dns-upstream-schliessen`](slice-v1-egress-dns-upstream-schliessen/README.md) | `tools/messung/egress-dns-upstream/` |
