// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/SignatureScopeWithTransformations.java (DSS 6.5.RC1).
package scope

import (
	"fmt"

	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/utils"
)

// SignatureScopeWithTransformationsBase is the signature scope with the performed
// transforms. It is the Go counterpart of the Java abstract class
// SignatureScopeWithTransformations: like SignatureScopeBase itself, it implements no
// Description()/Type() (those stay abstract, per the embedding concrete type - e.g.
// ManifestSignatureScope, outside this manifest).
type SignatureScopeWithTransformationsBase struct {
	mscope.SignatureScopeBase

	// transformations is the list of transform definitions.
	transformations []string
}

// NewSignatureScopeWithTransformationsBase is the default constructor with a name extracted
// from a document. Port of the protected SignatureScopeWithTransformations(DSSDocument,
// List<String>) constructor.
func NewSignatureScopeWithTransformationsBase(document model.DSSDocument, transformations []string) SignatureScopeWithTransformationsBase {
	return NewSignatureScopeWithTransformationsBaseWithName(document.Name(), document, transformations)
}

// NewSignatureScopeWithTransformationsBaseWithName is the default constructor with a name
// provided. Port of the protected SignatureScopeWithTransformations(String, DSSDocument,
// List<String>) constructor.
func NewSignatureScopeWithTransformationsBaseWithName(name string, document model.DSSDocument, transformations []string) SignatureScopeWithTransformationsBase {
	return SignatureScopeWithTransformationsBase{
		SignatureScopeBase: mscope.NewSignatureScopeBaseWithName(name, document),
		transformations:    transformations,
	}
}

// AddTransformationIfNeeded adds a description to the signature scope if needed. Port of
// addTransformationIfNeeded(String).
func (s *SignatureScopeWithTransformationsBase) AddTransformationIfNeeded(description string) string {
	if utils.IsCollectionNotEmpty(s.transformations) {
		description += " with transformations."
	}
	return description
}

// Transformations returns the list of transformations. Port of getTransformations().
func (s *SignatureScopeWithTransformationsBase) Transformations() []string {
	return s.transformations
}

// String returns the Java toString() form of SignatureScopeWithTransformations. Port of
// toString().
func (s *SignatureScopeWithTransformationsBase) String() string {
	return fmt.Sprintf("SignatureScopeWithTransformations{transformations=%v} %s", s.transformations, s.SignatureScopeBase.String())
}
