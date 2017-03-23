package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
)

func Container(container dbng.CreatedContainer, meta dbng.ContainerMetadata) atc.Container {
	return atc.Container{
		ID:         container.Handle(),
		WorkerName: container.WorkerName(),

		Type: string(meta.Type),

		PipelineID:     meta.PipelineID,
		JobID:          meta.JobID,
		BuildID:        meta.BuildID,
		ResourceID:     meta.ResourceID,
		ResourceTypeID: meta.ResourceTypeID,

		StepName: meta.StepName,
		Attempt:  meta.Attempt,

		WorkingDirectory: meta.WorkingDirectory,
		User:             meta.User,
	}
}
