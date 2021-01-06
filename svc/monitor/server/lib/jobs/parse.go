package jobs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"dfl/lib/cher"
	"dfl/svc/monitor"
)

func ParseData(jsonString []byte) ([]*monitor.Job, error) {
	var out []*monitor.Job

	if err := json.Unmarshal(jsonString, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func ParseFromFile(filePath string) ([]*monitor.Job, error) {
	if !fileExists(filePath) {
		return nil, cher.New("jobs_file_missing", nil)
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParseData(contents)
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
