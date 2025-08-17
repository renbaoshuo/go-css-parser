package selector

import (
	"errors"
	"strings"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/token_stream"
)

// ConsumeSimpleSelector consumes a simple selector from the token stream.
func (sp *SelectorParser) consumeSimpleSelector() (*SimpleSelector, error) {
	token := sp.tokenStream.Peek()
	switch token.Type {
	case csslexer.HashToken:
		return sp.consumeId()

	case csslexer.DelimiterToken:
		switch token.Value {
		case ".":
			return sp.consumeClass()
		case "&":
			return sp.consumeNestingParent()
		default:
			return nil, errors.New("invalid selector: unknown delimiter")
		}

	case csslexer.LeftBracketToken:
		return sp.consumeAttribute()

	case csslexer.ColonToken:
		return sp.consumePseudo()

	default:
		return nil, errors.New("invalid selector: expected simple selector")
	}
}

// consumeId consumes an ID selector from the token stream.
//
// The caller makes sure that the token stream is positioned at a hash token
// before calling this method.
//
// Returns a SimpleSelector representing the ID selector.
func (sp *SelectorParser) consumeId() (*SimpleSelector, error) {
	token := sp.tokenStream.Consume() // Consume the hash token
	return &SimpleSelector{
		Match: SelectorMatchId,
		Value: token.Value,
	}, nil
}

// consumeClass consumes a class selector from the token stream.
//
// The caller makes sure that the token stream is positioned at a delimiter
// token with a single character '.' before calling this method.
//
// Returns a SimpleSelector representing the class selector.
//
// If the next token is not an identifier, it returns an error.
func (sp *SelectorParser) consumeClass() (*SimpleSelector, error) {
	sp.tokenStream.Consume() // Consume the delimiter token ('.')
	token := sp.tokenStream.Peek()
	if token.Type != csslexer.IdentToken {
		return nil, errors.New("invalid selector: expected class name after '.'")
	}
	sp.tokenStream.Consume() // Consume the class name token
	return &SimpleSelector{
		Match: SelectorMatchClass,
		Value: token.Value,
	}, nil
}

// consumeName consumes a name from the token stream.
//
// The caller makes sure that the token stream is positioned at a valid left
// bracket token before calling this method.
//
// Returns a SimpleSelector representing the attribute selector.
//
// If the tokens cannot be consumed correctly, it returns an error, and skip
// to the token after right bracket token.
func (sp *SelectorParser) consumeAttribute() (*SimpleSelector, error) {
	var sel *SimpleSelector

	err := sp.tokenStream.ConsumeBlock(func(_ *token_stream.TokenStream) error {
		// consume the whitespace before the attribute selector
		sp.tokenStream.ConsumeWhitespace()

		name, namespace, hasQName := sp.consumeName()

		if !hasQName {
			return errors.New("invalid attribute selector: missing name")
		}

		if name == "" {
			return errors.New("invalid attribute selector: name cannot be empty")
		}

		// TODO: Handle namespace uri
		// FIXME: when serializing, the "|" will be wrongly serialized as "\|"
		nameStr := name
		if namespace != "" {
			nameStr = namespace + "|" + name
		}

		if sp.tokenStream.AtEnd() {
			sel = &SimpleSelector{
				Match: SelectorMatchAttributeSet,
				Value: nameStr,
			}

			return nil
		}

		matchType := sp.consumeAttributeMatch()
		if matchType == SelectorMatchUnknown {
			return errors.New("invalid attribute selector: unknown match type")
		}

		valueToken := sp.tokenStream.Peek()
		if valueToken.Type != csslexer.IdentToken && valueToken.Type != csslexer.StringToken {
			return errors.New("invalid attribute selector: expected value after match type")
		}
		sp.tokenStream.ConsumeIncludingWhitespace() // Consume the value token

		attrMatchType := sp.consumeAttributeFlags()

		if !sp.tokenStream.AtEnd() {
			return errors.New("invalid attribute selector: unexpected tokens after value")
		}

		sel = &SimpleSelector{
			Match:     matchType,
			Value:     nameStr,
			AttrValue: valueToken.Value,
			AttrMatch: attrMatchType,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sel, nil
}

// consumeAttributeMatch consumes the attribute match type from the token stream.
func (sp *SelectorParser) consumeAttributeMatch() SelectorMatchType {
	token := sp.tokenStream.Peek()
	switch token.Type {
	case csslexer.IncludeMatchToken:
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorMatchAttributeList

	case csslexer.DashMatchToken:
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorMatchAttributeHyphen

	case csslexer.PrefixMatchToken:
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorMatchAttributeBegin

	case csslexer.SuffixMatchToken:
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorMatchAttributeEnd

	case csslexer.SubstringMatchToken:
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorMatchAttributeContain

	case csslexer.DelimiterToken:
		if token.Value == "=" {
			sp.tokenStream.ConsumeIncludingWhitespace()
			return SelectorMatchAttributeExact
		}

		return SelectorMatchUnknown // Invalid attribute match type

	default:
		return SelectorMatchUnknown
	}
}

// consumeAttributeFlags consumes the attribute flags from the token stream.
func (sp *SelectorParser) consumeAttributeFlags() SelectorAttributeMatchType {
	if sp.tokenStream.Peek().Type != csslexer.IdentToken {
		return SelectorAttributeMatchCaseSensitive // Default to case-sensitive if no flags are specified
	}

	token := sp.tokenStream.ConsumeIncludingWhitespace() // Consume the identifier token

	if strings.ToLower(token.Value) == "i" {
		return SelectorAttributeMatchCaseInsensitive
	} else if strings.ToLower(token.Value) == "s" {
		return SelectorAttributeMatchCaseSensitiveAlways
	} else {
		return SelectorAttributeMatchCaseSensitive // Default to case-sensitive if unknown flag
	}
}

// consumePseudo consumes a pseudo-class or pseudo-element selector from
// the token stream.
func (sp *SelectorParser) consumePseudo() (*SimpleSelector, error) {
	return nil, errors.New("not implemented: SelectorParser.consumePseudo")
}

func (sp *SelectorParser) consumeNestingParent() (*SimpleSelector, error) {
	return nil, errors.New("not implemented: SelectorParser.consumeNestingParent")
}
