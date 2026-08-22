// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CadesLevelBaselineLTATimestampExtractor.java (DSS 6.5.RC1).
//
// This file owns the byte-exact side of the CAdES archive time-stamp: the ats-hash-index
// attribute in all three of its versions and the archive-time-stamp-v3 message imprint. Both
// are digested by a TSA and later re-computed by a verifier, so every octet produced here is
// part of the format, not an implementation detail; cades_level_baseline_lta_timestamp_extractor_kat_test.go
// pins them against upstream DSS 6.5.RC1 itself (see testdata/gen/AtsHashIndexOracle.java).
//
// BouncyCastle replacements used here (see PORTING.md):
//
//   - org.bouncycastle.asn1.cms.Attribute / AttributeTable -> cmscore.Attribute / cmscore.Attributes.
//   - org.bouncycastle.cms.SignerInformation and its toASN1Structure() -> *cmscore.SignerInfo,
//     which is both at once in this port.
//   - ASN1Sequence -> its encoding, []byte, the representation dss-spi already settled on for
//     the ats-hash-index tables (DSSASN1Utils.getDEROctetStrings, getAlgorithmIdentifier).
//   - DERSequence / DERSet / DEROctetString / DERTaggedObject -> the asn1ber writers; the
//     components of a DERSequence are normalised to DER on the way in, exactly as
//     BouncyCastle's DER output does, see cadesLTAToDER.
//   - org.bouncycastle.util.Store<X509CRLHolder> and the two OCSP stores of CMSSignedData ->
//     CMS.CRLs(), CMS.OcspResponseStore() and CMS.OcspBasicStore(), each already the encoding
//     of the member Java re-encodes with getEncoded().
//
// The one non-obvious thing the bytes depend on is the order AttributeTable#toASN1EncodableVector()
// hands the unsigned attributes out in, which is NOT their order in the SignerInfo but the order
// of the java.util.Hashtable BouncyCastle groups them by attrType in. It is reproduced by
// cadesLTAAttributeTableOrder at the bottom of this file, and pinned by its own Java oracle.
//
// DEVIATIONS:
//
//   - An archive time-stamp carrying no ats-hash-index at all makes upstream's
//     getVerifiedAtsHashIndex raise a NullPointerException while encoding
//     `new Attribute(null, ...)`. Here the data-dependent throw becomes an error, raised after
//     the ArchiveTimestampHashIndexStatus has been set on the token exactly as upstream sets it.
//   - getUnsignedAttributesHashIndex reads SignerInformation#getUnsignedAttributes() directly,
//     so a SignerInfo carrying no unsignedAttrs field raises a NullPointerException upstream.
//     Here it yields an empty SEQUENCE instead. The case cannot arise in the augmentation flow
//     that calls it (an archive time-stamp is only ever added to a signature that already
//     carries at least a signature-time-stamp), and no byte of a reachable output changes.
//   - The extractor holds the narrow cadesLTASignature view of Signature rather than the
//     concrete type, so that the byte-exact core can be driven by the Java-generated fixtures
//     without a full document analyzer. *Signature is its only production implementation
//     and the exported constructor still takes it.
//   - slf4j logging is dropped (PORTING.md); the LOG.warn calls that accompany an
//     ArchiveTimestampHashIndexStatus error message are covered by that message itself.
package cades

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// cadesLTASignature is the part of *CAdESSignature the extractor reads. Java holds the concrete
// Signature; see the DEVIATION note in the file header for why this port names the subset.
type cadesLTASignature interface {
	// CertificateSource gets the certificate source of the signature. Port of
	// getCertificateSource().
	CertificateSource() *spi.SignatureCertificateSource
	// CompleteCertificateSource gets the merged certificate source of the signature and of its
	// timestamps. Port of getCompleteCertificateSource().
	CompleteCertificateSource() *spi.ListCertificateSource
	// CRLSource gets the CRL source of the signature. Port of getCRLSource().
	CRLSource() spi.OfflineRevocationSource[revocation.CRL]
	// OCSPSource gets the OCSP source of the signature. Port of getOCSPSource().
	OCSPSource() spi.OfflineRevocationSource[revocation.OCSP]
	// CompleteCRLSource gets the merged CRL source of the signature and of its timestamps.
	// Port of getCompleteCRLSource().
	CompleteCRLSource() *spi.ListRevocationSource[revocation.CRL]
	// CompleteOCSPSource gets the merged OCSP source of the signature and of its timestamps.
	// Port of getCompleteOCSPSource().
	CompleteOCSPSource() *spi.ListRevocationSource[revocation.OCSP]
	// CMS gets the CMS document the signature belongs to. Port of getCMS().
	CMS() *cms.CMS
}

// LevelBaselineLTATimestampExtractor extracts the necessary information to compute the
// CAdES Archive Timestamp V3.
//
// See "5.5.2 The ats-hash-index-v3 attribute":
//
// The ats-hash-index-v3 is invalid if it contains a reference for which the original value is
// not found, i.e.:
//   - a reference represented by an entry in certificatesHashIndex which corresponds to no
//     instance of CertificateChoices within certificates field of the root SignedData;
//   - a reference represented by an entry in crlsHashIndex which corresponds to no instance of
//     RevocationInfoChoice within crls field of the root SignedData; or
//   - a reference represented by an entry in unsignedAttrValuesHashIndex which corresponds to
//     no octet stream resulting from concatenating one of the AttributeValue instances within
//     field Attribute.attrValues and the corresponding Attribute.attrType within one Attribute
//     instance in unsignedAttrs field of the SignerInfo.
type LevelBaselineLTATimestampExtractor struct {
	// signature is the Signature.
	signature cadesLTASignature
}

