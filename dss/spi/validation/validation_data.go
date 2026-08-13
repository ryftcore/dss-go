// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/ValidationData.java (DSS 6.5.RC1).
//
// Java's Set<CertificateToken>/Set<CRLToken>/Set<OCSPToken> (HashSet, keyed by equals/
// hashCode) become maps keyed by DSSIDAsString(), per the convention established by
// spi.ListCertificateSource (Set<T> -> map[T]struct{}/map[string]*T behind small helpers, see
// PORTING.md's Collections section and dss-model's CertificateToken.Equals).
//
// slf4j logging is dropped per the phase 2a handoff fact ("slf4j dropped unless
// load-bearing"); the TRACE-level messages here carry no behaviour.
package validation

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// ValidationData contains a validation data to be included into the signature.
type ValidationData struct {
	// certificateTokens is the set of certificate tokens, keyed by DSSIDAsString().
	certificateTokens map[string]*model.CertificateToken

	// crlTokens is the set of CRL tokens, keyed by DSSIDAsString().
	crlTokens map[string]*spi.CRLToken

	// ocspTokens is the set of OCSP tokens, keyed by DSSIDAsString().
	ocspTokens map[string]*spi.OCSPToken

	// storedPublicKeys is the internal set of containing public keys, keyed by
	// EntityIdentifier.String().
	storedPublicKeys map[string]*model.EntityIdentifier
}

// NewValidationData is the default constructor instantiating empty maps of tokens.
func NewValidationData() *ValidationData {
	return &ValidationData{
		certificateTokens: make(map[string]*model.CertificateToken),
		crlTokens:         make(map[string]*spi.CRLToken),
		ocspTokens:        make(map[string]*spi.OCSPToken),
		storedPublicKeys:  make(map[string]*model.EntityIdentifier),
	}
}

// CertificateTokens gets certificate tokens to be included into the signature. Port of
// getCertificateTokens(); the returned slice is a snapshot, standing in for Java's
// Collections.unmodifiableSet.
func (v *ValidationData) CertificateTokens() []*model.CertificateToken {
	tokens := make([]*model.CertificateToken, 0, len(v.certificateTokens))
	for _, t := range v.certificateTokens {
		tokens = append(tokens, t)
	}
	return tokens
}

// CrlTokens gets CRL tokens to be included into the signature. Port of getCrlTokens(); the
// returned slice is a snapshot, standing in for Java's Collections.unmodifiableSet.
func (v *ValidationData) CrlTokens() []*spi.CRLToken {
	tokens := make([]*spi.CRLToken, 0, len(v.crlTokens))
	for _, t := range v.crlTokens {
		tokens = append(tokens, t)
	}
	return tokens
}

// OcspTokens gets OCSP tokens to be included into the signature. Port of getOcspTokens(); the
// returned slice is a snapshot, standing in for Java's Collections.unmodifiableSet.
func (v *ValidationData) OcspTokens() []*spi.OCSPToken {
	tokens := make([]*spi.OCSPToken, 0, len(v.ocspTokens))
	for _, t := range v.ocspTokens {
		tokens = append(tokens, t)
	}
	return tokens
}

// AddToken adds a validation data token and returns whether the token has been added
// successfully. Port of addToken(Token).
//
// Java's `instanceof RevocationToken` check (a raw-typed, wildcard instanceof) becomes a type
// assertion against AnyRevocationToken (this package's non-generic stand-in for
// RevocationToken<?>, defined in revocation_data_loading_strategy.go).
func (v *ValidationData) AddToken(token model.Token) bool {
	switch t := token.(type) {
	case *model.CertificateToken:
		if v.addCertificateToken(t) {
			return true
		}
	default:
		if revocationToken, ok := token.(AnyRevocationToken); ok {
			if v.addRevocationToken(revocationToken) {
				return true
			}
		} else {
			panic(model.NewDSSError(fmt.Sprintf("Unexpected token with Id '%s'", token.DSSIDAsString())))
		}
	}
	return false
}

func (v *ValidationData) addCertificateToken(certificateToken *model.CertificateToken) bool {
	if !v.containsCertificateToken(certificateToken) {
		v.certificateTokens[certificateToken.DSSIDAsString()] = certificateToken
		v.storedPublicKeys[certificateToken.EntityKey().String()] = certificateToken.EntityKey()
		return true
	}
	return false
}

func (v *ValidationData) addRevocationToken(revocationToken AnyRevocationToken) bool {
	switch revocationToken.RevocationType() {
	case enumerations.RevocationType_CRL:
		crlToken, ok := revocationToken.(*spi.CRLToken)
		if !ok {
			panic(model.NewDSSError(fmt.Sprintf("Unexpected RevocationToken with Id '%s'", revocationToken.DSSIDAsString())))
		}
		if !v.containsCRLToken(crlToken) {
			v.crlTokens[crlToken.DSSIDAsString()] = crlToken
			return true
		}
	case enumerations.RevocationType_OCSP:
		ocspToken, ok := revocationToken.(*spi.OCSPToken)
		if !ok {
			panic(model.NewDSSError(fmt.Sprintf("Unexpected RevocationToken with Id '%s'", revocationToken.DSSIDAsString())))
		}
		if !v.containsOCSPToken(ocspToken) {
			v.ocspTokens[ocspToken.DSSIDAsString()] = ocspToken
			return true
		}
	default:
		panic(model.NewDSSError(fmt.Sprintf("Unexpected RevocationToken with Id '%s'", revocationToken.DSSIDAsString())))
	}
	return false
}

