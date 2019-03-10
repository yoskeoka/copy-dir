package cpdir

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirContents copy directory contents
func CopyDirContents(src, dest string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcinfo.IsDir() {
		return fmt.Errorf("src %s is not direcotory", src)
	}

	destinfo, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !destinfo.IsDir() {
		return fmt.Errorf("dest %s is not direcotory", dest)
	}

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if src == path {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(dest, relPath), os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}

		err = CopyFile(path, filepath.Join(dest, relPath))
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

// CopyFile copy file
func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