// NewCadesLevelBaselineLTATimestampExtractor is the default constructor for the
// LevelBaselineLTATimestampExtractor, taking the Signature related to the archive
// timestamp. Panics with the Java message when it is nil (Objects.requireNonNull upstream).
func NewCadesLevelBaselineLTATimestampExtractor(cadesSignature *Signature) *LevelBaselineLTATimestampExtractor {
	if cadesSignature == nil {
		panic("CAdESSignature cannot be null!")
	}
	return newCadesLevelBaselineLTATimestampExtractor(cadesSignature)
}

// newCadesLevelBaselineLTATimestampExtractor builds the extractor over the narrow view of the
// signature; see the DEVIATION note in the file header.
func newCadesLevelBaselineLTATimestampExtractor(signature cadesLTASignature) *LevelBaselineLTATimestampExtractor {
	return &LevelBaselineLTATimestampExtractor{signature: signature}
}

// AtsHashIndex builds the ats-hash-index unsigned attribute, which provides an unambiguous
// imprint of the essential components of a CAdES signature for use in the archive time-stamp
// (see 6.4.3). These essential components are elements of the following ASN.1 SET OF
// structures: unsignedAttrs, SignedData.certificates, and SignedData.crls.
//
// The ats-hash-index attribute value has the ASN.1 syntax ATSHashIndex:
//
//	ATSHashIndex ::= SEQUENCE {
//	    hashIndAlgorithm      AlgorithmIdentifier DEFAULT {algorithm id-sha256},
//	    certificatesHashIndex SEQUENCE OF OCTET STRING,
//	    crlsHashIndex         SEQUENCE OF OCTET STRING,
//	    unsignedAttrsHashIndex SEQUENCE OF OCTET STRING }
//
// Port of getAtsHashIndex(SignerInformation, DigestAlgorithm, ASN1ObjectIdentifier).
func (e *LevelBaselineLTATimestampExtractor) AtsHashIndex(signerInformation *cmscore.SignerInfo,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier) (*cmscore.Attribute, error) {
	algorithmIdentifier, err := e.hashIndexDigestAlgorithmIdentifier(hashIndexDigestAlgorithm)
	if err != nil {
		return nil, err
	}
	certificatesHashIndex, err := e.certificatesHashIndex(hashIndexDigestAlgorithm)
	if err != nil {
		return nil, err
	}
	crLsHashIndex, err := e.crlsHashIndex(hashIndexDigestAlgorithm)
	if err != nil {
		return nil, err
	}
	unsignedAttributesHashIndex, err := e.unsignedAttributesHashIndex(
		signerInformation, atsHashIndexVersionIdentifier, hashIndexDigestAlgorithm)
	if err != nil {
		return nil, err
	}
	return cadesLTAComposedAtsHashIndex(algorithmIdentifier, certificatesHashIndex, crLsHashIndex,
		unsignedAttributesHashIndex, atsHashIndexVersionIdentifier), nil
}

// VerifiedAtsHashIndex re-builds the ats-hash-index for verification of the provided token, and
// records what it found in the token's ArchiveTimestampHashIndexStatus.
// Port of getVerifiedAtsHashIndex(SignerInformation, TimestampToken).
func (e *LevelBaselineLTATimestampExtractor) VerifiedAtsHashIndex(signerInformation *cmscore.SignerInfo,
	timestampToken *validation.TimestampToken) (*cmscore.Attribute, error) {
	unsignedAttributes := timestampToken.UnsignedAttributes()
	atsHashIndexVersionIdentifier := UtilsAtsHashIndexVersionIdentifier(unsignedAttributes)
	atsHashIndex := UtilsAtsHashIndexByVersion(unsignedAttributes, atsHashIndexVersionIdentifier)
	// Upstream logs "A valid atsHashIndex [oid: {}] has not been found for a timestamp with id
	// {}" when atsHashIndex is null; the status below carries the same finding.
	atsHashIndexStatus := cadesLTAArchiveTimestampHashIndexStatus(atsHashIndexVersionIdentifier)
	derObjectAlgorithmIdentifier := spi.DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(atsHashIndex)
	hashIndexDigestAlgorithm, err := cadesLTAHashIndexDigestAlgorithm(derObjectAlgorithmIdentifier)
	if err != nil {
		return nil, err
	}
	certificatesHashIndex, err := e.verifiedCertificatesHashIndex(atsHashIndex, hashIndexDigestAlgorithm, atsHashIndexStatus)
	if err != nil {
		return nil, err
	}
	crLsHashIndex, err := e.verifiedCRLsHashIndex(atsHashIndex, hashIndexDigestAlgorithm, atsHashIndexStatus)
	if err != nil {
		return nil, err
	}
	verifiedAttributesHashIndex, err := e.verifiedUnsignedAttributesHashIndex(
		signerInformation, atsHashIndex, atsHashIndexVersionIdentifier, hashIndexDigestAlgorithm, atsHashIndexStatus)
	if err != nil {
		return nil, err
	}
	timestampToken.SetAtsHashIndexStatus(atsHashIndexStatus)
	if atsHashIndexVersionIdentifier == nil {
		// An archive time-stamp carrying no ats-hash-index at all leaves the version identifier
		// null, and upstream then builds `new Attribute(null, ...)`, whose encoding raises a
		// NullPointerException - after the status above has been set on the token, which is
		// where the finding survives. The data-dependent throw becomes an error (PORTING.md).
		return nil, model.NewDSSError("The ats-hash-index was not found or not supported.")
	}
	return cadesLTAComposedAtsHashIndex(derObjectAlgorithmIdentifier, certificatesHashIndex, crLsHashIndex,
		verifiedAttributesHashIndex, atsHashIndexVersionIdentifier), nil
}

