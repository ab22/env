# env

## Example

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
}

func main() {
	config := &AppConfig{}

	env.Parse(config)
	fmt.Println(config)
}
```
