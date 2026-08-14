package xpath10

import "strings"

// parser turns a token stream into the AST of ast.go.
//
// It is a recursive-descent parser over the XPath 1.0 grammar in which every production this
// package does not implement is replaced by an *UnsupportedError naming it. That is the point
// of the package: an expression is either evaluated exactly as Xalan would, or refused by
// name. Nothing is approximated and nothing is silently ignored.
//
// The grammar actually accepted, which is exactly what testdata/expressions.txt needs:
//
//	Expr      := OrExpr
//	OrExpr    := EqExpr ('or' EqExpr)*
//	EqExpr    := Operand ('=' Operand)*
//	Operand   := Literal | 'not' '(' Expr ')' | 'local-name' '(' ')' | Path
//	Path      := ('/' | '//')? Step (('/' | '//') Step)*
//	Step      := '.' | ('@' | ('child'|'parent'|'self') '::')? NodeTest Predicate*
//	NodeTest  := QName | '*' | 'node' '(' ')' | 'text' '(' ')'
//	Predicate := '[' Expr ']'
type parser struct {
	expr string
	toks []token
	i    int
	ns   NamespaceContext

	// transform widens the grammar to the XML-DSig transform subset; see transform.go. It is
	// false for Compile, which keeps the inventory subset exactly as documented.
	transform bool
}

// parse compiles expression to an AST, resolving prefixes against ns.
//
// The top level must be a location path. Every expression DSS compiles is one - it is always
// evaluated as a NODESET by JavaXmlXPathQueryExecutor - so an expression whose value is a
// boolean or a string is refused rather than given an evaluation nothing has verified.
func parse(expr string, ns NamespaceContext) (*pathExpr, error) {
	toks, err := lex(expr)
	if err != nil {
		return nil, err
	}
	p := &parser{expr: expr, toks: toks, ns: ns}
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tokEOF {
		return nil, p.syntax("unexpected trailing input")
	}
	path, ok := n.(*pathExpr)
	if !ok {
		return nil, &UnsupportedError{expr, 0, "a top-level expression that is not a location path"}
	}
	return path, nil
}

// ------------------------------------------------------------------ token helpers

func (p *parser) peek() token       { return p.toks[p.i] }
func (p *parser) next() token       { t := p.toks[p.i]; p.i++; return t }
func (p *parser) at(k tokKind) bool { return p.toks[p.i].kind == k }

func (p *parser) syntax(msg string) error {
	return &SyntaxError{Expression: p.expr, Pos: p.peek().pos, Msg: msg}
}

func (p *parser) unsupported(construct string) error {
	return &UnsupportedError{Expression: p.expr, Pos: p.peek().pos, Construct: construct}
}

func (p *parser) unsupportedAt(pos int, construct string) error {
	return &UnsupportedError{Expression: p.expr, Pos: pos, Construct: construct}
}

// ------------------------------------------------------------------ expressions

func (p *parser) parseExpr() (node, error) { return p.parseOr() }

func (p *parser) parseOr() (node, error) {
	lhs, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	for p.at(tokOr) {
		p.next()
		rhs, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		lhs = &orExpr{lhs: lhs, rhs: rhs}
	}
	// An operand followed by an operator this package does not implement - "@URI != ''",
	// "a and b", "x | y" - has parsed fine up to here and would otherwise be reported as a
	// missing ']' or trailing input, which blames the wrong thing. parseOr is the top of the
	// expression grammar, so nothing above it can legitimately be an operator.
	if p.at(tokOperator) {
		return nil, p.unsupported(operatorName(p.peek().text))
	}
	return lhs, nil
}

// parseUnion parses "a | b". It sits between parseEquality and parseOperand so that "|" binds
// tighter than "=" and "or", which is XPath 1.0's precedence. Only the transform grammar
// reaches it; parseEquality calls parseOperand directly otherwise, so Compile still refuses
// "|" by name from parseOr's trailing-operator check.
func (p *parser) parseUnion() (node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	for p.at(tokOperator) && p.peek().text == "|" {
		p.next()
		rhs, err := p.parseOperand()
		if err != nil {
			return nil, err
		}
		lhs = &unionExpr{lhs: lhs, rhs: rhs}
	}
	return lhs, nil
}

func (p *parser) parseEquality() (node, error) {
	operand := p.parseOperand
	if p.transform {
		operand = p.parseUnion
	}
	lhs, err := operand()
	if err != nil {
		return nil, err
	}
	for p.at(tokEq) {
		p.next()
		rhs, err := operand()
		if err != nil {
			return nil, err
		}
		lhs = &eqExpr{lhs: lhs, rhs: rhs}
	}
	return lhs, nil
}

