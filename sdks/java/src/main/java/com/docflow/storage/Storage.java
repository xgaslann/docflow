package com.docflow.storage;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Optional;

/**
 * Interface for file storage backends.
 */
public interface Storage {

    /**
     * Save data to the given path.
     *
     * @param path Relative path where to save the file.
     * @param data Binary data to save.
     * @throws IOException If saving fails.
     */
    void save(String path, byte[] data) throws IOException;

    /**
     * Save data from an input stream.
     *
     * @param path        Relative path where to save the file.
     * @param inputStream Stream to read from.
     * @throws IOException If saving fails.
     */
    void save(String path, InputStream inputStream) throws IOException;

    /**
     * Load data from the given path.
     *
     * @param path Relative path to the file.
     * @return Binary file contents.
     * @throws IOException If loading fails or file doesn't exist.
     */
    byte[] load(String path) throws IOException;

    /**
     * Delete file at the given path.
     *
     * @param path Relative path to the file.
     * @throws IOException If deletion fails.
     */
    void delete(String path) throws IOException;

    /**
     * Check if file exists at the given path.
     *
     * @param path Relative path to check.
     * @return true if file exists.
     */
    boolean exists(String path);

    /**
     * List files in the given directory.
     *
     * @param directory Relative path to directory.
     * @return List of filenames in the directory.
     * @throws IOException If listing fails.
     */
    List<String> list(String directory) throws IOException;

    /**
     * Get URL for accessing the file.
     *
     * @param path Relative path to the file.
     * @return URL string if supported, empty otherwise.
     */
    default Optional<String> getUrl(String path) {
        return Optional.empty();
    }
}
