package main

import (
	"github.com/xxRanger/bridge/backEnd/manager"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	PRIVATE_CHAIN_CONFIG_FILE = "etc/privateChain.json"
	PUBLIC_CHAIN_CONFIG_FILE  = "etc/publicChain.json"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)


	// log to file
	//
	//f,err:= os.OpenFile(LOG_FILE,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	//if err!=nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//log.SetOutput(f)
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main () {
	r := mux.NewRouter()

	log.Println("start manager")
	m:= manager.NewManager()

	log.Println("set manager socket handler")
	m.InitHandlers()

	log.Println("set manager chain handler")
	pbcChainHandler,err:= manager.NewPublicChainHandler(PUBLIC_CHAIN_CONFIG_FILE)
	if err!=nil {
		panic(err)
	}

	pvcChainHandler,err:=manager.NewPrivateChainHandler(PRIVATE_CHAIN_CONFIG_FILE)
	if err!=nil {
		panic(err)
	}

	chainHandler:= manager.NewChainHandler(pvcChainHandler,pbcChainHandler)
	m.SetChainHandler(chainHandler)

	go m.Start()
	go m.Subscribe()

	log.Println("start server")
	r.HandleFunc("/",m.WebsocketHandler).Methods("GET")
	r.HandleFunc("/{chain:public|private}/{user}/nonce",m.GetNonce).Methods("GET")
	r.HandleFunc("/{chain:public|private}/chainId",m.GetChainId).Methods("GET")
	http.ListenAndServe(
		"0.0.0.0:4000", handlers.CORS(
			handlers.AllowedMethods([]string{"get", "options", "post", "put", "head"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
	)
}