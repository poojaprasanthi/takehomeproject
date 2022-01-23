package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const mb = 1024 * 1024
func check(e error) {
	if e != nil {
		panic(e)
	}
}
var Client *redis.Client
var linecount int
func main() {
	start := time.Now()
	ctx := context.TODO()
	InputFilePath := os.Args[1]
	if  InputFilePath == "" {
		panic(fmt.Errorf("Please provide input file path"))
	}
	fmt.Println("InputFilePath:", InputFilePath)
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	pong, err := Client.Ping(ctx).Result()
	fmt.Println(pong, err)
	readfile, err := os.Open(InputFilePath)
	check(err)
	writefile1, err := os.OpenFile("results/wordcountperline.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
	writefile2, err := os.OpenFile("results/wordmap.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	check(err)
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
		WriteToFile1(words, writefile1)
		linecount++
		UpdateRedis(words, Client)
	}
	WritetoFile2(Client,writefile2)

	elapsed := time.Since(start)
	log.Printf("Execution took %s", elapsed)
}

func WriteToFile1(words []string, writefile1 *os.File) {
	_, err := writefile1.WriteString(fmt.Sprintf("%d - %d\n",linecount,len(words)))
	check(err)
}

func UpdateRedis(words []string,client *redis.Client) {
	ctx := context.TODO()
	for _, w := range words {
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
}
func WritetoFile2(client *redis.Client, writefile2 *os.File) {
	ctx := context.TODO()
	cm := client.Scan(ctx,0,"",10)
	res,cursor,err := cm.Result()
	check(err)
	for cursor != 0 {
		for i:=0;i<len(res);i++ {
			val, err := client.Get(ctx, res[i]).Result()
			check(err)
			_, err = writefile2.WriteString(fmt.Sprintf("%s - %s\n",res[i], val))
			check(err)
		}
		cm = client.Scan(ctx, cursor,"",10)
		res,cursor,err = cm.Result()
	}
}
//type trie struct {
//	Child [26]*trie
//	Word string
//	Count int
//}
//func buildtrie(words []string) *trie {
//	root := &trie{}
//	for _, word := range words {
//		fmt.Println("word:", word)
//		temp := root
//		for _, w := range word {
//			if temp.Child[w-'a'] == nil {
//				temp.Child[w-'a'] = &trie{}
//			}
//			temp = temp.Child[w-'a']
//		}
//		if len(temp.Word)==0 {
//			temp.Word = word
//			temp.Count = 1
//		} else {
//			temp.Count++
//		}
//	}
//	return root
//}
//
//func searchtrie(root *trie, client *redis.Client) {
//	ctx := context.TODO()
//	if len(root.Word) != 0 {
//		fmt.Println("word:", root.Word)
//		val, err := client.Get(ctx, root.Word).Result()
//		if err == nil {
//			v, _ := strconv.Atoi(val)
//			v+=root.Count
//			err = client.Set(ctx, root.Word, strconv.Itoa(v), 0).Err()
//			check(err)
//		}
//		if err.Error() == "redis: nil" {
//			err = client.Set(ctx, root.Word, 1, 0).Err()
//			check(err)
//		}
//	}
//	for i:=0;i<26;i++ {
//		if root.Child[i] !=nil {
//			searchtrie(root.Child[i],client)
//		}
//	}
//}
