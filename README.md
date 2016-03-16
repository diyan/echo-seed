echo_seed
=========

Seed project for RESTful API powered by LabStack's Echo. http://labstack.com/echo/

Build
-----
```bash
go get github.com/labstack/echo
go get github.com/stretchr/testify
go get github.com/smartystreets/goconvey
```

Run
---
```bash
go run main.go
```

Continuous testing
------------------
```bash
$GOPATH/bin/goconvey &
firefox http://localhost:8080
```