// cadesLTAArchiveTimestampHashIndexStatus ports the private
// buildArchiveTimestampHashIndexStatus(ASN1ObjectIdentifier).
func cadesLTAArchiveTimestampHashIndexStatus(atsHashIndexVersionIdentifier asn1.ObjectIdentifier) *validation.ArchiveTimestampHashIndexStatus {
	status := validation.NewArchiveTimestampHashIndexStatus()
	var version enumerations.ArchiveTimestampHashIndexVersion
	if atsHashIndexVersionIdentifier != nil {
		version = enumerations.ArchiveTimestampHashIndexVersionForOID(atsHashIndexVersionIdentifier.String())
	}
	if version != "" {
		status.SetVersion(version)
	} else {
		status.AddErrorMessage("The ats-hash-index was not found or not supported.")
	}
	return status
}

// cadesLTAHashIndexDigestAlgorithm reads the hash algorithm used to compute the hash values
// contained in certificatesHashIndex, crlsHashIndex, and unsignedAttrsHashIndex. This algorithm
// shall be the same as the hash algorithm used for computing the archive time-stamp's message
// imprint.
//
//	hashIndAlgorithm AlgorithmIdentifier DEFAULT {algorithm id-sha256}
//
// Port of the private getHashIndexDigestAlgorithm(AlgorithmIdentifier); Java's forOID raises an
// IllegalArgumentException for an unknown OID, which is data-dependent and therefore an error
// here (PORTING.md).
func cadesLTAHashIndexDigestAlgorithm(algorithmIdentifier *spi.AlgorithmIdentifier) (enumerations.DigestAlgorithm, error) {
	if algorithmIdentifier != nil {
		return enumerations.DigestAlgorithmForOID(algorithmIdentifier.Algorithm.String())
	}
	return CAdESUtilsDefaultArchiveTimestampHashAlgo, nil
}

// cadesLTAComposedAtsHashIndex assembles the ATSHashIndex SEQUENCE and wraps it in its
// Attribute. Port of the private getComposedAtsHashIndex.
func cadesLTAComposedAtsHashIndex(algorithmIdentifiers *spi.AlgorithmIdentifier,
	certificatesHashIndex, crLsHashIndex, unsignedAttributesHashIndex []byte,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier) *cmscore.Attribute {
	var vector []byte
	if algorithmIdentifiers != nil {
		vector = append(vector, algorithmIdentifiers.DER()...)
	} else if spi.OIDIdAaATSHashIndexV2.Equal(atsHashIndexVersionIdentifier) ||
		spi.OIDIdAaATSHashIndexV3.Equal(atsHashIndexVersionIdentifier) {
		// for id_aa_ATSHashIndexV2 and id_aa_ATSHashIndexV3, the algorithmIdentifier is required
		//
		// NOTE: this is `new AlgorithmIdentifier(oid)` upstream, i.e. an AlgorithmIdentifier
		// with absent parameters - not DSSASN1Utils.getAlgorithmIdentifier(SHA256), which would
		// agree here but not for every algorithm.
		sha256AlgorithmIdentifier := spi.NewAlgorithmIdentifier(cadesLTASHA256OID)
		vector = append(vector, sha256AlgorithmIdentifier.DER()...)
	}
	if certificatesHashIndex != nil {
		vector = append(vector, cadesLTAToDER(certificatesHashIndex)...)
	}
	if crLsHashIndex != nil {
		vector = append(vector, cadesLTAToDER(crLsHashIndex)...)
	}
	if unsignedAttributesHashIndex != nil {
		vector = append(vector, cadesLTAToDER(unsignedAttributesHashIndex)...)
	}
	derSequence := asn1ber.WriteSequence(vector)
	return cmscore.NewAttribute(atsHashIndexVersionIdentifier, derSequence)
}

// certificatesHashIndex builds the certificatesHashIndex field, a sequence of octet strings.
// Each one contains the hash value of one instance of CertificateChoices within certificates
// field of the root SignedData. A hash value for every instance of CertificateChoices, as
// present at the time when the corresponding archive time-stamp is requested, shall be included
// in certificatesHashIndex. No other hash value shall be included in this field.
//
// Port of the private getCertificatesHashIndex(DigestAlgorithm).
func (e *LevelBaselineLTATimestampExtractor) certificatesHashIndex(
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	var certificatesHashIndexVector []byte
	signedDataCertificates := e.signature.CertificateSource().SignedDataCertificates()
	for _, certificateToken := range signedDataCertificates {
		digest, err := certificateToken.Digest(hashIndexDigestAlgorithm)
		if err != nil {
			return nil, err
		}
		// Upstream logs "Adding to CertificatesHashIndex DSS-Identifier: {} with hash {}".
		certificatesHashIndexVector = append(certificatesHashIndexVector, cadesLTADEROctetString(digest)...)
	}
	return asn1ber.WriteSequence(certificatesHashIndexVector), nil
}

