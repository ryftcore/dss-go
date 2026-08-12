// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/DataLoader.java (DSS 6.5.RC1).
package http

// DataLoader is a component that allows retrieving data using any protocol:
// HTTP, HTTPS, FTP, LDAP.
type DataLoader interface {
	// Get executes an HTTP GET operation. Returns the obtained data, or nil.
	Get(url string) []byte

	// GetFromURLs executes an HTTP GET operation. Used when many URLs are
	// available to access the same resource: the operation stops after the
	// first successful download. Ports the get(List<String>) overload.
	GetFromURLs(urlStrings []string) (*DataAndURL, error)

	// Post executes an HTTP POST operation.
	Post(url string, content []byte) []byte

	// SetContentType allows setting the content type, e.g. Content-Type
	// "application/ocsp-request".
	SetContentType(contentType string)
}

// DataAndURL is an internal class used to model the couple data and url used
// to obtain this data.
type DataAndURL struct {
	// urlString is the URL used to obtain data.
	urlString string

	// data is the obtained data.
	data []byte
}

// NewDataAndURL is the default constructor.
func NewDataAndURL(urlString string, data []byte) *DataAndURL {
	return &DataAndURL{urlString: urlString, data: data}
}

// URLString gets the URL string used to download the data.
func (d *DataAndURL) URLString() string {
	return d.urlString
}

// Data gets the downloaded data.
func (d *DataAndURL) Data() []byte {
	return d.data
}
