// Ported from dss-model/.../claim/ClaimAgeOverNN.java (DSS 6.5.RC1).
package claim

// AgeOverNN defines a claim containing a boolean value whether the
// age of EAA holder is over or less a defined value.
type AgeOverNN interface {
	Claim

	// Age gets the value of the age corresponding to the claim
	// definition. Nil-able (mirrors Java's Integer). Ports
	// ClaimAgeOverNN#getAge.
	Age() *int
}
