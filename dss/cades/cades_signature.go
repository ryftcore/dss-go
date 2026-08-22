// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignature.java
// (DSS 6.5.RC1).
//
// # BouncyCastle replacements (see cades_utils.go's header)
//
//   - org.bouncycastle.cms.SignerInformation      -> *cmscore.SignerInfo
//   - org.bouncycastle.cms.SignerInformationStore -> []*cmscore.SignerInfo
//   - org.bouncycastle.cms.SignerId               -> *cmscore.SignerIdentifier (the SID field)
//   - org.bouncycastle.asn1.cms.AttributeTable    -> cmscore.Attributes
//   - org.bouncycastle.asn1.cms.Attribute         -> *cmscore.Attribute
//   - eu.europa.esig.dss.cms.CMS                  -> *cms.CMS
//   - every ASN.1 object BouncyCastle supplies with no DSS class of its own (SignerLocation,
//     CommitmentTypeIndication, SignaturePolicyIdentifier/Id, OtherHashAlgAndValue,
//     SigPolicyQualifiers/Info, SPUserNotice, NoticeReference, ContentHints, ContentIdentifier,
//     RSASSAPSSparams, the RFC 5126 v1 SignerAttribute, the RFC 5755 AttributeCertificate and
//     RoleSyntax) -> parsed here directly with internal/asn1ber, next to their single reader,
//     the same choice cades_level_baseline_b.go made for the write side of several of the same
//     structures (SignerLocation, CommitmentTypeIndication, the v1 SignerAttribute ClaimedAttributes
//     branch).
//
// # BC SignerInformation accessors with no cmscore.SignerInfo counterpart
//
// internal/cmscore is frozen (see PORTING.md); every BC SignerInformation getter this file
// needs but cmscore.SignerInfo does not expose is computed here from the fields cmscore.SignerInfo
// does carry (DigestAlgorithm, SignatureAlgorithm, Signature, SID): encryptionAlgOID/digestAlgOID
// read the two AlgorithmIdentifier.Algorithm OIDs directly, encryptionAlgParams reads
// SignatureAlgorithm.Parameters, and getSID() is the SID field itself. contentDigest (used only
// by the very rare non-BASELINE "message-digest absent" fallback) has no accessible equivalent
// without re-running BC's own digest-over-content computation; contentReferenceValidation below
// documents the narrowed behaviour this causes.
//
// cades_timestamp_source.go additionally documents its own assumptions about this file's public
// surface (SignerInformation, CMS, DetachedContents, CertificateSource, CRLSource, OCSPSource,
// CounterSignatures, ID) - all provided below, matching that file's header exactly.
//
// # Deviations
//
//   - isCounterSignature(): Java queries BC's SignerInformation#isCounterSignature(), a flag BC
//     sets only when it extracts a SignerInformation from another SignerInfo's countersignature
//     unsigned attribute (SignerInformationStore#getCounterSignatures) - not derivable from the
//     encoded SignerInfo itself, and cmscore.SignerInfo (frozen) carries no such flag. Every
//     code path that builds a counter-signature CAdESSignature (CounterSignatures below) calls
//     SetMasterSignature immediately, and nothing observes IsCounterSignature before that call
//     completes, so DefaultAdvancedSignature's own masterSignature-based IsCounterSignature() is
//     behaviourally equivalent here; IsCounterSignature is still redefined explicitly below
//     (rather than left to promotion) so the equivalence is visible at the call site rather than
//     silently relied upon.
//   - getContentReferenceValidation's contentDigest branch (RFC 5652 "message-digest absent",
//     reachable only for a non-ETSI CMS signature with no signed attributes at all) needs BC's
//     SignerInformation#getContentDigest(), which recomputes and caches the digest BC's verifier
//     already streamed the content through. This port has no such cached value and does not
//     redo the streaming digest computation here (it would duplicate CheckSignatureIntegrity's
//     own verification pass for a code path every BASELINE-B+ profile check already excludes via
//     the message-digest cardinality==1 requirement). contentReferenceValidation therefore always
//     reports "not found", matching the *outcome* Java reaches whenever getContentDigest() cannot
//     produce a value, without attempting the recomputation. Flagged for integrator follow-up.
//   - getCertifiedSignerRoles walks a hand-decoded RFC 5755 AttributeCertificate/RoleSyntax
//     (BC's org.bouncycastle.asn1.x509.{AttributeCertificate, AttributeCertificateInfo,
//     AttCertValidityPeriod, RoleSyntax}, none ported elsewhere). AttributeCertificateInfo's
//     holder/issuer/signature/serialNumber fields are opaque here and skipped positionally
//     (nothing this method reads needs their content); RoleSyntax.roleName is assumed to carry
//     the common case, a GeneralName uniformResourceIdentifier ([6] IMPLICIT IA5String) - the
//     only alternative RFC 3281's role-identifier convention actually uses in practice. A
//     roleName of a different GeneralName choice is skipped (best effort), matching how every
//     other malformed-input branch in this file degrades.
//   - getPSSHashAlgorithm's RSASSAPSSparams.hashAlgorithm ([0] EXPLICIT, DEFAULT sha1Identifier)
//     is likewise hand-decoded (org.bouncycastle.asn1.pkcs.RSASSAPSSparams has no DSS port).
//
// slf4j logging is dropped per PORTING.md; every LOG.warn/LOG.debug call site is called out in
// the surrounding comment instead. DSS's checked/unchecked exception split is preserved as
// documented at each method: a Java try/catch becomes a Go error check with the same fallback,
// and the absence of one becomes a Go panic carrying the Java message, reproducing the unchecked
// exception that would otherwise propagate out of the interface method uncaught (none of the
// AdvancedSignature accessors this file implements has an error return to use instead).
package cades

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// The PKCS#9 attribute types this file reads that no other landed file of this package has
// already declared (see cades_utils.go / cades_level_baseline_b.go for the rest), i.e. the
// org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers constants CAdESSignature static-imports.
// They keep their exact Java field name behind the "OID_" prefix, the convention every other
// file of this package already established.
var (
	// OIDPkcs9AtContentType is 1.2.840.113549.1.9.3.
	OIDPkcs9AtContentType = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}

	// OIDPkcs9AtMessageDigest is 1.2.840.113549.1.9.4.
	OIDPkcs9AtMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}

	// OIDIdAaContentReference is 1.2.840.113549.1.9.16.2.10.
	OIDIdAaContentReference = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 10}
)

// CAdESSignature is the CAdES Signature class helper. Port of the CAdESSignature class.
//
// serialVersionUID and java.io.Serializable are dropped (no Go counterpart).
type CAdESSignature struct {
	validation.DefaultAdvancedSignature

	// cmsDocument is the CMS of the signature. Port of the private final CMS cms field.
	cmsDocument *cms.CMS

	// signerInformation is the corresponding SignerInformation to the signature. Port of the
	// private final SignerInformation signerInformation field.
	signerInformation *cmscore.SignerInfo

	// counterSignaturesStore caches the counter-signature SignerInfos so that a unique
	// identifier for counter signatures can be properly computed.
	// Port of the private SignerInformationStore counterSignaturesStore field.
	counterSignaturesStore []*cmscore.SignerInfo

	// cachedCryptoVerification is this port's idempotency guard for CheckSignatureIntegrity, a
	// replacement for reading the (inaccessible, base-package-private)
	// signatureCryptographicVerification field directly the way Java's own method does for its
	// "already computed" early-return check. It is set to the same
	// *signature.SignatureCryptographicVerification value handed to
	// SetSignatureCryptographicVerification, so the two never disagree.
	cachedCryptoVerification *signature.SignatureCryptographicVerification
}

// NewCAdESSignature is the default constructor for CAdESSignature.
// Port of the public CAdESSignature(CMS, SignerInformation) constructor.
//
// Panics with the Java message when cmsDocument or signerInformation is missing
// (Objects.requireNonNull).
func NewCAdESSignature(cmsDocument *cms.CMS, signerInformation *cmscore.SignerInfo) *CAdESSignature {
	if cmsDocument == nil {
		panic("CMS cannot be null!")
	}
	if signerInformation == nil {
		panic("SignerInformation must be provided!")
	}
	s := &CAdESSignature{
		DefaultAdvancedSignature: validation.NewDefaultAdvancedSignatureBase(),
		cmsDocument:              cmsDocument,
		signerInformation:        signerInformation,
	}
	s.InitDefaultAdvancedSignature(s)
	return s
}

// SignatureForm specifies the format of the signature. Port of getSignatureForm().
func (s *CAdESSignature) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormCAdES
}

// CertificateSource gets a certificate source which contains ALL certificates embedded in the
// signature. Port of getCertificateSource().
//
// Panics with the underlying error's message when the source cannot be built: the Java
// constructor call is not guarded by a try/catch, so an unchecked exception here would
// propagate out of the Java getter uncaught.
func (s *CAdESSignature) CertificateSource() *spi.SignatureCertificateSource {
	if s.OfflineCertificateSource() == nil {
		cadesCertificateSource, err := NewCAdESCertificateSource(s.cmsDocument, s.signerInformation)
		if err != nil {
			panic(err)
		}
		s.SetOfflineCertificateSource(&cadesCertificateSource.SignatureCertificateSource)
	}
	return s.OfflineCertificateSource()
}

// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature.
// Port of getCRLSource().
func (s *CAdESSignature) CRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	if s.SignatureCRLSource() == nil {
		crlSource, err := NewCAdESCRLSource(s.cmsDocument, s.signerInformation.UnsignedAttributes)
		if err != nil {
			// Upstream logs "Error in computing or in format of the algorithm: just
			// continue..." and leaves signatureCRLSource null (will try to get online
			// information later).
			return nil
		}
		s.SetSignatureCRLSource(crlSource)
	}
	return s.SignatureCRLSource()
}

// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the signature.
// Port of getOCSPSource().
//
// Panics with the underlying error's message when the source cannot be built: the Java
// constructor call is not guarded by a try/catch, so an unchecked exception here would
// propagate out of the Java getter uncaught.
func (s *CAdESSignature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	if s.SignatureOCSPSource() == nil {
		ocspSource, err := NewCAdESOCSPSource(s.cmsDocument, s.signerInformation.UnsignedAttributes)
		if err != nil {
			panic(err)
		}
		s.SetSignatureOCSPSource(ocspSource)
	}
	return s.SignatureOCSPSource()
}

// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded in
// the signature. Port of getTimestampSource(), covariant in Java (returns CAdESTimestampSource);
// Go interface satisfaction needs the exact validation.TimestampSource return type, so CAdES
// callers needing the concrete type assert on the result (as counterSignatureTimestampSource
// already does in cades_timestamp_source.go via the plain AdvancedSignature.TimestampSource()).
func (s *CAdESSignature) TimestampSource() validation.TimestampSource {
	if s.SignatureTimestampSource() == nil {
		s.SetSignatureTimestampSource(NewCAdESTimestampSource(s))
	}
	return s.SignatureTimestampSource()
}

// SignerId returns the SignerIdentifier of the related signerInformation.
// Port of the public SignerId getSignerId().
func (s *CAdESSignature) SignerId() *cmscore.SignerIdentifier {
	return s.signerInformation.SID
}

// FindSignatureScopes finds signature scopes. Port of the protected findSignatureScopes().
func (s *CAdESSignature) FindSignatureScopes() []scope.SignatureScope {
	return NewCAdESSignatureScopeFinder().FindSignatureScope(s)
}

// BuildSignaturePolicy extracts a signature policy from a signature and builds the object.
// Port of the protected buildSignaturePolicy().
func (s *CAdESSignature) BuildSignaturePolicy() *signature.SignaturePolicy {
	attribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaEtsSigPolicyId)
	if attribute == nil {
		return nil
	}

	attrValue := spi.DSSASN1UtilsAsn1Encodable(attribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signature policy identifier attribute. Skip
		// processing.".
		return nil
	}

	if attrValue.IsUniversal(asn1ber.TagNull) {
		return signature.NewSignaturePolicy()
	}

	// SignaturePolicyIdentifier ::= CHOICE {
	//     signaturePolicyId      SignaturePolicyId,
	//     signaturePolicyImplied SignaturePolicyImplied } -- SignaturePolicyImplied ::= NULL
	//
	// A NULL alternative was already handled above; every other case is the SignaturePolicyId
	// SEQUENCE, i.e. org.bouncycastle.asn1.esf.SignaturePolicyId(ASN1ObjectIdentifier,
	// OtherHashAlgAndValue, SigPolicyQualifiers).
	if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() || len(attrValue.Children()) < 2 {
		return nil
	}
	children := attrValue.Children()

	policyIdOID, err := children[0].ObjectIdentifier()
	if err != nil {
		return nil
	}
	sigPolicy := signature.NewSignaturePolicyWithIdentifier(policyIdOID.String())

	// OtherHashAlgAndValue ::= SEQUENCE { hashAlgorithm AlgorithmIdentifier, hashValue OCTET STRING }
	sigPolicyHash := children[1]
	if !sigPolicyHash.IsUniversal(asn1ber.TagSequence) || !sigPolicyHash.IsConstructed() || len(sigPolicyHash.Children()) != 2 {
		return nil
	}
	digestValueBytes := sigPolicyHash.Children()[1].Octets()
	zeroHash := cadesIsZeroHash(digestValueBytes)
	sigPolicy.SetZeroHash(zeroHash)

	if !zeroHash {
		digestAlgorithmIdentifier, err := asn1ber.AlgorithmIdentifierFromElement(sigPolicyHash.Children()[0])
		if err == nil {
			digestAlgorithm := cadesDigestAlgorithmForOID(digestAlgorithmIdentifier.Algorithm.String())
			if digestAlgorithm != "" {
				sigPolicy.SetDigest(model.NewDigest(digestAlgorithm, digestValueBytes))
			} else {
				// Upstream logs "Signature policy identifier hash is not found or wrongly
				// encoded!".
			}
		}
	}

	if len(children) > 2 {
		s.buildSigPolicyQualifiers(sigPolicy, children[2])
	}

	return sigPolicy
}

// buildSigPolicyQualifiers ports the sigPolicyQualifiers loop of buildSignaturePolicy().
// Every entry is processed independently: a malformed SigPolicyQualifierInfo is skipped, as
// Java's per-entry try/catch(Exception) skips it (LOG.warn "Unable to read SigPolicyQualifierInfo
// {} : {}" dropped).
func (s *CAdESSignature) buildSigPolicyQualifiers(sigPolicy *signature.SignaturePolicy, sigPolicyQualifiers *asn1ber.Element) {
	if !sigPolicyQualifiers.IsUniversal(asn1ber.TagSequence) || !sigPolicyQualifiers.IsConstructed() {
		return
	}
	for _, qualifierInfo := range sigPolicyQualifiers.Children() {
		s.buildSigPolicyQualifierInfo(sigPolicy, qualifierInfo)
	}
}

// buildSigPolicyQualifierInfo ports one iteration's try body: SigPolicyQualifierInfo ::=
// SEQUENCE { sigPolicyQualifierId SigPolicyQualifierId, sigQualifier ANY DEFINED BY
// sigPolicyQualifierId }.
func (s *CAdESSignature) buildSigPolicyQualifierInfo(sigPolicy *signature.SignaturePolicy, qualifierInfo *asn1ber.Element) {
	if !qualifierInfo.IsUniversal(asn1ber.TagSequence) || !qualifierInfo.IsConstructed() || len(qualifierInfo.Children()) != 2 {
		return
	}
	policyQualifierInfoID, err := qualifierInfo.Children()[0].ObjectIdentifier()
	if err != nil {
		return
	}
	sigQualifier := qualifierInfo.Children()[1]

	switch {
	case OIDIdSpqEtsUri.Equal(policyQualifierInfoID):
		sigPolicy.SetURI(asn1ber.ASN1ToString(sigQualifier))

	case OIDIdSpqEtsUnotice.Equal(policyQualifierInfoID):
		userNotice := cadesBuildSPUserNotice(sigQualifier)
		if userNotice != nil {
			sigPolicy.SetUserNotice(userNotice)
		}

	case spi.OIDIdSpDocSpecification.Equal(policyQualifierInfoID):
		spDocSpecification := model.NewSpDocSpecification()
		spDocSpecification.SetId(asn1ber.ASN1ToString(sigQualifier))
		sigPolicy.SetDocSpecification(spDocSpecification)

	default:
		// Upstream logs "Unknown signature policy qualifier id: {} with value: {}".
	}
}

// cadesBuildSPUserNotice ports the private buildSPUserNoticeString(SPUserNotice).
//
//	SPUserNotice ::= SEQUENCE {
//	    noticeRef    NoticeReference OPTIONAL,
//	    explicitText DisplayText     OPTIONAL }
//	NoticeReference ::= SEQUENCE {
//	    organization  DisplayText,
//	    noticeNumbers SEQUENCE OF INTEGER }
//
// i.e. org.bouncycastle.asn1.esf.SPUserNotice/org.bouncycastle.asn1.x509.{NoticeReference,DisplayText}.
// A malformed value is skipped like every other qualifier (nil, dropping the "unable to build"
// path this file has no separate log for).
func cadesBuildSPUserNotice(sigQualifier *asn1ber.Element) *model.UserNotice {
	if !sigQualifier.IsUniversal(asn1ber.TagSequence) || !sigQualifier.IsConstructed() {
		return nil
	}
	userNotice := model.NewUserNotice()
	for _, child := range sigQualifier.Children() {
		switch {
		case child.IsUniversal(asn1ber.TagSequence) && child.IsConstructed():
			// noticeRef: NoticeReference ::= SEQUENCE { organization, noticeNumbers }
			refChildren := child.Children()
			if len(refChildren) != 2 {
				continue
			}
			userNotice.SetOrganization(asn1ber.ASN1ToString(refChildren[0]))
			if !refChildren[1].IsUniversal(asn1ber.TagSequence) || !refChildren[1].IsConstructed() {
				continue
			}
			numbers := refChildren[1].Children()
			noticeNumbers := make([]int, len(numbers))
			for index, number := range numbers {
				noticeNumbers[index] = int(number.Integer().Int64())
			}
			userNotice.SetNoticeNumbers(noticeNumbers...)
		default:
			// explicitText: DisplayText is a CHOICE of string types.
			userNotice.SetExplicitText(asn1ber.ASN1ToString(child))
		}
	}
	return userNotice
}

