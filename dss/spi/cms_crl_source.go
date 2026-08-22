// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CMSCRLSource.java (DSS 6.5.RC1).
//
// BouncyCastle replacements used here (see PORTING.md):
//
//   - org.bouncycastle.util.Store<X509CRLHolder> for SignedData.crls -> [][]byte, the
//     encoding of every CertificateList of the field, i.e. what cmscore.CMS.CRLs() hands out;
//     Store#getMatches(null) is then iterating the slice, and X509CRLHolder#getEncoded() is
//     the member itself.
//   - org.bouncycastle.asn1.cms.AttributeTable -> cmscore.Attributes.
//   - org.bouncycastle.asn1.esf.{RevocationValues, CrlOcspRef, CrlListID, OcspListID} -> the
//     types defined below. CMSOCSPSource reads the same two structures from the same two
//     attributes (RevocationValues.ocspVals and CrlOcspRef.ocspids), so they live here, next
//     to the first of their two consumers, the way crl_ref.go already hosts CrlValidatedID.
//
// This file also carries DSSASN1Utils.getRevocationValues(ASN1Encodable), which takes an ESF
// type only this file introduces; it keeps the flattened static-utility naming so that folding
// it back into dss_asn1_utils.go later is a pure move.
package spi

import (
	"encoding/asn1"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

// The unsigned attribute types the two revocation sources read. They are
// org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers constants upstream, hanging off
// 1.2.840.113549.1.9.16.2, and keep their exact Java field name behind the "OID_" prefix.
var (
	// OIDIdAaEtsRevocationRefs is
	// id-aa-ets-revocationRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 22 }
	OIDIdAaEtsRevocationRefs = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 22}

	// OIDIdAaEtsRevocationValues is
	// id-aa-ets-revocationValues OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
	// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 24 }
	OIDIdAaEtsRevocationValues = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 24}
)

// -----------------------------------------------------------------------------
// The ESF structures the revocation-values and revocation-references attributes carry,
// replacing org.bouncycastle.asn1.esf.
//
// The ETSI TS 101 733 ASN.1 module is EXPLICIT TAGS, so every context-specific field below
// wraps its value rather than replacing its tag - which is also how BouncyCastle reads them
// (ASN1TaggedObject#getExplicitBaseObject).
// -----------------------------------------------------------------------------

// RevocationValues is
//
//	RevocationValues ::= SEQUENCE {
//	    crlVals      [0] SEQUENCE OF CertificateList OPTIONAL,
//	    ocspVals     [1] SEQUENCE OF BasicOCSPResponse OPTIONAL,
//	    otherRevVals [2] OtherRevVals OPTIONAL }
//
// replacing org.bouncycastle.asn1.esf.RevocationValues. The members are kept as their own
// encodings, since a CRL binary and an OCSP response binary are both identified by their
// exact bytes.
type RevocationValues struct {
	// CrlVals holds the encoding of every CertificateList of the crlVals field; an absent
	// field yields an empty slice, as BouncyCastle's getCrlVals() yields an empty array.
	CrlVals [][]byte
	// OcspVals holds the encoding of every BasicOCSPResponse of the ocspVals field.
	OcspVals [][]byte
	// OtherRevVals is the encoding of the otherRevVals field, nil when absent.
	OtherRevVals []byte
}

