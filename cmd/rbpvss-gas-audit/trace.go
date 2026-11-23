package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type opcodeEntry struct {
	Count uint64 `json:"count"`
	Gas   uint64 `json:"charged_gas"`
}
type precompileEntry struct {
	Calls        uint64 `json:"calls"`
	InputBytes   uint64 `json:"input_bytes"`
	Gas          uint64 `json:"traced_gas"`
	EIPGas       uint64 `json:"eip_predicted_gas"`
	PairingTerms uint64 `json:"pairing_terms,omitempty"`
}
type slotKey struct {
	Address common.Address
	Slot    common.Hash
}
type callFrame struct {
	Name        string
	InputBytes  uint64
	ExpectedGas uint64
	Terms       uint64
}

type gasAudit struct {
	N                      int                         `json:"n"`
	L                      int                         `json:"l"`
	T                      int                         `json:"t"`
	Run                    int                         `json:"run"`
	AcceptedCount          int                         `json:"accepted_count"`
	ReceiptGas             uint64                      `json:"receipt_gas_used"`
	PredictedGas           uint64                      `json:"eip_opcode_reconstructed_gas"`
	Difference             int64                       `json:"difference"`
	TxHash                 string                      `json:"transaction_hash"`
	CalldataZero           uint64                      `json:"calldata_zero_bytes"`
	CalldataNonzero        uint64                      `json:"calldata_nonzero_bytes"`
	MemoryBytes            uint64                      `json:"peak_root_memory_bytes"`
	Components             map[string]uint64           `json:"components_gas"`
	Precompiles            map[string]*precompileEntry `json:"precompiles"`
	Opcodes                map[string]*opcodeEntry     `json:"opcodes"`
	SloadExpected          uint64                      `json:"eip2929_predicted_sload"`
	SstoreExpected         uint64                      `json:"eip2200_2929_predicted_sstore"`
	SstoreZeroToNonzero    uint64                      `json:"sstore_zero_to_nonzero"`
	SstoreNonzeroToNonzero uint64                      `json:"sstore_nonzero_to_nonzero"`
	SstoreNoop             uint64                      `json:"sstore_noop"`
	ExpectedRefund         int64                       `json:"eip3529_refund_counter"`
	ExecutionGross         uint64                      `json:"opcode_precompile_execution_gas_before_refund"`
	CallOverhead           uint64                      `json:"staticcall_overhead_including_call_memory"`
	vmContext              *tracing.VMContext
	seen                   map[slotKey]bool
	original               map[slotKey]common.Hash
	frames                 map[int]callFrame
	callCosts              map[int]uint64
	otherOpcode            uint64
	sload                  uint64
	sstore                 uint64
	problem                error
}

func newGasAudit(n, l, t, run, m int, data []byte) *gasAudit {
	a := &gasAudit{N: n, L: l, T: t, Run: run, AcceptedCount: m, Components: map[string]uint64{}, Precompiles: map[string]*precompileEntry{}, Opcodes: map[string]*opcodeEntry{}, seen: map[slotKey]bool{}, original: map[slotKey]common.Hash{}, frames: map[int]callFrame{}, callCosts: map[int]uint64{}}
	for _, b := range data {
		if b == 0 {
			a.CalldataZero++
		} else {
			a.CalldataNonzero++
		}
	}
	return a
}

func (a *gasAudit) fail(format string, args ...interface{}) {
	if a.problem == nil {
		a.problem = fmt.Errorf(format, args...)
	}
}

