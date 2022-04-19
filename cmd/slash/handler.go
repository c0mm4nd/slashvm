// Package main
// an example of using go-slashvm to implement and exec evm calls
package main

import (
	"fmt"
	"time"

	"encoding/hex"

	"github.com/c0mm4nd/slashvm"
	"github.com/c0mm4nd/slashvm/configs"
	"github.com/c0mm4nd/slashvm/core"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
	"github.com/c0mm4nd/slashvm/gasometer"
	"github.com/c0mm4nd/slashvm/utils"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

// SimpleHandler implements a Handler for handling chain-related calls
type SimpleHandler struct {
	*configs.Config
	*gasometer.Gasometer
	Height     uint64
	Depth      uint64
	Balances   map[string]*uint256.Int
	Codes      map[string][]byte
	CodeHashes map[string][]byte
	Storages   map[string]map[string]*uint256.Int
}

func (cs *SimpleHandler) Balance(b *uint256.Int) *uint256.Int {
	return cs.Balances[utils.U256ToHashHex(b)]
}

func (cs *SimpleHandler) CodeSize(b *uint256.Int) *uint256.Int {
	return uint256.NewInt(uint64(len(cs.Codes[utils.U256ToHashHex(b)])))
}

func (cs *SimpleHandler) CodeHash(b *uint256.Int) *uint256.Int {
	return new(uint256.Int).SetBytes32(cs.CodeHashes[utils.U256ToHashHex(b)])
}

func (cs *SimpleHandler) ChainID() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) Origin() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) GasPrice() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) BlockBaseFeePerGas() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) BlockHash(number uint64) *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) BlockCoinbase() *uint256.Int {
	txHash, _ := hex.DecodeString("deadbeef")
	return new(uint256.Int).SetBytes32(txHash)
}

func (cs *SimpleHandler) BlockTimestamp() *uint256.Int {
	now := time.Now().Unix()
	return uint256.NewInt(uint64(now))
}

func (cs *SimpleHandler) BlockNumber() *uint256.Int {
	return uint256.NewInt(cs.Height)
}

func (cs *SimpleHandler) BlockDifficulty() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) GasLimit() *uint256.Int {
	return uint256.NewInt(0)
}

func (cs *SimpleHandler) SetStorage(address *uint256.Int, index *uint256.Int, value *uint256.Int) *exits.Error {
	_, exists := cs.Storages[utils.U256ToHashHex(address)]
	if !exists {
		cs.Storages[utils.U256ToHashHex(address)] = make(map[string]*uint256.Int)
	}
	cs.Storages[utils.U256ToHashHex(address)][utils.U256ToHashHex(index)] = value

	return nil
}

func (cs *SimpleHandler) Storage(address *uint256.Int, index *uint256.Int) *uint256.Int {
	return cs.Storages[utils.U256ToHashHex(address)][utils.U256ToHashHex(index)]
}

func (cs *SimpleHandler) GasLeft() *uint256.Int {
	return new(uint256.Int).SetAllOne()
}

func (cs *SimpleHandler) Create(address *uint256.Int, scheme *slashvm.CreateSchemeLegacy, value *uint256.Int, code []byte) (exits.Reason, *uint256.Int, []byte, *slashvm.CreateInterrupt) {
	return nil, nil, nil, nil //TODO
}

func (cs *SimpleHandler) Create2(address *uint256.Int, scheme *slashvm.CreateSchemeCreate2, value *uint256.Int, code []byte) (exits.Reason, *uint256.Int, []byte, *slashvm.CreateInterrupt) {
	return nil, nil, nil, nil //TODO
}

func (cs *SimpleHandler) Log(address *uint256.Int, topics []*uint256.Int, data []byte) *exits.Error {
	fmt.Println(address, topics, string(data))
	return nil
}

func (cs *SimpleHandler) MarkDelete(address *uint256.Int, target *uint256.Int) *exits.Error {
	return nil
}

func (cs *SimpleHandler) Code(address *uint256.Int) []byte {
	return cs.Codes[utils.U256ToHashHex(address)]
}

func (cs *SimpleHandler) Call(addr *uint256.Int, transfer *slashvm.Transfer, input []byte, gas uint64, context *slashvm.Context) (exits.Reason, []byte, *slashvm.CallInterrupt) {
	return cs.call_inner(addr, transfer, input, &gas, false, true, true, context)
}

