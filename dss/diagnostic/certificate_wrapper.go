// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/CertificateWrapper.java (DSS 6.5.RC1).
//
// getCertificateExtensionForOid(String, Class<T>) is a generic instance method in Java; Go has
// no generic methods, so it is ported as the package-level generic function
// CertificateExtensionForOid[T]. The xsi:type polymorphism of the schema's
// List<XmlCertificateExtension> is represented on the Go side (see jaxb/xml.go) by the
// jaxb.XmlCertificateExtensionItem interface, satisfied by every concrete extension struct (e.g.
// *jaxb.XmlSubjectAlternativeNames) via its embedded XmlCertificateExtensionContent/Attrs; that
// interface's OID accessor is ExtensionOID() *string (not OID() string).
package diagnostic

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// CertificateWrapper provides a user-friendly interface of dealing with JAXB XmlCertificate.
type CertificateWrapper struct {
	AbstractTokenProxyBase

	// certificate is the wrapped XmlCertificate instance.
	certificate *jaxb.XmlCertificate
}

// NewCertificateWrapper is the default constructor. Port of CertificateWrapper(XmlCertificate);
// panics per Objects.requireNonNull(certificate, "XMLCertificate cannot be null!").
func NewCertificateWrapper(certificate *jaxb.XmlCertificate) *CertificateWrapper {
	if certificate == nil {
		panic("XMLCertificate cannot be null!")
	}
	w := &CertificateWrapper{certificate: certificate}
	w.InitTokenProxy(w)
	return w
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *CertificateWrapper) Id() string {
	if w.certificate.Id != nil {
		return string(*w.certificate.Id)
	}
	return ""
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of
// getCurrentBasicSignature().
func (w *CertificateWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	return w.certificate.BasicSignature
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *CertificateWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	return w.certificate.CertificateChain.All()
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *CertificateWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	return w.certificate.SigningCertificate
}

// FoundCertificates is the AbstractTokenProxy default (not overridden in Java). Port of
// foundCertificates().
func (w *CertificateWrapper) FoundCertificates() *FoundCertificatesProxy {
	return DefaultFoundCertificates()
}

// FoundRevocations is the AbstractTokenProxy default (not overridden in Java). Port of
// foundRevocations().
func (w *CertificateWrapper) FoundRevocations() *FoundRevocationsProxy {
	return DefaultFoundRevocations()
}

// DigestMatchers is the AbstractTokenProxy default (not overridden in Java). Port of
// getDigestMatchers().
func (w *CertificateWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return nil
}

// IsTrusted reports whether the certificate is trusted. Port of isTrusted().
func (w *CertificateWrapper) IsTrusted() bool {
	return w.certificate.Trusted != nil && w.certificate.Trusted.Value
}

// TrustStartDate returns a certificate's trust start date, when available. If nil is
// returned and the certificate is trusted, the certificate is considered indefinitely trusted.
// Port of getTrustStartDate().
func (w *CertificateWrapper) TrustStartDate() *time.Time {
	if w.certificate.Trusted != nil && w.certificate.Trusted.StartDate != nil {
		t := w.certificate.Trusted.StartDate.Time()
		return &t
	}
	return nil
}

// TrustSunsetDate returns a certificate's trust end date, when available. If nil is
// returned and the certificate is trusted, the certificate is considered indefinitely trusted.
// Port of getTrustSunsetDate().
func (w *CertificateWrapper) TrustSunsetDate() *time.Time {
	if w.certificate.Trusted != nil && w.certificate.Trusted.SunsetDate != nil {
		t := w.certificate.Trusted.SunsetDate.Time()
		return &t
	}
	return nil
}

// IsSelfSigned reports whether the certificate is self-signed. Port of isSelfSigned().
func (w *CertificateWrapper) IsSelfSigned() bool {
	return w.certificate.SelfSigned
}

// CertificateExtensions returns a list of all certificate extensions. Port of
// getCertificateExtensions().
func (w *CertificateWrapper) CertificateExtensions() []jaxb.XmlCertificateExtensionItem {
	return append([]jaxb.XmlCertificateExtensionItem(nil), w.certificate.CertificateExtensions.All()...)
}

// CertificateExtensionForOid returns a certificate extension with the given oid when present.
// Port of the generic <T extends XmlCertificateExtension> T
// getCertificateExtensionForOid(String, Class<T>); panics (Java throws
// UnsupportedOperationException) when a match is found but does not have type T.
func CertificateExtensionForOid[T jaxb.XmlCertificateExtensionItem](w *CertificateWrapper, oid string) T {
	var zero T
	for _, certificateExtension := range w.CertificateExtensions() {
		extensionOID := certificateExtension.ExtensionOID()
		if extensionOID != nil && oid == *extensionOID {
			if typed, ok := certificateExtension.(T); ok {
				return typed
			}
			panic(fmt.Sprintf("A certificate extension with OID '%s' shall be in instance of '%T' class!", oid, zero))
		}
	}
	return zero
}

// CertificateExtensionsOids returns a list of all certificate extensions OIDs. Port of
// getCertificateExtensionsOids().
func (w *CertificateWrapper) CertificateExtensionsOids() []string {
	certificateExtensions := w.CertificateExtensions()
	var result []string
	if len(certificateExtensions) != 0 {
		for _, ext := range certificateExtensions {
			if extensionOID := ext.ExtensionOID(); extensionOID != nil {
				result = append(result, *extensionOID)
			}
		}
	}
	return result
}

// SubjectAlternativeNames returns subject alternative names. Port of
// getSubjectAlternativeNames().
func (w *CertificateWrapper) SubjectAlternativeNames() []*jaxb.XmlGeneralName {
	subjectAlternativeNames := w.getXmlSubjectAlternativeNames()
	if subjectAlternativeNames != nil {
		return subjectAlternativeNames.SubjectAlternativeName
	}
	return nil
}

func (w *CertificateWrapper) getXmlSubjectAlternativeNames() *jaxb.XmlSubjectAlternativeNames {
	return CertificateExtensionForOid[*jaxb.XmlSubjectAlternativeNames](w, enumerations.CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME.OID())
}

// IsCA reports whether the certificate defines BasicConstraints.cA extension set to TRUE. Port
// of isCA().
func (w *CertificateWrapper) IsCA() bool {
	basicConstraints := w.getXmlBasicConstraints()
	return basicConstraints != nil && basicConstraints.CA
}

// PathLenConstraint returns value of BasicConstraints.PathLenConstraint if present and
// BasicConstraints.cA is set to true. Port of getPathLenConstraint().
func (w *CertificateWrapper) PathLenConstraint() int {
	basicConstraints := w.getXmlBasicConstraints()
	if basicConstraints != nil && basicConstraints.CA && basicConstraints.PathLenConstraint != nil {
		return *basicConstraints.PathLenConstraint
	}
	return -1
}

func (w *CertificateWrapper) getXmlBasicConstraints() *jaxb.XmlBasicConstraints {
	return CertificateExtensionForOid[*jaxb.XmlBasicConstraints](w, enumerations.CertificateExtensionEnum_BASIC_CONSTRAINTS.OID())
}

// RequireExplicitPolicy returns value of the requireExplicitPolicy field of
// policyConstraints certificate extension. Port of getRequireExplicitPolicy().
func (w *CertificateWrapper) RequireExplicitPolicy() int {
	policyConstraints := w.getXmlPolicyConstraints()
	if policyConstraints != nil && policyConstraints.RequireExplicitPolicy != nil {
		return *policyConstraints.RequireExplicitPolicy
	}
	return -1
}

// InhibitPolicyMapping returns value of the inhibitPolicyMapping field of
// policyConstraints certificate extension. Port of getInhibitPolicyMapping().
func (w *CertificateWrapper) InhibitPolicyMapping() int {
	policyConstraints := w.getXmlPolicyConstraints()
	if policyConstraints != nil && policyConstraints.InhibitPolicyMapping != nil {
		return *policyConstraints.InhibitPolicyMapping
	}
	return -1
}

func (w *CertificateWrapper) getXmlPolicyConstraints() *jaxb.XmlPolicyConstraints {
	return CertificateExtensionForOid[*jaxb.XmlPolicyConstraints](w, enumerations.CertificateExtensionEnum_POLICY_CONSTRAINTS.OID())
}

// InhibitAnyPolicy returns value of the inhibitAnyPolicy certificate extension's value.
// Port of getInhibitAnyPolicy().
func (w *CertificateWrapper) InhibitAnyPolicy() int {
	inhibitAnyPolicy := w.getXmlInhibitAnyPolicy()
	if inhibitAnyPolicy != nil && inhibitAnyPolicy.Value != nil {
		return *inhibitAnyPolicy.Value
	}
	return -1
}

func (w *CertificateWrapper) getXmlInhibitAnyPolicy() *jaxb.XmlInhibitAnyPolicy {
	return CertificateExtensionForOid[*jaxb.XmlInhibitAnyPolicy](w, enumerations.CertificateExtensionEnum_INHIBIT_ANY_POLICY.OID())
}

// PermittedSubtrees returns value of the permittedSubtrees field of nameConstraints
// certificate extension, when present. Port of getPermittedSubtrees().
func (w *CertificateWrapper) PermittedSubtrees() []*jaxb.XmlGeneralSubtree {
	nameConstraints := w.getXmlNameConstraints()
	if nameConstraints != nil {
		return nameConstraints.PermittedSubtree
	}
	return nil
}

// ExcludedSubtrees returns value of the excludedSubtrees field of nameConstraints
// certificate extension, when present. Port of getExcludedSubtrees().
func (w *CertificateWrapper) ExcludedSubtrees() []*jaxb.XmlGeneralSubtree {
	nameConstraints := w.getXmlNameConstraints()
	if nameConstraints != nil {
		return nameConstraints.ExcludedSubtree
	}
	return nil
}

func (w *CertificateWrapper) getXmlNameConstraints() *jaxb.XmlNameConstraints {
	return CertificateExtensionForOid[*jaxb.XmlNameConstraints](w, enumerations.CertificateExtensionEnum_NAME_CONSTRAINTS.OID())
}

// KeyUsages returns the defined key-usages for the certificate. Port of getKeyUsages().
func (w *CertificateWrapper) KeyUsages() []enumerations.KeyUsageBit {
	keyUsage := w.getXmlKeyUsage()
	if keyUsage == nil {
		return nil
	}
	result := make([]enumerations.KeyUsageBit, len(keyUsage.KeyUsageBit))
	for i, v := range keyUsage.KeyUsageBit {
		result[i] = enumerations.KeyUsageBit(v)
	}
	return result
}

func (w *CertificateWrapper) getXmlKeyUsage() *jaxb.XmlKeyUsages {
	return CertificateExtensionForOid[*jaxb.XmlKeyUsages](w, enumerations.CertificateExtensionEnum_KEY_USAGE.OID())
}

// IsRevocationDataAvailable reports whether the revocation data is available for the
// certificate. Port of isRevocationDataAvailable().
func (w *CertificateWrapper) IsRevocationDataAvailable() bool {
	return len(w.certificate.Revocations.All()) != 0
}

// Sources returns a list of sources the certificate has been obtained from (e.g.
// TRUSTED_LIST, SIGNATURE, AIA, etc.). Port of getSources().
func (w *CertificateWrapper) Sources() []enumerations.CertificateSourceType {
	values := w.certificate.Sources.All()
	if values == nil {
		return nil
	}
	result := make([]enumerations.CertificateSourceType, len(values))
	for i, v := range values {
		result[i] = enumerations.CertificateSourceType(v)
	}
	return result
}

// CertificateRevocationData returns a list of revocation data relevant to the certificate.
// Port of getCertificateRevocationData().
func (w *CertificateWrapper) CertificateRevocationData() []*CertificateRevocationWrapper {
	var certRevocationWrappers []*CertificateRevocationWrapper
	for _, xmlCertificateRevocation := range w.certificate.Revocations.All() {
		certRevocationWrappers = append(certRevocationWrappers, NewCertificateRevocationWrapper(xmlCertificateRevocation))
	}
	return certRevocationWrappers
}

// RevocationDataById returns revocation data by its id. Port of getRevocationDataById(String).
func (w *CertificateWrapper) RevocationDataById(revocationId string) *CertificateRevocationWrapper {
	for _, revocationData := range w.CertificateRevocationData() {
		if revocationId == revocationData.Id() {
			return revocationData
		}
	}
	return nil
}

// IsIdPkixOcspNoCheck reports whether the certificate has id-pkix-ocsp-no-check attribute.
// Port of isIdPkixOcspNoCheck().
func (w *CertificateWrapper) IsIdPkixOcspNoCheck() bool {
	ocspNoCheck := w.getXmlIdPkixOcspNoCheck()
	return ocspNoCheck != nil && ocspNoCheck.Present != nil && *ocspNoCheck.Present
}

func (w *CertificateWrapper) getXmlIdPkixOcspNoCheck() *jaxb.XmlIdPkixOcspNoCheck {
	return CertificateExtensionForOid[*jaxb.XmlIdPkixOcspNoCheck](w, enumerations.CertificateExtensionEnum_OCSP_NOCHECK.OID())
}

// IsIdKpOCSPSigning checks if the certificate has an extended-key-usage "ocspSigning"
// (1.3.6.1.5.5.7.3.9). Port of isIdKpOCSPSigning().
func (w *CertificateWrapper) IsIdKpOCSPSigning() bool {
	extendedKeyUsage := w.getXmlExtendedKeyUsages()
	if extendedKeyUsage != nil {
		for _, xmlOID := range extendedKeyUsage.ExtendedKeyUsageOid {
			if enumerations.ExtendedKeyUsage_OCSP_SIGNING.OID() == xmlOID.Value {
				return true
			}
		}
	}
	return false
}

// IsValAssuredShortTermCertificate reports whether the certificate contains
// id-etsi-ext-valassured-ST-certs extension, as defined in ETSI EN 319 412-1 "5.2 Certificate
// Extensions regarding Validity Assured Certificate". Port of
// isValAssuredShortTermCertificate().
func (w *CertificateWrapper) IsValAssuredShortTermCertificate() bool {
	valAssuredShortTermCertificate := w.getXmlValAssuredShortTermCertificate()
	return valAssuredShortTermCertificate != nil && valAssuredShortTermCertificate.Present != nil && *valAssuredShortTermCertificate.Present
}

func (w *CertificateWrapper) getXmlValAssuredShortTermCertificate() *jaxb.XmlValAssuredShortTermCertificate {
	return CertificateExtensionForOid[*jaxb.XmlValAssuredShortTermCertificate](w, enumerations.CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM.OID())
}

// IsNoRevAvail reports whether the certificate contains noRevAvail extension, as defined in
// RFC 9608 "No Revocation Available for X.509 Public Key Certificates". Port of
// isNoRevAvail().
func (w *CertificateWrapper) IsNoRevAvail() bool {
	noRevAvail := w.getXmlNoRevAvail()
	return noRevAvail != nil && noRevAvail.Present != nil && *noRevAvail.Present
}

func (w *CertificateWrapper) getXmlNoRevAvail() *jaxb.XmlNoRevAvail {
	return CertificateExtensionForOid[*jaxb.XmlNoRevAvail](w, enumerations.CertificateExtensionEnum_NO_REVOCATION_AVAILABLE.OID())
}

// ExtendedKeyUsages returns a list of extended-key-usages. Port of getExtendedKeyUsages().
func (w *CertificateWrapper) ExtendedKeyUsages() []*jaxb.XmlOID {
	extendedKeyUsage := w.getXmlExtendedKeyUsages()
	if extendedKeyUsage != nil {
		return extendedKeyUsage.ExtendedKeyUsageOid
	}
	return nil
}

func (w *CertificateWrapper) getXmlExtendedKeyUsages() *jaxb.XmlExtendedKeyUsages {
	return CertificateExtensionForOid[*jaxb.XmlExtendedKeyUsages](w, enumerations.CertificateExtensionEnum_EXTENDED_KEY_USAGE.OID())
}

// NotBefore returns the certificate's notBefore date. Port of getNotBefore().
func (w *CertificateWrapper) NotBefore() *time.Time {
	if w.certificate.NotBefore == nil {
		return nil
	}
	t := w.certificate.NotBefore.Time()
	return &t
}

// NotAfter returns the certificate's notAfter date. Port of getNotAfter().
func (w *CertificateWrapper) NotAfter() *time.Time {
	if w.certificate.NotAfter == nil {
		return nil
	}
	t := w.certificate.NotAfter.Time()
	return &t
}

// EntityKey returns a string identifier of the certificate's entity key. Port of
// getEntityKey().
func (w *CertificateWrapper) EntityKey() string {
	if w.certificate.EntityKey != nil {
		return *w.certificate.EntityKey
	}
	return ""
}

// IssuerEntityKey returns a string identifier of the certificate's issuer entity key. Port
// of getIssuerEntityKey().
func (w *CertificateWrapper) IssuerEntityKey() string {
	if w.certificate.IssuerEntityKey != nil {
		return w.certificate.IssuerEntityKey.Value
	}
	return ""
}

// IsMatchingIssuerKey checks whether the issuer's public key matches to the key used to sign
// this token. Port of isMatchingIssuerKey().
func (w *CertificateWrapper) IsMatchingIssuerKey() bool {
	return w.certificate.IssuerEntityKey != nil && w.certificate.IssuerEntityKey.Key
}

// IsMatchingIssuerSubjectName checks whether the issuer's subject name matches to the key used
// to sign this token. Port of isMatchingIssuerSubjectName().
func (w *CertificateWrapper) IsMatchingIssuerSubjectName() bool {
	return w.certificate.IssuerEntityKey != nil && w.certificate.IssuerEntityKey.SubjectName
}

// CertificateTSPServiceExpiredCertsRevocationInfo returns expiredCertsRevocationInfo
// extension from TL Trusted Serviced. Port of
// getCertificateTSPServiceExpiredCertsRevocationInfo().
func (w *CertificateWrapper) CertificateTSPServiceExpiredCertsRevocationInfo() *time.Time {
	trustServiceProviders := w.certificate.TrustServiceProviders.All()
	if trustServiceProviders != nil {
		for _, trustServiceProvider := range trustServiceProviders {
			for _, xmlTrustService := range trustServiceProvider.TrustServices.All() {
				if xmlTrustService.ExpiredCertsRevocationInfo != nil {
					t := xmlTrustService.ExpiredCertsRevocationInfo.Time() // TODO improve
					return &t
				}
			}
		}
	}
	return nil
}

// SerialNumber returns the serial number of the certificate. Port of getSerialNumber().
func (w *CertificateWrapper) SerialNumber() string {
	serialNumber := w.certificate.SerialNumber
	if serialNumber == nil {
		return ""
	}
	return serialNumber.String()
}

// SubjectSerialNumber returns the subject serial number of the certificate. Port of
// getSubjectSerialNumber().
func (w *CertificateWrapper) SubjectSerialNumber() string {
	return stringOrEmpty(w.certificate.SubjectSerialNumber)
}

// Title returns the title. Port of getTitle().
func (w *CertificateWrapper) Title() string {
	return stringOrEmpty(w.certificate.Title)
}

// CommonName returns the common name. Port of getCommonName().
func (w *CertificateWrapper) CommonName() string {
	return stringOrEmpty(w.certificate.CommonName)
}

// CountryName returns the country code. Port of getCountryName().
func (w *CertificateWrapper) CountryName() string {
	return stringOrEmpty(w.certificate.CountryName)
}

// GivenName returns the given name. Port of getGivenName().
func (w *CertificateWrapper) GivenName() string {
	return stringOrEmpty(w.certificate.GivenName)
}

// OrganizationIdentifier returns the organization identifier. Port of
// getOrganizationIdentifier().
func (w *CertificateWrapper) OrganizationIdentifier() string {
	return stringOrEmpty(w.certificate.OrganizationIdentifier)
}

// OrganizationName returns the organization name. Port of getOrganizationName().
func (w *CertificateWrapper) OrganizationName() string {
	return stringOrEmpty(w.certificate.OrganizationName)
}

// OrganizationalUnit returns the organization unit. Port of getOrganizationalUnit().
func (w *CertificateWrapper) OrganizationalUnit() string {
	return stringOrEmpty(w.certificate.OrganizationalUnit)
}

// Email returns the email. Port of getEmail().
func (w *CertificateWrapper) Email() string {
	return stringOrEmpty(w.certificate.Email)
}

// Locality returns the locality. Port of getLocality().
func (w *CertificateWrapper) Locality() string {
	return stringOrEmpty(w.certificate.Locality)
}

// State returns the state. Port of getState().
func (w *CertificateWrapper) State() string {
	return stringOrEmpty(w.certificate.State)
}

// Surname returns the surname. Port of getSurname().
func (w *CertificateWrapper) Surname() string {
	return stringOrEmpty(w.certificate.Surname)
}

// Pseudo returns the pseudo. Port of getPseudo().
func (w *CertificateWrapper) Pseudo() string {
	return stringOrEmpty(w.certificate.Pseudonym)
}

// stringOrEmpty dereferences a *string field, standing in for Java null with the empty string
// (Go strings cannot be nil). Used throughout this file for the many optional xs:string
// JAXB elements/attributes on XmlCertificate and its extensions.
func stringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// DigestAlgoAndValue returns the certificate's Digest if present. Port of
// getDigestAlgoAndValue().
func (w *CertificateWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.certificate.DigestAlgoAndValue
}

// IsTrustedListReached reports whether the Trusted List has been reached for the particular
// certificate. Port of isTrustedListReached().
func (w *CertificateWrapper) IsTrustedListReached() bool {
	return len(w.certificate.TrustServiceProviders.All()) != 0
}

// TrustServiceProviders returns a list of XmlTrustServiceProvider. Port of
// getTrustServiceProviders().
func (w *CertificateWrapper) TrustServiceProviders() []*jaxb.XmlTrustServiceProvider {
	return w.certificate.TrustServiceProviders.All()
}

// TrustServices returns a list of TrustServiceWrapper. Port of getTrustServices().
//
// NOTE (cross-chunk assumption): TrustServiceWrapper is owned by DIAGWRAP_B; this constructs it
// as a plain struct literal with exported fields named after the Java setter names (minus
// "set"), per this port's convention elsewhere. Verify/adjust field names once DIAGWRAP_B lands.
func (w *CertificateWrapper) TrustServices() []*TrustServiceWrapper {
	var result []*TrustServiceWrapper
	tsps := w.certificate.TrustServiceProviders.All()
	for _, tsp := range tsps {
		tspNames := langAndValueValues(tsp.TSPNames.All())
		tspTradeNames := langAndValueValues(tsp.TSPTradeNames.All())
		for _, trustService := range tsp.TrustServices.All() {
			wrapper := &TrustServiceWrapper{
				TrustedList:              tsp.TL,
				ListOfTrustedLists:       tsp.LOTL,
				TspNames:                 tspNames,
				TspTradeNames:            tspTradeNames,
				ServiceDigitalIdentifier: NewCertificateWrapper(trustService.ServiceDigitalIdentifier),
				ServiceNames:             langAndValueValues(trustService.ServiceNames.All()),
				Status:                   stringOrEmpty(trustService.Status),
				Type:                     stringOrEmpty(trustService.ServiceType),
				StartDate:                xsDateTimePtr(trustService.StartDate),
				EndDate:                  xsDateTimePtr(trustService.EndDate),
				CapturedQualifiers:       append([]*jaxb.XmlQualifier(nil), trustService.CapturedQualifiers.All()...),
				AdditionalServiceInfos:   append([]string(nil), trustService.AdditionalServiceInfoUris.All()...),
				EnactedMRA:               trustService.EnactedMRA,
			}

			mraTrustServiceMapping := trustService.MRATrustServiceMapping
			if mraTrustServiceMapping != nil {
				wrapper.MraTrustServiceLegalIdentifier = stringOrEmpty(mraTrustServiceMapping.TrustServiceLegalIdentifier)
				wrapper.MraTrustServiceEquivalenceStatusStartingTime = xsDateTimePtr(mraTrustServiceMapping.EquivalenceStatusStartingTime)
				wrapper.MraTrustServiceEquivalenceStatusEndingTime = xsDateTimePtr(mraTrustServiceMapping.EquivalenceStatusEndingTime)
				originalThirdCountryMapping := mraTrustServiceMapping.OriginalThirdCountryMapping
				if originalThirdCountryMapping != nil {
					wrapper.OriginalTCType = stringOrEmpty(originalThirdCountryMapping.ServiceType)
					wrapper.OriginalTCStatus = stringOrEmpty(originalThirdCountryMapping.Status)
					wrapper.OriginalCapturedQualifiers = append([]*jaxb.XmlQualifier(nil), originalThirdCountryMapping.CapturedQualifiers.All()...)
					wrapper.OriginalTCAdditionalServiceInfos = append([]string(nil), originalThirdCountryMapping.AdditionalServiceInfoUris.All()...)
				}
			}

			result = append(result, wrapper)
		}
	}
	return result
}

// xsDateTimePtr converts an optional jaxb.XSDateTime field to *time.Time, tolerating nil.
func xsDateTimePtr(d *jaxb.XSDateTime) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time()
	return &t
}

