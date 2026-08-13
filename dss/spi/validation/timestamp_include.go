// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampInclude.java (DSS 6.5.RC1).
package validation

// TimestampInclude represents the XAdES Include tag in case of an IndividualDataObjectsTimeStamp.
type TimestampInclude struct {
	// uri is the reference URI.
	uri string
	// referencedData is the referencedData attribute, which shall be present in each and every
	// Include element and set to "true".
	referencedData bool
}

// NewTimestampInclude builds an empty include. Port of the empty constructor.
func NewTimestampInclude() *TimestampInclude {
	return &TimestampInclude{}
}

// NewTimestampIncludeWithURI builds an include for the given reference URI.
// Port of the TimestampInclude(String, boolean) constructor.
func NewTimestampIncludeWithURI(uri string, referencedData bool) *TimestampInclude {
	return &TimestampInclude{uri: uri, referencedData: referencedData}
}

// URI gets the reference URI. Port of getURI().
func (i *TimestampInclude) URI() string {
	return i.uri
}

// SetURI sets the reference URI. Port of setURI(String).
func (i *TimestampInclude) SetURI(uri string) {
	i.uri = uri
}

// IsReferencedData reports whether the reference is timestamped. Port of isReferencedData().
func (i *TimestampInclude) IsReferencedData() bool {
	return i.referencedData
}

// SetReferencedData sets whether the reference is timestamped.
// Port of setReferencedData(boolean).
func (i *TimestampInclude) SetReferencedData(referencedData bool) {
	i.referencedData = referencedData
}
