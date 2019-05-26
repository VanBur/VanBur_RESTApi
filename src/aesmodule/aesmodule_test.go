package aesmodule

import (
	"testing"
)

type teststruct struct {
	key        string
	inputData  string
	outputData string
}

var cbcDecTests = []teststruct{
	{"superpass", "U2FsdGVkX18PJILwscA+WPkF9jB+vtBMH4hjEVhQU1Wl+Zbi75xtwQuOhKVEuyEh", "wow_it_super_password"},
	{"qwerty", "U2FsdGVkX1/Tn5fh1t1MErXQxqw9ed4Pdei/QBF6ncc=", "too_easy_pass"},
}

var ecbDecTests = []teststruct{
	{"ZaXsCdVf", "U2FsdGVkX18hNXfWpBEaBlrolavOTUQ9zo1MSQm0JAU=", "powerman3000"},
	{"96-69", "U2FsdGVkX1/N28CyZRIJ5V5ZXMeXcXYk00/Rtqkq/cEh2G3twzFZYeIAH9iF4wb1", "master_calculator"},
}

func TestCBCDecrypter(t *testing.T) {
	for _, pair := range cbcDecTests {
		v := Decrypter(pair.key, pair.inputData, TYPE_128_CBC)
		if v != pair.outputData {
			t.Error(
				"For data", pair.inputData,
				"with key", pair.key,
				"expected", pair.outputData,
				"got", v,
			)
		}
	}
}

func TestECBDecrypter(t *testing.T) {
	for _, pair := range ecbDecTests {
		v := Decrypter(pair.key, pair.inputData, TYPE_128_ECB)
		if v != pair.outputData {
			t.Error(
				"For data", pair.inputData,
				"with key", pair.key,
				"expected", pair.outputData,
				"got", v,
			)
		}
	}
}
