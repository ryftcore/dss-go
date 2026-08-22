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
	// KeyUsageBitDigitalSignature is digitalSignature.
	KeyUsageBitDigitalSignature KeyUsageBit = "DIGITAL_SIGNATURE"
	// KeyUsageBitNonRepudiation is nonRepudiation.
	KeyUsageBitNonRepudiation KeyUsageBit = "NON_REPUDIATION"
	// KeyUsageBitKeyEncipherment is keyEncipherment.
	KeyUsageBitKeyEncipherment KeyUsageBit = "KEY_ENCIPHERMENT"
	// KeyUsageBitDataEncipherment is dataEncipherment.
	KeyUsageBitDataEncipherment KeyUsageBit = "DATA_ENCIPHERMENT"
	// KeyUsageBitKeyAgreement is keyAgreement.
	KeyUsageBitKeyAgreement KeyUsageBit = "KEY_AGREEMENT"
	// KeyUsageBitKeyCertSign is keyCertSign.
	KeyUsageBitKeyCertSign KeyUsageBit = "KEY_CERT_SIGN"
	// KeyUsageBitCRLSign is crlSign.
	KeyUsageBitCRLSign KeyUsageBit = "CRL_SIGN"
	// KeyUsageBitEncipherOnly is encipherOnly.
	KeyUsageBitEncipherOnly KeyUsageBit = "ENCIPHER_ONLY"
	// KeyUsageBitDecipherOnly is decipherOnly.
	KeyUsageBitDecipherOnly KeyUsageBit = "DECIPHER_ONLY"
)

type keyUsageBitFields struct {
	value string
	index int
	bit   int
}

// keyUsageBitData holds the (value, index, bit) triple for each KeyUsageBit.
var keyUsageBitData = map[KeyUsageBit]keyUsageBitFields{
	KeyUsageBitDigitalSignature: {"digitalSignature", 0, 128},
	KeyUsageBitNonRepudiation:   {"nonRepudiation", 1, 64},
	KeyUsageBitKeyEncipherment:  {"keyEncipherment", 2, 32},
	KeyUsageBitDataEncipherment: {"dataEncipherment", 3, 16},
	KeyUsageBitKeyAgreement:     {"keyAgreement", 4, 8},
	KeyUsageBitKeyCertSign:      {"keyCertSign", 5, 4},
	KeyUsageBitCRLSign:          {"crlSign", 6, 2},
	KeyUsageBitEncipherOnly:     {"encipherOnly", 7, 1},
	KeyUsageBitDecipherOnly:     {"decipherOnly", 8, 32768},
}

// KeyUsageBitValues returns all KeyUsageBit constants in declaration order.
func KeyUsageBitValues() []KeyUsageBit {
	return []KeyUsageBit{
		KeyUsageBitDigitalSignature,
		KeyUsageBitNonRepudiation,
		KeyUsageBitKeyEncipherment,
		KeyUsageBitDataEncipherment,
		KeyUsageBitKeyAgreement,
		KeyUsageBitKeyCertSign,
		KeyUsageBitCRLSign,
		KeyUsageBitEncipherOnly,
		KeyUsageBitDecipherOnly,
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
