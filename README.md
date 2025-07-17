# go-css-parser

## Installation

```bash
go get github.com/renbaoshuo/go-css-parser
```

## Usage

```go
package main

import "github.com/renbaoshuo/go-css-parser"

func main() {
  // Parse CSS declarations (e.g., from style attribute)
  parser.ParseDeclarations(`
    color: red;
    font-size: 16px;
  `)
}
```

## Credits

- https://github.com/tdewolff/parse
- https://github.com/csstree/csstree
- https://github.com/aymerick/douceur

## Author

**go-css-parser** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
