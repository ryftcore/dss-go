// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/FileRevocationSource.java (DSS 6.5.RC1).
//
// Java's DEFAULT_REVOCATION_CACHE_SUBDIRECTORY constant is "/dss-cache-revocation" (a leading
// slash) and is passed as the child of `new File(tmpdir, child)`; the Go port joins the bare
// "dss-cache-revocation" name with the OS temp dir via filepath.Join, which normalizes away the
// leading-slash artifact without changing the resulting path structure.
package spi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/utils"
)

// fileRevocationSourceDefaultCacheSubdirectory is the subdirectory used by default for
// revocation data caching. Port of DEFAULT_REVOCATION_CACHE_SUBDIRECTORY (without its leading
// slash, see the file header).
const fileRevocationSourceDefaultCacheSubdirectory = "dss-cache-revocation"

// fileCacheEntryCertFileExtension is the file extension used to store a revocation data
// issuer certificate. Port of FileCacheEntry.CERT_FILE_EXTENSION.
const fileCacheEntryCertFileExtension = ".cer"

// fileCacheEntryURIFileExtension is the file extension used to define the original revocation
// data URI. Port of FileCacheEntry.URI_FILE_EXTENSION.
const fileCacheEntryURIFileExtension = ".uri"

// FileRevocationSourceOverrides captures what FileRevocationSourceBase needs to reach through
// virtual dispatch: the two operations Java declares abstract on FileRevocationSource<R>
// (reconstructTokenFromEncodedData, getRevocationFileExtension).
type FileRevocationSourceOverrides[R revocation.Revocation] interface {
	// ReconstructTokenFromEncodedData creates a revocation token from cached encoded data, nil
	// if creation fails. Port of the abstract reconstructTokenFromEncodedData(FileCacheEntry,
	// CertificateToken, CertificateToken).
	ReconstructTokenFromEncodedData(revocationCache *FileCacheEntry[R], certificateToken, issuerCertToken *model.CertificateToken) RevocationToken[R]
	// RevocationFileExtension gets the file extension used for cached revocation files (e.g.
	// ".crl" or ".ocsp"). Port of the abstract getRevocationFileExtension().
	RevocationFileExtension() string
}

// FileRevocationSourceBase extends RepositoryRevocationSourceBase to provide file-based caching
// functionality for revocation data. A concrete source embeds it and registers itself with
// InitFileRevocationSource; the outstanding RepositoryRevocationSourceOverrides methods
// (InitRevocationTokenKeys, RevocationAccessURLs, RevocationTokenKey) still need to come from
// that concrete source, since FileRevocationSource itself is abstract on those in Java too.
type FileRevocationSourceBase[R revocation.Revocation] struct {
	RepositoryRevocationSourceBase[R]

	// overrides points back at the concrete source; see InitFileRevocationSource.
	overrides FileRevocationSourceOverrides[R]

	// fileCacheDirectory is the directory where cached revocation files are stored. Default: a
	// "dss-cache-revocation" subdirectory of the OS temp directory.
	fileCacheDirectory string
}

// NewFileRevocationSourceBase builds an empty file revocation source; the proxied revocation
// source can be provided with SetProxySource. Port of the protected default constructor.
func NewFileRevocationSourceBase[R revocation.Revocation]() FileRevocationSourceBase[R] {
	return FileRevocationSourceBase[R]{
		RepositoryRevocationSourceBase: NewRepositoryRevocationSourceBase[R](),
		fileCacheDirectory:             filepath.Join(os.TempDir(), fileRevocationSourceDefaultCacheSubdirectory),
	}
}

// NewFileRevocationSourceBaseWithProxiedSource builds a file revocation source initialized with
// proxiedSource, used to load revocation data when the corresponding document is not available
// in the file system. Port of the protected FileRevocationSource(RevocationSource) constructor.
func NewFileRevocationSourceBaseWithProxiedSource[R revocation.Revocation](proxiedSource RevocationSource[R]) FileRevocationSourceBase[R] {
	base := NewFileRevocationSourceBase[R]()
	base.SetProxySource(proxiedSource)
	return base
}

