package asic

import (
	"errors"
	"fmt"
	"testing"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
)

// The existing threshold tests pin the DEFAULT values and exercise the guards at those defaults.
// That leaves a hole: a guard that ignored its configured limit and hard-coded the default would
// pass every one of them. The tests below close it by mutation - each moves one limit off its
// default in both directions and requires the outcome to follow the limit, so a guard that stopped
// reading its field would fail. Each also drives the value through SecureContainerHandlerBuilder,
// since that is how upstream callers configure a handler.

// TestSecureContainerHandlerThresholdIsLoadBearing mutates `threshold`, the inflated-size floor
// below which the compression-ratio check is not applied at all (upstream ANDs the two).
func TestSecureContainerHandlerThresholdIsLoadBearing(t *testing.T) {
	// ~4 MB of zeroes: deflates tiny, so the ratio check trips whenever the threshold lets it.
	payload := make([]byte, 4*1024*1024)
	archive := zipCoreBuildArchive(t, []zipCoreTestEntry{{Name: "payload.bin", Content: payload}})
	if int64(len(archive))*100 >= int64(len(payload)) {
		t.Fatalf("fixture archive is not compressible enough to exercise the ratio guard (archive %d, inflated %d)", len(archive), len(payload))
	}
	doc := model.NewInMemoryDocumentWithName(archive, "payload.zip")

	// Raised above the inflated size, the guard must let the same archive through.
	raised := NewSecureContainerHandler()
	raised.SetThreshold(int64(len(payload)) + 1)
	if _, err := raised.ExtractContainerContent(doc); err != nil {
		t.Errorf("threshold above the inflated size must disarm the guard, got %v", err)
	}

	// Lowered below it, the guard must trip - and it must trip through the builder too.
	for label, handler := range map[string]ZipContainerHandler{
		"setter": func() ZipContainerHandler {
			h := NewSecureContainerHandler()
			h.SetThreshold(1024)
			return h
		}(),
		"builder": NewSecureContainerHandlerBuilder().SetThreshold(1024).Build(),
	} {
		var illegalInput *exception.IllegalInputException
		if _, err := handler.ExtractContainerContent(doc); err == nil {
			t.Errorf("%s: lowered threshold must trip the zip-bomb guard", label)
		} else if !errors.As(err, &illegalInput) ||
			illegalInput.Message != "Zip Bomb detected in the ZIP container. Validation is interrupted." {
			t.Errorf("%s: error = %v, want the zip-bomb IllegalInputException", label, err)
		}
	}
}

// TestSecureContainerHandlerMaxCompressionRatioIsLoadBearing mutates `maxCompressionRatio`, the
// other half of the AND: the inflated size must also exceed containerSize * ratio.
func TestSecureContainerHandlerMaxCompressionRatioIsLoadBearing(t *testing.T) {
	payload := make([]byte, 4*1024*1024)
	archive := zipCoreBuildArchive(t, []zipCoreTestEntry{{Name: "payload.bin", Content: payload}})
	doc := model.NewInMemoryDocumentWithName(archive, "payload.zip")

	actualRatio := int64(len(payload)) / int64(len(archive))
	if actualRatio < 2 {
		t.Fatalf("fixture ratio %d is too low to mutate around", actualRatio)
	}

	// A ratio ceiling above what this archive achieves disarms the guard...
	permissive := NewSecureContainerHandler()
	permissive.SetMaxCompressionRatio(actualRatio * 2)
	if _, err := permissive.ExtractContainerContent(doc); err != nil {
		t.Errorf("ratio ceiling above the archive's actual ratio must disarm the guard, got %v", err)
	}

	// ...and one below it arms the guard, via setter and builder alike.
	for label, handler := range map[string]ZipContainerHandler{
		"setter": func() ZipContainerHandler {
			h := NewSecureContainerHandler()
			h.SetMaxCompressionRatio(2)
			return h
		}(),
		"builder": NewSecureContainerHandlerBuilder().SetMaxCompressionRatio(2).Build(),
	} {
		var illegalInput *exception.IllegalInputException
		if _, err := handler.ExtractContainerContent(doc); err == nil {
			t.Errorf("%s: ratio ceiling below the archive's actual ratio must trip the guard", label)
		} else if !errors.As(err, &illegalInput) {
			t.Errorf("%s: error = %v, want an IllegalInputException", label, err)
		}
	}
}

