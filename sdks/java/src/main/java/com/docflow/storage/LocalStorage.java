package com.docflow.storage;

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Local filesystem storage implementation.
 *
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * Storage storage = new LocalStorage("./output");
 * storage.save("doc.pdf", pdfBytes);
 * }</pre>
 */
public class LocalStorage implements Storage {

    private final Path basePath;

    /**
     * Create local storage at the given path.
     *
     * @param basePath Root directory for all operations.
     * @throws IOException If directory creation fails.
     */
    public LocalStorage(String basePath) throws IOException {
        this.basePath = Paths.get(basePath).toAbsolutePath();
        Files.createDirectories(this.basePath);
    }

    @Override
    public void save(String path, byte[] data) throws IOException {
        Path fullPath = fullPath(path);
        Files.createDirectories(fullPath.getParent());
        Files.write(fullPath, data);
    }

    @Override
    public void save(String path, InputStream inputStream) throws IOException {
        Path fullPath = fullPath(path);
        Files.createDirectories(fullPath.getParent());
        Files.copy(inputStream, fullPath, StandardCopyOption.REPLACE_EXISTING);
    }

    @Override
    public byte[] load(String path) throws IOException {
        Path fullPath = fullPath(path);
        if (!Files.exists(fullPath)) {
            throw new FileNotFoundException("File not found: " + path);
        }
        return Files.readAllBytes(fullPath);
    }

    @Override
    public void delete(String path) throws IOException {
        Path fullPath = fullPath(path);
        Files.deleteIfExists(fullPath);
    }

    @Override
    public boolean exists(String path) {
        return Files.exists(fullPath(path));
    }

    @Override
    public List<String> list(String directory) throws IOException {
        Path dirPath = fullPath(directory);
        if (!Files.exists(dirPath)) {
            return Collections.emptyList();
        }

        try (Stream<Path> stream = Files.list(dirPath)) {
            return stream
                    .filter(Files::isRegularFile)
                    .map(p -> p.getFileName().toString())
                    .collect(Collectors.toList());
        }
    }

    @Override
    public Optional<String> getUrl(String path) {
        return Optional.of("file://" + fullPath(path).toString());
    }

    /**
     * Get the absolute filesystem path.
     *
     * @param path Relative path.
     * @return Absolute path string.
     */
    public String getAbsolutePath(String path) {
        return fullPath(path).toString();
    }

    /**
     * Get the base path of this storage.
     *
     * @return Base path.
     */
    public Path getBasePath() {
        return basePath;
    }

    private Path fullPath(String path) {
        return basePath.resolve(path);
    }
}
