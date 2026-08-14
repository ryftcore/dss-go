import jakarta.xml.bind.annotation.adapters.XmlAdapter;

import java.io.PrintStream;
import java.lang.reflect.ParameterizedType;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

/**
 * Ground truth for the Go port of the generated JAXB adapters Adapter1..Adapter39
 * of eu.europa.esig.dss.diagnostic.jaxb: the lexical form each adapter prints for
 * every constant of the enumeration it binds. TSV, one row per constant:
 *
 *   adapter &lt;TAB&gt; enumeration &lt;TAB&gt; constant name &lt;TAB&gt; lexical form
 *
 * A constant whose adapter prints null - it has no lexical form - is written
 * with \N in the last column.
 *
 * Usage: EnumLexicalOracle &lt;outFile&gt;
 */
public class EnumLexicalOracle {

    private static final String PKG = "eu.europa.esig.dss.diagnostic.jaxb.";

    @SuppressWarnings({"unchecked", "rawtypes"})
    public static void main(String[] args) throws Exception {
        List<String> rows = new ArrayList<>();
        for (int i = 1; i <= 39; i++) {
            Class<?> adapterClass;
            try {
                adapterClass = Class.forName(PKG + "Adapter" + i);
            } catch (ClassNotFoundException e) {
                continue;
            }
            XmlAdapter adapter = (XmlAdapter) adapterClass.getDeclaredConstructor().newInstance();
            ParameterizedType sup = (ParameterizedType) adapterClass.getGenericSuperclass();
            Class<?> valueType = (Class<?>) sup.getActualTypeArguments()[1];
            if (!valueType.isEnum()) {
                continue;
            }
            for (Object constant : valueType.getEnumConstants()) {
                Object printed = adapter.marshal(constant);
                // A constant the adapter prints as null has no lexical form at
                // all: JAXB then leaves the property out of the document.
                String lexical = printed == null ? "\\N" : printed.toString();
                rows.add(adapterClass.getSimpleName() + "\t" + valueType.getSimpleName() + "\t"
                        + ((Enum<?>) constant).name() + "\t" + lexical);
            }
        }
        StringBuilder sb = new StringBuilder();
        for (String row : rows) {
            sb.append(row).append('\n');
        }
        Files.write(Paths.get(args[0]), sb.toString().getBytes(StandardCharsets.UTF_8));
        PrintStream out = System.out;
        out.println("rows: " + rows.size());
    }
}
