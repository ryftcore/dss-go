// Oracle generator for cades/testdata/attribute-table-order-oracle.txt.
//
// The unsignedAttrsHashIndex of an ats-hash-index (ETSI TS 101 733 archive-time-stamp-v3) is built
// by iterating org.bouncycastle.asn1.cms.AttributeTable#toASN1EncodableVector(), whose order is the
// bucket order of the java.util.Hashtable the table keeps its attributes in - see
// cadesLTAAttributeTableOrder in cades_level_baseline_lta_timestamp_extractor.go. This program asks
// the real BouncyCastle for that order, so the ordering the port reproduces is checked against the
// library rather than against a hand-written replica.
//
// Usage:
//
//	BC=bcprov-jdk18on-<version>.jar:bcutil-jdk18on-<version>.jar
//	javac -cp $BC -d /tmp/ats gen/AtsHashIndexOracle.java
//	java  -cp /tmp/ats:$BC AtsHashIndexOracle attribute-table-order-oracle.txt > /tmp/regenerated.txt
//	diff attribute-table-order-oracle.txt /tmp/regenerated.txt
//
// The input is the oracle file itself: every "IN <attrType>..." line is replayed and its "OUT" line
// is recomputed; comment and blank lines are copied through and the input's own "OUT" lines are
// ignored. The i-th attribute of a case carries the INTEGER i as its value, so that repeated
// attribute types stay distinguishable; each output entry is "<attrType>#<i>".
//
// The committed file was produced with BouncyCastle 1.84 and re-verified byte for byte with 1.78.1:
// the ordering is java.util.Hashtable's and ASN1ObjectIdentifier#hashCode()'s
// (org.bouncycastle.util.Arrays.hashCode over the content octets), neither of which differs
// between those releases.
import java.io.BufferedReader;
import java.io.FileReader;
import java.util.ArrayList;
import java.util.List;

import org.bouncycastle.asn1.ASN1EncodableVector;
import org.bouncycastle.asn1.ASN1Integer;
import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.DERSet;
import org.bouncycastle.asn1.cms.Attribute;
import org.bouncycastle.asn1.cms.AttributeTable;

public class AtsHashIndexOracle {

    public static void main(String[] args) throws Exception {
        try (BufferedReader reader = new BufferedReader(new FileReader(args[0]))) {
            String line;
            while ((line = reader.readLine()) != null) {
                String trimmed = line.trim();
                if (trimmed.startsWith("OUT ")) {
                    continue;
                }
                System.out.println(line);
                if (!trimmed.startsWith("IN ")) {
                    continue;
                }
                String[] types = trimmed.substring(3).trim().split("\\s+");
                ASN1EncodableVector input = new ASN1EncodableVector();
                for (int index = 0; index < types.length; index++) {
                    input.add(new Attribute(new ASN1ObjectIdentifier(types[index]),
                            new DERSet(new ASN1Integer(index))));
                }
                ASN1EncodableVector ordered = new AttributeTable(input).toASN1EncodableVector();
                List<String> out = new ArrayList<>();
                for (int index = 0; index < ordered.size(); index++) {
                    Attribute attribute = Attribute.getInstance(ordered.get(index));
                    ASN1Integer position = ASN1Integer.getInstance(attribute.getAttrValues().getObjectAt(0));
                    out.add(attribute.getAttrType().getId() + "#" + position.getValue());
                }
                System.out.println("OUT " + String.join(" ", out));
            }
        }
    }
}
