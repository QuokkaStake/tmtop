package embed

import (
	"embed"
)

//go:embed frontend/dist
var Frontend embed.FS
