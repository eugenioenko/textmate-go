module github.com/eugenioenko/textmate-go/benchmarks/chroma

go 1.25

require (
	github.com/alecthomas/chroma/v2 v2.24.1
	github.com/eugenioenko/textmate-go v0.0.0
)

require (
	github.com/dlclark/regexp2 v1.12.0 // indirect
	github.com/dlclark/regexp2/v2 v2.8.0 // indirect
)

replace github.com/eugenioenko/textmate-go => ../..

replace github.com/dlclark/regexp2/v2 => ../../../regexp2
