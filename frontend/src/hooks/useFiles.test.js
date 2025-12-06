import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useFiles } from '../useFiles';

describe('useFiles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('starts with empty files', () => {
    const { result } = renderHook(() => useFiles());

    expect(result.current.files).toEqual([]);
    expect(result.current.hasFiles).toBe(false);
    expect(result.current.fileCount).toBe(0);
  });

  it('adds files correctly', async () => {
    const { result } = renderHook(() => useFiles());

    const mockFile = new File(['# Hello'], 'test.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([mockFile]);
    });

    expect(result.current.files).toHaveLength(1);
    expect(result.current.files[0].name).toBe('test.md');
    expect(result.current.hasFiles).toBe(true);
    expect(result.current.fileCount).toBe(1);
  });

  it('adds multiple files with correct order', async () => {
    const { result } = renderHook(() => useFiles());

    const file1 = new File(['# First'], 'first.md', { type: 'text/markdown' });
    const file2 = new File(['# Second'], 'second.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([file1, file2]);
    });

    expect(result.current.files).toHaveLength(2);
    expect(result.current.files[0].order).toBe(0);
    expect(result.current.files[1].order).toBe(1);
  });

  it('removes file by id', async () => {
    const { result } = renderHook(() => useFiles());

    const mockFile = new File(['# Hello'], 'test.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([mockFile]);
    });

    const fileId = result.current.files[0].id;

    act(() => {
      result.current.removeFile(fileId);
    });

    expect(result.current.files).toHaveLength(0);
    expect(result.current.hasFiles).toBe(false);
  });

  it('clears all files', async () => {
    const { result } = renderHook(() => useFiles());

    const file1 = new File(['# First'], 'first.md', { type: 'text/markdown' });
    const file2 = new File(['# Second'], 'second.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([file1, file2]);
    });

    expect(result.current.files).toHaveLength(2);

    act(() => {
      result.current.clearFiles();
    });

    expect(result.current.files).toHaveLength(0);
  });

  it('selects file', async () => {
    const { result } = renderHook(() => useFiles());

    const mockFile = new File(['# Hello'], 'test.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([mockFile]);
    });

    const fileId = result.current.files[0].id;

    act(() => {
      result.current.selectFile(fileId);
    });

    expect(result.current.selectedFileId).toBe(fileId);
    expect(result.current.selectedFile).toBeTruthy();
    expect(result.current.selectedFile.name).toBe('test.md');
  });

  it('reorders files', async () => {
    const { result } = renderHook(() => useFiles());

    const file1 = new File(['# A'], 'a.md', { type: 'text/markdown' });
    const file2 = new File(['# B'], 'b.md', { type: 'text/markdown' });
    const file3 = new File(['# C'], 'c.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([file1, file2, file3]);
    });

    const ids = result.current.files.map((f) => f.id);

    // Move first to last
    act(() => {
      result.current.reorderFiles(ids[0], ids[2]);
    });

    const ordered = result.current.getOrderedFiles();
    expect(ordered[0].name).toBe('b.md');
    expect(ordered[1].name).toBe('c.md');
    expect(ordered[2].name).toBe('a.md');
  });

  it('getOrderedFiles returns sorted by order', async () => {
    const { result } = renderHook(() => useFiles());

    const file1 = new File(['# A'], 'a.md', { type: 'text/markdown' });
    const file2 = new File(['# B'], 'b.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([file1, file2]);
    });

    const ordered = result.current.getOrderedFiles();

    // Should be in order 0, 1
    for (let i = 0; i < ordered.length - 1; i++) {
      expect(ordered[i].order).toBeLessThan(ordered[i + 1].order);
    }
  });

  it('auto-selects first file when added', async () => {
    const { result } = renderHook(() => useFiles());

    const mockFile = new File(['# Hello'], 'test.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([mockFile]);
    });

    expect(result.current.selectedFileId).toBe(result.current.files[0].id);
  });

  it('clears selection when all files removed', async () => {
    const { result } = renderHook(() => useFiles());

    const mockFile = new File(['# Hello'], 'test.md', { type: 'text/markdown' });

    await act(async () => {
      await result.current.addFiles([mockFile]);
    });

    act(() => {
      result.current.clearFiles();
    });

    expect(result.current.selectedFileId).toBeNull();
    expect(result.current.selectedFile).toBeNull();
  });
});
