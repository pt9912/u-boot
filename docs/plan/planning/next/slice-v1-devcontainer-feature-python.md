# Slice V1: Devcontainer-Feature `python` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.9.0 (Devcontainer-Features Homebrew, SDKMAN, Python).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog-Regel: offizielle Einträge brauchen das Präfix `ghcr.io/devcontainers/features/`).

**Berührte Spec-Stellen:** `spezifikation.md` [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (eingebaute Features) und [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Egress-Hosts je Feature); das Lastenheft bleibt unverändert (Beispielliste, keine Zusage zu einzelnen Schlüsseln).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `python` ist ein eingebautes Devcontainer-Feature und rendert das offizielle Feature `ghcr.io/devcontainers/features/python:1`; die Egress-Allowlist trägt `pypi.org` und `files.pythonhosted.org`, wenn es aktiv ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Feature-Optionen (Python-Version, Toolchain-Pakete) — der Katalog kennt keine Optionen ([ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md), Re-Evaluierungs-Trigger); die Version überschreibt der Nutzer über `devcontainer.features.python.version`.
- Weitere Schlüssel — jeder Schlüssel braucht eigene Hosts und Tests; neue Wünsche kommen als eigener Slice.

## 2. Definition of Done

- [ ] **Katalog und Generator:** `python` im Katalog (offizielle Quelle, Standardversion `1`); Egress-Hosts `pypi.org` und `files.pythonhosted.org` nur bei aktivem Feature; Golden Cases belegen Render und Hosts.
- [ ] **Doku:** Spezifikation ([SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)), Benutzerhandbuch (inkl. Hinweis auf das Betriebssystem-Python nach PEP 668), READMEs und CHANGELOG nennen den Schlüssel.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_features.go` | update | Katalogeintrag `python` |
| `internal/hexagon/application/devcontainer_egress.go` | update | Hosts für `python` |
| `internal/hexagon/application/*_test.go` | update | Katalog-Schlüssel, Golden Cases |
| `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice des Implementers in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): das offizielle Feature `python` ist nicht mehr unter `ghcr.io/devcontainers/features/python` verfügbar.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Das Feature installiert das Betriebssystem-Python; `pip install` ohne venv lehnt PEP 668 ab — **Ausgang:** entfallen: gemessen (Python 3.13.5, `pip install` im venv läuft mit `pypi.org` und `files.pythonhosted.org`); das Handbuch nennt den venv-Hinweis.
- Weitere Paketquellen (private Indizes) brauchen eigene Hosts — **Ausgang:** weiter offen: der Nutzer ergänzt sie über `devcontainer.sandbox.egress.allow`.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/hexagon/application` (Greenfield laut `harness/conventions.md`), `docs/user/`, `spec/`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.
