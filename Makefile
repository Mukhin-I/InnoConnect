MAKEFLAGS += --silent # Убрать стандартные warnings из-за повторяюбщихся кодовых слов и встроенной логики make 

ENV_FILE = .env

COMPOSE = docker compose -f ./backend/deployment/docker/docker-compose.yml --env-file $(ENV_FILE) # Подготовленные команды docker, сохранённые в переменных
COMPOSE_DEV = docker compose -f ./backend/deployment/docker/docker-compose-dev.yml --env-file $(ENV_FILE)

SERVICES = all up down stop switch selfDestroySequence db dev gateway meeting meetings request user chat log logs help cache remove init # объявление кодовых слов
.PHONY: $(SERVICES) # не запускать как файлы с подобными именами

uniq = $(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$(wordlist 2,$(words $1),$1)))) # функция для фильтрации слов по одному разу

input = $(if $(MAKECMDGOALS),$(MAKECMDGOALS), up) # Если есть ввод - брать ввод, если make без аргументов - подразумевать подъём всего

 #* ARGS = $(call uniq, $(filter $(SERVICES), $(input)))
ARGS = $(filter $(SERVICES), $(input)) # оставляем только кодовые фразы без опечаток и мусора

 # $(info info: $(ARGS))
 # $(info Количество слов: $(words $(ARGS)))  

%: # заглушка для слов с опечатками(make запускает пайплайн для каждого слова, в том числе и неправильно написанного, вне зависимости хотим мы того или нет)
	@:
