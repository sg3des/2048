# 2048

Cross-platform 2048 game on opengl.

build on linux: 

```sh
make build
#or
go build
```

build on windows:

```
set CGO_ENABLED=1
set GOARCH=386
go build -ldflags -H=windowsgui
```

build on osx:
```sh
go build -ldflags=-s
```

![screenshot](screenshots/2048.png)

