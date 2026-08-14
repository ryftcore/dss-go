// Ported from dss-crl-parser-x509crl/src/main/java/eu/europa/esig/dss/crl/x509/impl/CRLUtilsX509CRLImpl.java (DSS 6.5.RC1).
//
// DEVIATION: upstream loads the CRL through java.security.cert.CertificateFactory, preferring
// a BouncyCastle provider and falling back to the JDK's own. This port parses with
// crypto/x509.ParseRevocationList directly: there is no provider concept in Go, and
// ParseRevocationList already performs the equivalent structural decoding (RFC 5280 CRL
// syntax) without depending on any registered algorithm provider.
package crlparser

import (
	"crypto/x509"
	"crypto/x509/pkix"
	encoding_asn1 "encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// OIDs of the CRL (crl_ext) and CRL entry (crl entry extensions) extensions this file reads
// directly off RevocationList(Entry).Extensions - either because crypto/x509 leaves them raw
// (issuingDistributionPoint, expiredCertsOnCRL, invalidityDate, certificateIssuer) or because
// their presence, not just their value, carries meaning (cRLReason; see CRLEntry.RevocationReason).
const (
	crlUtilsOIDIssuingDistributionPoint = "2.5.29.28"
	crlUtilsOIDExpiredCertsOnCRL        = "2.5.29.60"
	crlUtilsOIDCRLReason                = "2.5.29.21"
	crlUtilsOIDInvalidityDate           = "2.5.29.24"
	crlUtilsOIDCertificateIssuer        = "2.5.29.29"
)

// crlUtilsX509CRLImplBuildCRLValidity is the native implementation of
// CRLUtilsX509CRLImpl#buildCRLValidity(CRLBinary, CertificateToken).
//
// DEVIATION (parser strictness): crypto/x509.ParseRevocationList decodes the cRLNumber and
// deltaCRLIndicator extensions eagerly and rejects the whole CRL when either is malformed,
// where BouncyCastle's CertificateFactory leaves every extension raw and upstream's
// extractCrlNumber only logs a warning for one it cannot read. A CRL carrying a malformed
// cRLNumber therefore yields a usable CRLValidity (with a null CRL number) in Java but an
// error here. Working around it would mean re-encoding the CRL, which would break the
// original-bytes guarantee CRLBinary/getDerEncoded rests on.
// crlUtilsParseRevocationList decodes a CertificateList, accepting X.509 v1 CRLs as well as v2
// ones.
//
// DEVIATION REPAIR (v1 CRLs): crypto/x509.ParseRevocationList insists on the OPTIONAL
// TBSCertList.version field being present and equal to 1 (v2) and rejects everything else with
// "x509: unsupported crl version". RFC 5280 section 5.1 makes that field OPTIONAL with a DEFAULT
// of v1, so a CRL with no extensions - which is what every conforming CRL issuer emits when it
// has nothing to put in crlExtensions - legally omits it. Upstream reaches
// java.security.cert.CertificateFactory#generateCRL, which parses v1 and v2 alike, so refusing
// v1 here is a Go-side regression against upstream, not a deliberate port decision: it silently
// drops real revocation data (found by the XAdES cross-validation harness on
// Signature-X-BE_ECON-3.xml, whose chain is covered by two v1 CRLs - Belgium Root CA2's and
// Belgium Root CA3's - alongside two v2 ones, so upstream reaches LT where this port stopped at
// T; the effect is not XAdES-specific but reaches every CRL-consuming path in the port).
//
// The repair keeps stdlib as the only structural decoder: when, and only when, the input is a
// well-formed CertificateList whose TBSCertList omits the version field, this splices a
// "version v2" INTEGER into a THROWAWAY copy of the DER purely so that stdlib's parser will
// accept it, then restores RevocationList.Raw and .RawTBSRevocationList to the caller's
// original, untouched bytes. Nothing re-encoded ever leaves this function, and in particular the
// signature is still verified over the original RawTBSRevocationList, so CRLBinary's
// original-bytes guarantee holds. Every other field (issuer, thisUpdate/nextUpdate, the revoked
// entries, extensions) decodes byte-identically either way, because inserting a field at the
// front of the SEQUENCE changes only offsets, never content.
func crlUtilsParseRevocationList(der []byte) (*x509.RevocationList, error) {
	revocationList, err := x509.ParseRevocationList(der)
	if err == nil {
		return revocationList, nil
	}

	patched, originalTBS, ok := crlUtilsSpliceV2VersionForParsing(der)
	if !ok {
		// Not a versionless CRL: report stdlib's own diagnosis unchanged.
		return nil, err
	}
	revocationList, patchedErr := x509.ParseRevocationList(patched)
	if patchedErr != nil {
		// Splicing did not help; the CRL is broken for some other reason, so the caller should
		// still see the failure stdlib reported for the bytes it was actually given.
		return nil, err
	}
	// Hand back the caller's own bytes: everything downstream (signature verification over the
	// TBS, CRLValidity.DerEncoded, the signature-algorithm re-read off Raw) must see the DER as
	// it arrived, never the throwaway copy.
	revocationList.Raw = der
	revocationList.RawTBSRevocationList = originalTBS
	return revocationList, nil
}

// crlUtilsSpliceV2VersionForParsing returns a copy of a CertificateList with an explicit
// "version v2" INTEGER inserted as the first TBSCertList field, together with the original
// TBSCertList element (tag, length and content) it was derived from. ok is false - and the other
// results meaningless - unless der really is a SEQUENCE whose first element is a SEQUENCE
// (TBSCertList) that itself starts with something other than an INTEGER, i.e. exactly the
// versionless v1 shape; in every other case the caller keeps stdlib's original error.
func crlUtilsSpliceV2VersionForParsing(der []byte) (patched, originalTBS []byte, ok bool) {
	outer := cryptobyte.String(der)
	var certificateList cryptobyte.String
	if !outer.ReadASN1(&certificateList, cryptobyte_asn1.SEQUENCE) || !outer.Empty() {
		return nil, nil, false
	}
	tbsReader := certificateList
	var tbsCertList cryptobyte.String
	if !tbsReader.ReadASN1Element(&tbsCertList, cryptobyte_asn1.SEQUENCE) {
		return nil, nil, false
	}
	originalTBS = []byte(tbsCertList)
	// Peek at the first TBSCertList field: an INTEGER means the version is already there and the
	// CRL was rejected for some other reason, which is not this repair's business.
	var tbsContent cryptobyte.String
	inner := cryptobyte.String(originalTBS)
	if !inner.ReadASN1(&tbsContent, cryptobyte_asn1.SEQUENCE) {
		return nil, nil, false
	}
	if tbsContent.PeekASN1Tag(cryptobyte_asn1.INTEGER) {
		return nil, nil, false
	}

	// v1 CRLs carry no crlExtensions, so a v1 CRL that already has them is malformed; leave it.
	var patchedTBS cryptobyte.Builder
	patchedTBS.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
		child.AddASN1Int64(crlUtilsX509V2Version)
		child.AddBytes(tbsContent)
	})
	patchedTBSBytes, err := patchedTBS.Bytes()
	if err != nil {
		return nil, nil, false
	}

	// The trailing signatureAlgorithm + signatureValue are copied over verbatim.
	trailing := []byte(tbsReader)
	var patchedList cryptobyte.Builder
	patchedList.AddASN1(cryptobyte_asn1.SEQUENCE, func(child *cryptobyte.Builder) {
		child.AddBytes(patchedTBSBytes)
		child.AddBytes(trailing)
	})
	patchedBytes, err := patchedList.Bytes()
	if err != nil {
		return nil, nil, false
	}
	return patchedBytes, originalTBS, true
}

