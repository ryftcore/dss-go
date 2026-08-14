// Ported from dss-enumerations/.../MimeTypeLoader.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestRegisterMimeTypeLoader_PreservesOrder(t *testing.T) {
	saved := mimeTypeLoaderRegistry
	mimeTypeLoaderRegistry = nil
	t.Cleanup(func() { mimeTypeLoaderRegistry = saved })

	a := &fakeMimeTypeLoader{}
	b := &fakeMimeTypeLoader{}
	RegisterMimeTypeLoader(a)
	RegisterMimeTypeLoader(b)

	loaders := mimeTypeLoaders()
	if len(loaders) != 2 || loaders[0] != MimeTypeLoader(a) || loaders[1] != MimeTypeLoader(b) {
		t.Errorf("mimeTypeLoaders() = %v, want [a, b] in registration order", loaders)
	}
}
