1. Create a local PostgreSQL DB (default named **bus** but you can change it)
2. Create its tables with `postgresql/schema.sql`
3. Replace areguments in `env/local.env` with your own values
3. Run Server
```
$ make install
$ ./BusServer
```