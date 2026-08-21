package jaxb

import (
	"os"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// adapterLexical maps the name of an enumeration bound by a generated adapter to
// the round-trip the corresponding Go adapter type performs: the lexical form it
// marshals a constant to, and the constant an unmarshalled lexical form yields.
var adapterLexical = map[string]struct {
	print func(name string) (string, error)
	parse func(lexical string) (string, error)
}{
	"ASiCContainerType": {
		print: func(name string) (string, error) {
			b, err := ASiCContainerTypeValue(enumerations.ASiCContainerType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v ASiCContainerTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.ASiCContainerType(v)), err
		},
	},
	"ArchiveTimestampHashIndexVersion": {
		print: func(name string) (string, error) {
			b, err := ArchiveTimestampHashIndexVersionValue(enumerations.ArchiveTimestampHashIndexVersion(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v ArchiveTimestampHashIndexVersionValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.ArchiveTimestampHashIndexVersion(v)), err
		},
	},
	"ArchiveTimestampType": {
		print: func(name string) (string, error) {
			b, err := ArchiveTimestampTypeValue(enumerations.ArchiveTimestampType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v ArchiveTimestampTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.ArchiveTimestampType(v)), err
		},
	},
	"COSESignatureType": {
		print: func(name string) (string, error) {
			b, err := COSESignatureTypeValue(enumerations.COSESignatureType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v COSESignatureTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.COSESignatureType(v)), err
		},
	},
	"CertificateOrigin": {
		print: func(name string) (string, error) {
			b, err := CertificateOriginValue(enumerations.CertificateOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v CertificateOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.CertificateOrigin(v)), err
		},
	},
	"CertificateRefOrigin": {
		print: func(name string) (string, error) {
			b, err := CertificateRefOriginValue(enumerations.CertificateRefOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v CertificateRefOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.CertificateRefOrigin(v)), err
		},
	},
	"CertificateSourceType": {
		print: func(name string) (string, error) {
			b, err := CertificateSourceTypeValue(enumerations.CertificateSourceType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v CertificateSourceTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.CertificateSourceType(v)), err
		},
	},
	"CertificateStatus": {
		print: func(name string) (string, error) {
			b, err := CertificateStatusValue(enumerations.CertificateStatus(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v CertificateStatusValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.CertificateStatus(v)), err
		},
	},
	"CertificationPermission": {
		print: func(name string) (string, error) {
			b, err := CertificationPermissionValue(enumerations.CertificationPermission(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v CertificationPermissionValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.CertificationPermission(v)), err
		},
	},
	"DigestAlgorithm": {
		print: func(name string) (string, error) {
			b, err := DigestAlgorithmValue(enumerations.DigestAlgorithm(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v DigestAlgorithmValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.DigestAlgorithm(v)), err
		},
	},
	"DigestMatcherType": {
		print: func(name string) (string, error) {
			b, err := DigestMatcherTypeValue(enumerations.DigestMatcherType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v DigestMatcherTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.DigestMatcherType(v)), err
		},
	},
	"EAAPresentationType": {
		print: func(name string) (string, error) {
			b, err := EAAPresentationTypeValue(enumerations.EAAPresentationType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EAAPresentationTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EAAPresentationType(v)), err
		},
	},
	"EAARevocationOrigin": {
		print: func(name string) (string, error) {
			b, err := EAARevocationOriginValue(enumerations.EAARevocationOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EAARevocationOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EAARevocationOrigin(v)), err
		},
	},
	"EAAStatus": {
		print: func(name string) (string, error) {
			b, err := EAAStatusValue(enumerations.EAAStatus(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EAAStatusValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EAAStatus(v)), err
		},
	},
	"EAAType": {
		print: func(name string) (string, error) {
			b, err := EAATypeValue(enumerations.EAAType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EAATypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EAAType(v)), err
		},
	},
	"EncryptionAlgorithm": {
		print: func(name string) (string, error) {
			b, err := EncryptionAlgorithmValue(enumerations.EncryptionAlgorithm(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EncryptionAlgorithmValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EncryptionAlgorithm(v)), err
		},
	},
	"EndorsementType": {
		print: func(name string) (string, error) {
			b, err := EndorsementTypeValue(enumerations.EndorsementType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EndorsementTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EndorsementType(v)), err
		},
	},
	"EvidenceRecordIncorporationType": {
		print: func(name string) (string, error) {
			b, err := EvidenceRecordIncorporationTypeValue(enumerations.EvidenceRecordIncorporationType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EvidenceRecordIncorporationTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EvidenceRecordIncorporationType(v)), err
		},
	},
	"EvidenceRecordOrigin": {
		print: func(name string) (string, error) {
			b, err := EvidenceRecordOriginValue(enumerations.EvidenceRecordOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EvidenceRecordOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EvidenceRecordOrigin(v)), err
		},
	},
	"EvidenceRecordTimestampType": {
		print: func(name string) (string, error) {
			b, err := EvidenceRecordTimestampTypeValue(enumerations.EvidenceRecordTimestampType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EvidenceRecordTimestampTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EvidenceRecordTimestampType(v)), err
		},
	},
	"EvidenceRecordTypeEnum": {
		print: func(name string) (string, error) {
			b, err := EvidenceRecordTypeEnumValue(enumerations.EvidenceRecordTypeEnum(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v EvidenceRecordTypeEnumValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.EvidenceRecordTypeEnum(v)), err
		},
	},
	"GeneralNameType": {
		print: func(name string) (string, error) {
			b, err := GeneralNameTypeValue(enumerations.GeneralNameType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v GeneralNameTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.GeneralNameType(v)), err
		},
	},
	"JWSSerializationType": {
		print: func(name string) (string, error) {
			b, err := JWSSerializationTypeValue(enumerations.JWSSerializationType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v JWSSerializationTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.JWSSerializationType(v)), err
		},
	},
	"KeyUsageBit": {
		print: func(name string) (string, error) {
			b, err := KeyUsageBitValue(enumerations.KeyUsageBit(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v KeyUsageBitValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.KeyUsageBit(v)), err
		},
	},
	"PdfLockAction": {
		print: func(name string) (string, error) {
			b, err := PdfLockActionValue(enumerations.PdfLockAction(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v PdfLockActionValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.PdfLockAction(v)), err
		},
	},
	"PdfObjectModificationType": {
		print: func(name string) (string, error) {
			b, err := PdfObjectModificationTypeValue(enumerations.PdfObjectModificationType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v PdfObjectModificationTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.PdfObjectModificationType(v)), err
		},
	},
	"RevocationOrigin": {
		print: func(name string) (string, error) {
			b, err := RevocationOriginValue(enumerations.RevocationOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v RevocationOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.RevocationOrigin(v)), err
		},
	},
	"RevocationReason": {
		print: func(name string) (string, error) {
			b, err := RevocationReasonValue(enumerations.RevocationReason(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v RevocationReasonValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.RevocationReason(v)), err
		},
	},
	"RevocationRefOrigin": {
		print: func(name string) (string, error) {
			b, err := RevocationRefOriginValue(enumerations.RevocationRefOrigin(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v RevocationRefOriginValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.RevocationRefOrigin(v)), err
		},
	},
	"RevocationType": {
		print: func(name string) (string, error) {
			b, err := RevocationTypeValue(enumerations.RevocationType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v RevocationTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.RevocationType(v)), err
		},
	},
	"SignatureLevel": {
		print: func(name string) (string, error) {
			b, err := SignatureLevelValue(enumerations.SignatureLevel(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v SignatureLevelValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.SignatureLevel(v)), err
		},
	},
	"SignatureScopeType": {
		print: func(name string) (string, error) {
			b, err := SignatureScopeTypeValue(enumerations.SignatureScopeType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v SignatureScopeTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.SignatureScopeType(v)), err
		},
	},
	"TimestampType": {
		print: func(name string) (string, error) {
			b, err := TimestampTypeValue(enumerations.TimestampType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v TimestampTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.TimestampType(v)), err
		},
	},
	"TimestampedObjectType": {
		print: func(name string) (string, error) {
			b, err := TimestampedObjectTypeValue(enumerations.TimestampedObjectType(name)).MarshalText()
			return string(b), err
		},
		parse: func(lexical string) (string, error) {
			var v TimestampedObjectTypeValue
			err := v.UnmarshalText([]byte(lexical))
			return string(enumerations.TimestampedObjectType(v)), err
		},
	},
}

// TestAdapterLexicalForms is the exhaustive table test of the generated JAXB
// adapters: for every constant of every enumeration the model binds, the Go
// adapter must print the lexical form the Java adapter prints, and must read it
// back to the same constant.
func TestAdapterLexicalForms(t *testing.T) {
	data, err := os.ReadFile(corpustest.Path(t, "adapters.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) < 300 {
		t.Fatalf("oracle holds only %d rows", len(lines))
	}
	covered := map[string]bool{}
	for _, line := range lines {
		cols := strings.Split(line, "\t")
		if len(cols) != 4 {
			t.Fatalf("malformed oracle row %q", line)
		}
		adapter, enum, constant, lexical := cols[0], cols[1], cols[2], cols[3]
		codec, ok := adapterLexical[enum]
		if !ok {
			t.Errorf("%s: no Go adapter for enumeration %s", adapter, enum)
			continue
		}
		covered[enum] = true
		if lexical == `\N` {
			// The Java adapter prints null for this constant, so JAXB leaves the
			// property out of the document altogether. encoding/xml has no such
			// spelling: the Go adapter writes the empty string, and the
			// deviation is documented in jaxb_adapters.go.
			got, err := codec.print(constant)
			if err != nil {
				t.Errorf("%s.%s: marshal: %v", enum, constant, err)
			} else if got != "" {
				t.Errorf("%s.%s: marshalled %q, the Java adapter has no lexical form for it", enum, constant, got)
			}
			continue
		}
		got, err := codec.print(constant)
		if err != nil {
			t.Errorf("%s.%s: marshal: %v", enum, constant, err)
			continue
		}
		if got != lexical {
			t.Errorf("%s.%s: marshalled %q, Java adapter prints %q", enum, constant, got, lexical)
		}
		back, err := codec.parse(lexical)
		if err != nil {
			t.Errorf("%s.%s: unmarshal %q: %v", enum, constant, lexical, err)
			continue
		}
		if back != constant {
			t.Errorf("%s: %q unmarshalled to %s, want %s", enum, lexical, back, constant)
		}
	}
	if len(covered) != len(adapterLexical) {
		for enum := range adapterLexical {
			if !covered[enum] {
				t.Errorf("%s is bound by the model but absent from the oracle", enum)
			}
		}
	}
}
