MAKEFLAGS += --silent

ENV_FILE = .env

COMPOSE = docker compose -f ./backend/deployment/docker/docker-compose.yml --env-file $(ENV_FILE)
COMPOSE_DEV = docker compose -f ./backend/deployment/docker/docker-compose-dev.yml --env-file $(ENV_FILE)

SERVICES = all up down stop db dev gateway meeting meetings request user chat log logs help clean
.PHONY: $(SERVICES)

uniq = $(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$(wordlist 2,$(words $1),$1))))

input = $(if $(MAKECMDGOALS),$(MAKECMDGOALS), up)

#* ARGS = $(call uniq, $(filter $(SERVICES), $(input)))
ARGS = $(filter $(SERVICES), $(input))

# $(info info: $(ARGS))
# $(info Количество слов: $(words $(ARGS)))  
# TODO: добавить обработку вкл/выкл для make ""

# TODO: fix bug with пользователями db
%:
	@:
$(firstword $(ARGS)):
	@if [ "$@" = "$(firstword $(ARGS))" ]; then \
		db="false";\
		dev="false";\
		for word in $(ARGS); do \
		echo "\n"; \
		case "$$word" in \
			up|all) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю ВСЕ сервисы"; \
						$(COMPOSE) up -d --build; \
					else \
						echo "Запускаю ВСЕ сервисы, НО в сервисном режиме"; \
						$(COMPOSE_DEV) up -d --build; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю ВСЕ бдшки"; \
						$(COMPOSE) up -d --build db_meeting db_request db_user db_chat; \
					else \
						echo "Запускаю ВСЕ бдшки в сервисном режиме"; \
						$(COMPOSE_DEV) up -d --build db_meeting db_request db_user db_chat; \
						dev="false"; \
					fi; \
					db="false"; \
				fi \
				;; \
			down|stop) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Останавливаю всё, что есть"; \
						$(COMPOSE) down; \
					else \
						echo "Останавливаю всё, что есть, даже в сервисном режиме!"; \
						$(COMPOSE_DEV) down; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo "Останавливаю всe бдшки"; \
						$(COMPOSE) down db_meeting db_request db_user db_chat; \
					else \
						echo "Останавливаю всe бдшки, сервисный режим им не поможет хе-хе"; \
						$(COMPOSE_DEV) down db_meeting db_request db_user db_chat; \
						dev="false"; \
					fi; \
					db="false"; \
				fi \
				;; \
			log|logs) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Открываю логи"; \
						$(COMPOSE) logs -f; \
					else \
						echo "Открываю логи разработчика)"; \
						$(COMPOSE_DEV) logs -f; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo "Открываю логи бдшек"; \
						$(COMPOSE) logs -f db_meeting db_request db_user db_chat; \
					else \
						echo "Открываю логи разрабных бдшек"; \
						$(COMPOSE_DEV) logs -f db_meeting db_request db_user db_chat; \
						dev="false"; \
					fi; \
					db="false"; \
				fi \
				;; \
			meetings) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю meetinG.."; \
						$(COMPOSE) up -d --build meeting; \
					else \
						echo "Запускаю meetinG... (в сервисном режиме!)"; \
						$(COMPOSE_DEV) up -d --build meeting; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю бдшку для meetinG.."; \
						$(COMPOSE) up -d --build db_meeting; \
					else \
						echo "Запускаю бдшку для meetinG.. с портами, все дела - сервисный режм"; \
						$(COMPOSE_DEV) up -d --build db_meeting; \
						dev="false"; \
					fi; \
					db="false"; \
				fi \
				;; \
			gateway) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю gateway"; \
						$(COMPOSE) up -d --build gateway; \
					else \
						echo "Запускаю gateway в супер режиме разработчика"; \
						$(COMPOSE_DEV) up -d --build gateway; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo 'А нету для gateway бдшки 8)'; \
						db="false"; \
					else \
						echo 'А нету для gateway бдшки 8), даже если ты разраб :('; \
						db="false"; \
						dev="false"; \
					fi; \
				fi \
				;; \
			meeting|request|user|chat) \
				if [ "$$db" = "false" ]; then \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю $$word"; \
						$(COMPOSE) up -d --build $$word; \
					else \
						echo "Запускаю $$word в режиме разработчика"; \
						$(COMPOSE_DEV) up -d --build $$word; \
						dev="false"; \
					fi; \
				else \
					if [ "$$dev" = "false" ]; then \
						echo "Запускаю бдшку для $$word"; \
						$(COMPOSE) up -d --build db_$$word; \
					else \
						echo "Запускаю разработческую бдшку для $$word"; \
						$(COMPOSE_DEV) up -d --build $$word; \
						dev="false"; \
					fi; \
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