// crlUtilsX509V2Version is the value RFC 5280 gives TBSCertList.version for a v2 CRL, and the
// only one crypto/x509.ParseRevocationList accepts.
const crlUtilsX509V2Version = 1

func crlUtilsX509CRLImplBuildCRLValidity(crlBinary *CRLBinary, issuerToken *model.CertificateToken) (*CRLValidity, error) {
	crlValidity := NewCRLValidity(crlBinary)

	revocationList, err := crlUtilsParseRevocationList(crlBinary.Binaries())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to parse the CRL : %s", err.Error()), err)
	}
	crlValidity.SetX509CRL(revocationList)

	sigAlgOID, sigAlgParams, err := crlUtilsSigAlgOIDAndParams(revocationList.Raw)
	if err != nil {
		return nil, err
	}
	signatureAlgorithm, err := enumerations.SignatureAlgorithmForOIDAndParams(sigAlgOID, sigAlgParams)
	if err != nil {
		return nil, err
	}
	crlValidity.SetSignatureAlgorithm(signatureAlgorithm)

	thisUpdate := revocationList.ThisUpdate
	crlValidity.SetThisUpdate(&thisUpdate)
	if !revocationList.NextUpdate.IsZero() {
		nextUpdate := revocationList.NextUpdate
		crlValidity.SetNextUpdate(&nextUpdate)
	}

	crlIssuer, err := model.NewX500Principal(revocationList.RawIssuer)
	if err != nil {
		return nil, err
	}
	if crlIssuer.Equals(issuerToken.Subject().Principal()) {
		crlValidity.SetIssuerX509PrincipalMatches(true)
	}

	criticalExtensionsOid := make([]string, 0, len(revocationList.Extensions))
	var issuingDistributionPointValue, expiredCertsOnCRLValue []byte
	for _, ext := range revocationList.Extensions {
		if ext.Critical {
			criticalExtensionsOid = append(criticalExtensionsOid, ext.Id.String())
		}
		switch ext.Id.String() {
		case crlUtilsOIDIssuingDistributionPoint:
			issuingDistributionPointValue = ext.Value
		case crlUtilsOIDExpiredCertsOnCRL:
			expiredCertsOnCRLValue = ext.Value
		}
	}
	crlValidity.SetCriticalExtensionsOid(criticalExtensionsOid)

	if err := crlUtilsExtractIssuingDistributionPointBinary(crlValidity, issuingDistributionPointValue); err != nil {
		return nil, err
	}
	crlUtilsExtractExpiredCertsOnCRL(crlValidity, expiredCertsOnCRLValue)
	crlUtilsExtractCrlNumber(crlValidity, revocationList.Number)

	crlUtilsCheckSignatureValue(revocationList, issuerToken, crlValidity)
	if crlValidity.IsSignatureIntact() {
		crlSign := issuerToken.CheckKeyUsage(enumerations.KeyUsageBit_CRL_SIGN)
		if !crlSign {
			crlValidity.SetSignatureInvalidityReason(
				fmt.Sprintf("CRL issuer does not have '%s' key usage!", enumerations.KeyUsageBit_CRL_SIGN.Value()))
		}
		crlValidity.SetCrlSignKeyUsage(crlSign)
	}

	return crlValidity, nil
}

