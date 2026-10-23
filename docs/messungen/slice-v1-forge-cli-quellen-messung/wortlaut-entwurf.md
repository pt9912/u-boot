# Wortlaut-Entwurf für den Lastenheft-CR: `gh` und `glab` als Katalog-Features

Entwurf des Slice [`slice-v1-forge-cli-quellen-messung`](../../plan/planning/in-progress/slice-v1-forge-cli-quellen-messung.md); **kein Eintrag im Lastenheft**. Der Projektinhaber gibt
den Wortlaut frei; erst danach entsteht der CR-Slice der Welle ([`welle-forge-cli-gh-glab`](../../plan/planning/welle-forge-cli-gh-glab.md) §4), der das Lastenheft auf die nächste freie Minor-Version ändert
(gelesen am 2026-10-23: `spec/lastenheft.md` trägt `0.21.0`; nächste freie Version `0.22.0`). Messwerte und Rohdaten: [`README.md`](README.md) (Aussagen A bis F).

**Voraussetzung dieses Entwurfs:** Der Wortlaut folgt der Quellenwahl, die [`ADR-0035`](../../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Proposed`) vorschlägt: je Werkzeug ein gepinntes Release-Archiv
des Upstream-Projekts, im erzeugten Dockerfile mit geprüfter SHA-256-Prüfsumme (Option B). Entscheidet der Projektinhaber anders (Freigabe-Punkte a und b der ADR), ändern sich die Sätze S2 bis S4; die
übrigen Sätze bleiben.

**Konvention** (wie [`slice-cr-devcontainer-ollama`](../../plan/planning/done/slice-cr-devcontainer-ollama.md)): Das Lastenheft nennt keine Slice-, Wellen- oder ADR-Kennung; Aussagen über das Verhalten des Generators sind Zielbild, das der
Code-Slice einlöst, und tragen **keinen Messwert**. Die Spalte *Nicht gemessen* nennt je Satz, was der Satz mehr sagt, als die Messung trägt.

## Stelle 1 — [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beispiele und neuer Absatz

Alt (Spiegelstriche, bleiben):

> - Git
> - Docker CLI
> - Node.js
> - Java
> - SDKMAN
> - PostgreSQL Client
> - Kubernetes Tools
> - Ollama (lokale Sprachmodelle, nur CPU)

Neu (zwei zusätzliche Spiegelstriche am Ende der Liste):

> - GitHub CLI (`gh`)
> - GitLab CLI (`glab`)

Neu (zusätzlicher Absatz hinter dem Ollama-Absatz, vor „Für optionale externe Feature-Quellen gilt:“):

> `gh` (GitHub CLI) und `glab` (GitLab CLI) sind als eingebaute Features `gh` und `glab` wählbar (`u-boot config set devcontainer.features.gh.enabled true`, `u-boot config set devcontainer.features.glab.enabled true`). [S1] Jedes installiert sein Werkzeug im Container aus einem gepinnten Release-Archiv des Upstream-Projekts (`cli/cli`, `gitlab-org/cli`) mit geprüfter SHA-256-Prüfsumme. [S2] Die Archive sind nicht signiert; für `gh` veröffentlicht das Projekt eine GitHub-Attestierung, die `u-boot` nicht prüft. `u-boot` weist das beim Erzeugen aus. [S3] Mit der Egress-Restriktion ([`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) stehen die Hosts von `gh` (`api.github.com`, `github.com`) in der Default-Allowlist und kommt für `glab` `gitlab.com` hinzu. [S4] Anmeldung und Zugangsdaten: siehe [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer); `u-boot` führt weder `gh auth login` noch `glab auth login` aus und speichert kein Token. [S5] Gemessen ist das auf Docker 29.9.2, Linux `amd64`, `gh` 2.103.0 und `glab` 1.124.0; nicht zugesagt sind Aufrufe, die ein gültiges Token brauchen und deren Hosts, Telemetrie- und Update-Hosts einer interaktiven Sitzung, selbst gehostete Instanzen (GitHub Enterprise, GitLab Self-Managed), andere Architekturen und Container-Engines als gemessen sowie andere Versionen der Werkzeuge als die gepinnten. [S6]

