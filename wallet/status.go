package wallet

import (
	b "tway/blockchain"
	conf "tway/config"
	"tway/util"
)

type WalletInfo struct {
	Ws     []WalletStatus
	Amount int
}

//Structure représentant les informations basique d'une adresse
type WalletStatus struct {
	Address                []byte
	Amount                 int
	AmountLockedByMultiSig int
	W                      *Wallet
}

//Retourne une structure WalletInfo
//permettant d'obtenir les informations concernant
//les wallets enregistrés localement.
//Les informations sont le montant de coins disponible
// pour chaque adresse
func GetWalletInfo() *WalletInfo {
	utxo := b.UTXO

	wInfo := &WalletInfo{}

	//pour chaque wallet
	for _, w := range WalletList {
		//on récupère le montant disponible pour le wallet
		amount, list := utxo.GetUnspentOutputsByPubKOrPubKH([][]byte{HashPubKey(w.PublicKey), w.PublicKey}, conf.MAX_COIN)

		var amountLocked int
		for _, us := range list {
			if us.MultiSig == true {
				amountLocked += util.DecodeInt(us.Output.Value)
			}
		}
		amount -= amountLocked
		ws := WalletStatus{w.GetAddress(), amount, amountLocked, w}
		wInfo.Ws = append(wInfo.Ws, ws)

		wInfo.Amount += amount
	}
	return wInfo
}
