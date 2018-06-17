package main

import (
	"flag"
	"github.com/Hotrook/dining_cryptographers/cryptographer"
)

var(
	id = flag.Int("id", 1, "Client number")
	payed = flag.Bool("payed", false, "Did client pay?")
)

func main(){

	flag.Parse()
	cryptographer := cryptographer.Cryptographer{Id: *id, Payed: *payed }
	cryptographer.Run()

}