// IsListOfTrustedEntitiesReached reports whether the List of Trusted Entities has been reached
// for the particular certificate. Port of isListOfTrustedEntitiesReached().
func (w *CertificateWrapper) IsListOfTrustedEntitiesReached() bool {
	return len(w.certificate.TrustedEntities.All()) != 0
}

// TrustedEntities returns a list of XmlTrustedEntity. Port of getTrustedEntities().
func (w *CertificateWrapper) TrustedEntities() []*jaxb.XmlTrustedEntity {
	return w.certificate.TrustedEntities.All()
}

// TrustedEntityServices returns a list of TrustedEntityServiceWrapper. Port of
// getTrustedEntityServices().
//
// NOTE (cross-chunk assumption): TrustedEntityServiceWrapper is owned by DIAGWRAP_B; see the
// same note on TrustServices().
func (w *CertificateWrapper) TrustedEntityServices() []*TrustedEntityServiceWrapper {
	var result []*TrustedEntityServiceWrapper
	tes := w.TrustedEntities()
	for _, te := range tes {
		entityNames := langAndValueValues(te.Names.All())
		tradeNames := langAndValueValues(te.TradeNames.All())
		for _, trustedService := range te.TrustedEntityServices.All() {
			wrapper := &TrustedEntityServiceWrapper{
				TrustedSourceList:        te.LoTE,
				ListOfTrustedSourceList:  te.LoLoTE,
				EntityNames:              entityNames,
				TradeNames:               tradeNames,
				ServiceDigitalIdentifier: NewCertificateWrapper(trustedService.ServiceDigitalIdentifier),
				ServiceNames:             langAndValueValues(trustedService.ServiceNames.All()),
				Status:                   stringOrEmpty(trustedService.Status),
				Type:                     stringOrEmpty(trustedService.ServiceType),
				StartDate:                xsDateTimePtr(trustedService.StartDate),
				EndDate:                  xsDateTimePtr(trustedService.EndDate),
				CapturedQualifiers:       append([]*jaxb.XmlQualifier(nil), trustedService.CapturedQualifiers.All()...),
				AdditionalServiceInfos:   append([]string(nil), trustedService.AdditionalServiceInfoUris.All()...),
			}
			result = append(result, wrapper)
		}
	}
	return result
}

