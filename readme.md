Для запуска выполнить  
go build .\cmd\server.go для Windows  
или go build ./cmd/server.go для Linux  
в корневой директории проекта.    
Задать в переменную окружения PORT.
Номер порта задать больше 1000.  
Запустить получившийся бинарник,  
открыть http://localhost:указанный номер порта  
По умолчанию данные хранятся в Mongo Atlas,  
для инмемори БД добавить ключ -inmemory