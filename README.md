# go-ssa2ast

An experimental implementation of converting from [go/ssa](https://pkg.go.dev/golang.org/x/tools/go/ssa) back to [go/ast](https://pkg.go.dev/go/ast)

## Known problems:

- Anonymous functions/lambdas are **not** supported
- Convert breaks "lazy" `for range` for map to go because there is no publicly available iterator for map