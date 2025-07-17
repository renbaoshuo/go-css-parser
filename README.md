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
  decls, err := parser.ParseDeclarations(`
    color: red;
    font-size: 16px;
  `)

  // Parse CSS stylesheet (e.g., from a <style> tag)
  stylesheet, err := parser.ParseStylesheet(`
    .example {
      color: blue;
      font-size: 14px;
    }
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
