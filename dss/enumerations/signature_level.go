// Ported from dss-enumerations/.../SignatureLevel.java (DSS 6.5.RC1).
//
// NOTE: SignatureForm and SignatureProfile (with their exported constants,
// e.g. SignatureForm_CAdES, SignatureProfile_BASELINE_B, preserving the
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
// AdES form, and SignatureLevel_UNKNOWN marks a level DSS could not
// determine.
const (
	SignatureLevel_XML_NOT_ETSI       SignatureLevel = "XML_NOT_ETSI"
	SignatureLevel_XAdES_BES          SignatureLevel = "XAdES_BES"
	SignatureLevel_XAdES_EPES         SignatureLevel = "XAdES_EPES"
	SignatureLevel_XAdES_T            SignatureLevel = "XAdES_T"
	SignatureLevel_XAdES_LT           SignatureLevel = "XAdES_LT"
	SignatureLevel_XAdES_C            SignatureLevel = "XAdES_C"
	SignatureLevel_XAdES_X            SignatureLevel = "XAdES_X"
	SignatureLevel_XAdES_XL           SignatureLevel = "XAdES_XL"
	SignatureLevel_XAdES_A            SignatureLevel = "XAdES_A"
	SignatureLevel_XAdES_ERS          SignatureLevel = "XAdES_ERS"
	SignatureLevel_XAdES_BASELINE_B   SignatureLevel = "XAdES_BASELINE_B"
	SignatureLevel_XAdES_BASELINE_T   SignatureLevel = "XAdES_BASELINE_T"
	SignatureLevel_XAdES_BASELINE_LT  SignatureLevel = "XAdES_BASELINE_LT"
	SignatureLevel_XAdES_BASELINE_LTA SignatureLevel = "XAdES_BASELINE_LTA"

	SignatureLevel_CMS_NOT_ETSI       SignatureLevel = "CMS_NOT_ETSI"
	SignatureLevel_CAdES_BES          SignatureLevel = "CAdES_BES"
	SignatureLevel_CAdES_EPES         SignatureLevel = "CAdES_EPES"
	SignatureLevel_CAdES_T            SignatureLevel = "CAdES_T"
	SignatureLevel_CAdES_LT           SignatureLevel = "CAdES_LT"
	SignatureLevel_CAdES_C            SignatureLevel = "CAdES_C"
	SignatureLevel_CAdES_X            SignatureLevel = "CAdES_X"
	SignatureLevel_CAdES_XL           SignatureLevel = "CAdES_XL"
	SignatureLevel_CAdES_A            SignatureLevel = "CAdES_A"
	SignatureLevel_CAdES_ERS          SignatureLevel = "CAdES_ERS"
	SignatureLevel_CAdES_BASELINE_B   SignatureLevel = "CAdES_BASELINE_B"
	SignatureLevel_CAdES_BASELINE_T   SignatureLevel = "CAdES_BASELINE_T"
	SignatureLevel_CAdES_BASELINE_LT  SignatureLevel = "CAdES_BASELINE_LT"
	SignatureLevel_CAdES_BASELINE_LTA SignatureLevel = "CAdES_BASELINE_LTA"

	SignatureLevel_PDF_NOT_ETSI       SignatureLevel = "PDF_NOT_ETSI"
	SignatureLevel_PKCS7_B            SignatureLevel = "PKCS7_B"
	SignatureLevel_PKCS7_T            SignatureLevel = "PKCS7_T"
	SignatureLevel_PKCS7_LT           SignatureLevel = "PKCS7_LT"
	SignatureLevel_PKCS7_LTA          SignatureLevel = "PKCS7_LTA"
	SignatureLevel_PAdES_BES          SignatureLevel = "PAdES_BES"
	SignatureLevel_PAdES_EPES         SignatureLevel = "PAdES_EPES"
	SignatureLevel_PAdES_LTV          SignatureLevel = "PAdES_LTV"
	SignatureLevel_PAdES_BASELINE_B   SignatureLevel = "PAdES_BASELINE_B"
	SignatureLevel_PAdES_BASELINE_T   SignatureLevel = "PAdES_BASELINE_T"
	SignatureLevel_PAdES_BASELINE_LT  SignatureLevel = "PAdES_BASELINE_LT"
	SignatureLevel_PAdES_BASELINE_LTA SignatureLevel = "PAdES_BASELINE_LTA"

	SignatureLevel_JSON_NOT_ETSI      SignatureLevel = "JSON_NOT_ETSI"
	SignatureLevel_JAdES              SignatureLevel = "JAdES"
	SignatureLevel_JAdES_BASELINE_B   SignatureLevel = "JAdES_BASELINE_B"
	SignatureLevel_JAdES_BASELINE_T   SignatureLevel = "JAdES_BASELINE_T"
	SignatureLevel_JAdES_BASELINE_LT  SignatureLevel = "JAdES_BASELINE_LT"
	SignatureLevel_JAdES_BASELINE_LTA SignatureLevel = "JAdES_BASELINE_LTA"

	SignatureLevel_CBOR_NOT_ETSI        SignatureLevel = "CBOR_NOT_ETSI"
	SignatureLevel_CB_AdES              SignatureLevel = "CB_AdES"
	SignatureLevel_CB_AdES_BASELINE_B   SignatureLevel = "CB_AdES_BASELINE_B"
	SignatureLevel_CB_AdES_BASELINE_T   SignatureLevel = "CB_AdES_BASELINE_T"
	SignatureLevel_CB_AdES_BASELINE_LT  SignatureLevel = "CB_AdES_BASELINE_LT"
	SignatureLevel_CB_AdES_BASELINE_LTA SignatureLevel = "CB_AdES_BASELINE_LTA"

	SignatureLevel_UNKNOWN SignatureLevel = "UNKNOWN"
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
	SignatureLevel_XML_NOT_ETSI:       {SignatureForm_XAdES, SignatureProfile_NOT_ETSI},
	SignatureLevel_XAdES_BES:          {SignatureForm_XAdES, SignatureProfile_EXTENDED_BES},
	SignatureLevel_XAdES_EPES:         {SignatureForm_XAdES, SignatureProfile_EXTENDED_EPES},
	SignatureLevel_XAdES_T:            {SignatureForm_XAdES, SignatureProfile_EXTENDED_T},
	SignatureLevel_XAdES_LT:           {SignatureForm_XAdES, SignatureProfile_EXTENDED_LT},
	SignatureLevel_XAdES_C:            {SignatureForm_XAdES, SignatureProfile_EXTENDED_C},
	SignatureLevel_XAdES_X:            {SignatureForm_XAdES, SignatureProfile_EXTENDED_X},
	SignatureLevel_XAdES_XL:           {SignatureForm_XAdES, SignatureProfile_EXTENDED_XL},
	SignatureLevel_XAdES_A:            {SignatureForm_XAdES, SignatureProfile_EXTENDED_A},
	SignatureLevel_XAdES_ERS:          {SignatureForm_XAdES, SignatureProfile_EXTENDED_ERS},
	SignatureLevel_XAdES_BASELINE_B:   {SignatureForm_XAdES, SignatureProfile_BASELINE_B},
	SignatureLevel_XAdES_BASELINE_T:   {SignatureForm_XAdES, SignatureProfile_BASELINE_T},
	SignatureLevel_XAdES_BASELINE_LT:  {SignatureForm_XAdES, SignatureProfile_BASELINE_LT},
	SignatureLevel_XAdES_BASELINE_LTA: {SignatureForm_XAdES, SignatureProfile_BASELINE_LTA},

	SignatureLevel_CMS_NOT_ETSI:       {SignatureForm_CAdES, SignatureProfile_NOT_ETSI},
	SignatureLevel_CAdES_BES:          {SignatureForm_CAdES, SignatureProfile_EXTENDED_BES},
	SignatureLevel_CAdES_EPES:         {SignatureForm_CAdES, SignatureProfile_EXTENDED_EPES},
	SignatureLevel_CAdES_T:            {SignatureForm_CAdES, SignatureProfile_EXTENDED_T},
	SignatureLevel_CAdES_LT:           {SignatureForm_CAdES, SignatureProfile_EXTENDED_LT},
	SignatureLevel_CAdES_C:            {SignatureForm_CAdES, SignatureProfile_EXTENDED_C},
	SignatureLevel_CAdES_X:            {SignatureForm_CAdES, SignatureProfile_EXTENDED_X},
	SignatureLevel_CAdES_XL:           {SignatureForm_CAdES, SignatureProfile_EXTENDED_XL},
	SignatureLevel_CAdES_A:            {SignatureForm_CAdES, SignatureProfile_EXTENDED_A},
	SignatureLevel_CAdES_ERS:          {SignatureForm_CAdES, SignatureProfile_EXTENDED_ERS},
	SignatureLevel_CAdES_BASELINE_B:   {SignatureForm_CAdES, SignatureProfile_BASELINE_B},
	SignatureLevel_CAdES_BASELINE_T:   {SignatureForm_CAdES, SignatureProfile_BASELINE_T},
	SignatureLevel_CAdES_BASELINE_LT:  {SignatureForm_CAdES, SignatureProfile_BASELINE_LT},
	SignatureLevel_CAdES_BASELINE_LTA: {SignatureForm_CAdES, SignatureProfile_BASELINE_LTA},

	SignatureLevel_PDF_NOT_ETSI:       {SignatureForm_PAdES, SignatureProfile_NOT_ETSI},
	SignatureLevel_PKCS7_B:            {SignatureForm_PKCS7, SignatureProfile_NOT_ETSI},
	SignatureLevel_PKCS7_T:            {SignatureForm_PKCS7, SignatureProfile_NOT_ETSI},
	SignatureLevel_PKCS7_LT:           {SignatureForm_PKCS7, SignatureProfile_NOT_ETSI},
	SignatureLevel_PKCS7_LTA:          {SignatureForm_PKCS7, SignatureProfile_NOT_ETSI},
	SignatureLevel_PAdES_BES:          {SignatureForm_PAdES, SignatureProfile_EXTENDED_BES},
	SignatureLevel_PAdES_EPES:         {SignatureForm_PAdES, SignatureProfile_EXTENDED_EPES},
	SignatureLevel_PAdES_LTV:          {SignatureForm_PAdES, SignatureProfile_EXTENDED_LTV},
	SignatureLevel_PAdES_BASELINE_B:   {SignatureForm_PAdES, SignatureProfile_BASELINE_B},
	SignatureLevel_PAdES_BASELINE_T:   {SignatureForm_PAdES, SignatureProfile_BASELINE_T},
	SignatureLevel_PAdES_BASELINE_LT:  {SignatureForm_PAdES, SignatureProfile_BASELINE_LT},
	SignatureLevel_PAdES_BASELINE_LTA: {SignatureForm_PAdES, SignatureProfile_BASELINE_LTA},

	SignatureLevel_JSON_NOT_ETSI:      {SignatureForm_JAdES, SignatureProfile_NOT_ETSI},
	SignatureLevel_JAdES:              {SignatureForm_JAdES, SignatureProfile_AdES},
	SignatureLevel_JAdES_BASELINE_B:   {SignatureForm_JAdES, SignatureProfile_BASELINE_B},
	SignatureLevel_JAdES_BASELINE_T:   {SignatureForm_JAdES, SignatureProfile_BASELINE_T},
	SignatureLevel_JAdES_BASELINE_LT:  {SignatureForm_JAdES, SignatureProfile_BASELINE_LT},
	SignatureLevel_JAdES_BASELINE_LTA: {SignatureForm_JAdES, SignatureProfile_BASELINE_LTA},

	SignatureLevel_CBOR_NOT_ETSI:        {SignatureForm_CBAdES, SignatureProfile_NOT_ETSI},
	SignatureLevel_CB_AdES:              {SignatureForm_CBAdES, SignatureProfile_AdES},
	SignatureLevel_CB_AdES_BASELINE_B:   {SignatureForm_CBAdES, SignatureProfile_BASELINE_B},
	SignatureLevel_CB_AdES_BASELINE_T:   {SignatureForm_CBAdES, SignatureProfile_BASELINE_T},
	SignatureLevel_CB_AdES_BASELINE_LT:  {SignatureForm_CBAdES, SignatureProfile_BASELINE_LT},
	SignatureLevel_CB_AdES_BASELINE_LTA: {SignatureForm_CBAdES, SignatureProfile_BASELINE_LTA},

	// UNKNOWN(null, NOT_ETSI): signatureForm is left as the zero value ("")
	// to mirror the Java null.
	SignatureLevel_UNKNOWN: {"", SignatureProfile_NOT_ETSI},
}

