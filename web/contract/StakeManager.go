// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakecontract

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

// StakecontractMetaData contains all meta data concerning the Stakecontract contract.
var StakecontractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_initiatorEscrowWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_voterStakeWei\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_tallierStakeWei\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_initiatorRewardPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_voterRewardPercent\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_tallierRewardPercent\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InitiatorEscrowFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewardPool\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"honestVoters\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"honestTalliers\",\"type\":\"uint256\"}],\"name\":\"RewardsSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tallierId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tallier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TallierStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"voterId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"VoterStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PERCENT_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tallierId\",\"type\":\"uint256\"}],\"name\":\"depositTallierStake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"voterId\",\"type\":\"uint256\"}],\"name\":\"depositVoterStake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fundInitiatorEscrow\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEscrowOverview\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"escrowFunded\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"settled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"totalEscrow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewardPool\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"voterCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tallierCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"settledTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInitiatorState\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposited\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimable\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"staked\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"withdrawn\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakedTallierIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakedVoterIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tallierId\",\"type\":\"uint256\"}],\"name\":\"getTallier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposited\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimable\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"staked\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"honest\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"withdrawn\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"voterId\",\"type\":\"uint256\"}],\"name\":\"getVoter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deposited\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimable\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"staked\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"honest\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"withdrawn\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiatorEscrowFunded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiatorEscrowWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initiatorRewardPercent\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardPoolAtSettlement\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsSettled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"honestVoterIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"honestTallierIds\",\"type\":\"uint256[]\"}],\"name\":\"settleRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settledAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tallierRewardPercent\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tallierStakeWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalEscrowed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"voterRewardPercent\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"voterStakeWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawInitiator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tallierId\",\"type\":\"uint256\"}],\"name\":\"withdrawTallier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"voterId\",\"type\":\"uint256\"}],\"name\":\"withdrawVoter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101603461017c57601f6119b138819003918201601f19168301916001600160401b038311848410176101815780849260c09460405283398101031261017c57805190602081015191604082015161005960608401610197565b9361007260a061006b60808701610197565b9501610197565b92606461009261008860ff881660ff8a166101a5565b60ff8716906101a5565b03610137573360805260a05260c05260e052610100928352610120918252610140908152604051916117e893846101c9853960805184818161046b0152818161074d015281816109430152610b03015260a0518481816104b6015261098a015260c0518481816111d0015261121f015260e05184818161065f015261101d0152518381816103000152610c7401525182818161017b0152610ca4015251816102c20152f35b60405162461bcd60e51b815260206004820152601b60248201527f7265776172642073706c6974206d75737420657175616c2031303000000000006044820152606490fd5b600080fd5b634e487b7160e01b600052604160045260246000fd5b519060ff8216820361017c57565b919082018092116101b257565b634e487b7160e01b600052601160045260246000fdfe6040608081526004908136101561001557600080fd5b600091823560e01c90816306709e51146111f3578163071415c5146111b8578163087d276314610ff157816309aa1bf914610aac57816313500fcc14610a5a57816321def7f0146109d2578163447ff822146109ad578163584d7a47146109725781635c39fcc11461092e57816377b4671d1461072d5781638147eb5a1461070b5781638481e076146106825781638c8f30e31461064757816395722851146105cc578163958c0298146104585781639e6c29591461043c578163a267e8c0146103dc578163a99a67e6146103bd578163a9ad160a14610343578163ac4d62df14610324578163b3cadede146102e6578163c14b171d146102a8578163d07bff0c1461022a57508063e7714bb41461019f578063ec56d44c146101625763f91682311461014157600080fd5b3461015e578160031936011261015e576020906001549051908152f35b5080fd5b503461015e578160031936011261015e576020905160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b5090346102275780600319360112610227578151918291600b54808552602080950194600b83527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db992905b8282106102105761020c8686610202828b03836115b5565b519182918261139e565b0390f35b8354875295860195600193840193909101906101ea565b80fd5b9050346102a45760203660031901126102a457358252600860208181529282902080546001820154600283015460039093015494516001600160a01b03909216825294810194909452604084015260ff808316151560608501529082901c81161515608084015260109190911c16151560a082015260c090f35b8280fd5b50503461015e578160031936011261015e576020905160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b50503461015e578160031936011261015e576020905160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b50503461015e578160031936011261015e576020906002549051908152f35b8284346102275780600319360112610227578151918291600a54808552602080950194600a83527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a892905b8282106103a65761020c8686610202828b03836115b5565b83548752958601956001938401939091019061038e565b50503461015e578160031936011261015e576020906003549051908152f35b50503461015e578160031936011261015e5761010091549060015460025447600a5491600b54936003549560ff81519881811615158a5260081c16151560208901528701526060860152608085015260a084015260c083015260e0820152f35b50503461015e578160031936011261015e576020905160648152f35b9050826003193601126102a457610499337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316146114cd565b8254906104ac60ff8360081c16156113da565b60ff8216610589577f0000000000000000000000000000000000000000000000000000000000000000340361054657906001918260ff198093161785556104f43484546114aa565b8355336bffffffffffffffffffffffff60a01b825416179055346005556007541617600755513481527f61e204da9a6dffdf6068c69850105af9ac413532055f3fc7d5c56c3e50a1e2de60203392a280f35b606490602084519162461bcd60e51b8352820152601a60248201527f696e636f727265637420696e69746961746f7220657363726f770000000000006044820152fd5b606490602084519162461bcd60e51b8352820152601f60248201527f696e69746961746f7220657363726f7720616c72656164792066756e646564006044820152fd5b9050346102a45760203660031901126102a457358252600960209081529181902080546001820154600283015460039093015493516001600160a01b03909216825293810193909352604083015260ff80821615156060840152600882901c81161515608084015260109190911c16151560a082015260c090f35b50503461015e578160031936011261015e57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b8391503461015e57602036600319011261015e576106a660ff835460081c16611573565b80358252600860205282822080546001600160a01b0316913383036106d35750906106d0916115ed565b80f35b606490602086519162461bcd60e51b835282015260126024820152711bdb9b1e481d9bdd195c881858d8dbdd5b9d60721b6044820152fd5b50503461015e578160031936011261015e5760ff602092541690519015158152f35b9050346102a457826003193601126102a4576001600160a01b03610774337f00000000000000000000000000000000000000000000000000000000000000008316146114cd565b61078460ff855460081c16611573565b815416916007549160ff8316156108f35760ff8360101c166108bd5760065492831561088457600686905562ff0000191662010000176007558480808086885af13d1561087f573d67ffffffffffffffff811161086c578351906107f2601f8201601f1916602001836115b5565b81528660203d92013e5b1561082f57507f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d59160209151908152a280f35b6020606492519162461bcd60e51b835282015260186024820152771dda5d1a191c985dc81d1c985b9cd9995c8819985a5b195960421b6044820152fd5b634e487b7160e01b875260418352602487fd5b6107fc565b506020606492519162461bcd60e51b835282015260136024820152726e6f7468696e6720746f20776974686472617760681b6044820152fd5b6020606492519162461bcd60e51b8352820152601160248201527030b63932b0b23c903bb4ba34323930bbb760791b6044820152fd5b6020606492519162461bcd60e51b835282015260166024820152751c185c9d1a58da5c185b9d081b9bdd081cdd185ad95960521b6044820152fd5b50503461015e578160031936011261015e57517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b50503461015e578160031936011261015e57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b50503461015e578160031936011261015e5760ff6020925460081c1690519015158152f35b8391503461015e57602036600319011261015e576109f660ff835460081c16611573565b80358252600960205282822080546001600160a01b031691338303610a205750906106d0916115ed565b606490602086519162461bcd60e51b835282015260146024820152731bdb9b1e481d185b1b1a595c881858d8dbdd5b9d60621b6044820152fd5b8284346102275780600319360112610227575060ff60a092600180851b0390541691600554906006549060075492815195865260208601528401528181161515606084015260101c1615156080820152f35b9050823461022757826003193601126102275767ffffffffffffffff82358181116102a457610ade9036908501611368565b939091602490602435908111610fed57610afb9036908401611368565b9290610b31337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316146114cd565b85549160ff9384600894610b498282881c16156113da565b1615610faa57600554958897899a8c8b80915b838310610f265750925050505b818110610ea957505050505084935b600a54851015610bed57610b8b85611426565b9054600391821b1c8752826020528389882091820154841c16600014610bd55760019181610bc5600285610bcc95015492019182546114aa565b905561151a565b945b0193610b78565b600191959482610be7920154906114aa565b93610bce565b96958594505b600b54851015610c6057610c0685611473565b9054600391821b1c8752600960205287872090810154831c841615610c485760019181610bc5600285610c3f95015492019182546114aa565b945b0193610bf3565b600191959482610c5a920154906114aa565b93610c41565b9692949593506002856002556064610c9a847f00000000000000000000000000000000000000000000000000000000000000001688611529565b04926064610cca827f00000000000000000000000000000000000000000000000000000000000000001689611529565b0492610cdf84610cda878b61155c565b61155c565b93610ced60069687546114aa565b865586610e1857610cff9086546114aa565b85555b89610d63575050509181610d3e7f2bf582ae7b4e23267f9742193453efcc823164c91b245c2d83075e61e4c5c3e99798946060979694546114aa565b90555b865461ff0019166101001787554260035581519384526020840152820152a180f35b929690939795919498610d8081610d7a818b61153c565b99611569565b99875b600b54811015610dd6578088610d9a600193611473565b9054600391821b1c8c52600960205289888d20918201548a1c16610dc1575b505001610d83565b01610dcd8c82546114aa565b9055888e610db9565b50935094509497606096507f2bf582ae7b4e23267f9742193453efcc823164c91b245c2d83075e61e4c5c3e997919350610e119082546114aa565b9055610d41565b979194610e3687610e30818c9d9b979a95989d61153c565b9a611569565b99885b600a54811015610e8c5780898b8a89610e53600196611426565b939054600394851b1c81528c60205220918201548b1c16610e77575b505001610e39565b01610e838d82546114aa565b9055898f610e6f565b50949196939799610ea2919693995086546114aa565b8555610d02565b610eb481838761150a565b358a5260206009815260038d8c20019081549089821615610ee3575061ff001916610100179055600101610b69565b8e5162461bcd60e51b8152808801919091526019818701527f686f6e6573742074616c6c696572206e6f74207374616b6564000000000000006044820152606490fd5b600390610f3484868861150a565b3583526020928b84522001908154908b821615610f61575061ff0019166101001790556001018b8e610b5c565b9450505050506064945060179291508b519362461bcd60e51b85528401528201527f686f6e65737420766f746572206e6f74207374616b65640000000000000000006044820152fd5b895162461bcd60e51b8152602081840152601b60248201527f696e69746961746f7220657363726f77206e6f742066756e64656400000000006044820152606490fd5b8480fd5b9050602091826003193601126111b45781359161101560ff865460081c16156113da565b821561117e577f0000000000000000000000000000000000000000000000000000000000000000340361113d5782855260098452818520600381019081549060ff82166111015780546001600160a01b031916331781553460019182015560ff19909116179055600b5490600160401b8210156110ee5750906110bf8260017f7a1efa3a5e76c30a8b2b2be3cfda3ea7fe7206df75794d952683ff5b8d703f249401600b55611473565b81549060031b9085821b91600019901b19161790556110e0346001546114aa565b60015551923484523393a380f35b634e487b7160e01b865260419052602485fd5b845162461bcd60e51b815280850188905260166024820152751d185b1b1a595c88185b1c9958591e481cdd185ad95960521b6044820152606490fd5b83606492519162461bcd60e51b8352820152601760248201527f696e636f72726563742074616c6c696572207374616b650000000000000000006044820152fd5b83606492519162461bcd60e51b835282015260126024820152711a5b9d985b1a59081d185b1b1a595c881a5960721b6044820152fd5b8380fd5b50503461015e578160031936011261015e57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b9050602091826003193601126111b45781359161121760ff865460081c16156113da565b8215611334577f000000000000000000000000000000000000000000000000000000000000000034036112fb5782855260088452818520600381019081549060ff82166112c15780546001600160a01b031916331781553460019182015560ff19909116179055600a5490600160401b8210156110ee5750906110bf8260017f93abcbcdc77902bad39329f481f8da6e20cf5da7a1a13c7b95fdc4fa19da005e9401600a55611426565b845162461bcd60e51b815280850188905260146024820152731d9bdd195c88185b1c9958591e481cdd185ad95960621b6044820152606490fd5b83606492519162461bcd60e51b83528201526015602482015274696e636f727265637420766f746572207374616b6560581b6044820152fd5b83606492519162461bcd60e51b8352820152601060248201526f1a5b9d985b1a59081d9bdd195c881a5960821b6044820152fd5b9181601f840112156113995782359167ffffffffffffffff8311611399576020808501948460051b01011161139957565b600080fd5b602090602060408183019282815285518094520193019160005b8281106113c6575050505090565b8351855293810193928101926001016113b8565b156113e157565b60405162461bcd60e51b815260206004820152601760248201527f7265776172647320616c726561647920736574746c65640000000000000000006044820152606490fd5b600a5481101561145d57600a6000527fc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a80190600090565b634e487b7160e01b600052603260045260246000fd5b600b5481101561145d57600b6000527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db90190600090565b919082018092116114b757565b634e487b7160e01b600052601160045260246000fd5b156114d457565b60405162461bcd60e51b815260206004820152600e60248201526d37b7363c9034b734ba34b0ba37b960911b6044820152606490fd5b919081101561145d5760051b0190565b60001981146114b75760010190565b818102929181159184041417156114b757565b8115611546570490565b634e487b7160e01b600052601260045260246000fd5b919082039182116114b757565b8115611546570690565b1561157a57565b60405162461bcd60e51b81526020600482015260136024820152721c995dd85c991cc81b9bdd081cd95d1d1b1959606a1b6044820152606490fd5b90601f8019910116810190811067ffffffffffffffff8211176115d757604052565b634e487b7160e01b600052604160045260246000fd5b600381019182549160ff8316156117745760ff8360101c1661173b576002019283549283156117005760009485905562ff00001916620100001790556001600160a01b0316918080808085875af1903d156116fa573d9067ffffffffffffffff82116116e6576040519161166b601f8201601f1916602001846115b5565b825260203d92013e5b156116a65760207f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d591604051908152a2565b60405162461bcd60e51b81526020600482015260186024820152771dda5d1a191c985dc81d1c985b9cd9995c8819985a5b195960421b6044820152606490fd5b634e487b7160e01b81526041600452602490fd5b50611674565b60405162461bcd60e51b81526020600482015260136024820152726e6f7468696e6720746f20776974686472617760681b6044820152606490fd5b60405162461bcd60e51b815260206004820152601160248201527030b63932b0b23c903bb4ba34323930bbb760791b6044820152606490fd5b60405162461bcd60e51b81526020600482015260166024820152751c185c9d1a58da5c185b9d081b9bdd081cdd185ad95960521b6044820152606490fdfea2646970667358221220c08f2f9c20c0b4bd297c5e1a7c3b2a94b5d2a043c148c52d53b0cb1219cb7b1c64736f6c63430008190033",
}

