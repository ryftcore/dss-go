package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/utain/esig/dss"
)

// cmdSign implements "esig sign".
func cmdSign(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("sign", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage: esig sign <file> -format <format> -level <level> -p12 <keystore> -p12-pass env:VAR [flags]

Signs <file>, producing a new signed document.

  -format string
    	signature format: cades, xades, pades, jades, asice or asics (required)
  -level string
    	baseline level to produce: B, T, LT or LTA (required)
  -p12 string
    	PKCS#12 (.p12/.pfx) key store holding the signing key (required)
  -p12-pass string
    	key store password, as env:VARNAME - the literal env var name, never the password itself (required)
  -tsa string
    	RFC 3161 time-stamp authority URL; required for -level T, LT or LTA
  -digest string
    	digest algorithm: sha256 (default), sha384 or sha512
  -packaging string
    	enveloped, enveloping or detached; defaults to the format's own default (ignored for pades and the asic formats)
  -detached
    	shorthand for -packaging detached
  -asic-format string
    	for -format asice/asics only: the signature format the container carries, cades or xades (default xades)
  -out string
    	output file; defaults to the signed document's own computed name, written to the current directory
`)
	}

	var format, level, p12, p12Pass, tsaURL, digest, packaging, asicFormat, out string
	var detached bool
	fs.StringVar(&format, "format", "", "")
	fs.StringVar(&level, "level", "", "")
	fs.StringVar(&p12, "p12", "", "")
	fs.StringVar(&p12Pass, "p12-pass", "", "")
	fs.StringVar(&tsaURL, "tsa", "", "")
	fs.StringVar(&digest, "digest", "", "")
	fs.StringVar(&packaging, "packaging", "", "")
	fs.BoolVar(&detached, "detached", false, "")
	fs.StringVar(&asicFormat, "asic-format", "", "")
	fs.StringVar(&out, "out", "", "")
	leading, hadLeading, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	file, ok := resolveOnePositional(fs, hadLeading, leading)
	if !ok {
		fmt.Fprintf(stderr, "esig sign: exactly one input file is required\n\n")
		fs.Usage()
		return exitUsage
	}
	if format == "" || level == "" || p12 == "" || p12Pass == "" {
		fmt.Fprintf(stderr, "esig sign: -format, -level, -p12 and -p12-pass are all required\n\n")
		fs.Usage()
		return exitUsage
	}

	sf, err := parseSignFormat(format, asicFormat)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n\n", err)
		fs.Usage()
		return exitUsage
	}
	lvl, err := parseLevel(level)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n\n", err)
		fs.Usage()
		return exitUsage
	}
	digestAlgorithm, err := parseDigestAlgorithm(digest)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n\n", err)
		fs.Usage()
		return exitUsage
	}
	var pkg dss.SignaturePackaging
	switch {
	case detached && packaging != "":
		fmt.Fprintf(stderr, "esig sign: -detached and -packaging are mutually exclusive\n\n")
		fs.Usage()
		return exitUsage
	case detached:
		pkg = dss.PackagingDetached
	case packaging != "":
		p, err := parsePackaging(packaging)
		if err != nil {
			fmt.Fprintf(stderr, "esig sign: %v\n\n", err)
			fs.Usage()
			return exitUsage
		}
		pkg = p
	}

	if lvl.NeedsTimestamp() && tsaURL == "" {
		fmt.Fprintf(stderr, "esig sign: -level %s needs a time-stamp authority; pass -tsa\n\n", level)
		fs.Usage()
		return exitUsage
	}

	password, err := resolvePassword(p12Pass)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: -p12-pass: %v\n", err)
		return exitRuntime
	}

	signer, err := dss.OpenPKCS12(p12, password)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n", err)
		return exitRuntime
	}
	defer signer.Close()

	doc, err := dss.OpenDocument(file)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n", err)
		return exitRuntime
	}

	opts := dss.SignOptions{
		Format:          sf.format,
		Level:           lvl,
		DigestAlgorithm: digestAlgorithm,
		Packaging:       pkg,
		ContainerType:   sf.container,
	}
	if lvl.NeedsTimestamp() {
		opts.TSPSource = newHTTPTSPSource(tsaURL)
	}

	signed, err := dss.Sign(doc, signer, opts)
	if err != nil {
		fmt.Fprintf(stderr, "esig sign: %v\n", err)
		return exitRuntime
	}

	outPath := out
	if outPath == "" {
		outPath = signed.Name()
	}
	if err := signed.Save(outPath); err != nil {
		fmt.Fprintf(stderr, "esig sign: writing %s: %v\n", outPath, err)
		return exitRuntime
	}
	fmt.Fprintf(stdout, "%s\n", outPath)
	return exitOK
}

// parsePackaging maps the CLI's -packaging value onto [dss.SignaturePackaging].
func parsePackaging(name string) (dss.SignaturePackaging, error) {
	switch name {
	case "enveloped":
		return dss.PackagingEnveloped, nil
	case "enveloping":
		return dss.PackagingEnveloping, nil
	case "detached":
		return dss.PackagingDetached, nil
	default:
		return "", fmt.Errorf("unsupported -packaging %q (want enveloped, enveloping or detached)", name)
	}
}
