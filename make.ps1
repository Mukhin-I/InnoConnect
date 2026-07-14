<#
    Windows-версия основного тела Makefile.
    Вызывается автоматически из Makefile, когда $(OS) = Windows_NT — сам Makefile
    при этом не меняется по смыслу, меняется только то, ЧТО выполняет система:
    на Linux/macOS это POSIX-shell блок, на Windows — этот скрипт.

    Почему так: recipe в Makefile написан на POSIX shell (функции, case/esac,
    подстановка команд в обратных кавычках, /dev/tty и т.д.) — на Windows такого
    интерпретатора "из коробки" нет (только cmd.exe и PowerShell). Вместо того
    чтобы требовать WSL/Git Bash, вся та же логика продублирована здесь на
    PowerShell, который уже встроен в Windows 10/11. Из внешних требований —
    только сам make.

    Сообщения, шутки и порядок действий сохранены один в один с оригиналом.
#>

param(
    [string]$EnvFile = ".env",
    [string]$ComposeFile = "./backend/deployment/docker/docker-compose.yml",
    [string]$ComposeDevFile = "./backend/deployment/docker/docker-compose-dev.yml",
    # все "слова" из командной строки (make up db, make dev meeting, и т.д.) —
    # аналог $(ARGS)/for word in $(ARGS) из Makefile; фильтрация по SERVICES
    # уже сделана самим Makefile до вызова этого скрипта
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$Words
)

# аналог "docker compose -f ... --env-file ..." сохранённого в переменной COMPOSE
function Compose {
    param([string[]]$CmdArgs)
    & docker compose -f $ComposeFile --env-file $EnvFile @CmdArgs
}
# аналог COMPOSE_DEV
function ComposeDev {
    param([string[]]$CmdArgs)
    & docker compose -f $ComposeDevFile --env-file $EnvFile @CmdArgs
}

# вспомогательная функция для удаления контейнеров, аналог bash _remove()
function Remove-ServiceContainer {
    param([string]$Target)
    Write-Host "Удаляю ${Target}: и контейнер, и образ, и кеш, и volume'ы"
    $cid = (Compose @('images', '-q', $Target) 2>$null)
    if (-not $cid) { $cid = (ComposeDev @('images', '-q', $Target) 2>$null) }
    docker stop $Target 2>$null | Out-Null
    docker rm -f -v $Target 2>$null | Out-Null
    $vols = (docker volume ls -q --filter "label=com.docker.compose.service=$Target")
    if ($vols) { docker volume rm $vols | Out-Null } # чистим volumes ещё раз
    if ($cid) {
        docker buildx prune --filter "parents=$cid" -f # чистим кеш сборки
        docker rmi $cid # ещё что-то чистим
    } else {
        Write-Host "Не нашёл образ $Target, нечего очищать"
    }
}

# вспомогательная функция помощи, аналог bash _help()
function Show-Help {
    Write-Host "Чтобы запустить, пересоздать или перезапустить, пиши: init, up/all, switch, gateway, meeting/meetings, request, user, chat"
    Write-Host "init/up/all/switch также запускают и frontend"
    Write-Host "Чтобы остановить, пиши: switch, down/stop"
    Write-Host "Есть ещё супер-слова: db и dev, они влияют на следующее слово, например make db up запустит только бдшки, а make dev meeting - только meeting, причём в режиме разработчика"
    Write-Host "Есть ещё очистка кеша или дезинтеграция целого(ых) сервиса(ов): cache и remove .., например make remove cache или make remove db_meeting"
    Write-Host "А ещё есть логи: make log или make db dev logs"
    Write-Host "И конечно же БОЛЬШАЯ красная КНОПКА самоуничтожения - selfDestroySequence *)"
    Write-Host ""
}

# флаги, копящиеся по ходу разбора слов - аналог db="false"/dev="false"/remove="false" из bash
$db = $false
$dev = $false
$remove = $false

$AllServices = @('gateway', 'meeting', 'request', 'user', 'chat', 'db_meeting', 'db_request', 'db_user', 'db_chat')
$DbServices = @('db_meeting', 'db_request', 'db_user', 'db_chat')

