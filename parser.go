package cssparser

import (
	"errors"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type Parser struct {
	parser     *css.Parser
	embedLevel int

	inline bool
	loose  bool // Whether to allow loose parsing, which is more permissive and allows for some errors in the CSS syntax.
}

func NewParser(content string, options ...ParserOption) *Parser {
	parser := &Parser{
		embedLevel: 0,
	}

	for _, option := range options {
		option(parser)
	}

	parser.parser = css.NewParser(parse.NewInputString(content), parser.inline)

	return parser
}

func ParseStylesheet(content string, options ...ParserOption) (*Stylesheet, error) {
	return NewParser(content, options...).ParseStylesheet()
}

func ParseDeclarations(content string, options ...ParserOption) ([]*Declaration, error) {
	fullOptions := []ParserOption{WithInline(true)}
	fullOptions = append(fullOptions, options...)
	return NewParser(content, fullOptions...).ParseDeclarations()
}

func (p *Parser) ParseStylesheet() (*Stylesheet, error) {
	s := NewStylesheet()

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
	rules := []*CssRule{}

	for {
		rule, err, eof := p.ParseRule()

		if eof {
			break
		}

		if err != nil {
			return nil, err
		}

		if rule != nil {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func (p *Parser) ParseRule() (*CssRule, error, bool) {
	gt, tt, data := p.parser.Next()

	switch gt {
	case css.ErrorGrammar:
		err := p.parser.Err()
		if err.Error() == "EOF" {
			return nil, nil, true
		}
		// In loose mode, skip error tokens and continue parsing
		if p.loose {
			return nil, nil, false
		}
		return nil, err, false

	case css.CommentGrammar, css.TokenGrammar:
		return nil, nil, false // Skip comments

	case css.AtRuleGrammar:
		rule, err := p.parseSingleAtRuleBase(gt, tt, data, p.parser.Values())
		if err != nil {
			return nil, err, false
		}
		return rule, nil, false

	case css.BeginAtRuleGrammar:
		rule, err := p.parseAtRuleBase(gt, tt, data, p.parser.Values())
		if err != nil {
			return nil, err, false
		}
		return rule, nil, false

	case css.BeginRulesetGrammar, css.QualifiedRuleGrammar:
		rule, err := p.parseRuleBase(gt, tt, data, p.parser.Values())
		if err != nil {
			return nil, err, false
		}
		return rule, nil, false

	case css.EndAtRuleGrammar, css.EndRulesetGrammar:
		return nil, nil, false

	default:
		// In loose mode, skip unexpected grammar types and continue parsing
		if p.loose {
			return nil, nil, false
		}
		return nil, errors.New("Unexpected grammar type: " + gt.String()), false
	}
}

// ParseDeclarations parses a list of CSS declarations from the input, now it uses only in inline mode.
func (p *Parser) ParseDeclarations() ([]*Declaration, error) {
	ds := []*Declaration{}

	for {
		gt, tt, data := p.parser.Next()

		switch gt {
		case css.ErrorGrammar:
			err := p.parser.Err()
			if p.inline && err.Error() == "EOF" {
				return ds, nil
			}
			// In loose mode, skip error tokens and continue parsing
			if p.loose {
				continue
			}
			return nil, err // Return error if not inline EOF

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			d, err := p.parseDeclarationBase(gt, tt, data, p.parser.Values())
			if err != nil {
				return nil, err
			}
			ds = append(ds, d)

		case css.CommentGrammar, css.TokenGrammar:
			// Skip comments
			continue

		default:
			// In loose mode, skip unexpected grammar types and continue parsing
			if p.loose {
				continue
			}
			return nil, errors.New("Unexpected grammar type: " + gt.String())
		}
	}
}

func (p *Parser) parseSubRuleOrDeclarations(rule *CssRule) error {
ScanLoop:
	for {
		gt, tt, data := p.parser.Next()
		values := p.parser.Values()

		switch gt {
		case css.ErrorGrammar:
			err := p.parser.Err()
			if err.Error() == "EOF" {
				return nil // EOF is not an error in this context
			}
			// In loose mode, skip error tokens and continue parsing
			if p.loose {
				continue
			}
			return err

		case css.CommentGrammar, css.TokenGrammar:
			continue // Skip comments

		case css.EndAtRuleGrammar, css.EndRulesetGrammar:
			break ScanLoop // End of current rule

		default:
			err := p.parseSubRuleOrDeclarationBase(rule, gt, tt, data, values)
			if err != nil {
				// In loose mode, skip errors and continue parsing
				if p.loose {
					continue
				}
				return err
			}
		}
	}

	return nil
}
