DOCKERCOMPOSECMD=docker-compose

dc-up:
	$(DOCKERCOMPOSECMD) up -d --build

dc-down:
	docker-compose down --remove-orphans

dc-restart: dc-down dc-up