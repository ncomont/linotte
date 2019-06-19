package runner

import proto_job "git.ve.home/nicolasc/linotte/services/job/proto"

// Runner is the generic interface for every runners
type Runner interface {
	Run(job *proto_job.JobReply)
	Configure(channels *Channels) Runner
}

// Channels contains the channels used to communicate with the runners
type Channels struct {
	Progression chan uint32
	Results     chan []*proto_job.ResultReply
	Errors      chan error
}

// Get returns the appropriate runner for the given job type
func Get(t string) Runner {
	switch t {
	case "RL":
		return NewRedListRunner()
	}
	return nil
}
