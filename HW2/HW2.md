# Лабораторная работа №2 - Observability

## Часть 0 - Создание своего сервиса.

Для создания своего сервиса я воспользовался клодом. В итоге он мне сгенерировал что-то такое (более подробнее можно увидеть в каталоге `/app`):

<img width="1777" height="1021" alt="image" src="https://github.com/user-attachments/assets/c5c95b1d-c167-4d78-bdd4-a4c3ad1c855e" />

После создания приложения я написал базовый докерфайл и сбилдил его:

<img width="2877" height="451" alt="image" src="https://github.com/user-attachments/assets/547cc9f8-97ba-4125-a8c4-1a9af7940541" />

А после этого прописал его в kind:

<img width="2878" height="192" alt="image" src="https://github.com/user-attachments/assets/59ff7950-4451-426d-bfdb-60d08c1f32d1" />

Сначала я решил набросать обычный манифест для того, чтоб проверить его корректность, а потом перевести его в helm:

<img width="1084" height="1334" alt="image" src="https://github.com/user-attachments/assets/7d6e10a2-6c44-4471-878b-f45ffce73d46" />

Дальше мы сделаем `kubectl apply` и проверим, что он отрабатываем корректно:

<img width="1350" height="516" alt="image" src="https://github.com/user-attachments/assets/9fce6480-ef75-41ba-a673-dea973ac5ee9" />

Проверим, что веб-приложение отрабатываем корректно:

<img width="1752" height="678" alt="image" src="https://github.com/user-attachments/assets/a2d26048-568b-4b4a-841e-45325f9dc456" />

Когда мы выяснили, что мой манифест отрабатывает корректно, можно начать писать helm:

<img width="1500" height="1354" alt="image" src="https://github.com/user-attachments/assets/2379b1c3-c9d4-4f86-95b4-420124644e4e" />

Теперь попробуем применить наш чарт:

<img width="1206" height="406" alt="image" src="https://github.com/user-attachments/assets/7aa5b57f-5ee7-4f27-b56d-85258640a646" />

Проверяем, что наш helm применился:

<img width="1388" height="554" alt="image" src="https://github.com/user-attachments/assets/a4ab3838-2b4c-4408-be0c-177fb8ce551c" />

## Часть 1. Установка Prometheus и Grafana

Теперь переходим к установке и настройке prometheus. Скачиваем helm с prometheus

<img width="2850" height="1374" alt="image" src="https://github.com/user-attachments/assets/6ec54460-5e55-44f6-909f-d82e50ea9ee7" />

Проверим, что все поднялось успешно

<img width="1704" height="504" alt="image" src="https://github.com/user-attachments/assets/9bc58735-bf28-4c0b-9e57-5a24fd3fdf2e" />

Теперь настроим отправку логов на prometheus с нашего веб приложения:

<img width="1520" height="1068" alt="image" src="https://github.com/user-attachments/assets/3ee9d614-e244-4ea3-9db9-f6bfce7df145" />

Применяем изменения:

<img width="1726" height="406" alt="image" src="https://github.com/user-attachments/assets/623d81b2-2268-4b23-b682-b6fcdd67c8d4" />

Проверим, что все работает:

<img width="2880" height="806" alt="image" src="https://github.com/user-attachments/assets/50c3de4f-bba9-4dcb-a712-8d5b7e5231b5" />

Теперь ставим графану, также через helm:

<img width="2092" height="322" alt="image" src="https://github.com/user-attachments/assets/5c9b23ea-1cc1-4c27-97d1-824b6348d932" />

Установим графану:

<img width="2866" height="1100" alt="image" src="https://github.com/user-attachments/assets/e2b28730-5d8e-4e6a-892b-0f32817333d7" />

Проверим, что все работает корректно:

<img width="1882" height="504" alt="image" src="https://github.com/user-attachments/assets/9675797f-1b51-4215-aef2-95afc073f25c" />

