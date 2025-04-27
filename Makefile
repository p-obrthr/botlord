NAME=botlord

build:
	go build -o $(NAME)

docker-build:
	docker build -t $(NAME) .

docker-run:
	docker run -e DISCORD_BOT_TOKEN=$(DISCORD_BOT_TOKEN) $(NAME) 

clean:
	rm -f $(NAME)
