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

Этап 4. Права capabilities и seccomp
Этап 5. mydocker.sh и сравнение с docker run
Этап 6. Образы
Этап 7. gVisor
Этап 8. Мониторинг
