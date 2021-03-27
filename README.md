# go-yapla

[![Go Report Card](https://goreportcard.com/badge/git.iglou.eu/Production/go-yapla)](https://goreportcard.com/report/git.iglou.eu/Production/go-yapla)
[![Go Reference](https://img.shields.io/badge/api-reference-blue)](https://pkg.go.dev/git.iglou.eu/Production/go-yapla)
[![Go coverage](https://img.shields.io/badge/coverage-87.7%-green)](https://img.shields.io)
[![Apache V2 License](https://img.shields.io/badge/licence-MIT-9cf)](https://opensource.org/licenses/MIT)

Support Yapla **API 2.0**    
Official documentation available [**here**](https://app.swaggerhub.com/apis/yapla/yapla/2.0.0)

## Getting Started

### Dependencies
No extra dependency are needed

### Go mod
To get a specific release version use `@<tag>` in your go get command.
```sh
go get git.iglou.eu/Production/go-yapla@v0.0.1
```

Or `@latest` to get the latest repository change
```sh
go get git.iglou.eu/Production/go-yapla@latest
```

## Quick Example

This example shows a complete working Go file which will login a member and showing their information map
```go
package main

import (
	"fmt"
	"log"

	"git.iglou.eu/Production/go-yapla"
)

func main() {
    // Create a new required instance
    //
    // Need to provide a valid Yapla API Key 
    //
    // You can overwrite default configuration
    // For this provide an extra argument like
    // yapla.Config{...}
	y, err := yapla.NewSession("HP1ST252NFKX6Z6RVJ4RKEU23WS2QXSTQHTVYA1JAFWYX306")
	if err != nil {
		log.Fatalln(err)
	}

    // Login a Yapla member of your organisation
    // Need "only" user email and user password
    //
    // Return a struct with the API reply state and data on map format
    //
    // LoginContact() function is also avilable
	yRep, err := y.LoginMember("moncompte@macompagnie.com", "monp4ssW0R4!")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(yRep.Data)
}
```
