package css_parser

import (
	"github.com/tdewolff/parse/v2/css"
)

func valuesToString(values []css.Token) string {
	result := ""
	for _, val := range values {
		result += string(val.Data)
	}
	return result
}