// ParseRevocationValues decodes a RevocationValues from its encoding.
// Port of RevocationValues.getInstance(Object).
func ParseRevocationValues(encoded []byte) (*RevocationValues, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the RevocationValues")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, errors.New("malformed RevocationValues: not a SEQUENCE")
	}
	if len(element.Children()) > 3 {
		return nil, fmt.Errorf("bad sequence size in RevocationValues: %d", len(element.Children()))
	}
	values := &RevocationValues{CrlVals: [][]byte{}, OcspVals: [][]byte{}}
	for _, child := range element.Children() {
		if child.Class() != asn1ber.ClassContextSpecific {
			return nil, errors.New("malformed RevocationValues: a member is not tagged")
		}
		base, err := cmsCRLSourceExplicitBase(child, "RevocationValues")
		if err != nil {
			return nil, err
		}
		switch child.TagNumber() {
		case 0:
			values.CrlVals, err = cmsCRLSourceSequenceOf(base, "RevocationValues.crlVals")
		case 1:
			values.OcspVals, err = cmsCRLSourceSequenceOf(base, "RevocationValues.ocspVals")
		case 2:
			values.OtherRevVals = base.Encoded()
		default:
			return nil, fmt.Errorf("invalid tag in RevocationValues: %d", child.TagNumber())
		}
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

// CrlListID is
//
//	CRLListID ::= SEQUENCE { crls SEQUENCE OF CrlValidatedID }
//
// replacing org.bouncycastle.asn1.esf.CrlListID.
type CrlListID struct {
	// crls holds the encoding of every member, decoded on demand by Crls().
	crls [][]byte
}

// Crls decodes and returns the members of the crls field.
// Port of CrlListID#getCrls(), which likewise decodes on each call and lets one malformed
// member abort the whole conversion.
func (c *CrlListID) Crls() ([]*CrlValidatedID, error) {
	crls := make([]*CrlValidatedID, 0, len(c.crls))
	for _, encoded := range c.crls {
		crl, err := ParseCrlValidatedID(encoded)
		if err != nil {
			return nil, err
		}
		crls = append(crls, crl)
	}
	return crls, nil
}

// OcspListID is
//
//	OcspListID ::= SEQUENCE { ocspResponses SEQUENCE OF OcspResponsesID }
//
// replacing org.bouncycastle.asn1.esf.OcspListID.
type OcspListID struct {
	// ocspResponses holds the encoding of every member, decoded on demand by OcspResponses().
	ocspResponses [][]byte
}

// OcspResponses decodes and returns the members of the ocspResponses field.
// Port of OcspListID#getOcspResponses().
func (o *OcspListID) OcspResponses() ([]*OcspResponsesID, error) {
	responses := make([]*OcspResponsesID, 0, len(o.ocspResponses))
	for _, encoded := range o.ocspResponses {
		response, err := ParseOcspResponsesID(encoded)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// CrlOcspRef is
//
//	CrlOcspRef ::= SEQUENCE {
//	    crlids   [0] CRLListID OPTIONAL,
//	    ocspids  [1] OcspListID OPTIONAL,
//	    otherRev [2] OtherRevRefs OPTIONAL }
//
// replacing org.bouncycastle.asn1.esf.CrlOcspRef. It is the member type of both
// CompleteRevocationRefs and AttributeRevocationRefs.
type CrlOcspRef struct {
	// Crlids references the CRLs, nil when the field is absent.
	Crlids *CrlListID
	// Ocspids references the OCSP responses, nil when the field is absent.
	Ocspids *OcspListID
	// OtherRev is the encoding of the otherRev field, nil when absent.
	OtherRev []byte
}

// ParseCrlOcspRef decodes a CrlOcspRef from its encoding.
// Port of CrlOcspRef.getInstance(Object).
func ParseCrlOcspRef(encoded []byte) (*CrlOcspRef, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the CrlOcspRef")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, errors.New("malformed CrlOcspRef: not a SEQUENCE")
	}
	ref := &CrlOcspRef{}
	for _, child := range element.Children() {
		if child.Class() != asn1ber.ClassContextSpecific {
			return nil, errors.New("illegal tag in CrlOcspRef")
		}
		base, err := cmsCRLSourceExplicitBase(child, "CrlOcspRef")
		if err != nil {
			return nil, err
		}
		switch child.TagNumber() {
		case 0:
			members, err := cmsCRLSourceInnerSequenceOf(base, "CRLListID.crls")
			if err != nil {
				return nil, err
			}
			ref.Crlids = &CrlListID{crls: members}
		case 1:
			members, err := cmsCRLSourceInnerSequenceOf(base, "OcspListID.ocspResponses")
			if err != nil {
				return nil, err
			}
			ref.Ocspids = &OcspListID{ocspResponses: members}
		case 2:
			ref.OtherRev = base.Encoded()
		default:
			return nil, fmt.Errorf("illegal tag in CrlOcspRef: %d", child.TagNumber())
		}
	}
	return ref, nil
}

// cmsCRLSourceExplicitBase unwraps an explicitly tagged field, i.e. returns the single value
// the tagged object holds. Port of ASN1TaggedObject#getExplicitBaseObject().
func cmsCRLSourceExplicitBase(element *asn1ber.Element, name string) (*asn1ber.Element, error) {
	if !element.IsConstructed() || len(element.Children()) != 1 {
		return nil, fmt.Errorf("malformed %s: a field is not explicitly tagged", name)
	}
	return element.Children()[0], nil
}

// cmsCRLSourceSequenceOf returns the encoding of every member of a SEQUENCE OF.
func cmsCRLSourceSequenceOf(element *asn1ber.Element, name string) ([][]byte, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed %s: not a SEQUENCE OF", name)
	}
	members := make([][]byte, 0, len(element.Children()))
	for _, child := range element.Children() {
		// BouncyCastle validates each member's shape (CertificateList.getInstance /
		// BasicOCSPResponse.getInstance) while walking the field; both are SEQUENCEs, and a
		// member that is not one makes the whole structure unreadable there as it does here.
		if !child.IsConstructed() {
			return nil, fmt.Errorf("malformed %s: a member is not a SEQUENCE", name)
		}
		members = append(members, child.Encoded())
	}
	return members, nil
}

// cmsCRLSourceInnerSequenceOf unwraps a SEQUENCE holding a single SEQUENCE OF field - the
// shape of both CRLListID and OcspListID - and returns the encoding of every member.
func cmsCRLSourceInnerSequenceOf(element *asn1ber.Element, name string) ([][]byte, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() || len(element.Children()) != 1 {
		return nil, fmt.Errorf("malformed %s: not a SEQUENCE holding one field", name)
	}
	return cmsCRLSourceSequenceOf(element.Children()[0], name)
}

// DSSASN1UtilsRevocationValues builds a RevocationValues from the encoding of an attribute
// value, returning nil when it is nil or cannot be read.
// Port of getRevocationValues(ASN1Encodable), whose ASN1Encodable argument is here that
// value's encoding.
func DSSASN1UtilsRevocationValues(encodable []byte) *RevocationValues {
	if encodable != nil {
		revocationValues, err := ParseRevocationValues(encodable)
		if err == nil {
			return revocationValues
		}
		// Upstream logs "Unable to parse RevocationValues".
	}
	return nil
}

// -----------------------------------------------------------------------------
// CMSCRLSource.
// -----------------------------------------------------------------------------

// CMSCRLSource is a CRLSource that retrieves information from a CMS SignedData container.
// Port of the abstract class CMSCRLSource; the concrete sources of the later phases
// (CAdESCRLSource, TimestampCRLSource) embed the *CMSCRLSource NewCMSCRLSource returns.
//
// As with OfflineCRLSourceBase, this abstract base does not call
// InitOfflineRevocationSource: the outermost concrete source registers itself, so that
// RevocationToken (singular) dispatches to the right RevocationTokens implementation.
type CMSCRLSource struct {
	OfflineCRLSourceBase

	// crls holds the SignedData.crls values, each as its own encoding.
	crls [][]byte

	// unsignedAttributes represents the unsigned properties, nil when the signer carries no
	// unsignedAttrs field.
	//
	// Java distinguishes a null AttributeTable from an empty one; a nil cmscore.Attributes
	// stands for both, which changes nothing here: every branch guarded by that null check
	// only looks attributes up, and an empty table yields none either way.
	unsignedAttributes cmscore.Attributes
}

// NewCMSCRLSource creates a CRL source over the CRLs of a CMS SignedData and the
// revocation-values / revocation-references attributes of a signer.
// Port of the protected CMSCRLSource(Store<X509CRLHolder>, AttributeTable) constructor.
//
// unsignedAttributes is the signer's unsignedAttrs, nil when it has none. The DSSException
// addX509CRLHolder raises for an unreadable CRL of SignedData.crls propagates out of the Java
// constructor; here it is returned as an error.
func NewCMSCRLSource(crls [][]byte, unsignedAttributes cmscore.Attributes) (*CMSCRLSource, error) {
	source := &CMSCRLSource{
		OfflineCRLSourceBase: NewOfflineCRLSourceBase(),
		crls:                 crls,
		unsignedAttributes:   unsignedAttributes,
	}
	if err := source.extract(); err != nil {
		return nil, err
	}
	return source, nil
}

// extract ports the private extract().
func (s *CMSCRLSource) extract() error {
	// Adds CRLs contained in SignedData
	if err := s.collectFromSignedData(); err != nil {
		return err
	}

	if s.unsignedAttributes != nil {
		/*
		 * ETSI TS 101 733 V2.2.1 (2013-04) page 43 6.3.4 revocation-values Attribute
		 * Definition. It holds the values of CRLs and OCSP referenced in the
		 * complete-revocation-references attribute.
		 *
		 * RevocationValues ::= SEQUENCE { crlVals [0] SEQUENCE OF CertificateList
		 * OPTIONAL, ocspVals [1] SEQUENCE OF BasicOCSPResponse OPTIONAL, otherRevVals
		 * [2] OtherRevVals OPTIONAL}
		 */
		s.collectRevocationValues(s.unsignedAttributes, OIDIdAaEtsRevocationValues,
			enumerations.RevocationOriginRevocationValues)

		/*
		 * ETSI TS 101 733 V2.2.1 (2013-04) pages 39,41 6.2.2
		 * complete-revocation-references Attribute Definition and 6.2.4
		 * attribute-revocation-references Attribute Definition.
		 *
		 * CompleteRevocationRefs ::= SEQUENCE OF CrlOcspRef CrlOcspRef ::= SEQUENCE {
		 * crlids [0] CRLListID OPTIONAL, ocspids [1] OcspListID OPTIONAL, otherRev [2]
		 * OtherRevRefs OPTIONAL } AttributeRevocationRefs ::= SEQUENCE OF CrlOcspRef
		 * (the same as for CompleteRevocationRefs)
		 */
		s.collectRevocationRefs(OIDIdAaEtsRevocationRefs, enumerations.RevocationRefOriginCompleteRevocationRefs)

		/*
		 * id-aa-ets-attrRevocationRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
		 * us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 45}
		 */
		s.collectRevocationRefs(OIDAttributeRevocationRefsOid, enumerations.RevocationRefOriginAttributeRevocationRefs)
	}
	return nil
}

// collectFromSignedData ports the private collectFromSignedData().
func (s *CMSCRLSource) collectFromSignedData() error {
	for _, x509CRLHolder := range s.crls {
		if err := s.AddX509CRLHolder(x509CRLHolder, enumerations.RevocationOriginCMSSignedData); err != nil {
			return err
		}
	}
	return nil
}

// collectRevocationValues ports the private
// collectRevocationValues(AttributeTable, ASN1ObjectIdentifier, RevocationOrigin).
func (s *CMSCRLSource) collectRevocationValues(attributeTable cmscore.Attributes,
	revocationValuesAttribute asn1.ObjectIdentifier, origin enumerations.RevocationOrigin) {
	attributes := DSSASN1UtilsAsn1Attributes(attributeTable, revocationValuesAttribute)
	for _, attribute := range attributes {
		attributeValues := attribute.Values
		for _, attrValue := range attributeValues {
			s.extractRevocationValues(attrValue.Encoded(), origin)
		}
	}
}

// extractRevocationValues ports the private extractRevocationValues(ASN1Encodable, RevocationOrigin).
func (s *CMSCRLSource) extractRevocationValues(attrValue []byte, origin enumerations.RevocationOrigin) {
	revValues := DSSASN1UtilsRevocationValues(attrValue)
	if revValues != nil {
		for _, revValue := range revValues.CrlVals {
			if err := s.AddX509CRLHolder(revValue, origin); err != nil {
				// Upstream logs "Unable to process CRL binary : {}".
				continue
			}
		}
	}
}

// AddX509CRLHolder computes and stores a CRLBinary from the encoding of a CertificateList.
// Port of the protected addX509CRLHolder(X509CRLHolder, RevocationOrigin); the DSSException
// it raises is returned as an error carrying the same message.
func (s *CMSCRLSource) AddX509CRLHolder(crlHolder []byte, origin enumerations.RevocationOrigin) error {
	crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(crlHolder)
	if err != nil {
		return model.NewDSSErrorMessageCause(fmt.Sprintf(
			"Unable to parse CRL binaries from origin '%s'. Reason : %s", origin, err.Error()), err)
	}
	s.AddBinary(crlBinary, origin)
	return nil
}

// collectRevocationRefs ports the private
// collectRevocationRefs(ASN1ObjectIdentifier, RevocationRefOrigin).
//
// Upstream wraps the whole method in one try - "When error in computing or in format, the
// algorithm just continues" - so a value that is not a SEQUENCE ends the collection for that
// origin instead of being skipped; the early return reproduces that.
func (s *CMSCRLSource) collectRevocationRefs(revocationRefsAttribute asn1.ObjectIdentifier,
	origin enumerations.RevocationRefOrigin) {
	attributes := DSSASN1UtilsAsn1Attributes(s.unsignedAttributes, revocationRefsAttribute)
	for _, attribute := range attributes {
		attributeValues := attribute.Values
		for _, attrValue := range attributeValues {
			// Java casts to ASN1Sequence; a ClassCastException lands in the method-wide catch.
			if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
				// Upstream logs "An error occurred during extraction of revocation
				// references from signature unsigned properties. Revocations for origin {}
				// were not stored".
				return
			}
			for _, member := range attrValue.Children() {
				s.collectRevocationRefFromASN1Encodable(member, origin)
			}
		}
	}
}

// collectRevocationRefFromASN1Encodable ports the private
// collectRevocationRefFromASN1Encodable(ASN1Encodable, RevocationRefOrigin).
func (s *CMSCRLSource) collectRevocationRefFromASN1Encodable(asn1Encodable *asn1ber.Element,
	origin enumerations.RevocationRefOrigin) {
	crlOcspRef, err := ParseCrlOcspRef(asn1Encodable.Encoded())
	if err != nil {
		// Upstream logs "Unable to process CRL reference : {}".
		return
	}
	crlIds := crlOcspRef.Crlids
	if crlIds != nil {
		ids, err := crlIds.Crls()
		if err != nil {
			return
		}
		for _, id := range ids {
			crlRef, err := NewCRLRefFromCrlValidatedID(id)
			if err != nil {
				return
			}
			s.AddRevocationReference(crlRef, origin)
		}
	}
}
