package xpath10

import "fmt"

// SyntaxError reports an expression this package could not parse at all: a stray bracket, an
// unterminated literal, a name where an operator belongs. Pos is a 0-based byte offset into
// the expression.
//
// It is distinct from UnsupportedError on purpose. A SyntaxError says the input is not XPath;
// an UnsupportedError says the input is XPath this package deliberately does not implement.
type SyntaxError struct {
	Expression string
	Pos        int
	Msg        string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("xpath10: %s at offset %d in %q", e.Msg, e.Pos, e.Expression)
}

// UnsupportedError reports a well-formed XPath 1.0 construct outside the subset this package
// implements. Construct names the construct in the words of the XPath 1.0 grammar, so the
// message says what to add rather than merely that something failed.
//
// Every rejection is one of these: the engine never approximates a construct it does not
// implement, and never silently drops one.
type UnsupportedError struct {
	Expression string
	Pos        int
	Construct  string
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("xpath10: unsupported XPath construct: %s at offset %d in %q",
		e.Construct, e.Pos, e.Expression)
}

// EvalError reports a failure that only shows up with a context node in hand: a nil context
// node, or a result the caller's contract does not admit, as SelectOne reports for a query
// that matched more than one node.
type EvalError struct {
	Expression string
	Msg        string
}

func (e *EvalError) Error() string {
	return fmt.Sprintf("xpath10: %s in %q", e.Msg, e.Expression)
}
