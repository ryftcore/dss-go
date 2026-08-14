package common

import "testing"

func TestValidatorConfigurator_ConfigureWiresErrorHandler(t *testing.T) {
	c := GetSecureValidatorConfigurator()
	v := &Validator{}
	if err := c.Configure(v); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if v.ErrorHandler == nil {
		t.Fatalf("Configure() did not wire an ErrorHandler")
	}
	if !v.ErrorHandler.IsValid() {
		t.Errorf("fresh ErrorHandler should be valid")
	}
}

func TestValidatorConfigurator_ConfigurePanicsOnNilValidator(t *testing.T) {
	c := GetSecureValidatorConfigurator()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil validator")
		}
	}()
	_ = c.Configure(nil)
}

func TestValidatorConfigurator_PostProcessAlertsOnInvalidHandler(t *testing.T) {
	c := GetSecureValidatorConfigurator()
	v := &Validator{ErrorHandler: NewDSSErrorHandler()}
	v.ErrorHandler.RecordError(&XMLParseError{Message: "bad"})

	err := c.PostProcess(v)
	if err == nil {
		t.Fatalf("PostProcess() error = nil, want an alert error for an invalid handler")
	}
}

func TestXmlDefinerUtils_Singleton(t *testing.T) {
	a := GetXmlDefinerUtilsInstance()
	b := GetXmlDefinerUtilsInstance()
	if a != b {
		t.Errorf("GetXmlDefinerUtilsInstance() did not return the same singleton instance")
	}
	opts, err := a.GetSecureDocumentBuilderFactory()
	if err != nil {
		t.Fatalf("GetSecureDocumentBuilderFactory() error = %v", err)
	}
	if opts.AllowDoctype {
		t.Errorf("AllowDoctype = true, want false")
	}
}

func TestXPathQueryBuilderFromXPathQuery_RoundTrips(t *testing.T) {
	original := FromCurrentPositionAttribute(XMLDSigElement_SIGNATURE, XMLDSigAttribute_ID)
	rebuilt := XPathQueryBuilderFromXPathQuery(original).Build()
	if got, want := rebuilt.QueryString(), original.QueryString(); got != want {
		t.Errorf("rebuilt QueryString() = %q, want %q", got, want)
	}
}

func TestXPathQueryBuilderFromXPathQuery_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil XPathQuery")
		}
	}()
	XPathQueryBuilderFromXPathQuery(nil)
}
