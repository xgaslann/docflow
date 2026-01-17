"""HTML template for PDF conversion."""


class Template:
    """Generates HTML templates for PDF conversion with print-optimized styling."""

    def generate(self, body_content: str) -> str:
        """Wrap HTML content in a complete HTML document.
        
        Args:
            body_content: HTML content for the body.
            
        Returns:
            Complete HTML document string.
        """
        return f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        {self._get_styles()}
    </style>
</head>
<body>
    <article class="document">
        {body_content}
    </article>
</body>
</html>"""

    def _get_styles(self) -> str:
        return """
        * { margin: 0; padding: 0; box-sizing: border-box; }
        @page { size: A4; margin: 20mm; }
        html { font-size: 11pt; -webkit-print-color-adjust: exact; print-color-adjust: exact; }
        body {
            font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6; color: #1a1a1a; background: #ffffff;
        }
        .document { max-width: 100%; margin: 0; padding: 0; }
        
        h1, h2, h3, h4, h5, h6 {
            font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
            font-weight: 600; line-height: 1.3; color: #111111;
            margin-top: 1.5em; margin-bottom: 0.5em;
            page-break-after: avoid; break-after: avoid;
        }
        h1 { font-size: 1.8em; color: #000; border-bottom: 2px solid #2563eb; padding-bottom: 0.3em; margin-top: 0; }
        h2 { font-size: 1.4em; border-bottom: 1px solid #d1d5db; padding-bottom: 0.2em; }
        h3 { font-size: 1.2em; }
        h4 { font-size: 1.1em; }
        h5, h6 { font-size: 1em; }
        
        p { margin-bottom: 0.8em; text-align: justify; hyphens: auto; orphans: 3; widows: 3; }
        a { color: #2563eb; text-decoration: none; }
        
        code {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.9em; background-color: #f3f4f6; color: #be185d;
            padding: 0.15em 0.4em; border-radius: 3px; border: 1px solid #e5e7eb;
        }
        
        pre {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.85em; background-color: #1e293b; color: #e2e8f0;
            padding: 1em; border-radius: 6px; margin: 1em 0;
            border-left: 4px solid #2563eb; page-break-inside: avoid;
            white-space: pre-wrap; word-wrap: break-word;
        }
        pre code { background: none; color: inherit; padding: 0; border: none; }
        
        blockquote {
            margin: 1em 0; padding: 0.8em 1.2em;
            border-left: 4px solid #2563eb; background-color: #f0f9ff;
            color: #374151; font-style: italic; page-break-inside: avoid;
        }
        
        ul, ol { margin: 0.8em 0; padding-left: 2em; }
        li { margin-bottom: 0.3em; line-height: 1.5; }
        
        table {
            width: 100%; border-collapse: collapse; margin: 1em 0;
            font-size: 0.9em; page-break-inside: avoid;
        }
        th, td { border: 1px solid #d1d5db; padding: 0.6em 0.8em; text-align: left; }
        th { background-color: #2563eb; color: #ffffff; font-weight: 600; }
        tr:nth-child(even) { background-color: #f9fafb; }
        
        img { max-width: 100%; height: auto; margin: 1em 0; page-break-inside: avoid; }
        hr { border: none; border-top: 1px solid #d1d5db; margin: 1.5em 0; }
        
        .file-separator { page-break-before: always; break-before: page; }
        .file-separator::before { content: ''; display: block; border-top: 2px solid #2563eb; margin-bottom: 1em; }
        .file-separator span, .file-header span {
            display: inline-block; background: #2563eb; color: #fff;
            padding: 0.3em 0.8em; border-radius: 4px; font-size: 0.8em;
            font-weight: 500; margin-bottom: 1em;
        }
        .file-header { margin-bottom: 1em; }
        
        @media print {
            body { font-size: 10pt; line-height: 1.5; }
            pre { font-size: 8pt; }
            h1, h2, h3, h4, h5, h6 { page-break-after: avoid; }
            pre, blockquote, table, img { page-break-inside: avoid; }
        }
        """
