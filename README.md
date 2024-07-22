# README

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Installing Wails

To install Wails, you first need gcc. Download the newest version here: https://winlibs.com/#download-release  
Then extract the archive into C:/ and add the the path (for example C:\mingw64\bin) to the Windows Environment Variables.  
Now proceed to install NPM: https://nodejs.org/en/download/prebuilt-installer  
After that simply run the following command:
```
go install github.com/wailsapp/wails/v2/cmd/wails@lates
```

Configure wails by running  
```
wails dev
```

To debug the program, some changes to the build config are needed. You need to add the following variables to the build environment:  
```
GOARCH=amd64;GOOS=windows;CGO_ENABLED=1
```

In case the debug isn't working, you can generate the build config for your ide with the following files:  
Goland: ``wails init -n ComicCollector -ide goland``  
VSCode: ``wails init -n ComicCollector -ide vscode``