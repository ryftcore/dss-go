// Ported from dss-policy-jaxb/.../policy/EtsiValidationPolicy.java (DSS 6.5.RC1).
//
// This class encapsulates the constraint file that controls the policy to be
// used during the validation process. It adds the functions to direct access
// to the file data. It is the implementation of the ETSI TS 102 853
// standard.
//
// Every accessor mirrors Java's null-safe navigation (a nil intermediate
// jaxb node makes the whole accessor return the interface's nil, not a
// wrapper holding a nil pointer - see the toXXX helpers at the bottom of
// this file) with one exception: like the Java class itself, this port never
// nil-checks the top-level policy field before dereferencing it (e.g.
// PolicyName reads p.policy.Name directly) - EtsiValidationPolicy's Java
// constructor does not validate its argument either, so a nil policy panics
// on first use here exactly as it NPEs there.
package policy

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// etsiValidationPolicyDefaultValidationModel is the default validation model
// (SHELL). Ports EtsiValidationPolicy#DEFAULT_VALIDATION_MODEL.
const etsiValidationPolicyDefaultValidationModel = enumerations.ValidationModelShell

// EtsiValidationPolicy encapsulates the constraint file that controls the
// policy used during the validation process. Ports EtsiValidationPolicy.
type EtsiValidationPolicy struct {
	// policy holds the validation policy constraints.
	policy *jaxb.ConstraintsParameters
}

var _ modelpolicy.ValidationPolicy = (*EtsiValidationPolicy)(nil)

// NewEtsiValidationPolicy is the default constructor.
func NewEtsiValidationPolicy(policy *jaxb.ConstraintsParameters) *EtsiValidationPolicy {
	return &EtsiValidationPolicy{policy: policy}
}

// PolicyName gets the name of the policy. Ports EtsiValidationPolicy#getPolicyName.
func (p *EtsiValidationPolicy) PolicyName() string {
	if p.policy.Name != nil {
		return *p.policy.Name
	}
	return ""
}

// PolicyDescription gets the policy description. Ports
// EtsiValidationPolicy#getPolicyDescription.
func (p *EtsiValidationPolicy) PolicyDescription() string {
	return p.policy.Description
}

// SignaturePolicyConstraint indicates if the signature policy should be
// checked. Ports EtsiValidationPolicy#getSignaturePolicyConstraint.
func (p *EtsiValidationPolicy) SignaturePolicyConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toMultiValuesRule(sc.AcceptablePolicies)
	}
	return nil
}

// SignaturePolicyIdentifiedConstraint indicates if the signature policy
// validation should be processed. Ports
// EtsiValidationPolicy#getSignaturePolicyIdentifiedConstraint.
func (p *EtsiValidationPolicy) SignaturePolicyIdentifiedConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toLevelRule(sc.PolicyAvailable)
	}
	return nil
}

// SignaturePolicyStorePresentConstraint indicates if a SignaturePolicyStore
// unsigned attribute presence shall be checked. Ports
// EtsiValidationPolicy#getSignaturePolicyStorePresentConstraint.
func (p *EtsiValidationPolicy) SignaturePolicyStorePresentConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toLevelRule(sc.SignaturePolicyStorePresent)
	}
	return nil
}

// SignaturePolicyPolicyHashValid indicates if the SignaturePolicyIdentifier
// digest shall match the extracted policy content. Ports
// EtsiValidationPolicy#getSignaturePolicyPolicyHashValid.
func (p *EtsiValidationPolicy) SignaturePolicyPolicyHashValid(context enumerations.Context) modelpolicy.LevelRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toLevelRule(sc.PolicyHashMatch)
	}
	return nil
}

// SignatureFormatConstraint returns the SignatureFormat constraint, if
// present. Ports EtsiValidationPolicy#getSignatureFormatConstraint.
func (p *EtsiValidationPolicy) SignatureFormatConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toMultiValuesRule(sc.AcceptableFormats)
	}
	return nil
}

// SignerInformationStoreConstraint checks if only one SignerInfo is present
// in a SignerInformationStore (PAdES only). Ports
// EtsiValidationPolicy#getSignerInformationStoreConstraint.
func (p *EtsiValidationPolicy) SignerInformationStoreConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.SignerInformationStore)
	}
	return nil
}

// ByteRangeConstraint checks if the ByteRange dictionary is valid (PAdES
// only). Ports EtsiValidationPolicy#getByteRangeConstraint.
func (p *EtsiValidationPolicy) ByteRangeConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ByteRange)
	}
	return nil
}

// ByteRangeCollisionConstraint checks if ByteRange does not collide with
// other signature byte ranges (PAdES only). Ports
// EtsiValidationPolicy#getByteRangeCollisionConstraint.
func (p *EtsiValidationPolicy) ByteRangeCollisionConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ByteRangeCollision)
	}
	return nil
}

// ByteRangeAllDocumentConstraint checks if ByteRange is valid for all
// signatures and document timestamps (PAdES only). Ports
// EtsiValidationPolicy#getByteRangeAllDocumentConstraint.
func (p *EtsiValidationPolicy) ByteRangeAllDocumentConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ByteRangeAllDocument)
	}
	return nil
}

// PdfSignatureDictionaryConstraint checks if signature dictionary is
// consistent across PDF revisions (PAdES only). Ports
// EtsiValidationPolicy#getPdfSignatureDictionaryConstraint.
func (p *EtsiValidationPolicy) PdfSignatureDictionaryConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.PdfSignatureDictionary)
	}
	return nil
}

// PdfPageDifferenceConstraint indicates if a PDF page difference check
// should be proceeded. Ports EtsiValidationPolicy#getPdfPageDifferenceConstraint.
func (p *EtsiValidationPolicy) PdfPageDifferenceConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.PdfPageDifference)
	}
	return nil
}

// PdfAnnotationOverlapConstraint indicates if a PDF annotation overlapping
// check should be proceeded. Ports EtsiValidationPolicy#getPdfAnnotationOverlapConstraint.
func (p *EtsiValidationPolicy) PdfAnnotationOverlapConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.PdfAnnotationOverlap)
	}
	return nil
}

// PdfVisualDifferenceConstraint indicates if a PDF visual difference check
// should be proceeded. Ports EtsiValidationPolicy#getPdfVisualDifferenceConstraint.
func (p *EtsiValidationPolicy) PdfVisualDifferenceConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.PdfVisualDifference)
	}
	return nil
}

// DocMDPConstraint checks for changes against permission rules identified
// within a /DocMDP dictionary. Ports EtsiValidationPolicy#getDocMDPConstraint.
func (p *EtsiValidationPolicy) DocMDPConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.DocMDP)
	}
	return nil
}

// FieldMDPConstraint checks for changes against permission rules identified
// within a /FieldMDP dictionary. Ports EtsiValidationPolicy#getFieldMDPConstraint.
func (p *EtsiValidationPolicy) FieldMDPConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.FieldMDP)
	}
	return nil
}

// SigFieldLockConstraint checks for changes against permission rules
// identified within a /SigFieldLock dictionary. Ports
// EtsiValidationPolicy#getSigFieldLockConstraint.
func (p *EtsiValidationPolicy) SigFieldLockConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.SigFieldLock)
	}
	return nil
}

// FormFillChangesConstraint checks whether a PDF document contains form fill
// or signing modifications after the current signature's revisions. Ports
// EtsiValidationPolicy#getFormFillChangesConstraint.
func (p *EtsiValidationPolicy) FormFillChangesConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.FormFillChanges)
	}
	return nil
}

// AnnotationChangesConstraint checks whether a PDF document contains
// annotation modifications after the current signature's revisions. Ports
// EtsiValidationPolicy#getAnnotationChangesConstraint.
func (p *EtsiValidationPolicy) AnnotationChangesConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.AnnotationChanges)
	}
	return nil
}

// UndefinedChangesConstraint checks whether a PDF document contains
// undefined object modifications after the current signature's revisions.
// Ports EtsiValidationPolicy#getUndefinedChangesConstraint.
func (p *EtsiValidationPolicy) UndefinedChangesConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.UndefinedChanges)
	}
	return nil
}

// StructuralValidationConstraint indicates if the structural validation
// should be checked. Ports EtsiValidationPolicy#getStructuralValidationConstraint.
func (p *EtsiValidationPolicy) StructuralValidationConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sc := p.signatureConstraintsByContext(context); sc != nil {
		return toLevelRule(sc.StructuralValidation)
	}
	return nil
}

// SigningCertificateRefersCertificateChainConstraint indicates if the
// Signing Certificate attribute should be checked against the certificate
// chain. Ports EtsiValidationPolicy#getSigningCertificateRefersCertificateChainConstraint.
func (p *EtsiValidationPolicy) SigningCertificateRefersCertificateChainConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SigningCertificateRefersCertificateChain)
	}
	return nil
}

// ReferencesToAllCertificateChainPresentConstraint indicates if the whole
// certificate chain is covered by the Signing Certificate attribute. Ports
// EtsiValidationPolicy#getReferencesToAllCertificateChainPresentConstraint.
func (p *EtsiValidationPolicy) ReferencesToAllCertificateChainPresentConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.ReferencesToAllCertificateChainPresent)
	}
	return nil
}

// SigningCertificateDigestAlgorithmConstraint checks the DigestAlgorithm
// used in signing-certificate-reference creation. Ports
// EtsiValidationPolicy#getSigningCertificateDigestAlgorithmConstraint.
func (p *EtsiValidationPolicy) SigningCertificateDigestAlgorithmConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SigningCertificateDigestAlgorithm)
	}
	return nil
}

// SigningTimeConstraint indicates if the signed property signing-time
// should be checked. Ports EtsiValidationPolicy#getSigningTimeConstraint.
func (p *EtsiValidationPolicy) SigningTimeConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SigningTime)
	}
	return nil
}

// SigningTimeInCertRangeConstraint indicates if signing-time should be
// checked against the signing-certificate's validity period. Ports
// EtsiValidationPolicy#getSigningTimeInCertRangeConstraint.
func (p *EtsiValidationPolicy) SigningTimeInCertRangeConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SigningTimeInCertRange)
	}
	return nil
}

// SignatureTypeConstraint indicates if the signed property signature type
// should be checked. Ports EtsiValidationPolicy#getSignatureTypeConstraint.
func (p *EtsiValidationPolicy) SignatureTypeConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.SignatureType)
	}
	return nil
}

// ContentTypeConstraint indicates if the signed property content-type
// should be checked. Ports EtsiValidationPolicy#getContentTypeConstraint.
func (p *EtsiValidationPolicy) ContentTypeConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.ContentType)
	}
	return nil
}

// ContentHintsConstraint indicates if the signed property content-hints
// should be checked. Ports EtsiValidationPolicy#getContentHintsConstraint.
func (p *EtsiValidationPolicy) ContentHintsConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.ContentHints)
	}
	return nil
}

// ContentIdentifierConstraint indicates if the signed property
// content-identifier should be checked. Ports
// EtsiValidationPolicy#getContentIdentifierConstraint.
func (p *EtsiValidationPolicy) ContentIdentifierConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.ContentIdentifier)
	}
	return nil
}

