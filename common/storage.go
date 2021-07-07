package common

import (
	"io/ioutil"
	"os"
	"path"
)

type iStorageUtils interface {
	GetFreeSpace(path string) (uint64, error)
}

var platformApi = platformStorageApi()

// GetFreeSpace returns total number of free bytes on a path's volume
func GetFreeSpace(dataDir string) (uint64, error) {
	return platformApi.GetFreeSpace(dataDir)
}

// ValidatePath validates file permissions for the current user for the directory in the provided path
func ValidatePath(dataDir string) bool {

	// check path is valid os path
	pathInfo, err := os.Stat(dataDir)
	if err != nil {
		println("Invalid target directory. Please provide a valid directory")
		return false
	}

	// check path is a dir
	if !pathInfo.IsDir() {
		println("Invalid target directory. Please provide a valid directory")
		return false
	}

	tempFilePath := path.Join(dataDir, "temp.bin")

	// check write file perms
	err = ioutil.WriteFile(tempFilePath, []byte{0xff}, 0644)
	if err != nil {
		println("You don't have write permissions for this directory. Please enter a directory you have permissions to write to")
		return false
	}

	// check read file perms
	_, err = ioutil.ReadFile(tempFilePath)
	if err != nil {
		println("You don't have read permissions to read from this directory. Please enter a directory you have permissions to read from.")
		return false
	}

	// check delete file perms
	err = os.Remove(tempFilePath)
	if err != nil {
		println("You don't have permissions to delete files in this directory. Please enter a directory you have permissions to delete files from")
		return false
	}

	return true
}
