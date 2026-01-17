package com.docflow.config;

// No additional imports needed - all types are primitives or enums

/**
 * Configuration for vector storage.
 */
public class VectorStoreConfig {

    public enum Provider {
        POSTGRESQL("postgresql"),
        MONGODB("mongodb");

        private final String value;

        Provider(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    public enum IndexType {
        HNSW("hnsw"),
        IVFFLAT("ivfflat"),
        FLAT("flat");

        private final String value;

        IndexType(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    public enum DistanceMetric {
        COSINE("cosine"),
        EUCLIDEAN("euclidean"),
        DOT_PRODUCT("dot");

        private final String value;

        DistanceMetric(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    private Provider provider = Provider.POSTGRESQL;
    private String connectionString = "";
    private String database = "docflow";
    private String collection = "chunks";

    // Embedding
    private String embeddingProvider = "openai";
    private String embeddingModel = "text-embedding-3-small";
    private String embeddingApiKey = "";
    private int embeddingDimensions = 1536;
    private int embeddingBatchSize = 100;

    // Index
    private IndexType indexType = IndexType.HNSW;
    private DistanceMetric distanceMetric = DistanceMetric.COSINE;

    // PostgreSQL specific
    private String host = "localhost";
    private int port = 5432;
    private String user = "postgres";
    private String password = "";
    private String sslMode = "prefer";
    private String schema = "public";
    private String tableName = "chunks";
    private int m = 16;
    private int efConstruction = 64;
    private int efSearch = 40;
    private int lists = 100;
    private int probes = 10;

    // MongoDB specific
    private String atlasCluster = "";
    private int numCandidates = 100;
    private String indexName = "vector_index";

    // Performance
    private int poolSize = 5;
    private int timeout = 30;

    public VectorStoreConfig() {
    }

    public static VectorStoreConfig defaultConfig() {
        return new VectorStoreConfig();
    }

    public String getDsn() {
        if (!connectionString.isEmpty()) {
            return connectionString;
        }
        return String.format("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
                user, password, host, port, database, sslMode);
    }

    // Getters and setters
    public Provider getProvider() {
        return provider;
    }

    public void setProvider(Provider provider) {
        this.provider = provider;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public void setConnectionString(String connectionString) {
        this.connectionString = connectionString;
    }

    public String getDatabase() {
        return database;
    }

    public void setDatabase(String database) {
        this.database = database;
    }

    public String getCollection() {
        return collection;
    }

    public void setCollection(String collection) {
        this.collection = collection;
    }

    public String getEmbeddingProvider() {
        return embeddingProvider;
    }

    public void setEmbeddingProvider(String embeddingProvider) {
        this.embeddingProvider = embeddingProvider;
    }

    public String getEmbeddingModel() {
        return embeddingModel;
    }

    public void setEmbeddingModel(String embeddingModel) {
        this.embeddingModel = embeddingModel;
    }

    public String getEmbeddingApiKey() {
        return embeddingApiKey;
    }

    public void setEmbeddingApiKey(String embeddingApiKey) {
        this.embeddingApiKey = embeddingApiKey;
    }

    public int getEmbeddingDimensions() {
        return embeddingDimensions;
    }

    public void setEmbeddingDimensions(int embeddingDimensions) {
        this.embeddingDimensions = embeddingDimensions;
    }

    public IndexType getIndexType() {
        return indexType;
    }

    public void setIndexType(IndexType indexType) {
        this.indexType = indexType;
    }

    public DistanceMetric getDistanceMetric() {
        return distanceMetric;
    }

    public void setDistanceMetric(DistanceMetric distanceMetric) {
        this.distanceMetric = distanceMetric;
    }

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public String getUser() {
        return user;
    }

    public void setUser(String user) {
        this.user = user;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getSslMode() {
        return sslMode;
    }

    public void setSslMode(String sslMode) {
        this.sslMode = sslMode;
    }

    public String getSchema() {
        return schema;
    }

    public void setSchema(String schema) {
        this.schema = schema;
    }

    public String getTableName() {
        return tableName;
    }

    public void setTableName(String tableName) {
        this.tableName = tableName;
    }

    public int getM() {
        return m;
    }

    public void setM(int m) {
        this.m = m;
    }

    public int getEfConstruction() {
        return efConstruction;
    }

    public void setEfConstruction(int efConstruction) {
        this.efConstruction = efConstruction;
    }

    public int getEfSearch() {
        return efSearch;
    }

    public void setEfSearch(int efSearch) {
        this.efSearch = efSearch;
    }

    public String getAtlasCluster() {
        return atlasCluster;
    }

    public void setAtlasCluster(String atlasCluster) {
        this.atlasCluster = atlasCluster;
    }

    public int getNumCandidates() {
        return numCandidates;
    }

    public void setNumCandidates(int numCandidates) {
        this.numCandidates = numCandidates;
    }

    public String getIndexName() {
        return indexName;
    }

    public void setIndexName(String indexName) {
        this.indexName = indexName;
    }

    public int getPoolSize() {
        return poolSize;
    }

    public void setPoolSize(int poolSize) {
        this.poolSize = poolSize;
    }

    public int getTimeout() {
        return timeout;
    }

    public void setTimeout(int timeout) {
        this.timeout = timeout;
    }

    public int getEmbeddingBatchSize() {
        return embeddingBatchSize;
    }

    public void setEmbeddingBatchSize(int embeddingBatchSize) {
        this.embeddingBatchSize = embeddingBatchSize;
    }

    public int getLists() {
        return lists;
    }

    public void setLists(int lists) {
        this.lists = lists;
    }

    public int getProbes() {
        return probes;
    }

    public void setProbes(int probes) {
        this.probes = probes;
    }
}
