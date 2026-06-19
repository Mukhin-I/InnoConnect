MAKEFLAGS += --silent

ENV_FILE = .env

COMPOSE = docker compose -f ./backend/deployment/docker/docker-compose.yml --env-file $(ENV_FILE)
COMPOSE_DEV = docker compose -f ./backend/deployment/docker/docker-compose-dev.yml --env-file $(ENV_FILE)

SERVICES = up down db gateway meeting meetings request user chat log logs help
.PHONY: $(SERVICES)

uniq = $(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$(wordlist 2,$(words $1),$1))))

input = $(if $(MAKECMDGOALS),$(MAKECMDGOALS), up)

#* ARGS = $(call uniq, $(filter $(SERVICES), $(input)))
ARGS = $(filter $(SERVICES), $(input))

# $(info info: $(ARGS))
# $(info Количество слов: $(words $(ARGS)))  
# TODO: добавить обработку dev 
# TODO: добавить обработку вкл/выкл для make ""

# TODO: fix bug with пользователями db
%:
	@:
$(firstword $(ARGS)):
	@if [ "$@" = "$(firstword $(ARGS))" ]; then \
		db="false";\
		dev="false";\
		docker="$(COMPOSE)"; \
		for word in $(ARGS); do \
		echo "\n"; \
		case "$$word" in \
			up) \
				if [ "$$db" = "false" ]; then \
					echo "Запускаю ВСЕ сервисы"; \
					$$docker up -d --build; \
				else \
					echo "Запускаю ВСЕ бдшки"; \
					$$docker up -d --build db_meeting db_request db_user db_chat; \
					db="false"; \
				fi \
				;; \
			down) \
				if [ "$$db" = "false" ]; then \
					echo "Останавливаю всё, что есть"; \
					$$docker down; \
				else \
					echo "Останавливаю всe бдшки"; \
					$$docker down db_meeting db_request db_user db_chat; \
					db="false"; \
				fi \
				;; \
			log|logs) \
				if [ "$$db" = "false" ]; then \
					echo "Открываю логи"; \
					$$docker logs -f; \
				else \
					echo "Открываю логи бдшек"; \
					$$docker logs -f db_meeting db_request db_user db_chat; \
					db="false"; \
				fi \
				;; \
			meetings) \
				if [ "$$db" = "false" ]; then \
					echo "Запускаю meetinG.."; \
					$$docker up -d --build meeting; \
				else \
					echo "Запускаю бдшку для meetinG.."; \
					$$docker up -d --build db_meeting; \
					db="false"; \
				fi \
				;; \
			gateway) \
				if [ "$$db" = "false" ]; then \
					echo "Запускаю gateway"; \
					$$docker up -d --build gateway; \
				else \
					echo 'А нету для gateway бдшки 8)'; \
					db="false"; \
				fi \
				;; \
			meeting|request|user|chat) \
				if [ "$$db" = "false" ]; then \
					echo "Запускаю $$word"; \
					$$docker up -d --build $$word; \
				else \
					echo "Запускаю $$word"; \
					$$docker up -d --build db_$$word; \
					db="false"; \
				fi \
				;; \
			db) \
				if [ "$$db" = "false" ]; then \
					echo "Жду бд"; \
					db="true"; \
				else \
					echo "Очень жду бд"; \
					db="true"; \
				fi \
				;; \
			dev) \
				if [ "$$dev" = "false" ]; then \
					echo "Жду dev"; \
					dev="true"; \
				else \
					echo "Сильно жду dev"; \
					dev="true"; \
				fi \
				;; \
			*) \
				echo "Неизвестный сервис: $$word"; \
				;; \
		esac; \
		done;\
	fi;\