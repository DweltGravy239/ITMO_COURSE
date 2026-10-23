
## Часть 0 — Сервисы

Сгенерируем сервис api и worker

<img width="813" height="171" alt="image" src="https://github.com/user-attachments/assets/65488226-a923-448e-9efd-284eb9041fab" />

Напишем докерфайл для нашего api, после чего соберем его:

<img width="2827" height="1431" alt="image" src="https://github.com/user-attachments/assets/91cbf505-44e5-4a0a-bfa6-34d1f97c9890" />

Теперь соберем образ для worker:

<img width="2851" height="1404" alt="image" src="https://github.com/user-attachments/assets/fac9d290-7028-4e27-8e0a-08afcc7ca996" />

Проверим, что образы собрались и лежат в docker images:

<img width="2311" height="685" alt="image" src="https://github.com/user-attachments/assets/705654fd-e8b1-4d04-9fc8-b4cf64c338bd" />

## Часть 1 — Ограждения на кластер

Для лабораторной работы я выбрал Kyverno, потому что я считаю, что использование Gatekeeper избыточно в данный момент. Сам Kyverno используется для того, чтоб перед загрузкой какого-либо манифеста написать правила, которые необходимо соблюдать.

Теперь установим Kyverno через helm:

<img width="2206" height="915" alt="image" src="https://github.com/user-attachments/assets/4213ad82-0aec-471a-8890-0f0332b2e010" />

Проверим, что поды поднялись:

<img width="1619" height="320" alt="image" src="https://github.com/user-attachments/assets/85eadb5e-5502-4eee-868b-f5c5804302c0" />

Сделаем первое барьерное правило, которое будет блокировать поды, в которых не указан лимит:

<img width="1739" height="674" alt="image" src="https://github.com/user-attachments/assets/15209980-038c-4d35-89b8-1563a0f764d6" />

Применим изменения:

<img width="1571" height="266" alt="image" src="https://github.com/user-attachments/assets/72f594b8-b115-463a-8e09-a072b6ea33a5" />

Теперь проверим, что правила отрабатывает корректно. Добавим манифест для создания подов, где в одном укажем лимиты, а в другом нет:

<img width="952" height="883" alt="image" src="https://github.com/user-attachments/assets/00af31b2-3390-4e93-b0f4-b068a781f94a" />

Проверим, как они отработают:

<img width="2460" height="221" alt="image" src="https://github.com/user-attachments/assets/73cf9805-5d7f-4831-8856-af5edd42691e" />

Как мы видим, все отработало корректно. С лимитами под отрабатывает, а без лимитов сообщает об ошибке. Если бы у нас был режим не Аудит, то тогда нам бы не дало создать его.

Теперь начнем создавать нужные правила:

### 1. У каждого контейнера задан лимит памяти

(для того чтоб рандомный разработчик не смог загрузить под, который сожрал бы у нас всю память)

<img width="1736" height="1166" alt="image" src="https://github.com/user-attachments/assets/9a3a2133-6b90-481d-9f0b-c55a5cb2d8e5" />

Теперь создадим собственный namespace, чтоб применять правила только для shop(это мы возьмем за правило и будем использовать всегда)

<img width="2140" height="1149" alt="image" src="https://github.com/user-attachments/assets/dd30dc79-5c67-47d2-b91e-9045137c055a" />

### 2. Запрещаем тег latest

(для того, чтоб фиксировать определенную версию и не попасть на то, что версия поменялась, а мы полетели с зависимостями)

Создадим это правило:

<img width="1678" height="767" alt="image" src="https://github.com/user-attachments/assets/5305e2bd-53e1-41c0-8e08-2123c2567465" />

Теперь проверим, что правило отрабатывает корректно, для этого создадим нужные поды:

<img width="885" height="878" alt="image" src="https://github.com/user-attachments/assets/6f6b0bc4-3d62-4a76-a25d-b5dfdf3ab4ce" />

Применим их и проверим, что все отрабатывает корректно:

<img width="1979" height="433" alt="image" src="https://github.com/user-attachments/assets/250a49eb-783d-4430-96f8-48c7776cf0a0" />

