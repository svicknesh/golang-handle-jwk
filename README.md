# Handle Key Generator

Generates private & public key in PEM format and public key in JSON Web Key (JWK) format for inclusion in an HS_ADMIN record.

## Pre-requisites

- Golang 1.10
- Go libraries
  - `go get -u golang.org/x/sys/unix`
  - `go get -u github.com/lestrrat-go/jwx/jwk`

## Building application

The beauty of Go is it allows a binary to be generated for any of the supported platforms from a single development machine.

### Building for Linux

`env GOOS=linux GOARCH=amd64 go build -o golang-handle-linux`

### Building for OSX

`env GOOS=darwin GOARCH=amd64 go build -o golang-handle-osx`

### Building for Windows

`env GOOS=windows GOARCH=amd64 go build -o golang-handle.exe`

## Releases

The binary releases for this application are available at the following [link](https://github.com/svicknesh/golang-handle-jwk/releases).

## Application usage

```bash
  -keysize int
    	Key size for the generated keys. (default 4096)
  -no-password
    	Do not encrypt generated private key.
  -output-file string
    	Output key basename. -private, -public and -jwk will be appended. (default "handle")
```

## Example

### Generating a passwordless private key

`./golang-handle-osx --no-password`

### Specify basename for generated keys

`./golang-handle-osx --output-file vicknesh`

### Specify keysize for generated keys

`./golang-handle-osx --keysize 2048`

## License & Warranty

Feel free to use this application and source as you see fit. This application does not come with any warranty, implied or otherwise.
