package css_parser

import (
	"errors"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type Parser struct {
	inline     bool
	parser     *css.Parser
	embedLevel int
}

func NewParser(content string, inline bool) *Parser {
	return &Parser{
		inline:     inline,
		parser:     css.NewParser(parse.NewInputString(content), inline),
		embedLevel: 0,
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
		return nil, err, false

	case css.CommentGrammar:
		return nil, nil, true // Skip comments

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
		return nil, errors.New("Unexpected grammar type: " + gt.String()), false
	}
}

func (p *Parser) ParseDeclarations() ([]*CssDeclaration, error) {
	ds := []*CssDeclaration{}

	for {
		gt, tt, data := p.parser.Next()

		switch gt {
		case css.ErrorGrammar:
			err := p.parser.Err()
			if p.inline && err.Error() == "EOF" {
				return ds, nil
			}
			return nil, err // Return error if not inline EOF

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			d, err := p.parseDeclarationBase(gt, tt, data, p.parser.Values())
			if err != nil {
				return nil, err
			}
			ds = append(ds, d)

		case css.CommentGrammar:
			// Skip comments
			continue

		default:
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
			return err

		case css.CommentGrammar:
			continue // Skip comments

		case css.EndAtRuleGrammar, css.EndRulesetGrammar:
			break ScanLoop // End of current rule

		default:
			err := p.parseSubRuleOrDeclarationBase(rule, gt, tt, data, values)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