// MessageDigestOrSignedPropertiesConstraint indicates if message-digest
// (CAdES) or SignedProperties (XAdES) should be checked. Ports
// EtsiValidationPolicy#getMessageDigestOrSignedPropertiesConstraint.
func (p *EtsiValidationPolicy) MessageDigestOrSignedPropertiesConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.MessageDigestOrSignedPropertiesPresent)
	}
	return nil
}

// EllipticCurveKeySizeConstraint checks whether a JWA signature has a valid
// elliptic curve key size. Ports EtsiValidationPolicy#getEllipticCurveKeySizeConstraint.
func (p *EtsiValidationPolicy) EllipticCurveKeySizeConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.EllipticCurveKeySize)
	}
	return nil
}

// CommitmentTypeIndicationConstraint indicates if the signed property
// commitment-type-indication should be checked. Ports
// EtsiValidationPolicy#getCommitmentTypeIndicationConstraint.
func (p *EtsiValidationPolicy) CommitmentTypeIndicationConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.CommitmentTypeIndication)
	}
	return nil
}

// SignerLocationConstraint indicates if the signed property signer-location
// should be checked. Ports EtsiValidationPolicy#getSignerLocationConstraint.
func (p *EtsiValidationPolicy) SignerLocationConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SignerLocation)
	}
	return nil
}

// ContentTimeStampConstraint indicates if the signed property
// content-time-stamp should be checked. Ports
// EtsiValidationPolicy#getContentTimeStampConstraint.
func (p *EtsiValidationPolicy) ContentTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.ContentTimeStamp)
	}
	return nil
}

// ContentTimeStampMessageImprintConstraint indicates if the
// content-time-stamp message-imprint should be checked. Ports
// EtsiValidationPolicy#getContentTimeStampMessageImprintConstraint.
func (p *EtsiValidationPolicy) ContentTimeStampMessageImprintConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.ContentTimeStampMessageImprint)
	}
	return nil
}

// ClaimedRoleConstraint indicates if the unsigned property claimed-role
// should be checked. Ports EtsiValidationPolicy#getClaimedRoleConstraint.
func (p *EtsiValidationPolicy) ClaimedRoleConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.ClaimedRoles)
	}
	return nil
}

// CertifiedRolesConstraint returns the mandated signer role. Ports
// EtsiValidationPolicy#getCertifiedRolesConstraint.
func (p *EtsiValidationPolicy) CertifiedRolesConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toMultiValuesRule(sac.CertifiedRoles)
	}
	return nil
}

// SignatureCryptographicConstraint creates the CryptographicSuite
// corresponding to the context parameter. Ports
// EtsiValidationPolicy#getSignatureCryptographicConstraint.
func (p *EtsiValidationPolicy) SignatureCryptographicConstraint(context enumerations.Context) modelpolicy.CryptographicSuite {
	return toCryptographicSuite(p.signatureCryptographic(context))
}

// signatureCryptographic ports the private getSignatureCryptographic(Context)
// helper.
func (p *EtsiValidationPolicy) signatureCryptographic(context enumerations.Context) *jaxb.CryptographicConstraint {
	sigCryptographic := &jaxb.CryptographicConstraint{}
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil && bsc.Cryptographic != nil {
		sigCryptographic = bsc.Cryptographic
	}
	initializeCryptographicSuite(sigCryptographic, p.Cryptographic())
	return sigCryptographic
}

// CertificateCryptographicConstraint creates the CryptographicSuite
// corresponding to the context/subContext parameters. Ports
// EtsiValidationPolicy#getCertificateCryptographicConstraint.
func (p *EtsiValidationPolicy) CertificateCryptographicConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.CryptographicSuite {
	certCryptographic := &jaxb.CryptographicConstraint{}
	if cc := p.certificateConstraints(context, subContext); cc != nil && cc.Cryptographic != nil {
		certCryptographic = cc.Cryptographic
	}
	initializeCryptographicSuite(certCryptographic, p.signatureCryptographic(context))
	return toCryptographicSuite(certCryptographic)
}

// initializeCryptographicSuite overrides all empty fields for the given
// CryptographicConstraint by the default CryptographicConstraint. Ports the
// private initializeCryptographicSuite(CryptographicConstraint,
// CryptographicConstraint) helper.
func initializeCryptographicSuite(cryptographicConstraint, defaultConstraint *jaxb.CryptographicConstraint) {
	if defaultConstraint != nil {
		if cryptographicConstraint.AcceptableDigestAlgo == nil {
			cryptographicConstraint.AcceptableDigestAlgo = defaultConstraint.AcceptableDigestAlgo
		}
		if cryptographicConstraint.AcceptableEncryptionAlgo == nil {
			cryptographicConstraint.AcceptableEncryptionAlgo = defaultConstraint.AcceptableEncryptionAlgo
		}
		if cryptographicConstraint.AlgoExpirationDate == nil {
			cryptographicConstraint.AlgoExpirationDate = defaultConstraint.AlgoExpirationDate
		}
		if cryptographicConstraint.Level.Level() == "" {
			cryptographicConstraint.Level = defaultConstraint.Level
		}
		if cryptographicConstraint.MiniPublicKeySize == nil {
			cryptographicConstraint.MiniPublicKeySize = defaultConstraint.MiniPublicKeySize
		}
	}
}

// CertificateCAConstraint returns certificate CA constraint. Ports
// EtsiValidationPolicy#getCertificateCAConstraint.
func (p *EtsiValidationPolicy) CertificateCAConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.CA)
	}
	return nil
}

// CertificateMaxPathLengthConstraint returns certificate MaxPathLength
// constraint. Ports EtsiValidationPolicy#getCertificateMaxPathLengthConstraint.
func (p *EtsiValidationPolicy) CertificateMaxPathLengthConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.MaxPathLength)
	}
	return nil
}

// CertificateKeyUsageConstraint returns certificate key usage constraint.
// Ports EtsiValidationPolicy#getCertificateKeyUsageConstraint.
func (p *EtsiValidationPolicy) CertificateKeyUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.KeyUsage)
	}
	return nil
}

// CertificateExtendedKeyUsageConstraint returns certificate extended key
// usage constraint. Ports EtsiValidationPolicy#getCertificateExtendedKeyUsageConstraint.
func (p *EtsiValidationPolicy) CertificateExtendedKeyUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.ExtendedKeyUsage)
	}
	return nil
}

// CertificatePolicyTreeConstraint returns certificate PolicyTree constraint.
// Ports EtsiValidationPolicy#getCertificatePolicyTreeConstraint.
func (p *EtsiValidationPolicy) CertificatePolicyTreeConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.PolicyTree)
	}
	return nil
}

// CertificateNameConstraintsConstraint returns certificate NameConstraints
// constraint. Ports EtsiValidationPolicy#getCertificateNameConstraintsConstraint.
func (p *EtsiValidationPolicy) CertificateNameConstraintsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NameConstraints)
	}
	return nil
}

// CertificateAuthorityKeyIdentifierPresentConstraint returns certificate
// AuthorityKeyIdentifierPresent constraint. Ports
// EtsiValidationPolicy#getCertificateAuthorityKeyIdentifierPresentConstraint.
func (p *EtsiValidationPolicy) CertificateAuthorityKeyIdentifierPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.AuthorityKeyIdentifierPresent)
	}
	return nil
}

// CertificateSubjectKeyIdentifierPresentConstraint returns certificate
// SubjectKeyIdentifierPresent constraint. Ports
// EtsiValidationPolicy#getCertificateSubjectKeyIdentifierPresentConstraint.
func (p *EtsiValidationPolicy) CertificateSubjectKeyIdentifierPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.SubjectKeyIdentifierPresent)
	}
	return nil
}

// CertificateNoRevAvailConstraint returns certificate NoRevAvail constraint.
// Ports EtsiValidationPolicy#getCertificateNoRevAvailConstraint.
func (p *EtsiValidationPolicy) CertificateNoRevAvailConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NoRevAvail)
	}
	return nil
}

// CertificateSupportedCriticalExtensionsConstraint returns certificate
// supported critical extensions constraint. Ports
// EtsiValidationPolicy#getCertificateSupportedCriticalExtensionsConstraint.
func (p *EtsiValidationPolicy) CertificateSupportedCriticalExtensionsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.SupportedCriticalExtensions)
	}
	return nil
}

// CertificateForbiddenExtensionsConstraint returns certificate forbidden
// extensions constraint. Ports
// EtsiValidationPolicy#getCertificateForbiddenExtensionsConstraint.
func (p *EtsiValidationPolicy) CertificateForbiddenExtensionsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.ForbiddenExtensions)
	}
	return nil
}

// CertificateNotExpiredConstraint returns certificate's validity range
// constraint. Ports EtsiValidationPolicy#getCertificateNotExpiredConstraint.
func (p *EtsiValidationPolicy) CertificateNotExpiredConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NotExpired)
	}
	return nil
}

// CertificateSunsetDateConstraint returns certificate's sunset date
// constraint. Ports EtsiValidationPolicy#getCertificateSunsetDateConstraint.
func (p *EtsiValidationPolicy) CertificateSunsetDateConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.SunsetDate)
	}
	return nil
}

// ProspectiveCertificateChainConstraint requests the presence of the trust
// anchor in the certificate chain. Ports
// EtsiValidationPolicy#getProspectiveCertificateChainConstraint.
func (p *EtsiValidationPolicy) ProspectiveCertificateChainConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ProspectiveCertificateChain)
	}
	return nil
}

// CertificateSignatureConstraint returns certificate's signature
// constraint. Ports EtsiValidationPolicy#getCertificateSignatureConstraint.
func (p *EtsiValidationPolicy) CertificateSignatureConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.Signature)
	}
	return nil
}

// UnknownStatusConstraint returns the UnknownStatus constraint. Ports
// EtsiValidationPolicy#getUnknownStatusConstraint.
func (p *EtsiValidationPolicy) UnknownStatusConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.UnknownStatus)
	}
	return nil
}

// ThisUpdatePresentConstraint returns the ThisUpdatePresent constraint.
// Ports EtsiValidationPolicy#getThisUpdatePresentConstraint.
func (p *EtsiValidationPolicy) ThisUpdatePresentConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.ThisUpdatePresent)
	}
	return nil
}

// RevocationIssuerKnownConstraint returns the RevocationIssuerKnown
// constraint. Ports EtsiValidationPolicy#getRevocationIssuerKnownConstraint.
func (p *EtsiValidationPolicy) RevocationIssuerKnownConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.RevocationIssuerKnown)
	}
	return nil
}

// RevocationIssuerValidAtProductionTimeConstraint returns the
// RevocationIssuerValidAtProductionTime constraint. Ports
// EtsiValidationPolicy#getRevocationIssuerValidAtProductionTimeConstraint.
func (p *EtsiValidationPolicy) RevocationIssuerValidAtProductionTimeConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.RevocationIssuerValidAtProductionTime)
	}
	return nil
}

// RevocationAfterCertificateIssuanceConstraint returns the
// RevocationIssuerKnowsCertificate constraint. Ports
// EtsiValidationPolicy#getRevocationAfterCertificateIssuanceConstraint.
func (p *EtsiValidationPolicy) RevocationAfterCertificateIssuanceConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.RevocationAfterCertificateIssuance)
	}
	return nil
}

