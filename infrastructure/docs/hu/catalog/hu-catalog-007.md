# HU-CATALOG-007: Gestión de Presentaciones de Productos

## 📌 Información General
- ID: HU-CATALOG-007
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2026-04-17

---

## 👤 Historia de Usuario

**Como** administrador/gerente de una empresa  
**Quiero** gestionar las presentaciones (variantes) de los productos  
**Para** poder manejar diferentes presentaciones como kilo, libra, etc. con precios distintos

---

## 🧠 Descripción Funcional

El sistema debe permitir crear y gestionar presentaciones (variantes) de productos, donde cada presentación puede tener:
- Nombre (kilo, libra, unidad, etc.)
- Factor de conversión a la unidad base
- Código de barras opcional
- Precio de costo y precio de venta individuales
- Banderas para compra y venta por defecto

---

## ✅ Criterios de Aceptación

### Escenario 1: Crear presentaciones exitosamente
- Dado que estoy autenticado como usuario con rol ADMIN o MANAGER
- Cuando envío una solicitud POST a `/catalog/products/:id/presentations` con una lista de presentaciones
- Entonces el sistema crea todas las presentaciones asociadas al producto
- Y retorna mensaje de éxito con la cantidad creada

### Escenario 2: Obtener presentaciones por producto
- Dado que las presentaciones existen para un producto
- Cuando consulto GET `/catalog/products/:id/presentations`
- Entonces el sistema devuelve todas las presentaciones del producto

### Escenario 3: Listar presentaciones con filtro por producto
- Dado que hay presentaciones en el sistema
- Cuando consulto GET `/presentations?product_id=123`
- Entonces el sistema devuelve solo las presentaciones del producto especificado

### Escenario 4: Crear o actualizar presentaciones (Upsert)
- Dado que quiero gestionar presentaciones de un producto
- Cuando envío una solicitud PUT a `/catalog/products/:id/presentations` con una lista de presentaciones
- Y algunas tienen `id` y otras no
- Entonces el sistema actualiza las que tienen `id` y crea las que no tienen
- Y retorna mensaje de éxito con la cantidad procesada

### Escenario 5: Actualizar presentación existente
- Dado que existe una presentación con ID=5 para un producto
- Cuando envío PUT con `{"id": 5, "name": "Nueva Libra", "factor": 0.5}`
- Entonces el sistema actualiza los datos de esa presentación

### Escenario 6: Crear nueva presentación en Upsert
- Dado que quiero agregar una nueva presentación
- Cuando envío PUT con `{"name": "Bulto", "factor": 10}`
- Entonces el sistema crea una nueva presentación

---

## ❌ Casos de Error

- **Producto no encontrado**: Si el product_id no existe, retorna error 400
- **Nombre duplicado**: Si el nombre de la presentación ya existe para el producto, retorna error 400
- **Campos requeridos faltantes**: Si name, factor, cost_price o sale_price faltan, retorna error 400

---

## 🔐 Reglas de Negocio

- Cada producto puede tener múltiples presentaciones
- El barcode es opcional pero debe ser único si se proporciona
- Solo name, factor, cost_price y sale_price son obligatorios
- Las presentaciones se crean asociadas a un product_id específico

---

## 📡 Requisitos Técnicos

- **Endpoints**:
  - `POST /catalog/products/:id/presentations` - Crear presentaciones
  - `GET /catalog/products/:id/presentations` - Listar por producto
  - `GET /presentations` - Listar todas (soporta filtro product_id)
  - `POST /presentations/page` - Paginado
  - `GET /presentations/:id` - Obtener por ID
  - `PUT /presentations/:id` - Actualizar
  - `DELETE /presentations/:id` - Eliminar

- **Request (crear)**:
```json
{
  "presentations": [
    {
      "name": "Kilo",
      "factor": 1,
      "barcode": "",
      "cost_price": 1100,
      "sale_price": 1200,
      "default_purchase": true,
      "default_sale": true
    },
    {
      "name": "Libra",
      "factor": 0.453592,
      "cost_price": 500,
      "sale_price": 550,
      "default_purchase": false,
      "default_sale": false
    }
  ]
}
```

- **Response**: 201 Created con mensaje y count

---

## 🧪 Criterios de Testing

- Tests unitarios para validaciones de servicio
- Tests de repositorio para CRUD
- Tests de handler para validación de request

---

## 📎 Dependencias

- Migración: 000006_presentations.up.sql
- Servicios: presentations.Service
- Módulo relacionado: EPICA-016 (Gestión de Catálogo) - productos