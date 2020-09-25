**dlock** is a client side distributed locking library.

You can implement _Store_ interface with the storage client of your choice, eg. redis, aerospike, etc.
By default it uses an in-memory map store.

Usage:

```
// initialize a lock object with redis storage client
lock := dlock.New(context.Background(), redisClientImpl)

// take lock on key "myKey" for 10s
err := lock.Take("myKey", 10)

// release lock after function returns
defer lock.Release()

// do something
```

To implement a _Store_ interface, you will need to write 3 methods -

```
func (c *myCustomStoreClientImpl) Set(ctx context.Context, key, uid string, expiry time.Time) error

func (c *myCustomStoreClientImpl) Get(ctx context.Context, key string) string

func (c *myCustomStoreClientImpl) Delete(ctx context.Context, key string) error
```