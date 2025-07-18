/*
Package `cssparser` provides a parser for CSS (Cascading Style Sheets) files.

It allows parsing of stylesheets and declarations, handling both inline and embedded styles.

For example, you can use it to parse a CSS stylesheet and retrieve its rules or declarations.

Here's a brief overview of the main functions:
- ParseStylesheet(content string) (*Stylesheet, error): Parses a complete CSS stylesheet.
- ParseDeclarations(content string) ([]*Declaration, error): Parses CSS declarations, typically used for inline styles.

Here are the available options for the parser:
- WithInline(bool): Whether to parse inline styles.
- WithLooseParsing(bool): Whether to allow loose parsing, which is more permissive and allows for some errors in the CSS syntax.

The source code of this package is hosted on GitHub: https://github.com/renbaoshuo/go-css-parser
*/
package cssparser
