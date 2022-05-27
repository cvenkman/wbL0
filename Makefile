APP = subscriber
PUB = publisher

all: $(APP) $(PUB)

$(APP):
	go build -o $@ cmd/server/main.go

$(PUB):
	go build -o $@ cmd/publisher/publisher.go

fclean:
	@rm $(APP) $(PUB)
	@echo "rm $(APP) and $(PUB)"

re: fclean all
