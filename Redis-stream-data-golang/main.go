package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	streamName  = "Redis_Streams_Data"
)

func init() {
	// Khởi tạo kết nối Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "192.168.2.45:6379", // Địa chỉ Redis server
		DB:   0,                   // Sử dụng DB mặc định
	})

	// Kiểm tra xem stream có tồn tại không, nếu không thì tạo
	exists := redisClient.Exists(ctx, streamName).Val()
	if exists == 0 {
		// Tạo Redis Stream
		redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"init": "start"},
		})
	}
}

func pushOrdersToRedis(orders []map[string]int) {
	for _, order := range orders {
		// Chuyển đổi đơn hàng thành JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshalling order: %v", err)
			continue
		}
		// Đẩy đơn hàng vào Redis Stream
		redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"order": string(orderJSON)},
		})
	}
}

func pushOrders(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Chuyển đổi value thành số
	var numOrders int
	fmt.Sscanf(value, "%d", &numOrders)

	// Tạo danh sách đơn hàng
	orders := make([]map[string]int, numOrders)
	for i := 1; i <= numOrders; i++ {
		orders[i-1] = map[string]int{"order_id": i}
	}

	// Thời điểm bắt đầu
	startTime := time.Now()

	// Gửi đơn hàng vào Redis Stream
	pushOrdersToRedis(orders)

	// Tính toán thời gian hoàn thành
	completionTime := time.Since(startTime)

	// Log thời gian hoàn thành
	log.Printf("Request completed in %s", completionTime)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Orders pushed to Redis Stream successfully",
		"completion_time": completionTime.Seconds(),
	})
}

func main() {
	r := gin.Default()

	// Thiết lập route
	r.GET("/push_orders", pushOrders)

	// Khởi chạy server
	r.Run(":8080")
}