// crlUtilsCheckSignatureValue verifies the CRL's signature with the issuer's public key,
// recording the outcome. Port of the private checkSignatureValue(X509CRL, CertificateToken,
// CRLValidity).
//
// DEVIATION: upstream reaches java.security.cert.X509CRL#verify(PublicKey) through the JCA,
// which raises a GeneralSecurityException for both "wrong key" and "no provider for this
// algorithm". This port verifies directly with crypto/x509.Certificate#CheckSignature (the
// same primitive model.CertificateToken's own signature checks use - see certificate_token.go
// in dss-model), wrapping a throwaway Certificate around the issuer's public key; the missing
// -key guard mirrors that file's InvalidKeyException branch. As documented for
// CertificateToken, crypto/x509 refuses MD5-signed input outright, where the JCA would verify
// it.
func crlUtilsCheckSignatureValue(revocationList *x509.RevocationList, issuerToken *model.CertificateToken, crlValidity *CRLValidity) {
	key := issuerToken.PublicKey().Key()
	if key == nil {
		crlValidity.SetSignatureInvalidityReason("CRL Signature cannot be validated : InvalidKeyException: the public key has not been parsed")
		return
	}
	signer := &x509.Certificate{PublicKey: key}
	if err := signer.CheckSignature(revocationList.SignatureAlgorithm, revocationList.RawTBSRevocationList, revocationList.Signature); err != nil {
		crlValidity.SetSignatureInvalidityReason(fmt.Sprintf("CRL Signature cannot be validated : %s", err.Error()))
		return
	}
	crlValidity.SetSignatureIntact(true)
	crlValidity.SetIssuerToken(issuerToken)
}

