package enumerations

import "testing"

func TestCOSESignatureTypeFields(t *testing.T) {
	tests := []struct {
		v                   COSESignatureType
		context             string
		label               string
		tag                 int64
		hasTag              bool
		counterSig          bool
		counterSigV2        bool
		counterSigHdrKey    int64
		hasCounterSigHdrKey bool
	}{
		{COSESignatureTypeCoseSign, "Signature", "COSE_Sign", 98, true, false, false, 0, false},
		{COSESignatureTypeCoseSign1, "Signature1", "COSE_Sign1", 18, true, false, false, 0, false},
		{COSESignatureTypeCoseSignature, "", "", 0, false, false, false, 0, false},
		{COSESignatureTypeCoseCounterSignature, "CounterSignature", "COSE_Countersignature", 19, true, true, false, 7, true},
		{COSESignatureTypeCoseCounterSignature0, "CounterSignature0", "COSE_Countersignature0", 0, false, true, false, 9, true},
		{COSESignatureTypeCoseCounterSignatureV2, "CounterSignatureV2", "COSE_Countersignature_V2", 19, true, true, true, 11, true},
		{COSESignatureTypeCoseCounterSignature0V2, "CounterSignature0V2", "COSE_Countersignature0_V2", 0, false, true, true, 12, true},
	}
	for _, tt := range tests {
		if got := tt.v.Context(); got != tt.context {
			t.Errorf("%v.Context() = %q, want %q", tt.v, got, tt.context)
		}
		if got := tt.v.Label(); got != tt.label {
			t.Errorf("%v.Label() = %q, want %q", tt.v, got, tt.label)
		}
		if tt.hasTag {
			if got := tt.v.Tag(); got != tt.tag {
				t.Errorf("%v.Tag() = %d, want %d", tt.v, got, tt.tag)
			}
		}
		if got := tt.v.IsCounterSignature(); got != tt.counterSig {
			t.Errorf("%v.IsCounterSignature() = %v, want %v", tt.v, got, tt.counterSig)
		}
		if got := tt.v.IsCounterSignatureV2(); got != tt.counterSigV2 {
			t.Errorf("%v.IsCounterSignatureV2() = %v, want %v", tt.v, got, tt.counterSigV2)
		}
		key, ok := tt.v.CounterSignatureHeaderKey()
		if ok != tt.hasCounterSigHdrKey {
			t.Errorf("%v.CounterSignatureHeaderKey() ok = %v, want %v", tt.v, ok, tt.hasCounterSigHdrKey)
		}
		if ok && key != tt.counterSigHdrKey {
			t.Errorf("%v.CounterSignatureHeaderKey() = %d, want %d", tt.v, key, tt.counterSigHdrKey)
		}
	}
}

func TestCOSESignatureTypeTagPanicsWhenUnavailable(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing tag, got none")
		}
	}()
	COSESignatureTypeCoseCounterSignature0.Tag()
}

func TestCOSESignatureTypeForLabel(t *testing.T) {
	for _, v := range COSESignatureTypeValues() {
		label := v.Label()
		if label == "" {
			continue
		}
		if got := COSESignatureTypeForLabel(label); got != v {
			t.Errorf("COSESignatureTypeForLabel(%q) = %q, want %q", label, got, v)
		}
	}
	if got := COSESignatureTypeForLabel("does-not-exist"); got != "" {
		t.Errorf("COSESignatureTypeForLabel(unknown) = %q, want empty", got)
	}
}

func TestCOSESignatureTypeCounterSignatureContextByHeaderKey(t *testing.T) {
	tests := []struct {
		key  int64
		want COSESignatureType
	}{
		{7, COSESignatureTypeCoseCounterSignature},
		{9, COSESignatureTypeCoseCounterSignature0},
		{11, COSESignatureTypeCoseCounterSignatureV2},
		{12, COSESignatureTypeCoseCounterSignature0V2},
		{999, ""},
	}
	for _, tt := range tests {
		if got := COSESignatureTypeCounterSignatureContextByHeaderKey(tt.key); got != tt.want {
			t.Errorf("COSESignatureTypeCounterSignatureContextByHeaderKey(%d) = %q, want %q", tt.key, got, tt.want)
		}
	}
}