// parseOperand parses one operand of an equality test: a string literal, one of the two
// permitted function calls, or a location path. The XPath grammar also puts UnionExpr,
// FilterExpr and number literals here; each is named and refused.
func (p *parser) parseOperand() (node, error) {
	t := p.peek()
	switch t.kind {
	case tokOperator:
		return nil, p.unsupported(operatorName(t.text))
	case tokLiteral:
		p.next()
		return &literal{value: t.text}, nil
	case tokNumber:
		return nil, p.unsupported("number literal")
	case tokLParen:
		return nil, p.unsupported("parenthesized expression")
	case tokName:
		if p.toks[p.i+1].kind == tokLParen && !isNodeTypeName(t.text) {
			fn, err := p.parseFunctionCall()
			if err != nil {
				return nil, err
			}
			// XPath's FilterExpr '/' RelativeLocationPath: "id('x')/node()". Only the
			// transform grammar admits it; without it the "/" would be trailing input.
			if p.transform && (p.at(tokSlash) || p.at(tokDoubleSlash)) {
				return p.parseRelativeFrom(fn)
			}
			return fn, nil
		}
	}
	return p.parsePath()
}

// parseRelativeFrom continues a location path from an already-parsed filter expression.
func (p *parser) parseRelativeFrom(start node) (node, error) {
	path := &pathExpr{start: start}
	for {
		switch {
		case p.at(tokSlash):
			p.next()
		case p.at(tokDoubleSlash):
			p.next()
			path.steps = append(path.steps, descendantOrSelfStep())
		default:
			return path, nil
		}
		if !p.startsStep() {
			return nil, p.syntax("expected a step after '/'")
		}
		s, err := p.parseStep()
		if err != nil {
			return nil, err
		}
		path.steps = append(path.steps, s)
	}
}

// parseFunctionCall parses not(Expr) and local-name(). Those are the only two core functions
// the inventory uses; every other one is refused by name, so a caller who reaches for, say,
// contains() gets told what is missing rather than a wrong answer.
func (p *parser) parseFunctionCall() (node, error) {
	name := p.next()
	p.next() // '('

	switch name.text {
	case "not":
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.at(tokComma) {
			return nil, p.unsupportedAt(name.pos, "not() with more than one argument")
		}
		if !p.at(tokRParen) {
			return nil, p.syntax("expected ')' to close not(")
		}
		p.next()
		return &notCall{arg: arg}, nil

	case "here":
		// Xalan's here() extension, which xmlsec 3.0 lost when it dropped Xalan: upstream's
		// JDKXPathAPI registers no extension functions, so here() fails there too.
		return nil, p.unsupportedAt(name.pos, "the here() extension function, which Apache Santuario 3.0 no longer provides either")

	case "local-name":
		// The optional node-set argument is not implemented: DSS only ever writes the
		// no-argument form, which means "the context node".
		if !p.at(tokRParen) {
			return nil, p.unsupportedAt(name.pos, "local-name() with an argument")
		}
		p.next()
		return &localNameCall{}, nil
	}
	if p.transform {
		if n, handled, err := p.parseTransformFunction(name); handled {
			return n, err
		}
	}
	return nil, p.unsupportedAt(name.pos, "function "+name.text+"()")
}

// ------------------------------------------------------------------ location paths

func (p *parser) parsePath() (node, error) {
	path := &pathExpr{}
	switch {
	case p.at(tokSlash):
		p.next()
		path.absolute = true
		if !p.startsStep() {
			return path, nil // "/" alone selects the root node
		}
	case p.at(tokDoubleSlash):
		p.next()
		path.absolute = true
		path.steps = append(path.steps, descendantOrSelfStep())
		if !p.startsStep() {
			return nil, p.syntax("'//' must be followed by a step")
		}
	}
	for {
		s, err := p.parseStep()
		if err != nil {
			return nil, err
		}
		path.steps = append(path.steps, s)
		switch {
		case p.at(tokSlash):
			p.next()
		case p.at(tokDoubleSlash):
			p.next()
			path.steps = append(path.steps, descendantOrSelfStep())
		default:
			return path, nil
		}
		if !p.startsStep() {
			return nil, p.syntax("expected a step after '/'")
		}
	}
}

// descendantOrSelfStep is what "//" abbreviates: "/descendant-or-self::node()/", per XPath 1.0
// clause 2.5. It is built here and never parsed, which is why axisByName does not know the
// name "descendant-or-self".
func descendantOrSelfStep() step {
	return step{axis: axisDescendantOrSelf, test: nodeTest{kind: testNode}}
}

// startsStep reports whether the next token can begin a location step. ".." is included so
// parseStep can refuse it by name rather than leaving parsePath to report a syntax error.
func (p *parser) startsStep() bool {
	switch p.peek().kind {
	case tokName, tokStar, tokAt, tokDot, tokDotDot:
		return true
	}
	return false
}

func (p *parser) parseStep() (step, error) {
	switch {
	case p.at(tokDot):
		p.next()
		if p.at(tokLBracket) {
			return step{}, p.syntax("'.' cannot carry a predicate")
		}
		return step{axis: axisSelf, test: nodeTest{kind: testNode}}, nil
	case p.at(tokDotDot):
		return step{}, p.unsupported("the '..' parent abbreviation")
	}

	s := step{axis: axisChild}
	if p.at(tokAt) {
		p.next()
		s.axis = axisAttribute
	} else if p.at(tokName) && p.toks[p.i+1].kind == tokColonColon {
		name := p.next()
		p.next() // '::'
		a, ok := axisByName(name.text)
		if !ok && p.transform {
			a, ok = transformAxisByName(name.text)
		}
		if !ok {
			return step{}, p.unsupportedAt(name.pos, "the "+name.text+":: axis")
		}
		s.axis = a
	}

	test, err := p.parseNodeTest()
	if err != nil {
		return step{}, err
	}
	s.test = test

	for p.at(tokLBracket) {
		p.next()
		pred, err := p.parseExpr()
		if err != nil {
			return step{}, err
		}
		if !p.at(tokRBracket) {
			return step{}, p.syntax("expected ']'")
		}
		p.next()
		s.preds = append(s.preds, pred)
	}
	return s, nil
}

