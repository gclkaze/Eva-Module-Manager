package utils

import "os"

func FolderExists(folder string) bool {
	_, err := os.Stat(folder)
	return !os.IsNotExist(err)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
