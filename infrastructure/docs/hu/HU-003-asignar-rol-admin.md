# HU-003: Asignar rol de administrador al usuario inicial

**Como** sistema  
**Quiero** asignar el rol ADMIN al primer usuario creado  
**Para** garantizar que el fundador tenga todos los permisos

---

## Criterios de Aceptación

- [x] Al crear el usuario inicial, insertar automáticamente en `public.user_roles`
- [x] Asignar el rol con ID 1 (ADMIN) al usuario fundador
- [x] Verificar que la tabla `roles` tenga el rol ADMIN seedeado

---

## Estado: ✅ 3/3 IMPLEMENTADO

---

## Notas Técnicas

- En `tenant/manager.go`, después de crear el usuario, se inserta en `public.user_roles`
- El rol ADMIN tiene ID 1
