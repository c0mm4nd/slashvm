package configs

type Config struct {
	GasExtCode                uint64
	GasExtCodeHash            uint64
	GasSStoreSet              uint64
	GasSStoreReset            uint64
	RefundSStoreClears        uint64
	MaxRefundQuotient         uint64 /// EIP-3529
	GasBalance                uint64
	GasSLoad                  uint64
	GasSLoadCold              uint64
	GasSuicide                uint64
	GasSuicideNewAccount      uint64
	GasCall                   uint64
	GasExpByte                uint64
	GasTransactionCreate      uint64
	GasTransactionCall        uint64
	GasTransactionZeroData    uint64
	GasTransactionNonZeroData uint64
	GasAccessListAddress      uint64
	GasAccessListStorageKey   uint64
	GasAccountAccessCold      uint64
	GasStorageReadWarm        uint64

	SStoreGasMetering        bool
	SStoreRevertUnderStipend bool
	IncreaseStateAccessGas   bool
	DecreaseClearsRefund     bool
	DisallowExecutableFormat bool
	ErrOnCallWithMoreGas     bool
	CallL64AfterGas          bool
	EmptyConsideredExists    bool
	CreateIncreaseNonce      bool

	StackLimit     uint
	MemoryLimit    uint
	CallStackLimit uint

	CreateContractLimit int
	CallStipend         uint64
	HasDelegateCall     bool
	HasCreate2          bool
	HasRevert           bool
	HasReturnData       bool
	HasBitwiseShifting  bool
	HasChainID          bool
	HasSelfBalance      bool
	HasExtCodeHash      bool
	HasBaseFee          bool
	Estimate            bool
}

var FrontierConfig = &Config{
	GasExtCode:                20,
	GasExtCodeHash:            20,
	GasBalance:                20,
	GasSLoad:                  50,
	GasSLoadCold:              0,
	GasSStoreSet:              20000,
	GasSStoreReset:            5000,
	RefundSStoreClears:        15000,
	MaxRefundQuotient:         2,
	GasSuicide:                0,
	GasSuicideNewAccount:      0,
	GasCall:                   40,
	GasExpByte:                10,
	GasTransactionCreate:      21000,
	GasTransactionCall:        21000,
	GasTransactionZeroData:    4,
	GasTransactionNonZeroData: 68,
	GasAccessListAddress:      0,
	GasAccessListStorageKey:   0,
	GasAccountAccessCold:      0,
	GasStorageReadWarm:        0,
	SStoreGasMetering:         false,
	SStoreRevertUnderStipend:  false,
	IncreaseStateAccessGas:    false,
	DecreaseClearsRefund:      false,
	DisallowExecutableFormat:  false,
	ErrOnCallWithMoreGas:      true,
	EmptyConsideredExists:     true,
	CreateIncreaseNonce:       false,
	CallL64AfterGas:           false,
	StackLimit:                1024,
	MemoryLimit:               ^uint(0),
	CallStackLimit:            1024,
	CreateContractLimit:       -1, // None
	CallStipend:               2300,
	HasDelegateCall:           false,
	HasCreate2:                false,
	HasRevert:                 false,
	HasReturnData:             false,
	HasBitwiseShifting:        false,
	HasChainID:                false,
	HasSelfBalance:            false,
	HasExtCodeHash:            false,
	HasBaseFee:                false,
	Estimate:                  false,
}

