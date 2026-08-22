// Ported from dss-enumerations/.../EllipticCurve.java (DSS 6.5.RC1).
//
// NOTE: Java's ECParameterSpec (java.security.spec) has no direct Go stdlib
// counterpart carrying named prime, A/B curve coefficients, generator point
// and order as arbitrary-precision integers; EllipticCurveParameter below
// is a minimal port of that value object (field prime P, coefficients A/B,
// generator Gx/Gy, order N) sufficient to reproduce Java's getParameter,
// forParameter and equalCurves behavior. Values for these parameters are
// taken from the same jose4j EllipticCurves.java source cited by the Java
// Javadoc; the values are also present in FIPS PUB 186-3.
package enumerations

import "math/big"

// EllipticCurve represents an elliptic curve.
type EllipticCurve string

const (
	// EllipticCurveP256 is the P-256 curve.
	EllipticCurveP256 EllipticCurve = "P_256"
	// EllipticCurveP384 is the P-384 curve.
	EllipticCurveP384 EllipticCurve = "P_384"
	// EllipticCurveP521 is the P-512 curve.
	EllipticCurveP521 EllipticCurve = "P_521"
	// EllipticCurveX25519 is X25519.
	EllipticCurveX25519 EllipticCurve = "X25519"
	// EllipticCurveX448 is X448.
	EllipticCurveX448 EllipticCurve = "X448"
	// EllipticCurveED25519 is EdDSA 25519.
	EllipticCurveED25519 EllipticCurve = "ED25519"
	// EllipticCurveED448 is EdDSA 448.
	EllipticCurveED448 EllipticCurve = "ED448"
	// EllipticCurveSECP256K1 is the SECP-256k1 curve.
	EllipticCurveSECP256K1 EllipticCurve = "SECP_256K1"
	// EllipticCurveBrainpoolP256R1 is the Brainpool P-256 R1 curve.
	EllipticCurveBrainpoolP256R1 EllipticCurve = "BRAINPOOL_P256_R1"
	// EllipticCurveBrainpoolP320R1 is the Brainpool P-320 R1 curve.
	EllipticCurveBrainpoolP320R1 EllipticCurve = "BRAINPOOL_P320_R1"
	// EllipticCurveBrainpoolP384R1 is the Brainpool P-384 R1 curve.
	EllipticCurveBrainpoolP384R1 EllipticCurve = "BRAINPOOL_P384_R1"
	// EllipticCurveBrainpoolP512R1 is the Brainpool P-512 R1 curve.
	EllipticCurveBrainpoolP512R1 EllipticCurve = "BRAINPOOL_P512_R1"
)

// EllipticCurveParameter is a minimal port of java.security.spec.ECParameterSpec:
// the field prime P, curve coefficients A and B, generator point (Gx, Gy),
// and the order N. The cofactor is always 1 for the curves in this file
// (matching Java's COFACTOR constant) and is not represented.
type EllipticCurveParameter struct {
	P, A, B, Gx, Gy, N *big.Int
}

func bigDec(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("elliptic_curve.go: invalid decimal literal: " + s)
	}
	return n
}

func bigHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("elliptic_curve.go: invalid hex literal: " + s)
	}
	return n
}

var ellipticCurveP256 = EllipticCurveParameter{
	P:  bigDec("115792089210356248762697446949407573530086143415290314195533631308867097853951"),
	A:  bigDec("115792089210356248762697446949407573530086143415290314195533631308867097853948"),
	B:  bigDec("41058363725152142129326129780047268409114441015993725554835256314039467401291"),
	Gx: bigDec("48439561293906451759052585252797914202762949526041747995844080717082404635286"),
	Gy: bigDec("36134250956749795798585127919587881956611106672985015071877198253568414405109"),
	N:  bigDec("115792089210356248762697446949407573529996955224135760342422259061068512044369"),
}

var ellipticCurveP384 = EllipticCurveParameter{
	P:  bigDec("39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319"),
	A:  bigDec("39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112316"),
	B:  bigDec("27580193559959705877849011840389048093056905856361568521428707301988689241309860865136260764883745107765439761230575"),
	Gx: bigDec("26247035095799689268623156744566981891852923491109213387815615900925518854738050089022388053975719786650872476732087"),
	Gy: bigDec("8325710961489029985546751289520108179287853048861315594709205902480503199884419224438643760392947333078086511627871"),
	N:  bigDec("39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643"),
}

