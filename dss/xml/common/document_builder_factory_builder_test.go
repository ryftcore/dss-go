package common

import "testing"

func TestDocumentBuilderFactoryBuilder_DefaultsToSecureParseOptions(t *testing.T) {
	opts, err := GetSecureDocumentBuilderFactoryBuilder().Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if opts.AllowDoctype {
		t.Errorf("AllowDoctype = true, want false (DOCTYPE must be disallowed by default)")
	}
}

func TestDocumentBuilderFactoryBuilder_DisablingDisallowDoctypeDeclAllowsDoctype(t *testing.T) {
	b := GetSecureDocumentBuilderFactoryBuilder()
	b.DisableFeature(documentBuilderFeatureDisallowDoctypeDecl)
	opts, err := b.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if !opts.AllowDoctype {
		t.Errorf("AllowDoctype = false, want true after disabling disallow-doctype-decl")
	}
}

func TestDocumentBuilderFactoryBuilder_FluentChainReturnsConcreteType(t *testing.T) {
	b := GetSecureDocumentBuilderFactoryBuilder().
		EnableFeature("custom-feature").
		DisableFeature("another-feature").
		SetAttribute("custom-attribute", "value").
		RemoveAttribute("custom-attribute")
	if b == nil {
		t.Fatalf("fluent chain returned nil")
	}
	if _, err := b.Build(); err != nil {
		t.Fatalf("Build() error = %v, want nil (unrecognized features/attributes are no-ops)", err)
	}
}
