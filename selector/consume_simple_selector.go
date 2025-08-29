package selector

import (
	"errors"
	"strings"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/token_stream"
)

// ConsumeSimpleSelector consumes a simple selector from the token stream.
func (sp *SelectorParser) consumeSimpleSelector() (*SimpleSelector, SelectorListFlagType, error) {
	token := sp.tokenStream.Peek()
	switch token.Type {
	case csslexer.HashToken:
		ss, err := sp.consumeId()
		return ss, 0, err

	case csslexer.DelimiterToken:
		switch token.Value {
		case ".":
			ss, err := sp.consumeClass()
			return ss, 0, err
		case "&":
			return sp.consumeNestingParent()
		default:
			return nil, 0, errors.New("invalid selector: unknown delimiter")
		}

	case csslexer.LeftBracketToken:
		ss, err := sp.consumeAttribute()
		return ss, 0, err

	case csslexer.ColonToken:
		ss, flags, err := sp.consumePseudo()
		if err != nil {
			return nil, 0, err
		}
		flags.Set(SelectorFlagContainsPseudo)
		return ss, flags, nil

	default:
		return nil, 0, errors.New("invalid selector: expected simple selector")
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
		Data:  NewSelectorData(token.Value),
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
		Data:  NewSelectorData(token.Value),
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
				Data:  NewSelectorDataAttr(nameStr, "", SelectorAttrMatchCaseSensitive),
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
			Match: matchType,
			Data:  NewSelectorDataAttr(nameStr, valueToken.Value, attrMatchType),
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
func (sp *SelectorParser) consumeAttributeFlags() SelectorAttrMatchType {
	if sp.tokenStream.Peek().Type != csslexer.IdentToken {
		return SelectorAttrMatchCaseSensitive // Default to case-sensitive if no flags are specified
	}

	token := sp.tokenStream.ConsumeIncludingWhitespace() // Consume the identifier token

	if strings.ToLower(token.Value) == "i" {
		return SelectorAttrMatchCaseInsensitive
	} else if strings.ToLower(token.Value) == "s" {
		return SelectorAttrMatchCaseSensitiveAlways
	} else {
		return SelectorAttrMatchCaseSensitive // Default to case-sensitive if unknown flag
	}
}

// consumePseudo consumes a pseudo-class or pseudo-element selector from
// the token stream.
func (sp *SelectorParser) consumePseudo() (*SimpleSelector, SelectorListFlagType, error) {
	sp.tokenStream.Consume() // Consume the colon token

	colons := 1
	if sp.tokenStream.Peek().Type == csslexer.ColonToken {
		sp.tokenStream.Consume() // Consume the second colon for pseudo-elements
		colons++
	}

	token := sp.tokenStream.Peek()
	if token.Type != csslexer.IdentToken && token.Type != csslexer.FunctionToken {
		return nil, 0, errors.New("invalid pseudo selector: expected ident-token or function-token after colon")
	}

	var flags SelectorListFlagType
	sel := &SimpleSelector{}

	switch colons {
	case 1:
		// Pseudo-class
		sel.Match = SelectorMatchPseudoClass

	case 2:
		// Pseudo-element
		sel.Match = SelectorMatchPseudoElement

	default:
		return nil, 0, errors.New("invalid pseudo selector: too many colons")
	}

	pseudoName := strings.ToLower(token.Value)
	pseudoType := parsePseudoType(pseudoName, token.Type == csslexer.FunctionToken)
	sel.Data = NewSelectorDataPseudo(pseudoName, pseudoType)

	if sel.Match == SelectorMatchPseudoElement {
		// TODO: check if current state disallows pseudo element selectors
	}

	if token.Type == csslexer.IdentToken {
		sp.tokenStream.Consume() // Consume the ident token

		pseudoData := sel.Data.(*SelectorDataPseudo)
		switch pseudoData.PseudoType {
		case SelectorPseudoUnknown:
			return nil, 0, errors.New("invalid pseudo selector: unknown pseudo type")
		case SelectorPseudoHost:
			// TODO: found_host_in_compound_ = true;
		case SelectorPseudoScope:
			flags.Set(SelectorFlagContainsScopeOrParent)
		}

		return sel, flags, nil
	}

	// for function tokens

	err := sp.tokenStream.ConsumeBlockToEnd(csslexer.RightParenthesisToken, func(ts *token_stream.TokenStream) error {
		ts.ConsumeWhitespace() // Consume any whitespace before the function arguments

		pseudoData := sel.Data.(*SelectorDataPseudo)
		switch pseudoData.PseudoType {
		case SelectorPseudoIs:
			// TODO
			return nil

		case SelectorPseudoWhere:
			// TODO
			return nil

		case SelectorPseudoHost, SelectorPseudoHostContext:
			// found_host_in_compound_ = true
			fallthrough

		case SelectorPseudoAny, SelectorPseudoCue:
			// TODO
			return nil

		case SelectorPseudoHas:
			// TODO
			return nil

		case SelectorPseudoNot:
			// TODO
			return nil

		case SelectorPseudoPicker, SelectorPseudoDir, SelectorPseudoState:
			// TODO
			return nil

		case SelectorPseudoPart:
			// TODO
			return nil

		case SelectorPseudoActiveViewTransitionType:
			// TODO
			return nil

		case SelectorPseudoViewTransitionGroup,
			SelectorPseudoViewTransitionGroupChildren,
			SelectorPseudoViewTransitionImagePair,
			SelectorPseudoViewTransitionOld,
			SelectorPseudoViewTransitionNew:
			// TODO
			return nil

		case SelectorPseudoSlotted:
			// TODO
			return nil

		case SelectorPseudoLang:
			// TODO
			return nil

		case SelectorPseudoNthChild,
			SelectorPseudoNthLastChild,
			SelectorPseudoNthOfType,
			SelectorPseudoNthLastOfType:
			// TODO
			return nil

		case SelectorPseudoScrollButton:
			// TODO
			return nil

		case SelectorPseudoHighlight:
			// TODO
			return nil

		default:
			return errors.New("invalid pseudo selector: unknown pseudo type")
		}
	})
	if err != nil {
		return nil, 0, err
	}
	return sel, flags, nil
}

func (sp *SelectorParser) consumeNestingParent() (*SimpleSelector, SelectorListFlagType, error) {
	return nil, 0, errors.New("not implemented: SelectorParser.consumeNestingParent")
}
