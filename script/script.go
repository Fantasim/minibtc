package script

import (
	"tway/util"
	"encoding/hex"
)

var Script = new(script)

type script struct {}

func init(){
}

//Generation d'un script de type PayToPubKeyHash (ScriptPubKey)
//OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
func (s *script) LockingScript(pubKeyHash []byte) [][]byte {
	return util.DupByteDoubleArray(
		append([]byte{}, OP_DUP),
		append([]byte{}, OP_HASH160),
		[]byte(pubKeyHash), 
		append([]byte{}, OP_EQUALVERIFY),
		append([]byte{}, OP_CHECKSIG),
	)
}

//Generation d'un script de locking script pour une transaction coinbase (Output) (ScriptPubkey)
//prend en paramètre la clé publique de son wallet
func (s *script) CoinbaseLockingScript(pubKey []byte) [][]byte {
	return util.DupByteDoubleArray(
		[]byte(pubKey), 
		append([]byte{}, OP_CHECKSIG),
	)
}

//Generation d'un script d'input : //ScriptSig
//<signature> <pubKey>
func (s *script) UnlockingScript(signature, pubKey []byte) [][]byte{ 
	return util.DupByteDoubleArray(
		[]byte(signature), 
		[]byte(pubKey),
	)
}

//Generation d'un script d' unlocking script pour une transaction coinbase (ScriptSig)
//prend en paramètre la signature de sa clé privée
func (s *script) CoinbaseUnlockingScript(signature []byte) [][]byte {
	return [][]byte{signature}
}


//Addition 1 + 4 = 5
//Script correct
func (s *script) FiveEqualFive() [][]byte {
	return util.DupByteDoubleArray(
		append([]byte{}, OP_DATA_1),
		append([]byte{}, OP_DATA_4),
		append([]byte{}, OP_ADD),
		append([]byte{}, OP_DATA_5),
		append([]byte{}, OP_EQUALVERIFY),
	)
}

//Adition 1 + 3 = 5
//Script incorrect
func (s *script) FourEqualFive() [][]byte {
	return util.DupByteDoubleArray(
		append([]byte{}, OP_DATA_1),
		append([]byte{}, OP_DATA_3),
		append([]byte{}, OP_ADD),
		append([]byte{}, OP_DATA_5),
		append([]byte{}, OP_EQUALVERIFY),
	)
}

func (s *script) TxScript() [][]byte {
	sig, _ := hex.DecodeString("4e4a1458c0e5346edfb63b5b5e2d5c96fb40bf20b7c226d158abcd975bb1a2157e6c2a46ef60c956703e54552db546840be13ddaff97d065617a70422bcbf2e1")
	pubk, _ := hex.DecodeString("bdd01e59dbc5103a1972c01adc90af9d319f768c7dcd35b107d9b0022067069e29465d2fe8c55680946263977821a6c12d7b1dfcd58e8258f8986d7759f38158")
	pubkHash, _ := hex.DecodeString("c3cd8e22f0e4d5d8c51490ecdc548213f4e3086a")

	return util.DupByteDoubleArray(
		append([]byte{}, sig...),
		append([]byte{}, pubk...),
		append([]byte{}, OP_DUP),
		append([]byte{}, OP_HASH160),
		append([]byte{}, pubkHash...),
		append([]byte{}, OP_EQUALVERIFY),
		append([]byte{}, OP_CHECKSIG),
	)
}

func (s *script) String(srpt [][]byte) string {
	ret := ""
	for _, elem := range srpt {
		if len(elem) == 1 {
			ret += opcodeArray[int(elem[0])].name
		} else {
			ret += hex.EncodeToString(elem)
		}
		ret += " "
	}
	return ret
}
/*
func TestScript(s func() [][]byte){
	scrpt := s()
	engine := NewEngine(util.Transaction{}, -1)
	err := engine.Run(scrpt)
	if err == nil {
		fmt.Println("Script correct")
	} else {
		fmt.Println(err)
		fmt.Println("Script incorrect")
	}
}*/



