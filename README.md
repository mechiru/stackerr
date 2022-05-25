# stackerr

[![ci](https://github.com/mechiru/stackerr/workflows/ci/badge.svg)](https://github.com/mechiru/stackerr/actions?query=workflow:ci)

This library provides an error with a stack trace.

## Example

Simple error:
```go
err := stackerr.New("error")
fmt.Println("------ error without stack ------")
fmt.Println(err.Error())
fmt.Println("------ error with stack ------")
fmt.Printf("%+v\n", err)
```

```console
------ error without stack ------
error
------ error with stack ------
error
goroutine 1 [running]:
main.simple()
        /Users/mechiru/src/github.com/mechiru/stackerr/_examples/console/main.go:16 +0x2c
main.main()
        /Users/mechiru/src/github.com/mechiru/stackerr/_examples/console/main.go:11 +0x20
```

Wrap error:
```go
err := errors.New("error1")
err = stackerr.Errorf("error2: %w", err)
fmt.Println("------ error without stack ------")
fmt.Println(err.Error())
fmt.Println("------ error with stack ------")
fmt.Printf("%+v\n", err)

fmt.Printf("errors.Is: %v\n", errors.Is(e2, e1))
```

```console
------ error without stack ------
error2: error1
------ error with stack ------
error2: error1
goroutine 1 [running]:
main.wrap()
        /Users/mechiru/src/github.com/mechiru/stackerr/_examples/console/main.go:26 +0x70
main.main()
        /Users/mechiru/src/github.com/mechiru/stackerr/_examples/console/main.go:12 +0x24
errors.Is: true
```

## Lisence
Licensed under the [MIT license](./LICENSE).
