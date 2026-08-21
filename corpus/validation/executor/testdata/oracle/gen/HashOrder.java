import java.util.*;

public class HashOrder {
    public static void main(String[] args) {
        // ids shaped like the DSS token ids the report builders hash
        String[][] cases = {
            {"C-AAA", "C-BBB", "C-CCC"},
            {"R-67877e78139a499c1a04375c3899f9f99eaf368f75895328dc077929fb1705ff",
             "R-cff2e3b95d7ebf7d89b1d586102d09632ce5e39bfa31115c3c9577461dd856ad",
             "R-1a07376222bf46318884ea94953003ecf8e35635c373305c973c1b39169202f0",
             "R-e07f11472bc72bb04a6e3c012dfc838cd323d130bdb3f0e5b1e9558de1b723d0"},
            {"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r"},
            {"ESTEID-SK 2015: Qualified certificates", "ESTEID qualified certificates", "LuxTrust Qualified Time Stamping"},
        };
        for (String[] c : cases) {
            // token-proxy hashCode: 31 + id.hashCode()
            Set<Key> set = new LinkedHashSet<>();
            Set<Key> hs = new HashSet<>();
            for (String s : c) { hs.add(new Key(s)); }
            StringBuilder sb = new StringBuilder();
            for (Key k : hs) { if (sb.length() > 0) sb.append('|'); sb.append(k.id); }
            // plain string HashSet
            Set<String> ss = new HashSet<>(Arrays.asList(c));
            StringBuilder sb2 = new StringBuilder();
            for (String s : ss) { if (sb2.length() > 0) sb2.append('|'); sb2.append(s); }
            System.out.println("TOKEN\t" + String.join("|", c) + "\t" + sb);
            System.out.println("STRING\t" + String.join("|", c) + "\t" + sb2);
        }
    }

    static final class Key {
        final String id;
        Key(String id) { this.id = id; }
        @Override public int hashCode() { return 31 + id.hashCode(); }
        @Override public boolean equals(Object o) { return o instanceof Key && ((Key) o).id.equals(id); }
    }
}
