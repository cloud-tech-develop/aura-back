# HU-CATALOG-003: Obtener Producto por ID

## 📌 Información General
- ID: HU-CATALOG-003
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** usuario del sistema  
**Quiero** obtener los detalles de un producto específico  
**Para** poder ver toda la información del producto

---

## 🧠 Descripción Funcional

El sistema debe retornar la información completa de un producto específico incluyendo todos sus atributos.

---

## ✅ Criterios de Aceptación

### Escenario 1: Producto encontrado
- Dado que existe un producto con ID 5 en la empresa
- Cuando solicito GET a `/products/5`
- Entonces el sistema retorna el producto con todos sus campos

### Escenario 2: Producto no encontrado
- Dado que no existe un producto con ese ID
- Cuando solicito GET a `/products/999`
- Entonces el sistema retorna error 404

---

## ❌ Casos de Error

- **Producto no encontrado**: Retorna 404 con mensaje "Producto no encontrado"
- **ID inválido**: Si el ID no es un número válido, retorna 400

---

## 🔐 Reglas de Negocio

- Solo se retorna el producto si pertenece a la empresa del usuario
- No se retornan productos con status DELETED

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products/:id`
- **Método HTTP**: GET
- **Response**: 200 OK con el producto

---

## 📎 Dependencias

- Servicios: products.Service.GetByID