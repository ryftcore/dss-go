// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/merge/AbstractASiCWithXAdESContainerMerger.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.merge lands in this same Go
// package (dss/asic/xades) per S7_BRIEF.md's package layout table.
package xades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
)

// AbstractASiCWithXAdESContainerMerger contains common code for ASiC with XAdES container merger
// classes.
type AbstractASiCWithXAdESContainerMerger struct {
	*asic.DefaultContainerMerger

	// asicFilenameFactory defines rules for filename creation for new ZIP entries (e.g.
	// signature files, etc.). Java declares the field protected; exported here since Go
	// subclasses live in the same package but a caller-visible setter is needed too.
	asicFilenameFactory ASiCWithXAdESFilenameFactory
}

// NewAbstractASiCWithXAdESContainerMergerBase is the empty constructor. Port of the
// package-private empty constructor. The subclass constructor must follow it with
// InitDefaultContainerMerger and then InitFromDocuments/InitFromASiCContents.
func NewAbstractASiCWithXAdESContainerMergerBase() AbstractASiCWithXAdESContainerMerger {
	return AbstractASiCWithXAdESContainerMerger{
		DefaultContainerMerger: asic.NewDefaultContainerMergerBase(),
		asicFilenameFactory:    NewDefaultASiCWithXAdESFilenameFactory(),
	}
}

// SetAsicFilenameFactory sets the ASiCWithXAdESFilenameFactory defining a set of rules for
// naming of newly created ZIP entries, such as signature files.
//
// Panics with the Java message when asicFilenameFactory is nil (Objects.requireNonNull).
func (m *AbstractASiCWithXAdESContainerMerger) SetAsicFilenameFactory(asicFilenameFactory ASiCWithXAdESFilenameFactory) {
	if asicFilenameFactory == nil {
		panic("ASiCWithXAdESFilenameFactory cannot be null!")
	}
	m.asicFilenameFactory = asicFilenameFactory
}

// IsSupportedDocument ports the @Override protected isSupported(DSSDocument).
func (m *AbstractASiCWithXAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	return NewASiCWithXAdESFormatDetector().IsSupportedZip(container)
}

// IsSupportedContent ports the @Override protected isSupported(ASiCContent).
func (m *AbstractASiCWithXAdESContainerMerger) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	return NewASiCWithXAdESFormatDetector().IsSupportedZipContent(asicContent)
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(DSSDocument).
func (m *AbstractASiCWithXAdESContainerMerger) GetContainerExtractor(container model.DSSDocument) *asic.DefaultASiCContainerExtractor {
	return &NewASiCWithXAdESContainerExtractor(container).DefaultASiCContainerExtractor
}
