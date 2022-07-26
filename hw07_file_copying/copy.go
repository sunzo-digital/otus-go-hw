package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrUnknownFileSize       = errors.New("unknown size of the file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	srcFileSize := srcFileInfo.Size()

	if srcFileSize == 0 {
		return ErrUnknownFileSize
	}

	if offset > srcFileSize {
		return ErrOffsetExceedsFileSize
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	_, err = srcFile.Seek(offset, io.SeekStart)

	if err != nil {
		return err
	}

	if limit == 0 || limit > srcFileSize {
		limit = srcFileSize
	}

	bar := pb.Full.Start64(limit - offset)

	barReader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, barReader, limit)

	bar.Finish()

	return err
}
