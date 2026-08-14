// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSPublicKeySecurityFactory.java (DSS 6.5.RC1).
//
// Upstream's SUBJECT_PUBLIC_KEY_INFO_INSTANCE resolves the SubjectPublicKeyInfo's algorithm
// OID to an EncryptionAlgorithm (falling back to the signature-algorithm OID table when the
// encryption-algorithm one misses) and calls KeyFactory.generatePublic(...) with it. In Go,
// model.NewPublicKey(subjectPublicKeyInfo) already performs the equivalent resolve-and-parse
// in one step via crypto/x509.ParsePKIXPublicKey (RSA/EC/Ed25519 OIDs are all recognised
// directly by the stdlib parser, so the EncryptionAlgorithm OID fallback chain upstream needs
// has no Go counterpart to reimplement), and - like PublicKey itself - keeps the original
// SubjectPublicKeyInfo bytes rather than re-encoding, per PORTING.md's crypto identity rule.
//
// DEVIATION: no JCA provider registry to retry against, same as the sibling factories in this
// package (see dss_security_factory.go).
package spi

import (
	"io"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// dssPublicKeySecurityFactoryClassName is KeyFactory.class.getSimpleName().
const dssPublicKeySecurityFactoryClassName = "KeyFactory"

// DSSPublicKeySecurityFactoryInputStreamInstance builds a PublicKey from an io.Reader. Port of
// DSSPublicKeySecurityFactory.INPUT_STREAM_INSTANCE.
var DSSPublicKeySecurityFactoryInputStreamInstance = &DSSSecurityFactory[io.Reader, *model.PublicKey]{
	FactoryClassName:  dssPublicKeySecurityFactoryClassName,
	ToString:          func(io.Reader) string { return "InputStream" },
	BuildWithProvider: dssPublicKeySecurityFactoryBuildFromStream,
}

func dssPublicKeySecurityFactoryBuildFromStream(r io.Reader) (*model.PublicKey, error) {
	if r == nil {
		panic("Input cannot be null")
	}
	data, err := utils.ToByteArray(r)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to read InputStream", err)
	}
	return model.NewPublicKey(data)
}

// DSSPublicKeySecurityFactoryBinaryInstance builds a PublicKey from the raw, DER-encoded
// SubjectPublicKeyInfo bytes. Port of DSSPublicKeySecurityFactory.BINARY_INSTANCE.
var DSSPublicKeySecurityFactoryBinaryInstance = &DSSSecurityFactory[[]byte, *model.PublicKey]{
	FactoryClassName:  dssPublicKeySecurityFactoryClassName,
	ToString:          dssPublicKeySecurityFactoryBinaryToString,
	BuildWithProvider: model.NewPublicKey,
}

func dssPublicKeySecurityFactoryBinaryToString(input []byte) string {
	if input == nil {
		return ""
	}
	return utils.ToBase64(input)
}

// DSSPublicKeySecurityFactorySubjectPublicKeyInfoInstance builds a PublicKey from the raw
// SubjectPublicKeyInfo DER bytes, standing in for
// DSSPublicKeySecurityFactory.SUBJECT_PUBLIC_KEY_INFO_INSTANCE - Go has no separate
// SubjectPublicKeyInfo type distinct from its DER encoding, so this is an alias of
// DSSPublicKeySecurityFactoryBinaryInstance (both ultimately call model.NewPublicKey).
var DSSPublicKeySecurityFactorySubjectPublicKeyInfoInstance = DSSPublicKeySecurityFactoryBinaryInstance