var ellipticCurveP521 = EllipticCurveParameter{
	P:  bigDec("6864797660130609714981900799081393217269435300143305409394463459185543183397656052122559640661454554977296311391480858037121987999716643812574028291115057151"),
	A:  bigDec("6864797660130609714981900799081393217269435300143305409394463459185543183397656052122559640661454554977296311391480858037121987999716643812574028291115057148"),
	B:  bigDec("1093849038073734274511112390766805569936207598951683748994586394495953116150735016013708737573759623248592132296706313309438452531591012912142327488478985984"),
	Gx: bigDec("2661740802050217063228768716723360960729859168756973147706671368418802944996427808491545080627771902352094241225065558662157113545570916814161637315895999846"),
	Gy: bigDec("3757180025770020463545507224491183603594455134769762486694567779615544477440556316691234405012945539562144444537289428522585666729196580810124344277578376784"),
	N:  bigDec("6864797660130609714981900799081393217269435300143305409394463459185543183397655394245057746333217197532963996371363321113864768612440380340372808892707005449"),
}

var ellipticCurveSecp256k1 = EllipticCurveParameter{
	P:  bigDec("115792089237316195423570985008687907853269984665640564039457584007908834671663"),
	A:  bigDec("0"),
	B:  bigDec("7"),
	Gx: bigDec("55066263022277343669578718895168534326250603453777594175500187360389116729240"),
	Gy: bigDec("32670510020758816978083085130507043184471273380659243275938904335757337482424"),
	N:  bigDec("115792089237316195423570985008687907852837564279074904382605163141518161494337"),
}

var ellipticCurveBP256 = EllipticCurveParameter{
	P:  bigHex("A9FB57DBA1EEA9BC3E660A909D838D726E3BF623D52620282013481D1F6E5377"),
	A:  bigHex("7D5A0975FC2C3057EEF67530417AFFE7FB8055C126DC5C6CE94A4B44F330B5D9"),
	B:  bigHex("26DC5C6CE94A4B44F330B5D9BBD77CBF958416295CF7E1CE6BCCDC18FF8C07B6"),
	Gx: bigHex("8BD2AEB9CB7E57CB2C4B482FFC81B7AFB9DE27E1E3BD23C23A4453BD9ACE3262"),
	Gy: bigHex("547EF835C3DAC4FD97F8461A14611DC9C27745132DED8E545C1D54C72F046997"),
	N:  bigHex("A9FB57DBA1EEA9BC3E660A909D838D718C397AA3B561A6F7901E0E82974856A7"),
}

var ellipticCurveBP320 = EllipticCurveParameter{
	P:  bigHex("D35E472036BC4FB7E13C785ED201E065F98FCFA6F6F40DEF4F92B9EC7893EC28FCD412B1F1B32E27"),
	A:  bigHex("3EE30B568FBAB0F883CCEBD46D3F3BB8A2A73513F5EB79DA66190EB085FFA9F492F375A97D860EB4"),
	B:  bigHex("520883949DFDBC42D3AD198640688A6FE13F41349554B49ACC31DCCD884539816F5EB4AC8FB1F1A6"),
	Gx: bigHex("43BD7E9AFB53D8B85289BCC48EE5BFE6F20137D10A087EB6E7871E2A10A599C710AF8D0D39E20611"),
	Gy: bigHex("14FDD05545EC1CC8AB4093247F77275E0743FFED117182EAA9C77877AAAC6AC7D35245D1692E8EE1"),
	N:  bigHex("D35E472036BC4FB7E13C785ED201E065F98FCFA5B68F12A32D482EC7EE8658E98691555B44C59311"),
}

var ellipticCurveBP384 = EllipticCurveParameter{
	P:  bigHex("8CB91E82A3386D280F5D6F7E50E641DF152F7109ED5456B412B1DA197FB71123ACD3A729901D1A71874700133107EC53"),
	A:  bigHex("7BC382C63D8C150C3C72080ACE05AFA0C2BEA28E4FB22787139165EFBA91F90F8AA5814A503AD4EB04A8C7DD22CE2826"),
	B:  bigHex("04A8C7DD22CE28268B39B55416F0447C2FB77DE107DCD2A62E880EA53EEB62D57CB4390295DBC9943AB78696FA504C11"),
	Gx: bigHex("1D1C64F068CF45FFA2A63A81B7C13F6B8847A3E77EF14FE3DB7FCAFE0CBD10E8E826E03436D646AAEF87B2E247D4AF1E"),
	Gy: bigHex("8ABE1D7520F9C2A45CB1EB8E95CFD55262B70B29FEEC5864E19C054FF99129280E4646217791811142820341263C5315"),
	N:  bigHex("8CB91E82A3386D280F5D6F7E50E641DF152F7109ED5456B31F166E6CAC0425A7CF3AB6AF6B7FC3103B883202E9046565"),
}

