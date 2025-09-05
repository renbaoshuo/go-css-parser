package selector

import (
	"strconv"
	"strings"

	"go.baoshuo.dev/cssutil"
)

// ===== SelectorDataPseudo =====

type SelectorDataPseudo struct {
	PseudoType   SelectorPseudoType     // The type of pseudo-selector.
	PseudoName   string                 // The original pseudo name as parsed.
	Argument     string                 // Used for :contains, :lang, :dir, etc.
	ArgumentList []string               // Used for :lang
	SelectorList []*Selector            // For pseudo-classes that take a selector list as an argument, e.g., :is(), :not()
	IdentList    []string               // Used for ::part(), :active-view-transition-type().
	NthData      *SelectorPseudoNthData // Used for :nth-child, :nth-last-child, etc.
}

func NewSelectorDataPseudo(pseudoName string, pseudoType SelectorPseudoType) *SelectorDataPseudo {
	return &SelectorDataPseudo{
		PseudoType: pseudoType,
		PseudoName: pseudoName,
	}
}

func (d *SelectorDataPseudo) String(match SelectorMatchType) string {
	var prefix string
	switch match {
	case SelectorMatchPseudoClass:
		prefix = ":"
	case SelectorMatchPseudoElement:
		prefix = "::"
	case SelectorMatchPagePseudoClass:
		prefix = "@page :"
	default:
		prefix = ":"
	}

	result := prefix + cssutil.SerializeIdentifier(d.PseudoName)

	// Handle different pseudo types with their specific arguments
	switch d.PseudoType {
	case SelectorPseudoNthChild, SelectorPseudoNthLastChild, SelectorPseudoNthOfType, SelectorPseudoNthLastOfType:
		// Use NthData for nth-* selectors
		if d.NthData != nil {
			result += "(" + d.NthData.String() + ")"
		}

	case SelectorPseudoIs, SelectorPseudoNot, SelectorPseudoWhere, SelectorPseudoHas:
		// Use SelectorList for :is(), :not(), :where(), :has()
		if len(d.SelectorList) > 0 {
			selectorStrs := make([]string, 0, len(d.SelectorList))
			for _, sel := range d.SelectorList {
				selectorStrs = append(selectorStrs, sel.String())
			}
			result += "(" + cssutil.SerializeCommaSeparatedList(selectorStrs) + ")"
		}

	case SelectorPseudoPart, SelectorPseudoActiveViewTransitionType:
		// Use IdentList for ::part(), :active-view-transition-type()
		if len(d.IdentList) > 0 {
			identStrs := make([]string, 0, len(d.IdentList))
			for _, ident := range d.IdentList {
				identStrs = append(identStrs, cssutil.SerializeIdentifier(ident))
			}
			result += "(" + cssutil.SerializeWhitespaceSeparatedList(identStrs) + ")"
		}

	case SelectorPseudoLang:
		// Use ArgumentList for :lang() with multiple arguments
		if len(d.ArgumentList) > 0 {
			argsStrs := make([]string, 0, len(d.ArgumentList))
			for _, arg := range d.ArgumentList {
				argsStrs = append(argsStrs, cssutil.SerializeString(arg))
			}
			result += "(" + cssutil.SerializeCommaSeparatedList(argsStrs) + "))"
		} else if d.Argument != "" {
			// Fallback to single Argument for backward compatibility
			result += "(" + cssutil.SerializeString(d.Argument) + ")"
		}

	case SelectorPseudoDir:
		// Use single Argument for :dir()
		if d.Argument != "" {
			result += "(" + cssutil.SerializeString(d.Argument) + ")"
		}

	default:
		// For other pseudo types, check if they have any arguments
		if d.Argument != "" {
			result += "(" + cssutil.SerializeString(d.Argument) + ")"
		}
	}

	return result
}

func (d *SelectorDataPseudo) Equals(other SelectorDataType) bool {
	otherData, ok := other.(*SelectorDataPseudo)
	if !ok {
		return false
	}

	if d.PseudoType != otherData.PseudoType ||
		d.PseudoName != otherData.PseudoName ||
		d.Argument != otherData.Argument ||
		len(d.ArgumentList) != len(otherData.ArgumentList) ||
		len(d.SelectorList) != len(otherData.SelectorList) ||
		len(d.IdentList) != len(otherData.IdentList) {
		return false
	}

	// compare NthData
	if !d.NthData.Equals(otherData.NthData) {
		return false
	}

	// compare ArgumentList
	for i, arg := range d.ArgumentList {
		if arg != otherData.ArgumentList[i] {
			return false
		}
	}

	// compare SelectorList
	for i, sel := range d.SelectorList {
		if !sel.Equals(otherData.SelectorList[i]) {
			return false
		}
	}

	// compare IdentList
	for i, ident := range d.IdentList {
		if ident != otherData.IdentList[i] {
			return false
		}
	}

	return true
}

