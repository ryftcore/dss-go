// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/CertificateExtensions.java (DSS 6.5.RC1).
package extension

// CertificateExtensions contains a set of certificate extensions processed by the
// application.
type CertificateExtensions struct {
	// keyUsage is RFC 5280: 4.2.1.3. Key Usage.
	keyUsage *KeyUsage

	// certificatePolicies is RFC 5280: 4.2.1.4. Certificate Policies.
	certificatePolicies *CertificatePolicies

	// subjectAlternativeNames is RFC 5280: 4.2.1.6. Subject Alternative Name.
	subjectAlternativeNames *SubjectAlternativeNames

	// basicConstraints is RFC 5280: 4.2.1.9. Basic Constraints.
	basicConstraints *BasicConstraints

	// nameConstraints is RFC 5280: 4.2.1.10. Name Constraints.
	nameConstraints *NameConstraints

	// policyConstraints is RFC 5280: 4.2.1.11. Policy Constraints.
	policyConstraints *PolicyConstraints

	// extendedKeyUsage is RFC 5280: 4.2.1.12. Extended Key Usage.
	extendedKeyUsage *ExtendedKeyUsages

	// crlDistributionPoints is RFC 5280: 4.2.1.13. CRL Distribution Points.
	crlDistributionPoints *CRLDistributionPoints

	// inhibitAnyPolicy is RFC 5280: 4.2.1.14. Inhibit anyPolicy.
	inhibitAnyPolicy *InhibitAnyPolicy

	// freshestCRL is RFC 5280: 4.2.1.15. Freshest CRL (a.k.a. Delta CRL Distribution
	// Point).
	freshestCRL *FreshestCRL

	// authorityKeyIdentifier is RFC 5280: 4.2.1.1. Authority Key Identifier.
	authorityKeyIdentifier *AuthorityKeyIdentifier

	// subjectKeyIdentifier is RFC 5280: 4.2.1.2. Subject Key Identifier.
	subjectKeyIdentifier *SubjectKeyIdentifier

	// authorityInformationAccess is RFC 5280: 4.2.2.1. Authority Information Access.
	authorityInformationAccess *AuthorityInformationAccess

	// ocspNoCheck is RFC 6960: 4.2.2.2.1. Revocation Checking of an Authorized
	// Responder.
	ocspNoCheck *OCSPNoCheck

	// validityAssuredShortTerm is ETSI EN 319 412-1: 5.2.2 Validity Assured - Short
	// Term.
	validityAssuredShortTerm *ValidityAssuredShortTerm

	// noRevAvail is RFC 9608: No Revocation Available for X.509 Public Key
	// Certificates.
	noRevAvail *NoRevAvail

	// qcStatements is ETSI EN 319 412-1/5: QCStatements.
	qcStatements *QcStatements

	// otherExtensions lists other extensions.
	otherExtensions []*CertificateExtension

	// allExtensions lists all certificate extensions.
	allExtensions []*CertificateExtension
}

// NewCertificateExtensions instantiates the object with null values. Ports the default
// constructor.
func NewCertificateExtensions() *CertificateExtensions {
	return &CertificateExtensions{
		otherExtensions: []*CertificateExtension{},
		allExtensions:   []*CertificateExtension{},
	}
}

// KeyUsage returns the key usage.
func (c *CertificateExtensions) KeyUsage() *KeyUsage {
	return c.keyUsage
}

// SetKeyUsage sets the key usage.
func (c *CertificateExtensions) SetKeyUsage(keyUsage *KeyUsage) {
	c.keyUsage = keyUsage
	if keyUsage != nil {
		c.addToAllExtensionsList(&keyUsage.CertificateExtension)
	}
}

// CertificatePolicies returns the certificate policies.
func (c *CertificateExtensions) CertificatePolicies() *CertificatePolicies {
	return c.certificatePolicies
}

// SetCertificatePolicies sets the certificate policies.
func (c *CertificateExtensions) SetCertificatePolicies(certificatePolicies *CertificatePolicies) {
	c.certificatePolicies = certificatePolicies
	if certificatePolicies != nil {
		c.addToAllExtensionsList(&certificatePolicies.CertificateExtension)
	}
}

// SubjectAlternativeNames returns the subject alternative names.
func (c *CertificateExtensions) SubjectAlternativeNames() *SubjectAlternativeNames {
	return c.subjectAlternativeNames
}

