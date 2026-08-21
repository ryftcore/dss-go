// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/scope/SignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureScope describes the scope of the signature. It is the polymorphic contract of the
// Java abstract class SignatureScope; the state and concrete method bodies live in
// SignatureScopeBase.
type SignatureScope interface {
	model.IdentifierBasedObject

	// DocumentName gets the name of the document. Port of getDocumentName().
	DocumentName() string
	// Name returns a signature scope name, extracting a token identifier through
	// tokenIdentifierProvider when required. Port of getName(TokenIdentifierProvider).
	Name(tokenIdentifierProvider model.TokenIdentifierProvider) string
	// Digest gets the digest of the document for digestAlgorithm, or the zero model.Digest
	// when the scope carries no document (Java's null). Port of getDigest(DigestAlgorithm).
	//
	// Java's getDigest never declares a checked exception (DSSException is unchecked); the
	// error return surfaces the same data-dependent failure the Go DSSDocument methods can
	// raise (e.g. a failing OpenStream) that Java would propagate as an unchecked exception.
	Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error)
	// Description gets the signature scope description, extracting a token identifier
	// through tokenIdentifierProvider when required. Port of the abstract
	// getDescription(TokenIdentifierProvider).
	Description(tokenIdentifierProvider model.TokenIdentifierProvider) string
	// Transformations returns a list of transformations on the original document when
	// applicable, nil by default. Port of getTransformations().
	Transformations() []string
	// Type returns the type of the signature scope. Port of the abstract getType().
	Type() enumerations.SignatureScopeType
	// Children returns a list of dependent signature scopes (e.g. Manifest entries). Port of
	// getChildren().
	Children() []SignatureScope
	// AddChildSignatureScope adds a new child SignatureScope. Port of
	// addChildSignatureScope(SignatureScope).
	AddChildSignatureScope(child SignatureScope)
	// DSSIDAsString returns a String representation of the DSS Identifier. Port of
	// getDSSIdAsString().
	DSSIDAsString() string
	// String returns the Java toString() form. Port of toString().
	String() string
	// Equals reports whether both signature scopes carry an equal DSS Identifier. Port of
	// equals(Object).
	Equals(other SignatureScope) bool
}

// SignatureScopeBase carries the state and the concrete behaviour of the Java abstract class
// SignatureScope. Concrete signature scopes embed it; unlike TokenBase, no method on the base
// needs to call back into the concrete type (Description and Type have no base logic that
// depends on them), so - unlike Token/TokenBase - there is no overrides-registration step:
// the concrete signature scope simply implements Description and Type directly to satisfy
// the SignatureScope interface, exactly the way Java's abstract methods must be implemented
// by a subclass.
type SignatureScopeBase struct {
	// name is the name of the item on which this signature scope applies.
	name string
	// document is the original signer data.
	document model.DSSDocument

	// dssID caches the default DSS Identifier.
	dssID *model.DataIdentifier
	// dssIDBuilt records whether dssID has been computed, distinguishing "not yet built" from
	// a build that legitimately failed and must be retried.
	dssIDBuilt bool

	// children is a list of dependent signature scopes (e.g. Manifest entries).
	children []SignatureScope
}

// NewSignatureScopeBase is the default constructor. Port of the protected
// SignatureScope(DSSDocument) constructor.
func NewSignatureScopeBase(document model.DSSDocument) SignatureScopeBase {
	return NewSignatureScopeBaseWithName(document.Name(), document)
}

// NewSignatureScopeBaseWithName is the default constructor with the name provided. Port of
// the protected SignatureScope(String, DSSDocument) constructor.
func NewSignatureScopeBaseWithName(name string, document model.DSSDocument) SignatureScopeBase {
	return SignatureScopeBase{name: name, document: document}
}

// DocumentName gets name of the document. Port of getDocumentName().
func (s *SignatureScopeBase) DocumentName() string {
	return s.name
}

// Name returns a signature scope name. The base implementation ignores
// tokenIdentifierProvider and returns DocumentName(); subclasses needing the provider
// override this behaviour on the embedding concrete type. Port of
// getName(TokenIdentifierProvider).
func (s *SignatureScopeBase) Name(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.DocumentName()
}

// Digest gets digests of the document. Port of getDigest(DigestAlgorithm).
func (s *SignatureScopeBase) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	if s.document == nil {
		return model.Digest{}, nil
	}
	if digestDocument, ok := s.document.(*model.DigestDocument); ok {
		return digestDocument.ExistingDigest()
	}
	value, err := s.document.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, value), nil
}

// Transformations returns a list of transformations on the original document when
// applicable; not implemented by default. Port of getTransformations().
func (s *SignatureScopeBase) Transformations() []string {
	return nil
}

// Children returns a list of dependent signature scopes (e.g. Manifest entries), lazily
// initialising the backing slice as Java lazily initialises its ArrayList. Port of
// getChildren().
func (s *SignatureScopeBase) Children() []SignatureScope {
	if s.children == nil {
		s.children = make([]SignatureScope, 0)
	}
	return s.children
}

// AddChildSignatureScope adds a new child SignatureScope. Port of
// addChildSignatureScope(SignatureScope).
func (s *SignatureScopeBase) AddChildSignatureScope(child SignatureScope) {
	s.children = append(s.Children(), child)
}

// DSSID returns the unique DSS Identifier, building it lazily. Port of getDSSId().
//
// Java's getDSSId() declares no checked exception (DSSException is unchecked); building the
// identifier can fail only when the wrapped document's content cannot be read, which Java
// would surface as an unchecked exception too, so a genuine failure here panics with a
// model.DSSError instead of silently caching a broken value.
func (s *SignatureScopeBase) DSSID() model.Identifier {
	if !s.dssIDBuilt {
		id, err := model.NewDataIdentifierForDocument(s.name, s.document)
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		s.dssID = id
		s.dssIDBuilt = true
	}
	return s.dssID
}

// DSSIDAsString returns a String representation of the DSS Identifier. Port of
// getDSSIdAsString().
func (s *SignatureScopeBase) DSSIDAsString() string {
	return s.DSSID().AsXmlID()
}

// String returns the Java toString() form of SignatureScope, listing name, document, dssId
// and children in declaration order.
func (s *SignatureScopeBase) String() string {
	var dssID any
	if s.dssIDBuilt {
		dssID = s.dssID
	}
	return fmt.Sprintf("SignatureScope{name='%s', document=%v, dssId=%v, children=%v}", s.name, s.document, dssID, s.children)
}

// Equals reports whether both signature scopes carry an equal DSS Identifier. Port of
// equals(Object); unlike model.Identifier#Equals, Java's SignatureScope#equals accepts any
// SignatureScope subclass (no getClass() check), so the comparison here is on the DSS
// Identifier value alone. Java's "this == obj" fast path has no useful Go counterpart here
// (s is the embedded base, not the enclosing concrete value), so it is omitted; the
// DSSID().Equals comparison below already reports true for that case.
func (s *SignatureScopeBase) Equals(other SignatureScope) bool {
	if other == nil {
		return false
	}
	return s.DSSID().Equals(other.DSSID())
}

// compile-time interface assertion.
var _ model.IdentifierBasedObject = (*SignatureScopeBase)(nil)
