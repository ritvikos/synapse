// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package spooler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

// Manages batch directories and size.
type batcher struct {
	currentDir  string
	config      BatchConfig
	currentSize atomic.Int64
}

// TODO: Resume batching logic
func newBatcher(config BatchConfig) (*batcher, error) {
	var err error

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("batcher: batcher config: %w", err)
	}

	currentDir, err := newBatch(config.BaseDir)
	if err != nil {
		return nil, err
	}

	batcher := &batcher{
		config:     config,
		currentDir: currentDir,
	}
	batcher.currentSize.Store(0)

	return batcher, nil
}

// Current batch size (in bytes)
func (b *batcher) Size() int64 {
	return b.currentSize.Load()
}

// Current batch directory
func (b *batcher) CurrentDir() string {
	return b.currentDir
}

// Increment current batch size (in bytes)
func (b *batcher) AddBytes(delta int) {
	b.currentSize.Add(int64(delta))
}

// Rotate the batch if it exceeds the maximum size
func (b *batcher) Rotate() error {
	currentDir := b.currentDir
	currentSize := b.currentSize.Load()

	if currentSize <= b.config.MaxBatchSize {
		return nil
	}

	// TODO: Make this robust
	if b.config.Processor.Async {
		go b.processBatchWorker(currentDir, currentSize)
	} else {
		if err := b.config.Processor.Hooks.OnBatchReady(currentDir, currentSize); err != nil {
			return err
		}

		if b.config.Processor.DeleteSource {
			if err := os.RemoveAll(currentDir); err != nil {
				return fmt.Errorf("batcher: failed to delete batch directory: %w", err)
			}
		}
	}

	return b.reset()
}

// Create a new batch directory
func newBatch(baseDir string) (string, error) {
	// Generate 'n' cryptographically secure random bytes, with a maximum of 'maxRetries' reattempts in case of failure.
	// Returns an error if 'n' is equal to zero, or if the maximum retries are exceeded.
	generateRandomBytes := func(n uint, maxRetries uint) ([]byte, error) {
		if n == 0 {
			return nil, errors.New("CSPRNG: 'n' must be greater than zero")
		}
		b := make([]byte, n)
		for range maxRetries {
			if _, err := rand.Read(b); err != nil {
				continue
			}
			return b, nil
		}
		return nil, fmt.Errorf("CSPRNG: max retries exceeded (%d)", maxRetries)
	}

	createFileID := func() (string, error) {
		b, err := generateRandomBytes(4, 3)
		if err != nil {
			return "", fmt.Errorf("batcher: failed to create file ID: %w", err)
		}
		return hex.EncodeToString(b), nil
	}

	fileID, err := createFileID()
	if err != nil {
		// TODO: Decide fallback approach:
		// 1. Generate a new ID with non-cryptographic RNG
		// 2. Fail the batch creation and return the error
		return "", err
	}

	batchDir := filepath.Join(baseDir, "batch-"+fileID)
	if err := os.MkdirAll(batchDir, filePermissions); err != nil {
		return "", fmt.Errorf("batcher: failed to create batch directory: %w", err)
	}
	return batchDir, nil
}

// Reset the batcher to start a new batch
func (b *batcher) reset() error {
	currentDir, err := newBatch(b.config.BaseDir)
	if err != nil {
		return err
	}
	b.currentDir = currentDir
	b.currentSize.Store(0)
	return nil
}

func (b *batcher) processBatchWorker(dir string, size int64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "batch worker: panic in async batch processing: %v\n", r)
		}
	}()

	if err := b.config.Processor.Hooks.OnBatchReady(dir, size); err != nil {
		b.config.Processor.Hooks.OnBatchError(dir, &err)
		return
	}

	if b.config.Processor.DeleteSource {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Fprintf(os.Stderr, "batch worker: failed to delete batch directory %s: %v\n", dir, err)
		}
	}
}