func (p *parser) parseNodeTest() (nodeTest, error) {
	if p.at(tokStar) {
		p.next()
		return nodeTest{kind: testAny}, nil
	}

	// A QName with an explicit-but-empty prefix (a stray leading colon, e.g. the "::" axis
	// separator immediately followed by another colon: "descendant:::Signature"). This is not
	// valid XPath 1.0 grammar - NCName forbids an empty production - but it is exactly what
	// XPath2FilterEnvelopedSignatureTransform's Java source generates when built with a
	// DSSNamespace whose prefix is "" (the ds elements use the default, unprefixed xmlns="u"
	// declaration instead of xmlns:ds="u"), and upstream's javax.xml.xpath evaluates it as a
	// name test in the default namespace currently in scope - not as a syntax error, and not
	// as the null namespace a genuinely unprefixed name test would get (that is a different
	// grammar production; see NamespaceContextOf's own doc comment on that distinction).
	// Reproduced only for CompileTransform's grammar (p.transform), since only the generated
	// transform text is ever malformed this way; a hand-written XPath elsewhere in the corpus
	// stays a syntax error. Pinned byte-exact against the Java oracle by
	// xades_reference_kat_test.go's "si-enveloped-default-prefix" case.
	if p.transform && p.at(tokColon) {
		p.next()
		if !p.at(tokName) {
			return nodeTest{}, p.syntax("expected a node test")
		}
		local := p.next()
		if p.at(tokColon) {
			return nodeTest{}, p.unsupportedAt(local.pos, "the 'prefix:*' node test")
		}
		return nodeTest{kind: testName, space: p.ns.NamespaceURI(defaultNamespaceKey), local: local.text}, nil
	}

	if !p.at(tokName) {
		return nodeTest{}, p.syntax("expected a node test")
	}
	name := p.next()

	// "prefix:*" arrives as three tokens because "*" is not an NCName. It is valid XPath and
	// this package does not implement it, so it is named rather than mistaken for a syntax
	// error at the "*".
	if p.at(tokColon) {
		return nodeTest{}, p.unsupportedAt(name.pos, "the 'prefix:*' node test")
	}

	if p.at(tokLParen) && isNodeTypeName(name.text) {
		p.next()
		var kind testKind
		switch name.text {
		case "node":
			kind = testNode
		case "text":
			kind = testText
		default:
			// comment() and processing-instruction(): valid, not implemented.
			return nodeTest{}, p.unsupportedAt(name.pos, "the "+name.text+"() node test")
		}
		if !p.at(tokRParen) {
			return nodeTest{}, p.syntax("expected ')' after " + name.text + "(")
		}
		p.next()
		return nodeTest{kind: kind}, nil
	}

	prefix, local := splitQName(name.text)
	return nodeTest{kind: testName, space: p.ns.NamespaceURI(prefix), local: local}, nil
}

func splitQName(name string) (prefix, local string) {
	if i := strings.IndexByte(name, ':'); i >= 0 {
		return name[:i], name[i+1:]
	}
	return "", name
}

// isNodeTypeName reports whether a name followed by "(" is a NodeType rather than a function
// call. All four are listed, including the two that are refused: "comment()" must be read as
// a node test so it can be refused as one, not as an unknown function.
func isNodeTypeName(s string) bool {
	switch s {
	case "node", "text", "comment", "processing-instruction":
		return true
	}
	return false
}

// axisByName maps an axis specifier written in the source. It deliberately knows only the
// three axes the inventory writes by name. "attribute::" is absent because the inventory only
// ever writes its "@" abbreviation, and "descendant-or-self::" is absent because it is only
// ever reached through "//"; both are refused with their own name, which is the honest
// message - the construct is real XPath that this package does not implement.
func axisByName(s string) (axis, bool) {
	switch s {
	case "child":
		return axisChild, true
	case "parent":
		return axisParent, true
	case "self":
		return axisSelf, true
	}
	return 0, false
}

// operatorName gives an *UnsupportedError a name for an operator token, so the message reads
// "unsupported XPath construct: the '|' union operator" rather than echoing a symbol.
func operatorName(text string) string {
	switch text {
	case "|":
		return "the '|' union operator"
	case "$":
		return "variable reference"
	case "and":
		return "the 'and' operator"
	case "!=":
		return "the '!=' operator"
	case "*", "+", "-", "div", "mod":
		return "the '" + text + "' arithmetic operator"
	case "<", "<=", ">", ">=":
		return "the '" + text + "' relational operator"
	}
	return "the '" + text + "' operator"
}