// verifiedCertificatesHashIndex returns the certificatesHashIndex of the time-stamp, having
// checked every entry against SignedData.certificates and, failing that, against every
// certificate the signature knows of (lax processing).
//
// Port of the private getVerifiedCertificatesHashIndex.
func (e *LevelBaselineLTATimestampExtractor) verifiedCertificatesHashIndex(timestampHashIndex []byte,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm,
	atsHashIndexStatus *validation.ArchiveTimestampHashIndexStatus) ([]byte, error) {
	certHashes := UtilsCertificatesHashIndex(timestampHashIndex)
	certHashesList, err := spi.DSSASN1UtilsDEROctetStrings(certHashes)
	if err != nil {
		return nil, err
	}

	// Evaluate SignedData.certificates
	signedDataCertificates := e.signature.CertificateSource().SignedDataCertificates()
	for _, certificateToken := range signedDataCertificates {
		digest, err := certificateToken.Digest(hashIndexDigestAlgorithm)
		if err != nil {
			return nil, err
		}
		// Upstream logs whether the certificate is present in the timestamp.
		certHashesList = cadesLTARemoveDigest(certHashesList, digest)
	}
	// Evaluate against other certificate entries (lax processing)
	if len(certHashesList) != 0 {
		allCertificates := e.signature.CompleteCertificateSource().Certificates()
		for _, certificateToken := range allCertificates {
			digest, err := certificateToken.Digest(hashIndexDigestAlgorithm)
			if err != nil {
				return nil, err
			}
			// Upstream logs "ats-hash-index attribute contains certificate '{}' present
			// outside SignedData.certificates" on a match.
			certHashesList = cadesLTARemoveDigest(certHashesList, digest)
		}

		if len(certHashesList) == 0 {
			atsHashIndexStatus.AddErrorMessage(
				"ats-hash-index attribute contains certificates present outside of SignedData.certificates.")
		} else {
			// Upstream logs "{} attribute(s) hash in Cert Hashes has not been found in
			// document attributes: {}".
			atsHashIndexStatus.AddErrorMessage(
				"Some ats-hash-index attribute certificates have not been found in document attributes.")
		}
	}

	return certHashes, nil
}

// crlsHashIndex builds the crlsHashIndex field, a sequence of octet strings. Each one contains
// the hash value of one instance of RevocationInfoChoice within crls field of the root
// SignedData. A hash value for every instance of RevocationInfoChoice, as present at the time
// when the corresponding archive time-stamp is requested, shall be included in crlsHashIndex.
// No other hash values shall be included in this field.
//
// Port of the private getCRLsHashIndex(DigestAlgorithm).
func (e *LevelBaselineLTATimestampExtractor) crlsHashIndex(
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	var crlsHashIndex []byte

	cmsDocument := e.signature.CMS()
	for _, x509CRLHolder := range cmsDocument.CRLs() {
		digested, err := cadesLTADigestAndAddToList(x509CRLHolder, hashIndexDigestAlgorithm)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to read CRL. Reason : %s", err.Error()), err)
		}
		crlsHashIndex = append(crlsHashIndex, digested...)
	}

	for _, otherRevocationInfoMatch := range cmsDocument.OcspResponseStore() {
		ocspResponseDigest, err := cadesLTAOcspResponseDigest(
			otherRevocationInfoMatch, cmscore.OIDRIOCSPResponse, hashIndexDigestAlgorithm)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to read OCSP Basic. Reason : %s", err.Error()), err)
		}
		crlsHashIndex = append(crlsHashIndex, ocspResponseDigest...)
	}

	for _, otherRevocationInfoMatch := range cmsDocument.OcspBasicStore() {
		ocspResponseDigest, err := cadesLTAOcspResponseDigest(
			otherRevocationInfoMatch, cmscore.OIDPKIXOCSPBasic, hashIndexDigestAlgorithm)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to read OCSP Basic. Reason : %s", err.Error()), err)
		}
		crlsHashIndex = append(crlsHashIndex, ocspResponseDigest...)
	}

	return asn1ber.WriteSequence(crlsHashIndex), nil
}

// cadesLTADigestAndAddToList ports the private digestAndAddToList, returning the DEROctetString
// rather than appending it to a vector.
func cadesLTADigestAndAddToList(encoded []byte,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	digest, err := spi.DSSUtilsDigest(hashIndexDigestAlgorithm, encoded)
	if err != nil {
		return nil, err
	}
	// Upstream logs "Adding to crlsHashIndex with hash {}".
	return cadesLTADEROctetString(digest), nil
}

// verifiedCRLsHashIndex returns the crlsHashIndex of the time-stamp, having checked every entry
// against SignedData.crls and, failing that, against every revocation binary the signature
// knows of (lax processing).
//
// Port of the private getVerifiedCRLsHashIndex.
func (e *LevelBaselineLTATimestampExtractor) verifiedCRLsHashIndex(timestampHashIndex []byte,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm,
	atsHashIndexStatus *validation.ArchiveTimestampHashIndexStatus) ([]byte, error) {
	crlHashes := UtilsCRLHashIndex(timestampHashIndex)
	crlHashesList, err := spi.DSSASN1UtilsDEROctetStrings(crlHashes)
	if err != nil {
		return nil, err
	}

	crlHashesList, err = e.findCRLMatches(crlHashesList, hashIndexDigestAlgorithm,
		cadesLTACMSSignedDataRevocationBinaries[revocation.CRL](e.signature.CRLSource()),
		cadesLTACMSSignedDataRevocationBinaries[revocation.OCSP](e.signature.OCSPSource()))
	if err != nil {
		return nil, err
	}

	// Evaluate against other certificate entries (lax processing)
	if len(crlHashesList) != 0 {

		crlHashesList, err = e.findCRLMatches(crlHashesList, hashIndexDigestAlgorithm,
			e.signature.CompleteCRLSource().AllRevocationBinaries(),
			e.signature.CompleteOCSPSource().AllRevocationBinaries())
		if err != nil {
			return nil, err
		}

		if len(crlHashesList) == 0 {
			atsHashIndexStatus.AddErrorMessage(
				"ats-hash-index attribute contains crls present outside of SignedData.crls.")
		} else {
			// Upstream logs "{} attribute(s) hash in CRL Hashes has not been found in
			// SignedData.crls: {}".
			atsHashIndexStatus.AddErrorMessage(
				"Some ats-hash-index attribute crls have not been found in document attributes.")
		}

	}

	return crlHashes, nil
}

