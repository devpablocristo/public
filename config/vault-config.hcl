# Configuración del backend de almacenamiento
storage "file" {
  path = "/vault/data"
}

# Configuración del listener
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

# Habilitar la interfaz de usuario
ui = true
