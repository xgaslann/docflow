# Contributing to DocFlow

Thanks for wanting to contribute! Here's how to do it properly.

## Ground Rules

1. **Tests are mandatory.** No exceptions. If your PR doesn't include tests, it won't be merged.
2. **Keep it simple.** Don't over-engineer. If a feature needs 500 lines of code, think again.
3. **Follow the existing style.** Look at the codebase before writing code.

## Repository Structure

```
docflow/
├── app/                    # Web Application (Go + React)
│   ├── backend/           # Go API server
│   └── frontend/          # React UI
├── sdks/                   # Standalone SDKs
│   ├── go/                # Go SDK
│   │   └── docflow/       # Core modules
│   ├── python/            # Python SDK
│   │   └── docflow/       # Core modules
│   └── java/              # Java SDK
│       └── src/main/java/ # Core modules
└── examples/               # Usage examples
```

## Getting Started

### 1. Fork & Clone

```bash
git clone https://github.com/YOUR_USERNAME/docflow.git
cd docflow
```

### 2. Set Up Development Environment

**Go SDK:**
```bash
cd sdks/go
go mod tidy
go test ./...
```

**Python SDK:**
```bash
cd sdks/python
pip install -e ".[dev]"
pytest
```

**Java SDK:**
```bash
cd sdks/java
mvn clean compile
mvn test
```

**Web App:**
```bash
cd app/backend && go mod tidy
cd ../frontend && npm install
```

### 3. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

## Code Style

### Go

- Use `gofmt` - no exceptions
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Keep functions small (<50 lines ideally)
- Error messages should be lowercase, no punctuation
- Add comments for exported functions

```go
// Good
func (p *Processor) ProcessFile(ctx context.Context, path string) (*Result, error) {
    if path == "" {
        return nil, errors.New("empty path")
    }
    // ...
}
```

### Python

- Follow PEP 8
- Use type hints for public APIs
- Use dataclasses or Pydantic for models
- Keep functions focused

```python
# Good
def process_file(path: str, config: RAGConfig) -> RAGDocument:
    """Process a file and return a RAG document."""
    if not path:
        raise ValueError("path cannot be empty")
    # ...
```

### Java

- Follow Google Java Style Guide
- Use meaningful variable names
- Add Javadoc for public methods
- Prefer composition over inheritance

```java
// Good
/**
 * Process a file and return a RAG document.
 * @param path Path to the file
 * @return Processed RAG document
 * @throws IOException if file cannot be read
 */
public RAGDocument processFile(String path) throws IOException {
    if (path == null || path.isEmpty()) {
        throw new IllegalArgumentException("path cannot be empty");
    }
    // ...
}
```

## Writing Tests

### Go Tests

```go
func TestRAGProcessor_ProcessFile(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid file", "testdata/sample.md", false},
        {"empty path", "", true},
        {"non-existent", "testdata/missing.md", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := NewRAGProcessor(DefaultConfig())
            _, err := p.ProcessFile(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Python Tests

```python
import pytest
from docflow.rag import RAGProcessor, RAGConfig

class TestRAGProcessor:
    def test_process_file_success(self, tmp_path):
        file = tmp_path / "test.md"
        file.write_text("# Hello World")
        
        processor = RAGProcessor(RAGConfig())
        result = processor.process_file(str(file))
        
        assert result.content == "# Hello World"
        assert len(result.chunks) > 0
    
    def test_process_file_empty_path(self):
        processor = RAGProcessor(RAGConfig())
        with pytest.raises(ValueError):
            processor.process_file("")
```

### Java Tests

```java
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

class RAGProcessorTest {
    private RAGProcessor processor;
    
    @BeforeEach
    void setUp() {
        processor = new RAGProcessor(RAGConfig.defaultConfig());
    }
    
    @Test
    void processFile_validFile_returnsDocument() throws Exception {
        RAGDocument doc = processor.processFile("src/test/resources/sample.md");
        assertNotNull(doc);
        assertFalse(doc.getChunks().isEmpty());
    }
    
    @Test
    void processFile_emptyPath_throwsException() {
        assertThrows(IllegalArgumentException.class, () -> {
            processor.processFile("");
        });
    }
}
```

## Pull Request Process

### 1. Before Submitting

- [ ] All tests pass across all SDKs you modified
- [ ] Code is formatted (gofmt, black, google-java-format)
- [ ] No debug statements or commented-out code
- [ ] Documentation updated if needed
- [ ] Commit messages are clear

### 2. PR Title Format

```
feat(go): add vector store support
fix(python): handle unicode in filenames
docs: update RAG pipeline documentation
refactor(java): simplify batch processor
test(all): increase coverage for chunker
```

### 3. PR Description Template

```markdown
## What
Brief description of what this PR does.

## Why
Why is this change needed?

## How
How did you implement it? Any notable design decisions?

## Testing
How did you test this? What test cases did you add?

## SDK Changes
- [ ] Go SDK modified
- [ ] Python SDK modified  
- [ ] Java SDK modified

## Breaking Changes
List any breaking changes (if applicable).
```

### 4. Review Process

1. Submit PR
2. CI runs tests for all SDKs
3. Maintainer reviews
4. Address feedback
5. Get merged

## Adding New Features

### Adding to All SDKs

When adding a new feature, implement it in all three SDKs:

1. **Design first** - Create an issue describing the API
2. **Python first** - Usually easiest to prototype
3. **Go second** - Port with idiomatic Go patterns
4. **Java third** - Port with enterprise patterns
5. **Update docs** - Add to all README files

### Adding a New Format Converter

1. Create `formats/new_format.go/py/java`
2. Implement `toMarkdown(data, filename)` method
3. Add tests with sample files
4. Register in format detection logic
5. Update feature matrix in READMEs

### Adding a New Vector Store

1. Create `storage/vector/new_store.go/py/java`
2. Implement the VectorStore interface
3. Add configuration class
4. Add integration tests
5. Document environment variables

## Questions?

Open an issue with the `question` label.

## Recognition

Contributors are added to the README. Significant contributors get mentioned in release notes.

---

Keep it clean, test your code, and have fun!
