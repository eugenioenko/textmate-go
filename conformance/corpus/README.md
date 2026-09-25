# Differential corpus

This corpus combines compact regression files written for textmate-go with the
real language samples maintained by the pinned `textmate-grammars-themes`
project. The local files are licensed under this repository's MIT license. The
external samples are read from the exact clean sibling checkout recorded in
`manifest.json`; Shiki documents all files outside its redistributed grammars
and themes as MIT under its root `LICENSE`, so no third-party source text is
copied into this repository.

`manifest.json` pins every file to a tm-grammars grammar name so extension
guessing cannot change the oracle. In particular, the corpus includes the
state-carrying cases behind ttt#670: multiline HTML and TSX attributes, JSX
expressions, block comments, JavaScript template literals, and Python
docstrings. The 40 pinned upstream samples add broader real-file syntax for
every embedded root grammar. YAML and TOML remain local supplemental cases because
they are intentionally outside the license-reviewed embedded set.
