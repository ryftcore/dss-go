// Ported from dss-enumerations/.../SignatureLevel.java (DSS 6.5.RC1).
//
// NOTE: SignatureForm and SignatureProfile (with their exported constants,
// e.g. SignatureFormCAdES, SignatureProfileBaselineB, preserving the
// exact Java constant spelling/casing) are defined outside this file's
// manifest and are assumed to exist per the porting brief.
package enumerations

import (
	"fmt"
	"strings"
)

// SignatureLevel represents signature profiles (form+level) handled by the
// SD-DSS framework.
type SignatureLevel string

// SignatureLevel constants, grouped by underlying signature format (XAdES,
// CAdES, PAdES/PKCS7, JAdES, CB_AdES/CBOR). Baseline profiles
// (BASELINE_B/T/LT/LTA) are the ETSI EN 319 122/132/142 recommended
// profiles; the non-baseline levels (e.g. XAdES_C, XAdES_X) are the legacy
// ETSI TS 101 903/CAdES-equivalent extended forms. The _NOT_ETSI values mark
// a signature container whose content is not recognized as any supported
// AdES form, and SignatureLevelUnknown marks a level DSS could not
// determine.
const (
	SignatureLevelXMLNotETSI       SignatureLevel = "XML_NOT_ETSI"
	SignatureLevelXAdESBES         SignatureLevel = "XAdES_BES"
	SignatureLevelXAdESEPES        SignatureLevel = "XAdES_EPES"
	SignatureLevelXAdEST           SignatureLevel = "XAdES_T"
	SignatureLevelXAdESLT          SignatureLevel = "XAdES_LT"
	SignatureLevelXAdESC           SignatureLevel = "XAdES_C"
	SignatureLevelXAdESX           SignatureLevel = "XAdES_X"
	SignatureLevelXAdESXL          SignatureLevel = "XAdES_XL"
	SignatureLevelXAdESA           SignatureLevel = "XAdES_A"
	SignatureLevelXAdESERS         SignatureLevel = "XAdES_ERS"
	SignatureLevelXAdESBaselineB   SignatureLevel = "XAdES_BASELINE_B"
	SignatureLevelXAdESBaselineT   SignatureLevel = "XAdES_BASELINE_T"
	SignatureLevelXAdESBaselineLT  SignatureLevel = "XAdES_BASELINE_LT"
	SignatureLevelXAdESBaselineLTA SignatureLevel = "XAdES_BASELINE_LTA"

	SignatureLevelCMSNotETSI       SignatureLevel = "CMS_NOT_ETSI"
	SignatureLevelCAdESBES         SignatureLevel = "CAdES_BES"
	SignatureLevelCAdESEPES        SignatureLevel = "CAdES_EPES"
	SignatureLevelCAdEST           SignatureLevel = "CAdES_T"
	SignatureLevelCAdESLT          SignatureLevel = "CAdES_LT"
	SignatureLevelCAdESC           SignatureLevel = "CAdES_C"
	SignatureLevelCAdESX           SignatureLevel = "CAdES_X"
	SignatureLevelCAdESXL          SignatureLevel = "CAdES_XL"
	SignatureLevelCAdESA           SignatureLevel = "CAdES_A"
	SignatureLevelCAdESERS         SignatureLevel = "CAdES_ERS"
	SignatureLevelCAdESBaselineB   SignatureLevel = "CAdES_BASELINE_B"
	SignatureLevelCAdESBaselineT   SignatureLevel = "CAdES_BASELINE_T"
	SignatureLevelCAdESBaselineLT  SignatureLevel = "CAdES_BASELINE_LT"
	SignatureLevelCAdESBaselineLTA SignatureLevel = "CAdES_BASELINE_LTA"

	SignatureLevelPDFNotETSI       SignatureLevel = "PDF_NOT_ETSI"
	SignatureLevelPKCS7B           SignatureLevel = "PKCS7_B"
	SignatureLevelPKCS7T           SignatureLevel = "PKCS7_T"
	SignatureLevelPKCS7LT          SignatureLevel = "PKCS7_LT"
	SignatureLevelPKCS7LTA         SignatureLevel = "PKCS7_LTA"
	SignatureLevelPAdESBES         SignatureLevel = "PAdES_BES"
	SignatureLevelPAdESEPES        SignatureLevel = "PAdES_EPES"
	SignatureLevelPAdESLTV         SignatureLevel = "PAdES_LTV"
	SignatureLevelPAdESBaselineB   SignatureLevel = "PAdES_BASELINE_B"
	SignatureLevelPAdESBaselineT   SignatureLevel = "PAdES_BASELINE_T"
	SignatureLevelPAdESBaselineLT  SignatureLevel = "PAdES_BASELINE_LT"
	SignatureLevelPAdESBaselineLTA SignatureLevel = "PAdES_BASELINE_LTA"

	SignatureLevelJSONNotETSI      SignatureLevel = "JSON_NOT_ETSI"
	SignatureLevelJAdES            SignatureLevel = "JAdES"
	SignatureLevelJAdESBaselineB   SignatureLevel = "JAdES_BASELINE_B"
	SignatureLevelJAdESBaselineT   SignatureLevel = "JAdES_BASELINE_T"
	SignatureLevelJAdESBaselineLT  SignatureLevel = "JAdES_BASELINE_LT"
	SignatureLevelJAdESBaselineLTA SignatureLevel = "JAdES_BASELINE_LTA"

	SignatureLevelCBORNotETSI       SignatureLevel = "CBOR_NOT_ETSI"
	SignatureLevelCBAdES            SignatureLevel = "CB_AdES"
	SignatureLevelCBAdESBaselineB   SignatureLevel = "CB_AdES_BASELINE_B"
	SignatureLevelCBAdESBaselineT   SignatureLevel = "CB_AdES_BASELINE_T"
	SignatureLevelCBAdESBaselineLT  SignatureLevel = "CB_AdES_BASELINE_LT"
	SignatureLevelCBAdESBaselineLTA SignatureLevel = "CB_AdES_BASELINE_LTA"

	SignatureLevelUnknown SignatureLevel = "UNKNOWN"
)