// StakecontractABI is the input ABI used to generate the binding from.
// Deprecated: Use StakecontractMetaData.ABI instead.
var StakecontractABI = StakecontractMetaData.ABI

// StakecontractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakecontractMetaData.Bin instead.
var StakecontractBin = StakecontractMetaData.Bin

// DeployStakecontract deploys a new Ethereum contract, binding an instance of Stakecontract to it.
func DeployStakecontract(auth *bind.TransactOpts, backend bind.ContractBackend, _initiatorEscrowWei *big.Int, _voterStakeWei *big.Int, _tallierStakeWei *big.Int, _initiatorRewardPercent uint8, _voterRewardPercent uint8, _tallierRewardPercent uint8) (common.Address, *types.Transaction, *Stakecontract, error) {
	parsed, err := StakecontractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakecontractBin), backend, _initiatorEscrowWei, _voterStakeWei, _tallierStakeWei, _initiatorRewardPercent, _voterRewardPercent, _tallierRewardPercent)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Stakecontract{StakecontractCaller: StakecontractCaller{contract: contract}, StakecontractTransactor: StakecontractTransactor{contract: contract}, StakecontractFilterer: StakecontractFilterer{contract: contract}}, nil
}

// Stakecontract is an auto generated Go binding around an Ethereum contract.
type Stakecontract struct {
	StakecontractCaller     // Read-only binding to the contract
	StakecontractTransactor // Write-only binding to the contract
	StakecontractFilterer   // Log filterer for contract events
}

// StakecontractCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakecontractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakecontractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakecontractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakecontractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakecontractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakecontractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakecontractSession struct {
	Contract     *Stakecontract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakecontractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakecontractCallerSession struct {
	Contract *StakecontractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StakecontractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakecontractTransactorSession struct {
	Contract     *StakecontractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StakecontractRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakecontractRaw struct {
	Contract *Stakecontract // Generic contract binding to access the raw methods on
}

// StakecontractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakecontractCallerRaw struct {
	Contract *StakecontractCaller // Generic read-only contract binding to access the raw methods on
}

// StakecontractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakecontractTransactorRaw struct {
	Contract *StakecontractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakecontract creates a new instance of Stakecontract, bound to a specific deployed contract.
func NewStakecontract(address common.Address, backend bind.ContractBackend) (*Stakecontract, error) {
	contract, err := bindStakecontract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stakecontract{StakecontractCaller: StakecontractCaller{contract: contract}, StakecontractTransactor: StakecontractTransactor{contract: contract}, StakecontractFilterer: StakecontractFilterer{contract: contract}}, nil
}

// NewStakecontractCaller creates a new read-only instance of Stakecontract, bound to a specific deployed contract.
func NewStakecontractCaller(address common.Address, caller bind.ContractCaller) (*StakecontractCaller, error) {
	contract, err := bindStakecontract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakecontractCaller{contract: contract}, nil
}

// NewStakecontractTransactor creates a new write-only instance of Stakecontract, bound to a specific deployed contract.
func NewStakecontractTransactor(address common.Address, transactor bind.ContractTransactor) (*StakecontractTransactor, error) {
	contract, err := bindStakecontract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakecontractTransactor{contract: contract}, nil
}

// NewStakecontractFilterer creates a new log filterer instance of Stakecontract, bound to a specific deployed contract.
func NewStakecontractFilterer(address common.Address, filterer bind.ContractFilterer) (*StakecontractFilterer, error) {
	contract, err := bindStakecontract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakecontractFilterer{contract: contract}, nil
}

// bindStakecontract binds a generic wrapper to an already deployed contract.
func bindStakecontract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakecontractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakecontract *StakecontractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakecontract.Contract.StakecontractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakecontract *StakecontractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakecontract.Contract.StakecontractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakecontract *StakecontractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakecontract.Contract.StakecontractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakecontract *StakecontractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakecontract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakecontract *StakecontractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakecontract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakecontract *StakecontractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakecontract.Contract.contract.Transact(opts, method, params...)
}

// PERCENTDENOMINATOR is a free data retrieval call binding the contract method 0x9e6c2959.
//
// Solidity: function PERCENT_DENOMINATOR() view returns(uint256)
func (_Stakecontract *StakecontractCaller) PERCENTDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "PERCENT_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PERCENTDENOMINATOR is a free data retrieval call binding the contract method 0x9e6c2959.
//
// Solidity: function PERCENT_DENOMINATOR() view returns(uint256)
func (_Stakecontract *StakecontractSession) PERCENTDENOMINATOR() (*big.Int, error) {
	return _Stakecontract.Contract.PERCENTDENOMINATOR(&_Stakecontract.CallOpts)
}

// PERCENTDENOMINATOR is a free data retrieval call binding the contract method 0x9e6c2959.
//
// Solidity: function PERCENT_DENOMINATOR() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) PERCENTDENOMINATOR() (*big.Int, error) {
	return _Stakecontract.Contract.PERCENTDENOMINATOR(&_Stakecontract.CallOpts)
}

