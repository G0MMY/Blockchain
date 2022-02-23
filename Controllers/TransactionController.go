package Controllers

import (
	"blockchain/Models"
	"fmt"
)

func CreateTransaction(from []byte, to []byte, amount int, fee int, outputs []Models.Output, memPoolId int) (Models.MemPoolTransaction, []Models.MemPoolInput, []Models.MemPoolOutput) {
	var memPooloutputs []Models.MemPoolOutput
	for _, output := range outputs {
		if amount-output.Amount >= 0 {

			amount -= output.Amount
		}
	}
	if amount != 0 {
		return Models.UnspentTransaction{}, []Models.Input{}, []Models.Output{}
	}

}

func filterOutputs(outputs []Models.Output, publicKey []byte) []Models.Output {
	var result []Models.Output
	for _, output := range outputs {
		if fmt.Sprintf("%s", output.PublicKey) == fmt.Sprintf("%s", publicKey) && output.InputId == -1 {
			result = append(result, output)
		}
	}
	return result
}
