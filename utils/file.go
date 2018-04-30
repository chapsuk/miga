package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/chapsuk/miga/logger"
)

func CreateMigrationsFiles(dir, name, ext string) error {
	timestamp := time.Now().Unix()
	base := fmt.Sprintf("%v/%v_%v.", dir, timestamp, name)
	os.MkdirAll(dir, os.ModePerm)

	err := createFile(base + "up." + ext)
	if err != nil {
		return err
	}

	return createFile(base + "down." + ext)
}

func createFile(fname string) error {
	_, err := os.Create(fname)
	if err != nil {
		return err
	}
	logger.G().Infof("Create migrations file: %s", fname)
	return nil
}
