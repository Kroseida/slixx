nodejs 16.18.1
yarn serve

using:
https://github.com/valyala/bytebufferpool/

Testing:
go test ./... -coverprofile coverage.out
go tool cover -html -func=coverage.out
