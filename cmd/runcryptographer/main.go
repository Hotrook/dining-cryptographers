package main

import (
	"flag"
	"github.com/Hotrook/dining_cryptographers/cryptographer"
	"github.com/Hotrook/dining_cryptographers/logutils"
)

var(
	id = flag.Int("id", 1, "Client number")
	payed = flag.Bool("payed", false, "Did client pay?")
)

func main(){

	flag.Parse()
	logutils.InitLogger()

	cryptographer := cryptographer.Cryptographer{Id: *id, Payed: *payed }
	cryptographer.Run()

}


