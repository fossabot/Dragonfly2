/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storedriver

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"d7y.io/dragonfly/v2/pkg/unit"
)

// Driver defines an interface to manage the data stored in the driver.
//
// NOTE:
// It is recommended that the lock granularity of the driver should be in piece.
// That means that the storage driver could read and write
// the different pieces of the same file concurrently.
type Driver interface {

	// Get data from the storage based on raw information.
	// If the length<=0, the driver should return all data from the raw.offset.
	// Otherwise, just return the data which starts from raw.offset and the length is raw.length.
	Get(ctx context.Context, raw *Raw) (io.ReadCloser, error)

	// Get data from the storage based on raw information.
	// The data should be returned in bytes.
	// If the length<=0, the storage driver should return all data from the raw.offset.
	// Otherwise, just return the data which starts from raw.offset and the length is raw.length.
	GetBytes(ctx context.Context, raw *Raw) ([]byte, error)

	// Put the data into the storage with raw information.
	// The storage will get data from io.Reader as io stream.
	// If the offset>0, the storage driver should starting at byte raw.offset off.
	Put(ctx context.Context, raw *Raw, data io.Reader) error

	// PutBytes puts the data into the storage with raw information.
	// The data is passed in bytes.
	// If the offset>0, the storage driver should starting at byte raw.offset off.
	PutBytes(ctx context.Context, raw *Raw, data []byte) error

	// Remove the data from the storage based on raw information.
	Remove(ctx context.Context, raw *Raw) error

	// Stat determines whether the data exists based on raw information.
	// If that, and return some info that in the form of struct StorageInfo.
	// If not, return the ErrFileNotExist.
	Stat(ctx context.Context, raw *Raw) (*StorageInfo, error)

	// GetAvailSpace returns the available disk space in B.
	GetAvailSpace(ctx context.Context) (unit.Bytes, error)

	// GetTotalAndFreeSpace
	GetTotalAndFreeSpace(ctx context.Context) (unit.Bytes, unit.Bytes, error)

	// GetTotalSpace
	GetTotalSpace(ctx context.Context) (unit.Bytes, error)

	// Walk walks the file tree rooted at root which determined by raw.Bucket and raw.Key,
	// calling walkFn for each file or directory in the tree, including root.
	Walk(ctx context.Context, raw *Raw) error

	// CreateBaseDir
	CreateBaseDir(ctx context.Context) error

	// GetPath
	GetPath(raw *Raw) string

	// MoveFile
	MoveFile(src string, dst string) error

	// Exits
	Exits(ctx context.Context, raw *Raw) bool

	// GetHomePath
	GetHomePath(ctx context.Context) string

	// GetGcConfig
	GetGcConfig(ctx context.Context) *GcConfig
}

// Raw identifies a piece of data uniquely.
// If the length<=0, it represents all data.
type Raw struct {
	Bucket    string
	Key       string
	Offset    int64
	Length    int64
	Trunc     bool
	TruncSize int64
	Append    bool
	WalkFn    filepath.WalkFunc
}

// StorageInfo includes partial meta information of the data.
type StorageInfo struct {
	Path       string    // file path
	Size       int64     // file size
	CreateTime time.Time // create time
	ModTime    time.Time // modified time
}

// GcConfig
type GcConfig struct {
	YoungGCThreshold  unit.Bytes    `yaml:"youngGCThreshold"`
	FullGCThreshold   unit.Bytes    `yaml:"fullGCThreshold"`
	CleanRatio        int           `yaml:"cleanRatio"`
	IntervalThreshold time.Duration `yaml:"intervalThreshold"`
}
