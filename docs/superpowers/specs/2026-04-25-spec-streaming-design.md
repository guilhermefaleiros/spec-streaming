# Spec Streaming Design

## Contexto

`spec-streaming` e um prototipo simples de plataforma de streaming com tres capacidades principais:

- ingestao de videos MP4 via interface web,
- transcoding assincrono para MPEG-DASH em multiplas qualidades,
- listagem e reproducao de videos processados.

O objetivo e manter o MVP pequeno, com arquitetura clara o suficiente para evoluir sem introduzir infraestrutura desnecessaria logo no inicio.

## Escopo

### Incluido

- upload de videos MP4 pelo front-end
- persistencia de metadados de videos e jobs no PostgreSQL
- processamento assincrono de transcoding no back-end
- geracao de artefatos MPEG-DASH com qualidades 480p, 720p e 1080p
- listagem simples de videos
- pagina de reproducao de video
- acompanhamento de status por polling
- suporte de storage com implementacoes local e S3-compatible
- serving do manifesto e segmentos DASH pela API

### Fora do escopo

- autenticacao, login ou cadastro
- geracao de legendas
- cache
- upload direto para o storage no MVP
- fila externa dedicada
- microsservicos separados

## Stack

### Back-end

- Go
- Echo
- PostgreSQL
- ffmpeg para transcoding

### Front-end

- React com Vite
- TypeScript
- shadcn/ui

### Testes E2E

- Playwright com `@playwright/test`

## Arquitetura

O sistema sera implementado como um monolito modular com dois processos principais no mesmo projeto:

1. API back-end
2. worker de transcoding

O front-end fica separado como aplicacao web consumindo a API.

### Componentes

#### Front-end

Responsavel por:

- enviar uploads de MP4
- consultar lista de videos
- atualizar status por polling
- abrir a pagina de reproducao
- inicializar o player com o manifesto DASH

#### API

Responsavel por:

- receber upload e validar entrada basica
- salvar o arquivo original via abstracao de storage
- criar registros de video e job
- expor endpoints de listagem, detalhe e status
- servir manifesto e segmentos DASH ao player

#### Worker

Responsavel por:

- buscar jobs pendentes no banco
- marcar job como `processing`
- ler o MP4 original do storage
- executar transcoding com `ffmpeg`
- gerar manifesto e segmentos DASH
- publicar artefatos no storage
- atualizar estado final de job e video

#### Storage Service

O sistema deve depender de uma interface de storage em vez de acoplamento a filesystem ou S3.

Operacoes esperadas:

- `SaveSource`
- `OpenSource`
- `SaveArtifact`
- `OpenArtifact`
- `ArtifactExists`

Implementacoes do MVP:

- local filesystem para desenvolvimento
- S3-compatible para ambientes que exigirem object storage

#### Services de aplicacao

Camadas de aplicacao encapsulam regras de negocio e mantem handlers HTTP e worker enxutos.

Servicos principais:

- `VideoService`
- `TranscodingJobService`
- `TranscoderService`

## Fluxo ponta a ponta

1. O usuario envia um MP4 pelo front-end.
2. A API valida o payload e grava o arquivo original no storage.
3. A API cria um registro em `videos` com status inicial e um registro em `transcoding_jobs` com status `pending`.
4. O worker consulta o banco por jobs pendentes.
5. Ao capturar um job, o worker marca o job como `processing` e o video como `processing`.
6. O worker le o arquivo original do storage e executa o transcoding para 480p, 720p e 1080p, gerando manifesto e segmentos DASH.
7. O worker salva os artefatos gerados no storage.
8. Em caso de sucesso, o worker marca o job como `completed` e o video como `ready`, registrando o caminho logico do manifesto.
9. Em caso de falha, o worker marca job e video como `failed`, persistindo a mensagem de erro.
10. O front-end consulta status por polling e atualiza a listagem.
11. Quando um video estiver `ready`, o usuario abre a pagina do player.
12. O player requisita manifesto e segmentos via API, que le os arquivos do storage e os entrega ao cliente.

## Modelo de dominio

### Video

Representa o ativo principal do produto.

Campos esperados:

- `id`
- `title`
- `original_filename`
- `status`
- `source_storage_key`
- `manifest_storage_key`
- `duration_seconds` opcional
- `error_message` opcional
- `created_at`
- `updated_at`

### TranscodingJob

Representa a unidade operacional assincrona de processamento.

Campos esperados:

- `id`
- `video_id`
- `status`
- `attempts`
- `started_at` opcional
- `finished_at` opcional
- `error_message` opcional
- `created_at`
- `updated_at`

Separar video e job evita misturar estado de catalogo com estado operacional.