// findCRLMatches removes from crlHashesList every hash matched by one of the provided
// revocation binaries, and returns what is left. Port of the private findCRLMatches, whose
// cmsSignedDataMode flag only selected a log message and is therefore dropped.
func (e *LevelBaselineLTATimestampExtractor) findCRLMatches(crlHashesList [][]byte,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm,
	crlBinaries []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL],
	ocspBinaries []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]) ([][]byte, error) {
	for _, crl := range crlBinaries {
		digest, err := crl.DigestValue(hashIndexDigestAlgorithm)
		if err != nil {
			return nil, err
		}
		crlHashesList = cadesLTARemoveDigest(crlHashesList, digest)
	}

	for _, ocsp := range ocspBinaries {
		binary, ok := ocsp.(*spi.OCSPResponseBinary)
		if !ok {
			// Java casts unconditionally; a non-OCSPResponseBinary would be a
			// ClassCastException there and cannot occur in either port.
			continue
		}
		objectIdentifier := binary.ASN1ObjectIdentifier()
		if cmscore.OIDPKIXOCSPBasic.Equal(objectIdentifier) {
			// OCSPObjectIdentifiers.id_pkix_ocsp_basic
			basicResponseDigest, err := cadesLTAOcspResponseDigest(
				binary.BasicOCSPRespContent(), objectIdentifier, hashIndexDigestAlgorithm)
			if err != nil {
				return nil, err
			}
			crlHashesList = cadesLTARemoveDEROctetString(crlHashesList, basicResponseDigest)

		} else {
			// OCSPObjectIdentifiers.id_ri_ocsp_response case
			if objectIdentifier == nil {
				objectIdentifier = cadesLTAPKIXOCSPResponseOID
			}
			fullBinaryDigest, err := cadesLTAOcspResponseDigest(
				binary.Binaries(), objectIdentifier, hashIndexDigestAlgorithm)
			if err != nil {
				return nil, err
			}
			basicResponseDigest, err := cadesLTAOcspResponseDigest(
				binary.BasicOCSPRespContent(), objectIdentifier, hashIndexDigestAlgorithm)
			if err != nil {
				return nil, err
			}
			// Java's `remove(a) || remove(b)` short-circuits, so the second removal only runs
			// when the first found nothing.
			before := len(crlHashesList)
			crlHashesList = cadesLTARemoveDEROctetString(crlHashesList, fullBinaryDigest)
			if len(crlHashesList) == before {
				crlHashesList = cadesLTARemoveDEROctetString(crlHashesList, basicResponseDigest)
			}
		}
	}
	return crlHashesList, nil
}

// cadesLTACMSSignedDataRevocationBinaries is
// OfflineRevocationSource#getCMSSignedDataRevocationBinaries().
//
// INTEGRATOR NOTE: spi.OfflineRevocationSourceBase implements that method, and every concrete
// source therefore has it, but the spi.OfflineRevocationSource interface AdvancedSignature
// hands out does not declare it - unlike the Java abstract class, which does. The assertion
// below is the workaround while dss-spi stays frozen; declaring the method on the interface
// would remove it and cannot break any implementation.
func cadesLTACMSSignedDataRevocationBinaries[R revocation.Revocation](
	source spi.OfflineRevocationSource[R]) []spi.EncapsulatedRevocationTokenIdentifier[R] {
	type cmsSignedDataRevocationBinaries interface {
		CMSSignedDataRevocationBinaries() []spi.EncapsulatedRevocationTokenIdentifier[R]
	}
	if cmsSource, ok := source.(cmsSignedDataRevocationBinaries); ok {
		return cmsSource.CMSSignedDataRevocationBinaries()
	}
	return nil
}

// cadesLTAOcspResponseDigest ports the private getOcspResponseDigest, returning the
// DEROctetString of the digest.
func cadesLTAOcspResponseDigest(binaries []byte, objectIdentifier asn1.ObjectIdentifier,
	digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	encoded, err := UtilsSignedDataEncodedOCSPResponse(binaries, objectIdentifier)
	if err != nil {
		return nil, err
	}
	digest, err := spi.DSSUtilsDigest(digestAlgorithm, encoded)
	if err != nil {
		return nil, err
	}
	return cadesLTADEROctetString(digest), nil
}

// unsignedAttributesHashIndex builds the unsignedAttrsHashIndex field, a sequence of octet
// strings. Each one contains the hash value of one instance of Attribute within unsignedAttrs
// field of the SignerInfo. A hash value for every instance of Attribute, as present at the time
// when the corresponding archive time-stamp is requested, shall be included in
// unsignedAttrsHashIndex. No other hash values shall be included in this field.
//
// Port of the private getUnsignedAttributesHashIndex; see the DEVIATION note in the file header
// on an absent unsignedAttrs field.
func (e *LevelBaselineLTATimestampExtractor) unsignedAttributesHashIndex(signerInformation *cmscore.SignerInfo,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	var unsignedAttributesHashIndex []byte
	unsignedAttributes := signerInformation.UnsignedAttributes
	for _, attribute := range cadesLTAAttributeTableOrder(unsignedAttributes) {
		attributeDerOctetStringHashes, err := cadesLTAAttributeDerOctetStringHashes(
			attribute, atsHashIndexVersionIdentifier, hashIndexDigestAlgorithm)
		if err != nil {
			return nil, err
		}
		for _, derOctetStringDigest := range attributeDerOctetStringHashes {
			unsignedAttributesHashIndex = append(unsignedAttributesHashIndex, derOctetStringDigest...)
		}
	}
	return asn1ber.WriteSequence(unsignedAttributesHashIndex), nil
}

