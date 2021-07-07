package common

import (
	"golang.org/x/sys/windows"
)

type windowsStorageAPI struct{}

func platformStorageApi() iStorageUtils {
	return &windowsStorageAPI{}
}

// GetFreeSpace returns the number of bytes free on a volume based on path
func (w *windowsStorageAPI) GetFreeSpace(dataDir string) (uint64, error) {

	var free, total, avail uint64
	pathPtr, err := windows.UTF16PtrFromString(dataDir)
	if err != nil {
		return 0, err
	}

	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)
	return free, err
}
