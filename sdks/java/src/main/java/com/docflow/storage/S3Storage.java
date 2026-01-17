package com.docflow.storage;

import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.*;

import java.io.*;
import java.net.URI;
import java.util.*;
import java.util.stream.Collectors;

/**
 * AWS S3 storage implementation.
 *
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * Storage storage = new S3Storage("my-bucket", "eu-west-1");
 * storage.save("doc.pdf", pdfBytes);
 * }</pre>
 */
public class S3Storage implements Storage {

    private final S3Client client;
    private final String bucket;
    private final String prefix;
    private final String region;

    /**
     * Create S3 storage with default credentials.
     *
     * @param bucket S3 bucket name.
     * @param region AWS region.
     */
    public S3Storage(String bucket, String region) {
        this(bucket, region, "", null);
    }

    /**
     * Create S3 storage with optional prefix and endpoint.
     *
     * @param bucket      S3 bucket name.
     * @param region      AWS region.
     * @param prefix      Optional prefix for all keys.
     * @param endpointUrl Optional custom endpoint (for MinIO, LocalStack).
     */
    public S3Storage(String bucket, String region, String prefix, String endpointUrl) {
        this.bucket = bucket;
        this.region = region;
        this.prefix = prefix != null ? prefix.replaceAll("/$", "") : "";

        var builder = S3Client.builder()
                .region(Region.of(region))
                .credentialsProvider(DefaultCredentialsProvider.create());

        if (endpointUrl != null && !endpointUrl.isEmpty()) {
            builder.endpointOverride(URI.create(endpointUrl));
            builder.forcePathStyle(true);
        }

        this.client = builder.build();
    }

    @Override
    public void save(String path, byte[] data) throws IOException {
        String key = fullKey(path);
        try {
            client.putObject(
                    PutObjectRequest.builder()
                            .bucket(bucket)
                            .key(key)
                            .build(),
                    RequestBody.fromBytes(data));
        } catch (Exception e) {
            throw new IOException("Failed to upload to S3: " + e.getMessage(), e);
        }
    }

    @Override
    public void save(String path, InputStream inputStream) throws IOException {
        byte[] data = inputStream.readAllBytes();
        save(path, data);
    }

    @Override
    public byte[] load(String path) throws IOException {
        String key = fullKey(path);
        try {
            return client.getObjectAsBytes(
                    GetObjectRequest.builder()
                            .bucket(bucket)
                            .key(key)
                            .build())
                    .asByteArray();
        } catch (NoSuchKeyException e) {
            throw new FileNotFoundException("File not found: " + path);
        } catch (Exception e) {
            throw new IOException("Failed to download from S3: " + e.getMessage(), e);
        }
    }

    @Override
    public void delete(String path) throws IOException {
        String key = fullKey(path);
        try {
            client.deleteObject(
                    DeleteObjectRequest.builder()
                            .bucket(bucket)
                            .key(key)
                            .build());
        } catch (Exception e) {
            // Ignore errors
        }
    }

    @Override
    public boolean exists(String path) {
        String key = fullKey(path);
        try {
            client.headObject(
                    HeadObjectRequest.builder()
                            .bucket(bucket)
                            .key(key)
                            .build());
            return true;
        } catch (Exception e) {
            return false;
        }
    }

    @Override
    public List<String> list(String directory) throws IOException {
        String keyPrefix = fullKey(directory);
        if (!keyPrefix.isEmpty() && !keyPrefix.endsWith("/")) {
            keyPrefix += "/";
        }

        try {
            String finalPrefix = keyPrefix;
            ListObjectsV2Response response = client.listObjectsV2(
                    ListObjectsV2Request.builder()
                            .bucket(bucket)
                            .prefix(keyPrefix)
                            .build());

            return response.contents().stream()
                    .map(S3Object::key)
                    .map(key -> key.substring(finalPrefix.length()))
                    .filter(name -> !name.isEmpty() && !name.contains("/"))
                    .collect(Collectors.toList());
        } catch (Exception e) {
            throw new IOException("Failed to list S3 objects: " + e.getMessage(), e);
        }
    }

    @Override
    public Optional<String> getUrl(String path) {
        String key = fullKey(path);
        return Optional.of(String.format("s3://%s/%s", bucket, key));
    }

    /**
     * Get HTTP URL for the file.
     */
    public String getHttpUrl(String path) {
        String key = fullKey(path);
        return String.format("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key);
    }

    private String fullKey(String path) {
        if (prefix.isEmpty()) {
            return path;
        }
        return prefix + "/" + path;
    }
}
