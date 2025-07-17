# gorm-paginator

**gorm-paginator** es una librería ligera y genérica para agregar paginación a tus consultas usando [GORM](https://gorm.io/) en Go.
Permite paginar cualquier modelo, incluir relaciones con `Preload` y devolver metadata útil para construir APIs RESTful con soporte para navegación y ordenamiento.

---

## ✨ Características

- Compatible con cualquier modelo GORM
- Soporta paginación con parámetros `page` y `limit`
- Incluye metadata: total de registros, páginas totales, si hay siguiente o anterior página
- Permite ordenar resultados dinámicamente (`OrderBy`)
- Soporta preload de relaciones
- Basado en `reflect` para manejar tipos dinámicos sin perder seguridad de tipos

---

## 🚀 Instalación

```bash
go get github.com/tuusuario/gorm-paginator
```

---

## 📦 Uso Básico
```go
import (
    "github.com/hiro-riveros/gorm-paginator"
    "gorm.io/gorm"
)

type User struct {
    ID    uint
    Name  string
    Email string
}

func ListUsers(db *gorm.DB, page, limit int) ([]*User, paginator.Metadata, error) {
    params := paginator.Params{
        Page:  page,
        Limit: limit,
        OrderBy: "created_at desc",
    }
    usersInterface, meta, err := paginator.Paginate(&User{}, db.Model(&User{}), params)
    if err != nil {
        return nil, meta, err
    }

    users := usersInterface.([]*User)
    return users, meta, nil
}
```

---

## ⚙️ Parámetros
| Campo     | Tipo     | Descripción                                |
| --------- | -------- | ------------------------------------------ |
| `Page`    | `int`    | Número de página (comienza en 1)           |
| `Limit`   | `int`    | Número de registros por página             |
| `OrderBy` | `string` | Ordenamiento SQL (ej: `"created_at desc"`) |

## 📊 Metadata devuelta
La estructura Metadata contiene información útil para la paginación:
```go
type Metadata struct {
    Page         int   `json:"page"`
    Limit        int   `json:"limit"`
    TotalRecords int64 `json:"total_records"`
    TotalPages   int   `json:"total_pages"`
    HasNext      bool  `json:"has_next"`
    HasPrev      bool  `json:"has_prev"`
}
```

---

## 🔄 Soporte para Preload
Puedes pasar nombres de relaciones para cargar con Preload:
```go
users, meta, err := paginator.Paginate(&User{}, db.Model(&User{}), params, "Wallets", "Transactions")
```

---

## Ejemplo completo

```go
package main

import (
    "fmt"
    "log"

    "github.com/hiro-riveros/gorm-paginator"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID    uint
    Name  string
    Email string
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&User{})

    for i := 1; i <= 50; i++ {
        db.Create(&User{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)})
    }

    params := paginator.Params{Page: 1, Limit: 10, OrderBy: "id desc"}

    usersInterface, meta, err := paginator.Paginate(&User{}, db.Model(&User{}), params)
    if err != nil {
        log.Fatal(err)
    }

    users := usersInterface.([]*User)
    fmt.Printf("Página %d de %d (Total: %d usuarios)\n", meta.Page, meta.TotalPages, meta.TotalRecords)
    for _, u := range users {
        fmt.Println(u.ID, u.Name)
    }
}
```

---

## 🤝 Contribuciones
¡Las contribuciones son bienvenidas!
Si quieres sugerir mejoras, reportar bugs o añadir funcionalidades, abre un issue o pull request.

## 📄 Licencia
MIT License © [@hiro-riveros]
