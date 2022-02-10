package attack

import "github.com/cbot/targets"

type Attack interface {


	Run(target targets.Target)

	PubProcess(process *AttackProcess)

}
