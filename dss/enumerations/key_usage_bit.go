// Ported from dss-enumerations/.../KeyUsageBit.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// KeyUsageBit lists the KeyUsage bit values.
//
//	KeyUsage ::= BIT STRING {
//		digitalSignature (0),
//		nonRepudiation (1),
//		keyEncipherment (2),
//		dataEncipherment (3),
//		keyAgreement (4),
//		keyCertSign (5),
//		cRLSign (6),
//		encipherOnly (7),
//		decipherOnly (8)
//	}
type KeyUsageBit string

const (
	// KeyUsageBit_DIGITAL_SIGNATURE is digitalSignature.
	KeyUsageBit_DIGITAL_SIGNATURE KeyUsageBit = "DIGITAL_SIGNATURE"
	// KeyUsageBit_NON_REPUDIATION is nonRepudiation.
	KeyUsageBit_NON_REPUDIATION KeyUsageBit = "NON_REPUDIATION"
	// KeyUsageBit_KEY_ENCIPHERMENT is keyEncipherment.
	KeyUsageBit_KEY_ENCIPHERMENT KeyUsageBit = "KEY_ENCIPHERMENT"
	// KeyUsageBit_DATA_ENCIPHERMENT is dataEncipherment.
	KeyUsageBit_DATA_ENCIPHERMENT KeyUsageBit = "DATA_ENCIPHERMENT"
	// KeyUsageBit_KEY_AGREEMENT is keyAgreement.
	KeyUsageBit_KEY_AGREEMENT KeyUsageBit = "KEY_AGREEMENT"
	// KeyUsageBit_KEY_CERT_SIGN is keyCertSign.
	KeyUsageBit_KEY_CERT_SIGN KeyUsageBit = "KEY_CERT_SIGN"
	// KeyUsageBit_CRL_SIGN is crlSign.
	KeyUsageBit_CRL_SIGN KeyUsageBit = "CRL_SIGN"
	// KeyUsageBit_ENCIPHER_ONLY is encipherOnly.
	KeyUsageBit_ENCIPHER_ONLY KeyUsageBit = "ENCIPHER_ONLY"
	// KeyUsageBit_DECIPHER_ONLY is decipherOnly.
	KeyUsageBit_DECIPHER_ONLY KeyUsageBit = "DECIPHER_ONLY"
)

type keyUsageBitFields struct {
	value string
	index int
	bit   int
}

// keyUsageBitData holds the (value, index, bit) triple for each KeyUsageBit.
var keyUsageBitData = map[KeyUsageBit]keyUsageBitFields{
	KeyUsageBit_DIGITAL_SIGNATURE: {"digitalSignature", 0, 128},
	KeyUsageBit_NON_REPUDIATION:   {"nonRepudiation", 1, 64},
	KeyUsageBit_KEY_ENCIPHERMENT:  {"keyEncipherment", 2, 32},
	KeyUsageBit_DATA_ENCIPHERMENT: {"dataEncipherment", 3, 16},
	KeyUsageBit_KEY_AGREEMENT:     {"keyAgreement", 4, 8},
	KeyUsageBit_KEY_CERT_SIGN:     {"keyCertSign", 5, 4},
	KeyUsageBit_CRL_SIGN:          {"crlSign", 6, 2},
	KeyUsageBit_ENCIPHER_ONLY:     {"encipherOnly", 7, 1},
	KeyUsageBit_DECIPHER_ONLY:     {"decipherOnly", 8, 32768},
}

// KeyUsageBitValues returns all KeyUsageBit constants in declaration order.
func KeyUsageBitValues() []KeyUsageBit {
	return []KeyUsageBit{
		KeyUsageBit_DIGITAL_SIGNATURE,
		KeyUsageBit_NON_REPUDIATION,
		KeyUsageBit_KEY_ENCIPHERMENT,
		KeyUsageBit_DATA_ENCIPHERMENT,
		KeyUsageBit_KEY_AGREEMENT,
		KeyUsageBit_KEY_CERT_SIGN,
		KeyUsageBit_CRL_SIGN,
		KeyUsageBit_ENCIPHER_ONLY,
		KeyUsageBit_DECIPHER_ONLY,
	}
}

// KeyUsageBitValueOf returns the KeyUsageBit matching the given Java enum name.
func KeyUsageBitValueOf(name string) (KeyUsageBit, error) {
	for _, v := range KeyUsageBitValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant KeyUsageBit.%s", name)
}

// Value returns the key usage name.
func (k KeyUsageBit) Value() string {
	return keyUsageBitData[k].value
}

// Index returns the key usage index.
func (k KeyUsageBit) Index() int {
	return keyUsageBitData[k].index
}

// Bit returns the key usage bit value.
func (k KeyUsageBit) Bit() int {
	return keyUsageBitData[k].bit
}
