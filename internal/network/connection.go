package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	cpProto "github.com/eager/cyberpunkmp/proto" // (0) Import protobuf généré
	"google.golang.org/protobuf/proto"           // (0) Librairie officielle protobuf
)

// (2.c) Fonction lancée par le serveur pour gérer une connexion client
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", addr)

	// (2.c.i) Dès la connexion, envoie un handshake serveur au client
	if err := sendHandshake(conn); err != nil {
		log.Printf("Failed to send handshake: %v", err)
		return
	}

	// (5) Boucle qui attend et traite les messages envoyés par le client
	for {
		msg, err := readClientMessage(conn)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", addr)
			} else {
				log.Printf("Failed to read client message: %v", err)
			}
			return
		}

		// (6) Traitement du message selon son type
		switch payload := msg.Payload.(type) {
		case *cpProto.ClientMessage_Handshake:
			log.Printf("Received handshake from player %s", payload.Handshake.PlayerName)
		case *cpProto.ClientMessage_PlayerUpdate:
			log.Printf("Player %s moved to %f, %f, %f", payload.PlayerUpdate.PlayerId, payload.PlayerUpdate.X, payload.PlayerUpdate.Y, payload.PlayerUpdate.Z)
		default:
			log.Printf("Unknown message")
		}
	}
}

// (2.c.i.a) Fonction qui envoie un handshake au client
func sendHandshake(conn net.Conn) error {
	msg := &cpProto.Handshake{
		PlayerId:   "player123",
		PlayerName: "Emilien",
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Envoi taille + message binaire
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

// (5.a) Fonction qui lit un message ClientMessage envoyé par le client
func readClientMessage(conn net.Conn) (*cpProto.ClientMessage, error) {
	lengthBuf := make([]byte, 4)
	// Lire taille (4 octets)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// Lire le message complet
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	// Désérialiser le message protobuf
	msg := &cpProto.ClientMessage{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