// cadesIsZeroHash ports isZeroHash/isZeroHashEmpty/doesZeroHashContainSigneZeroByte: the
// hashValue within sigPolicyHash may be set to zero or be empty to indicate that the policy
// hash value is not known.
func cadesIsZeroHash(hashValue []byte) bool {
	if len(hashValue) == 0 {
		return true
	}
	return len(hashValue) == 1 && (hashValue[0] == '0' || hashValue[0] == 0x00)
}

// SignaturePolicyStore returns the Signature Policy Store from the signature.
// Port of getSignaturePolicyStore().
func (s *CAdESSignature) SignaturePolicyStore() *model.SignaturePolicyStore {
	sigPolicyStore := CAdESUtilsUnsignedAttribute(s.signerInformation, spi.OIDIdAaEtsSigPolicyStore)
	if sigPolicyStore == nil || len(sigPolicyStore.Values) == 0 {
		return nil
	}

	attrValue := spi.DSSASN1UtilsAsn1Encodable(sigPolicyStore)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signature policy store attribute. Skip
		// processing.".
		return nil
	}
	if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
		return nil
	}

	children := attrValue.Children()
	if len(children) != 2 {
		// Upstream logs "Unable to extract a signature-policy-store. The element shall
		// contain two attributes.".
		return nil
	}

	signaturePolicyStore := model.NewSignaturePolicyStore()
	spDocSpecification := model.NewSpDocSpecification()
	spDocSpecification.SetId(asn1ber.ASN1ToString(children[0]))

	spDocument := children[1]
	switch {
	case spDocument.IsUniversal(asn1ber.TagOctetString):
		signaturePolicyStore.SetSignaturePolicyContent(model.NewInMemoryDocument(spDocument.Octets()))
	case spDocument.IsUniversal(asn1ber.TagIA5String):
		signaturePolicyStore.SetSigPolDocLocalURI(spDocument.AsString())
	default:
		// Upstream logs "Unable to extract a signature-policy-store spDocument. One of
		// 'sigPolicyEncoded' or 'sigPolicyLocalURI' is expected!".
	}
	signaturePolicyStore.SetSpDocSpecification(spDocSpecification)
	return signaturePolicyStore
}

// SigningTime returns the signing time included within the signature, or nil.
// Port of getSigningTime().
func (s *CAdESSignature) SigningTime() *time.Time {
	attr := CAdESUtilsSignedAttribute(s.signerInformation, OIDPkcs9AtSigningTime)
	if attr == nil {
		return nil
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(attr)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signing-time attribute. Skip processing.".
		return nil
	}
	signingDate := CAdESUtilsReadSigningDate(attrValue.Encoded())
	if signingDate.IsZero() {
		return nil
	}
	return &signingDate
}

// CMS gets the CMS. Port of the public CMS getCMS().
func (s *CAdESSignature) CMS() *cms.CMS {
	return s.cmsDocument
}

// SignatureProductionPlace returns information about the place where the signature was
// generated. Port of getSignatureProductionPlace().
func (s *CAdESSignature) SignatureProductionPlace() *signature.SignatureProductionPlace {
	attribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaEtsSignerLocation)
	if attribute == nil {
		return nil
	}

	attrValue := spi.DSSASN1UtilsAsn1Encodable(attribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signer-location attribute. Skip processing.".
		return nil
	}
	if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
		// Upstream logs "Unable to build a SignerLocation instance. Reason : {}".
		return nil
	}

	// SignerLocation ::= SEQUENCE {
	//     countryName   [0] EXPLICIT DirectoryString OPTIONAL,
	//     localityName  [1] EXPLICIT DirectoryString OPTIONAL,
	//     postalAddress [2] EXPLICIT PostalAddress OPTIONAL }
	// PostalAddress ::= SEQUENCE SIZE(1..6) OF DirectoryString
	var countryName, localityName string
	var postalAddressElement *asn1ber.Element
	for _, child := range attrValue.Children() {
		if !child.IsConstructed() || len(child.Children()) != 1 {
			continue
		}
		base := child.Children()[0]
		switch {
		case child.IsContextSpecific(0):
			countryName = spi.DSSASN1UtilsDirectoryStringValue(base.Encoded())
		case child.IsContextSpecific(1):
			localityName = spi.DSSASN1UtilsDirectoryStringValue(base.Encoded())
		case child.IsContextSpecific(2):
			postalAddressElement = base
		}
	}

	signatureProductionPlace := signature.NewSignatureProductionPlace()
	if countryName != "" {
		signatureProductionPlace.SetCountryName(countryName)
	}
	if localityName != "" {
		signatureProductionPlace.SetCity(localityName)
	}
	if postalAddressElement != nil {
		postalAddress := signatureProductionPlace.PostalAddress()
		for _, entry := range postalAddressElement.Children() {
			value := spi.DSSASN1UtilsDirectoryStringValue(entry.Encoded())
			if utils.IsStringNotEmpty(value) {
				postalAddress = append(postalAddress, value)
			}
		}
		signatureProductionPlace.SetPostalAddress(postalAddress)
	}
	return signatureProductionPlace
}

// CommitmentTypeIndications obtains the information concerning commitment type indication
// linked to the signature. Port of getCommitmentTypeIndications().
func (s *CAdESSignature) CommitmentTypeIndications() []*signature.CommitmentTypeIndication {
	attribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaEtsCommitmentType)
	if attribute == nil {
		return []*signature.CommitmentTypeIndication{}
	}

	if len(attribute.Values) == 0 {
		// Java: commitmentTypeIndications stays null, returned as-is by the try block.
		return nil
	}

	commitmentTypeIndications := make([]*signature.CommitmentTypeIndication, 0, len(attribute.Values))
	for _, value := range attribute.Values {
		if !value.IsUniversal(asn1ber.TagSequence) || !value.IsConstructed() {
			// Upstream logs "Unsupported type for CommitmentType : {}".
			continue
		}
		// CommitmentTypeIndication ::= SEQUENCE {
		//     commitmentTypeId         CommitmentTypeIdentifier,
		//     commitmentTypeQualifier  SEQUENCE OF CommitmentTypeQualifier OPTIONAL }
		if len(value.Children()) == 0 {
			// A malformed member makes CommitmentTypeIndication.getInstance throw upstream,
			// caught by the surrounding catch(Exception) which discards every result
			// gathered so far.
			return []*signature.CommitmentTypeIndication{}
		}
		commitmentTypeID, err := value.Children()[0].ObjectIdentifier()
		if err != nil {
			return []*signature.CommitmentTypeIndication{}
		}
		commitmentTypeIndications = append(commitmentTypeIndications,
			signature.NewCommitmentTypeIndication(commitmentTypeID.String()))
	}
	return commitmentTypeIndications
}

// SignedAssertions returns the list of embedded signed assertions.
// Port of getSignedAssertions().
func (s *CAdESSignature) SignedAssertions() []*signature.SignerRole {
	var result []*signature.SignerRole
	signerAttrV2 := s.signerAttributeV2()
	if signerAttrV2 != nil && signerAttrV2.SignedAssertions() != nil {
		for _, sa := range signerAttrV2.SignedAssertions().Assertions() {
			result = append(result, signature.NewSignerRole(sa.String(), enumerations.EndorsementTypeSigned))
		}
	}
	return result
}

// ClaimedSignerRoles returns the claimed roles of the signer. Port of getClaimedSignerRoles().
func (s *CAdESSignature) ClaimedSignerRoles() []*signature.SignerRole {
	signerAttr := s.signerAttributeV1()
	signerAttrV2 := s.signerAttributeV2()

	var claimedAttributes [][]byte
	if signerAttr != nil {
		claimedAttributes = signerAttr.claimedAttributes
	} else if signerAttrV2 != nil {
		claimedAttributes = signerAttrV2.ClaimedAttributes()
	} else {
		return []*signature.SignerRole{}
	}

	claimedRoles := []*signature.SignerRole{}
	for _, attributeDER := range claimedAttributes {
		claimedRoles = append(claimedRoles, cadesClaimedSignerRolesFromAttribute(attributeDER)...)
	}
	return claimedRoles
}

// cadesClaimedSignerRolesFromAttribute ports the private
// getClaimedSignerRoles(org.bouncycastle.asn1.x509.Attribute).
//
//	Attribute ::= SEQUENCE { type OBJECT IDENTIFIER, values SET OF AttributeValue }
func cadesClaimedSignerRolesFromAttribute(attributeDER []byte) []*signature.SignerRole {
	var claimedRoles []*signature.SignerRole
	element, rest, err := asn1ber.Parse(attributeDER)
	if err != nil || len(rest) != 0 || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() || len(element.Children()) != 2 {
		return claimedRoles
	}
	values := element.Children()[1]
	if !values.IsUniversal(asn1ber.TagSet) || !values.IsConstructed() {
		return claimedRoles
	}
	for _, value := range values.Children() {
		if value.IsASN1String() {
			claimedRoles = append(claimedRoles, signature.NewSignerRole(value.AsString(), enumerations.EndorsementTypeClaimed))
		}
	}
	return claimedRoles
}