// RevocationHasInformationAboutCertificateConstraint returns the
// RevocationIssuerHasInformationAboutCertificate constraint. Ports
// EtsiValidationPolicy#getRevocationHasInformationAboutCertificateConstraint.
func (p *EtsiValidationPolicy) RevocationHasInformationAboutCertificateConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.RevocationHasInformationAboutCertificate)
	}
	return nil
}

// OCSPResponseResponderIdMatchConstraint returns the OCSPResponderIdMatch
// constraint. Ports EtsiValidationPolicy#getOCSPResponseResponderIdMatchConstraint.
func (p *EtsiValidationPolicy) OCSPResponseResponderIdMatchConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.OCSPResponderIdMatch)
	}
	return nil
}

// OCSPResponseCertHashPresentConstraint returns the OCSPCertHashPresent
// constraint. Ports EtsiValidationPolicy#getOCSPResponseCertHashPresentConstraint.
func (p *EtsiValidationPolicy) OCSPResponseCertHashPresentConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.OCSPCertHashPresent)
	}
	return nil
}

// OCSPResponseCertHashMatchConstraint returns the OCSPCertHashMatch
// constraint. Ports EtsiValidationPolicy#getOCSPResponseCertHashMatchConstraint.
func (p *EtsiValidationPolicy) OCSPResponseCertHashMatchConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.OCSPCertHashMatch)
	}
	return nil
}

// SelfIssuedOCSPConstraint returns the SelfIssuedOCSP constraint. Ports
// EtsiValidationPolicy#getSelfIssuedOCSPConstraint.
func (p *EtsiValidationPolicy) SelfIssuedOCSPConstraint() modelpolicy.LevelRule {
	if rc := p.RevocationConstraints(); rc != nil {
		return toLevelRule(rc.SelfIssuedOCSP)
	}
	return nil
}

// RevocationDataAvailableConstraint returns revocation data available
// constraint. Ports EtsiValidationPolicy#getRevocationDataAvailableConstraint.
func (p *EtsiValidationPolicy) RevocationDataAvailableConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.RevocationDataAvailable)
	}
	return nil
}

// AcceptableRevocationDataFoundConstraint returns acceptable revocation data
// available constraint. Ports
// EtsiValidationPolicy#getAcceptableRevocationDataFoundConstraint.
func (p *EtsiValidationPolicy) AcceptableRevocationDataFoundConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.AcceptableRevocationDataFound)
	}
	return nil
}

// CRLNextUpdatePresentConstraint returns CRL's nextUpdate present
// constraint. Ports EtsiValidationPolicy#getCRLNextUpdatePresentConstraint.
func (p *EtsiValidationPolicy) CRLNextUpdatePresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.CRLNextUpdatePresent)
	}
	return nil
}

// OCSPNextUpdatePresentConstraint returns OCSP's nextUpdate present
// constraint. Ports EtsiValidationPolicy#getOCSPNextUpdatePresentConstraint.
func (p *EtsiValidationPolicy) OCSPNextUpdatePresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.OCSPNextUpdatePresent)
	}
	return nil
}

// RevocationFreshnessConstraint returns revocation data's freshness
// constraint. Ports EtsiValidationPolicy#getRevocationFreshnessConstraint.
func (p *EtsiValidationPolicy) RevocationFreshnessConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.DurationRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toDurationRule(cc.RevocationFreshness)
	}
	return nil
}

// RevocationFreshnessNextUpdateConstraint returns revocation data's
// freshness for nextUpdate check constraint. Ports
// EtsiValidationPolicy#getRevocationFreshnessNextUpdateConstraint.
func (p *EtsiValidationPolicy) RevocationFreshnessNextUpdateConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.RevocationFreshnessNextUpdate)
	}
	return nil
}

// CertificateNotRevokedConstraint returns certificate's not revoked
// constraint. Ports EtsiValidationPolicy#getCertificateNotRevokedConstraint.
func (p *EtsiValidationPolicy) CertificateNotRevokedConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NotRevoked)
	}
	return nil
}

// CertificateNotOnHoldConstraint returns certificate's not onHold
// constraint. Ports EtsiValidationPolicy#getCertificateNotOnHoldConstraint.
func (p *EtsiValidationPolicy) CertificateNotOnHoldConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NotOnHold)
	}
	return nil
}

// RevocationIssuerNotExpiredConstraint returns revocation issuer's validity
// range constraint. Ports EtsiValidationPolicy#getRevocationIssuerNotExpiredConstraint.
func (p *EtsiValidationPolicy) RevocationIssuerNotExpiredConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.RevocationIssuerNotExpired)
	}
	return nil
}

// CertificateNotSelfSignedConstraint returns certificate's not self-signed
// constraint. Ports EtsiValidationPolicy#getCertificateNotSelfSignedConstraint.
func (p *EtsiValidationPolicy) CertificateNotSelfSignedConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.NotSelfSigned)
	}
	return nil
}

// CertificateSelfSignedConstraint returns certificate's self-signed
// constraint. Ports EtsiValidationPolicy#getCertificateSelfSignedConstraint.
func (p *EtsiValidationPolicy) CertificateSelfSignedConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.SelfSigned)
	}
	return nil
}

// TrustServiceTypeIdentifierConstraint returns trusted service type
// identifier constraint. Ports
// EtsiValidationPolicy#getTrustServiceTypeIdentifierConstraint.
func (p *EtsiValidationPolicy) TrustServiceTypeIdentifierConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toMultiValuesRule(bsc.TrustServiceTypeIdentifier)
	}
	return nil
}

// TrustServiceStatusConstraint returns trusted service status constraint.
// Ports EtsiValidationPolicy#getTrustServiceStatusConstraint.
func (p *EtsiValidationPolicy) TrustServiceStatusConstraint(context enumerations.Context) modelpolicy.MultiValuesRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toMultiValuesRule(bsc.TrustServiceStatus)
	}
	return nil
}

// CertificatePolicyIdsConstraint returns the CertificatePolicyIds
// constraint, if present. Ports EtsiValidationPolicy#getCertificatePolicyIdsConstraint.
func (p *EtsiValidationPolicy) CertificatePolicyIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.PolicyIds)
	}
	return nil
}

// CertificatePolicyQualificationIdsConstraint indicates if the
// CertificatePolicyIds declare the certificate as qualified. Ports
// EtsiValidationPolicy#getCertificatePolicyQualificationIdsConstraint.
func (p *EtsiValidationPolicy) CertificatePolicyQualificationIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.PolicyQualificationIds)
	}
	return nil
}

// CertificatePolicySupportedByQSCDIdsConstraint indicates if the
// CertificatePolicyIds mandate QSCD support. Ports
// EtsiValidationPolicy#getCertificatePolicySupportedByQSCDIdsConstraint.
func (p *EtsiValidationPolicy) CertificatePolicySupportedByQSCDIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.PolicySupportedByQSCDIds)
	}
	return nil
}

// CertificateQCComplianceConstraint indicates if the end user certificate
// is QC Compliant. Ports EtsiValidationPolicy#getCertificateQCComplianceConstraint.
func (p *EtsiValidationPolicy) CertificateQCComplianceConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.QcCompliance)
	}
	return nil
}

// CertificateQcEuLimitValueCurrencyConstraint indicates the allowed
// currency for the QCLimitValue statement. Ports
// EtsiValidationPolicy#getCertificateQcEuLimitValueCurrencyConstraint.
func (p *EtsiValidationPolicy) CertificateQcEuLimitValueCurrencyConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.ValueRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toValueRule(cc.QcEuLimitValueCurrency)
	}
	return nil
}

// CertificateMinQcEuLimitValueConstraint indicates the minimal allowed
// QcEuLimitValue transaction limit. Ports
// EtsiValidationPolicy#getCertificateMinQcEuLimitValueConstraint.
func (p *EtsiValidationPolicy) CertificateMinQcEuLimitValueConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.NumericValueRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toNumericValueRule(cc.MinQcEuLimitValue)
	}
	return nil
}

// CertificateMinQcEuRetentionPeriodConstraint indicates the minimal allowed
// QC retention period. Ports
// EtsiValidationPolicy#getCertificateMinQcEuRetentionPeriodConstraint.
func (p *EtsiValidationPolicy) CertificateMinQcEuRetentionPeriodConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.NumericValueRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toNumericValueRule(cc.MinQcEuRetentionPeriod)
	}
	return nil
}

// CertificateQcSSCDConstraint indicates if the end user certificate is
// mandated to be QSCD-supported. Ports EtsiValidationPolicy#getCertificateQcSSCDConstraint.
func (p *EtsiValidationPolicy) CertificateQcSSCDConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.QcSSCD)
	}
	return nil
}

// CertificateQcEuPDSLocationConstraint indicates the location(s) of PKI
// Disclosure Statements. Ports EtsiValidationPolicy#getCertificateQcEuPDSLocationConstraint.
func (p *EtsiValidationPolicy) CertificateQcEuPDSLocationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcEuPDSLocation)
	}
	return nil
}

// CertificateQcTypeConstraint indicates the claimed certificate type(s).
// Ports EtsiValidationPolicy#getCertificateQcTypeConstraint.
func (p *EtsiValidationPolicy) CertificateQcTypeConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcType)
	}
	return nil
}

// CertificateQcCCLegislationConstraint indicates the country/countries
// under whose legislation the certificate is issued as qualified. Ports
// EtsiValidationPolicy#getCertificateQcCCLegislationConstraint.
func (p *EtsiValidationPolicy) CertificateQcCCLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcLegislationCountryCodes)
	}
	return nil
}

// CertificateIssuedToNaturalPersonConstraint indicates if the end user
// certificate is issued to a natural person. Ports
// EtsiValidationPolicy#getCertificateIssuedToNaturalPersonConstraint.
func (p *EtsiValidationPolicy) CertificateIssuedToNaturalPersonConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.IssuedToNaturalPerson)
	}
	return nil
}

// CertificateIssuedToLegalPersonConstraint indicates if the end user
// certificate is issued to a legal person. Ports
// EtsiValidationPolicy#getCertificateIssuedToLegalPersonConstraint.
func (p *EtsiValidationPolicy) CertificateIssuedToLegalPersonConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.IssuedToLegalPerson)
	}
	return nil
}

// CertificateSemanticsIdentifierConstraint indicates the acceptable
// QCStatement semantics identifier. Ports
// EtsiValidationPolicy#getCertificateSemanticsIdentifierConstraint.
func (p *EtsiValidationPolicy) CertificateSemanticsIdentifierConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.SemanticsIdentifier)
	}
	return nil
}

// CertificatePS2DQcTypeRolesOfPSPConstraint indicates the acceptable QC
// PS2D roles. Ports EtsiValidationPolicy#getCertificatePS2DQcTypeRolesOfPSPConstraint.
func (p *EtsiValidationPolicy) CertificatePS2DQcTypeRolesOfPSPConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.PSD2QcTypeRolesOfPSP)
	}
	return nil
}