func (cs *SimpleHandler) CallCode(addr *uint256.Int, transfer *slashvm.Transfer, input []byte, gas uint64, context *slashvm.Context) (exits.Reason, []byte, *slashvm.CallInterrupt) {
	return nil, nil, nil
}
func (cs *SimpleHandler) DelegateCall(addr *uint256.Int, input []byte, gas uint64, context *slashvm.Context) (exits.Reason, []byte, *slashvm.CallInterrupt) {
	return nil, nil, nil
}
func (cs *SimpleHandler) StaticCode(addr *uint256.Int, input []byte, gas uint64, context *slashvm.Context) (exits.Reason, []byte, *slashvm.CallInterrupt) {
	return cs.call_inner(addr, nil, input, &gas, true, true, true, context)
}
func (cs *SimpleHandler) CallFeedback() *exits.Error {
	return nil
}

func (cs *SimpleHandler) PreValidate(ctx *slashvm.Context, op opcodes.Opcode, stack *core.Stack) *exits.Error {
	return nil
}
func (cs *SimpleHandler) Other(op opcodes.Opcode, m *core.Machine) *exits.Error {
	panic(fmt.Sprintf("unknown op %d", op))
}

// transact_call
// func (cs *SimpleHandler) CallOnTx(caller string, address string, value *uint256.Int,
// 	data []byte, accessList map[string]*uint256.Int) {
// 	cs.Call()
// }

func (cs *SimpleHandler) call_inner(
	code_address *uint256.Int,
	transfer *slashvm.Transfer,
	input []byte,
	target *uint64,
	is_static bool,
	take_l64 bool,
	take_stipend bool,
	ctx *slashvm.Context) (exits.Reason, []byte, *slashvm.CallInterrupt) {

	// event!(Call {
	// 	code_address,
	// 	transfer: &transfer,
	// 	input: &input,
	// 	target_gas,
	// 	is_static,
	// 	context: &context,
	// });
	var l64 = func(gas uint64) uint64 {
		return gas - gas/64
	}

	var after_gas uint64
	if take_l64 && cs.Config.CallL64AfterGas {
		if cs.Config.Estimate {
			initial_after_gas := cs.Gasometer.Gas()
			diff := initial_after_gas - l64(initial_after_gas)
			cs.Gasometer.RecordCost(diff)
			after_gas = cs.Gasometer.Gas()
		} else {
			after_gas = l64(cs.Gasometer.Gas())
		}
	} else {
		after_gas = cs.Gasometer.Gas()
	}

	target_gas := after_gas
	if target != nil {
		target_gas = *target
	}

	gasLimit := uint256.NewInt(utils.Min([]uint64{after_gas, target_gas})) // generics looks silly
	err := cs.Gasometer.RecordCost(gasLimit.Uint64())
	if err != nil {
		return err, nil, nil
	}

	if transfer != nil {
		if take_stipend && !transfer.Value.IsZero() {
			gasLimit = utils.SaturatingAdd(gasLimit, uint256.NewInt(cs.Config.CallStipend))
		}
	}

	code := cs.Code(code_address)

	cs.enter_substate(gas_limit, is_static)
	cs.touch(context.address)

	if cs.Depth != nil {
		if *cs.Depth > cs.Config.CallStackLimit {
			cs.exit_substate(exits.Revert)
			return exits.CallTooDeep, nil, nil
		}
	}

	if transfer != nil {
		err := cs.Transfer(transfer)
		if err != nil {
			cs.exit_substate(exits.Revert)
			return err, nil, nil
		}
	}

	if result := cs.precompile_set.execute(code_address, input, gasLimit, context, is_static); result != nil {
		switch result.(type) {

		}
	}

	runtime := slashvm.NewRuntime(code, input, context, cs.Config)
	reason := cs.Execuate(runtime)
	switch reason.(type) {

	}
}

func (cs *SimpleHandler) create(address *uint256.Int,
	scheme *slashvm.CreateSchemeCreate2, value *uint256.Int, code []byte) {
	// create_address
	contractAddrBytes := sha3.Sum256([]byte{01, 02})
	contractAddr := contractAddrBytes[:]
	contractAddrStr := hex.EncodeToString(contractAddr)
	cs.Codes[contractAddrStr] = code
}
