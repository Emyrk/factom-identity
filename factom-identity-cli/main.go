package main

import (
	"flag"
	_ "flag"
	"fmt"

	"encoding/json"

	"bytes"

	"github.com/Emyrk/factom-identity"
	"github.com/FactomProject/factomd/common/identity"
	"github.com/FactomProject/factomd/common/primitives"
)

func main() {
	var (
		all     = flag.Bool("all", false, "Parse for all identities")
		rootHex = flag.String("id", "", "Root Chain ID starting with '888888...'")
		factomd = flag.String("s", "localhost:8088", "Factomd api location")
		pretty  = flag.Bool("p", false, "Make the printout pretty for us mere humans")
	)

	flag.Parse()

	var data []byte
	c := factom_identity.NewAPIController(*factomd)
	if !c.IsWorking() {
		fmt.Println("Factomd location is not working")
		return
	}

	if *all {
		ids, err := parseAll(c)
		if err != nil {
			fmt.Println(err)
			return
		}

		data, err = json.Marshal(ids)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// Single
		if *rootHex == "" {
			fmt.Println("factom-identity-cli -id=888888....")
			return
		}

		id, err := parseSingle(*rootHex, c)
		if err != nil {
			fmt.Println(err)
			return
		}

		data, err = json.Marshal(id)
		if err != nil {
			fmt.Println(err)
			return
		}
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

func parseSingle(rootHex string, c *factom_identity.Controller) (*identity.Identity, error) {
	root, err := primitives.HexToHash(rootHex)
	if err != nil {
		return nil, err
	}

	id, err := c.FindIdentity(root)
	return id, err
}

func parseAll(c *factom_identity.Controller) (map[string]*identity.Identity, error) {
	ids, err := c.FindAllIdentities()
	return ids, err
}
