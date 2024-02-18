# Установка
0. Для начала у вас должен стоять брокер сообщений [Apache Kafka](https://kafka.apache.org/quickstart).  
   Содержание docker-compose файла можно взять [отсюда](https://betterdatascience.com/how-to-install-apache-kafka-using-docker-the-easy-way/#:~:text=version%3A%20%273,%3A%20zookeeper%3A2181).
1. Следующим этапам необходимо настроить систему. 
В файле `.env` можно настроить следующие параметры:
    - KAFKA_URL  
	URL кафки
    - TIME_<operation>  
	Время на выполнение каждой операции. 
	Операции: ADD - сложение, SUBSTRACT - вычитание, MULT - умножение, DIVISION - деление.
	- AGENTS_CNT  
	Количество вычислятовров
2. Начинаем запускаться.   
   Дла начала запустите вычисляторы. Для каждого i-го вычислятора запуск будет выглядеть так:  
   `go run worker.go <i>`  
   Пример для 4 вычисляторов.
   ```
   ~ go run worker 1
   Ready!
   ~ go run worker 2
   Ready!
   ~ go run worker 3
   Ready!
   ~ go run worker 4
   Ready!
   ```  
   Дальше запустим орекстратор 
   ```
   ~ go run main.go
   I am listening at http://localhost:8080!
   ```
Всё готово! Теперь перейдём к API

# API
Обрабатываются мат. выражения как с целыми числами, так и с числами с плавающей точкой. 
Доступные операции: +, -, *, /.
Скобок нет.  

Методы
- Положить новое задание  
  POST /expression  
  Content-Type: application/json  
  {  
    "str_expr": <Математическое выражение>  
  }  
- Проверить готовность  
  GET /expression?id=<идентификатор выражения>

# Примеры
- Положить новое задание:
  ```
  curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"str_expr": "2 * 3 + 5 / 2"}' \
  http://localhost:8080/expression
  ```
  ```
  curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"str_expr": "2*3 + 5.98/2 - 0.001 + 21.1"}' \
  http://localhost:8080/expression
  ```
- Проверить готовность:
    ```
    curl -X GET http://localhost:8080/expression?id=123456789
    ```