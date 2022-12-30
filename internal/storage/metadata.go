// metadata is a key-value pair table that contains assortments of data, like system parameters or system states
// todo: could have been mongodb

package storage

import "context"

type Metadata interface {
	// CanUsersBeMade controls whether the CreateUser endpoint can be used
	CanUsersBeMade(ctx context.Context) (bool, error)
	// SetCanUsersBeMade sets the current value for CanUsersBeMade
	SetCanUsersBeMade(ctx context.Context, canBeMade bool) (bool, error)
}
