NAME=botlord

build:
	go build -o $(NAME)

docker-build:
	docker build -t $(NAME) .

docker-run-api:
	docker run -d --name botlord \
		-e DISCORD_BOT_TOKEN=$(DISCORD_BOT_TOKEN) \
		-e ENABLE_API=1 \
		-p 8080:8080 \
		botlord   

docker-run:
	docker run -d --name botlord \
		-e DISCORD_BOT_TOKEN=$(DISCORD_BOT_TOKEN) \
		-e ENABLE_API=0 \
		botlord
	
clean:
	rm -f $(NAME)
