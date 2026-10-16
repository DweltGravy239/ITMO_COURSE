Лаба 1 — Свой Docker
Этап 1. Сервис и запуск напрямую
Создадим свой сервис и запустим его. Сам сервис написал ИИ.
<img width="1677" height="1059" alt="image" src="https://github.com/user-attachments/assets/951e6243-b18c-48b8-8b32-916663ac9be7" />
Запустим наш сервис:
<img width="1238" height="172" alt="image" src="https://github.com/user-attachments/assets/8ba3fe63-2669-4663-b59d-b7f2af4780f5" />
Проверим, что все отрабатывает корректно:
<img width="1529" height="1238" alt="image" src="https://github.com/user-attachments/assets/a62796e8-87cf-4526-98a9-4a4b55058c89" />
Этап 2. Namespaces
Создадим namespace и посмотрит ps aux изнутри. 
<img width="1515" height="400" alt="image" src="https://github.com/user-attachments/assets/64348d6c-6603-4b5e-ab5f-20fc8476e075" />
Теперь попробуем посмотреть на наш namespace извне:
<img width="2225" height="256" alt="image" src="https://github.com/user-attachments/assets/c62abd43-e9dc-499f-ab6d-791081790326" />
Добавим сеть, DNS и разделение памяти в наш namespace и сравним, что теперь параметры отличаются от основного сервера:
<img width="2496" height="666" alt="image" src="https://github.com/user-attachments/assets/24de320b-f4d7-431c-9142-1e2d4d4a8658" />
Проверим, что все работает так, как мы думали с веб-приложением:
<img width="2300" height="1028" alt="image" src="https://github.com/user-attachments/assets/4e2af098-a1e7-4071-b58f-fd6f4cf6a9ec" />
Теперь запустим окружение, но не из под рута, а из под обычного пользователя:
<img width="2876" height="522" alt="image" src="https://github.com/user-attachments/assets/ebc35589-47fb-4bd8-a64c-c98b2a5dfa90" />
Этап 3. Cgroups
Теперь выделим лимиты на использование ресурсов нашей основной машины
Для того чтоб мы могли прописать параметры memory.max, cpu.max, pids.max нам нужно будет включить параметр cgroup.subtree_control:
Создадим cgroup для нашего namespace:
<img width="2600" height="662" alt="image" src="https://github.com/user-attachments/assets/18f89ce4-6dc8-4d92-b97b-400d30c163c3" />
Теперь выставим лимиты и попробуем проверить их работоспособность:
<img width="2588" height="1412" alt="image" src="https://github.com/user-attachments/assets/d9d78b0a-1187-48d9-a39e-6f9c0ac8b37f" />
Теперь проверим логи:
<img width="2850" height="400" alt="image" src="https://github.com/user-attachments/assets/47775eac-dbcc-4554-a3c7-b8e5ae03837e" />
Теперь выставим лимиты для CPU:
<img width="2104" height="1112" alt="image" src="https://github.com/user-attachments/assets/6cb68927-c198-48ef-b788-37d80dabc9da" />
Проверим, что все отрабатывает корректно:
<img width="1950" height="698" alt="image" src="https://github.com/user-attachments/assets/14002214-240b-449c-9322-01a47b92d6d5" />
Настроим максимальное количество потоков у нас в namespace и проверим результат:
<img width="2190" height="408" alt="image" src="https://github.com/user-attachments/assets/0534e68c-7ba5-4154-8b8e-a5a7c90f4388" />
Этап 4. Права capabilities и seccomp
Для начала проверим текущие права у рута:
<img width="2880" height="542" alt="image" src="https://github.com/user-attachments/assets/0140ab6e-5b66-407c-a4eb-011ca118452b" />
Теперь попробуем поиграться с правами CAP_SYS_TIME и проверим, что права от окружения могут отличаться:
<img width="2120" height="276" alt="image" src="https://github.com/user-attachments/assets/774d50c7-3f7a-4c3e-a0d8-52d5a6704868" />
Проверим, какие системы вызовы разрешены процессу:
<img width="1216" height="434" alt="image" src="https://github.com/user-attachments/assets/9452d786-e1f8-4b95-941c-51541d2fcc40" />
Проверим, что докер автоматически настраивает seccomp в режиме BPF-фильтра:
<img width="2164" height="532" alt="image" src="https://github.com/user-attachments/assets/7c57d441-e95f-4b3d-89cb-5a59efbf88af" />
Проверим, какое название syscall у нас:
<img width="1640" height="114" alt="image" src="https://github.com/user-attachments/assets/b4da9baf-d0c2-4a37-bf21-56a0504ac094" />
Теперь приступим к написанию своего правила для режима BPF-фильтра:
<img width="1546" height="334" alt="image" src="https://github.com/user-attachments/assets/345fbdfd-a4f8-4346-9023-af0ec369b2cf" />
Проверим, что все работает:
<img width="2404" height="216" alt="image" src="https://github.com/user-attachments/assets/612f8740-2d98-4bd8-a8bd-0f66be0cbc41" />
Этап 5. mydocker.sh и сравнение с docker run
Настало время создания скрипта mydocker.sh, для того чтоб в дальнейшем сравнить его с docker run
Наш скрипт будет состоять из 7 этапов:
- Создать cgroup
- Записать себя в cgroup.procs
- Запустить unshare
- Выставить hostname
- Сбросить capabilities
- Наложить seccomp (сделаем на последнем шаге)
- Запустить веб приложение
Попробуем выполнить первый этап, создать cgroup и записать себя в cgroup.procs:
<img width="1384" height="800" alt="image" src="https://github.com/user-attachments/assets/50406e30-1e8f-48c7-a844-b575cb7ba120" />
Для теста добавим проверку файлов memory.max и pids.max и запустим скрипт:
<img width="1340" height="446" alt="image" src="https://github.com/user-attachments/assets/a694cdae-26c0-4fd4-9582-74add28d815c" />
Добавим namespace и проверим, что все отрабатывает корректно:
<img width="1784" height="1298" alt="image" src="https://github.com/user-attachments/assets/ccd8091d-d2e1-4cc4-99f9-f1d45d030f37" />
Теперь сбросим capabilities и проверим возможность выполнения команд:
<img width="2520" height="1156" alt="image" src="https://github.com/user-attachments/assets/be53796f-47f3-4b57-8564-270aab6c11a2" />
Теперь запустим веб приложение
<img width="2860" height="1222" alt="image" src="https://github.com/user-attachments/assets/93fd568a-fc26-44b2-881a-ef85f33ae16a" />
Осталось только обернуть все в seccomp и процесс создания скрипта закончен:
Скачаем питоновский модуль для создания обретки:
<img width="1312" height="138" alt="image" src="https://github.com/user-attachments/assets/979121ca-b72e-4033-b1d9-5a6a679ebd0b" />
Напишем обертку и проверим, что она отрабатывает. На mkdir выдает ошибку, а на другую команду, все хорошо:
<img width="2712" height="544" alt="image" src="https://github.com/user-attachments/assets/9044aafa-527e-434c-a143-f7d93962cb03" />
Теперь скорректируем наш скрипт, добавив нужные параметры, после чего проверим работоспособность:
<img width="2846" height="1240" alt="image" src="https://github.com/user-attachments/assets/251c7847-dc04-41c1-b72b-932aa5be25c3" />
Теперь соберем образ через docker:
<img width="2846" height="1056" alt="image" src="https://github.com/user-attachments/assets/6413c37f-54f0-4f25-9e7b-dba90cb710c3" />
Из сравнения двух вариантов, мы можем понять, что докер дополняет сам контейнер, а именно:
- Создает дополнительный интерфейс и выдает ип адрес
- Собирает файловую систему из образа, которая видна только в контейнере
- Сам отвечает за безопасность и выдаем минимальные права
- Для того чтоб останосить, можно использовать команды docker stop, start и т.д., а скриптом надо было бы дописать в скрипт.
Этап 6. Образы
Напишем свой докерфайл и соберем свой образ:
<img width="1638" height="1226" alt="image" src="https://github.com/user-attachments/assets/48439143-a597-478f-a023-2f01ee6131ff" />
Теперь сделаем мульти-стейдж сборку с минимальной базой:
<img width="2880" height="1264" alt="image" src="https://github.com/user-attachments/assets/9db0427a-d6ea-43a0-aaf7-937d3810e7fc" />
Теперь сравним число слоев и что переиспользователось в кэше при повторной сборке:
<img width="1492" height="676" alt="image" src="https://github.com/user-attachments/assets/e24bfdfa-55a9-4065-b182-1e03e0f7364e" />
На скриншоте мы видим, что первый образ имеет 11 слоев и в сумме весит около 500 МБ, это все потому что мы ставили толстый образ golang. В мулти-стейдж сборке один слой с данными - бинарник. Все что нужно для сборки (компилятор, кэш, исходники), остались в первой стадии и в итоге в образ не попали. Это и есть смысл мульти-стейдж.
Теперь попробуем записать какой-то файл в контейнер и пересоздать его:
<img width="1446" height="246" alt="image" src="https://github.com/user-attachments/assets/158d816a-fcd3-4c08-bab5-afbbd034c36e" />
Теперь попробуем провести подобный эксперемент, но только используя том:
<img width="1428" height="522" alt="image" src="https://github.com/user-attachments/assets/49a82b77-9dc3-4cab-83ab-529be7aaeb79" />
Мы можем наблюдать, что топ создался у нас локально и поэтому при создание/удаление контейнера у нас никаких проблем не возникло и данные сохранились. 
Этап 7. gVisor
Скачиваем gVisor:
<img width="2878" height="838" alt="image" src="https://github.com/user-attachments/assets/4a28f39a-9fc9-4cdf-b031-906490cc2d11" />
Распакуем архив и проверим, что все отрабатываем корректно и установка прошла успешно:
<img width="1202" height="340" alt="image" src="https://github.com/user-attachments/assets/2151f506-56da-4b8a-bad5-f154e48f8ff1" />
Зарегистрируем рантайм:
<img width="708" height="406" alt="image" src="https://github.com/user-attachments/assets/6f0bdcee-f4ab-4ebd-8af8-858c5e5f0a3d" />
Теперь проверим, что gVisor отрабатывает корректно:
<img width="1076" height="550" alt="image" src="https://github.com/user-attachments/assets/4100c9b9-e4b1-419b-b935-5da1ccb6ce6d" />

Этап 8. Мониторинг