// SignatureLevelValues returns all constants in declaration order.
func SignatureLevelValues() []SignatureLevel {
	return []SignatureLevel{
		SignatureLevel_XML_NOT_ETSI, SignatureLevel_XAdES_BES, SignatureLevel_XAdES_EPES, SignatureLevel_XAdES_T, SignatureLevel_XAdES_LT,
		SignatureLevel_XAdES_C, SignatureLevel_XAdES_X, SignatureLevel_XAdES_XL, SignatureLevel_XAdES_A, SignatureLevel_XAdES_ERS,
		SignatureLevel_XAdES_BASELINE_B, SignatureLevel_XAdES_BASELINE_T, SignatureLevel_XAdES_BASELINE_LT, SignatureLevel_XAdES_BASELINE_LTA,

		SignatureLevel_CMS_NOT_ETSI, SignatureLevel_CAdES_BES, SignatureLevel_CAdES_EPES, SignatureLevel_CAdES_T, SignatureLevel_CAdES_LT,
		SignatureLevel_CAdES_C, SignatureLevel_CAdES_X, SignatureLevel_CAdES_XL, SignatureLevel_CAdES_A, SignatureLevel_CAdES_ERS,
		SignatureLevel_CAdES_BASELINE_B, SignatureLevel_CAdES_BASELINE_T, SignatureLevel_CAdES_BASELINE_LT, SignatureLevel_CAdES_BASELINE_LTA,

		SignatureLevel_PDF_NOT_ETSI, SignatureLevel_PKCS7_B, SignatureLevel_PKCS7_T, SignatureLevel_PKCS7_LT, SignatureLevel_PKCS7_LTA,
		SignatureLevel_PAdES_BES, SignatureLevel_PAdES_EPES, SignatureLevel_PAdES_LTV,
		SignatureLevel_PAdES_BASELINE_B, SignatureLevel_PAdES_BASELINE_T, SignatureLevel_PAdES_BASELINE_LT, SignatureLevel_PAdES_BASELINE_LTA,

		SignatureLevel_JSON_NOT_ETSI, SignatureLevel_JAdES, SignatureLevel_JAdES_BASELINE_B,
		SignatureLevel_JAdES_BASELINE_T, SignatureLevel_JAdES_BASELINE_LT, SignatureLevel_JAdES_BASELINE_LTA,

		SignatureLevel_CBOR_NOT_ETSI, SignatureLevel_CB_AdES, SignatureLevel_CB_AdES_BASELINE_B,
		SignatureLevel_CB_AdES_BASELINE_T, SignatureLevel_CB_AdES_BASELINE_LT, SignatureLevel_CB_AdES_BASELINE_LTA,

		SignatureLevel_UNKNOWN,
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
// signature level does not support one (i.e. SignatureLevel_UNKNOWN).
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
	return "", fmt.Errorf("the signature level '%s' is not supported", SignatureLevel_UNKNOWN)
}
