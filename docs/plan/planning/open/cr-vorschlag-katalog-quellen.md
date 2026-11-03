# CR-Vorschlag: Katalogquellen `cpp`/`postgres-client` entfernen und unverarbeitbare Quellen-Formen abweisen — Alt/Neu je Stelle

**Status:** vorgelegt (2026-11-03); wartet auf die Freigabe des Projektinhabers.

**Gegenstand:** die Lieferung zu DoD-Zeile 1 des CR-Slices
[`slice-cr-devcontainer-katalog-quellen`](slice-cr-devcontainer-katalog-quellen.md).
Mess-Basis: [`docs/messungen/slice-v1-katalog-features-quellen-und-freigabe/`](../../../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md)
(M1: die Quellen `ghcr.io/devcontainers/features/{cpp,postgresql-client}` existieren
nicht — Upstream 29 Features ohne beide, Token-Verweigerung 403 DENIED;
M2: die CLI `0.89.0` verarbeitet `oci://…:<tag>` und `https://…:<tag>` nicht).
Entscheidungen des Projektinhabers 2026-11-03: **N1** Einträge entfernen (mit
CR für die Lastenheft-Beispiel-Zeile), **N2** Abweisung unverarbeitbarer
Formen (Exit-Code `10`, Abhilfe in der Meldung).

## Stelle (a) — Lastenheft [`LH-FA-DEV-003`](../../../../spec/lastenheft.md): Beispiele-Zeile

**Alt:** die Beispiele-Zeile nennt „PostgreSQL Client“ (Zeile ~405).
**Neu:** die Zeile ohne „PostgreSQL Client“; keine andere Zusage-Änderung.

## Stelle (b) — Lastenheft [`LH-FA-DEV-003`](../../../../spec/lastenheft.md): Abweisungs-Sätze (externer Quellen)

**Neu (Vorschlag):** „Eine externe Quell-Angabe, deren Form die Dev-Container-CLI
nicht verarbeitet (gemessen am 2026-11-03 mit der CLI des Repos: `oci://…:<tag>`
und `https://…:<tag>` unverarbeitbar, Registry-Tarball-Form offen), beendet
`u-boot generate devcontainer` mit Exit-Code `10`, ehe nichts geschrieben wird,
und die Meldung nennt die Abhilfe (die verarbeitbare Registry-Form); `u-boot
config set …source` wendet dieselbe Prüfung an ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)). Der Tarball-Weg der CLI bleibt Offenpunkt; sein Messergebnis kommt in die Historie.“

## Stelle (c) — Technik

[`SPEC-023`](../../../../spec/spezifikation.md): die Katalogtabelle verliert die
Zeilen `cpp`/`postgresql-client`; die Quelle-Form-Aufzählung trägt die
gemessene Verarbeitungstabelle (Registry-Form verarbeitet; `oci://`/`https://`
nicht; Tarball offen). [`SPEC-040`](../../../../spec/spezifikation.md): der
Schreib-Pipeline-Schritt 6 (Fachprüfung) trägt die Form-Prüfung der `.source`-Werte.
Beleg-Aufträge statt Test-Namen.

## Stelle (d) — Historie

`| 0.33.0 | 2026-11-03 | Entfernung der Katalog-Einträge `cpp`/`postgres-client` (Quellen existieren nicht, gemessen) und die Abweisung unverarbeitbarer Formen externer Quellen je Exit-Code `10` (gemessen mit der CLI des Repos) |`
