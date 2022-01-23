# takehomeproject


Prerequisite:
1. You need to have docker installed.
2. Pull Redis image and run it in the background.
```
	docker pull redis
        docker run --name redis -p 6379:6379 -d redis
        docker start redis
```
	
**How to run**:

*executable*:

make path="small.txt" service

*test*:

make test

**Algorithm**:
1. Setup a redis in the background. This will be used for storing key values, where key is the word that is read from the input file and value will be the number of occurances in that input file.
2. read each line from file	
	* count all words in the current line and write to the file1
	* update key in redis
		* Check if the key already exists, if yes increment the counter else set to 0
4. Scan thru all the keys from redis and write it to file2.
5. It will calculate the time taken to execute the program and display the log as follows
This is the output for 3MB inputfile.
```
madapuri@Poojas-MBP takehomeproject % make path="SampleInput.txt" service          
go get github.com/go-redis/redis
go build readwrite.go
go run readwrite.go SampleInput.txt
InputFilePath: SampleInput.txt
PONG <nil>
2022/01/22 18:56:08 Execution took 10m58.553909379s

```
The resulting files are stored in results dir.
Example output data :

```
"results/wordcountperline.txt"

0 - 15
1 - 20
2 - 27
```

```
results/wordmap.txt
she - 4
delivery - 4
ready - 4
m - 4
residency - 4
so - 4
but - 4
her - 8
finish - 4
```

Analysis:

I first thought of using `trie`, it will be easy for building and searching. But the downside of it is , it will occupy too much memory with huge input data file.
Also, it can be used only for chars `a-z` OR `A-Z`. So this was not having flexibility of input content.
I first implemented the following logic - 
	1. build trie for 1000 words at a stretch and update these keys in redis
```//type trie struct {
//      Child [26]*trie
//      Word string
//      Count int
//}
//func buildtrie(words []string) *trie {
//      root := &trie{}
//      for _, word := range words {
//              fmt.Println("word:", word)
//              temp := root
//              for _, w := range word {
//                      if temp.Child[w-'a'] == nil {
//                              temp.Child[w-'a'] = &trie{}
//                      }
//                      temp = temp.Child[w-'a']
//              }
//              if len(temp.Word)==0 {
//                      temp.Word = word
//                      temp.Count = 1
//              } else {
//                      temp.Count++
//              }
//      }
//      return root
//}
//func searchtrie(root *trie, client *redis.Client) {
//      ctx := context.TODO()
//      if len(root.Word) != 0 {
//              fmt.Println("word:", root.Word)
//              val, err := client.Get(ctx, root.Word).Result()
//              if err == nil {
//                      v, _ := strconv.Atoi(val)
//                      v+=root.Count
//                      err = client.Set(ctx, root.Word, strconv.Itoa(v), 0).Err()
//                      check(err)
//              }
//              if err.Error() == "redis: nil" {
//                      err = client.Set(ctx, root.Word, 1, 0).Err()
//                      check(err)
//              }
//      }
//      for i:=0;i<26;i++ {
//              if root.Child[i] !=nil {
//                      searchtrie(root.Child[i],client)
//              }
//      }
//}

```
I am not sure how to use AVL trees for this purpose, so did not consider that.

Linked list is not a good choice for this assignment, as it is best used for an assignment which has more insert and deletion operations. Our assignment involves more of searching, like search if the key exists then update. Linked list is not good option when we need more search operations especially when the list size is big.

With all the above thoughts, I stick to using `redis` in the background to manage the key values. Redis is an in-memory store. It can therefore use data structures which are adapted to memory storage and the time omplexity is O(1).

**Requirements that could not be fulfilled**
I Admit, that the current executable works fine for data ~3 MB and not 10GB.
I feel that there is no simple way of processing 10GB data. One thought to achieve this is following the below design.

<img width="830" alt="Screen Shot 2022-01-23 at 11 26 53 AM" src="https://user-images.githubusercontent.com/98238678/150694854-d1845142-fd96-4313-a8c9-f903eaa7093d.png">

1. chunker chunks 10GB to 1MB chunk.
2. Chunk is process using our current executable, which will update file1 and update Redis.
3. After the 1MB data is processed , RDBMS will persist all the keys from Redis and Flushes all keys out of Redis to reset for the next set.
4. Worker groups will pick a key from RDBMS, combine duplicate keys to 1 and total their values.
5. When all 10GB is processed, write file method will write the keys from DB to file2.

The second requirement I could not do is, parallelizing the process. We can achieve this using worker groups and chanels, but I could not finish it with in the timeboxed time.
