package com.docflow;

/**
 * HTML template generator for PDF conversion.
 */
public class Template {

    /**
     * Wrap HTML content in a complete HTML document.
     *
     * @param bodyContent HTML content for the body.
     * @return Complete HTML document string.
     */
    public String generate(String bodyContent) {
        return String.format("""
                <!DOCTYPE html>
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
                </html>
                """, getStyles(), bodyContent);
    }

    private String getStyles() {
        return """
                * { margin: 0; padding: 0; box-sizing: border-box; }
                @page { size: A4; margin: 20mm; }
                html { font-size: 11pt; }
                body {
                    font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
                    line-height: 1.6; color: #1a1a1a; background: #ffffff;
                }
                .document { max-width: 100%%; margin: 0; padding: 0; }

                h1, h2, h3, h4, h5, h6 {
                    font-weight: 600; line-height: 1.3; color: #111111;
                    margin-top: 1.5em; margin-bottom: 0.5em;
                }
                h1 { font-size: 1.8em; border-bottom: 2px solid #2563eb; padding-bottom: 0.3em; margin-top: 0; }
                h2 { font-size: 1.4em; border-bottom: 1px solid #d1d5db; padding-bottom: 0.2em; }
                h3 { font-size: 1.2em; }

                p { margin-bottom: 0.8em; text-align: justify; }
                a { color: #2563eb; text-decoration: none; }

                code {
                    font-family: 'Consolas', monospace; font-size: 0.9em;
                    background-color: #f3f4f6; color: #be185d;
                    padding: 0.15em 0.4em; border-radius: 3px;
                }

                pre {
                    font-family: 'Consolas', monospace; font-size: 0.85em;
                    background-color: #1e293b; color: #e2e8f0;
                    padding: 1em; border-radius: 6px; margin: 1em 0;
                    border-left: 4px solid #2563eb;
                    white-space: pre-wrap; word-wrap: break-word;
                }
                pre code { background: none; color: inherit; padding: 0; }

                blockquote {
                    margin: 1em 0; padding: 0.8em 1.2em;
                    border-left: 4px solid #2563eb; background-color: #f0f9ff;
                    font-style: italic;
                }

                ul, ol { margin: 0.8em 0; padding-left: 2em; }
                li { margin-bottom: 0.3em; }

                table { width: 100%%; border-collapse: collapse; margin: 1em 0; }
                th, td { border: 1px solid #d1d5db; padding: 0.6em 0.8em; text-align: left; }
                th { background-color: #2563eb; color: #ffffff; font-weight: 600; }
                tr:nth-child(even) { background-color: #f9fafb; }

                img { max-width: 100%%; height: auto; margin: 1em 0; }
                hr { border: none; border-top: 1px solid #d1d5db; margin: 1.5em 0; }

                .file-separator { page-break-before: always; }
                .file-separator span, .file-header span {
                    display: inline-block; background: #2563eb; color: #fff;
                    padding: 0.3em 0.8em; border-radius: 4px; font-size: 0.8em;
                    margin-bottom: 1em;
                }
                """;
    }
}