func (a *gasAudit) hooks() *tracing.Hooks {
	return &tracing.Hooks{
		OnTxStart: func(context *tracing.VMContext, _ *types.Transaction, _ common.Address) { a.vmContext = context },
		OnOpcode: func(_ uint64, op byte, _ uint64, cost uint64, scope tracing.OpContext, _ []byte, depth int, err error) {
			if err != nil {
				a.fail("opcode error: %v", err)
				return
			}
			name := vm.OpCode(op).String()
			entry := a.Opcodes[name]
			if entry == nil {
				entry = &opcodeEntry{}
				a.Opcodes[name] = entry
			}
			entry.Count++
			if depth == 1 && uint64(len(scope.MemoryData())) > a.MemoryBytes {
				a.MemoryBytes = uint64(len(scope.MemoryData()))
			}
			if op == byte(vm.STATICCALL) {
				a.callCosts[depth] = cost
				return
			}
			if op == byte(vm.CALL) || op == byte(vm.DELEGATECALL) || op == byte(vm.CREATE) || op == byte(vm.CREATE2) {
				a.fail("unexpected non-precompile call opcode %s", name)
				return
			}
			entry.Gas += cost
			if op != byte(vm.SLOAD) && op != byte(vm.SSTORE) {
				a.otherOpcode += cost
				return
			}
			stack := scope.StackData()
			slot := slotKey{scope.Address(), common.Hash(stack[len(stack)-1].Bytes32())}
			current := a.vmContext.StateDB.GetState(slot.Address, slot.Slot)
			if _, ok := a.original[slot]; !ok {
				a.original[slot] = current
			}
			cold := !a.seen[slot]
			a.seen[slot] = true
			if op == byte(vm.SLOAD) {
				expected := uint64(100)
				if cold {
					expected = 2100
				}
				a.SloadExpected += expected
				a.sload += cost
				if cost != expected {
					a.fail("SLOAD charged %d, EIP-2929 predicts %d", cost, expected)
				}
				return
			}
			value := common.Hash(stack[len(stack)-2].Bytes32())
			original := a.original[slot]
			expected := uint64(100)
			if value == current {
				a.SstoreNoop++
			} else if original == current {
				if original == (common.Hash{}) {
					expected = 20000
					a.SstoreZeroToNonzero++
				} else {
					expected = 2900
					if value == (common.Hash{}) {
						a.ExpectedRefund += 4800
					} else {
						a.SstoreNonzeroToNonzero++
					}
				}
			} else {
				if original != (common.Hash{}) {
					if current == (common.Hash{}) {
						a.ExpectedRefund -= 4800
					}
					if value == (common.Hash{}) {
						a.ExpectedRefund += 4800
					}
				}
				if original == value {
					if original == (common.Hash{}) {
						a.ExpectedRefund += 19900
					} else {
						a.ExpectedRefund += 2800
					}
				}
			}
			if cold {
				expected += 2100
			}
			a.SstoreExpected += expected
			a.sstore += cost
			if cost != expected {
				a.fail("SSTORE charged %d, EIP-2200/2929 predicts %d", cost, expected)
			}
		},
		OnEnter: func(depth int, _ byte, _ common.Address, to common.Address, input []byte, gas uint64, _ *big.Int) {
			if depth == 0 {
				return
			}
			cost, ok := a.callCosts[depth]
			if !ok || cost < gas {
				a.fail("missing parent STATICCALL cost at depth %d", depth)
				return
			}
			a.CallOverhead += cost - gas
			name, expected, terms, err := precompilePrice(to, input)
			if err != nil {
				a.fail("%v", err)
				return
			}
			a.frames[depth] = callFrame{name, uint64(len(input)), expected, terms}
		},
		OnExit: func(depth int, _ []byte, gasUsed uint64, err error, reverted bool) {
			if depth == 0 {
				return
			}
			frame, ok := a.frames[depth]
			if !ok {
				a.fail("missing precompile frame")
				return
			}
			delete(a.frames, depth)
			if err != nil || reverted {
				a.fail("precompile %s failed: %v", frame.Name, err)
			}
			entry := a.Precompiles[frame.Name]
			if entry == nil {
				entry = &precompileEntry{}
				a.Precompiles[frame.Name] = entry
			}
			entry.Calls++
			entry.InputBytes += frame.InputBytes
			entry.Gas += gasUsed
			entry.EIPGas += frame.ExpectedGas
			entry.PairingTerms += frame.Terms
			if gasUsed != frame.ExpectedGas {
				a.fail("precompile %s charged %d, EIP predicts %d", frame.Name, gasUsed, frame.ExpectedGas)
			}
		},
	}
}

