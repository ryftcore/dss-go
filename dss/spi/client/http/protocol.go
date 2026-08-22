// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/Protocol.java (DSS 6.5.RC1).
package http

import "strings"

// Protocol lists all network protocols that can be used during the signature
// creation or validation: OCSP, CRL, AIA, TSL...
type Protocol string

const (
	// ProtocolFile is the file protocol.
	ProtocolFile Protocol = "file"
	// ProtocolHTTP is the http protocol.
	ProtocolHTTP Protocol = "http"
	// ProtocolHTTPS is the https protocol.
	ProtocolHTTPS Protocol = "https"
	// ProtocolLDAP is the ldap protocol.
	ProtocolLDAP Protocol = "ldap"
	// ProtocolFTP is the ftp protocol.
	ProtocolFTP Protocol = "ftp"
)

// Name gets the name of the protocol.
func (p Protocol) Name() string {
	return string(p)
}

// ProtocolIsHttps indicates if the given string represents the HTTPS protocol.
func ProtocolIsHttps(name string) bool {
	return strings.EqualFold(string(ProtocolHTTPS), name)
}

// ProtocolIsHttp indicates if the given string represents the HTTP protocol.
func ProtocolIsHttp(name string) bool {
	return strings.EqualFold(string(ProtocolHTTP), name)
}

// ProtocolIsFileURL indicates if the given URL uses the FILE protocol.
func ProtocolIsFileURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(ProtocolFile))
}

// ProtocolIsHttpURL indicates if the given URL uses the HTTP protocol.
func ProtocolIsHttpURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(ProtocolHTTP))
}

// ProtocolIsFtpURL indicates if the given URL uses the FTP protocol.
func ProtocolIsFtpURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(ProtocolFTP))
}

// ProtocolIsLdapURL indicates if the given URL uses the LDAP protocol.
func ProtocolIsLdapURL(urlString string) bool {
	return strings.HasPrefix(urlString, string(ProtocolLDAP))
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
	case ProtocolHTTP.IsTheSame(urlString):
		return ProtocolHTTP
	case ProtocolHTTPS.IsTheSame(urlString):
		return ProtocolHTTPS
	case ProtocolLDAP.IsTheSame(urlString):
		return ProtocolLDAP
	case ProtocolFTP.IsTheSame(urlString):
		return ProtocolFTP
	case ProtocolFile.IsTheSame(urlString):
		return ProtocolFile
	}
	return ""
}
