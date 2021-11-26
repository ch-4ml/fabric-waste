#!/bin/bash

set -ev

#chaincode install
docker exec cli peer chaincode install -n recycle -v 1.7 -p github.com/recycle/go

#chaincode instatiate
docker exec cli peer chaincode instantiate -n recycle -v 1.7 -C mychannel2 -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org3MSP.member")'
sleep 10

#chaincode invoke user1
docker exec cli peer chaincode invoke -n recycle -C mychannel2 -c '{"Args":["initLedger"]}'
sleep 5

#chaincode query user1
# docker exec cli peer chaincode query -n recycle -C mychannel2 -c '{"Args":["",""]}'


echo '-------------------------------------END-------------------------------------'