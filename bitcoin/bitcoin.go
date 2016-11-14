package bitcoin

type BitcoinDOpts struct {
}

type BitcoinD interface {
	Start() error
	Stop() error
}

type bitcoinD struct {
	opts BitcoinDOpts
}

func NewBitcoinD(opts BitcoinDOpts) {
	return &bitcoinD{
		opts: opts,
	}
}

func (b *bitcoinD) Start() error {
}

func (b *bitcoinD) Stop() error {
}
