setup:
	-docker stop etcd
	-docker stop postgres
	-docker rm etcd
	-docker rm postgres
	docker build -t etcd images/etcd
	docker build -t gotooling images/gotooling
	docker build -t postgres images/postgres
	docker run -d -p 7001:7001 -p 4001:4001 --name etcd etcd
	docker run -d -p 5432:5432 --name postgres postgres
