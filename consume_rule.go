package cssparser

import (
	"errors"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/rule"
	"go.baoshuo.dev/cssparser/selector"
	"go.baoshuo.dev/cssparser/token_stream"
)

// ConsumeRuleList consumes a list of CSS rules from the lexer.
//
// It handles both at-rules (like @media) and qualified rules.
// The function continues until it reaches the end of the input (EOF).
// It ignores whitespace and comments, and processes each rule accordingly.
//
// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-list-of-rules
func (p *Parser) consumeRuleList(
	allowedRules allowedRuleType,
	allowCdoCdcTokens bool,
	nestingType nesting.NestingTypeType,
	parentRuleForNesting *rule.StyleRule,
) ([]*rule.StyleRule, error) {
	var rules []*rule.StyleRule

loop:
	for {
		token := p.s.Peek()

		switch token.Type {
		case csslexer.EOFToken:
			break loop

		case csslexer.WhitespaceToken:
			// Ignore whitespace
			p.s.Consume()
			continue

		case csslexer.CDCToken, csslexer.CDOToken:
			// TODO: Handle CDCToken and CDOToken if needed, now we just ignore them
			p.s.Consume()
			continue

		case csslexer.AtKeywordToken:
			// Handle at-rules like @media
			rule, err := p.consumeAtRule()
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)

		default:
			// Handle qualified rules
			rule, err := p.consumeQualifiedRule(allowedRules, nestingType, parentRuleForNesting)
			if err != nil {
				return nil, err
			}
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

// consumeAtRule consumes an at-rule from the lexer.
//
// https://www.w3.org/TR/2021/CRD-css-syntax-3-20211224/#consume-at-rule
func (p *Parser) consumeAtRule() (*rule.StyleRule, error) {
	return nil, errors.New("not implemented yet")
}

// consumeQualifiedRule consumes a qualified rule from the lexer.
//
// https://drafts.csswg.org/css-syntax/#consume-qualified-rule
func (p *Parser) consumeQualifiedRule(
	allowedRules allowedRuleType,
	nestingType nesting.NestingTypeType,
	parentRuleForNesting *rule.StyleRule,
) (*rule.StyleRule, error) {
	if allowedRules.Has(qualifiedRuleTypeStyle) {
		return p.consumeStyleRule(nestingType, parentRuleForNesting, false)
	}

	if allowedRules.Has(qualifiedRuleTypeKeyframes) {
		return p.consumeKeyframeStyleRule()
	}

	return nil, errors.New("no qualified rule parsed")
}

func (p *Parser) consumeStyleRule(
	nestingType nesting.NestingTypeType,
	parentRuleForNesting *rule.StyleRule,
	nested bool,
) (*rule.StyleRule, error) {
	selectors, err := selector.ConsumeSelector(p.s, nestingType, parentRuleForNesting)

	if err != nil || len(selectors) == 0 {
		// Read the rest of the prelude if there was an error
		if nested {
			p.s.SkipUntil(csslexer.LeftBraceToken, csslexer.SemicolonToken)
		} else {
			p.s.SkipUntil(csslexer.LeftBraceToken)
		}
	}

	if p.s.Peek().Type != csslexer.LeftBraceToken {
		return nil, errors.New("expected '{' after selector")
	}

	if len(selectors) == 0 {
		err := p.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
			return nil
		})
		if err != nil {
			return nil, err
		}

		return nil, errors.New("invalid selector")
	}

	rule := &rule.StyleRule{
		Type: rule.StyleRuleTypeQualifiedRule,
	}

	err = p.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
		rules, err := p.consumeStyleRuleContents()
		if err != nil {
			return err
		}
		rule.Rules = rules
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func (p *Parser) consumeStyleRuleContents() ([]*rule.Rule, error) {
	return nil, errors.New("not implemented yet")
}

func (p *Parser) consumeKeyframeStyleRule() (*rule.StyleRule, error) {
	return nil, errors.New("not implemented yet")
}
