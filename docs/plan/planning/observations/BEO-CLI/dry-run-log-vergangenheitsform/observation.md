# Dry-Run: Logzeilen auf stderr stehen in Vergangenheitsform

**Sub-Area:** `internal/adapter/driving/cli` (Kürzel CLI, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Ausgabe-Zusagen des Dry-Runs ([`LH-FA-CLI-007`](../../../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)) sind dort Konvention, die Logzeilen entstehen in den Use-Cases der Anwendungsschicht.

Im Dry-Run schreiben die Use-Cases weiterhin `level=INFO`-Logzeilen auf stderr, die ein Ergebnis in Vergangenheitsform melden, das nicht eingetreten ist (`config set: updated`, `generate: updated block`, `generate devcontainer: created`). Die Zusammenfassung auf stdout sagt seit [`slice-v1-cli-dry-run-would-formulierung`](../../../done/slice-v1-cli-dry-run-would-formulierung.md) „Would …“; die Logzeilen sind nicht Teil jenes Slice.
