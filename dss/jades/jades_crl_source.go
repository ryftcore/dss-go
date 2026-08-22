// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESCRLSource.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: JAdESEtsiUHeader/EtsiUComponent - see jades_certificate_source.go's header.
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
// JAdESCRLSource, extending spi.OfflineCRLSourceBase.
type JAdESCRLSource struct {
	spi.OfflineCRLSourceBase

	// etsiUHeader represents the unsigned 'etsiU' header. Port of the private transient final
	// JAdESEtsiUHeader etsiUHeader field.
	etsiUHeader *JAdESEtsiUHeader
}

// NewJAdESCRLSource is the default constructor. Port of the public JAdESCRLSource(JAdESEtsiUHeader)
// constructor.
//
// Panics with the Java message when etsiUHeader is missing (Objects.requireNonNull).
func NewJAdESCRLSource(etsiUHeader *JAdESEtsiUHeader) *JAdESCRLSource {
	if etsiUHeader == nil {
		panic("etsiUComponents cannot be null")
	}
	s := &JAdESCRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		etsiUHeader:          etsiUHeader,
	}
	s.InitOfflineRevocationSource(s)

	s.extractEtsiU()

	return s
}

func (s *JAdESCRLSource) extractEtsiU() {
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

func (s *JAdESCRLSource) extractRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRVals == attribute.HeaderName() {
		s.extractCRLValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRVals),
			enumerations.RevocationOrigin_REVOCATION_VALUES)
	}
}

func (s *JAdESCRLSource) extractAttributeRevocationValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArVals == attribute.HeaderName() {
		s.extractCRLValues(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArVals),
			enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES)
	}
}

func (s *JAdESCRLSource) extractTimestampValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesTstVD, enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA)
}

func (s *JAdESCRLSource) extractAnyValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesAnyValData, enumerations.RevocationOrigin_ANY_VALIDATION_DATA)
}

func (s *JAdESCRLSource) extractValidationData(attribute *EtsiUComponent, headerName string, origin enumerations.RevocationOrigin) {
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

func (s *JAdESCRLSource) extractCompleteRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesRRefs == attribute.HeaderName() {
		s.extractCRLReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesRRefs),
			enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)
	}
}

func (s *JAdESCRLSource) extractAttributeRevocationRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesArRefs == attribute.HeaderName() {
		s.extractCRLReferences(DSSJsonUtilsToMap(attribute.Value(), JAdESHeaderParameterNamesArRefs),
			enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)
	}
}

func (s *JAdESCRLSource) extractCRLValues(rVals *jose.Object, origin enumerations.RevocationOrigin) {
	crlVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesCrlVals)
	for _, item := range crlVals {
		pkiOb := DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesPkiOb)
		s.extractCRLFromPkiOb(pkiOb, origin)
	}
}

func (s *JAdESCRLSource) extractCRLFromPkiOb(pkiOb *jose.Object, origin enumerations.RevocationOrigin) {
	if pkiOb.Size() != 0 {
		encoding := DSSJsonUtilsGetAsString(pkiOb, JAdESHeaderParameterNamesEncoding)
		if utils.IsStringEmpty(encoding) || utils.AreStringsEqual(enumerations.PKIEncoding_DER.URI(), encoding) {
			val := DSSJsonUtilsGetAsString(pkiOb, JAdESHeaderParameterNamesVal)
			if utils.IsStringNotEmpty(val) {
				s.add(val, origin)
			}
		} else {
			// Upstream logs "Unsupported encoding '{}'".
		}
	}
}

func (s *JAdESCRLSource) add(crlValueDerB64 string, origin enumerations.RevocationOrigin) {
	crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(utils.FromBase64(crlValueDerB64))
	if err != nil {
		// Upstream logs "Unable to extract CRL from '{}'. Reason : {}".
		return
	}
	s.AddBinary(crlBinary, origin)
}

func (s *JAdESCRLSource) extractCRLReferences(rRefs *jose.Object, origin enumerations.RevocationRefOrigin) {
	crlRefs := DSSJsonUtilsGetAsList(rRefs, JAdESHeaderParameterNamesCrlRefs)
	for _, item := range crlRefs {
		crlRefMap := DSSJsonUtilsToMapValue(item)
		if crlRefMap.Size() != 0 {
			crlRef := JAdESRevocationRefExtractionUtilsCreateCRLRef(crlRefMap)
			if crlRef != nil {
				s.AddRevocationReference(crlRef, origin)
			}
		}
	}
}

// compile-time assertion: a JAdESCRLSource satisfies spi.OfflineRevocationSourceOverrides.
var _ spi.OfflineRevocationSourceOverrides[revocation.CRL] = (*JAdESCRLSource)(nil)
