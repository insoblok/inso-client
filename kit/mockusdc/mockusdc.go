// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockusdc

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MockusdcMetaData contains all meta data concerning the Mockusdc contract.
var MockusdcMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"Debug\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040526040518060400160405280600d81526020017f4d6f636b2055534420436f696e000000000000000000000000000000000000008152505f908161004791906104ec565b506040518060400160405280600481526020017f55534443000000000000000000000000000000000000000000000000000000008152506001908161008c91906104ec565b50600660025f6101000a81548160ff021916908360ff1602179055503480156100b3575f5ffd5b507f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a6040516100e190610615565b60405180910390a15f64e8d4a5100090507f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a60405161011f9061067d565b60405180910390a18060045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055507f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a604051610196906106e5565b60405180910390a1806003819055507f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a6040516101d29061074d565b60405180910390a13373ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610237919061077a565b60405180910390a37f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a60405161026c906107dd565b60405180910390a17f7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a6040516102a190610845565b60405180910390a150610863565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061032a57607f821691505b60208210810361033d5761033c6102e6565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261039f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610364565b6103a98683610364565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6103ed6103e86103e3846103c1565b6103ca565b6103c1565b9050919050565b5f819050919050565b610406836103d3565b61041a610412826103f4565b848454610370565b825550505050565b5f5f905090565b610431610422565b61043c8184846103fd565b505050565b5b8181101561045f576104545f82610429565b600181019050610442565b5050565b601f8211156104a45761047581610343565b61047e84610355565b8101602085101561048d578190505b6104a161049985610355565b830182610441565b50505b505050565b5f82821c905092915050565b5f6104c45f19846008026104a9565b1980831691505092915050565b5f6104dc83836104b5565b9150826002028217905092915050565b6104f5826102af565b67ffffffffffffffff81111561050e5761050d6102b9565b5b6105188254610313565b610523828285610463565b5f60209050601f831160018114610554575f8415610542578287015190505b61054c85826104d1565b8655506105b3565b601f19841661056286610343565b5f5b8281101561058957848901518255600182019150602085019450602081019050610564565b868310156105a657848901516105a2601f8916826104b5565b8355505b6001600288020188555050505b505050505050565b5f82825260208201905092915050565b7f537461727420636f6e7374727563746f720000000000000000000000000000005f82015250565b5f6105ff6011836105bb565b915061060a826105cb565b602082019050919050565b5f6020820190508181035f83015261062c816105f3565b9050919050565b7f496e697469616c20737570706c792063616c63756c61746564000000000000005f82015250565b5f6106676019836105bb565b915061067282610633565b602082019050919050565b5f6020820190508181035f8301526106948161065b565b9050919050565b7f42616c616e63652061737369676e656420746f206d73672e73656e64657200005f82015250565b5f6106cf601e836105bb565b91506106da8261069b565b602082019050919050565b5f6020820190508181035f8301526106fc816106c3565b9050919050565b7f546f74616c20737570706c7920736574000000000000000000000000000000005f82015250565b5f6107376010836105bb565b915061074282610703565b602082019050919050565b5f6020820190508181035f8301526107648161072b565b9050919050565b610774816103c1565b82525050565b5f60208201905061078d5f83018461076b565b92915050565b7f5472616e73666572206576656e7420656d6974746564000000000000000000005f82015250565b5f6107c76016836105bb565b91506107d282610793565b602082019050919050565b5f6020820190508181035f8301526107f4816107bb565b9050919050565b7f436f6e7374727563746f7220646f6e65000000000000000000000000000000005f82015250565b5f61082f6010836105bb565b915061083a826107fb565b602082019050919050565b5f6020820190508181035f83015261085c81610823565b9050919050565b6101f9806108705f395ff3fe608060405234801561000f575f5ffd5b5060043610610029575f3560e01c806306fdde031461002d575b5f5ffd5b61003561004b565b6040516100429190610146565b60405180910390f35b5f805461005790610193565b80601f016020809104026020016040519081016040528092919081815260200182805461008390610193565b80156100ce5780601f106100a5576101008083540402835291602001916100ce565b820191905f5260205f20905b8154815290600101906020018083116100b157829003601f168201915b505050505081565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610118826100d6565b61012281856100e0565b93506101328185602086016100f0565b61013b816100fe565b840191505092915050565b5f6020820190508181035f83015261015e818461010e565b905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806101aa57607f821691505b6020821081036101bd576101bc610166565b5b5091905056fea2646970667358221220c730ffc979af4b344ba94bc3f0b3a6b38aa16076aaa4b7108b074ea2013d1d4864736f6c634300081d0033",
}

// MockusdcABI is the input ABI used to generate the binding from.
// Deprecated: Use MockusdcMetaData.ABI instead.
var MockusdcABI = MockusdcMetaData.ABI

// MockusdcBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockusdcMetaData.Bin instead.
var MockusdcBin = MockusdcMetaData.Bin

// DeployMockusdc deploys a new Ethereum contract, binding an instance of Mockusdc to it.
func DeployMockusdc(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Mockusdc, error) {
	parsed, err := MockusdcMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockusdcBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mockusdc{MockusdcCaller: MockusdcCaller{contract: contract}, MockusdcTransactor: MockusdcTransactor{contract: contract}, MockusdcFilterer: MockusdcFilterer{contract: contract}}, nil
}

// Mockusdc is an auto generated Go binding around an Ethereum contract.
type Mockusdc struct {
	MockusdcCaller     // Read-only binding to the contract
	MockusdcTransactor // Write-only binding to the contract
	MockusdcFilterer   // Log filterer for contract events
}

// MockusdcCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockusdcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockusdcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockusdcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockusdcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockusdcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockusdcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockusdcSession struct {
	Contract     *Mockusdc         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockusdcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockusdcCallerSession struct {
	Contract *MockusdcCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MockusdcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockusdcTransactorSession struct {
	Contract     *MockusdcTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MockusdcRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockusdcRaw struct {
	Contract *Mockusdc // Generic contract binding to access the raw methods on
}

// MockusdcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockusdcCallerRaw struct {
	Contract *MockusdcCaller // Generic read-only contract binding to access the raw methods on
}

// MockusdcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockusdcTransactorRaw struct {
	Contract *MockusdcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockusdc creates a new instance of Mockusdc, bound to a specific deployed contract.
func NewMockusdc(address common.Address, backend bind.ContractBackend) (*Mockusdc, error) {
	contract, err := bindMockusdc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mockusdc{MockusdcCaller: MockusdcCaller{contract: contract}, MockusdcTransactor: MockusdcTransactor{contract: contract}, MockusdcFilterer: MockusdcFilterer{contract: contract}}, nil
}

// NewMockusdcCaller creates a new read-only instance of Mockusdc, bound to a specific deployed contract.
func NewMockusdcCaller(address common.Address, caller bind.ContractCaller) (*MockusdcCaller, error) {
	contract, err := bindMockusdc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockusdcCaller{contract: contract}, nil
}

// NewMockusdcTransactor creates a new write-only instance of Mockusdc, bound to a specific deployed contract.
func NewMockusdcTransactor(address common.Address, transactor bind.ContractTransactor) (*MockusdcTransactor, error) {
	contract, err := bindMockusdc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockusdcTransactor{contract: contract}, nil
}

// NewMockusdcFilterer creates a new log filterer instance of Mockusdc, bound to a specific deployed contract.
func NewMockusdcFilterer(address common.Address, filterer bind.ContractFilterer) (*MockusdcFilterer, error) {
	contract, err := bindMockusdc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockusdcFilterer{contract: contract}, nil
}

// bindMockusdc binds a generic wrapper to an already deployed contract.
func bindMockusdc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockusdcMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockusdc *MockusdcRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockusdc.Contract.MockusdcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockusdc *MockusdcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockusdc.Contract.MockusdcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockusdc *MockusdcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockusdc.Contract.MockusdcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockusdc *MockusdcCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockusdc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockusdc *MockusdcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockusdc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockusdc *MockusdcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockusdc.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Mockusdc *MockusdcCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Mockusdc.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Mockusdc *MockusdcSession) Name() (string, error) {
	return _Mockusdc.Contract.Name(&_Mockusdc.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Mockusdc *MockusdcCallerSession) Name() (string, error) {
	return _Mockusdc.Contract.Name(&_Mockusdc.CallOpts)
}

// MockusdcApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Mockusdc contract.
type MockusdcApprovalIterator struct {
	Event *MockusdcApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockusdcApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockusdcApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockusdcApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockusdcApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockusdcApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockusdcApproval represents a Approval event raised by the Mockusdc contract.
type MockusdcApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Mockusdc *MockusdcFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MockusdcApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Mockusdc.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MockusdcApprovalIterator{contract: _Mockusdc.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Mockusdc *MockusdcFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MockusdcApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Mockusdc.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockusdcApproval)
				if err := _Mockusdc.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Mockusdc *MockusdcFilterer) ParseApproval(log types.Log) (*MockusdcApproval, error) {
	event := new(MockusdcApproval)
	if err := _Mockusdc.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockusdcDebugIterator is returned from FilterDebug and is used to iterate over the raw logs and unpacked data for Debug events raised by the Mockusdc contract.
type MockusdcDebugIterator struct {
	Event *MockusdcDebug // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockusdcDebugIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockusdcDebug)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockusdcDebug)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockusdcDebugIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockusdcDebugIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockusdcDebug represents a Debug event raised by the Mockusdc contract.
type MockusdcDebug struct {
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDebug is a free log retrieval operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string message)
func (_Mockusdc *MockusdcFilterer) FilterDebug(opts *bind.FilterOpts) (*MockusdcDebugIterator, error) {

	logs, sub, err := _Mockusdc.contract.FilterLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return &MockusdcDebugIterator{contract: _Mockusdc.contract, event: "Debug", logs: logs, sub: sub}, nil
}

// WatchDebug is a free log subscription operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string message)
func (_Mockusdc *MockusdcFilterer) WatchDebug(opts *bind.WatchOpts, sink chan<- *MockusdcDebug) (event.Subscription, error) {

	logs, sub, err := _Mockusdc.contract.WatchLogs(opts, "Debug")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockusdcDebug)
				if err := _Mockusdc.contract.UnpackLog(event, "Debug", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDebug is a log parse operation binding the contract event 0x7cdb51e9dbbc205231228146c3246e7f914aa6d4a33170e43ecc8e3593481d1a.
//
// Solidity: event Debug(string message)
func (_Mockusdc *MockusdcFilterer) ParseDebug(log types.Log) (*MockusdcDebug, error) {
	event := new(MockusdcDebug)
	if err := _Mockusdc.contract.UnpackLog(event, "Debug", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockusdcTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Mockusdc contract.
type MockusdcTransferIterator struct {
	Event *MockusdcTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockusdcTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockusdcTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockusdcTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockusdcTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockusdcTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockusdcTransfer represents a Transfer event raised by the Mockusdc contract.
type MockusdcTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Mockusdc *MockusdcFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockusdcTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Mockusdc.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockusdcTransferIterator{contract: _Mockusdc.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Mockusdc *MockusdcFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MockusdcTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Mockusdc.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockusdcTransfer)
				if err := _Mockusdc.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Mockusdc *MockusdcFilterer) ParseTransfer(log types.Log) (*MockusdcTransfer, error) {
	event := new(MockusdcTransfer)
	if err := _Mockusdc.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
