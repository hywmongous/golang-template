package snapshots

import "github.com/hywmongous/example-service/internal/domain/identity/aggregate"

type Snapshot struct {
	State aggregate.Identity
}

func RecreateSnapshot(
	state aggregate.Identity,
) Snapshot {
	return Snapshot{
		State: state,
	}
}
