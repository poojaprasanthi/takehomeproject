package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"testing"
)
func TestWriteToFile1(t *testing.T) {
	words := []string{"check1", "check2", "check3"}
	tesfile1, _ := os.OpenFile("results/test1.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer tesfile1.Close()
	WriteToFile1(words,tesfile1)
	scanner := bufio.NewScanner(tesfile1)
	scanner.Split(bufio.ScanWords)
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	i:=0
	for scanner.Scan() {
		word := scanner.Text()
		if word != words[i] {
			t.Errorf("Words are not correctly written into file")
		}
		i++
	}
}

func TestUpdateRedis(t *testing.T) {
	ctx := context.TODO()
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	val, _ := Client.Get(ctx, "check").Result()
	expected, _ := strconv.Atoi(val)

	UpdateRedis([]string{"check"},Client)
	val, _ = Client.Get(ctx, "check").Result()
	actual, _ := strconv.Atoi(val)
	if expected+1 != actual {
		t.Errorf("update to redis failed %d, %d", expected+1,actual)
	}
}
