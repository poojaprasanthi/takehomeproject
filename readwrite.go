package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const mb = 1024 * 1024
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	ctx := context.TODO()
	InputFilePath := os.Args[1]
	if  InputFilePath == "" {
		panic(fmt.Errorf("Please provide input file path"))
	}
	fmt.Println("InputFilePath:", InputFilePath)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)
	readfile, err := os.Open(InputFilePath)
	check(err)
	writefile1, err := os.OpenFile("results/wordcountperline.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	writefile2, err := os.OpenFile("results/wordmap.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	linecounter := 0
	defer readfile.Close()
	defer writefile1.Close()
	defer writefile2.Close()

	scanner := bufio.NewScanner(readfile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			break
		}
		line := scanner.Text()
		f := func(c rune) bool {
			return unicode.IsSpace(c) || unicode.IsPunct(c)
		}
		words := strings.FieldsFunc(line,f)
		_, err := writefile1.WriteString(fmt.Sprintf("%d - %d\n",linecounter,len(words)))
		check(err)
		for _, w := range words {
			//fmt.Println("w:", w)
			//get first
			val, err := client.Get(ctx, strings.ToLower(w)).Result()
			if err == nil {
				v, _ := strconv.Atoi(val)
				v++
				err = client.Set(ctx, w, strconv.Itoa(v), 0).Err()
				check(err)
				continue
			}
			if err.Error() == "redis: nil" {
				err = client.Set(ctx, w, 1, 0).Err()
				check(err)
			}
		}
		linecounter++
	}

	//initial scan
	cm := client.Scan(ctx,0,"",10)
	res,cursor,err := cm.Result()

	for cursor != 0 {
		for i:=0;i<len(res);i++ {
			val, err := client.Get(ctx, res[i]).Result()
			check(err)
			//fmt.Println("val:", val, res[i])
			_, err = writefile2.WriteString(fmt.Sprintf("%s - %s\n",res[i], val))
			check(err)
		}
		cm = client.Scan(ctx, cursor,"",10)
		res,cursor,err = cm.Result()
	}
}