// signatureLevelFields holds the (signatureForm, signatureProfile) pair for
// each constant. A zero-value ("") signatureForm mirrors Java's null
// signatureForm on the UNKNOWN constant.
type signatureLevelFields struct {
	signatureForm    SignatureForm
	signatureProfile SignatureProfile
}

// signatureLevelData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var signatureLevelData = map[SignatureLevel]signatureLevelFields{
	SignatureLevelXMLNotETSI:       {SignatureFormXAdES, SignatureProfileNotETSI},
	SignatureLevelXAdESBES:         {SignatureFormXAdES, SignatureProfileExtendedBES},
	SignatureLevelXAdESEPES:        {SignatureFormXAdES, SignatureProfileExtendedEPES},
	SignatureLevelXAdEST:           {SignatureFormXAdES, SignatureProfileExtendedT},
	SignatureLevelXAdESLT:          {SignatureFormXAdES, SignatureProfileExtendedLT},
	SignatureLevelXAdESC:           {SignatureFormXAdES, SignatureProfileExtendedC},
	SignatureLevelXAdESX:           {SignatureFormXAdES, SignatureProfileExtendedX},
	SignatureLevelXAdESXL:          {SignatureFormXAdES, SignatureProfileExtendedXL},
	SignatureLevelXAdESA:           {SignatureFormXAdES, SignatureProfileExtendedA},
	SignatureLevelXAdESERS:         {SignatureFormXAdES, SignatureProfileExtendedERS},
	SignatureLevelXAdESBaselineB:   {SignatureFormXAdES, SignatureProfileBaselineB},
	SignatureLevelXAdESBaselineT:   {SignatureFormXAdES, SignatureProfileBaselineT},
	SignatureLevelXAdESBaselineLT:  {SignatureFormXAdES, SignatureProfileBaselineLT},
	SignatureLevelXAdESBaselineLTA: {SignatureFormXAdES, SignatureProfileBaselineLTA},

	SignatureLevelCMSNotETSI:       {SignatureFormCAdES, SignatureProfileNotETSI},
	SignatureLevelCAdESBES:         {SignatureFormCAdES, SignatureProfileExtendedBES},
	SignatureLevelCAdESEPES:        {SignatureFormCAdES, SignatureProfileExtendedEPES},
	SignatureLevelCAdEST:           {SignatureFormCAdES, SignatureProfileExtendedT},
	SignatureLevelCAdESLT:          {SignatureFormCAdES, SignatureProfileExtendedLT},
	SignatureLevelCAdESC:           {SignatureFormCAdES, SignatureProfileExtendedC},
	SignatureLevelCAdESX:           {SignatureFormCAdES, SignatureProfileExtendedX},
	SignatureLevelCAdESXL:          {SignatureFormCAdES, SignatureProfileExtendedXL},
	SignatureLevelCAdESA:           {SignatureFormCAdES, SignatureProfileExtendedA},
	SignatureLevelCAdESERS:         {SignatureFormCAdES, SignatureProfileExtendedERS},
	SignatureLevelCAdESBaselineB:   {SignatureFormCAdES, SignatureProfileBaselineB},
	SignatureLevelCAdESBaselineT:   {SignatureFormCAdES, SignatureProfileBaselineT},
	SignatureLevelCAdESBaselineLT:  {SignatureFormCAdES, SignatureProfileBaselineLT},
	SignatureLevelCAdESBaselineLTA: {SignatureFormCAdES, SignatureProfileBaselineLTA},

	SignatureLevelPDFNotETSI:       {SignatureFormPAdES, SignatureProfileNotETSI},
	SignatureLevelPKCS7B:           {SignatureFormPKCS7, SignatureProfileNotETSI},
	SignatureLevelPKCS7T:           {SignatureFormPKCS7, SignatureProfileNotETSI},
	SignatureLevelPKCS7LT:          {SignatureFormPKCS7, SignatureProfileNotETSI},
	SignatureLevelPKCS7LTA:         {SignatureFormPKCS7, SignatureProfileNotETSI},
	SignatureLevelPAdESBES:         {SignatureFormPAdES, SignatureProfileExtendedBES},
	SignatureLevelPAdESEPES:        {SignatureFormPAdES, SignatureProfileExtendedEPES},
	SignatureLevelPAdESLTV:         {SignatureFormPAdES, SignatureProfileExtendedLTV},
	SignatureLevelPAdESBaselineB:   {SignatureFormPAdES, SignatureProfileBaselineB},
	SignatureLevelPAdESBaselineT:   {SignatureFormPAdES, SignatureProfileBaselineT},
	SignatureLevelPAdESBaselineLT:  {SignatureFormPAdES, SignatureProfileBaselineLT},
	SignatureLevelPAdESBaselineLTA: {SignatureFormPAdES, SignatureProfileBaselineLTA},

	SignatureLevelJSONNotETSI:      {SignatureFormJAdES, SignatureProfileNotETSI},
	SignatureLevelJAdES:            {SignatureFormJAdES, SignatureProfileAdES},
	SignatureLevelJAdESBaselineB:   {SignatureFormJAdES, SignatureProfileBaselineB},
	SignatureLevelJAdESBaselineT:   {SignatureFormJAdES, SignatureProfileBaselineT},
	SignatureLevelJAdESBaselineLT:  {SignatureFormJAdES, SignatureProfileBaselineLT},
	SignatureLevelJAdESBaselineLTA: {SignatureFormJAdES, SignatureProfileBaselineLTA},

	SignatureLevelCBORNotETSI:       {SignatureFormCBAdES, SignatureProfileNotETSI},
	SignatureLevelCBAdES:            {SignatureFormCBAdES, SignatureProfileAdES},
	SignatureLevelCBAdESBaselineB:   {SignatureFormCBAdES, SignatureProfileBaselineB},
	SignatureLevelCBAdESBaselineT:   {SignatureFormCBAdES, SignatureProfileBaselineT},
	SignatureLevelCBAdESBaselineLT:  {SignatureFormCBAdES, SignatureProfileBaselineLT},
	SignatureLevelCBAdESBaselineLTA: {SignatureFormCBAdES, SignatureProfileBaselineLTA},

	// UNKNOWN(null, NOT_ETSI): signatureForm is left as the zero value ("")
	// to mirror the Java null.
	SignatureLevelUnknown: {"", SignatureProfileNotETSI},
}

