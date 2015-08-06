# ghost
A set of go tools with simple construction and powerful features.  
As _ghost_ is forced to be external depandency free, 
it is easy to run go codes with _ghost_ without downloading varieties of extra libs.
This implies that _ghost_ might be a best solution for internal enterprise servers with internet cut off
or some users who have harsh conditions to access abroad, especially CHN developers.
## Connection Pool

### Example

```go
// Create a factory() to be used with the Pool
factory := func() (net.Conn, error) {
    return net.Dial("tcp", "127.0.0.1:4000")
}

// Create a new blocking pool with an init capacity of 5 and max capacity of 30.
// Currently no new connections would be made when the pool is busy, so
// the number of connections of the pool is kept no more than 5 and 30 does not
// make sense though the api is reserved. 
// The timeout to block Get() is set to 3 by default.
p, err := pool.NewBlockingPool(5, 30, factory)

// Get a connection from the pool, if there is no connection available
// it blocks for a recycled one is available til timeout.
conn, err := p.Get()

// Use net.Conn's Write() method
// This conn is marked as unusable automatically when an error counted,
// and will be later abandoned when you close the conn, which means
// the underlying connection will be closed.
// Yet a new net.Conn will be created through factory() and put back to the pool. 
lens, err := conn.Write(bytes)

// You can also use other net.Conn's methods
// The underlying connection is only wrapped with all features reserved.

// Put it back to the pool by closing the connection
conn.Close()

```

## Shared File

### Example

## Retry

### Delegate

```go
type Operation func() error
type Recursion func(nSub1 time.Duration) (n time.Duration)
```

### Simple

```go
//Retry at most 3 times.
//Sleeps for 1 second before first retry, and sleep time doubles after each time it retries
retries, errors := retry.Attempt(3, 1*time.Second, func() error{
	//if error is counted, 
	//	return the error
	//if operation is success
	//	return nil
})
```

### Custom

```go
r := &retry.Retry{
	//switch to activate randomize
	Randomize: false,

	//time to sleep before first retry
	FirstSleep: 1 * time.Second,
	//range of sleep time
	MinSleep: 0 * time.Second,
	MaxSleep: 3 * time.Second,

	//sleep time increase or decline
	//linear:
	//	n = nsub1 + 1
	//exponent:
	//	n = nsub1 * 2
	Recursion: retry.Double,
	//Recursion: func(nsub1 time.Duration) time.Duration {
	//	return nsub1+1 * time.Second
	//},

	//max time to retry
	Retries: 5,
}
//return nil to announce success
//return a counted error to announce failure, requesting for retry
//see example/custom.go for detail
r.Attempt(yourOperation)
```
