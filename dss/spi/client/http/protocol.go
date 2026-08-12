// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/Protocol.java (DSS 6.5.RC1).
package http

import "strings"

// Protocol lists all network protocols that can be used during the signature
// creation or validation: OCSP, CRL, AIA, TSL...
type Protocol string

const (
	// Protocol_FILE is the file protocol.
	Protocol_FILE Protocol = "file"
	// Protocol_HTTP is the http protocol.
	Protocol_HTTP Protocol = "http"
	// Protocol_HTTPS is the https protocol.
	Protocol_HTTPS Protocol = "https"
	// Protocol_LDAP is the ldap protocol.
	Protocol_LDAP Protocol = "ldap"
	// Protocol_FTP is the ftp protocol.
	Protocol_FTP Protocol = "ftp"
)

// Name gets the name of the protocol.
func (p Protocol) Name() string {
	return string(p)
}

// ProtocolIsHttps indicates if the given string represents the HTTPS protocol.
func ProtocolIsHttps(name string) bool {
	return strings.EqualFold(string(Protocol_HTTPS), name)
}

// ProtocolIsHttp indicates if the given string represents the HTTP protocol.
func ProtocolIsHttp(name string) bool {
	return strings.EqualFold(string(Protocol_HTTP), name)
}

// ProtocolIsFileURL indicates if the given URL uses the FILE protocol.
func ProtocolIsFileURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(Protocol_FILE))
}

// ProtocolIsHttpURL indicates if the given URL uses the HTTP protocol.
func ProtocolIsHttpURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(Protocol_HTTP))
}

// ProtocolIsFtpURL indicates if the given URL uses the FTP protocol.
func ProtocolIsFtpURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(Protocol_FTP))
}

// ProtocolIsLdapURL indicates if the given URL uses the LDAP protocol.
func ProtocolIsLdapURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(Protocol_LDAP))
}

// IsTheSame indicates if the given URL uses the current protocol.
func (p Protocol) IsTheSame(urlString string) bool {
	return strings.HasPrefix(urlString, string(p))
}

// ProtocolGetFrom tries to retrieve the protocol indicated by the given URL
// string. Returns "" (the zero Protocol) when none matches, mirroring the
// Java method's null return.
func ProtocolGetFrom(urlString string) Protocol {
	switch {
	case Protocol_HTTP.IsTheSame(urlString):
		return Protocol_HTTP
	case Protocol_HTTPS.IsTheSame(urlString):
		return Protocol_HTTPS
	case Protocol_LDAP.IsTheSame(urlString):
		return Protocol_LDAP
	case Protocol_FTP.IsTheSame(urlString):
		return Protocol_FTP
	case Protocol_FILE.IsTheSame(urlString):
		return Protocol_FILE
	}
	return ""
}
