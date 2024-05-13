package wallet

func init() {
	coins[BTC] = newBTC
	coins[BTCTestnet] = newBTC
}

type btc struct {
	name   string
	symbol string
	key    *Key
}

func newBTC(key *Key) Wallet {
	return &btc{
		name:   "Bitcoin",
		symbol: "BTC",
		key:    key,
	}
}

func (c *btc) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *btc) GetPath() []uint32 {
	return c.key.Opt.GetPath()
}

func (c *btc) GetName() string {
	return c.name
}

func (c *btc) GetSymbol() string {
	return c.symbol
}

func (c *btc) GetKey() *Key {
	return c.key
}

func (c *btc) GetAddress() (string, error) {
	return c.key.AddressBTCLegacy()
}

func (c *btc) GetPrivateKey() (string, error) {
	return c.key.PrivateWIF(true)
}
