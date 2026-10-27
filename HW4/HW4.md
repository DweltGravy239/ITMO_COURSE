Лабораторная работа № 4 - SLA критического пути под перегрузкой
Подготовка
Для начала нашей лабораторной работы нам необходимо создать новый кластер. Для этого, мы сначала остановим старую ноду
<img width="2843" height="433" alt="image" src="https://github.com/user-attachments/assets/6dc6420e-2b8a-4ffd-a743-a9f0c5311659" />
Теперь напишем манифест для создания нового кластера и применим изменения:
<img width="1940" height="1093" alt="image" src="https://github.com/user-attachments/assets/47203aad-0b6f-4c0b-9ff0-6906248a405a" />
Теперь снимем taint с control-plane, чтоб поды могли распределятся на обе ноды.
<img width="2354" height="180" alt="image" src="https://github.com/user-attachments/assets/8baa0ae7-55ed-49cb-ba61-e476491cfdfa" />
Теперь скопируем уже написанные хелмы и политики из предыдущей лабораторной работы:
<img width="1605" height="338" alt="image" src="https://github.com/user-attachments/assets/db1b8a63-cc65-46dc-b566-5d685661f135" />
Установим kyverno:
<img width="2005" height="1155" alt="image" src="https://github.com/user-attachments/assets/73f84e03-1d89-4eee-949c-3051b0bb2af6" />
Теперь применим все политики:
<img width="1067" height="531" alt="image" src="https://github.com/user-attachments/assets/e3e9c23c-6a81-400f-a20b-c100bacce429" />
Установим оператор CloudNativePG:
<img width="1721" height="1042" alt="image" src="https://github.com/user-attachments/assets/fedd838d-741a-4fbe-8208-2e288ed49c85" />
Теперь загрузим оба образа в наш кластер:
<img width="2876" height="1031" alt="image" src="https://github.com/user-attachments/assets/bc5d4bb2-5a4f-4c22-8118-28931062d999" />
Теперь перенесем конфигурации мониторинга и немного скорректируем сам конфигурационный файл:
<img width="1489" height="1354" alt="image" src="https://github.com/user-attachments/assets/b9437395-678e-4cb9-9153-6bf8f6045771" />
Теперь установим сам мониторинг: 
<img width="2826" height="1444" alt="image" src="https://github.com/user-attachments/assets/58126c4f-2182-44dc-9eb4-0cff78e9cea2" />
Теперь осталось поставить само наше приложение из лабораторной работы 2:
<img width="1466" height="562" alt="image" src="https://github.com/user-attachments/assets/cd454029-76fd-490e-84c4-34bd5783f3e6" />
Проверим, что база поднялась:
<img width="1615" height="436" alt="image" src="https://github.com/user-attachments/assets/dfbe6d8f-b527-4128-b377-4344c3a492bf" />
Проверим, что все отрабатывает корректно:
<img width="2796" height="359" alt="image" src="https://github.com/user-attachments/assets/224478da-e1b2-437b-b542-d07fb4cd20b9" />
Часть 0 — Сервис batch
Сгенерируем сервис и пропишем его в каталоге:
<img width="1598" height="1376" alt="image" src="https://github.com/user-attachments/assets/7b8b9dda-70e1-4799-9cf4-519de5f449c8" />
Соберем докерфайл и сбилдим его:
<img width="2857" height="1178" alt="image" src="https://github.com/user-attachments/assets/aa663754-a594-4dce-b101-109cf005db90" />
Теперь проверим, что образ корректно запускается:
<img width="2466" height="473" alt="image" src="https://github.com/user-attachments/assets/15c444f5-c297-4b69-960d-2387cf06630f" />
Загрузим образ:
<img width="2870" height="237" alt="image" src="https://github.com/user-attachments/assets/1ed52ab2-db82-4413-8095-0902318610b1" />
Часть 1 — Выбери механизм разнесения
В качестве параметра разъяснения выберем вариант topologySpreadConstraints, т.к. при использование podAntiAffinity и при обновление деплоймента не получится создавать временный под, для того чтоб бесшовно провести обновление. Правило topologySpreadConstraints с maxSkew 1 и DoNotSchedule допускает только раскладки 2+1 и 1+2 и запрещает 3+0, поэтому при потере любой из двух нод минимум одна реплика api остаётся работать. podAntiAffinity не подходит, потому что он не гарантирует расселение и из-за этого может все поды положить на одну ноду. 
Скорректируем values.yaml:
<img width="1206" height="874" alt="image" src="https://github.com/user-attachments/assets/a8a272f1-c783-49f6-a6b9-177b2532a4c3" />
Часть 2 — Тесная площадка
Откроем values и добавим раздел с batch:
<img width="1038" height="1282" alt="image" src="https://github.com/user-attachments/assets/4de9ef1f-b232-47bf-bb09-90678ad63556" />
Теперь сделаем шаблон shop-banch и проверим, что ошибок у нас нету:
<img width="1430" height="1434" alt="image" src="https://github.com/user-attachments/assets/1cf22438-63ce-4ec0-a404-a3d8d76b0d15" />
Применим наш хелм:
<img width="2364" height="736" alt="image" src="https://github.com/user-attachments/assets/940a33f7-3b54-437a-bf62-a380becd6d31" />
Часть 3 — Реальные числа вместо угадывания
Перед прогоном я бы выставил лимиты такие:
- api - cpu 1 - memory 128Mi - limits - cpu 1 - memory 128Mi
- worker - cpu 1 - memory 200Mi - limits - cpu 1.5 - memory 260Mi
- postgres - cpu 1 - memory 512Mi - limits - cpu 1 - memory 512Mi
Теперь прогоним наш hey в течении 20 минут:
<img width="2836" height="362" alt="image" src="https://github.com/user-attachments/assets/b803cc9e-2cdf-4e71-a378-77f8b32e8ff2" />
Пока у нас работает hey, мы успеем поставить krr:
<img width="1542" height="834" alt="image" src="https://github.com/user-attachments/assets/dc499c5b-ab25-46ed-85d2-c40b9483e114" />
Посмотрит итоговую сводку hey:
<img width="718" height="1160" alt="image" src="https://github.com/user-attachments/assets/c22ef6bf-dde6-414b-9368-b693a0438eb0" />
Так как у нас база данных разворачивается через оператора, krr не сможет брать оттуда метрики, поэтому мы сделаем так, чтоб он брал через прометеус. Для этого проверим, есть ли вообще эти метрики в нем:
<img width="2758" height="178" alt="image" src="https://github.com/user-attachments/assets/40f95183-0e26-427f-a9d6-1582439e617a" />
Запустим krr и обратимся к прометеус, для того чтоб понять нагрузку на поды:
<img width="1146" height="1050" alt="image" src="https://github.com/user-attachments/assets/d5ad5c37-0e53-4706-b08b-3543427a04e3" />
Как мы видим, с использованием krr и информации из prometheus мы видим, что значения:
- api - cpu 73m - memory 20Mi - limits - cpu не задан - memory 20Mi
- worker - cpu 10m - memory 20Mi - limits - cpu не задан - memory 20Mi
- postgres - cpu 180m - memory - 127Mi - cpu не задан - memory 127Mi
Теперь внесем информацию в наш values.yaml:
<img width="646" height="1038" alt="image" src="https://github.com/user-attachments/assets/35a8d700-119d-4bcc-b55a-3d02f21ed3b1" />
Теперь добавим указанную информацию в темплейтс:
<img width="1546" height="1122" alt="image" src="https://github.com/user-attachments/assets/659630cc-b7b4-4253-aaf9-e83a11472d81" />
Применим изменения:
<img width="1704" height="310" alt="image" src="https://github.com/user-attachments/assets/05cd64c2-1dc3-4604-9499-e6f7e50102b7" />
Проверим, что все поды у нас живут:
<img width="1034" height="308" alt="image" src="https://github.com/user-attachments/assets/03c45dbd-fabf-4719-b90b-e1edd891ecbb" />
Часть 4 — Разнеси критический сервис
Обновим наш манифест, который добавляет выбранну методику из части 1:
<img width="1216" height="802" alt="image" src="https://github.com/user-attachments/assets/e50f05f6-6ca5-4844-9619-2492d868f25d" />
Применим изменения:
<img width="1936" height="496" alt="image" src="https://github.com/user-attachments/assets/5dcf60a2-ae80-47a3-bf05-216209601d0e" />
Часть 5 — Приоритеты и гарантии
Для того чтоб статус QoS был Guaranteed нам нужно, чтоб limits и requests были равны.
Значения, который я выставил:
api - cpu: 200 - memory: 64 (повысил лимиты и request для того, чтоб был запас и был запас по потреблению)
postgres - cpu: 400 - memory: 256 (также повысил лимиты по причине выше)
batch - cpu: 100 - memory: как было (для того чтоб защитить api, т.к. batch получит статус "Burstable" и будет не в приоритете)
Теперь добавим priorityClass в наш values и напишем манифест:
<img width="904" height="1422" alt="image" src="https://github.com/user-attachments/assets/80893585-41fe-4897-a6b8-1b633bff6a85" />
Добавим блок pdb в values и в манифесты:
<img width="2094" height="486" alt="image" src="https://github.com/user-attachments/assets/1a151102-540e-42c1-9b01-b02949286e89" />
Обновим хелм:
<img width="2066" height="366" alt="image" src="https://github.com/user-attachments/assets/a60374d6-2b67-453d-bcc6-56d929986809" />
Проверим, что все отрабатывает корректно:
<img width="1578" height="780" alt="image" src="https://github.com/user-attachments/assets/e1b74c95-20da-405e-8891-41282510394b" />
Часть 6 — Создай нехватку, поймай preemption
Для начала узнаем, сколько планировщик считает свободным местом на нодах:
<img width="2004" height="904" alt="image" src="https://github.com/user-attachments/assets/539a9245-064b-43f6-a526-369b9a886579" />
Для того чтоб словить ошибку мы заложим больше ресурсов, чем нужно. Это мы сделаем повышением числа реплик на shop-branch:
<img width="792" height="502" alt="image" src="https://github.com/user-attachments/assets/c29135a2-8751-46f0-85e3-8a08e1b94011" />
Применим изменения и проверим результат:
<img width="2410" height="1184" alt="image" src="https://github.com/user-attachments/assets/5fd8513d-1191-4d47-a3b4-740702c9ffab" />
Как мы видим, все поды поднялись и в лимиты мы не ушли, поэтому поднимаем количество репликаций до 14 и применим изменения:
<img width="2372" height="922" alt="image" src="https://github.com/user-attachments/assets/30f07450-590c-4d4a-a71e-78e72e7b6f11" />
Посмотрим лог пода для того чтоб подтвердить причину:
<img width="2874" height="460" alt="image" src="https://github.com/user-attachments/assets/ee1a7a35-4010-466b-ad45-6937b90d474b" />
Теперь посмотрим на занятые ресурсы на нодах и проверим, что значения изменились:
<img width="2048" height="884" alt="image" src="https://github.com/user-attachments/assets/02db78ad-7328-4619-b4b8-e78d89025ad0" />
Теперь проверим, что на самом деле CPU пойти свободен и планировщик берет ресурсы из requests, а не из фактического потребления:
<img width="1406" height="278" alt="image" src="https://github.com/user-attachments/assets/7ac6330e-d9a7-4e08-8cda-1a1284f62538" />
Следующим шагом попробуем увеличить число критичных подов (api) и посмотрим, как он должен будет вывести ненужные поды branch:
<img width="772" height="470" alt="image" src="https://github.com/user-attachments/assets/77347957-6504-4fc7-8d84-37836dd62339" />
Когда мы проверяли, мы столкнулись с интересным событием. После переполнения ресурсов вывели не branch, как мы ожидали, а kiverno. Это все из-за того, что мы там не выставляли приоритеты и там был 0, а на branch 100.
<img width="1668" height="344" alt="image" src="https://github.com/user-attachments/assets/02f27cf7-33e5-495d-a046-1616e04d4faf" />
Попробуем уменьшить число репликаций с 14 до 10 и проверить, восстановится ли кто-то:
<img width="1696" height="510" alt="image" src="https://github.com/user-attachments/assets/c37b78a9-e44a-4f64-935b-4961e87881ec" />
Для того чтоб у нас точно выводился branch мы им дадим отрицательный приоритет:
<img width="568" height="264" alt="image" src="https://github.com/user-attachments/assets/1a38398a-19f6-4785-8505-d67187457488" />
Применим изменения и проверим, что все применилось:
<img width="2606" height="976" alt="image" src="https://github.com/user-attachments/assets/8133a33a-3924-4a6f-b782-044ed8bd2f13" />
Снова забиваем кластер и поднимаем число репликаций до 14:
<img width="2136" height="662" alt="image" src="https://github.com/user-attachments/assets/a2fa1de5-e0a9-48a6-9fd3-d7db5fe68653" />
Теперь пробуем поднять число подов для api и проверить, что начнут умирать поды brench:
<img width="2846" height="1396" alt="image" src="https://github.com/user-attachments/assets/d5fc2d65-4d88-47bc-87a8-00de2007eb7d" />
Как мы видим, мы поймали именно то, что хотели.
Осталось вернуть все на прежние места и применить изменения:
<img width="1938" height="678" alt="image" src="https://github.com/user-attachments/assets/9b54ab73-0090-47a5-96e1-62746095d22d" />
Часть 7 — Давление по памяти, порядок жертв
Для начала посмотрим, что кубелет считает доступной памяти на каждой из нод:
<img width="1752" height="576" alt="image" src="https://github.com/user-attachments/assets/12ecc0e1-9e05-4dba-9c3f-c597ea4dbe61" />
Теперь добавим порог в конфиг кубелета и проверим результаты
<img width="1820" height="242" alt="image" src="https://github.com/user-attachments/assets/183e1a6a-c440-498b-bb24-35c29143dc7a" />
Теперь внесем изменения на другой же ноде:
<img width="1878" height="252" alt="image" src="https://github.com/user-attachments/assets/5c93dea2-ae65-4035-9a1b-3a1efc4d2260" />
Теперь перезагрузим кубелет для того чтоб конфиг прихватился:
<img width="1578" height="320" alt="image" src="https://github.com/user-attachments/assets/8ae5d0e0-786d-4840-8d7f-a38bc207f0e4" />
Теперь выполним на следующей ноде:
<img width="1728" height="254" alt="image" src="https://github.com/user-attachments/assets/5fb3d683-2b7f-4d90-b6f5-893cbda00174" />
Для того чтоб у нас начались вытеснения, нам необходимо поменять values в блоке batch:
<img width="594" height="374" alt="image" src="https://github.com/user-attachments/assets/c81e0a96-a47c-4fc0-8a2f-ed71d2240a87" />
Проверим возможность применения изменений:
<img width="2840" height="1164" alt="image" src="https://github.com/user-attachments/assets/45ecc287-e407-43e1-925c-497cb7c70006" />
Как мы видим, на sla-worker стало меньше порога, поэтому kubelet выставил MemoryPressure=True и начал вытеснять поды batch. Он начал вытеснять их, из-за того что их приоритет равен отрицательному значению и они являются главной цели при возможной нехватки места.
Теперь посмотрим состояние пода:
<img width="1890" height="106" alt="image" src="https://github.com/user-attachments/assets/0d6237bd-94e8-432f-8b5c-830e5eae12ee" />
Мы видим, что у ноды MemoryPressure стоит True с причиной KubeletHasInsufficientMemory, это значит, что под давлением на эту ноду больше не садятся другие поды.
Теперь вернем все значения на прежние:
<img width="526" height="336" alt="image" src="https://github.com/user-attachments/assets/19751a48-a34f-4ed3-bba3-80d77afb961c" />
Применим изменения:
<img width="1434" height="226" alt="image" src="https://github.com/user-attachments/assets/38f38736-4404-457e-ad86-9e44dfc0e572" />
Уберем следы вытеснения:
<img width="1072" height="356" alt="image" src="https://github.com/user-attachments/assets/536a7246-7422-4ec6-856e-8a0c6dfa3f78" />
Проверим, что все нормально и давление снято:
<img width="1864" height="426" alt="image" src="https://github.com/user-attachments/assets/f42c4122-8ac5-4a76-a00b-5164e2bb106c" />
Часть 8 — Докажи SLA под нагрузкой
Для того чтоб сформировать SLA, нужно прогнать еще раз нагрузку через hey и посмотреть результаты:
<img width="674" height="1092" alt="image" src="https://github.com/user-attachments/assets/d5e2dbae-148e-4206-a7ae-284a599b0956" />
По метрикам для формирования SLA, нужно обратить внимание на p95 (строка Latency distribution) и долю успешных ответов (поле Status code distribution).
Теперь нужно построить график в prometheus, для того чтоб это сделать, нам нужно понять, какие переменные брать, для этого выведем все переменные и выберем нужную:
<img width="2400" height="1176" alt="image" src="https://github.com/user-attachments/assets/86c9a4c7-e7ab-4951-8d17-cf4134bbc4aa" />
Теперь построим график в прометеус:
<img width="878" height="382" alt="image" src="https://github.com/user-attachments/assets/24b22846-7b03-4547-8e9b-bb9a965bb797" />
На графике видно, что что во время вытеснения batch критический путь отвечал стабильно, и значение p95 составляло 0.159. Единственный скачок был во время helm upgrade. По нему мы можем понять, что наш SLA соблюдается. А SLA формируем вот так: "за любое 30-минутное окно 95% запросов POST /order выполняются быстрее 300 мс и не менее 99.9% ответов успешны". Если опираться на наш SLA, то у нас все соблюдается, т.к. при p95 = 159мс у нас 100 процентов успешных ответов
Теперь построим график доли успеха:
<img width="1710" height="1244" alt="image" src="https://github.com/user-attachments/assets/cb6d89f2-80ad-49e6-b77d-aac49315eebc" />
Как мы видим, наш график прямо сопоставим с словами выше. Так как у нас доля успешных запросов равна 100 процентов, то и весь график держится на единице. 
Часть 9 — Мониторинг
Мы будем собирать определенные алерты. Эти алерты:
- p95 выше порога, который мы определили в SLA и он держится несколько минут.
- доступных реплик shop-api стало меньше 2
- memorypressure равен true. Это предвестник того, что скоро kubelet начнет вытеснять поды.
Для этого мы добавим новые параметры в values:
<img width="2334" height="966" alt="image" src="https://github.com/user-attachments/assets/4efbb591-dc2b-4317-8737-69ea2a3d62e1" />
Обновим наш хелм:
<img width="1902" height="362" alt="image" src="https://github.com/user-attachments/assets/90c83493-58a6-4a4d-af32-8a819d040ad9" />