// verifiedUnsignedAttributesHashIndex returns the unsignedAttrsHashIndex of the time-stamp,
// having checked every entry against the unsignedAttrs of the SignerInfo.
//
// We check that every hash attribute found in the timestamp token is found in the
// signerInformation.
//
// If there are more unsigned attributes in the signerInformation than present in the hash
// attributes list (and there is at least the archiveTimestampAttributeV3), we don't report any
// error nor which attributes are signed by the timestamp. If there are some attributes that are
// not present or altered in the signerInformation, we just return some empty sequence to make
// sure that the timestamped data will not match. We do not report which attributes hash are
// present if any.
//
// If there is no attribute at all in the archive timestamp hash index, that would mean we
// didn't check anything.
//
// Port of the private getVerifiedUnsignedAttributesHashIndex.
func (e *LevelBaselineLTATimestampExtractor) verifiedUnsignedAttributesHashIndex(
	signerInformation *cmscore.SignerInfo, timestampHashIndex []byte,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier, hashIndexDigestAlgorithm enumerations.DigestAlgorithm,
	atsHashIndexStatus *validation.ArchiveTimestampHashIndexStatus) ([]byte, error) {

	unsignedAttributesHashes := UtilsUnsignedAttributesHashIndex(timestampHashIndex)
	timestampUnsignedAttributesHashesList, err := spi.DSSASN1UtilsDEROctetStrings(unsignedAttributesHashes)
	if err != nil {
		return nil, err
	}

	unsignedAttributes := UtilsUnsignedAttributes(signerInformation)
	for _, attribute := range cadesLTAAttributeTableOrder(unsignedAttributes) {
		attributeDerOctetStringHashes, err := cadesLTAAttributeDerOctetStringHashes(
			attribute, atsHashIndexVersionIdentifier, hashIndexDigestAlgorithm)
		if err != nil {
			return nil, err
		}
		for _, derOctetStringDigest := range attributeDerOctetStringHashes {
			// Upstream logs whether the attribute is present in the timestamp.
			timestampUnsignedAttributesHashesList = cadesLTARemoveDEROctetString(
				timestampUnsignedAttributesHashesList, derOctetStringDigest)
		}
	}
	if len(timestampUnsignedAttributesHashesList) != 0 {
		// Upstream logs "{} attribute(s) hash in Timestamp has not been found in unsignedAttrs: {}".
		atsHashIndexStatus.AddErrorMessage(
			"Some ats-hash-index attribute entries have not been found in unsignedAttrs.")
	}
	// return the original DERSequence
	return unsignedAttributesHashes, nil
}

// cadesLTAAttributeDerOctetStringHashes ports the private getAttributeDerOctetStringHashes,
// returning the DEROctetString encoding of every digest.
func cadesLTAAttributeDerOctetStringHashes(attribute *cmscore.Attribute,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier,
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) ([][]byte, error) {
	octets, err := UtilsOctetStringForAtsHashIndex(attribute, atsHashIndexVersionIdentifier)
	if err != nil {
		return nil, err
	}
	if len(octets) != 0 {
		derOctetStrings := make([][]byte, 0, len(octets))
		for _, bytesToDigest := range octets {
			digest, err := spi.DSSUtilsDigest(hashIndexDigestAlgorithm, bytesToDigest)
			if err != nil {
				return nil, err
			}
			derOctetStrings = append(derOctetStrings, cadesLTADEROctetString(digest))
			// Upstream logs "Digest string [{}] has been added to the hash table".
		}
		return derOctetStrings, nil
	}
	return nil, nil
}

// hashIndexDigestAlgorithmIdentifier ports the private
// getHashIndexDigestAlgorithmIdentifier(DigestAlgorithm).
func (e *LevelBaselineLTATimestampExtractor) hashIndexDigestAlgorithmIdentifier(
	hashIndexDigestAlgorithm enumerations.DigestAlgorithm) (*spi.AlgorithmIdentifier, error) {
	// If the algorithm identifier in ATSHashIndex has the default value, then it can be omitted
	if hashIndexDigestAlgorithm.OID() == CAdESUtilsDefaultArchiveTimestampHashAlgo.OID() {
		return nil, nil
	}
	return spi.DSSASN1UtilsAlgorithmIdentifierForDigest(hashIndexDigestAlgorithm)
}

