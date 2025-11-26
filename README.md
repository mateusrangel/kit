# Kit [![Go Reference](https://pkg.go.dev/badge/github.com/mateusrangel/kit.svg)](https://pkg.go.dev/github.com/mateusrangel/kit)
A set of common packages for building resilient applications in Go.



## Development
### Run tests and generate coverage report
```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
```