// CertificatePS2DQcCompetentAuthorityNameConstraint indicates the
// acceptable QC PS2D names. Ports
// EtsiValidationPolicy#getCertificatePS2DQcCompetentAuthorityNameConstraint.
func (p *EtsiValidationPolicy) CertificatePS2DQcCompetentAuthorityNameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.PSD2QcCompetentAuthorityName)
	}
	return nil
}

// CertificatePS2DQcCompetentAuthorityIdConstraint indicates the acceptable
// QC PS2D ids. Ports EtsiValidationPolicy#getCertificatePS2DQcCompetentAuthorityIdConstraint.
func (p *EtsiValidationPolicy) CertificatePS2DQcCompetentAuthorityIdConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.PSD2QcCompetentAuthorityId)
	}
	return nil
}

// CertificateQcQSCDLegislationConstraint indicates the country/countries
// under whose legislation the signature creation device has a qualified
// status. Ports EtsiValidationPolicy#getCertificateQcQSCDLegislationConstraint.
func (p *EtsiValidationPolicy) CertificateQcQSCDLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcQSCDLegislation)
	}
	return nil
}

// CertificateQcIdentificationMethodConstraint indicates the verification
// method used for issuance per eIDAS Article 24. Ports
// EtsiValidationPolicy#getCertificateQcIdentificationMethodConstraint.
func (p *EtsiValidationPolicy) CertificateQcIdentificationMethodConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcIdentificationMethod)
	}
	return nil
}

// CertificateQcPSBCountryOfLegislationConstraint indicates the verification
// method for country of legislation of the QcPSB QcStatement. Ports
// EtsiValidationPolicy#getCertificateQcPSBCountryOfLegislationConstraint.
func (p *EtsiValidationPolicy) CertificateQcPSBCountryOfLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcPSBCountryOfLegislation)
	}
	return nil
}

// CertificateQcPSBAuthSourceIdentificationConstraint indicates the
// verification method for authentic source identification of the QcPSB
// QcStatement. Ports EtsiValidationPolicy#getCertificateQcPSBAuthSourceIdentificationConstraint.
func (p *EtsiValidationPolicy) CertificateQcPSBAuthSourceIdentificationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcPSBAuthSourceIdentification)
	}
	return nil
}

// CertificateQcPSBLegislationIdentificationConstraint indicates the
// verification method for legislation identification of the QcPSB
// QcStatement. Ports EtsiValidationPolicy#getCertificateQcPSBLegislationIdentificationConstraint.
func (p *EtsiValidationPolicy) CertificateQcPSBLegislationIdentificationConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.QcPSBLegislationIdentification)
	}
	return nil
}

// SigningCertificateRecognitionConstraint indicates if signing-certificate
// has been identified. Ports EtsiValidationPolicy#getSigningCertificateRecognitionConstraint.
func (p *EtsiValidationPolicy) SigningCertificateRecognitionConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if cc := p.signingCertificateByContext(context); cc != nil {
		return toLevelRule(cc.Recognition)
	}
	return nil
}

// SigningCertificateAttributePresentConstraint indicates if the signing
// certificate attribute is present. Ports
// EtsiValidationPolicy#getSigningCertificateAttributePresentConstraint.
func (p *EtsiValidationPolicy) SigningCertificateAttributePresentConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.SigningCertificatePresent)
	}
	return nil
}

// UnicitySigningCertificateAttributeConstraint indicates if the signing
// certificate is not ambiguously determined. Ports
// EtsiValidationPolicy#getUnicitySigningCertificateAttributeConstraint.
func (p *EtsiValidationPolicy) UnicitySigningCertificateAttributeConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.UnicitySigningCertificate)
	}
	return nil
}

// SigningCertificateDigestValuePresentConstraint indicates if the signing
// certificate reference's digest value is present. Ports
// EtsiValidationPolicy#getSigningCertificateDigestValuePresentConstraint.
func (p *EtsiValidationPolicy) SigningCertificateDigestValuePresentConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.CertDigestPresent)
	}
	return nil
}

// SigningCertificateDigestValueMatchConstraint indicates if the signing
// certificate reference's digest value matches. Ports
// EtsiValidationPolicy#getSigningCertificateDigestValueMatchConstraint.
func (p *EtsiValidationPolicy) SigningCertificateDigestValueMatchConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.CertDigestMatch)
	}
	return nil
}

// SigningCertificateIssuerSerialMatchConstraint indicates if the signing
// certificate reference's issuer serial matches. Ports
// EtsiValidationPolicy#getSigningCertificateIssuerSerialMatchConstraint.
func (p *EtsiValidationPolicy) SigningCertificateIssuerSerialMatchConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.IssuerSerialMatch)
	}
	return nil
}

// KeyIdentifierPresent indicates if the 'kid' header parameter is present.
// Ports EtsiValidationPolicy#getKeyIdentifierPresent.
func (p *EtsiValidationPolicy) KeyIdentifierPresent(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.KeyIdentifierPresent)
	}
	return nil
}

// KeyIdentifierMatch indicates if the value of 'kid' matches the
// signing-certificate. Ports EtsiValidationPolicy#getKeyIdentifierMatch.
func (p *EtsiValidationPolicy) KeyIdentifierMatch(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.KeyIdentifierMatch)
	}
	return nil
}

// X509UrlPresent indicates if the value of 'x5u' header parameter is
// present. Ports EtsiValidationPolicy#getX509UrlPresent.
func (p *EtsiValidationPolicy) X509UrlPresent(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.X509UrlPresent)
	}
	return nil
}

// X509UrlMatch indicates if the signing-certificate can be derived from the
// 'x5u' header parameter. Ports EtsiValidationPolicy#getX509UrlMatch.
func (p *EtsiValidationPolicy) X509UrlMatch(context enumerations.Context) modelpolicy.LevelRule {
	if sac := p.signedAttributeConstraints(context); sac != nil {
		return toLevelRule(sac.X509UrlMatch)
	}
	return nil
}

// ReferenceDataExistenceConstraint indicates if the referenced data is
// found. Ports EtsiValidationPolicy#getReferenceDataExistenceConstraint.
func (p *EtsiValidationPolicy) ReferenceDataExistenceConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ReferenceDataExistence)
	}
	return nil
}

// ReferenceDataIntactConstraint indicates if the referenced data is
// intact. Ports EtsiValidationPolicy#getReferenceDataIntactConstraint.
func (p *EtsiValidationPolicy) ReferenceDataIntactConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ReferenceDataIntact)
	}
	return nil
}

// ReferenceDataNameMatchConstraint indicates if the referenced document
// names match the manifest entry references. Ports
// EtsiValidationPolicy#getReferenceDataNameMatchConstraint.
func (p *EtsiValidationPolicy) ReferenceDataNameMatchConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ReferenceDataNameMatch)
	}
	return nil
}

// ManifestEntryObjectExistenceConstraint indicates if the manifested
// document is found. Ports EtsiValidationPolicy#getManifestEntryObjectExistenceConstraint.
func (p *EtsiValidationPolicy) ManifestEntryObjectExistenceConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ManifestEntryObjectExistence)
	}
	return nil
}

// ManifestEntryObjectIntactConstraint indicates if the manifested document
// is intact. Ports EtsiValidationPolicy#getManifestEntryObjectIntactConstraint.
func (p *EtsiValidationPolicy) ManifestEntryObjectIntactConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ManifestEntryObjectIntact)
	}
	return nil
}

// ManifestEntryObjectGroupConstraint indicates if all manifest entries have
// been found. Ports EtsiValidationPolicy#getManifestEntryObjectGroupConstraint.
func (p *EtsiValidationPolicy) ManifestEntryObjectGroupConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ManifestEntryObjectGroup)
	}
	return nil
}

// ManifestEntryNameMatchConstraint indicates if names of all matching
// documents match the manifest entry names. Ports
// EtsiValidationPolicy#getManifestEntryNameMatchConstraint.
func (p *EtsiValidationPolicy) ManifestEntryNameMatchConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.ManifestEntryNameMatch)
	}
	return nil
}

// SignatureIntactConstraint indicates if the signature is intact. Ports
// EtsiValidationPolicy#getSignatureIntactConstraint.
func (p *EtsiValidationPolicy) SignatureIntactConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.SignatureIntact)
	}
	return nil
}

// SignatureDuplicatedConstraint indicates if the signature is not
// ambiguous. Ports EtsiValidationPolicy#getSignatureDuplicatedConstraint.
func (p *EtsiValidationPolicy) SignatureDuplicatedConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if bsc := p.basicSignatureConstraintsByContext(context); bsc != nil {
		return toLevelRule(bsc.SignatureDuplicated)
	}
	return nil
}

// BestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint
// checks if the certificate is not expired on best-signature-time. Ports
// EtsiValidationPolicy#getBestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint.
func (p *EtsiValidationPolicy) BestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.BestSignatureTimeBeforeExpirationDateOfSigningCertificate)
	}
	return nil
}

// TimestampCoherenceConstraint checks if the timestamp order is coherent.
// Ports EtsiValidationPolicy#getTimestampCoherenceConstraint.
func (p *EtsiValidationPolicy) TimestampCoherenceConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.Coherence)
	}
	return nil
}

// TimestampDelayConstraint returns the TimestampDelay constraint, if
// present. Ports EtsiValidationPolicy#getTimestampDelayConstraint.
func (p *EtsiValidationPolicy) TimestampDelayConstraint() modelpolicy.DurationRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toDurationRule(ts.TimestampDelay)
	}
	return nil
}

// TimestampValidConstraint returns whether the time-stamp is valid. Ports
// EtsiValidationPolicy#getTimestampValidConstraint.
func (p *EtsiValidationPolicy) TimestampValidConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.TimestampValid)
	}
	return nil
}

// TimestampTSAGeneralNamePresent indicates if the timestamp's
// TSTInfo.tsa field is present. Ports
// EtsiValidationPolicy#getTimestampTSAGeneralNamePresent.
func (p *EtsiValidationPolicy) TimestampTSAGeneralNamePresent() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.TSAGeneralNamePresent)
	}
	return nil
}

// TimestampTSAGeneralNameContentMatch indicates if TSTInfo.tsa matches the
// timestamp's issuer distinguishing name. Ports
// EtsiValidationPolicy#getTimestampTSAGeneralNameContentMatch.
func (p *EtsiValidationPolicy) TimestampTSAGeneralNameContentMatch() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.TSAGeneralNameContentMatch)
	}
	return nil
}

// TimestampTSAGeneralNameOrderMatch indicates if TSTInfo.tsa value and
// order match the timestamp's issuer distinguishing name. Ports
// EtsiValidationPolicy#getTimestampTSAGeneralNameOrderMatch.
func (p *EtsiValidationPolicy) TimestampTSAGeneralNameOrderMatch() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.TSAGeneralNameOrderMatch)
	}
	return nil
}

// AtsHashIndexConstraint returns the timestamp AtsHashIndex constraint, if
// present. Ports EtsiValidationPolicy#getAtsHashIndexConstraint.
func (p *EtsiValidationPolicy) AtsHashIndexConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.AtsHashIndex)
	}
	return nil
}