func langAndValueValues(langAndValues []*jaxb.XmlLangAndValue) []string {
	var result []string
	for _, v := range langAndValues {
		result = append(result, v.Value)
	}
	return result
}

// CertificateDN returns the certificate's Distinguished Name (by RFC 2253). Port of
// getCertificateDN().
func (w *CertificateWrapper) CertificateDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.certificate.SubjectDistinguishedName)
	return distinguishedNameListWrapper.Value("RFC2253")
}

// CertificateIssuerDN returns the certificate issuer's Distinguished Name (by RFC 2253).
// Port of getCertificateIssuerDN().
func (w *CertificateWrapper) CertificateIssuerDN() string {
	distinguishedNameListWrapper := NewDistinguishedNameListWrapper(w.certificate.IssuerDistinguishedName)
	return distinguishedNameListWrapper.Value("RFC2253")
}

// CRLDistributionPoints returns the CRL Distribution Points URLs. Port of
// getCRLDistributionPoints().
func (w *CertificateWrapper) CRLDistributionPoints() []string {
	crlDistributionPoints := w.getXmlCRLDistributionPoints()
	if crlDistributionPoints != nil {
		return crlDistributionPoints.CrlUrl
	}
	return nil
}

func (w *CertificateWrapper) getXmlCRLDistributionPoints() *jaxb.XmlCRLDistributionPoints {
	return CertificateExtensionForOid[*jaxb.XmlCRLDistributionPoints](w, enumerations.CertificateExtensionEnum_CRL_DISTRIBUTION_POINTS.OID())
}

