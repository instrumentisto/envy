package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tyranron/envigo"
)

type Config struct {
	DebugMode    bool `env:"DEBUG_MODE"`
	WorkersCount int  `env:"WORKERS_COUNT"`
	Timeouts     struct {
		Default time.Duration `env:"TIMEOUT"`
	}
}

func main() {
	conf := &Config{}
	if err := (envigo.Parser{}).Parse(conf); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", conf)
}
