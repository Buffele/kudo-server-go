set GOARCH=amd64
set GOOS=linux
cls & go build -o server main.go & cmd /k