package cssparser

import (
	"errors"
	"strings"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/selector"
	"go.baoshuo.dev/cssparser/token_stream"
	"go.baoshuo.dev/cssparser/variable"
)

// consumeRuleList consumes a list of CSS rules from the lexer.
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
	parentRuleForNesting *css.StyleRule,
) ([]*css.StyleRule, error) {
	var rules []*css.StyleRule

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
func (p *Parser) consumeAtRule() (*css.StyleRule, error) {
	return nil, errors.New("not implemented yet")
}

// consumeQualifiedRule consumes a qualified rule from the lexer.
//
// https://drafts.csswg.org/css-syntax/#consume-qualified-rule
func (p *Parser) consumeQualifiedRule(
	allowedRules allowedRuleType,
	nestingType nesting.NestingTypeType,
	parentRuleForNesting *css.StyleRule,
) (*css.StyleRule, error) {
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
	parentRuleForNesting *css.StyleRule,
	nested bool,
) (*css.StyleRule, error) {
	// Check for custom property ambiguity - style rules that look like custom property declarations
	// are not allowed by css-syntax
	// https://drafts.csswg.org/css-syntax/#consume-qualified-rule
	customPropertyAmbiguity := p.startsCustomPropertyDeclaration()

	// Parse the prelude of the style rule (selectors)
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
		// Parse error, EOF instead of qualified rule block
		// (or we went into error recovery above)
		// NOTE: If we aborted due to a semicolon, don't consume it here;
		// the caller will do that for us
		return nil, errors.New("expected '{' after selector")
	}

	if customPropertyAmbiguity {
		if nested {
			// https://drafts.csswg.org/css-syntax/#consume-the-remnants-of-a-bad-declaration
			// Note that the caller consumes the bad declaration remnants
			return nil, errors.New("custom property ambiguity in nested context")
		}
		// "If nested is false, consume a block from input, and return nothing."
		// https://drafts.csswg.org/css-syntax/#consume-qualified-rule
		err := p.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
			return nil
		})
		if err != nil {
			return nil, err
		}
		return nil, errors.New("custom property ambiguity")
	}

	// Check if rule is "valid in current context"
	// https://drafts.csswg.org/css-syntax/#consume-qualified-rule
	// This means checking if the selector parsed successfully
	if len(selectors) == 0 {
		err := p.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
			return nil
		})
		if err != nil {
			return nil, err
		}
		return nil, errors.New("invalid selector")
	}

	// Create the style rule and consume its contents
	styleRule := &css.StyleRule{
		Type:      css.StyleRuleTypeQualifiedRule,
		Selectors: selectors,
	}

	err = p.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
		return p.consumeStyleRuleContents(styleRule, nestingType)
	})
	if err != nil {
		return nil, err
	}

	return styleRule, nil
}

// startsCustomPropertyDeclaration checks if the current token stream starts with a custom property declaration
func (p *Parser) startsCustomPropertyDeclaration() bool {
	return variable.StartsCustomPropertyDeclaration(*p.s)
}

// consumeStyleRuleContents consumes the contents of a style rule block
func (p *Parser) consumeStyleRuleContents(styleRule *css.StyleRule, nestingType nesting.NestingTypeType) error {
	var childRules []*css.StyleRule
	var declarations []*css.Declaration

	for {
		// Skip whitespace and comments
		p.s.ConsumeWhitespace()

		if p.s.AtEnd() {
			break
		}

		token := p.s.Peek()

		switch token.Type {
		case csslexer.SemicolonToken:
			// Skip semicolons
			p.s.Consume()
			continue

		case csslexer.AtKeywordToken:
			// Handle at-rules (nested @media, @supports, etc.)
			if nestingType != nesting.NestingTypeNone {
				nestedRule, err := p.consumeNestedAtRule(nestingType, styleRule)
				if err == nil && nestedRule != nil {
					childRules = append(childRules, nestedRule)
				}
			} else {
				// Skip unknown at-rules in regular style blocks
				p.skipToNextDeclarationOrRule()
			}

		case csslexer.IdentToken:
			// Try to parse as CSS declaration first
			state := p.s.State()
			decl, err := p.consumeDeclaration()
			if err == nil && decl != nil {
				declarations = append(declarations, decl)
			} else {
				// If declaration parsing failed, try as nested style rule
				state.Restore()
				if nestingType != nesting.NestingTypeNone {
					nestedRule, err := p.consumeNestedStyleRule(nestingType, styleRule)
					if err == nil && nestedRule != nil {
						childRules = append(childRules, nestedRule)
					} else {
						// Skip to next valid token if nested rule parsing also failed
						p.skipToNextDeclarationOrRule()
					}
				} else {
					// Skip invalid declaration in non-nested context
					p.skipToNextDeclarationOrRule()
				}
			}

		default:
			// Handle other tokens that might start nested rules
			if nestingType != nesting.NestingTypeNone {
				state := p.s.State()
				nestedRule, err := p.consumeNestedStyleRule(nestingType, styleRule)
				if err == nil && nestedRule != nil {
					childRules = append(childRules, nestedRule)
				} else {
					state.Restore()
					p.skipToNextDeclarationOrRule()
				}
			} else {
				// Skip unknown tokens in regular style blocks
				p.skipToNextDeclarationOrRule()
			}
		}
	}

	// Store the parsed declarations and child rules
	styleRule.Declarations = declarations
	for _, childRule := range childRules {
		// TODO: Properly convert childRule to css.GenericRule
		// For now, we create a placeholder GenericRule
		_ = childRule // Avoid unused variable error
		styleRule.Rules = append(styleRule.Rules, &css.GenericRule{})
	}

	return nil
}

