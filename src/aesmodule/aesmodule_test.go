package aesmodule

import (
	"testing"
)

type goodteststruct struct {
	key           string
	encryptedData string
	decryptedData string
}

type failteststruct struct {
	key       string
	inputData string
	aesType   string
}

var cbcDecTests = []goodteststruct{
	{"superpass", "U2FsdGVkX18PJILwscA+WPkF9jB+vtBMH4hjEVhQU1Wl+Zbi75xtwQuOhKVEuyEh", "wow_it_super_password"},
	{"qwerty", "U2FsdGVkX1/Tn5fh1t1MErXQxqw9ed4Pdei/QBF6ncc=", "too_easy_pass"},
}

var ecbDecTests = []goodteststruct{
	{"ZaXsCdVf", "U2FsdGVkX18hNXfWpBEaBlrolavOTUQ9zo1MSQm0JAU=", "powerman3000"},
	{"96-69", "U2FsdGVkX1/N28CyZRIJ5V5ZXMeXcXYk00/Rtqkq/cEh2G3twzFZYeIAH9iF4wb1", "master_calculator"},
}

var failDecTests = []failteststruct{
	{"ZzZzZz", "U2FsdGVkX18hNXfWpBEaBlrolavOTUQ9zo1MSQm0JAU=", TYPE_128_CBC},
	{"qwwerty", "U2FsdGVkX1/Tn5fh1t1MErXQxqw9ed4Pdei/QBF6ncc=", TYPE_128_ECB},
}

func TestCBCDecrypter(t *testing.T) {
	for _, pair := range cbcDecTests {
		v, _ := Decrypter(pair.key, pair.encryptedData, TYPE_128_CBC)
		if v != pair.decryptedData {
			t.Error(
				"For data", pair.encryptedData,
				"with key", pair.key,
				"expected", pair.decryptedData,
				"got", v,
			)
		}
	}
}

func TestECBDecrypter(t *testing.T) {
	for _, pair := range ecbDecTests {
		v, _ := Decrypter(pair.key, pair.encryptedData, TYPE_128_ECB)
		if v != pair.decryptedData {
			t.Error(
				"For data", pair.encryptedData,
				"with key", pair.key,
				"expected", pair.decryptedData,
				"got", v,
			)
		}
	}
}

func TestBadDecrypter(t *testing.T) {
	for _, pair := range failDecTests {
		_, e := Decrypter(pair.key, pair.inputData, pair.aesType)
		if e == nil {
			t.Error(
				"For data", pair.inputData,
				"with key", pair.key,
				"expected error, got 'nil'",
			)
		}
	}
}

func TestECBEcnrypter(t *testing.T) {
	for _, pair := range ecbDecTests {
		enc, _ := Encrypter(pair.key, pair.decryptedData, TYPE_128_ECB)
		dec, _ := Decrypter(pair.key, enc, TYPE_128_ECB)
		if dec != pair.decryptedData {
			t.Error(
				"After enc + dec",
				"with key", pair.key,
				"expected", pair.encryptedData,
				"got", dec,
			)
		}
	}
}

func TestCBCEcnrypter(t *testing.T) {
	for _, pair := range cbcDecTests {
		enc, _ := Encrypter(pair.key, pair.decryptedData, TYPE_128_CBC)
		dec, _ := Decrypter(pair.key, enc, TYPE_128_CBC)
		if dec != pair.decryptedData {
			t.Error(
				"After enc + dec",
				"with key", pair.key,
				"expected", pair.encryptedData,
				"got", dec,
			)
		}
	}
}

/*
func TestCBCEcnrypter(t *testing.T) {
	for _, pair := range cbcDecTests {
		v, _ := Encrypter(pair.key, pair.decryptedData, TYPE_128_CBC)
		if v != pair.encryptedData {
			t.Error(
				"For data", pair.decryptedData,
				"with key", pair.key,
				"expected", pair.encryptedData,
				"got", v,
			)
		}
	}
}


*/
