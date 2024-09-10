package env

import (
	"fmt"
	"io/fs"
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
		DbUrl:        fmt.Sprintf("file:%s", escaped),
	}
}

type ConfigStore struct {
	ConfigPath string
}

func NewConfigStore() (*ConfigStore, error) {
	configFilePath, err := xdg.ConfigFile("jamkstrecka/config.yml")
	if err != nil {
		return nil, fmt.Errorf("could not resolve path for config file: %w", err)
	}

	return &ConfigStore{
		ConfigPath: configFilePath,
	}, nil
}

func (s *ConfigStore) Config() (Config, error) {
	_, err := os.Stat(s.ConfigPath)
	if os.IsNotExist(err) {
		cfg, _ := yaml.Marshal(DefaultConfig())
		os.WriteFile(s.ConfigPath, cfg, 0644)
		fmt.Printf("No config file found, creating and using defaults at %s\n", s.ConfigPath)
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

	fmt.Printf("config: %v\n", cfg)
	return cfg, nil

}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func ReadFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {

		fmt.Println(os.Getwd())

		fmt.Println("No config file found, creating and using defaults")

		os.WriteFile("config.yml", []byte(`
				discord_token: ""\n
				guild: ""
			`), 0644)

		processError(err)

		return
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func ReadEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
