package com.docflow.storage;

import com.azure.storage.blob.*;
import com.azure.storage.blob.models.*;

import java.io.*;
import java.util.*;
import java.util.stream.Collectors;

/**
 * Azure Blob Storage implementation.
 *
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * Storage storage = new AzureStorage("myaccount", "mycontainer", "accountkey");
 * storage.save("doc.pdf", pdfBytes);
 * }</pre>
 */
public class AzureStorage implements Storage {

    private final BlobContainerClient containerClient;
    private final String prefix;
    private final String accountName;
    private final String containerName;

    /**
     * Create Azure storage with account key.
     *
     * @param accountName   Azure storage account name.
     * @param containerName Container name.
     * @param accountKey    Account key.
     */
    public AzureStorage(String accountName, String containerName, String accountKey) {
        this(accountName, containerName, accountKey, "");
    }

    /**
     * Create Azure storage with optional prefix.
     *
     * @param accountName   Azure storage account name.
     * @param containerName Container name.
     * @param accountKey    Account key.
     * @param prefix        Optional prefix for all blob names.
     */
    public AzureStorage(String accountName, String containerName, String accountKey, String prefix) {
        this.accountName = accountName;
        this.containerName = containerName;
        this.prefix = prefix != null ? prefix.replaceAll("/$", "") : "";

        String connectionString = String.format(
                "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net",
                accountName, accountKey);

        BlobServiceClient serviceClient = new BlobServiceClientBuilder()
                .connectionString(connectionString)
                .buildClient();

        this.containerClient = serviceClient.getBlobContainerClient(containerName);
    }

    /**
     * Create Azure storage from connection string.
     *
     * @param connectionString Full Azure connection string.
     * @param containerName    Container name.
     * @param prefix           Optional prefix.
     */
    public static AzureStorage fromConnectionString(String connectionString, String containerName, String prefix) {
        AzureStorage storage = new AzureStorage("", containerName, "");
        // Override with connection string
        BlobServiceClient serviceClient = new BlobServiceClientBuilder()
                .connectionString(connectionString)
                .buildClient();

        try {
            var field = AzureStorage.class.getDeclaredField("containerClient");
            field.setAccessible(true);
            field.set(storage, serviceClient.getBlobContainerClient(containerName));
        } catch (Exception e) {
            throw new RuntimeException("Failed to create from connection string", e);
        }

        return storage;
    }

    @Override
    public void save(String path, byte[] data) throws IOException {
        String blobName = fullKey(path);
        try {
            BlobClient blobClient = containerClient.getBlobClient(blobName);
            blobClient.upload(new ByteArrayInputStream(data), data.length, true);
        } catch (Exception e) {
            throw new IOException("Failed to upload to Azure: " + e.getMessage(), e);
        }
    }

    @Override
    public void save(String path, InputStream inputStream) throws IOException {
        byte[] data = inputStream.readAllBytes();
        save(path, data);
    }

    @Override
    public byte[] load(String path) throws IOException {
        String blobName = fullKey(path);
        try {
            BlobClient blobClient = containerClient.getBlobClient(blobName);
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            blobClient.downloadStream(baos);
            return baos.toByteArray();
        } catch (BlobStorageException e) {
            if (e.getStatusCode() == 404) {
                throw new FileNotFoundException("File not found: " + path);
            }
            throw new IOException("Failed to download from Azure: " + e.getMessage(), e);
        }
    }

    @Override
    public void delete(String path) throws IOException {
        String blobName = fullKey(path);
        try {
            BlobClient blobClient = containerClient.getBlobClient(blobName);
            blobClient.deleteIfExists();
        } catch (Exception e) {
            // Ignore
        }
    }

    @Override
    public boolean exists(String path) {
        String blobName = fullKey(path);
        BlobClient blobClient = containerClient.getBlobClient(blobName);
        return blobClient.exists();
    }

    @Override
    public List<String> list(String directory) throws IOException {
        String blobPrefix = fullKey(directory);
        if (!blobPrefix.isEmpty() && !blobPrefix.endsWith("/")) {
            blobPrefix += "/";
        }

        try {
            String finalPrefix = blobPrefix;
            return containerClient.listBlobs(
                    new ListBlobsOptions().setPrefix(blobPrefix), null).stream()
                    .map(BlobItem::getName)
                    .map(name -> name.substring(finalPrefix.length()))
                    .filter(name -> !name.isEmpty() && !name.contains("/"))
                    .collect(Collectors.toList());
        } catch (Exception e) {
            throw new IOException("Failed to list Azure blobs: " + e.getMessage(), e);
        }
    }

    @Override
    public Optional<String> getUrl(String path) {
        String blobName = fullKey(path);
        return Optional.of(String.format(
                "https://%s.blob.core.windows.net/%s/%s",
                accountName, containerName, blobName));
    }

    private String fullKey(String path) {
        if (prefix.isEmpty()) {
            return path;
        }
        return prefix + "/" + path;
    }
}
