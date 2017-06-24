package histories

type Synchronizer interface {
	SyncRecent() (int, error)
	SyncAll() (int, error)
}
