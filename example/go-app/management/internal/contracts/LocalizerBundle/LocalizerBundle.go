package LocalizerBundle

import (
	"embed"
)

//go:embed locale.*.toml
//go:generate go run ../../tools/gen_string_enum.go -input=locale.en.toml -type=LocaleKey -output=locale_enum.go
var LocalFS embed.FS
var LocaleFiles = []string{
	"locale.en.toml",
}
