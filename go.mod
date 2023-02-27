module gee_cache

go 1.20

require (
	lru v1.0.0 // indirect
	cache v1.0.0
	httppool v1.0.0
)

replace (
	lru v1.0.0 => ./lru
	cache v1.0.0 => ./cache
	httppool v1.0.0 => ./httppool
)