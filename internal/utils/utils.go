package utils

import "os"

func FolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
