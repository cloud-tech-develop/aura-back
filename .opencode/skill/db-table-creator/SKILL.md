---
name: db-table-creator
description: Asistente para la creación de tablas y migraciones siguiendo estándares del proyecto
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: database-migrations
---

## What I do

- Genero sentencias SQL para la creación de tablas y sus respectivos campos.
- Aseguro que los nombres de tablas y campos estén en **snake_case** e **inglés**.
- Agrego comentarios descriptivos en **español** para cada columna/tabla.
- Valido las relaciones entre el esquema **público** (`tenant/migrations/public`) y el esquema **privado** (`tenant/migrations/tenant`).
- Proporciono los comandos de `golang-migrate/v4` necesarios para generar los archivos de migración.

## When to use me

Utilízame cuando necesites añadir una nueva tabla o modificar una existente, asegurando que se mantenga la consistencia con la arquitectura multi-tenant del proyecto.

## Guidelines

1. **Naming**:
   - Tablas y columnas: `snake_case` e Inglés (ej. `order_items`, `quantity`).
   - Índices: `idx_schema_table_column`.
   - Constraints: `table_column_fk` o `table_column_unique`.

2. **Comments**:
   - Cada columna importante debe tener un comentario en español: `COMMENT ON COLUMN schema.table.column IS 'Descripción en español';`.

3. **Multi-tenancy and Relationships**:
   - **Public Schema**: Tables in `tenant/migrations/public` are shared across all tenants.
   - **Tenant Schema**: Tables in `tenant/migrations/tenant` are specific to each tenant.
   - **Validation Rule**: Before creating a foreign key in a tenant table that points to a public table, verify that the public table exists in `tenant/migrations/public`.
   - **Schema Prefix**: Always use `public.` when referencing a table in the public schema from a tenant migration.
   - **FK Naming**: ForeignKey constraints should follow `[table]_[column]_fk`.

4. **Migration Tool**:
   - Para crear la migración: `migrate create -ext sql -dir tenant/migrations/[public|tenant] -seq migration_name`.

## Example Interaction

**User**: "Necesito una tabla para inventario de productos en el esquema de tenant."
**Agent**: (Siguiendo esta skill) "Generaré la tabla `product_inventory` en `tenant/migrations/tenant`. Los campos serán `id`, `product_id` (ref a `product.id`), `quantity`, etc. Los comentarios estarán en español."
