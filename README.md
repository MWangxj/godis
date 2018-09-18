# godis

## 使用说明

    获取 go get -u -v github.com/MWangxj/godis

    package godis // import "github.com/MWangxj/godis"

    func Dail(addr, passwd string, maxidle, maxac, db int, to time.Duration) *redis.Pool
    func Del(key string) error
    func Get(key string, v interface{}) error
    func HDel(key, region string) error
    func HGet(key, region string, v interface{}) error
    func HGetAll(key string, v *map[string]string) error
    func HSet(key, region string, v interface{}) error
    func SAdd(key string, v interface{}) error
    func SMembers(key string, v []string) error
    func Set(key string, v interface{}) error
    func SetKeyPrefix(prefix string)
    func SetKeySubfix(subfix string)
