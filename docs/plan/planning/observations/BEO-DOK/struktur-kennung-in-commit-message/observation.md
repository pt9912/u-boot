# Struktur-Kennung in der Commit-Message ohne Sensor

**Sub-Area:** `harness/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — [`AGENTS.md`](../../../../../../AGENTS.md) §5 Regel 1 und das Commit-Gate sind die Konvention und ihr Sensor.

Eine Commit-Message nennt eine Struktur-Kennung (`SPEC-<NNN>` oder `ARC-<NNN>`), obwohl `AGENTS.md` §5 Regel 1 sie aus der Commit-Message hält: Sie adressiert nur innerhalb der Spec und wandert mit deren Umbau. Das Commit-Gate prüft die Anwesenheit einer Anforderungs-, ADR-, MR- oder Slice-Kennung (`tools/harness/commit-msg-traceability.sh`, `make doc-commits`), nicht die Abwesenheit einer Struktur-Kennung; die Message geht mit Exit 0 durch. Gefunden hat es der Review, nicht ein Gate.
