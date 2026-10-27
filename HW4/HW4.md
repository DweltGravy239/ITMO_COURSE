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



Часть 6 — Создай нехватку, поймай preemption
Часть 7 — Давление по памяти, порядок жертв
Часть 8 — Докажи SLA под нагрузкой
Часть 9 — Мониторинг


Теперь осталось поставить само наше приложение из лабораторной работы 2:
<img width="1466" height="562" alt="image" src="https://github.com/user-attachments/assets/cd454029-76fd-490e-84c4-34bd5783f3e6" />
Проверим, что база поднялась:
<img width="1615" height="436" alt="image" src="https://github.com/user-attachments/assets/dfbe6d8f-b527-4128-b377-4344c3a492bf" />
Проверим, что все отрабатывает корректно:
<img width="2796" height="359" alt="image" src="https://github.com/user-attachments/assets/224478da-e1b2-437b-b542-d07fb4cd20b9" />
