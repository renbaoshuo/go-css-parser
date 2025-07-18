package cssparser

type Stylesheet struct {
	Rules []*CssRule
}

func NewStylesheet() *Stylesheet {
	return &Stylesheet{}
}

func (s *Stylesheet) String() string {
	result := ""

	for _, rule := range s.Rules {
		if result != "" {
			result += "\n"
		}
		result += rule.String()
	}

	return result
}