### 3. Необходимо использовать только registry.local

(для того, чтоб брать образы только из доверенного источника)

Создаем правило:

<img width="1532" height="768" alt="image" src="https://github.com/user-attachments/assets/100f3e00-0d7e-4f80-87eb-1ee0ec0346d3" />

Создаем поды и проверяем, что правило работает:

<img width="2006" height="1321" alt="image" src="https://github.com/user-attachments/assets/54736727-50dd-4cc2-930b-ed83ad544727" />

### 4. Не использовать root

(Задаем такое правило, чтоб контейнер был более безопасный и пользователь в контейнере не имел прав на целевом хосте)

Создадим правило:

<img width="2818" height="838" alt="image" src="https://github.com/user-attachments/assets/044fee83-dce1-4150-9202-6358520fdbcb" />

Проверим, что он отрабатывает:

<img width="2767" height="1093" alt="image" src="https://github.com/user-attachments/assets/3ffa305d-eddc-408f-9ccb-136de99c99cc" />

### 5. Обязательно указывать метки

(задаем, т.к. по ним склеивается вся логика в кубере)

Создадим правило:

<img width="2258" height="926" alt="image" src="https://github.com/user-attachments/assets/97a610db-87fa-4cdb-b584-d481d0535c76" />

Теперь создадим эталонный манифест, который будет подходить по всем правилам:

<img width="2031" height="968" alt="image" src="https://github.com/user-attachments/assets/9e6bac68-f561-42ff-906a-fb9485c67f2e" />

Создадим 5 манифестов-нарушителей:

<img width="852" height="1281" alt="image" src="https://github.com/user-attachments/assets/4aeca7ad-7df6-4331-813c-ffc58b60872a" />

Включим режим Deny, вместо Audit:

<img width="1200" height="122" alt="image" src="https://github.com/user-attachments/assets/7e181e18-8fa2-4fa8-9df0-9caf8cf9336b" />

Применим изменения и проверим, что все отрабатывает корректно:

<img width="2848" height="727" alt="image" src="https://github.com/user-attachments/assets/c6efa87f-55b1-4607-8e45-deb362d98822" />

Также выводим правила:

<img width="2855" height="990" alt="image" src="https://github.com/user-attachments/assets/c3fa2aaa-ec52-4224-8e01-f1e10398b1c2" />

## Часть 2 — Чарт api и worker

Загружаем наши образы в kind:

<img width="2847" height="930" alt="image" src="https://github.com/user-attachments/assets/f48e07eb-5c1e-47cb-af63-84d0d37d520e" />

Сделаем каркас для helm:

<img width="1564" height="779" alt="image" src="https://github.com/user-attachments/assets/c87cd724-e652-4c2c-9d63-393061ccef36" />

Для создания хелма, сначала сделаем обычный манифест, который потом переделаем в хелм:

<img width="1293" height="1392" alt="image" src="https://github.com/user-attachments/assets/07ad9191-dc5a-4ee6-ac28-dddc0f160dd7" />

Добавим values и скорректируем наш манифест:

<img width="1052" height="1328" alt="image" src="https://github.com/user-attachments/assets/5c8dbc65-f549-4b1d-9b9d-59aa9ba23b23" />

Проверим, что все отрабатывает корректно:

<img width="692" height="966" alt="image" src="https://github.com/user-attachments/assets/70df51ec-dacd-4567-97b9-236079c0a577" />

Добавим еще парочку параметров и проверим, что все отрабатывает корректно:

<img width="738" height="1274" alt="image" src="https://github.com/user-attachments/assets/f1a540c8-fd7d-4e7e-afa6-6180420db0c4" />

Теперь создадим хелм для worker и проверим его:

<img width="414" height="1226" alt="image" src="https://github.com/user-attachments/assets/5a8b09b2-e8b1-4a45-989b-891a0f222f7a" />

Установим чарты:

<img width="2118" height="1029" alt="image" src="https://github.com/user-attachments/assets/cc06d839-6916-408a-9fde-d195c1579ee1" />

