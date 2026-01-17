package com.docflow;

import com.docflow.models.*;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import okhttp3.*;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * Client for the DocFlow API.
 * 
 * <p>
 * Provides methods for converting between Markdown and PDF formats.
 * 
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * DocFlowClient client = new DocFlowClient("http://localhost:8080");
 * 
 * List<MDFile> files = List.of(
 *         new MDFile("doc.md", "# Hello World"));
 * 
 * PDFResult result = client.convertMdToPdf(files, ConvertOptions.separate());
 * if (result.isSuccess()) {
 *     System.out.println("PDFs created: " + result.getFilePaths());
 * }
 * }</pre>
 */
public class DocFlowClient {

    private static final MediaType JSON = MediaType.get("application/json; charset=utf-8");

    private final String baseUrl;
    private final OkHttpClient httpClient;
    private final Gson gson;

    /**
     * Creates a new DocFlow client with default timeout.
     *
     * @param baseUrl Base URL of the DocFlow server (e.g., "http://localhost:8080")
     */
    public DocFlowClient(String baseUrl) {
        this(baseUrl, 60);
    }

    /**
     * Creates a new DocFlow client with custom timeout.
     *
     * @param baseUrl        Base URL of the DocFlow server
     * @param timeoutSeconds Request timeout in seconds
     */
    public DocFlowClient(String baseUrl, int timeoutSeconds) {
        this.baseUrl = baseUrl.replaceAll("/$", "");
        this.httpClient = new OkHttpClient.Builder()
                .connectTimeout(timeoutSeconds, TimeUnit.SECONDS)
                .readTimeout(timeoutSeconds, TimeUnit.SECONDS)
                .writeTimeout(timeoutSeconds, TimeUnit.SECONDS)
                .build();
        this.gson = new Gson();
    }

    /**
     * Checks server health.
     *
     * @return HealthResponse with status, version, and timestamp
     * @throws IOException If the request fails
     */
    public HealthResponse health() throws IOException {
        Request request = new Request.Builder()
                .url(baseUrl + "/api/health")
                .get()
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Health check failed with status: " + response.code());
            }
            return gson.fromJson(response.body().string(), HealthResponse.class);
        }
    }

    /**
     * Generates an HTML preview of markdown content.
     *
     * @param content Markdown content to preview
     * @return HTML string
     * @throws IOException If the request fails
     */
    public String preview(String content) throws IOException {
        JsonObject json = new JsonObject();
        json.addProperty("content", content);

        RequestBody body = RequestBody.create(gson.toJson(json), JSON);
        Request request = new Request.Builder()
                .url(baseUrl + "/api/preview")
                .post(body)
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Preview failed with status: " + response.code());
            }
            JsonObject result = gson.fromJson(response.body().string(), JsonObject.class);
            return result.get("html").getAsString();
        }
    }

    /**
     * Converts Markdown files to PDF.
     *
     * @param files   List of MDFile objects to convert
     * @param options Conversion options (merge mode, output name)
     * @return PDFResult with success status and file paths
     * @throws IOException If the request fails
     */
    public PDFResult convertMdToPdf(List<MDFile> files, ConvertOptions options) throws IOException {
        if (files == null || files.isEmpty()) {
            throw new IllegalArgumentException("At least one file is required");
        }

        if (options == null) {
            options = ConvertOptions.separate();
        }

        // Build request JSON
        JsonArray filesArray = new JsonArray();
        for (int i = 0; i < files.size(); i++) {
            MDFile file = files.get(i);
            JsonObject fileJson = new JsonObject();
            fileJson.addProperty("id", file.getId() != null ? file.getId() : String.valueOf(i + 1));
            fileJson.addProperty("name", file.getName());
            fileJson.addProperty("content", file.getContent());
            fileJson.addProperty("order", i);
            filesArray.add(fileJson);
        }

        JsonObject json = new JsonObject();
        json.add("files", filesArray);
        json.addProperty("mergeMode", options.getMergeMode());
        if (options.getOutputName() != null) {
            json.addProperty("outputName", options.getOutputName());
        }

        RequestBody body = RequestBody.create(gson.toJson(json), JSON);
        Request request = new Request.Builder()
                .url(baseUrl + "/api/convert")
                .post(body)
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            JsonObject result = gson.fromJson(response.body().string(), JsonObject.class);

            PDFResult pdfResult = new PDFResult();
            pdfResult.setSuccess(result.get("success").getAsBoolean());

            if (pdfResult.isSuccess() && result.has("files")) {
                List<String> filePaths = new ArrayList<>();
                result.get("files").getAsJsonArray().forEach(f -> filePaths.add(f.getAsString()));
                pdfResult.setFilePaths(filePaths);
            } else if (result.has("error")) {
                pdfResult.setError(result.get("error").getAsString());
            }

            return pdfResult;
        }
    }

    /**
     * Extracts text from a PDF and converts to Markdown.
     *
     * @param pdfData  Raw PDF file bytes
     * @param filename Original filename of the PDF
     * @return MDResult with success status and markdown content
     * @throws IOException If the request fails
     */
    public MDResult extractPdfToMd(byte[] pdfData, String filename) throws IOException {
        if (pdfData == null || pdfData.length == 0) {
            throw new IllegalArgumentException("PDF data is required");
        }
        if (filename == null || filename.isEmpty()) {
            throw new IllegalArgumentException("Filename is required");
        }

        // Encode PDF to base64
        String encoded = Base64.getEncoder().encodeToString(pdfData);

        JsonObject json = new JsonObject();
        json.addProperty("fileName", filename);
        json.addProperty("content", encoded);

        RequestBody body = RequestBody.create(gson.toJson(json), JSON);
        Request request = new Request.Builder()
                .url(baseUrl + "/api/pdf/extract")
                .post(body)
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            JsonObject result = gson.fromJson(response.body().string(), JsonObject.class);

            MDResult mdResult = new MDResult();
            mdResult.setSuccess(result.get("success").getAsBoolean());

            if (mdResult.isSuccess()) {
                mdResult.setMarkdown(result.has("markdown") ? result.get("markdown").getAsString() : "");
                mdResult.setFilePath(result.has("filePath") ? result.get("filePath").getAsString() : "");
                mdResult.setFileName(result.has("fileName") ? result.get("fileName").getAsString() : "");
            } else if (result.has("error")) {
                mdResult.setError(result.get("error").getAsString());
            }

            return mdResult;
        }
    }

    /**
     * Downloads a PDF file from the server.
     *
     * @param filePath Path returned from convertMdToPdf (e.g., "/output/doc.pdf")
     * @return Raw PDF bytes
     * @throws IOException If the download fails
     */
    public byte[] downloadPdf(String filePath) throws IOException {
        Request request = new Request.Builder()
                .url(baseUrl + filePath)
                .get()
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Download failed with status: " + response.code());
            }
            return response.body().bytes();
        }
    }
}
