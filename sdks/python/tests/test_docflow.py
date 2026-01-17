"""Unit tests for DocFlow Python SDK."""

import os
import tempfile
from pathlib import Path

import pytest

from docflow import (
    MDFile,
    ConvertOptions,
    PDFResult,
    ExtractResult,
    LocalStorage,
    MarkdownParser,
    Template,
)


class TestMDFile:
    """Tests for MDFile dataclass."""

    def test_create_mdfile(self):
        file = MDFile(name="test.md", content="# Hello")
        assert file.name == "test.md"
        assert file.content == "# Hello"
        assert file.id == "test.md"
        assert file.order == 0

    def test_to_dict(self):
        file = MDFile(name="test.md", content="# Hello", order=5)
        d = file.to_dict()
        assert d["name"] == "test.md"
        assert d["content"] == "# Hello"
        assert d["order"] == 5


class TestMarkdownParser:
    """Tests for MarkdownParser."""

    @pytest.fixture
    def parser(self):
        return MarkdownParser()

    def test_heading(self, parser):
        html = parser.to_html("# Hello World")
        assert "<h1>" in html or "<h1" in html
        assert "Hello World" in html

    def test_paragraph(self, parser):
        html = parser.to_html("This is a paragraph.")
        assert "<p>" in html
        assert "This is a paragraph." in html

    def test_bold(self, parser):
        html = parser.to_html("This is **bold** text.")
        assert "<strong>" in html
        assert "bold" in html

    def test_italic(self, parser):
        html = parser.to_html("This is *italic* text.")
        assert "<em>" in html
        assert "italic" in html

    def test_code(self, parser):
        html = parser.to_html("This is `code` text.")
        assert "<code>" in html
        assert "code" in html

    def test_list(self, parser):
        html = parser.to_html("- Item 1\n- Item 2")
        assert "<ul>" in html
        assert "<li>" in html
        assert "Item 1" in html

    def test_table(self, parser):
        md = "| A | B |\n|---|---|\n| 1 | 2 |"
        html = parser.to_html(md)
        assert "<table>" in html
        assert "<th>" in html or "<td>" in html

    def test_merge_files(self, parser):
        files = [
            MDFile(name="second.md", content="# Second", order=1),
            MDFile(name="first.md", content="# First", order=0),
            MDFile(name="third.md", content="# Third", order=2),
        ]
        merged = parser.merge_files(files)

        assert "# First" in merged
        assert "# Second" in merged
        assert "# Third" in merged
        # First should appear before Second
        assert merged.index("# First") < merged.index("# Second")

    def test_estimate_page_count(self, parser):
        assert parser.estimate_page_count("x" * 100) == 1
        assert parser.estimate_page_count("x" * 6000) == 2
        assert parser.estimate_page_count("x" * 15000) == 5


class TestTemplate:
    """Tests for Template."""

    def test_generate(self):
        template = Template()
        html = template.generate("<h1>Test</h1>")

        assert "<!DOCTYPE html>" in html
        assert "<h1>Test</h1>" in html
        assert "@page" in html
        assert "size: A4" in html


class TestLocalStorage:
    """Tests for LocalStorage."""

    @pytest.fixture
    def storage(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            yield LocalStorage(tmpdir)

    def test_save_and_load(self, storage):
        data = b"Hello, World!"
        storage.save("test.txt", data)
        loaded = storage.load("test.txt")
        assert loaded == data

    def test_exists(self, storage):
        assert not storage.exists("nonexistent.txt")
        storage.save("exists.txt", b"data")
        assert storage.exists("exists.txt")

    def test_delete(self, storage):
        storage.save("to_delete.txt", b"data")
        assert storage.exists("to_delete.txt")
        storage.delete("to_delete.txt")
        assert not storage.exists("to_delete.txt")

    def test_list(self, storage):
        storage.save("file1.txt", b"data1")
        storage.save("file2.txt", b"data2")
        files = storage.list("")
        assert "file1.txt" in files
        assert "file2.txt" in files

    def test_nested_directory(self, storage):
        storage.save("subdir/nested.txt", b"nested data")
        assert storage.exists("subdir/nested.txt")
        loaded = storage.load("subdir/nested.txt")
        assert loaded == b"nested data"

    def test_get_url(self, storage):
        storage.save("file.txt", b"data")
        url = storage.get_url("file.txt")
        assert url.startswith("file://")
        assert "file.txt" in url


class TestConvertOptions:
    """Tests for ConvertOptions."""

    def test_default_values(self):
        opts = ConvertOptions()
        assert opts.merge_mode == "separate"
        assert opts.output_name is None

    def test_custom_values(self):
        opts = ConvertOptions(merge_mode="merged", output_name="output")
        assert opts.merge_mode == "merged"
        assert opts.output_name == "output"


class TestResults:
    """Tests for result types."""

    def test_pdf_result_success(self):
        result = PDFResult(success=True, file_paths=["file.pdf"])
        assert result.success
        assert "file.pdf" in result.file_paths

    def test_pdf_result_error(self):
        result = PDFResult(success=False, error="Something went wrong")
        assert not result.success
        assert result.error == "Something went wrong"

    def test_extract_result(self):
        result = ExtractResult(
            success=True,
            markdown="# Title",
            page_count=5,
        )
        assert result.success
        assert result.markdown == "# Title"
        assert result.page_count == 5