## Estados

### `video.status`

- `uploaded`
- `processing`
- `ready`
- `failed`

### `transcoding_jobs.status`

- `pending`
- `processing`
- `completed`
- `failed`

Transicoes invalidas devem ser rejeitadas pela camada de servico.

## API

### `POST /videos`

Recebe upload do MP4 e metadados simples, como titulo.

Responsabilidades:

- validar entrada minima
- salvar arquivo original
- criar video
- criar job pendente

Resposta esperada:

- identificador do video
- status inicial
- metadados basicos

### `GET /videos`

Lista videos com metadados basicos e status atual.

### `GET /videos/:id`

Retorna detalhes do video, incluindo status, mensagens de erro e disponibilidade logica do manifesto.

### `GET /videos/:id/status`

Endpoint leve para polling de status.

### `GET /videos/:id/stream/manifest.mpd`

Serve o manifesto DASH para um video pronto.

### `GET /videos/:id/stream/*`

Serve segmentos e demais arquivos associados ao streaming DASH.

## Decisoes de arquitetura

### Upload pelo back-end no MVP

O upload passara pela API neste primeiro ciclo. Isso simplifica a implementacao e reduz o numero de componentes externos. A interface de storage deve, no entanto, preservar uma evolucao futura para upload direto ao storage sem quebrar o dominio.

### Worker separado da API

O processamento assincrono rodara em processo separado do servidor HTTP. Isso evita misturar carga de CPU e I/O pesado do transcoding com trafego da API e torna o comportamento operacional mais previsivel.

### Polling no front-end

O acompanhamento de status sera feito por polling, nao por WebSocket ou SSE. Para o MVP, isso reduz complexidade sem prejudicar a experiencia principal.

### Serving pela API

Mesmo com suporte a storage local e S3-compatible, o cliente consumira manifesto e segmentos via API. Isso mantem um contrato unico para reproducao e simplifica o suporte aos dois tipos de storage.

## Tratamento de erros

Regras principais:

1. se o upload falhar, a API nao deve deixar video ou job orfaos persistidos sem consistencia
2. se o transcoding falhar, `video.status` e `job.status` devem refletir falha e expor `error_message`
3. se o video nao estiver `ready`, endpoints de streaming devem retornar erro de negocio adequado
4. o worker deve impedir processamento simultaneo do mesmo job
5. arquivos temporarios de transcoding devem ser removidos no fim da execucao, com sucesso ou falha

## Estrategia de testes

### Unitarios

- transicoes de status
- regras de servico
- abstracao de storage
- composicao de paths e chaves logicas

### Integracao back-end

- upload criando video e job
- persistencia correta no PostgreSQL
- worker consumindo job e atualizando estado
- serving do manifesto e segmentos via API

### Front-end

- fluxo de upload
- polling de status
- listagem de videos
- navegacao para pagina do player

### E2E com Playwright

O fluxo end-to-end deve usar `@playwright/test` como runner principal.

Diretrizes:

- configurar `baseURL` para a aplicacao web
- configurar `webServer` no `playwright.config.ts` para subir a aplicacao antes dos testes
- iniciar com projeto `chromium` no MVP
- confiar no isolamento por teste provido por `BrowserContext`
- habilitar `trace: 'on-first-retry'`
- habilitar `screenshot: 'only-on-failure'`
- habilitar `video: 'retain-on-failure'`

Cenario E2E principal:

1. abrir a aplicacao
2. enviar um MP4 curto
3. acompanhar o status ate `ready`
4. abrir a pagina do video
5. validar carregamento do player
6. validar inicio do playback

## Estrutura sugerida

Uma estrutura inicial coerente com o design seria:

- `backend/api`
- `backend/worker`
- `backend/internal/videos`
- `backend/internal/transcoding`
- `backend/internal/storage`
- `frontend`

Os nomes exatos podem mudar para refletir convencoes do projeto, desde que a separacao entre dominio, infraestrutura e interfaces externas seja preservada.

## Riscos e limites do MVP

- transcoding pode ser lento em ambiente local dependendo do tamanho do arquivo
- serving pela API adiciona carga no back-end, mas e aceitavel para prototipo
- sem fila externa, a recuperacao de jobs depende de logica cuidadosa de consulta e travamento no banco
- sem autenticacao, qualquer usuario do ambiente consegue enviar e reproduzir videos

## Resultado esperado

Ao final do MVP, o sistema deve permitir:

1. enviar um video MP4 por uma interface web
2. acompanhar seu processamento assincrono
3. listar videos cadastrados
4. abrir um video pronto
5. reproduzir streaming MPEG-DASH com selecao de qualidade entre 480p, 720p e 1080p
