package constants

// Network represents a bitcoin-like coin, housing all the special magic
// numbers which are used for encoding addresses, and public and private keys.
type Network struct {
	// The ticker symbol and name of this coin's network.
	Symbol string
	Name   string

	// Only availble for segwit-enabled coins.
	Bech32 string

	// Constants used for base58-check encoding of addresses.
	ScriptHash uint16
	PubkeyHash uint16

	// Used for base58-check encoding of WIF private keys.
	WIF byte

	// Used for base58-check encoding of extended public and private keys.
	ExtendedPublic  uint32
	ExtendedPrivate uint32
}

var (
	BitcoinNetwork = Network{
		Symbol:          "BTC",
		Name:            "Bitcoin",
		Bech32:          "bc",
		ScriptHash:      5,
		PubkeyHash:      0,
		WIF:             128,
		ExtendedPublic:  76067358,
		ExtendedPrivate: 76066276,
	}

	BitcoinTestnet = Network{
		Symbol:          "tBTC",
		Name:            "Testnet Bitcoin",
		Bech32:          "tb",
		ScriptHash:      196,
		PubkeyHash:      111,
		WIF:             239,
		ExtendedPublic:  70617039,
		ExtendedPrivate: 70615956,
	}

	LitecoinNetwork = Network{
		Symbol:          "LTC",
		Name:            "Litecoin",
		Bech32:          "ltc",
		ScriptHash:      50,
		PubkeyHash:      48,
		WIF:             176,
		ExtendedPublic:  27108450,
		ExtendedPrivate: 27106558,
	}

	ZcashNetwork = Network{
		Symbol:          "ZEC",
		Name:            "Zcash",
		ScriptHash:      7357,
		PubkeyHash:      7352,
		WIF:             128,
		ExtendedPublic:  76067358,
		ExtendedPrivate: 76066276,
	}
)

// NetworksByName exposes network data sorted by capitalized coin names.
var NetworksByName = map[string]Network{
	"Bitcoin":  BitcoinNetwork,
	"Litecoin": LitecoinNetwork,
	"Zcash":    ZcashNetwork,
}

// NetworksBySymbol exposes network data sorted by coin ticker symbol.
var NetworksBySymbol = map[string]Network{
	"BTC": BitcoinNetwork,
	"LTC": LitecoinNetwork,
	"ZEC": ZcashNetwork,
}

// CurrentNetwork is a global variable that defines which network callers should use.
var CurrentNetwork = BitcoinNetwork
