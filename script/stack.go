package script

import (
	"errors"
	"tway/util"
)

type stack struct {
	stk [][]byte
}

//ajoute un []byte dans la stack
func (s *stack) Push(elem []byte) error {
	s.stk = append(s.stk, elem)
	return nil
}

//ajoute un int dans la stack
func (s *stack) PushInt(n int) error {
	s.stk = append(s.stk, []byte{util.IntToByte(n)})
	return nil
}

//ajoute un bool dans la stack
func (s *stack) PushBool(b bool) error {
	if b {
		s.stk = append(s.stk, []byte{0x01})
	} else {
		s.stk = append(s.stk, nil)
	}
	return nil
}

//recupere et supprime le dernier element ajouté dans la stack
//format : []byte
func (s *stack) Pop() ([]byte, error) {
	if len(s.stk) == 0 {
		return []byte{}, errors.New("empty stack")
	}
	var idx = len(s.stk) - 1
	var ret = s.stk[idx]
	s.stk = append(s.stk[:idx], s.stk[idx+1:]...)
	return ret, nil
}

//recupere et supprime le dernier element ajouté dans la stack
//format : int
func (s *stack) PopInt() (int, error) {
	if len(s.stk) == 0 {
		return 0, errors.New("empty stack")
	}
	var idx = len(s.stk) - 1
	var ret int
	var err error
	if len(s.stk[idx]) == 1 {
		ret = util.ByteToInt(s.stk[idx][0])
	} else {
		ret, err = util.ArrayByteToInt(s.stk[idx])
		if err != nil {
			return 0, err
		}
	}
	s.stk = append(s.stk[:idx], s.stk[idx+1:]...)
	return ret, nil
}

//recupere et supprime le dernier element ajouté dans la stack
//format : bool
func (s *stack) PopBool() (bool, error) {
	if len(s.stk) == 0 {
		return false, errors.New("empty stack")
	}
	var idx = len(s.stk) - 1
	var ret int
	var err error
	if len(s.stk[idx]) == 1 {
		ret = util.ByteToInt(s.stk[idx][0])
	} else {
		ret, err = util.ArrayByteToInt(s.stk[idx])
		if err != nil {
			return false, err
		}
	}
	s.stk = append(s.stk[:idx], s.stk[idx+1:]...)
	if ret == 1 {
		return true, nil
	}
	return false, nil
}

//Duplique les n derniers elements de la stack
func (s *stack) DupN(n int) error {
	for n > 0 {
		s.stk = append(s.stk, s.stk[len(s.stk)-1])
		n--
	}
	return nil
}
