package jenkins


import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"time"
)

type JenkinsClient struct {

	jclient *gojenkins.Jenkins
	buildPollTime uint64
	ctx context.Context
}

func NewJenkinsClient(host string,port int,proto,user,passwd string,buildPollTime uint64) (*JenkinsClient,error) {


	jclient := gojenkins.CreateJenkins(nil,fmt.Sprintf("%s://%s:%d/",proto,host,port),user,passwd)
	_, err := jclient.Init(context.Background())

	if err!=nil {

		return nil,err
	}

	return &JenkinsClient{jclient:jclient,
		buildPollTime:buildPollTime,
		ctx:context.Background(),
	},nil

}


func (j *JenkinsClient) CreateJob(name,config string)(string,error) {

	job,err := j.jclient.CreateJob(j.ctx,config,name)

	if err!=nil {
		return "",err
	}

	return job.Base,err
}

// A task in queue will be assigned a build number in a job after a few seconds.
// this function will return the build object.
func  (jc *JenkinsClient)GetBuildFromQueueID(queueid int64) (*gojenkins.Build, error) {

	j := jc.jclient
	ctx := jc.ctx

	task, err := j.GetQueueItem(ctx, queueid)
	if err != nil {
		return nil, err
	}
	// Jenkins queue API has about 4.7second quiet period
	for task.Raw.Executable.Number == 0 {
		time.Sleep(time.Duration(jc.buildPollTime) * time.Millisecond)
		_, err = task.Poll(ctx)
		if err != nil {
			return nil, err
		}
	}

	buildid := task.Raw.Executable.Number
	job, err := task.GetJob(ctx)
	if err != nil {
		return nil, err
	}
	build, err := job.GetBuild(ctx, buildid)
	if err != nil {
		return nil, err
	}
	return build, nil
}

func (j *JenkinsClient) BuildJobWaitResult(name string)(string,error) {

	qid,err :=j.jclient.BuildJob(context.Background(),name,map[string]string{})

	if err!=nil {
		return "",err
	}

	build, err := j.GetBuildFromQueueID(qid)
	if err != nil {

		return "",err
	}

	// Wait for build to finish
	for build.IsRunning(j.ctx) {
		time.Sleep(time.Duration(j.buildPollTime) * time.Millisecond)
		build.Poll(j.ctx)
	}

	return build.GetConsoleOutput(j.ctx),nil
}

func (j *JenkinsClient) BuildJob(name string) (int64, error) {

	return j.jclient.BuildJob(context.Background(),name,map[string]string{})

}

func (j *JenkinsClient) DeleteJob(name string) {

	j.jclient.DeleteJob(j.ctx,name)

}