var ellipticCurveBP512 = EllipticCurveParameter{
	P:  bigHex("AADD9DB8DBE9C48B3FD4E6AE33C9FC07CB308DB3B3C9D20ED6639CCA703308717D4D9B009BC66842AECDA12AE6A380E62881FF2F2D82C68528AA6056583A48F3"),
	A:  bigHex("7830A3318B603B89E2327145AC234CC594CBDD8D3DF91610A83441CAEA9863BC2DED5D5AA8253AA10A2EF1C98B9AC8B57F1117A72BF2C7B9E7C1AC4D77FC94CA"),
	B:  bigHex("3DF91610A83441CAEA9863BC2DED5D5AA8253AA10A2EF1C98B9AC8B57F1117A72BF2C7B9E7C1AC4D77FC94CADC083E67984050B75EBAE5DD2809BD638016F723"),
	Gx: bigHex("81AEE4BDD82ED9645A21322E9C4C6A9385ED9F70B5D916C1B43B62EEF4D0098EFF3B1F78E2D0D48D50D1687B93B97D5F7C6D5047406A5E688B352209BCB9F822"),
	Gy: bigHex("7DDE385D566332ECC0EABFA9CF7822FDF209F70024A57B1AA000C55B881F8111B2DCDE494A5F485E5BCA4BD88A2763AED1CA2B2FA8F0540678CD1E0F3AD80892"),
	N:  bigHex("AADD9DB8DBE9C48B3FD4E6AE33C9FC07CB308DB3B3C9D20ED6639CCA70330870553E5C414CA92619418661197FAC10471DB1D381085DDADDB58796829CA90069"),
}

// ellipticCurveParameters mirrors Java's ELLIPTIC_CURVE_PARAMETERS map.
// Only curves with an ECParameterSpec registered in Java are present here
// (X25519, X448, ED25519, ED448 have none).
var ellipticCurveParameters = map[EllipticCurve]EllipticCurveParameter{
	EllipticCurveP256:            ellipticCurveP256,
	EllipticCurveP384:            ellipticCurveP384,
	EllipticCurveP521:            ellipticCurveP521,
	EllipticCurveSECP256K1:       ellipticCurveSecp256k1,
	EllipticCurveBrainpoolP256R1: ellipticCurveBP256,
	EllipticCurveBrainpoolP320R1: ellipticCurveBP320,
	EllipticCurveBrainpoolP384R1: ellipticCurveBP384,
	EllipticCurveBrainpoolP512R1: ellipticCurveBP512,
}

var ellipticCurveLabel = map[EllipticCurve]string{
	EllipticCurveP256:            "P-256",
	EllipticCurveP384:            "P-384",
	EllipticCurveP521:            "P-521",
	EllipticCurveX25519:          "X25519",
	EllipticCurveX448:            "X448",
	EllipticCurveED25519:         "Ed25519",
	EllipticCurveED448:           "Ed448",
	EllipticCurveSECP256K1:       "secp256k1",
	EllipticCurveBrainpoolP256R1: "brainpoolP256r1",
	EllipticCurveBrainpoolP320R1: "brainpoolP320r1",
	EllipticCurveBrainpoolP384R1: "brainpoolP384r1",
	EllipticCurveBrainpoolP512R1: "brainpoolP512r1",
}

var ellipticCurveSize = map[EllipticCurve]int{
	EllipticCurveP256:            32,
	EllipticCurveP384:            48,
	EllipticCurveP521:            66,
	EllipticCurveX25519:          32,
	EllipticCurveX448:            56,
	EllipticCurveED25519:         32,
	EllipticCurveED448:           57,
	EllipticCurveSECP256K1:       32,
	EllipticCurveBrainpoolP256R1: 32,
	EllipticCurveBrainpoolP320R1: 40,
	EllipticCurveBrainpoolP384R1: 48,
	EllipticCurveBrainpoolP512R1: 64,
}

// ellipticCurveCOSEValues mirrors Java's ELLIPTIC_CURVE_COSE_VALUES map;
// values from the IANA COSE registry.
var ellipticCurveCOSEValues = map[EllipticCurve]int64{
	EllipticCurveP256:            1,
	EllipticCurveP384:            2,
	EllipticCurveP521:            3,
	EllipticCurveX25519:          4,
	EllipticCurveX448:            5,
	EllipticCurveED25519:         6,
	EllipticCurveED448:           7,
	EllipticCurveSECP256K1:       8,
	EllipticCurveBrainpoolP256R1: 256,
	EllipticCurveBrainpoolP320R1: 257,
	EllipticCurveBrainpoolP384R1: 258,
	EllipticCurveBrainpoolP512R1: 259,
}

