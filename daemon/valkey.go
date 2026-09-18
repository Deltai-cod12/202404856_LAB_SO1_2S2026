package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ValkeyClient struct {
	client *redis.Client
}

func newValkeyClient() *ValkeyClient {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &ValkeyClient{
		client: client,
	}
}

func (v *ValkeyClient) ping() error {

	result, err := v.client.Ping(ctx).Result()

	if err != nil {
		return err
	}

	fmt.Println("Valkey:", result)

	return nil
}

func (v *ValkeyClient) close() error {
	return v.client.Close()
}

func (v *ValkeyClient) saveRAM(ram RAMInfo) error {

	timestamp := time.Now().UnixMilli()

	data, err := json.Marshal(ram)

	if err != nil {
		return err
	}

	// Estado actual de RAM.
	err = v.client.Set(
		ctx,
		"so1:ram:latest",
		data,
		0,
	).Err()

	if err != nil {
		return err
	}

	// Historial de RAM.
	err = v.client.ZAdd(
		ctx,
		"so1:ram:history",
		redis.Z{
			Score:  float64(timestamp),
			Member: fmt.Sprintf("%d|%s", timestamp, string(data)),
		},
	).Err()

	return err
}

func (v *ValkeyClient) saveContainers(containers []ContainerInfo) error {

	timestamp := time.Now().UnixMilli()

	for _, container := range containers {

		if container.Process == nil {
			continue
		}

		data := map[string]interface{}{
			"timestamp":        timestamp,
			"id":               container.ID,
			"name":             container.Name,
			"pid":              container.PID,
			"profile":          container.Profile,
			"resource":         container.Resource,
			"vsz_kb":           container.Process.VSZKB,
			"rss_kb":           container.Process.RSSKB,
			"mem_percent_x100": container.Process.MemPercentX100,
			"cpu_percent_x100": container.Process.CPUPercentX100,
			"status":           "active",
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			return err
		}

		// Ultimo estado conocido del contenedor.
		latestKey := fmt.Sprintf(
			"so1:container:latest:%s",
			container.ID,
		)

		err = v.client.Set(
			ctx,
			latestKey,
			jsonData,
			0,
		).Err()

		if err != nil {
			return err
		}

		// Historial de metricas de contenedores.
		err = v.client.ZAdd(
			ctx,
			"so1:containers:history",
			redis.Z{
				Score: float64(timestamp),
				Member: fmt.Sprintf(
					"%d|%s|%s",
					timestamp,
					container.ID,
					string(jsonData),
				),
			},
		).Err()

		if err != nil {
			return err
		}
	}

	return nil
}