# Offline Module - API Documentation

## Resumen

Módulo para sincronización offline de empresas. Estos endpoints solo funcionan en modo offline (SQLite).

## Endpoints

### 1. GET /offline/ping

Sincroniza empresas desde producción hacia SQLite local.

#### OpenAPI Specification

```yaml
/offline/ping:
  get:
    summary: "Sincronizar empresas desde producción"
    description: |
      Consume la ruta /enterprises desde URL_PROD (configurada en .env) y sincroniza 
      empresas que no existan en SQLite local. Solo funciona en modo offline (SQLITE).
    tags:
      - offline
    operationId: ping
    responses:
      200:
        description: "Sincronización completada"
        content:
          application/json:
            schema:
              type: object
              properties:
                synced:
                  type: integer
                  description: "Número de empresas sincronizadas"
                source:
                  type: string
                  description: "URL de origen de producción"
                mode:
                  type: string
                  description: "Modo de funcionamiento"
                  enum: [offline]
                message:
                  type: string
                  description: "Mensaje de estado"
            example:
              {
                "synced": 5,
                "source": "http://localhost:8081",
                "mode": "offline",
                "message": "Sincronización completada"
              }
      400:
        description: "Error al sincronizar"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Error'
      403:
        description: "Endpoint no disponible en modo online"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Error'
```

#### Response Codes

| Status Code | Description | Body |
|------------|-------------|------|
| 200 | OK - Sincronización exitosa | `{synced: int, source: string, mode: "offline", message: string}` |
| 400 | Bad Request - Error al sincronizar | `{success: false, message: string}` |
| 403 | Forbidden - No disponible en modo online | `{success: false, message: string}` |

#### Example Request

```http
GET /offline/ping HTTP/1.1
Host: localhost:8080
Accept: application/json
```

#### Example Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": {
    "synced": 5,
    "source": "http://localhost:8081",
    "mode": "offline",
    "message": "Sincronización completada"
  },
  "message": "Operación exitosa",
  "success": true
}
```

---

### 2. GET /offline/enterprises

Lista empresas guardadas localmente en SQLite.

#### OpenAPI Specification

```yaml
/offline/enterprises:
  get:
    summary: "Listar empresas guardadas localmente"
    description: |
      Retorna todas las empresas almacenadas en la base de datos SQLite local.
      Solo funciona en modo offline (SQLite).
    tags:
      - offline
    operationId: listEnterprises
    responses:
      200:
        description: "Lista de empresas"
        content:
          application/json:
            schema:
              type: object
              properties:
                data:
                  type: array
                  items:
                    $ref: '#/components/schemas/Enterprise'
                total:
                  type: integer
                  description: "Total de empresas"
                source:
                  type: string
                  description: "Fuente de los datos"
                  enum: [local]
            example:
              {
                "data": [
                  {
                    "id": 1,
                    "tenant_id": 100,
                    "name": "Empresa Uno SAS",
                    "commercial_name": "Tienda Uno",
                    "slug": "empresa_uno",
                    "sub_domain": "empresa-uno",
                    "email": "contacto@empresauno.com",
                    "document": "900123456",
                    "dv": "1",
                    "phone": "+57300123456",
                    "municipality_id": "05001",
                    "municipality": "Medellín",
                    "status": "ACTIVE",
                    "settings": {},
                    "created_at": "2024-01-15T10:00:00Z",
                    "updated_at": "2024-01-15T10:00:00Z",
                    "deleted_at": null
                  }
                ],
                "total": 1,
                "source": "local"
              }
      400:
        description: "Error al listar empresas"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Error'
      403:
        description: "Endpoint no disponible en modo online"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Error'
```

#### Response Codes

| Status Code | Description | Body |
|------------|-------------|------|
| 200 | OK - Lista exitosa | `{data: [...Enterprise], total: int, source: "local"}` |
| 400 | Bad Request - Error al listar | `{success: false, message: string}` |
| 403 | Forbidden - No disponible en modo online | `{success: false, message: string}` |

#### Example Request

```http
GET /offline/enterprises HTTP/1.1
Host: localhost:8080
Accept: application/json
```

#### Example Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": {
    "data": [
      {
        "id": 1,
        "tenant_id": 100,
        "name": "Empresa Uno SAS",
        "commercial_name": "Tienda Uno",
        "slug": "empresa_uno",
        "sub_domain": "empresa-uno",
        "email": "contacto@empresauno.com",
        "document": "900123456",
        "dv": "1",
        "phone": "+57300123456",
        "municipality_id": "05001",
        "municipality": "Medellín",
        "status": "ACTIVE",
        "settings": {},
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z",
        "deleted_at": null
      }
    ],
    "total": 1,
    "source": "local"
  },
  "message": "Operación exitosa",
  "success": true
}
```

---

## Schemas

### Enterprise

```yaml
components:
  schemas:
    Enterprise:
      type: object
      properties:
        id:
          type: integer
          description: "ID único de la empresa"
        tenant_id:
          type: integer
          description: "ID del tenant en producción"
        name:
          type: string
          description: "Razón social"
        commercial_name:
          type: string
          description: "Nombre comercial"
        slug:
          type: string
          description: "Slug único para la empresa"
        sub_domain:
          type: string
          description: "Subdominio"
        email:
          type: string
          description: "Correo electrónico"
        document:
          type: string
          description: "Número de documento"
        dv:
          type: string
          description: "Dígito de verificación"
        phone:
          type: string
          description: "Teléfono"
        municipality_id:
          type: string
          description: "ID del municipio"
        municipality:
          type: string
          description: "Nombre del municipio"
        status:
          type: string
          description: "Estado de la empresa"
          enum: [ACTIVE, INACTIVE, SUSPENDED]
        settings:
          type: object
          description: "Configuración adicional"
        created_at:
          type: string
          format: date-time
          description: "Fecha de creación"
        updated_at:
          type: string
          format: date-time
          description: "Fecha de actualización"
        deleted_at:
          type: string
          format: date-time
          nullable: true
          description: "Fecha de eliminación"
```

### Error

```yaml
components:
  schemas:
    Error:
      type: object
      properties:
        success:
          type: boolean
          description: "Indica si la operación fue exitosa"
          enum: [false]
        message:
          type: string
          description: "Mensaje de error"
```

---

## Modo de Funcionamiento

### Offline Mode (SQLite)

- Driver: `sqlite`
- Archivo: Definido en `DATABASE_URL`
- Endpoints disponibles: `/offline/ping`, `/offline/enterprises`
- Fuente de datos: `URL_PROD` (configurada en `.env`)

### Online Mode (PostgreSQL) - No Disponible

- Driver: `postgres`
- Endpoints offline retornan `403 Forbidden`

---

## Environment Variables

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| DATABASE_DRIVER | Driver de base de datos | `sqlite` o `postgres` |
| DATABASE_URL | URL de conexión | `./offline.db` |
| URL_PROD | URL de producción | `http://localhost:8081` |