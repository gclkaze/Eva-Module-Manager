package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FolderExists(folder string) bool {
	info, err := os.Stat(folder)
	return !os.IsNotExist(err) && info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func CreateFolder(folder string) error {
	if folder == "." {
		return nil
	}
	err := os.Mkdir(folder, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
	}
	return err
}

func UntarGzToDir(srcTarGz, destDir string) error {
	f, err := os.Open(srcTarGz)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	destDir, err = filepath.Abs(destDir)
	if err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Clean and secure path
		relPath := filepath.Clean(hdr.Name)
		if strings.HasPrefix(relPath, "..") {
			return errors.New("tar entry attempts path traversal: " + hdr.Name)
		}

		target := filepath.Join(destDir, relPath)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			out, err := os.OpenFile(
				target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				os.FileMode(hdr.Mode),
			)
			if err != nil {
				return err
			}

			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}

	return nil
}

func DeleteFile(p string) error {
	return os.Remove(p)
}

func IsFolderEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read at most one entry
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil // empty
	}
	return false, err // not empty or error
}

func FolderIsEmpty(path string) bool {
	if !FolderExists(path) {
		return true
	}

	res, err := IsFolderEmpty(path)
	if err != nil {
		return true
	}
	return res
}

func CleanFolder(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		err := os.RemoveAll(path) // Removes file or directory recursively
		if err != nil {
			return err
		}
	}
	return nil
}
func RemoveFolder(dir string) error {
	return os.Remove(dir)
}

func CleanAndRemoveFolder(dir string) error {
	err := CleanFolder(dir)
	if err != nil {
		return err
	}
	return RemoveFolder(dir)
}