// crlUtilsSigAlgOIDAndParams extracts the OID and parameters of the outer CertificateList's
// signatureAlgorithm AlgorithmIdentifier, standing in for X509CRL#getSigAlgOID()/
// #getSigAlgParams(). Parallels certificateTokenSigAlgOIDAndParams in dss-model's
// certificate_token.go, which reads the same shape (AlgorithmIdentifier as the second field of
// an outer SEQUENCE) off a Certificate instead of a CertificateList; the two are not shared
// because they live in different packages.
func crlUtilsSigAlgOIDAndParams(raw []byte) (string, []byte, error) {
	var certificateList struct {
		TBSCertList        encoding_asn1.RawValue
		SignatureAlgorithm encoding_asn1.RawValue
		SignatureValue     encoding_asn1.BitString
	}
	if _, err := encoding_asn1.Unmarshal(raw, &certificateList); err != nil {
		return "", nil, model.NewDSSErrorMessageCause("Unable to parse the CRL signature algorithm", err)
	}
	var algorithmIdentifier struct {
		Algorithm  encoding_asn1.ObjectIdentifier
		Parameters encoding_asn1.RawValue `asn1:"optional"`
	}
	if _, err := encoding_asn1.Unmarshal(certificateList.SignatureAlgorithm.FullBytes, &algorithmIdentifier); err != nil {
		return "", nil, model.NewDSSErrorMessageCause("Unable to parse the CRL signature algorithm", err)
	}
	params := algorithmIdentifier.Parameters.FullBytes
	if len(params) == 0 ||
		(algorithmIdentifier.Parameters.Class == encoding_asn1.ClassUniversal && algorithmIdentifier.Parameters.Tag == encoding_asn1.TagNull) {
		params = nil
	}
	return algorithmIdentifier.Algorithm.String(), params, nil
}

