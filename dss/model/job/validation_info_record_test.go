// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ValidationInfoRecord.java (DSS 6.5.RC1).
package job

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

type fakeValidationInfoRecord struct {
	fakeInfoRecord
	indication         enumerations.Indication
	subIndication      enumerations.SubIndication
	signingTime        time.Time
	signingCertificate *model.CertificateToken
	potentialSigners   []*model.CertificateToken
	valid              bool
	indeterminate      bool
	invalid            bool
}

func (f *fakeValidationInfoRecord) Indication() enumerations.Indication       { return f.indication }
func (f *fakeValidationInfoRecord) SubIndication() enumerations.SubIndication { return f.subIndication }
func (f *fakeValidationInfoRecord) SigningTime() time.Time                    { return f.signingTime }
func (f *fakeValidationInfoRecord) SigningCertificate() *model.CertificateToken {
	return f.signingCertificate
}
func (f *fakeValidationInfoRecord) PotentialSigners() []*model.CertificateToken {
	return f.potentialSigners
}
func (f *fakeValidationInfoRecord) IsValid() bool         { return f.valid }
func (f *fakeValidationInfoRecord) IsIndeterminate() bool { return f.indeterminate }
func (f *fakeValidationInfoRecord) IsInvalid() bool       { return f.invalid }

var _ ValidationInfoRecord = (*fakeValidationInfoRecord)(nil)

func TestValidationInfoRecord_RoundTrip(t *testing.T) {
	now := time.Now()
	rec := &fakeValidationInfoRecord{
		indication:    enumerations.IndicationTotalPassed,
		subIndication: enumerations.SubIndicationFormatFailure,
		signingTime:   now,
		valid:         true,
	}

	var vir ValidationInfoRecord = rec
	if vir.Indication() != enumerations.IndicationTotalPassed {
		t.Fatalf("Indication() = %v", vir.Indication())
	}
	if vir.SubIndication() != enumerations.SubIndicationFormatFailure {
		t.Fatalf("SubIndication() = %v", vir.SubIndication())
	}
	if !vir.SigningTime().Equal(now) {
		t.Fatalf("SigningTime() did not round-trip")
	}
	if !vir.IsValid() || vir.IsIndeterminate() || vir.IsInvalid() {
		t.Fatalf("status booleans did not round-trip")
	}
}
