// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSP7CCertificatesSecurityFactory.java (DSS 6.5.RC1).
//
// Same DEVIATION as dss_certificate_token_security_factory.go: no JCA provider registry to
// retry against, so BuildWithProvider delegates straight to the (already ported)
// DSSUtilsLoadCertificateFromP7c family, which reimplements CertificateFactory#generateCertificates
// by hand-parsing the PKCS#7 SignedData degenerate case (see dss_utils.go).
package spi

import (
	"io"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// dssP7CCertificatesSecurityFactoryClassName is CertificateFactory.class.getSimpleName().
const dssP7CCertificatesSecurityFactoryClassName = "CertificateFactory"

// DSSP7CCertificatesSecurityFactoryInputStreamInstance builds a p7c certificate chain from an
// io.Reader. Port of DSSP7CCertificatesSecurityFactory.INPUT_STREAM_INSTANCE.
var DSSP7CCertificatesSecurityFactoryInputStreamInstance = &DSSSecurityFactory[io.Reader, []*model.CertificateToken]{
	FactoryClassName:  dssP7CCertificatesSecurityFactoryClassName,
	ToString:          func(io.Reader) string { return "InputStream" },
	BuildWithProvider: DSSUtilsLoadCertificateFromP7cStream,
}

// DSSP7CCertificatesSecurityFactoryBinaryInstance builds a p7c certificate chain from the raw
// p7c bytes. Port of DSSP7CCertificatesSecurityFactory.BINARY_INSTANCE.
var DSSP7CCertificatesSecurityFactoryBinaryInstance = &DSSSecurityFactory[[]byte, []*model.CertificateToken]{
	FactoryClassName:  dssP7CCertificatesSecurityFactoryClassName,
	ToString:          dssP7CCertificatesSecurityFactoryBinaryToString,
	BuildWithProvider: DSSUtilsLoadCertificateFromP7cBinary,
}

func dssP7CCertificatesSecurityFactoryBinaryToString(input []byte) string {
	if input == nil {
		return ""
	}
	return utils.ToBase64(input)
}

// DSSP7CCertificatesSecurityFactoryFileInstance builds a p7c certificate chain from a file
// path. Port of DSSP7CCertificatesSecurityFactory.FILE_INSTANCE.
var DSSP7CCertificatesSecurityFactoryFileInstance = &DSSSecurityFactory[string, []*model.CertificateToken]{
	FactoryClassName:  dssP7CCertificatesSecurityFactoryClassName,
	ToString:          func(input string) string { return input },
	BuildWithProvider: DSSUtilsLoadCertificateFromP7c,
}