// FreshestCRLUrls returns the Freshest CRL URLs. Port of getFreshestCRLUrls().
func (w *CertificateWrapper) FreshestCRLUrls() []string {
	freshestCRL := w.getXmlFreshestCRL()
	if freshestCRL != nil {
		return freshestCRL.CrlUrl
	}
	return nil
}

func (w *CertificateWrapper) getXmlFreshestCRL() *jaxb.XmlFreshestCRL {
	return CertificateExtensionForOid[*jaxb.XmlFreshestCRL](w, enumerations.CertificateExtensionEnum_FRESHEST_CRL.OID())
}

// CAIssuersAccessUrls returns the Authority Information Access URLs. Port of
// getCAIssuersAccessUrls().
func (w *CertificateWrapper) CAIssuersAccessUrls() []string {
	authorityInformationAccess := w.getXmlAuthorityInformationAccess()
	if authorityInformationAccess != nil {
		return authorityInformationAccess.CaIssuersUrl
	}
	return nil
}

// OCSPAccessUrls returns the OCSP Access URLs. Port of getOCSPAccessUrls().
func (w *CertificateWrapper) OCSPAccessUrls() []string {
	authorityInformationAccess := w.getXmlAuthorityInformationAccess()
	if authorityInformationAccess != nil {
		return authorityInformationAccess.OcspUrl
	}
	return nil
}

