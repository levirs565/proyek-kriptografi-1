package main

import (
	"slices"
	"testing"
)

func TestDecodeHex(t *testing.T) {
	in := "603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4"
	out := [32]uint8{0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
		0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4}
	out_actual, err := decodeHexString(in)
	if err != nil {
		t.Error(err)
		return
	}

	if !slices.Equal(out[:], out_actual) {
		t.Logf("%x\n", out_actual[:])
		t.Errorf("Not match")
	}
}

func TestEncodeHex(t *testing.T) {
	in := [32]uint8{0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
		0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4}
	out := "603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4"
	out_actual, err := encodeHexString(in[:])
	if err != nil {
		t.Error(err)
		return
	}

	if out != out_actual {
		t.Logf("%s", out_actual)
		t.Errorf("Not match")
	}
}

var encodeDecodeNameMap = []string{"Encode", "Decode"}

func TestBase64(t *testing.T) {
	list := []struct {
		input  string
		output string
	}{
		{
			input:  "Hello World!",
			output: "SGVsbG8gV29ybGQh",
		}, {
			input:  "hai halo",
			output: "aGFpIGhhbG8=",
		}, {
			input:  "hai halohh",
			output: "aGFpIGhhbG9oaA==",
		},
	}

	for _, test := range list {
		for i := range 2 {
			encode := i == 0
			t.Run(encodeDecodeNameMap[i], func(t *testing.T) {
				input := test.input
				output := test.output

				if !encode {
					input, output = output, input
				}

				var output_actual string
				if encode {
					output_actual = encodeBase64([]uint8(input))
				} else {
					o, err := decodeBase64(input)
					if err != nil {
						t.Error(err)
						return
					}
					output_actual = string(o)
				}

				if output_actual != output {
					t.Log(output)
					t.Log(output_actual)
					t.Errorf("Not match")
				}
			})
		}
	}
}