Вытаскиваем пароль от админской учетки и пробуем зайти в графану:

<img width="2880" height="1024" alt="image" src="https://github.com/user-attachments/assets/80aee688-1603-488e-9d05-c9352592a681" />

Для того, чтоб каждый раз не прописывать источник в графане, пропишет его сразу же в конфигурациях:

<img width="1220" height="466" alt="image" src="https://github.com/user-attachments/assets/e44241dc-c02e-4916-9e75-1e673235ca77" />

Применим изменения:

<img width="2868" height="1070" alt="image" src="https://github.com/user-attachments/assets/bec3962c-7995-4060-b98d-8933e9e43667" />

Проверим в графане, что все отработало корректно:

<img width="2260" height="730" alt="image" src="https://github.com/user-attachments/assets/dab9b459-c3ef-427f-b0f0-65d092f16632" />

Проверяем, что графики отображаются корректно(для этого построим график)

<img width="2200" height="1120" alt="image" src="https://github.com/user-attachments/assets/42aa34b7-4fcc-4aa7-8c58-82a95dce1e44" />

Теперь сохраняем наш дашборд

<img width="1482" height="1316" alt="image" src="https://github.com/user-attachments/assets/13b2a003-491d-4bef-afc6-aa96e5ffcef5" />

Загрузим через kubectl:

<img width="2224" height="192" alt="image" src="https://github.com/user-attachments/assets/9942fc06-030c-44bc-8070-73da7f0a435d" />

Теперь пропишем в values, чтоб сразу же загружать информацию о дашборде в графану:

<img width="1084" height="924" alt="image" src="https://github.com/user-attachments/assets/e0e103ed-a94b-400f-899f-8c287d347b3c" />

Сохраним изменения:

<img width="2860" height="1168" alt="image" src="https://github.com/user-attachments/assets/c63866d0-32c3-45ca-a4a1-55ae88594e25" />

Проверим, что графана сразу же отображает дашборд:

<img width="2262" height="1058" alt="image" src="https://github.com/user-attachments/assets/ff3c246a-326d-44d2-aa8c-a6e2aa063091" />

## Часть 2. Установка Loki

Для установки локи настраиваем values для loki:

<img width="610" height="1370" alt="image" src="https://github.com/user-attachments/assets/c918e96e-64e0-4700-a0d3-79411a9a7538" />

После этого устанавливаем helm (его также возьмем из онлайн репозитория):

<img width="2000" height="1014" alt="image" src="https://github.com/user-attachments/assets/35213781-2375-48d0-a008-69919362407b" />

Проверим, что поды поднялись успешно:

<img width="1806" height="546" alt="image" src="https://github.com/user-attachments/assets/0340d8f6-19ae-4bf3-8bde-ebd1862368fc" />

Теперь начинаем ставить агентов, для них будем использовать grafana alloy.

Для того чтоб проверить, что все поднимается и отправляется корректно, напишем неполный values и попробуем запустить:

<img width="1908" height="334" alt="image" src="https://github.com/user-attachments/assets/8391f5e8-4e80-4fd8-a061-bcf34ed0904e" />

Установим helm:

<img width="2050" height="382" alt="image" src="https://github.com/user-attachments/assets/534a8420-c882-4a96-903e-e4fbe26c5f99" />

Проверим, что все отрабатывает корректно:

<img width="1808" height="550" alt="image" src="https://github.com/user-attachments/assets/a013aea9-db6f-446b-a943-a1ba3288bcee" />

Теперь, когда мы убедились, что все отрабатывает корректно, мы добавим источник в наш конфиг:

<img width="1914" height="1386" alt="image" src="https://github.com/user-attachments/assets/a17bda5b-c1e1-4e3c-949d-0ea95e0be5da" />

Применим изменения:

<img width="1990" height="980" alt="image" src="https://github.com/user-attachments/assets/238553e9-b7f3-42ba-ab4d-2f6bc8d93326" />

Проверим, что логи нормальные:

