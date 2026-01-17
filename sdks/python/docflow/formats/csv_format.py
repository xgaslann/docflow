"""CSV format converter for DocFlow."""

import csv
import io
from typing import List, Optional, Union

from ..types import ConvertResult


class CSVConverter:
    """Converts between CSV and Markdown formats.
    
    Example:
        >>> converter = CSVConverter()
        >>> md = converter.to_markdown(csv_data, "data.csv")
        >>> csv_data = converter.from_markdown(md_content)
    """

    def __init__(
        self,
        delimiter: str = ",",
        has_header: bool = True,
        table_title: Optional[str] = None,
    ) -> None:
        """Initialize CSV converter.
        
        Args:
            delimiter: CSV delimiter character.
            has_header: Whether first row is header.
            table_title: Optional title for the table in markdown.
        """
        self.delimiter = delimiter
        self.has_header = has_header
        self.table_title = table_title

    def to_markdown(
        self,
        csv_data: Union[str, bytes],
        filename: str = "data.csv",
        include_metadata: bool = True,
    ) -> ConvertResult:
        """Convert CSV to Markdown table.
        
        Args:
            csv_data: CSV content as string or bytes.
            filename: Original filename for metadata.
            include_metadata: Include YAML frontmatter.
            
        Returns:
            ConvertResult with markdown content.
        """
        try:
            if isinstance(csv_data, bytes):
                csv_data = csv_data.decode("utf-8")

            reader = csv.reader(io.StringIO(csv_data), delimiter=self.delimiter)
            rows = list(reader)

            if not rows:
                return ConvertResult(success=False, error="Empty CSV file")

            md_parts = []

            # Metadata frontmatter
            if include_metadata:
                md_parts.append("---")
                md_parts.append(f"source: {filename}")
                md_parts.append("format: csv")
                md_parts.append(f"rows: {len(rows)}")
                md_parts.append(f"columns: {len(rows[0]) if rows else 0}")
                md_parts.append("---\n")

            # Title
            title = self.table_title or filename.replace(".csv", "").replace("_", " ").title()
            md_parts.append(f"# {title}\n")

            # Table
            if self.has_header and len(rows) > 0:
                header = rows[0]
                data_rows = rows[1:]
            else:
                header = [f"Column {i+1}" for i in range(len(rows[0]))]
                data_rows = rows

            # Header row
            md_parts.append("| " + " | ".join(self._escape_cell(c) for c in header) + " |")
            # Separator
            md_parts.append("| " + " | ".join("---" for _ in header) + " |")
            # Data rows
            for row in data_rows:
                # Pad row if shorter than header
                padded = row + [""] * (len(header) - len(row))
                md_parts.append("| " + " | ".join(self._escape_cell(c) for c in padded[:len(header)]) + " |")

            return ConvertResult(
                success=True,
                content="\n".join(md_parts),
                format="csv",
                metadata={"rows": len(rows), "columns": len(rows[0]) if rows else 0},
            )

        except Exception as e:
            return ConvertResult(success=False, error=str(e))

    def from_markdown(
        self,
        markdown: str,
        include_header: bool = True,
    ) -> ConvertResult:
        """Extract tables from Markdown and convert to CSV.
        
        Args:
            markdown: Markdown content.
            include_header: Include header row in CSV.
            
        Returns:
            ConvertResult with CSV content.
        """
        try:
            tables = self._extract_tables(markdown)
            
            if not tables:
                return ConvertResult(success=False, error="No tables found in markdown")

            # Convert first table to CSV
            table = tables[0]
            output = io.StringIO()
            writer = csv.writer(output, delimiter=self.delimiter)

            start_idx = 0 if include_header else 1
            for row in table[start_idx:]:
                writer.writerow(row)

            return ConvertResult(
                success=True,
                content=output.getvalue(),
                format="csv",
                metadata={"tables_found": len(tables)},
            )

        except Exception as e:
            return ConvertResult(success=False, error=str(e))

    def _extract_tables(self, markdown: str) -> List[List[List[str]]]:
        """Extract all tables from markdown."""
        tables = []
        lines = markdown.split("\n")
        
        current_table = []
        in_table = False
        
        for line in lines:
            line = line.strip()
            
            if line.startswith("|") and line.endswith("|"):
                # Skip separator rows
                if set(line.replace("|", "").replace("-", "").replace(":", "").strip()) == set():
                    continue
                
                cells = [c.strip() for c in line.split("|")[1:-1]]
                current_table.append(cells)
                in_table = True
            else:
                if in_table and current_table:
                    tables.append(current_table)
                    current_table = []
                in_table = False
        
        if current_table:
            tables.append(current_table)
        
        return tables

    def _escape_cell(self, cell: str) -> str:
        """Escape pipe characters in cell content."""
        return cell.replace("|", "\\|").replace("\n", " ")
