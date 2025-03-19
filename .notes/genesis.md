Hola Claude, te tengo un desafio para ti. Y necesito que lo pienses con mucho cuidado.

La cosa va así:

- Necesitamos un server http que tenga soporte para templ (libreria de templates en go).
- El server debe tener un endpoint /dashboard, este endpoint debe renderizar un template html con la data que le pasemos.
- La data debe ser la siguiente:
    - Nombre de crypto: string
    - balance: string
    - last_updated_at: string
    muestra esta información en un tabla ordenada desde el mas reciente hasta el mas antiguo.
    la tabla debe tener un buscador para poder buscar por nombre de crypto.

Usa html, templ, tailwing css, go, y sqlite.