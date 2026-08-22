package claim

import (
	"testing"
	"time"
)

// TestAbstractClaimStringMatchesJava pins AbstractClaim#toString, which every
// concrete claim inherits:
//
//	getClass().getSimpleName() + " {" + "'" + name + "'"
//	    + (selectivelyDisclosable ? " (disclosure)" : "") + ": " + getValueAsString() + '}'
func TestAbstractClaimStringMatchesJava(t *testing.T) {
	date := time.Date(2024, time.March, 5, 6, 7, 8, 0, time.UTC)
	yes := true

	tests := []struct {
		name  string
		claim Claim
		want  string
	}{
		{
			name:  "string claim",
			claim: NewClaimStringWithName("given_name", "Alice"),
			want:  "ClaimString {'given_name': Alice}",
		},
		{
			name:  "selectively disclosable claim carries the disclosure marker",
			claim: NewClaimStringWithDisclosable("family_name", "Smith", true),
			want:  "ClaimString {'family_name' (disclosure): Smith}",
		},
		{
			name:  "boolean claim",
			claim: NewClaimBooleanWithName("adult", &yes),
			want:  "ClaimBoolean {'adult': true}",
		},
		{
			name:  "date claim uses the RFC 3339 value form",
			claim: NewClaimDateWithName("birth_date", &date),
			want:  "ClaimDate {'birth_date': 2024-03-05T06:07:08Z}",
		},
		{
			// ClaimNull#getValueAsString returns the literal "null".
			name:  "null claim",
			claim: NewClaimNullWithName("nothing"),
			want:  "ClaimNull {'nothing': null}",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := AbstractClaimString(tc.claim); got != tc.want {
				t.Errorf("AbstractClaimString() = %q, want %q", got, tc.want)
			}
			// Every concrete claim's String() must forward to the same rendering.
			stringer, ok := tc.claim.(interface{ String() string })
			if !ok {
				t.Fatalf("%T does not implement String()", tc.claim)
			}
			if got := stringer.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestEveryConcreteClaimImplementsString guards against a concrete claim class
// losing the toString() it inherits from AbstractClaim in Java.
func TestEveryConcreteClaimImplementsString(t *testing.T) {
	// The Java classes that extend AbstractClaim, and therefore inherit toString().
	concrete := []any{
		&Array{}, &Boolean{}, &ByteString{}, &Date{},
		&Map{}, &Null{}, &Number{}, &String{},
	}
	for _, c := range concrete {
		if _, ok := c.(interface{ String() string }); !ok {
			t.Errorf("%T must implement String() (inherited AbstractClaim#toString)", c)
		}
	}
}