// CertifiedSignerRoles returns the certified roles of the signer.
// Port of getCertifiedSignerRoles().
func (s *CAdESSignature) CertifiedSignerRoles() []*signature.SignerRole {
	signerAttr := s.signerAttributeV1()
	signerAttrV2 := s.signerAttributeV2()

	var certifiedAttributes [][]byte
	if signerAttr != nil {
		if signerAttr.certifiedAttributes != nil {
			certifiedAttributes = [][]byte{signerAttr.certifiedAttributes}
		}
	} else if signerAttrV2 != nil && signerAttrV2.CertifiedAttributes() != nil {
		for _, attributeCertificate := range signerAttrV2.CertifiedAttributes().AttributeCertificates() {
			certifiedAttributes = append(certifiedAttributes, attributeCertificate)
		}
	} else {
		return []*signature.SignerRole{}
	}

	roles := []*signature.SignerRole{}
	for _, attributeCertificateDER := range certifiedAttributes {
		roles = append(roles, cadesCertifiedSignerRolesFromAttributeCertificate(attributeCertificateDER)...)
	}
	return roles
}

// cadesCertifiedSignerRolesFromAttributeCertificate ports the private
// getCertifiedSignerRoles(AttributeCertificate); see the file header DEVIATION note for the
// scope of the RFC 5755 AttributeCertificate/RoleSyntax decoding performed here.
func cadesCertifiedSignerRolesFromAttributeCertificate(attributeCertificateDER []byte) []*signature.SignerRole {
	notBefore, notAfter, attributes, err := cadesAttributeCertificateValidityAndAttributes(attributeCertificateDER)
	if err != nil {
		// Upstream would let the parse failure propagate to the caller's catch(Exception),
		// which discards every result gathered across the whole certifiedAttributes loop; the
		// same behaviour is unreachable to reproduce exactly at this granularity (one
		// AttributeCertificate at a time) without restructuring the caller, so a malformed
		// AttributeCertificate is instead skipped on its own, matching every other
		// best-effort degradation in this file.
		return nil
	}

	var roles []*signature.SignerRole
	for _, attributeDER := range attributes {
		element, rest, perr := asn1ber.Parse(attributeDER)
		if perr != nil || len(rest) != 0 || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() || len(element.Children()) != 2 {
			continue
		}
		attrValue := cadesSingleAttributeValue(element)
		if attrValue == nil {
			// Upstream logs "Invalid encoding for a certified signer role attribute. Skip
			// processing." and returns null for the WHOLE method (dropping every role found
			// so far); reproduced here since it applies to a single AttributeCertificate.
			return nil
		}
		if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
			// Upstream logs "Unsupported type for RoleSyntax : {}".
			continue
		}
		roleName, ok := cadesRoleSyntaxRoleName(attrValue)
		if !ok {
			continue
		}
		certifiedRole := signature.NewSignerRole(roleName, enumerations.EndorsementTypeCertified)
		certifiedRole.SetNotBefore(cadesGeneralizedTime(notBefore))
		certifiedRole.SetNotAfter(cadesGeneralizedTime(notAfter))
		roles = append(roles, certifiedRole)
	}
	return roles
}

// cadesSingleAttributeValue mirrors DSSASN1Utils.getAsn1Encodable(Attribute) for an
// already-parsed X.509 Attribute element (attrType OID, attrValues SET OF).
func cadesSingleAttributeValue(attributeElement *asn1ber.Element) *asn1ber.Element {
	values := attributeElement.Children()[1]
	if !values.IsUniversal(asn1ber.TagSet) || !values.IsConstructed() || len(values.Children()) != 1 {
		return nil
	}
	return values.Children()[0]
}

// cadesAttributeCertificateValidityAndAttributes decodes the notBeforeTime/notAfterTime
// (AttCertValidityPeriod) and attributes (SEQUENCE OF Attribute) fields of an RFC 5755
// AttributeCertificate. holder/issuer/signature/serialNumber are opaque and skipped
// positionally - see the file header DEVIATION note.
//
//	AttributeCertificate     ::= SEQUENCE { acinfo AttributeCertificateInfo, ... }
//	AttributeCertificateInfo ::= SEQUENCE {
//	    version, holder, issuer, signature, serialNumber,
//	    attrCertValidityPeriod AttCertValidityPeriod,
//	    attributes             SEQUENCE OF Attribute,
//	    ... OPTIONAL }
//	AttCertValidityPeriod    ::= SEQUENCE { notBeforeTime GeneralizedTime, notAfterTime GeneralizedTime }
func cadesAttributeCertificateValidityAndAttributes(encoded []byte) (notBefore, notAfter *asn1ber.Element, attributes [][]byte, err error) {
	element, rest, perr := asn1ber.Parse(encoded)
	if perr != nil {
		return nil, nil, nil, perr
	}
	if len(rest) != 0 {
		return nil, nil, nil, model.NewDSSError("extra data found after AttributeCertificate")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() || len(element.Children()) < 1 {
		return nil, nil, nil, model.NewDSSError("malformed AttributeCertificate: not a SEQUENCE")
	}
	acInfo := element.Children()[0]
	if !acInfo.IsUniversal(asn1ber.TagSequence) || !acInfo.IsConstructed() || len(acInfo.Children()) < 7 {
		return nil, nil, nil, model.NewDSSError("malformed AttributeCertificateInfo")
	}
	children := acInfo.Children()

	validityPeriod := children[5]
	if !validityPeriod.IsUniversal(asn1ber.TagSequence) || !validityPeriod.IsConstructed() || len(validityPeriod.Children()) != 2 {
		return nil, nil, nil, model.NewDSSError("malformed AttCertValidityPeriod")
	}

	attrs := children[6]
	if !attrs.IsUniversal(asn1ber.TagSequence) || !attrs.IsConstructed() {
		return nil, nil, nil, model.NewDSSError("malformed AttributeCertificateInfo.attributes")
	}
	for _, attribute := range attrs.Children() {
		attributes = append(attributes, attribute.Encoded())
	}
	return validityPeriod.Children()[0], validityPeriod.Children()[1], attributes, nil
}

// cadesRoleSyntaxRoleName extracts the roleName of a RoleSyntax value, assuming it carries the
// common uniformResourceIdentifier GeneralName alternative; see the file header DEVIATION note.
//
//	RoleSyntax ::= SEQUENCE {
//	    roleAuthority [0] GeneralNames OPTIONAL,
//	    roleName      [1] GeneralName }
func cadesRoleSyntaxRoleName(roleSyntax *asn1ber.Element) (string, bool) {
	for _, child := range roleSyntax.Children() {
		if !child.IsContextSpecific(1) || !child.IsConstructed() || len(child.Children()) != 1 {
			continue
		}
		generalName := child.Children()[0]
		if generalName.IsContextSpecific(6) {
			return string(generalName.Content()), true
		}
		if generalName.IsASN1String() {
			return generalName.AsString(), true
		}
	}
	return "", false
}

// cadesGeneralizedTime decodes a GeneralizedTime element to a time.Time, the zero value (Java:
// null) on a malformed one.
func cadesGeneralizedTime(element *asn1ber.Element) time.Time {
	if element == nil || !element.IsUniversal(asn1ber.TagGeneralizedTime) {
		return time.Time{}
	}
	value, err := asn1ber.ParseGeneralizedTime(string(element.Content()))
	if err != nil {
		return time.Time{}
	}
	return value
}

// signerAttributeV1 ports the private getSignerAttributeV1(): SignerAttribute ::= SEQUENCE OF
// CHOICE { claimedAttributes [0] ClaimedAttributes, certifiedAttributes [1] CertifiedAttributes },
// ClaimedAttributes ::= SEQUENCE OF Attribute, CertifiedAttributes ::= AttributeCertificate (a
// single value, per RFC 5126) - i.e. org.bouncycastle.asn1.esf.SignerAttribute, which has no DSS
// port; see cadesLevelBaselineBSignerAttribute for the matching write side.
func (s *CAdESSignature) signerAttributeV1() *cadesSignerAttributeV1 {
	idAaEtsSignerAttr := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaEtsSignerAttr)
	if idAaEtsSignerAttr == nil {
		return nil
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(idAaEtsSignerAttr)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signer-attribute attribute. Skip processing.".
		return nil
	}
	result, err := parseCadesSignerAttributeV1(attrValue.Encoded())
	if err != nil {
		// Upstream logs "Unable to parse signerAttr - [{}]. Reason : {}".
		return nil
	}
	return result
}

// cadesSignerAttributeV1 is the parsed content of an RFC 5126/BouncyCastle
// org.bouncycastle.asn1.esf.SignerAttribute value.
type cadesSignerAttributeV1 struct {
	// claimedAttributes holds the DER encoding of each Attribute of the claimedAttributes
	// field, nil when the field is absent.
	claimedAttributes [][]byte
	// certifiedAttributes holds the DER encoding of the single certifiedAttributes
	// AttributeCertificate, nil when the field is absent.
	certifiedAttributes []byte
}

