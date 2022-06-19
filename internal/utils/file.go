package utils

import "path/filepath"

func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}
