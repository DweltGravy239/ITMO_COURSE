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
Часть 4 — Разнеси критический сервис
Часть 5 — Приоритеты и гарантии
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