// SetSubjectAlternativeNames sets the subject alternative names.
func (c *CertificateExtensions) SetSubjectAlternativeNames(subjectAlternativeNames *SubjectAlternativeNames) {
	c.subjectAlternativeNames = subjectAlternativeNames
	if subjectAlternativeNames != nil {
		c.addToAllExtensionsList(&subjectAlternativeNames.CertificateExtension)
	}
}

// BasicConstraints returns the basic constraints.
func (c *CertificateExtensions) BasicConstraints() *BasicConstraints {
	return c.basicConstraints
}

// SetBasicConstraints sets the basic constraints.
func (c *CertificateExtensions) SetBasicConstraints(basicConstraints *BasicConstraints) {
	c.basicConstraints = basicConstraints
	if basicConstraints != nil {
		c.addToAllExtensionsList(&basicConstraints.CertificateExtension)
	}
}

// NameConstraints returns the name constraints.
func (c *CertificateExtensions) NameConstraints() *NameConstraints {
	return c.nameConstraints
}

// SetNameConstraints sets the name constraints.
func (c *CertificateExtensions) SetNameConstraints(nameConstraints *NameConstraints) {
	c.nameConstraints = nameConstraints
	if nameConstraints != nil {
		c.addToAllExtensionsList(&nameConstraints.CertificateExtension)
	}
}

// PolicyConstraints returns the policy constraints.
func (c *CertificateExtensions) PolicyConstraints() *PolicyConstraints {
	return c.policyConstraints
}

// SetPolicyConstraints sets the policy constraints.
func (c *CertificateExtensions) SetPolicyConstraints(policyConstraints *PolicyConstraints) {
	c.policyConstraints = policyConstraints
	if policyConstraints != nil {
		c.addToAllExtensionsList(&policyConstraints.CertificateExtension)
	}
}

// ExtendedKeyUsage returns the extended key usages.
func (c *CertificateExtensions) ExtendedKeyUsage() *ExtendedKeyUsages {
	return c.extendedKeyUsage
}

// SetExtendedKeyUsage sets the extended key usages.
func (c *CertificateExtensions) SetExtendedKeyUsage(extendedKeyUsage *ExtendedKeyUsages) {
	c.extendedKeyUsage = extendedKeyUsage
	if extendedKeyUsage != nil {
		c.addToAllExtensionsList(&extendedKeyUsage.CertificateExtension)
	}
}

// CRLDistributionPoints returns the CRL distribution points.
func (c *CertificateExtensions) CRLDistributionPoints() *CRLDistributionPoints {
	return c.crlDistributionPoints
}

// SetCRLDistributionPoints sets the CRL distribution points.
func (c *CertificateExtensions) SetCRLDistributionPoints(crlDistributionPoints *CRLDistributionPoints) {
	c.crlDistributionPoints = crlDistributionPoints
	if crlDistributionPoints != nil {
		c.addToAllExtensionsList(&crlDistributionPoints.CertificateExtension)
	}
}

// InhibitAnyPolicy returns the InhibitAnyPolicy extension.
func (c *CertificateExtensions) InhibitAnyPolicy() *InhibitAnyPolicy {
	return c.inhibitAnyPolicy
}

// SetInhibitAnyPolicy sets the InhibitAnyPolicy extension.
func (c *CertificateExtensions) SetInhibitAnyPolicy(inhibitAnyPolicy *InhibitAnyPolicy) {
	c.inhibitAnyPolicy = inhibitAnyPolicy
	if inhibitAnyPolicy != nil {
		c.addToAllExtensionsList(&inhibitAnyPolicy.CertificateExtension)
	}
}

// FreshestCRL returns the FreshestCRL extension.
func (c *CertificateExtensions) FreshestCRL() *FreshestCRL {
	return c.freshestCRL
}

// SetFreshestCRL sets the FreshestCRL extension.
func (c *CertificateExtensions) SetFreshestCRL(freshestCRL *FreshestCRL) {
	c.freshestCRL = freshestCRL
	if freshestCRL != nil {
		c.addToAllExtensionsList(&freshestCRL.CertificateExtension)
	}
}

// AuthorityKeyIdentifier returns the authority key identifier.
func (c *CertificateExtensions) AuthorityKeyIdentifier() *AuthorityKeyIdentifier {
	return c.authorityKeyIdentifier
}

// SetAuthorityKeyIdentifier sets the authority key identifier.
func (c *CertificateExtensions) SetAuthorityKeyIdentifier(authorityKeyIdentifier *AuthorityKeyIdentifier) {
	c.authorityKeyIdentifier = authorityKeyIdentifier
	if authorityKeyIdentifier != nil {
		c.addToAllExtensionsList(&authorityKeyIdentifier.CertificateExtension)
	}
}

