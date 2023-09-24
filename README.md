# Проектная работа по курсу Golang Developer. Professional(OTUS) "Системный мониторинг"

## Общее описание
Реализация скраппера, собирающего информацию о системе, на которой запущен, и отправляющего её своим клиентам по GRPC.

## Собираемые данные скраппером
- Средняя загрузка системы (load average).

- Средняя загрузка CPU (%user_mode, %system_mode, %idle).

- Загрузка дисков:
    - tps (transfers per second);
    - KB/s (kilobytes (read+write) per second);

- Информация о дисках по каждой файловой системе:
    - использовано мегабайт, % от доступного количества;
    - использовано inode, % от доступного количества.

- Top talkers по сети:
    - по протоколам: protocol (TCP, UDP, ICMP, etc), bytes, % от sum(bytes) за последние **M**), сортируем по убыванию процента;
    - по трафику: source ip:port, destination ip:port, protocol, bytes per second (bps), сортируем по убыванию bps.

- Статистика по сетевым соединениям:
    - слушающие TCP & UDP сокеты: command, pid, user, protocol, port;
    - количество TCP соединений, находящихся в разных состояниях (ESTAB, FIN_WAIT, SYN_RCV и пр.).


## Запуск 

Клонируем репозиторий и переходим в папку:
```
git clone https://github.com/filatkinen/sysmon
cd sysmon

```

### Скраппер



 -Используя Makefile
```
#Без возможности сбора статистики(Top talkers)с использованием tcpdump
make run 
#Чтобы был доступен сбор через статистики tcpdump нужено запускать через sudo
make run-sudo 
```
 - Через командную строку:
 ```
 go run ./cmd/service/ -config configs/service.yaml
 ```

### Клиент
 - Используя Makefile
```
make run-client
```

 - Через командную строку:
 ```
go run ./cmd/client/ -M 5 -N 15 -address localhost -port 50051
 ```


## Описание  настроек скраппера
Файл с настройками в формате *.yaml
```
subsystems: #Сбор метрик: 
  la: true  # Средняя загрузка системы (load average)
  avgcpu: true #  Средняя загрузка CPU (%user_mode, %system_mode, %idle)
  disksload: true #  Загрузка дисков
  disksuse: true # Информация о дисках
  networkstat: true # Top talkers по сети
  networktop: true #Статистика по сетевым соединениям

scrapeinterval: 1s  #default time=5s  - интервал опроса метрик
cleaninterval: 1m  #default time=5m - интервал запуска очистки устаревших данных
depth: 2m  #default 1h - глубина хранения данных с метриками

bindings:
  port: 50051 # default = 50051  GRPC порт
  address: 0.0.0.0 #default 0.0.0.0 адресс для биндинга

``` 

## Описание работы

Скраппер собирает данные каждые **scrapeinterval** секунд.
Также работает очистка старых данных, глубина храннения определяется параметром **depth**


При появлении клиента, скраппер отдаем ему данных каждые **N** секунд.
Передаваемые данные кленту усреднены за **M** секунд.
Параметры **N** и **M** определяются клиентом при старте и передаются скрапперу.



## Поддерживаемые ОС
- Скраппер:
    - Linux - реализован сбор всех параметров
    - OS Mac(Darwin) - реализован сбор метрики средняя загрузка системы (load average). [ссылка на скриншоты запуска из Darwin](https://github.com/filatkinen/sysmon/tree/main/assets/approve_darwin)
- Клиент: поддерживается любая ОС, на которой скомплилируется бинарник

## Тестирование
#### Юнит-тесты

```
make test
```

#### Интеграционные тесты

```
make test-integration
```