package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ryftcore/dss-go/dss"
)

// stringList collects a repeatable flag (e.g. -detached, -trust) into a
// slice, in the order given on the command line. Implements flag.Value.
type stringList []string

func (l *stringList) String() string {
	if l == nil {
		return ""
	}
	return strings.Join(*l, ",")
}

func (l *stringList) Set(v string) error {
	*l = append(*l, v)
	return nil
}

// resolvePassword resolves a password flag's value (-p12-pass, -pdf-pass).
// Only the "env:VAR" form is accepted - never a literal password - so a
// password given on the command line does not end up in the shell history
// or in the process list any longer than it takes to read the environment
// variable. example is the variable name the error messages suggest for
// the flag at hand.
//
// The rejection message never quotes spec back. Typing the password where
// env:VARNAME belongs is the very mistake this form exists to catch, so spec
// is a secret exactly when the error fires, and every caller prints the error
// to stderr - which a CI log or a scrollback keeps at least as durably as the
// shell history and the process list this flag is guarding against.
func resolvePassword(spec, example string) (string, error) {
	name, ok := strings.CutPrefix(spec, "env:")
	if !ok {
		return "", fmt.Errorf("password must be given as env:VARNAME (e.g. env:%s), not the password itself", example)
	}
	if name == "" {
		return "", fmt.Errorf("env: form needs a variable name, e.g. env:%s", example)
	}
	value, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", name)
	}
	return value, nil
}

// resolveOptionalPassword is [resolvePassword] for a flag that need not be
// given at all (-pdf-pass, which only an encrypted PDF needs): an empty spec
// means "no password" and resolves to nil, which is what the facade's
// PasswordProtection fields take for a document that is not encrypted.
func resolveOptionalPassword(spec string) ([]byte, error) {
	if spec == "" {
		return nil, nil
	}
	value, err := resolvePassword(spec, "PDF_PASSWORD")
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

// parseDigestAlgorithm maps the CLI's lower-case digest names onto the
// facade's [dss.DigestAlgorithm] constants. An empty name returns the
// facade's own default (SHA-256).
func parseDigestAlgorithm(name string) (dss.DigestAlgorithm, error) {
	switch strings.ToLower(name) {
	case "":
		return "", nil
	case "sha256":
		return dss.DigestSHA256, nil
	case "sha384":
		return dss.DigestSHA384, nil
	case "sha512":
		return dss.DigestSHA512, nil
	default:
		return "", fmt.Errorf("unsupported -digest %q (want sha256, sha384 or sha512)", name)
	}
}

// parseLevel maps the CLI's level names onto [dss.Level], accepting the
// bare letters case-insensitively (b, t, lt, lta).
func parseLevel(name string) (dss.Level, error) {
	switch strings.ToUpper(name) {
	case "B":
		return dss.LevelB, nil
	case "T":
		return dss.LevelT, nil
	case "LT":
		return dss.LevelLT, nil
	case "LTA":
		return dss.LevelLTA, nil
	default:
		return "", fmt.Errorf("unsupported -level %q (want B, T, LT or LTA)", name)
	}
}

// signFormat is the outcome of parsing "sign"'s -format value: the facade
// Format to use, plus - for the two ASiC formats only - the container
// flavour (ASiC-S or ASiC-E) the format name itself selects.
type signFormat struct {
	format    dss.Format
	container dss.ContainerType
}

// parseSignFormat maps the CLI's short format names onto a [dss.Format].
// "asice" and "asics" select an ASiC container; which signature format it
// holds (CAdES or XAdES) is a separate choice the facade itself requires
// (SignOptions.Format), so it is read from asicFormat here (the -asic-format
// flag), defaulting to XAdES, the more common choice for containers holding
// more than one document.
func parseSignFormat(name, asicFormat string) (signFormat, error) {
	switch strings.ToLower(name) {
	case "cades":
		return signFormat{format: dss.FormatCAdES}, nil
	case "xades":
		return signFormat{format: dss.FormatXAdES}, nil
	case "pades":
		return signFormat{format: dss.FormatPAdES}, nil
	case "jades":
		return signFormat{format: dss.FormatJAdES}, nil
	case "asice", "asics":
		container := dss.ContainerASiCE
		if strings.ToLower(name) == "asics" {
			container = dss.ContainerASiCS
		}
		switch strings.ToLower(asicFormat) {
		case "", "xades":
			return signFormat{format: dss.FormatASiCWithXAdES, container: container}, nil
		case "cades":
			return signFormat{format: dss.FormatASiCWithCAdES, container: container}, nil
		default:
			return signFormat{}, fmt.Errorf("unsupported -asic-format %q (want cades or xades)", asicFormat)
		}
	default:
		return signFormat{}, fmt.Errorf("unsupported -format %q (want cades, xades, pades, jades, asice or asics)", name)
	}
}

// parseExtendFormat maps the CLI's short format names onto a [dss.Format]
// for "extend". Extension has no container-type choice to make (the
// container the document already is stays what it is), so unlike
// [parseSignFormat] it returns a bare [dss.Format]; "asice" and "asics" are
// accepted as synonyms of each other here, since asicFormat (-asic-format)
// alone decides between the two ASiC signature formats.
func parseExtendFormat(name, asicFormat string) (dss.Format, error) {
	switch strings.ToLower(name) {
	case "cades":
		return dss.FormatCAdES, nil
	case "xades":
		return dss.FormatXAdES, nil
	case "pades":
		return dss.FormatPAdES, nil
	case "jades":
		return dss.FormatJAdES, nil
	case "asice", "asics":
		switch strings.ToLower(asicFormat) {
		case "", "xades":
			return dss.FormatASiCWithXAdES, nil
		case "cades":
			return dss.FormatASiCWithCAdES, nil
		default:
			return "", fmt.Errorf("unsupported -asic-format %q (want cades or xades)", asicFormat)
		}
	default:
		return "", fmt.Errorf("unsupported -format %q (want cades, xades, pades, jades, asice or asics)", name)
	}
}

// splitPositional pops a single leading positional argument - the "<file>
// [flags]" form every subcommand's usage text documents - off args, so the
// remainder can be handed to a flag.FlagSet's Parse: the standard library's
// flag package stops parsing at the first argument that does not start with
// "-", so a positional argument ahead of the flags has to be removed by hand
// first, or every flag after it would be left unparsed.
func splitPositional(args []string) (leading string, hadLeading bool, rest []string) {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return args[0], true, args[1:]
	}
	return "", false, args
}

