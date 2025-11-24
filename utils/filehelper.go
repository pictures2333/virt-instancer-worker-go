package utils

import (
	"fmt"
	"io"
	"os"
)

func CheckDir(path string) (err error) {
	var info os.FileInfo

	if info, err = os.Stat(path); err != nil {
		return err
	}

	// no err, continue
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}

func MustMkdir(path string) (err error) {
	var info os.FileInfo

	if info, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// if not exists -> mkdir
			return os.Mkdir(path, 0o770)
		} else {
			// other err
			return err
		}
	}

	// no err -> exists
	if !info.IsDir() {
		return os.Mkdir(path, 0o770)
	}

	return nil
}

func MustRmdir(path string) (err error) {
	return os.RemoveAll(path)
}

func CopyFile(srcPath string, dstPath string) (err error) {
	var (
		src *os.File
		dst *os.File
	)

	if src, err = os.Open(srcPath); err != nil {
		return err
	}
	defer src.Close()

	if dst, err = os.Create(dstPath); err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return dst.Sync()
}