func (w *CertificateWrapper) getXmlAuthorityInformationAccess() *jaxb.XmlAuthorityInformationAccess {
	return CertificateExtensionForOid[*jaxb.XmlAuthorityInformationAccess](w, enumerations.CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS.OID())
}

// AuthorityKeyIdentifier returns the Authority Key Identifier certificate extension's
// value, when present. Port of getAuthorityKeyIdentifier().
func (w *CertificateWrapper) AuthorityKeyIdentifier() []byte {
	xmlAuthorityKeyIdentifier := w.getXmlAuthorityKeyIdentifier()
	if xmlAuthorityKeyIdentifier != nil && xmlAuthorityKeyIdentifier.KeyIdentifier != nil {
		return []byte(*xmlAuthorityKeyIdentifier.KeyIdentifier)
	}
	return nil
}

// AuthorityKeyIdentifierIssuerSerial returns the Authority Key Identifier certificate
// extension's value, which is a combination of authorityCertIssuer and
// authorityCertSerialNumber fields, when present. Port of
// getAuthorityKeyIdentifierIssuerSerial().
func (w *CertificateWrapper) AuthorityKeyIdentifierIssuerSerial() []byte {
	xmlAuthorityKeyIdentifier := w.getXmlAuthorityKeyIdentifier()
	if xmlAuthorityKeyIdentifier != nil && xmlAuthorityKeyIdentifier.AuthorityCertIssuerSerial != nil {
		return []byte(*xmlAuthorityKeyIdentifier.AuthorityCertIssuerSerial)
	}
	return nil
}

