package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mingi3442/mingicoin/blockchain"
	"github.com/mingi3442/mingicoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port :3000 will upgrade the request from :4000
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string, broadcast bool) {
	// Port :4000 is requesting an upgrade from the port :3000
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		BroadCastNewPeer(p)
		return
	}
	sendNewestBlock(p)
}

func BroadcastNewBlcok(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func BroadCastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
