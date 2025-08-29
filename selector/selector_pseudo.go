package selector

import (
	"go.baoshuo.dev/cssutil"
)

// ===== SelectorDataPseudo =====

type SelectorDataPseudo struct {
	PseudoType   SelectorPseudoType // The type of pseudo-selector.
	PseudoName   string             // The original pseudo name as parsed.
	Argument     string             // Used for :contains, :lang, :dir, etc.
	ArgumentList []string           // Used for :lang
	SelectorList []*Selector        // For pseudo-classes that take a selector list as an argument, e.g., :is(), :not()
	IdentList    []string           // Used for ::part(), :active-view-transition-type().
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

	return prefix + cssutil.SerializeIdentifier(d.PseudoName)
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

	for i, arg := range d.ArgumentList {
		if arg != otherData.ArgumentList[i] {
			return false
		}
	}

	for i, sel := range d.SelectorList {
		if !sel.Equals(otherData.SelectorList[i]) {
			return false
		}
	}

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