// SignatureLevelValues returns all constants in declaration order.
func SignatureLevelValues() []SignatureLevel {
	return []SignatureLevel{
		SignatureLevelXMLNotETSI, SignatureLevelXAdESBES, SignatureLevelXAdESEPES, SignatureLevelXAdEST, SignatureLevelXAdESLT,
		SignatureLevelXAdESC, SignatureLevelXAdESX, SignatureLevelXAdESXL, SignatureLevelXAdESA, SignatureLevelXAdESERS,
		SignatureLevelXAdESBaselineB, SignatureLevelXAdESBaselineT, SignatureLevelXAdESBaselineLT, SignatureLevelXAdESBaselineLTA,

		SignatureLevelCMSNotETSI, SignatureLevelCAdESBES, SignatureLevelCAdESEPES, SignatureLevelCAdEST, SignatureLevelCAdESLT,
		SignatureLevelCAdESC, SignatureLevelCAdESX, SignatureLevelCAdESXL, SignatureLevelCAdESA, SignatureLevelCAdESERS,
		SignatureLevelCAdESBaselineB, SignatureLevelCAdESBaselineT, SignatureLevelCAdESBaselineLT, SignatureLevelCAdESBaselineLTA,

		SignatureLevelPDFNotETSI, SignatureLevelPKCS7B, SignatureLevelPKCS7T, SignatureLevelPKCS7LT, SignatureLevelPKCS7LTA,
		SignatureLevelPAdESBES, SignatureLevelPAdESEPES, SignatureLevelPAdESLTV,
		SignatureLevelPAdESBaselineB, SignatureLevelPAdESBaselineT, SignatureLevelPAdESBaselineLT, SignatureLevelPAdESBaselineLTA,

		SignatureLevelJSONNotETSI, SignatureLevelJAdES, SignatureLevelJAdESBaselineB,
		SignatureLevelJAdESBaselineT, SignatureLevelJAdESBaselineLT, SignatureLevelJAdESBaselineLTA,

		SignatureLevelCBORNotETSI, SignatureLevelCBAdES, SignatureLevelCBAdESBaselineB,
		SignatureLevelCBAdESBaselineT, SignatureLevelCBAdESBaselineLT, SignatureLevelCBAdESBaselineLTA,

		SignatureLevelUnknown,
	}
}

