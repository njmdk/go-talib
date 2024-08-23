# go-talib

本库由原始库(https://github.com/markcheno/go-talib)修改

主要增加一些不支持或者有误的指标

A pure [Go](http://golang.org/) port of [TA-Lib](http://ta-lib.org)

## Install

Install the package with:

```bash
go get github.com/njmdk/go-talib
```

Import it with:

```go
import "github.com/njmdk/go-talib"
```

and use `talib` as the package name inside the code.

## Example

```go
package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
	"github.com/njmdk/go-talib"
)

func main() {
	spy, _ := quote.NewQuoteFromYahoo("spy", "2016-01-01", "2016-04-01", quote.Daily, true)
	fmt.Print(spy.CSV())
	rsi2 := talib.Rsi(spy.Close, 2)
	fmt.Println(rsi2)
}
```

## License

MIT License  - see LICENSE for more details

# Contributors

- [代码](https://github.com/njmdk) 
- [Alessandro Sanino AKA saniales](https://github.com/saniales)
