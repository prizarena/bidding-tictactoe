pragma solidity ^0.4.17;


contract BiddingTicTacToe {

    // Play at:
    //
    //   * Web: https://BiddingTicTacToe.com
    //
    //   * Telegram: https://t.me/BiddingTicTacToeBot

    address public owner; // Creator of the contract
    address public fees; // 1% fee is accumulated here on top-ups
    //    address public server; // Controls & changes balances

    //        uint txPrice = tx.gasprice * gasCostToWithdraw;

    enum Player {None, X, O}

    struct Game {
    address PlayerX;
    address PlayerO;
    uint bank;
    Player winner;
    }

    mapping (uint => Game) games;

    uint constant gasCostToStartGame = 1000;

    uint constant gasCostToJoinGame = 1000;

    uint constant gasCostToBid = 1000;

    uint constant gasCostToWithdraw = 22000;

    mapping (address => uint) balances;


    function topUp() public payable returns (uint balance) {
        require(msg.value > 100);
        uint fee = msg.value / 100;
        balance = balances[msg.sender] + msg.value - fee;
        balances[msg.sender] = balance;
        fees.transfer(fee);
    }

    function withdraw(uint value) public {
        var balance = balances[msg.sender];
        require(balance > 0 && value <= balance);
        if (value == 0) {
            value = balance;
        }
        balances[msg.sender] = balance - value;
        msg.sender.transfer(value);
    }

    function userBalance(address user) public constant returns (uint balance) {
        return balances[user];
    }


    uint lastGameID;

    //    function uintToBytes(uint v) private pure returns (bytes32 ret) { // https://github.com/pipermerriam/ethereum-string-utils
    //        if (v == 0) {
    //            ret = '0';
    //        }
    //        else {
    //            while (v > 0) {
    //                ret = bytes32(uint(ret) / (2 ** 8));
    //                ret |= bytes32(((v % 10) + 48) * 2 ** (8 * 31));
    //                v /= 10;
    //            }
    //        }
    //        return ret;
    //    }
}
