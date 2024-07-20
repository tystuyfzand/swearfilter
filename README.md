# swearfilter

[![GoDoc](https://godoc.org/github.com/JoshuaDoes/gofuckyourself?status.svg)](https://godoc.org/github.com/tystuyfzand/swearfilter)
[![Go Report Card](https://goreportcard.com/badge/github.com/tystuyfzand/swearfilter)](https://goreportcard.com/report/github.com/JoshuaDoes/gofuckyourself)

A sanitization-based swear filter for Go.

# Installing
`go get github.com/tystuyfzand/swearfilter`

# Example
```Go
package main

import (
	"fmt"

	"github.com/tystuyfzand/swearfilter"
)

var message = "This is a fooing message with barring swear words."
var swears = []string{"foo", "bar"}

func main() {
	filter := swearfilter.New(false, swears...)
	swearsFound, err := filter.Check(message)
	fmt.Println("Swears tripped: ", swearsFound)
	fmt.Println("Error: ", err)
}
```
### Output
```
> go run main.go
Swears tripped:  [foo bar]
Error:  <nil>
```

## License
The source code for swearfilter is released under the MIT License. See LICENSE for more details.

## Donations

*The below is the original creator's donation link. Please support them, this fork does nothing but clean it up!*

All donations are appreciated and help me stay awake at night to work on this more. Even if it's not much, it helps a lot in the long run!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)