// InitFileRevocationSource registers the concrete source with its base, and forwards the
// registration to the embedded RepositoryRevocationSourceBase (overrides must also implement
// RepositoryRevocationSourceOverrides - the concrete source's own InitRevocationTokenKeys,
// RevocationAccessURLs and RevocationTokenKey combine with FindRevocations/InsertRevocation/
// UpdateRevocation/RemoveRevocation promoted from this struct). It must be called by the
// outermost concrete source's constructor before the source is used.
func (s *FileRevocationSourceBase[R]) InitFileRevocationSource(overrides FileRevocationSourceOverrides[R]) {
	s.overrides = overrides
	repositoryOverrides, ok := overrides.(RepositoryRevocationSourceOverrides[R])
	if !ok {
		panic("FileRevocationSource was not initialised: the concrete source must also implement " +
			"RepositoryRevocationSourceOverrides (InitRevocationTokenKeys, RevocationAccessURLs, RevocationTokenKey)")
	}
	s.RepositoryRevocationSourceBase.InitRepositoryRevocationSource(repositoryOverrides)
}

// fileRevocationSourceBaseOverrides returns the registered overrides, panicking when the
// concrete source forgot to call InitFileRevocationSource.
func (s *FileRevocationSourceBase[R]) fileRevocationSourceBaseOverrides() FileRevocationSourceOverrides[R] {
	if s.overrides == nil {
		panic("FileRevocationSource was not initialised: the concrete source must call InitFileRevocationSource in its constructor")
	}
	return s.overrides
}

// SetFileCacheDirectory sets the file cache directory, creating it if it does not exist.
// Default: a "dss-cache-revocation" subdirectory of the OS temp directory. Port of
// setFileCacheDirectory(File).
//
// Panics with the Java messages: fileCacheDirectory == "" mirrors Objects.requireNonNull("File
// cache directory cannot be null!"); a directory that fails to be created mirrors the
// IllegalStateException("Unable to create cache directory '%s'"); an existing, non-directory
// path mirrors the IllegalArgumentException("Cache path '%s' is not a directory").
func (s *FileRevocationSourceBase[R]) SetFileCacheDirectory(fileCacheDirectory string) {
	if fileCacheDirectory == "" {
		panic("File cache directory cannot be null!")
	}
	s.fileCacheDirectory = fileRevocationSourceInitializeCacheDirectory(fileCacheDirectory)
}

// fileRevocationSourceInitializeCacheDirectory initializes the cache directory by creating it
// if it doesn't exist. Port of the private initializeCacheDirectory(File).
func fileRevocationSourceInitializeCacheDirectory(fileCacheDirectory string) string {
	info, err := os.Stat(fileCacheDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("Unable to create cache directory '%s'", fileCacheDirectory))
		}
		if mkErr := os.MkdirAll(fileCacheDirectory, 0o755); mkErr != nil {
			panic(fmt.Sprintf("Unable to create cache directory '%s'", fileCacheDirectory))
		}
		// Upstream logs "Cache directory '{}' created successfully"
		return fileCacheDirectory
	}
	if !info.IsDir() {
		panic(fmt.Sprintf("Cache path '%s' is not a directory", fileCacheDirectory))
	}
	return fileCacheDirectory
}

// FileCacheDirectory gets the cache directory. Port of getFileCacheDirectory().
func (s *FileRevocationSourceBase[R]) FileCacheDirectory() string {
	return s.fileCacheDirectory
}

