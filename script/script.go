package script

import (
	"tway/util"
	"fmt"
	"encoding/hex"
)

var Script = new(script)

type script struct {}


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

func TestScript(s func() [][]byte){
	scrpt := s()
	engine := NewEngine()
	err := engine.Run(scrpt)
	if err == nil {
		fmt.Println("Script correct")
	} else {
		fmt.Println(err)
		fmt.Println("Script incorrect")
	}
}
