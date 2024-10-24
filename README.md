# go-musthave-shortener-tpl
# Сервис сокращения URL

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. Запустите main.go файл.
3. Укажите флаги или переменные среды окружения
   - ```-a``` - адрес вместе с портом для запуска сервера (default - ```localhost:8080```)
   - ```-b``` - base url. Например, ```http://myshorturl.ru```
   - ```-f``` - сохранять данные в памяти
   - ```-d```- database dsn
   - ```-s``` - включить генерацию самоподписывующихся сертификатов
   - ```-c``` - чтения конфига из json файла. ```path/to/config.json```
   - ```-t``` - ограничения доступа к хендлеру GET /api/internal/stats (trusted subnet)
   - ```-g``` - запустить gRPC сервер
Или ENV:
  - SERVER_ADDRESS
  - BASE_URL
  - FILE_STORAGE_PATH
  - DATABASE_DSN
  - ENABLE_HTTPS
  - TRUSTED_SUBNET
  - GRPC_SERVER
  - CONFIG

## URLs
```POST / ``` - short url  
```POST /api/shorten``` - request body ```{"url":"<some_url>"}``` => result ```{"result":"<short_url>"}```  
```POST /api/shorten/batch``` - req body ```[{"correlation_id": "<строковый идентификатор>", "original_url": "<URL для сокращения>"},...]``` => ```[{"correlation_id": "<строковый идентификатор из объекта запроса>", "short_url": "<результирующий сокращённый URL>"},...]```  
```GET /{id}``` - normal url  
```GET /ping``` - check db connecction  
```GET /api/user/urls``` - all users short urls. Only authorized users.  
```GET /api/internal/stats``` - quantity of users and short urls  
```DELETE /api/user/urls``` - soft delete  
  
gRPC сервер полностью дублирует HTTP хендлеры. Используется middleware для gzip шифрования. Также используются JWT токены для авторизации.

