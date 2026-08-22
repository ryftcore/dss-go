// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CMSOCSPSource.java (DSS 6.5.RC1).
//
// BouncyCastle replacements used here (see PORTING.md):
//
//   - the two Store<?> of OCSP responses -> [][]byte, the encoding of the otherRevInfo value
//     of every OtherRevocationInfoFormat of SignedData.crls carrying the corresponding
//     format OID, i.e. what cmscore.CMS.OCSPResponses() (id-ri-ocsp-response) and
//     cmscore.CMS.OCSPBasicResponses() (id-pkix-ocsp-basic) hand out. Store#getMatches(null)
//     is then iterating the slice, and "instanceof ASN1Sequence" is a shape check on the
//     parsed value.
//   - org.bouncycastle.asn1.cms.AttributeTable -> cmscore.Attributes.
//   - org.bouncycastle.asn1.esf.{RevocationValues, CrlOcspRef, OcspListID} -> the types
//     cms_crl_source.go defines, since both revocation sources read them.
package spi

import (
	"encoding/asn1"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
)

// CMSObjectIdentifierIDRIOcspResponse is id-ri-ocsp-response, the OtherRevocationInfoFormat
// under which SignedData.crls carries a complete RFC 6960 OCSPResponse. Corresponds to
// CMSObjectIdentifiers.id_ri_ocsp_response.
var CMSObjectIdentifierIDRIOcspResponse = cmscore.OIDRIOCSPResponse

// CMSOCSPSource is an OCSPSource that retrieves information from a CMS SignedData container.
// Port of the abstract class CMSOCSPSource; the concrete sources of the later phases
// (OCSPSource, TimestampOCSPSource) embed the *CMSOCSPSource NewCMSOCSPSource returns.
//
// As with OfflineOCSPSourceBase, this abstract base does not call
// InitOfflineRevocationSource: the outermost concrete source registers itself, so that
// RevocationToken (singular) dispatches to the right RevocationTokens implementation.
type CMSOCSPSource struct {
	OfflineOCSPSourceBase

	// ocspResponsesStore holds the encoding of each OCSPResponse carried under
	// id-ri-ocsp-response, nil when the CMS has none.
	ocspResponsesStore [][]byte

	// ocspBasicStore holds the encoding of each BasicOCSPResponse carried under
	// id-pkix-ocsp-basic, nil when the CMS has none.
	ocspBasicStore [][]byte

	// UnsignedAttributes represents the unsigned properties. It mirrors the protected field
	// of the same name, which subclasses read; it is nil when the signer carries no
	// unsignedAttrs field.
	//
	// Java distinguishes a null AttributeTable from an empty one; a nil cmscore.Attributes
	// stands for both, which changes nothing here: every branch guarded by that null check
	// only looks attributes up, and an empty table yields none either way.
	UnsignedAttributes cmscore.Attributes
}

// NewCMSOCSPSource creates an OCSP source over the OCSP responses of a CMS SignedData and the
// revocation-values / revocation-references attributes of a signer.
// Port of the protected CMSOCSPSource(Store<?>, Store<?>, AttributeTable) constructor.
//
// unsignedAttributes is the signer's unsignedAttrs, nil when it has none. The runtime
// exceptions appendContainedOCSPResponses lets escape - a response of SignedData.crls that
// cannot be turned into a binary, and the ClassCastException of collectRevocationRefs -
// propagate out of the Java constructor; here they are returned as errors.
func NewCMSOCSPSource(ocspResponsesStore, ocspBasicStore [][]byte,
	unsignedAttributes cmscore.Attributes) (*CMSOCSPSource, error) {
	source := &CMSOCSPSource{
		OfflineOCSPSourceBase: NewOfflineOCSPSourceBase(),
		ocspResponsesStore:    ocspResponsesStore,
		ocspBasicStore:        ocspBasicStore,
		UnsignedAttributes:    unsignedAttributes,
	}
	if err := source.appendContainedOCSPResponses(); err != nil {
		return nil, err
	}
	return source, nil
}

