#!/bin/bash

C_YELLOW='\033[1;33m'
C_BLUE='\033[0;34m'
C_RESET='\033[0m'

# subinfoln echos in blue color
function infoln() {
  echo -e "${C_YELLOW}${1}${C_RESET}"
}

function subinfoln() {
  echo -e "${C_BLUE}${1}${C_RESET}"
}

# add PATH to ensure we are picking up the correct binaries
export PATH=${HOME}/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/network/config

# Test Setting
CC_NAME="myBooking"
CHANNEL_NAME="mychannel"
PEER_CONN_PARMS="--peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/network/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt"
ORDERER_CA=${PWD}/network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## Set the peer on peer0.org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# test 수행
# ## TEST1 : Invoking the chaincode
# infoln "TEST1-1 : Invoking the chaincode (RegisterAccommodation)"
# set -x
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"RegisterAccommodation","Args":["test1", "admin", "hotel1", "room1", "100", "울산"]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt
# sleep 3

## TEST1 : Invoking the chaincode
# infoln "TEST1-2 : Invoking the chaincode (RegisterAccommodation)"
# set -x
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"RegisterAccommodation","Args":["test2", "admin", "hotel1", "room2", "200", "울산"]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt
# sleep 3

## TEST1 : Invoking the chaincode
# infoln "TEST1-3 : Invoking the chaincode (UpdateAccommodation)"
# set -x
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"UpdateAccommodation","Args":["test2", "admin", "accommodation4",  "hotel1", "room2", "300", "부산"]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt
# sleep 3

## TEST1 : Invoking the chaincode
# infoln "TEST1-3 : Invoking the chaincode (DeleteAccommodationInfo)"
# set -x
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"DeleteAccommodation","Args":["test2", "user", "accommodation2"]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt
# sleep 3

# ## TEST1 : Query the chaincode
# infoln "TEST2-1 : Query the chaincode (GetAllAccommodations)"
# set -x
# peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function":"GetAllAccommodations","Args":[]}' >&log.txt
# { set +x; } 2>/dev/null
# cat log.txt
# sleep 3

# TEST1 : Invoking the chaincode
infoln "TEST1-3 : Invoking the chaincode (Booking)"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS -c '{"function":"Booking","Args":["test2", "user","tom", "accommodation4",  "hotel1", "room2", "20231129", "20231130", "300", "card"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3

