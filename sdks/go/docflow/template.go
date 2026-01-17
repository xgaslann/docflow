package docflow

import "fmt"

// Template generates HTML templates for PDF conversion.
type Template struct{}

// NewTemplate creates a new template generator.
func NewTemplate() *Template {
	return &Template{}
}

// Generate wraps HTML content in a complete HTML document with print-optimized styling.
func (t *Template) Generate(bodyContent string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        %s
    </style>
</head>
<body>
    <article class="document">
        %s
    </article>
</body>
</html>`, t.getStyles(), bodyContent)
}

func (t *Template) getStyles() string {
	return `
        /* Reset and base */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        /* Page setup for print */
        @page {
            size: A4;
            margin: 20mm;
        }

        html {
            font-size: 11pt;
            -webkit-print-color-adjust: exact;
            print-color-adjust: exact;
        }

        body {
            font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #1a1a1a;
            background: #ffffff;
            max-width: 100%;
            padding: 0;
            margin: 0;
        }

        .document {
            max-width: 100%;
            margin: 0;
            padding: 0;
        }

        /* Typography */
        h1, h2, h3, h4, h5, h6 {
            font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
            font-weight: 600;
            line-height: 1.3;
            color: #111111;
            margin-top: 1.5em;
            margin-bottom: 0.5em;
            page-break-after: avoid;
            break-after: avoid;
        }

        h1 {
            font-size: 1.8em;
            color: #000000;
            border-bottom: 2px solid #2563eb;
            padding-bottom: 0.3em;
            margin-top: 0;
        }

        h2 {
            font-size: 1.4em;
            color: #1a1a1a;
            border-bottom: 1px solid #d1d5db;
            padding-bottom: 0.2em;
        }

        h3 { font-size: 1.2em; }
        h4 { font-size: 1.1em; }
        h5, h6 { font-size: 1em; }

        /* Paragraphs */
        p {
            margin-bottom: 0.8em;
            text-align: justify;
            hyphens: auto;
            orphans: 3;
            widows: 3;
        }

        /* Links */
        a {
            color: #2563eb;
            text-decoration: none;
        }

        /* Inline code */
        code {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
            background-color: #f3f4f6;
            color: #be185d;
            padding: 0.15em 0.4em;
            border-radius: 3px;
            border: 1px solid #e5e7eb;
        }

        /* Code blocks */
        pre {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.85em;
            background-color: #1e293b;
            color: #e2e8f0;
            padding: 1em;
            border-radius: 6px;
            overflow-x: auto;
            margin: 1em 0;
            border-left: 4px solid #2563eb;
            page-break-inside: avoid;
            break-inside: avoid;
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        pre code {
            background: none;
            color: inherit;
            padding: 0;
            border: none;
            font-size: 1em;
        }

        /* Blockquotes */
        blockquote {
            margin: 1em 0;
            padding: 0.8em 1.2em;
            border-left: 4px solid #2563eb;
            background-color: #f0f9ff;
            color: #374151;
            font-style: italic;
            page-break-inside: avoid;
            break-inside: avoid;
        }

        blockquote p { margin-bottom: 0; }
        blockquote p + p { margin-top: 0.5em; }

        /* Lists */
        ul, ol {
            margin: 0.8em 0;
            padding-left: 2em;
        }

        li {
            margin-bottom: 0.3em;
            line-height: 1.5;
        }

        li > ul, li > ol {
            margin-top: 0.3em;
            margin-bottom: 0;
        }

        /* Tables */
        table {
            width: 100%%;
            border-collapse: collapse;
            margin: 1em 0;
            font-size: 0.9em;
            page-break-inside: avoid;
            break-inside: avoid;
        }

        th, td {
            border: 1px solid #d1d5db;
            padding: 0.6em 0.8em;
            text-align: left;
            vertical-align: top;
        }

        th {
            background-color: #2563eb;
            color: #ffffff;
            font-weight: 600;
        }

        tr:nth-child(even) { background-color: #f9fafb; }
        tr:hover { background-color: #f3f4f6; }

        /* Images */
        img {
            max-width: 100%%;
            height: auto;
            margin: 1em 0;
            page-break-inside: avoid;
            break-inside: avoid;
        }

        /* Horizontal rules */
        hr {
            border: none;
            border-top: 1px solid #d1d5db;
            margin: 1.5em 0;
        }

        /* File separators */
        .file-separator {
            page-break-before: always;
            break-before: page;
            margin-top: 0;
            padding-top: 0;
        }

        .file-separator::before {
            content: '';
            display: block;
            border-top: 2px solid #2563eb;
            margin-bottom: 1em;
        }

        .file-separator span,
        .file-header span {
            display: inline-block;
            background: #2563eb;
            color: #ffffff;
            padding: 0.3em 0.8em;
            border-radius: 4px;
            font-size: 0.8em;
            font-weight: 500;
            margin-bottom: 1em;
        }

        .file-header { margin-bottom: 1em; }
        .file-content { margin-bottom: 0; }

        /* Print styles */
        @media print {
            body {
                font-size: 10pt;
                line-height: 1.5;
            }
            .document { padding: 0; }
            pre {
                white-space: pre-wrap;
                word-wrap: break-word;
                font-size: 8pt;
            }
            h1, h2, h3, h4, h5, h6 {
                page-break-after: avoid;
                break-after: avoid;
            }
            pre, blockquote, table, figure, img {
                page-break-inside: avoid;
                break-inside: avoid;
            }
            .file-separator {
                page-break-before: always;
                break-before: page;
            }
        }
    `
}