// parseCadesSignerAttributeV1 decodes a SignerAttribute from its encoding.
func parseCadesSignerAttributeV1(encoded []byte) (*cadesSignerAttributeV1, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, model.NewDSSError("extra data found after the SignerAttribute")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, model.NewDSSError("malformed SignerAttribute: not a SEQUENCE")
	}
	result := &cadesSignerAttributeV1{}
	for _, child := range element.Children() {
		if child.Class() != asn1ber.ClassContextSpecific || !child.IsConstructed() || len(child.Children()) != 1 {
			return nil, model.NewDSSError("illegal tag in SignerAttribute")
		}
		base := child.Children()[0]
		switch child.TagNumber() {
		case 0:
			if !base.IsUniversal(asn1ber.TagSequence) || !base.IsConstructed() {
				return nil, model.NewDSSError("malformed SignerAttribute.claimedAttributes: not a SEQUENCE")
			}
			for _, attribute := range base.Children() {
				result.claimedAttributes = append(result.claimedAttributes, attribute.Encoded())
			}
		case 1:
			result.certifiedAttributes = base.Encoded()
		default:
			return nil, model.NewDSSError(fmt.Sprintf("illegal tag: %d", child.TagNumber()))
		}
	}
	return result, nil
}

// signerAttributeV2 ports the private getSignerAttributeV2().
func (s *CAdESSignature) signerAttributeV2() *cms.SignerAttributeV2 {
	idAaEtsSignerAttrV2 := CAdESUtilsSignedAttribute(s.signerInformation, spi.OIDIdAaEtsSignerAttrV2)
	if idAaEtsSignerAttrV2 == nil {
		return nil
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(idAaEtsSignerAttrV2)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a signer-attribute-v2 attribute. Skip
		// processing.".
		return nil
	}
	result, err := cms.ParseSignerAttributeV2(attrValue.Encoded())
	if err != nil {
		// Upstream logs "Unable to parse signerAttrV2 : {}".
		return nil
	}
	return result
}

// EncryptionAlgorithm retrieves the encryption algorithm used for generating the signature.
// Port of getEncryptionAlgorithm().
func (s *CAdESSignature) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	oid := s.encryptionAlgOID()
	if encryptionAlgorithm, err := enumerations.EncryptionAlgorithmForOID(oid); err == nil {
		return encryptionAlgorithm
	}
	// fallback to identify via signature algorithm
	if signatureAlgorithm, err := enumerations.SignatureAlgorithmForOID(oid); err == nil {
		return signatureAlgorithm.EncryptionAlgorithm()
	}
	// Upstream logs "Unable to identify encryption algorithm for OID '{}'. Reason : {}".
	return ""
}

// DigestAlgorithm retrieves the digest algorithm used for generating the signature.
// Port of getDigestAlgorithm().
func (s *CAdESSignature) DigestAlgorithm() enumerations.DigestAlgorithm {
	signatureAlgorithm := s.encryptedDigestAlgo()
	if signatureAlgorithm != "" {
		if enumerations.EncryptionAlgorithmRSASSAPSS == signatureAlgorithm.EncryptionAlgorithm() {
			return s.pssHashAlgorithm()
		}
		return signatureAlgorithm.DigestAlgorithm()
	}
	return cadesDigestAlgorithmForOID(s.digestAlgOID())
}

// encryptedDigestAlgo ports the private getEncryptedDigestAlgo().
func (s *CAdESSignature) encryptedDigestAlgo() enumerations.SignatureAlgorithm {
	signatureAlgorithm, err := enumerations.SignatureAlgorithmForOID(s.encryptionAlgOID())
	if err != nil {
		return ""
	}
	return signatureAlgorithm
}

// pssHashAlgorithm ports the private getPSSHashAlgorithm().
func (s *CAdESSignature) pssHashAlgorithm() enumerations.DigestAlgorithm {
	encryptionAlgParams := s.encryptionAlgParams()
	if utils.IsArrayNotEmpty(encryptionAlgParams) && !bytes.Equal(asn1ber.DERNull, encryptionAlgParams) {
		hashAlgorithm, err := cadesRSASSAPSSParamsHashAlgorithm(encryptionAlgParams)
		if err != nil {
			// Upstream logs "Unable to analyze EncryptionAlgParams".
			return ""
		}
		return cadesDigestAlgorithmForOID(hashAlgorithm.Algorithm.String())
	}
	return ""
}

// cadesRSASSAPSSParamsHashAlgorithm decodes the hashAlgorithm [0] EXPLICIT field of an
// RSASSA-PSS-params SEQUENCE (RFC 4055 / PKCS#1), defaulting to SHA-1 - the ASN.1 DEFAULT - when
// the field is absent, i.e. org.bouncycastle.asn1.pkcs.RSASSAPSSparams#getHashAlgorithm, which
// has no DSS port.
//
//	RSASSA-PSS-params ::= SEQUENCE {
//	    hashAlgorithm [0] AlgorithmIdentifier DEFAULT sha1Identifier, ... }
func cadesRSASSAPSSParamsHashAlgorithm(encoded []byte) (*asn1ber.AlgorithmIdentifier, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, model.NewDSSError("extra data found after RSASSA-PSS-params")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, model.NewDSSError("malformed RSASSA-PSS-params: not a SEQUENCE")
	}
	for _, child := range element.Children() {
		if child.IsContextSpecific(0) {
			if !child.IsConstructed() || len(child.Children()) != 1 {
				return nil, model.NewDSSError("malformed RSASSA-PSS-params.hashAlgorithm")
			}
			return asn1ber.AlgorithmIdentifierFromElement(child.Children()[0])
		}
	}
	return spi.DSSASN1UtilsAlgorithmIdentifierForDigest(enumerations.DigestAlgorithmSHA1)
}

// encryptionAlgOID ports SignerInformation#getEncryptionAlgOID(): the OID of the
// signatureAlgorithm field, exactly as encoded (RFC 3852: "a 'signature algorithm' (encryption
// + digest algorithms)").
func (s *CAdESSignature) encryptionAlgOID() string {
	return s.signerInformation.SignatureAlgorithm.Algorithm.String()
}

// digestAlgOID ports SignerInformation#getDigestAlgOID(): the OID of the digestAlgorithm field.
func (s *CAdESSignature) digestAlgOID() string {
	return s.signerInformation.DigestAlgorithm.Algorithm.String()
}

// encryptionAlgParams ports SignerInformation#getEncryptionAlgParams(): the DER encoding of the
// signatureAlgorithm's parameters, nil when absent.
func (s *CAdESSignature) encryptionAlgParams() []byte {
	return s.signerInformation.SignatureAlgorithm.Parameters
}

// cadesDigestAlgorithmForOID ports the private getDigestAlgorithmForOID(String).
func cadesDigestAlgorithmForOID(oid string) enumerations.DigestAlgorithm {
	if utils.IsStringEmpty(oid) {
		// Upstream logs "DigestAlgorithm cannot be defined with an empty OID! Skip
		// processing.".
		return ""
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(oid)
	if err != nil {
		// Upstream logs "Unable to identify DigestAlgorithm for OID '{}'. Reason : {}".
		return ""
	}
	return digestAlgorithm
}

// SignatureAlgorithm retrieves the signature algorithm (or cipher) used for generating the
// signature. Port of getSignatureAlgorithm().
func (s *CAdESSignature) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	return enumerations.SignatureAlgorithmGetAlgorithm(s.EncryptionAlgorithm(), s.DigestAlgorithm())
}

