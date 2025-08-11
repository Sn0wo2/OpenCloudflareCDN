package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
)

type JSONLoader struct{}

func NewJSONLoader() *JSONLoader {
	return &JSONLoader{}
}

func (j *JSONLoader) Load(cfg *config.Config, fileName string) error {
	file, err := os.ReadFile(fileName) //nolint:gosec
	if err != nil {
		return err
	}

	return json.Unmarshal(file, cfg)
}

func (j *JSONLoader) Save(cfg *config.Config, fileName string) error {
	file, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fileName), 0o750); err != nil {
		return err
	}

	return os.WriteFile(fileName, file, 0o600)
}

func (j *JSONLoader) GetAllowFileExtensions() []string {
	return []string{"json"}
}
