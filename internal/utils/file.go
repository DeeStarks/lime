package utils

import "path/filepath"

func GetFileExtension(fileName string) string {
	LogMessage("GetFileExtension: "+filepath.Ext(fileName))
	return filepath.Ext(fileName)
}