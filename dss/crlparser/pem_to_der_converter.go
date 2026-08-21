// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/PemToDerConverter.java (DSS 6.5.RC1).
package crlparser

import (
	"encoding/pem"

	"github.com/ryftcore/dss-go/dss/model"
)

// PemToDerConverterConvert converts PEM encoded binaries (CRL, Cert) to their DER encoded
// equivalent. Port of the static PemToDerConverter.convert(byte[]).
//
// DEVIATION: upstream reads the PEM through BouncyCastle's PemReader and distinguishes two
// failure modes - a well-formed stream with no PEM object (message "Unable to read PEM
// Object") and an I/O failure while reading it (message "Unable to convert the CRL to DER",
// wrapping the cause). encoding/pem works on an in-memory byte slice and has no I/O failure
// mode of its own; pem.Decode returns a nil block for anything it cannot parse as a PEM
// object, which this port reports with the first message.
func PemToDerConverterConvert(pemEncoded []byte) ([]byte, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, model.NewDSSError("Unable to read PEM Object")
	}
	return block.Bytes, nil
}
