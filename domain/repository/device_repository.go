package repository

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type redisDeviceCheching struct {
	redis *redis.Client
	mongo *mongo.Client
}

var ctx = context.Background()

func NewRepo(redis *redis.Client, mongo *mongo.Client) redisDeviceCheching {
	return redisDeviceCheching{redis: redis, mongo: mongo}
}

func (rdb redisDeviceCheching) SetData(devicename string, imei string) error {
	redis_key := imei
	fmt.Println("Repository : ", imei, " Values : ", devicename)
	err := rdb.redis.SetNX(ctx, redis_key, devicename, 10*time.Second).Err()
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}
func (rdb redisDeviceCheching) GetData(imei string) (string, error) {
	redis_key := imei
	result, err := rdb.redis.Get(ctx, redis_key).Result()
	if err != nil {
		return "", err
	}
	return result, nil

}
func (rdb redisDeviceCheching) SetBackUp(devicename string, imei string) error {

	device := DeviceDomain{
		Imei:       imei,
		DeviceName: devicename,
		CreatedAt:  time.Now().UTC(),
	}

	collection := rdb.mongo.Database("memory").Collection("backup_gateway")
	_, err := collection.InsertOne(ctx, device)
	if err != nil {
		return err
	}
	return nil

}
func (rdb redisDeviceCheching) GetBackUp(imei string) (string, error) {
	collection := rdb.mongo.Database("memory").Collection("backup_gateway")
	//result := make(map[any]any)
	var device DeviceDomain
	err := collection.FindOne(ctx, bson.M{"imei": imei}).Decode(&device)
	if err == mongo.ErrNoDocuments {
		return "", mongo.ErrNoDocuments
	}
	if err != nil {
		return "", err
	}
	fmt.Println(" Repo Mognodb : ", device.Imei)
	return device.Imei, nil

}
