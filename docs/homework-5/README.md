# Домашние задания по модулю "Многопоточное программирование"
Дедлайн:  6 апреля, 23:59 (сдача) / 9 апреля, 23:59 (проверка)

## Основное задание
Уменьшить время ответа cart.list.

    - Распараллельные вызовы ручки http://route256.pavl.uk:8080/docs/#/ProductService/ProductService_GetProduct продакт
      сервиса
    - Самим написать аналог https://pkg.go.dev/golang.org/x/sync/errgroup и использовать его
    - В случае ошибки - отменять все текущие запросы и вернуть ошибку из errgroup
    - при общении с Product Service необходимо использовать лимит 10 RPS на клиентской стороне
    - група живет в рамках одного запроса = группа не переиспользуется между запросами
    - инмемори репозитории защитить мьютексами
    - требуется наличие комментариев в коде

## Дополнительное задание (на 10 баллов)
    - тесты на многопоточность in-memory репозитория
    - грейсфул завершение приложения (использвать signal.Notify)
    - контроль утечек (горутин через https://github.com/uber-go/goleak)
    - контроль рейзов (через тесты с флагом -race)
    - реализовать параллельный запуск тестов
