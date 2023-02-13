# Go Cache

[![codecov](https://codecov.io/gh/damianopetrungaro/go-cache/branch/main/graph/badge.svg?token=5ESXFZo2j2)](https://codecov.io/gh/damianopetrungaro/go-cache)

GoCache is an opinionated cache package
with simple APIs and provides multiple implementations.

## Why another cache package?

GoCache is designed with Go 1.18 and leverage the generics implementation,
reducing the amount of casting needed to provide an efficient cache implementation.

On top of that, GoCache also requires a context argument allowing context propagation,
which is not available in most of the cache package available today.

## Examples

Usage 
```go
var c Cache[string, string]
var ctx := context.Background()

const k = "key"
const val = "value"

// Adds a never expiring item to the cache
if err := c.Set(ctx, k, val, NoExpiration); err != nil {
  return err
}

// Adds an item which will expire after 10 seconds
if err := c.Set(ctx, k, val, 10*time.Second); err != nil {
  return err
}

// Retrieves an item from the cache
item, err := c.Get(ctx, k)
if err != nil {
  return err
}

// Deletes an item from the cache
if err := c.Delete(ctx, k); err != nil {
  return err
}
```

GoCache provides also some errors
```go
ErrNotSet    = errors.New("could not set cache value")
ErrNotGet    = errors.New("could not get cache value")
ErrNotFound  = fmt.Errorf("%w: could not find cache value", ErrNotGet)
ErrExpired   = fmt.Errorf("%w: could not get expired cache value", ErrNotGet)
ErrNotDelete = errors.New("could not delete cache value")
```

### In Memory

Create an InMemory implementation

```go
// the first generic is a comparable type used as key
// the second generic is any type used as value
// the first argument of the factory function represent the max item capacity of the cache
// when the max capacity gets hit, then all the expired items get deleted and if none is expired 
// then the one closest to the expiry get deleted 
inmem := NewInMemory[string, int](100_000)
```

### Redis

```go
import (
    goRedis "github.com/go-redis/redis/v9"
)

var redisClient *goRedis.Client

// the first generic is a comparable type used as key
// the second generic is any type used as value
// the first argument of the factory function is the client used to communicate with the redis server
redisCache := redis.New[string, int](redisClient)
```

```go
import (
    goRedis "github.com/go-redis/redis/v9"
)

var redisClient *goRedis.Client

// the first generic is a comparable type used as key
// the second generic is any type used as value
// the first argument of the factory function is the client used to communicate with the redis server
redisCache := redis.New[string, int](redisClient)

// When complex types are passed as values, 
// you can pass an EncodeDecodeOption to specify how your type needs to be serialized
// a default implementation is given as part of the library which relies on the encoding/json package.
redisCache := redis.New[string, user](
    redisClient,
    EncodeDecodeOption[string, user](DefaultEncoder[user], DefaultDecoder[*user]),
)
```

### Multi Level

```go
var c1, c2, c3 cache.Cache

// the first generic is a comparable type used as key
// the second generic is any type used as value
// the arguments passed are the cache levels
multilvl := NewMultiLevel[string, int](c1,c2, c3)

// The behavior of the multilvel cache are documented in the method:

// Get traverse all the caches, if all of them fail it returns a generic ErrNotGet
// Set traverse all the caches, if all of them fail it returns a generic ErrNotSet
// Delete traverse all the caches, if all of them fail it returns a generic ErrNotDelete
```

## Performances

GoCache is a really fast caching solution,
with 0 allocations as well as crazy performances.

Benchmarks comparing it to [dgraph/ristretto](https://github.com/dgraph-io/ristretto),
[allegro/bigcache](https://github.com/allegro/bigcache),
and [patrickmn/go-cache](https://github.com/patrickmn/go-cache)

```text
// dgraph/ristretto

goos: darwin
goarch: amd64
pkg: github.com/damianopetrungaro/go-cache/benchmarks/cache/dgraph
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLogger/dgraph/ristretto.empty-12                4819352              1222 ns/op             240 B/op          6 allocs/op
BenchmarkLogger/dgraph/ristretto.prefilled-12            4267911              1328 ns/op             240 B/op          6 allocs/op
PASS
ok      github.com/damianopetrungaro/go-cache/benchmarks/cache/dgraph   14.753s
```

```
// patrickmn/go-cache

cd ./patrickmn && go1.17.11 test ./... -bench=. -benchmem -benchtime=5s
goos: darwin
goarch: amd64
pkg: github.com/damianopetrungaro/go-cache/benchmarks/cache/patrickmn
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLogger/patrickmn/go-cache.empty-12             16924348               340.0 ns/op            24 B/op          1 allocs/op
BenchmarkLogger/patrickmn/go-cache.prefilled-12         17133223               367.0 ns/op            24 B/op          1 allocs/op
BenchmarkLogger/patrickmn/go-cache.prefilled_with_cleanup-12            16663183               366.7 ns/op            24 B/op          1 allocs/op
PASS
ok      github.com/damianopetrungaro/go-cache/benchmarks/cache/patrickmn        19.469s
```

```
// allegro/bigcache

cd ./allegro && go1.17.11 test ./... -bench=. -benchmem -benchtime=5s
goos: darwin
goarch: amd64
pkg: github.com/damianopetrungaro/go-cache/benchmarks/cache/allegro
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLogger/allegro/bigcache.empty-12               13124148               593.7 ns/op            91 B/op          0 allocs/op
BenchmarkLogger/allegro/bigcache.prefilled-12           10756856               464.9 ns/op            55 B/op          0 allocs/op
BenchmarkLogger/allegro/bigcache.prefilled_with_cleanup-12              14478454               425.3 ns/op            10 B/op          0 allocs/op
PASS
ok      github.com/damianopetrungaro/go-cache/benchmarks/cache/allegro  22.240s

```

```
// damianopetrungaro/go-cache

cd ./damianopetrungaro && go1.18.3 test ./... -bench=. -benchmem -benchtime=5s
goos: darwin
goarch: amd64
pkg: github.com/damianopetrungaro/go-cache/benchmarks/cache
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLogger/damianopetrungaro/go-cache.empty-12             22228510               270.5 ns/op             0 B/op          0 allocs/op
BenchmarkLogger/damianopetrungaro/go-cache.prefilled-12         20715753               268.8 ns/op             0 B/op          0 allocs/op
BenchmarkLogger/damianopetrungaro/go-cache.prefilled_with_cleanup-12            21777249               269.2 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/damianopetrungaro/go-cache/benchmarks/cache  19.755s
```
