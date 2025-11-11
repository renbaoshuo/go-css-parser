package css

import (
	"testing"
)

func TestSelectorDataPseudoString(t *testing.T) {
	tests := []struct {
		name     string
		pseudo   *SelectorDataPseudo
		match    SelectorMatchType
		expected string
	}{
		{
			name: "simple pseudo class",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoHover,
				PseudoName: "hover",
			},
			match:    SelectorMatchPseudoClass,
			expected: ":hover",
		},
		{
			name: "pseudo element",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoBefore,
				PseudoName: "before",
			},
			match:    SelectorMatchPseudoElement,
			expected: "::before",
		},
		{
			name: "page pseudo class",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoFirstPage,
				PseudoName: "first",
			},
			match:    SelectorMatchPagePseudoClass,
			expected: "@page :first",
		},
		{
			name: "nth-child with simple An+B",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoNthChild,
				PseudoName: "nth-child",
				NthData:    NewSelectorPseudoNthData(2, 1),
			},
			match:    SelectorMatchPseudoClass,
			expected: ":nth-child(2n+1)",
		},
		{
			name: "nth-child with selector list",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoNthChild,
				PseudoName: "nth-child",
				NthData: &SelectorPseudoNthData{
					A: 1,
					B: 0,
					SelectorList: []*Selector{
						{
							Selectors: []*SimpleSelector{
								{
									Match:    SelectorMatchClass,
									Relation: SelectorRelationSubSelector,
									Data:     NewSelectorData("item"),
								},
							},
						},
					},
				},
			},
			match:    SelectorMatchPseudoClass,
			expected: ":nth-child(n of .item)",
		},
		{
			name: "is() with selector list",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoIs,
				PseudoName: "is",
				SelectorList: []*Selector{
					{
						Selectors: []*SimpleSelector{
							{
								Match:    SelectorMatchTag,
								Relation: SelectorRelationSubSelector,
								Data:     NewSelectorDataTag("", "h1"),
							},
						},
					},
					{
						Selectors: []*SimpleSelector{
							{
								Match:    SelectorMatchTag,
								Relation: SelectorRelationSubSelector,
								Data:     NewSelectorDataTag("", "h2"),
							},
						},
					},
				},
			},
			match:    SelectorMatchPseudoClass,
			expected: ":is(h1, h2)",
		},
		{
			name: "not() with selector",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoNot,
				PseudoName: "not",
				SelectorList: []*Selector{
					{
						Selectors: []*SimpleSelector{
							{
								Match:    SelectorMatchClass,
								Relation: SelectorRelationSubSelector,
								Data:     NewSelectorData("disabled"),
							},
						},
					},
				},
			},
			match:    SelectorMatchPseudoClass,
			expected: ":not(.disabled)",
		},
		{
			name: "lang() with single argument",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoLang,
				PseudoName: "lang",
				Argument:   "en",
			},
			match:    SelectorMatchPseudoClass,
			expected: ":lang(\"en\")",
		},
		{
			name: "lang() with multiple arguments",
			pseudo: &SelectorDataPseudo{
				PseudoType:   SelectorPseudoLang,
				PseudoName:   "lang",
				ArgumentList: []string{"en", "fr", "de"},
			},
			match:    SelectorMatchPseudoClass,
			expected: ":lang(\"en\", \"fr\", \"de\")",
		},
		{
			name: "dir() with argument",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoDir,
				PseudoName: "dir",
				Argument:   "ltr",
			},
			match:    SelectorMatchPseudoClass,
			expected: ":dir(\"ltr\")",
		},
		{
			name: "part() with ident list",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoPart,
				PseudoName: "part",
				IdentList:  []string{"button", "primary"},
			},
			match:    SelectorMatchPseudoElement,
			expected: "::part(button primary)",
		},
		{
			name: "unknown pseudo with argument",
			pseudo: &SelectorDataPseudo{
				PseudoType: SelectorPseudoUnknown,
				PseudoName: "custom",
				Argument:   "value",
			},
			match:    SelectorMatchPseudoClass,
			expected: ":custom(\"value\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pseudo.String(tt.match)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataPseudoEquals(t *testing.T) {
	pseudo1 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoHover,
		PseudoName: "hover",
	}

	pseudo2 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoHover,
		PseudoName: "hover",
	}

	pseudo3 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoFocus,
		PseudoName: "focus",
	}

	pseudo4 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoLang,
		PseudoName: "lang",
		Argument:   "en",
	}

	pseudo5 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoLang,
		PseudoName: "lang",
		Argument:   "en",
	}

	pseudo6 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoLang,
		PseudoName: "lang",
		Argument:   "fr",
	}

	nthData1 := NewSelectorPseudoNthData(2, 1)
	nthData2 := NewSelectorPseudoNthData(2, 1)
	nthData3 := NewSelectorPseudoNthData(3, 0)

	pseudo7 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoNthChild,
		PseudoName: "nth-child",
		NthData:    nthData1,
	}

	pseudo8 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoNthChild,
		PseudoName: "nth-child",
		NthData:    nthData2,
	}

	pseudo9 := &SelectorDataPseudo{
		PseudoType: SelectorPseudoNthChild,
		PseudoName: "nth-child",
		NthData:    nthData3,
	}

	tests := []struct {
		name     string
		pseudo1  SelectorDataType
		pseudo2  SelectorDataType
		expected bool
	}{
		{"identical simple pseudos", pseudo1, pseudo2, true},
		{"same object", pseudo1, pseudo1, true},
		{"different pseudo types", pseudo1, pseudo3, false},
		{"identical with arguments", pseudo4, pseudo5, true},
		{"different arguments", pseudo4, pseudo6, false},
		{"identical nth data", pseudo7, pseudo8, true},
		{"different nth data", pseudo7, pseudo9, false},
		{"different data types", pseudo1, NewSelectorData("hover"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pseudo1.Equals(tt.pseudo2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataPseudoEqualsComplex(t *testing.T) {
	// Test pseudos with different list lengths
	pseudo1 := &SelectorDataPseudo{
		PseudoType:   SelectorPseudoLang,
		PseudoName:   "lang",
		ArgumentList: []string{"en"},
	}

	pseudo2 := &SelectorDataPseudo{
		PseudoType:   SelectorPseudoLang,
		PseudoName:   "lang",
		ArgumentList: []string{"en", "fr"},
	}

	if pseudo1.Equals(pseudo2) {
		t.Error("expected pseudos with different ArgumentList lengths to not be equal")
	}

	// Test pseudos with selector lists
	selector1 := &Selector{
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	pseudo3 := &SelectorDataPseudo{
		PseudoType:   SelectorPseudoIs,
		PseudoName:   "is",
		SelectorList: []*Selector{selector1},
	}

	pseudo4 := &SelectorDataPseudo{
		PseudoType:   SelectorPseudoIs,
		PseudoName:   "is",
		SelectorList: []*Selector{selector1, selector1},
	}

	if pseudo3.Equals(pseudo4) {
		t.Error("expected pseudos with different SelectorList lengths to not be equal")
	}
}

func TestNewSelectorDataPseudo(t *testing.T) {
	pseudoName := "hover"
	pseudoType := SelectorPseudoHover

	pseudo := NewSelectorDataPseudo(pseudoName, pseudoType)

	if pseudo == nil {
		t.Error("expected NewSelectorDataPseudo to return non-nil")
	}
	if pseudo.PseudoName != pseudoName {
		t.Errorf("expected PseudoName %q, got %q", pseudoName, pseudo.PseudoName)
	}
	if pseudo.PseudoType != pseudoType {
		t.Errorf("expected PseudoType %v, got %v", pseudoType, pseudo.PseudoType)
	}
}

func TestSelectorPseudoNthDataString(t *testing.T) {
	tests := []struct {
		name     string
		nthData  *SelectorPseudoNthData
		expected string
	}{
		{
			name:     "nil data",
			nthData:  nil,
			expected: "",
		},
		{
			name:     "A=0, B=5 (just number)",
			nthData:  NewSelectorPseudoNthData(0, 5),
			expected: "5",
		},
		{
			name:     "A=1, B=0 (just n)",
			nthData:  NewSelectorPseudoNthData(1, 0),
			expected: "n",
		},
		{
			name:     "A=1, B=3 (n+3)",
			nthData:  NewSelectorPseudoNthData(1, 3),
			expected: "n+3",
		},
		{
			name:     "A=1, B=-2 (n-2)",
			nthData:  NewSelectorPseudoNthData(1, -2),
			expected: "n-2",
		},
		{
			name:     "A=2, B=0 (2n)",
			nthData:  NewSelectorPseudoNthData(2, 0),
			expected: "2n",
		},
		{
			name:     "A=2, B=1 (2n+1)",
			nthData:  NewSelectorPseudoNthData(2, 1),
			expected: "2n+1",
		},
		{
			name:     "A=3, B=-4 (3n-4)",
			nthData:  NewSelectorPseudoNthData(3, -4),
			expected: "3n-4",
		},
		{
			name:     "A=-1, B=5 (-1n+5)",
			nthData:  NewSelectorPseudoNthData(-1, 5),
			expected: "-1n+5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nthData.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorPseudoNthDataStringWithSelectorList(t *testing.T) {
	selector := &Selector{
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("item"),
			},
		},
	}

	nthData := &SelectorPseudoNthData{
		A:            2,
		B:            1,
		SelectorList: []*Selector{selector},
	}

	expected := "2n+1 of .item"
	result := nthData.String()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSelectorPseudoNthDataEquals(t *testing.T) {
	nth1 := NewSelectorPseudoNthData(2, 1)
	nth2 := NewSelectorPseudoNthData(2, 1)
	nth3 := NewSelectorPseudoNthData(3, 1)
	nth4 := NewSelectorPseudoNthData(2, 0)

	selector1 := &Selector{
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	nth5 := &SelectorPseudoNthData{
		A:            2,
		B:            1,
		SelectorList: []*Selector{selector1},
	}

	nth6 := &SelectorPseudoNthData{
		A:            2,
		B:            1,
		SelectorList: []*Selector{selector1},
	}

	nth7 := &SelectorPseudoNthData{
		A:            2,
		B:            1,
		SelectorList: []*Selector{},
	}

	tests := []struct {
		name     string
		nth1     *SelectorPseudoNthData
		nth2     *SelectorPseudoNthData
		expected bool
	}{
		{"identical basic", nth1, nth2, true},
		{"same object", nth1, nth1, true},
		{"different A", nth1, nth3, false},
		{"different B", nth1, nth4, false},
		{"both nil", nil, nil, true},
		{"one nil", nth1, nil, false},
		{"other nil", nil, nth1, false},
		{"identical with selectors", nth5, nth6, true},
		{"different selector lists", nth5, nth7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nth1.Equals(tt.nth2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewSelectorPseudoNthData(t *testing.T) {
	a := 2
	b := 1

	nthData := NewSelectorPseudoNthData(a, b)

	if nthData == nil {
		t.Error("expected NewSelectorPseudoNthData to return non-nil")
	}
	if nthData.A != a {
		t.Errorf("expected A %d, got %d", a, nthData.A)
	}
	if nthData.B != b {
		t.Errorf("expected B %d, got %d", b, nthData.B)
	}
	if len(nthData.SelectorList) != 0 {
		t.Errorf("expected empty SelectorList, got %d items", len(nthData.SelectorList))
	}
}
