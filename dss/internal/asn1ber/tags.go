package asn1ber

// Universal ASN.1 tag numbers, i.e. the org.bouncycastle.asn1.BERTags constants.
const (
	TagBoolean         = 0x01
	TagInteger         = 0x02
	TagBitString       = 0x03
	TagOctetString     = 0x04
	TagNull            = 0x05
	TagOID             = 0x06
	TagObjectDesc      = 0x07
	TagUTF8String      = 0x0C
	TagSequence        = 0x10
	TagSet             = 0x11
	TagNumericString   = 0x12
	TagPrintableString = 0x13
	TagT61String       = 0x14
	TagVideotexString  = 0x15
	TagIA5String       = 0x16
	TagUTCTime         = 0x17
	TagGeneralizedTime = 0x18
	TagGraphicString   = 0x19
	TagVisibleString   = 0x1A
	TagGeneralString   = 0x1B
	TagUniversalString = 0x1C
	TagBMPString       = 0x1E
)

// ASN.1 identifier octet bit masks and class bits.
const (
	// ClassMask selects the class bits of an identifier octet.
	ClassMask = 0xC0
	// Constructed is the identifier octet bit marking a constructed element.
	Constructed = 0x20
	// TagMask selects the low tag number bits of an identifier octet.
	TagMask = 0x1F

	// ClassUniversal is the universal class.
	ClassUniversal = 0x00
	// ClassApplication is the application class.
	ClassApplication = 0x40
	// ClassContextSpecific is the context-specific class.
	ClassContextSpecific = 0x80
	// ClassPrivate is the private class.
	ClassPrivate = 0xC0
)

// DERNull is the DER encoding of ASN.1 NULL, i.e. BouncyCastle's DERNull.INSTANCE.
var DERNull = []byte{TagNull, 0x00}
