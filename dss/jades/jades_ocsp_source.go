// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESOCSPSource.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: JAdESEtsiUHeader/EtsiUComponent - see jades_certificate_source.go's header.
package jades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESOCSPSource extracts and stores OCSPs from a JAdES signature. Port of the class
// JAdESOCSPSource, extending spi.OfflineOCSPSourceBase.
type JAdESOCSPSource struct {
	spi.OfflineOCSPSourceBase

	// etsiUHeader represents the unsigned 'etsiU' header. Port of the private transient final
	// JAdESEtsiUHeader etsiUHeader field.
	etsiUHeader *JAdESEtsiUHeader
}

// NewJAdESOCSPSource is the default constructor. Port of the public
// JAdESOCSPSource(JAdESEtsiUHeader) constructor.
//
// Panics with the Java message when etsiUHeader is missing (Objects.requireNonNull).
func NewJAdESOCSPSource(etsiUHeader *JAdESEtsiUHeader) *JAdESOCSPSource {
	if etsiUHeader == nil {
		panic("etsiUHeader cannot be null")
	}
	s := &JAdESOCSPSource{
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

func (s *JAdESOCSPSource) extractEtsiU() {
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

func (s *JAdESOCSPSource) extractRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRVals == attribute.HeaderName() {
		s.extractOCSPValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRVals),
			enumerations.RevocationOriginRevocationValues)
	}
}

func (s *JAdESOCSPSource) extractAttributeRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArVals == attribute.HeaderName() {
		s.extractOCSPValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArVals),
			enumerations.RevocationOriginAttributeRevocationValues)
	}
}

func (s *JAdESOCSPSource) extractTimestampValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesTstVD, enumerations.RevocationOriginTimestampValidationData)
}

func (s *JAdESOCSPSource) extractAnyValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesAnyValData, enumerations.RevocationOriginAnyValidationData)
}

func (s *JAdESOCSPSource) extractValidationData(attribute *EtsiUComponent, headerName string, origin enumerations.RevocationOrigin) {
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

func (s *JAdESOCSPSource) extractCompleteRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRRefs == attribute.HeaderName() {
		s.extractOCSPReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRRefs),
			enumerations.RevocationRefOriginCompleteRevocationRefs)
	}
}

func (s *JAdESOCSPSource) extractAttributeRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArRefs == attribute.HeaderName() {
		s.extractOCSPReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArRefs),
			enumerations.RevocationRefOriginAttributeRevocationRefs)
	}
}

func (s *JAdESOCSPSource) extractOCSPValues(rVals *jose.Object, origin enumerations.RevocationOrigin) {
	ocspVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesOcspVals)
	for _, item := range ocspVals {
		pkiOb := DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesPkiOb)
		s.extractOCSPFromPkiOb(pkiOb, origin)
	}
}

func (s *JAdESOCSPSource) extractOCSPFromPkiOb(pkiOb *jose.Object, origin enumerations.RevocationOrigin) {
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

func (s *JAdESOCSPSource) add(ocspValueDerB64 string, origin enumerations.RevocationOrigin) {
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

func (s *JAdESOCSPSource) extractOCSPReferences(rRefs *jose.Object, origin enumerations.RevocationRefOrigin) {
	ocspRefs := DSSJsonUtilsGetAsList(rRefs, JAdESHeaderParameterNamesOcspRefs)
	for _, item := range ocspRefs {
		ocspRefMap := DSSJsonUtilsToMapValue(item)
		if ocspRefMap.Size() != 0 {
			ocspRef := JAdESRevocationRefExtractionUtilsCreateOCSPRef(ocspRefMap)
			if ocspRef != nil {
				s.AddRevocationReference(ocspRef, origin)
			}
		}
	}
}