// GetEscrowOverview is a free data retrieval call binding the contract method 0xa267e8c0.
//
// Solidity: function getEscrowOverview() view returns(bool escrowFunded, bool settled, uint256 totalEscrow, uint256 rewardPool, uint256 contractBalance, uint256 voterCount, uint256 tallierCount, uint256 settledTimestamp)
func (_Stakecontract *StakecontractCaller) GetEscrowOverview(opts *bind.CallOpts) (struct {
	EscrowFunded     bool
	Settled          bool
	TotalEscrow      *big.Int
	RewardPool       *big.Int
	ContractBalance  *big.Int
	VoterCount       *big.Int
	TallierCount     *big.Int
	SettledTimestamp *big.Int
}, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getEscrowOverview")

	outstruct := new(struct {
		EscrowFunded     bool
		Settled          bool
		TotalEscrow      *big.Int
		RewardPool       *big.Int
		ContractBalance  *big.Int
		VoterCount       *big.Int
		TallierCount     *big.Int
		SettledTimestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EscrowFunded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Settled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.TotalEscrow = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.RewardPool = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.ContractBalance = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.VoterCount = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TallierCount = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.SettledTimestamp = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetEscrowOverview is a free data retrieval call binding the contract method 0xa267e8c0.
//
// Solidity: function getEscrowOverview() view returns(bool escrowFunded, bool settled, uint256 totalEscrow, uint256 rewardPool, uint256 contractBalance, uint256 voterCount, uint256 tallierCount, uint256 settledTimestamp)
func (_Stakecontract *StakecontractSession) GetEscrowOverview() (struct {
	EscrowFunded     bool
	Settled          bool
	TotalEscrow      *big.Int
	RewardPool       *big.Int
	ContractBalance  *big.Int
	VoterCount       *big.Int
	TallierCount     *big.Int
	SettledTimestamp *big.Int
}, error) {
	return _Stakecontract.Contract.GetEscrowOverview(&_Stakecontract.CallOpts)
}

// GetEscrowOverview is a free data retrieval call binding the contract method 0xa267e8c0.
//
// Solidity: function getEscrowOverview() view returns(bool escrowFunded, bool settled, uint256 totalEscrow, uint256 rewardPool, uint256 contractBalance, uint256 voterCount, uint256 tallierCount, uint256 settledTimestamp)
func (_Stakecontract *StakecontractCallerSession) GetEscrowOverview() (struct {
	EscrowFunded     bool
	Settled          bool
	TotalEscrow      *big.Int
	RewardPool       *big.Int
	ContractBalance  *big.Int
	VoterCount       *big.Int
	TallierCount     *big.Int
	SettledTimestamp *big.Int
}, error) {
	return _Stakecontract.Contract.GetEscrowOverview(&_Stakecontract.CallOpts)
}

// GetInitiatorState is a free data retrieval call binding the contract method 0x13500fcc.
//
// Solidity: function getInitiatorState() view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool withdrawn)
func (_Stakecontract *StakecontractCaller) GetInitiatorState(opts *bind.CallOpts) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Withdrawn bool
}, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getInitiatorState")

	outstruct := new(struct {
		Account   common.Address
		Deposited *big.Int
		Claimable *big.Int
		Staked    bool
		Withdrawn bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Account = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Deposited = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Claimable = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Staked = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.Withdrawn = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// GetInitiatorState is a free data retrieval call binding the contract method 0x13500fcc.
//
// Solidity: function getInitiatorState() view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool withdrawn)
func (_Stakecontract *StakecontractSession) GetInitiatorState() (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetInitiatorState(&_Stakecontract.CallOpts)
}

// GetInitiatorState is a free data retrieval call binding the contract method 0x13500fcc.
//
// Solidity: function getInitiatorState() view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool withdrawn)
func (_Stakecontract *StakecontractCallerSession) GetInitiatorState() (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetInitiatorState(&_Stakecontract.CallOpts)
}

// GetStakedTallierIds is a free data retrieval call binding the contract method 0xe7714bb4.
//
// Solidity: function getStakedTallierIds() view returns(uint256[])
func (_Stakecontract *StakecontractCaller) GetStakedTallierIds(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getStakedTallierIds")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetStakedTallierIds is a free data retrieval call binding the contract method 0xe7714bb4.
//
// Solidity: function getStakedTallierIds() view returns(uint256[])
func (_Stakecontract *StakecontractSession) GetStakedTallierIds() ([]*big.Int, error) {
	return _Stakecontract.Contract.GetStakedTallierIds(&_Stakecontract.CallOpts)
}

// GetStakedTallierIds is a free data retrieval call binding the contract method 0xe7714bb4.
//
// Solidity: function getStakedTallierIds() view returns(uint256[])
func (_Stakecontract *StakecontractCallerSession) GetStakedTallierIds() ([]*big.Int, error) {
	return _Stakecontract.Contract.GetStakedTallierIds(&_Stakecontract.CallOpts)
}

// GetStakedVoterIds is a free data retrieval call binding the contract method 0xa9ad160a.
//
// Solidity: function getStakedVoterIds() view returns(uint256[])
func (_Stakecontract *StakecontractCaller) GetStakedVoterIds(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getStakedVoterIds")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetStakedVoterIds is a free data retrieval call binding the contract method 0xa9ad160a.
//
// Solidity: function getStakedVoterIds() view returns(uint256[])
func (_Stakecontract *StakecontractSession) GetStakedVoterIds() ([]*big.Int, error) {
	return _Stakecontract.Contract.GetStakedVoterIds(&_Stakecontract.CallOpts)
}

// GetStakedVoterIds is a free data retrieval call binding the contract method 0xa9ad160a.
//
// Solidity: function getStakedVoterIds() view returns(uint256[])
func (_Stakecontract *StakecontractCallerSession) GetStakedVoterIds() ([]*big.Int, error) {
	return _Stakecontract.Contract.GetStakedVoterIds(&_Stakecontract.CallOpts)
}

// GetTallier is a free data retrieval call binding the contract method 0x95722851.
//
// Solidity: function getTallier(uint256 tallierId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractCaller) GetTallier(opts *bind.CallOpts, tallierId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getTallier", tallierId)

	outstruct := new(struct {
		Account   common.Address
		Deposited *big.Int
		Claimable *big.Int
		Staked    bool
		Honest    bool
		Withdrawn bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Account = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Deposited = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Claimable = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Staked = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.Honest = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.Withdrawn = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// GetTallier is a free data retrieval call binding the contract method 0x95722851.
//
// Solidity: function getTallier(uint256 tallierId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractSession) GetTallier(tallierId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetTallier(&_Stakecontract.CallOpts, tallierId)
}

// GetTallier is a free data retrieval call binding the contract method 0x95722851.
//
// Solidity: function getTallier(uint256 tallierId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractCallerSession) GetTallier(tallierId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetTallier(&_Stakecontract.CallOpts, tallierId)
}

// GetVoter is a free data retrieval call binding the contract method 0xd07bff0c.
//
// Solidity: function getVoter(uint256 voterId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractCaller) GetVoter(opts *bind.CallOpts, voterId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "getVoter", voterId)

	outstruct := new(struct {
		Account   common.Address
		Deposited *big.Int
		Claimable *big.Int
		Staked    bool
		Honest    bool
		Withdrawn bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Account = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Deposited = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Claimable = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Staked = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.Honest = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.Withdrawn = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// GetVoter is a free data retrieval call binding the contract method 0xd07bff0c.
//
// Solidity: function getVoter(uint256 voterId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractSession) GetVoter(voterId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetVoter(&_Stakecontract.CallOpts, voterId)
}

// GetVoter is a free data retrieval call binding the contract method 0xd07bff0c.
//
// Solidity: function getVoter(uint256 voterId) view returns(address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
func (_Stakecontract *StakecontractCallerSession) GetVoter(voterId *big.Int) (struct {
	Account   common.Address
	Deposited *big.Int
	Claimable *big.Int
	Staked    bool
	Honest    bool
	Withdrawn bool
}, error) {
	return _Stakecontract.Contract.GetVoter(&_Stakecontract.CallOpts, voterId)
}

// Initiator is a free data retrieval call binding the contract method 0x5c39fcc1.
//
// Solidity: function initiator() view returns(address)
func (_Stakecontract *StakecontractCaller) Initiator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "initiator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Initiator is a free data retrieval call binding the contract method 0x5c39fcc1.
//
// Solidity: function initiator() view returns(address)
func (_Stakecontract *StakecontractSession) Initiator() (common.Address, error) {
	return _Stakecontract.Contract.Initiator(&_Stakecontract.CallOpts)
}

// Initiator is a free data retrieval call binding the contract method 0x5c39fcc1.
//
// Solidity: function initiator() view returns(address)
func (_Stakecontract *StakecontractCallerSession) Initiator() (common.Address, error) {
	return _Stakecontract.Contract.Initiator(&_Stakecontract.CallOpts)
}

// InitiatorEscrowFunded is a free data retrieval call binding the contract method 0x8147eb5a.
//
// Solidity: function initiatorEscrowFunded() view returns(bool)
func (_Stakecontract *StakecontractCaller) InitiatorEscrowFunded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "initiatorEscrowFunded")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// InitiatorEscrowFunded is a free data retrieval call binding the contract method 0x8147eb5a.
//
// Solidity: function initiatorEscrowFunded() view returns(bool)
func (_Stakecontract *StakecontractSession) InitiatorEscrowFunded() (bool, error) {
	return _Stakecontract.Contract.InitiatorEscrowFunded(&_Stakecontract.CallOpts)
}

// InitiatorEscrowFunded is a free data retrieval call binding the contract method 0x8147eb5a.
//
// Solidity: function initiatorEscrowFunded() view returns(bool)
func (_Stakecontract *StakecontractCallerSession) InitiatorEscrowFunded() (bool, error) {
	return _Stakecontract.Contract.InitiatorEscrowFunded(&_Stakecontract.CallOpts)
}

// InitiatorEscrowWei is a free data retrieval call binding the contract method 0x584d7a47.
//
// Solidity: function initiatorEscrowWei() view returns(uint256)
func (_Stakecontract *StakecontractCaller) InitiatorEscrowWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "initiatorEscrowWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitiatorEscrowWei is a free data retrieval call binding the contract method 0x584d7a47.
//
// Solidity: function initiatorEscrowWei() view returns(uint256)
func (_Stakecontract *StakecontractSession) InitiatorEscrowWei() (*big.Int, error) {
	return _Stakecontract.Contract.InitiatorEscrowWei(&_Stakecontract.CallOpts)
}

// InitiatorEscrowWei is a free data retrieval call binding the contract method 0x584d7a47.
//
// Solidity: function initiatorEscrowWei() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) InitiatorEscrowWei() (*big.Int, error) {
	return _Stakecontract.Contract.InitiatorEscrowWei(&_Stakecontract.CallOpts)
}

// InitiatorRewardPercent is a free data retrieval call binding the contract method 0xb3cadede.
//
// Solidity: function initiatorRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCaller) InitiatorRewardPercent(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "initiatorRewardPercent")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// InitiatorRewardPercent is a free data retrieval call binding the contract method 0xb3cadede.
//
// Solidity: function initiatorRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractSession) InitiatorRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.InitiatorRewardPercent(&_Stakecontract.CallOpts)
}

// InitiatorRewardPercent is a free data retrieval call binding the contract method 0xb3cadede.
//
// Solidity: function initiatorRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCallerSession) InitiatorRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.InitiatorRewardPercent(&_Stakecontract.CallOpts)
}

// RewardPoolAtSettlement is a free data retrieval call binding the contract method 0xac4d62df.
//
// Solidity: function rewardPoolAtSettlement() view returns(uint256)
func (_Stakecontract *StakecontractCaller) RewardPoolAtSettlement(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "rewardPoolAtSettlement")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPoolAtSettlement is a free data retrieval call binding the contract method 0xac4d62df.
//
// Solidity: function rewardPoolAtSettlement() view returns(uint256)
func (_Stakecontract *StakecontractSession) RewardPoolAtSettlement() (*big.Int, error) {
	return _Stakecontract.Contract.RewardPoolAtSettlement(&_Stakecontract.CallOpts)
}

// RewardPoolAtSettlement is a free data retrieval call binding the contract method 0xac4d62df.
//
// Solidity: function rewardPoolAtSettlement() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) RewardPoolAtSettlement() (*big.Int, error) {
	return _Stakecontract.Contract.RewardPoolAtSettlement(&_Stakecontract.CallOpts)
}

// RewardsSettled is a free data retrieval call binding the contract method 0x447ff822.
//
// Solidity: function rewardsSettled() view returns(bool)
func (_Stakecontract *StakecontractCaller) RewardsSettled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "rewardsSettled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RewardsSettled is a free data retrieval call binding the contract method 0x447ff822.
//
// Solidity: function rewardsSettled() view returns(bool)
func (_Stakecontract *StakecontractSession) RewardsSettled() (bool, error) {
	return _Stakecontract.Contract.RewardsSettled(&_Stakecontract.CallOpts)
}

// RewardsSettled is a free data retrieval call binding the contract method 0x447ff822.
//
// Solidity: function rewardsSettled() view returns(bool)
func (_Stakecontract *StakecontractCallerSession) RewardsSettled() (bool, error) {
	return _Stakecontract.Contract.RewardsSettled(&_Stakecontract.CallOpts)
}

// SettledAt is a free data retrieval call binding the contract method 0xa99a67e6.
//
// Solidity: function settledAt() view returns(uint256)
func (_Stakecontract *StakecontractCaller) SettledAt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "settledAt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SettledAt is a free data retrieval call binding the contract method 0xa99a67e6.
//
// Solidity: function settledAt() view returns(uint256)
func (_Stakecontract *StakecontractSession) SettledAt() (*big.Int, error) {
	return _Stakecontract.Contract.SettledAt(&_Stakecontract.CallOpts)
}

// SettledAt is a free data retrieval call binding the contract method 0xa99a67e6.
//
// Solidity: function settledAt() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) SettledAt() (*big.Int, error) {
	return _Stakecontract.Contract.SettledAt(&_Stakecontract.CallOpts)
}

// TallierRewardPercent is a free data retrieval call binding the contract method 0xc14b171d.
//
// Solidity: function tallierRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCaller) TallierRewardPercent(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "tallierRewardPercent")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TallierRewardPercent is a free data retrieval call binding the contract method 0xc14b171d.
//
// Solidity: function tallierRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractSession) TallierRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.TallierRewardPercent(&_Stakecontract.CallOpts)
}