// ArchiveTimestampV3MessageImprint computes a message-imprint for an archive-time-stamp-v3.
// Port of getArchiveTimestampV3MessageImprint(SignerInformation, Attribute, DSSDocument,
// DigestAlgorithm).
func (e *LevelBaselineLTATimestampExtractor) ArchiveTimestampV3MessageImprint(
	signerInformation *cmscore.SignerInfo, atsHashIndexAttribute *cmscore.Attribute,
	originalDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) (model.DSSMessageDigest, error) {
	/*
	 * The input for the archive-time-stamp-v3's message imprint computation shall be the
	 * concatenation (in the order shown by the list below) of the signed data hash (see bullet
	 * 2 below) and certain fields in their binary encoded form without any modification and
	 * including the tag, length and value octets:
	 */
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	digestCalculator.Update(e.encodedContentType())

	documentDigest, err := originalDocument.DigestValue(digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	digestCalculator.Update(documentDigest)

	if err := e.writeSignedFields(signerInformation, digestCalculator); err != nil {
		return model.DSSMessageDigest{}, err
	}

	digestCalculator.Update(cadesLTAAtsHashIndexAttributeEncoding(atsHashIndexAttribute))

	return digestCalculator.MessageDigest(digestAlgorithm), nil
}

// encodedContentType returns 1) The SignedData.encapContentInfo.eContentType, DER encoded.
// Port of the private getEncodedContentType().
func (e *LevelBaselineLTATimestampExtractor) encodedContentType() []byte {
	return asn1ber.EncodeOID(e.signature.CMS().SignedContentType())
}

// writeSignedFields feeds the digest calculator with 3) Fields version, sid, digestAlgorithm,
// signedAttrs, signatureAlgorithm, and signature within the SignedData.signerInfos's item
// corresponding to the signature being archive time-stamped, in their order of appearance.
//
// Port of the private writeSignedFields.
func (e *LevelBaselineLTATimestampExtractor) writeSignedFields(signerInformation *cmscore.SignerInfo,
	digestCalculator *spi.DSSMessageDigestCalculator) error {
	signerInfo := signerInformation

	digestCalculator.Update(cadesLTAEncodedVersion(signerInfo.Version))

	digestCalculator.Update(signerInfo.SID.DER())

	digestCalculator.Update(signerInfo.DigestAlgorithm.DER())

	signedAttributes, err := UtilsDERSignedAttributes(signerInformation)
	if err != nil {
		return err
	}
	digestCalculator.Update(signedAttributes)

	digestCalculator.Update(signerInfo.SignatureAlgorithm.DER())

	digestCalculator.Update(cadesLTADEROctetString(signerInfo.Signature))

	return nil
}

// cadesLTAAtsHashIndexAttributeEncoding ports the private getAtsHashIndexAttributeEncoding.
func cadesLTAAtsHashIndexAttributeEncoding(atsHashIndexAttribute *cmscore.Attribute) []byte {
	attrValue := cadesLTAFirstAttributeValue(atsHashIndexAttribute)
	if attrValue != nil {
		return attrValue
	}
	// Upstream logs "Invalid ats-hash-table-index attribute encoding! The value is skipped."
	return spi.DSSUtilsEmptyByteArray
}

// -----------------------------------------------------------------------------
// Small encoding helpers, all of them one BouncyCastle expression each.
// -----------------------------------------------------------------------------

// cadesLTASHA256OID is DigestAlgorithm.SHA256.getOid() as an ASN.1 object identifier; it is the
// default hash-index algorithm the ats-hash-index-v2/v3 have to spell out.
var cadesLTASHA256OID = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

// cadesLTAPKIXOCSPResponseOID is OCSPObjectIdentifiers.id_pkix_ocsp_response, the fallback
// format identifier findCRLMatches uses for an OCSPResponseBinary that carries none.
var cadesLTAPKIXOCSPResponseOID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 1}

// cadesLTADEROctetString is `new DEROctetString(content).getEncoded(DER)`.
func cadesLTADEROctetString(content []byte) []byte {
	return asn1ber.WriteTLV(asn1ber.TagOctetString, content)
}

// cadesLTAEncodedVersion is `new ASN1Integer(version).getEncoded(DER)` for the small non-negative
// CMSVersion of a SignerInfo.
func cadesLTAEncodedVersion(version int) []byte {
	return asn1ber.EncodeInteger(big.NewInt(int64(version)))
}

// cadesLTAToDER re-encodes an ASN.1 element in DER, which is what adding it to a DERSequence
// does upstream. Unparsable input is passed through unchanged: it can only come from an
// ats-hash-index field this port already parsed, and dropping it would lose octets a verifier
// has to see.
func cadesLTAToDER(encoded []byte) []byte {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil || len(rest) != 0 {
		return encoded
	}
	return element.DEREncoded()
}

// cadesLTAFirstAttributeValue is DSSASN1Utils.getAsn1Encodable(Attribute) followed by
// getDEREncoded: the DER encoding of the first attribute value, nil when there is none.
func cadesLTAFirstAttributeValue(attribute *cmscore.Attribute) []byte {
	if attribute == nil {
		return nil
	}
	if attribute.Element() == nil {
		values := attribute.ValueEncodings()
		if len(values) == 0 {
			return nil
		}
		return cadesLTAToDER(values[0])
	}
	if len(attribute.Values) == 0 {
		return nil
	}
	return attribute.Values[0].DEREncoded()
}

// -----------------------------------------------------------------------------
// AttributeTable#toASN1EncodableVector() ordering.
//
// The unsignedAttrsHashIndex is built by iterating that vector, so the order it hands the
// attributes out in is part of the ats-hash-index bytes a TSA signs. It is NOT the order the
// attributes appear in the SignerInfo: BouncyCastle's AttributeTable groups the attributes into
// a java.util.Hashtable keyed by attrType, and toASN1EncodableVector() returns
// Hashtable#elements() - buckets from the last to the first, each chain most-recently-inserted
// first, the attributes sharing one attrType staying in their received order inside their entry.
// The corpus of testdata/ats-hash-index-oracle.txt pins the reproduction (7 attributes over
// 6 attribute types, plus a repeated attrType).
//
// INTEGRATOR NOTE: this is BouncyCastle machinery with no DSS class of its own, so per
// PORTING.md it belongs in internal/cmscore next to cmscore.Attributes (as, say,
// Attributes.ASN1EncodableVector()) rather than here. It is written here because internal/ is
// frozen for this phase; every other caller of AttributeTable#toASN1EncodableVector() in
// dss-cades needs the same function.
// -----------------------------------------------------------------------------

// cadesLTAAttributeTableOrder returns the attributes in the order
// AttributeTable#toASN1EncodableVector() yields them.
func cadesLTAAttributeTableOrder(attributes cmscore.Attributes) cmscore.Attributes {
	if len(attributes) < 2 {
		return attributes
	}
	table := newCadesLTAAttributeHashtable()
	for _, attribute := range attributes {
		table.add(attribute)
	}
	return table.elements()
}

