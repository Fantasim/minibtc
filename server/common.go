package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"
	conf "tway/config"
	"tway/serverutil"
	"tway/util"
)

func commandToBytes(command string) []byte {
	var bytes [conf.CommandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

//retourne un []byte representant la cmd en string
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:conf.CommandLength]
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//Cette fonction recupere le contenu d'une requete envoyé par un noeud selon la requete
func getPayload(request []byte, payload interface{}) error {
	var buff bytes.Buffer
	//ecrit dans le buffeur le payload de la request de commandLength bit jusqu'à la fin
	buff.Write(request[conf.CommandLength:])
	//[]byte en structure verzion
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(payload)
	if err != nil {
		return err
	}
	return nil
}

//Cette fonction recupere l'adresse locale du noeud
func GetLocalNetAddr() *serverutil.NetAddress {
	ip, err := util.GetIP()
	if err != nil {
		log.Panic(err)
	}
	port, _ := strconv.Atoi(conf.NODE_ID)
	addrMe := serverutil.NewNetAddressIPPort(ip, uint16(port))
	return addrMe
}

//Cette fonction recupère l'adresse du noeud principale
func GetMainNode() *serverutil.NetAddress {
	return serverutil.NewNetAddressIPPort(conf.MainNodeIP, conf.MainNodePort)
}

func (s *Server) LocalWaiting() {
	if s.prod == false {
		time.Sleep(time.Millisecond * 5)
	}
}
