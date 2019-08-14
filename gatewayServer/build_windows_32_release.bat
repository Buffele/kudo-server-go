set GOARCH=386
set GOOS=windows
cls & go build -o server.exe -ldflags "-s -w" main.go & cmd /k