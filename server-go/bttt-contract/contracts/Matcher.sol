pragma solidity ^0.4.0;


contract Matcher {
    address owner;

    uint public minimumMid = 100;

    struct Stranger {
    address user;
    uint bid;
    }

    // Key indicates an order or max bid
    mapping (uint8 => Stranger) strangers;

    function Matcher(){

    }

    function getOrderMag(uint input) constant returns (uint8 counter){
        input = input / 10;
        while (input >= 1) {
            counter++;
            input = input / 10;
        }
        return counter;
    }

    function playWithStranger(uint8 minBidOrder) public payable returns (uint gameID) {
        require((msg.value == 0 && minBidOrder == 0) || msg.value > minimumMid);
        uint8 bidMagnitudeOrder = getOrderMag(msg.value());
        require(minBidOrder <= maxBidOrder);

        while (bidMagnitudeOrder >= 0) {
            Stranger stranger = strangers[bidMagnitudeOrder];
            if (stranger.user != address(0x0)) {
                stranger.user = 0x0;
                stranger.bid = 0;
            }
            bidMagnitudeOrder -= 1;
            return;
        }
        strangers[bidMagnitudeOrder] = Stranger({user : msg.sender, bid : msg.value});
    }
}
