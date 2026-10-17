# Host-Werkzeug außerhalb Docker

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Docker-only-Regel ([`AGENTS.md`](../../../../../../AGENTS.md) §3.1) ist die Konvention des Werkzeug-Verbunds; ein Aufruf ohne `make` umgeht ihn.

Ein Agent führt ein Host-Werkzeug (`python3`, Host-Go-Toolchain) aus, statt über `make` im Docker zu arbeiten. Eine Guard-Sperre für den Aufruf fehlt oder greift nicht; der Verstoß fällt erst im Bericht des Implementers oder im Review auf. Der Nachweis der Änderung läuft dann nur über `make test` und `make gates`.
