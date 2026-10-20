Часть 0 — Сервисы
Сгенерируем сервис api и worker
<img width="813" height="171" alt="image" src="https://github.com/user-attachments/assets/65488226-a923-448e-9efd-284eb9041fab" />
Напишем докерфайл для нашего api, после чего соберем его:
<img width="2827" height="1431" alt="image" src="https://github.com/user-attachments/assets/91cbf505-44e5-4a0a-bfa6-34d1f97c9890" />
Теперь соберем образ для worker:
<img width="2851" height="1404" alt="image" src="https://github.com/user-attachments/assets/fac9d290-7028-4e27-8e0a-08afcc7ca996" />
Проверим, что образы собрались и лежат в docker images:
<img width="2311" height="685" alt="image" src="https://github.com/user-attachments/assets/705654fd-e8b1-4d04-9fc8-b4cf64c338bd" />
Часть 1 — Ограждения на кластер
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
1. У каждого контейнера задан лимит памяти (для того чтоб рандомный разработчик не смог загрузить под, который сожрал бы у нас всю память)
<img width="1736" height="1166" alt="image" src="https://github.com/user-attachments/assets/9a3a2133-6b90-481d-9f0b-c55a5cb2d8e5" />
Теперь создадим собственный namespace, чтоб применять правила только для shop(это мы возьмем за правило и будем использовать всегда)
<img width="2140" height="1149" alt="image" src="https://github.com/user-attachments/assets/dd30dc79-5c67-47d2-b91e-9045137c055a" />
2. Запрещаем тег latest (для того, чтоб фиксировать определенную версию и не попасть на то, что версия поменялась, а мы полетели с зависимостями)
Создадим это правило:
<img width="1678" height="767" alt="image" src="https://github.com/user-attachments/assets/5305e2bd-53e1-41c0-8e08-2123c2567465" />
Теперь проверим, что правило отрабатывает корректно, для этого создадим нужные поды:
<img width="885" height="878" alt="image" src="https://github.com/user-attachments/assets/6f6b0bc4-3d62-4a76-a25d-b5dfdf3ab4ce" />
Применим их и проверим, что все отрабатывает корректно: 
<img width="1979" height="433" alt="image" src="https://github.com/user-attachments/assets/250a49eb-783d-4430-96f8-48c7776cf0a0" />
3. Необходимо использовать только  registry.local (для того, чтоб брать образы только из доверенного источника)
Создаем правило:
<img width="1532" height="768" alt="image" src="https://github.com/user-attachments/assets/100f3e00-0d7e-4f80-87eb-1ee0ec0346d3" />
Создаем поды и проверяем, что правило работает:
<img width="2006" height="1321" alt="image" src="https://github.com/user-attachments/assets/54736727-50dd-4cc2-930b-ed83ad544727" />
4. Не использовать root (Задаем такое правило, чтоб контейнер был более безопасный и пользователь в контейнере не имел прав на целевом хосте)
Создадим правило:
<img width="2818" height="838" alt="image" src="https://github.com/user-attachments/assets/044fee83-dce1-4150-9202-6358520fdbcb" />
Проверим, что он отрабатывает:
<img width="2767" height="1093" alt="image" src="https://github.com/user-attachments/assets/3ffa305d-eddc-408f-9ccb-136de99c99cc" />
5. Обязательно указывать метки (задаем, т.к. по ним склеивается вся логика в кубере)
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
Часть 2 — Чарт api и worker
Часть 3 — Postgres через оператора
Часть 4 — Падение control plane
Часть 5 — Мониторинг
