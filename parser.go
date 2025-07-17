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
		rule, err := p.ParseRule()

		if err != nil {
			return nil, err
		}

		if rule == nil {
			break // EOF reached, exit the loop
		}

		rules = append(rules, rule)
	}

	return rules, nil // TODO: Implement this method
}

func (p *Parser) ParseRule() (*CssRule, error) {
	firstRun := true
	selectors := []string{}

	for {
		gt, _, _ := p.parser.Next()

		switch gt {
		case css.ErrorGrammar:
			err := p.parser.Err()
			if firstRun && err.Error() == "EOF" {
				return nil, nil
			}
			return nil, err

		case css.CommentGrammar:
			continue

		case css.QualifiedRuleGrammar:
			selector := valuesToString(p.parser.Values())
			selectors = append(selectors, selector)

		case css.AtRuleGrammar:
			// TODO
			return nil, errors.New("AtRuleGrammar parsing not implemented yet")

		case css.BeginAtRuleGrammar:
			// TODO
			return nil, errors.New("BeginAtRuleGrammar parsing not implemented yet")

		case css.BeginRulesetGrammar:
			selector := valuesToString(p.parser.Values())
			selectors = append(selectors, selector)

			rule := NewCssRule(QualifiedRule)
			rule.EmbedLevel = p.embedLevel
			rule.Selectors = selectors

			p.embedLevel++
			declarations, err := p.ParseDeclarations()
			p.embedLevel--
			if err != nil {
				return nil, err
			}

			rule.Declarations = declarations

			return rule, nil

		default:
			return nil, errors.New("Unexpected grammar type: " + gt.String())
		}

		firstRun = false
	}
}

func (p *Parser) ParseDeclarations() ([]*CssDeclaration, error) {
	ds := []*CssDeclaration{}

	for {
		gt, _, data := p.parser.Next()

		switch gt {
		case css.ErrorGrammar:
			err := p.parser.Err()
			if p.inline && err.Error() == "EOF" {
				break // If inline and EOF, break the loop
			}
			return nil, err // Return error if not EOF

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			d := NewCssDeclaration()
			d.Property = string(data)

			importantFlag := false

			for _, val := range p.parser.Values() {
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

		case css.CommentGrammar:
			// Skip comments
			continue

		case css.EndRulesetGrammar, css.EndAtRuleGrammar:
			return ds, nil // Return declarations when end of ruleset or at-rule is reached

		default:
			return nil, errors.New("Unexpected grammar type: " + gt.String())
		}
	}
}

func valuesToString(values []css.Token) string {
	result := ""
	for _, val := range values {
		result += string(val.Data)
	}
	return result
}
