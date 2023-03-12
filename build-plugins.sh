# Build Engine Plugins
go build -buildmode=plugin -o engine/ebiten/plugin.so engine/ebiten/app.go
go build -buildmode=plugin -o engine/g3n/plugin.so engine/g3n/app.go
go build -buildmode=plugin -o engine/raylib/plugin.so engine/raylib/app.go