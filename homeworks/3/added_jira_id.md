# Домашнее задание №3

P.S. Не описал кейсы и что делать, если вдруг что-то пойдет не так. В уроке это упоминалось, но перед дедлайном не хватило времени на это.

## Описание

Мы отдали часть системы из авторизации и таск трекера бизнесу. Спустя 15 минут бизнес попросил внести необходимые доработки после фидбэка.

Во время работы с таск трекером, мы поняли, что попуги часто в title задачи пишут конструкцию [jira-id] - Title (Пример: "UBERPOP-42 — Поменять оттенок зелёного на кнопке"). В результате чего поэтому нам необходимо разделить title на два поля: title (string) + jira_id (string). При этом, для всех новых событий, необходимо убедиться что jira-id не присутствует в title. Для этого достаточно сделать валидацию на наличие квадратных скобок (] или [)

1. Добавить jira_id в события связанные с задачами и опционально поменять код в сервисах

2. Расписать процесс миграции на новую логику в текстовом виде (можно в отдельном файле, можно в описании ПР-а)

## Процесс реализации

1. Добавляется новая версия события TaskCreated в schema-registry:

    1.1 Добавляется новая версия протофайла
    1.2 Генерируются pb-файлы в go-module
    1.3 Коммитится код и отправляется в github

2. Готовим новую версию analytics:
    Так как в текущем виде analytics просто консюмит события то:

     - импортируется новая версия go-модуля где хранится сгенеренный код pb с новой версией TaskCreatedV2
     - добавляется обработка новой версии TaskCreatedV2 параллельно версии TaskCreatedV1
     - правятся метрики (не в домашке :-( )
     - пишутся unit-тесты и e2e-тесты (не в домашке :-( )
     - тестируется при необходимости вручную, что обе версии событий обрабатываются

3. Выкатываем analytics в test stage и потом в prod

4. Готовим новую версию accounting:
    4.1 Скрипт миграции таблицы "tasks" в БД "accounting".

    Сейчас только столбец "description". Добавляется новый столбец "jira_id".
    При запуске сервис автоматически отмигрирует БД. Сервис также при миграции распарсит текущие записи и при возможности разнесет jira_id и description по разным столбцам

    4.2 Правится код:

     - модель данных в модуле для работы с БД
     - импортируется новая версия go-модуля где хранится сгенеренный код pb с новой версией TaskCreatedV2
     - добавляется обработка новой версии TaskCreatedV2 параллельно версии TaskCreatedV1
     - также будем изменена схема api самого сервиса accounting, в частности добавится новое поле jira_id. Там версию api не поднимаю, обратная совместимость не ломается, но по факту если клиент не обновит proto-файл и код, сгенеренный по нему, то он не увидит распарсенный jira_id, только title.
     - правятся метрики (не в домашке :-( )
     - пишутся unit-тесты и e2e-тесты (не в домашке :-( )
     - тестируется при необходимости вручную, что обе версии событий обрабатываются

5. Выкатываем accounting в test stage и потом в prod

6. Готовим новую версию taks-tracker:

    6.1 Скрипт миграции таблицы "tasks" в БД task-tracker.

    Сейчас столбец "description". Добавляется новый столбец "jira_id"

    6.2 Правится код:

     - модель данных в модуле для работы с БД
     - импортируется новая версия go-модуля где хранится сгенеренный код pb с новой версией TaskCreatedV2
     - при создании таски вместо TaskCreatedV1 отправляется TaskCreatedV2
     - правятся метрики (не в домашке :-( )
     - пишутся unit-тесты и e2e-тесты (не в домашке :-( )
     - тестируется при необходимости вручную, что обе версии событий отправляются корректно

7. Выкатываем taks-tracker в test stage и потом в prod

8. Тестируем в stage и prod

9. Убеждаемся, что в системе более нет событий TaskCreatedV1

10. Удаляем обработку старой версии TaskCreatedV1 из accounting

11. Удаляем обработку старой версии TaskCreatedV1 из analytics