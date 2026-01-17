#!/usr/bin/env python3
"""
DocFlow Python SDK Examples

Demonstrates standalone usage of the DocFlow library.
No server required.
"""

import os
from pathlib import Path


def main():
    # Create output directory
    output_dir = Path("./output")
    output_dir.mkdir(exist_ok=True)

    print("=== Example 1: Basic Conversion ===")
    basic_conversion(output_dir)

    print("\n=== Example 2: Merge Multiple Files ===")
    merge_files(output_dir)

    print("\n=== Example 3: PDF Extraction ===")
    extract_pdf()

    print("\n=== Example 4: Get PDF as Bytes ===")
    get_pdf_bytes(output_dir)

    print("\n=== Example 5: Preview Markdown ===")
    preview_markdown()


def basic_conversion(output_dir: Path):
    """Convert a single markdown file to PDF."""
    from docflow import Converter, MDFile, LocalStorage

    # Create converter with storage
    storage = LocalStorage(str(output_dir))
    converter = Converter(storage=storage)

    # Create markdown file
    files = [
        MDFile(
            name="hello.md",
            content="""# Hello World

This is a **bold** statement and this is *italic*.

## Features

- Easy to use
- Standalone library
- No server required

## Code Example

```python
print("Hello, World!")
```
"""
        )
    ]

    # Convert to PDF
    result = converter.convert_to_pdf(files)

    if result.success:
        print(f"✓ PDF created: {result.file_paths}")
    else:
        print(f"✗ Error: {result.error}")


def merge_files(output_dir: Path):
    """Merge multiple markdown files into a single PDF."""
    from docflow import Converter, MDFile, ConvertOptions, LocalStorage

    storage = LocalStorage(str(output_dir))
    converter = Converter(storage=storage)

    files = [
        MDFile(name="chapter1.md", content="# Chapter 1\n\nIntroduction to the topic.", order=0),
        MDFile(name="chapter2.md", content="# Chapter 2\n\nDeeper exploration.", order=1),
        MDFile(name="chapter3.md", content="# Chapter 3\n\nConclusion.", order=2),
    ]

    result = converter.convert_to_pdf(
        files,
        ConvertOptions(merge_mode="merged", output_name="combined_document")
    )

    if result.success:
        print(f"✓ Merged PDF created: {result.file_paths}")


def extract_pdf():
    """Extract text from PDF and convert to markdown."""
    from docflow import Extractor

    sample_pdf = Path("./sample.pdf")
    if not sample_pdf.exists():
        print("⊙ Skipping: sample.pdf not found")
        print("  Create a sample.pdf file to test extraction")
        return

    extractor = Extractor()

    with open(sample_pdf, "rb") as f:
        result = extractor.extract_to_markdown(f.read(), "sample.pdf")

    if result.success:
        print(f"✓ Extracted {result.page_count} pages")
        preview = result.markdown[:200] + "..." if len(result.markdown) > 200 else result.markdown
        print(f"  Preview: {preview}")


def get_pdf_bytes(output_dir: Path):
    """Get PDF as bytes for direct use."""
    from docflow import Converter, MDFile

    converter = Converter()

    files = [
        MDFile(name="inline.md", content="# Inline PDF\n\nThis PDF is generated as bytes.")
    ]

    pdf_bytes = converter.convert_to_bytes(files)

    # Save bytes to file
    output_path = output_dir / "from_bytes.pdf"
    output_path.write_bytes(pdf_bytes)

    print(f"✓ PDF bytes saved: {output_path} ({len(pdf_bytes)} bytes)")


def preview_markdown():
    """Generate HTML preview of markdown."""
    from docflow import Converter

    converter = Converter()

    html = converter.preview("# Preview\n\nThis is a **preview** of the markdown.")

    print(f"✓ HTML Preview ({len(html)} bytes):")
    preview = html[:100] + "..." if len(html) > 100 else html
    print(f"  {preview}")


if __name__ == "__main__":
    main()
