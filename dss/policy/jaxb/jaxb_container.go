// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// ContainerConstraints is the Go form of the generated JAXB class
// ContainerConstraints (complexType ContainerConstraints): ASiC specific
// constraints.
type ContainerConstraints struct {
	AcceptableContainerTypes      *MultiValuesConstraint `xml:"AcceptableContainerTypes,omitempty"`
	ZipCommentPresent             *LevelConstraint       `xml:"ZipCommentPresent,omitempty"`
	AcceptableZipComment          *MultiValuesConstraint `xml:"AcceptableZipComment,omitempty"`
	MimeTypeFilePresent           *LevelConstraint       `xml:"MimeTypeFilePresent,omitempty"`
	AcceptableMimeTypeFileContent *MultiValuesConstraint `xml:"AcceptableMimeTypeFileContent,omitempty"`
	ManifestFilePresent           *LevelConstraint       `xml:"ManifestFilePresent,omitempty"`
	SignedFilesPresent            *LevelConstraint       `xml:"SignedFilesPresent,omitempty"`
	FilenameAdherence             *LevelConstraint       `xml:"FilenameAdherence,omitempty"`
	AllFilesSigned                *LevelConstraint       `xml:"AllFilesSigned,omitempty"`
}

// PDFAConstraints is the Go form of the generated JAXB class PDFAConstraints
// (complexType PDFAConstraints): a group of constraints used for a PDF
// document validation against a PDF/A specification.
type PDFAConstraints struct {
	AcceptablePDFAProfiles *MultiValuesConstraint `xml:"AcceptablePDFAProfiles,omitempty"`
	PDFACompliant          *LevelConstraint       `xml:"PDFACompliant,omitempty"`
}
