envigo [![GitHub release](https://img.shields.io/github/release/tyranron/envigo.svg)](https://github.com/tyranron/envigo/releases)
======

[![Build Status](https://travis-ci.org/tyranron/envigo.svg?branch=master)](https://travis-ci.org/tyranron/envigo)
[![GoCover](https://gocover.io/_badge/github.com/tyranron/envigo)](https://gocover.io/github.com/tyranron/envigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/tyranron/envigo)](https://goreportcard.com/report/github.com/tyranron/envigo)
[![GoDoc](https://godoc.org/github.com/tyranron/envigo?status.svg)](https://godoc.org/github.com/tyranron/envigo)
[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue.svg)](#license)

Yet another Go package for parsing environment variables in `struct`s.




## Usage

Having some environment variables:
```bash
export DEBUG_MODE=true
export TIMEOUT="4s300ms"
```

And struct with tagged fields:
```go
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
```

Results in:
```
&{DebugMode:true WorkersCount:0 Timeouts:{Default:4.3s}}
```




## Supported Types

These types are supported by envigo be default:

- `bool`
- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `byte`, `rune`,
- `float32`, `float64`
- [`time.Duration`][1]
- [`net.IP`][3]
- anything that implements [`encoding.TextUnmarshaler`][2]




## Custom Parsing

To make parser be able to parse a type that is not supported by default, or to change default parsing behaviour,you just need to implement [`encoding.TextUnmarshaler`][2] interface for your type.
 
```go
type MyType struct {
	Value int
}

func (t *MyType) UnmarshalText(envVarValue []byte) error {
	val, err := strconv.ParseInt(string(envVarValue), 10, 64)
	if err != nil {
		return err
	}
	t.Value = val
	return nil
}
```




## TODO

- parsing slices, maps
- parsing `time.Time` (?)
- different parsing modes (strict, etc)




## License

Licensed under either of

- Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT License ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.




[1]: https://golang.org/pkg/time/#Duration
[2]: https://golang.org/pkg/encoding/#TextUnmarshaler
[3]: https://golang.org/pkg/net/#IP
