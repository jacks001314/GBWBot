package node

import "time"

func (n *Node) Ping() {

	tchan := time.Tick(5 * time.Minute)

	go func() {

		for {

			select {
			case <-tchan:

				n.nodeClient.Ping()
			}
		}
	}()
}