// TimestampContainerSignedAndTimestampedFilesCoveredConstraint returns the
// timestamp ContainerSignedAndTimestampedFilesCovered constraint, if
// present. Ports EtsiValidationPolicy#getTimestampContainerSignedAndTimestampedFilesCoveredConstraint.
func (p *EtsiValidationPolicy) TimestampContainerSignedAndTimestampedFilesCoveredConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.ContainerSignedAndTimestampedFilesCovered)
	}
	return nil
}

// RevocationTimeAgainstBestSignatureTimeConstraint returns the
// RevocationTimeAgainstBestSignatureTime constraint, if present. Ports
// EtsiValidationPolicy#getRevocationTimeAgainstBestSignatureTimeConstraint.
func (p *EtsiValidationPolicy) RevocationTimeAgainstBestSignatureTimeConstraint() modelpolicy.LevelRule {
	if ts := p.TimestampConstraints(); ts != nil {
		return toLevelRule(ts.RevocationTimeAgainstBestSignatureTime)
	}
	return nil
}

// EvidenceRecordValidConstraint returns whether the evidence record is
// valid. Ports EtsiValidationPolicy#getEvidenceRecordValidConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordValidConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.EvidenceRecordValid)
	}
	return nil
}

// EvidenceRecordDataObjectExistenceConstraint returns the
// DataObjectExistence constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordDataObjectExistenceConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordDataObjectExistenceConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.DataObjectExistence)
	}
	return nil
}

// EvidenceRecordDataObjectIntactConstraint returns the DataObjectIntact
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordDataObjectIntactConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordDataObjectIntactConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.DataObjectIntact)
	}
	return nil
}

// EvidenceRecordDataObjectFoundConstraint returns the DataObjectFound
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordDataObjectFoundConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordDataObjectFoundConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.DataObjectFound)
	}
	return nil
}

// EvidenceRecordDataObjectGroupConstraint returns the DataObjectGroup
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordDataObjectGroupConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordDataObjectGroupConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.DataObjectGroup)
	}
	return nil
}

// EvidenceRecordSignedFilesCoveredConstraint returns the SignedFilesCovered
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordSignedFilesCoveredConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordSignedFilesCoveredConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.SignedFilesCovered)
	}
	return nil
}

// EvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint
// returns the evidence record ContainerSignedAndTimestampedFilesCovered
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.ContainerSignedAndTimestampedFilesCovered)
	}
	return nil
}

// EvidenceRecordHashTreeRenewalConstraint returns the HashTreeRenewal
// constraint, if present. Ports
// EtsiValidationPolicy#getEvidenceRecordHashTreeRenewalConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordHashTreeRenewalConstraint() modelpolicy.LevelRule {
	if er := p.EvidenceRecordConstraints(); er != nil {
		return toLevelRule(er.HashTreeRenewal)
	}
	return nil
}

// EvidenceRecordCryptographicConstraint returns cryptographic constraints
// for validation of Evidence Record. Ports
// EtsiValidationPolicy#getEvidenceRecordCryptographicConstraint.
func (p *EtsiValidationPolicy) EvidenceRecordCryptographicConstraint() modelpolicy.CryptographicSuite {
	evidenceRecordCryptographic := &jaxb.CryptographicConstraint{}
	if er := p.EvidenceRecordConstraints(); er != nil && er.Cryptographic != nil {
		evidenceRecordCryptographic = er.Cryptographic
	}
	initializeCryptographicSuite(evidenceRecordCryptographic, p.Cryptographic())
	return toCryptographicSuite(evidenceRecordCryptographic)
}

// EAACryptographicConstraint returns cryptographic constraints for
// validation of EAA. Ports EtsiValidationPolicy#getEAACryptographicConstraint.
func (p *EtsiValidationPolicy) EAACryptographicConstraint() modelpolicy.CryptographicSuite {
	eaaPresentationCryptographic := &jaxb.CryptographicConstraint{}
	if eaa := p.EAAConstraints(); eaa != nil && eaa.Cryptographic != nil {
		eaaPresentationCryptographic = eaa.Cryptographic
	}
	initializeCryptographicSuite(eaaPresentationCryptographic, p.Cryptographic())
	return toCryptographicSuite(eaaPresentationCryptographic)
}

// CounterSignatureConstraint returns the CounterSignature constraint, if
// present. Ports EtsiValidationPolicy#getCounterSignatureConstraint.
func (p *EtsiValidationPolicy) CounterSignatureConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.CounterSignature)
	}
	return nil
}

// SignatureTimeStampConstraint indicates if the presence of the
// signature-time-stamp unsigned property should be checked. Ports
// EtsiValidationPolicy#getSignatureTimeStampConstraint.
func (p *EtsiValidationPolicy) SignatureTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.SignatureTimeStamp)
	}
	return nil
}

// ValidationDataTimeStampConstraint indicates if the presence of the
// validation data timestamp unsigned property should be checked. Ports
// EtsiValidationPolicy#getValidationDataTimeStampConstraint.
func (p *EtsiValidationPolicy) ValidationDataTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.ValidationDataTimeStamp)
	}
	return nil
}

// ValidationDataRefsOnlyTimeStampConstraint indicates if the presence of
// the validation data references only timestamp unsigned property should
// be checked. Ports EtsiValidationPolicy#getValidationDataRefsOnlyTimeStampConstraint.
func (p *EtsiValidationPolicy) ValidationDataRefsOnlyTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.ValidationDataRefsOnlyTimeStamp)
	}
	return nil
}

// ArchiveTimeStampConstraint indicates if the presence of the
// archive-time-stamp unsigned property should be checked. Ports
// EtsiValidationPolicy#getArchiveTimeStampConstraint.
func (p *EtsiValidationPolicy) ArchiveTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.ArchiveTimeStamp)
	}
	return nil
}

// DocumentTimeStampConstraint indicates if the presence of the document
// timestamp unsigned property should be checked. Ports
// EtsiValidationPolicy#getDocumentTimeStampConstraint.
func (p *EtsiValidationPolicy) DocumentTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.DocumentTimeStamp)
	}
	return nil
}

// TLevelTimeStampConstraint indicates if the presence of the
// signature-time-stamp or document timestamp should be checked. Ports
// EtsiValidationPolicy#getTLevelTimeStampConstraint.
func (p *EtsiValidationPolicy) TLevelTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.TLevelTimeStamp)
	}
	return nil
}

// LTALevelTimeStampConstraint indicates if the presence of the
// archive-time-stamp or document timestamp covering the validation data
// should be checked. Ports EtsiValidationPolicy#getLTALevelTimeStampConstraint.
func (p *EtsiValidationPolicy) LTALevelTimeStampConstraint(context enumerations.Context) modelpolicy.LevelRule {
	if uac := p.unsignedAttributeConstraints(context); uac != nil {
		return toLevelRule(uac.LTALevelTimeStamp)
	}
	return nil
}

// CertificateCountryConstraint returns the CertificateCountry constraint,
// if present. Ports EtsiValidationPolicy#getCertificateCountryConstraint.
func (p *EtsiValidationPolicy) CertificateCountryConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Country)
	}
	return nil
}

// CertificateLocalityConstraint returns the CertificateLocality constraint,
// if present. Ports EtsiValidationPolicy#getCertificateLocalityConstraint.
func (p *EtsiValidationPolicy) CertificateLocalityConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Locality)
	}
	return nil
}

// CertificateStateConstraint returns the CertificateState constraint, if
// present. Ports EtsiValidationPolicy#getCertificateStateConstraint.
func (p *EtsiValidationPolicy) CertificateStateConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.State)
	}
	return nil
}

// CertificateOrganizationIdentifierConstraint returns the
// CertificateOrganizationIdentifier constraint, if present. Ports
// EtsiValidationPolicy#getCertificateOrganizationIdentifierConstraint.
func (p *EtsiValidationPolicy) CertificateOrganizationIdentifierConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.OrganizationIdentifier)
	}
	return nil
}

// CertificateOrganizationNameConstraint returns the
// CertificateOrganizationName constraint, if present. Ports
// EtsiValidationPolicy#getCertificateOrganizationNameConstraint.
func (p *EtsiValidationPolicy) CertificateOrganizationNameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.OrganizationName)
	}
	return nil
}

// CertificateOrganizationUnitConstraint returns the
// CertificateOrganizationUnit constraint, if present. Ports
// EtsiValidationPolicy#getCertificateOrganizationUnitConstraint.
func (p *EtsiValidationPolicy) CertificateOrganizationUnitConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.OrganizationUnit)
	}
	return nil
}

// CertificateIssuerNameConstraint returns certificate IssuerName
// constraint. Ports EtsiValidationPolicy#getCertificateIssuerNameConstraint.
func (p *EtsiValidationPolicy) CertificateIssuerNameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.IssuerName)
	}
	return nil
}

// CertificateSurnameConstraint returns the CertificateSurname constraint,
// if present. Ports EtsiValidationPolicy#getCertificateSurnameConstraint.
func (p *EtsiValidationPolicy) CertificateSurnameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Surname)
	}
	return nil
}

// CertificateGivenNameConstraint returns the CertificateGivenName
// constraint, if present. Ports EtsiValidationPolicy#getCertificateGivenNameConstraint.
func (p *EtsiValidationPolicy) CertificateGivenNameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.GivenName)
	}
	return nil
}

// CertificateCommonNameConstraint returns the CertificateCommonName
// constraint, if present. Ports EtsiValidationPolicy#getCertificateCommonNameConstraint.
func (p *EtsiValidationPolicy) CertificateCommonNameConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.CommonName)
	}
	return nil
}

// CertificatePseudonymConstraint returns the CertificatePseudonym
// constraint, if present. Ports EtsiValidationPolicy#getCertificatePseudonymConstraint.
func (p *EtsiValidationPolicy) CertificatePseudonymConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Pseudonym)
	}
	return nil
}

// CertificatePseudoUsageConstraint returns the CertificatePseudoUsage
// constraint, if present. Ports EtsiValidationPolicy#getCertificatePseudoUsageConstraint.
func (p *EtsiValidationPolicy) CertificatePseudoUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.UsePseudonym)
	}
	return nil
}

// CertificateTitleConstraint returns the CertificateTitle constraint, if
// present. Ports EtsiValidationPolicy#getCertificateTitleConstraint.
func (p *EtsiValidationPolicy) CertificateTitleConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Title)
	}
	return nil
}

// CertificateEmailConstraint returns the CertificateEmail constraint, if
// present. Ports EtsiValidationPolicy#getCertificateEmailConstraint.
func (p *EtsiValidationPolicy) CertificateEmailConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.MultiValuesRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toMultiValuesRule(cc.Email)
	}
	return nil
}

// CertificateSerialNumberConstraint returns the CertificateSerialNumber
// constraint, if present. Ports EtsiValidationPolicy#getCertificateSerialNumberConstraint.
func (p *EtsiValidationPolicy) CertificateSerialNumberConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.SerialNumberPresent)
	}
	return nil
}

// CertificateAuthorityInfoAccessPresentConstraint returns the
// CertificateAuthorityInfoAccessPresent constraint, if present. Ports
// EtsiValidationPolicy#getCertificateAuthorityInfoAccessPresentConstraint.
func (p *EtsiValidationPolicy) CertificateAuthorityInfoAccessPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.AuthorityInfoAccessPresent)
	}
	return nil
}

