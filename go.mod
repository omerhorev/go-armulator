module github.com/omerhorev/goarmulator

go 1.18

require github.com/pkg/errors v0.9.1

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/omerhorev/go-armulator => ./

// replace github.com/omerhorev/go-armulator/elf => ./src/elf