// CheckSignatureIntegrity verifies the signature integrity; checks if the signed content has
// not been tampered with. Port of checkSignatureIntegrity().
//
// The early "already computed" return reads s.cachedCryptoVerification rather than the base's
// own (inaccessible, cross-package-private) signatureCryptographicVerification field, the way
// Java's own method reads its protected field directly; see that field's doc comment.
func (s *CAdESSignature) CheckSignatureIntegrity() {
	if s.cachedCryptoVerification != nil {
		return
	}
	verification := signature.NewSignatureCryptographicVerification()
	s.cachedCryptoVerification = verification
	s.SetSignatureCryptographicVerification(verification)

	var signerInformationToCheck *cmscore.SignerInfo
	if s.cmsDocument.IsDetachedSignature() && !s.IsCounterSignature() {
		if utils.IsCollectionEmpty(s.DetachedContents()) {
			verification.SetErrorMessage("Detached file not found!")
			s.ReferenceValidationsForSignerInformation(signerInformationToCheck)
			return
		}
		recreated, err := s.recreateSignerInformation()
		if err != nil {
			// catch (CMSException e)
			verification.SetErrorMessage(err.Error())
			return
		}
		signerInformationToCheck = recreated
	} else {
		signerInformationToCheck = s.signerInformation
	}

	candidatesForSigningCertificate := s.CandidatesForSigningCertificate()

	signedContent, err := s.signedContentForIntegrityCheck(signerInformationToCheck)
	if err != nil {
		verification.SetErrorMessage(err.Error())
		return
	}

	// Computed up front (and reused below) so it can gate CAdESSignatureIntegrityValidator.Verify
	// the way BouncyCastle's SignerInformation#verify itself does: BC recomputes the digest of
	// the associated content and compares it against the message-digest signed attribute BEFORE
	// checking the raw signature bytes, throwing CMSSignerDigestMismatchException - and thus
	// failing every candidate uniformly - on a mismatch. See the contentDigestMismatch field doc
	// on CAdESSignatureIntegrityValidator for the fixture this fixes.
	refValidations := s.ReferenceValidationsForSignerInformation(signerInformationToCheck)
	contentDigestMismatch := false
	for _, referenceValidation := range refValidations {
		if referenceValidation.Type() == enumerations.DigestMatcherTypeMessageDigest &&
			referenceValidation.IsFound() && !referenceValidation.IsIntact() {
			contentDigestMismatch = true
		}
	}

	signingCertificateValidator := NewCAdESSignatureIntegrityValidator(signerInformationToCheck, signedContent, contentDigestMismatch)
	certificateValidity := signingCertificateValidator.Validate(candidatesForSigningCertificate)
	if certificateValidity != nil {
		if err := candidatesForSigningCertificate.SetTheCertificateValidity(certificateValidity); err != nil {
			panic(err)
		}
	}

	verification.SetErrorMessages(signingCertificateValidator.ErrorMessages())
	verification.SetSignatureIntact(certificateValidity != nil)

	referenceDataFound := true
	referenceDataIntact := true
	for _, referenceValidation := range refValidations {
		referenceDataFound = referenceDataFound && referenceValidation.IsFound()
		referenceDataIntact = referenceDataIntact && referenceValidation.IsIntact()
	}
	verification.SetReferenceDataFound(referenceDataFound)
	verification.SetReferenceDataIntact(referenceDataIntact)
}

// signedContentForIntegrityCheck computes the bytes CAdESSignatureIntegrityValidator.Verify
// checks the signature against, matching that file's own documented contract (see
// cades_signature_integrity_validator.go's header): the DER encoding of the SignedAttrs value
// as a SET OF (RFC 5652 clause 5.4) when signerInformationToCheck carries signed attributes -
// every CAdES baseline profile requires them - or the signed content itself otherwise (BC's
// SignerInformation#verify falls back to the CMSSignedData's own encapsulated/detached content
// in that case, a back-reference cmscore.SignerInfo does not carry).
func (s *CAdESSignature) signedContentForIntegrityCheck(signerInformationToCheck *cmscore.SignerInfo) ([]byte, error) {
	if signerInformationToCheck.HasSignedAttributes() {
		return signerInformationToCheck.SignedAttributes.DERSetEncoded(), nil
	}
	originalDocument, err := s.SignerDocumentContent()
	if err != nil {
		return nil, err
	}
	if originalDocument == nil {
		return nil, model.NewDSSError("Detached file not found!")
	}
	return spi.DSSUtilsToByteArrayOfDocument(originalDocument)
}

// ReferenceValidationsForSignerInformation returns the reference validation.
// Port of the public getReferenceValidations(SignerInformation).
func (s *CAdESSignature) ReferenceValidationsForSignerInformation(signerInformationToCheck *cmscore.SignerInfo) []*model.ReferenceValidation {
	if s.CachedReferenceValidations() == nil {
		originalDocument, err := s.SignerDocumentContent()
		if err != nil {
			// Upstream logs "Original document not found".
			originalDocument = nil
		}

		var refValidation *model.ReferenceValidation
		messageDigestValue := s.MessageDigestValue()
		if messageDigestValue != nil {
			refValidation = s.messageDigestReferenceValidation(originalDocument, messageDigestValue)
		} else {
			// Upstream logs "message-digest is not present in SignedData! Extracting digests
			// from content SignatureValue...".
			refValidation = s.contentReferenceValidation(originalDocument, signerInformationToCheck)
		}
		s.SetCachedReferenceValidations([]*model.ReferenceValidation{refValidation})
	}
	return s.CachedReferenceValidations()
}

// SignerDocumentContent extracts a document content that was signed.
//
// NOTE: Some differences are possible with PAdES.
//
// Port of the protected getSignerDocumentContent().
func (s *CAdESSignature) SignerDocumentContent() (model.DSSDocument, error) {
	return s.OriginalDocument()
}

// verifyDigestAlgorithm ports the private verifyDigestAlgorithm(DSSDocument,
// Set<DigestAlgorithm>, Digest).
func (s *CAdESSignature) verifyDigestAlgorithm(originalDocument model.DSSDocument, messageDigestAlgorithms []enumerations.DigestAlgorithm, messageDigest *model.Digest) bool {
	if utils.IsCollectionNotEmpty(messageDigestAlgorithms) {
		for _, digestAlgorithm := range messageDigestAlgorithms {
			base64Digest, err := originalDocument.DigestValue(digestAlgorithm)
			if err != nil {
				continue
			}
			if bytes.Equal(messageDigest.Value(), base64Digest) {
				messageDigest.SetAlgorithm(digestAlgorithm)
				return true
			}
		}
	} else {
		// Upstream logs "Message DigestAlgorithms not found in SignedData! Reference
		// validation is not possible.".
	}
	return false
}

// getManifestEntryValidation ports the private getManifestEntryValidation().
func (s *CAdESSignature) getManifestEntryValidation() []*model.ReferenceValidation {
	manifestEntryValidations := []*model.ReferenceValidation{}
	manifestFile := s.ManifestFile()
	if manifestFile == nil {
		// Upstream logs "No related manifest file found for a signature with name [{}]".
		return manifestEntryValidations
	}
	for _, entry := range manifestFile.Entries() {
		entryValidation := model.NewReferenceValidation()
		entryValidation.SetType(enumerations.DigestMatcherTypeManifestEntry)
		entryValidation.SetUri(entry.Uri())
		entryValidation.SetDocument(entry.Document())
		entryValidation.SetDigest(entry.Digest())
		entryValidation.SetFound(entry.IsFound())
		entryValidation.SetIntact(entry.IsIntact())
		manifestEntryValidations = append(manifestEntryValidations, entryValidation)
	}
	return manifestEntryValidations
}

// ReferenceValidations returns individual validation for each reference (XAdES, JAdES) or for
// the message-imprint (CAdES). Port of getReferenceValidations().
func (s *CAdESSignature) ReferenceValidations() []*model.ReferenceValidation {
	s.CheckSignatureIntegrity()
	return s.CachedReferenceValidations()
}

// messageDigestReferenceValidation verifies a message-digest of a CMS, when applicable.
// Port of the private getMessageDigestReferenceValidation(DSSDocument, byte[]).
//
// Java's Digest is a mutable object shared by reference between this method's local variable
// and messageDigestValidation once setDigest(messageDigest) is called, so every later
// messageDigest.setAlgorithm(...) - including the one inside verifyDigestAlgorithm - is
// automatically visible through messageDigestValidation.getDigest() too. model.Digest is a Go
// value type with no such aliasing, so SetDigest is called again after every mutation that must
// be observable afterward, reproducing the same externally-visible state at each point Java's
// object graph would show it.
func (s *CAdESSignature) messageDigestReferenceValidation(originalDocument model.DSSDocument, messageDigestValue []byte) *model.ReferenceValidation {
	messageDigestValidation := model.NewReferenceValidation()
	messageDigestValidation.SetType(enumerations.DigestMatcherTypeMessageDigest)

	messageDigest := model.NewDigest("", messageDigestValue)

	var digestAlgorithmCandidates []enumerations.DigestAlgorithm
	signerInformationDigestAlgorithm := s.DigestAlgorithm()
	if signerInformationDigestAlgorithm != "" {
		digestAlgorithmCandidates = cadesAppendDigestAlgorithm(digestAlgorithmCandidates, signerInformationDigestAlgorithm)
	}
	for _, digestAlgorithm := range s.MessageDigestAlgorithms() {
		digestAlgorithmCandidates = cadesAppendDigestAlgorithm(digestAlgorithmCandidates, digestAlgorithm)
	}

	if len(digestAlgorithmCandidates) == 1 {
		messageDigest.SetAlgorithm(digestAlgorithmCandidates[0])
	}
	messageDigestValidation.SetDigest(messageDigest)

	if originalDocument != nil {
		messageDigestValidation.SetDocument(originalDocument)
		messageDigestValidation.SetFound(true)
		intact := s.verifyDigestAlgorithm(originalDocument, digestAlgorithmCandidates, &messageDigest)
		messageDigestValidation.SetIntact(intact)
		messageDigestValidation.SetDigest(messageDigest)

		manifestFile := s.ManifestFile()
		if manifestFile != nil {
			manifestDigestValue, err := manifestFile.DigestValue(messageDigest.Algorithm())
			if err == nil && bytes.Equal(messageDigest.Value(), manifestDigestValue) {
				// get references to documents contained in the manifest file (for ASiC-E
				// container)
				//
				// GAP flagged for integrator: model.ReferenceValidation (frozen, dss-model) has
				// no mutator for its dependentReferenceValidations field beyond the lazily
				// initialising DependentValidations() getter (see model/reference_validation.go),
				// so the manifest-entry validations gathered by getManifestEntryValidation()
				// cannot be attached in place the way Java's mutable
				// getDependentValidations().addAll(...) does; model/... is frozen per PORTING.md
				// and is not edited here. A one-method, additive accessor
				//
				//	func (r *ReferenceValidation) SetDependentValidations(dependentReferenceValidations []*ReferenceValidation) {
				//		r.dependentReferenceValidations = dependentReferenceValidations
				//	}
				//
				// would close the gap (used the moment it lands: replace the two lines below
				// with messageDigestValidation.SetDependentValidations(append(
				// messageDigestValidation.DependentValidations(),
				// s.getManifestEntryValidation()...))). Until then the manifest-entry
				// dependent validations are computed (real logic, not a stub) but not
				// observable through messageDigestValidation - an ASiC-E ".getDependentValidations()"
				// caller of this reference validation gets an empty list rather than upstream's
				// populated one.
				_ = s.getManifestEntryValidation()
			}
		}
	} else {
		// Upstream logs "The original document is not found or cannot be extracted.
		// Reference validation is not possible.".
	}
	return messageDigestValidation
}

