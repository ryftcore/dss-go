// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESOCSPSource.java
// (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESOCSPSource extracts and stores OCSPs from a JAdES signature. Port of the class
// OCSPSource, extending spi.OfflineOCSPSourceBase.
type OCSPSource struct {
	spi.OfflineOCSPSourceBase

	// etsiUHeader represents the unsigned 'etsiU' header. Port of the private transient final
	// EtsiUHeader etsiUHeader field.
	etsiUHeader *EtsiUHeader
}

// NewJAdESOCSPSource is the default constructor. Port of the public
// OCSPSource(EtsiUHeader) constructor.
//
// Panics with the Java message when etsiUHeader is missing (Objects.requireNonNull).
func NewJAdESOCSPSource(etsiUHeader *EtsiUHeader) *OCSPSource {
	if etsiUHeader == nil {
		panic("etsiUHeader cannot be null")
	}
	s := &OCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
		etsiUHeader:           etsiUHeader,
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches OfflineOCSPSourceBase.RevocationTokens; this class
	// does not override it.
	s.InitOfflineRevocationSource(s)

	s.extractEtsiU()

	return s
}

func (s *OCSPSource) extractEtsiU() {
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

func (s *OCSPSource) extractRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRVals == attribute.HeaderName() {
		s.extractOCSPValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRVals),
			enumerations.RevocationOriginRevocationValues)
	}
}

func (s *OCSPSource) extractAttributeRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArVals == attribute.HeaderName() {
		s.extractOCSPValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArVals),
			enumerations.RevocationOriginAttributeRevocationValues)
	}
}

func (s *OCSPSource) extractTimestampValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesTstVD, enumerations.RevocationOriginTimestampValidationData)
}

func (s *OCSPSource) extractAnyValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesAnyValData, enumerations.RevocationOriginAnyValidationData)
}

func (s *OCSPSource) extractValidationData(attribute *EtsiUComponent, headerName string, origin enumerations.RevocationOrigin) {
	if headerName == attribute.HeaderName() {
		tstVd := DSSJsonUtilsToMap(attribute.Value(), headerName)
		if tstVd.Size() != 0 {
			rVals := DSSJsonUtilsGetAsMap(tstVd, JAdESHeaderParameterNamesRVals)
			if rVals.Size() != 0 {
				s.extractOCSPValues(rVals, origin)
			}
		}
	}
}

func (s *OCSPSource) extractCompleteRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRRefs == attribute.HeaderName() {
		s.extractOCSPReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRRefs),
			enumerations.RevocationRefOriginCompleteRevocationRefs)
	}
}

func (s *OCSPSource) extractAttributeRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArRefs == attribute.HeaderName() {
		s.extractOCSPReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArRefs),
			enumerations.RevocationRefOriginAttributeRevocationRefs)
	}
}

func (s *OCSPSource) extractOCSPValues(rVals *jose.Object, origin enumerations.RevocationOrigin) {
	ocspVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesOcspVals)
	for _, item := range ocspVals {
		pkiOb := DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesPkiOb)
		s.extractOCSPFromPkiOb(pkiOb, origin)
	}
}

func (s *OCSPSource) extractOCSPFromPkiOb(pkiOb *jose.Object, origin enumerations.RevocationOrigin) {
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

func (s *OCSPSource) add(ocspValueDerB64 string, origin enumerations.RevocationOrigin) {
	basicOCSPResp, err := spi.DSSRevocationUtilsLoadOCSPBase64Encoded(ocspValueDerB64)
	if err != nil {
		// Upstream logs "Unable to extract OCSP from '{}'".
		return
	}
	ocspResponseBinary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
	if err != nil {
		// Upstream logs "Unable to extract OCSP from '{}'".
		return
	}
	s.AddBinary(ocspResponseBinary, origin)
}

func (s *OCSPSource) extractOCSPReferences(rRefs *jose.Object, origin enumerations.RevocationRefOrigin) {
	ocspRefs := DSSJsonUtilsGetAsList(rRefs, JAdESHeaderParameterNamesOcspRefs)
	for _, item := range ocspRefs {
		ocspRefMap := DSSJsonUtilsToMapValue(item)
		if ocspRefMap.Size() != 0 {
			ocspRef := RevocationRefExtractionUtilsCreateOCSPRef(ocspRefMap)
			if ocspRef != nil {
				s.AddRevocationReference(ocspRef, origin)
			}
		}
	}
}
