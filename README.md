# 2048

2048 cross-platform opengl.

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

![screenshot](screenshots/2048.png)

