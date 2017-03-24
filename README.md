# Symphony Go Code Sample

## Pre requirements ##
Requirements for this project are localy installed Golang and mongodb.
You should create database with any name and create one user with read/write permisions.
(Later in examples I used database called symphonyGopher and user with credentials user:password on port 27017)

## Configuration ##

Application configuration is stored in environment variables.

Example script to configure environment (added to $HOME/.bash_profile):
- Go variables
`export GOPATH="/Users/stefanilijevski/go"`
`export GOBIN="/Users/stefanilijevski/go/bin"`

- SympohonyGopher variables
`export SG_DB="mongodb://user:password@localhost:27017/symphonyGopher"`
`export SG_PORT=8082`

## Installation ##

Application uses vendoring, and [govendor](https://github.com/kardianos/govendor) tool to manage dependencies.
Dependencies are kept into vendor/ folder, aren't version controlled. Folders under vendor/ folder are ignored.

### Steps ###

- Install govendor tool `go get -u github.com/kardianos/govendor`
- Run `govendor get github.com/IlijevskiSymphony/symphonyGopher`
- Run `govendor sync` in project folder
- To install project, run `go install`
- Change directory to bin: `cd $GOPATH/bin`
- Run application server with `./symphonyGopher`


### Examples of use ###
Server endpoints and their purpose:
- '/register' (POST)
  Used for registration of new user, expects next form data parameters in JSON format.
  {
  	"email":"random123@gmail.com",
  	"password": "Ilijevski1234"
  }

- '/signup/accept/{reference uuid}' (GET)
  Used to confirm registration and should be sent by email. Reference uuid is stored in registrations document and is referenced to user by 'partnerID'

- '/login' (GET)
  Used to check if there is current session active. Session is stored in cookie.

- '/login' (POST)
  Used to login to server and create new session, expects next form data parameters in JSON format.
  {
  	"email":"random123@gmail.com",
  	"password": "Ilijevski1234"
  }

- '/logout' (GET)
  Used to terminate user session

- '/addLink' (POST)
  Add link with tags to repository for logged in user, expects next form data parameters in JSON format.
  {
  	"link":"www.symphony.is",
  	"tags": [
              "da",
              "ud"
            ]
  }

- '/verification/resend' (POST)
  Resends verification email for registrated credentials, expects next form data parameters in JSON format.
  {
    "email":"random123@gmail.com"
  }

- '/resetPasswordEmail' (POST)
  Sends an email with reference to user to reset password, expects next form data parameters in JSON format.
  {
    "email":"random123@gmail.com"
  }

- '/resetPasswordCheck/{reference}' (GET)
  Checks if reference used from reset password email is valid

- '/resetPassword' (POST)
  Resets password for user that used link from the email, expects next form data parameters in JSON format.
  {
    "password":"newPassword123"
    "id":"e1955c74-8211-43fb-a1d1-b43e24ef6d1e"
  }
