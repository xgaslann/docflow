package com.docflow;

import com.docflow.models.*;
import com.docflow.storage.LocalStorage;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class DocFlowTest {

    @TempDir
    Path tempDir;

    // MarkdownParser Tests
    @Nested
    class MarkdownParserTests {

        private MarkdownParser parser;

        @BeforeEach
        void setUp() {
            parser = new MarkdownParser();
        }

        @Test
        void testHeading() {
            String html = parser.toHtml("# Hello World");
            assertTrue(html.contains("<h1"));
            assertTrue(html.contains("Hello World"));
        }

        @Test
        void testParagraph() {
            String html = parser.toHtml("This is a paragraph.");
            assertTrue(html.contains("<p>"));
            assertTrue(html.contains("This is a paragraph."));
        }

        @Test
        void testBold() {
            String html = parser.toHtml("This is **bold** text.");
            assertTrue(html.contains("<strong>"));
            assertTrue(html.contains("bold"));
        }

        @Test
        void testItalic() {
            String html = parser.toHtml("This is *italic* text.");
            assertTrue(html.contains("<em>"));
        }

        @Test
        void testCode() {
            String html = parser.toHtml("This is `code` text.");
            assertTrue(html.contains("<code>"));
        }

        @Test
        void testList() {
            String html = parser.toHtml("- Item 1\n- Item 2");
            assertTrue(html.contains("<ul>"));
            assertTrue(html.contains("<li>"));
        }

        @Test
        void testMergeFiles() {
            List<MDFile> files = List.of(
                    new MDFile("second.md", "# Second") {
                        {
                            setOrder(1);
                        }
                    },
                    new MDFile("first.md", "# First") {
                        {
                            setOrder(0);
                        }
                    },
                    new MDFile("third.md", "# Third") {
                        {
                            setOrder(2);
                        }
                    });

            String merged = parser.mergeFiles(files);

            assertTrue(merged.contains("# First"));
            assertTrue(merged.contains("# Second"));
            assertTrue(merged.contains("# Third"));
            assertTrue(merged.indexOf("# First") < merged.indexOf("# Second"));
        }

        @Test
        void testEstimatePageCount() {
            assertEquals(1, parser.estimatePageCount("x".repeat(100)));
            assertEquals(2, parser.estimatePageCount("x".repeat(6000)));
            assertEquals(5, parser.estimatePageCount("x".repeat(15000)));
        }
    }

    // Template Tests
    @Nested
    class TemplateTests {

        @Test
        void testGenerate() {
            Template template = new Template();
            String html = template.generate("<h1>Test</h1>");

            assertTrue(html.contains("<!DOCTYPE html>"));
            assertTrue(html.contains("<h1>Test</h1>"));
            assertTrue(html.contains("@page"));
        }
    }

    // LocalStorage Tests
    @Nested
    class LocalStorageTests {

        private LocalStorage storage;

        @BeforeEach
        void setUp() throws IOException {
            storage = new LocalStorage(tempDir.toString());
        }

        @Test
        void testSaveAndLoad() throws IOException {
            byte[] data = "Hello, World!".getBytes();
            storage.save("test.txt", data);
            byte[] loaded = storage.load("test.txt");
            assertArrayEquals(data, loaded);
        }

        @Test
        void testExists() throws IOException {
            assertFalse(storage.exists("nonexistent.txt"));
            storage.save("exists.txt", "data".getBytes());
            assertTrue(storage.exists("exists.txt"));
        }

        @Test
        void testDelete() throws IOException {
            storage.save("to_delete.txt", "data".getBytes());
            assertTrue(storage.exists("to_delete.txt"));
            storage.delete("to_delete.txt");
            assertFalse(storage.exists("to_delete.txt"));
        }

        @Test
        void testList() throws IOException {
            storage.save("file1.txt", "data1".getBytes());
            storage.save("file2.txt", "data2".getBytes());
            List<String> files = storage.list("");
            assertTrue(files.contains("file1.txt"));
            assertTrue(files.contains("file2.txt"));
        }

        @Test
        void testGetUrl() throws IOException {
            storage.save("file.txt", "data".getBytes());
            var url = storage.getUrl("file.txt");
            assertTrue(url.isPresent());
            assertTrue(url.get().startsWith("file://"));
        }
    }

    // MDFile Tests
    @Nested
    class MDFileTests {

        @Test
        void testConstructor() {
            MDFile file = new MDFile("test.md", "# Hello");
            assertEquals("test.md", file.getName());
            assertEquals("# Hello", file.getContent());
        }

        @Test
        void testBuilder() {
            MDFile file = MDFile.builder()
                    .name("test.md")
                    .content("# Hello")
                    .order(5)
                    .build();

            assertEquals("test.md", file.getName());
            assertEquals(5, file.getOrder());
        }
    }

    // ConvertOptions Tests
    @Nested
    class ConvertOptionsTests {

        @Test
        void testSeparate() {
            ConvertOptions opts = ConvertOptions.separate();
            assertEquals("separate", opts.getMergeMode());
        }

        @Test
        void testMerged() {
            ConvertOptions opts = ConvertOptions.merged("output");
            assertEquals("merged", opts.getMergeMode());
            assertEquals("output", opts.getOutputName());
        }
    }

    // Integration Test (if needed)
    @Test
    @Disabled("Requires full setup")
    void testConverterIntegration() throws IOException {
        Converter converter = new Converter(tempDir.toString());

        List<MDFile> files = List.of(
                new MDFile("test.md", "# Test Document\n\nThis is a test."));

        PDFResult result = converter.convertToPdf(files, ConvertOptions.separate());
        assertTrue(result.isSuccess());
        assertFalse(result.getFilePaths().isEmpty());
    }
}
