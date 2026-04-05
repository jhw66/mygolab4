package utils

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeMu   sync.RWMutex
	snowflakeNode *snowflake.Node
)

func InitSnowflake(nodeID int64) error {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}

	snowflakeMu.Lock()
	snowflakeNode = node
	snowflakeMu.Unlock()
	return nil
}

func GenerateID() (string, error) {
	snowflakeMu.RLock()
	node := snowflakeNode
	snowflakeMu.RUnlock()

	if node == nil {
		if err := InitSnowflake(1); err != nil {
			return "", err
		}
		snowflakeMu.RLock()
		node = snowflakeNode
		snowflakeMu.RUnlock()
	}

	return node.Generate().String(), nil
}
