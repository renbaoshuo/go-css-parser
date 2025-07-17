package css_parser

import (
	"errors"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type Parser struct {
	Parser *css.Parser
}

func NewParser(content string, inline bool) *Parser {
	return &Parser{
		Parser: css.NewParser(parse.NewInputString(content), inline),
	}
}

func ParseStylesheet(content string) (*CssStylesheet, error) {
	return NewParser(content, false).ParseStylesheet()
}

func ParseDeclarations(content string) ([]*CssDeclaration, error) {
	return NewParser(content, true).ParseDeclarations()
}

func (p *Parser) ParseStylesheet() (*CssStylesheet, error) {
	s := NewCssStylesheet()

	rules, err := p.ParseRules()
	if err != nil {
		return nil, err
	}
	s.Rules = rules

	if len(s.Rules) == 0 {
		return nil, nil
	}

	return s, nil
}

func (p *Parser) ParseRules() ([]*CssRule, error) {
	return nil, nil // TODO: Implement this method

}

func (p *Parser) ParseDeclarations() ([]*CssDeclaration, error) {
	ds := []*CssDeclaration{}

	for {
		gt, _, data := p.Parser.Next()

		if gt == css.ErrorGrammar {
			err := p.Parser.Err()

			if err.Error() == "EOF" {
				return ds, nil // No rules found
			}

			return nil, err // Return error if not EOF
		}

		if gt == css.DeclarationGrammar || gt == css.CustomPropertyGrammar {
			d := NewCssDeclaration()
			d.Property = string(data)

			importantFlag := false

			for _, val := range p.Parser.Values() {
				val := string(val.Data)

				if importantFlag {
					importantFlag = false

					if val == "important" {
						d.Important = true
						continue
					}

					d.Value += "!" // Append the exclamation mark
				}

				if val == "!" {
					importantFlag = true
					continue
				}

				d.Value += val
			}

			ds = append(ds, d)

			continue
		}

		if gt == css.CommentGrammar {
			continue // Skip comments
		}

		return nil, errors.New("Unexpected grammar type: " + gt.String())
	}
}
