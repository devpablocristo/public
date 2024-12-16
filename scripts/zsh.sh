#!/bin/bash

# Actualizar e instalar zsh
sudo apt update
sudo apt install -y zsh

# Cambiar el shell por defecto a zsh
chsh -s $(which zsh)

# Verificar si zsh está en /etc/shells
if ! grep -Fxq "$(which zsh)" /etc/shells
then
    echo "$(which zsh)" | sudo tee -a /etc/shells
fi

# Hacer zsh persistente para todos los usuarios (opcional)
echo "export SHELL=$(which zsh)" >> ~/.zshrc
echo "[[ -z "$ZSH_VERSION" ]] && exec $(which zsh) -l" >> ~/.bashrc

# Mensaje de finalización
echo "zsh se ha configurado como el shell por defecto de forma persistente. Cierra la sesión y vuelve a iniciarla para aplicar los cambios."

# (Opcional) Instalar Oh My Zsh
read -p "¿Deseas instalar Oh My Zsh? (s/n): " install_oh_my_zsh
if [ "$install_oh_my_zsh" = "s" ]; then
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
fi
