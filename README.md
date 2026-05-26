# foodtosave-watcher

Monitor periódico de sacolas (gondolas) disponíveis na API do **Food To Save** por loja (merchant), com notificações via **ntfy** quando há novas sacolas, redução de quantidade ou esgotamento.

## Requisitos

- Go 1.24+
- Arquivo `config.yaml` (ou caminho via `--config` / `CONFIG_PATH`)

## Configuração (`config.yaml`)

Copie o exemplo e ajuste:

```bash
cp config.example.yaml config.yaml
```

Exemplo:

```yaml
interval: 5m

cache_path: ./cache.json

ntfy:
  topic: foodtosave-alerts
  server_url: https://ntfy.sh

merchants:
  - id: "merchant-id-1"
    name: "Loja Centro"
  - id: "merchant-id-2"
    name: "Loja Bairro"
```

- **`interval`**: opcional; padrão `5m` se omitido.
- **`cache_path`**: opcional; padrão `./cache.json`. Caminho do arquivo JSON de cache em disco (veja abaixo).
- **`ntfy.topic`**: obrigatório (não vazio).
- **`ntfy.server_url`**: opcional; padrão `https://ntfy.sh`.
- **`merchants`**: pelo menos um item; cada um precisa de `id` e `name`.

Ordem de resolução do arquivo de config:

1. Flag `--config`
2. Variável de ambiente `CONFIG_PATH`
3. `./config.yaml` (local) ou `/app/config.yaml` no container (via `CMD` do Dockerfile)

## Rodar localmente

Build:

```bash
go build -o bin/foodtosave-watcher ./cmd/foodtosave-watcher
```

Executar:

```bash
./bin/foodtosave-watcher --config ./config.yaml
```

Ou com variável de ambiente:

```bash
CONFIG_PATH=./config.yaml ./bin/foodtosave-watcher
```

## Testes

```bash
go test ./...
```

## Docker

Build:

```bash
docker build -t foodtosave-watcher .
```

Run (montando o config):

```bash
docker run --rm \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  foodtosave-watcher
```

Com `CONFIG_PATH`:

```bash
docker run --rm \
  -e CONFIG_PATH=/app/config.yaml \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  foodtosave-watcher
```

## Docker Compose

```bash
docker compose up --build
```

O `docker-compose.yml` inclui um volume nomeado `foodtosave-data` montado em `/app/data` dentro do container. Para persistir o cache entre reinicializações, aponte `cache_path` para esse diretório no seu `config.yaml`:

```yaml
cache_path: /app/data/cache.json
```

## Tópicos no ntfy

1. Escolha um nome de tópico único em [ntfy.sh](https://ntfy.sh) (ou seu servidor ntfy) e configure em `ntfy.topic`.
2. No app ntfy (Android/iOS) ou navegador, inscreva-se nesse mesmo tópico para receber as notificações.
3. Como qualquer um que souber o nome do tópico pode enviar mensagens para tópicos públicos, use um nome longo e difícil de adivinhar, ou hospede ntfy com autenticação.

## Cache em disco

O estado das gondolas é salvo em um arquivo JSON (`cache_path`). A cada ciclo de monitoramento, após processar os eventos, o arquivo é gravado de forma atômica (escrita em `.tmp` + rename). Isso garante que:

- Se a aplicação reiniciar, ela não renotifica sacolas que já foram notificadas.
- Corrupção parcial do arquivo é evitada pelo mecanismo de rename atômico.

Se o arquivo estiver ausente na primeira execução (ou corrompido), a aplicação inicia com cache vazio e um aviso é logado.

## Tipos de notificação

| Situação | Título do header |
|----------|-------------------|
| Gondola nova para o merchant | `Nova sacola disponível` |
| Quantidade diminuiu (mesma gondola) | `Sacolas diminuíram` |
| Gondola sumiu da resposta da API | `Sacola esgotada` |

O corpo das mensagens segue os modelos descritos na documentação interna da aplicação (preço em BRL e datas em RFC3339).

## API monitorada

```
GET https://api.foodtosave.com.br/api/v1/merchants/{MERCHANT_ID}/gondolas
```

## Licença

Uso pessoal / exemplo.
