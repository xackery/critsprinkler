
goversioninfo -icon=assets/critsprinkler.ico -manifest=critsprinkler.exe.manifest -o=rsrc.syso versioninfo.json
go build -ldflags "-X main.Version=dev" || exit /b
move critsprinkler.exe bin/critsprinkler.exe || exit /b
cd bin || exit /b
critsprinkler.exe c:\games\eq\thj\Logs\eqlog_Shin_thj.txt || exit /b