### Auralis
Simple relational database with pgql compatible sql dialect

---

##### TODO
- [x] querying specific fields
  - [ ] any column order
  - [ ] aliasing
  - [x] query result set pretty print
- [x] basic data structures
  - [ ] variable character types
- [ ] data paging
- [ ] DML
- [ ] DDL
- [ ] indexing
- [ ] nullable columns
- [ ] TCL
- [ ] expose server
- [x] basic lexer
- [ ] custom schema

---
Example supported queries

```sql
CREATE TABLE users (  id    uniqueidentifier,  name  varchar,  age   smallint)
```

```sql
INSERT INTO users (id, name, age) VALUES ('e28c20d7-483d-4f6e-9b31-9d0d6819ba39', 'test', '18')
```

```sql
SELECT id, name, age FROM users

SELECT id, name, age FROM users WHERE age >= 1
```


```sql 
-- query metadata for tables
SELECT * FROM auralis.tables

+---+---------------+--------------+--------------+
|   | database_name | table_schema | table_name   |
+---+---------------+--------------+--------------+
| 1 | auralis       | auralis      | tables       |
| 2 | auralis       | auralis      | columns      |
| 3 | test-database | dbo          | users        |
+---+---------------+--------------+--------------+

-- query metadata for columns
SELECT * FROM auralis.columns

+---+--------------+--------------+---------------+-----------+----------+
|   | table_schema | table_name   | column_name   | data_type | position |
+---+--------------+--------------+---------------+-----------+----------+
| 1 | auralis      | tables       | database_name | varchar   | 1        |
| 2 | auralis      | tables       | table_schema  | varchar   | 2        |
| 3 | auralis      | tables       | table_name    | varchar   | 3        |
| 4 | auralis      | columns      | table_schema  | varchar   | 1        |
| 5 | auralis      | columns      | table_name    | varchar   | 2        |
| 6 | auralis      | columns      | column_name   | varchar   | 3        |
| 7 | auralis      | columns      | data_type     | varchar   | 4        |
| 8 | auralis      | columns      | position      | smallint  | 5        |
| 9 | dbo          | users        | age           | smallint  | 1        |
+---+--------------+--------------+---------------+-----------+----------+
```