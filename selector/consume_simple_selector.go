package selector

import (
	"errors"
	"strconv"
	"strings"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/nesting"
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
		case SelectorPseudoIs, SelectorPseudoWhere:
			// :is() and :where() use forgiving nested selector lists
			selectorList, err := sp.consumeForgivingNestedSelectorList()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse selector list for :is()/:where()")
			}
			pseudoData.SelectorList = selectorList
			return nil

		case SelectorPseudoHost, SelectorPseudoHostContext:
			// found_host_in_compound_ = true
			// :host() and :host-context() use compound selector lists
			selectorList, err := sp.consumeCompoundSelectorList()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse compound selector list for :host()/:host-context()")
			}

			// :host() can only have single complex selectors
			if pseudoData.PseudoType == SelectorPseudoHost && len(selectorList) > 1 {
				return errors.New("invalid pseudo selector: :host() can only contain single complex selector")
			}
			if pseudoData.PseudoType == SelectorPseudoHostContext && len(selectorList) > 1 {
				return errors.New("invalid pseudo selector: :host-context() can only contain single complex selector")
			}

			pseudoData.SelectorList = selectorList
			return nil

		case SelectorPseudoAny, SelectorPseudoCue:
			// :any() and :cue() use compound selector lists
			selectorList, err := sp.consumeCompoundSelectorList()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse compound selector list for :any()/:cue()")
			}
			pseudoData.SelectorList = selectorList
			return nil

		case SelectorPseudoHas:
			// :has() uses relative selector lists
			selectorList, err := sp.consumeRelativeSelectorList()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse relative selector list for :has()")
			}
			pseudoData.SelectorList = selectorList
			flags.Set(SelectorFlagContainsPseudo)
			flags.Set(SelectorFlagContainsComplexSelector)
			return nil

		case SelectorPseudoNot:
			// :not() uses nested selector lists
			selectorList, err := sp.consumeNestedSelectorList()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse nested selector list for :not()")
			}
			pseudoData.SelectorList = selectorList
			return nil

		case SelectorPseudoPicker, SelectorPseudoDir, SelectorPseudoState:
			// These pseudo-classes take a simple identifier argument
			token := ts.Peek()
			if token.Type != csslexer.IdentToken {
				return errors.New("invalid pseudo selector: expected identifier argument")
			}
			pseudoData.Argument = token.Value
			ts.ConsumeIncludingWhitespace()
			if !ts.AtEnd() {
				return errors.New("invalid pseudo selector: unexpected tokens after argument")
			}
			return nil

		case SelectorPseudoPart:
			// ::part() takes a space-separated list of identifiers
			var parts []string
			for !ts.AtEnd() {
				token := ts.Peek()
				if token.Type != csslexer.IdentToken {
					return errors.New("invalid pseudo selector: expected identifier in part list")
				}
				parts = append(parts, token.Value)
				ts.ConsumeIncludingWhitespace()
			}
			if len(parts) == 0 {
				return errors.New("invalid pseudo selector: part list cannot be empty")
			}
			pseudoData.IdentList = parts
			return nil

		case SelectorPseudoActiveViewTransitionType:
			// :active-view-transition-type() takes a comma-separated list of identifiers
			var types []string
			for {
				token := ts.Peek()
				if token.Type != csslexer.IdentToken {
					return errors.New("invalid pseudo selector: expected identifier in type list")
				}
				types = append(types, token.Value)
				ts.ConsumeIncludingWhitespace()

				if ts.AtEnd() {
					break
				}

				comma := ts.Peek()
				if comma.Type != csslexer.CommaToken {
					return errors.New("invalid pseudo selector: expected comma in type list")
				}
				ts.ConsumeIncludingWhitespace()
				if ts.AtEnd() {
					return errors.New("invalid pseudo selector: trailing comma in type list")
				}
			}
			pseudoData.IdentList = types
			return nil

		case SelectorPseudoViewTransitionGroup,
			SelectorPseudoViewTransitionGroupChildren,
			SelectorPseudoViewTransitionImagePair,
			SelectorPseudoViewTransitionOld,
			SelectorPseudoViewTransitionNew:
			// These pseudo-elements take a name and optional classes
			var nameAndClasses []string

			// Check for view transition class (starts with '.')
			if ts.Peek().Type == csslexer.DelimiterToken && ts.Peek().Value == "." {
				nameAndClasses = append(nameAndClasses, "*") // Universal selector for classes
			}

			if len(nameAndClasses) == 0 {
				token := ts.Peek()
				if token.Type == csslexer.DelimiterToken && token.Value == "*" {
					nameAndClasses = append(nameAndClasses, "*")
					ts.Consume()
				} else if token.Type == csslexer.IdentToken {
					nameAndClasses = append(nameAndClasses, token.Value)
					ts.Consume()
				} else {
					return errors.New("invalid pseudo selector: expected name or * for view transition")
				}
			}

			// Parse view transition classes
			for !ts.AtEnd() && ts.Peek().Type != csslexer.WhitespaceToken {
				if ts.Peek().Type != csslexer.DelimiterToken || ts.Consume().Value != "." {
					return errors.New("invalid pseudo selector: expected '.' before class name")
				}

				token := ts.Peek()
				if token.Type != csslexer.IdentToken {
					return errors.New("invalid pseudo selector: expected class name after '.'")
				}
				nameAndClasses = append(nameAndClasses, token.Value)
				ts.Consume()
			}

			ts.ConsumeWhitespace()
			if !ts.AtEnd() {
				return errors.New("invalid pseudo selector: unexpected tokens after view transition")
			}

			pseudoData.IdentList = nameAndClasses
			return nil

		case SelectorPseudoSlotted:
			// ::slotted() takes a single compound selector
			sel, err := sp.consumeCompoundSelectorAsComplexSelector()
			ts.ConsumeWhitespace()
			if err != nil || !ts.AtEnd() {
				return errors.New("invalid pseudo selector: failed to parse compound selector for ::slotted()")
			}
			pseudoData.SelectorList = []*Selector{sel}
			return nil

		case SelectorPseudoLang:
			// :lang() supports extended language ranges
			var langs []string

			for !ts.AtEnd() {
				var langRange strings.Builder
				token := ts.Peek()

				// Initial subtag: identifier, string, or wildcard
				if token.Type == csslexer.IdentToken {
					value := token.Value
					// Reject identifiers starting with hyphen
					if len(value) > 0 && value[0] == '-' {
						return errors.New("invalid pseudo selector: lang range cannot start with hyphen")
					}
					langRange.WriteString(value)
					ts.Consume()
				} else if token.Type == csslexer.StringToken {
					langRange.WriteString(token.Value)
					ts.Consume()
				} else if token.Type == csslexer.DelimiterToken && token.Value == "*" {
					langRange.WriteString("*")
					ts.Consume()
				} else {
					return errors.New("invalid pseudo selector: invalid lang range start")
				}

				// Parse hyphen-separated subtags
				for !ts.AtEnd() {
					next := ts.Peek()
					if next.Type == csslexer.DelimiterToken && next.Value == "-" {
						langRange.WriteString("-")
						ts.Consume()

						if ts.AtEnd() {
							return errors.New("invalid pseudo selector: trailing hyphen in lang range")
						}

						afterHyphen := ts.Peek()
						if afterHyphen.Type == csslexer.IdentToken {
							langRange.WriteString(afterHyphen.Value)
							ts.Consume()
						} else if afterHyphen.Type == csslexer.StringToken {
							langRange.WriteString(afterHyphen.Value)
							ts.Consume()
						} else if afterHyphen.Type == csslexer.DelimiterToken && afterHyphen.Value == "*" {
							langRange.WriteString("*")
							ts.Consume()
						} else if afterHyphen.Type == csslexer.NumberToken {
							// TODO: Check if it's a positive integer
							langRange.WriteString(afterHyphen.Value)
							ts.Consume()
						} else {
							return errors.New("invalid pseudo selector: unexpected token after hyphen")
						}
					} else if next.Type == csslexer.DelimiterToken && next.Value == "*" && strings.HasSuffix(langRange.String(), "-") {
						langRange.WriteString("*")
						ts.Consume()
					} else if next.Type == csslexer.IdentToken && len(next.Value) > 0 && next.Value[0] == '-' {
						langRange.WriteString(next.Value)
						ts.Consume()
					} else {
						break
					}
				}

				langs = append(langs, langRange.String())
				ts.ConsumeWhitespace()

				if !ts.AtEnd() {
					if ts.Peek().Type != csslexer.CommaToken {
						return errors.New("invalid pseudo selector: expected comma in lang list")
					}
					ts.ConsumeIncludingWhitespace()
					if ts.AtEnd() {
						return errors.New("invalid pseudo selector: trailing comma in lang list")
					}
				}
			}

			if len(langs) == 0 {
				return errors.New("invalid pseudo selector: empty lang list")
			}

			pseudoData.ArgumentList = langs
			return nil

		case SelectorPseudoNthChild,
			SelectorPseudoNthLastChild,
			SelectorPseudoNthOfType,
			SelectorPseudoNthLastOfType:
			// Parse An+B notation
			a, b, err := sp.consumeANPlusB()
			if err != nil {
				return err
			}
			ts.ConsumeWhitespace()

			nthData := NewSelectorPseudoNthData(a, b)

			if ts.AtEnd() {
				// Simple An+B case
				pseudoData.NthData = nthData
				return nil
			}

			// Check for "of <selectors>" syntax (only for :nth-child and :nth-last-child)
			if pseudoData.PseudoType != SelectorPseudoNthChild &&
				pseudoData.PseudoType != SelectorPseudoNthLastChild {
				return errors.New("invalid pseudo selector: unexpected tokens after An+B")
			}

			subSelectors, err := sp.consumeNthChildOfSelectors()
			if err != nil {
				return err
			}
			ts.ConsumeWhitespace()
			if !ts.AtEnd() {
				return errors.New("invalid pseudo selector: unexpected tokens after selector list")
			}

			nthData.SelectorList = subSelectors
			pseudoData.NthData = nthData
			return nil

		case SelectorPseudoScrollButton:
			// ::scroll-button() takes a direction keyword or *
			token := ts.Peek()
			if token.Type == csslexer.IdentToken {
				// Check if it's a valid scroll button direction keyword
				switch strings.ToLower(token.Value) {
				case "up", "down", "left", "right", "block-start", "block-end", "inline-start", "inline-end":
					pseudoData.Argument = token.Value
				default:
					return errors.New("invalid pseudo selector: invalid scroll button direction")
				}
			} else if token.Type == csslexer.DelimiterToken && token.Value == "*" {
				pseudoData.Argument = "*"
			} else {
				return errors.New("invalid pseudo selector: expected direction or * for scroll button")
			}
			ts.ConsumeIncludingWhitespace()
			if !ts.AtEnd() {
				return errors.New("invalid pseudo selector: unexpected tokens after scroll button argument")
			}
			return nil

		case SelectorPseudoHighlight:
			// ::highlight() takes a simple identifier argument
			token := ts.Peek()
			if token.Type != csslexer.IdentToken {
				return errors.New("invalid pseudo selector: expected identifier for highlight")
			}
			pseudoData.Argument = token.Value
			ts.ConsumeIncludingWhitespace()
			if !ts.AtEnd() {
				return errors.New("invalid pseudo selector: unexpected tokens after highlight argument")
			}
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
	sp.tokenStream.Consume() // Consume the '&' delimiter token

	var flags SelectorListFlagType
	flags.Set(SelectorFlagContainsScopeOrParent)
	flags.Set(SelectorFlagContainsPseudo)          // Nesting parent can contain pseudo selectors
	flags.Set(SelectorFlagContainsComplexSelector) // Nesting parent can contain complex selectors

	sel := &SimpleSelector{
		Match: SelectorMatchPseudoClass,
		Data:  NewSelectorDataPseudo("parent", SelectorPseudoParent), // & is represented as the parent pseudo-class
	}

	return sel, flags, nil
}

