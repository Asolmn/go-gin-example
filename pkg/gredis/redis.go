package gredis

import (
	"encoding/json"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,     // 最大空闲连接数
		MaxActive:   setting.RedisSetting.MaxActive,   //在给定时间，允许分配的最大连接数
		IdleTimeout: setting.RedisSetting.IdleTimeout, // 在给定时间将回保持空闲状态
		Dial: func() (redis.Conn, error) { // 创建和配置应用程序连接函数
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil { // 向Redis服务器发送AUTH命令并返回收到的答复
					errclose := c.Close()
					if errclose != nil {
						return nil, errclose
					}
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 可选的应用程序检查健康功能
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// 设置key
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Close Fail")
		}
	}(conn)

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value) // 绑定键值
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time) // 设置键过期时间
	if err != nil {
		return err
	}
	return nil

}

// 检测key是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Close Fail")
		}
	}(conn)

	exists, err := redis.Bool(conn.Do("EXISTS", key)) // 检查key是否存在，返回布尔值
	if err != nil {
		return false
	}
	return exists
}

// 返回key对应的缓存，并为[]bytes结构
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Close Fail")
		}
	}(conn)

	reply, err := redis.Bytes(conn.Do("GET", key)) // 获取key所对应的value，返回转为Bytes
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// 删除key对应的缓存
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Close Fail")
		}
	}(conn)

	return redis.Bool(conn.Do("DEL", key)) // 删除key，返回布尔值
}

// 删除所有具备参数关键字的key
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			panic("Close Fail")
		}
	}(conn)
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*")) // 寻找所有符合给定模式的key，将命令返回转为[]string
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