// RevocationDataSkipConstraint returns the RevocationDataSkip constraint,
// if present. Ports EtsiValidationPolicy#getRevocationDataSkipConstraint.
func (p *EtsiValidationPolicy) RevocationDataSkipConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.CertificateApplicabilityRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil && cc.RevocationDataSkip != nil {
		return toCertificateApplicabilityRule(cc.RevocationDataSkip)
	}
	return nil
}

// CertificateRevocationInfoAccessPresentConstraint returns the
// CertificateRevocationInfoAccessPresent constraint, if present. Ports
// EtsiValidationPolicy#getCertificateRevocationInfoAccessPresentConstraint.
func (p *EtsiValidationPolicy) CertificateRevocationInfoAccessPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) modelpolicy.LevelRule {
	if cc := p.certificateConstraints(context, subContext); cc != nil {
		return toLevelRule(cc.RevocationInfoAccessPresent)
	}
	return nil
}

// AcceptedContainerTypesConstraint returns the AcceptedContainerTypes
// constraint, if present. Ports EtsiValidationPolicy#getAcceptedContainerTypesConstraint.
func (p *EtsiValidationPolicy) AcceptedContainerTypesConstraint() modelpolicy.MultiValuesRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toMultiValuesRule(cc.AcceptableContainerTypes)
	}
	return nil
}

// ZipCommentPresentConstraint returns the ZipCommentPresent constraint, if
// present. Ports EtsiValidationPolicy#getZipCommentPresentConstraint.
func (p *EtsiValidationPolicy) ZipCommentPresentConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.ZipCommentPresent)
	}
	return nil
}

// AcceptedZipCommentsConstraint returns the AcceptedZipComments constraint,
// if present. Ports EtsiValidationPolicy#getAcceptedZipCommentsConstraint.
func (p *EtsiValidationPolicy) AcceptedZipCommentsConstraint() modelpolicy.MultiValuesRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toMultiValuesRule(cc.AcceptableZipComment)
	}
	return nil
}

// MimeTypeFilePresentConstraint returns the MimeTypeFilePresent
// constraint, if present. Ports EtsiValidationPolicy#getMimeTypeFilePresentConstraint.
func (p *EtsiValidationPolicy) MimeTypeFilePresentConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.MimeTypeFilePresent)
	}
	return nil
}

// AcceptedMimeTypeContentsConstraint returns the AcceptedMimeTypeContents
// constraint, if present. Ports EtsiValidationPolicy#getAcceptedMimeTypeContentsConstraint.
func (p *EtsiValidationPolicy) AcceptedMimeTypeContentsConstraint() modelpolicy.MultiValuesRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toMultiValuesRule(cc.AcceptableMimeTypeFileContent)
	}
	return nil
}

// ManifestFilePresentConstraint returns the ManifestFilePresent
// constraint, if present. Ports EtsiValidationPolicy#getManifestFilePresentConstraint.
func (p *EtsiValidationPolicy) ManifestFilePresentConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.ManifestFilePresent)
	}
	return nil
}

// SignedFilesPresentConstraint returns the SignedFilesPresent constraint,
// if present. Ports EtsiValidationPolicy#getSignedFilesPresentConstraint.
func (p *EtsiValidationPolicy) SignedFilesPresentConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.SignedFilesPresent)
	}
	return nil
}

// FilenameAdherenceConstraint returns the FilenameAdherence constraint, if
// present. Ports EtsiValidationPolicy#getFilenameAdherenceConstraint.
func (p *EtsiValidationPolicy) FilenameAdherenceConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.FilenameAdherence)
	}
	return nil
}

// AllFilesSignedConstraint returns the AllFilesSigned constraint, if
// present. Ports EtsiValidationPolicy#getAllFilesSignedConstraint.
func (p *EtsiValidationPolicy) AllFilesSignedConstraint() modelpolicy.LevelRule {
	if cc := p.ContainerConstraints(); cc != nil {
		return toLevelRule(cc.AllFilesSigned)
	}
	return nil
}

// FullScopeConstraint returns the FullScope constraint, if present. Ports
// EtsiValidationPolicy#getFullScopeConstraint.
func (p *EtsiValidationPolicy) FullScopeConstraint() modelpolicy.LevelRule {
	if sc := p.SignatureConstraints(); sc != nil {
		return toLevelRule(sc.FullScope)
	}
	return nil
}

// AcceptablePDFAProfilesConstraint returns the AcceptablePDFAProfiles
// constraint, if present. Ports EtsiValidationPolicy#getAcceptablePDFAProfilesConstraint.
func (p *EtsiValidationPolicy) AcceptablePDFAProfilesConstraint() modelpolicy.MultiValuesRule {
	if pc := p.PDFAConstraints(); pc != nil {
		return toMultiValuesRule(pc.AcceptablePDFAProfiles)
	}
	return nil
}

// PDFACompliantConstraint returns the PDFACompliant constraint, if
// present. Ports EtsiValidationPolicy#getPDFACompliantConstraint.
func (p *EtsiValidationPolicy) PDFACompliantConstraint() modelpolicy.LevelRule {
	if pc := p.PDFAConstraints(); pc != nil {
		return toLevelRule(pc.PDFACompliant)
	}
	return nil
}

// EAASignatureUnicityConstraint returns the EAASignatureUnicity
// constraint, if present. Ports EtsiValidationPolicy#getEAASignatureUnicityConstraint.
func (p *EtsiValidationPolicy) EAASignatureUnicityConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAASignatureUnicity)
	}
	return nil
}

// EAASignatureValidConstraint returns the EAASignatureValid constraint, if
// present. Ports EtsiValidationPolicy#getEAASignatureValidConstraint.
func (p *EtsiValidationPolicy) EAASignatureValidConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAASignatureValid)
	}
	return nil
}

// EAADisclosurePresentConstraint returns the DisclosurePresent constraint,
// if present. Ports EtsiValidationPolicy#getEAADisclosurePresentConstraint.
func (p *EtsiValidationPolicy) EAADisclosurePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.DisclosurePresent)
	}
	return nil
}

// EAADisclosureFoundConstraint returns the DisclosureFound constraint, if
// present. Ports EtsiValidationPolicy#getEAADisclosureFoundConstraint.
func (p *EtsiValidationPolicy) EAADisclosureFoundConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.DisclosureFound)
	}
	return nil
}

// EAADisclosureIntactConstraint returns the DisclosureIntact constraint,
// if present. Ports EtsiValidationPolicy#getEAADisclosureIntactConstraint.
func (p *EtsiValidationPolicy) EAADisclosureIntactConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.DisclosureIntact)
	}
	return nil
}

// EAADisclosureListExhaustiveConstraint returns the
// DisclosureListExhaustive constraint, if present. Ports
// EtsiValidationPolicy#getEAADisclosureListExhaustiveConstraint.
func (p *EtsiValidationPolicy) EAADisclosureListExhaustiveConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.DisclosureListExhaustive)
	}
	return nil
}

// EAAKeyBindingSignaturePresentConstraint returns the
// KeyBindingSignaturePresent constraint, if present. Ports
// EtsiValidationPolicy#getEAAKeyBindingSignaturePresentConstraint.
func (p *EtsiValidationPolicy) EAAKeyBindingSignaturePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.KeyBindingSignaturePresent)
	}
	return nil
}

// EAAKeyBindingSignatureValidConstraint returns the KeyBindingSignatureValid
// constraint, if present. Ports EtsiValidationPolicy#getEAAKeyBindingSignatureValidConstraint.
func (p *EtsiValidationPolicy) EAAKeyBindingSignatureValidConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.KeyBindingSignatureValid)
	}
	return nil
}

// EAATypeIntegrityPresentConstraint returns the EAATypeIntegrityPresent
// constraint, if present. Ports EtsiValidationPolicy#getEAATypeIntegrityPresentConstraint.
func (p *EtsiValidationPolicy) EAATypeIntegrityPresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAATypeIntegrityPresent)
	}
	return nil
}

// EAAIdentifierPresentConstraint returns the EAAIdentifierPresent
// constraint, if present. Ports EtsiValidationPolicy#getEAAIdentifierPresentConstraint.
func (p *EtsiValidationPolicy) EAAIdentifierPresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAIdentifierPresent)
	}
	return nil
}

// EAAIssuanceDatePresentConstraint returns the EAAIssuanceDatePresent
// constraint, if present. Ports EtsiValidationPolicy#getEAAIssuanceDatePresentConstraint.
func (p *EtsiValidationPolicy) EAAIssuanceDatePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAIssuanceDatePresent)
	}
	return nil
}

// EAACategoryConstraint returns the EAACategory constraint, if present.
// Ports EtsiValidationPolicy#getEAACategoryConstraint.
func (p *EtsiValidationPolicy) EAACategoryConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAACategory)
	}
	return nil
}

// EAASubjectConstraint returns the EAASubject constraint, if present.
// Ports EtsiValidationPolicy#getEAASubjectConstraint.
func (p *EtsiValidationPolicy) EAASubjectConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAASubject)
	}
	return nil
}

// EAASubjectPseudonymConstraint returns the EAASubjectPseudonym
// constraint, if present. Ports EtsiValidationPolicy#getEAASubjectPseudonymConstraint.
func (p *EtsiValidationPolicy) EAASubjectPseudonymConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAASubjectPseudonym)
	}
	return nil
}

// EAAIssuingCountryConstraint returns the EAAIssuingCountry constraint, if
// present. Ports EtsiValidationPolicy#getEAAIssuingCountryConstraint.
func (p *EtsiValidationPolicy) EAAIssuingCountryConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAAIssuingCountry)
	}
	return nil
}

// EAAIssuingAuthorityConstraint returns the EAAIssuingAuthority
// constraint, if present. Ports EtsiValidationPolicy#getEAAIssuingAuthorityConstraint.
func (p *EtsiValidationPolicy) EAAIssuingAuthorityConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAAIssuingAuthority)
	}
	return nil
}

// EAAIssuingAuthorityRegistrationIdentifierConstraint returns the
// EAAIssuingAuthorityRegistrationIdentifier constraint, if present. Ports
// EtsiValidationPolicy#getEAAIssuingAuthorityRegistrationIdentifierConstraint.
func (p *EtsiValidationPolicy) EAAIssuingAuthorityRegistrationIdentifierConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAAIssuingAuthorityRegistrationIdentifier)
	}
	return nil
}

// EAARevocationPresentConstraint returns the EAARevocationPresent
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationPresentConstraint.
func (p *EtsiValidationPolicy) EAARevocationPresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAARevocationPresent)
	}
	return nil
}

// EAAShortLivedConstraint returns the EAAShortLived constraint, if
// present. Ports EtsiValidationPolicy#getEAAShortLivedConstraint.
func (p *EtsiValidationPolicy) EAAShortLivedConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAShortLived)
	}
	return nil
}

// EAAOneTimeUseConstraint returns the EAAOneTimeUse constraint, if
// present. Ports EtsiValidationPolicy#getEAAOneTimeUseConstraint.
func (p *EtsiValidationPolicy) EAAOneTimeUseConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAOneTimeUse)
	}
	return nil
}

