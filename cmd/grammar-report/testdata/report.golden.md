# Grammar compatibility report

This report scans the pinned `textmate-grammars-themes` grammar corpus at revision `37edd1b26f18838050661d912334aba0ca7f4931`. It parses and registry-loads every grammar, then sends every regex field reachable through the parsed grammar model through the same `oniguruma.NewScanner` translation and compilation path used by the tokenizer.

This file is the golden output for `TestGrammarReportGolden`. Update it from the `textmate-go` repository root after an intentional translator or grammar change:

```sh
TM_GRAMMARS_DIR=../tm-grammars/packages/tm-grammars \
  go test -count=1 ./cmd/grammar-report -run '^TestGrammarReportGolden$' -update
```

Numeric capture references in `end` and `while` fields are replaced with a safe literal before compilation. At runtime those references are replaced with escaped text captured by the corresponding `begin`; compiling them as standalone backreferences would report false failures. The report always shows the original grammar pattern.

## Summary

| Check | Count |
|---|---:|
| Grammar JSON files | 260 |
| Parsed grammars | 260 |
| Registry load successes | 260 |
| Grammar load failures | 0 |
| Regex fields scanned | 34698 |
| Patterns with diagnostics | 161 |
| Expected unsupported-syntax diagnostics | 159 |
| Pattern compile failures | 4 |
| Other diagnostics | 0 |

`unsupported_syntax` entries are known Oniguruma constructs that regexp2 cannot faithfully execute; the adapter degrades the affected construct or pattern and records it explicitly. `compile_error` entries are separate: translation completed, but regexp2 rejected the result. A static compile scan cannot produce match-time timeout or match-error diagnostics.

## Grammar load failures

None.

## Expected unsupported syntax

