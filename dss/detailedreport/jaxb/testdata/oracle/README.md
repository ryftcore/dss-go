# Marshal-parity oracle corpus

Every `*.xml` file in this directory is byte-for-byte JAXB reference
implementation output, captured by re-marshalling
`dss-detailed-report-jaxb`'s own test fixtures
(`src/test/resources/dr*.xml`) through `DetailedReportFacade` -- see
`../ReserializeDetailedReport.java`.

## Marshal through the facade, never straight to an OutputStream

The facade is not an incidental detail. `DetailedReportFacade.marshall(jaxb,
validate)` is the exact call `AbstractReports.getXmlDetailedReport()` makes,
and it marshals into a `StringWriter`. Handing the same tree to
`Marshaller.marshal(element, OutputStream)` instead selects a different RI
output backend (`IndentingUTF8XmlOutput`) which differs in two visible ways:

* it writes the document element's `xmlns` declaration *before* the other
  attributes (`<DetailedReport xmlns="..." ValidationTime="..."/>`), where
  the writer backend puts it last, and
* its indentation wraps every 8 levels, so an element at depth 8 prints at
  the same zero-space indent as the root. `DetailedReport.xsd` recurses
  (`CRS -> RAC -> CRS`, `SubXCV -> RFC -> ...`) deeply enough to hit this;
  `dr-eaa-status.xml` is the fixture that reaches those depths.

Neither shape is ever produced by DSS itself, so an oracle captured that way
would calibrate the Go marshaller against bytes no DSS consumer will see.
Regenerate through `ReserializeDetailedReport`, not through a raw
`JAXBContext`.

## Regenerating

From a `dss-upstream` checkout with the modules built, compose a classpath of
the `eu.europa.ec.joinup.sd-dss` module classes and their dependencies, then:

    javac -cp "$CP" -d /tmp/oracle ReserializeDetailedReport.java
    java -cp "$CP:/tmp/oracle" ReserializeDetailedReport oracle <dr*.xml>...

The output is a fixed point: feeding a dump back through the same program
reproduces it byte for byte.

`DSS_DETAILEDREPORT_ORACLE_DIR` replays the KAT over a larger local corpus
without committing it.

## Fixture notes

One source fixture (`dr1.xml`) carries a malformed `ValidationTime`
(`2020-02-10T06:53:57`, missing the trailing `Z` `DateParser` requires); the
RI's `unmarshal` silently drops the unparsable value rather than throwing, so
`dr1.xml`'s oracle has no `ValidationTime` attribute on the root at all. This
port's `XSDateTime.UnmarshalText` is stricter (it returns an error), which is
why the KAT feeds these oracle dumps -- already RI-canonical -- back through
Go's `Unmarshal`/`Marshal`, never the malformed originals.

`dr-fill.xml` has no counterpart under `src/test/resources`: it is a
schema-coverage fixture built to reach elements and attributes upstream's own
fixtures never populate. It is dumped through the same facade path as the
rest, so it is RI output like every other file here.