(Die Marken `[S1]` bis `[S6]` stehen nur in diesem Entwurf, nicht im Lastenheft.)

| Satz | Beleg | Nicht gemessen |
|---|---|---|
| S1 — wählbare Katalog-Features, Konfigurationsschlüssel | Zielbild des Code-Slice; Katalog-Mechanismus für Schlüssel ohne OCI-Quelle ist der des Features `ollama` (Typ „Dockerfile-Block“, gelesen in `featureCatalogueEntry`) | der Katalogeintrag und der Dockerfile-Block für `gh`/`glab` existieren nicht; **ob** der Block mit und ohne Sandbox-Profil wirkt (im Dockerfile steht `USER root` heute nur bei Sandbox-Profil oder UID-Überschreibung, `Dockerfile.tmpl`) ist nicht gemessen; gemessen sind nur die Blöcke als Handvariante in einem Sandbox-Dockerfile (`ghtgz`, `glabtgz`, Aussage C10) |
| S2 — gepinntes Archiv, SHA-256 vor dem Entpacken geprüft | C10 (Bau OK, Zeile `…tar.gz: OK`), C11 (verfälschter Pin bricht den Bau ab, kein Image) | `arm64` (der Pin je Architektur ist aus den Prüfsummen-Dateien gelesen, nicht gebaut); Verhalten bei einem Release, das nachträglich ersetzt wurde (`gh`: Release `immutable`, A10; `glab`: nicht untersucht) |
| S3 — nicht signiert; `gh`-Attestierung nicht von `u-boot` geprüft; Ausweisung | A10 (keine Signatur-Assets, zwei Attestierungen über die API, Inhalt nicht dekodiert), A12 (`glab`: keine Signatur-/Provenance-Assets, Signierschritte nur für macOS/Windows nach Suchmuster) | ob `gh attestation verify` das Archiv bestätigt (nicht gemessen); ob `glab` außerhalb des Suchmusters signiert; die Ausweisung beim Erzeugen ist Zielbild |
| S4 — Hosts | D1 (`gh`: `api.github.com`, `github.com`, beide in der Basis-Liste B6), D2/D3 (`glab`: nur `gitlab.com`, nicht in der Basis-Liste; ohne Eintrag scheitern die Aufrufe, mit Eintrag gelingen sie) | alle Hosts eines Aufrufs mit gültigem Token (`gh release download` leitet auf weitere Hosts, nicht gemessen), Telemetrie- und Update-Hosts (D5), Hosts der Anmeldung; die Eintragung in die Allowlist ist Zielbild |
| S5 — keine Anmeldung durch `u-boot`, kein gespeichertes Token | Zusage über das Produkt, nicht über die Werkzeuge: E4 (der Platzhalter-Wert steht in keiner Datei unter `~/.config`, `~/.cache`, `~/.local`, `~/.gitconfig`; Umgebungsvariable genügt dem Werkzeug als Quelle, E2) | was die Werkzeuge nach einem **eigenen** `auth login` des Nutzers speichern (Schlüsselbund oder `hosts.yml`); Verhalten mit gültigem Token |
| S6 — Messbereich und nicht Zugesagtes | Messstand und F | — (der Satz benennt selbst die Grenze) |

## Stelle 2 — [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Zugangsdaten eines Forge-Werkzeugs

Alt (Einleitungssatz):

> *Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten (etwa der Zugangsschlüssel eines Coding-Agenten, den ein aktiviertes Devcontainer-Feature nutzt) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.

Neu:

> *Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten oder eines Forge-Werkzeugs (etwa der Zugangsschlüssel eines Coding-Agenten oder das Token für `gh` oder `glab`, die ein aktiviertes Devcontainer-Feature nutzt) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.

