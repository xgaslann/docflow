package com.docflow.storage.vector;

import com.docflow.config.VectorStoreConfig;
import com.docflow.models.Chunk;
import com.docflow.models.RAGDocument;

import java.util.*;

/**
 * MongoDB Atlas Vector Store using vector search.
 * Requires mongodb-driver-sync dependency.
 */
public class MongoDBVectorStore {

    private final VectorStoreConfig config;
    private Object mongoClient; // com.mongodb.client.MongoClient
    private Object collection; // com.mongodb.client.MongoCollection

    public MongoDBVectorStore(VectorStoreConfig config) {
        this.config = config;
    }

    public MongoDBVectorStore(String uri, String database, String collection) {
        this.config = new VectorStoreConfig();
        this.config.setConnectionString(uri);
        this.config.setDatabase(database);
        this.config.setCollection(collection);
        this.config.setProvider(VectorStoreConfig.Provider.MONGODB);
    }

    /**
     * Initialize the MongoDB connection.
     */
    public void initialize() throws Exception {
        try {
            Class<?> clientClass = Class.forName("com.mongodb.client.MongoClients");
            mongoClient = clientClass.getMethod("create", String.class)
                    .invoke(null, config.getConnectionString());

            Object database = mongoClient.getClass().getMethod("getDatabase", String.class)
                    .invoke(mongoClient, config.getDatabase());

            collection = database.getClass().getMethod("getCollection", String.class)
                    .invoke(database, config.getCollection());
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("MongoDB driver not found. Add mongodb-driver-sync dependency.");
        }
    }

    /**
     * Upsert a document and its chunks.
     */
    public void upsert(RAGDocument doc) throws Exception {
        if (collection == null) {
            throw new IllegalStateException("Not initialized. Call initialize() first.");
        }

        for (Chunk chunk : doc.getChunks()) {
            Map<String, Object> document = new HashMap<>();
            document.put("_id", doc.getId() + "_" + chunk.getIndex());
            document.put("document_id", doc.getId());
            document.put("chunk_index", chunk.getIndex());
            document.put("content", chunk.getContent());
            document.put("embedding", new ArrayList<>()); // Would need to generate
            document.put("metadata", chunk.getMetadata());
            document.put("created_at", new Date());

            // Use replaceOne with upsert
            Class<?> filtersClass = Class.forName("com.mongodb.client.model.Filters");
            Class<?> optionsClass = Class.forName("com.mongodb.client.model.ReplaceOptions");

            Object filter = filtersClass.getMethod("eq", String.class, Object.class)
                    .invoke(null, "_id", document.get("_id"));

            Object options = optionsClass.getConstructor().newInstance();
            options.getClass().getMethod("upsert", boolean.class).invoke(options, true);

            // Convert Map to Document
            Class<?> docClass = Class.forName("org.bson.Document");
            Object bsonDoc = docClass.getConstructor(Map.class).newInstance(document);

            collection.getClass()
                    .getMethod("replaceOne", Class.forName("org.bson.conversions.Bson"), docClass, optionsClass)
                    .invoke(collection, filter, bsonDoc, options);
        }
    }

    /**
     * Vector search.
     */
    public List<SearchResult> search(List<Float> queryVector, int topK) throws Exception {
        if (collection == null) {
            throw new IllegalStateException("Not initialized. Call initialize() first.");
        }

        List<SearchResult> results = new ArrayList<>();

        try {
            // Build vector search pipeline
            List<Map<String, Object>> pipeline = new ArrayList<>();

            // $vectorSearch stage
            Map<String, Object> vectorSearch = new HashMap<>();
            Map<String, Object> vectorSearchParams = new HashMap<>();
            vectorSearchParams.put("index", config.getIndexName());
            vectorSearchParams.put("path", "embedding");
            vectorSearchParams.put("queryVector", queryVector);
            vectorSearchParams.put("numCandidates", config.getNumCandidates());
            vectorSearchParams.put("limit", topK);
            vectorSearch.put("$vectorSearch", vectorSearchParams);
            pipeline.add(vectorSearch);

            // $project stage
            Map<String, Object> project = new HashMap<>();
            Map<String, Object> projectFields = new HashMap<>();
            projectFields.put("content", 1);
            projectFields.put("metadata", 1);
            projectFields.put("score", Map.of("$meta", "vectorSearchScore"));
            project.put("$project", projectFields);
            pipeline.add(project);

            // Execute aggregation using reflection
            Class<?> docClass = Class.forName("org.bson.Document");
            List<Object> bsonPipeline = new ArrayList<>();
            for (Map<String, Object> stage : pipeline) {
                bsonPipeline.add(docClass.getConstructor(Map.class).newInstance(stage));
            }

            Object aggregateResult = collection.getClass()
                    .getMethod("aggregate", List.class)
                    .invoke(collection, bsonPipeline);

            // Iterate results
            Object iterator = aggregateResult.getClass().getMethod("iterator").invoke(aggregateResult);
            while ((boolean) iterator.getClass().getMethod("hasNext").invoke(iterator)) {
                Object doc = iterator.getClass().getMethod("next").invoke(iterator);

                SearchResult result = new SearchResult();
                result.setId((String) doc.getClass().getMethod("getString", String.class).invoke(doc, "_id"));
                result.setContent((String) doc.getClass().getMethod("getString", String.class).invoke(doc, "content"));
                result.setScore(
                        ((Number) doc.getClass().getMethod("get", String.class).invoke(doc, "score")).floatValue());
                results.add(result);
            }
        } catch (Exception e) {
            throw new RuntimeException("Vector search failed: " + e.getMessage(), e);
        }

        return results;
    }

    /**
     * Delete document chunks.
     */
    public void delete(String documentId) throws Exception {
        if (collection == null) {
            throw new IllegalStateException("Not initialized. Call initialize() first.");
        }

        Class<?> filtersClass = Class.forName("com.mongodb.client.model.Filters");
        Object filter = filtersClass.getMethod("eq", String.class, Object.class)
                .invoke(null, "document_id", documentId);

        collection.getClass()
                .getMethod("deleteMany", Class.forName("org.bson.conversions.Bson"))
                .invoke(collection, filter);
    }

    /**
     * Close the connection.
     */
    public void close() throws Exception {
        if (mongoClient != null) {
            mongoClient.getClass().getMethod("close").invoke(mongoClient);
        }
    }

    /**
     * Search result.
     */
    public static class SearchResult {
        private String id;
        private String content;
        private float score;
        private Map<String, Object> metadata;

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public float getScore() {
            return score;
        }

        public void setScore(float score) {
            this.score = score;
        }

        public Map<String, Object> getMetadata() {
            return metadata;
        }

        public void setMetadata(Map<String, Object> metadata) {
            this.metadata = metadata;
        }
    }
}
