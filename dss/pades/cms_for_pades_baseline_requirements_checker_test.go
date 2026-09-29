package pades

import (
	"os"
	"testing"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CMSForPAdESBaselineRequirementsChecker#isValidForPAdESBaselineBProfile ends in the protected
// cmsBaselineBRequirements(), which knows nothing of CAdES requirement (k) (a
// signature-policy-store needs a signature-policy-identifier defining sigPolicyHash); (k) is only
// part of hasBaselineBProfile(). A PAdES CMS carrying such a store must therefore still be
// accepted here.
func TestCMSForPAdESBaselineBDoesNotApplyCAdESRequirementK(t *testing.T) {
	data, err := os.ReadFile(padesFixturePath(t, "upstream/validation/pades3_Baseline_B.pdf"))
	if err != nil {
		t.Fatal(err)
	}
	var signatureCMS *cms.CMS
	for _, revision := range NewDefaultPdfObjFactory().NewPAdESSignatureService().
		GetRevisions(model.NewInMemoryDocument(data), nil) {
		if signatureRevision, ok := revision.(*PdfSignatureRevision); ok {
			signatureCMS = signatureRevision.PdfSigDictInfo().CMS()
			break
		}
	}
	if signatureCMS == nil {
		t.Fatal("no signature in the fixture")
	}

	messageDigest := model.NewDSSMessageDigestWithValue(enumerations.DigestAlgorithmSHA256, make([]byte, 32))
	checkerFor := func(t *testing.T, signedData *cms.CMS) *CMSForPAdESBaselineRequirementsChecker {
		t.Helper()
		return NewCMSForPAdESBaselineRequirementsChecker(padesWithExternalCMSServiceToCAdESSignature(signedData, messageDigest))
	}

	if checker := checkerFor(t, signatureCMS); !checker.IsValidForPAdESBaselineBProfile() || !checker.HasBaselineBProfile() {
		t.Fatal("the fixture's CMS is not PAdES-B compliant, the test has no baseline to start from")
	}

	// signature-policy-store ::= SEQUENCE { spDocSpec OID 2.2.25.1, spDocument OCTET STRING "p" },
	// added as an unsigned attribute to a signature that has no signature-policy-identifier.
	store := []byte{0x30, 0x08, 0x06, 0x03, 0x52, 0x19, 0x01, 0x04, 0x01, 0x70}
	signerInformation := spi.DSSASN1UtilsFirstSignerInformation(signatureCMS.SignerInfos())
	unsignedAttributes := append(append(cmscore.Attributes{}, cades.UtilsUnsignedAttributes(signerInformation)...),
		cmscore.NewAttribute(spi.OIDIdAaEtsSigPolicyStore, store))
	newSignerInformation, err := cms.UtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributes)
	if err != nil {
		t.Fatal(err)
	}
	withStore, err := cms.UtilsReplaceSigners(signatureCMS, []*cmscore.SignerInfo{newSignerInformation})
	if err != nil {
		t.Fatal(err)
	}

	checker := checkerFor(t, withStore)
	if checker.Signature().SignaturePolicyStore() == nil {
		t.Fatal("test bug: the signature policy store was not picked up")
	}
	if checker.HasBaselineBProfile() {
		t.Fatal("test bug: requirement (k) should reject a store without a signature-policy-identifier hash")
	}
	if !checker.IsValidForPAdESBaselineBProfile() {
		t.Error("IsValidForPAdESBaselineBProfile rejected a CMS that only fails CAdES requirement (k)")
	}
}
