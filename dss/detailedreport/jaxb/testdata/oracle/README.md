# Marshal-parity oracle corpus

Every `*.xml` file in this directory is byte-for-byte JAXB reference
implementation output (`org.glassfish.jaxb:jaxb-runtime:3.0.2`, the version
`dss-detailed-report-jaxb`'s `xjc` plugin generated the model with), captured
by re-marshalling `dss-detailed-report-jaxb`'s own test fixtures
(`src/test/resources/dr*.xml`) through a fresh `JAXBContext` for
`eu.europa.esig.dss.detailedreport.jaxb.ObjectFactory` with
`JAXB_FORMATTED_OUTPUT=true` - the same configuration
`DetailedReportFacade`/`AbstractJaxbFacade` use.

Regeneration recipe (from a `dss-upstream` checkout with the module already
built, so `target/classes` and `target/generated-sources` exist):

```java
JAXBContext jc = JAXBContext.newInstance(ObjectFactory.class);
Unmarshaller u = jc.createUnmarshaller();
Marshaller m = jc.createMarshaller();
m.setProperty(Marshaller.JAXB_FORMATTED_OUTPUT, Boolean.TRUE);
ObjectFactory of = new ObjectFactory();
XmlDetailedReport report = (XmlDetailedReport) u.unmarshal(new File(inPath));
m.marshal(of.createDetailedReport(report), new FileOutputStream(outPath));
```

compiled and run against:

```
dss-detailed-report-jaxb/target/classes
jakarta.xml.bind:jakarta.xml.bind-api:3.0.1
org.glassfish.jaxb:jaxb-runtime:3.0.2
org.glassfish.jaxb:jaxb-core:3.0.2
org.glassfish.jaxb:txw2:3.0.2
com.sun.istack:istack-commons-runtime:4.0.1
jakarta.activation:jakarta.activation-api:2.1.4
com.sun.activation:jakarta.activation:2.0.1
dss-jaxb-parsers/target/classes   (DateParser and friends, referenced by
                                    XmlDetailedReport's ValidationTime adapter)
dss-enumerations/target/classes
```

One source fixture (`dr1.xml`) carries a malformed `ValidationTime`
(`2020-02-10T06:53:57`, missing the trailing `Z` `DateParser` requires); the
RI's `unmarshal` silently drops the unparsable value rather than throwing, so
`dr1.xml`'s oracle has no `ValidationTime` attribute on the root at all. This
port's `XSDateTime.UnmarshalText` is stricter (it returns an error), which is
why the KAT feeds these oracle dumps - already RI-canonical - back through
Go's `Unmarshal`/`Marshal`, never the malformed originals.
