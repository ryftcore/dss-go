// Ported from the number handling of org.jose4j.json.internal.json_simple.parser.Yylex (which
// materialises JSON numbers as java.lang.Long, java.math.BigInteger or java.lang.Double) and
// from the java.lang.Long.toString / java.math.BigInteger.toString / java.lang.Double.toString
// renderings that JSONValue.writeJSONString then emits (jose4j 0.9.6, JDK 21).
package jose

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// NumberKind distinguishes the three Java classes json-simple can produce for a JSON number.
// It has to be kept because the class decides how the value is rendered again: the token "1e3"
// parses to a Double and serializes back as "1000.0", never as "1e3" and never as "1000".
type NumberKind int

const (
	// NumberLong is java.lang.Long: an integer token that fits in a signed 64-bit value.
	NumberLong NumberKind = iota
	// NumberBigInteger is java.math.BigInteger: an integer token too large for a Long, which
	// Yylex falls back to inside its NumberFormatException handler.
	NumberBigInteger
	// NumberDouble is java.lang.Double: any token with a fraction or an exponent.
	NumberDouble
)

// Number is a JSON number with the Java class json-simple would have given it. jose4j has no
// such type - it just holds a java.lang.Number - but Go needs one, because "1" and "1.0" are
// the same float64 and must serialize differently, and because Long and Double render through
// entirely different algorithms.
type Number struct {
	kind NumberKind
	i    int64
	b    *big.Int
	f    float64
}

// NewLong returns a Number backed by a java.lang.Long.
func NewLong(v int64) *Number { return &Number{kind: NumberLong, i: v} }

// NewBigInteger returns a Number backed by a java.math.BigInteger. v is copied.
func NewBigInteger(v *big.Int) *Number {
	return &Number{kind: NumberBigInteger, b: new(big.Int).Set(v)}
}

// NewDouble returns a Number backed by a java.lang.Double.
func NewDouble(v float64) *Number { return &Number{kind: NumberDouble, f: v} }

// Kind reports which Java class backs the number.
func (n *Number) Kind() NumberKind { return n.kind }

// IsIntegral reports whether the number came from an integer token (Long or BigInteger).
func (n *Number) IsIntegral() bool { return n.kind != NumberDouble }

// Int64 returns the value as an int64, the way JsonHelp.getLong does with ((Number) o)
// .longValue(): a Double is truncated toward zero and a BigInteger is taken modulo 2^64, both
// exactly as the JVM's narrowing conversions do.
func (n *Number) Int64() int64 {
	switch n.kind {
	case NumberLong:
		return n.i
	case NumberBigInteger:
		// BigInteger.longValue() returns the low-order 64 bits, two's complement.
		return n.b.Int64()
	default:
		// Java's (long) double: NaN is 0, out-of-range saturates at MIN/MAX_VALUE.
		if math.IsNaN(n.f) {
			return 0
		}
		if n.f >= math.MaxInt64 {
			return math.MaxInt64
		}
		if n.f <= math.MinInt64 {
			return math.MinInt64
		}
		return int64(n.f)
	}
}

// Float64 returns the value as a float64 (Java's Number.doubleValue()).
func (n *Number) Float64() float64 {
	switch n.kind {
	case NumberLong:
		return float64(n.i)
	case NumberBigInteger:
		f, _ := new(big.Float).SetInt(n.b).Float64()
		return f
	default:
		return n.f
	}
}

// String renders the number the way JSONValue.writeJSONString does, i.e. via the backing class's
// toString(). For a Double that is infinite or NaN the writer emits "null" instead and never
// calls this; see writeJSONValue.
func (n *Number) String() string {
	switch n.kind {
	case NumberLong:
		return strconv.FormatInt(n.i, 10)
	case NumberBigInteger:
		return n.b.String()
	default:
		return JavaDoubleToString(n.f)
	}
}