// cadesAppendDigestAlgorithm appends digestAlgorithm to candidates if not already present,
// reproducing java.util.Set<DigestAlgorithm> insertion semantics with deterministic
// (first-seen) ordering rather than Go's nondeterministic map iteration.
func cadesAppendDigestAlgorithm(candidates []enumerations.DigestAlgorithm, digestAlgorithm enumerations.DigestAlgorithm) []enumerations.DigestAlgorithm {
	for _, candidate := range candidates {
		if candidate == digestAlgorithm {
			return candidates
		}
	}
	return append(candidates, digestAlgorithm)
}

// contentReferenceValidation verifies a content digest, when applicable.
// Port of the private getContentReferenceValidation(DSSDocument, SignerInformation).
//
// See the file header DEVIATION note: this port has no equivalent of BC's
// SignerInformation#getContentDigest() and always reports "not found" here.
func (s *CAdESSignature) contentReferenceValidation(_ model.DSSDocument, _ *cmscore.SignerInfo) *model.ReferenceValidation {
	contentValidation := model.NewReferenceValidation()
	contentValidation.SetType(enumerations.DigestMatcherTypeContentDigest)
	return contentValidation
}

// buildSignatureDigestReference returns a signature reference element as defined in
// TS 119 442 - V1.1.1 - Electronic Signatures and Infrastructures (ESI), ch. 5.1.4.2.1.3 XML
// component: in case of CAdES signatures, the input to the digest value computation shall be
// one of the DER-encoded instances of SignedInfo type present within the CMS structure.
// Port of buildSignatureDigestReference(DigestAlgorithm).
func (s *CAdESSignature) BuildSignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.SignatureDigestReference {
	derEncodedSignerInfo := s.signerInformation.DER()
	digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, derEncodedSignerInfo)
	if err != nil {
		panic(err)
	}
	return signature.NewSignatureDigestReference(model.NewDigest(digestAlgorithm, digestValue))
}

// DataToBeSignedRepresentation returns the DTBSR, which is then used to create the signature.
// Port of getDataToBeSignedRepresentation().
func (s *CAdESSignature) DataToBeSignedRepresentation() model.Digest {
	referenceValidations := s.ReferenceValidations()
	// only one is allowed for CMS
	referenceValidation := referenceValidations[0]
	switch referenceValidation.Type() {
	case enumerations.DigestMatcherTypeMessageDigest:
		digestAlgorithm := s.DigestAlgorithm()
		if digestAlgorithm != "" {
			signedAttributes := CAdESUtilsSignedAttributes(s.signerInformation)
			derEncoded := signedAttributes.DERSetEncoded()
			digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, derEncoded)
			if err != nil {
				panic(err)
			}
			return model.NewDigest(digestAlgorithm, digestValue)
		}
		return model.Digest{}
	case enumerations.DigestMatcherTypeContentDigest:
		return referenceValidation.Digest()
	default:
		panic(model.NewDSSError(fmt.Sprintf("The found referenceValidation type '%s' is not supported! "+
			"Unable to compute DTBSR.", referenceValidation.Type())))
	}
}

// recreateSignerInformation recreates a SignerInformation with the content using a CMSParser.
// Port of the private recreateSignerInformation(), declared "throws CMSException".
func (s *CAdESSignature) recreateSignerInformation() (*cmscore.SignerInfo, error) {
	dssDocument := s.DetachedContents()[0] // only one element for CAdES Signature
	digestCalculatorProvider := cms.NewPrecomputedDigestCalculatorProvider(dssDocument)
	return cms.CMSUtilsRecomputeSignerInformation(s.cmsDocument, s.SignerId(), digestCalculatorProvider, CAdESUtilsDefaultResourcesHandlerBuilder)
}

// MessageDigestAlgorithms returns a set of used DigestAlgorithms incorporated into the CMS.
// Port of the public getMessageDigestAlgorithms().
func (s *CAdESSignature) MessageDigestAlgorithms() []enumerations.DigestAlgorithm {
	var result []enumerations.DigestAlgorithm
	for _, algorithmIdentifier := range s.cmsDocument.DigestAlgorithmIDs() {
		digestAlgorithm := cadesDigestAlgorithmForOID(algorithmIdentifier.Algorithm.String())
		if digestAlgorithm != "" {
			result = cadesAppendDigestAlgorithm(result, digestAlgorithm)
		}
	}
	return result
}

// MessageDigestValue returns a digest value incorporated in an attribute "message-digest" in
// CMS Signed Data. Port of the public getMessageDigestValue().
func (s *CAdESSignature) MessageDigestValue() []byte {
	messageDigestAttribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDPkcs9AtMessageDigest)
	if messageDigestAttribute == nil {
		return nil
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(messageDigestAttribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a message-digest attribute. Skip processing.".
		return nil
	}
	if !attrValue.IsUniversal(asn1ber.TagOctetString) {
		// Upstream logs "Message-digest attribute shall be an instance of type
		// ASN1OctetString. Found : {}".
		return nil
	}
	return attrValue.Octets()
}

// ContentType returns the value of the signed attribute content-type.
// Port of getContentType().
func (s *CAdESSignature) ContentType() string {
	contentTypeAttribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDPkcs9AtContentType)
	if contentTypeAttribute == nil {
		return ""
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(contentTypeAttribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a content-type attribute. Skip processing.".
		return ""
	}
	if !attrValue.IsUniversal(asn1ber.TagOID) {
		// Upstream logs "content-type attribute shall be an instance of type
		// ASN1ObjectIdentifier. Found : {}".
		return ""
	}
	oid, err := attrValue.ObjectIdentifier()
	if err != nil {
		return ""
	}
	return oid.String()
}

// MimeType returns the value of the signed attribute mime-type. Port of getMimeType().
func (s *CAdESSignature) MimeType() string {
	mimeTypeAttribute := CAdESUtilsSignedAttribute(s.signerInformation, spi.OIDIdAaEtsMimeType)
	if mimeTypeAttribute == nil {
		return ""
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(mimeTypeAttribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a mime-type attribute. Skip processing.".
		return ""
	}
	return spi.DSSASN1UtilsString(attrValue.Encoded())
}

// SignatureType returns the value of the signature type protected header (JAdES, CB-AdES).
// Port of getSignatureType(); not supported for CAdES.
func (s *CAdESSignature) SignatureType() string {
	return ""
}

// ContentIdentifier gets ContentIdentifier String. Port of the public getContentIdentifier().
func (s *CAdESSignature) ContentIdentifier() string {
	contentIdentifierAttribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaContentIdentifier)
	if contentIdentifierAttribute == nil {
		return ""
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(contentIdentifierAttribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a content-identifier attribute. Skip
		// processing.".
		return ""
	}
	// ContentIdentifier ::= OCTET STRING (org.bouncycastle.asn1.ess.ContentIdentifier)
	result, err := spi.DSSASN1UtilsToString(attrValue.Encoded())
	if err != nil {
		return ""
	}
	return result
}

// ContentHints gets Content Hints. Port of the public getContentHints().
func (s *CAdESSignature) ContentHints() string {
	contentHintAttribute := CAdESUtilsSignedAttribute(s.signerInformation, OIDIdAaContentHint)
	if contentHintAttribute == nil {
		return ""
	}
	attrValue := spi.DSSASN1UtilsAsn1Encodable(contentHintAttribute)
	if attrValue == nil {
		// Upstream logs "Invalid encoding for a content-hint attribute. Skip processing.".
		return ""
	}

	// ContentHints ::= SEQUENCE { contentDescription UTF8String OPTIONAL, contentType ContentType }
	if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
		// Upstream logs "Unable to parse ContentHints - [{}]. Reason : {}".
		return ""
	}
	children := attrValue.Children()
	var contentTypeElement *asn1ber.Element
	var description string
	hasDescription := false
	switch len(children) {
	case 1:
		contentTypeElement = children[0]
	case 2:
		if !children[0].IsUniversal(asn1ber.TagUTF8String) {
			return ""
		}
		description = children[0].AsString()
		hasDescription = true
		contentTypeElement = children[1]
	default:
		return ""
	}
	if !contentTypeElement.IsUniversal(asn1ber.TagOID) {
		return ""
	}
	oid, err := contentTypeElement.ObjectIdentifier()
	if err != nil {
		return ""
	}
	contentHint := oid.String()
	if hasDescription {
		contentHint += " [" + description + "]"
	}
	return contentHint
}

