package css_parser

type CssStylesheet struct {
	Rules []*CssRule
}

func NewCssStylesheet() *CssStylesheet {
	return &CssStylesheet{}
}

func (s *CssStylesheet) String() string {
	result := ""

	for _, rule := range s.Rules {
		if result != "" {
			result += "\n"
		}
		result += rule.String()
	}

	return result
}