func (v *ValidationData) containsCertificateToken(certificateTokenToAdd *model.CertificateToken) bool {
	if _, ok := v.certificateTokens[certificateTokenToAdd.DSSIDAsString()]; ok {
		return true
	}
	_, ok := v.storedPublicKeys[certificateTokenToAdd.EntityKey().String()]
	return ok
}

func (v *ValidationData) containsCRLToken(crlTokenToAdd *spi.CRLToken) bool {
	_, ok := v.crlTokens[crlTokenToAdd.DSSIDAsString()]
	return ok
}

func (v *ValidationData) containsOCSPToken(ocspTokenToAdd *spi.OCSPToken) bool {
	_, ok := v.ocspTokens[ocspTokenToAdd.DSSIDAsString()]
	return ok
}

// AddValidationData allows to add all tokens from a provided validation data to the current
// collection. Port of addValidationData(...).
func (v *ValidationData) AddValidationData(validationData *ValidationData) {
	for _, token := range validationData.CertificateTokens() {
		v.AddToken(token)
	}
	for _, token := range validationData.CrlTokens() {
		v.AddToken(token)
	}
	for _, token := range validationData.OcspTokens() {
		v.AddToken(token)
	}
}

// ExcludeValidationData excludes validationDataToExclude from the current validation data
// set. Port of excludeValidationData(...).
func (v *ValidationData) ExcludeValidationData(validationDataToExclude *ValidationData) {
	v.ExcludeCertificateTokens(validationDataToExclude.CertificateTokens())
	v.ExcludeCRLTokensCollection(validationDataToExclude.CrlTokens())
	v.ExcludeOCSPTokensCollection(validationDataToExclude.OcspTokens())
}

// ExcludeCertificateTokens removes all certificate token entries matching the provided
// collection. Port of excludeCertificateTokens(...).
func (v *ValidationData) ExcludeCertificateTokens(certificateTokensToExclude []*model.CertificateToken) {
	for _, certificateToken := range certificateTokensToExclude {
		if v.containsCertificateToken(certificateToken) {
			entityKey := certificateToken.EntityKey()
			delete(v.storedPublicKeys, entityKey.String())
			v.excludeWithEntityKey(entityKey)
		}
	}
}

func (v *ValidationData) excludeWithEntityKey(entityIdentifier *model.EntityIdentifier) {
	for id, certToken := range v.certificateTokens {
		if entityIdentifier.Equals(certToken.EntityKey()) {
			delete(v.certificateTokens, id)
		}
	}
}

// ExcludeCRLTokens removes all CRL token entries matching the provided collection of
// encapsulated CRL binaries. Port of excludeCRLTokens(...).
func (v *ValidationData) ExcludeCRLTokens(crlTokensToExclude []model.Identifier) {
	if len(crlTokensToExclude) == 0 {
		return
	}
	tokenIDsToExclude := make(map[string]struct{}, len(crlTokensToExclude))
	for _, identifier := range crlTokensToExclude {
		tokenIDsToExclude[identifier.AsXmlID()] = struct{}{}
	}
	for id, crlToken := range v.crlTokens {
		if _, ok := tokenIDsToExclude[crlToken.DSSIDAsString()]; ok {
			delete(v.crlTokens, id)
		}
	}
}

// ExcludeCRLTokensCollection removes all CRL token entries matching crlTokensToExclude from
// the current CRL data set. Port of excludeCRLTokensCollection(...).
func (v *ValidationData) ExcludeCRLTokensCollection(crlTokensToExclude []*spi.CRLToken) {
	if len(crlTokensToExclude) == 0 {
		return
	}
	identifiers := make([]model.Identifier, 0, len(crlTokensToExclude))
	for _, t := range crlTokensToExclude {
		identifiers = append(identifiers, t.DSSID())
	}
	v.ExcludeCRLTokens(identifiers)
}

// ExcludeOCSPTokens removes all OCSP token entries matching the provided collection of
// encapsulated OCSP binaries. Port of excludeOCSPTokens(...).
func (v *ValidationData) ExcludeOCSPTokens(ocspTokensToExclude []model.Identifier) {
	if len(ocspTokensToExclude) == 0 {
		return
	}
	tokenIDsToExclude := make(map[string]struct{}, len(ocspTokensToExclude))
	for _, identifier := range ocspTokensToExclude {
		tokenIDsToExclude[identifier.AsXmlID()] = struct{}{}
	}
	for id, ocspToken := range v.ocspTokens {
		if _, ok := tokenIDsToExclude[ocspToken.DSSIDAsString()]; ok {
			delete(v.ocspTokens, id)
		}
	}
}

// ExcludeOCSPTokensCollection removes all OCSP token entries matching ocspTokensToExclude from
// the current OCSP data set. Port of excludeOCSPTokensCollection(...).
func (v *ValidationData) ExcludeOCSPTokensCollection(ocspTokensToExclude []*spi.OCSPToken) {
	if len(ocspTokensToExclude) == 0 {
		return
	}
	identifiers := make([]model.Identifier, 0, len(ocspTokensToExclude))
	for _, t := range ocspTokensToExclude {
		identifiers = append(identifiers, t.DSSID())
	}
	v.ExcludeOCSPTokens(identifiers)
}

// IsEmpty checks if the validation data is empty. Port of isEmpty().
func (v *ValidationData) IsEmpty() bool {
	return len(v.certificateTokens) == 0 && len(v.crlTokens) == 0 && len(v.ocspTokens) == 0
}