// TallierRewardPercent is a free data retrieval call binding the contract method 0xc14b171d.
//
// Solidity: function tallierRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCallerSession) TallierRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.TallierRewardPercent(&_Stakecontract.CallOpts)
}

// TallierStakeWei is a free data retrieval call binding the contract method 0x8c8f30e3.
//
// Solidity: function tallierStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractCaller) TallierStakeWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "tallierStakeWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TallierStakeWei is a free data retrieval call binding the contract method 0x8c8f30e3.
//
// Solidity: function tallierStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractSession) TallierStakeWei() (*big.Int, error) {
	return _Stakecontract.Contract.TallierStakeWei(&_Stakecontract.CallOpts)
}

// TallierStakeWei is a free data retrieval call binding the contract method 0x8c8f30e3.
//
// Solidity: function tallierStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) TallierStakeWei() (*big.Int, error) {
	return _Stakecontract.Contract.TallierStakeWei(&_Stakecontract.CallOpts)
}

// TotalEscrowed is a free data retrieval call binding the contract method 0xf9168231.
//
// Solidity: function totalEscrowed() view returns(uint256)
func (_Stakecontract *StakecontractCaller) TotalEscrowed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "totalEscrowed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalEscrowed is a free data retrieval call binding the contract method 0xf9168231.
//
// Solidity: function totalEscrowed() view returns(uint256)
func (_Stakecontract *StakecontractSession) TotalEscrowed() (*big.Int, error) {
	return _Stakecontract.Contract.TotalEscrowed(&_Stakecontract.CallOpts)
}

// TotalEscrowed is a free data retrieval call binding the contract method 0xf9168231.
//
// Solidity: function totalEscrowed() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) TotalEscrowed() (*big.Int, error) {
	return _Stakecontract.Contract.TotalEscrowed(&_Stakecontract.CallOpts)
}

