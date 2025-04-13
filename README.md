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
CREATE TABLE users (
  id    uniqueidentifier,
  name  varchar,
  age   smallint
)
```

```sql
INSERT INTO users (id, name, age) VALUES ('e28c20d7-483d-4f6e-9b31-9d0d6819ba39', 'test', '18')
```

```sql
SELECT id, name, age FROM users

SELECT id, name, age FROM users WHERE age >= 1
```