### 1. `apex.json` — `$.repository["array-creation-expression"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(new)\b\s*(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)?\s*(?=\[)
```

- Adapter translation:

```text
\b(new)\b\s*(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)?\s*(?=\[)
```

### 2. `apex.json` — `$.repository["cast-expression"].match`

- Scope: `source.apex`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(\()\s*(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(\))(?=\s*@?[(_[:alnum:]])
```

- Adapter translation:

```text
(\()\s*(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(\))(?=\s*@?[(_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}])
```

### 3. `apex.json` — `$.repository["catch-clause"].patterns[0].patterns[0].match`

- Scope: `source.apex`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:(\g<identifier>)\b)?
```

- Adapter translation:

```text
(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?
```

### 4. `apex.json` — `$.repository["field-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+(\g<identifier>)\s*(?!=[=>])(?=[,;=]|$)
```

- Adapter translation:

```text
(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?!=[=>])(?=[,;=]|$)
```

### 5. `apex.json` — `$.repository["indexer-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(?<indexer_name>this)\s*(?=\[)
```

- Adapter translation:

```text
(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?(?<indexer_name>this)\s*(?=\[)
```

### 6. `apex.json` — `$.repository["invocation-expression"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(\??\.)\s*)?(@?[_[:alpha:]][_[:alnum:]]*)\s*(?<type_args>\s*<([^<>]|\g<type_args>)+>\s*)?\s*(?=\()
```

- Adapter translation:

```text
(?:(\??\.)\s*)?(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<([^<>]|(?!))+>\s*)?\s*(?=\()
```

### 7. `apex.json` — `$.repository["local-constant-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(?<const_keyword>const)\b\s*(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+(\g<identifier>)\s*(?=[,;=])
```

- Adapter translation:

```text
\b(?<const_keyword>const)\b\s*(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=[,;=])
```

### 8. `apex.json` — `$.repository["local-variable-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:\b(ref)\s+)?\b(var)\b|(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+(\g<identifier>)\s*(?=[),;=])
```

- Adapter translation:

```text
(?:(?:\b(ref)\s+)?\b(var)\b|(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=[),;=])
```

### 9. `apex.json` — `$.repository["member-access-expression"].patterns[1].match`

- Scope: `source.apex`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_params&gt;&#34; was neutralized
- Original pattern:

```text
(\??\.)?\s*(@?[_[:alpha:]][_[:alnum:]]*)(?<type_params>\s*<([^<>]|\g<type_params>)+>\s*)(?=(\s*\?)?\s*\.\s*@?[_[:alpha:]][_[:alnum:]]*)
```

- Adapter translation:

```text
(\??\.)?\s*(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)(?<type_params>\s*<([^<>]|(?!))+>\s*)(?=(\s*\?)?\s*\.\s*@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)
```

### 10. `apex.json` — `$.repository["method-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(\g<identifier>)\s*(<([^<>]+)>)?\s*(?=\()
```

- Adapter translation:

```text
(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(<([^<>]+)>)?\s*(?=\()
```

### 11. `apex.json` — `$.repository["object-creation-expression-with-no-parameters"].match`

- Scope: `source.apex`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(delete|insert|undelete|update|upsert)?\s*(new)\s+(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?=\{|$)
```

- Adapter translation:

```text
(delete|insert|undelete|update|upsert)?\s*(new)\s+(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?=\{|$)
```

### 12. `apex.json` — `$.repository["object-creation-expression-with-parameters"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(delete|insert|undelete|update|upsert)?\s*(new)\s+(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?=\()
```

- Adapter translation:

```text
(delete|insert|undelete|update|upsert)?\s*(new)\s+(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?=\()
```

### 13. `apex.json` — `$.repository["parameter"].match`

- Scope: `source.apex`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(this|final)\b\s+)?(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+(\g<identifier>)
```

- Adapter translation:

```text
(?:\b(this|final)\b\s+)?(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))
```

### 14. `apex.json` — `$.repository["property-declaration"].begin`

- Scope: `source.apex`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?!.*\b(?:class|interface|enum)\b)\s*(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(?<property_name>\g<identifier>)\s*(?=\{|=>|$)
```

- Adapter translation:

```text
(?!.*\b(?:class|interface|enum)\b)\s*(?<return_type>(?<type_name>(?:ref\s+)?(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?(?<property_name>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=\{|=>|$)
```

### 15. `blade.json` — `$.repository["function-parameters"].patterns[3].match`

- Scope: `text.html.php.blade`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;8&gt;&#34; was neutralized
- Original pattern:

```text
(?i)(array|callable)\s+((&)?\s*(\$+)[_a-z\x7F-ÿ][0-9_a-z\x7F-ÿ]*)(?:\s*(=)\s*(?:(null)|(\[)((?>[^]\[]+|\[\g<8>])*)(])|(\S*?\(\)|\S*?)))?\s*(?=[),]|/[*/]|#|$)
```

- Adapter translation:

```text
(?i)(array|callable)\s+((&)?\s*(\$+)[_a-z\x7F-ÿ][0-9_a-z\x7F-ÿ]*)(?:\s*(=)\s*(?:(null)|(\[)((?>[^]\[]+|\[(?!)])*)(])|(\S*?\(\)|\S*?)))?\s*(?=[),]|/[*/]|#|$)
```

### 16. `blade.json` — `$.repository["function-parameters"].patterns[6].begin`

- Scope: `text.html.php.blade`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;6&gt;&#34; was neutralized
- Original pattern:

```text
(?i)((&)?\s*(\.\.\.)?(\$+)[_a-z\x7F-ÿ][0-9_a-z\x7F-ÿ]*)\s*(=)\s*(?:(\[)((?>[^]\[]+|\[\g<6>])*)(]))?
```

- Adapter translation:

```text
(?i)((&)?\s*(\.\.\.)?(\$+)[_a-z\x7F-ÿ][0-9_a-z\x7F-ÿ]*)\s*(=)\s*(?:(\[)((?>[^]\[]+|\[(?!)])*)(]))?
```

### 17. `cpp-macro.json` — `$.repository["constructor_inline"].patterns[0].patterns[3].patterns[0].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 18. `cpp-macro.json` — `$.repository["constructor_root"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;8&gt;&#34; was neutralized
- Original pattern:

```text
\s*+((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<8>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)::((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)\10((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\())
```

- Adapter translation:

```text
(?>\s*)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)::((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)\10((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\())
```

### 19. `cpp-macro.json` — `$.repository["constructor_root"].patterns[0].patterns[3].patterns[0].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 20. `cpp-macro.json` — `$.repository["curly_initializer"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\{)
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\{)
```

### 21. `cpp-macro.json` — `$.repository["destructor_root"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)::((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)~\14((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\())
```

- Adapter translation:

```text
((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)::((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)~\14((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\())
```

### 22. `cpp-macro.json` — `$.repository["enum_block"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)enum(?!\w))(?:\s+(class|struct))?(?:(?:\s+|((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\))))|(?=\{))\s+{0,1}((?:(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))?)(?:\s+{0,1}(:)\s+{0,1}(?:((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::))?\s+{0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))?
```

- Adapter translation:

```text
((?<!\w)enum(?!\w))(?:\s+(class|struct))?(?:(?:\s+|((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\))))|(?=\{))(?:\s+){0,1}((?:(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))?)(?:(?:\s+){0,1}(:)(?:\s+){0,1}(?:((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::))?(?:\s+){0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))?
```

### 23. `cpp-macro.json` — `$.repository["function_call"].patterns[0].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;11&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)([A-Z][0-9A-Z_]*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))([A-Z][0-9A-Z_]*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 24. `cpp-macro.json` — `$.repository["function_call"].patterns[1].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;11&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 25. `cpp-macro.json` — `$.repository["function_definition"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;52&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:^|\G|(?<=[;}]))|(?<=>|\*/))\s*+(?:((?<!\w)template(?!\w))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:((?<!\w)(?:(?:constexpr|consteval|explicit|mutable|virtual|inline|friend)|(?:thread_local|volatile|register|restrict|static|extern|const))(?!\w))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\()
```

- Adapter translation:

```text
(?:(?:(?:(?<![\s\S])|^(?=[\s\S]))|\G|(?<=[;}]))|(?<=>|\*/))(?>\s*)(?:((?<!\w)template(?!\w))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:((?<!\w)(?:(?:constexpr|consteval|explicit|mutable|virtual|inline|friend)|(?:thread_local|volatile|register|restrict|static|extern|const))(?!\w))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\()
```

### 26. `cpp-macro.json` — `$.repository["function_definition"].patterns[0].patterns[2].match`

- Scope: `source.cpp.embedded.macro`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;23&gt;&#34; was neutralized
- Original pattern:

```text
(?<=^|\))\s+{0,1}(->)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<23>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<23>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))
```

- Adapter translation:

```text
(?<=(?:(?<![\s\S])|^(?=[\s\S]))|\))(?:\s+){0,1}(->)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))
```

### 27. `cpp-macro.json` — `$.repository["function_pointer"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 28. `cpp-macro.json` — `$.repository["function_pointer_parameter"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 29. `cpp-macro.json` — `$.repository["inheritance_context"].patterns[4].match`

- Scope: `source.cpp.embedded.macro`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(?<=protected|virtual|private|public|[,:])\s+{0,1}(?!p(?:rotected|rivate|ublic)|virtual)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))
```

- Adapter translation:

```text
(?<=protected|virtual|private|public|[,:])(?:\s+){0,1}(?!p(?:rotected|rivate|ublic)|virtual)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))
```

### 30. `cpp-macro.json` — `$.repository["lambdas"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?<=\S|^)(?<![]"\&)*>\[\w])|(?<=(?:\W|^)return))\s+{0,1}(\[(?!\[| *+"| *+\d))((?:[^]\[]|((?<!\[)\[(?!\[)(?:[^]\[]*+\g<3>?)++]))*+)(](?!((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)[];=\[]))
```

- Adapter translation:

```text
(?:(?<=\S|(?:(?<![\s\S])|^(?=[\s\S])))(?<![]"&)*>\[\w])|(?<=(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))return))(?:\s+){0,1}(\[(?!\[|(?> *)"|(?> *)\d))((?>(?:[^]\[]|((?<!\[)\[(?!\[)(?>(?:(?>[^]\[]*)(?!)?)+)]))*))(](?!((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)[];=\[]))
```

### 31. `cpp-macro.json` — `$.repository["namespace_block"].patterns[0].patterns[4].match`

- Scope: `source.cpp.embedded.macro`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;4&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<4>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)\s+{0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s+{0,1}(?:(::)\s+{0,1}(inline))?
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))(?:\s+){0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?:\s+){0,1}(?:(::)(?:\s+){0,1}(inline))?
```

### 32. `cpp-macro.json` — `$.repository["normal_variable_assignment"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;31&gt;&#34; was neutralized
- Original pattern:

```text
^((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[\&^]|<<|>>|\|)=)|(=)))
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[&^]|<<|>>|\|)=)|(=)))
```

### 33. `cpp-macro.json` — `$.repository["normal_variable_declaration"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;31&gt;&#34; was neutralized
- Original pattern:

```text
^((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[,;\[])(?![^=]++=))
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[,;\[])(?!(?>[^=]+)=))
```

### 34. `cpp-macro.json` — `$.repository["operator_overload"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;60&gt;&#34; was neutralized
- Original pattern:

```text
((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(operator)(?:((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?<!\w)const(?!\w)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(?:(?:(delete\[]|delete|new\[]|<=>|<<=|new|>>=|->\*|/=|%=|&=|>=|\|=|\+\+|--|\(\)|\[]|->|\+\+|<<|>>|--|<=|\^=|==|!=|&&|\|\||\+=|-=|\*=|[!%\&*-\-/<=>^|~])|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:\[])?)))|("")((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[(;<])
```

- Adapter translation:

```text
((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(operator)(?:((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?<!\w)const(?!\w)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(?:(?:(delete\[]|delete|new\[]|<=>|<<=|new|>>=|->\*|/=|%=|&=|>=|\|=|\+\+|--|\(\)|\[]|->|\+\+|<<|>>|--|<=|\^=|==|!=|&&|\|\||\+=|-=|\*=|[!%&*-\-/<=>^|~])|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:\[])?)))|("")((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[(;<])
```

### 35. `cpp-macro.json` — `$.repository["range_for_inner"].patterns[0].match`

- Scope: `source.cpp.embedded.macro`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;27&gt;&#34; was neutralized
- Original pattern:

```text
((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(:)(?!:)
```

- Adapter translation:

```text
((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(:)(?!:)
```

### 36. `cpp-macro.json` — `$.repository["range_for_inner"].patterns[1].match`

- Scope: `source.cpp.embedded.macro`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;27&gt;&#34; was neutralized
- Original pattern:

```text
((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\[)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)(?:((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(,)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*))*((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(])((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(:)(?!:)
```

- Adapter translation:

```text
((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\[)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)(?:((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(,)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*))*((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(])((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(:)(?!:)
```

### 37. `cpp-macro.json` — `$.repository["typedef_function_pointer"].patterns[0].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 38. `cpp-macro.json` — `$.repository["using_namespace"].begin`

- Scope: `source.cpp.embedded.macro`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;6&gt;&#34; was neutralized
- Original pattern:

```text
(?<!\w)(using)\s+(namespace)\s+((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<6>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)?((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(?=[\n;])
```

- Adapter translation:

```text
(?<!\w)(using)\s+(namespace)\s+((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))?((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?=[\n;])
```

### 39. `cpp.json` — `$.repository["constructor_bracket_call"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(?!class|struct|union|enum|explicit|new|delete|operator|template|throw|decltype|typename|override|final)\b(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\{)
```

- Adapter translation:

```text
(?!class|struct|union|enum|explicit|new|delete|operator|template|throw|decltype|typename|override|final)\b((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\{)
```

### 40. `cpp.json` — `$.repository["constructor_inline"].patterns[0].patterns[3].patterns[0].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 41. `cpp.json` — `$.repository["constructor_root"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;8&gt;&#34; was neutralized
- Original pattern:

```text
\s*+((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<8>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)::((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)\10((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\())
```

- Adapter translation:

```text
(?>\s*)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)::((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)\10((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\())
```

### 42. `cpp.json` — `$.repository["constructor_root"].patterns[0].patterns[3].patterns[0].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 43. `cpp.json` — `$.repository["curly_initializer"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\{)
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\{)
```

### 44. `cpp.json` — `$.repository["destructor_root"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)::((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)~\14((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\())
```

- Adapter translation:

```text
((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?>(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)::((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)~\14((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\())
```

### 45. `cpp.json` — `$.repository["enum_block"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
((?<!\w)enum(?!\w))(?:\s+(class|struct))?(?:(?:\s+|((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\))))|(?=\{))\s+{0,1}((?:(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))?)(?:\s+{0,1}(:)\s+{0,1}(?:((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::))?\s+{0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))?
```

- Adapter translation:

```text
((?<!\w)enum(?!\w))(?:\s+(class|struct))?(?:(?:\s+|((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\))))|(?=\{))(?:\s+){0,1}((?:(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))?)(?:(?:\s+){0,1}(:)(?:\s+){0,1}(?:((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::))?(?:\s+){0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))?
```

### 46. `cpp.json` — `$.repository["function_call"].patterns[0].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;11&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)([A-Z][0-9A-Z_]*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))([A-Z][0-9A-Z_]*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 47. `cpp.json` — `$.repository["function_call"].patterns[1].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;11&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(\()
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(\()
```

### 48. `cpp.json` — `$.repository["function_definition"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;52&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:^|\G|(?<=[;}]))|(?<=>|\*/))\s*+(?:((?<!\w)template(?!\w))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:((?<!\w)(?:(?:constexpr|consteval|explicit|mutable|virtual|inline|friend)|(?:thread_local|volatile|register|restrict|static|extern|const))(?!\w))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<52>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)\b(?<!(?:\W|^)(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=\()
```

- Adapter translation:

```text
(?:(?:(?:(?<![\s\S])|^(?=[\s\S]))|\G|(?<=[;}]))|(?<=>|\*/))(?>\s*)(?:((?<!\w)template(?!\w))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:((?<!\w)(?:(?:constexpr|consteval|explicit|mutable|virtual|inline|friend)|(?:thread_local|volatile|register|restrict|static|extern|const))(?!\w))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)\b(?<!(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))(?:reinterpret_cast|atomic_noexcept|uint_least16_t|uint_least32_t|uint_least64_t|atomic_cancel|atomic_commit|uint_least8_t|uint_fast16_t|uint_fast32_t|int_least16_t|int_least32_t|int_least64_t|uint_fast64_t|thread_local|int_fast16_t|int_fast32_t|int_fast64_t|synchronized|uint_fast8_t|dynamic_cast|int_least8_t|int_fast8_t|static_cast|suseconds_t|const_cast|useconds_t|constinit|co_return|uintmax_t|constexpr|consteval|constexpr|consteval|protected|namespace|blksize_t|co_return|in_addr_t|in_port_t|uintptr_t|template|noexcept|continue|co_await|co_yield|unsigned|u_quad_t|blkcnt_t|uint16_t|uint32_t|uint64_t|intptr_t|intmax_t|volatile|register|restrict|explicit|volatile|noexcept|operator|decltype|typename|requires|co_await|co_yield|reflexpr|swblk_t|virtual|ssize_t|concept|mutable|fixpt_t|int16_t|int32_t|int64_t|uint8_t|typedef|daddr_t|caddr_t|qaddr_t|default|nlink_t|segsz_t|u_short|wchar_t|private|__asm__|alignas|alignof|mutable|nullptr|clock_t|mode_t|public|size_t|double|quad_t|static|time_t|module|import|export|extern|inline|xor_eq|and_eq|return|friend|not_eq|signed|struct|int8_t|ushort|switch|u_long|typeid|u_char|sizeof|bitand|delete|ino_t|key_t|pid_t|off_t|uid_t|short|break|catch|compl|while|false|class|union|const|or_eq|const|throw|bitor|u_int|using|div_t|dev_t|gid_t|float|long|goto|uint|id_t|case|auto|void|enum|true|char|id_t|NULL|this|bool|else|for|new|not|xor|and|asm|int|try|do|if|or))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=\()
```

### 49. `cpp.json` — `$.repository["function_definition"].patterns[0].patterns[2].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;23&gt;&#34; was neutralized
- Original pattern:

```text
(?<=^|\))\s+{0,1}(->)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<23>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<23>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))
```

- Adapter translation:

```text
(?<=(?:(?<![\s\S])|^(?=[\s\S]))|\))(?:\s+){0,1}(->)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))
```

### 50. `cpp.json` — `$.repository["function_pointer"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 51. `cpp.json` — `$.repository["function_pointer_parameter"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 52. `cpp.json` — `$.repository["inheritance_context"].patterns[4].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(?<=protected|virtual|private|public|[,:])\s+{0,1}(?!p(?:rotected|rivate|ublic)|virtual)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))
```

- Adapter translation:

```text
(?<=protected|virtual|private|public|[,:])(?:\s+){0,1}(?!p(?:rotected|rivate|ublic)|virtual)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))
```

### 53. `cpp.json` — `$.repository["lambdas"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?<=\S|^)(?<![]"\&)*>\[\w])|(?<=(?:\W|^)return))\s+{0,1}(\[(?!\[| *+"| *+\d))((?:[^]\[]|((?<!\[)\[(?!\[)(?:[^]\[]*+\g<3>?)++]))*+)(](?!((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)[];=\[]))
```

- Adapter translation:

```text
(?:(?<=\S|(?:(?<![\s\S])|^(?=[\s\S])))(?<![]"&)*>\[\w])|(?<=(?:\W|(?:(?<![\s\S])|^(?=[\s\S])))return))(?:\s+){0,1}(\[(?!\[|(?> *)"|(?> *)\d))((?>(?:[^]\[]|((?<!\[)\[(?!\[)(?>(?:(?>[^]\[]*)(?!)?)+)]))*))(](?!((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)[];=\[]))
```

### 54. `cpp.json` — `$.repository["namespace_alias"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;8&gt;&#34; was neutralized
- Original pattern:

```text
(?<!\w)(namespace)\s+((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s+{0,1}(=)\s+{0,1}(((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<8>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)\s+{0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s+{0,1}(?:(;)|\n))
```

- Adapter translation:

```text
(?<!\w)(namespace)\s+((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?:\s+){0,1}(=)(?:\s+){0,1}(((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))(?:\s+){0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?:\s+){0,1}(?:(;)|\n))
```

### 55. `cpp.json` — `$.repository["namespace_block"].patterns[0].patterns[4].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;4&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<4>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)\s+{0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s+{0,1}(?:(::)\s+{0,1}(inline))?
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))(?:\s+){0,1}((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?:\s+){0,1}(?:(::)(?:\s+){0,1}(inline))?
```

### 56. `cpp.json` — `$.repository["normal_variable_assignment"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;31&gt;&#34; was neutralized
- Original pattern:

```text
^((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[\&^]|<<|>>|\|)=)|(=)))
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[&^]|<<|>>|\|)=)|(=)))
```

### 57. `cpp.json` — `$.repository["normal_variable_declaration"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;31&gt;&#34; was neutralized
- Original pattern:

```text
^((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<31>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[,;\[])(?![^=]++=))
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[,;\[])(?!(?>[^=]+)=))
```

### 58. `cpp.json` — `$.repository["operator_overload"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;60&gt;&#34; was neutralized
- Original pattern:

```text
((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(operator)(?:((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?<!\w)const(?!\w)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<60>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(?:(?:(delete\[]|delete|new\[]|<=>|<<=|new|>>=|->\*|/=|%=|&=|>=|\|=|\+\+|--|\(\)|\[]|->|\+\+|<<|>>|--|<=|\^=|==|!=|&&|\|\||\+=|-=|\*=|[!%\&*-\-/<=>^|~])|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:\[])?)))|("")((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[(;<])
```

- Adapter translation:

```text
((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?:__(?:cdec|clrcal|stdcal|fastcal|thiscal|vectorcal)l)?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(operator)(?:((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?<!\w)const(?!\w)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(?:(?:(delete\[]|delete|new\[]|<=>|<<=|new|>>=|->\*|/=|%=|&=|>=|\|=|\+\+|--|\(\)|\[]|->|\+\+|<<|>>|--|<=|\^=|==|!=|&&|\|\||\+=|-=|\*=|[!%&*-\-/<=>^|~])|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:\[])?)))|("")((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[(;<])
```

### 59. `cpp.json` — `$.repository["qualified_type"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;11&gt;&#34; was neutralized
- Original pattern:

```text
\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<11>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w])
```

- Adapter translation:

```text
(?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w])
```

### 60. `cpp.json` — `$.repository["range_for_inner"].patterns[0].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;27&gt;&#34; was neutralized
- Original pattern:

```text
((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(:)(?!:)
```

- Adapter translation:

```text
((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(:)(?!:)
```

### 61. `cpp.json` — `$.repository["range_for_inner"].patterns[1].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;27&gt;&#34; was neutralized
- Original pattern:

```text
((?:((?:(?:(?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<27>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\[)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)(?:((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(,)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*))*((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(])((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(:)(?!:)
```

- Adapter translation:

```text
((?:((?:(?:(?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\[)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)(?:((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(,)((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*))*((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(])((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(:)(?!:)
```

### 62. `cpp.json` — `$.repository["scope_resolution"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 63. `cpp.json` — `$.repository["scope_resolution_function_call"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 64. `cpp.json` — `$.repository["scope_resolution_function_call_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 65. `cpp.json` — `$.repository["scope_resolution_function_definition"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 66. `cpp.json` — `$.repository["scope_resolution_function_definition_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 67. `cpp.json` — `$.repository["scope_resolution_function_definition_operator_overload"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 68. `cpp.json` — `$.repository["scope_resolution_function_definition_operator_overload_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 69. `cpp.json` — `$.repository["scope_resolution_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 70. `cpp.json` — `$.repository["scope_resolution_namespace_alias"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 71. `cpp.json` — `$.repository["scope_resolution_namespace_alias_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 72. `cpp.json` — `$.repository["scope_resolution_namespace_block"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 73. `cpp.json` — `$.repository["scope_resolution_namespace_block_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 74. `cpp.json` — `$.repository["scope_resolution_namespace_using"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 75. `cpp.json` — `$.repository["scope_resolution_namespace_using_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 76. `cpp.json` — `$.repository["scope_resolution_parameter"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 77. `cpp.json` — `$.repository["scope_resolution_parameter_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 78. `cpp.json` — `$.repository["scope_resolution_template_call"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 79. `cpp.json` — `$.repository["scope_resolution_template_call_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 80. `cpp.json` — `$.repository["scope_resolution_template_definition"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<3>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+
```

- Adapter translation:

```text
(::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*)
```

### 81. `cpp.json` — `$.repository["scope_resolution_template_definition_inner_generated"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;7&gt;&#34; was neutralized
- Original pattern:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))\s*+(((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<7>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?(::)
```

- Adapter translation:

```text
((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))((?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?>\s*)(((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?(::)
```

### 82. `cpp.json` — `$.repository["simple_array_assignment"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*](((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))(\[) *(])(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[\&^]|<<|>>|\|)=)|(=))
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*](((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))(\[) *(])(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[&^]|<<|>>|\|)=)|(=))
```

### 83. `cpp.json` — `$.repository["simple_constructor_call"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(?!class|struct|union|enum|explicit|new|delete|operator|template|throw|decltype|typename|override|final)\b(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(?=(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)(?=(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[({]))))
```

- Adapter translation:

```text
(?!class|struct|union|enum|explicit|new|delete|operator|template|throw|decltype|typename|override|final)\b((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(?=(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?=(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[({]))))
```

### 84. `cpp.json` — `$.repository["simple_type"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;12&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<12>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*](((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*](((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?
```

### 85. `cpp.json` — `$.repository["template_call_innards"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;1&gt;&#34; was neutralized
- Original pattern:

```text
((?<!<)<(?!<)(?:(/\*)((?:[^*]++|\*+(?!/))*+(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<1>|(?:[^"'/<>]|/[^*])++)*>)\s*+
```

- Adapter translation:

```text
((?<!<)<(?!<)(?:(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*)
```

### 86. `cpp.json` — `$.repository["type_alias"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;19&gt;&#34; was neutralized
- Original pattern:

```text
(using)\s+(?!namespace)((?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)\s+{0,1}((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?\s+{0,1}(=)\s+{0,1}((?:typename)?)\s+{0,1}((?:(?:((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)?(?:(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<19>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<19>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))|(.*(?<!;)))(?:((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*](((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?:(\[)(\w*)(])\s+{0,1})?\s+{0,1}(?:(;)|\n)
```

- Adapter translation:

```text
(using)\s+(?!namespace)((?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)(?:\s+){0,1}((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(?:\s+){0,1}(=)(?:\s+){0,1}((?:typename)?)(?:\s+){0,1}((?:(?:((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)?(?:((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))|(.*(?<!;)))(?:((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*](((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?:(\[)(\w*)(])(?:\s+){0,1})?(?:\s+){0,1}(?:(;)|\n)
```

### 87. `cpp.json` — `$.repository["typedef_function_pointer"].patterns[0].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;18&gt;&#34; was neutralized
- Original pattern:

```text
(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<18>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))(((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*]((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?((?:\s*+(/\*)((?:[^*]++|\*+(?!/))*+(\*/))\s*+)+|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\()(\*)\s+{0,1}((?:(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*)?)\s+{0,1}(?:(\[)(\w*)(])\s+{0,1})*(\))\s+{0,1}(\()
```

- Adapter translation:

```text
((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))(((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*]((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?((?:(?>\s*)(/\*)((?>(?:(?>[^*]+)|\*+(?!/))*)(\*/))(?>\s*))+|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(\()(\*)(?:\s+){0,1}((?:(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*)?)(?:\s+){0,1}(?:(\[)(\w*)(])(?:\s+){0,1})*(\))(?:\s+){0,1}(\()
```

### 88. `cpp.json` — `$.repository["typename"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;17&gt;&#34; was neutralized
- Original pattern:

```text
((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)typename(?!\w))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<17>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<17>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))
```

- Adapter translation:

```text
((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)typename(?!\w))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))
```

### 89. `cpp.json` — `$.repository["using_namespace"].begin`

- Scope: `source.cpp`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;6&gt;&#34; was neutralized
- Original pattern:

```text
(?<!\w)(using)\s+(namespace)\s+((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<6>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*\s*+)?((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w))(?=[\n;])
```

- Adapter translation:

```text
(?<!\w)(using)\s+(namespace)\s+((::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*(?>\s*))?((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w))(?=[\n;])
```

### 90. `cpp.json` — `$.repository["variable_assignment"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;16&gt;&#34; was neutralized
- Original pattern:

```text
(?:((?:(?:((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<16>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<16>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*](((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[\&^]|<<|>>|\|)=)|(=))
```

- Adapter translation:

```text
(?:((?:(?:((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*](((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:((?:[-%*+]|(?<!\()/)=)|((?:[&^]|<<|>>|\|)=)|(=))
```

### 91. `cpp.json` — `$.repository["variable_declare"].match`

- Scope: `source.cpp`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;16&gt;&#34; was neutralized
- Original pattern:

```text
(?:((?:(?:((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(\s*+((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*(?:((?:::)?(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)\s*+(((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<16>|(?:[^"'/<>]|/[^*])++)*>)\s*+)?::)*+)(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*\b((?<!<)<(?!<)(?:/\*(?:[^*]++|\*+(?!/))*+\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|\g<16>|(?:[^"'/<>]|/[^*])++)*>)?(?![.:<\w]))((((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)?(?:[\&*](((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z))*[\&*])?(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u\h{4}|U\h{8}))(?:[0-9A-Z_a-z]|\\(?:u\h{4}|U\h{8}))*(?!\w)))(((?:\s*+/\*(?:[^*]++|\*+(?!/))*+\*/\s*+)+)|\s++|(?<=\W)|(?=\W)|^|\n?$|\A|\Z)(?=[,;\[])(?![^=]++=)
```

- Adapter translation:

```text
(?:((?:(?:((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?<!\w)(?:thread_local|volatile|register|restrict|static|extern|const)(?!\w)\s+)+)(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?((?>\s*)((?:(?:(?:\[\[.*?]]|__attribute(?:__)?\s*\(\s*\(.*?\)\s*\))|__declspec\(.*?\))|alignas\(.*?\))(?!\)))?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:(?:(?:unsigned|signed|short|long)|(?:struct|class|union|enum))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*(?:((?:::)?(?>(?:(?!\b(?:__has_cpp_attribute|reinterpret_cast|atomic_noexcept|atomic_commit|atomic_cancel|__has_include|thread_local|dynamic_cast|synchronized|static_cast|const_cast|consteval|co_return|protected|constinit|constexpr|co_return|consteval|namespace|constexpr|co_await|explicit|volatile|noexcept|co_yield|noexcept|requires|typename|decltype|operator|template|continue|co_await|co_yield|volatile|register|restrict|reflexpr|mutable|alignof|include|private|defined|typedef|_Pragma|__asm__|concept|mutable|warning|default|virtual|alignas|public|sizeof|delete|not_eq|bitand|and_eq|xor_eq|typeid|switch|return|struct|static|extern|inline|friend|ifndef|define|pragma|export|import|module|catch|throw|const|or_eq|compl|while|ifdef|const|bitor|union|class|undef|error|break|using|endif|goto|line|enum|this|case|else|elif|else|not|try|for|asm|and|xor|new|do|if|or|if)\b)(?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)(?>\s*)(((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)(?>\s*))?::)*))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))?(?!(?:transaction_safe_dynamic|__has_cpp_attribute|reinterpret_cast|transaction_safe|atomic_noexcept|atomic_commit|__has_include|atomic_cancel|synchronized|thread_local|dynamic_cast|static_cast|const_cast|constexpr|co_return|constinit|namespace|protected|consteval|constexpr|co_return|consteval|co_await|continue|template|reflexpr|volatile|register|co_await|co_yield|restrict|noexcept|volatile|override|explicit|decltype|operator|noexcept|typename|requires|co_yield|nullptr|alignof|alignas|default|mutable|virtual|mutable|private|include|warning|_Pragma|defined|typedef|__asm__|concept|define|module|sizeof|switch|delete|pragma|and_eq|inline|xor_eq|typeid|import|extern|public|bitand|static|export|return|friend|ifndef|not_eq|false|final|break|const|catch|endif|ifdef|undef|error|audit|while|using|axiom|or_eq|compl|throw|bitor|const|line|case|else|this|true|goto|else|NULL|elif|new|asm|xor|and|try|not|for|do|if|or|if)\b)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*\b((?<!<)<(?!<)(?:/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/|"(?:[^"]*|\\")"|'(?:[^']*|\\')'|(?!)|(?>(?:[^"'/<>]|/[^*])+))*>)?(?![.:<\w]))((((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)?(?:[&*](((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z))*[&*])?(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?:\b(?:(?:(?:([0-9a-z]+)|([0-9A-Za-z]+_[0-9A-Za-z]*))|([a-z]+[A-Z][0-9A-Za-z]*))|([A-Z][0-9A-Z_]*))\b|((?<!\w)(?:[A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))(?:[0-9A-Z_a-z]|\\(?:u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8}))*(?!\w)))(((?:(?>\s*)/\*(?>(?:(?>[^*]+)|\*+(?!/))*)\*/(?>\s*))+)|(?>\s+)|(?<=\W)|(?=\W)|(?:(?<![\s\S])|^(?=[\s\S]))|\n?$|\A|\Z)(?=[,;\[])(?!(?>[^=]+)=)
```

### 92. `csharp.json` — `$.repository["anonymous-method-expression"].patterns[0].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
((?:\b(?:async|static)\b\s*)*)(?:(@?[_[:alpha:]][_[:alnum:]]*)\b|(\()(?<tuple>(?:[^()]|\(\g<tuple>\))*)(\)))\s*(=>)
```

- Adapter translation:

```text
((?:\b(?:async|static)\b\s*)*)(?:(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\b|(\()(?<tuple>(?:[^()]|\((?!)\))*)(\)))\s*(=>)
```

### 93. `csharp.json` — `$.repository["array-creation-expression"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(new|stackalloc)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\*\s*)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s*(?=\[)
```

- Adapter translation:

```text
\b(new|stackalloc)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\*\s*)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s*(?=\[)
```

### 94. `csharp.json` — `$.repository["as-expression"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<!\.)\b(as)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?(?!\?))?(?:\s*\[\s*(?:,\s*)*](?:\s*\?(?!\?))?)*)?
```

- Adapter translation:

```text
(?<!\.)\b(as)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?(?!\?))?(?:\s*\[\s*(?:,\s*)*](?:\s*\?(?!\?))?)*)?
```

### 95. `csharp.json` — `$.repository["base-class-constructor-call"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(@?[_[:alpha:]][_[:alnum:]]*)\s*(\.))*(@?[_[:alpha:]][_[:alnum:]]*)\s*(<(?<type_args>[^()<>]|\((?:[^()<>]|<[^()<>]*>|\([^()<>]*\))*\)|<\g<type_args>*>)*>\s*)?(?=\()
```

- Adapter translation:

```text
(?:(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(\.))*(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(<(?<type_args>[^()<>]|\((?:[^()<>]|<[^()<>]*>|\([^()<>]*\))*\)|<(?!)*>)*>\s*)?(?=\()
```

### 96. `csharp.json` — `$.repository["cast-expression"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(\()\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(\))(?=\s*-*!*@?[(_[:alnum:]])
```

- Adapter translation:

```text
(\()\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(\))(?=\s*-*!*@?[(_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}])
```

### 97. `csharp.json` — `$.repository["catch-clause"].patterns[0].patterns[0].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?:(\g<identifier>)\b)?
```

- Adapter translation:

```text
(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?:((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?
```

### 98. `csharp.json` — `$.repository["conversion-operator-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(?<explicit_or_implicit_keyword>(?:ex|im)plicit)\s*\b(?<operator_keyword>operator)\s*(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?=\()
```

- Adapter translation:

```text
\b(?<explicit_or_implicit_keyword>(?:ex|im)plicit)\s*\b(?<operator_keyword>operator)\s*(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?=\()
```

### 99. `csharp.json` — `$.repository["declaration-expression-local"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(var)\b|(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+(\g<identifier>)\b\s*(?=[]),])
```

- Adapter translation:

```text
(?:\b(var)\b|(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b\s*(?=[]),])
```

### 100. `csharp.json` — `$.repository["declaration-expression-tuple"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(var)\b|(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+(\g<identifier>)\b\s*(?=[),])
```

- Adapter translation:

```text
(?:\b(var)\b|(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b\s*(?=[),])
```

### 101. `csharp.json` — `$.repository["delegate-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(delegate)\b\s+(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+(\g<identifier>)\s*(<([^<>]+)>)?\s*(?=\()
```

- Adapter translation:

```text
\b(delegate)\b\s+(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(<([^<>]+)>)?\s*(?=\()
```

### 102. `csharp.json` — `$.repository["event-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(event)\b\s*(?<return_type>(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(\g<identifier>)\s*(?=[,;={]|//|/\*|$)
```

- Adapter translation:

```text
\b(event)\b\s*(?<return_type>(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=[,;={]|//|/\*|$)
```

### 103. `csharp.json` — `$.repository["explicit-anonymous-function-parameter"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(ref|params|out|in)\b\s*)?(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args><(?:[^<>]|\g<type_args>)*>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)*\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*\b(\g<identifier>)\b
```

- Adapter translation:

```text
(?:\b(ref|params|out|in)\b\s*)?(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args><(?:[^<>]|(?!))*>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))*\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*\b((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b
```

### 104. `csharp.json` — `$.repository["field-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+(\g<identifier>)\s*(?!=[=>])(?=[,;=]|$)
```

- Adapter translation:

```text
(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?!=[=>])(?=[,;=]|$)
```

### 105. `csharp.json` — `$.repository["fixed-size-buffer-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(fixed)\b\s+(?<type_name>(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*)\s+(\g<identifier>)\s*(?=\[)
```

- Adapter translation:

```text
\b(fixed)\b\s+(?<type_name>(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=\[)
```

### 106. `csharp.json` — `$.repository["foreach-statement"].patterns[1].patterns[1].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:\b(ref)\s+)?\b(var)\b|(?<type_name>(?:ref\s+)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+(\g<identifier>)\s+\b(in)\b
```

- Adapter translation:

```text
(?:(?:\b(ref)\s+)?\b(var)\b|(?<type_name>(?:ref\s+)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s+\b(in)\b
```

### 107. `csharp.json` — `$.repository["foreach-statement"].patterns[1].patterns[2].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(var)\b\s*)?(?<tuple>\((?:[^()]|\g<tuple>)+\))\s+\b(in)\b
```

- Adapter translation:

```text
(?:\b(var)\b\s*)?(?<tuple>\((?:[^()]|(?!))+\))\s+\b(in)\b
```

### 108. `csharp.json` — `$.repository["indexer-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(?<indexer_name>this)\s*(?=\[)
```

- Adapter translation:

```text
(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?(?<indexer_name>this)\s*(?=\[)
```

### 109. `csharp.json` — `$.repository["invocation-expression"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:(\?)\s*)?(\.)\s*|(->)\s*)?(@?[_[:alpha:]][_[:alnum:]]*)\s*(<(?<type_args>[^()<>]|\((?:[^()<>]|<[^()<>]*>|\([^()<>]*\))*\)|<\g<type_args>*>)*>\s*)?(?=\()
```

- Adapter translation:

```text
(?:(?:(\?)\s*)?(\.)\s*|(->)\s*)?(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(<(?<type_args>[^()<>]|\((?:[^()<>]|<[^()<>]*>|\([^()<>]*\))*\)|<(?!)*>)*>\s*)?(?=\()
```

### 110. `csharp.json` — `$.repository["join-clause"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(join)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s+(\g<identifier>)\b\s*\b(in)\b\s*
```

- Adapter translation:

```text
\b(join)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b\s*\b(in)\b\s*
```

### 111. `csharp.json` — `$.repository["local-constant-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(?<const_keyword>const)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+(\g<identifier>)\s*(?=[,;=])
```

- Adapter translation:

```text
\b(?<const_keyword>const)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=[,;=])
```

### 112. `csharp.json` — `$.repository["local-function-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b((?:(?:async|unsafe|static|extern)\s+)*)(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?)?(?:\s*\[\s*(?:,\s*)*](?:\s*\?)?)*)\s+(\g<identifier>)\s*(<[^<>]+>)?\s*(?=\()
```

- Adapter translation:

```text
\b((?:(?:async|unsafe|static|extern)\s+)*)(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?)?(?:\s*\[\s*(?:,\s*)*](?:\s*\?)?)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(<[^<>]+>)?\s*(?=\()
```

### 113. `csharp.json` — `$.repository["local-tuple-declaration-deconstruction"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
(?<tuple>\((?:[^()]|\g<tuple>)+\))\s*(?!=[=>])(?==)
```

- Adapter translation:

```text
(?<tuple>\((?:[^()]|(?!))+\))\s*(?!=[=>])(?==)
```

### 114. `csharp.json` — `$.repository["local-tuple-var-deconstruction"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
\b(var)\b\s*(?<tuple>\((?:[^()]|\g<tuple>)+\))\s*(?=[);=])
```

- Adapter translation:

```text
\b(var)\b\s*(?<tuple>\((?:[^()]|(?!))+\))\s*(?=[);=])
```

### 115. `csharp.json` — `$.repository["local-variable-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:(?:\b(ref)\s+(?:\b(readonly)\s+)?)?\b(var)\b|(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\*\s*)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+(\g<identifier>)\s*(?!=>)(?=[),;=])
```

- Adapter translation:

```text
(?:(?:\b(ref)\s+(?:\b(readonly)\s+)?)?\b(var)\b|(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\*\s*)*(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?!=>)(?=[),;=])
```

### 116. `csharp.json` — `$.repository["member-access-expression"].patterns[1].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_params&gt;&#34; was neutralized
- Original pattern:

```text
(\.)?\s*(@?[_[:alpha:]][_[:alnum:]]*)(?<type_params>\s*<([^<>]|\g<type_params>)+>\s*)(?=(\s*\?)?\s*\.\s*@?[_[:alpha:]][_[:alnum:]]*)
```

- Adapter translation:

```text
(\.)?\s*(@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)(?<type_params>\s*<([^<>]|(?!))+>\s*)(?=(\s*\?)?\s*\.\s*@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)
```

### 117. `csharp.json` — `$.repository["method-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(\g<identifier>)\s*(<([^<>]+)>)?\s*(?=\()
```

- Adapter translation:

```text
(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(<([^<>]+)>)?\s*(?=\()
```

### 118. `csharp.json` — `$.repository["object-creation-expression-with-no-parameters"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(new)\s+(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?=\{|//|/\*|$)
```

- Adapter translation:

```text
(new)\s+(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*(?=\{|//|/\*|$)
```

### 119. `csharp.json` — `$.repository["object-creation-expression-with-parameters"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(new)(?:\s+(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))?\s*(?=\()
```

- Adapter translation:

```text
(new)(?:\s+(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*))?\s*(?=\()
```

### 120. `csharp.json` — `$.repository["operator-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*\b(?<operator_keyword>operator)\b\s*(?<operator>[-!%\&*+/<=>^|~]+|true|false)\s*(?=\()
```

- Adapter translation:

```text
(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s*\b(?<operator_keyword>operator)\b\s*(?<operator>[-!%&*+/<=>^|~]+|true|false)\s*(?=\()
```

### 121. `csharp.json` — `$.repository["parameter"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(ref|params|out|in|this)\b\s+)?(?<type_name>(?:ref\s+)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+(\g<identifier>)
```

- Adapter translation:

```text
(?:\b(ref|params|out|in|this)\b\s+)?(?<type_name>(?:ref\s+)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))
```

### 122. `csharp.json` — `$.repository["property-declaration"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?![[:word:]\s]*\b(?:class|interface|struct|union|enum|event)\b)(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>\g<type_name>\s*\.\s*)?(?<property_name>\g<identifier>)\s*(?=\{|=>|//|/\*|$)
```

- Adapter translation:

```text
(?![\p{L}\p{M}\p{N}\p{Pc}\s]*\b(?:class|interface|struct|union|enum|event)\b)(?<return_type>(?<type_name>(?:ref\s+(?:readonly\s+)?)?(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)\s+)(?<interface_name>(?!)\s*\.\s*)?(?<property_name>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s*(?=\{|=>|//|/\*|$)
```

### 123. `csharp.json` — `$.repository["query-expression"].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
\b(from)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s+(\g<identifier>)\b\s*\b(in)\b\s*
```

- Adapter translation:

```text
\b(from)\b\s*(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)?\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b\s*\b(in)\b\s*
```

### 124. `csharp.json` — `$.repository["tuple-deconstruction-assignment"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
(?<tuple>\s*\((?:[^()]|\g<tuple>)+\))\s*(?!=[=>])(?==)
```

- Adapter translation:

```text
(?<tuple>\s*\((?:[^()]|(?!))+\))\s*(?!=[=>])(?==)
```

### 125. `csharp.json` — `$.repository["tuple-element"].match`

- Scope: `source.cs`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type_args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type_name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name_and_type_args>\g<identifier>\s*(?<type_args>\s*<(?:[^<>]|\g<type_args>)+>\s*)?)(?:\s*\.\s*\g<name_and_type_args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)(?:(?<tuple_name>\g<identifier>)\b)?
```

- Adapter translation:

```text
(?<type_name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name_and_type_args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type_args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*\??\s*)*)(?:(?<tuple_name>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?
```

### 126. `go.json` — `$.repository["other_struct_interface_expressions"].patterns[1].match`

- Scope: `source.go`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;brackets&gt;&#34; was neutralized
- Original pattern:

```text
\b(?!(?:struct|interface)\b)([.\w]+)(?<brackets>\[(?:[^]\[]|\g<brackets>)*])?(?=\{)
```

- Adapter translation:

```text
\b(?!(?:struct|interface)\b)([.\w]+)(?<brackets>\[(?:[^]\[]|(?:\[(?:[^]\[]|(?!))*]))*])?(?=\{)
```

### 127. `go.json` — `$.repository["support_functions"].match`

- Scope: `source.go`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;brackets&gt;&#34; was neutralized
- Original pattern:

```text
(?:((?<=\.)\b\w+)|\b(\w+))(?<brackets>\[(?:[^]\[]|\g<brackets>)*])?(?=\()
```

- Adapter translation:

```text
(?:((?<=\.)\b\w+)|\b(\w+))(?<brackets>\[(?:[^]\[]|(?:\[(?:[^]\[]|(?!))*]))*])?(?=\()
```

### 128. `haskell.json` — `$.repository["adt_constructor"].patterns[1].end`

- Scope: `source.haskell`
- Field: `end`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;paren&gt;&#34; was neutralized
- Original pattern:

```text
(?:\G|^)\s*(?:(?<!')\b(['._\p{Ll}\p{Lu}\p{Lt}\d]+)|('?(?<paren>\((?:[^()]?|\g<paren>)*\)))|('?(?<brac>\((?:[^]\[]?|\g<brac>)*])))\s*(?:(?<![[\p{S}\p{P}]&&[^]"'(),;\[_`{}]])(:[[\p{S}\p{P}]&&[^]"'(),;\[_`{}]]*)|(`)([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)(`))|(?<!')\b([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)|(\()\s*(:[[\p{S}\p{P}]&&[^]"'(),;\[_`{}]]*)\s*(\))
```

- Adapter translation:

```text
(?:\G|(?:(?<![\s\S])|^(?=[\s\S])))\s*(?:(?<!')\b(['._\p{Ll}\p{Lu}\p{Lt}\d]+)|('?(?<paren>\((?:[^()]?|(?!))*\)))|('?(?<brac>\((?:[^]\[]?|(?!))*])))\s*(?:(?<!(?:(?=[^]"'(),;\[_`{}])[\p{S}\p{P}]))(:(?:(?=[^]"'(),;\[_`{}])[\p{S}\p{P}])*)|(`)([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)(`))|(?<!')\b([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)|(\()\s*(:(?:(?=[^]"'(),;\[_`{}])[\p{S}\p{P}])*)\s*(\))
```

### 129. `haskell.json` — `$.repository["fun_decl"].begin`

- Scope: `source.haskell`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;fn&gt;&#34; was neutralized
- Original pattern:

```text
^(\s*)(?<fn>(?:[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*#*|\(\s*(?!--+\))[[\p{S}\p{P}]&&[^]"'(),:;\[_`{}]][[\p{S}\p{P}]&&[^]"'(),;\[_`{}]]*\s*\))(?:\s*,\s*\g<fn>)?)\s*(?<![[\p{S}\p{P}]&&[^]"'),;_`}]])(::|∷)(?![[\p{S}\p{P}]&&[^"'(,;\[_`{]])
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))(\s*)(?<fn>(?:[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*#*|\(\s*(?!--+\))(?:(?=[^]"'(),:;\[_`{}])[\p{S}\p{P}])(?:(?=[^]"'(),;\[_`{}])[\p{S}\p{P}])*\s*\))(?:\s*,\s*(?!))?)\s*(?<!(?:(?=[^]"'),;_`}])[\p{S}\p{P}]))(::|∷)(?!(?:(?=[^"'(,;\[_`{])[\p{S}\p{P}]))
```

### 130. `haskell.json` — `$.repository["module_name"].match`

- Scope: `source.haskell`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;conid&gt;&#34; was neutralized
- Original pattern:

```text
(?<conid>[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(\.\g<conid>)?)
```

- Adapter translation:

```text
(?<conid>[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(\.(?!))?)
```

### 131. `kotlin.json` — `$.repository["class-declaration"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
\b(class|(?:fun\s+)?interface)\s+(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|\g<GROUP>)+>)?
```

- Adapter translation:

```text
\b(class|(?:fun\s+)?interface)\s+(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|(?!))+>)?
```

### 132. `kotlin.json` — `$.repository["function"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
\b(fun)\b\s*(?<GROUP><([^<>]|\g<GROUP>)+>)?\s*(?:(?:(\w+)\.)?(\b\w+\b|`[^`]+`))?
```

- Adapter translation:

```text
\b(fun)\b\s*(?<GROUP><([^<>]|(?!))+>)?\s*(?:(?:(\w+)\.)?(\b\w+\b|`[^`]+`))?
```

### 133. `kotlin.json` — `$.repository["function-call"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
\??\.?(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|\g<GROUP>)+>)?\s*(?=[({])
```

- Adapter translation:

```text
\??\.?(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|(?!))+>)?\s*(?=[({])
```

### 134. `kotlin.json` — `$.repository["type-alias"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
\b(typealias)\s+(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|\g<GROUP>)+>)?
```

- Adapter translation:

```text
\b(typealias)\s+(\b\w+\b|`[^`]+`)\s*(?<GROUP><([^<>]|(?!))+>)?
```

### 135. `kotlin.json` — `$.repository["type-annotation"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
(?<![:?]):\s*([?\w\s]|->|(?<GROUP>[(<]([^"'()<>]|\g<GROUP>)+[)>]))+
```

- Adapter translation:

```text
(?<![:?]):\s*([?\w\s]|->|(?<GROUP>[(<]([^"'()<>]|(?!))+[)>]))+
```

### 136. `kotlin.json` — `$.repository["variable-declaration"].match`

- Scope: `source.kotlin`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;GROUP&gt;&#34; was neutralized
- Original pattern:

```text
\b(va[lr])\b\s*(?<GROUP><([^<>]|\g<GROUP>)+>)?
```

- Adapter translation:

```text
\b(va[lr])\b\s*(?<GROUP><([^<>]|(?!))+>)?
```

### 137. `markdown.json` — `$.repository["bold"].begin`

- Scope: `text.html.markdown`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(?<open>(\*\*(?=\w)|(?<!\w)\*\*|(?<!\w)\b__))(?=\S)(?=(<[^>]*+>|(?<raw>`+)([^`]|(?!(?<!`)\k<raw>(?!`))`)*+\k<raw>|\\[-\]!#(-+.>\[\\_`{}]?+|\[((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+](( ?\[[^]]*+])|(\([\t ]*+<?(.*?)>?[\t ]*+((?<title>["'])(.*?)\k<title>)?\))))|(?!(?<=\S)\k<open>).)++(?<=\S)(?=__\b|\*\*)\k<open>)
```

- Adapter translation:

```text
(?<open>(\*\*(?=\w)|(?<!\w)\*\*|(?<!\w)\b__))(?=\S)(?=(?>(<(?>[^>]*)>|(?<raw>`+)(?>([^`]|(?!(?<!`)\k<raw>(?!`))`)*)\k<raw>|\\(?>[-\]!#(-+.>\[\\_`{}]?)|\[((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*)](( ?\[(?>[^]]*)])|(\((?>[\t ]*)<?(.*?)>?(?>[\t ]*)((?<title>["'])(.*?)\k<title>)?\))))|(?!(?<=\S)\k<open>).)+)(?<=\S)(?=__\b|\*\*)\k<open>)
```

### 138. `markdown.json` — `$.repository["image-inline"].match`

- Scope: `text.html.markdown`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(!\[)((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+)(])(\()[\t ]*((<)((?:\\[<>]|[^\n<>])*)(>)|((?<url>(?>[^()\s]+)|\(\g<url>*\))*))[\t ]*(?:((\().+?(\)))|((").+?("))|((').+?(')))?\s*(\))
```

- Adapter translation:

```text
(!\[)((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*))(])(\()[\t ]*((<)((?:\\[<>]|[^\n<>])*)(>)|((?<url>(?>[^()\s]+)|\((?!)*\))*))[\t ]*(?:((\().+?(\)))|((").+?("))|((').+?(')))?\s*(\))
```

### 139. `markdown.json` — `$.repository["image-ref"].match`

- Scope: `text.html.markdown`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(!\[)((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+)(]) ?(\[)(.*?)(])
```

- Adapter translation:

```text
(!\[)((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*))(]) ?(\[)(.*?)(])
```

### 140. `markdown.json` — `$.repository["italic"].begin`

- Scope: `text.html.markdown`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(?<open>(\*(?=\w)|(?<!\w)\*|(?<!\w)\b_))(?=\S)(?=(<[^>]*+>|(?<raw>`+)([^`]|(?!(?<!`)\k<raw>(?!`))`)*+\k<raw>|\\[-\]!#(-+.>\[\\_`{}]?+|\[((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+](( ?\[[^]]*+])|(\([\t ]*+<?(.*?)>?[\t ]*+((?<title>["'])(.*?)\k<title>)?\))))|\k<open>\k<open>|(?!(?<=\S)\k<open>).)++(?<=\S)(?=_\b|\*)\k<open>)
```

- Adapter translation:

```text
(?<open>(\*(?=\w)|(?<!\w)\*|(?<!\w)\b_))(?=\S)(?=(?>(<(?>[^>]*)>|(?<raw>`+)(?>([^`]|(?!(?<!`)\k<raw>(?!`))`)*)\k<raw>|\\(?>[-\]!#(-+.>\[\\_`{}]?)|\[((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*)](( ?\[(?>[^]]*)])|(\((?>[\t ]*)<?(.*?)>?(?>[\t ]*)((?<title>["'])(.*?)\k<title>)?\))))|\k<open>\k<open>|(?!(?<=\S)\k<open>).)+)(?<=\S)(?=_\b|\*)\k<open>)
```

### 141. `markdown.json` — `$.repository["link-inline"].match`

- Scope: `text.html.markdown`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(\[)((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+)(])(\()[\t ]*((<)((?:\\[<>]|[^\n<>])*)(>)|((?<url>(?>[^()\s]+)|\(\g<url>*\))*))[\t ]*(?:((\()[^()]*(\)))|((")[^"]*("))|((')[^']*(')))?\s*(\))
```

- Adapter translation:

```text
(\[)((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*))(])(\()[\t ]*((<)((?:\\[<>]|[^\n<>])*)(>)|((?<url>(?>[^()\s]+)|\((?!)*\))*))[\t ]*(?:((\()[^()]*(\)))|((")[^"]*("))|((')[^']*(')))?\s*(\))
```

### 142. `markdown.json` — `$.repository["link-ref"].match`

- Scope: `text.html.markdown`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(?<![]\\])(\[)((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+)(])(\[)([^]]*+)(])
```

- Adapter translation:

```text
(?<![]\\])(\[)((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*))(])(\[)((?>[^]]*))(])
```

### 143. `markdown.json` — `$.repository["link-ref-literal"].match`

- Scope: `text.html.markdown`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;square&gt;&#34; was neutralized
- Original pattern:

```text
(?<![]\\])(\[)((?<square>[^]\[\\]|\\.|\[\g<square>*+])*+)(]) ?(\[)(])
```

- Adapter translation:

```text
(?<![]\\])(\[)((?>(?<square>[^]\[\\]|\\.|\[(?>(?:[^]\[\\]|\\.|\[(?>(?!)*)])*)])*))(]) ?(\[)(])
```

### 144. `mdx.json` — `$.repository["commonmark-definition"].match`

- Scope: `source.mdx`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;destination_raw&gt;&#34; was neutralized
- Original pattern:

```text
(?:^|\G)[\t ]*(\[)((?:[^]\[\\]|\\[]\[\\]?)+?)(])(:)[\t ]*(?:(<)((?:[^\n<>\\]|\\[<>\\]?)*)(>)|(\g<destination_raw>))(?:[\t ]+(?:(")((?:[^"\\]|\\["\\]?)*)(")|(')((?:[^'\\]|\\['\\]?)*)(')|(\()((?:[^)\\]|\\[)\\]?)*)(\))))?$(?<destination_raw>(?!<)(?:(?:[^ ()\\\p{Cc}]|\\[()\\]?)|\(\g<destination_raw>*\))+){0}
```

- Adapter translation:

```text
(?:(?:(?<![\s\S])|^(?=[\s\S]))|\G)[\t ]*(\[)((?:[^]\[\\]|\\[]\[\\]?)+?)(])(:)[\t ]*(?:(<)((?:[^\n<>\\]|\\[<>\\]?)*)(>)|((?!)))(?:[\t ]+(?:(")((?:[^"\\]|\\["\\]?)*)(")|(')((?:[^'\\]|\\['\\]?)*)(')|(\()((?:[^)\\]|\\[)\\]?)*)(\))))?$((?!)){0}
```

### 145. `mdx.json` — `$.repository["commonmark-label-end"].patterns[0].match`

- Scope: `source.mdx`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;destination_raw&gt;&#34; was neutralized
- Original pattern:

```text
(])(\()[\t ]*(?:(?:(<)((?:[^\n<>\\]|\\[<>\\]?)*)(>)|(\g<destination_raw>))(?:[\t ]+(?:(")((?:[^"\\]|\\["\\]?)*)(")|(')((?:[^'\\]|\\['\\]?)*)(')|(\()((?:[^)\\]|\\[)\\]?)*)(\))))?)?[\t ]*(\))(?<destination_raw>(?!<)(?:(?:[^ ()\\\p{Cc}]|\\[()\\]?)|\(\g<destination_raw>*\))+){0}
```

- Adapter translation:

```text
(])(\()[\t ]*(?:(?:(<)((?:[^\n<>\\]|\\[<>\\]?)*)(>)|((?!)))(?:[\t ]+(?:(")((?:[^"\\]|\\["\\]?)*)(")|(')((?:[^'\\]|\\['\\]?)*)(')|(\()((?:[^)\\]|\\[)\\]?)*)(\))))?)?[\t ]*(\))((?!)){0}
```

### 146. `mdx.json` — `$.repository["extension-gfm-autolink-literal"].patterns[0].match`

- Scope: `source.mdx`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;path&gt;&#34; was neutralized
- Original pattern:

```text
(?<=^|[]\t\n\r (*\[_~])(?=(?i:www)\.[^\n\r])(?:(?:[-\p{L}\p{N}]|[._](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))+\g<path>?)?(?<path>(?:(?:[^]\t\n\r !"\&-*,.:;<?_~]|&(?![A-Za-z]*;[!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[]))|[!"')*,.:;?_~](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))|\(\g<path>*\))+){0}
```

- Adapter translation:

```text
(?<=(?:(?<![\s\S])|^(?=[\s\S]))|[]\t\n\r (*\[_~])(?=(?i:www)\.[^\n\r])(?:(?:[-\p{L}\p{N}]|[._](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))+(?!)?)?((?!)){0}
```

### 147. `mdx.json` — `$.repository["extension-gfm-autolink-literal"].patterns[1].match`

- Scope: `source.mdx`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;path&gt;&#34; was neutralized
- Original pattern:

```text
(?<=^|[^A-Za-z])(?i:https?://)(?=[\p{L}\p{N}])(?:(?:[-\p{L}\p{N}]|[._](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))+\g<path>?)?(?<path>(?:(?:[^]\t\n\r !"\&-*,.:;<?_~]|&(?![A-Za-z]*;[!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[]))|[!"')*,.:;?_~](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))|\(\g<path>*\))+){0}
```

- Adapter translation:

```text
(?<=(?:(?<![\s\S])|^(?=[\s\S]))|[^A-Za-z])(?i:https?://)(?=[\p{L}\p{N}])(?:(?:[-\p{L}\p{N}]|[._](?![!"')*,.:;<?_~]*(?:[<\s]|][\t\n (\[])))+(?!)?)?((?!)){0}
```

### 148. `powershell.json` — `$.repository["unicodeEscape"].patterns[0].match`

- Scope: `source.powershell`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;1&gt;&#34; was neutralized
- Original pattern:

```text
`u\{(?:(?:10)?(\h){1,4}|0?\g<1>{1,5})}
```

- Adapter translation:

```text
`u\{(?:(?:10)?([0-9A-Fa-f]){1,4}|0?(?!){1,5})}
```

### 149. `purescript.json` — `$.repository["double_colon_parens"].patterns[0].match`

- Scope: `source.purescript`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;paren&gt;&#34; was neutralized
- Original pattern:

```text
\((?<paren>(?:[^()]|\(\g<paren>\))*)(::|∷)(?<paren2>(?:[^()}]|\(\g<paren2>\))*)\)
```

- Adapter translation:

```text
\((?<paren>(?:[^()]|\((?!)\))*)(::|∷)(?<paren2>(?:[^()}]|\((?!)\))*)\)
```

### 150. `purescript.json` — `$.repository["type_signature"].patterns[1].match`

- Scope: `source.purescript`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: duplicate named capture &#34;classConstraint&#34; is not supported
- Original pattern:

```text
\((?<classConstraints>([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*)\s+(?<classConstraint>(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*|(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*\.)?[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)(?:\s*\s+\s*(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*|(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*\.)?[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*))*)(?:\s*,\s*([\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*)\s+(?<classConstraint>(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*|(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*\.)?[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)(?:\s*\s+\s*(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*|(?:[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*(?:\.[\p{Lu}\p{Lt}]['_\p{Ll}\p{Lu}\p{Lt}\d]*)*\.)?[_\p{Ll}]['_\p{Ll}\p{Lu}\p{Lt}\d]*))*))*)\)\s*(=>|<=|[⇐⇒])
```

### 151. `razor.json` — `$.repository["catch-condition"].patterns[0].match`

- Scope: `text.aspnetcorerazor`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type-args&gt;&#34; was neutralized
- Original pattern:

```text
(?<type-name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name-and-type-args>\g<identifier>\s*(?<type-args>\s*<(?:[^<>]|\g<type-args>)+>\s*)?)(?:\s*\.\s*\g<name-and-type-args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:(\g<identifier>)\b)?
```

- Adapter translation:

```text
(?<type-name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name-and-type-args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type-args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?
```

### 152. `razor.json` — `$.repository["foreach-condition"].patterns[0].match`

- Scope: `text.aspnetcorerazor`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;type-args&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(var)\b|(?<type-name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name-and-type-args>\g<identifier>\s*(?<type-args>\s*<(?:[^<>]|\g<type-args>)+>\s*)?)(?:\s*\.\s*\g<name-and-type-args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+(\g<identifier>)\s+\b(in)\b
```

- Adapter translation:

```text
(?:\b(var)\b|(?<type-name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name-and-type-args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type-args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s+\b(in)\b
```

### 153. `razor.json` — `$.repository["foreach-condition"].patterns[1].match`

- Scope: `text.aspnetcorerazor`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;tuple&gt;&#34; was neutralized
- Original pattern:

```text
(?:\b(var)\b\s*)?(?<tuple>\((?:[^()]|\g<tuple>)+\))\s+\b(in)\b
```

- Adapter translation:

```text
(?:\b(var)\b\s*)?(?<tuple>\((?:[^()]|(?!))+\))\s+\b(in)\b
```

### 154. `stata.json` — `$.repository["unicode-regex-internals"].patterns[1].match`

- Scope: `source.stata`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma character-class intersection is not supported
- Original pattern:

```text
\$(?![,013_{|}[\w&&[^0-9_]]\w])
```

- Adapter translation:

```text
\$(?!
```

### 155. `swift.json` — `$.repository["literals-regular-expression-literal"].patterns[1].match`

- Scope: `source.swift`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;guts&gt;&#34; was neutralized
- Original pattern:

```text
(/)(?!\s)(?!/)(?:\\\s(?=/)|(?<guts>(?>(?:\\Q(?:(?!\\E)(?!/).)*+(?:\\E|(?=/))|\\.|\(\?#[^)]*\)|\(\?(?>\{(?:[^{].*?|\{[^{].*?}|\{\{[^{].*?}}|\{\{\{[^{].*?}}}|\{\{\{\{[^{].*?}}}}|\{\{\{\{\{.+?}}}}})})(?:\[(?!\d)\w+])?[<>X]?\)|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\])+])+])+])+]|\(\g<guts>?+\)|(?:(?!/)[^()\[\\])+)+))?+(?<!\s))(/)
```

- Adapter translation:

```text
(/)(?!\s)(?!/)(?:\\\s(?=/)|(?>(?<guts>(?>(?:\\Q(?>(?:(?!\\E)(?!/).)*)(?:\\E|(?=/))|\\.|\(\?#[^)]*\)|\(\?(?>\{(?:[^{].*?|\{[^{].*?}|\{\{[^{].*?}}|\{\{\{[^{].*?}}}|\{\{\{\{[^{].*?}}}}|\{\{\{\{\{.+?}}}}})})(?:\[(?!\d)\w+])?[<>X]?\)|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\])+])+])+])+]|\((?>(?!)?)\)|(?:(?!/)[^()\[\\])+)+))?)(?<!\s))(/)
```

### 156. `swift.json` — `$.repository["literals-regular-expression-literal"].patterns[2].match`

- Scope: `source.swift`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;guts&gt;&#34; was neutralized
- Original pattern:

```text
((#+)/)(?<guts>(?>(?:\\Q(?:(?!\\E)(?!/\2).)*+(?:\\E|(?=/\2))|\\.|\(\?#[^)]*\)|\(\?(?>\{(?:[^{].*?|\{[^{].*?}|\{\{[^{].*?}}|\{\{\{[^{].*?}}}|\{\{\{\{[^{].*?}}}}|\{\{\{\{\{.+?}}}}})})(?:\[(?!\d)\w+])?[<>X]?\)|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\])+])+])+])+]|\(\g<guts>?+\)|(?:(?!/\2)[^()\[\\])+)+))?+(/\2)|#+/.+(\n)
```

- Adapter translation:

```text
((#+)/)(?>(?<guts>(?>(?:\\Q(?>(?:(?!\\E)(?!/\2).)*)(?:\\E|(?=/\2))|\\.|\(\?#[^)]*\)|\(\?(?>\{(?:[^{].*?|\{[^{].*?}|\{\{[^{].*?}}|\{\{\{[^{].*?}}}|\{\{\{\{[^{].*?}}}}|\{\{\{\{\{.+?}}}}})})(?:\[(?!\d)\w+])?[<>X]?\)|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\]|\[(?:\\.|[^]\[\\])+])+])+])+]|\((?>(?!)?)\)|(?:(?!/\2)[^()\[\\])+)+))?)(/\2)|#+/.+(\n)
```

### 157. `swift.json` — `$.repository["literals-regular-expression-literal-callout"].match`

- Scope: `source.swift`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: duplicate named capture &#34;name&#34; is not supported
- Original pattern:

```text
(\()(?<keyw>\?C)(?:(?<num>\d+)|`(?<name>(?:[^`]|``)*)`|'(?<name>(?:[^']|'')*)'|"(?<name>(?:[^"]|"")*)"|\^(?<name>(?:[^^]|\^\^)*)\^|%(?<name>(?:[^%]|%%)*)%|#(?<name>(?:[^#]|##)*)#|\$(?<name>(?:[^$]|\$\$)*)\$|\{(?<name>(?:[^}]|}})*)})?(\))|(\()(?<keyw>\*)(?<name>(?!\d)\w+)(?:\[(?<tag>(?!\d)\w+)])?(?:\{[^,}]+(?:,[^,}]+)*})?(\))|(\()(?<keyw>\?)(?>(\{(?:\g<20>|(?!\{).*?)}))(?:\[(?<tag>(?!\d)\w+)])?(?<keyw>[<>X]?)(\))
```

- Adapter translation:

```text
(\()(?<keyw>\?C)(?:(?<num>\d+)|`(?<name>(?:[^`]|``)*)`|'(?<name>(?:[^']|'')*)'|"(?<name>(?:[^"]|"")*)"|\^(?<name>(?:[^^]|\^\^)*)\^|%(?<name>(?:[^%]|%%)*)%|#(?<name>(?:[^#]|##)*)#|\$(?<name>(?:[^$]|\$\$)*)\$|\{(?<name>(?:[^}]|}})*)})?(\))|(\()(?<keyw>\*)(?<name>(?!\d)\w+)(?:\[(?<tag>(?!\d)\w+)])?(?:\{[^,}]+(?:,[^,}]+)*})?(\))|(\()(?<keyw>\?)(?>(\{(?:(?!)|(?!\{).*?)}))(?:\[(?<tag>(?!\d)\w+)])?(?<keyw>[<>X]?)(\))
```

### 158. `swift.json` — `$.repository["literals-regular-expression-literal-group-or-conditional"].patterns[1].begin`

- Scope: `source.swift`
- Field: `begin`
- Kind: `unsupported_syntax`
- Reason: duplicate named capture &#34;num&#34; is not supported
- Original pattern:

```text
(\()(?<cond>\?\()(?:(?<NumberRef>(?<num>[-+]?\d+)(?:(?<op>[-+])(?<num>\d+))?)|(?<cond>R)\g<NumberRef>?|(?<cond>R&)(?<NamedRef>(?<name>(?!\d)\w+)(?:(?<op>[-+])(?<num>\d+))?)|(?<cond><)(?:\g<NamedRef>|\g<NumberRef>)(?<cond>>)|(?<cond>')(?:\g<NamedRef>|\g<NumberRef>)(?<cond>')|(?<cond>DEFINE)|(?<cond>VERSION)(?<compar>>?=)(?<num>\d+\.\d+))(?<cond>\))|(\()(?<cond>\?)(?=\()
```

- Adapter translation:

```text
(\()(?<cond>\?\()(?:(?<NumberRef>(?<num>[-+]?\d+)(?:(?<op>[-+])(?<num>\d+))?)|(?<cond>R)(?:(?<num>[-+]?\d+)(?:(?<op>[-+])(?<num>\d+))?)?|(?<cond>R&)(?<NamedRef>(?<name>(?!\d)\w+)(?:(?<op>[-+])(?<num>\d+))?)|(?<cond><)(?:(?:(?<name>(?!\d)\w+)(?:(?<op>[-+])(?<num>\d+))?)|(?:(?<num>[-+]?\d+)(?:(?<op>[-+])(?<num>\d+))?))(?<cond>>)|(?<cond>')(?:(?:(?<name>(?!\d)\w+)(?:(?<op>[-+])(?<num>\d+))?)|(?:(?<num>[-+]?\d+)(?:(?<op>[-+])(?<num>\d+))?))(?<cond>')|(?<cond>DEFINE)|(?<cond>VERSION)(?<compar>>?=)(?<num>\d+\.\d+))(?<cond>\))|(\()(?<cond>\?)(?=\()
```

### 159. `wolfram.json` — `$.repository["simple-toplevel-definitions"].patterns[1].match`

- Scope: `source.wolfram`
- Field: `match`
- Kind: `unsupported_syntax`
- Reason: oniguruma subroutine call &#34;\\g&lt;3&gt;&#34; was neutralized
- Original pattern:

```text
^\s*(`?(?:[$[:alpha:]][$[:alnum:]]*`)*)([$[:alpha:]][$[:alnum:]]*)(?=\s*(\[(?>[^]\[]+|\g<3>)*])\s*(?:/;.*)?(?::=|=(?![!.=])))
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))\s*(`?(?:[$\p{L}\p{Nl}\p{Other_Alphabetic}][$\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*`)*)([$\p{L}\p{Nl}\p{Other_Alphabetic}][$\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)(?=\s*(\[(?>[^]\[]+|(?!))*])\s*(?:/;.*)?(?::=|=(?![!.=])))
```

## Pattern compile failures

### 1. `csharp.json` — `$.repository["comment"].patterns[1].patterns[0].begin`

- Scope: `source.cs`
- Field: `begin`
- Kind: `compile_error`
- Reason: error parsing regexp: unrecognized grouping construct: (?~ in `(?!)(?=(?~\*/)$)`
- Original pattern:

```text
\G(?=(?~\*/)$)
```

### 2. `csharp.json` — `$.repository["comment"].patterns[1].patterns[0].while`

- Scope: `source.cs`
- Field: `while`
- Kind: `compile_error`
- Reason: error parsing regexp: unrecognized grouping construct: (?~ in `(?:(?&lt;![\s\S])|^(?=[\s\S]))((?&gt;\s*))(\*(?!/))?(?=(?~\*/)$)`
- Original pattern:

```text
^(\s*+)(\*(?!/))?(?=(?~\*/)$)
```

- Adapter translation:

```text
(?:(?<![\s\S])|^(?=[\s\S]))((?>\s*))(\*(?!/))?(?=(?~\*/)$)
```

### 3. `razor.json` — `$.repository["catch-condition"].patterns[0].match`

- Scope: `text.aspnetcorerazor`
- Field: `match`
- Kind: `compile_error`
- Reason: error parsing regexp: reference to undefined group name and in `(?&lt;type-name&gt;(?:(?:(?&lt;identifier&gt;@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?&lt;name-and-type-args&gt;(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?&lt;type-args&gt;\s*&lt;(?:[^&lt;&gt;]|(?!))+&gt;\s*)?)(?:\s*\.\s*(?!))*|(?&lt;tuple&gt;\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?`
- Original pattern:

```text
(?<type-name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name-and-type-args>\g<identifier>\s*(?<type-args>\s*<(?:[^<>]|\g<type-args>)+>\s*)?)(?:\s*\.\s*\g<name-and-type-args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:(\g<identifier>)\b)?
```

- Adapter translation:

```text
(?<type-name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name-and-type-args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type-args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*)\s*(?:((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\b)?
```

### 4. `razor.json` — `$.repository["foreach-condition"].patterns[0].match`

- Scope: `text.aspnetcorerazor`
- Field: `match`
- Kind: `compile_error`
- Reason: error parsing regexp: reference to undefined group name and in `(?:\b(var)\b|(?&lt;type-name&gt;(?:(?:(?&lt;identifier&gt;@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?&lt;name-and-type-args&gt;(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?&lt;type-args&gt;\s*&lt;(?:[^&lt;&gt;]|(?!))+&gt;\s*)?)(?:\s*\.\s*(?!))*|(?&lt;tuple&gt;\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s+\b(in)\b`
- Original pattern:

```text
(?:\b(var)\b|(?<type-name>(?:(?:(?<identifier>@?[_[:alpha:]][_[:alnum:]]*)\s*::\s*)?(?<name-and-type-args>\g<identifier>\s*(?<type-args>\s*<(?:[^<>]|\g<type-args>)+>\s*)?)(?:\s*\.\s*\g<name-and-type-args>)*|(?<tuple>\s*\((?:[^()]|\g<tuple>)+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+(\g<identifier>)\s+\b(in)\b
```

- Adapter translation:

```text
(?:\b(var)\b|(?<type-name>(?:(?:(?<identifier>@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*::\s*)?(?<name-and-type-args>(?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*)\s*(?<type-args>\s*<(?:[^<>]|(?!))+>\s*)?)(?:\s*\.\s*(?!))*|(?<tuple>\s*\((?:[^()]|(?!))+\)))(?:\s*\?\s*)?(?:\s*\[(?:\s*,\s*)*]\s*)*))\s+((?:@?[_\p{L}\p{Nl}\p{Other_Alphabetic}][_\p{L}\p{Nl}\p{Other_Alphabetic}\p{Nd}]*))\s+\b(in)\b
```

## Other regex diagnostics

None.
