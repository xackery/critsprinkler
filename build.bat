
goversioninfo -icon=library/assets/critsprinkler.ico -manifest=critsprinkler.exe.manifest -o=rsrc.syso versioninfo.json
go build -ldflags "-H windowsgui -s -w -X main.Version=0.1.1.2" || exit /b
move critsprinkler.exe bin/critsprinkler.exe || exit /b
cd bin || exit /b