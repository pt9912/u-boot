# internal/

Nicht-exportierbare Go-Pakete für u-boot (Coverage-Scope, siehe
`LH-FA-BUILD-009` in [../spec/lastenheft.md](../spec/lastenheft.md)).

Im MVP-Bootstrap noch leer. Erste Pakete entstehen mit dem ersten
fachlichen Slice (voraussichtlich `internal/cli/`, `internal/config/`,
`internal/project/`).

Solange `./internal/...` keinen produktiven Code enthält, läuft das
Coverage-Gate im Bootstrap-Modus mit Schwellwert `0`
(`LH-FA-BUILD-008`).
