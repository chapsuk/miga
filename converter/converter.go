package converter

import (
	"fmt"
	"os"

	"github.com/chapsuk/miga/driver"
)

type Formatter interface {
	Read(src string) ([]Task, error)
	Write(dest string, tasks []Task) error
}

type Task struct {
	Up      []byte
	Down    []byte
	Name    string
	Version int
}

func Convert(from, to, src, dest string) error {
	f, err := createFormatter(from)
	if err != nil {
		return err
	}

	t, err := createFormatter(to)
	if err != nil {
		return err
	}

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return err
	}

	tasks, err := f.Read(src)
	if err != nil {
		return err
	}

	_, err = os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.Mkdir(dest, 0755)
		if err != nil {
			return err
		}
	}

	return t.Write(dest, tasks)
}

func createFormatter(name string) (Formatter, error) {
	switch name {
	case driver.Goose:
		return &GooseFormatter{}, nil
	case driver.Migrate:
		return &MigrateFormatter{}, nil
	case driver.Impg:
		return &MigrateFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported driver name %s", name)
	}
}
