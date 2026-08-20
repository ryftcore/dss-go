import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.spi.validation.RevocationDataVerifier;
import eu.europa.esig.dss.spi.validation.TrustAnchorVerifier;
import eu.europa.esig.dss.validation.RevocationDataVerifierFactory;
import eu.europa.esig.dss.validation.TrustAnchorVerifierFactory;
import eu.europa.esig.dss.validation.policy.ValidationPolicyLoader;

import java.io.InputStream;
import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Date;
import java.util.List;
import java.util.Locale;
import java.util.TreeSet;

public class VerifierFactoryOracle {
    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        Date t = new Date(1700000000000L);
        String[] policies = {"", "/policy/certificate-constraint.xml", "/policy/qwac-constraint.xml", "/policy/eaa-constraint.xml"};
        for (String p : policies) {
            ValidationPolicy policy;
            if (p.isEmpty()) {
                policy = ValidationPolicyLoader.fromDefaultValidationPolicy().create();
            } else {
                try (InputStream is = VerifierFactoryOracle.class.getResourceAsStream(p)) {
                    policy = ValidationPolicyLoader.fromValidationPolicy(is).create();
                }
            }
            TrustAnchorVerifier tav = new TrustAnchorVerifierFactory(policy).create();
            RevocationDataVerifier rdv = new RevocationDataVerifierFactory(policy).setValidationTime(t).create();
            System.out.printf("%s\ttav.acceptRevocationUntrusted=%s\ttav.acceptTimestampUntrusted=%s\ttav.useSunsetDate=%s%n",
                    p.isEmpty() ? "default" : p, tav.isAcceptRevocationUntrustedCertificateChains(),
                    tav.isAcceptTimestampUntrustedCertificateChains(), tav.isUseSunsetDate());
            System.out.printf("%s\trdv.digestAlgs=%s\trdv.sigAlgs=%d\trdv.skipExt=%s\trdv.skipPol=%s\trdv.sigFresh=%s\trdv.tstFresh=%s\trdv.revFresh=%s\trdv.nextUpdate=%s\trdv.acceptRevWithout=%s\trdv.acceptTstWithout=%s%n",
                    p.isEmpty() ? "default" : p,
                    sorted(get(rdv, "acceptableDigestAlgorithms")),
                    ((java.util.Map<?, ?>) getRaw(rdv, "acceptableSignatureAlgorithmKeyLength")).size(),
                    sorted(get(rdv, "revocationSkipCertificateExtensions")),
                    sorted(get(rdv, "revocationSkipCertificatePolicies")),
                    getRaw(rdv, "signatureMaximumRevocationFreshness"),
                    getRaw(rdv, "timestampMaximumRevocationFreshness"),
                    getRaw(rdv, "revocationMaximumRevocationFreshness"),
                    getRaw(rdv, "checkRevocationFreshnessNextUpdate"),
                    getRaw(rdv, "acceptRevocationCertificatesWithoutRevocation"),
                    getRaw(rdv, "acceptTimestampCertificatesWithoutRevocation"));
        }
    }

    private static Object getRaw(Object o, String name) throws Exception {
        Field f = o.getClass().getDeclaredField(name);
        f.setAccessible(true);
        return f.get(o);
    }

    @SuppressWarnings("unchecked")
    private static Collection<Object> get(Object o, String name) throws Exception {
        return (Collection<Object>) getRaw(o, name);
    }

    private static String sorted(Collection<Object> c) {
        if (c == null) return "null";
        List<String> l = new ArrayList<>();
        for (Object o : c) l.add(String.valueOf(o));
        return new TreeSet<>(l).toString();
    }
}