<img width="2878" height="1202" alt="image" src="https://github.com/user-attachments/assets/e28ec107-6019-47b6-bf8c-7cea6ce8a474" />

Теперь добавим его как источник в графану и посмотрим, все ли он корректно выводит:

<img width="2194" height="1028" alt="image" src="https://github.com/user-attachments/assets/ac6b9ab7-921a-4e64-aa57-ede9574d67ed" />

<img width="2148" height="1172" alt="image" src="https://github.com/user-attachments/assets/e524a019-649d-44b5-a240-a8a145242937" />

Осталось только отредактировать конфиг графаны, чтоб источник мы добавляли не руками, а он сразу был доступен:

<img width="1152" height="1180" alt="image" src="https://github.com/user-attachments/assets/5e32fbdc-2dd8-40f8-887f-f239198b5636" />

Обновим helm:

<img width="2182" height="378" alt="image" src="https://github.com/user-attachments/assets/39eed7df-bb36-4540-b203-9f98873151a0" />

## Часть 3. Установка и настройка Jaeger

Добавляем репозиторий jaeger и смотрим, какие версии там есть

<img width="2068" height="580" alt="image" src="https://github.com/user-attachments/assets/598350f0-83a2-4dcd-b009-7bfd766d6d76" />

Установим jaeger:

<img width="2862" height="950" alt="image" src="https://github.com/user-attachments/assets/3edb08f4-fb43-447f-9f29-9819c042c48b" />

Скорректируем наш helm приложения, чтоб отправлять логи сразу же джагеру:

<img width="1518" height="908" alt="image" src="https://github.com/user-attachments/assets/a1d45ccb-be5a-43e5-90cf-d7955644f4d1" />

Обновим хелм приложения:

<img width="1314" height="772" alt="image" src="https://github.com/user-attachments/assets/8787cbbc-9fbc-489a-8a63-dd39f030bc16" />

Проверим, что веб морда работает:

<img width="1726" height="616" alt="image" src="https://github.com/user-attachments/assets/d4fe6b06-7336-42b1-806a-31418d298177" />

<img width="2870" height="972" alt="image" src="https://github.com/user-attachments/assets/0714762f-8b04-4831-a610-7c86cc402db4" />

Для проверки работоспособности, нагенерим трафик и проверим, появился ли он в jaeger:

<img width="2880" height="346" alt="image" src="https://github.com/user-attachments/assets/68cd3abb-97ff-4486-85b4-c4aff81fefa4" />

<img width="2130" height="446" alt="image" src="https://github.com/user-attachments/assets/8137ac55-d8ed-417f-861c-b9b537c6ce62" />

Теперь проверим, что каждый статус нашего приложения отображается корректно:

### slow:

<img width="2870" height="914" alt="image" src="https://github.com/user-attachments/assets/d0fa505c-9a9a-48fd-ab1d-e01b5c8d9f92" />

### error:

<img width="2872" height="778" alt="image" src="https://github.com/user-attachments/assets/0cf19d64-98d0-48da-ac08-ad9c9f4e6ec6" />

Теперь нам необходимо связать логи и трейсы. Делать мы это будем в графане:

Проверим руками, существует ли какое-то событие в Loki, а потом поищем это событие в jaeger:

<img width="2878" height="1284" alt="image" src="https://github.com/user-attachments/assets/cb9d137d-f49c-4b9f-afc1-c2bcae33a1af" />

<img width="2870" height="730" alt="image" src="https://github.com/user-attachments/assets/e7852a66-24a8-40d2-905e-63a911275160" />

<img width="2870" height="730" alt="image" src="https://github.com/user-attachments/assets/07bed5ea-b73d-4b04-8b02-4d134a21db89" />

Теперь сделаем так, чтоб графана знала о jaeger и при каждом логе в loki появлялась ссылка в jaeger:

Скорректируем values у хелма для графаны:

