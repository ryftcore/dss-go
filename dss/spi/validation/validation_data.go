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

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ValidationData contains a validation data to be included into the signature.
//
// Java's HashMap iteration order is arbitrary but stable within a JVM run; a bare Go map is
// randomized on every run instead, and CertificateTokens()/CrlTokens()/OcspTokens() return in
// these maps' iteration order, so all three are kept insertion-ordered (slice + index map,
// PORTING.md's Collections rule) rather than bare maps.
type ValidationData struct {
	// certificateTokens is the set of certificate tokens, keyed by DSSIDAsString().
	certificateTokens *utils.OrderedMap[string, *model.CertificateToken]

	// crlTokens is the set of CRL tokens, keyed by DSSIDAsString().
	crlTokens *utils.OrderedMap[string, *spi.CRLToken]

	// ocspTokens is the set of OCSP tokens, keyed by DSSIDAsString().
	ocspTokens *utils.OrderedMap[string, *spi.OCSPToken]

	// storedPublicKeys is the internal set of containing public keys, keyed by
	// EntityIdentifier.String().
	storedPublicKeys map[string]*model.EntityIdentifier
}

// NewValidationData is the default constructor instantiating empty maps of tokens.
func NewValidationData() *ValidationData {
	return &ValidationData{
		certificateTokens: utils.NewOrderedMap[string, *model.CertificateToken](),
		crlTokens:         utils.NewOrderedMap[string, *spi.CRLToken](),
		ocspTokens:        utils.NewOrderedMap[string, *spi.OCSPToken](),
		storedPublicKeys:  make(map[string]*model.EntityIdentifier),
	}
}

// CertificateTokens gets certificate tokens to be included into the signature. Port of
// getCertificateTokens(); the returned slice is a snapshot, standing in for Java's
// Collections.unmodifiableSet.
func (v *ValidationData) CertificateTokens() []*model.CertificateToken {
	return v.certificateTokens.Values()
}

// CrlTokens gets CRL tokens to be included into the signature. Port of getCrlTokens(); the
// returned slice is a snapshot, standing in for Java's Collections.unmodifiableSet.
func (v *ValidationData) CrlTokens() []*spi.CRLToken {
	return v.crlTokens.Values()
}

// OcspTokens gets OCSP tokens to be included into the signature. Port of getOcspTokens(); the
// returned slice is a snapshot, standing in for Java's Collections.unmodifiableSet.
func (v *ValidationData) OcspTokens() []*spi.OCSPToken {
	return v.ocspTokens.Values()
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
		v.certificateTokens.Set(certificateToken.DSSIDAsString(), certificateToken)
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
			v.crlTokens.Set(crlToken.DSSIDAsString(), crlToken)
			return true
		}
	case enumerations.RevocationType_OCSP:
		ocspToken, ok := revocationToken.(*spi.OCSPToken)
		if !ok {
			panic(model.NewDSSError(fmt.Sprintf("Unexpected RevocationToken with Id '%s'", revocationToken.DSSIDAsString())))
		}
		if !v.containsOCSPToken(ocspToken) {
			v.ocspTokens.Set(ocspToken.DSSIDAsString(), ocspToken)
			return true
		}
	default:
		panic(model.NewDSSError(fmt.Sprintf("Unexpected RevocationToken with Id '%s'", revocationToken.DSSIDAsString())))
	}
	return false
}

func (v *ValidationData) containsCertificateToken(certificateTokenToAdd *model.CertificateToken) bool {
	if _, ok := v.certificateTokens.Get(certificateTokenToAdd.DSSIDAsString()); ok {
		return true
	}
	_, ok := v.storedPublicKeys[certificateTokenToAdd.EntityKey().String()]
	return ok
}

func (v *ValidationData) containsCRLToken(crlTokenToAdd *spi.CRLToken) bool {
	_, ok := v.crlTokens.Get(crlTokenToAdd.DSSIDAsString())
	return ok
}

func (v *ValidationData) containsOCSPToken(ocspTokenToAdd *spi.OCSPToken) bool {
	_, ok := v.ocspTokens.Get(ocspTokenToAdd.DSSIDAsString())
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
	for _, id := range v.certificateTokens.Keys() {
		certToken, _ := v.certificateTokens.Get(id)
		if entityIdentifier.Equals(certToken.EntityKey()) {
			v.certificateTokens.Delete(id)
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
	for _, id := range v.crlTokens.Keys() {
		crlToken, _ := v.crlTokens.Get(id)
		if _, ok := tokenIDsToExclude[crlToken.DSSIDAsString()]; ok {
			v.crlTokens.Delete(id)
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
	for _, id := range v.ocspTokens.Keys() {
		ocspToken, _ := v.ocspTokens.Get(id)
		if _, ok := tokenIDsToExclude[ocspToken.DSSIDAsString()]; ok {
			v.ocspTokens.Delete(id)
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
	return v.certificateTokens.Len() == 0 && v.crlTokens.Len() == 0 && v.ocspTokens.Len() == 0
}