func (w *CertificateWrapper) getXmlAuthorityKeyIdentifier() *jaxb.XmlAuthorityKeyIdentifier {
	return CertificateExtensionForOid[*jaxb.XmlAuthorityKeyIdentifier](w, enumerations.CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER.OID())
}

// SubjectKeyIdentifier returns the Subject Key Identifier certificate extension's value,
// when present. Port of getSubjectKeyIdentifier().
func (w *CertificateWrapper) SubjectKeyIdentifier() []byte {
	xmlSubjectKeyIdentifier := w.getXmlSubjectKeyIdentifier()
	if xmlSubjectKeyIdentifier != nil && xmlSubjectKeyIdentifier.Ski != nil {
		return []byte(*xmlSubjectKeyIdentifier.Ski)
	}
	return nil
}

func (w *CertificateWrapper) getXmlSubjectKeyIdentifier() *jaxb.XmlSubjectKeyIdentifier {
	return CertificateExtensionForOid[*jaxb.XmlSubjectKeyIdentifier](w, enumerations.CertificateExtensionEnum_SUBJECT_KEY_IDENTIFIER.OID())
}

// CpsUrls returns the certificate policies URLs. Port of getCpsUrls().
func (w *CertificateWrapper) CpsUrls() []string {
	var result []string
	xmlCertificatePolicies := w.getXmlCertificatePolicies()
	if xmlCertificatePolicies != nil {
		for _, xmlCertificatePolicy := range xmlCertificatePolicies.CertificatePolicy {
			cpsUrl := xmlCertificatePolicy.CpsUrl
			if cpsUrl != nil && *cpsUrl != "" {
				result = append(result, *cpsUrl)
			}
		}
	}
	return result
}

// PolicyIds returns the certificate policies Ids. Port of getPolicyIds(); Java's
// XmlCertificatePolicy extends XmlOID, ported here as XmlCertificatePolicy embedding
// jaxb.XmlOID so its Value field is promoted directly.
func (w *CertificateWrapper) PolicyIds() []string {
	xmlCertificatePolicies := w.getXmlCertificatePolicies()
	if xmlCertificatePolicies == nil {
		return nil
	}
	var result []string
	for _, p := range xmlCertificatePolicies.CertificatePolicy {
		result = append(result, p.Value)
	}
	return result
}

