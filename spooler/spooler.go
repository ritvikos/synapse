// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package spooler

import (
	"context"
	"errors"
)

// TODO:
// - Lightweight embedded database for book-keeping config state and current working batch directory, for recovery after crashes or restarts.

// Spooler manages local batches of files for remote storage or local persistence.
//
// It accumulates files into batches up to a configured size limit. When a batch
// reaches its size threshold, it triggers a rotation and invokes the configured hooks
// for processing, if any.
//
// # Example Usage
//
//	config := spooler.SpoolerConfig{
//		FileWriterConfig: spooler.FileWriterConfig{
//			MaxFileSize: 5 * 1024 * 1024, // 5 MB per file
//		},
//		BatchConfig: spooler.BatchConfig{
//			BaseDir:      "/var/spool/synapse",
//			MaxBatchSize: 1 * 1024 * 1024 * 1024, // 1 GB per batch
//			Processor: spooler.BatchProcessor{
//				Async:        true,
//				DeleteSource: true,
//				Hooks: &spooler.BatchHooks{
//					OnBatchReady: func(batchDir string, totalBytes int64) error {
//						return nil
//					},
//					OnBatchError: func(batchDir string, err *error) {
//						// Handle processing errors
//					},
//				},
//			},
//		},
//	}
//
//	s, err := spooler.NewSpooler(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// // Create a new file writer
//
//	if err := s.NewWriter(ctx, "output.txt"); err != nil {
//		log.Fatal(err)
//	}
//
// //  Write data chunks
//
//	if err := s.WriteChunk([]byte("some data...")); err != nil {
//		log.Fatal(err)
//	}
//
// // Commit the current file into the batch.
//
//	if err := s.Commit(); err != nil {
//		log.Fatal(err)
//	}
type Spooler struct {
	config        SpoolerConfig
	writer        writerFactory
	batcher       *batcher
	currentWriter *fileWriter
}

// Create a new Spooler instance
func NewSpooler(config SpoolerConfig) (*Spooler, error) {
	batcher, err := newBatcher(config.BatchConfig)
	if err != nil {
		return nil, err
	}

	writerFactory, err := newWriterFactory(config.FileWriterConfig)
	if err != nil {
		return nil, err
	}

	return &Spooler{
		config:  config,
		writer:  writerFactory,
		batcher: batcher,
	}, nil
}

// Current batch size (in bytes)
func (s *Spooler) BatchSize() int {
	return int(s.batcher.Size())
}

// Create a new file writer
func (s *Spooler) NewWriter(ctx context.Context, fileName string) error {
	if s.currentWriter != nil {
		if err := s.Commit(); err != nil {
			return err
		}
	}

	writer, err := s.writer.NewFileWriter(s.batcher.CurrentDir(), fileName)
	if err != nil {
		return err
	}

	s.currentWriter = writer
	return nil
}

// Write a chunk of data to the current file writer
func (s *Spooler) WriteChunk(data []byte) error {
	var err error

	if err = s.currentWriter.Write(data); err != nil {
		if errors.Is(err, ErrWriteLimitReached) {
			if err = s.currentWriter.Abort(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Commit the current file writer
func (s *Spooler) Commit() error {
	var err error

	totalWritten, err := s.currentWriter.Commit()
	if err != nil {
		return err
	}
	s.currentWriter = nil

	s.batcher.AddBytes(totalWritten)
	if err = s.batcher.Rotate(); err != nil {
		return err
	}

	return nil
}