// VoterRewardPercent is a free data retrieval call binding the contract method 0xec56d44c.
//
// Solidity: function voterRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCaller) VoterRewardPercent(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "voterRewardPercent")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// VoterRewardPercent is a free data retrieval call binding the contract method 0xec56d44c.
//
// Solidity: function voterRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractSession) VoterRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.VoterRewardPercent(&_Stakecontract.CallOpts)
}

// VoterRewardPercent is a free data retrieval call binding the contract method 0xec56d44c.
//
// Solidity: function voterRewardPercent() view returns(uint8)
func (_Stakecontract *StakecontractCallerSession) VoterRewardPercent() (uint8, error) {
	return _Stakecontract.Contract.VoterRewardPercent(&_Stakecontract.CallOpts)
}

// VoterStakeWei is a free data retrieval call binding the contract method 0x071415c5.
//
// Solidity: function voterStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractCaller) VoterStakeWei(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakecontract.contract.Call(opts, &out, "voterStakeWei")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VoterStakeWei is a free data retrieval call binding the contract method 0x071415c5.
//
// Solidity: function voterStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractSession) VoterStakeWei() (*big.Int, error) {
	return _Stakecontract.Contract.VoterStakeWei(&_Stakecontract.CallOpts)
}

// VoterStakeWei is a free data retrieval call binding the contract method 0x071415c5.
//
// Solidity: function voterStakeWei() view returns(uint256)
func (_Stakecontract *StakecontractCallerSession) VoterStakeWei() (*big.Int, error) {
	return _Stakecontract.Contract.VoterStakeWei(&_Stakecontract.CallOpts)
}

// DepositTallierStake is a paid mutator transaction binding the contract method 0x087d2763.
//
// Solidity: function depositTallierStake(uint256 tallierId) payable returns()
func (_Stakecontract *StakecontractTransactor) DepositTallierStake(opts *bind.TransactOpts, tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "depositTallierStake", tallierId)
}

// DepositTallierStake is a paid mutator transaction binding the contract method 0x087d2763.
//
// Solidity: function depositTallierStake(uint256 tallierId) payable returns()
func (_Stakecontract *StakecontractSession) DepositTallierStake(tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.DepositTallierStake(&_Stakecontract.TransactOpts, tallierId)
}

// DepositTallierStake is a paid mutator transaction binding the contract method 0x087d2763.
//
// Solidity: function depositTallierStake(uint256 tallierId) payable returns()
func (_Stakecontract *StakecontractTransactorSession) DepositTallierStake(tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.DepositTallierStake(&_Stakecontract.TransactOpts, tallierId)
}

// DepositVoterStake is a paid mutator transaction binding the contract method 0x06709e51.
//
// Solidity: function depositVoterStake(uint256 voterId) payable returns()
func (_Stakecontract *StakecontractTransactor) DepositVoterStake(opts *bind.TransactOpts, voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "depositVoterStake", voterId)
}

// DepositVoterStake is a paid mutator transaction binding the contract method 0x06709e51.
//
// Solidity: function depositVoterStake(uint256 voterId) payable returns()
func (_Stakecontract *StakecontractSession) DepositVoterStake(voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.DepositVoterStake(&_Stakecontract.TransactOpts, voterId)
}

// DepositVoterStake is a paid mutator transaction binding the contract method 0x06709e51.
//
// Solidity: function depositVoterStake(uint256 voterId) payable returns()
func (_Stakecontract *StakecontractTransactorSession) DepositVoterStake(voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.DepositVoterStake(&_Stakecontract.TransactOpts, voterId)
}

// FundInitiatorEscrow is a paid mutator transaction binding the contract method 0x958c0298.
//
// Solidity: function fundInitiatorEscrow() payable returns()
func (_Stakecontract *StakecontractTransactor) FundInitiatorEscrow(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "fundInitiatorEscrow")
}

// FundInitiatorEscrow is a paid mutator transaction binding the contract method 0x958c0298.
//
// Solidity: function fundInitiatorEscrow() payable returns()
func (_Stakecontract *StakecontractSession) FundInitiatorEscrow() (*types.Transaction, error) {
	return _Stakecontract.Contract.FundInitiatorEscrow(&_Stakecontract.TransactOpts)
}

// FundInitiatorEscrow is a paid mutator transaction binding the contract method 0x958c0298.
//
// Solidity: function fundInitiatorEscrow() payable returns()
func (_Stakecontract *StakecontractTransactorSession) FundInitiatorEscrow() (*types.Transaction, error) {
	return _Stakecontract.Contract.FundInitiatorEscrow(&_Stakecontract.TransactOpts)
}

// SettleRewards is a paid mutator transaction binding the contract method 0x09aa1bf9.
//
// Solidity: function settleRewards(uint256[] honestVoterIds, uint256[] honestTallierIds) returns()
func (_Stakecontract *StakecontractTransactor) SettleRewards(opts *bind.TransactOpts, honestVoterIds []*big.Int, honestTallierIds []*big.Int) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "settleRewards", honestVoterIds, honestTallierIds)
}

// SettleRewards is a paid mutator transaction binding the contract method 0x09aa1bf9.
//
// Solidity: function settleRewards(uint256[] honestVoterIds, uint256[] honestTallierIds) returns()
func (_Stakecontract *StakecontractSession) SettleRewards(honestVoterIds []*big.Int, honestTallierIds []*big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.SettleRewards(&_Stakecontract.TransactOpts, honestVoterIds, honestTallierIds)
}

// SettleRewards is a paid mutator transaction binding the contract method 0x09aa1bf9.
//
// Solidity: function settleRewards(uint256[] honestVoterIds, uint256[] honestTallierIds) returns()
func (_Stakecontract *StakecontractTransactorSession) SettleRewards(honestVoterIds []*big.Int, honestTallierIds []*big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.SettleRewards(&_Stakecontract.TransactOpts, honestVoterIds, honestTallierIds)
}

// WithdrawInitiator is a paid mutator transaction binding the contract method 0x77b4671d.
//
// Solidity: function withdrawInitiator() returns()
func (_Stakecontract *StakecontractTransactor) WithdrawInitiator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "withdrawInitiator")
}

// WithdrawInitiator is a paid mutator transaction binding the contract method 0x77b4671d.
//
// Solidity: function withdrawInitiator() returns()
func (_Stakecontract *StakecontractSession) WithdrawInitiator() (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawInitiator(&_Stakecontract.TransactOpts)
}

// WithdrawInitiator is a paid mutator transaction binding the contract method 0x77b4671d.
//
// Solidity: function withdrawInitiator() returns()
func (_Stakecontract *StakecontractTransactorSession) WithdrawInitiator() (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawInitiator(&_Stakecontract.TransactOpts)
}

// WithdrawTallier is a paid mutator transaction binding the contract method 0x21def7f0.
//
// Solidity: function withdrawTallier(uint256 tallierId) returns()
func (_Stakecontract *StakecontractTransactor) WithdrawTallier(opts *bind.TransactOpts, tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "withdrawTallier", tallierId)
}

// WithdrawTallier is a paid mutator transaction binding the contract method 0x21def7f0.
//
// Solidity: function withdrawTallier(uint256 tallierId) returns()
func (_Stakecontract *StakecontractSession) WithdrawTallier(tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawTallier(&_Stakecontract.TransactOpts, tallierId)
}