// CertificatePolicies returns the certificate policies Ids. Port of
// getCertificatePolicies().
func (w *CertificateWrapper) CertificatePolicies() []*jaxb.XmlCertificatePolicy {
	xmlCertificatePolicies := w.getXmlCertificatePolicies()
	if xmlCertificatePolicies != nil {
		return xmlCertificatePolicies.CertificatePolicy
	}
	return nil
}

func (w *CertificateWrapper) getXmlCertificatePolicies() *jaxb.XmlCertificatePolicies {
	return CertificateExtensionForOid[*jaxb.XmlCertificatePolicies](w, enumerations.CertificateExtensionEnum_CERTIFICATE_POLICIES.OID())
}

// CertificatePoliciesOids returns the certificate policies OIDs. Port of
// getCertificatePoliciesOids().
func (w *CertificateWrapper) CertificatePoliciesOids() []string {
	var result []string
	certificatePolicies := w.CertificatePolicies()
	if len(certificatePolicies) != 0 {
		for _, cp := range certificatePolicies {
			result = append(result, cp.Value)
		}
	}
	return result
}

// IsQcCompliance reports whether the certificate is QC compliant (has id-etsi-qcs-QcCompliance
// extension). Port of isQcCompliance().
func (w *CertificateWrapper) IsQcCompliance() bool {
	xmlQcStatements := w.getXmlQcStatements()
	return xmlQcStatements != nil && xmlQcStatements.QcCompliance != nil && xmlQcStatements.QcCompliance.Present
}

// IsSupportedByQSCD reports whether the certificate is supported by QSCD (has
// id-etsi-qcs-QcSSCD extension). Port of isSupportedByQSCD().
func (w *CertificateWrapper) IsSupportedByQSCD() bool {
	xmlQcStatements := w.getXmlQcStatements()
	return xmlQcStatements != nil && xmlQcStatements.QcSSCD != nil && xmlQcStatements.QcSSCD.Present
}

// QcTypes returns a list of QCTypes (present inside id-etsi-qcs-QcType extension). Port of
// getQcTypes().
func (w *CertificateWrapper) QcTypes() []enumerations.QCType {
	var result []enumerations.QCType
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		for _, oid := range xmlQcStatements.QcTypes.All() {
			result = append(result, enumerations.QCTypeFromOID(oid.Value))
		}
	}
	return result
}

// QcLegislationCountryCodes returns a list of QCLegislation Country Codes (present inside
// id-etsi-qcs-QcCClegislation extension). Port of getQcLegislationCountryCodes().
func (w *CertificateWrapper) QcLegislationCountryCodes() []string {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return xmlQcStatements.QcCClegislation.All()
	}
	return nil
}

// PSD2Info returns the PSD2 QCStatement (id-etsi-psd2-qcStatement extension, ETSI TS 119
// 495). Port of getPSD2Info().
func (w *CertificateWrapper) PSD2Info() *PSD2InfoWrapper {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.PSD2QcInfo != nil {
		return NewPSD2InfoWrapper(xmlQcStatements.PSD2QcInfo)
	}
	return nil
}

// QCLimitValue returns the QCEuLimitValue. Port of getQCLimitValue().
func (w *CertificateWrapper) QCLimitValue() *QCLimitValueWrapper {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.QcEuLimitValue != nil {
		return NewQCLimitValueWrapper(xmlQcStatements.QcEuLimitValue)
	}
	return nil
}

// QCEuRetentionPeriod returns QcEuRetentionPeriod. Port of getQCEuRetentionPeriod().
func (w *CertificateWrapper) QCEuRetentionPeriod() *int {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return xmlQcStatements.QcEuRetentionPeriod
	}
	return nil
}

// QCPDSLocations returns QcEuPDS Locations. Port of getQCPDSLocations().
func (w *CertificateWrapper) QCPDSLocations() []*jaxb.XmlLangAndValue {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return xmlQcStatements.QcEuPDS.All()
	}
	return nil
}

// SemanticsIdentifier returns the semantics identifier. Port of getSemanticsIdentifier().
func (w *CertificateWrapper) SemanticsIdentifier() enumerations.SemanticsIdentifier {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.SemanticsIdentifier != nil {
		xmlOID := xmlQcStatements.SemanticsIdentifier
		if xmlOID != nil {
			return enumerations.SemanticsIdentifierFromOID(xmlOID.Value)
		}
	}
	return ""
}

// QcQSCDLegislation returns a list of QcQSCDlegislation country codes (present inside
// id-etsi-qcs-QcQSCDlegislation extension). Port of getQcQSCDLegislation().
func (w *CertificateWrapper) QcQSCDLegislation() []string {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return xmlQcStatements.QcQSCDlegislation.All()
	}
	return nil
}

// QcIdentMethod returns a QCIdentMethod (present inside id-etsi-qcs-QcIdentMethod
// extension). Port of getQcIdentMethod().
func (w *CertificateWrapper) QcIdentMethod() enumerations.QCIdentMethod {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.QcIdentMethod != nil {
		return enumerations.QCIdentMethodFromOID(xmlQcStatements.QcIdentMethod.Value)
	}
	return nil
}

// QcPSB gets QcPSB defined in the QcStatements of the certificate, when present. Port of
// getQcPSB().
func (w *CertificateWrapper) QcPSB() *QCPSBWrapper {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.QcPSB != nil {
		return NewQCPSBWrapper(xmlQcStatements.QcPSB)
	}
	return nil
}

// OtherQcStatements returns a list of QcStatements OIDs not supported by the
// implementation. Port of getOtherQcStatements().
func (w *CertificateWrapper) OtherQcStatements() []string {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return oidValues(xmlQcStatements.OtherOIDs.All())
	}
	return nil
}

// IsEnactedMRA reports whether the MRA has been enacted. Port of isEnactedMRA().
func (w *CertificateWrapper) IsEnactedMRA() bool {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil {
		return xmlQcStatements.EnactedMRA != nil && *xmlQcStatements.EnactedMRA
	}
	return false
}

