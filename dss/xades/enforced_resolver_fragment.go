// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/EnforcedResolverFragment.java (DSS 6.5.RC1).
//
// NOT reimplemented here: internal/xmldsig/doc.go's Santuario-replacement mapping table lists
// this exact Java class ("EnforcedResolverFragment  dss-xades EnforcedResolverFragment") as
// already ported to internal/xmldsig/resolver.go's xmldsig.EnforcedResolverFragment - the
// org.apache.xml.security.utils.resolver.implementations.ResolverFragment subclass this Java
// file itself extends is xmldsig.ResolverFragment there, and the XPath-injection guard
// (checkValueForXpathInjection, the XPATH_CHAR_FILTER constant, and the engineCanResolveURI
// override) is reproduced verbatim as xmldsig.EnforcedResolverFragment.CanResolve /
// checkValueForXPathInjection / xpathCharFilter. xmldsig.DefaultResolvers() already returns it
// alongside ResolverXPointer (per internal/xmldsig's doc.go), which is how every landed XAdES
// signature-building/validating call site gets its XPath-injection protection - see
// xades_signature.go's file header, which documents this instead of assuming a
// package-`xades`-local EnforcedResolverFragment forward dependency.
//
// PORTING.md forbids editing frozen packages, and internal/xmldsig is frozen; declaring a second,
// unused eu.europa.esig.dss.xades.EnforcedResolverFragment type in package xades here would only
// shadow-duplicate that already-landed, already-wired implementation with dead code, so this file
// intentionally carries no Go declarations of its own - it exists to preserve the "one Go file
// per Java class" mapping and to record where the port actually lives, per S4D_BRIEF.md's
// instruction to flag frozen-package situations for the integrator instead of duplicating them.
package xades
