# foodtosave-watcher

Monitora sacolas disponíveis na API do **Food To Save** e envia notificações via **Gotify** quando há novas sacolas, redução de quantidade ou esgotamento.

## Configuração

Copie o exemplo e ajuste os valores:

```bash
cp config.example.yaml config.yaml
```

| Campo | Descrição |
|-------|-----------|
| `interval` | Intervalo entre verificações (padrão: `5m`) |
| `cache_path` | Arquivo JSON de cache em disco (padrão: `./cache.json`) |
| `gotify.server_url` | URL do servidor Gotify |
| `gotify.app_token` | Token da Application (Apps no WebUI), não o token do Client |
| `merchants` | Lista de lojas com `id` e `name` |

O caminho do config pode ser definido pela flag `--config`, pela variável `CONFIG_PATH` ou pelo padrão `./config.yaml`.

## Como executar

**Local:**

```bash
make build
./bin/foodtosave-watcher --config ./config.yaml
```

**Desenvolvimento com hot reload ([Air](https://github.com/air-verse/air)):**

```bash
go install github.com/air-verse/air@latest
air
```

O projeto já inclui [`.air.toml`](.air.toml). O Air recompila e reinicia o watcher automaticamente ao salvar arquivos `.go`.

## Docker

A imagem é publicada automaticamente no [GitHub Container Registry](https://github.com/marcelotk15/foodtosave-watcher/pkgs/container/foodtosave-watcher) a cada push na branch `main`:

```bash
docker pull ghcr.io/marcelotk15/foodtosave-watcher:latest
```

**Executar com `docker run`:**

```bash
docker run -d \
  --name foodtosave-watcher \
  --restart unless-stopped \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v foodtosave-data:/app/data \
  ghcr.io/marcelotk15/foodtosave-watcher:latest
```

**Executar com Docker Compose** (build local a partir do [repositório](https://github.com/marcelotk15/foodtosave-watcher)):

```bash
git clone https://github.com/marcelotk15/foodtosave-watcher.git
cd foodtosave-watcher
cp config.example.yaml config.yaml
# edite config.yaml
docker compose up -d --build
```

No container, use `cache_path: /app/data/cache.json` no `config.yaml` para persistir o cache entre reinicializações.

## Como testar

```bash
go test ./...
```

Ou via Makefile:

```bash
make test
```

## Como contribuir

1. Faça fork do repositório e crie uma branch para sua mudança.
2. Antes do commit, rode `make fmt` e `make lint`.
3. Abra um PR descrevendo o motivo da mudança.