<img width="1062" height="1386" alt="image" src="https://github.com/user-attachments/assets/5b7f895d-b06d-4c4d-8526-1b8d892e3ea0" />

Выполним helm upgrade:

<img width="2880" height="1136" alt="image" src="https://github.com/user-attachments/assets/2e352137-3faf-425f-9a18-4188cf2c78c9" />

Проверим, что на дашборде все отрабатывает корректно и появилась ссылка на джагер:

<img width="1064" height="710" alt="image" src="https://github.com/user-attachments/assets/dc52280d-5d35-4765-b2a5-474b0d69a966" />

Теперь нажмем на ссылку и перейдем в джагер для просмотра трейса:

<img width="1422" height="1212" alt="image" src="https://github.com/user-attachments/assets/5d0d6048-ed8f-4d44-a665-4c1cabdae247" />

## Часть 4. Алерты

Теперь настроим алерты, для этого мы скорректируем конфиг prometheus:

<img width="780" height="654" alt="image" src="https://github.com/user-attachments/assets/469075a6-5b3f-4e1e-ae0d-2ae5e9ba959c" />

Обновим конфигурации нашего хелма:

<img width="2858" height="1144" alt="image" src="https://github.com/user-attachments/assets/9baf01f2-2545-4ad2-abde-afcf3ca344fb" />

Теперь напишем конфиг для prometheus, чтоб указать, при каких моментах нужно присылать алерт:

<img width="2342" height="1316" alt="image" src="https://github.com/user-attachments/assets/fe6c8fae-704c-49e4-8a6b-78976e9bb6c0" />

Обновим хелм:

<img width="2834" height="1244" alt="image" src="https://github.com/user-attachments/assets/60325ba9-5a9b-4f66-88c3-29d7072d147a" />

Проверим, что все нормально загрузилось:

<img width="2874" height="1024" alt="image" src="https://github.com/user-attachments/assets/c757c037-2bab-4eaf-af31-c3f8c1a93e66" />

Теперь попробуем запустить curl заведомо алертный запрос и проверим алерт в prometheus:

<img width="2880" height="1120" alt="image" src="https://github.com/user-attachments/assets/53509df5-bd0e-4dbe-b090-b8bfa8f3abca" />

Теперь запустим алертменеджер и посмотрим алерты там:

<img width="1468" height="570" alt="image" src="https://github.com/user-attachments/assets/6d46923b-4bd0-4321-80dd-3163c9fde6ed" />

Теперь начнем устанавливать карму, чтоб был удобный вывод, настроим env файл:

<img width="1152" height="106" alt="image" src="https://github.com/user-attachments/assets/c790382f-65d9-4bf7-a89e-72a1f6b0febf" />

<img width="2870" height="514" alt="image" src="https://github.com/user-attachments/assets/48fb8146-294f-4297-8d7c-9ee2bf562274" />

Проверим, что алерты отрабатывают корректно:

<img width="998" height="576" alt="image" src="https://github.com/user-attachments/assets/a2f32caa-021b-4e27-8f3d-9cef3c5586be" />

Проверим в алертменеджере, что алерты пришли:

<img width="2376" height="760" alt="image" src="https://github.com/user-attachments/assets/6555c04d-8ffc-4b6c-9782-8cb3ddf60745" />

<img width="2376" height="760" alt="image" src="https://github.com/user-attachments/assets/a9ef6827-d521-4a5c-92a5-444652e1074e" />

Теперь уменьшим количество подов и проверим работоспособность алертов:

<img width="1502" height="164" alt="image" src="https://github.com/user-attachments/assets/d88837a5-d419-4799-b43a-0d12e83d39a1" />

Теперь смотрим в alertmanager и karma:

<img width="1318" height="300" alt="image" src="https://github.com/user-attachments/assets/d574963f-1b7e-4077-a49e-83d1f9bc72c5" />

<img width="998" height="334" alt="image" src="https://github.com/user-attachments/assets/ed639b89-a676-41af-8e25-13c96cd4ac51" />