// crlUtilsX509CRLImplRevocationInfo is the native implementation of
// CRLUtilsX509CRLImpl#getRevocationInfo(CRLValidity, BigInteger).
//
// Upstream re-parses the CRL (loadCRL) only when the given CRLValidity is not an
// X509CRLValidity or otherwise carries no cached X509CRL, wrapping a parse failure in an
// unchecked DSSException; ICRLUtils#getRevocationInfo declares no checked exception, so that
// failure propagates to the caller same as any other RuntimeException would. Every CRLValidity
// this port produces (CRLUtilsBuildCRLValidity, above) already carries its parsed
// RevocationList, so this path is not expected to run in practice; it panics rather than
// silently returning nil, matching the unchecked-propagation semantics.
func crlUtilsX509CRLImplRevocationInfo(crlValidity *CRLValidity, serialNumber *big.Int) *CRLEntry {
	revocationList := crlValidity.X509CRL()
	if revocationList == nil {
		var err error
		revocationList, err = crlUtilsParseRevocationList(crlValidity.DerEncoded())
		if err != nil {
			panic(model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to get revocation info. Reason : %s", err.Error()), err))
		}
	}
	indirect := crlUtilsIsIndirectCRL(revocationList)
	// previousCertificateIssuer carries the certificateIssuer of the last preceding entry that
	// declared one, as RFC 5280 section 5.3.3 prescribes for indirect CRLs; it stays nil for a
	// direct CRL, where getCertificateIssuer() answers null for every entry.
	var previousCertificateIssuer *model.X500Principal
	for i := range revocationList.RevokedCertificateEntries {
		entry := &revocationList.RevokedCertificateEntries[i]
		if entry.SerialNumber != nil && entry.SerialNumber.Cmp(serialNumber) == 0 {
			return newCRLEntry(entry, indirect, previousCertificateIssuer)
		}
		if indirect {
			if issuer := crlEntryDeclaredCertificateIssuer(entry); issuer != nil {
				previousCertificateIssuer = issuer
			}
		}
	}
	return nil
}

// crlUtilsIsIndirectCRL reports whether the CRL declares indirectCRL TRUE in its
// issuingDistributionPoint extension, which is what gates certificateIssuer resolution in
// BouncyCastle's X509CRLObject. A malformed extension answers false, as
// IssuingDistributionPoint.getInstance's failure leaves BouncyCastle's flag unset.
func crlUtilsIsIndirectCRL(revocationList *x509.RevocationList) bool {
	for _, ext := range revocationList.Extensions {
		if ext.Id.String() != crlUtilsOIDIssuingDistributionPoint {
			continue
		}
		idp, err := crlUtilsParseIssuingDistributionPoint(ext.Value)
		if err != nil {
			return false
		}
		return idp.indirectCrl
	}
	return false
}

// crlEntryDeclaredCertificateIssuer returns the certificateIssuer the entry declares itself,
// ignoring any inherited one; nil when the entry carries no readable certificateIssuer.
//
// DEVIATION: when it refreshes the inherited certificateIssuer, BouncyCastle's
// X509CRLObject#getRevokedCertificate takes getNames()[0] of the GeneralNames unconditionally
// and lets X500Name.getInstance throw an IllegalArgumentException if that first alternative is
// not a directoryName - even though the entry's own getCertificateIssuer() searches the list
// for the first directoryName instead. This port searches in both places, so an indirect CRL
// whose certificateIssuer lists a non-directoryName first resolves here instead of aborting
// the whole lookup.
func crlEntryDeclaredCertificateIssuer(entry *x509.RevocationListEntry) *model.X500Principal {
	for _, ext := range entry.Extensions {
		if ext.Id.String() == crlUtilsOIDCertificateIssuer {
			return crlEntryParseCertificateIssuer(ext.Value)
		}
	}
	return nil
}

// CRLEntry is a revoked certificate entry of a CRL, the Go stand-in for
// java.security.cert.X509CRLEntry as CRLUtilsRevocationInfo returns it.
//
// DEVIATION: crypto/x509.RevocationListEntry already exposes SerialNumber, RevocationTime and
// ReasonCode, but collapses "reasonCode extension absent" and "reasonCode extension present
// with value 0 (Unspecified)" to the same zero int (its own documented ambiguity), where
// X509CRLEntry#getRevocationReason() returns null only for the former. CRLEntry resolves that
// by checking the entry's Extensions for the reasonCode OID itself. invalidityDate and
// certificateIssuer (indirect CRLs) have no dedicated RevocationListEntry field at all and are
// decoded here from the raw extension bytes.
type CRLEntry struct {
	serialNumber      *big.Int
	revocationDate    time.Time
	revocationReason  *int
	certificateIssuer *model.X500Principal
	invalidityDate    *time.Time
	extensions        []pkix.Extension
}

// newCRLEntry builds a CRLEntry from a parsed crypto/x509.RevocationListEntry.
//
// indirect tells whether the enclosing CRL is an indirect CRL and previousCertificateIssuer is
// the certificateIssuer in force for this entry, i.e. the one declared by the closest preceding
// entry: together they reproduce BouncyCastle's X509CRLEntryObject#loadCertificateIssuer, where
// a direct CRL answers null for every entry and an entry of an indirect CRL that declares no
// certificateIssuer of its own inherits the preceding one (RFC 5280 section 5.3.3).
func newCRLEntry(entry *x509.RevocationListEntry, indirect bool, previousCertificateIssuer *model.X500Principal) *CRLEntry {
	crlEntry := &CRLEntry{
		serialNumber:   entry.SerialNumber,
		revocationDate: entry.RevocationTime,
		extensions:     entry.Extensions,
	}
	if indirect {
		crlEntry.certificateIssuer = previousCertificateIssuer
	}
	for _, ext := range entry.Extensions {
		switch ext.Id.String() {
		case crlUtilsOIDCRLReason:
			reason := entry.ReasonCode
			crlEntry.revocationReason = &reason
		case crlUtilsOIDInvalidityDate:
			if t, ok := crlUtilsParseGeneralizedTime(ext.Value); ok {
				crlEntry.invalidityDate = &t
			}
		case crlUtilsOIDCertificateIssuer:
			if indirect {
				// A certificateIssuer the entry cannot be read out of answers null, it does
				// NOT fall back to the inherited one.
				crlEntry.certificateIssuer = crlEntryParseCertificateIssuer(ext.Value)
			}
		}
	}
	return crlEntry
}

// SerialNumber returns the serial number of the revoked certificate.
// Port of X509CRLEntry#getSerialNumber().
func (e *CRLEntry) SerialNumber() *big.Int {
	return e.serialNumber
}

// RevocationDate returns the date on which the certificate was revoked.
// Port of X509CRLEntry#getRevocationDate().
func (e *CRLEntry) RevocationDate() time.Time {
	return e.revocationDate
}

// RevocationReason returns the RFC 5280 reasonCode ordinal, or nil when the reasonCode
// extension is absent. Port of X509CRLEntry#getRevocationReason().
func (e *CRLEntry) RevocationReason() *int {
	return e.revocationReason
}

// CertificateIssuer returns the certificateIssuer extension value (indirect CRLs), or nil when
// absent. Port of X509CRLEntry#getCertificateIssuer().
func (e *CRLEntry) CertificateIssuer() *model.X500Principal {
	return e.certificateIssuer
}

// InvalidityDate returns the invalidityDate extension value, or nil when absent.
// Port of reading the invalidityDate extension off X509CRLEntry#getExtensionValue(...); the
// JDK exposes no dedicated accessor for it either.
func (e *CRLEntry) InvalidityDate() *time.Time {
	return e.invalidityDate
}

// HasExtensions reports whether the entry carries any extensions.
// Port of X509CRLEntry#hasExtensions().
func (e *CRLEntry) HasExtensions() bool {
	return len(e.extensions) > 0
}

// CriticalExtensionOIDs returns the OIDs of the entry's critical extensions.
// Port of X509Extension#getCriticalExtensionOIDs().
func (e *CRLEntry) CriticalExtensionOIDs() []string {
	oids := make([]string, 0, len(e.extensions))
	for _, ext := range e.extensions {
		if ext.Critical {
			oids = append(oids, ext.Id.String())
		}
	}
	return oids
}

// crlEntryParseCertificateIssuer decodes the certificateIssuer entry extension:
//
//	certificateIssuer ::= GeneralNames
//
// looking for the first directoryName [4] Name GeneralName. Like DistributionPointName's
// fullName, GeneralName's directoryName alternative wraps a CHOICE (Name) and is therefore
// EXPLICIT, so the content of the [4] wrapper is the Name's own RDNSequence SEQUENCE encoding,
// ready to hand to model.NewX500Principal as-is.
//
// Unlike the IDP's distributionPoint field - where GeneralNames sits under an IMPLICIT [0]
// tag contributed by the enclosing fullName alternative - certificateIssuer's ASN.1 type is
// GeneralNames itself, which keeps its own natural SEQUENCE tag; the extension value must be
// unwrapped one more level than the IDP's fullName content before the GeneralName elements
// can be walked.
func crlEntryParseCertificateIssuer(value []byte) *model.X500Principal {
	input := cryptobyte.String(value)
	var generalNames cryptobyte.String
	if !input.ReadASN1(&generalNames, cryptobyte_asn1.SEQUENCE) {
		return nil
	}
	for !generalNames.Empty() {
		var content cryptobyte.String
		var tag cryptobyte_asn1.Tag
		if !generalNames.ReadAnyASN1(&content, &tag) {
			return nil
		}
		if tag == cryptobyte_asn1.Tag(4).Constructed().ContextSpecific() {
			principal, err := model.NewX500Principal(content)
			if err != nil {
				return nil
			}
			return principal
		}
	}
	return nil
}
