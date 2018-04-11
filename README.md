# Tool to parse identities from the blockchain

Running the command line tool uses factomd api to build an identity from the blockchain. the identity will be returned in json

## Usage

`factom-identity-cli -s localhost:8088 -id=8888888000000000000000000000000000000000000000000000000000000000`

Optional `-p` flag to _pretty print_ the json response.