$(firstword $(ARGS)): # берём из отфильтрованных слов первое и на нём триггерим проход по всей входной строчке 
	@if [ "$@" = "$(firstword $(ARGS))" ]; then # на остальных кодовых словах ничего не делаем, чтобы один и тот же цикл не выполнять N слов=раз\
		db="false"; # со знака @ тут начинается bash синтаксис, объявляем флаги\
		dev="false"; # из полезного: $() - считать переменную make, $$.. - считать переменную bash\
		remove="false";\
		_remove() { # вспомогательная bash функция для удаления контейнеров\
			target="$$1"; \
			echo "Удаляю $$target: и контейнер, и образ, и кеш, и volume'ы"; \
			cid=`$(COMPOSE) images -q $$target 2>/dev/null`; \
			if [ -z "$$cid" ]; then cid=`$(COMPOSE_DEV) images -q $$target 2>/dev/null`; fi; \
			docker stop $$target 2>/dev/null; # стопаем контейнеры \
			docker rm -f -v $$target 2>/dev/null; # чистим volumes \
			vols=`docker volume ls -q --filter "label=com.docker.compose.service=$$target"`; \
			if [ -n "$$vols" ]; then docker volume rm $$vols; fi; # чистим volumes ещё раз \
			if [ -n "$$cid" ]; then \
				docker buildx prune --filter "parents=$$cid" -f; # чистим кеш сборки \
				docker rmi $$cid; # ещё что-то чистим \
			else \
				echo "Не нашёл образ $$target, нечего очищать"; \
			fi; \
		}; \
		for word in $(ARGS); do # тепербь для каждого слова в строке ввода(после фильтрации, оставив только знакомые слова) \
		echo "\n"; \
		case "$$word" in #switch case по словам\
			up|all) # поднимаем *всё* с обработкой флагов db => поднять *всё* бд, dev => поднять *всё* [бд] в dev режиме, remove => удалить *всё* [бд] \
				if [ "$$remove" = "true" ]; then \
					if [ "$$db" = "true" ]; then \
						echo "!Удаляю ВСЕ бдшки"; \
						for s in db_meeting db_request db_user db_chat; do _remove $$s; done; \
					else \
						echo "!ЁУдаляю ВСЕ сервисы"; \
						for s in gateway meeting request user chat db_meeting db_request db_user db_chat; do _remove $$s; done; \
					fi; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then # 4 варианта соответственно для комбинаций db=true/false, dev=true/false \
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
				fi; \
				echo "Запускаю фронтend"; # в конце prodная история, поднимаем ещё и фронтend \
				(cd frontend && npm run dev) \
				;; \
			down|stop) # то же самое, но вместа запуска, останавливаем \
				if [ "$$remove" = "true" ]; then \
					echo "удалить ВСЁ??? точно? я боюсь... если надо, напиши тогда remove all"; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
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
			switch) # если все сервисы подняты -> выключаем, если не все -> доподнимаем \
				if [ "$$remove" = "true" ]; then \
					echo "switch не дружит с remove, забыл стереть флаг?"; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				else \
					all_up="true"; \
					for s in gateway meeting request user chat db_meeting db_request db_user db_chat; do \
						state=`$(COMPOSE) ps -q $$s 2>/dev/null`; \
						if [ -z "$$state" ]; then all_up="false"; fi; \
					done; \
					if [ "$$all_up" = "true" ]; then \
						echo "Все сервисы запущены, выключаю всё"; \
						$(COMPOSE) down; \
					else \
						echo "Не все сервисы запущены, запускаю всё"; \
						$(COMPOSE) up -d --build; \
						echo "Запускаю фронтend"; \
						(cd frontend && npm run dev) \
					fi; \
					db="false"; \
					dev="false"; \
				fi \
				;; \
			log|logs) # открыть логи все или бдшек \
				if [ "$$remove" = "true" ]; then \
					echo "Что написано пером, уже не вырубишь топором"; \
					echo "Всё, что будет в будущем, будет в будущем, всё, что было - уже история, живи настоящим. Зачем хочешь удалить историю)?"; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
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
			meetings) # отдельный кейс с неправильным названием сервиса с буквой С в конце - популярная моя ошибка со стёбными комментариями \
				if [ "$$remove" = "true" ]; then # remove meeting или db meeting (логично?) удаляет это контейнер \
					if [ "$$db" = "true" ]; then \
						_remove db_meeting; \
					else \
						_remove meeting; \
					fi; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
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
			gateway) # отдельный кейс для gateway, потому что у него нет бд =) обработку этого флага нужно писать иначе \
				if [ "$$remove" = "true" ]; then \
					_remove gateway; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
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
			meeting|request|user|chat) # для основной группы сервисок всё одинаково, парсим название через $$word и передаём в запуск, удаление, бдшки и бла бла бла \
				if [ "$$remove" = "true" ]; then \
					if [ "$$db" = "true" ]; then \
						_remove db_$$word; \
					else \
						_remove $$word; \
					fi; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
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
			db) # УраЁ!!! сервисы закончились, остались флаги, кеш и инициализация \
				if [ "$$db" = "false" ]; then \
					echo "Жду бд"; \
					db="true"; \
				else \
					echo "Очень жду бд"; # когда встретилось db при уже активированном db, например make db db meeting \
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
			remove) \
				if [ "$$remove" = "false" ]; then \
					echo "Жду сервис для удаления"; \
					remove="true"; \
				else \
					echo "ДАЙТЕ мне сервис для удаления"; \
					remove="true"; \
				fi \
				;; \
			cache) # очистка кеша, TODO: обработка dev и db не сделана, чистит сразу всё \
				if [ "$$remove" = "true" ]; then # когда make remove cache = make cache\
					echo "Всё, понял, если что, можешь просто cache писать\nНо так не оч красиво((("; \
					remove="false"; \
				fi;\
				echo "Чищу кеш сборки..."; \
				ids=""; \
				for s in gateway meeting request user chat db_meeting db_request db_user db_chat; do \
					id=`$(COMPOSE) images -q $$s 2>/dev/null`; \
					if [ -z "$$id" ]; then id=`$(COMPOSE_DEV) images -q $$s 2>/dev/null`; fi; # если не нашли id('') то заполнить нулями и потом проверить   P.s.[не 100%]TODO \
					if [ -n "$$id" ]; then \
						echo "Нашёл $$s: ($$id)"; \
						ids="$$ids;$$id"; \
					else \
						echo "Образа $$s нет, пропускаю"; \
					fi; \
				done; \
				ids=`echo $$ids | sed 's/^;//'`; \
				if [ -n "$$ids" ]; then \
					docker buildx prune --filter "parents=$$ids" -f; \
				else \
					echo "Нечего чистить — ещё ничего не было"; \
				fi \
				;; \
			selfDestroySequence) # команда самоуничтожения для удобного удаления проекта[нужно знать пароль] TODO спрятать пароль или хранить его хеш \
				echo "Так-так-так. Кнопка самоуничтожения. Серьёзно?"; \
				echo "Это удалит ВСЕ docker-контейнеры, образы, volume'ы, кеш сборки И все файлы проекта с диска. Без права на отмену."; \
				echo "Чтобы убедиться, что ты делаешь это осознанно, ответь правильно на 3 вопроса (или хотя бы попытайся)."; \
				echo "\n"; \
				echo "Вопрос 1 из 3: ты точно уверен, что это то, чего ты хочешь? \nКакой твой любимый фрукт?"; \
				read -r answer1 < /dev/tty; \
				echo "Вопрос 2 из 3: ты осознаёшь, что обратной дороги не будет, и я не сохраняю бэкапы?\n Ты можешь написать картину?"; \
				read -r answer2 < /dev/tty; \
				echo "Вопрос 3 из 3: последний шанс передумать. Точно удаляем ВСЁ? \nЧто такое Чёрный кофе?"; \
				read -r answer3 < /dev/tty; \
				if [ "$$answer1" = "яблоко" ] && [ "$$answer2" = "робот" ] && [ "$$answer3" = "Агата Кристи" ]; then \
					echo "Ну ты сам напросился. Поехали..."; \
					sleep 1; \
					echo "Гашу и стираю все сервисы..."; \
					for s in gateway meeting request user chat db_meeting db_request db_user db_chat; do _remove $$s; done; \
					echo "Чищу остатки buildx кеша..."; \
					docker buildx prune -af 2>/dev/null; \
					echo "Сношу volume'ы проекта..."; \
					vols=`docker volume ls -q --filter "label=com.docker.compose.project"`; \
					if [ -n "$$vols" ]; then docker volume rm $$vols 2>/dev/null; fi; \
					echo "Сношу сети проекта, а то останутся висеть как призраки..."; \
					nets=`docker network ls -q --filter "label=com.docker.compose.project"`; \
					if [ -n "$$nets" ]; then docker network rm $$nets 2>/dev/null; fi; \
					echo "Чищу pgdata-папки изнутри докера, иначе обычный rm на них обижается (root есть root)..."; \
					for pg in pgdata_meeting pgdata_request pgdata_user pgdata_chat; do \
						if [ -d "./backend/db/$$pg" ]; then \
							docker run --rm -v "$$(pwd)/backend/db/$$pg:/target" alpine sh -c "rm -rf /target/*" 2>/dev/null; \
						fi; \
					done; \
					echo "Удаляю файлы проекта с диска..."; \
					cd .. && rm -rf "$$OLDPWD"; \
					echo "Готово. Меня (и проекта) больше нет. Это было хорошее время вместе.\n Спасибо за всё. \nНе плачь, живи дальше...\nи прощай..."; \
				else \
					echo "А вот и нет;) Фух, пронесло. Самоуничтожение отменено."; \
					echo "Значит, в другой раз..."; \
				fi \
				;; \
			help) # команда помощи (увы, не психологической)\
				_help() { \
					echo "Чтобы запустить, пересоздать или перезапустить, пиши: init, up/all, switch, gateway, meeting/meetings, request, user, chat"; \
					echo "init/up/all/switch также запускают и frontend"; \
					echo "Чтобы остановить, пиши: switch, down/stop"; \
					echo "Есть ещё супер-слова: db и dev, они влияют на следующее слово, например make db up запустит только бдшки, а make dev meeting — только meeting, причём в режиме разработчика"; \
					echo "Есть ещё очистка кеша или дезинтеграция целого(ых) сервиса(ов): cache и remove .., например make remove cache или make remove db_meeting"; \
					echo "А ещё есть логи: make log или make db dev logs"; \
					echo "И конечно же БОЛЬШАЯ красная КНОПКА самоуничтожения - selfDestroySequence *)"; \
					echo "\n"; \
				}; \
				if [ "$$remove" = "true" ]; then \
					echo "Убрать помощь? Я лучше оставлю... на всякий случай"; \
					remove="false"; \
					db="false"; \
					dev="false"; \
				elif [ "$$db" = "false" ]; then \
    				_help; \
					if [ "$$dev" = "true" ]; then \
						echo "Доступные слова (с флагом dev): те же самые, но сервисы поднимутся через docker-compose-dev.yml — с пробросом портов и всем необходимым для разработки"; \
						dev="false"; \
					fi; \
				else \
    				_help; \
					if [ "$$dev" = "false" ]; then \
						echo "Доступные бдшки: meeting, request, user, chat (как отдельные слова, например 'make db meeting')"; \
						echo "'make db up' поднимет все бдшки разом"; \
					else \
						echo "Доступные бдшки (в dev режиме): те же, но через docker-compose-dev.yml и с открытыми портами"; \
						dev="false"; \
					fi; \
					db="false"; \
				fi \
				;; \
			init) # дошли до конца=вернулись к началу всего приложения, проверка зависимостей(docker и node.js) в системе, скачивание библиотек и запуск всего и вся \
				if [ "$$remove" = "true" ]; then \
					echo "Убрать начало непросто, но возможно вы хотите ВЗОРВАТЬ ВСЁ!?!??!?! => selfDestructSequence"; \
					remove="false"; \
				else \
					if [ "$$db" = "true" ]; then \
						echo "И fronend поднимем, и бд поднимем"; \
						db="false"; \
					fi; \
					command -v npm >/dev/null 2>&1 || { \
						echo "Не работает npm"; \
						echo "Плиииз поставь Node.js (в него уже встроен npm): https://nodejs.org/ или через nvm:"; \
						echo "  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash"; \
						echo "  nvm install --lts"; \
						exit 1; \
					}; \
					NODE_VER=$$(node -e "process.stdout.write(String(parseInt(process.versions.node)))" 2>/dev/null); \
					if [ -z "$$NODE_VER" ] || [ "$$NODE_VER" -lt 12 ]; then \
						echo "Node.js слишком старый (нашёл: $$(node -v 2>/dev/null || echo 'вообще непонятно что')), нужна версия 12+"; \
						echo "Обнови через nvm:"; \
						echo "  nvm install --lts && nvm use --lts"; \
						echo "Или скачай свежий с https://nodejs.org/"; \
						exit 1; \
					fi; \
					command -v docker >/dev/null 2>&1 || { \
						echo "А где docker? Не вижу в системе..."; \
						echo "Поставь Docker пж: https://docs.docker.com/get-docker/"; \
						exit 1; \
					}; \
					echo "npm и docker на месте, погнали"; \
					echo "Копирую .env'ы..."; \
					cp -n .env.example .env 2>/dev/null && echo "  .env создан" || echo "  .env уже есть, не трогаю"; \
					cp -n frontend/.env.example frontend/.env 2>/dev/null && echo "  frontend/.env создан" || echo "  frontend/.env уже есть, не трогаю"; \
					echo "Ставлю зависимости фронтенда..."; \
					(cd frontend && npm install && npm install mapbox-gl && npm install react-router-dom); \
					if [ "$$dev" = "true" ]; then \
						dev="false"; \
						echo "Поднимаю ВСЕ сервисы бэкенда в сервисном режиме"; \
						$(COMPOSE_ENV) up -d --build; \
					else \
						echo "Поднимаю ВСЕ сервисы бэкенда"; \
						$(COMPOSE) up -d --build; \
					fi; \
					echo "Запускаю фронтend"; \
					(cd frontend && npm run dev) \
				fi; \
				;; \
			*) #ну мало ли сюда как-то попадёт что-то кроме отфильтрованных слов\
				echo "Неизвестный сервис: $$word"; \
				;; \
		esac; \
		done;\
	fi;\
 # ВСЁ! Ты герой, ты прошёл этот файл вместе со мной!