#!/bin/bash

set -ev

#chaincode install
docker exec cli peer chaincode install -n waste -v 1.0 -p github.com/waste/go
sleep 3

#chaincode instatiate
docker exec cli peer chaincode instantiate -n waste -v 1.0 -C mychannel1 -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member")'
sleep 5

#chaincode invoke user1
docker exec cli peer chaincode invoke -n waste -C mychannel1 -c '{"Args":["initLedger"]}'
sleep 5

#chaincode query user1
# Waste{ObjectType, Phase: args[1], Name: args[2], Color: args[3], CasNo: args[4]}
docker exec cli peer chaincode query -n waste -C mychannel1 -c '{"Args":["createWaste","4", "Gas", "Cyanogen chioride", "None", "506-77-4"]}'


echo '-------------------------------------END-------------------------------------'