// appendContainedOCSPResponses ports the private method of the same name.
func (s *CMSOCSPSource) appendContainedOCSPResponses() error {
	// Add OCSPs from SignedData
	if err := s.collectFromSignedData(); err != nil {
		return err
	}

	if s.UnsignedAttributes != nil {
		/*
		   ETSI TS 101 733 V2.2.1 (2013-04) page 43
		   6.3.4 revocation-values Attribute Definition
		   The revocation-values attribute is an unsigned attribute. Only a single instance
		   of this attribute shall occur with an electronic signature. It holds the values
		   of CRLs and OCSP referenced in the complete-revocation-references attribute.

		   RevocationValues ::= SEQUENCE {
		   crlVals [0] SEQUENCE OF CertificateList OPTIONAL,
		   ocspVals [1] SEQUENCE OF BasicOCSPResponse OPTIONAL,
		   otherRevVals [2] OtherRevVals OPTIONAL}
		*/
		s.collectRevocationValues(s.UnsignedAttributes, OIDIdAaEtsRevocationValues,
			enumerations.RevocationOriginRevocationValues)

		/*
		 * ETSI TS 101 733 V2.2.1 (2013-04) pages 39,41
		 * 6.2.2 complete-revocation-references Attribute Definition and
		 * 6.2.4 attribute-revocation-references Attribute Definition
		 *
		 * CompleteRevocationRefs ::= SEQUENCE OF CrlOcspRef
		 * CrlOcspRef ::= SEQUENCE {
		 *  crlids [0] CRLListID OPTIONAL,
		 *  ocspids [1] OcspListID OPTIONAL,
		 *  otherRev [2] OtherRevRefs OPTIONAL
		 * }
		 * AttributeRevocationRefs ::= SEQUENCE OF CrlOcspRef (the same as for
		 * CompleteRevocationRefs)
		 */
		if err := s.collectRevocationRefs(s.UnsignedAttributes, OIDIdAaEtsRevocationRefs,
			enumerations.RevocationRefOriginCompleteRevocationRefs); err != nil {
			return err
		}
		/*
		 * id-aa-ets-attrRevocationRefs OBJECT IDENTIFIER ::= { iso(1) member-body(2)
		 * us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 45}
		 */
		if err := s.collectRevocationRefs(s.UnsignedAttributes, OIDAttributeRevocationRefsOid,
			enumerations.RevocationRefOriginAttributeRevocationRefs); err != nil {
			return err
		}
	}
	return nil
}

// collectFromSignedData ports the private collectFromSignedData().
func (s *CMSOCSPSource) collectFromSignedData() error {
	if err := s.addBasicOcspRespFromIDRIOcspResponse(); err != nil {
		return err
	}
	return s.addBasicOcspRespFromIDPkixOcspBasic()
}

// addBasicOcspRespFromIDRIOcspResponse ports the private
// addBasicOcspRespFrom_id_ri_ocsp_response().
//
// DEVIATION: upstream builds the OCSPResponseBinary without checking that the response could
// be decoded, so an unreadable entry raises a NullPointerException out of the constructor;
// the error returned here stands in for it, with a message naming the actual problem.
func (s *CMSOCSPSource) addBasicOcspRespFromIDRIOcspResponse() error {
	if s.ocspResponsesStore == nil {
		return nil
	}
	for _, object := range s.ocspResponsesStore {
		otherRevocationInfoMatch, err := cmsOCSPSourceAsSequence(object)
		if err != nil {
			// Upstream logs "Unsupported object type for id_ri_ocsp_response (SHALL be an
			// ASN1Sequence) : {}".
			continue
		}
		var basicOCSPResp *BasicOCSPResp
		if len(otherRevocationInfoMatch.Children()) == 4 {
			basicOCSPResp = DSSRevocationUtilsBasicOcspResp(object)
		} else {
			// NOTE: the expected encoding
			ocspResp := DSSRevocationUtilsOcspResp(object)
			if ocspResp == nil {
				return model.NewDSSError("Unable to process an OCSP response of SignedData.crls : " +
					"the id_ri_ocsp_response entry is not an OCSPResponse")
			}
			basicOCSPResp = DSSRevocationUtilsFromRespToBasic(ocspResp)
		}

		ocspResponseIdentifier, err := OCSPResponseBinaryBuild(basicOCSPResp)
		if err != nil {
			return err
		}
		ocspResponseIdentifier.SetASN1ObjectIdentifier(CMSObjectIdentifierIDRIOcspResponse)
		s.AddBinary(ocspResponseIdentifier, enumerations.RevocationOriginCMSSignedData)
	}
	return nil
}

// addBasicOcspRespFromIDPkixOcspBasic ports the private
// addBasicOcspRespFrom_id_pkix_ocsp_basic().
func (s *CMSOCSPSource) addBasicOcspRespFromIDPkixOcspBasic() error {
	if s.ocspBasicStore == nil {
		return nil
	}
	for _, object := range s.ocspBasicStore {
		if _, err := cmsOCSPSourceAsSequence(object); err != nil {
			// Upstream logs "Unsupported object type for id_pkix_ocsp_basic (SHALL be an
			// ASN1Sequence) : {}".
			continue
		}
		basicOCSPResp := DSSRevocationUtilsBasicOcspResp(object)
		if basicOCSPResp != nil {
			ocspResponseIdentifier, err := OCSPResponseBinaryBuild(basicOCSPResp)
			if err != nil {
				return err
			}
			ocspResponseIdentifier.SetASN1ObjectIdentifier(OCSPObjectIdentifierIDPkixOcspBasic)
			s.AddBinary(ocspResponseIdentifier, enumerations.RevocationOriginCMSSignedData)
		}
		// Upstream logs "Unable to create an OCSP response from an objects. The entry is
		// skipped." otherwise.
	}
	return nil
}

