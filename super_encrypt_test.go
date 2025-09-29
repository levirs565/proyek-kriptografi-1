package main

import (
	"kriptografi1/aes"
	"slices"
	"testing"
)

func TestSuperEncrypt(t *testing.T) {
	aes.AESInit()
	privateKey := "SPRIV-r4GoNFLgC3jUMHcpeulLpiT2HJbGcS0l+4TlfRqnOKHcAGsAJd66kZ20UrCHDRPCRU9coU2Hp1lXNxTtQ0MdP6B1Q9SbPE72d1PtX0uzDel4U8ywF5yvKP2+W/vOl3xDkgTyW16zdKs4mW7KWZiBuuwznOqeCuzjZTvwWGZupNYZto1bATmtCVvfAEMlBnYodRoJEwtarwdDJOXinpagX537SaX1MK8inLZkhHz5ssbwDIWbb6yOUlM+Mw6F8I+7zhUmrPfItjd1Pz68nC9M2mkp5ejkfXoWBDbmT2Im2FVUDEsvcy0+0kP/jDipNTlmqnTzFYqlDYwM+KgpC+XAbQ==-PEzK8vdKfQaNuP0oJk/Dqxjl6lqvYbFGbI1zn7EppR6PjSE2Usycw20wyHxYWljglZm31L/jjM74VTd+cW68vOCviZehom3q5oLSerxyj6QYsIoSMnql/+p2nfn9ODT8rk2+yz/VXo95idXz9iac5BFFhnNhy5UqZ9zZLjY7a6yF1IGf8ESY55eNehvet+hxHgjcb8pbojVReRMi75nlsLuTcXq+VzF93SNdnecbfIBASaMyS+qhXAdD4rGVSkjedYrT9RoP2NIbcRoy/tVKbIdtHCVF2/ftdNyiD2nsLFXvXd9dDUp+oRst9wBQu24o4zIVpy4/P+a52Ev6MF+UgQ=="
	publicKey := "SPUB-r4GoNFLgC3jUMHcpeulLpiT2HJbGcS0l+4TlfRqnOKHcAGsAJd66kZ20UrCHDRPCRU9coU2Hp1lXNxTtQ0MdP6B1Q9SbPE72d1PtX0uzDel4U8ywF5yvKP2+W/vOl3xDkgTyW16zdKs4mW7KWZiBuuwznOqeCuzjZTvwWGZupNYZto1bATmtCVvfAEMlBnYodRoJEwtarwdDJOXinpagX537SaX1MK8inLZkhHz5ssbwDIWbb6yOUlM+Mw6F8I+7zhUmrPfItjd1Pz68nC9M2mkp5ejkfXoWBDbmT2Im2FVUDEsvcy0+0kP/jDipNTlmqnTzFYqlDYwM+KgpC+XAbQ==-AQAB"
	message := "Halo dunia. Aku baru bagun. Tapi ingin tidur lagi. Gak ada semangat ngapain"

	privateKeyRSA, err := SuperDecodePrivateKey(privateKey)

	if err != nil {
		t.Error(err)
		return
	}

	publicKeyRSA, err := SuperDecodePublicKey(publicKey)

	if err != nil {
		t.Error(err)
		return
	}

	cipher, err := SuperEncrypt(publicKeyRSA, []uint8(message))

	if err != nil {
		t.Error(err)
		return
	}

	plain, err := SuperDecrypt(privateKeyRSA, cipher)

	if err != nil {
		t.Error(err)
		return
	}

	if !slices.Equal(plain, []uint8(message)) {
		t.Errorf("Not match")
		return
	}
}
