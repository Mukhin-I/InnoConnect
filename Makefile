ENV_FILE = .env

COMPOSE = docker compose -f ./backend/deployment/docker/docker-compose.yml --env-file $(ENV_FILE)
COMPOSE_DEV = docker compose -f ./backend/deployment/docker/docker-compose-dev.yml --env-file $(ENV_FILE)

SERVICES = up down db gateway meeting meetings request user chat log logs help
.PHONY: $(SERVICES)

uniq = $(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$(wordlist 2,$(words $1),$1))))

input = $(if $(MAKECMDGOALS),$(MAKECMDGOALS), up)

#* ARGS = $(call uniq, $(filter $(SERVICES), $(input)))
ARGS = $(filter $(SERVICES), $(input))

$(info info: $(ARGS))
# TODO: ДОбавть  разную обработку db
# TODO: добавить обработку вкл/выкл для make ""
# TODO: fix bug with пользователями db
# TODO: добавить обработку dev 
# TODO добавить обработку дублирующихся команд


$(SERVICES):
	@echo "hi $@"
	@case "$@" in \
		up) \
			echo "Запускаю ВСЕ сервисы"; \
			$(COMPOSE) up -d; \
			;; \
		down) \
			echo "Останавливаю всё, что есть"; \
			$(COMPOSE) down; \
			;; \
		log|logs) \
			echo "Открываю логи"; \
			$(COMPOSE) logs -f; \
			;; \
		meetings) \
			echo "Запускаю meetinG.."; \
			$(COMPOSE) up -d meeting; \
			;; \
		gateway|meeting|request|user|chat) \
			echo "Запускаю $@"; \
			$(COMPOSE) up -d $@; \
			;; \
		*) \
			echo "Неизвестный сервис: $@"; \
			;; \
	esac; \