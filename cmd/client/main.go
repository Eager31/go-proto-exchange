package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	cpProto "github.com/eager/cyberpunkmp/proto" // (0) Import package protobuf généré
	"google.golang.org/protobuf/proto"           // (0) Librairie officielle protobuf
)

func main() {
	// (1) Connexion TCP au serveur
	conn, err := net.Dial("tcp", "localhost:11778")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// (3) Lire le handshake envoyé par le serveur dès la connexion
	serverMsg, err := readServerMessage(conn)
	if err != nil {
		log.Fatalf("Failed to read server message: %v", err)
	}
	fmt.Printf("Server handshake: PlayerId=%s PlayerName=%s\n", serverMsg.PlayerId, serverMsg.PlayerName)

	// (4) Envoyer son propre handshake au serveur
	clientHandshake := &cpProto.Handshake{
		PlayerId:   "client123",
		PlayerName: "GoClient",
	}
	if err := sendClientMessage(conn, &cpProto.ClientMessage{Payload: &cpProto.ClientMessage_Handshake{Handshake: clientHandshake}}); err != nil {
		log.Fatalf("Failed to send handshake: %v", err)
	}

	fmt.Println("Handshake envoyé, client prêt.")
}

// (3.a) Fonction qui lit un message handshake venant du serveur
func readServerMessage(conn net.Conn) (*cpProto.Handshake, error) {
	lengthBuf := make([]byte, 4)
	// Lire d'abord la taille (4 octets)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// Lire le message complet selon la taille
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	// Désérialiser les données en message protobuf Handshake
	msg := &cpProto.Handshake{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// (4.a) Fonction pour envoyer un message ClientMessage au serveur
func sendClientMessage(conn net.Conn, msg *cpProto.ClientMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Envoi taille puis message
	length := uint32(len(data))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, length)

	if _, err := conn.Write(lengthBuf); err != nil {
		return err
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}
	return nil
}