// FindRevocations reconstructs the RevocationToken cached for key, when present. Port of the
// findRevocations(String, CertificateToken, CertificateToken) override.
//
// Java wraps the reconstruction in a try/catch that logs and swallows any Exception; the Go
// port lets a panicking ReconstructTokenFromEncodedData propagate instead, matching the
// "a source that panics is left to propagate" convention used throughout this port (see
// composite_revocation_source.go).
func (s *FileRevocationSourceBase[R]) FindRevocations(key string, certificateToken, issuerCertToken *model.CertificateToken) []RevocationToken[R] {
	revocationCache := s.revocationCache(key)
	if revocationCache.Exists() {
		token := s.fileRevocationSourceBaseOverrides().ReconstructTokenFromEncodedData(revocationCache, certificateToken, issuerCertToken)
		if token != nil {
			return []RevocationToken[R]{token}
		}
		// Upstream logs "Failed to reconstruct revocation token from cache for key: {}"
	}
	return nil
}

// revocationCache builds the FileCacheEntry for key. Port of the private getRevocationCache(String).
func (s *FileRevocationSourceBase[R]) revocationCache(key string) *FileCacheEntry[R] {
	return NewFileCacheEntry[R](s, key, s.fileRevocationSourceBaseOverrides().RevocationFileExtension())
}

// InsertRevocation writes token into the cache file for revocationKey. Port of the
// insertRevocation(String, RevocationToken) override.
func (s *FileRevocationSourceBase[R]) InsertRevocation(revocationKey string, token RevocationToken[R]) {
	revocationCache := s.revocationCache(revocationKey)
	s.saveRevocationToken(revocationCache, token)
	// Upstream logs "Revocation token inserted into cache file for key: {}"
}

// saveRevocationToken writes token within the file system. Port of the protected
// saveRevocationToken(FileCacheEntry, RevocationToken).
//
// Panics on an underlying I/O failure: Java's DSSUtils.saveToFile wraps IOException into an
// unchecked DSSException, and InsertRevocation/UpdateRevocation (RepositoryRevocationSourceOverrides)
// have no error channel to return the Go FileCacheEntry.SaveRevocationToken error through.
func (s *FileRevocationSourceBase[R]) saveRevocationToken(revocationCache *FileCacheEntry[R], token RevocationToken[R]) {
	if err := revocationCache.SaveRevocationToken(token); err != nil {
		panic(err)
	}
}

// UpdateRevocation writes token within the cache file for revocationKey. For a file-based
// cache, update is the same as insert (replace the file content). Port of the
// updateRevocation(String, RevocationToken) override.
func (s *FileRevocationSourceBase[R]) UpdateRevocation(revocationKey string, token RevocationToken[R]) {
	s.InsertRevocation(revocationKey, token)
	// Upstream logs "Revocation token updated in cache file for key: {}"
}

// RemoveRevocation removes the cache entry for revocationKey. Port of the
// removeRevocation(String) override.
func (s *FileRevocationSourceBase[R]) RemoveRevocation(revocationKey string) {
	revocationCache := s.revocationCache(revocationKey)
	if revocationCache.Exists() {
		revocationCache.Clean()
		// Upstream logs whether all associated files were removed or only some/none.
	}
}

// DeleteCacheFile deletes cacheFile, silently doing nothing if it does not exist. Port of the
// protected deleteCacheFile(File).
func (s *FileRevocationSourceBase[R]) DeleteCacheFile(cacheFile string) {
	// Upstream logs a warning on failure; a missing file is not an error (Files.deleteIfExists).
	_ = os.Remove(cacheFile)
}

// ClearCache clears all cached files from the cache directory. Port of clearCache().
func (s *FileRevocationSourceBase[R]) ClearCache() {
	info, err := os.Stat(s.fileCacheDirectory)
	if err != nil {
		// Upstream logs "Cache directory '{}' does not exist!"
		return
	}
	if !info.IsDir() {
		// Upstream logs "Cache directory '{}' is not a directory"
		return
	}
	if err := utils.CleanDirectory(s.fileCacheDirectory); err != nil {
		// Upstream logs "Failed to clean the directory '{}' : {}"
		return
	}
	// Upstream logs "Cache cleared for directory: {}"
}

