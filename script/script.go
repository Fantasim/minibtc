package script

import (
	"tway/util"
	"encoding/hex"
	"tway/config"
	"errors"
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

func (s *script) IsPayToPubKeyHash(scriptBytes [][]byte)bool{
	engine := new(Engine)
	engine.scripts = make([][]parsedOpcode, 1)
	engine.ParseScript(scriptBytes)

	if len(scriptBytes) == config.P2PKHSize {
		sigSize := len(scriptBytes[0]) == 64 
		pubKeySize := len(scriptBytes[1]) == 64
		opDup := engine.scripts[0][2].opcode.value == OP_DUP
		opHash160 := engine.scripts[0][3].opcode.value == OP_HASH160
		pubKeyHashSize := len(scriptBytes[4]) == 20
		opEqualVerify := engine.scripts[0][5].opcode.value == OP_EQUALVERIFY
		opCheckSig := engine.scripts[0][6].opcode.value == OP_CHECKSIG
		
		if sigSize && pubKeySize && pubKeyHashSize && opDup && opHash160 && opEqualVerify && opCheckSig {
			return true
		}
	}
	return false
}

func (s *script) GetPubKeyHash(pubKeyScript [][]byte) ([]byte, error) {
	if len(pubKeyScript) == 5 && len(pubKeyScript[2]) != 20 {
		return []byte{}, errors.New("not a P2PKH script") 
	}
	return pubKeyScript[2], nil
} 




