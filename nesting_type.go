package cssparser

type NestingTypeType int

const (
	// We are not in a nesting context, and '&' resolves like :scope instead.
	NestingTypeNone NestingTypeType = iota

	// We are in a nesting context as defined by @scope.
	//
	// https://www.w3.org/TR/2024/WD-css-cascade-6-20240906/#scope-atrule
	// https://www.w3.org/TR/2022/WD-selectors-4-20221111/#scope-pseudo
	NestingTypeScope // Scoped nesting, e.g., `:scope`

	// We are in a css-nesting nesting context, and '&' resolves according to:
	// https://www.w3.org/TR/2023/WD-css-nesting-1-20230214/#nest-selector
	NestingTypeNesting

	// We are inside @function. The parsing behavior is generally the same as
	// kNesting, except we don't allow qualified rules, and we emit
	// CSSFunctionDeclarations instead of CSSNestedDeclarations.
	NestingTypeFunction
)
