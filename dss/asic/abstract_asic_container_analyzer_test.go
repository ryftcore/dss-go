package asic

import (
	"testing"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// analyzerLeafWithoutOverride stands in for ASiCContainerWithXAdESAnalyzer: it embeds the base
// and does NOT declare AttachExternalTimestamps, so it can only satisfy
// AbstractASiCContainerAnalyzerOverrides through the base's promoted default body - the Go
// equivalent of a Java subclass inheriting a non-abstract protected method.
type analyzerLeafWithoutOverride struct{ *AbstractASiCContainerAnalyzer }

func (l *analyzerLeafWithoutOverride) IsSupportedASiCContent(*ASiCContent) bool { return true }
func (l *analyzerLeafWithoutOverride) GetContainerExtractor() *DefaultASiCContainerExtractor {
	return nil
}
func (l *analyzerLeafWithoutOverride) GetManifestFilesDescriptions() []*model.ManifestFile {
	return nil
}
func (l *analyzerLeafWithoutOverride) GetSignatureAnalyzers() []analyzer.DocumentAnalyzer {
	return nil
}

// analyzerLeafWithOverride stands in for ASiCContainerWithCAdESAnalyzer, which shadows the
// base's empty body with one that really attaches container-level timestamps.
type analyzerLeafWithOverride struct{ analyzerLeafWithoutOverride }

func (l *analyzerLeafWithOverride) AttachExternalTimestamps(
	[]validation.AdvancedSignature) []*validation.TimestampToken {
	return []*validation.TimestampToken{{}}
}

var _ AbstractASiCContainerAnalyzerOverrides = (*analyzerLeafWithoutOverride)(nil)
var _ AbstractASiCContainerAnalyzerOverrides = (*analyzerLeafWithOverride)(nil)

// TestAttachExternalTimestampsIsVirtual guards the virtual-dispatch contract that
// AbstractASiCContainerAnalyzer.GetAllSignatures depends on. Java calls
// attachExternalTimestamps() virtually and only ASiCContainerWithCAdESAnalyzer overrides it, to
// attach ASiC-S container / ASiC-E archive timestamps to the signatures they cover. Dispatching
// on the embedded base instead of through AbstractASiCContainerAnalyzerOverrides would bind to
// the empty default and silently drop every such timestamp, so both directions are pinned here:
// a leaf that overrides must win, and a leaf that does not must still inherit the default.
func TestAttachExternalTimestampsIsVirtual(t *testing.T) {
	withoutOverride := &analyzerLeafWithoutOverride{
		AbstractASiCContainerAnalyzer: NewAbstractASiCContainerAnalyzerBase()}
	withOverride := &analyzerLeafWithOverride{analyzerLeafWithoutOverride{
		AbstractASiCContainerAnalyzer: NewAbstractASiCContainerAnalyzerBase()}}

	var inherited AbstractASiCContainerAnalyzerOverrides = withoutOverride
	if got := inherited.AttachExternalTimestamps(nil); len(got) != 0 {
		t.Errorf("leaf without an override should inherit the empty default, got %d timestamps", len(got))
	}

	var overridden AbstractASiCContainerAnalyzerOverrides = withOverride
	if got := overridden.AttachExternalTimestamps(nil); len(got) != 1 {
		t.Errorf("leaf with an override should reach it, got %d timestamps, want 1", len(got))
	}

	// And the base really routes through the interface rather than calling its own method: with
	// the overriding leaf registered, the base must see the override's result.
	withOverride.InitAbstractASiCContainerAnalyzer(withOverride)
	if got := withOverride.requireOverrides().AttachExternalTimestamps(nil); len(got) != 1 {
		t.Errorf("base dispatch reached the default body instead of the leaf override (got %d, want 1)", len(got))
	}
}
