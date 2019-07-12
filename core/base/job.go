package base

import (
	pb "core/jobrequest"
)

// IJob is an interface for utron Jobs
type IJob interface {
	New(*Context)
	Handle(req *pb.JobRequest) error
}

// Job implements the Controller interface, It is recommended all
// user defined Jobs should embed *BaseJob.
type Job struct {
	Ctx *Context
}

// New sets ctx as the active context
func (b *Job) New(ctx *Context) {
	b.Ctx = ctx
}

// Handle commits the changes made in the active context.
func (b *Job) Handle(req *pb.JobRequest) error {
	return nil
}
