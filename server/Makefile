NAME=botlord

build:
	go mod tidy
	go build -o $(NAME)

docker-build:
	docker build -t $(NAME) .

docker-registry:
	docker build -t $(REGISTRY_IP):5000/$(NAME) .
	docker push $(REGISTRY_IP):5000/$(NAME)

docker-run-api:
	docker run -d --name botlord \
		-e DISCORD_BOT_TOKEN=$(DISCORD_BOT_TOKEN) \
		-e TEXT_CHANNEL_ID=$(TEXT_CHANNEL_ID) \
		-e ENABLE_API=1 \
		-p 8080:8080 \
		-v botlord-db-volume:/app/db \
		botlord   

docker-run:
	docker run -d --name botlord \
		-e DISCORD_BOT_TOKEN=$(DISCORD_BOT_TOKEN) \
		-e TEXT_CHANNEL_ID=$(TEXT_CHANNEL_ID) \
		-e ENABLE_API=0 \
		-v botlord-db-volume:/app/db \
		botlord
	
clean:
	rm -f $(NAME)