// FileCacheEntry represents a cache entry related to a single revocation token. Port of the
// non-static inner class FileRevocationSource.FileCacheEntry; Java's implicit outer-class
// reference (used to read fileCacheDirectory, which can change after the entry is built via
// SetFileCacheDirectory) becomes an explicit back-reference to the owning source.
type FileCacheEntry[R revocation.Revocation] struct {
	// source is the owning FileRevocationSourceBase, read live for its fileCacheDirectory.
	source *FileRevocationSourceBase[R]
	// key is the unique identifier of the revocation data (e.g. a normalized URI location).
	key string
	// revocationExtension is the filename extension for the revocation data document.
	revocationExtension string
}

// NewFileCacheEntry builds a cache entry for key, within source's file cache directory,
// documents stored with revocationExtension. Port of the public FileCacheEntry(String, String)
// constructor.
func NewFileCacheEntry[R revocation.Revocation](source *FileRevocationSourceBase[R], key, revocationExtension string) *FileCacheEntry[R] {
	return &FileCacheEntry[R]{source: source, key: key, revocationExtension: revocationExtension}
}

// RevocationDataBinaries gets the revocation data binaries, nil when absent. Port of
// getRevocationDataBinaries().
func (e *FileCacheEntry[R]) RevocationDataBinaries() []byte {
	return fileCacheEntryContent(e.cacheRevocationFile())
}

// RevocationDataSourceURL gets the URL originally used to retrieve the revocation data, "" when
// absent (Java returns null). Port of getRevocationDataSourceUrl().
func (e *FileCacheEntry[R]) RevocationDataSourceURL() string {
	content := fileCacheEntryContent(e.cacheURIFile())
	if content != nil {
		return string(content)
	}
	return ""
}

// IssuerCertificateToken gets a revocation data issuer's certificate, nil when absent in the
// filesystem or unreadable. Port of getIssuerCertificateToken().
func (e *FileCacheEntry[R]) IssuerCertificateToken() *model.CertificateToken {
	issuerCertificateFile := e.cacheRevocationIssuerCertificateFile()
	if _, err := os.Stat(issuerCertificateFile); err != nil {
		return nil
	}
	encodedCertificate := fileCacheEntryContent(issuerCertificateFile)
	if encodedCertificate == nil {
		return nil
	}
	certificateToken, err := DSSUtilsLoadCertificateFromBinary(encodedCertificate)
	if err != nil {
		// Upstream logs "Unable to load revocation data issuer certificate from file with
		// filename '{}' : {}"
		return nil
	}
	return certificateToken
}

// fileCacheEntryContent reads file's content, nil when it is missing or unreadable. Port of
// the private getFileContent(File), with its exists() pre-check folded into DSSUtilsToByteArray's
// own not-found handling.
func fileCacheEntryContent(file string) []byte {
	data, err := DSSUtilsToByteArray(file)
	if err != nil {
		// Upstream logs "The file '{}' does not exist or has been removed." or "Failed to read
		// revocation cache file for key '{}': {}"
		return nil
	}
	return data
}

// cacheRevocationFile gets the cached revocation file. Port of the private getCacheRevocationFile().
func (e *FileCacheEntry[R]) cacheRevocationFile() string {
	return e.cacheFileFromKey(e.revocationExtension)
}

// cacheURIFile gets the cached URI file. Port of the private getCacheUriFile().
func (e *FileCacheEntry[R]) cacheURIFile() string {
	return e.cacheFileFromKey(fileCacheEntryURIFileExtension)
}

// cacheRevocationIssuerCertificateFile gets the cached revocation data issuer certificate file.
// Port of the private getCacheRevocationIssuerCertificateFile().
func (e *FileCacheEntry[R]) cacheRevocationIssuerCertificateFile() string {
	return e.cacheFileFromKey(fileCacheEntryCertFileExtension)
}

