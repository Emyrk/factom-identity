package main

import (
	"flag"
	_ "flag"
	"fmt"

	"encoding/json"

	"bytes"

	"github.com/Emyrk/factom-identity"
	"github.com/FactomProject/factomd/common/primitives"
)

func main() {
	var (
		rootHex = flag.String("id", "", "Root Chain ID starting with '888888...'")
		factomd = flag.String("s", "localhost:8088", "Factomd api location")
		pretty  = flag.Bool("p", false, "Make the printout pretty for us mere humans")
	)

	flag.Parse()

	if *rootHex == "" {
		fmt.Println("factom-identity-cli -id=888888....")
		return
	}

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

	data, err := json.Marshal(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *pretty {
		var dst bytes.Buffer
		err := json.Indent(&dst, data, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(dst.Bytes()))
	} else {
		fmt.Println(string(data))
	}
}
