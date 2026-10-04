Лабораторная работа №2 - Observability
Цель:
Часть 0 - Создание своего сервиса.
Для создания своего сервиса я воспользовался клодом. В итоге он мне сгенерировал что-то такое (более подробнее можно увидеть в каталоге /app):
<img width="1777" height="1021" alt="image" src="https://github.com/user-attachments/assets/c5c95b1d-c167-4d78-bdd4-a4c3ad1c855e" />
После создания приложения я написал базовый докерфайл и сбилдил его:
<img width="2877" height="451" alt="image" src="https://github.com/user-attachments/assets/547cc9f8-97ba-4125-a8c4-1a9af7940541" />
А после этого прописал его в kind:
<img width="2878" height="192" alt="image" src="https://github.com/user-attachments/assets/59ff7950-4451-426d-bfdb-60d08c1f32d1" />
Сначала я решил набросать обычный манифест для того, чтоб проверить его корректность, а потом перевести его в helm:
<img width="1084" height="1334" alt="image" src="https://github.com/user-attachments/assets/7d6e10a2-6c44-4471-878b-f45ffce73d46" />
Дальше мы сделаем kubectl apply и проверим, что он отрабатываем корректно:
<img width="1350" height="516" alt="image" src="https://github.com/user-attachments/assets/9fce6480-ef75-41ba-a673-dea973ac5ee9" />
Проверим, что веб-приложение отрабатываем корректно:
<img width="1752" height="678" alt="image" src="https://github.com/user-attachments/assets/a2d26048-568b-4b4a-841e-45325f9dc456" />
Когда мы выяснили, что мой манифест отрабатывает корректно, можно начать писать helm:
<img width="1500" height="1354" alt="image" src="https://github.com/user-attachments/assets/2379b1c3-c9d4-4f86-95b4-420124644e4e" />
Теперь попробуем применить наш чарт:
<img width="1206" height="406" alt="image" src="https://github.com/user-attachments/assets/7aa5b57f-5ee7-4f27-b56d-85258640a646" />
Проверяем, что наш helm применился:
<img width="1388" height="554" alt="image" src="https://github.com/user-attachments/assets/a4ab3838-2b4c-4408-be0c-177fb8ce551c" />
Часть 1. Установка Prometheus и Grafana
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
Часть 2. Установка Loki




