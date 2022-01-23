# takehomeproject

How to run the executable:
Prerequisite:
1. You need to have docker installed.
2. Pull Redis image and run it in the background.
	docker pull redis
        docker run --name redis -p 6379:6379 -d redis
        docker start redis
Command to run:
make path="small.txt" service

