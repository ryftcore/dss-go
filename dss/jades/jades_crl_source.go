// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESCRLSource.java
// (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESCRLSource extracts and stores CRLs from a JAdES signature. Port of the class
// CRLSource, extending spi.OfflineCRLSourceBase.
type CRLSource struct {
	spi.OfflineCRLSourceBase

	// etsiUHeader represents the unsigned 'etsiU' header. Port of the private transient final
	// EtsiUHeader etsiUHeader field.
	etsiUHeader *EtsiUHeader
}

// NewJAdESCRLSource is the default constructor. Port of the public JAdESCRLSource(JAdESEtsiUHeader)
// constructor.
//
// Panics with the Java message when etsiUHeader is missing (Objects.requireNonNull).
func NewCRLSource(etsiUHeader *EtsiUHeader) *CRLSource {
	if etsiUHeader == nil {
		panic("etsiUComponents cannot be null")
	}
	s := &CRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		etsiUHeader:          etsiUHeader,
	}
	s.InitOfflineRevocationSource(s)

	s.extractEtsiU()

	return s
}

func (s *CRLSource) extractEtsiU() {
	if !s.etsiUHeader.IsExist() {
		return
	}

	for _, attribute := range s.etsiUHeader.Attributes() {
		s.extractRevocationValues(attribute)
		s.extractAttributeRevocationValues(attribute)
		s.extractTimestampValidationData(attribute)
		s.extractAnyValidationData(attribute)

		s.extractCompleteRevocationRefs(attribute)
		s.extractAttributeRevocationRefs(attribute)
	}
}

func (s *CRLSource) extractRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRVals == attribute.HeaderName() {
		s.extractCRLValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRVals),
			enumerations.RevocationOriginRevocationValues)
	}
}

func (s *CRLSource) extractAttributeRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArVals == attribute.HeaderName() {
		s.extractCRLValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArVals),
			enumerations.RevocationOriginAttributeRevocationValues)
	}
}

func (s *CRLSource) extractTimestampValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesTstVD, enumerations.RevocationOriginTimestampValidationData)
}

func (s *CRLSource) extractAnyValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesAnyValData, enumerations.RevocationOriginAnyValidationData)
}

func (s *CRLSource) extractValidationData(attribute *EtsiUComponent, headerName string, origin enumerations.RevocationOrigin) {
	if headerName == attribute.HeaderName() {
		tstVd := DSSJsonUtilsToMap(attribute.Value(), headerName)
		if tstVd.Size() != 0 {
			rVals := DSSJsonUtilsGetAsMap(tstVd, JAdESHeaderParameterNamesRVals)
			if rVals.Size() != 0 {
				s.extractCRLValues(rVals, origin)
			}
		}
	}
}

func (s *CRLSource) extractCompleteRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRRefs == attribute.HeaderName() {
		s.extractCRLReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRRefs),
			enumerations.RevocationRefOriginCompleteRevocationRefs)
	}
}

func (s *CRLSource) extractAttributeRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArRefs == attribute.HeaderName() {
		s.extractCRLReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArRefs),
			enumerations.RevocationRefOriginAttributeRevocationRefs)
	}
}

func (s *CRLSource) extractCRLValues(rVals *jose.Object, origin enumerations.RevocationOrigin) {
	crlVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesCrlVals)
	for _, item := range crlVals {
		pkiOb := DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesPkiOb)
		s.extractCRLFromPkiOb(pkiOb, origin)
	}
}

func (s *CRLSource) extractCRLFromPkiOb(pkiOb *jose.Object, origin enumerations.RevocationOrigin) {
	if pkiOb.Size() != 0 {
		encoding := DSSJsonUtilsGetAsString(pkiOb, JAdESHeaderParameterNamesEncoding)
		if utils.IsStringEmpty(encoding) || utils.AreStringsEqual(enumerations.PKIEncodingDER.URI(), encoding) {
			val := DSSJsonUtilsGetAsString(pkiOb, JAdESHeaderParameterNamesVal)
			if utils.IsStringNotEmpty(val) {
				s.add(val, origin)
			}
		} else {
			// Upstream logs "Unsupported encoding '{}'".
		}
	}
}

func (s *CRLSource) add(crlValueDerB64 string, origin enumerations.RevocationOrigin) {
	crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(utils.FromBase64(crlValueDerB64))
	if err != nil {
		// Upstream logs "Unable to extract CRL from '{}'. Reason : {}".
		return
	}
	s.AddBinary(crlBinary, origin)
}

func (s *CRLSource) extractCRLReferences(rRefs *jose.Object, origin enumerations.RevocationRefOrigin) {
	crlRefs := DSSJsonUtilsGetAsList(rRefs, JAdESHeaderParameterNamesCrlRefs)
	for _, item := range crlRefs {
		crlRefMap := DSSJsonUtilsToMapValue(item)
		if crlRefMap.Size() != 0 {
			crlRef := RevocationRefExtractionUtilsCreateCRLRef(crlRefMap)
			if crlRef != nil {
				s.AddRevocationReference(crlRef, origin)
			}
		}
	}
}

// compile-time assertion: a CRLSource satisfies spi.OfflineRevocationSourceOverrides.
var _ spi.OfflineRevocationSourceOverrides[revocation.CRL] = (*CRLSource)(nil)
