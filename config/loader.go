package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sn0wo2/OpenCloudflareCDN/debug"
)

var Instance *Config

var DefaultConfig = &Config{
	Log: Log{
		Level: "debug",
		Dir:   "./logs",
	},
	Server: Server{
		Address: ":3000",
		Header:  "OpenCloudflareCDN",
	},
	StaticPath:     "./static",
	StaticIndex:    "index.html",
	OriginalServer: "https://www.example.com",
	JWTSecret:      "your-super-secret-and-long-key",
}

type Config struct {
	IsDefault  bool   `json:"-" optional:"true" yaml:"-"`
	ConfigPath string `json:"-" optional:"true" yaml:"-"`

	Log                Log    `json:"log"                optional:"true"           yaml:"log"`
	Server             Server `json:"server"             yaml:"server"`
	StaticPath         string `json:"staticPath"         yaml:"staticPath"`
	StaticIndex        string `json:"staticIndex"        optional:"true"           yaml:"staticIndex"`
	OriginalServer     string `json:"OriginalServer"     yaml:"OriginalServer"`
	TurnstileSecretKey string `json:"turnstileSecretKey" yaml:"turnstileSecretKey"`
	JWTSecret          string `json:"jwtSecret"          yaml:"jwtSecret"`
}

type Log struct {
	Level string `json:"level" optional:"true" yaml:"level"`
	Dir   string `json:"dir"   optional:"true" yaml:"dir"`
}

type Server struct {
	Address string `json:"address" yaml:"address"`
	Header  string `json:"header"  optional:"true" yaml:"header"`
	TLS     TLS    `json:"tls"     optional:"true" yaml:"tls"`
}

type TLS struct {
	Cert string `json:"cert" yaml:"cert"`
	Key  string `json:"key"  yaml:"key"`
}

type Loader interface {
	Load(cfg *Config, fileName string) error
	Save(cfg *Config, fileName string) error
	// GetAllowFileExtensions lowercase
	GetAllowFileExtensions() []string
}

var ErrConfigNotFound = errors.New("config file not found")

func Init(loaders ...Loader) error {
	var err error

	Instance, err = NewConfig(loaders...)

	return err
}

func NewConfig(loaders ...Loader) (*Config, error) {
	if len(loaders) == 0 {
		return nil, errors.New("no loaders provided")
	}

	loaderByExt := make(map[string]Loader)

	for _, l := range loaders {
		for _, ext := range l.GetAllowFileExtensions() {
			loaderByExt["."+strings.ToLower(ext)] = l
		}
	}

	envPath := os.Getenv("CONFIG_PATH")
	if debug.IsDebugging() {
		if p := os.Getenv("DEBUG_CONFIG_PATH"); p != "" {
			envPath = p
		}
	}

	var foundPath string

	if envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			foundPath = envPath
		} else {
			base := strings.TrimSuffix(envPath, filepath.Ext(envPath))
			for ext := range loaderByExt {
				tryPath := base + ext
				if _, err := os.Stat(tryPath); err == nil {
					foundPath = tryPath

					break
				}
			}
		}
	}

	if foundPath == "" {
		searchPaths := []string{"./data/"}

	searchLoop:
		for _, p := range searchPaths {
			for ext := range loaderByExt {
				fullPath := filepath.Join(p, "config"+ext)
				if _, err := os.Stat(fullPath); err == nil {
					foundPath = fullPath

					break searchLoop
				}
			}
		}
	}

	if foundPath == "" {
		DefaultConfig.ConfigPath = envPath

		return DefaultConfig, ErrConfigNotFound
	}

	ext := strings.ToLower(filepath.Ext(foundPath))

	loader, ok := loaderByExt[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported config file extension: %s", ext)
	}

	var fileCfg Config
	if err := loader.Load(&fileCfg, foundPath); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", foundPath, err)
	}

	if err := validate(&fileCfg); err != nil {
		return nil, fmt.Errorf("validation failed for config file %s: %w", foundPath, err)
	}

	merge(DefaultConfig, &fileCfg)

	DefaultConfig.ConfigPath = foundPath

	return DefaultConfig, nil
}
