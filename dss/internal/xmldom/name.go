package xmldom

const (
	// XMLNamespace is the namespace name permanently bound to the "xml" prefix.
	XMLNamespace = "http://www.w3.org/XML/1998/namespace" // the "xml" prefix
	// XMLNSNamespace is the namespace name of namespace declaration attributes.
	XMLNSNamespace = "http://www.w3.org/2000/xmlns/" // namespace declarations
)

// Name is an expanded name plus the literal prefix as it appeared in the source.
// Canonicalization emits QNames verbatim, so Prefix is load-bearing and must never
// be regenerated. Space is "" for a name in no namespace.
//
// Namespace declarations are modelled as attributes, exactly as org.w3c.dom does:
//
//	xmlns:p="u"  ->  {XMLNSNamespace, "p",     "xmlns"}   QName "xmlns:p"
//	xmlns="u"    ->  {XMLNSNamespace, "xmlns", ""}        QName "xmlns"
//	p:a="v"      ->  {resolved URI,   "a",     "p"}       QName "p:a"
//	a="v"        ->  {"",             "a",     ""}        QName "a"
//
// Java reports namespaceURI == null for an unprefixed attribute where we use "".
// Santuario's AttrCompare sorts null before any non-null URI and "" is
// lexicographically least among all strings, so the two are equivalent and no
// separate "has namespace" flag is needed.
type Name struct {
	Space  string
	Local  string
	Prefix string
}

// QName returns "prefix:local", or "local" when Prefix is empty.
func (n Name) QName() string {
	if n.Prefix == "" {
		return n.Local
	}
	return n.Prefix + ":" + n.Local
}
