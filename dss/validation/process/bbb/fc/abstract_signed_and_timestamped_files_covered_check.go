// Ported from dss-validation/.../validation/process/bbb/fc/checks/AbstractSignedAndTimestampedFilesCoveredCheck.java (DSS 6.5.RC1).
package fc

import (
	"fmt"
	"slices"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractSignedAndTimestampedFilesCoveredCheck checks whether all files signed by the covered
// signatures or timestamped by covered timestamps are covered by the current timestamp as well.
// T is the result type (extends XmlConstraintsConclusion in Java).
type AbstractSignedAndTimestampedFilesCoveredCheck[T any] struct {
	*process.ChainItemBase[T]

	DiagnosticData    *diagnostic.Data
	TimestampFilename string
}

// InitAbstractSignedAndTimestampedFilesCoveredCheck wires the shared state; called by the
// concrete constructor before InitChainItem.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) InitAbstractSignedAndTimestampedFilesCoveredCheck(
	i18nProvider *i18n.Provider, result *process.Result[T], diagnosticData *diagnostic.Data, timestampFilename string,
	constraint policy.LevelRule) {
	c.DiagnosticData = diagnosticData
	c.TimestampFilename = timestampFilename
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
}

// Process performs the check.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) Process() bool {
	manifestFile := c.DiagnosticData.ManifestFileForFilename(c.TimestampFilename)
	if manifestFile != nil {
		return c.CheckManifestFilesCovered(manifestFile.Entries.All()) && c.isAnyRootLevelDocumentCovered(manifestFile)
	}
	// no manifest -> no check is required (ASiC-S case)
	return true
}

// CheckManifestFilesCovered runs the validation process for the manifest entries coverage.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) CheckManifestFilesCovered(entries []string) bool {
	covered := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		covered[entry] = struct{}{}
	}
	return c.checkManifestFilesCoveredRecursively(covered, entries, map[string]struct{}{}, map[string]struct{}{}, true)
}

// checkManifestFilesCoveredRecursively ports the recursive method of the same
// name. coveredEntries is the List of the Java method as a set (List#contains
// is a membership test, and the list is never modified), which keeps the walk
// linear in the number of manifest entries instead of quadratic.
//
// expanding holds the manifest entries whose nested manifest is currently being
// walked (i.e. those on the recursion stack); it has no Java counterpart.
// DIVERGENCE, deliberate: Java records an entry in checkedEntries only after its
// nested manifest has been walked, so a manifest that lists its own time-stamp
// (or two manifests listing each other's) recurses without end and Java dies
// with a StackOverflowError - an Error the JVM survives. A Go stack overflow is
// a fatal runtime error that no recover() can catch and that ends the process,
// which a crafted container must not be able to cause. The cycle is therefore
// reported by an ordinary panic at the exact point Java would recurse into it
// again, so the facade's recovered() turns it into an error just as it does for
// the other unchecked exceptions of this port. No verdict is produced or
// changed: every input that terminates in Java yields the same result here.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) checkManifestFilesCoveredRecursively(
	coveredEntries map[string]struct{}, manifestEntries []string, checkedEntries map[string]struct{},
	expanding map[string]struct{}, rootProcess bool) bool {
	for _, manifestEntry := range manifestEntries {
		// skip validation for the first loop (same manifest is evaluated)
		if !rootProcess {
			if _, ok := coveredEntries[manifestEntry]; !ok {
				return false
			}
			if _, ok := checkedEntries[manifestEntry]; ok {
				continue
			}
		}
		entryManifest := c.DiagnosticData.ManifestFileForFilename(manifestEntry)
		if entryManifest != nil {
			if _, ok := expanding[manifestEntry]; ok {
				panic(fmt.Sprintf("Cyclic manifest reference detected for the entry '%s'", manifestEntry))
			}
			expanding[manifestEntry] = struct{}{}
			covered := c.checkManifestFilesCoveredRecursively(coveredEntries, entryManifest.Entries.All(),
				checkedEntries, expanding, false)
			delete(expanding, manifestEntry)
			if !covered {
				return false
			}
		}
		checkedEntries[manifestEntry] = struct{}{}
	}
	return true
}

func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) isAnyRootLevelDocumentCovered(timestampManifest *diagjaxb.XmlManifestFile) bool {
	root := rootLevelFiles(c.DiagnosticData.ContainerContentFilenames())
	entries := timestampManifest.Entries.All()
	for _, r := range root {
		if slices.Contains(entries, r) {
			return true
		}
	}
	return false
}

// MessageTag returns the constraint message i18n key.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCISFPASTFORAMC
}

// ErrorMessageTag returns the error message i18n key.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCISFPASTFORAMCANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
