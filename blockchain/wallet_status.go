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
		amount, _ := UTXO.GetUnspentOutputsByPubKeyHash(wallet.HashPubKey(w.PublicKey), MAX_COIN)
		ws := WalletStatus{w.GetAddress(), amount, w}
		wInfo.Ws = append(wInfo.Ws, ws)

		wInfo.Amount += amount
	}
	return wInfo
}