Попробуем удалить под и проверим, что он пересоздаться:

<img width="1641" height="618" alt="image" src="https://github.com/user-attachments/assets/10ab00ce-1912-4786-8061-b949ed4b2684" />

Теперь попробуем сделать обновление нашего образа:

<img width="2878" height="356" alt="image" src="https://github.com/user-attachments/assets/0f38dde8-ce01-44e7-9135-059286589212" />

Теперь обновим наш хелм, для только чтоб сделать успешный rollingupdate:

<img width="768" height="408" alt="image" src="https://github.com/user-attachments/assets/eadf3cb7-9412-4ff5-9fa2-53690e6ac90e" />

Загрузим в кубер новый образ, для того чтоб в дальнейшем сгенерировать нагрузку:

<img width="1914" height="226" alt="image" src="https://github.com/user-attachments/assets/36531815-0959-478e-8976-b46ce4e55410" />

Запустим нагрузку на наши поды и пробуем провести плановое обновление:

<img width="2734" height="1360" alt="image" src="https://github.com/user-attachments/assets/b2302e94-ebd1-4e4e-b108-8491856f552c" />

### Ошибка при раскатке

Теперь пробуем сделать ошибку при раскатке

Даем нагрузку и пытаемся сделать заведомо плохую раскатку:

<img width="2880" height="1408" alt="image" src="https://github.com/user-attachments/assets/3650b479-2ce3-4bdf-860e-ee87ce24527c" />

Мы видим, что даже при ошибке в развертывание, хелм не дал нам обновится и нагрузка на наше приложение не переключилось на сломанные ноды

Делаем ролбек:

<img width="2244" height="630" alt="image" src="https://github.com/user-attachments/assets/1c555f96-4385-4943-b7e6-0ac3465455d7" />

## Часть 3 — Postgres через оператора

Поставим оператор и добавим репозиторий чарта CloudNativePG:

<img width="1610" height="1254" alt="image" src="https://github.com/user-attachments/assets/c27321f7-75e1-4b09-9252-57c83856c9c5" />

Напишем хелм для базы данных:

<img width="914" height="802" alt="image" src="https://github.com/user-attachments/assets/dfa6932d-09eb-4057-a210-6158ee1573f4" />

Применим изменения:

<img width="1052" height="616" alt="image" src="https://github.com/user-attachments/assets/a3f99c43-b5c9-43af-8172-561d4369ac69" />

Проверим, что все отрабатывает корректно и поды создались:

<img width="1246" height="356" alt="image" src="https://github.com/user-attachments/assets/df3345ec-d402-4f54-bcab-ce8608cdce13" />

Подключаем наше приложение к новосозданной базе данных:

<img width="1161" height="1226" alt="image" src="https://github.com/user-attachments/assets/e2f83713-da6f-41be-808b-2c6dd0c4b870" />

Проверяем, что все отрабатывает корректно:

<img width="1716" height="906" alt="image" src="https://github.com/user-attachments/assets/be0804c7-9204-479d-8eb1-b72e050d0658" />

Теперь попробуем удалить один из подов БД:

<img width="2012" height="940" alt="image" src="https://github.com/user-attachments/assets/6dd8c037-fd10-499f-aaba-0ea9b8793e37" />

Как мы видим, все отработало корректно.

Теперь попробуем убить primary ноду:

<img width="2057" height="497" alt="image" src="https://github.com/user-attachments/assets/48f15e75-cae8-48d9-839c-5d0c213e29a5" />

Видим, что primary нода поменялась и все отработало корректно.

Найдем spec и status:

<img width="1248" height="195" alt="image" src="https://github.com/user-attachments/assets/9f67156d-cd6a-412c-a37b-2f693cbfc294" />

controller-manager это встроенная часть кубера. Он следит только за обычными объектами вроде деплоев и умеет простые вещи, например удалиять или добавлять поды. Оператор CloudNativePG работает примерно также, но он разбирается именно в PostgreSQL, это было видно, когда у нас упала база и он автоматически переключил ее на другую.

## Часть 4 — Падение control plane

