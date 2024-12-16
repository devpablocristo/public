NOTA: como 'stty -ixon', a principio del zsh, para evitar conflicos con Powerlevel10k

En `zsh`, cuando los logs se detienen al presionar `Ctrl + S`, es porque ese comando pausa la salida de la terminal debido al control de flujo XON/XOFF. Para reanudar la salida, usa `Ctrl + Q`.

### Solución

1. **Evita presionar `Ctrl + S`** para no pausar la salida.
2. **Presiona `Ctrl + Q`** para reanudar la salida si se ha detenido.
3. **Desactiva XON/XOFF** añadiendo `stty -ixon` a tu archivo `.zshrc`:

   ```sh
   echo 'stty -ixon' >> ~/.zshrc
   source ~/.zshrc
   ```

Esto evitará que `Ctrl + S` pause la terminal en el futuro.