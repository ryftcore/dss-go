// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/resources/TempFileResourcesHandler.java (DSS 6.5.RC1).
//
// DEVIATION: Java's File#deleteOnExit() registers a JVM shutdown hook so the temp file is
// removed even if the process is killed before Close/ForceDelete runs; Go has no portable
// equivalent (os.Exit bypasses deferred cleanup entirely) and this port does not attempt to
// simulate one. Callers are expected to Close (or have WriteToDSSDocument/ForceDelete take
// ownership of) every handler they create, exactly as upstream's own javadoc already asks
// ("do not forget to preserve the output documents...").
package document

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/utain/esig/dss/model"
)

// TempFileResourcesHandler stores temporary documents to a unique file within the filesystem.
// Removes temporary files on Close, unless they have been used within a FileDocument (returned
// by WriteToDSSDocument).
type TempFileResourcesHandler struct {
	AbstractResourcesHandler

	// tempFilePath is the temporary file used for streaming the data.
	tempFilePath string

	// file is the *os.File opened by buildOutputStream, kept to Close/Flush it directly (Go's
	// io.Writer, unlike Java's OutputStream, does not expose Close/Flush itself).
	file *os.File

	// writer buffers writes to file, mirroring Java's BufferedOutputStream wrapping.
	writer *bufio.Writer

	// toBeDeleted indicates whether the File should be removed on calling #Close.
	toBeDeleted bool
}

// NewTempFileResourcesHandler is the default constructor. Port of
// TempFileResourcesHandler(String, String, File): Files.createTempFile(dir, prefix, suffix)
// becomes os.CreateTemp(dir, prefix+"*"+suffix), Go's equivalent placeholder-based naming.
func NewTempFileResourcesHandler(fileNamePrefix, fileNameSuffix, tempFileDirectory string) (*TempFileResourcesHandler, error) {
	f, err := os.CreateTemp(tempFileDirectory, fileNamePrefix+"*"+fileNameSuffix)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create a temporary file. Reason : %s", err.Error()), err)
	}
	path := f.Name()
	if err := f.Close(); err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create a temporary file. Reason : %s", err.Error()), err)
	}

	h := &TempFileResourcesHandler{tempFilePath: path, toBeDeleted: true}
	h.BuildOutputStream = h.buildOutputStream
	return h, nil
}

// buildOutputStream ports the protected #buildOutputStream.
func (h *TempFileResourcesHandler) buildOutputStream() (io.Writer, error) {
	h.assertFileExists()
	f, err := os.OpenFile(h.tempFilePath, os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	h.file = f
	h.writer = bufio.NewWriter(f)
	return h.writer, nil
}

// WriteToDSSDocument ports #writeToDSSDocument.
func (h *TempFileResourcesHandler) WriteToDSSDocument() (model.DSSDocument, error) {
	h.assertFileExists()
	// Avoid deletion of the File on exit.
	h.toBeDeleted = false
	// close to finalize the OutputStream.
	if err := h.Close(); err != nil {
		return nil, err
	}
	// create the document.
	return model.NewFileDocument(h.tempFilePath)
}

// assertFileExists ports the private #assertFileExists.
func (h *TempFileResourcesHandler) assertFileExists() {
	if _, err := os.Stat(h.tempFilePath); err != nil {
		panic(fmt.Sprintf("The file '%s' does not exists!", h.tempFilePath))
	}
}

// Close ports #close: flushes and closes the underlying file (Java's super.close(), applied to
// the buffered OutputStream), then removes the temporary file, unless #toBeDeleted was cleared.
func (h *TempFileResourcesHandler) Close() error {
	var flushErr, closeErr error
	if h.writer != nil {
		flushErr = h.writer.Flush()
	}
	if h.file != nil {
		closeErr = h.file.Close()
		h.file = nil
	}
	if h.toBeDeleted {
		h.ForceDelete()
	}
	if flushErr != nil {
		return flushErr
	}
	return closeErr
}

// ForceDelete is used to delete the temporary File forcibly, even with a flag #toBeDeleted set
// to false. Method should be called responsively and the temp file should be preserved when
// needed (e.g. output of #signDocument()). Port of #forceDelete: the slf4j warning upstream logs
// on a failed removal is dropped (not load-bearing, per PORTING.md).
func (h *TempFileResourcesHandler) ForceDelete() {
	if h.tempFilePath == "" {
		return
	}
	_ = os.Remove(h.tempFilePath)
}

// compile-time interface assertion is provided by temp_file_resources_handler_builder.go, which
// imports the resources package this handler must satisfy.
