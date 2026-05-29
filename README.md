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

**Docker Compose (recomendado para produção):**

```bash
docker compose up --build
```

No container, use `cache_path: /app/data/cache.json` para persistir o cache entre reinicializações.

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