// cmsOCSPSourceAsSequence parses a member of one of the two stores and checks that it is an
// ASN.1 SEQUENCE, standing in for Java's "object instanceof ASN1Sequence".
func cmsOCSPSourceAsSequence(object []byte) (*asn1ber.Element, error) {
	element, rest, err := asn1ber.Parse(object)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("extra data found after the revocation info value")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("the revocation info value is not an ASN1Sequence")
	}
	return element, nil
}

// collectRevocationValues ports the private
// collectRevocationValues(AttributeTable, ASN1ObjectIdentifier, RevocationOrigin).
func (s *CMSOCSPSource) collectRevocationValues(attributeTable cmscore.Attributes,
	revocationValueAttributes asn1.ObjectIdentifier, origin enumerations.RevocationOrigin) {
	attributes := DSSASN1UtilsAsn1Attributes(attributeTable, revocationValueAttributes)
	for _, attribute := range attributes {
		attributeValues := attribute.Values
		for _, attrValue := range attributeValues {
			s.extractRevocationValues(attrValue.Encoded(), origin)
		}
	}
}

// extractRevocationValues ports the private
// extractRevocationValues(ASN1Encodable, RevocationOrigin).
//
// TODO (upstream): should add also OtherRevVals, but "The syntax and semantics of the other
// revocation values (OtherRevVals) are outside the scope of the present document."
func (s *CMSOCSPSource) extractRevocationValues(attrValue []byte, origin enumerations.RevocationOrigin) {
	revocationValues := DSSASN1UtilsRevocationValues(attrValue)
	if revocationValues != nil {
		for _, basicOCSPResponse := range revocationValues.OcspVals {
			basicOCSPResp := DSSRevocationUtilsBasicOcspResp(basicOCSPResponse)
			if basicOCSPResp == nil {
				// Upstream logs "Unable to process OCSP binary : {}"; the BasicOCSPResp
				// constructor is what raises there.
				continue
			}
			ocspResponseIdentifier, err := OCSPResponseBinaryBuild(basicOCSPResp)
			if err != nil {
				continue
			}
			s.AddBinary(ocspResponseIdentifier, origin)
		}
	}
}

// collectRevocationRefs ports the private
// collectRevocationRefs(AttributeTable, ASN1ObjectIdentifier, RevocationRefOrigin).
//
// Two upstream quirks are reproduced verbatim: a value the attribute holds in anything other
// than a single copy makes the method return - abandoning the attributes that follow, not
// just that one - and the cast of that value to an ASN1Sequence is outside the try, so a
// ClassCastException escapes the constructor. The latter is an error here.
func (s *CMSOCSPSource) collectRevocationRefs(unsignedAttributes cmscore.Attributes,
	revocationReferencesAttribute asn1.ObjectIdentifier, origin enumerations.RevocationRefOrigin) error {
	attributes := DSSASN1UtilsAsn1Attributes(unsignedAttributes, revocationReferencesAttribute)
	if len(attributes) == 0 {
		return nil
	}
	for _, attribute := range attributes {
		attrValue := DSSASN1UtilsAsn1Encodable(attribute)
		if attrValue == nil {
			return nil
		}

		if !attrValue.IsUniversal(asn1ber.TagSequence) || !attrValue.IsConstructed() {
			return model.NewDSSError("The revocation references attribute value is not an ASN1Sequence")
		}
		for _, member := range attrValue.Children() {
			otherCertID, err := ParseCrlOcspRef(member.Encoded())
			if err != nil {
				// Upstream logs "Unable to process OCSP reference : {}".
				continue
			}
			ocspListID := otherCertID.Ocspids
			if ocspListID != nil {
				ocspResponses, err := ocspListID.OcspResponses()
				if err != nil {
					continue
				}
				for _, ocspResponsesID := range ocspResponses {
					ocspRef, err := NewOCSPRefFromOcspResponsesID(ocspResponsesID)
					if err != nil {
						break
					}
					s.AddRevocationReference(ocspRef, origin)
				}
			}
		}
	}
	return nil
}
