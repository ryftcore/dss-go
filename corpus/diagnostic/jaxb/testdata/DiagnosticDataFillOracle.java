import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;

import java.io.File;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;
import java.util.TimeZone;

/**
 * Exhaustive ground truth for the Go port of the dss-diagnostic-jaxb generated
 * model: builds an XmlDiagnosticData in which every property of every reachable
 * generated class is set, then marshals it with the JAXB reference
 * implementation. Real validation dumps only ever exercise a fraction of the
 * schema; this dump exercises the model itself, so the Go round-trip KAT covers
 * every element, attribute and lexical form the model can produce.
 *
 * Values are derived from a counter, so the dump is reproducible.
 *
 * Usage: DiagnosticDataFillOracle &lt;outFile&gt; [depth] [full|empty|specials] [counterStart]
 *
 * In "specials" mode every string carries the characters whose XML spelling
 * differs between the reference implementation and encoding/xml.
 *
 * In "empty" mode every property with an empty lexical form (xs:string,
 * xs:base64Binary) is set to that empty form while everything else is filled as
 * usual, so every text-bearing element of the model appears empty at once: that
 * is what pins down the &lt;X&gt;&lt;/X&gt; spelling the reference
 * implementation uses for them, as opposed to the &lt;X/&gt; it uses for an
 * empty complexType. (Numbers, dates and enumerations keep a value: several of
 * the generated adapters dereference them without a null check.)
 */
public class DiagnosticDataFillOracle {

    private static final String PKG = "eu.europa.esig.dss.diagnostic.jaxb.";

    /** Concrete stand-ins for the abstract types an IDREF property is declared with. */
    private static final String ABSTRACT_TOKEN = PKG + "XmlAbstractToken";

    private static int counter = 0;
    private static int maxDepth = 6;
    /** "full" sets every property; "empty" sets the emptiable ones to their empty lexical form. */
    private static String mode = "full";
    private static final Set<String> instantiated = new LinkedHashSet<>();

    private static int next() {
        return ++counter;
    }

    private static Date date() {
        Calendar c = Calendar.getInstance(TimeZone.getTimeZone("UTC"));
        c.clear();
        c.set(2020, Calendar.JANUARY, 1, 0, 0, 0);
        c.add(Calendar.SECOND, next());
        return c.getTime();
    }

    private static boolean isModelClass(Class<?> c) {
        return c.getName().startsWith(PKG) && !c.isEnum();
    }

    /** The concrete class to instantiate for a (possibly abstract) declared type. */
    private static Class<?> concrete(Class<?> c) throws Exception {
        if (!Modifier.isAbstract(c.getModifiers())) {
            return c;
        }
        if (ABSTRACT_TOKEN.equals(c.getName())) {
            return Class.forName(PKG + "XmlCertificate");
        }
        return null;
    }

    /** Builds an instance carrying only its xs:ID, for an IDREF property. */
    private static Object stub(Class<?> declared) throws Exception {
        Class<?> c = concrete(declared);
        if (c == null) {
            return null;
        }
        Object o = c.getDeclaredConstructor().newInstance();
        instantiated.add(c.getSimpleName());
        Field id = idField(c);
        if (id != null) {
            id.setAccessible(true);
            id.set(o, "REF-" + next());
        }
        return o;
    }

    private static Field idField(Class<?> c) {
        for (Class<?> k = c; k != null; k = k.getSuperclass()) {
            for (Field f : k.getDeclaredFields()) {
                if (f.isAnnotationPresent(jakarta.xml.bind.annotation.XmlID.class)) {
                    return f;
                }
            }
        }
        return null;
    }

    /**
     * Every character whose XML spelling differs between the reference
     * implementation and encoding/xml, in character data and in attribute
     * values alike.
     */
    private static final String SPECIALS = "a&b<c>d\"e'f g\th\ni\rj\u20ac\u2211]]>k\uD83D\uDE00";

    private static final String[] CERT_EXTENSIONS = {
        "XmlKeyUsages", "XmlExtendedKeyUsages", "XmlCertificatePolicies",
        "XmlSubjectAlternativeNames", "XmlBasicConstraints", "XmlPolicyConstraints",
        "XmlInhibitAnyPolicy", "XmlNameConstraints", "XmlCRLDistributionPoints",
        "XmlFreshestCRL", "XmlAuthorityKeyIdentifier", "XmlSubjectKeyIdentifier",
        "XmlAuthorityInformationAccess", "XmlIdPkixOcspNoCheck",
        "XmlValAssuredShortTermCertificate", "XmlNoRevAvail", "XmlQcStatements",
        "XmlCertificateExtension",
    };

