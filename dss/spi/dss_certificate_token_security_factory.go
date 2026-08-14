// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSCertificateTokenSecurityFactory.java (DSS 6.5.RC1).
//
// Upstream builds a CertificateToken from an InputStream, byte array, File or BouncyCastle
// X509CertificateHolder by retrying CertificateFactory.getInstance("X.509", securityProvider)
// across every configured JCA security provider (see dss_security_factory.go). This port has
// no provider registry to retry against (see DSSSecurityProviderInitSystemProviders), so
// BuildWithProvider below parses directly with crypto/x509 - reusing the exact same
// dssUtilsParseCertificate logic DSSUtils.loadCertificate(...) already exposes as
// DSSUtilsLoadCertificate/-FromBinary/-FromStream (see the DEVIATION note at the top of
// dss_utils.go, which flagged this file as the reconciliation point).
//
// DEVIATION: Go has no counterpart to BouncyCastle's X509CertificateHolder (a thin wrapper
// around a certificate's DER encoding); DSSCertificateTokenSecurityFactoryX509CertificateHolderInstance
// is therefore an alias of DSSCertificateTokenSecurityFactoryBinaryInstance, taking the raw DER
// bytes directly instead of a holder object.
package spi

import (
	"io"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// dssCertificateTokenSecurityFactoryClassName is CertificateFactory.class.getSimpleName().
const dssCertificateTokenSecurityFactoryClassName = "CertificateFactory"

// DSSCertificateTokenSecurityFactoryInputStreamInstance builds a CertificateToken from an
// io.Reader. Port of DSSCertificateTokenSecurityFactory.INPUT_STREAM_INSTANCE.
var DSSCertificateTokenSecurityFactoryInputStreamInstance = &DSSSecurityFactory[io.Reader, *model.CertificateToken]{
	FactoryClassName:  dssCertificateTokenSecurityFactoryClassName,
	ToString:          func(io.Reader) string { return "InputStream" },
	BuildWithProvider: DSSUtilsLoadCertificateFromStream,
}

// DSSCertificateTokenSecurityFactoryBinaryInstance builds a CertificateToken from the raw,
// DER-encoded certificate bytes. Port of DSSCertificateTokenSecurityFactory.BINARY_INSTANCE.
var DSSCertificateTokenSecurityFactoryBinaryInstance = &DSSSecurityFactory[[]byte, *model.CertificateToken]{
	FactoryClassName:  dssCertificateTokenSecurityFactoryClassName,
	ToString:          dssCertificateTokenSecurityFactoryBinaryToString,
	BuildWithProvider: DSSUtilsLoadCertificateFromBinary,
}

func dssCertificateTokenSecurityFactoryBinaryToString(input []byte) string {
	if input == nil {
		return ""
	}
	return utils.ToBase64(input)
}

// DSSCertificateTokenSecurityFactoryFileInstance builds a CertificateToken from a file path.
// Port of DSSCertificateTokenSecurityFactory.FILE_INSTANCE.
var DSSCertificateTokenSecurityFactoryFileInstance = &DSSSecurityFactory[string, *model.CertificateToken]{
	FactoryClassName:  dssCertificateTokenSecurityFactoryClassName,
	ToString:          func(input string) string { return input },
	BuildWithProvider: DSSUtilsLoadCertificate,
}

// DSSCertificateTokenSecurityFactoryX509CertificateHolderInstance builds a CertificateToken
// from raw certificate DER bytes, standing in for
// DSSCertificateTokenSecurityFactory.X509_CERTIFICATE_HOLDER_INSTANCE (see the file DEVIATION).
var DSSCertificateTokenSecurityFactoryX509CertificateHolderInstance = DSSCertificateTokenSecurityFactoryBinaryInstance
