package cssparser

import (
	"errors"

	"go.baoshuo.dev/csslexer"
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
	nestingType NestingTypeType,
	parentRuleForNesting *Rule,
) ([]*Rule, error) {
	var rules []*Rule

loop:
	for {
		token := p.s.Peek()

		switch token.Type {
		case csslexer.EOFToken:
			break loop

		case csslexer.WhitespaceToken, csslexer.CommentToken:
			// Ignore whitespace and comments
			continue

		case csslexer.CDCToken, csslexer.CDOToken:
			// TODO: Handle CDCToken and CDOToken if needed, now we just ignore them
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

func (p *Parser) consumeAtRule() (*Rule, error) {
	return nil, errors.New("not implemented yet")
}

func (p *Parser) consumeQualifiedRule(
	allowedRules allowedRuleType,
	nestingType NestingTypeType,
	parentRuleForNesting *Rule,
) (*Rule, error) {
	if allowedRules.Has(qualifiedRuleTypeStyle) {
		return p.consumeStyleRule(nestingType, parentRuleForNesting, false)
	}

	if allowedRules.Has(qualifiedRuleTypeKeyframes) {
		return p.consumeKeyframeStyleRule()
	}

	return nil, errors.New("no qualified rule parsed")
}

func (p *Parser) consumeStyleRule(
	nestingType NestingTypeType,
	parentRuleForNesting *Rule,
	nested bool,
) (*Rule, error) {
	return nil, errors.New("not implemented yet")
}

func (p *Parser) consumeKeyframeStyleRule() (*Rule, error) {
	return nil, errors.New("not implemented yet")
}
