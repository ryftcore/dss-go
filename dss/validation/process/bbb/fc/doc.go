// Package fc ports dss-validation/.../validation/process/bbb/fc (DSS 6.5.RC1):
// the EN 319 102-1 "5.2.2 Format Checking" building block, for signatures,
// timestamps, EAAs and EAA revocation tokens.
//
// # Dependency on validation/process and validation/process/bbb
//
// This package relies on validation/process's Chain/ChainItem framework (see
// chain.go, chain_item.go) and validation/process/bbb's AbstractMultiValuesCheckItem
// and friends (bbb/abstract_multi_values_check_item.go). The key shapes this
// package relies on:
//
//	// process.Result[T] is the Go stand-in for Java's "T extends
//	// XmlConstraintsConclusion" bound: it pairs the concrete jaxb result
//	// (e.g. *drjaxb.XmlFC) with pointers to the two embedded base structs
//	// every ConstraintsConclusion extension carries, since Go generics can't
//	// reach an embedded field through a type parameter.
//	xmlFC := &drjaxb.XmlFC{}
//	result := process.NewResult(xmlFC, &xmlFC.XmlConstraintsConclusionContent, &xmlFC.XmlConstraintsConclusionAttrs)
//
//	// process.ChainItem[T] is the interface chain LINKING operates on
//	// (SetNextItem/Execute); process.ChainItemBase[T] carries the state and
//	// concrete method bodies (the Java "template" side), with defaults for
//	// every ChainItemOverrides member Java doesn't force a subclass to
//	// override. A leaf check embeds *process.ChainItemBase[*drjaxb.XmlFC]
//	// (or *bbb.AbstractMultiValuesCheckItem[*drjaxb.XmlFC], which itself
//	// embeds it), defines only the methods its Java class actually
//	// overrides (Process/FailedIndicationForConclusion/
//	// FailedSubIndicationForConclusion are the true abstracts; everything
//	// else - MessageTag, ErrorMessageTag, BuildErrorMessage,
//	// BuildAdditionalInfo *string, BlockType, SuccessIndication, ... - is
//	// promoted from ChainItemBase unless overridden), and calls
//	// InitChainItem(self) as the last constructor step.
//	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
//	c.InitChainItem(c)
//
//	// process.Chain[T] / process.ChainBase[T] mirror the same pattern one
//	// level up: AbstractFormatChecking embeds *process.ChainBase[*drjaxb.XmlFC],
//	// only InitChain() has no default (Java's sole abstract method on
//	// Chain), and the registration method is InitChainBase (not InitChain -
//	// that name is reserved for the ChainOverrides member a concrete leaf
//	// defines; Go has no overloading, so a same-named zero-arg method on the
//	// leaf would otherwise shadow the promoted registrar).
//	c.ChainBase = process.NewChainBase(i18nProvider, result)
//	...
//	c.InitChainBase(c) // in the leaf's constructor, once InitChain() is defined
//
// A generic check's own "result" parameter is *process.Result[T]; a fixed
// (non-generic) check's is *process.Result[*drjaxb.XmlFC]. bbb.NewAbstractMultiValuesCheckItem
// does NOT call InitChainItem itself - the leaf must, exactly as for a plain
// ChainItemBase embedder (see container_type_check.go for the pattern).
//
// # Unported cross-package dependency: eaa
//
// EAAFormatChecking and EAARevocationFormatChecking wire checks from
// eu.europa.esig.dss.validation.process.eaa.checks and .eaa.status.
// eaa_format_checking.go and eaa_revocation_format_checking.go import
// "github.com/ryftcore/dss-go/dss/validation/process/eaa/checks" and call
// constructors (NewEAASignatureUnicityCheck, NewDisclosurePresentCheck,
// NewDisclosureListExhaustiveCheck, NewKeyBindingSignaturePresentCheck,
// NewEAARevocationTokenTypeCheck), each shaped like every other check
// constructor in this file (i18nProvider, result *process.Result[*drjaxb.XmlFC],
// token, constraint) -> process.ChainItem[*drjaxb.XmlFC].
//
// Both files are gated `//go:build eaa` (see eaa_format_checking.go's header
// for the rationale), so the default, untagged `go build ./...`/`go vet
// ./...`/`go test ./...` for the whole module is green without that package.
// `go build -tags eaa ./...` is green too.
package fc
