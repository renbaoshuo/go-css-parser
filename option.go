package cssparser

type ParserOption func(*Parser)

func WithInline(inline bool) ParserOption {
	return func(p *Parser) {
		p.inline = inline
	}
}

func WithLooseParsing(loose bool) ParserOption {
	return func(p *Parser) {
		p.loose = loose
	}
}
