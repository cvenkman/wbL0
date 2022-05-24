
APP = subscriber
PUB = publisher

all: $(APP) $(PUB)

$(APP):
	go build -o $@ cmd/subscriber/main.go

$(PUB):
	go build -o $@ cmd/publisher/publisher.go

fclean:
	@rm $(APP) $(PUB)
	@echo "rm $(APP) and $(PUB)"

re: fclean all

# clean:
# 	@bash close.sh
# 	@echo "порты 4222 и 8080 освобождены"

# run: $(APP) $(PUB)
# 	bash run.sh
	
