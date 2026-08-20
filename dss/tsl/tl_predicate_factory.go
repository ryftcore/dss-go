// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TLPredicateFactory.java (DSS 6.5.RC1).
//
// Java's private constructor (a non-instantiable utility class) has no Go counterpart; the port
// is simply a set of package-level functions.
package tsl

// TLPredicateFactoryCreateEULOTLPredicate creates a Predicate used to filter the XML European
// list of trusted list (LOTL). Port of createEULOTLPredicate().
func TLPredicateFactoryCreateEULOTLPredicate() OtherTSLPointerPredicate {
	return otherTSLPointerPredicateAndOf(NewEULOTLOtherTSLPointer(), NewXMLOtherTSLPointer())
}

// TLPredicateFactoryCreateEUTLPredicate creates a Predicate used to filter the XML European
// Trusted List (TL). Port of createEUTLPredicate().
func TLPredicateFactoryCreateEUTLPredicate() OtherTSLPointerPredicate {
	return otherTSLPointerPredicateAndOf(NewEUTLOtherTSLPointer(), NewXMLOtherTSLPointer())
}

// TLPredicateFactoryCreatePredicateWithCustomTSLType creates a Predicate used to filter an XML
// Trusted List (TL) defined with a custom TSLType. Port of createPredicateWithCustomTSLType(String).
func TLPredicateFactoryCreatePredicateWithCustomTSLType(tslType string) OtherTSLPointerPredicate {
	return otherTSLPointerPredicateAndOf(NewTypeOtherTSLPointer(tslType), NewXMLOtherTSLPointer())
}

// TLPredicateFactoryCreateEUTLCountryCodePredicate creates a Predicate used to filter XML
// European Trusted Lists (TL) with the defined Scheme Territory codes. Port of
// createEUTLCountryCodePredicate(String...).
func TLPredicateFactoryCreateEUTLCountryCodePredicate(countryCodes ...string) OtherTSLPointerPredicate {
	return otherTSLPointerPredicateAndOf(
		otherTSLPointerPredicateAndOf(NewSchemeTerritoryOtherTSLPointerCollection(countryCodes), NewEUTLOtherTSLPointer()),
		NewXMLOtherTSLPointer())
}

// TLPredicateFactoryCreateXMLOtherTSLPointerPredicate creates a predicate used to filter all XML
// Trusted Lists (TL). Port of createXMLOtherTSLPointerPredicate().
func TLPredicateFactoryCreateXMLOtherTSLPointerPredicate() OtherTSLPointerPredicate {
	return NewXMLOtherTSLPointer()
}

// TLPredicateFactoryCreatePDFOtherTSLPointerPredicate creates a predicate used to filter all PDF
// Trusted Lists (TL). Port of createPDFOtherTSLPointerPredicate().
func TLPredicateFactoryCreatePDFOtherTSLPointerPredicate() OtherTSLPointerPredicate {
	return NewPDFOtherTSLPointer()
}

// TLPredicateFactoryCreatePredicateWithCustomMimeType creates a predicate used to filter all
// Trusted Lists (TL) with the defined mimetype. Port of createPredicateWithCustomMimeType(String).
func TLPredicateFactoryCreatePredicateWithCustomMimeType(mimetype string) OtherTSLPointerPredicate {
	return NewMimetypeOtherTSLPointer(mimetype)
}
