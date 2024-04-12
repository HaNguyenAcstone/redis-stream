package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "192.168.2.45:6379", // Địa chỉ Redis server
		DB:   0,
	})

	// Gọi hàm lắng nghe Redis Stream
	streamKey := "Redis_Streams_Data" // tên của stream
	go listenRedisStream(rdb, streamKey)

	select {} // Chặn main không cho thoát
}

// Hàm lắng nghe Redis Stream
func listenRedisStream(rdb *redis.Client, streamKey string) {
	lastID := "0" // Bắt đầu từ ID đầu tiên trong stream

	for {
		// Đọc dữ liệu từ stream
		entries, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamKey, lastID},
			Block:   0,
		}).Result()

		if err != nil {
			fmt.Println("Error reading from stream:", err)
			return
		}

		for _, entry := range entries {
			for _, message := range entry.Messages {
				fmt.Printf("Received message: ID=%s, Values=%v\n", message.ID, message.Values)
				lastID = message.ID // Cập nhật lastID để tiếp tục từ message cuối cùng
			}
		}
	}
}
