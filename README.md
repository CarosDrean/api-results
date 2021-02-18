# API RESULTS

Api Results es el que se encarga de toda la gestion de datos.

# Instalacion

```
go get
```

# Configuracion

Para las configuracionde Base de Datos debera crear el archivo **configuration.json** con los siguientes campos:

```json
{
  "engine": "mssql",
  "server": "DREAN",
  "port": "1433",
  "user": "sa",
  "password": "123456",
  "database": "AMACHAY",
  "databaseaux": "HoloCovid"
}
```

# Compilacion

Para compilar el proyecto use:

```
go build
```
Para que el proyecto pueda funcionar correctamente no olvide el archivo **configuration.json**.

Ademas debera generar las firmas (**private.rsa** y **public.rsa.pub**) para generar el token con OpenSSh y acompañarlo del .exe.