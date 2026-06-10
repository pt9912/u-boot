"""Markdown-Link-Validator fuer u-boot.

Scant alle Markdown-Dateien unter `docs/`, `spec/`, `harness/` und im
Repo-Root nach relativen `[text](path)`-Links und meldet alle nicht
aufloesbaren Pfade, nicht vorhandenen Markdown-Anker oder unzulaessigen
Referenzmodell-Kanten. Externe Links (`http://`, `https://`, `mailto:`)
werden uebersprungen.

Stdlib-only, kein Runtime-Dep. Aufruf:

    make docs-check              # Docker-gekapselt (kanonisch)
    python tools/check_refs.py   # direkt, wenn Python am Host vorhanden

Scope-Abgrenzung:
- Diese Variante deckt Markdown-Link-Pfade und Markdown-Heading-Anker ab.
- Zusaetzlich erzwingt sie die Referenzmatrix fuer Spec-Straten, ADR,
  Slice, Carveout und Roadmap/Welle.
- Nackte `ADR-NNNN`-Kennungen in normalem Markdown-Text muessen als
  Markdown-Link formuliert werden, damit sie als Dokumentkante pruefbar
  sind. Inline-Code zaehlt nicht als Link.
- Nackte `LH-*`-Kennungen in Sicht-Specs, ADR-Dokumenten, Planning-
  Dokumenten, Root-Markdown-Dateien, `harness/`, `docs/user/` und
  `docs/archive/` muessen ebenfalls als Markdown-Link formuliert werden,
  weil sie aufwaerts zum Vertrag referenzieren.
- Im Lastenheft selbst sind `LH-*`-Kennungen ausserhalb von
  Ueberschriften linkpflichtig, wenn die ID einen eigenen
  Lastenheft-Heading-Anker hat. Ueberschriften definieren die Anker;
  reine Index-Kennungen ohne eigenen Abschnitt bleiben unverlinkt.
- Konkrete Slice-/Tranche-IDs duerfen in Lastenheft, Sicht-Specs und ADRs
  nicht auftauchen, weil diese Artefakte nicht Richtung Planning zeigen.
  In allen anderen gescannten Markdown-Dateien muessen eindeutig
  aufloesbare Planning-IDs als Markdown-Link formuliert werden.
  Markdown-Ueberschriften sind ausgenommen, damit bestehende
  Section-Anker stabil bleiben.
- Konkrete `PH-*`-, `TC-*`- und `CO-*`-Kennungen muessen als
  Markdown-Link formuliert werden. Heute existieren konkrete `PH-*`- und
  `TC-*`-Kennungen nur als Traceability-Aliase im Lastenheft; ihre Links
  zeigen auf die zugehoerige `LH-*`-Anforderung derselben Matrixzeile.
- Verschachtelte Markdown-Link-Artefakte werden abgelehnt, weil sie fuer
  Menschen wie Links aussehen, aber vom Parser nur teilweise geprueft
  werden.
- LH-Kurzformen direkt hinter einem LH-Link (z. B. `.../008` oder
  `...-001..-007`) werden abgelehnt; der Zielpunkt muss selbst ein
  Markdown-Link sein.
- Weitere Kennungs-Aufloesung jenseits der aktivierten `ADR-*`-, `LH-*`-,
  Planning-ID- und Traceability-Alias-Linkpflicht bleibt fuer spaetere
  Erweiterungen.

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
from functools import lru_cache
from pathlib import Path
from urllib.parse import unquote

# Markdown-Link `[text](target)`. `text` darf escapte Klammern enthalten,
# aber fuer den Minimal-Linter reicht non-greedy. Bilder `![alt](src)`
# fangen wir mit demselben Pattern.
_LINK_PATTERN = re.compile(r"!?\[[^\]]*\]\(([^)]+)\)")

_REFERENCE_LINK_DEFINITION_PATTERN = re.compile(r"^\s{0,3}\[[^\]]+\]:\s*(<[^>]+>|\S+)")

_HEADING_PATTERN = re.compile(r"^(#{1,6})\s+(.+?)\s*#*\s*$")

_ADR_FILENAME_PATTERN = re.compile(r"^\d{4}-.+\.md$")

_ADR_ID_PATTERN = re.compile(r"\bADR-\d{4}\b")

_LH_ID_PATTERN = re.compile(r"\bLH(?:-[A-Z0-9]+)+-\d{3}[A-Z]?\b")

_PLANNING_ID_PATTERN = re.compile(
    r"\b(?:slice|tranche)-[a-z0-9](?:[a-z0-9-]|\.(?!md\b))*[a-z0-9]\b",
)

_TRACE_ALIAS_ID_PATTERN = re.compile(r"\b(?:PH|TC|CO)(?:-[A-Z0-9]+)+-\d{3}[A-Z]?\b")

_NESTED_LINK_ARTIFACT_PATTERN = re.compile(r"\]\([^)]+\)\]\(")

_LH_LINK_SHORTHAND_PATTERN = re.compile(
    r"\[`LH(?:-[A-Z0-9]+)+-\d{3}[A-Z]?`\]\([^)]+\)(?:\.\.|/)(?:`?-?\d{3}[A-Z]?`?)"
)

_ADR_INACTIVE_STATUS_PREFIXES = (
    "superseded by",
    "deprecated",
)

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
        print("[check_refs] all markdown references valid")
        return 0
    for violation in violations:
        print(violation.format(repo_root), file=sys.stderr)
    print(
        f"[check_refs] {len(violations)} broken markdown reference(s) — see stderr",
        file=sys.stderr,
    )
    return 1


def _iter_markdown_files(repo_root: Path) -> Iterator[Path]:
    """Liefert alle `*.md`-Dateien in `docs/`, `spec/`, `harness/` und im Root.

    Andere Top-Level-Verzeichnisse (`scripts/`, `tools/`, `cmd/`,
    `internal/`) enthalten heute keine fachlichen Markdown-Querverweise;
    bei Bedarf spaeter erweitern. Root-`*.md` deckt README.md /
    README.de.md ab.
    """
    for top in ("docs", "spec", "harness"):
        root = repo_root / top
        if not root.exists():
            continue
        yield from sorted(root.rglob("*.md"))
    for path in sorted(repo_root.glob("*.md")):
        yield path


def _check_file(repo_root: Path, source: Path) -> Iterator[BrokenRef]:
    text = source.read_text(encoding="utf-8")
    require_lh_links = _requires_lh_links(repo_root, source)
    require_planning_links = _requires_planning_links(repo_root, source)
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
        nested_match = _NESTED_LINK_ARTIFACT_PATTERN.search(scannable)
        if nested_match is not None:
            yield BrokenRef(
                source=source,
                lineno=lineno,
                target=nested_match.group(0),
                reason="nested markdown link artifact",
            )
        inline_code_spans = _inline_code_spans(line)
        for shorthand_match in _LH_LINK_SHORTHAND_PATTERN.finditer(line):
            if _point_is_contained(shorthand_match.start(), inline_code_spans):
                continue
            yield BrokenRef(
                source=source,
                lineno=lineno,
                target=shorthand_match.group(0),
                reason="LH shorthand suffix must be a markdown link",
            )
        reference_target = _reference_definition_target(line)
        if reference_target is not None:
            problem = _classify_target(repo_root, source, reference_target)
            if problem is not None:
                yield BrokenRef(
                    source=source,
                    lineno=lineno,
                    target=reference_target,
                    reason=problem,
                )
        for match in _LINK_PATTERN.finditer(scannable):
            target = match.group(1).strip()
            problem = _classify_target(repo_root, source, target)
            if problem is not None:
                yield BrokenRef(source=source, lineno=lineno, target=target, reason=problem)
        linkless = "" if reference_target is not None else _LINK_PATTERN.sub("", line)
        for match in _ADR_ID_PATTERN.finditer(linkless):
            yield BrokenRef(
                source=source,
                lineno=lineno,
                target=match.group(0),
                reason="ADR id must be a markdown link",
            )
        for match in _TRACE_ALIAS_ID_PATTERN.finditer(linkless):
            yield BrokenRef(
                source=source,
                lineno=lineno,
                target=match.group(0),
                reason="traceability alias id must be a markdown link",
            )
        if require_lh_links and _HEADING_PATTERN.match(line) is None:
            for match in _LH_ID_PATTERN.finditer(linkless):
                lh_id = match.group(0)
                if _artifact_kind(repo_root, source) == "contract_spec" and not _lh_id_has_heading(
                    source,
                    lh_id,
                ):
                    continue
                yield BrokenRef(
                    source=source,
                    lineno=lineno,
                    target=lh_id,
                    reason="LH id must be a markdown link",
                )
        if _HEADING_PATTERN.match(line) is None:
            for match in _PLANNING_ID_PATTERN.finditer(linkless):
                planning_id = match.group(0)
                if _resolve_planning_id(repo_root, planning_id) is None:
                    continue
                if _forbids_planning_references(repo_root, source):
                    yield BrokenRef(
                        source=source,
                        lineno=lineno,
                        target=planning_id,
                        reason="planning id must not be referenced from normative specs or ADRs",
                    )
                    continue
                if not require_planning_links:
                    continue
                yield BrokenRef(
                    source=source,
                    lineno=lineno,
                    target=planning_id,
                    reason="planning id must be a markdown link",
                )


def _classify_target(repo_root: Path, source: Path, target: str) -> str | None:
    """`None` wenn der Target aufloesbar (oder bewusst extern) ist,
    sonst Begruendungs-String.
    """
    if _is_external(target):
        return None
    path_part, anchor = _split_target(target)
    if not path_part:
        target_path = source
    else:
        if path_part.startswith("/"):
            return "absolute path; expected repo-relative"
        target_path, problem = _resolve_relative_path(repo_root, source, path_part)
        if problem is not None:
            return problem
    semantic_problem = _check_reference_model(repo_root, source, target_path)
    if semantic_problem is not None:
        return semantic_problem
    if anchor is None:
        return None
    return _check_anchor(target_path, anchor)


def _check_reference_model(repo_root: Path, source: Path, target: Path) -> str | None:
    """Prueft die normative Referenzmatrix fuer kanonische Artefakte.

    README-, Harness- und User-Dokumentation bleiben ausserhalb dieser
    Matrix: Links dort sind Navigations-/Onboarding-Kontext, keine
    normative Ableitung.
    """
    if source == target:
        return None

    source_kind = _artifact_kind(repo_root, source)
    target_kind = _artifact_kind(repo_root, target)
    if source_kind is None or target_kind is None:
        return None

    if source_kind == "contract_spec":
        return (
            "semantic reference violation: contract spec may only link intra-spec, "
            f"not {_artifact_label(target_kind)}"
        )

    if source_kind == "adr":
        if target_kind in {"contract_spec", "view_spec", "adr"}:
            return None
        return (
            "semantic reference violation: ADR may only link spec strata or ADR lineage, "
            f"not {_artifact_label(target_kind)}"
        )

    if source_kind == "view_spec":
        if target_kind in {"contract_spec", "view_spec"}:
            return None
        return (
            "semantic reference violation: view spec may not link down to "
            f"{_artifact_label(target_kind)}"
        )

    if source_kind in {"slice", "carveout"} and target_kind == "adr":
        if _is_active_adr(target):
            return None
        return (
            "semantic reference violation: slices and carveouts may reference "
            "only active ADRs"
        )

    return None


def _artifact_kind(repo_root: Path, path: Path) -> str | None:
    try:
        rel = path.relative_to(repo_root)
    except ValueError:
        return None

    if rel == Path("spec/lastenheft.md"):
        return "contract_spec"
    if rel == Path("spec/architecture.md"):
        return "view_spec"

    parts = rel.parts
    if (
        len(parts) == 4
        and parts[:3] == ("docs", "plan", "adr")
        and _ADR_FILENAME_PATTERN.match(parts[3])
    ):
        return "adr"

    if len(parts) >= 5 and parts[:3] == ("docs", "plan", "planning"):
        filename = parts[-1]
        lifecycle = parts[3]
        if lifecycle == "in-progress" and filename == "carveouts.md":
            return "carveout"
        if lifecycle == "in-progress" and filename == "roadmap.md":
            return "roadmap"
        if filename.startswith(("slice-", "tranche-")) and filename.endswith(".md"):
            return "slice"

    return None


def _artifact_label(kind: str) -> str:
    return {
        "contract_spec": "Vertrag/Lastenheft",
        "view_spec": "Sicht-Spec",
        "adr": "ADR",
        "slice": "Slice",
        "carveout": "Carveout",
        "roadmap": "Roadmap/Welle",
    }.get(kind, kind)


def _requires_lh_links(repo_root: Path, path: Path) -> bool:
    try:
        rel = path.relative_to(repo_root)
    except ValueError:
        return False

    if _artifact_kind(repo_root, path) in {"contract_spec", "view_spec"}:
        return True
    if len(rel.parts) >= 3 and rel.parts[:3] == ("docs", "plan", "adr"):
        return True
    if len(rel.parts) >= 3 and rel.parts[:3] == ("docs", "plan", "planning"):
        return True
    if len(rel.parts) == 1 and rel.suffix == ".md":
        return True
    if len(rel.parts) >= 1 and rel.parts[0] == "harness":
        return True
    if len(rel.parts) >= 2 and rel.parts[:2] == ("docs", "archive"):
        return True
    return len(rel.parts) >= 2 and rel.parts[:2] == ("docs", "user")


def _requires_planning_links(repo_root: Path, path: Path) -> bool:
    try:
        rel = path.relative_to(repo_root)
    except ValueError:
        return False

    return not _forbids_planning_references(repo_root, path)


def _forbids_planning_references(repo_root: Path, path: Path) -> bool:
    kind = _artifact_kind(repo_root, path)
    return kind in {"contract_spec", "view_spec", "adr"}


def _lh_id_has_heading(path: Path, lh_id: str) -> bool:
    return lh_id in _lh_heading_ids(path)


@lru_cache(maxsize=None)
def _lh_heading_ids(path: Path) -> frozenset[str]:
    found: set[str] = set()
    in_fenced = False
    for line in path.read_text(encoding="utf-8").splitlines():
        stripped = line.lstrip()
        if stripped.startswith("```") or stripped.startswith("~~~"):
            in_fenced = not in_fenced
            continue
        if in_fenced:
            continue
        match = _HEADING_PATTERN.match(line)
        if match is None:
            continue
        id_match = _LH_ID_PATTERN.search(match.group(2))
        if id_match is not None:
            found.add(id_match.group(0))
    return frozenset(found)


def _resolve_planning_id(repo_root: Path, planning_id: str) -> Path | None:
    matches = _planning_id_targets(repo_root).get(planning_id)
    if matches is None or len(matches) != 1:
        return None
    return matches[0]


@lru_cache(maxsize=None)
def _planning_id_targets(repo_root: Path) -> dict[str, tuple[Path, ...]]:
    root = repo_root / "docs/plan/planning"
    found: dict[str, list[Path]] = {}
    if not root.exists():
        return {}
    for path in sorted(root.rglob("*.md")):
        if not path.name.startswith(("slice-", "tranche-")):
            continue
        found.setdefault(path.stem, []).append(path)
    return {key: tuple(paths) for key, paths in found.items()}


@lru_cache(maxsize=None)
def _is_active_adr(path: Path) -> bool:
    status = _adr_status(path)
    if status is None:
        return False
    normalized = status.lower()
    return not normalized.startswith(_ADR_INACTIVE_STATUS_PREFIXES)


@lru_cache(maxsize=None)
def _adr_status(path: Path) -> str | None:
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()
    for index, line in enumerate(lines):
        if line.strip().lower() != "## status":
            continue
        for candidate in lines[index + 1 :]:
            stripped = candidate.strip()
            if not stripped:
                continue
            if stripped.startswith("#"):
                return None
            if stripped.startswith(">"):
                continue
            return stripped
        return None
    return None


def _split_target(target: str) -> tuple[str, str | None]:
    """Split `path#anchor` links while preserving anchor-only refs."""
    if "#" not in target:
        return target, None
    path_part, anchor = target.split("#", 1)
    return path_part, unquote(anchor)


def _is_external(target: str) -> bool:
    """`True` fuer Targets, die nicht als Repo-Pfad zu pruefen sind:
    leerer String, `http://`, `https://`, `mailto:`, `ftp://`.
    """
    if not target:
        return True
    return target.startswith(_EXTERNAL_PREFIXES)


def _resolve_relative_path(
    repo_root: Path,
    source: Path,
    path_part: str,
) -> tuple[Path | None, str | None]:
    """Resolve `path_part` relativ zu `source.parent` und prueft Existenz.

    Defense-in-depth: Symlinks werden ausdruecklich abgelehnt, bevor
    `Path.resolve()` ihnen folgen koennte. Verhindert Info-Leak ueber
    boesartige Symlinks auf Pseudo-Dateisysteme (z. B. /proc/self/...).
    """
    candidate_raw = source.parent / path_part
    if candidate_raw.is_symlink():
        return None, "target is a symlink (validator refuses to follow)"
    candidate = candidate_raw.resolve()
    try:
        candidate.relative_to(repo_root)
    except ValueError:
        return None, f"escapes repo root ({candidate})"
    if not candidate.exists():
        return None, "target file does not exist"
    return candidate, None


def _check_anchor(target_path: Path, anchor: str) -> str | None:
    if not anchor:
        return None
    if target_path.suffix.lower() != ".md":
        return "anchor target is not a markdown file"
    anchors = _markdown_heading_anchors(target_path)
    if anchor not in anchors:
        return f"anchor not found in target markdown ({anchor})"
    return None


@lru_cache(maxsize=None)
def _markdown_heading_anchors(path: Path) -> frozenset[str]:
    anchors: set[str] = set()
    counts: dict[str, int] = {}
    in_fenced = False
    text = path.read_text(encoding="utf-8")
    for line in text.splitlines():
        stripped = line.lstrip()
        if stripped.startswith("```") or stripped.startswith("~~~"):
            in_fenced = not in_fenced
            continue
        if in_fenced:
            continue
        match = _HEADING_PATTERN.match(line)
        if match is None:
            continue
        base_slug = _github_heading_slug(match.group(2))
        if not base_slug:
            continue
        index = counts.get(base_slug, 0)
        counts[base_slug] = index + 1
        anchors.add(base_slug if index == 0 else f"{base_slug}-{index}")
    return frozenset(anchors)


def _github_heading_slug(heading: str) -> str:
    """Best-effort GitHub/GFM heading slug for project documentation.

    Keeps unicode word characters and ASCII hyphens, strips common
    Markdown punctuation, lowercases, and maps each whitespace character
    to one `-`.
    """
    text = heading.strip()
    text = re.sub(r"<[^>]+>", "", text)
    text = re.sub(r"!\[([^\]]*)\]\([^)]+\)", r"\1", text)
    text = re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", text)
    text = text.replace("`", "")
    text = text.lower()
    text = re.sub(r"[^\w\s-]", "", text, flags=re.UNICODE)
    return re.sub(r"\s", "-", text.strip())


def _reference_definition_target(line: str) -> str | None:
    match = _REFERENCE_LINK_DEFINITION_PATTERN.match(line)
    if match is None:
        return None
    target = match.group(1).strip()
    if target.startswith("<") and target.endswith(">"):
        return target[1:-1]
    return target


def _inline_code_spans(line: str) -> tuple[tuple[int, int], ...]:
    return tuple(match.span() for match in _INLINE_CODE_PATTERN.finditer(line))


def _point_is_contained(point: int, containers: tuple[tuple[int, int], ...]) -> bool:
    return any(start <= point < end for start, end in containers)


def _check_relative_path(repo_root: Path, source: Path, path_part: str) -> str | None:
    """Backward-compatible wrapper for path-only references."""
    _, problem = _resolve_relative_path(repo_root, source, path_part)
    return problem


if __name__ == "__main__":
    raise SystemExit(main())