// consumeDeclaration parses a single CSS declaration (property: value)
func (p *Parser) consumeDeclaration() (*css.Declaration, error) {
	// Expect an identifier token (property name)
	token := p.s.Peek()
	if token.Type != csslexer.IdentToken {
		return nil, errors.New("expected property name")
	}

	propertyName := strings.TrimSpace(token.Value)
	p.s.ConsumeIncludingWhitespace()

	// Expect colon
	if p.s.Peek().Type != csslexer.ColonToken {
		return nil, errors.New("expected ':' after property name")
	}
	p.s.Consume() // Consume colon
	p.s.ConsumeWhitespace()

	// Parse property value
	var valueTokens []string
	var important bool

	// Consume tokens until we hit semicolon, EOF, or closing brace
	for {
		token := p.s.Peek()
		if token.Type == csslexer.SemicolonToken ||
			token.Type == csslexer.EOFToken ||
			token.Type == csslexer.RightBraceToken {
			break
		}

		// Check for !important
		if token.Type == csslexer.DelimiterToken && token.Value == "!" {
			p.s.Consume()
			p.s.ConsumeWhitespace()
			nextToken := p.s.Peek()
			if nextToken.Type == csslexer.IdentToken &&
				strings.ToLower(nextToken.Value) == "important" {
				important = true
				p.s.ConsumeIncludingWhitespace()
				break
			} else {
				// Not !important, add the "!" back to value
				valueTokens = append(valueTokens, "!")
			}
		} else {
			valueTokens = append(valueTokens, token.Value)
			p.s.Consume()
		}
	}

	if len(valueTokens) == 0 {
		return nil, errors.New("empty property value")
	}

	value := strings.TrimSpace(strings.Join(valueTokens, ""))
	if value == "" {
		return nil, errors.New("empty property value")
	}

	// Consume semicolon if present
	if p.s.Peek().Type == csslexer.SemicolonToken {
		p.s.Consume()
	}

	return &css.Declaration{
		Property:  propertyName,
		Value:     value,
		Important: important,
	}, nil
}

// skipToNextDeclarationOrRule skips tokens until the next declaration or rule
func (p *Parser) skipToNextDeclarationOrRule() {
	for !p.s.AtEnd() {
		token := p.s.Peek()
		if token.Type == csslexer.SemicolonToken ||
			token.Type == csslexer.RightBraceToken ||
			token.Type == csslexer.EOFToken {
			if token.Type == csslexer.SemicolonToken {
				p.s.Consume()
			}
			break
		}
		p.s.Consume()
	}
}

// consumeNestedAtRule handles nested at-rules like @media, @supports within style rules
func (p *Parser) consumeNestedAtRule(nestingType nesting.NestingTypeType, parentRule *css.StyleRule) (*css.StyleRule, error) {
	// For now, skip nested at-rules as they require more complex parsing
	// In a full implementation, this would handle @media, @supports, @container, etc.
	p.skipToNextDeclarationOrRule()
	return nil, errors.New("nested at-rules not implemented yet")
}

// consumeNestedStyleRule handles nested style rules within CSS nesting
func (p *Parser) consumeNestedStyleRule(nestingType nesting.NestingTypeType, parentRule *css.StyleRule) (*css.StyleRule, error) {
	// Parse nested style rule with the current nesting context
	return p.consumeStyleRule(nestingType, parentRule, true)
}

func (p *Parser) consumeKeyframeStyleRule() (*css.StyleRule, error) {
	return nil, errors.New("not implemented yet")
}
