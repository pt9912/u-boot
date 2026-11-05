# CR-Vorschlag: Homebrew-Formeln im Devcontainer — Schlüssel `brew` mit `packages` und Tap-Freigabeliste, Lastenheft `LH-FA-DEV-003` und `LH-FA-CONF-002` — Alt/Neu je Stelle

**Status:** vorgelegt (2026-11-05); wartet auf die Freigabe des Projektinhabers.

**Gegenstand:** die Lieferung zur DoD-Zeile 1 des CR-Slices [`slice-cr-devcontainer-brew-pakete`](../planning/in-progress/slice-cr-devcontainer-brew-pakete.md) (Lastenheft-Wortlaut); Umsetzung im Folge-Slice [`slice-v1-devcontainer-homebrew-pakete`](../planning/next/slice-v1-devcontainer-homebrew-pakete.md).

**Grundlage:** [`ADR-0046`](../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (`Accepted`, 2026-11-05: Schlüssel `brew`, Freigabeliste `devcontainer.brewTaps.allow`, Commit-Pin ohne Übersteuerung) und die zwei Entscheidungen des Projektinhabers vom 2026-11-05 (neuer Schlüssel neben dem OCI-Feature `homebrew`; Taps nur über eine eigene Freigabeliste). **Messbasis:** keine; der Block ist noch nicht gebaut. Jede Aussage über Bau, Plattform und Netz steht deshalb als Zielbild mit ausgewiesener Grenze, nicht als gemessene Zusage (Muster der Instanz-Zusage; die Messung folgt in DoD 3 des Umsetzungs-Slices).

---

## Stelle (a) — [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beispielliste

**Alt:** `spec/lastenheft.md` Zeile 418, letzter Punkt der Liste:

> - Anschluss an einen Ollama-Dienst des Hosts (`ollama-host`)

**Neu (Vorschlag):** ein weiterer Punkt darunter:

> - Homebrew-Formeln im Image (`brew`)

## Stelle (b) — [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), neuer Absatz nach dem Absatz zu `ollama-host`

**Alt:** `spec/lastenheft.md` Zeile 430 ist der Absatz „Auf Wunsch des Nutzers erzeugt `u-boot` einen Anschluss des Containers an einen Ollama-Dienst des Hosts …“; Zeile 432 beginnt mit „Für optionale externe Feature-Quellen gilt:“. Dazwischen steht nichts.

**Neu (Vorschlag):** ein neuer Absatz zwischen beiden:

> `brew` (Homebrew-Formeln im Image) ist in beiden Profilen (Standard und Sandbox) als eingebautes Feature wählbar (`u-boot config set devcontainer.features.brew.enabled true`) und nimmt die Liste `devcontainer.features.brew.packages` entgegen. Der Schlüssel installiert im erzeugten Devcontainer-Dockerfile Homebrew aus einem gepinnten Commit des Upstream-Projekts, ohne Installer-Skript, und danach die Formeln der Liste. Die Formeln selbst sind nicht gepinnt: Ihre Version bestimmt Homebrew zum Zeitpunkt des Baus, und der Pin deckt nur Homebrew. Ein Eintrag der Liste ist `<formel>` oder `<user>/<tap>/<formel>`. `u-boot` weist beim Erzeugen den Block, die Quelle von Homebrew, die ungepinnte Formelliste und jeden Tap aus. Das Community-Feature `homebrew` bleibt unverändert wählbar; beide zusammen aktiv sind ein fachlicher Validierungsfehler (Exit-Code `10`). Zugesagt ist die Erzeugung des Dockerfile-Blocks und der Konfiguration; der Bau des Images und der Betrieb sind nur auf einem Linux-Host (`amd64`) mit lokaler, nicht verschachtelter Docker-Engine und nur soweit zugesagt, wie ein Test mit Build-Tag `docker` es belegt. Nicht zugesagt: Version, Verfügbarkeit und Arbeitsverhalten einer Formel; die Hosts, die eine Formel beim Bau oder im Betrieb braucht (hinter der Egress-Restriktion trägt der Nutzer sie in `devcontainer.sandbox.egress.allow` ein, `u-boot` leitet sie aus Paketnamen nicht ab); Linux-`arm64`; das Überstehen eines Neubaus des Images, soweit kein Test es belegt.

## Stelle (c) — [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Freigabe externer Quellen: Taps

**Alt:** `spec/lastenheft.md` Zeilen 432 bis 440, Liste „Für optionale externe Feature-Quellen gilt:“, endet mit den Punkten „`--yes` allein gilt nicht als Zustimmung für externe Quellen.“ (Zeile 439) und der Quell-Form-Regel (Zeile 440). Ein Tap kommt nicht vor.

**Neu (Vorschlag):** ein weiterer Punkt am Ende der Liste:

> - Ein Tap in `devcontainer.features.brew.packages` (Eintrag der Form `<user>/<tap>/<formel>`) ist eine eigene externe Quelle: Er braucht einen Eintrag `<user>/<tap>` in `devcontainer.brewTaps.allow`. Die Freigabe erfolgt über `u-boot config set devcontainer.brewTaps.allow <user>/<tap>[,<user>/<tap>...]`; `--allow-external-feature-sources` gilt dafür nicht, und `--yes` ersetzt die Freigabe nicht. Ohne Eintrag endet `u-boot generate devcontainer` mit Exit-Code `10`, ehe etwas geschrieben wird.

## Stelle (d) — [`LH-FA-CONF-002`](../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), Hinweise zur Konfiguration

**Alt:** `spec/lastenheft.md` Zeilen 868 bis 870: „`devcontainer.featureSources.allow` ist optional; fehlt das Feld, ist die Liste leer …“, „Erlaubte Einträge in `devcontainer.featureSources.allow` müssen gültige, non-empty Quell-Strings … sein.“, „Beim Schreiben wird die Liste dedupliziert.“

**Neu (Vorschlag):** zwei weitere Punkte nach Zeile 870:

> - `devcontainer.brewTaps.allow` ist optional; fehlt das Feld, ist die Liste leer (kein Tap freigegeben). Einträge haben die Form `<user>/<tap>`; beim Schreiben wird die Liste dedupliziert. Ein ungültiger Eintrag ist ein fachlicher Validierungsfehler (Exit-Code `10`).
> - `devcontainer.features.brew.packages` ist optional (Liste von Zeichenketten der Form `<formel>` oder `<user>/<tap>/<formel>`) und nur am Schlüssel `brew` ohne eigene `source` zulässig; ein ungültiger Eintrag oder eine Liste an einem anderen Schlüssel ist ein fachlicher Validierungsfehler (Exit-Code `10`).

## Stelle (e) — Versionsfeld und Historie (Übernahme)

**Alt:** `spec/lastenheft.md` Zeile 9: Version `0.35.0`; §16 Historie endet mit der Zeile `0.35.0` (2026-11-04).

**Neu (Vorschlag):** Version `0.36.0` und eine Historie-Zeile mit dem Datum der Übernahme, die die Stellen (a) bis (d) in einem Satz nennt (Schlüssel `brew` als Dockerfile-Block mit ungepinnter Formelliste, Tap-Freigabe `devcontainer.brewTaps.allow`, Konflikt mit dem Feature `homebrew`, Zusage der Erzeugung, Bau nur nach Test mit Build-Tag `docker`) und mit „Entscheidung des Projektinhabers 2026-11-05“.

## Nicht geändert

- [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte): Der Wortlaut bleibt. Die Auslegung (ein Tap ist eine vom Nutzer genannte Quelle, die Ausnahme für Paketmanager-Pakete deckt, was der Paketmanager mitbringt) steht in [`ADR-0046`](../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md).
- `devcontainer down`, `up --recreate` und der Katalog-Eintrag `homebrew`: unverändert.

---

**Abweichungen der Übernahme je Stelle:** *(nach der Freigabe zu tragen.)*

**Offen zur Freigabe:**

- **F1 — Zielbild vor Messung (Stelle (b)).** Der Satz zur Bau-Zusage (Linux-`amd64`, nur soweit ein Docker-Test es belegt) steht, bevor der Block gebaut ist. Gegenposition: den Satz streichen und mit der Messung nachziehen (eigene Lastenheft-Änderung). Empfehlung: lassen, weil er die Zusage an den Test bindet und ohne Test nichts zusagt.
- **F2 — kein Flag für die Tap-Freigabe (Stelle (c)).** Der Weg ist `config set`; ein Flag wie bei `--allow-external-feature-sources` gibt es nicht. Gegenposition: ein Flag nur für `generate devcontainer`. Empfehlung: kein Flag, weil jedes Flag eine Prüfstrecke an drei Befehlen kostet.
- **F3 — Versionsnummer.** `0.36.0` ist der nächste freie Minor-Schritt; die Zahl vergibt die Übernahme.

**Regeln dieser Form:** die Kennungen sind verlinkt; der Neu-Wortlaut trägt Zielbild nur mit ausgewiesener Grenze; die Historie-Zeile entsteht bei der Übernahme mit dem Datum des Commits.
