package model

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

// TestDigestAlgorithmSupportIsConsistentAcrossThePackage pins that a digest
// algorithm usable for a document is also usable for a token or an identifier.
// In Java both paths call DigestAlgorithm#getMessageDigest, so they can never
// disagree; the Go port keeps one table per file and they must stay in step.
func TestDigestAlgorithmSupportIsConsistentAcrossThePackage(t *testing.T) {
	for _, algorithm := range enumerations.DigestAlgorithmValues() {
		_, documentErr := commonDocumentHash(algorithm)
		_, identifierErr := identifierMessageDigest(algorithm)
		if (documentErr == nil) != (identifierErr == nil) {
			t.Errorf("algorithm %s: CommonDocument supported=%t but Identifier supported=%t",
				algorithm, documentErr == nil, identifierErr == nil)
		}
	}
}

// TestDigestDocumentMessagesMatchJava pins the exact messages upstream raises,
// which surface in DSS validation reports.
func TestDigestDocumentMessagesMatchJava(t *testing.T) {
	document := NewDigestDocument()

	if _, err := document.ExistingDigest(); err == nil {
		t.Fatalf("ExistingDigest() on an empty document must fail")
	} else if !errors.Is(err, ErrNoDigest) {
		t.Errorf("ExistingDigest() error = %v, want ErrNoDigest", err)
	} else if want := "The DigestDocument does not contain any digest! You must specify it by using addDigest() method."; err.Error() != want {
		t.Errorf("ExistingDigest() error = %q, want %q", err.Error(), want)
	}

	if _, err := document.OpenStream(); err == nil {
		t.Fatalf("OpenStream() must fail for a digest-only document")
	} else if want := "Not possible with Digest document"; err.Error() != want {
		t.Errorf("OpenStream() error = %q, want %q", err.Error(), want)
	}
	if err := document.Save("/does/not/matter"); err == nil {
		t.Fatalf("Save() must fail for a digest-only document")
	} else if want := "Not possible with Digest document"; err.Error() != want {
		t.Errorf("Save() error = %q, want %q", err.Error(), want)
	}

	if _, err := document.DigestValue(enumerations.DigestAlgorithm_SHA256); err == nil {
		t.Fatalf("DigestValue() on an empty document must fail")
	} else if want := "The digest document does not contain a digest value for the algorithm : SHA256"; err.Error() != want {
		t.Errorf("DigestValue() error = %q, want %q", err.Error(), want)
	}

	err := document.AddDigestBase64(enumerations.DigestAlgorithm_SHA256, "not base64!!")
	if err == nil {
		t.Fatalf("AddDigestBase64() must reject invalid base64")
	}
	if want := "Unable to base64-decode string 'not base64!!' : "; !strings.HasPrefix(err.Error(), want) {
		t.Errorf("AddDigestBase64() error = %q, want prefix %q", err.Error(), want)
	}
}

// TestDigestDocumentAssertionsMatchJava pins the Objects.requireNonNull checks
// upstream performs, including the order in which they fire: the algorithm is
// validated before the value, so an empty Digest reports the algorithm first.
func TestDigestDocumentAssertionsMatchJava(t *testing.T) {
	tests := []struct {
		name string
		call func(*DigestDocument)
		want string
	}{
		{
			name: "missing algorithm",
			call: func(d *DigestDocument) { d.AddDigestValue("", []byte{0x01}) },
			want: "The Digest Algorithm is not defined",
		},
		{
			name: "missing value",
			call: func(d *DigestDocument) { d.AddDigestValue(enumerations.DigestAlgorithm_SHA256, nil) },
			want: "The digest value is not defined",
		},
		{
			name: "empty Digest reports the algorithm first",
			call: func(d *DigestDocument) { d.AddDigest(Digest{}) },
			want: "The Digest Algorithm is not defined",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					t.Fatalf("expected a panic")
				}
				if got, ok := recovered.(string); !ok || got != tc.want {
					t.Errorf("panic = %v, want %q", recovered, tc.want)
				}
			}()
			tc.call(NewDigestDocument())
		})
	}
}

// TestDataIdentifierPropagatesDocumentDigestFailure pins that a failure to
// digest the document is propagated as-is. Upstream's build() catches
// IOException only, and getDigestValue raises a DSSException, so the
// "Unable to build a JAdESAttributeIdentifier" wrapper never applies.
func TestDataIdentifierPropagatesDocumentDigestFailure(t *testing.T) {
	_, err := NewDataIdentifierForDocument("name", &documentContractFailingDocument{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	var dssError *DSSError
	if !errors.As(err, &dssError) {
		t.Fatalf("error = %T, want *DSSError", err)
	}
	if want := "Unable to compute the digest"; dssError.Message != want {
		t.Errorf("error message = %q, want %q", dssError.Message, want)
	}
	if strings.Contains(err.Error(), "JAdESAttributeIdentifier") {
		t.Errorf("the IOException wrapper must not be applied to this path: %q", err.Error())
	}
}

// TestDataIdentifierPropagatesMissingExistingDigest pins that the DigestDocument
// branch propagates the IllegalStateException, which upstream's catch(IOException)
// likewise does not swallow.
func TestDataIdentifierPropagatesMissingExistingDigest(t *testing.T) {
	_, err := NewDataIdentifierForDocument("name", NewDigestDocument())
	if !errors.Is(err, ErrNoDigest) {
		t.Errorf("error = %v, want ErrNoDigest", err)
	}
}

// documentContractFailingDocument is a DSSDocument whose stream always fails, so
// that CommonDocument#getDigestValue takes its DSSException branch.
type documentContractFailingDocument struct {
	CommonDocument
}

func (d *documentContractFailingDocument) OpenStream() (io.ReadCloser, error) {
	return nil, errors.New("stream unavailable")
}

func (d *documentContractFailingDocument) WriteTo(w io.Writer) (int64, error) {
	return commonDocumentWriteTo(d, w)
}

func (d *documentContractFailingDocument) Save(filePath string) error {
	return commonDocumentSave(d, filePath)
}

func (d *documentContractFailingDocument) Digest(a enumerations.DigestAlgorithm) (Digest, error) {
	return commonDocumentDigest(d, &d.CommonDocument, a)
}

func (d *documentContractFailingDocument) DigestValue(a enumerations.DigestAlgorithm) ([]byte, error) {
	return commonDocumentDigestValue(d, &d.CommonDocument, a)
}

var _ DSSDocument = (*documentContractFailingDocument)(nil)
