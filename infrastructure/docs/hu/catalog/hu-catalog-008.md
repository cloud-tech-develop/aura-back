# HU-CATALOG-008: Obtener Categoría por ID

## 📌 Información General
- ID: HU-CATALOG-008
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** usuario del sistema  
**Quiero** obtener los detalles de una categoría específica  
**Para** ver información completa de la categoría

---

## 🧠 Descripción Funcional

El sistema debe retornar la información de una categoría específica.

---

## ✅ Criterios de Aceptación

### Escenario 1: Categoría encontrada
- Dado que existe una categoría con ID 3
- Cuando solicito GET a `/categories/3`
-Entonces retorna la categoría

### Escenario 2: No encontrada
- Cuando solicito un ID inexistente
-Entonces retorna 404

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/categories/:id`
- **Método HTTP**: GET

---

## 📎 Dependencias

- Servicios: categories.Service.GetByID