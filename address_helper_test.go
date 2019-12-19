package cnlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func addressHelperTestHelpers() *AddressHelper {
	bc := NewBaseCoin(84, 0, 0)
	ah := NewAddressHelper(bc)
	return ah
}

func TestBase58CheckEncoding_ValidAddress_ReturnsTrue(t *testing.T) {
	addresses := []string{
		"12vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
		"3EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
	}

	for _, addr := range addresses {
		err := AddressIsBase58CheckEncoded(addr)
		assert.Nil(t, err)
	}

}

func TestBase58CheckEncoding_InvalidAddresses_ReturnFalse(t *testing.T) {
	addresses := []string{
		"12vRFewBpbdiS5HXDDLEfVFtJnpA2",
		"12vRFewBpbdiS5HXDDLEfVFt",
		"diS5HXDDLEfVFtJnpA2x8NV8",
		"212vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"42vRFewBpbdiS5HXDDLEfVFtJnpA2x8NV8",
		"3EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk",
		"j6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"ZBaYXhVCa8ZrwpHJYtk44bGX",
		"23EH9Wj6KWaZBaYXhVCa8ZrwpHJYtk44bGX",
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz3",
		"3Cd4xEu2VvM352BVgd9cb1Ct",
		"Eu2VvM352BVgd9cb1Ct5vxz318tVT",
		"M352BVgd9cb1Ct5vxz318tVT",
		"23Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
		"4Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
		"0xF26C29D25a1E1696c5CC54DE4bf2AEc906EB4F79",
		"qr45rul6luexjgg5h8p26c0cs6rrhwzrkg6e0hdvrf",
		"Jenny86753098675309IgotIt",
		"31415926535ILikePi89793238462643",
		"foo",
		"",
		"com.coinninja.CoinKeeper.beta://google/link/",
	}

	for _, addr := range addresses {
		err := AddressIsBase58CheckEncoded(addr)
		assert.NotNil(t, err)
	}
}

func TestSegwitAddress_ValidAddresses_ReturnTrue(t *testing.T) {
	addresses := []string{
		"bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",                     // demo wallet first receive address
		"bc1q8c6fshw2dlwun7ekn9qwf37cu2rn755upcp6el",                     // demo wallet first change address
		"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",                     // p2wpkh sipa demo
		"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3", // p2wsh sipa demo
	}

	for _, addr := range addresses {
		err := AddressIsValidSegwitAddress(addr)
		assert.Nil(t, err)
	}
}

func TestSegwitAddress_InvalidAddresses_ReturnFalse(t *testing.T) {
	addresses := []string{
		"BC1QW508D6QEJXTDG4Y5R3ZARVAYR0C5XW7KV8F3T4", // p2wsh sipa demo invalid p2wpkh, YR transposed
		"3Cd4xEu2VvM352BVgd9cb1Ct5vxz318tVT",
		"com.coinninja.CoinKeeper.beta://google/link/",
		"",
	}

	for _, addr := range addresses {
		err := AddressIsValidSegwitAddress(addr)
		assert.NotNil(t, err)
	}
}

func TestSegwitAddressHRP(t *testing.T) {
	bcAddr := "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu"
	rtAddr := "bcrt1q6rz28mcfaxtmd6v789l9rrlrusdprr9pz3cppk"
	legacyAddr := "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf"
	helper := addressHelperTestHelpers()

	bcHrp, err := helper.HRPFromAddress(bcAddr)
	assert.Nil(t, err)
	assert.Equal(t, "bc", bcHrp)

	rtHrp, err := helper.HRPFromAddress(rtAddr)
	assert.Nil(t, err)
	assert.Equal(t, "bcrt", rtHrp)

	laHrp, err := helper.HRPFromAddress(legacyAddr)
	assert.NotNil(t, err)
	assert.Equal(t, "", laHrp)
}

func TestBytesPerInputBIP84Input(t *testing.T) {
	path := NewDerivationPath(84, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, 1, path, nil, true)
	bpi := addressHelperTestHelpers().bytesPerInput(utxo)
	assert.Equal(t, p2wpkhSegwitInputSize, bpi)
}

func TestBytesPerInputBIP49Input(t *testing.T) {
	bc := NewBaseCoin(49, 0, 0)
	ah := NewAddressHelper(bc)
	path := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, 1, path, nil, true)
	bpi := ah.bytesPerInput(utxo)
	assert.Equal(t, p2shSegwitInputSize, bpi)
}

func TestBytesPerChangeOuptutBIP84(t *testing.T) {
	bpco := addressHelperTestHelpers().bytesPerChangeOuptut()
	assert.Equal(t, p2wpkhOutputSize, bpco)
}

func TestBytesPerChangeOuptutBIP49(t *testing.T) {
	bc := NewBaseCoin(49, 0, 0)
	ah := NewAddressHelper(bc)
	bpco := ah.bytesPerChangeOuptut()
	assert.Equal(t, p2shOutputSize, bpco)
}

func TestTotalBytes_SingleBIP49Input_TwoBIP49Outputs(t *testing.T) {
	address := "35dKN7xvHH3xnBWUrWzJtkjfrAFXk6hyH8"
	bc := NewBaseCoin(49, 0, 0)
	expectedBytes := 166
	ah := NewAddressHelper(bc)
	path := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, 1, path, nil, true)
	utxos := []*UTXO{utxo}

	bytes, err := ah.totalBytes(utxos, address, true)
	assert.Nil(t, err)
	assert.Equal(t, expectedBytes, bytes)
}

func TestTotalBytes_SingleBIP84Input_TwoBIP84Outputs(t *testing.T) {
	address := "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
	expectedBytes := 141
	path := NewDerivationPath(84, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, 1, path, nil, true)
	utxos := []*UTXO{utxo}

	bytes, err := addressHelperTestHelpers().totalBytes(utxos, address, true)
	assert.Nil(t, err)
	assert.Equal(t, expectedBytes, bytes)
}

func TestTotalBytes_SingleBIP49Input_LegacyOutput_BIP49Change(t *testing.T) {
	address := "1LqBGSKuX5yYUonjxT5qGfpUsXKYYWeabA"
	bc := NewBaseCoin(49, 0, 0)
	expectedBytes := 168
	ah := NewAddressHelper(bc)
	path := NewDerivationPath(49, 0, 0, 0, 0)
	utxo := NewUTXO("previous txid", 0, 1, path, nil, true)
	utxos := []*UTXO{utxo}

	bytes, err := ah.totalBytes(utxos, address, true)
	assert.Nil(t, err)
	assert.Equal(t, expectedBytes, bytes)
}
