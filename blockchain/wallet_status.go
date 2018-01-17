package blockchain

import (
	"letsgo/wallet"
)


type WalletInfo struct {
	Ws []WalletStatus
	Amount int
}

//Structure représentant les informations basique d'une adresse
type WalletStatus struct {
	Address []byte
	Amount int
	Wallet *wallet.Wallet
}

//Structure représentant les informations liés à un UTXO
type UnspentOutput struct {
	TxID []byte
	Output int
	Amount int
}

//Structure représentant les informations liées 
//à un UTXO présent dans un wallet local
type LocalUnspentOutput struct {
	TxID []byte
	Output int
	Amount int
	Wallet *wallet.Wallet
}

//Retourne une structure WalletInfo
//permettant d'obtenir les informations concernant
//les wallets enregistrés localement.
//Les informations sont le montant de coins disponible
// pour chaque adresse
func GetWalletInfo() *WalletInfo {
	wInfo := &WalletInfo{}

	//pour chaque wallet
	for _, w := range wallet.WalletList {
		//on récupère le montant disponible pour le wallet
		amount, _ := UTXO.FindSpendableOutputsByPubKeyHash(wallet.HashPubKey(w.PublicKey), MAX_COIN)
		ws := WalletStatus{w.GetAddress(), amount, w}
		wInfo.Ws = append(wInfo.Ws, ws)

		wInfo.Amount += amount
	}
	return wInfo
}

//Récupère une liste UTXO sur des wallets 
//enregistrés localement.
func (wInfo *WalletInfo) GetLocalUnspentOutputs(amount int) (int, []LocalUnspentOutput)  {
	var total = 0
	var localUnSpents []LocalUnspentOutput

	for _, ws := range Walletinfo.Ws {
		
		if amount < total {
			break
		}

		a, outs := UTXO.FindSpendableOutputsByPubKeyHash(wallet.HashPubKey(ws.Wallet.PublicKey), amount - total)
		total += a
		for _, uo := range outs {
			uo := LocalUnspentOutput{uo.TxID, uo.Output, uo.Amount, ws.Wallet}
			localUnSpents = append(localUnSpents, uo)
		} 
	}
	return total, localUnSpents
}