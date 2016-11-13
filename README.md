PoC for golang and Cassandra
============================

POST http://localhost/api/notes/ 

```
{
    "title": "buy the book" 
    "description": "Down the Highway: The Life of Bob Dylan"
}
```


GET http://localhost/api/notes/384eef71-3173-40a9-aed5-223d326e6fe4
```
{
    "createdon": 1479051483, 
    "description": "Down the Highway: The Life of Bob Dylan", 
    "id": "384eef71-3173-40a9-aed5-223d326e6fe4", 
    "title": "buy the book"
}
```

Cassandra scheme:
```
CREATE TABLE razvan.table1 (
	id uuid,
	createdon timestamp,
	title text,
	text text,
	PRIMARY KEY (id)
);

```