// Independently encoded EIP prices; this does not call the EVM's RequiredGas.
func precompilePrice(address common.Address, input []byte) (string, uint64, uint64, error) {
	id := new(big.Int).SetBytes(address[:])
	if !id.IsUint64() {
		return "", 0, 0, fmt.Errorf("unexpected call address")
	}
	switch id.Uint64() {
	case 2:
		return "sha256", 60 + 12*((uint64(len(input))+31)/32), 0, nil
	case 5:
		if len(input) != 192 {
			return "", 0, 0, fmt.Errorf("unexpected modexp input length")
		}
		for i := 0; i < 3; i++ {
			if new(big.Int).SetBytes(input[i*32:(i+1)*32]).Uint64() != 32 {
				return "", 0, 0, fmt.Errorf("unexpected modexp word length")
			}
		}
		iterations := new(big.Int).SetBytes(input[128:160]).BitLen() - 1
		if iterations < 1 {
			iterations = 1
		}
		gas := uint64(16 * iterations / 3)
		if gas < 200 {
			gas = 200
		}
		return "modexp", gas, 0, nil
	case 6:
		return "ecadd", 150, 0, nil
	case 7:
		return "ecmul", 6000, 0, nil
	case 8:
		if len(input)%192 != 0 {
			return "", 0, 0, fmt.Errorf("invalid pairing input length")
		}
		terms := uint64(len(input) / 192)
		return "pairing", 45000 + 34000*terms, terms, nil
	default:
		return "", 0, 0, fmt.Errorf("unexpected precompile address %s", address)
	}
}

func (a *gasAudit) finish(receipt *types.Receipt) error {
	if a.problem != nil {
		return a.problem
	}
	expectedCalls := map[string]uint64{"ecmul": uint64(5*a.N + 9*a.L), "ecadd": uint64(4*a.N + 5*a.L), "pairing": uint64(a.L), "modexp": uint64(a.N + a.L), "sha256": uint64(2 + a.L + a.N + a.L - a.T)}
	for name, count := range expectedCalls {
		if a.Precompiles[name] == nil || a.Precompiles[name].Calls != count {
			return fmt.Errorf("%s operation count mismatch: %+v, want %d", name, a.Precompiles[name], count)
		}
	}
	if a.Precompiles["pairing"].PairingTerms != uint64(4*a.L) {
		return fmt.Errorf("pairing term count mismatch")
	}
	var precompileGas, predictedPrecompile uint64
	for name, entry := range a.Precompiles {
		a.Components[name] = entry.Gas
		precompileGas += entry.Gas
		predictedPrecompile += entry.EIPGas
	}
	a.Components["transaction_base"] = 21000
	a.Components["calldata"] = 4*a.CalldataZero + 16*a.CalldataNonzero
	a.Components["sload"] = a.sload
	a.Components["sstore"] = a.sstore
	words := (a.MemoryBytes + 31) / 32
	memoryGas := 3*words + words*words/512
	if memoryGas > a.otherOpcode+a.CallOverhead {
		return fmt.Errorf("memory expansion exceeds non-storage opcode gas")
	}
	a.Components["memory_expansion"] = memoryGas
	a.Components["other_execution"] = a.otherOpcode + a.CallOverhead - memoryGas
	a.ExecutionGross = precompileGas + a.otherOpcode + a.CallOverhead + a.sload + a.sstore
	gross := 21000 + a.Components["calldata"] + a.ExecutionGross
	if a.ExpectedRefund < 0 {
		return fmt.Errorf("negative final refund")
	}
	refund := uint64(a.ExpectedRefund)
	if cap := gross / 5; refund > cap {
		refund = cap
	}
	a.Components["refund_subtracted"] = refund
	a.ReceiptGas = receipt.GasUsed
	a.TxHash = receipt.TxHash.Hex()
	a.PredictedGas = 21000 + a.Components["calldata"] + predictedPrecompile + a.otherOpcode + a.CallOverhead + a.SloadExpected + a.SstoreExpected - refund
	a.Difference = int64(a.PredictedGas) - int64(a.ReceiptGas)
	if a.Difference != 0 {
		return fmt.Errorf("gas reconstruction differs: modeled %d, receipt %d, delta %d", a.PredictedGas, a.ReceiptGas, a.Difference)
	}
	return nil
}