Alt (zweiter Spiegelstrich):

> - Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote; für den Agenten nennt sie die Umgebungsvariable, aus der der Zugangsschlüssel stammt.

Neu:

> - Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote; für den Agenten und für `gh` und `glab` nennt sie die Umgebungsvariable, aus der der Zugangsschlüssel beziehungsweise das Token stammt (`GH_TOKEN` für `gh`, `GITLAB_TOKEN` für `glab`).

Alt (dritter Spiegelstrich):

> - `u-boot doctor` prüft die Quelle der Zugangsdaten (Git-Token-Quelle, bei aktivem Agenten-Feature die Umgebungsvariable des Agenten): fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

Neu:

> - `u-boot doctor` prüft die Quelle der Zugangsdaten (Git-Token-Quelle, bei aktivem Agenten-Feature die Umgebungsvariable des Agenten, bei aktivem `gh`- oder `glab`-Feature deren Token-Variable): fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

| Satz | Beleg | Nicht gemessen |
|---|---|---|
| Zugangsdaten von `gh`/`glab` nicht in Image und Dateien | E4 (nur Umgebungsvariablen; keine Datei trägt den Platzhalter), B5 (die erzeugte `devcontainer.json` trägt `GIT_TOKEN` als Referenz, keinen Wert) | Verhalten der Werkzeuge mit gültigem Token; die Referenz `${localEnv:GH_TOKEN}`/`${localEnv:GITLAB_TOKEN}` in `remoteEnv` ist Zielbild (Code-Slice) |
| `GH_TOKEN` (`gh`), `GITLAB_TOKEN` (`glab`) als die Variablen der Dokumentation | E1 (Hilfe der Werkzeuge nennt beide an erster Stelle der Rangfolge), E2 (`401` mit der Variable gesetzt: das Werkzeug verwendet sie ohne `auth login`) | dass dieselbe Anfrage mit gültigem Token gelingt; `GITHUB_TOKEN` und `GITLAB_ACCESS_TOKEN`/`OAUTH_TOKEN` sind ebenfalls wirksam (E2), die Wahl der einen Variable je Werkzeug ist Entscheidung des Projektinhabers ([`ADR-0035`](../../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md), Punkt c) |
| `doctor` prüft die Variable | Zielbild (Code-Slice); Muster `agentCredentialEnv`/`checkSandboxCredentials` gelesen (`devcontainer_sandbox.go`, `doctor_sandbox.go`) | alles am Prüfverhalten für `gh`/`glab` |

## Stelle 3 — keine Änderung vorgeschlagen

- [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion): Die Hosts je Feature stehen in der Spezifikation ([`SPEC-019`](../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)), nicht im Lastenheft (Muster `ollama`); S4 nennt sie im Absatz zu [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nur, soweit der Projektinhaber die Nennung wünscht — sonst entfällt der zweite Teil von S4.
- [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte): keine Änderung. Ein eingebautes Katalog-Feature ist hinterlegt; der Download eines Release-Archivs im Dockerfile steht unter derselben Auslegung wie bei `ollama` (Archiv mit geprüfter Prüfsumme, Grenze ausgewiesen).
- [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): keine Änderung; beide Werkzeuge binden keinen Socket und fügen keine Lockerung hinzu (Metadaten von Features ohne Eintrag sind nicht anwendbar; ein Dockerfile-Block hat keine Metadaten).

## Stelle 4 — Version

Alt: `| Version          | 0.21.0 |`. Neu: `| Version          | 0.22.0 |` und eine Zeile in der Historie des Lastenhefts (Lastenheft-Historie führt der CR-Slice).

## Offene Punkte für die Freigabe

Die Punkte stehen mit Empfehlung in [`ADR-0035`](../../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (Freigabe-Punkte a bis g); der Wortlaut hängt an a (Quelle von `gh`), b (Quelle von `glab`), c (Token-Variable), d (Wirkung mit und ohne Sandbox-Profil).
