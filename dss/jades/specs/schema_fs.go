// schema_fs.go embeds the JSON schema resources this package's validators reference. There is no
// single Java source: upstream ships these files as Maven resources
// (specs-jades/src/main/resources/schema/, plus the specs-jws RFC 7515/7517 schemas
// specs-jades depends on transitively through the specs-jws Maven dependency - see the notes to
// the porting lead for why rfc7515.json/rfc7517.json/rfc7515-jws.json, which live in the
// specs-jws Maven module rather than specs-jades, are copied in here too).
package specs

import "embed"

//go:embed schema
var schemaFS embed.FS
