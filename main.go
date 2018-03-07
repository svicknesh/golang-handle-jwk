package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/crypto/ssh/terminal"
)

/*
	make sure to run
		go get -u golang.org/x/sys/unix
		go get -u github.com/lestrrat-go/jwx/jwk
	before building this app
*/

func main() {

	keySize := flag.Int("keysize", 4096, "Key size for the generated keys.")
	outputFile := flag.String("output-file", "handle", "Output key basename. -private, -public and -jwk will be appended.")
	noPasswd := flag.Bool("no-password", false, "Do not encrypt generated private key.")

	flag.Parse()

	var password1 string

	if !*noPasswd {
		// private key is meant to be encrypted

		fmt.Print("Enter Password: ")
		bytePassword1, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err.Error())
		}
		password1 = strings.TrimSpace(string(bytePassword1))

		fmt.Println()
		fmt.Print("Re-enter Password: ")
		bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err.Error())
		}
		password2 := strings.TrimSpace(string(bytePassword2))

		fmt.Println()
		if password1 != password2 {
			fmt.Println("Password entered does not match.")
			os.Exit(1)
		}
	}

	privkey, err := rsa.GenerateKey(rand.Reader, *keySize)
	if err != nil {
		fmt.Printf("failed to generate private key: %s\n", err.Error())
		os.Exit(2)
	}

	// Convert it to pem
	blockPrivate := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privkey),
	}

	blockPublic := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privkey.PublicKey),
	}

	// Encrypt the pem
	if len(password1) > 0 {
		blockPrivate, err = x509.EncryptPEMBlock(rand.Reader, blockPrivate.Type, blockPrivate.Bytes, []byte(password1), x509.PEMCipherAES256)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(3)
		}
	}

	pemBytesPrivate := pem.EncodeToMemory(blockPrivate)
	pemBytesPublic := pem.EncodeToMemory(blockPublic)

	outputFilePrivate := fmt.Sprintf("%s-private.pem", *outputFile)
	outputFilePublic := fmt.Sprintf("%s-public.pem", *outputFile)

	ioutil.WriteFile(outputFilePrivate, pemBytesPrivate, 0600)
	ioutil.WriteFile(outputFilePublic, pemBytesPublic, 0600)

	key, err := jwk.New(&privkey.PublicKey)
	if err != nil {
		fmt.Printf("failed to create JWK: %s\n", err.Error())
		os.Exit(4)
	}

	jsonbuf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		fmt.Printf("failed to generate JSON: %s\n", err)
		os.Exit(5)
	}
	jsonbuf = append(jsonbuf, []byte("\n")...)

	//os.Stdout.Write(jsonbuf)

	outputFileJWK := fmt.Sprintf("%s-jwk.txt", *outputFile)
	ioutil.WriteFile(outputFileJWK, jsonbuf, 0600)

	fmt.Printf("Done generating files %s, %s and %s\n", outputFilePrivate, outputFilePublic, outputFileJWK)
}
