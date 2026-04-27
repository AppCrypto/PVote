// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract StakeManager {
    uint256 public constant PERCENT_DENOMINATOR = 100;

    address public immutable initiator;
    uint256 public immutable initiatorEscrowWei;
    uint256 public immutable voterStakeWei;
    uint256 public immutable tallierStakeWei;
    uint8 public immutable initiatorRewardPercent;
    uint8 public immutable voterRewardPercent;
    uint8 public immutable tallierRewardPercent;

    bool public initiatorEscrowFunded;
    bool public rewardsSettled;
    uint256 public totalEscrowed;
    uint256 public rewardPoolAtSettlement;
    uint256 public settledAt;

    struct Participant {
        address account;
        uint256 deposited;
        uint256 claimable;
        bool staked;
        bool honest;
        bool withdrawn;
    }

    Participant private initiatorState;
    mapping(uint256 => Participant) private voters;
    mapping(uint256 => Participant) private talliers;
    uint256[] private stakedVoterIds;
    uint256[] private stakedTallierIds;

    event InitiatorEscrowFunded(address indexed initiator, uint256 amount);
    event VoterStaked(uint256 indexed voterId, address indexed voter, uint256 amount);
    event TallierStaked(uint256 indexed tallierId, address indexed tallier, uint256 amount);
    event RewardsSettled(uint256 rewardPool, uint256 honestVoters, uint256 honestTalliers);
    event Withdrawn(address indexed account, uint256 amount);

    modifier onlyInitiator() {
        require(msg.sender == initiator, "only initiator");
        _;
    }

    modifier beforeSettlement() {
        require(!rewardsSettled, "rewards already settled");
        _;
    }

    modifier afterSettlement() {
        require(rewardsSettled, "rewards not settled");
        _;
    }

    constructor(
        uint256 _initiatorEscrowWei,
        uint256 _voterStakeWei,
        uint256 _tallierStakeWei,
        uint8 _initiatorRewardPercent,
        uint8 _voterRewardPercent,
        uint8 _tallierRewardPercent
    ) {
        require(
            uint256(_initiatorRewardPercent) + uint256(_voterRewardPercent) + uint256(_tallierRewardPercent) == PERCENT_DENOMINATOR,
            "reward split must equal 100"
        );

        initiator = msg.sender;
        initiatorEscrowWei = _initiatorEscrowWei;
        voterStakeWei = _voterStakeWei;
        tallierStakeWei = _tallierStakeWei;
        initiatorRewardPercent = _initiatorRewardPercent;
        voterRewardPercent = _voterRewardPercent;
        tallierRewardPercent = _tallierRewardPercent;
    }

    function fundInitiatorEscrow() external payable onlyInitiator beforeSettlement {
        require(!initiatorEscrowFunded, "initiator escrow already funded");
        require(msg.value == initiatorEscrowWei, "incorrect initiator escrow");

        initiatorEscrowFunded = true;
        totalEscrowed += msg.value;

        initiatorState.account = msg.sender;
        initiatorState.deposited = msg.value;
        initiatorState.staked = true;

        emit InitiatorEscrowFunded(msg.sender, msg.value);
    }

    function depositVoterStake(uint256 voterId) external payable beforeSettlement {
        require(voterId > 0, "invalid voter id");
        require(msg.value == voterStakeWei, "incorrect voter stake");

        Participant storage voter = voters[voterId];
        require(!voter.staked, "voter already staked");

        voter.account = msg.sender;
        voter.deposited = msg.value;
        voter.staked = true;
        stakedVoterIds.push(voterId);
        totalEscrowed += msg.value;

        emit VoterStaked(voterId, msg.sender, msg.value);
    }

    function depositTallierStake(uint256 tallierId) external payable beforeSettlement {
        require(tallierId > 0, "invalid tallier id");
        require(msg.value == tallierStakeWei, "incorrect tallier stake");

        Participant storage tallier = talliers[tallierId];
        require(!tallier.staked, "tallier already staked");

        tallier.account = msg.sender;
        tallier.deposited = msg.value;
        tallier.staked = true;
        stakedTallierIds.push(tallierId);
        totalEscrowed += msg.value;

        emit TallierStaked(tallierId, msg.sender, msg.value);
    }

    function settleRewards(uint256[] calldata honestVoterIds, uint256[] calldata honestTallierIds)
        external
        onlyInitiator
        beforeSettlement
    {
        require(initiatorEscrowFunded, "initiator escrow not funded");

        uint256 rewardPool = initiatorState.deposited;
        uint256 honestVoterCount = 0;
        uint256 honestTallierCount = 0;

        for (uint256 i = 0; i < honestVoterIds.length; i++) {
            Participant storage voter = voters[honestVoterIds[i]];
            require(voter.staked, "honest voter not staked");
            voter.honest = true;
        }

        for (uint256 i = 0; i < honestTallierIds.length; i++) {
            Participant storage tallier = talliers[honestTallierIds[i]];
            require(tallier.staked, "honest tallier not staked");
            tallier.honest = true;
        }

        for (uint256 i = 0; i < stakedVoterIds.length; i++) {
            Participant storage voter = voters[stakedVoterIds[i]];
            if (voter.honest) {
                voter.claimable += voter.deposited;
                honestVoterCount++;
            } else {
                rewardPool += voter.deposited;
            }
        }

        for (uint256 i = 0; i < stakedTallierIds.length; i++) {
            Participant storage tallier = talliers[stakedTallierIds[i]];
            if (tallier.honest) {
                tallier.claimable += tallier.deposited;
                honestTallierCount++;
            } else {
                rewardPool += tallier.deposited;
            }
        }

        rewardPoolAtSettlement = rewardPool;

        uint256 initiatorShare = (rewardPool * initiatorRewardPercent) / PERCENT_DENOMINATOR;
        uint256 voterSharePool = (rewardPool * voterRewardPercent) / PERCENT_DENOMINATOR;
        uint256 tallierSharePool = rewardPool - initiatorShare - voterSharePool;

        initiatorState.claimable += initiatorShare;

        if (honestVoterCount == 0) {
            initiatorState.claimable += voterSharePool;
        } else {
            uint256 eachVoterShare = voterSharePool / honestVoterCount;
            uint256 voterRemainder = voterSharePool % honestVoterCount;
            for (uint256 i = 0; i < stakedVoterIds.length; i++) {
                Participant storage voter = voters[stakedVoterIds[i]];
                if (voter.honest) {
                    voter.claimable += eachVoterShare;
                }
            }
            initiatorState.claimable += voterRemainder;
        }

        if (honestTallierCount == 0) {
            initiatorState.claimable += tallierSharePool;
        } else {
            uint256 eachTallierShare = tallierSharePool / honestTallierCount;
            uint256 tallierRemainder = tallierSharePool % honestTallierCount;
            for (uint256 i = 0; i < stakedTallierIds.length; i++) {
                Participant storage tallier = talliers[stakedTallierIds[i]];
                if (tallier.honest) {
                    tallier.claimable += eachTallierShare;
                }
            }
            initiatorState.claimable += tallierRemainder;
        }

        rewardsSettled = true;
        settledAt = block.timestamp;

        emit RewardsSettled(rewardPool, honestVoterCount, honestTallierCount);
    }

    function withdrawInitiator() external onlyInitiator afterSettlement {
        _withdraw(initiatorState, initiatorState.account);
    }

    function withdrawVoter(uint256 voterId) external afterSettlement {
        Participant storage voter = voters[voterId];
        require(voter.account == msg.sender, "only voter account");
        _withdraw(voter, voter.account);
    }

    function withdrawTallier(uint256 tallierId) external afterSettlement {
        Participant storage tallier = talliers[tallierId];
        require(tallier.account == msg.sender, "only tallier account");
        _withdraw(tallier, tallier.account);
    }

    function getInitiatorState()
        external
        view
        returns (address account, uint256 deposited, uint256 claimable, bool staked, bool withdrawn)
    {
        return (
            initiatorState.account,
            initiatorState.deposited,
            initiatorState.claimable,
            initiatorState.staked,
            initiatorState.withdrawn
        );
    }

    function getVoter(uint256 voterId)
        external
        view
        returns (address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
    {
        Participant storage voter = voters[voterId];
        return (voter.account, voter.deposited, voter.claimable, voter.staked, voter.honest, voter.withdrawn);
    }

    function getTallier(uint256 tallierId)
        external
        view
        returns (address account, uint256 deposited, uint256 claimable, bool staked, bool honest, bool withdrawn)
    {
        Participant storage tallier = talliers[tallierId];
        return (tallier.account, tallier.deposited, tallier.claimable, tallier.staked, tallier.honest, tallier.withdrawn);
    }

    function getStakedVoterIds() external view returns (uint256[] memory) {
        return stakedVoterIds;
    }

    function getStakedTallierIds() external view returns (uint256[] memory) {
        return stakedTallierIds;
    }

    function getEscrowOverview()
        external
        view
        returns (
            bool escrowFunded,
            bool settled,
            uint256 totalEscrow,
            uint256 rewardPool,
            uint256 contractBalance,
            uint256 voterCount,
            uint256 tallierCount,
            uint256 settledTimestamp
        )
    {
        return (
            initiatorEscrowFunded,
            rewardsSettled,
            totalEscrowed,
            rewardPoolAtSettlement,
            address(this).balance,
            stakedVoterIds.length,
            stakedTallierIds.length,
            settledAt
        );
    }

    function _withdraw(Participant storage participant, address account) private {
        require(participant.staked, "participant not staked");
        require(!participant.withdrawn, "already withdrawn");
        require(participant.claimable > 0, "nothing to withdraw");

        uint256 amount = participant.claimable;
        participant.claimable = 0;
        participant.withdrawn = true;

        (bool ok, ) = payable(account).call{value: amount}("");
        require(ok, "withdraw transfer failed");

        emit Withdrawn(account, amount);
    }
}
