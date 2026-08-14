// KAT test for manifest_attribute.go, manifest_element.go, manifest_namespace.go and
// manifest_path.go: every value below is dumped verbatim from a Java oracle run over
// dss-asic-xades 6.5.RC1 + dss-xml-common 6.5.RC1 (java.lang.reflect over
// ManifestAttribute/ManifestElement.getEnumConstants(), ManifestNamespace.NS's fields, and
// ManifestPath.FILE_ENTRY_PATH.getQueryString()/getFullPathAttribute(NS)/getMediaTypeAttribute(NS)),
// per PORTING.md's "exhaustive table test" rule for registry-like tables and this phase's
// explicit KAT requirement for definition enums (S7_BRIEF.md).
package xades

import (
	"testing"

	"github.com/utain/esig/dss/xml/common"
)

func TestManifestNamespace_KAT(t *testing.T) {
	if got := ManifestNS.Uri(); got != "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" {
		t.Errorf("ManifestNS.Uri() = %q, want %q", got, "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0")
	}
	if got := ManifestNS.Prefix(); got != "manifest" {
		t.Errorf("ManifestNS.Prefix() = %q, want %q", got, "manifest")
	}
}

func TestManifestAttribute_KAT(t *testing.T) {
	cases := []struct {
		attr ManifestAttribute
		name string
	}{
		{ManifestAttribute_VERSION, "version"},
		{ManifestAttribute_FULL_PATH, "full-path"},
		{ManifestAttribute_MEDIA_TYPE, "media-type"},
	}
	for _, c := range cases {
		if got := c.attr.AttributeName(); got != c.name {
			t.Errorf("%s.AttributeName() = %q, want %q", c.attr, got, c.name)
		}
	}
}

func TestManifestElement_KAT(t *testing.T) {
	cases := []struct {
		elem   ManifestElement
		tag    string
		uri    string
		prefix string
	}{
		{ManifestElement_MANIFEST, "manifest", "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0", "manifest"},
		{ManifestElement_FILE_ENTRY, "file-entry", "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0", "manifest"},
	}
	for _, c := range cases {
		if got := c.elem.TagName(); got != c.tag {
			t.Errorf("%s.TagName() = %q, want %q", c.elem, got, c.tag)
		}
		if got := c.elem.URI(); got != c.uri {
			t.Errorf("%s.URI() = %q, want %q", c.elem, got, c.uri)
		}
		if got := c.elem.Namespace().Prefix(); got != c.prefix {
			t.Errorf("%s.Namespace().Prefix() = %q, want %q", c.elem, got, c.prefix)
		}
		if !c.elem.IsSameTagName(c.tag) {
			t.Errorf("%s.IsSameTagName(%q) = false, want true", c.elem, c.tag)
		}
		if c.elem.IsSameTagName("not-a-tag") {
			t.Errorf("%s.IsSameTagName(%q) = true, want false", c.elem, "not-a-tag")
		}
	}
}

func TestManifestPath_KAT(t *testing.T) {
	// FILE_ENTRY_PATH.getQueryString() = "./manifest:manifest/manifest:file-entry"
	if got := ManifestPath_FILE_ENTRY_PATH.QueryString(); got != "./manifest:manifest/manifest:file-entry" {
		t.Errorf("ManifestPath_FILE_ENTRY_PATH.QueryString() = %q, want %q",
			got, "./manifest:manifest/manifest:file-entry")
	}
	if ManifestPath_FILE_ENTRY_PATH.IsAll() {
		t.Error("ManifestPath_FILE_ENTRY_PATH.IsAll() = true, want false")
	}
	if !ManifestPath_FILE_ENTRY_PATH.IsFromCurrentPosition() {
		t.Error("ManifestPath_FILE_ENTRY_PATH.IsFromCurrentPosition() = false, want true")
	}

	// getFullPathAttribute(NS) = "manifest:full-path"
	if got := ManifestPathGetFullPathAttribute(ManifestNS); got != "manifest:full-path" {
		t.Errorf("ManifestPathGetFullPathAttribute(ManifestNS) = %q, want %q", got, "manifest:full-path")
	}
	// getMediaTypeAttribute(NS) = "manifest:media-type"
	if got := ManifestPathGetMediaTypeAttribute(ManifestNS); got != "manifest:media-type" {
		t.Errorf("ManifestPathGetMediaTypeAttribute(ManifestNS) = %q, want %q", got, "manifest:media-type")
	}
}

func TestManifestPath_GetFullPathAttribute_NoPrefix(t *testing.T) {
	// A namespace with an empty prefix must not be prefixed - verifies
	// manifestPathAddPrefixIfNeeded's Utils.isStringEmpty(prefix) branch directly, since the
	// oracle's classpath lacked a usable Utils.isStringEmpty implementation (Apache
	// Commons/Guava utils jar wiring) to confirm this branch at runtime; the logic itself
	// (return attributeName unchanged when the prefix is empty) is a direct, unambiguous read
	// of ManifestPath.addPrefixIfNeeded's source.
	noPrefix := common.NewDSSNamespace("urn:x", "")
	if got := manifestPathAddPrefixIfNeeded("full-path", noPrefix); got != "full-path" {
		t.Errorf("addPrefixIfNeeded with empty prefix = %q, want %q", got, "full-path")
	}
}
