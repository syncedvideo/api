# Postgres naming conventions

> {tablename}_{columnname(s)}_{suffix}

- `pkey` for a Primary Key constraint
- `key` for a Unique constraint
- `excl` for an Exclusion constraint
- `idx` for any other kind of index
- `fkey` for a Foreign key
- `check` for a Check constraint

[Found here](https://stackoverflow.com/questions/4107915/postgresql-default-constraint-names/4108266#4108266)
