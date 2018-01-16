package blockchain

import (
	"letsgo/wallet"
	"letsgo/util"
)


type WalletInfo struct {
	Ws []WalletStatus
	Amount int
}

//Structure représent un output non dépensé
type WalletStatus struct {
	Address []byte
	Amount int
	Wallet *wallet.Wallet
}

type UnspentOutput struct {
	TxID []byte
	Output int
	Amount int
}

type LocalUnspentOutput struct {
	TxID []byte
	Output int
	Amount int
	Wallet *wallet.Wallet
}

func GetWalletInfo() *WalletInfo {
	wInfo := &WalletInfo{}

	for _, w := range wallet.WalletList {
		amount, _ := UTXO.FindSpendableOutputsByPubKeyHash(util.Sha256(w.PublicKey), MAX_COIN)
		ws := WalletStatus{w.GetAddress(), amount, w}
		
		wInfo.Ws = append(wInfo.Ws, ws)
		wInfo.Amount += amount
	}
	return wInfo
}

func (wInfo *WalletInfo) GetLocalUnspentOutputs(amount int) (int, []LocalUnspentOutput)  {
	var total = 0
	var localUnSpents []LocalUnspentOutput

	for _, ws := range Walletinfo.Ws {
		
		if amount < total {
			break
		}

		a, outs := UTXO.FindSpendableOutputsByPubKeyHash(util.Sha256(ws.Wallet.PublicKey), amount - total)
		total += a
		for _, uo := range outs {
			uo := LocalUnspentOutput{uo.TxID, uo.Output, uo.Amount, ws.Wallet}
			localUnSpents = append(localUnSpents, uo)
		} 
	}
	return total, localUnSpents
}