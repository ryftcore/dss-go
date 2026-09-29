// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESStructureValidatorFactory.java
// (DSS 6.5.RC1).
//
// Upstream's assertXAdESStructureValidatorLoaded() uses Class.forName to detect, at runtime,
// whether the optional 'dss-validation' module (which upstream bundles XAdESStructureValidator
// in) is present on the classpath, raising ExceptionInInitializerError otherwise. Go has no
// separate-module/classpath distinction - StructureValidator is compiled into this same
// package unconditionally - so that runtime presence check has nothing left to check and is
// dropped; StructureValidatorFactory reduces to the plain singleton + factory-method shape.
package xades

import "sync"

// StructureValidatorFactory creates a relevant implementation of StructureValidator.
// Port of the class XAdESStructureValidatorFactory.
type StructureValidatorFactory struct{}

// xadesStructureValidatorFactorySingleton is the current factory instance. Port of the private
// static singleton field. Initialised under a sync.Once: upstream's unsynchronized
// check-then-set is harmless in Java, but it is a data race in Go under concurrent validations.
var (
	xadesStructureValidatorFactorySingleton *StructureValidatorFactory
	xadesStructureValidatorFactoryOnce      sync.Once
)

// StructureValidatorFactoryGetInstance gets the instance of StructureValidatorFactory.
// Port of the static getInstance().
func StructureValidatorFactoryGetInstance() *StructureValidatorFactory {
	xadesStructureValidatorFactoryOnce.Do(func() {
		xadesStructureValidatorFactorySingleton = &StructureValidatorFactory{}
	})
	return xadesStructureValidatorFactorySingleton
}

// FromXAdESSignature creates a XAdESStructureValidator for the given XAdESSignature. Port of
// fromXAdESSignature(Signature).
func (f *StructureValidatorFactory) FromXAdESSignature(signature *Signature) *StructureValidator {
	return newStructureValidator(signature.SignatureElement(), signature.XAdESPaths())
}
