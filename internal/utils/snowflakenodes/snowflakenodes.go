package snowflakenodes

import (
	"github.com/bwmarrin/snowflake"
)

var (
	// NodesReport contains snowflake nodes for
	// transactions.
	NodesReport []*snowflake.Node

	NodeTransaction *snowflake.Node
	// NodeGuildLog is the snowflake node
	// for guild logs.
	NodeGuildLog *snowflake.Node

	// nodeMap maps snowflake node IDs with
	// their identifier strings.
	nodeMap map[int]string
)

// Setup initializes the snowflake nodes and
// nodesMap.
func Setup() (err error) {
	nodeMap = make(map[int]string)

	NodeTransaction, _ = RegisterNode(0, "transaction")
	NodeGuildLog, _ = RegisterNode(10, "karmarules")

	return
}

func RegisterNode(id int, name string) (*snowflake.Node, error) {
	nodeMap[id] = name
	return snowflake.NewNode(int64(id))
}

// GetNodeName returns the identifier name of
// the snowflake node by nodeID.
func GetNodeName(nodeID int64) string {
	if ident, ok := nodeMap[int(nodeID)]; ok {
		return ident
	}
	return "undefined"
}
