package utils

import (
	"fmt"
	"os"

	"miga/logger"
)

func CreateMigrationsFiles(
	version int64,
	dir, name, ext string,
) (upFileName, downFileName string, err error) {
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}

	upFileName = fmt.Sprintf("%v%v_%v.up.%s", dir, version, name, ext)
	downFileName = fmt.Sprintf("%v%v_%v.down.%s", dir, version, name, ext)

	os.MkdirAll(dir, os.ModePerm)

	err = createFile(upFileName)
	if err != nil {
		return
	}

	err = createFile(downFileName)
	return
}

func createFile(fname string) error {
	_, err := os.Create(fname)
	if err != nil {
		return err
	}
	logger.G().Infof("Create migrations file: %s", fname)
	return nil
}
