module root

go 1.24.2

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/rs/zerolog v1.34.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/rs/zerolog v1.34.0 => github.com/Gz3zFork/zerolog v0.0.0-20250403082219-16986a1f50be
