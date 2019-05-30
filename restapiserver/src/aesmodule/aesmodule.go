/*
	Selected module was created on OPENSSL-based logic.
Openssl was selected because it's free and used by everyone for generate some SSL-keys.
Openssl can enc/dec strings, binary data, files and etc.
If it needed - this module can be upgraded for new functional.

	link : https://www.openssl.org/
*/
package aesmodule

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

// Types of AES-based encryption modes
const (
	TYPE_128_CBC = "aes-128-cbc"
	TYPE_192_CBC = "aes-192-cbc"
	TYPE_256_CBC = "aes-256-cbc"
	TYPE_128_ECB = "aes-128-ecb"
	TYPE_192_ECB = "aes-192-ecb"
	TYPE_256_ECB = "aes-256-ecb"
)

const ERROR_BAD_DECRYPT = "aesmodule : bad decrypt"

//Encrypter is a function that generate encrypted string from simple string + enryption key
// 'encType' is a parameter of some AES - modes
func Encrypter(key, data, encType string) (string, error) {
	c1 := exec.Command("echo", data)
	c2 := exec.Command("openssl", encType, "-a", "-salt", "-k", key)
	return pipingExec(c1, c2)
}

//Encrypter is a function that generate decrypted string from encrypted with enryption key
// 'encType' is a parameter of some AES - modes
func Decrypter(key, data, encType string) (string, error) {
	c1 := exec.Command("echo", data)
	c2 := exec.Command("openssl", encType, "-a", "-d", "-salt", "-k", key)
	return pipingExec(c1, c2)
}

//pipingExec is a function that make piping nice and easy + anti-duplication code
func pipingExec(c1, c2 *exec.Cmd) (string, error) {
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2o bytes.Buffer
	var b2e bytes.Buffer
	c2.Stdout = &b2o
	c2.Stderr = &b2e

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()

	if len(b2e.String()) != 0 {
		return "", errors.New(ERROR_BAD_DECRYPT)
	}

	return strings.Replace(b2o.String(), "\n", "", 1), nil
}
