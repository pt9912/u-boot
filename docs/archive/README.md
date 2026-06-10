# docs/archive

Abgelöste oder veraltete Doku-Inhalte aus `docs/user/`, `docs/plan/` oder
anderen Bereichen. Statt zu löschen wird hierher verschoben, damit die
Historie auffindbar bleibt (siehe
[`LH-FA-PROJDOCS-004`](../../spec/lastenheft.md#lh-fa-projdocs-004-archivierung)).

Konventionen:

- Verschieben per `git mv`, damit die Datei-Historie erhalten bleibt.
- Am Anfang des verschobenen Dokuments einen Hinweis ergänzen, z. B.:

  ```markdown
  > Archiviert am 2026-05-21; ersetzt durch [docs/user/neue-anleitung.md](../user/neue-anleitung.md).
  ```

- Querverweise in lebendiger Doku auf das neue Ziel umbiegen.

## Archivierte Dokumente

| Datei | Inhalt |
| --- | --- |
| [`roadmap-history-v0.1-v0.3.md`](roadmap-history-v0.1-v0.3.md) | Ausgelagerte Release- und Roadmap-Historie fuer v0.1.0 bis v0.3.0. |
