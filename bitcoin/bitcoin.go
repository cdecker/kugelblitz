package bitcoin

type BitcoinDOpts struct {
	Datadir string
}

type BitcoinD interface {
	Start() error
	Stop() error
}

type bitcoinD struct {
	opts BitcoinDOpts
}

func NewBitcoinD(opts BitcoinDOpts) BitcoinD {
	return bitcoinD{
		opts: opts,
	}
}

func (b bitcoinD) Start() error {
	return nil
}

func (b bitcoinD) Stop() error {
	return nil
}