// consumeANPlusB parses An+B notation for nth-child selectors
func (sp *SelectorParser) consumeANPlusB() (int, int, error) {
	if sp.tokenStream.AtEnd() {
		return 0, 0, errors.New("unexpected end of input")
	}

	// Check for block tokens first - An+B notation cannot start with block tokens
	if token_stream.IsBlockToken(sp.tokenStream.Peek().Type) {
		return 0, 0, errors.New("An+B notation cannot start with block token")
	}

	token := sp.tokenStream.Consume()

	// Handle simple number case (just B)
	if token.Type == csslexer.NumberToken {
		value, err := strconv.Atoi(token.Value)
		if err != nil {
			return 0, 0, errors.New("invalid number in An+B notation")
		}
		return 0, value, nil
	}

	// Handle "odd" and "even" keywords
	if token.Type == csslexer.IdentToken {
		switch strings.ToLower(token.Value) {
		case "odd":
			return 2, 1, nil
		case "even":
			return 2, 0, nil
		}

		// Check if this is an ident starting with 'n' or '-n'
		lowerValue := strings.ToLower(token.Value)
		if strings.HasPrefix(lowerValue, "n") || strings.HasPrefix(lowerValue, "-n") {
			// Parse dimension-like tokens (e.g., "n", "-n", "n-3")
			if token.Value == "n" {
				sp.tokenStream.ConsumeWhitespace()
				return sp.parseOptionalB(1)
			}
			if token.Value == "-n" {
				sp.tokenStream.ConsumeWhitespace()
				return sp.parseOptionalB(-1)
			}

			// Parse coefficient and remainder if present
			return sp.parseNWithCoefficient(token.Value)
		}

		return 0, 0, errors.New("invalid An+B notation")
	}

	// Handle dimension tokens (e.g., "2n")
	if token.Type == csslexer.DimensionToken && strings.HasSuffix(strings.ToLower(token.Value), "n") {
		// For dimension tokens, we need to parse the numeric part ourselves
		// since the full token value includes the unit
		numPart := strings.TrimSuffix(token.Value, "n")
		numPart = strings.TrimSuffix(strings.ToLower(numPart), "n")

		if numPart == "" || numPart == "+" {
			sp.tokenStream.ConsumeWhitespace()
			return sp.parseOptionalB(1)
		} else if numPart == "-" {
			sp.tokenStream.ConsumeWhitespace()
			return sp.parseOptionalB(-1)
		}

		a, err := strconv.Atoi(numPart)
		if err != nil {
			return 0, 0, errors.New("invalid coefficient in dimension token")
		}

		// Check if next token is a signed number (like "+5" or "-3") which should be combined
		nextToken := sp.tokenStream.Peek()
		if nextToken.Type == csslexer.NumberToken {
			// Handle cases like "2n+5" where +5 comes as a single token
			bStr := nextToken.Value
			if strings.HasPrefix(bStr, "+") || strings.HasPrefix(bStr, "-") {
				sp.tokenStream.Consume()
				b, err := strconv.Atoi(bStr)
				if err != nil {
					return 0, 0, errors.New("invalid B value in dimension+number")
				}
				return a, b, nil
			}
		}

		sp.tokenStream.ConsumeWhitespace()
		return sp.parseOptionalB(a)
	}

	// Handle "+n" case
	if token.Type == csslexer.DelimiterToken && token.Value == "+" {
		nextToken := sp.tokenStream.Peek()
		if nextToken.Type == csslexer.IdentToken && strings.ToLower(nextToken.Value) == "n" {
			sp.tokenStream.Consume() // consume 'n'
			sp.tokenStream.ConsumeWhitespace()
			return sp.parseOptionalB(1)
		}
	}

	return 0, 0, errors.New("invalid An+B notation")
}

