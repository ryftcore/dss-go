// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCParameters.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart), as is hashCode() (nothing in the
// ported tree keys a hash container on these objects; equals() survives as Equals).
package asic

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// Parameters regroups the signature parameters related to ASiC form.
//
// Java subclasses (ASiCContainerEvidenceRecordParameters here, and the format-specific
// ASiCWith{CAdES,XAdES}SignatureParameters in the dss-asic-cades/dss-asic-xades ports) embed this
// struct by value rather than extending it.
type Parameters struct {
	// zipComment indicates if the ZIP comment should be used to store the signed content
	// mime-type.
	zipComment bool

	// mimeType indicates the mime-type to be set within the mimetype file. If empty (Java null)
	// the stored mime-type is that of the signed content.
	mimeType string

	// containerType is the form of the container -S or -E.
	containerType enumerations.ASiCContainerType
}

// NewASiCParameters instantiates an object with null values. Port of the default constructor.
func NewASiCParameters() *Parameters {
	return &Parameters{}
}

// IsZipComment indicates if the ZIP comment must include the mime-type. Port of isZipComment().
func (p *Parameters) IsZipComment() bool {
	return p.zipComment
}

// SetZipComment sets if the zip comment will contain the mime type. Port of
// setZipComment(boolean).
func (p *Parameters) SetZipComment(zipComment bool) {
	p.zipComment = zipComment
}

// MimeType gets the mimetype. Port of getMimeType().
func (p *Parameters) MimeType() string {
	return p.mimeType
}

// SetMimeType sets the mime-type within the mimetype file. Port of setMimeType(String).
func (p *Parameters) SetMimeType(mimeType string) {
	p.mimeType = mimeType
}

// ContainerType returns the expected type of the ASiC container. Port of getContainerType().
func (p *Parameters) ContainerType() enumerations.ASiCContainerType {
	return p.containerType
}

// SetContainerType sets the expected container type. Port of
// setContainerType(ASiCContainerType).
func (p *Parameters) SetContainerType(containerType enumerations.ASiCContainerType) {
	p.containerType = containerType
}

// String ports toString(). Java prints the null mimeType as the four characters 'null' between
// the single quotes and the null containerType as null; the Go zero values ("" and "") render as
// empty instead - a cosmetic difference in a debug-only string.
func (p *Parameters) String() string {
	return fmt.Sprintf("ASiCParameters [zipComment=%t, mimeType='%s', containerType=%s]",
		p.zipComment, p.mimeType, string(p.containerType))
}

// Equals ports equals(Object).
func (p *Parameters) Equals(other *Parameters) bool {
	if p == other {
		return true
	}
	if p == nil || other == nil {
		return false
	}
	return p.zipComment == other.zipComment &&
		p.mimeType == other.mimeType &&
		p.containerType == other.containerType
}
