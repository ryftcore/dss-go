// PasswordFixtures.java - generator for the password-charset known-answer
// fixtures in this directory. Run BY HAND with pdfbox 3.0.x, never from
// go test (internal/pdf/DESIGN.md §5.1); the PDFs it writes are the goldens.
//
//   M2=$HOME/.m2/repository
//   CP=$M2/org/apache/pdfbox/pdfbox/3.0.6/pdfbox-3.0.6.jar\
//   :$M2/org/apache/pdfbox/pdfbox-io/3.0.6/pdfbox-io-3.0.6.jar\
//   :$M2/org/apache/pdfbox/fontbox/3.0.6/fontbox-3.0.6.jar\
//   :$M2/commons-logging/commons-logging/1.3.5/commons-logging-1.3.5.jar
//   javac -nowarn -cp "$CP" -d /tmp/pwfixtures PasswordFixtures.java
//   java -cp "$CP:/tmp/pwfixtures" PasswordFixtures .
//
// Every document is one empty page protected by StandardSecurityHandler
// (StandardProtectionPolicy), so what the fixtures pin is only how pdfbox
// turns the Java String password into the bytes it hashes:
// String.getBytes(ISO_8859_1) for /R 2-4, SaslPrep + String.getBytes(UTF_8)
// for /R 6 (StandardSecurityHandler.prepareDocumentForEncryption /
// prepareForDecryption, PDFBOX-4155). See README.md for the password of each.
import java.io.File;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.encryption.AccessPermission;
import org.apache.pdfbox.pdmodel.encryption.StandardProtectionPolicy;

public final class PasswordFixtures {
    private PasswordFixtures() {}

    public static void main(String[] args) throws Exception {
        File dir = new File(args.length > 0 ? args[0] : ".");
        // /V 2 /R 3, RC4 128-bit: Latin-1 passwords on both sides.
        write(dir, "rc4_r3_latin1.pdf", "caf\u00E9", "own\u00E9r", 128, false);
        // /V 4 /R 4, AES-128 (/AESV2): the common case in the wild.
        write(dir, "aes128_r4_latin1.pdf", "caf\u00E9", "own\u00E9r", 128, true);
        // /V 4 /R 4 with a user password outside Latin-1: ISO-8859-1 encodes
        // U+20AC as '?', so the stored hash is that of "caf?".
        write(dir, "aes128_r4_unmappable.pdf", "caf\u20AC", "owner", 128, true);
        // /V 5 /R 6, AES-256 (/AESV3): UTF-8, no transcoding.
        write(dir, "aes256_r6_utf8.pdf", "caf\u00E9", "own\u00E9r", 256, true);
        // /V 5 /R 6 with a user password SASLprep changes: U+FB01 (fi ligature)
        // NFKC-decomposes to "fi", U+00AD (soft hyphen) maps to nothing, U+00A0
        // (no-break space) maps to ' '. pdfbox stores the hash of "fishca fé".
        write(dir, "aes256_r6_saslprep.pdf", "\uFB01sh\u00ADca\u00A0f\u00E9", "owner", 256, true);
    }

    private static void write(File dir, String name, String user, String owner,
                              int keyLength, boolean preferAES) throws Exception {
        try (PDDocument doc = new PDDocument()) {
            doc.addPage(new PDPage());
            StandardProtectionPolicy policy = new StandardProtectionPolicy(owner, user, new AccessPermission());
            policy.setEncryptionKeyLength(keyLength);
            policy.setPreferAES(preferAES);
            doc.protect(policy);
            doc.save(new File(dir, name));
        }
    }
}
