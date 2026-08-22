// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/evidencerecord/ASiCWithCAdESContainerEvidenceRecordBuilder.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// ASiCWithCAdESContainerEvidenceRecordBuilder validates and incorporates an existing Evidence
// Record within an ASiC with CAdES container.
type ASiCWithCAdESContainerEvidenceRecordBuilder struct {
	asic.AbstractASiCContainerEvidenceRecordBuilder
}

var _ asic.AbstractASiCContainerEvidenceRecordBuilderOverrides = (*ASiCWithCAdESContainerEvidenceRecordBuilder)(nil)

// NewASiCWithCAdESContainerEvidenceRecordBuilder is the default constructor. Ports
// ASiCWithCAdESContainerEvidenceRecordBuilder(CertificateVerifier,
// EvidenceRecordFilenameFactory).
func NewASiCWithCAdESContainerEvidenceRecordBuilder(certificateVerifier validation.CertificateVerifier,
	asicFilenameFactory asic.EvidenceRecordFilenameFactory) *ASiCWithCAdESContainerEvidenceRecordBuilder {
	b := &ASiCWithCAdESContainerEvidenceRecordBuilder{
		AbstractASiCContainerEvidenceRecordBuilder: asic.NewAbstractASiCContainerEvidenceRecordBuilderBase(certificateVerifier, asicFilenameFactory),
	}
	b.InitAbstractASiCContainerEvidenceRecordBuilder(b)
	return b
}

// GetASiCContentBuilder ports the @Override protected getASiCContentBuilder().
func (b *ASiCWithCAdESContainerEvidenceRecordBuilder) GetASiCContentBuilder() *asic.AbstractASiCContentBuilder {
	return NewASiCWithCAdESASiCContentBuilder().AbstractASiCContentBuilder
}

// AssertEvidenceRecordFilenameValid ports the @Override protected
// assertEvidenceRecordFilenameValid(String, EvidenceRecordTypeEnum, Content). Named exported
// here (rather than shadowing the embedded base's unexported assertEvidenceRecordFilenameValid)
// since the base dispatches self-calls through overrides per the virtual-dispatch precedent -
// see the caller in the base's Build().
func (b *ASiCWithCAdESContainerEvidenceRecordBuilder) AssertEvidenceRecordFilenameValid(evidenceRecordFilename string, evidenceRecordType enumerations.EvidenceRecordTypeEnum, asicContent *asic.Content) {
	b.AbstractASiCContainerEvidenceRecordBuilder.AssertEvidenceRecordFilenameValid(evidenceRecordFilename, evidenceRecordType, asicContent)

	if enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord == evidenceRecordType &&
		!asic.UtilsIsAsn1EvidenceRecord(evidenceRecordFilename) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("RFC 4998 Evidence Record's filename '%s' is "+
			"not compliant to the ASiC with CAdES filename convention!", evidenceRecordFilename)))
	} else if enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord == evidenceRecordType &&
		!asic.UtilsIsXmlEvidenceRecord(evidenceRecordFilename) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("RFC 6283 XML Evidence Record's filename '%s' is "+
			"not compliant to the ASiC with CAdES filename convention!", evidenceRecordFilename)))
	}
}
