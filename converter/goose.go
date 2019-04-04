package converter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"miga/logger"
)

var (
	gooseUpPrefix   = "-- +goose Up\n"
	gooseDownPrefix = "-- +goose Down\n"

	gooseStatementBegin = "-- +goose StatementBegin\n"
	gooseStatementEnd   = "-- +goose StatementEnd\n"
)

type GooseFormatter struct{}

func (f *GooseFormatter) Read(src string) ([]Task, error) {
	files, err := filepath.Glob(src + "/**.sql")
	if err != nil {
		return nil, err
	}

	result := make([]Task, 0, len(files))
	for _, file := range files {
		filename := filepath.Base(file)

		vidx := strings.Index(filename, "_")
		if vidx < 0 {
			logger.G().Warnf("(skip) incorrect goose migration file name: %s", filename)
			continue
		}

		version, err := strconv.Atoi(filename[:vidx])
		if err != nil || version <= 0 {
			logger.G().Warnf("(skip) incorrect goose migration file name: %s error: %s", filename, err)
			continue
		}

		nidx := strings.Index(filename, ".")
		if nidx < 0 {
			logger.G().Warnf("(skip) incorrect goose migration file name: %s", filename)
			continue
		}

		task := Task{
			Name:    filename[vidx+1 : nidx],
			Version: version,
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		p := bytes.Split(b, []byte(gooseDownPrefix))
		switch len(p) {
		case 2:
			task.Down = removeGooseArtifacts(p[1])
			fallthrough
		case 1:
			task.Up = removeGooseArtifacts(p[0])
		default:
			logger.G().Errorf("incorrect goose migration file %s body: %s", file, b)
			continue
		}

		result = append(result, task)
	}

	return result, nil
}

func (f *GooseFormatter) Write(dest string, tasks []Task) error {
	for _, t := range tasks {
		filename := fmt.Sprintf("%05v_%s.sql", t.Version, t.Name)

		body := bytes.NewBuffer([]byte(gooseUpPrefix))
		body.Write([]byte(gooseStatementBegin))
		body.Write(t.Up)
		body.Write([]byte("\n"))
		body.Write([]byte(gooseStatementEnd))
		body.Write([]byte(gooseDownPrefix))
		body.Write([]byte(gooseStatementBegin))
		body.Write(t.Down)
		body.Write([]byte("\n"))
		body.Write([]byte(gooseStatementEnd))

		fpath := dest + "/" + filename
		err := ioutil.WriteFile(fpath, body.Bytes(), 0755)
		if err != nil {
			return err
		}
		logger.G().Infof("%s file created", fpath)
	}
	return nil
}

func removeGooseArtifacts(body []byte) []byte {
	body = bytes.Replace(body, []byte(gooseUpPrefix), []byte(""), 1)
	body = bytes.Replace(body, []byte(gooseDownPrefix), []byte(""), 1)
	body = bytes.Replace(body, []byte(gooseStatementBegin), []byte(""), 1)
	body = bytes.Replace(body, []byte(gooseStatementEnd), []byte(""), 1)
	return body
}
