# jose4j oracle

`jose4j_oracle.tsv` is ground truth produced by running org.jose4j 0.9.6 itself - the exact
version `sd-dss/pom.xml` pins for dss-jades 6.5.RC1 - not by hand-deriving what it ought to do.
`JoseOracle.java` is the generator, kept here so the file can be regenerated and extended:

```sh
mvn dependency:get -Dartifact=org.bitbucket.b_c:jose4j:0.9.6
JAR=~/.m2/repository/org/bitbucket/b_c/jose4j/0.9.6/jose4j-0.9.6.jar
javac -cp "$JAR" -d /tmp/joseoracle JoseOracle.java
java -cp "/tmp/joseoracle:$JAR" JoseOracle > jose4j_oracle.tsv
```

The format is one case per line, tab separated: `SECTION`, a case name, an input column and an
output column. Every payload column is hex of UTF-8 bytes, so a control character in a test case
never has to survive a round trip through a text editor or a terminal.

| section     | input                                | output                                              |
|-------------|--------------------------------------|-----------------------------------------------------|
| `ESC`       | a string value                       | `JsonUtil.toJson` of `{"k": value}`                  |
| `ESCKEY`    | a member name                        | `JsonUtil.toJson` of `{name: "v"}`                   |
| `STRUCT`    | -                                    | `JsonUtil.toJson` of a LinkedHashMap built in code   |
| `VALUE`     | -                                    | `JSONValue.toJSONString` of a non-map value          |
| `HASHORDER` | comma-separated keys, in put order   | `JsonUtil.toJson` of the resulting `java.util.HashMap`|
| `B64ENC`    | raw bytes                            | `Base64Url.encode`                                   |
| `B64DEC`    | encoded text                         | `Base64Url.decode`                                   |
| `CSSER`     | -                                    | `CompactSerializer.serialize` of arguments in code   |
| `CSDESER`   | a compact serialization              | `CompactSerializer.deserialize`, parts joined by `\|` |
| `PARSE`     | JSON text                            | parse then re-serialize, or `ERR:<exception class>`  |
| `PARSEBAD`  | JSON text                            | `OK:<re-serialized>` or `ERR:<exception class>`      |
| `DBL`       | `Double.toHexString` of the input    | `Double.toString` of it                              |
| `ANY`       | JSON text                            | `<class>:<re-serialized>` from a bare `JSONParser`   |

The `HASHORDER` section deserves a note: it is not testing jose4j at all, it is testing
`java.util.HashMap`'s iteration order, which reaches serialized JAdES bytes because DSS's own
`JsonObject()` no-argument constructor wraps a bare HashMap. `many13`, `many20` and `many40`
cross the resize thresholds (16 -> 32 -> 64) on purpose.
