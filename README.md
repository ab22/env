# env

Use tag structures to parse environment variables into structure fields.

## Usage

```go
package main

import (
	"fmt"
	"github.com/ab22/env"
)

type AppConfig struct {
	Env        string `env:"ENV" envDefault:"PRODUCTION"`
	RiotApiKey string `env:"RIOT_API_KEY"`
	Port       int    `env:"APP_PORT" envDefault:"1337"`

	Smtp struct {
		Host     string `env:"SMTP_HOST"`
		Port     int    `env:"SMTP_PORT"`
		User     string `env:"SMTP_USER"`
		Password string `env:"SMTP_PASS"`
	}

	Db struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASS" envDefault:"1234"`
		Name     string `env:"DB_NAME" envDefault:"lol_db"`
	}
}

func main() {
	config := &AppConfig{}

	env.Parse(config)
	fmt.Println(config)
}
```