// cadesLTAAttributeHashtable is the java.util.Hashtable AttributeTable stores its attributes in,
// with its default capacity of 11 and load factor of 0.75.
type cadesLTAAttributeHashtable struct {
	// buckets holds one chain per table slot, each chain ordered head first, i.e. from the
	// most recently inserted entry to the oldest.
	buckets [][]*cadesLTAAttributeEntry
	// count is the number of entries, and threshold the count at which the table is rehashed.
	count, threshold int
}

// cadesLTAAttributeEntry is one Hashtable entry: an attrType and the attributes carrying it,
// which BouncyCastle keeps in a java.util.Vector as soon as there is more than one.
type cadesLTAAttributeEntry struct {
	hash       int32
	attrType   asn1.ObjectIdentifier
	attributes cmscore.Attributes
}

// newCadesLTAAttributeHashtable is `new Hashtable()`.
func newCadesLTAAttributeHashtable() *cadesLTAAttributeHashtable {
	const initialCapacity = 11
	return &cadesLTAAttributeHashtable{
		buckets:   make([][]*cadesLTAAttributeEntry, initialCapacity),
		threshold: initialCapacity * 3 / 4,
	}
}

// add is AttributeTable#addAttribute: it appends to the entry of the attrType when there is
// one - which leaves that entry where it is, Hashtable#put replacing a value in place - and
// inserts a new entry at the head of its bucket otherwise.
func (h *cadesLTAAttributeHashtable) add(attribute *cmscore.Attribute) {
	hash := cadesLTAOIDHashCode(attribute.Type)
	index := h.index(hash)
	for _, entry := range h.buckets[index] {
		if entry.attrType.Equal(attribute.Type) {
			entry.attributes = append(entry.attributes, attribute)
			return
		}
	}
	if h.count >= h.threshold {
		h.rehash()
		index = h.index(hash)
	}
	entry := &cadesLTAAttributeEntry{hash: hash, attrType: attribute.Type, attributes: cmscore.Attributes{attribute}}
	h.buckets[index] = append([]*cadesLTAAttributeEntry{entry}, h.buckets[index]...)
	h.count++
}

// index is `(hash & 0x7FFFFFFF) % tab.length`.
func (h *cadesLTAAttributeHashtable) index(hash int32) int {
	return int(uint32(hash)&0x7FFFFFFF) % len(h.buckets)
}

// rehash grows the table to `(oldCapacity << 1) + 1` slots, re-inserting the entries from the
// last old bucket to the first and each chain head first, every one of them at the head of its
// new bucket.
func (h *cadesLTAAttributeHashtable) rehash() {
	oldBuckets := h.buckets
	newCapacity := len(oldBuckets)*2 + 1
	h.buckets = make([][]*cadesLTAAttributeEntry, newCapacity)
	h.threshold = newCapacity * 3 / 4
	for index := len(oldBuckets) - 1; index >= 0; index-- {
		for _, entry := range oldBuckets[index] {
			newIndex := h.index(entry.hash)
			h.buckets[newIndex] = append([]*cadesLTAAttributeEntry{entry}, h.buckets[newIndex]...)
		}
	}
}

// elements is Hashtable#elements() flattened over the attributes of each entry, i.e. exactly
// what AttributeTable#toASN1EncodableVector() builds.
func (h *cadesLTAAttributeHashtable) elements() cmscore.Attributes {
	var vector cmscore.Attributes
	for index := len(h.buckets) - 1; index >= 0; index-- {
		for _, entry := range h.buckets[index] {
			vector = append(vector, entry.attributes...)
		}
	}
	return vector
}

// cadesLTAOIDHashCode is ASN1ObjectIdentifier#hashCode(), i.e.
// org.bouncycastle.util.Arrays.hashCode over the identifier's content octets:
//
//	int i = data.length; int hc = i + 1;
//	while (--i >= 0) { hc *= 257; hc ^= data[i]; }
//
// with Java's signed byte and int arithmetic.
func cadesLTAOIDHashCode(oid asn1.ObjectIdentifier) int32 {
	contents := cadesLTAOIDContentOctets(oid)
	index := len(contents)
	hashCode := uint32(index + 1)
	for index > 0 {
		index--
		hashCode *= 257
		// data[index] is a signed Java byte, sign-extended to int before the xor.
		hashCode ^= uint32(int32(int8(contents[index])))
	}
	return int32(hashCode)
}

// cadesLTAOIDContentOctets returns the content octets of the OBJECT IDENTIFIER, i.e. its DER
// encoding minus the identifier and length octets.
func cadesLTAOIDContentOctets(oid asn1.ObjectIdentifier) []byte {
	element, _, err := asn1ber.Parse(asn1ber.EncodeOID(oid))
	if err != nil {
		return nil
	}
	return element.Content()
}

// cadesLTARemoveDigest removes the first occurrence of the digest from the list, i.e.
// List<DEROctetString>#remove(new DEROctetString(digest)) where the list members are the
// content octets DSSASN1Utils.getDEROctetStrings hands out.
func cadesLTARemoveDigest(list [][]byte, digest []byte) [][]byte {
	for index, member := range list {
		if bytes.Equal(member, digest) {
			return append(list[:index:index], list[index+1:]...)
		}
	}
	return list
}

// cadesLTARemoveDEROctetString is cadesLTARemoveDigest for a digest still wrapped in its
// DEROctetString encoding, which is what the digest helpers of this file return.
func cadesLTARemoveDEROctetString(list [][]byte, derOctetString []byte) [][]byte {
	element, rest, err := asn1ber.Parse(derOctetString)
	if err != nil || len(rest) != 0 {
		return list
	}
	return cadesLTARemoveDigest(list, element.Octets())
}
