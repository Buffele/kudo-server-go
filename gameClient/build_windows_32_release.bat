set GOARCH=386
set GOOS=windows
cls & go build -o client.exe -ldflags "-s -w" main.go & cmd /k