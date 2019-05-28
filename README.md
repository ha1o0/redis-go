## redis-go
the redis server that is implemented by golang

### How to use
+ `git clone git@github.com:weifengsmile/redis-go.git`
+ Install official redis-cli
    + For Mac, just use `brew install redis`
    + For Linux, `apt-get`/`yum`...
    + Or use the source codes to make install manually
+ Start golang-redis-server of this repo(default port is **6378**).
+ Start the official redis client to connect this golang-redis-server by `redis-cli -p 6378`
+ Just enjoy redis commands like `ping,get,set,exists,del,rpush,rpop,lpop,llen……`