Проверим исходное состояние нашего кубера:

<img width="1254" height="192" alt="image" src="https://github.com/user-attachments/assets/78515edb-da85-4815-abbe-dbb62e76dca9" />

Остановим control plane. Убивать его мы будем через перенос манифестов в другую папку. Для начала создадим эту папку и перенесем туда все:

<img width="1256" height="168" alt="image" src="https://github.com/user-attachments/assets/e9961740-0654-4199-be41-6d3bf066931e" />

Проверим, что все контейнеры остановились и проверим, что мы ничего не можем сделать:

<img width="1260" height="242" alt="image" src="https://github.com/user-attachments/assets/9cc3c0e1-a224-4997-a21a-aca23329ee29" />

Теперь проверим, что контейнеры и приложение все равно живут без control plane:

<img width="1262" height="460" alt="image" src="https://github.com/user-attachments/assets/3807b77f-f254-4171-b524-2baef84e0ce8" />

Восстановим control plane и проверим это:

<img width="1268" height="552" alt="image" src="https://github.com/user-attachments/assets/c9fc589a-fec3-42f7-8d0a-23d48e12f6cd" />

## Часть 5 — Мониторинг

Подготовим новый values для нашего хелма, а именно отключим ненужные фитчи:

<img width="1496" height="1408" alt="image" src="https://github.com/user-attachments/assets/1384dc97-60f6-4d60-9fbb-82b646bfad1d" />

Применим изменения:

<img width="2864" height="1376" alt="image" src="https://github.com/user-attachments/assets/56f59d45-5666-4a4a-b38b-d2b199543ed9" />

Изменим уже существующие values и пропишем мониторинг:

<img width="726" height="287" alt="image" src="https://github.com/user-attachments/assets/1c7e9575-3e26-4f2c-9525-6aaa3b8174c9" />

<img width="562" height="295" alt="image" src="https://github.com/user-attachments/assets/be46ad28-d569-489f-be86-813d2e5941c9" />

Применим изменения:

<img width="1056" height="648" alt="image" src="https://github.com/user-attachments/assets/4d3b1bf8-b018-4d6e-83c5-910b0d1a7922" />

Проверим, что прометеус видит цели:

<img width="2544" height="909" alt="image" src="https://github.com/user-attachments/assets/3b8f6645-7f0c-41c1-9420-5671343dd631" />

Теперь напишем алерты:

<img width="884" height="668" alt="image" src="https://github.com/user-attachments/assets/35e44584-8995-4aaa-9e93-7ea174383a09" />

Применим изменения:

<img width="1005" height="305" alt="image" src="https://github.com/user-attachments/assets/5fc2e508-cfa2-4f8c-8723-2cfb08323702" />

### Провоцируем события

Теперь попробуем спровоцировать события:

Убьем ноды и проверим, что все отрабатывает корректно:

<img width="2399" height="452" alt="image" src="https://github.com/user-attachments/assets/e574bc32-e604-4b98-b01b-78ea83998ebe" />

Теперь проверим алерт с очередями на наших воркерах:

<img width="2541" height="1214" alt="image" src="https://github.com/user-attachments/assets/9c3a38ea-c217-450d-b6d8-9a2325605016" />

Вернем в работу и проверим, что очередь рассосалась:

<img width="1253" height="323" alt="image" src="https://github.com/user-attachments/assets/4f0d6be7-7368-4910-b21f-917137f49e0e" />

Теперь устроим контролируемую неисправность и проверим, что наши алерты отрабатывают корректно:

<img width="2537" height="372" alt="image" src="https://github.com/user-attachments/assets/22c76620-c4d7-4972-b835-43478bdc5005" />

Теперь восстановим и проверим, что все отрабатывает корректно:

<img width="1250" height="221" alt="image" src="https://github.com/user-attachments/assets/1ec08d79-aba8-4929-ab32-9cdd81a415c3" />

Финально проверяем, что все поднялось:

<img width="1619" height="1255" alt="image" src="https://github.com/user-attachments/assets/eb636770-98d8-4ebc-965d-13a43b81a0dc" />
