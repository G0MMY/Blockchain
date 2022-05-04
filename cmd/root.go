package cmd

import (
	"blockchain/Models"
	"blockchain/Network"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var port string
var neighbor string
var node string
var privateKey string
var to string
var amount int
var fee int
var timestamp int

var rootCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "A simple blockchain",
	Long:  `Blockchain is a simple blockchain that I did to learn about Golang and blockchains.`,
}

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Create a new transaction",
	Long:  "Transaction is to create a new transaction on the blockchain.",

	Run: func(cmd *cobra.Command, args []string) {
		node, _ := cmd.Flags().GetString("node")
		privateKey, _ := cmd.Flags().GetString("privateKey")
		to, _ := cmd.Flags().GetString("receiver")
		amount, _ := cmd.Flags().GetInt("amount")
		fee, _ := cmd.Flags().GetInt("fee")
		timestamp, _ := cmd.Flags().GetInt("timestamp")

		transactionRequest := Models.CreateTransactionRequest{
			privateKey,
			to,
			amount,
			fee,
			int64(timestamp),
		}

		transactionRequest.AddToNode(node)
	},
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Initialize a node",
	Long:  `Node is to initialize a new node. You will have to pass the port of the node.`,

	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		neighbor, _ := cmd.Flags().GetString("neighbor")

		Network.InitializeNode(port, neighbor)
	},
}

var minerCmd = &cobra.Command{
	Use:   "miner",
	Short: "Initialize a miner",
	Long:  `Miner is to initialize a new miner. You will have to pass the port of the miner.`,

	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		node, _ := cmd.Flags().GetString("node")
		privateKey, _ := cmd.Flags().GetString("privateKey")
		bytePrivateKey, _ := hex.DecodeString(privateKey)

		Network.InitializeMiner(port, node, bytePrivateKey)
	},
}

func Execute() {
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(minerCmd)
	rootCmd.AddCommand(transactionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	nodeCmd.PersistentFlags().StringVar(&port, "port", "", "The port that the node will run on")
	nodeCmd.PersistentFlags().StringVar(&neighbor, "neighbor", "", "A known node to get the current network of the blockchain")
	nodeCmd.MarkPersistentFlagRequired("port")

	minerCmd.PersistentFlags().StringVar(&port, "port", "", "The port that the miner will run on")
	minerCmd.PersistentFlags().StringVar(&node, "node", "", "The port of the node that the miner will talk to")
	minerCmd.PersistentFlags().StringVar(&privateKey, "privateKey", "", "The private key of the miner to receive the money for mining blocks")
	minerCmd.MarkPersistentFlagRequired("port")
	minerCmd.MarkPersistentFlagRequired("node")
	minerCmd.MarkPersistentFlagRequired("privateKey")

	transactionCmd.PersistentFlags().StringVar(&node, "node", "", "The port of the node that the transaction will be sent to")
	transactionCmd.PersistentFlags().StringVar(&privateKey, "privateKey", "", "The private key of the sender of the transaction")
	transactionCmd.PersistentFlags().StringVar(&to, "receiver", "", "The address or the public key of the receiver of the transaction")
	transactionCmd.PersistentFlags().IntVar(&amount, "amount", 0, "The amount of the transaction")
	transactionCmd.PersistentFlags().IntVar(&fee, "fee", 0, "The fee of the transaction")
	transactionCmd.PersistentFlags().IntVar(&timestamp, "timestamp", 0, "The timestamp of the transaction (when it can be accepted in the blockchain)")
	transactionCmd.MarkPersistentFlagRequired("node")
	transactionCmd.MarkPersistentFlagRequired("privateKey")
	transactionCmd.MarkPersistentFlagRequired("receiver")
	transactionCmd.MarkPersistentFlagRequired("amount")
	transactionCmd.MarkPersistentFlagRequired("fee")
	transactionCmd.MarkPersistentFlagRequired("timestamp")
}
