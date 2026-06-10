from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

try:
    from tools import check_refs
except ModuleNotFoundError:  # pragma: no cover - direct script execution fallback.
    import check_refs


class CheckRefsTest(unittest.TestCase):
    def test_heading_slug_uses_markdown_link_text_without_destination(self) -> None:
        heading = (
            "Slice V1: `u-boot logs` "
            "([`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005-logs-anzeigen))"
        )

        self.assertEqual(
            check_refs._github_heading_slug(heading),
            "slice-v1-u-boot-logs-lh-fa-up-005",
        )

    def test_reference_style_definition_counts_as_link(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            repo = Path(temp)
            target = repo / "docs/plan/adr/0001-foo.md"
            target.parent.mkdir(parents=True)
            target.write_text("# ADR 0001\n", encoding="utf-8")
            source = repo / "README.md"
            source.write_text("[ADR-0001]: docs/plan/adr/0001-foo.md\n", encoding="utf-8")

            violations = list(check_refs._check_file(repo, source))

        self.assertEqual(violations, [])

    def test_lh_shorthand_rule_ignores_inline_code_examples(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            repo = Path(temp)
            source = repo / "scratch.md"
            source.write_text(
                "Example: `[`LH-FA-001`](spec/lastenheft.md#x)/-002`\n",
                encoding="utf-8",
            )

            violations = list(check_refs._check_file(repo, source))

        self.assertEqual(violations, [])

    def test_lh_shorthand_rule_still_rejects_real_markdown_links(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            repo = Path(temp)
            target = repo / "target.md"
            target.write_text("# Target\n", encoding="utf-8")
            source = repo / "scratch.md"
            source.write_text("[`LH-FA-001`](target.md)/-002\n", encoding="utf-8")

            violations = list(check_refs._check_file(repo, source))

        self.assertEqual(len(violations), 1)
        self.assertEqual(violations[0].reason, "LH shorthand suffix must be a markdown link")


if __name__ == "__main__":
    unittest.main()
