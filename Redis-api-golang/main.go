package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func init() {
	// Khởi tạo kết nối Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "192.168.38.128:6379", // Địa chỉ Redis server
		DB:   0,                     // Sử dụng DB mặc định
	})
}

var orderIndex = 100000
var streamName = "Orders"

func pushOrdersToRedis() {
	var info = map[string]interface{}{
		"orderID":    orderIndex,
		"deliveryID": RandomDeliveryID(6),
		"status":     RandomStatus(),
	}

	ordInfoJSON, _ := json.Marshal(info)

	// Đẩy đơn hàng vào Redis Stream
	redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{"order_info": string(ordInfoJSON)},
	})

	fmt.Printf("Đã lưu thông tin đơn hàng thành công %d", orderIndex+1)
	orderIndex = orderIndex + 1
}

func RandomDeliveryID(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

var statuses = []string{"Processing", "Confirmed", "Delivery", "Completed"}

func RandomStatus() string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	return statuses[r.Intn(len(statuses))]
}

func getOrdersToRedis(key string) []map[string]interface{} {
	var dataReponse []map[string]interface{}
	var cursor uint64
	var keys []string

	var err error
	keys, cursor, err = redisClient.Scan(ctx, cursor, "order:*", 10).Result()

	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		info, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			fmt.Printf("Không thể lấy thông tin đơn hàng %s: %v\n", key, err)
			return dataReponse
		}

		item := map[string]interface{}{
			"OrderID":     strings.Split(key, ":")[1],
			"DeliveryID":  info["deliveryID"],
			"OrderStatus": info["status"],
		}

		dataReponse = append(dataReponse, item)
	}

	if cursor == 0 { // Khi SCAN hoàn thành
		fmt.Println("Handle Data Success")
	}

	return dataReponse
}

func pushOrders(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Thời điểm bắt đầu
	startTime := time.Now()

	// Gửi đơn hàng vào Redis Stream
	pushOrdersToRedis()

	// Tính toán thời gian hoàn thành
	completionTime := time.Since(startTime)

	// Log thời gian hoàn thành
	log.Printf("Request completed in %s", completionTime)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Orders pushed to Redis Stream successfully",
		"completion_time": completionTime.Seconds(),
	})
}

func getOrders(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	dataResponse := getOrdersToRedis(key)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    dataResponse,
	})
}

func startGinServer(port string) {
	r := gin.Default()

	// Thiết lập route
	r.GET("/push_orders", pushOrders)
	r.GET("/get-order", getOrders)

	// Log và khởi động server
	log.Printf("Starting Gin server on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Gin server on port %s: %v", port, err)
	}
}

func main() {
	var wg sync.WaitGroup
	ports := []string{"4000", "4001", "4002", "4003", "4004"}

	for _, port := range ports {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			startGinServer(port)
		}(port)
	}

	// Chờ cho tất cả servers kết thúc (trong trường hợp này thì không bao giờ xảy ra)
	wg.Wait()
}