// EAAUsePseudonymConstraint returns the EAAUsePseudonym constraint, if
// present. Ports EtsiValidationPolicy#getEAAUsePseudonymConstraint.
func (p *EtsiValidationPolicy) EAAUsePseudonymConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAUsePseudonym)
	}
	return nil
}

// EAAClaimsConstraint returns the EAAClaims constraint, if present. Ports
// EtsiValidationPolicy#getEAAClaimsConstraint.
func (p *EtsiValidationPolicy) EAAClaimsConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAAClaims)
	}
	return nil
}

// EAASupportedClaimsConstraint returns the EAASupportedClaims constraint,
// if present. Ports EtsiValidationPolicy#getEAASupportedClaimsConstraint.
func (p *EtsiValidationPolicy) EAASupportedClaimsConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAASupportedClaims)
	}
	return nil
}

// EAARevocationAvailableConstraint returns the EAARevocationAvailable
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationAvailableConstraint.
func (p *EtsiValidationPolicy) EAARevocationAvailableConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAARevocationAvailable)
	}
	return nil
}

// AcceptableEAARevocationFoundConstraint returns the
// AcceptableEAARevocationFound constraint, if present. Ports
// EtsiValidationPolicy#getAcceptableEAARevocationFoundConstraint.
func (p *EtsiValidationPolicy) AcceptableEAARevocationFoundConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.AcceptableEAARevocationFound)
	}
	return nil
}

// EAARevocationNotRevokedConstraint returns the NotRevoked constraint, if
// present. Ports EtsiValidationPolicy#getEAARevocationNotRevokedConstraint.
func (p *EtsiValidationPolicy) EAARevocationNotRevokedConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.NotRevoked)
	}
	return nil
}

// EAARevocationNotOnHoldConstraint returns the NotOnHold constraint, if
// present. Ports EtsiValidationPolicy#getEAARevocationNotOnHoldConstraint.
func (p *EtsiValidationPolicy) EAARevocationNotOnHoldConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.NotOnHold)
	}
	return nil
}

// EAATypeConstraint returns the EAAType constraint, if present. Ports
// EtsiValidationPolicy#getEAATypeConstraint.
func (p *EtsiValidationPolicy) EAATypeConstraint() modelpolicy.MultiValuesRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toMultiValuesRule(eaa.EAAType)
	}
	return nil
}

// EAANotBeforePresentConstraint returns the EAANotBeforePresent
// constraint, if present. Ports EtsiValidationPolicy#getEAANotBeforePresentConstraint.
func (p *EtsiValidationPolicy) EAANotBeforePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAANotBeforePresent)
	}
	return nil
}

// EAAExpirationPresentConstraint returns the EAAExpirationPresent
// constraint, if present. Ports EtsiValidationPolicy#getEAAExpirationPresentConstraint.
func (p *EtsiValidationPolicy) EAAExpirationPresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAExpirationPresent)
	}
	return nil
}

// EAANotExpiredConstraint returns the EAANotExpired constraint, if
// present. Ports EtsiValidationPolicy#getEAANotExpiredConstraint.
func (p *EtsiValidationPolicy) EAANotExpiredConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAANotExpired)
	}
	return nil
}

// EAAAdministrativeIssuanceDatePresentConstraint returns the
// EAAAdministrativeIssuanceDatePresent constraint, if present. Ports
// EtsiValidationPolicy#getEAAAdministrativeIssuanceDatePresentConstraint.
func (p *EtsiValidationPolicy) EAAAdministrativeIssuanceDatePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAAdministrativeIssuanceDatePresent)
	}
	return nil
}

// EAAAdministrativeExpirationDatePresentConstraint returns the
// EAAAdministrativeExpirationDatePresent constraint, if present. Ports
// EtsiValidationPolicy#getEAAAdministrativeExpirationDatePresentConstraint.
func (p *EtsiValidationPolicy) EAAAdministrativeExpirationDatePresentConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAAdministrativeExpirationDatePresent)
	}
	return nil
}

// EAAAdministrativePeriodNotExpiredConstraint returns the
// EAAAdministrativePeriodNotExpired constraint, if present. Ports
// EtsiValidationPolicy#getEAAAdministrativePeriodNotExpiredConstraint.
func (p *EtsiValidationPolicy) EAAAdministrativePeriodNotExpiredConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.EAAAdministrativePeriodNotExpired)
	}
	return nil
}

// EAAETSI194721ConformanceConstraint returns the ETSI194721Conformance
// constraint, if present. Ports EtsiValidationPolicy#getEAAETSI194721ConformanceConstraint.
func (p *EtsiValidationPolicy) EAAETSI194721ConformanceConstraint() modelpolicy.LevelRule {
	if eaa := p.EAAConstraints(); eaa != nil {
		return toLevelRule(eaa.ETSI194721Conformance)
	}
	return nil
}

// EAARevocationTokenTypeConstraint returns the StatusTokenType constraint,
// if present. Ports EtsiValidationPolicy#getEAARevocationTokenTypeConstraint.
func (p *EtsiValidationPolicy) EAARevocationTokenTypeConstraint() modelpolicy.MultiValuesRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toMultiValuesRule(eaaRev.Type)
	}
	return nil
}

// EAARevocationUnknownStatusConstraint returns the UnknownStatus
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationUnknownStatusConstraint.
func (p *EtsiValidationPolicy) EAARevocationUnknownStatusConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.UnknownStatus)
	}
	return nil
}

// EAARevocationIssuanceTimeConstraint returns the IssuanceTime constraint,
// if present. Ports EtsiValidationPolicy#getEAARevocationIssuanceTimeConstraint.
func (p *EtsiValidationPolicy) EAARevocationIssuanceTimeConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.IssuanceTime)
	}
	return nil
}

// EAARevocationExpirationTimeConstraint returns the ExpirationTime
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationExpirationTimeConstraint.
func (p *EtsiValidationPolicy) EAARevocationExpirationTimeConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.ExpirationTime)
	}
	return nil
}

// EAARevocationNotExpiredConstraint returns the NotExpired constraint, if
// present. Ports EtsiValidationPolicy#getEAARevocationNotExpiredConstraint.
func (p *EtsiValidationPolicy) EAARevocationNotExpiredConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.NotExpired)
	}
	return nil
}

// EAARevocationSubjectConstraint returns the EAARevocationSubject
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationSubjectConstraint.
func (p *EtsiValidationPolicy) EAARevocationSubjectConstraint() modelpolicy.MultiValuesRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toMultiValuesRule(eaaRev.Subject)
	}
	return nil
}

// EAARevocationSubjectMatchConstraint returns the EAARevocationSubjectMatch
// constraint, if present. Ports EtsiValidationPolicy#getEAARevocationSubjectMatchConstraint.
func (p *EtsiValidationPolicy) EAARevocationSubjectMatchConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.SubjectMatch)
	}
	return nil
}

// EAARevocationIssuerValidAtIssuanceTimeConstraint returns the
// EAARevocationIssuerValidAtIssuanceTime constraint, if present. Ports
// EtsiValidationPolicy#getEAARevocationIssuerValidAtIssuanceTimeConstraint.
func (p *EtsiValidationPolicy) EAARevocationIssuerValidAtIssuanceTimeConstraint() modelpolicy.LevelRule {
	if eaaRev := p.EAARevocationConstraints(); eaaRev != nil {
		return toLevelRule(eaaRev.IssuerValidAtIssuanceTime)
	}
	return nil
}

// EIDASConstraintPresent returns true if EIDAS constraints are present
// (qualification check shall be performed). Ports
// EtsiValidationPolicy#isEIDASConstraintPresent.
func (p *EtsiValidationPolicy) EIDASConstraintPresent() bool {
	return p.EIDASConstraints() != nil
}

// TLFreshnessConstraint returns the TLFreshness constraint, if present.
// Ports EtsiValidationPolicy#getTLFreshnessConstraint.
func (p *EtsiValidationPolicy) TLFreshnessConstraint() modelpolicy.DurationRule {
	if e := p.EIDASConstraints(); e != nil {
		return toDurationRule(e.TLFreshness)
	}
	return nil
}

// TLWellSignedConstraint returns the TLWellSigned constraint, if present.
// Ports EtsiValidationPolicy#getTLWellSignedConstraint.
func (p *EtsiValidationPolicy) TLWellSignedConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.TLWellSigned)
	}
	return nil
}

// TLNotExpiredConstraint returns the TLNotExpired constraint, if present.
// Ports EtsiValidationPolicy#getTLNotExpiredConstraint.
func (p *EtsiValidationPolicy) TLNotExpiredConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.TLNotExpired)
	}
	return nil
}

// TLVersionConstraint returns the TLVersion constraint, if present. Ports
// EtsiValidationPolicy#getTLVersionConstraint.
func (p *EtsiValidationPolicy) TLVersionConstraint() modelpolicy.MultiValuesRule {
	if e := p.EIDASConstraints(); e != nil {
		return toMultiValuesRule(e.TLVersion)
	}
	return nil
}

// TLStructureConstraint returns the TLStructure constraint, if present.
// Ports EtsiValidationPolicy#getTLStructureConstraint.
func (p *EtsiValidationPolicy) TLStructureConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.TLStructure)
	}
	return nil
}

// LoTEFreshnessConstraint returns the LoTEFreshness constraint, if
// present. Ports EtsiValidationPolicy#getLoTEFreshnessConstraint.
func (p *EtsiValidationPolicy) LoTEFreshnessConstraint() modelpolicy.DurationRule {
	if e := p.EIDASConstraints(); e != nil {
		return toDurationRule(e.LoTEFreshness)
	}
	return nil
}

// LoTEWellSignedConstraint returns the LoTEWellSigned constraint, if
// present. Ports EtsiValidationPolicy#getLoTEWellSignedConstraint.
func (p *EtsiValidationPolicy) LoTEWellSignedConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.LoTEWellSigned)
	}
	return nil
}

// LoTENotExpiredConstraint returns the LoTENotExpired constraint, if
// present. Ports EtsiValidationPolicy#getLoTENotExpiredConstraint.
func (p *EtsiValidationPolicy) LoTENotExpiredConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.LoTENotExpired)
	}
	return nil
}

// LoTEVersionConstraint returns the LoTEVersion constraint, if present.
// Ports EtsiValidationPolicy#getLoTEVersionConstraint.
func (p *EtsiValidationPolicy) LoTEVersionConstraint() modelpolicy.MultiValuesRule {
	if e := p.EIDASConstraints(); e != nil {
		return toMultiValuesRule(e.LoTEVersion)
	}
	return nil
}

// LoTEStructureConstraint returns the LoTEStructure constraint, if
// present. Ports EtsiValidationPolicy#getLoTEStructureConstraint.
func (p *EtsiValidationPolicy) LoTEStructureConstraint() modelpolicy.LevelRule {
	if e := p.EIDASConstraints(); e != nil {
		return toLevelRule(e.LoTEStructure)
	}
	return nil
}

// ValidationModel returns the used validation model (default SHELL;
// alternatives CHAIN and HYBRID). Ports EtsiValidationPolicy#getValidationModel.
func (p *EtsiValidationPolicy) ValidationModel() enumerations.ValidationModel {
	currentModel := etsiValidationPolicyDefaultValidationModel
	if mc := p.policy.Model; mc != nil && mc.Value.ValidationModel() != "" {
		currentModel = mc.Value.ValidationModel()
	}
	return currentModel
}

