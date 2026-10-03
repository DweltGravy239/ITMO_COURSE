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

