package file

import (
	"os"
	"path/filepath"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"gopkg.in/yaml.v3"
)

type YAMLLoader struct{}

func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

func (y *YAMLLoader) Load(cfg *config.Config, fileName string) error {
	file, err := os.ReadFile(fileName) //nolint:gosec
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, cfg)
}

func (y *YAMLLoader) Save(cfg *config.Config, fileName string) error {
	file, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fileName), 0o750); err != nil {
		return err
	}

	return os.WriteFile(fileName, file, 0o600)
}

func (y *YAMLLoader) GetAllowFileExtensions() []string {
	return []string{"yaml", "yml"}
}
