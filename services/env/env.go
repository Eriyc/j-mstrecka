package env

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DiscordToken string `yaml:"discord_token" envconfig:"DISCORD_TOKEN" required:"true"`
	Guild        string `yaml:"guild" envconfig:"GUILD" required:"false"`
	DbUrl        string `yaml:"db_url" envconfig:"DB_URL" required:"true"`
}

func DefaultConfig() Config {
	file, _ := xdg.DataFile("jamkstrecka/local.db")

	escaped := strings.ReplaceAll(file, " ", "\\ ")

	return Config{
		DiscordToken: "",
		Guild:        "",
		DbUrl:        escaped,
	}
}

type ConfigStore struct {
	ConfigPath string
	Logger     *slog.Logger
}

func NewConfigStore(logger *slog.Logger) (*ConfigStore, error) {
	configFilePath, err := xdg.ConfigFile("jamkstrecka/config.yml")
	if err != nil {
		return nil, fmt.Errorf("could not resolve path for config file: %w", err)
	}

	return &ConfigStore{
		ConfigPath: configFilePath,
		Logger:     logger,
	}, nil
}

func (s *ConfigStore) Config() (Config, error) {
	_, err := os.Stat(s.ConfigPath)
	if os.IsNotExist(err) {
		cfg, _ := yaml.Marshal(DefaultConfig())
		os.WriteFile(s.ConfigPath, cfg, 0644)
		s.Logger.Error("no config file found, creating and using defaults", "path", s.ConfigPath)
		return DefaultConfig(), nil
	}

	dir, fileName := filepath.Split(s.ConfigPath)
	if len(dir) == 0 {
		dir = "."
	}

	buf, err := fs.ReadFile(os.DirFS(dir), fileName)
	if err != nil {
		return Config{}, fmt.Errorf("could not read the configuration file: %w", err)
	}

	if len(buf) == 0 {
		return DefaultConfig(), nil
	}

	cfg := Config{}
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return Config{}, fmt.Errorf("configuration file does not have a valid format: %w", err)
	}

	s.Logger.Debug("config", "token", cfg)
	return cfg, nil

}

func ReadFile(cfg *Config) error {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println("No config file found, creating and using defaults")
		os.WriteFile("config.yml", []byte(`discord_token: ""\nguild: ""`), 0644)

		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)

	return err
}

func ReadEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	return err
}
