// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// +build !windows

package filestore

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func isBusy(err error) bool {
	err = underlyingError(err)
	return err == unix.EBUSY
}

func diskInfoFromPath(path string) (info DiskInfo, err error) {
	var stat unix.Statfs_t
	err = unix.Statfs(path, &stat)
	if err != nil {
		return DiskInfo{"", -1}, err
	}

	// the Bsize size depends on the OS and unconvert gives a false-positive
	availableSpace := int64(stat.Bavail) * int64(stat.Bsize) //nolint
	filesystemID := fmt.Sprintf("%08x%08x", stat.Fsid.Val[0], stat.Fsid.Val[1])

	return DiskInfo{filesystemID, availableSpace}, nil
}

// rename renames oldpath to newpath
func rename(oldpath, newpath string) error {
	inputFile, err := os.Open(oldpath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(newpath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(oldpath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

// openFileReadOnly opens the file with read only
func openFileReadOnly(path string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, perm)
}
