# Inline-`nolint` im Bestand gegen das Suppression-Verbot

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — Das Suppression-Verbot ([`AGENTS.md`](../../../../../../AGENTS.md) §3.2) ist eine repo-weite Konvention des Werkzeug-Verbunds, kein Teil einer einzelnen Go-Sub-Area; die Lücke ist ein fehlender Sensor im Lint-Profil (`nolintlint` ist in `.golangci.yml` nicht aktiviert). Die Treffer liegen in zwei Go-Sub-Areas (APP, CLI); eine Aufteilung nach Fundort ergäbe zwei Pfade für dieselbe Beobachtung.

Im Quellbaum stehen Inline-`nolint`-Kommentare, obwohl das Verbot sie ausschließt. Kein Gate prüft diese Zeilen: Ein Kommentar mit Leerzeichen (`// nolint:…`) ist keine Direktive und wirkt nicht, ein wirksamer unterdrückt einen Befund, ohne dass `make lint` es meldet. Die Drift bleibt unsichtbar, bis ein Review den Baum durchsucht.
