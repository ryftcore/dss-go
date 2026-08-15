package fc

import (
	"math/big"
	"testing"

	"github.com/utain/esig/dss/diagnostic"
	diagjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
)

// fakeSignatureWrapper is a minimal diagnostic.AbstractSignatureWrapperOverrides for tests
// that only exercise logic depending on Id() and PDFRevision().
type fakeSignatureWrapper struct {
	id          string
	pdfRevision *diagnostic.PDFRevisionWrapper
}

func (f *fakeSignatureWrapper) CurrentBasicSignature() *diagjaxb.XmlBasicSignature { return nil }
func (f *fakeSignatureWrapper) CurrentCertificateChain() []*diagjaxb.XmlChainItem  { return nil }
func (f *fakeSignatureWrapper) CurrentSigningCertificate() *diagjaxb.XmlSigningCertificate {
	return nil
}
func (f *fakeSignatureWrapper) FoundCertificates() *diagnostic.FoundCertificatesProxy { return nil }
func (f *fakeSignatureWrapper) FoundRevocations() *diagnostic.FoundRevocationsProxy   { return nil }
func (f *fakeSignatureWrapper) DigestMatchers() []*diagjaxb.XmlDigestMatcher          { return nil }
func (f *fakeSignatureWrapper) Binaries() []byte                                      { return nil }
func (f *fakeSignatureWrapper) Id() string                                            { return f.id }
func (f *fakeSignatureWrapper) Filename() string                                      { return "" }
func (f *fakeSignatureWrapper) PDFRevision() *diagnostic.PDFRevisionWrapper           { return f.pdfRevision }

func pdfRevisionWithByteRange(byteRange ...int64) *diagnostic.PDFRevisionWrapper {
	values := make([]*big.Int, len(byteRange))
	for i, v := range byteRange {
		values[i] = big.NewInt(v)
	}
	xmlRevision := &diagjaxb.XmlPDFRevision{
		PDFSignatureDictionary: &diagjaxb.XmlPDFSignatureDictionary{
			SignatureByteRange: &diagjaxb.XmlByteRange{Value: values},
		},
	}
	return diagnostic.NewPDFRevisionWrapper(xmlRevision)
}

// TestByteRangeCollisionCheck_collideRevisions exercises the pure byte-range collision math
// ByteRangeCollisionCheck.Process() delegates to, over the happy path (sequential, disjoint
// ranges) and the failure path (one signature's range starts inside the other's gap).
func TestByteRangeCollisionCheck_collideRevisions(t *testing.T) {
	tests := []struct {
		name          string
		current       []int64
		other         []int64
		wantCollision bool
	}{
		{"happy: sequential disjoint ranges", []int64{0, 100, 200, 50}, []int64{300, 50, 400, 50}, false},
		{"failure: other range starts inside current's gap", []int64{0, 100, 200, 50}, []int64{50, 10, 300, 50}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := &fakeSignatureWrapper{id: "current", pdfRevision: pdfRevisionWithByteRange(tt.current...)}
			other := &fakeSignatureWrapper{id: "other", pdfRevision: pdfRevisionWithByteRange(tt.other...)}

			if got := collideRevisions(current, other); got != tt.wantCollision {
				t.Errorf("collideRevisions() = %v, want %v", got, tt.wantCollision)
			}
		})
	}
}
