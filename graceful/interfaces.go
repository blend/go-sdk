package graceful

// Graceful is a server that can start and stop.
type Graceful interface {
	Start() error // this call must block
	Stop() error
}

// Updater is an optional interface that allows for updates.
type Updater interface {
	Update() error
}
