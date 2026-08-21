// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/evidencerecord/ASiCWithXAdESContainerEvidenceRecordBuilder.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.evidencerecord lands in
// this same Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithXAdESContainerEvidenceRecordBuilder validates and incorporates an existing Evidence
// Record within an ASiC with XAdES container.
type ASiCWithXAdESContainerEvidenceRecordBuilder struct {
	asic.AbstractASiCContainerEvidenceRecordBuilder
}

var _ asic.AbstractASiCContainerEvidenceRecordBuilderOverrides = (*ASiCWithXAdESContainerEvidenceRecordBuilder)(nil)

// NewASiCWithXAdESContainerEvidenceRecordBuilder is the default constructor. Ports
// ASiCWithXAdESContainerEvidenceRecordBuilder(CertificateVerifier,
// ASiCEvidenceRecordFilenameFactory).
func NewASiCWithXAdESContainerEvidenceRecordBuilder(certificateVerifier validation.CertificateVerifier,
	asicFilenameFactory asic.ASiCEvidenceRecordFilenameFactory) *ASiCWithXAdESContainerEvidenceRecordBuilder {
	b := &ASiCWithXAdESContainerEvidenceRecordBuilder{
		AbstractASiCContainerEvidenceRecordBuilder: asic.NewAbstractASiCContainerEvidenceRecordBuilderBase(certificateVerifier, asicFilenameFactory),
	}
	b.InitAbstractASiCContainerEvidenceRecordBuilder(b)
	return b
}

// GetASiCContentBuilder ports the @Override protected getASiCContentBuilder().
func (b *ASiCWithXAdESContainerEvidenceRecordBuilder) GetASiCContentBuilder() *asic.AbstractASiCContentBuilder {
	return NewASiCWithXAdESASiCContentBuilder().AbstractASiCContentBuilder
}

// AssertEvidenceRecordFilenameValid ports the @Override protected
// assertEvidenceRecordFilenameValid(String, EvidenceRecordTypeEnum, ASiCContent). Named exported
// here (rather than shadowing the embedded base's unexported assertEvidenceRecordFilenameValid)
// since the base dispatches self-calls through overrides per the virtual-dispatch precedent -
// see the caller in the base's Build().
//
// Panics with an *exception.IllegalInputException matching Java's IllegalInputException.
func (b *ASiCWithXAdESContainerEvidenceRecordBuilder) AssertEvidenceRecordFilenameValid(evidenceRecordFilename string, evidenceRecordType enumerations.EvidenceRecordTypeEnum, asicContent *asic.ASiCContent) {
	evidenceRecordDocuments := asicContent.EvidenceRecordDocuments()
	if utils.IsCollectionNotEmpty(evidenceRecordDocuments) {
		for _, name := range spi.DSSUtilsDocumentNames(evidenceRecordDocuments) {
			if name == evidenceRecordFilename {
				panic(exception.NewIllegalInputException(fmt.Sprintf("The ASiC container already contains a file with name '%s'! "+
					"Addition of an evidence record of the same type is not allowed for ASiC with XAdES container.", evidenceRecordFilename)))
			}
		}
	}

	if enumerations.EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD == evidenceRecordType &&
		asic.ASiCUtilsEvidenceRecordERS != evidenceRecordFilename {
		panic(exception.NewIllegalInputException(fmt.Sprintf("RFC 4998 Evidence Record's filename '%s' is "+
			"not compliant to the ASiC with XAdES filename convention!", evidenceRecordFilename)))
	} else if enumerations.EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD == evidenceRecordType &&
		asic.ASiCUtilsEvidenceRecordXML != evidenceRecordFilename {
		panic(exception.NewIllegalInputException(fmt.Sprintf("RFC 6283 XML Evidence Record's filename '%s' is "+
			"not compliant to the ASiC with XAdES filename convention!", evidenceRecordFilename)))
	}
}
