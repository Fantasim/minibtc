package script

import (
	"errors"
	"bytes"
	"letsgo/util"
)

type opcode struct {
	value  byte
	name   string
	length int
	opfunc func(*parsedOpcode, *Engine) error
}

// potential data associated with it.
type parsedOpcode struct {
	opcode *opcode
	data   []byte
}

type SigHashType byte

// Hash type bits from the end of a signature.
const (
	SigHashOld          SigHashType = 0x0
	SigHashAll          SigHashType = 0x1
	SigHashNone         SigHashType = 0x2
	SigHashSingle       SigHashType = 0x3
	SigHashAllValue     SigHashType = 0x4
	SigHashAnyOneCanPay SigHashType = 0x80

	// sigHashMask defines the number of bits of the hash type which is used
	// to identify which outputs are signed.
	sigHashMask = 0x1f
)


const (
	OP_0                   = 0x00 // 0
	OP_DATA_1              = 0x01 // 1
	OP_DATA_2              = 0x02 // 2
	OP_DATA_3              = 0x03 // 3
	OP_DATA_4              = 0x04 // 4
	OP_DATA_5              = 0x05 // 5
	OP_DATA_6              = 0x06 // 6
	OP_DATA_7              = 0x07 // 7
	OP_DATA_8              = 0x08 // 8
	OP_DATA_9              = 0x09 // 9
	OP_DATA_10             = 0x0a // 10
	OP_DATA_11             = 0x0b // 11
	OP_DATA_12             = 0x0c // 12
	OP_DATA_13             = 0x0d // 13
	OP_DATA_14             = 0x0e // 14
	OP_DATA_15             = 0x0f // 15
	OP_DATA_16             = 0x10 // 16
	OP_2                   = 0x52 // 82
	OP_3                   = 0x53 // 83
	OP_DUP                 = 0x76 // 118 
	OP_EQUALVERIFY         = 0x88 // 136
	OP_ADD                 = 0x93 // 147
	OP_SUB                 = 0x94 // 148
	OP_HASH160             = 0xa9 // 169
	OP_CHECKSIG            = 0xac // 172
	OP_CHECKSIGVERIFY      = 0xad // 173
	OP_CHECKMULTISIG       = 0xae // 174
)

var opcodeArray = [256]opcode{
	OP_DATA_1:    {OP_DATA_1, "OP_DATA_1", 2, opcodePushData},
	OP_DATA_2:    {OP_DATA_2, "OP_DATA_2", 3, opcodePushData},
	OP_DATA_3:    {OP_DATA_3, "OP_DATA_3", 4, opcodePushData},
	OP_DATA_4:    {OP_DATA_4, "OP_DATA_4", 5, opcodePushData},
	OP_DATA_5:    {OP_DATA_5, "OP_DATA_5", 6, opcodePushData},
	OP_DATA_6:    {OP_DATA_6, "OP_DATA_6", 7, opcodePushData},
	OP_DATA_7:    {OP_DATA_7, "OP_DATA_7", 8, opcodePushData},
	OP_DATA_8:    {OP_DATA_8, "OP_DATA_8", 9, opcodePushData},
	OP_DATA_9:    {OP_DATA_9, "OP_DATA_9", 10, opcodePushData},
	OP_DATA_10:   {OP_DATA_10, "OP_DATA_10", 11, opcodePushData},
	OP_DATA_11:   {OP_DATA_11, "OP_DATA_11", 12, opcodePushData},
	OP_DATA_12:   {OP_DATA_12, "OP_DATA_12", 13, opcodePushData},
	OP_DATA_13:   {OP_DATA_13, "OP_DATA_13", 14, opcodePushData},
	OP_DATA_14:   {OP_DATA_14, "OP_DATA_14", 15, opcodePushData},
	OP_DATA_15:   {OP_DATA_15, "OP_DATA_15", 16, opcodePushData},
	OP_DATA_16:   {OP_DATA_16, "OP_DATA_16", 17, opcodePushData},
	OP_2:         {OP_2, "OP_2", 1, opcodeN},
	OP_3:         {OP_3, "OP_3", 1, opcodeN},
	OP_DUP:          {OP_DUP, "OP_DUP", 1, opcodeDup},
	OP_EQUALVERIFY: {OP_EQUALVERIFY, "OP_EQUALVERIFY", 1, opcodeEqualVerify},
	OP_ADD:                {OP_ADD, "OP_ADD", 1, opcodeAdd},
	OP_SUB:                {OP_SUB, "OP_SUB", 1, opcodeSub},
	OP_HASH160:             {OP_HASH160, "OP_HASH160", 1, opcodeHash160},
	OP_CHECKSIG:            {OP_CHECKSIG, "OP_CHECKSIG", 1, opcodeCheckSig},
	OP_CHECKSIGVERIFY:      {OP_CHECKSIGVERIFY, "OP_CHECKSIGVERIFY", 1, nil},
	OP_CHECKMULTISIG:       {OP_CHECKMULTISIG, "OP_CHECKMULTISIG", 1, nil},
}

