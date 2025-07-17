package css_parser

import (
	"errors"

	"github.com/tdewolff/parse/v2/css"
)

func (p *Parser) parseSubRuleOrDeclarationBase(rule *CssRule, gt css.GrammarType, tt css.TokenType, data []byte, values []css.Token) error {
	switch gt {
	case css.DeclarationGrammar, css.CustomPropertyGrammar:
		d, err := p.parseDeclarationBase(gt, tt, data, values)
		if err != nil {
			return err
		}
		rule.Declarations = append(rule.Declarations, d)

	case css.BeginRulesetGrammar, css.QualifiedRuleGrammar:
		r, err := p.parseRuleBase(gt, tt, data, values)
		if err != nil {
			return err
		}
		rule.Rules = append(rule.Rules, r)

	case css.AtRuleGrammar:
		r, err := p.parseSingleAtRuleBase(gt, tt, data, values)
		if err != nil {
			return err
		}
		rule.Rules = append(rule.Rules, r)

	case css.BeginAtRuleGrammar:
		r, err := p.parseAtRuleBase(gt, tt, data, values)
		if err != nil {
			return err
		}
		rule.Rules = append(rule.Rules, r)

	default:
		return errors.New("Unexpected grammar type: " + gt.String())
	}

	return nil
}

func (p *Parser) parseSingleAtRuleBase(_ css.GrammarType, _ css.TokenType, data []byte, values []css.Token) (*CssRule, error) {
	rule := NewCssRule(AtRule)
	rule.EmbedLevel = p.embedLevel
	rule.Name = string(data)
	rule.Prelude = valuesToString(values)

	return rule, nil
}

func (p *Parser) parseAtRuleBase(_ css.GrammarType, _ css.TokenType, data []byte, values []css.Token) (*CssRule, error) {
	rule := NewCssRule(AtRule)
	rule.EmbedLevel = p.embedLevel
	rule.Name = string(data)
	rule.Prelude = valuesToString(values)

	p.embedLevel++
	err := p.parseSubRuleOrDeclarations(rule)
	p.embedLevel--

	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (p *Parser) parseRuleBase(gt css.GrammarType, _ css.TokenType, _ []byte, values []css.Token) (*CssRule, error) {
	var selectors []string

GetAllSelectorsLoop:
	for {
		switch gt {
		case css.QualifiedRuleGrammar:
			selector := valuesToString(values)
			selectors = append(selectors, selector)
			gt, _, _ = p.parser.Next()
			values = p.parser.Values()

		case css.BeginRulesetGrammar:
			selector := valuesToString(values)
			selectors = append(selectors, selector)
			break GetAllSelectorsLoop

		default:
			return nil, errors.New("Unexpected grammar type: " + gt.String())
		}
	}

	rule := NewCssRule(QualifiedRule)
	rule.EmbedLevel = p.embedLevel
	rule.Selectors = selectors

	p.embedLevel++
	err := p.parseSubRuleOrDeclarations(rule)
	p.embedLevel--

	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (*Parser) parseDeclarationBase(gt css.GrammarType, tt css.TokenType, data []byte, values []css.Token) (*CssDeclaration, error) {
	d := NewCssDeclaration()
	d.Property = string(data)

	importantFlag := false

	for _, val := range values {
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

	return d, nil
}