// SubjectKeyIdentifier returns the subject key identifier.
func (c *CertificateExtensions) SubjectKeyIdentifier() *SubjectKeyIdentifier {
	return c.subjectKeyIdentifier
}

// SetSubjectKeyIdentifier sets the subject key identifier.
func (c *CertificateExtensions) SetSubjectKeyIdentifier(subjectKeyIdentifier *SubjectKeyIdentifier) {
	c.subjectKeyIdentifier = subjectKeyIdentifier
	if subjectKeyIdentifier != nil {
		c.addToAllExtensionsList(&subjectKeyIdentifier.CertificateExtension)
	}
}

// AuthorityInformationAccess returns the authority information access.
func (c *CertificateExtensions) AuthorityInformationAccess() *AuthorityInformationAccess {
	return c.authorityInformationAccess
}

// SetAuthorityInformationAccess sets the authority information access.
func (c *CertificateExtensions) SetAuthorityInformationAccess(authorityInformationAccess *AuthorityInformationAccess) {
	c.authorityInformationAccess = authorityInformationAccess
	if authorityInformationAccess != nil {
		c.addToAllExtensionsList(&authorityInformationAccess.CertificateExtension)
	}
}

// OcspNoCheck returns the ocsp-nocheck value.
func (c *CertificateExtensions) OcspNoCheck() *OCSPNoCheck {
	return c.ocspNoCheck
}

// SetOcspNoCheck sets the ocsp-nocheck value.
func (c *CertificateExtensions) SetOcspNoCheck(ocspNoCheck *OCSPNoCheck) {
	c.ocspNoCheck = ocspNoCheck
	if ocspNoCheck != nil {
		c.addToAllExtensionsList(&ocspNoCheck.CertificateExtension)
	}
}

// ValidityAssuredShortTerm returns the ext-etsi-valassured-ST-certs value.
func (c *CertificateExtensions) ValidityAssuredShortTerm() *ValidityAssuredShortTerm {
	return c.validityAssuredShortTerm
}

// SetValidityAssuredShortTerm sets the ext-etsi-valassured-ST-certs value.
func (c *CertificateExtensions) SetValidityAssuredShortTerm(validityAssuredShortTerm *ValidityAssuredShortTerm) {
	c.validityAssuredShortTerm = validityAssuredShortTerm
	if validityAssuredShortTerm != nil {
		c.addToAllExtensionsList(&validityAssuredShortTerm.CertificateExtension)
	}
}

// NoRevAvail returns the noRevAvail value.
func (c *CertificateExtensions) NoRevAvail() *NoRevAvail {
	return c.noRevAvail
}

// SetNoRevAvail sets the noRevAvail value.
func (c *CertificateExtensions) SetNoRevAvail(noRevAvail *NoRevAvail) {
	c.noRevAvail = noRevAvail
	if noRevAvail != nil {
		c.addToAllExtensionsList(&noRevAvail.CertificateExtension)
	}
}

// QcStatements returns the QcStatements.
func (c *CertificateExtensions) QcStatements() *QcStatements {
	return c.qcStatements
}

// SetQcStatements sets the QcStatements.
func (c *CertificateExtensions) SetQcStatements(qcStatements *QcStatements) {
	c.qcStatements = qcStatements
	if qcStatements != nil {
		c.addToAllExtensionsList(&qcStatements.CertificateExtension)
	}
}

// OtherExtensions returns a list of other certificate extensions.
func (c *CertificateExtensions) OtherExtensions() []*CertificateExtension {
	return c.otherExtensions
}

// AddOtherExtension adds another certificate extension.
func (c *CertificateExtensions) AddOtherExtension(certificateExtension *CertificateExtension) {
	c.otherExtensions = append(c.otherExtensions, certificateExtension)
	c.addToAllExtensionsList(certificateExtension)
}

// AllCertificateExtensions returns a list of all certificate extensions.
func (c *CertificateExtensions) AllCertificateExtensions() []*CertificateExtension {
	return c.allExtensions
}

// addToAllExtensionsList mirrors CertificateExtensions#addToAllExtensionsList: it is a
// no-op for a nil certificateExtension.
func (c *CertificateExtensions) addToAllExtensionsList(certificateExtension *CertificateExtension) {
	if certificateExtension != nil {
		c.allExtensions = append(c.allExtensions, certificateExtension)
	}
}
