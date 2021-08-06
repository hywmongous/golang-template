package snapshots

import identity "github.com/hywmongous/example-service/internal/domain/identity/aggregate"

type Snapshot struct {
	State identity.Identity
}

func RecreateSnapshot(
	state identity.Identity,
) Snapshot {
	return Snapshot{
		State: state,
	}
}
