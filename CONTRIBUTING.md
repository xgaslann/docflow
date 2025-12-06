# Contributing to DocFlow

Thanks for wanting to contribute. Here's how to do it properly.

## Ground Rules

1. **Tests are mandatory.** No exceptions. If your PR doesn't include tests, it won't be merged.
2. **Keep it simple.** Don't over-engineer. If a feature needs 500 lines of code, think again.
3. **Follow the existing style.** Look at the codebase before writing code.

## Getting Started

### 1. Fork & Clone

```bash
git clone https://github.com/YOUR_USERNAME/docflow.git
cd docflow
```

### 2. Set Up Development Environment

```bash
# Backend
cd backend
go mod tidy

# Frontend
cd ../frontend
npm install
```

### 3. Run Tests

Before making any changes, make sure existing tests pass:

```bash
# Backend tests
cd backend
go test -v ./...

# Frontend tests
cd frontend
npm test
```

### 4. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

## Code Style

### Go (Backend)

- Use `gofmt` - no exceptions
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Keep functions small (< 50 lines ideally)
- Error messages should be lowercase, no punctuation
- Add comments for exported functions

```go
// Good
func (s *Service) ProcessFile(ctx context.Context, data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }
    // ...
}

// Bad
func (s *Service) ProcessFile(ctx context.Context, data []byte) error {
    if len(data) == 0 {
        return errors.New("Empty data.") // wrong style
    }
    // ...
}
```

### JavaScript/React (Frontend)

- Use functional components with hooks
- Use named exports (not default) for components
- Keep components focused - one job per component
- Use descriptive variable names

```jsx
// Good
export function FileUploader({ onFilesAdded, disabled }) {
  const handleDrop = useCallback((files) => {
    onFilesAdded(files);
  }, [onFilesAdded]);
  
  return <Dropzone onDrop={handleDrop} disabled={disabled} />;
}

// Bad
export default function FU({ cb, d }) {
  return <Dropzone onDrop={cb} disabled={d} />;
}
```

## Writing Tests

### Backend Tests

Put tests next to the code they test: `service.go` → `service_test.go`

```go
package service

import (
    "testing"
)

func TestMarkdownToHTML(t *testing.T) {
    svc := NewMarkdownService()
    
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple heading",
            input:    "# Hello",
            expected: "<h1>Hello</h1>",
        },
        {
            name:     "paragraph",
            input:    "Some text",
            expected: "<p>Some text</p>",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := svc.ToHTML(tt.input)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if !strings.Contains(result, tt.expected) {
                t.Errorf("expected %q to contain %q", result, tt.expected)
            }
        })
    }
}
```

### Frontend Tests

Use Vitest. Test files go in `__tests__` folders or with `.test.jsx` suffix.

```jsx
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { FileUploader } from '../FileUploader';

describe('FileUploader', () => {
  it('calls onFilesAdded when files are dropped', async () => {
    const onFilesAdded = vi.fn();
    render(<FileUploader onFilesAdded={onFilesAdded} />);
    
    // ... test implementation
  });
  
  it('shows disabled state correctly', () => {
    render(<FileUploader onFilesAdded={vi.fn()} disabled />);
    // ... assertions
  });
});
```

### What to Test

**Do test:**
- Business logic
- Edge cases
- Error handling
- User interactions

**Don't test:**
- External libraries
- Simple getters/setters
- CSS styling

## Pull Request Process

### 1. Before Submitting

- [ ] All tests pass (`go test ./...` and `npm test`)
- [ ] Code is formatted (`gofmt` and Prettier)
- [ ] No console.logs or debug statements
- [ ] Commit messages are clear

### 2. PR Title Format

```
feat: add batch processing support
fix: handle empty files correctly
docs: update API documentation
refactor: simplify conversion logic
test: add handler tests
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

## Screenshots (if UI changes)

Before/After screenshots if applicable.
```

### 4. Review Process

1. Submit PR
2. CI runs tests
3. Maintainer reviews
4. Address feedback
5. Get merged

PRs are usually reviewed within a few days. If it's been a week, feel free to ping.

## Commit Messages

Keep them short and descriptive. Use present tense.

```
# Good
feat: add PDF merge preview
fix: handle unicode in filenames
refactor: extract PDF service

# Bad
Fixed the bug
WIP
asdfasdf
```

## Project Structure Guidelines

### Adding a New Feature

1. **Backend changes?** Add to appropriate service in `internal/service/`
2. **New endpoint?** Add handler in `internal/handler/`, register in routes
3. **Frontend changes?** Create component in `components/`, hook in `hooks/`
4. **Always** add tests for new code

### File Naming

```
# Backend (Go)
service.go
service_test.go

# Frontend (React)
ComponentName/
  index.jsx
  ComponentName.test.jsx
```

## Questions?

Open an issue with the `question` label. Don't be shy.

## Recognition

Contributors are added to the README. Significant contributors get mentioned in release notes.

---

That's it. Keep it clean, test your code, and have fun.
