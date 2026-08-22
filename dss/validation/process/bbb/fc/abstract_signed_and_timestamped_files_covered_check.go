// Ported from dss-validation/.../validation/process/bbb/fc/checks/AbstractSignedAndTimestampedFilesCoveredCheck.java (DSS 6.5.RC1).
package fc

import (
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

	DiagnosticData    *diagnostic.DiagnosticData
	TimestampFilename string
}

// InitAbstractSignedAndTimestampedFilesCoveredCheck wires the shared state; called by the
// concrete constructor before InitChainItem.
func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) InitAbstractSignedAndTimestampedFilesCoveredCheck(
	i18nProvider *i18n.I18nProvider, result *process.Result[T], diagnosticData *diagnostic.DiagnosticData, timestampFilename string,
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
	return c.checkManifestFilesCoveredRecursively(entries, entries, map[string]struct{}{}, true)
}

func (c *AbstractSignedAndTimestampedFilesCoveredCheck[T]) checkManifestFilesCoveredRecursively(
	coveredEntries, manifestEntries []string, checkedEntries map[string]struct{}, rootProcess bool) bool {
	for _, manifestEntry := range manifestEntries {
		// skip validation for the first loop (same manifest is evaluated)
		if !rootProcess {
			if !slices.Contains(coveredEntries, manifestEntry) {
				return false
			}
			if _, ok := checkedEntries[manifestEntry]; ok {
				continue
			}
		}
		entryManifest := c.DiagnosticData.ManifestFileForFilename(manifestEntry)
		if entryManifest != nil &&
			!c.checkManifestFilesCoveredRecursively(coveredEntries, entryManifest.Entries.All(), checkedEntries, false) {
			return false
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
