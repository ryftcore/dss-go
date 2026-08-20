// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/SchemeInformationURIByLang.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/trustedlist/jaxb"

// SchemeInformationURIByLang filters scheme information by language.
type SchemeInformationURIByLang struct {
	// lang is the language code to filter by.
	lang string
}

// NewSchemeInformationURIByLang is the default constructor. Port of
// SchemeInformationURIByLang(String).
//
// Panics with the Java message when lang is empty (Objects.requireNonNull).
func NewSchemeInformationURIByLang(lang string) *SchemeInformationURIByLang {
	if lang == "" {
		panic("lang")
	}
	return &SchemeInformationURIByLang{lang: lang}
}

// Test ports test(NonEmptyMultiLangURIType).
func (p *SchemeInformationURIByLang) Test(schemeInformationURI *jaxb.NonEmptyMultiLangURIType) bool {
	return p.lang == schemeInformationURI.Lang
}