var IstanbulConfig = &Config{
	GasExtCode:                700,
	GasExtCodeHash:            700,
	GasBalance:                700,
	GasSLoad:                  800,
	GasSLoadCold:              0,
	GasSStoreSet:              20000,
	GasSStoreReset:            5000,
	RefundSStoreClears:        15000,
	MaxRefundQuotient:         2,
	GasSuicide:                5000,
	GasSuicideNewAccount:      25000,
	GasCall:                   700,
	GasExpByte:                50,
	GasTransactionCreate:      53000,
	GasTransactionCall:        21000,
	GasTransactionZeroData:    4,
	GasTransactionNonZeroData: 16,
	GasAccessListAddress:      0,
	GasAccessListStorageKey:   0,
	GasAccountAccessCold:      0,
	GasStorageReadWarm:        0,
	SStoreGasMetering:         true,
	SStoreRevertUnderStipend:  true,
	IncreaseStateAccessGas:    false,
	DecreaseClearsRefund:      false,
	DisallowExecutableFormat:  false,
	ErrOnCallWithMoreGas:      false,
	EmptyConsideredExists:     false,
	CreateIncreaseNonce:       true,
	CallL64AfterGas:           true,
	StackLimit:                1024,
	MemoryLimit:               ^uint(0),
	CallStackLimit:            1024,
	CreateContractLimit:       0x6000,
	CallStipend:               2300,
	HasDelegateCall:           true,
	HasCreate2:                true,
	HasRevert:                 true,
	HasReturnData:             true,
	HasBitwiseShifting:        true,
	HasChainID:                true,
	HasSelfBalance:            true,
	HasExtCodeHash:            true,
	HasBaseFee:                false,
	Estimate:                  false,
}

type DerivedConfigInputs struct {
	GasStorageReadWarm       uint64
	GasSLoadCold             uint64
	GasAccessListStorageKey  uint64
	DecreaseClearsRefund     bool
	HasBaseFee               bool
	DisallowExecutableFormat bool
}

func NewConfigWithDerivedValues(inputs DerivedConfigInputs) *Config {
	var (
		gasSStoreReset     = 5000 - inputs.GasSLoadCold
		refundSStoreClears = uint64(15000)
		maxRefundQuotient  = uint64(2)
	)
	if inputs.DecreaseClearsRefund {
		refundSStoreClears = uint64(gasSStoreReset + inputs.GasAccessListStorageKey)
		maxRefundQuotient = uint64(5)
	}

	return &Config{
		GasExtCode:                0,
		GasExtCodeHash:            0,
		GasBalance:                0,
		GasSLoad:                  inputs.GasStorageReadWarm,
		GasSLoadCold:              inputs.GasSLoadCold,
		GasSStoreSet:              20000,
		GasSStoreReset:            gasSStoreReset,
		RefundSStoreClears:        refundSStoreClears,
		MaxRefundQuotient:         maxRefundQuotient,
		GasSuicide:                5000,
		GasSuicideNewAccount:      25000,
		GasCall:                   0,
		GasExpByte:                50,
		GasTransactionCreate:      53000,
		GasTransactionCall:        21000,
		GasTransactionZeroData:    4,
		GasTransactionNonZeroData: 16,
		GasAccessListAddress:      2400,
		GasAccessListStorageKey:   inputs.GasAccessListStorageKey,
		GasAccountAccessCold:      2600,
		GasStorageReadWarm:        inputs.GasStorageReadWarm,
		SStoreGasMetering:         true,
		SStoreRevertUnderStipend:  true,
		IncreaseStateAccessGas:    true,
		DecreaseClearsRefund:      inputs.DecreaseClearsRefund,
		DisallowExecutableFormat:  inputs.DisallowExecutableFormat,
		ErrOnCallWithMoreGas:      false,
		EmptyConsideredExists:     false,
		CreateIncreaseNonce:       true,
		CallL64AfterGas:           true,
		StackLimit:                1024,
		MemoryLimit:               ^uint(0),
		CallStackLimit:            1024,
		CreateContractLimit:       0x6000,
		CallStipend:               2300,
		HasDelegateCall:           true,
		HasCreate2:                true,
		HasRevert:                 true,
		HasReturnData:             true,
		HasBitwiseShifting:        true,
		HasChainID:                true,
		HasSelfBalance:            true,
		HasExtCodeHash:            true,
		HasBaseFee:                inputs.HasBaseFee,
		Estimate:                  false,
	}
}
