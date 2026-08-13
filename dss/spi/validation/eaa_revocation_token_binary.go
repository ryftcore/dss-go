// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAARevocationTokenBinary.java (DSS 6.5.RC1).
package validation

import "github.com/utain/esig/dss/model"

// EAARevocationTokenBinary contains binaries of the EAA revocation token.
type EAARevocationTokenBinary struct {
	model.MultipleDigestIdentifier
}

// NewEAARevocationTokenBinary is the default constructor.
func NewEAARevocationTokenBinary(binaries []byte) *EAARevocationTokenBinary {
	return &EAARevocationTokenBinary{
		MultipleDigestIdentifier: model.NewMultipleDigestIdentifier("EAARevocationTokenBinary", "EAAR", binaries),
	}
}
