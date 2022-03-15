// Example of a butter dapp (decentralised application) where data is persistent: wiki. The basic functionality of the
// wiki is to be able to add an entry and read an entry.
package main

import (
	"bufio"
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"os"
	"strings"
)

func add(overlay *pcg.Peer) {
	fmt.Println("Input information:")
	in := bufio.NewReader(os.Stdin)
	data, _ := in.ReadString('\n') // Read string up to newline
	uuid := pcg.PCGStore(overlay, strings.TrimSpace(data))
	fmt.Println("UUID:", uuid)
	fmt.Println("Data:", strings.TrimSpace(data))
	fmt.Println("Enter to continue...")
	in.ReadString('\n')
}

func retrieve(overlay *pcg.Peer) {
	fmt.Println("Information UUID:")
	in := bufio.NewReader(os.Stdin)
	uuid, _ := in.ReadString('\n') // Read string up to newline
	data := pcg.NaiveRetrieve(overlay, strings.TrimSpace(uuid))
	fmt.Println(string(data))
	fmt.Println("Enter to continue...")
	in.ReadString('\n')
}

func printAll(peer *pcg.Peer) {
	fmt.Println(peer.String())
	fmt.Println("Enter to continue...")
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
}

func interact(overlayInterface node.Overlay) {
	peer := overlayInterface.(*pcg.Peer)
	for {
		// prompt to pcgStore or pcgRetrieve information
		var interactionType string
		fmt.Print("add(1), retrieve(2) or see all my groups(3) information?")
		fmt.Scanln(&interactionType)
		switch interactionType {
		case "1":
			add(peer)
		case "2":
			retrieve(peer)
		case "3":
			printAll(peer)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func main() {
	// Create a new node by specifying a port (or setting it to 0 to let the OS assign one) and defining an upper limit
	// on memory usage
	butterNode, _ := node.NewNode(0, 512)

	// PCG overlay network
	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	// App-level client behaviour (i.e. how are the users expected to interface with the dapp?)
	butterNode.RegisterClientBehaviour(interact)

	fmt.Println("Node is listening at", butterNode.Address())

	// Spawn your node into the butter network
	butter.Spawn(&overlay, false) // Blocking
}
