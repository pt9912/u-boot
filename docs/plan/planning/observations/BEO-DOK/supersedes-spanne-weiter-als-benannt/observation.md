# Eine ADR mit `Supersedes (teilweise)` ersetzt mehr, als ihr Kopf benennt

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Kopf einer ADR nennt die ersetzten Stellen, die Tabelle der geltenden Fassung gibt den neuen Wortlaut; nach `Accepted` gilt nur eine Folge-ADR ([`AGENTS.md`](../../../../../../AGENTS.md) §3.5).

Eine ADR, die eine ältere nur teilweise ersetzt, nennt im Kopf (`Supersedes (teilweise)`) für eine Stelle einen Satzteil, und die Tabelle der „geltenden Fassung“ schreibt mehr um als diesen Satzteil: den Schluss des Satzes oder eine neue Klammer dahinter. Wer die ältere ADR liest und nur den benannten Teil ersetzt, behält einen Rest, den die neue ADR meint zu ersetzen. Kein Gate prüft den Kopf gegen die Tabelle; der Nachlauf-Review fand es.

**Trigger:** ein zweiter Vorgang mit einer ADR, deren Tabelle mehr ersetzt als der Kopf benennt, oder die Entscheidung des Projektinhabers, den Abgleich von Kopf und Tabelle in die ADR-Form ([`MR-008`](../../../../../../harness/conventions/MR-008-adr-form.md)) zu schreiben.
