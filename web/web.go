// Package web contiene la página estática que se sirve en "/".
package web

import _ "embed"

//go:embed index.html
var IndexHTML []byte
