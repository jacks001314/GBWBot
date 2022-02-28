package node

func (n *Node) waitAttackProcess() {

	attackProcessChan := n.attackTasks.SubAttackProcess()

	go func() {
		for {

			select {

			case ap := <-attackProcessChan:

				//send to sbot
				n.nodeClient.SendAttackProcess(ap)
			}
		}
	}()
}
