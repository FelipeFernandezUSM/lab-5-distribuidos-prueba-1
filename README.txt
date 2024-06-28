
docker-broker: ## Initiates the Docker code for the broker server
	@docker build -f dockerfile.broker -t broker .
	@docker run -p 50051:50051 -d --name broker broker

docker-vanguardia: ## Initiates the Docker code for the vanguard server
	@docker build -f dockerfile.comandante -t comandante .
	@docker run -i -p 50050:50050 --name comandante comandante

docker-f1: ## Initiates the Docker code for fulcrum 1
	@docker build -f dockerfile.fulcrum --build-arg NUM=1 -t fulcrum1 .
	@docker run -p 50056:50056 -d --name fulcrum1 fulcrum1

docker-f2: ## Initiates the Docker code for fulcrum 2
	@docker build -f dockerfile.fulcrum --build-arg NUM=2 -t fulcrum2 .
	@docker run -p 50057:50057 -d --name fulcrum2 fulcrum2

docker-f3: ## Initiates the Docker code for fulcrum 3
	@docker build -f dockerfile.fulcrum --build-arg NUM=3 -t fulcrum3 .
	@docker run -p 50058:50058 -d --name fulcrum3 fulcrum3

docker-i1: ## Initiates the Docker code for Caiatl
	@docker build -f dockerfile.jeth -t jeth .
	@docker run -i -p 50052:50052 --name jeth jeth

docker-i2: ## Initiates the Docker code for Osiris
	@docker build -f dockerfile.malkor -t malkor .
	@docker run -i -p 50053:50053 --name malkor malkor