// resolveOnePositional resolves the single positional argument a subcommand
// takes, once fs.Parse(rest) has run on the rest [splitPositional] returned.
// When splitPositional popped a leading argument, that is it - provided
// nothing unparsed is left over (fs.NArg() == 0); a flag.FlagSet also
// accepts positional arguments trailing the flags, so when nothing was
// popped (the arguments started with a flag), fs's own single remaining
// argument is used instead. Either way exactly one positional argument is
// required; ok is false for zero or more than one.
func resolveOnePositional(fs *flag.FlagSet, hadLeading bool, leading string) (positional string, ok bool) {
	if hadLeading {
		return leading, fs.NArg() == 0
	}
	if fs.NArg() == 1 {
		return fs.Arg(0), true
	}
	return "", false
}

// loadDocuments opens each path in paths as a [dss.Document] with
// [dss.OpenDocument].
func loadDocuments(paths []string) ([]dss.Document, error) {
	docs := make([]dss.Document, 0, len(paths))
	for _, path := range paths {
		doc, err := dss.OpenDocument(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", path, err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// loadTrustAnchors reads each path in paths as a DER or PEM certificate with
// [dss.LoadCertificate]. Only individual certificate files are accepted -
// not PKCS#12 key stores, which hold private keys, not trust anchors.
func loadTrustAnchors(paths []string) ([]*dss.CertificateToken, error) {
	certs := make([]*dss.CertificateToken, 0, len(paths))
	for _, path := range paths {
		cert, err := dss.LoadCertificate(path)
		if err != nil {
			return nil, fmt.Errorf("loading trust anchor %s: %w", path, err)
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
