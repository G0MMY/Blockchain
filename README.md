# Blockchain

This is a simple blockchain. The language that I used is Golang and I used levelDB to store the chain state. I did this project to learn Golang and the fundamentals of blockchains. The main branch is my first attempt, which isn't really good.

## How does it work?

The blockchain that I did only has nodes and miners. The nodes hold the whole state of the blockchain and the miners are the ones creating the new nodes. When a node receives a new block or transaction, it checks if its correct. If it is, it sends it to all the other nodes. It's quite simple, but it works at really small scale. 

## Installation

1. Download the code
1. Go in your terminal and cd in the blockchain directory
1. Execute "go install blockchain"
1. Execute "export PATH=$PATH:GOPATH/bin" (GOPATH is the path to your go directory)
1. To get the list of the commands, run "blockchain"
1. To get more details on each command, you can type "blockchain theCommand -h"

## Example to test the blockchain

I recommend using Postman to do requests to the nodes

1. Open 3 termials
2. Execute "export PATH=$PATH:GOPATH/bin" in each one
3. Run "blockchain init --port 4000" in the first one
4. Go to this link to get a private key: http://localhost:4000/create/wallet
5. Run "blockchain miner --port 4001 --node 4000 --privateKey yourPrivateKey in the second terminal
6. Run "blockchain node --port 4002 --neighbor 4000" in the third one
7. You now have the blockchain running with two nodes and one miner
8. You can send transactions to one of the nodes with the "blockchain transaction" command
9. You can also view the chain at this link: http://localhost:portOfANode/chain
10. You can get a list of all the endpoints in the FullNode.go file in the Network folder
