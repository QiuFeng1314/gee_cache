module httppool

go 1.20

require (
	cache v1.0.0
	lru v1.0.0 // indirect
)

replace (
	cache v1.0.0 => ../cache
	lru v1.0.0 => ../lru
)