// ellipticCurveOIDs mirrors Java's ELLIPTIC_CURVE_OIDS map. Note P_256 and
// SECP_256K1 share the same OID value in the Java source (verbatim copy).
var ellipticCurveOIDs = map[EllipticCurve]string{
	EllipticCurveP256:            "1.2.840.10045.3.1.7",
	EllipticCurveP384:            "1.3.132.0.34",
	EllipticCurveP521:            "1.3.132.0.35",
	EllipticCurveX25519:          "1.3.101.110",
	EllipticCurveX448:            "1.3.101.111",
	EllipticCurveED25519:         "1.3.101.112",
	EllipticCurveED448:           "1.3.101.113",
	EllipticCurveSECP256K1:       "1.2.840.10045.3.1.7",
	EllipticCurveBrainpoolP256R1: "1.3.36.3.3.2.8.1.1.7",
	EllipticCurveBrainpoolP320R1: "1.3.36.3.3.2.8.1.1.9",
	EllipticCurveBrainpoolP384R1: "1.3.36.3.3.2.8.1.1.11",
	EllipticCurveBrainpoolP512R1: "1.3.36.3.3.2.8.1.1.13",
}

// EllipticCurveValues returns all constants in declaration order.
func EllipticCurveValues() []EllipticCurve {
	return []EllipticCurve{
		EllipticCurveP256,
		EllipticCurveP384,
		EllipticCurveP521,
		EllipticCurveX25519,
		EllipticCurveX448,
		EllipticCurveED25519,
		EllipticCurveED448,
		EllipticCurveSECP256K1,
		EllipticCurveBrainpoolP256R1,
		EllipticCurveBrainpoolP320R1,
		EllipticCurveBrainpoolP384R1,
		EllipticCurveBrainpoolP512R1,
	}
}

// Label returns a user-friendly label of the elliptic curve.
func (e EllipticCurve) Label() string {
	return ellipticCurveLabel[e]
}

// EllipticCurveForLabel gets an elliptic curve for the given label (same as
// COSE name). Returns "" (zero value) if label is empty or unknown,
// mirroring Java's null return.
func EllipticCurveForLabel(label string) EllipticCurve {
	if label == "" {
		return ""
	}
	for _, v := range EllipticCurveValues() {
		if ellipticCurveLabel[v] == label {
			return v
		}
	}
	return ""
}

// Size returns the coordinate byte size for the elliptic curve.
func (e EllipticCurve) Size() int {
	return ellipticCurveSize[e]
}

// Parameter returns the elliptic curve parameter specification, and whether
// one is registered for this curve (X25519, X448, ED25519, ED448 have
// none, mirroring Java's map miss returning null).
func (e EllipticCurve) Parameter() (EllipticCurveParameter, bool) {
	p, ok := ellipticCurveParameters[e]
	return p, ok
}

func ellipticCurveEqualParameters(a, b EllipticCurveParameter) bool {
	return a.P.Cmp(b.P) == 0 &&
		a.A.Cmp(b.A) == 0 &&
		a.B.Cmp(b.B) == 0 &&
		a.Gx.Cmp(b.Gx) == 0 &&
		a.Gy.Cmp(b.Gy) == 0 &&
		a.N.Cmp(b.N) == 0
}

// EllipticCurveForParameter gets the elliptic curve for the given parameter
// spec. Returns "" (zero value) if not found, mirroring Java's null return.
func EllipticCurveForParameter(parameter EllipticCurveParameter) EllipticCurve {
	for _, v := range EllipticCurveValues() {
		p, ok := ellipticCurveParameters[v]
		if ok && ellipticCurveEqualParameters(parameter, p) {
			return v
		}
	}
	return ""
}

// COSEValue returns the COSE value as defined in the IANA registry, and
// whether one is registered for this curve.
func (e EllipticCurve) COSEValue() (int64, bool) {
	v, ok := ellipticCurveCOSEValues[e]
	return v, ok
}

// EllipticCurveForCOSEValue gets an elliptic curve for the given COSE value.
// Returns "" (zero value) if not found, mirroring Java's null return.
func EllipticCurveForCOSEValue(value int64) EllipticCurve {
	for _, v := range EllipticCurveValues() {
		if cv, ok := ellipticCurveCOSEValues[v]; ok && cv == value {
			return v
		}
	}
	return ""
}

// OID returns the OID value.
func (e EllipticCurve) OID() string {
	return ellipticCurveOIDs[e]
}

// EllipticCurveForOID gets an elliptic curve for the given OID value.
// Returns "" (zero value) if not found, mirroring Java's null return.
// NOTE: since P_256 and SECP_256K1 share the same OID in the upstream Java
// source, this lookup (like Java's reversed-map construction, where the
// later map.put call for a duplicate key wins) resolves to whichever
// constant comes later in EllipticCurveValues() for that shared OID
// ("1.2.840.10045.3.1.7") — SECP_256K1.
func EllipticCurveForOID(oid string) EllipticCurve {
	// Scan in declaration order and keep the last match, reproducing Java's
	// HashMap put-order semantics (the later registerOIDs() put call for a
	// duplicate key wins).
	var result EllipticCurve
	for _, v := range EllipticCurveValues() {
		if ellipticCurveOIDs[v] == oid {
			result = v
		}
	}
	return result
}