func (op opcode) IsEmpty() bool {
	return op.name == ""
}

//Si l'opcode est une action a effectué sur la stack
//retourne true
func (code *parsedOpcode) IsAction() bool {
	if code.opcode.IsEmpty() == true {
		return false
	}

	switch uint(code.opcode.value) {
		case OP_DUP:
			return true
		case OP_EQUALVERIFY:
			return true
		case OP_ADD:
			return true
		case OP_SUB:
			return true
		case OP_HASH160:
			return true
		case OP_CHECKSIG:
			return true
		case OP_CHECKSIGVERIFY:
			return true
		case OP_CHECKMULTISIG:
			return true
		default:
			return false
	}
}


func opcodePushData(op *parsedOpcode, vm *Engine) error {
	vm.dstack.Push(op.data)
	return nil
}

// opcodeN is a common handler for the small integer data push opcodes.  It
// pushes the numeric value the opcode represents (which will be from 1 to 16)
// onto the data stack.
func opcodeN(op *parsedOpcode, vm *Engine) error {
	// The opcodes are all defined consecutively, so the numeric value is
	// the difference.
	vm.dstack.Push(op.data)
	return nil
}

// opcodeDup duplicates the top item on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x3]
func opcodeDup(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DupN(1)
}

// opcodeEqual removes the top 2 items of the data stack, compares them as raw
// bytes, and pushes the result, encoded as a boolean, back to the stack.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeEqual(op *parsedOpcode, engine *Engine) error {
	a, err := engine.dstack.Pop()
	if err != nil {
		return err
	}
	b, err := engine.dstack.Pop()
	if err != nil {
		return err
	}
	engine.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

// opcodeEqualVerify is a combination of opcodeEqual and opcodeVerify.
// Specifically, it removes the top 2 items of the data stack, compares them,
// and pushes the result, encoded as a boolean, back to the stack.  Then, it
// examines the top item on the data stack as a boolean value and verifies it
// evaluates to true.  An error is returned if it does not.
//
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeEqualVerify(op *parsedOpcode, engine *Engine) error {
	err := opcodeEqual(op, engine)
	if err == nil {
		verified, err := engine.dstack.PopBool()
		if err != nil {
			return err
		}
		if !verified {
			return errors.New("error from opcodeEqualVerify. Failed " + op.opcode.name)
		}
	}
	return err
}

// opcodeAdd treats the top two items on the data stack as integers and replaces
// them with their sum.
//
// Stack transformation: [... x1 x2] -> [... x1+x2]
func opcodeAdd(op *parsedOpcode, engine *Engine) error {
	v0, err := engine.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := engine.dstack.PopInt()
	if err != nil {
		return err
	}
	engine.dstack.PushInt(v0 + v1)
	return nil
}

// opcodeSub treats the top two items on the data stack as integers and replaces
// them with the result of subtracting the top entry from the second-to-top
// entry.
//
// Stack transformation: [... x1 x2] -> [... x1-x2]
func opcodeSub(op *parsedOpcode, engine *Engine) error {
	v0, err := engine.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := engine.dstack.PopInt()
	if err != nil {
		return err
	}

	engine.dstack.PushInt(v1 - v0)
	return nil
}

// opcodeHash160 treats the top item of the data stack as raw bytes and replaces
// it with ripemd160(sha256(data)).
//
// Stack transformation: [... x1] -> [... ripemd160(sha256(x1))]
func opcodeHash160(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.Pop()
	if err != nil {
		return err
	}

	hash := util.Ripemd160(util.Sha256(buf))
	vm.dstack.Push(hash)
	return nil
}

// Stack transformation: [... signature pubkey] -> [... bool]
func opcodeCheckSig(op *parsedOpcode, vm *Engine) error {
	/*pkBytes, err := vm.dstack.Pop()
	
	if err != nil {
		return err
	}

	fullSigBytes, err := vm.dstack.Pop()
	if err != nil {
		return err
	}

	// The signature actually needs needs to be longer than this, but at
	// least 1 byte is needed for the hash type below.  The full length is
	// checked depending on the script flags and upon parsing the signature.
	if len(fullSigBytes) < 1 {
		vm.dstack.PushBool(false)
		return nil
	}*/
	return nil
}