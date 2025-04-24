package types

import (
	"fmt"
	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

func printChainState(c *Chain) {
	f := colorjson.NewFormatter()
	f.Indent = 4

	state := map[string]interface{}{
		"chain_id":       c.ChainID,
		"started":        c.ChainStarted,
		"block_count":    len(c.Blocks),
		"pending_tx":     len(c.PendingTransactions),
		"last_block":     getLastBlockInfo(c),
		"pending_tx_ids": getPendingTxIds(c),
	}

	s, err := f.Marshal(state)
	if err != nil {
		color.Red("Error formatting state: %v", err)
		return
	}

	color.Cyan("\n=== Blockchain State ===")
	fmt.Println(string(s))
	color.Cyan("=======================\n")
}

func getLastBlockInfo(c *Chain) map[string]interface{} {
	if len(c.Blocks) == 0 {
		return nil
	}
	last := c.Blocks[len(c.Blocks)-1]
	return map[string]interface{}{
		"hash":      last.BlockHash,
		"tx_count":  len(last.Txs),
		"timestamp": last.Timestamp,
	}
}

func getPendingTxIds(c *Chain) []string {
	ids := make([]string, 0, len(c.PendingTransactions))
	for _, tx := range c.PendingTransactions {
		ids = append(ids, tx.Id[:8]+"...")
	}
	return ids
}