// ---------------------------------------------------------------- accessors
//
// These are not part of the model/policy.ValidationPolicy interface, but
// upstream exposes them as public methods on EtsiValidationPolicy; ported
// faithfully (see PORTING.md: "keep ALL methods").

// SignatureConstraints returns the constraint used for Signature
// validation. Ports EtsiValidationPolicy#getSignatureConstraints.
func (p *EtsiValidationPolicy) SignatureConstraints() *jaxb.SignatureConstraints {
	return p.policy.SignatureConstraints
}

// CounterSignatureConstraints returns the constraint used for Counter
// Signature validation. Ports EtsiValidationPolicy#getCounterSignatureConstraints.
func (p *EtsiValidationPolicy) CounterSignatureConstraints() *jaxb.SignatureConstraints {
	return p.policy.CounterSignatureConstraints
}

// KeyBindingSignatureConstraints returns the constraint used for Key
// Binding Signature validation. Ports EtsiValidationPolicy#getKeyBindingSignatureConstraints.
func (p *EtsiValidationPolicy) KeyBindingSignatureConstraints() *jaxb.SignatureConstraints {
	return p.policy.KeyBindingSignatureConstraints
}

// TimestampConstraints returns the constraint used for Timestamp
// validation. Ports EtsiValidationPolicy#getTimestampConstraints.
func (p *EtsiValidationPolicy) TimestampConstraints() *jaxb.TimestampConstraints {
	return p.policy.Timestamp
}

// RevocationConstraints returns the constraint used for Revocation
// validation. Ports EtsiValidationPolicy#getRevocationConstraints.
func (p *EtsiValidationPolicy) RevocationConstraints() *jaxb.RevocationConstraints {
	return p.policy.Revocation
}

// EvidenceRecordConstraints returns the constraint used for Evidence
// Record validation. Ports EtsiValidationPolicy#getEvidenceRecordConstraints.
func (p *EtsiValidationPolicy) EvidenceRecordConstraints() *jaxb.EvidenceRecordConstraints {
	return p.policy.EvidenceRecord
}

// EAAConstraints returns the constraint used for EAA validation. Ports
// EtsiValidationPolicy#getEAAConstraints.
func (p *EtsiValidationPolicy) EAAConstraints() *jaxb.EAAConstraints {
	return p.policy.EAA
}

// EAARevocationConstraints returns the constraint used for EAA revocation
// token validation. Ports EtsiValidationPolicy#getEAARevocationConstraints.
func (p *EtsiValidationPolicy) EAARevocationConstraints() *jaxb.EAARevocationConstraints {
	return p.policy.EAARevocation
}

// ContainerConstraints returns the constraint used for ASiC Container
// validation. Ports EtsiValidationPolicy#getContainerConstraints.
func (p *EtsiValidationPolicy) ContainerConstraints() *jaxb.ContainerConstraints {
	return p.policy.ContainerConstraints
}

// PDFAConstraints returns the constraint used for PDF/A validation. Ports
// EtsiValidationPolicy#getPDFAConstraints.
func (p *EtsiValidationPolicy) PDFAConstraints() *jaxb.PDFAConstraints {
	return p.policy.PDFAConstraints
}

// EIDASConstraints returns the constraint used for qualification
// validation. Ports EtsiValidationPolicy#getEIDASConstraints.
func (p *EtsiValidationPolicy) EIDASConstraints() *jaxb.EIDAS {
	return p.policy.EIDAS
}

// Cryptographic returns the common constraint used for cryptographic
// validation. Ports EtsiValidationPolicy#getCryptographic.
func (p *EtsiValidationPolicy) Cryptographic() *jaxb.CryptographicConstraint {
	return p.policy.Cryptographic
}

// ------------------------------------------------------------- navigation

// signingCertificateByContext ports the private getSigningCertificateByContext(Context) helper.
func (p *EtsiValidationPolicy) signingCertificateByContext(context enumerations.Context) *jaxb.CertificateConstraints {
	return p.certificateConstraints(context, enumerations.SubContextSigningCert)
}

// certificateConstraints ports the private getCertificateConstraints(Context, SubContext) helper.
func (p *EtsiValidationPolicy) certificateConstraints(context enumerations.Context, subContext enumerations.SubContext) *jaxb.CertificateConstraints {
	bsc := p.basicSignatureConstraintsByContext(context)
	if bsc != nil {
		switch subContext {
		case enumerations.SubContextSigningCert:
			return bsc.SigningCertificate
		case enumerations.SubContextCACertificate:
			return bsc.CACertificate
		}
	}
	return nil
}

// basicSignatureConstraintsByContext ports the private
// getBasicSignatureConstraintsByContext(Context) helper.
func (p *EtsiValidationPolicy) basicSignatureConstraintsByContext(context enumerations.Context) *jaxb.BasicSignatureConstraints {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCertificate: // TODO improve
		if mainSignature := p.SignatureConstraints(); mainSignature != nil {
			return mainSignature.BasicSignatureConstraints
		}
	case enumerations.ContextCounterSignature:
		if counterSignature := p.CounterSignatureConstraints(); counterSignature != nil {
			return counterSignature.BasicSignatureConstraints
		}
	case enumerations.ContextKeyBindingSignature:
		if keyBindingSignature := p.KeyBindingSignatureConstraints(); keyBindingSignature != nil {
			return keyBindingSignature.BasicSignatureConstraints
		}
	case enumerations.ContextTimestamp:
		if timestampConstraints := p.TimestampConstraints(); timestampConstraints != nil {
			return timestampConstraints.BasicSignatureConstraints
		}
	case enumerations.ContextRevocation:
		if revocationConstraints := p.RevocationConstraints(); revocationConstraints != nil {
			return revocationConstraints.BasicSignatureConstraints
		}
	case enumerations.ContextEAA:
		return nil
	case enumerations.ContextEAARevocation:
		if eaaRevocationConstraints := p.EAARevocationConstraints(); eaaRevocationConstraints != nil {
			return eaaRevocationConstraints.BasicSignatureConstraints
		}
	default:
		panic(fmt.Sprintf("Unsupported context '%s'", context))
	}
	return nil
}

// signedAttributeConstraints ports the private getSignedAttributeConstraints(Context)
// helper. Java logs a warning (slf4j dropped, per PORTING.md) and returns
// null for an unsupported context rather than throwing.
func (p *EtsiValidationPolicy) signedAttributeConstraints(context enumerations.Context) *jaxb.SignedAttributesConstraints {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCertificate: // TODO improve
		if mainSignature := p.SignatureConstraints(); mainSignature != nil {
			return mainSignature.SignedAttributes
		}
	case enumerations.ContextCounterSignature:
		if counterSignature := p.CounterSignatureConstraints(); counterSignature != nil {
			return counterSignature.SignedAttributes
		}
	case enumerations.ContextKeyBindingSignature:
		if keyBindingSignature := p.KeyBindingSignatureConstraints(); keyBindingSignature != nil {
			return keyBindingSignature.SignedAttributes
		}
	case enumerations.ContextTimestamp:
		if timestampConstraints := p.TimestampConstraints(); timestampConstraints != nil {
			return timestampConstraints.SignedAttributes
		}
	}
	return nil
}

// unsignedAttributeConstraints ports the private getUnsignedAttributeConstraints(Context)
// helper. Java logs a warning (slf4j dropped, per PORTING.md) and returns
// null for an unsupported context rather than throwing.
func (p *EtsiValidationPolicy) unsignedAttributeConstraints(context enumerations.Context) *jaxb.UnsignedAttributesConstraints {
	switch context {
	case enumerations.ContextSignature:
		if mainSignature := p.SignatureConstraints(); mainSignature != nil {
			return mainSignature.UnsignedAttributes
		}
	case enumerations.ContextCounterSignature:
		if counterSignature := p.CounterSignatureConstraints(); counterSignature != nil {
			return counterSignature.UnsignedAttributes
		}
	case enumerations.ContextKeyBindingSignature:
		if keyBindingSignature := p.KeyBindingSignatureConstraints(); keyBindingSignature != nil {
			return keyBindingSignature.UnsignedAttributes
		}
	}
	return nil
}

// signatureConstraintsByContext ports the private getSignatureConstraintsByContext(Context)
// helper. Java logs a warning (slf4j dropped, per PORTING.md) and returns
// null for an unsupported context rather than throwing.
func (p *EtsiValidationPolicy) signatureConstraintsByContext(context enumerations.Context) *jaxb.SignatureConstraints {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCertificate: // TODO improve
		return p.SignatureConstraints()
	case enumerations.ContextCounterSignature:
		return p.CounterSignatureConstraints()
	case enumerations.ContextKeyBindingSignature:
		return p.KeyBindingSignatureConstraints()
	}
	return nil
}

// ------------------------------------------------------------------- toXXX
//
// Ports the private toLevelRule/toRule/toCryptographicSuite overloads. Each
// returns the exported model/policy interface type, not the concrete
// wrapper, and - unlike the wrapper constructors themselves, which accept a
// nil constraint - returns a true nil interface value for a nil constraint:
// a nil *XxxConstraintWrapper wrapped in a non-nil interface would panic the
// moment a caller invoked one of its methods (e.g. LevelConstraintWrapper.Level
// dereferences its receiver's field), the classic Go "typed nil in an
// interface" trap. Returning nil directly, mirroring Java's null, avoids it.

func toLevelRule(constraint *jaxb.LevelConstraint) modelpolicy.LevelRule {
	if constraint == nil {
		return nil
	}
	return NewLevelConstraintWrapper(constraint)
}

func toMultiValuesRule(constraint *jaxb.MultiValuesConstraint) modelpolicy.MultiValuesRule {
	if constraint == nil {
		return nil
	}
	return NewMultiValuesConstraintWrapper(constraint)
}

func toValueRule(constraint *jaxb.ValueConstraint) modelpolicy.ValueRule {
	if constraint == nil {
		return nil
	}
	return NewValueConstraintWrapper(constraint)
}

func toNumericValueRule(constraint *jaxb.IntValueConstraint) modelpolicy.NumericValueRule {
	if constraint == nil {
		return nil
	}
	return NewIntValueConstraintWrapper(constraint)
}

func toDurationRule(constraint *jaxb.TimeConstraint) modelpolicy.DurationRule {
	if constraint == nil {
		return nil
	}
	return NewTimeConstraintWrapper(constraint)
}

func toCertificateApplicabilityRule(constraint *jaxb.CertificateValuesConstraint) modelpolicy.CertificateApplicabilityRule {
	if constraint == nil {
		return nil
	}
	return NewCertificateValuesConstraintWrapper(constraint)
}

func toCryptographicSuite(constraint *jaxb.CryptographicConstraint) modelpolicy.CryptographicSuite {
	return NewCryptographicConstraintWrapper(constraint)
}

// String ports EtsiValidationPolicy#toString.
func (p *EtsiValidationPolicy) String() string {
	return fmt.Sprintf("EtsiValidationPolicy [policyName=%s]", p.PolicyName())
}