// ===== SelectorPseudoType =====

type SelectorPseudoType int

const (
	SelectorPseudoUnknown SelectorPseudoType = iota
	SelectorPseudoActive
	SelectorPseudoActiveViewTransition
	SelectorPseudoActiveViewTransitionType
	SelectorPseudoAfter
	SelectorPseudoAny
	SelectorPseudoAnyLink
	SelectorPseudoAutofill
	SelectorPseudoAutofillPreviewed
	SelectorPseudoAutofillSelected
	SelectorPseudoBackdrop
	SelectorPseudoBefore
	SelectorPseudoCheckMark
	SelectorPseudoChecked
	SelectorPseudoCornerPresent
	SelectorPseudoCurrent
	SelectorPseudoDecrement
	SelectorPseudoDefault
	SelectorPseudoDetailsContent
	SelectorPseudoDialogInTopLayer
	SelectorPseudoDisabled
	SelectorPseudoDoubleButton
	SelectorPseudoDrag
	SelectorPseudoEmpty
	SelectorPseudoEnabled
	SelectorPseudoEnd
	SelectorPseudoFileSelectorButton
	SelectorPseudoFirstChild
	SelectorPseudoFirstLetter
	SelectorPseudoFirstLine
	SelectorPseudoFirstOfType
	SelectorPseudoFirstPage
	SelectorPseudoFocus
	SelectorPseudoFocusVisible
	SelectorPseudoFocusWithin
	SelectorPseudoFullPageMedia
	SelectorPseudoHasInterest
	SelectorPseudoHasSlotted
	SelectorPseudoHorizontal
	SelectorPseudoHover
	SelectorPseudoIncrement
	SelectorPseudoIndeterminate
	SelectorPseudoInterestHint
	SelectorPseudoInvalid
	SelectorPseudoIs
	SelectorPseudoLang
	SelectorPseudoLastChild
	SelectorPseudoLastOfType
	SelectorPseudoLeftPage
	SelectorPseudoLink
	SelectorPseudoMarker
	SelectorPseudoModal
	SelectorPseudoNoButton
	SelectorPseudoNot
	SelectorPseudoNthChild // Includes :nth-child(An+B of <selector>)
	SelectorPseudoNthLastChild
	SelectorPseudoNthLastOfType
	SelectorPseudoNthOfType
	SelectorPseudoOnlyChild
	SelectorPseudoOnlyOfType
	SelectorPseudoOptional
	SelectorPseudoParent // Written as & (in nested rules).
	SelectorPseudoPart
	SelectorPseudoPermissionElementInvalidStyle
	SelectorPseudoPermissionElementOccluded
	SelectorPseudoPermissionGranted
	SelectorPseudoPermissionIcon
	SelectorPseudoPlaceholder
	SelectorPseudoPlaceholderShown
	SelectorPseudoReadOnly
	SelectorPseudoReadWrite
	SelectorPseudoRequired
	SelectorPseudoResizer
	SelectorPseudoRightPage
	SelectorPseudoRoot
	SelectorPseudoScope
	SelectorPseudoScrollbar
	SelectorPseudoScrollbarButton
	SelectorPseudoScrollbarCorner
	SelectorPseudoScrollbarThumb
	SelectorPseudoScrollbarTrack
	SelectorPseudoScrollbarTrackPiece
	SelectorPseudoSearchText
	SelectorPseudoPickerIcon
	SelectorPseudoPicker
	SelectorPseudoSelection
	SelectorPseudoSelectorFragmentAnchor
	SelectorPseudoSingleButton
	SelectorPseudoStart
	SelectorPseudoState
	SelectorPseudoTarget
	SelectorPseudoTargetOfInterest

	// Something that was unparsable, but contained either a nesting
	// selector (&), or a :scope pseudo-class, and must therefore be kept
	// for serialization purposes.
	SelectorPseudoUnparsed
	SelectorPseudoUserInvalid
	SelectorPseudoUserValid
	SelectorPseudoValid
	SelectorPseudoVertical
	SelectorPseudoVisited
	SelectorPseudoWebKitAutofill
	SelectorPseudoWebkitAnyLink
	SelectorPseudoWhere
	SelectorPseudoWindowInactive
	SelectorPseudoFullScreen
	SelectorPseudoFullScreenAncestor
	SelectorPseudoFullscreen
	SelectorPseudoInRange
	SelectorPseudoOutOfRange
	SelectorPseudoPaused
	SelectorPseudoPictureInPicture
	SelectorPseudoPlaying
	SelectorPseudoXrOverlay
	// Pseudo-elements in UA ShadowRoots. Available in any stylesheets.
	SelectorPseudoWebKitCustomElement
	// Pseudo-elements in UA ShadowRoots. Available only in UA stylesheets.
	SelectorPseudoBlinkInternalElement
	// Pseudo-element for fragment styling
	SelectorPseudoColumn
	SelectorPseudoCue
	SelectorPseudoDefined
	SelectorPseudoDir
	SelectorPseudoFutureCue
	SelectorPseudoGrammarError
	SelectorPseudoHas
	SelectorPseudoHasDatalist
	SelectorPseudoHighlight
	SelectorPseudoHost
	SelectorPseudoHostContext
	SelectorPseudoHostHasNonAutoAppearance
	SelectorPseudoIsHtml
	SelectorPseudoListBox
	SelectorPseudoMultiSelectFocus
	SelectorPseudoOpen
	SelectorPseudoPastCue
	SelectorPseudoPatching
	SelectorPseudoPopoverInTopLayer
	SelectorPseudoPopoverOpen
	SelectorPseudoRelativeAnchor
	SelectorPseudoSlotted
	SelectorPseudoSpatialNavigationFocus
	SelectorPseudoSpellingError
	SelectorPseudoTargetText
	SelectorPseudoVideoPersistent
	SelectorPseudoVideoPersistentAncestor

	// Active ::scroll-marker styling.
	// https://drafts.csswg.org/css-overflow-5/#active-scroll-marker
	SelectorPseudoTargetCurrent

	// The following selectors are used to target pseudo-elements created for
	// ViewTransition.
	// See https://drafts.csswg.org/css-view-transitions-1/#pseudo
	// and https://drafts.csswg.org/css-view-transitions-2
	// for details.
	SelectorPseudoViewTransition
	SelectorPseudoViewTransitionGroup
	SelectorPseudoViewTransitionGroupChildren
	SelectorPseudoViewTransitionImagePair
	SelectorPseudoViewTransitionNew
	SelectorPseudoViewTransitionOld
	// Scroll markers pseudos for Carousel
	SelectorPseudoScrollMarker
	SelectorPseudoScrollMarkerGroup
	// Scroll button pseudo for Carousel
	SelectorPseudoScrollButton
)

