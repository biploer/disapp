package config

import (
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	AssetsDir   string `yaml:"assets_dir"`
	HTTPServer  `yaml:"http_server"`
	Certificate `yaml:"certificate"`
}

type HTTPServer struct {
	Scheme      string        `yaml:"scheme" env-default:"http"`
	Domain      string        `yaml:"domain" env-default:"localhost:8080"`
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Certificate struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// Load configuration file or call Fatalf
func MustLoad(configFs fs.FS, isProdEnv bool) Config {
	var config Config

	fileName := "local.yaml"
	if isProdEnv {
		fileName = "prod.yaml"
	}

	file, err := configFs.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := cleanenv.ParseYAML(file, &config); err != nil {
		log.Fatalf("can not parse config: %s", err)
	}
	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("can not read config depending on tags: %s", err)
	}
	fmt.Println(config)
	return config
}