    private static Object fill(Class<?> declared, int depth) throws Exception {
        Class<?> c = concrete(declared);
        if (c == null) {
            return null;
        }
        Object o = c.getDeclaredConstructor().newInstance();
        instantiated.add(c.getSimpleName());
        for (Class<?> k = c; k != null && k.getName().startsWith(PKG); k = k.getSuperclass()) {
            for (Field f : k.getDeclaredFields()) {
                if (Modifier.isStatic(f.getModifiers())) {
                    continue;
                }
                f.setAccessible(true);
                Object v = value(f, f.getType(), f.getGenericType(), depth);
                if (v != null) {
                    f.set(o, v);
                }
            }
        }
        return o;
    }

    private static Object value(Field f, Class<?> type, Type generic, int depth) throws Exception {
        boolean idref = f != null && f.isAnnotationPresent(jakarta.xml.bind.annotation.XmlIDREF.class);
        boolean isId = f != null && f.isAnnotationPresent(jakarta.xml.bind.annotation.XmlID.class);
        if (type == List.class) {
            Type arg = ((ParameterizedType) generic).getActualTypeArguments()[0];
            Class<?> item = (Class<?>) arg;
            List<Object> list = new ArrayList<>();
            if (f != null && "certificateExtensions".equals(f.getName())) {
                for (String name : CERT_EXTENSIONS) {
                    Object e = fill(Class.forName(PKG + name), depth + 1);
                    if (e != null) {
                        list.add(e);
                    }
                }
                return list;
            }
            for (int i = 0; i < 1; i++) {
                Object e;
                if (idref) {
                    e = stub(item);
                } else if (isModelClass(item)) {
                    if (depth >= maxDepth) {
                        return null;
                    }
                    e = fill(item, depth + 1);
                } else {
                    e = scalar(item, false);
                }
                if (e == null) {
                    return null;
                }
                list.add(e);
            }
            return list;
        }
        if (idref) {
            return stub(type);
        }
        if (isModelClass(type)) {
            if (depth >= maxDepth) {
                return null;
            }
            return fill(type, depth + 1);
        }
        return scalar(type, isId);
    }

    private static Object scalar(Class<?> type, boolean isId) {
        if (type == String.class) {
            if (isId) {
                return "ID-" + next();
            }
            if ("empty".equals(mode)) {
                return "";
            }
            if ("specials".equals(mode)) {
                return SPECIALS + next();
            }
            return "value-" + next();
        }
        if (type == boolean.class) {
            return next() % 2 == 0;
        }
        if (type == int.class) {
            return next();
        }
        if (type == long.class) {
            return (long) next();
        }
        if (type == Boolean.class) {
            return next() % 2 == 0;
        }
        if (type == Integer.class) {
            return next();
        }
        if (type == Long.class) {
            return (long) next();
        }
        if (type == BigInteger.class) {
            return BigInteger.valueOf(next());
        }
        if (type == byte[].class) {
            if ("empty".equals(mode)) {
                return new byte[0];
            }
            byte[] b = new byte[8];
            for (int i = 0; i < b.length; i++) {
                b[i] = (byte) (next() + i);
            }
            return b;
        }
        if (type == Date.class) {
            return date();
        }
        if (type.isEnum()) {
            Object[] values = type.getEnumConstants();
            return values[next() % values.length];
        }
        return null;
    }

    public static void main(String[] args) throws Exception {
        if (args.length > 1) {
            maxDepth = Integer.parseInt(args[1]);
        }
        if (args.length > 2) {
            mode = args[2];
        }
        if (args.length > 3) {
            counter = Integer.parseInt(args[3]);
        }
        XmlDiagnosticData dd = (XmlDiagnosticData) fill(XmlDiagnosticData.class, 0);
        String xml = DiagnosticDataFacade.newFacade().marshall(dd, false);
        Files.write(Paths.get(args[0]), xml.getBytes(StandardCharsets.UTF_8));
        System.out.println("classes instantiated: " + instantiated.size());
        List<String> sorted = new ArrayList<>(instantiated);
        java.util.Collections.sort(sorted);
        System.out.println(sorted);
        System.out.println("bytes: " + new File(args[0]).length());
    }
}
