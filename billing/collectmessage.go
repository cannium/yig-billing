package billing

import (
	. "github.com/journeymidnight/yig-billing/helper"
	"github.com/journeymidnight/yig-billing/messagebus"
	"github.com/journeymidnight/yig-billing/redis"
	"strconv"
	"time"
)

const (
	TimeLayout           = "2006-01-02 15:04:05"
	SEPARATOR            = ":"
	MinSizeForStandardIa = 64 << 10
)

func collectMessage() {
	go messagebus.StartConsumerReceiver()
	go collector()
}

func collector() {
	for {
		message := <-messagebus.ConsumerMessagePipe
		isEffective := isMessageEffective(message)
		if !isEffective {
			Logger.Println("[INVALID MESSAGE] UUID is:", message.Uuid)
			continue
		}
		if message.Messages[messagebus.OperationName] == "DeleteObject" {
			delDeleteObject(message)
		} else {
			Logger.Println("[INVALID MESSAGE] UUID is:", message.Uuid)
		}
	}
}

func isMessageEffective(message messagebus.ConsumerMessage) bool {
	if message.Messages[messagebus.HttpStatus] != "200" || message.Messages[messagebus.BucketName] == "-" || message.Messages[messagebus.ObjectName] == "-" {
		Logger.Println("[INVALID MESSAGE] Request err, uid is:", message.Uuid)
		return false
	} else {
		if message.Messages[messagebus.StorageClass] != "STANDARD" && message.Messages[messagebus.StorageClass] != "-" {
			return true
		}
		Logger.Println("[INVALID MESSAGE] Request storage class err, uid is:", message.Uuid)
		return false
	}
}

func delDeleteObject(msg messagebus.ConsumerMessage) {
	Logger.Println("[MESSAGE] Enter del delete object, Uuid is:", msg.Uuid)
	info := msg.Messages
	projectId := info[messagebus.ProjectId]
	bucketName := info[messagebus.BucketName]
	objectName := info[messagebus.ObjectName]
	storageClass := info[messagebus.StorageClass]
	objectSize := info[messagebus.ObjectSize]
	deadLine := info[messagebus.LastModifiedTime]
	expire := getDeadLine(deadLine)
	redisMsg := new(redis.MessageForRedis)
	redisMsg.Key = redis.BillingUsagePrefix + projectId + SEPARATOR + storageClass + SEPARATOR + bucketName + SEPARATOR + objectName
	redisMsg.Value = getBillingSize(objectSize)
	redis.RedisConn.SetToRedisWithExpire(*redisMsg, expire, msg.Uuid)
	Logger.Println("[MESSAGE] Delete object successful, Uuid is:", msg.Uuid, ",expire is:", expire, "s")
}

func getBillingSize(objectSize string) (countSize string) {
	size, _ := strconv.ParseInt(objectSize, 10, 64)
	count := (size + MinSizeForStandardIa - 1) / MinSizeForStandardIa * MinSizeForStandardIa
	countSize = strconv.FormatInt(count, 10)
	return
}

func getDeadLine(createTime string) (expire int){
	t, _ := time.Parse(TimeLayout, createTime)
	deadTime := t.Add(720 * time.Hour)
	nowTime := time.Now()
	if deadTime.Unix() > nowTime.Unix() {
		expire = int(deadTime.Unix() - nowTime.Unix())
		return 
	}
	return 0
}
