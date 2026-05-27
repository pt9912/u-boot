"""Markdown-Link-Validator fuer u-boot.

Scant alle Markdown-Dateien unter `docs/` und `spec/` nach relativen
`[text](path)`-Links und meldet alle nicht aufloesbaren Pfade. Externe
Links (`http://`, `https://`, `mailto:`, `#anchor`-only) werden
uebersprungen.

Stdlib-only, kein Runtime-Dep. Aufruf:

    make docs-check              # Docker-gekapselt (kanonisch)
    python tools/check_refs.py   # direkt, wenn Python am Host vorhanden

Scope-Abgrenzung:
- Diese Variante deckt nur Markdown-Link-Pfade ab.
- Kennungs-Aufloesung (`LH-*`-Querverweise) und `§`-Sektions-Verweise
  bleiben fuer eine spaetere Erweiterung.

Output-Format pro Verstoss: `{source_rel}:{lineno}\\t{target}\\t{reason}`.
Exit-Code 0 = alle Links aufloesbar, 1 = mindestens ein Verstoss.

Adaptiert von c-hsm-doc/tools/check_refs.py (gleiche Build-Familie,
gleiche Docker-only-Philosophie); u-boot-Spezifika: LH-* statt HSM-*.
"""

from __future__ import annotations

import re
import sys
from collections.abc import Iterator
from dataclasses import dataclass
from pathlib import Path

# Markdown-Link `[text](target)`. `text` darf escapte Klammern enthalten,
# aber fuer den Minimal-Linter reicht non-greedy. Bilder `![alt](src)`
# fangen wir mit demselben Pattern.
_LINK_PATTERN = re.compile(r"!?\[[^\]]*\]\(([^)]+)\)")

# Inline-Code `` `...` ``-Spans werden vor dem Link-Match aus der Zeile
# entfernt, damit Demo-Strings in Backticks (z. B. ``[text](path)``)
# nicht als echte Links interpretiert werden.
#
# Zwei Alternativen, laenger zuerst (Markdown erlaubt N-fache Backticks
# als Delimiter, sodass der Span N-1 einzelne Backticks enthalten kann):
# - ``...`` (doppelter Delimiter, darf einzelne ` enthalten — lazy match)
# - `...`  (einfacher Delimiter, keine ` im Inneren)
_INLINE_CODE_PATTERN = re.compile(r"``.+?``|`[^`]+`")

# Pfade mit einem dieser Praefixe sind keine relativen Filesystem-Refs.
_EXTERNAL_PREFIXES: tuple[str, ...] = (
    "http://",
    "https://",
    "mailto:",
    "ftp://",
)


@dataclass(frozen=True, slots=True)
class BrokenRef:
    """Eine festgestellte, nicht aufloesbare Pfad-Referenz."""

    source: Path
    lineno: int
    target: str
    reason: str

    def format(self, repo_root: Path) -> str:
        rel = self.source.relative_to(repo_root)
        return f"{rel}:{self.lineno}\t{self.target}\t{self.reason}"


def main() -> int:
    repo_root = Path(__file__).resolve().parents[1]
    sources = _iter_markdown_files(repo_root)
    violations: list[BrokenRef] = []
    for source in sources:
        violations.extend(_check_file(repo_root, source))
    if not violations:
        print("[check_refs] all markdown link targets resolved")
        return 0
    for violation in violations:
        print(violation.format(repo_root), file=sys.stderr)
    print(
        f"[check_refs] {len(violations)} broken markdown reference(s) — see stderr",
        file=sys.stderr,
    )
    return 1


def _iter_markdown_files(repo_root: Path) -> Iterator[Path]:
    """Liefert alle `*.md`-Dateien in `docs/`, `spec/` und im Root.

    Andere Top-Level-Verzeichnisse (`scripts/`, `tools/`, `cmd/`,
    `internal/`) enthalten heute keine fachlichen Markdown-Querverweise;
    bei Bedarf spaeter erweitern. Root-`*.md` deckt README.md /
    README.de.md ab.
    """
    for top in ("docs", "spec"):
        root = repo_root / top
        if not root.exists():
            continue
        yield from sorted(root.rglob("*.md"))
    for path in sorted(repo_root.glob("*.md")):
        yield path


def _check_file(repo_root: Path, source: Path) -> Iterator[BrokenRef]:
    text = source.read_text(encoding="utf-8")
    in_fenced = False
    for lineno, line in enumerate(text.splitlines(), start=1):
        # Fenced code blocks (```...``` or ~~~...~~~) sind keine Link-Quelle;
        # alles dazwischen ueberspringen. Die Fence-Zeile selbst zaehlt auch
        # als „in block" (sicherer: keine Links auf der Fence-Zeile selber).
        stripped = line.lstrip()
        if stripped.startswith("```") or stripped.startswith("~~~"):
            in_fenced = not in_fenced
            continue
        if in_fenced:
            continue
        # Inline-Code-Spans ausblenden.
        scannable = _INLINE_CODE_PATTERN.sub("", line)
        for match in _LINK_PATTERN.finditer(scannable):
            target = match.group(1).strip()
            problem = _classify_target(repo_root, source, target)
            if problem is not None:
                yield BrokenRef(source=source, lineno=lineno, target=target, reason=problem)


def _classify_target(repo_root: Path, source: Path, target: str) -> str | None:
    """`None` wenn der Target aufloesbar (oder bewusst extern) ist,
    sonst Begruendungs-String.
    """
    if _is_external_or_anchor(target):
        return None
    path_part = target.split("#", 1)[0]
    if not path_part:
        return None
    if path_part.startswith("/"):
        return "absolute path; expected repo-relative"
    return _check_relative_path(repo_root, source, path_part)


def _is_external_or_anchor(target: str) -> bool:
    """`True` fuer Targets, die nicht als Repo-Pfad zu pruefen sind:
    leerer String, reine `#anchor`-Refs, `http://`, `https://`,
    `mailto:`, `ftp://`.
    """
    if not target or target.startswith("#"):
        return True
    return target.startswith(_EXTERNAL_PREFIXES)


def _check_relative_path(repo_root: Path, source: Path, path_part: str) -> str | None:
    """Resolve `path_part` relativ zu `source.parent` und prueft Existenz.

    Defense-in-depth: Symlinks werden ausdruecklich abgelehnt, bevor
    `Path.resolve()` ihnen folgen koennte. Verhindert Info-Leak ueber
    boesartige Symlinks auf Pseudo-Dateisysteme (z. B. /proc/self/...).
    """
    candidate_raw = source.parent / path_part
    if candidate_raw.is_symlink():
        return "target is a symlink (validator refuses to follow)"
    candidate = candidate_raw.resolve()
    try:
        candidate.relative_to(repo_root)
    except ValueError:
        return f"escapes repo root ({candidate})"
    if not candidate.exists():
        return "target file does not exist"
    return None


if __name__ == "__main__":
    raise SystemExit(main())
