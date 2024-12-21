package policy

// EvictionPolicy defines the interface for cache eviction algorithms
type EvictionPolicy[K comparable, V any] interface {
	// RecordAccess updates access patterns for a key-value pair
	RecordAccess(key K, value V)
	// Evict removes and returns the entry chosen by the policy
	Evict() (K, V, bool)
	// Remove explicitly removes an entry
	Remove(key K) (V, bool)
	// Add adds a new entry to the policy tracker
	Add(key K, value V)
}