// parseNWithCoefficient parses tokens like "2n-3", "-n+1", etc.
func (sp *SelectorParser) parseNWithCoefficient(value string) (int, int, error) {
	lower := strings.ToLower(value)

	// Find the 'n'
	nIndex := strings.Index(lower, "n")
	if nIndex == -1 {
		return 0, 0, errors.New("invalid coefficient notation")
	}

	// Parse coefficient (A)
	var a int
	coeffPart := lower[:nIndex]
	if coeffPart == "" || coeffPart == "+" {
		a = 1
	} else if coeffPart == "-" {
		a = -1
	} else {
		var err error
		a, err = strconv.Atoi(coeffPart)
		if err != nil {
			return 0, 0, errors.New("invalid coefficient")
		}
	}

	// Parse remainder (B) if present
	remainder := lower[nIndex+1:]
	if remainder == "" {
		sp.tokenStream.ConsumeWhitespace()
		return sp.parseOptionalB(a)
	}

	// Direct B parsing from remainder like "n-3"
	b, err := strconv.Atoi(remainder)
	if err != nil {
		return 0, 0, errors.New("invalid constant")
	}

	return a, b, nil
}

// parseOptionalB parses optional +B or -B part after An
func (sp *SelectorParser) parseOptionalB(a int) (int, int, error) {
	// Check for optional + or - B
	token := sp.tokenStream.Peek()
	if token.Type != csslexer.DelimiterToken {
		return a, 0, nil
	}

	var sign int
	switch token.Value {
	case "+":
		sign = 1
		sp.tokenStream.ConsumeIncludingWhitespace()
	case "-":
		sign = -1
		sp.tokenStream.ConsumeIncludingWhitespace()
	default:
		return a, 0, nil
	}

	// Parse the B value
	bToken := sp.tokenStream.Peek()
	if bToken.Type != csslexer.NumberToken {
		return 0, 0, errors.New("expected number after +/- in An+B")
	}

	bValue, err := strconv.Atoi(bToken.Value)
	if err != nil {
		return 0, 0, errors.New("invalid number for B value in An+B")
	}

	b := bValue * sign
	sp.tokenStream.Consume()

	return a, b, nil
}

// consumeNthChildOfSelectors parses the "of <selector-list>" part
func (sp *SelectorParser) consumeNthChildOfSelectors() ([]*Selector, error) {
	token := sp.tokenStream.Peek()
	if token.Type != csslexer.IdentToken || strings.ToLower(token.Value) != "of" {
		return nil, errors.New("expected 'of' keyword")
	}
	sp.tokenStream.ConsumeIncludingWhitespace()

	return sp.consumeComplexSelectorList(nesting.NestingTypeNone)
}
