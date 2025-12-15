# Spooler

## Purpose
It buffers data streams locally, organizes them into batches with atomic writes, prior to processing and remote storage persistence.

Internally it groups files into directories in `batch-<id>` format into base directory with configurable threshold at file and batch levels, atomically commit files to the currently active batch, tracks batch size, rotate batches when thresholds are met, and provide read access to committed batches for processing.
