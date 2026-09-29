// Regression cover for SingleDssDict#extractVRIs on a hostile /DSS /VRI dictionary.
//
// Upstream walks the /VRI entries inside try { ... } catch (Exception e) and returns an empty
// list when anything throws. `new PdfVriDict(name, vriDict.getAsDict(name))` throws a
// NullPointerException when the value under a valid 40-hex key is not a dictionary, so a single
// malformed entry makes upstream drop every VRI dictionary of the document and carry on. The
// port let the equivalent nil dereference escape and abort the whole validation.
package pades

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/pdf"
)

const (
	testVRIKeyA = "0123456789ABCDEF0123456789ABCDEF01234567"
	testVRIKeyB = "FEDCBA9876543210FEDCBA9876543210FEDCBA98"
)

func testDssDictWithVRIs(vri *pdf.Dict) PdfDict {
	dss := pdf.NewDict()
	dss.Set(pdf.Name(PAdESConstantsVriDictionaryName), vri)
	return newNativePdfDict(&pdf.Document{}, dss, pdf.ObjectKey{})
}

func TestSingleDssDictVRIsOfWellFormedEntries(t *testing.T) {
	vri := pdf.NewDict()
	vri.Set(pdf.Name(testVRIKeyA), pdf.NewDict())
	vri.Set(pdf.Name(testVRIKeyB), pdf.NewDict())
	// not a SHA-1 digest key: skipped, as upstream's isDictionaryKey does
	vri.Set("NotADigest", pdf.Integer(1))

	vris := newSingleDssDict(testDssDictWithVRIs(vri)).VRIs()
	if len(vris) != 2 || vris[0].Name() != testVRIKeyA || vris[1].Name() != testVRIKeyB {
		t.Fatalf("VRIs = %v, want the two digest-keyed entries in order", vris)
	}
}

// A non-dictionary value under a digest key must not panic out of newSingleDssDict; upstream's
// catch (Exception) turns it into an empty VRI list, valid entries included.
func TestSingleDssDictVRIsSurviveNonDictionaryEntry(t *testing.T) {
	for name, value := range map[string]pdf.Object{
		"integer": pdf.Integer(42),
		"name":    pdf.Name("Oops"),
		"string":  pdf.String{Bytes: []byte("oops")},
		"null":    pdf.Null{},
	} {
		t.Run(name, func(t *testing.T) {
			vri := pdf.NewDict()
			vri.Set(pdf.Name(testVRIKeyA), pdf.NewDict())
			vri.Set(pdf.Name(testVRIKeyB), value)

			var dssDict *SingleDssDict
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("newSingleDssDict panicked on a non-dictionary VRI entry: %v", r)
					}
				}()
				dssDict = newSingleDssDict(testDssDictWithVRIs(vri))
			}()
			if got := dssDict.VRIs(); len(got) != 0 {
				t.Errorf("VRIs = %v, want none (upstream returns Collections.emptyList())", got)
			}
		})
	}
}
