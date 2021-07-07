// +build linux darwin

package common

import (
	"golang.org/x/sys/unix"
)

type linuxStorageAPI struct{}

func platformStorageApi() iStorageUtils {
	return &linuxStorageAPI{}
}

// GetFreeSpace returns the number of bytes free on a volume based on path
func (l *linuxStorageAPI) GetFreeSpace(dataDir string) (uint64, error) {

	var stat unix.Statfs_t

	err := unix.Statfs(dataDir, &stat)
	if err != nil {
		return 0, err
	}

	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize), nil
}