// SignatureLevelValueOf returns the constant matching the given Java enum name.
func SignatureLevelValueOf(name string) (SignatureLevel, error) {
	for _, v := range SignatureLevelValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignatureLevel.%s", name)
}

// SignatureLevelValueByName returns the SignatureLevel based on the name,
// accepting a dash-separated spelling (e.g. "XAdES-BASELINE-B") in addition
// to the underscore-separated Java enum name.
func SignatureLevelValueByName(name string) (SignatureLevel, error) {
	return SignatureLevelValueOf(strings.ReplaceAll(name, "-", "_"))
}

// String returns the dash-separated representation of the SignatureLevel,
// mirroring Java's overridden toString().
func (s SignatureLevel) String() string {
	return strings.ReplaceAll(string(s), "_", "-")
}

// SignatureForm returns the corresponding SignatureForm, or an error if the
// signature level does not support one (i.e. SignatureLevelUnknown).
func (s SignatureLevel) SignatureForm() (SignatureForm, error) {
	form := signatureLevelData[s].signatureForm
	if form == "" {
		return "", fmt.Errorf("the signature level '%s' is not supported", s)
	}
	return form, nil
}

// SignatureProfile returns the corresponding SignatureProfile, or an error if the
// signature level does not have one. Upstream guards this with the same "not supported"
// throw as getSignatureForm(); every declared level does carry a profile, so the error
// path is only reachable for a SignatureLevel value that is not one of the constants.
func (s SignatureLevel) SignatureProfile() (SignatureProfile, error) {
	profile := signatureLevelData[s].signatureProfile
	if profile == "" {
		return "", fmt.Errorf("the signature level '%s' is not supported", s)
	}
	return profile, nil
}

// GetSignatureLevel returns an applicable SignatureLevel for the given SignatureForm and
// SignatureProfile.
//
// Upstream returns null when nothing matches — but never actually reaches that return.
// It compares against currentLevel.getSignatureForm() while iterating values(), and
// UNKNOWN, declared last, has a null form whose getter throws
// UnsupportedOperationException. So a lookup that finds no match throws rather than
// returning null, which is the behaviour reproduced here as an error return.
func GetSignatureLevel(signatureForm SignatureForm, signatureProfile SignatureProfile) (SignatureLevel, error) {
	for _, current := range SignatureLevelValues() {
		data := signatureLevelData[current]
		if data.signatureForm == "" {
			// Upstream's getSignatureForm() throws here before any comparison happens.
			return "", fmt.Errorf("the signature level '%s' is not supported", current)
		}
		if data.signatureForm == signatureForm && data.signatureProfile == signatureProfile {
			return current, nil
		}
	}
	return "", fmt.Errorf("the signature level '%s' is not supported", SignatureLevelUnknown)
}