// MRAEnactedTrustServiceLegalIdentifier returns a name of a Trusted Service used to apply
// translation for the certificate QcStatements based on the defined Mutual Recognition
// Agreement scheme. Port of getMRAEnactedTrustServiceLegalIdentifier().
func (w *CertificateWrapper) MRAEnactedTrustServiceLegalIdentifier() string {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.MRACertificateMapping != nil &&
		xmlQcStatements.MRACertificateMapping.TrustServiceEquivalenceInformation != nil {
		return stringOrEmpty(xmlQcStatements.MRACertificateMapping.TrustServiceEquivalenceInformation.TrustServiceLegalIdentifier)
	}
	return ""
}

// MRACertificateContentEquivalenceList returns a XmlCertificateContentEquivalence list
// corresponding to the matching MRA information. Port of
// getMRACertificateContentEquivalenceList().
func (w *CertificateWrapper) MRACertificateContentEquivalenceList() []*jaxb.XmlCertificateContentEquivalence {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.MRACertificateMapping != nil &&
		xmlQcStatements.MRACertificateMapping.TrustServiceEquivalenceInformation != nil {
		return xmlQcStatements.MRACertificateMapping.TrustServiceEquivalenceInformation.CertificateContentEquivalenceList.All()
	}
	return nil
}

func (w *CertificateWrapper) getOriginalThirdCountryMapping() *jaxb.XmlOriginalThirdCountryQcStatementsMapping {
	xmlQcStatements := w.getXmlQcStatements()
	if xmlQcStatements != nil && xmlQcStatements.MRACertificateMapping != nil {
		return xmlQcStatements.MRACertificateMapping.OriginalThirdCountryMapping
	}
	return nil
}

func (w *CertificateWrapper) getXmlQcStatements() *jaxb.XmlQcStatements {
	return CertificateExtensionForOid[*jaxb.XmlQcStatements](w, enumerations.CertificateExtensionEnum_QC_STATEMENTS.OID())
}

func oidValues(xmlOids []*jaxb.XmlOID) []string {
	var result []string
	for _, xmlOID := range xmlOids {
		result = append(result, xmlOID.Value)
	}
	return result
}

// IsOriginalThirdCountryQcCompliance reports whether the certificate has been defined as QC
// compliant in a third-country Trusted List before MRA mapping. Port of
// isOriginalThirdCountryQcCompliance().
func (w *CertificateWrapper) IsOriginalThirdCountryQcCompliance() bool {
	originalThirdCountryMapping := w.getOriginalThirdCountryMapping()
	return originalThirdCountryMapping != nil && originalThirdCountryMapping.QcCompliance != nil && originalThirdCountryMapping.QcCompliance.Present
}

// IsOriginalThirdCountrySupportedByQSCD reports whether the certificate has been defined as
// supported by QSCD in a third-country Trusted List before MRA mapping. Port of
// isOriginalThirdCountrySupportedByQSCD().
func (w *CertificateWrapper) IsOriginalThirdCountrySupportedByQSCD() bool {
	originalThirdCountryMapping := w.getOriginalThirdCountryMapping()
	return originalThirdCountryMapping != nil && originalThirdCountryMapping.QcSSCD != nil && originalThirdCountryMapping.QcSSCD.Present
}

// OriginalThirdCountryQCTypes returns a list of QCTypes defined in a third-country Trusted
// List before MRA mapping. Port of getOriginalThirdCountryQCTypes().
func (w *CertificateWrapper) OriginalThirdCountryQCTypes() []enumerations.QCType {
	var result []enumerations.QCType
	originalThirdCountryMapping := w.getOriginalThirdCountryMapping()
	if originalThirdCountryMapping != nil {
		for _, oid := range originalThirdCountryMapping.QcTypes.All() {
			result = append(result, enumerations.QCTypeFromOID(oid.Value))
		}
	}
	return result
}

// OriginalThirdCountryQcLegislationCountryCodes returns a list of QCLegislation Country
// Codes defined in a third-country Trusted List before MRA mapping. Port of
// getOriginalThirdCountryQcLegislationCountryCodes().
func (w *CertificateWrapper) OriginalThirdCountryQcLegislationCountryCodes() []string {
	originalThirdCountryMapping := w.getOriginalThirdCountryMapping()
	if originalThirdCountryMapping != nil {
		return originalThirdCountryMapping.QcCClegislation.All()
	}
	return nil
}

// OriginalThirdCountryOtherQcStatements returns a list of QcStatements OIDs not supported
// by the implementation defined in a third-country Trusted List before MRA mapping. Port of
// getOriginalThirdCountryOtherQcStatements().
func (w *CertificateWrapper) OriginalThirdCountryOtherQcStatements() []string {
	originalThirdCountryMapping := w.getOriginalThirdCountryMapping()
	if originalThirdCountryMapping != nil {
		return oidValues(originalThirdCountryMapping.OtherOIDs.All())
	}
	return nil
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *CertificateWrapper) Binaries() []byte {
	if w.certificate.Base64Encoded == nil {
		return nil
	}
	return []byte(*w.certificate.Base64Encoded)
}

// ReadableCertificateName returns human-readable certificate name. Port of
// getReadableCertificateName().
func (w *CertificateWrapper) ReadableCertificateName() string {
	if w.certificate.CommonName != nil {
		return *w.certificate.CommonName
	}
	if w.certificate.GivenName != nil {
		return *w.certificate.GivenName
	}
	if w.certificate.Surname != nil {
		return *w.certificate.Surname
	}
	if w.certificate.Pseudonym != nil {
		return *w.certificate.Pseudonym
	}
	if w.certificate.OrganizationName != nil {
		return *w.certificate.OrganizationName
	}
	if w.certificate.OrganizationalUnit != nil {
		return *w.certificate.OrganizationalUnit
	}
	return "?"
}
