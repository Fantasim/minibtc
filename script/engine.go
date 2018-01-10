package script

import (
	"fmt"
	"letsgo/util"
	"errors"
)


// Engine is the virtual machine that executes scripts.
type Engine struct {
	scripts         [][]parsedOpcode

	dstack          stack // data stack
	astack          stack // alt stack
}

func (engine *Engine) PrintScript(idx int){
	fmt.Printf("script[%d]: ", idx)
	for _, code := range engine.scripts[idx] {
		fmt.Print(" ", code.opcode.name)
	}
	fmt.Print("\n")
}

//Genère un pointeur vers une structure engine
func NewEngine() *Engine {
	engine := new(Engine)
	engine.scripts = make([][]parsedOpcode, 1)
	return engine	
}

//Parse le double array byte en script
func (engine *Engine) ParseScript(script [][]byte) error {
	for _, opcodeByte := range script {
		idx, err := util.ArrayByteToInt(opcodeByte)
		if err != nil {
			return err
		}
		var op *opcode
		if idx > 0 && idx < 256 {
			op = &opcodeArray[idx]
		}

		engine.scripts[0] = append(
			engine.scripts[0],
			parsedOpcode{
				op, 
				opcodeByte,
			},
		)
	}
	return nil
}

//Si le script a l'index est vide
func (engine *Engine) IsScriptEmpty(index int) bool {
	return len(engine.scripts[index]) == 0
}

//Demarre la lecture du script
func (engine *Engine) Run(initialScript [][]byte) error {
	//parsing du script
	err := engine.ParseScript(initialScript)
	if err != nil {
		return err
	}
	//si le script est vide
	if engine.IsScriptEmpty(0) == true {
		return errors.New("empty")
	}
	//affichage du script
	engine.PrintScript(0)
	var i = 0
	for i < len(engine.scripts[0]) {
		//Pour chaque ordre du script, effectue la function correspondante
		//push une valeur a la stake || effectue une action
		err := engine.scripts[0][i].opcode.opfunc(&engine.scripts[0][i], engine)
		if err != nil {
			return err
		}
		var newLineScript []parsedOpcode
		var j = i + 1
		//on créer une copie de la ligne de script
		//en ajoutant chaque ordre n'ayant pas encore 
		//été push sur la stack ou ayant été exécuté sur la stack
		for j < len(engine.scripts[0]) {
			newLineScript = append(newLineScript, engine.scripts[0][j])
			j++
		}
		//on ajoute la copie dans la tableau
		engine.scripts = append(engine.scripts, newLineScript)
		//on affiche la stack
		engine.dstack.PrintStack()
		//on affiche la nouvelle copie du script
		engine.PrintScript(i + 1)
		i++
	}
	return nil
}