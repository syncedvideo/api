### Local development with VSCode Remote Container

#### Windows WSL2

1. Clone repository to WSL filesystem (`npm run serve` hot reloading in `web` container does not work on Windows filesystem)
2. Run `sudo docker-compose up -d` (`sudo` is required to load environment variables from `.env` file) from WSL
3. Run `F1 > Remote-Containers: Open Folder in Container...` in VSCode
