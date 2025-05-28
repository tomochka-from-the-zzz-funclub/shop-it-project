package entity

type ClusterHealth struct {
	Status           string
	ActiveShards     int64
	RelocatingShards int64
	UnassignedShards int64
	TimedOut         bool
}
