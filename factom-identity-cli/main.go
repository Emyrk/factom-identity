package main

import (
	"flag"
	_ "flag"
	"fmt"

	"github.com/Emyrk/factom-identity"
	"github.com/FactomProject/factomd/common/primitives"
)

func main() {
	var (
		rootHex = flag.String("id", "", "Root Chain ID starting with '888888...'")
		factomd = flag.String("s", "localhost:8088", "Factomd api location")
	)

	flag.Parse()

	root, err := primitives.HexToHash(*rootHex)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := factom_identity.NewAPIController(*factomd)
	if !c.IsWorking() {
		fmt.Println("Factomd location is not working")
		return
	}

	id, err := c.FindIdentity(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	j, _ := id.JSONString()
	fmt.Println(j)
}
