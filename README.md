# echo
A lite log package, provide [zap](https://github.com/uber-go/zap) style API.

## Install
```
go get -u -v github.com/vizee/echo
```

### Usage
example:
```go
package main

import (
	"errors"
	"os"

	"github.com/vizee/echo"
)

func main() {
	echo.SetOutput(os.Stdout)
	echo.SetLevel(echo.DebugLevel)
	echo.SetFormatter(&echo.PlainFormatter{})
	echo.Debug("debug message1", echo.Int("int", 126), echo.Echo("bytes", echo.BytesEchoer(`233`)))
	echo.SetLevel(echo.InfoLevel)
	echo.Debug("debug message2", echo.Int("int", 126))
	echo.Info("info message", echo.Var("var", map[string]int{"a": 1, "b": 2}))
	echo.Warn("warn message", echo.String("string", "blah\n\tblah"), echo.Quote("quote", "blah\n\tblah"))
	echo.Error("error message", echo.Errors("err", errors.New("oops!")))
	echo.Fatal("fatal message", echo.Stack(true))
}
```
example output:
```
2017-06-17/00:17:04 [D] debug message1 {int=126 bytes=233}
2017-06-17/00:17:04 [I] info message {var=map[a:1 b:2]}
2017-06-17/00:17:04 [W] warn message {string=blah
        blah quote="blah\n\tblah"}
2017-06-17/00:17:04 [E] error message {err=oops!}
2017-06-17/00:17:04 [F] fatal message {
goroutine 1 [running]:
github.com/vizee/echo.Stack(0x4caa01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
        /home/vizee/src/go/src/github.com/vizee/echo/field.go:128 +0xcc
main.main()
        /home/vizee/src/go/src/playground/main.go:20 +0x9d5
}
```