// TestSecureContainerHandlerMaxAllowedFilesAmountIsLoadBearing mutates the entry-count ceiling
// around a fixed 20-entry archive, so the pass/fail flip can only come from the configured value.
func TestSecureContainerHandlerMaxAllowedFilesAmountIsLoadBearing(t *testing.T) {
	entries := make([]zipCoreTestEntry, 0, 20)
	for i := 0; i < 20; i++ {
		entries = append(entries, zipCoreTestEntry{Name: fmt.Sprintf("f%02d.txt", i), Content: []byte("x")})
	}
	doc := model.NewInMemoryDocumentWithName(zipCoreBuildArchive(t, entries), "many.zip")
	want := "Too many files detected. Cannot extract ASiC content from the file."

	// Exactly at the entry count the archive holds, the guard must NOT trip: upstream compares
	// with a strict >, so 20 entries against a ceiling of 20 is allowed.
	atLimit := NewSecureContainerHandler()
	atLimit.SetMaxAllowedFilesAmount(20)
	if documents, err := atLimit.ExtractContainerContent(doc); err != nil {
		t.Errorf("a ceiling equal to the entry count must not trip the guard, got %v", err)
	} else if len(documents) != 20 {
		t.Errorf("documents = %d, want 20", len(documents))
	}

	// One below, it must trip - through setter and builder.
	for label, handler := range map[string]ZipContainerHandler{
		"setter": func() ZipContainerHandler {
			h := NewSecureContainerHandler()
			h.SetMaxAllowedFilesAmount(19)
			return h
		}(),
		"builder": NewSecureContainerHandlerBuilder().SetMaxAllowedFilesAmount(19).Build(),
	} {
		var illegalInput *exception.IllegalInputException
		if _, err := handler.ExtractContainerContent(doc); err == nil {
			t.Errorf("%s: a ceiling below the entry count must trip the file-count guard", label)
		} else if !errors.As(err, &illegalInput) || illegalInput.Message != want {
			t.Errorf("%s: error = %v, want IllegalInputException %q", label, err, want)
		}
	}
}

// TestSecureContainerHandlerMaxMalformedFilesIsLoadBearing mutates `maxMalformedFiles`, the cap on
// how many malformed entries getNextValidEntry may skip before giving up. Upstream's loop is
// `while (malformedFilesCounter < maxMalformedFiles)` followed by an unconditional throw, so a
// budget of 0 skips the loop entirely and raises DSSException("Unable to retrieve a valid ZipEntry
// (0 tries)") before the first entry is even read. The interpolated count in that message is what
// makes this a mutation test: it can only be right if the guard read the configured field.
func TestSecureContainerHandlerMaxMalformedFilesIsLoadBearing(t *testing.T) {
	// A real upstream fixture with a malformed entry name: the default budget of 100 recovers and
	// returns 5 documents (pinned by TestSecureContainerHandlerStopsAtMalformedEntryNames).
	const fixture = "dss-asic-cades/src/test/resources/validation/cp852encoded_signature.asice"
	doc := zipCoreFileDocument(t, fixture)

	baseline := NewSecureContainerHandler()
	documents, err := baseline.ExtractContainerContent(doc)
	if err != nil {
		t.Fatalf("the default budget must read this fixture: %v", err)
	}
	if len(documents) != 5 || !baseline.malformedEntriesDetected() {
		t.Fatalf("baseline = %d documents (malformed detected: %v), want 5 and true",
			len(documents), baseline.malformedEntriesDetected())
	}

	for _, budget := range []int{0, 1, 2} {
		want := fmt.Sprintf("Unable to retrieve a valid ZipEntry (%d tries)", budget)
		for label, handler := range map[string]ZipContainerHandler{
			"setter": func() ZipContainerHandler {
				h := NewSecureContainerHandler()
				h.SetMaxMalformedFiles(budget)
				return h
			}(),
			"builder": NewSecureContainerHandlerBuilder().SetMaxMalformedFiles(budget).Build(),
		} {
			_, err := handler.ExtractContainerContent(doc)
			if err == nil {
				t.Errorf("%s: a budget of %d must exhaust before the fixture's malformed entry", label, budget)
				continue
			}
			if err.Error() != want {
				t.Errorf("%s: error = %q, want %q", label, err.Error(), want)
			}
		}
	}

	// The fixture raises three malformed-entry events while its stream is walked, so 3 is the
	// smallest budget that clears it. Pinning that boundary exactly - 2 fails above, 3 succeeds
	// here - shows the failures are the budget biting rather than the fixture being unreadable,
	// and would catch an off-by-one in the loop's `counter < max` condition.
	generous := NewSecureContainerHandler()
	generous.SetMaxMalformedFiles(3)
	if got, err := generous.ExtractContainerContent(doc); err != nil {
		t.Errorf("a budget of 3 must clear this fixture's malformed entries: %v", err)
	} else if len(got) != 5 {
		t.Errorf("documents = %d, want 5", len(got))
	}
}
