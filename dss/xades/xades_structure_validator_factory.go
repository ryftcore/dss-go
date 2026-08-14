// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESStructureValidatorFactory.java
// (DSS 6.5.RC1).
//
// Upstream's assertXAdESStructureValidatorLoaded() uses Class.forName to detect, at runtime,
// whether the optional 'dss-validation' module (which upstream bundles XAdESStructureValidator
// in) is present on the classpath, raising ExceptionInInitializerError otherwise. Go has no
// separate-module/classpath distinction - XAdESStructureValidator is compiled into this same
// package unconditionally - so that runtime presence check has nothing left to check and is
// dropped; XAdESStructureValidatorFactory reduces to the plain singleton + factory-method shape.
package xades

// XAdESStructureValidatorFactory creates a relevant implementation of XAdESStructureValidator.
// Port of the class XAdESStructureValidatorFactory.
type XAdESStructureValidatorFactory struct{}

// xadesStructureValidatorFactorySingleton is the current factory instance. Port of the private
// static singleton field.
var xadesStructureValidatorFactorySingleton *XAdESStructureValidatorFactory

// XAdESStructureValidatorFactoryGetInstance gets the instance of XAdESStructureValidatorFactory.
// Port of the static getInstance().
func XAdESStructureValidatorFactoryGetInstance() *XAdESStructureValidatorFactory {
	if xadesStructureValidatorFactorySingleton == nil {
		xadesStructureValidatorFactorySingleton = &XAdESStructureValidatorFactory{}
	}
	return xadesStructureValidatorFactorySingleton
}

// FromXAdESSignature creates a XAdESStructureValidator for the given XAdESSignature. Port of
// fromXAdESSignature(XAdESSignature).
func (f *XAdESStructureValidatorFactory) FromXAdESSignature(signature *XAdESSignature) *XAdESStructureValidator {
	return newXAdESStructureValidator(signature.SignatureElement(), signature.XAdESPaths())
}