// cacheFileFromKey gets the cache file path for the entry's key and the given extension. Port
// of the private getCacheFileFromKey(String).
func (e *FileCacheEntry[R]) cacheFileFromKey(fileExtension string) string {
	return filepath.Join(e.source.fileCacheDirectory, e.key+fileExtension)
}

// SaveRevocationToken writes revocationToken to the corresponding cache document and associated
// documents. Port of saveRevocationToken(RevocationToken).
//
// Panics with the Java message when revocationToken is missing (Objects.requireNonNull). The
// DSSException DSSUtils.saveToFile raises upstream on an I/O failure is returned as an error.
func (e *FileCacheEntry[R]) SaveRevocationToken(revocationToken RevocationToken[R]) error {
	if revocationToken == nil {
		panic("RevocationToken cannot be null!")
	}
	if err := DSSUtilsSaveToFile(revocationToken.Encoded(), e.cacheRevocationFile()); err != nil {
		return err
	}
	if revocationToken.SourceURL() != "" {
		if err := DSSUtilsSaveToFile([]byte(revocationToken.SourceURL()), e.cacheURIFile()); err != nil {
			return err
		}
	}
	return nil
}

// SaveCertificateToken writes certificateToken to the corresponding cache document. Port of
// saveCertificateToken(CertificateToken).
//
// Panics with the Java message when certificateToken is missing (Objects.requireNonNull). The
// DSSException DSSUtils.saveToFile raises upstream on an I/O failure is returned as an error.
func (e *FileCacheEntry[R]) SaveCertificateToken(certificateToken *model.CertificateToken) error {
	if certificateToken == nil {
		panic("CertificateToken cannot be null!")
	}
	return DSSUtilsSaveToFile(certificateToken.Encoded(), e.cacheRevocationIssuerCertificateFile())
}

// Clean cleans all files within the file system associated with the current cache entry, TRUE
// if the cache has been cleaned successfully. Port of clean().
func (e *FileCacheEntry[R]) Clean() bool {
	cacheCleaned := fileCacheEntryRemoveFile(e.cacheRevocationFile())
	cacheURIFile := e.cacheURIFile()
	if _, err := os.Stat(cacheURIFile); err == nil {
		cacheCleaned = cacheCleaned != fileCacheEntryRemoveFile(cacheURIFile)
	}
	cacheCertificateFile := e.cacheRevocationIssuerCertificateFile()
	if _, err := os.Stat(cacheCertificateFile); err == nil {
		cacheCleaned = cacheCleaned != fileCacheEntryRemoveFile(cacheCertificateFile)
	}
	return cacheCleaned
}

// fileCacheEntryRemoveFile removes fileToRemove if it exists, reporting whether it did so. Port
// of the private removeFile(File).
func fileCacheEntryRemoveFile(fileToRemove string) bool {
	if _, err := os.Stat(fileToRemove); err != nil {
		// Upstream logs "Unable to remove the file with filename '{}'! The file does not exist."
		return false
	}
	if err := os.Remove(fileToRemove); err != nil {
		// Upstream logs "Unable to remove the cached file with name '%s'. Reason : %s"
		return false
	}
	return true
}

// Exists checks whether the revocation cache exists. Port of exists().
func (e *FileCacheEntry[R]) Exists() bool {
	_, err := os.Stat(e.cacheRevocationFile())
	return err == nil
}

// compile-time assertion: a FileRevocationSourceBase is a RevocationSource and a
// MultipleRevocationSource (promoted from RepositoryRevocationSourceBase).
var (
	_ RevocationSource[revocation.CRL]         = (*FileRevocationSourceBase[revocation.CRL])(nil)
	_ MultipleRevocationSource[revocation.CRL] = (*FileRevocationSourceBase[revocation.CRL])(nil)
)