// SignerInformation gets a SignedInformation. Port of the public getSignerInformation().
func (s *CAdESSignature) SignerInformation() *cmscore.SignerInfo {
	return s.signerInformation
}

// SignatureValue returns the digital signature value. Port of getSignatureValue().
func (s *CAdESSignature) SignatureValue() []byte {
	return s.signerInformation.Signature
}

// IsCounterSignature checks if the current signature is a counter signature (i.e. has a master
// signature). Port of isCounterSignature(); see the file header DEVIATION note.
func (s *CAdESSignature) IsCounterSignature() bool {
	return s.MasterSignature() != nil
}

// CounterSignatures returns a list of counter signatures applied to this signature.
// Port of getCounterSignatures().
func (s *CAdESSignature) CounterSignatures() []validation.AdvancedSignature {
	if s.CachedCounterSignatures() != nil {
		return s.CachedCounterSignatures()
	}

	var counterSignatures []validation.AdvancedSignature
	for _, counterSignerInformation := range s.CounterSignatureStore().SignerInfos() {
		counterSignature := NewCAdESSignature(s.cmsDocument, counterSignerInformation)
		counterSignature.SetFilename(s.Filename())
		counterSignature.SetMasterSignature(s)
		counterSignatures = append(counterSignatures, counterSignature)
	}
	s.SetCachedCounterSignatures(counterSignatures)
	return counterSignatures
}

// CAdESSignerInformationStore is this port's counterpart of BC's SignerInformationStore for the
// narrow purpose CounterSignatureStore below serves: a collection of SignerInfos exposing the
// single SignerInfos() accessor CAdESSignatureIdentifierBuilder needs. Building a full *cms.CMS
// instead would need a complete SignedData this counter-signature list does not have, exactly as
// Java's own SignerInformationStore carries no CMSSignedData of its own either; sharing the
// SignerInfos() method name is what CAdESSignatureIdentifierBuilder actually needs, so it is
// reproduced here on a lightweight named slice type instead.
type CAdESSignerInformationStore []*cmscore.SignerInfo

// SignerInfos returns the underlying SignerInfo slice. Port of SignerInformationStore#getSigners().
func (st CAdESSignerInformationStore) SignerInfos() []*cmscore.SignerInfo { return st }

// CounterSignatureStore returns a SignerInformationStore containing counter signatures.
// Port of the protected getCounterSignatureStore().
func (s *CAdESSignature) CounterSignatureStore() CAdESSignerInformationStore {
	if s.counterSignaturesStore == nil {
		s.counterSignaturesStore = cadesCounterSignaturesOf(s.signerInformation)
	}
	return CAdESSignerInformationStore(s.counterSignaturesStore)
}

// cadesCounterSignaturesOf ports SignerInformationStore SignerInformation#getCounterSignatures():
// the SignerInfos parsed from the id-countersignature unsigned attribute of signerInformation.
func cadesCounterSignaturesOf(signerInformation *cmscore.SignerInfo) []*cmscore.SignerInfo {
	if signerInformation == nil || !signerInformation.HasUnsignedAttributes() {
		return nil
	}
	var counterSignatures []*cmscore.SignerInfo
	for _, attribute := range signerInformation.UnsignedAttributes.GetAll(cmscore.OIDCounterSignature) {
		for _, value := range attribute.Values {
			if !value.IsUniversal(asn1ber.TagSequence) || !value.IsConstructed() {
				continue
			}
			counterSignerInfo, err := cmscore.SignerInfoFromElement(value)
			if err != nil {
				continue
			}
			counterSignatures = append(counterSignatures, counterSignerInfo)
		}
	}
	return counterSignatures
}

// OriginalDocument returns the original signed document. Port of the public
// getOriginalDocument().
//
// Panics with the underlying error's message on failure: Java's DSSException is unchecked and
// AdvancedSignature has no accessor named getOriginalDocument that could return an error; the
// only caller within this file (SignerDocumentContent) already forwards the error instead.
func (s *CAdESSignature) OriginalDocument() (model.DSSDocument, error) {
	// RFC 5652 ch 11.4.
	if s.IsCounterSignature() {
		return model.NewInMemoryDocument(s.MasterSignature().SignatureValue()), nil
	}
	return CAdESUtilsOriginalDocument(s.cmsDocument, s.DetachedContents())
}

// SignatureIdentifierBuilder returns a builder to define and build a signature Id.
// Port of the protected getSignatureIdentifierBuilder().
func (s *CAdESSignature) SignatureIdentifierBuilder() validation.SignatureIdentifierBuilder {
	return NewCAdESSignatureIdentifierBuilder(s)
}

// DAIdentifier returns an identifier provided by the Driving Application (DA); not applicable
// for CAdES. Port of getDAIdentifier().
func (s *CAdESSignature) DAIdentifier() string {
	return ""
}

// SignerInformationStoreInfos returns a Set of CertificateIdentifier extracted from a
// SignerInformationStore of CMS Signed Data. Port of the public getSignerInformationStoreInfos().
func (s *CAdESSignature) SignerInformationStoreInfos() []*spi.SignerIdentifier {
	return s.CertificateSource().AllCertificateIdentifiers()
}

// AddExternalTimestamp allows adding an external timestamp. The given timestamp must be
// processed before. NOTE: supported only for CAdES signatures. Port of
// addExternalTimestamp(TimestampToken).
//
// Panics with the Java message when the timestamp was not validated first (Java's DSSException
// is unchecked and AddExternalTimestamp has no error return per the AdvancedSignature interface).
func (s *CAdESSignature) AddExternalTimestamp(timestamp *validation.TimestampToken) {
	if !timestamp.IsProcessed() {
		panic("Timestamp token must be validated first !")
	}
	s.TimestampSource().AddExternalTimestamp(timestamp)
}

// DataFoundUpToLevel returns the signature level. Port of getDataFoundUpToLevel().
func (s *CAdESSignature) DataFoundUpToLevel() enumerations.SignatureLevel {
	if !s.HasBESProfile() {
		return enumerations.SignatureLevelCMSNotETSI
	}

	baselineProfile := s.HasBProfile()

	if !s.HasExtendedTProfile() {
		if baselineProfile {
			return enumerations.SignatureLevelCAdESBaselineB
		} else if s.HasEPESProfile() {
			return enumerations.SignatureLevelCAdESEPES
		}
		return enumerations.SignatureLevelCAdESBES
	}

	baselineProfile = baselineProfile && s.HasTProfile()

	if baselineProfile && s.HasLTProfile() {
		if s.HasERSProfile() {
			return enumerations.SignatureLevelCAdESERS
		}
		if s.HasLTAProfile() {
			return enumerations.SignatureLevelCAdESBaselineLTA
		}
		return enumerations.SignatureLevelCAdESBaselineLT

	} else if s.HasCProfile() {
		if s.HasXLProfile() {
			if s.HasERSProfile() {
				return enumerations.SignatureLevelCAdESERS
			}
			if s.HasAProfile() {
				return enumerations.SignatureLevelCAdESA
			}
			if s.HasXProfile() {
				return enumerations.SignatureLevelCAdESXL
			}
		}
		if s.HasXProfile() {
			return enumerations.SignatureLevelCAdESX
		}
		return enumerations.SignatureLevelCAdESC

	} else if s.HasXLProfile() {
		if s.HasERSProfile() {
			return enumerations.SignatureLevelCAdESERS
		}
		if s.HasAProfile() {
			// CAdES-E-A can be built on CAdES-E-T directly
			return enumerations.SignatureLevelCAdESA
		}
		return enumerations.SignatureLevelCAdESLT
	}

	if baselineProfile {
		return enumerations.SignatureLevelCAdESBaselineT
	}
	return enumerations.SignatureLevelCAdEST
}

// BaselineRequirementsChecker returns the cached instance of the CAdESBaselineRequirementsChecker.
// Port of the protected covariant-return override getBaselineRequirementsChecker().
func (s *CAdESSignature) BaselineRequirementsChecker() *CAdESBaselineRequirementsChecker {
	return s.DefaultAdvancedSignature.BaselineRequirementsChecker().(*CAdESBaselineRequirementsChecker)
}

// CreateBaselineRequirementsChecker instantiates a BaselineRequirementsChecker according to the
// signature format. Port of the protected createBaselineRequirementsChecker(CertificateVerifier).
func (s *CAdESSignature) CreateBaselineRequirementsChecker(certificateVerifier validation.CertificateVerifier) validation.BaselineRequirementsCheckerContract {
	return NewCAdESBaselineRequirementsChecker(s, certificateVerifier)
}
