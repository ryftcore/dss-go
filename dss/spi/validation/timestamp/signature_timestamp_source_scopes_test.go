package timestamp

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// scopeProbeAttr is the smallest SignatureAttribute the generic source accepts.
type scopeProbeAttr struct{ validation.SignatureAttribute }

// scopeProbeBuilder answers every message-imprint request with an empty digest, which matches
// no timestamp.
type scopeProbeBuilder struct{}

func (scopeProbeBuilder) ContentTimestampMessageDigest() model.DSSMessageDigest {
	return model.NewDSSMessageDigest()
}
func (scopeProbeBuilder) SignatureTimestampMessageDigest() model.DSSMessageDigest {
	return model.NewDSSMessageDigest()
}
func (scopeProbeBuilder) TimestampX1MessageDigest() model.DSSMessageDigest {
	return model.NewDSSMessageDigest()
}
func (scopeProbeBuilder) TimestampX2MessageDigest() model.DSSMessageDigest {
	return model.NewDSSMessageDigest()
}
func (scopeProbeBuilder) ArchiveTimestampMessageDigest() model.DSSMessageDigest {
	return model.NewDSSMessageDigest()
}

// scopeProbe implements just the two overrides validateTimestamps() reaches in this scenario;
// the embedded nil interface would panic on any other, which keeps the test honest about what
// the base calls.
type scopeProbe struct {
	SignatureTimestampSourceOverrides[validation.AdvancedSignature, scopeProbeAttr]
	scoped []*validation.TimestampToken
}

func (p *scopeProbe) GetTimestampMessageImprintDigestBuilderForToken(*validation.TimestampToken) MessageDigestBuilder {
	return scopeProbeBuilder{}
}

func (p *scopeProbe) GetTimestampScopes(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	p.scoped = append(p.scoped, timestampToken)
	return nil
}

// TestValidateTimestampsUsesTheFormatSpecificScopeFinder pins that the base's validateTimestamps
// asks the concrete source for the scopes of every content and archive timestamp, as Java's
// virtual getTimestampScopes does (XAdESTimestampSource overrides it to filter by the
// IndividualDataObjectsTimestamp includes; before GetTimestampScopes joined
// SignatureTimestampSourceOverrides the base always used the XAdES-agnostic finder and
// over-attributed the timestamp to every signature scope).
func TestValidateTimestampsUsesTheFormatSpecificScopeFinder(t *testing.T) {
	content := loadFixtureTimestampToken(t, enumerations.TimestampTypeContentTimestamp)
	archive := loadFixtureTimestampToken(t, enumerations.TimestampTypeArchiveTimestamp)
	signature := loadFixtureTimestampToken(t, enumerations.TimestampTypeSignatureTimestamp)

	probe := &scopeProbe{}
	source := &SignatureTimestampSource[validation.AdvancedSignature, scopeProbeAttr]{}
	source.InitSignatureTimestampSource(probe)
	source.contentTimestamps = []*validation.TimestampToken{content}
	source.signatureTimestamps = []*validation.TimestampToken{signature}
	source.sigAndRefsTimestamps = []*validation.TimestampToken{}
	source.refsOnlyTimestamps = []*validation.TimestampToken{}
	source.archiveTimestamps = []*validation.TimestampToken{archive}

	source.validateTimestamps()

	if len(probe.scoped) != 2 || probe.scoped[0] != content || probe.scoped[1] != archive {
		t.Fatalf("GetTimestampScopes was asked about %d timestamps, want the content then the archive timestamp only", len(probe.scoped))
	}
}