// WithdrawTallier is a paid mutator transaction binding the contract method 0x21def7f0.
//
// Solidity: function withdrawTallier(uint256 tallierId) returns()
func (_Stakecontract *StakecontractTransactorSession) WithdrawTallier(tallierId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawTallier(&_Stakecontract.TransactOpts, tallierId)
}

// WithdrawVoter is a paid mutator transaction binding the contract method 0x8481e076.
//
// Solidity: function withdrawVoter(uint256 voterId) returns()
func (_Stakecontract *StakecontractTransactor) WithdrawVoter(opts *bind.TransactOpts, voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.contract.Transact(opts, "withdrawVoter", voterId)
}

// WithdrawVoter is a paid mutator transaction binding the contract method 0x8481e076.
//
// Solidity: function withdrawVoter(uint256 voterId) returns()
func (_Stakecontract *StakecontractSession) WithdrawVoter(voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawVoter(&_Stakecontract.TransactOpts, voterId)
}

// WithdrawVoter is a paid mutator transaction binding the contract method 0x8481e076.
//
// Solidity: function withdrawVoter(uint256 voterId) returns()
func (_Stakecontract *StakecontractTransactorSession) WithdrawVoter(voterId *big.Int) (*types.Transaction, error) {
	return _Stakecontract.Contract.WithdrawVoter(&_Stakecontract.TransactOpts, voterId)
}