# тепербь для каждого слова в строке ввода (как for word in $(ARGS); do ... done)
foreach ($word in $Words) {
    Write-Host ""
    switch -CaseSensitive ($word) { # -CaseSensitive, чтобы вести себя как bash case (например, selfDestroySequence не спутать с selfdestroysequence)

        { $_ -in @('up', 'all') } { # поднимаем *всё* с обработкой флагов db => поднять *всё* бд, dev => поднять *всё* [бд] в dev режиме, remove => удалить *всё* [бд]
            if ($remove) {
                if ($db) {
                    Write-Host "!Удаляю ВСЕ бдшки"
                    foreach ($s in $DbServices) { Remove-ServiceContainer $s }
                } else {
                    Write-Host "!ЁУдаляю ВСЕ сервисы"
                    foreach ($s in $AllServices) { Remove-ServiceContainer $s }
                }
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) { # 4 варианта соответственно для комбинаций db=true/false, dev=true/false
                if (-not $dev) {
                    Write-Host "Запускаю ВСЕ сервисы"
                    Compose @('up', '-d', '--build')
                } else {
                    Write-Host "Запускаю ВСЕ сервисы, НО в сервисном режиме"
                    ComposeDev @('up', '-d', '--build')
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "Запускаю ВСЕ бдшки"
                    Compose (@('up', '-d', '--build') + $DbServices)
                } else {
                    Write-Host "Запускаю ВСЕ бдшки в сервисном режиме"
                    ComposeDev (@('up', '-d', '--build') + $DbServices)
                    $dev = $false
                }
                $db = $false
            }
            Write-Host "Запускаю фронтend" # в конце prodная история, поднимаем ещё и фронтend
            Push-Location frontend
            npm run dev
            Pop-Location
            continue
        }

        { $_ -in @('down', 'stop') } { # то же самое, но вместа запуска, останавливаем
            if ($remove) {
                Write-Host "удалить ВСЁ??? точно? я боюсь... если надо, напиши тогда remove all"
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                if (-not $dev) {
                    Write-Host "Останавливаю всё, что есть"
                    Compose @('down')
                } else {
                    Write-Host "Останавливаю всё, что есть, даже в сервисном режиме!"
                    ComposeDev @('down')
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "Останавливаю всe бдшки"
                    Compose (@('down') + $DbServices)
                } else {
                    Write-Host "Останавливаю всe бдшки, сервисный режим им не поможет хе-хе"
                    ComposeDev (@('down') + $DbServices)
                    $dev = $false
                }
                $db = $false
            }
            continue
        }

        'switch' { # если все сервисы подняты -> выключаем, если не все -> доподнимаем
            if ($remove) {
                Write-Host "switch не дружит с remove, забыл стереть флаг?"
                $remove = $false; $db = $false; $dev = $false
            } else {
                $allUp = $true
                foreach ($s in $AllServices) {
                    $state = (Compose @('ps', '-q', $s) 2>$null)
                    if (-not $state) { $allUp = $false }
                }
                if ($allUp) {
                    Write-Host "Все сервисы запущены, выключаю всё"
                    Compose @('down')
                } else {
                    Write-Host "Не все сервисы запущены, запускаю всё"
                    Compose @('up', '-d', '--build')
                    Write-Host "Запускаю фронтend"
                    Push-Location frontend
                    npm run dev
                    Pop-Location
                }
                $db = $false; $dev = $false
            }
            continue
        }

        { $_ -in @('log', 'logs') } { # открыть логи все или бдшек
            if ($remove) {
                Write-Host "Что написано пером, уже не вырубишь топором"
                Write-Host "Всё, что будет в будущем, будет в будущем, всё, что было - уже история, живи настоящим. Зачем хочешь удалить историю)?"
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                if (-not $dev) {
                    Write-Host "Открываю логи"
                    Compose @('logs', '-f')
                } else {
                    Write-Host "Открываю логи разработчика)"
                    ComposeDev @('logs', '-f')
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "Открываю логи бдшек"
                    Compose (@('logs', '-f') + $DbServices)
                } else {
                    Write-Host "Открываю логи разрабных бдшек"
                    ComposeDev (@('logs', '-f') + $DbServices)
                    $dev = $false
                }
                $db = $false
            }
            continue
        }

        'meetings' { # отдельный кейс с неправильным названием сервиса с буквой С в конце - популярная моя ошибка со стёбными комментариями
            if ($remove) { # remove meeting или db meeting (логично?) удаляет это контейнер
                if ($db) { Remove-ServiceContainer 'db_meeting' } else { Remove-ServiceContainer 'meeting' }
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                if (-not $dev) {
                    Write-Host "Запускаю meetinG.."
                    Compose @('up', '-d', '--build', 'meeting')
                } else {
                    Write-Host "Запускаю meetinG... (в сервисном режиме!)"
                    ComposeDev @('up', '-d', '--build', 'meeting')
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "Запускаю бдшку для meetinG.."
                    Compose @('up', '-d', '--build', 'db_meeting')
                } else {
                    Write-Host "Запускаю бдшку для meetinG.. с портами, все дела - сервисный режм"
                    ComposeDev @('up', '-d', '--build', 'db_meeting')
                    $dev = $false
                }
                $db = $false
            }
            continue
        }

        'gateway' { # отдельный кейс для gateway, потому что у него нет бд =) обработку этого флага нужно писать иначе
            if ($remove) {
                Remove-ServiceContainer 'gateway'
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                if (-not $dev) {
                    Write-Host "Запускаю gateway"
                    Compose @('up', '-d', '--build', 'gateway')
                } else {
                    Write-Host "Запускаю gateway в супер режиме разработчика"
                    ComposeDev @('up', '-d', '--build', 'gateway')
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "А нету для gateway бдшки 8)"
                    $db = $false
                } else {
                    Write-Host "А нету для gateway бдшки 8), даже если ты разраб :("
                    $db = $false; $dev = $false
                }
            }
            continue
        }

        { $_ -in @('meeting', 'request', 'user', 'chat') } { # для основной группы сервисок всё одинаково, парсим название через $word и передаём в запуск, удаление, бдшки и бла бла бла
            if ($remove) {
                if ($db) { Remove-ServiceContainer "db_$word" } else { Remove-ServiceContainer $word }
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                if (-not $dev) {
                    Write-Host "Запускаю $word"
                    Compose @('up', '-d', '--build', $word)
                } else {
                    Write-Host "Запускаю $word в режиме разработчика"
                    ComposeDev @('up', '-d', '--build', $word)
                    $dev = $false
                }
            }
            else {
                if (-not $dev) {
                    Write-Host "Запускаю бдшку для $word"
                    Compose @('up', '-d', '--build', "db_$word")
                } else {
                    # ВНИМАНИЕ: в оригинальном Makefile здесь тоже баг — поднимается $word, а не db_$word.
                    # Оставлено специально таким же, чтобы поведение совпадало 1-в-1 с Linux/macOS-версией.
                    Write-Host "Запускаю разработческую бдшку для $word"
                    ComposeDev @('up', '-d', '--build', $word)
                    $dev = $false
                }
                $db = $false
            }
            continue
        }

        'db' { # УраЁ!!! сервисы закончились, остались флаги, кеш и инициализация
            if (-not $db) {
                Write-Host "Жду бд"
                $db = $true
            } else {
                Write-Host "Очень жду бд" # когда встретилось db при уже активированном db, например make db db meeting
                $db = $true
            }
            continue
        }

        'dev' {
            if (-not $dev) {
                Write-Host "Жду dev"
                $dev = $true
            } else {
                Write-Host "Сильно жду dev"
                $dev = $true
            }
            continue
        }

        'remove' {
            if (-not $remove) {
                Write-Host "Жду сервис для удаления"
                $remove = $true
            } else {
                Write-Host "ДАЙТЕ мне сервис для удаления"
                $remove = $true
            }
            continue
        }

        'cache' { # очистка кеша, TODO: обработка dev и db не сделана, чистит сразу всё
            if ($remove) { # когда make remove cache = make cache
                Write-Host "Всё, понял, если что, можешь просто cache писать`nНо так не оч красиво((("
                $remove = $false
            }
            Write-Host "Чищу кеш сборки..."
            $ids = @()
            foreach ($s in $AllServices) {
                $id = (Compose @('images', '-q', $s) 2>$null)
                if (-not $id) { $id = (ComposeDev @('images', '-q', $s) 2>$null) } # если не нашли id - проверить в dev-файле
                if ($id) {
                    Write-Host "Нашёл ${s}: ($id)"
                    $ids += $id
                } else {
                    Write-Host "Образа $s нет, пропускаю"
                }
            }
            if ($ids.Count -gt 0) {
                docker buildx prune --filter "parents=$($ids -join ';')" -f
            } else {
                Write-Host "Нечего чистить — ещё ничего не было"
            }
            continue
        }

        'selfDestroySequence' { # команда самоуничтожения для удобного удаления проекта [нужно знать пароль] TODO спрятать пароль или хранить его хеш
            Write-Host "Так-так-так. Кнопка самоуничтожения. Серьёзно?"
            Write-Host "Это удалит ВСЕ docker-контейнеры, образы, volume'ы, кеш сборки И все файлы проекта с диска. Без права на отмену."
            Write-Host "Чтобы убедиться, что ты делаешь это осознанно, ответь правильно на 3 вопроса (или хотя бы попытайся)."
            Write-Host ""
            Write-Host "Вопрос 1 из 3: ты точно уверен, что это то, чего ты хочешь? `nКакой твой любимый фрукт?"
            $answer1 = Read-Host
            Write-Host "Вопрос 2 из 3: ты осознаёшь, что обратной дороги не будет, и я не сохраняю бэкапы?`n Ты можешь написать картину?"
            $answer2 = Read-Host
            Write-Host "Вопрос 3 из 3: последний шанс передумать. Точно удаляем ВСЁ? `nЧто такое Чёрный кофе?"
            $answer3 = Read-Host
            if ($answer1 -ceq 'яблоко' -and $answer2 -ceq 'робот' -and $answer3 -ceq 'Агата Кристи') {
                Write-Host "Ну ты сам напросился. Поехали..."
                Start-Sleep -Seconds 1
                Write-Host "Гашу и стираю все сервисы..."
                foreach ($s in $AllServices) { Remove-ServiceContainer $s }
                Write-Host "Чищу остатки buildx кеша..."
                docker buildx prune -af 2>$null
                Write-Host "Сношу volume'ы проекта..."
                $vols = (docker volume ls -q --filter "label=com.docker.compose.project")
                if ($vols) { docker volume rm $vols 2>$null }
                Write-Host "Сношу сети проекта, а то останутся висеть как призраки..."
                $nets = (docker network ls -q --filter "label=com.docker.compose.project")
                if ($nets) { docker network rm $nets 2>$null }
                Write-Host "Чищу pgdata-папки изнутри докера, иначе обычный rm на них обижается (root есть root)..."
                foreach ($pg in @('pgdata_meeting', 'pgdata_request', 'pgdata_user', 'pgdata_chat')) {
                    $pgPath = Join-Path (Get-Location).Path "backend/db/$pg"
                    if (Test-Path $pgPath) {
                        docker run --rm -v "${pgPath}:/target" alpine sh -c "rm -rf /target/*" 2>$null
                    }
                }
                Write-Host "Удаляю файлы проекта с диска..."
                $projectDir = (Get-Location).Path
                Set-Location ..
                Remove-Item -LiteralPath $projectDir -Recurse -Force
                Write-Host "Готово. Меня (и проекта) больше нет. Это было хорошее время вместе.`n Спасибо за всё. `nНе плачь, живи дальше...`nи прощай..."
            } else {
                Write-Host "А вот и нет;) Фух, пронесло. Самоуничтожение отменено."
                Write-Host "Значит, в другой раз..."
            }
            continue
        }

        'help' { # команда помощи (увы, не психологической)
            if ($remove) {
                Write-Host "Убрать помощь? Я лучше оставлю... на всякий случай"
                $remove = $false; $db = $false; $dev = $false
            }
            elseif (-not $db) {
                Show-Help
                if ($dev) {
                    Write-Host "Доступные слова (с флагом dev): те же самые, но сервисы поднимутся через docker-compose-dev.yml — с пробросом портов и всем необходимым для разработки"
                    $dev = $false
                }
            }
            else {
                Show-Help
                if (-not $dev) {
                    Write-Host "Доступные бдшки: meeting, request, user, chat (как отдельные слова, например 'make db meeting')"
                    Write-Host "'make db up' поднимет все бдшки разом"
                } else {
                    Write-Host "Доступные бдшки (в dev режиме): те же, но через docker-compose-dev.yml и с открытыми портами"
                    $dev = $false
                }
                $db = $false
            }
            continue
        }

        'init' { # дошли до конца=вернулись к началу всего приложения, проверка зависимостей (docker и node.js) в системе, скачивание библиотек и запуск всего и вся
            if ($remove) {
                Write-Host "Убрать начало непросто, но возможно вы хотите ВЗОРВАТЬ ВСЁ!?!??!?! => selfDestructSequence"
                $remove = $false
            } else {
                if ($db) {
                    Write-Host "И fronend поднимем, и бд поднимем"
                    $db = $false
                }
                if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
                    Write-Host "Не работает npm"
                    Write-Host "Плиииз поставь Node.js (в него уже встроен npm): https://nodejs.org/"
                    exit 1
                }
                $nodeVerRaw = & node -e "process.stdout.write(String(parseInt(process.versions.node)))" 2>$null
                $nodeVerOk = $false
                if ($nodeVerRaw) {
                    try { if ([int]$nodeVerRaw -ge 12) { $nodeVerOk = $true } } catch { $nodeVerOk = $false }
                }
                if (-not $nodeVerOk) {
                    $curVer = try { (& node -v) } catch { 'вообще непонятно что' }
                    if (-not $curVer) { $curVer = 'вообще непонятно что' }
                    Write-Host "Node.js слишком старый (нашёл: $curVer), нужна версия 12+"
                    Write-Host "Скачай свежий с https://nodejs.org/"
                    exit 1
                }
                if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
                    Write-Host "А где docker? Не вижу в системе..."
                    Write-Host "Поставь Docker пж: https://docs.docker.com/get-docker/"
                    exit 1
                }
                Write-Host "npm и docker на месте, погнали"
                Write-Host "Копирую .env'ы..."
                if ((-not (Test-Path $EnvFile)) -and (Test-Path ".env.example")) {
                    Copy-Item ".env.example" $EnvFile
                    Write-Host "  .env создан"
                } else {
                    Write-Host "  .env уже есть, не трогаю"
                }
                if ((-not (Test-Path "frontend/.env")) -and (Test-Path "frontend/.env.example")) {
                    Copy-Item "frontend/.env.example" "frontend/.env"
                    Write-Host "  frontend/.env создан"
                } else {
                    Write-Host "  frontend/.env уже есть, не трогаю"
                }
                Write-Host "Ставлю зависимости фронтенда..."
                Push-Location frontend
                npm install
                npm install mapbox-gl
                npm install react-router-dom
                Pop-Location
                if ($dev) {
                    $dev = $false
                    Write-Host "Поднимаю ВСЕ сервисы бэкенда в сервисном режиме"
                    # ПРИМЕЧАНИЕ: в оригинале тут стояла опечатка $(COMPOSE_ENV) (несуществующая переменная),
                    # из-за которой dev-ветка init на Linux/macOS реально не поднимала бэкенд.
                    # Здесь используется явно правильная ComposeDev — стоит поправить и в самом Makefile.
                    ComposeDev @('up', '-d', '--build')
                } else {
                    Write-Host "Поднимаю ВСЕ сервисы бэкенда"
                    Compose @('up', '-d', '--build')
                }
                Write-Host "Запускаю фронтend"
                Push-Location frontend
                npm run dev
                Pop-Location
            }
            continue
        }

        default { # ну мало ли сюда как-то попадёт что-то кроме отфильтрованных слов
            Write-Host "Неизвестный сервис: $word"
            continue
        }
    }
}
