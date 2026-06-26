package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis 缓存客户端
type Client struct {
	client *redis.Client
	ttl    time.Duration
}

// Config Redis 缓存配置
type Config struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
}

// NewClient 创建 Redis 缓存客户端
func NewClient(cfg Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败: %w", err)
	}

	ttl := cfg.TTL
	if ttl == 0 {
		ttl = 5 * time.Minute
	}

	return &Client{client: client, ttl: ttl}, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.client.Close()
}

// Get 获取缓存
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Set 设置缓存
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	d := c.ttl
	if len(ttl) > 0 {
		d = ttl[0]
	}
	return c.client.Set(ctx, key, data, d).Err()
}

// Del 删除缓存
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// BuildKey 构建缓存键
func BuildKey(parts ...string) string {
	key := "blog"
	for _, p := range parts {
		key += ":" + p
	}
	return key
}