// StakecontractInitiatorEscrowFundedIterator is returned from FilterInitiatorEscrowFunded and is used to iterate over the raw logs and unpacked data for InitiatorEscrowFunded events raised by the Stakecontract contract.
type StakecontractInitiatorEscrowFundedIterator struct {
	Event *StakecontractInitiatorEscrowFunded // Event containing the contract specifics and raw log

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
func (it *StakecontractInitiatorEscrowFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakecontractInitiatorEscrowFunded)
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
		it.Event = new(StakecontractInitiatorEscrowFunded)
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
func (it *StakecontractInitiatorEscrowFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakecontractInitiatorEscrowFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakecontractInitiatorEscrowFunded represents a InitiatorEscrowFunded event raised by the Stakecontract contract.
type StakecontractInitiatorEscrowFunded struct {
	Initiator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInitiatorEscrowFunded is a free log retrieval operation binding the contract event 0x61e204da9a6dffdf6068c69850105af9ac413532055f3fc7d5c56c3e50a1e2de.
//
// Solidity: event InitiatorEscrowFunded(address indexed initiator, uint256 amount)
func (_Stakecontract *StakecontractFilterer) FilterInitiatorEscrowFunded(opts *bind.FilterOpts, initiator []common.Address) (*StakecontractInitiatorEscrowFundedIterator, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}

	logs, sub, err := _Stakecontract.contract.FilterLogs(opts, "InitiatorEscrowFunded", initiatorRule)
	if err != nil {
		return nil, err
	}
	return &StakecontractInitiatorEscrowFundedIterator{contract: _Stakecontract.contract, event: "InitiatorEscrowFunded", logs: logs, sub: sub}, nil
}

// WatchInitiatorEscrowFunded is a free log subscription operation binding the contract event 0x61e204da9a6dffdf6068c69850105af9ac413532055f3fc7d5c56c3e50a1e2de.
//
// Solidity: event InitiatorEscrowFunded(address indexed initiator, uint256 amount)
func (_Stakecontract *StakecontractFilterer) WatchInitiatorEscrowFunded(opts *bind.WatchOpts, sink chan<- *StakecontractInitiatorEscrowFunded, initiator []common.Address) (event.Subscription, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}

	logs, sub, err := _Stakecontract.contract.WatchLogs(opts, "InitiatorEscrowFunded", initiatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakecontractInitiatorEscrowFunded)
				if err := _Stakecontract.contract.UnpackLog(event, "InitiatorEscrowFunded", log); err != nil {
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

// ParseInitiatorEscrowFunded is a log parse operation binding the contract event 0x61e204da9a6dffdf6068c69850105af9ac413532055f3fc7d5c56c3e50a1e2de.
//
// Solidity: event InitiatorEscrowFunded(address indexed initiator, uint256 amount)
func (_Stakecontract *StakecontractFilterer) ParseInitiatorEscrowFunded(log types.Log) (*StakecontractInitiatorEscrowFunded, error) {
	event := new(StakecontractInitiatorEscrowFunded)
	if err := _Stakecontract.contract.UnpackLog(event, "InitiatorEscrowFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakecontractRewardsSettledIterator is returned from FilterRewardsSettled and is used to iterate over the raw logs and unpacked data for RewardsSettled events raised by the Stakecontract contract.
type StakecontractRewardsSettledIterator struct {
	Event *StakecontractRewardsSettled // Event containing the contract specifics and raw log

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
func (it *StakecontractRewardsSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakecontractRewardsSettled)
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
		it.Event = new(StakecontractRewardsSettled)
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
func (it *StakecontractRewardsSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakecontractRewardsSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakecontractRewardsSettled represents a RewardsSettled event raised by the Stakecontract contract.
type StakecontractRewardsSettled struct {
	RewardPool     *big.Int
	HonestVoters   *big.Int
	HonestTalliers *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRewardsSettled is a free log retrieval operation binding the contract event 0x2bf582ae7b4e23267f9742193453efcc823164c91b245c2d83075e61e4c5c3e9.
//
// Solidity: event RewardsSettled(uint256 rewardPool, uint256 honestVoters, uint256 honestTalliers)
func (_Stakecontract *StakecontractFilterer) FilterRewardsSettled(opts *bind.FilterOpts) (*StakecontractRewardsSettledIterator, error) {

	logs, sub, err := _Stakecontract.contract.FilterLogs(opts, "RewardsSettled")
	if err != nil {
		return nil, err
	}
	return &StakecontractRewardsSettledIterator{contract: _Stakecontract.contract, event: "RewardsSettled", logs: logs, sub: sub}, nil
}

// WatchRewardsSettled is a free log subscription operation binding the contract event 0x2bf582ae7b4e23267f9742193453efcc823164c91b245c2d83075e61e4c5c3e9.
//
// Solidity: event RewardsSettled(uint256 rewardPool, uint256 honestVoters, uint256 honestTalliers)
func (_Stakecontract *StakecontractFilterer) WatchRewardsSettled(opts *bind.WatchOpts, sink chan<- *StakecontractRewardsSettled) (event.Subscription, error) {

	logs, sub, err := _Stakecontract.contract.WatchLogs(opts, "RewardsSettled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakecontractRewardsSettled)
				if err := _Stakecontract.contract.UnpackLog(event, "RewardsSettled", log); err != nil {
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

// ParseRewardsSettled is a log parse operation binding the contract event 0x2bf582ae7b4e23267f9742193453efcc823164c91b245c2d83075e61e4c5c3e9.
//
// Solidity: event RewardsSettled(uint256 rewardPool, uint256 honestVoters, uint256 honestTalliers)
func (_Stakecontract *StakecontractFilterer) ParseRewardsSettled(log types.Log) (*StakecontractRewardsSettled, error) {
	event := new(StakecontractRewardsSettled)
	if err := _Stakecontract.contract.UnpackLog(event, "RewardsSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakecontractTallierStakedIterator is returned from FilterTallierStaked and is used to iterate over the raw logs and unpacked data for TallierStaked events raised by the Stakecontract contract.
type StakecontractTallierStakedIterator struct {
	Event *StakecontractTallierStaked // Event containing the contract specifics and raw log

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
func (it *StakecontractTallierStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakecontractTallierStaked)
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
		it.Event = new(StakecontractTallierStaked)
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
func (it *StakecontractTallierStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakecontractTallierStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakecontractTallierStaked represents a TallierStaked event raised by the Stakecontract contract.
type StakecontractTallierStaked struct {
	TallierId *big.Int
	Tallier   common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTallierStaked is a free log retrieval operation binding the contract event 0x7a1efa3a5e76c30a8b2b2be3cfda3ea7fe7206df75794d952683ff5b8d703f24.
//
// Solidity: event TallierStaked(uint256 indexed tallierId, address indexed tallier, uint256 amount)
func (_Stakecontract *StakecontractFilterer) FilterTallierStaked(opts *bind.FilterOpts, tallierId []*big.Int, tallier []common.Address) (*StakecontractTallierStakedIterator, error) {

	var tallierIdRule []interface{}
	for _, tallierIdItem := range tallierId {
		tallierIdRule = append(tallierIdRule, tallierIdItem)
	}
	var tallierRule []interface{}
	for _, tallierItem := range tallier {
		tallierRule = append(tallierRule, tallierItem)
	}

	logs, sub, err := _Stakecontract.contract.FilterLogs(opts, "TallierStaked", tallierIdRule, tallierRule)
	if err != nil {
		return nil, err
	}
	return &StakecontractTallierStakedIterator{contract: _Stakecontract.contract, event: "TallierStaked", logs: logs, sub: sub}, nil
}

// WatchTallierStaked is a free log subscription operation binding the contract event 0x7a1efa3a5e76c30a8b2b2be3cfda3ea7fe7206df75794d952683ff5b8d703f24.
//
// Solidity: event TallierStaked(uint256 indexed tallierId, address indexed tallier, uint256 amount)
func (_Stakecontract *StakecontractFilterer) WatchTallierStaked(opts *bind.WatchOpts, sink chan<- *StakecontractTallierStaked, tallierId []*big.Int, tallier []common.Address) (event.Subscription, error) {

	var tallierIdRule []interface{}
	for _, tallierIdItem := range tallierId {
		tallierIdRule = append(tallierIdRule, tallierIdItem)
	}
	var tallierRule []interface{}
	for _, tallierItem := range tallier {
		tallierRule = append(tallierRule, tallierItem)
	}

	logs, sub, err := _Stakecontract.contract.WatchLogs(opts, "TallierStaked", tallierIdRule, tallierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakecontractTallierStaked)
				if err := _Stakecontract.contract.UnpackLog(event, "TallierStaked", log); err != nil {
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

// ParseTallierStaked is a log parse operation binding the contract event 0x7a1efa3a5e76c30a8b2b2be3cfda3ea7fe7206df75794d952683ff5b8d703f24.
//
// Solidity: event TallierStaked(uint256 indexed tallierId, address indexed tallier, uint256 amount)
func (_Stakecontract *StakecontractFilterer) ParseTallierStaked(log types.Log) (*StakecontractTallierStaked, error) {
	event := new(StakecontractTallierStaked)
	if err := _Stakecontract.contract.UnpackLog(event, "TallierStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakecontractVoterStakedIterator is returned from FilterVoterStaked and is used to iterate over the raw logs and unpacked data for VoterStaked events raised by the Stakecontract contract.
type StakecontractVoterStakedIterator struct {
	Event *StakecontractVoterStaked // Event containing the contract specifics and raw log

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
func (it *StakecontractVoterStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakecontractVoterStaked)
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
		it.Event = new(StakecontractVoterStaked)
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
func (it *StakecontractVoterStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakecontractVoterStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakecontractVoterStaked represents a VoterStaked event raised by the Stakecontract contract.
type StakecontractVoterStaked struct {
	VoterId *big.Int
	Voter   common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVoterStaked is a free log retrieval operation binding the contract event 0x93abcbcdc77902bad39329f481f8da6e20cf5da7a1a13c7b95fdc4fa19da005e.
//
// Solidity: event VoterStaked(uint256 indexed voterId, address indexed voter, uint256 amount)
func (_Stakecontract *StakecontractFilterer) FilterVoterStaked(opts *bind.FilterOpts, voterId []*big.Int, voter []common.Address) (*StakecontractVoterStakedIterator, error) {

	var voterIdRule []interface{}
	for _, voterIdItem := range voterId {
		voterIdRule = append(voterIdRule, voterIdItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _Stakecontract.contract.FilterLogs(opts, "VoterStaked", voterIdRule, voterRule)
	if err != nil {
		return nil, err
	}
	return &StakecontractVoterStakedIterator{contract: _Stakecontract.contract, event: "VoterStaked", logs: logs, sub: sub}, nil
}

// WatchVoterStaked is a free log subscription operation binding the contract event 0x93abcbcdc77902bad39329f481f8da6e20cf5da7a1a13c7b95fdc4fa19da005e.
//
// Solidity: event VoterStaked(uint256 indexed voterId, address indexed voter, uint256 amount)
func (_Stakecontract *StakecontractFilterer) WatchVoterStaked(opts *bind.WatchOpts, sink chan<- *StakecontractVoterStaked, voterId []*big.Int, voter []common.Address) (event.Subscription, error) {

	var voterIdRule []interface{}
	for _, voterIdItem := range voterId {
		voterIdRule = append(voterIdRule, voterIdItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _Stakecontract.contract.WatchLogs(opts, "VoterStaked", voterIdRule, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakecontractVoterStaked)
				if err := _Stakecontract.contract.UnpackLog(event, "VoterStaked", log); err != nil {
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

// ParseVoterStaked is a log parse operation binding the contract event 0x93abcbcdc77902bad39329f481f8da6e20cf5da7a1a13c7b95fdc4fa19da005e.
//
// Solidity: event VoterStaked(uint256 indexed voterId, address indexed voter, uint256 amount)
func (_Stakecontract *StakecontractFilterer) ParseVoterStaked(log types.Log) (*StakecontractVoterStaked, error) {
	event := new(StakecontractVoterStaked)
	if err := _Stakecontract.contract.UnpackLog(event, "VoterStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakecontractWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Stakecontract contract.
type StakecontractWithdrawnIterator struct {
	Event *StakecontractWithdrawn // Event containing the contract specifics and raw log

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
func (it *StakecontractWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakecontractWithdrawn)
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
		it.Event = new(StakecontractWithdrawn)
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
func (it *StakecontractWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakecontractWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakecontractWithdrawn represents a Withdrawn event raised by the Stakecontract contract.
type StakecontractWithdrawn struct {
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed account, uint256 amount)
func (_Stakecontract *StakecontractFilterer) FilterWithdrawn(opts *bind.FilterOpts, account []common.Address) (*StakecontractWithdrawnIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stakecontract.contract.FilterLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return &StakecontractWithdrawnIterator{contract: _Stakecontract.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed account, uint256 amount)
func (_Stakecontract *StakecontractFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *StakecontractWithdrawn, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Stakecontract.contract.WatchLogs(opts, "Withdrawn", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakecontractWithdrawn)
				if err := _Stakecontract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed account, uint256 amount)
func (_Stakecontract *StakecontractFilterer) ParseWithdrawn(log types.Log) (*StakecontractWithdrawn, error) {
	event := new(StakecontractWithdrawn)
	if err := _Stakecontract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
