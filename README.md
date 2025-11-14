## Conexion a la base

Para conectarte a la base de datos PostgreSQL usando psql:

### PostgreSQL

```
docker exec -it postgres-db psql -U <USER_DB> -d <DB_NAME>
```

(Reemplaza <USER_DB> y <DB_NAME> con los valores de tu archivo .env).