// ===== SelectorPseudoNthData =====

// SelectorPseudoNthData represents An+B notation for nth-child selectors
type SelectorPseudoNthData struct {
	A            int         // The 'A' coefficient in An+B
	B            int         // The 'B' constant in An+B
	SelectorList []*Selector // Optional selector list for :nth-child(An+B of <selectors>)
}

func NewSelectorPseudoNthData(a, b int) *SelectorPseudoNthData {
	return &SelectorPseudoNthData{
		A: a,
		B: b,
	}
}

func (d *SelectorPseudoNthData) String() string {
	if d == nil {
		return ""
	}

	var builder strings.Builder

	// Handle An+B notation
	if d.A == 0 {
		builder.WriteString(strconv.Itoa(d.B))
	} else if d.A == 1 && d.B == 0 {
		builder.WriteString("n")
	} else if d.A == 1 {
		builder.WriteString("n")
		if d.B > 0 {
			builder.WriteString("+")
			builder.WriteString(strconv.Itoa(d.B))
		} else if d.B < 0 {
			builder.WriteString(strconv.Itoa(d.B))
		}
	} else if d.B == 0 {
		builder.WriteString(strconv.Itoa(d.A))
		builder.WriteString("n")
	} else {
		builder.WriteString(strconv.Itoa(d.A))
		builder.WriteString("n")
		if d.B > 0 {
			builder.WriteString("+")
			builder.WriteString(strconv.Itoa(d.B))
		} else {
			builder.WriteString(strconv.Itoa(d.B))
		}
	}

	// Handle "of <selector-list>" for nth-child
	if len(d.SelectorList) > 0 {
		builder.WriteString(" of ")
		selectorStrs := make([]string, 0, len(d.SelectorList))
		for _, sel := range d.SelectorList {
			selectorStrs = append(selectorStrs, sel.String())
		}
		builder.WriteString(cssutil.SerializeCommaSeparatedList(selectorStrs))
	}

	return builder.String()
}

func (d *SelectorPseudoNthData) Equals(other *SelectorPseudoNthData) bool {
	if d == nil && other == nil {
		return true
	}
	if d == nil || other == nil {
		return false
	}

	if d.A != other.A || d.B != other.B {
		return false
	}

	if len(d.SelectorList) != len(other.SelectorList) {
		return false
	}

	// Compare each selector in the list
	for i, sel := range d.SelectorList {
		if !sel.Equals(other.SelectorList[i]) {
			return false
		}
	}

	return true
}
