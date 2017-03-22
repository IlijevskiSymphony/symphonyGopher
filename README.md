# Symphony Go Code Sample

## Configuration ##

Application configuration is stored in environment variables.

Example script to configure environment (added to $HOME/.bash_profile):
`export NODE_PATH="/usr/local/lib/node_modules"`
`export GOPATH="/Users/stefanilijevski/go"`
`export GOBIN="/Users/stefanilijevski/go/bin"`
`export SG_DB="mongodb://user:password@localhost:27017/symphonyGopher"`
`export SG_PORT=8082`

## Installation ##

Application uses vendoring, and [govendor](https://github.com/kardianos/govendor) tool to manage dependencies.
Dependencies are kept into vendor/ folder, aren't version controlled. Folders under vendor/ folder are ignored.

### Steps ###

- Install govendor tool `go get -u github.com/kardianos/govendor`
- Run `govendor get github.com/IlijevskiSymphony/symphonyGopher`
- To install project, run `go install`
- Change directory to bin: `cd $GOPATH/bin`
- Run application server with `./symphonyGopher`
