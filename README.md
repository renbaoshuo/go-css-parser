# Go CSS Parser

## Installation

```bash
go get go.baoshuo.dev/cssparser
```

## Usage

```go
package main

import "go.baoshuo.dev/cssparser"

func main() {
  // Parse CSS declarations (e.g., from style attribute)
  declarations, err := cssparser.ParseDeclarations(`
    color: red;
    font-size: 16px;
  `)

  // Parse CSS stylesheet (e.g., from a <style> tag)
  stylesheet, err := cssparser.ParseStylesheet(`
    .example {
      color: blue;
      font-size: 14px;
    }
  `)
}
```

## Credits

- https://github.com/tdewolff/parse
- https://github.com/aymerick/douceur
- https://github.com/csstree/csstree

## Author

**go-css-parser** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
