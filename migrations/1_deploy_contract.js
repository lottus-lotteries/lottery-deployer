const lotteryMachine = artifacts.require("lotteryMachine.sol");

module.exports = function(deployer) {
    deployer.deploy(lotteryMachine);
};
