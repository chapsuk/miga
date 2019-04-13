package converter

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"miga/logger"
	"miga/utils"
)

type MigrateFormatter struct{}

func (f *MigrateFormatter) Read(src string) ([]Task, error) {
	files, err := filepath.Glob(src + "/**.sql")
	if err != nil {
		return nil, err
	}

	plainTasks := make(map[string]*Task, len(files))
	for _, file := range files {
		filename := filepath.Base(file)

		vidx := strings.Index(filename, "_")
		if vidx < 0 {
			logger.G().Warnf("(skip) incorrect migarte migration file name: %s", filename)
			continue
		}

		version, err := strconv.Atoi(filename[:vidx])
		if err != nil || version <= 0 {
			logger.G().Warnf("(skip) incorrect migarte migration file name: %s error: %s", filename, err)
			continue
		}

		nidx := strings.Index(filename, ".")
		if nidx < 0 {
			logger.G().Warnf("(skip) incorrect migarte migration file name: %s", filename)
			continue
		}

		name := filename[vidx+1 : nidx]
		if len(name) == 0 {
			logger.G().Warnf("(skip) wrong migrate migration file name: %s", filename)
			continue
		}

		tidx := strings.Index(filename[nidx+1:], ".")
		if tidx < 0 {
			logger.G().Warnf("(skip) missing migrate migration type: %s", filename)
			continue
		}

		if _, ok := plainTasks[name]; !ok {
			plainTasks[name] = &Task{
				Name:    name,
				Version: version,
			}
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		mtype := filename[nidx+1 : nidx+tidx+1]
		switch mtype {
		case "up":
			plainTasks[name].Up = b
		case "down":
			plainTasks[name].Down = b
		default:
			logger.G().Warnf("(skip) incorrect migrate migration type: %s", mtype)
			continue
		}
	}

	result := make([]Task, 0, len(plainTasks))
	for _, t := range plainTasks {
		result = append(result, *t)
	}

	return result, nil
}

func (f *MigrateFormatter) Write(dest string, tasks []Task) error {
	for _, t := range tasks {
		upf, dwn, err := utils.CreateMigrationsFiles(int64(t.Version), dest, t.Name, "sql")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(upf, t.Up, 0755)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(dwn, t.Down, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