// JavaDoubleToString reproduces java.lang.Double.toString(double) as specified since JDK 19
// (JDK-4511638, "Double.toString(double) sometimes produces inaccurate results"), which is the
// behaviour of the JDK 21 that builds DSS 6.5.RC1.
//
// The specification, restated: let R be the set of decimals that read back as exactly v. Let p
// be the smallest length of any member of R. Take T = the members of R of length p, widened to
// also include the members of length 2 when p is 1. The rendered decimal is the member of T
// closest to v, preferring the shorter one on a tie. It is then laid out as
//
//	10^-3 <= |v| < 10^7   ->  plain,       e.g. 0.001, 0.1, 1234567.0
//	otherwise             ->  scientific,  e.g. 1.0E-4, 1.0E7, 4.9E-324
//
// with the invariant that there is always at least one digit on each side of the point.
//
// Go's strconv gives the length-p digits directly (FormatFloat with precision -1), so the only
// extra work is the "widen 1 digit to 2" clause - which is not academic: it is the difference
// between Java's "4.9E-324" and Go's "5e-324" for Double.MIN_VALUE.
//
// This matters far less than the rest of this package: JAdES numbers are NumericDate claims and
// counters, i.e. integers, so a Double reaches serialization only if a signature under
// validation carries one. It is still exact, and TestJavaDoubleToString pins it against the
// jose4j/JDK oracle.
func JavaDoubleToString(v float64) string {
	switch {
	case math.IsNaN(v):
		return "NaN"
	case math.IsInf(v, 1):
		return "Infinity"
	case math.IsInf(v, -1):
		return "-Infinity"
	}
	if v == 0 {
		if math.Signbit(v) {
			return "-0.0"
		}
		return "0.0"
	}

	neg := v < 0
	av := math.Abs(v)

	digits, exp10 := javaShortestDigits(av)

	var sb strings.Builder
	if neg {
		sb.WriteByte('-')
	}
	// exp10 is the power of ten of the first digit: av ~= 0.<digits> * 10^(exp10+1).
	// Java switches to scientific notation outside [10^-3, 10^7).
	if exp10 >= -3 && exp10 < 7 {
		writeJavaPlainDecimal(&sb, digits, exp10)
	} else {
		writeJavaScientificDecimal(&sb, digits, exp10)
	}
	return sb.String()
}

// javaShortestDigits returns the significant digits and the decimal exponent of the first digit
// for a positive, finite v, applying the "at least two digits when the shortest is one" rule.
func javaShortestDigits(av float64) (digits string, exp10 int) {
	shortest := strconv.FormatFloat(av, 'e', -1, 64) // e.g. "4.9406564584124654e-324" -> shortest form
	digits, exp10 = splitScientific(shortest)

	if len(digits) != 1 {
		return digits, exp10
	}
	// p == 1: Java also considers the length-2 decimals and picks whichever is closer to av,
	// keeping the shorter one when they are equidistant - which includes the common case where
	// the two decimals denote the same number (0.1 and 0.10), and excludes Double.MIN_VALUE,
	// where 4.9e-324 really is closer to the double than 5e-324 is.
	//
	// The two candidates must be compared as exact decimals, not as float64: both by
	// construction read back as av, so comparing their float64 values would always tie.
	two := strconv.FormatFloat(av, 'e', 1, 64)
	twoDigits, twoExp := splitScientific(two)
	oneDec, ok1 := new(big.Rat).SetString(shortest)
	twoDec, ok2 := new(big.Rat).SetString(two)
	if !ok1 || !ok2 {
		return digits, exp10
	}
	exact := new(big.Rat).SetFloat64(av)
	if exact == nil {
		return digits, exp10
	}
	distOne := new(big.Rat).Sub(oneDec, exact)
	distOne.Abs(distOne)
	distTwo := new(big.Rat).Sub(twoDec, exact)
	distTwo.Abs(distTwo)
	if distTwo.Cmp(distOne) < 0 {
		return twoDigits, twoExp
	}
	return digits, exp10
}

// splitScientific turns strconv's "d.dddde±dd" into the bare digit string and the decimal
// exponent of its leading digit.
func splitScientific(s string) (digits string, exp10 int) {
	mantissa := s
	if i := strings.IndexByte(s, 'e'); i >= 0 {
		mantissa = s[:i]
		exp10, _ = strconv.Atoi(s[i+1:])
	}
	digits = strings.Replace(mantissa, ".", "", 1)
	// strconv never emits leading zeros in the mantissa of %e for a non-zero value, but it can
	// emit trailing ones when an explicit precision was asked for; Java keeps them, because they
	// are what make the decimal two digits long.
	return digits, exp10
}

// writeJavaPlainDecimal renders digits*10^(exp10-len+1) without an exponent, always leaving at
// least one digit on each side of the point.
func writeJavaPlainDecimal(sb *strings.Builder, digits string, exp10 int) {
	switch {
	case exp10 < 0:
		// 0.00ddd
		sb.WriteString("0.")
		for i := 0; i < -exp10-1; i++ {
			sb.WriteByte('0')
		}
		sb.WriteString(digits)
	case exp10 >= len(digits)-1:
		// ddd000.0
		sb.WriteString(digits)
		for i := 0; i < exp10-len(digits)+1; i++ {
			sb.WriteByte('0')
		}
		sb.WriteString(".0")
	default:
		sb.WriteString(digits[:exp10+1])
		sb.WriteByte('.')
		sb.WriteString(digits[exp10+1:])
	}
}

// writeJavaScientificDecimal renders d.dddEn, forcing a fractional zero for a single digit.
func writeJavaScientificDecimal(sb *strings.Builder, digits string, exp10 int) {
	sb.WriteByte(digits[0])
	sb.WriteByte('.')
	if len(digits) == 1 {
		sb.WriteByte('0')
	} else {
		sb.WriteString(digits[1:])
	}
	sb.WriteByte('E')
	sb.WriteString(strconv.Itoa(exp10))
}
