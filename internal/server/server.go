package server

import (
	"log"
	"net"

	"github.com/eager/cyberpunkmp/internal/network" // (0) Import de la couche réseau
)

// (2) Fonction démarrant le serveur TCP
func Start() error {
	// (2.a) Ecoute sur toutes les interfaces sur le port 11778
	listener, err := net.Listen("tcp", "0.0.0.0:11778")
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Println("Server listening on :11778")

	// (2.b) Boucle infinie qui accepte les connexions entrantes
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		// (2.c) Pour chaque connexion, lance une goroutine pour gérer l'échange
		go network.HandleConnection(conn)
	}
}
