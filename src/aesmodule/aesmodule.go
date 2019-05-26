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

//Encrypter is a function that generate encrypted string from simple string + enryption key
// 'encType' is a parameter of some AES - modes
func Encrypter(key, data, encType string) string {
	c1 := exec.Command("echo", data)
	c2 := exec.Command("openssl", encType, "-a", "-salt", "-k", key)
	return pipingExec(c1, c2)
}

//Encrypter is a function that generate decrypted string from encrypted with enryption key
// 'encType' is a parameter of some AES - modes
func Decrypter(key, data, encType string) string {
	c1 := exec.Command("echo", data)
	//c2 := exec.Command("openssl", "aes-128-cbc", "-a", "-d", "-salt", "-k", key)
	c2 := exec.Command("openssl", encType, "-a", "-d", "-salt", "-k", key)
	return pipingExec(c1, c2)
}

//pipingExec is a function that make piping nice and easy + anti-duplication code
func pipingExec(c1, c2 *exec.Cmd) string {
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()

	return strings.Replace(b2.String(), "\n", "", 1)
}
