# CR-Vorschlag: Homebrew-Formeln im Devcontainer — neue Anforderung im Lastenheft — Alt/Neu je Stelle

**Status:** freigegeben wie vorgelegt (Projektinhaber, 2026-11-05); F1 Priorität `V2` und F2 Lastenheft-Version `0.36.0` wie vorgeschlagen.

**Gegenstand:** die Lieferung zur DoD-Zeile 1 des CR-Slices [`slice-cr-devcontainer-brew-pakete`](../planning/in-progress/slice-cr-devcontainer-brew-pakete.md) (Lastenheft-Wortlaut); Umsetzung im Folge-Slice [`slice-v1-devcontainer-homebrew-pakete`](../planning/next/slice-v1-devcontainer-homebrew-pakete.md).

**Grundlage:** die zwei Entscheidungen des Projektinhabers vom 2026-11-05 (Homebrew-Formeln als frei wählbare Liste im Devcontainer-Image; Formeln aus einem Tap nur mit ausdrücklicher Freigabe) und [`ADR-0046`](../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (`Accepted`). Das Lastenheft nimmt nur die Anforderung auf. Schlüsselnamen, Pfade der Konfiguration, das Verfahren der Freigabe, die Pin-Stufe von Homebrew sowie Bau-, Plattform- und Messgrenzen der Zusage gehören in die Spezifikation und in die Tests des Umsetzungs-Slices.

---

## Stelle (a) — neue Anforderung im Abschnitt der Devcontainer-Anforderungen

**Alt:** `spec/lastenheft.md` Zeile 531: [`LH-FA-DEV-012`](../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) ist die letzte Anforderung des Abschnitts; danach folgt die Trennlinie `---`. Eine Anforderung zu Homebrew-Formeln gibt es nicht.

**Neu (Vorschlag):** eine neue Anforderung nach `LH-FA-DEV-012` (die Kennung vergibt die Übernahme, nächste freie: `LH-FA-DEV-013`): <!-- d-check:ignore (Kennung der neuen Anforderung, entsteht mit der Übernahme) -->

> ### LH-FA-DEV-013 – Homebrew-Formeln im Devcontainer <!-- d-check:ignore (Kennung der neuen Anforderung, entsteht mit der Übernahme) -->
>
> *Priorität: V2.* Das Produkt soll Homebrew-Formeln, die der Benutzer in der Projektkonfiguration als Liste führt, im erzeugten Devcontainer-Image installieren, ohne dass `u-boot` die Werkzeuge einzeln kennen muss.
>
> - Ein Eintrag der Liste ist eine Formel oder eine Formel aus einem Tap (`<user>/<tap>/<formel>`).
> - Ein Tap ist eine externe Quelle ([`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)): Eine Formel aus einem Tap wird nur nach ausdrücklicher Freigabe des Benutzers installiert; `--yes` ersetzt sie nicht. Ohne Freigabe endet `u-boot generate devcontainer` mit einem fachlichen Fehler (Exit-Code `10`), ehe etwas geschrieben wird.
> - Die Formeln sind nicht gepinnt; `u-boot` weist das beim Erzeugen aus.

## Stelle (b) — Versionsfeld und Historie (Übernahme)

**Alt:** `spec/lastenheft.md` Zeile 9: Version `0.35.0`; §16 Historie endet mit der Zeile `0.35.0` (2026-11-04).

**Neu (Vorschlag):** Version `0.36.0` und eine Historie-Zeile mit dem Datum der Übernahme: neue Anforderung Homebrew-Formeln im Devcontainer (Liste in der Projektkonfiguration, Tap nur mit ausdrücklicher Freigabe, Formeln ungepinnt und ausgewiesen); Entscheidung des Projektinhabers 2026-11-05.

## Nicht geändert

- [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CONF-002`](../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) und [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte): Die Anforderung steht für sich. Katalogschlüssel, Konfigurationsschlüssel, die Freigabeliste für Taps und die Regel zur Kombination mit dem Katalog-Feature `homebrew` (Exit-Code 10 bei gleichzeitiger Aktivierung) stehen in der Spezifikation und präzisieren die neue Anforderung, ohne sie zu erweitern.
- Bau, Betrieb, Plattformen und die Grenzen der Zusage: keine Lastenheft-Aussage. Sie stehen in der Spezifikation und in den Tests des Umsetzungs-Slices, mit ihrer Messgrundlage.
- `devcontainer down`, `up --recreate` und das Community-Feature `homebrew`: unverändert.

---

**Abweichungen der Übernahme je Stelle:** *(bei der Übernahme zu tragen.)*

**Offen zur Freigabe:**

- **F1 — Priorität.** Vorschlag `V2` (wie [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und [`LH-FA-DEV-012`](../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers)). Gegenposition: `V1`.
- **F2 — Version des Lastenhefts.** `0.36.0` ist der nächste freie Schritt der Lastenheft-Version (nicht die Release-Version); die Zahl vergibt die Übernahme.

**Regeln dieser Form:** die Kennungen sind verlinkt; der Neu-Wortlaut nennt Anforderungen, keine Plattform-, Bau- oder Messgrenzen und keine Konfigurationsschlüssel; die Historie-Zeile entsteht bei der Übernahme mit dem Datum des Commits.
