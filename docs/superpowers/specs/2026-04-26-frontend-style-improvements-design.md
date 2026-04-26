# Frontend Style Improvements Design

## Contexto

O projeto `spec-streaming` ja possui fluxo funcional de upload, catalogo e reproducao, mas o front-end atual ainda esta com:

- estilos inline e linguagem visual basica,
- acoplamento de estado assicrono em `useEffect` manual,
- player baseado em `dashjs` (fora do requisito novo),
- falta de design system consistente para escalabilidade.

Com base em `docs/raw/STYLE.MD`, esta etapa define melhorias de UX/UI e arquitetura de front-end mantendo o back-end atual.

## Objetivo

Entregar uma experiencia moderna inspirada na paleta Netflix, com 3 paginas dedicadas e navegacao clara, adotando obrigatoriamente:

- Shaka Player para reproducao MPEG-DASH,
- shadcn/ui para componentes e estilização,
- TanStack Query (React Query) para integracao com APIs e estados assincros.

## Escopo

### Incluido

- layout application shell com sidebar no desktop e drawer no mobile,
- rotas separadas para upload, catalogo e player,
- tema visual consistente (tokens, tipografia, contraste, estados),
- migracao das chamadas de API para React Query,
- polling declarativo com `refetchInterval`,
- substituicao de `dashjs` por Shaka Player,
- melhoria de feedback visual para loading, erro, vazio e sucesso,
- cobertura de testes de componentes e fluxo principal de navegacao/player.

### Fora do escopo

- migracao para TanStack Start,
- alteracoes estruturais no contrato de API do back-end,
- autenticacao/autorizacao,
- recomendacao personalizada de catalogo,
- analytics avancado de playback.

## Decisoes validadas

- Estrutura final de paginas:
  - `/upload`
  - `/catalog`
  - `/videos/:id`
- Navegacao: sidebar fixo no desktop e drawer colapsavel no mobile.
- Stack base preservada: Vite + React Router.

## Arquitetura de front-end

### App Shell

Criar um shell unico para todas as rotas com:

- area de navegacao lateral,
- area principal de conteudo,
- cabecalho contextual por pagina,
- comportamento responsivo controlado por breakpoint.

Comportamento:

- desktop: sidebar persistente com links principais,
- mobile: botao de menu abre drawer (`Sheet`) com os mesmos links,
- rota ativa destacada visualmente.

### Roteamento

- `UploadPage` (`/upload`): formulario de envio e lista resumida de uploads recentes.
- `CatalogPage` (`/catalog`): lista principal de videos com status e acao de abrir player.
- `VideoPlayerPage` (`/videos/:id`): foco na reproducao e estado detalhado do video.
- `Navigate` para redirecionar `/` para `/upload` com contrato explicito.

### Componentizacao

Separar componentes em blocos pequenos e reutilizaveis:

- navegacao: `AppSidebar`, `MobileNavDrawer`, `NavItem`,
- estados de tela: `PageHeader`, `ErrorState`, `EmptyState`, `LoadingState`,
- dominio de video: `VideoCard`, `VideoStatusBadge`, `UploadPanel`, `VideoGrid`,
- player: `ShakaVideoPlayer`, `PlayerStatusPanel`.

Cada componente deve ter responsabilidade unica e interface de props explicita.

## Design system e UX

### Direcao visual

Aplicar tema inspirado na Netflix sem copiar interface 1:1:

- base escura para foco em midia,
- vermelho de acento para CTAs e destaques,
- neutros com contraste alto para leitura,
- superficies com variacao de tons para profundidade.

### Paleta Netflix (obrigatoria)

Aplicar explicitamente a paleta com predominancia de tons Netflix:

- `--background: #141414` (fundo principal),
- `--surface: #181818` (cards e blocos),
- `--surface-2: #232323` (hover e superficies secundarias),
- `--primary: #E50914` (acao principal),
- `--primary-foreground: #FFFFFF`,
- `--text: #FFFFFF` (texto principal),
- `--text-muted: #B3B3B3` (texto auxiliar),
- `--border: #2F2F2F`,
- `--success: #46D369`,
- `--warning: #F5C518`,
- `--danger: #FF4D4F`,
- `--ring: #E50914`.

Regras de uso:

- usar `--primary` apenas para elementos de alta prioridade (CTA, estado ativo, foco),
- evitar saturar vermelho em elementos neutros de suporte,
- manter contraste AA minimo em textos e badges de status.

### Tokens

Definir variaveis de tema (CSS variables) para:

- `--background`, `--surface`, `--surface-2`,
- `--text`, `--text-muted`,
- `--primary`, `--primary-foreground`,
- `--success`, `--warning`, `--danger`,
- `--border`, `--ring`.

Os componentes shadcn devem consumir tokens para manter consistencia global.

### Componentes shadcn previstos

- `Button`, `Input`, `Card`, `Badge`, `Separator`,
- `Skeleton`, `Alert`,
- `Sheet` (drawer mobile),
- `ScrollArea` (quando necessario em listas longas).

## Camada de dados com TanStack Query

### Setup global

No bootstrap do app:

- criar `QueryClient` unico,
- envolver app com `QueryClientProvider`,
- manter React Router como camada de navegacao.

### Query keys e hooks

Padrao de chaves:

- `['videos']` para catalogo/lista,
- `['video', id]` para detalhe.

Hooks:

- `useVideosQuery()`
- `useVideoQuery(id)`
- `useUploadVideoMutation()`

Observacao de stack:

- manter o projeto em Vite + React Router,
- aplicar React Query sem migracao para TanStack Start.

### Polling e invalidação

- catalogo: `refetchInterval` mais conservador (ex.: 5s),
- upload: intervalo menor enquanto houver item em processamento,
- player: polling condicional quando status nao e `ready`.

No upload concluido com sucesso:

- invalidar `['videos']`,
- invalidar `['video', id]` quando o detalhe estiver em cache.

## Reproducao com Shaka Player

### Fluxo tecnico

No componente de player:

1. instalar polyfills (`shaka.polyfill.installAll()`),
2. validar suporte de browser,
3. criar instancia do player,
4. anexar ao elemento `HTMLVideoElement`,
5. registrar listener de erro,
6. carregar manifesto DASH,
7. destruir player no cleanup.

### Regras de UX no player

- mostrar estado de inicializacao/buffering,
- exibir mensagens de erro compreensiveis,
- esconder controles de playback quando video nao esta `ready`,
- fornecer CTA de retorno ao catalogo.

## Tratamento de erros

### Upload

- erro de rede: alerta claro com tentativa novamente,
- erro de validacao: mensagem orientada para acao (arquivo/titulo).

### Listagem e detalhe

- fallback visual para falha de carregamento,
- estado vazio amigavel quando nao houver videos.

### Playback

- erro de manifest/segmento: estado dedicado de falha no player,
- log de erro tecnico para suporte sem expor stack ao usuario final.

## Responsividade e acessibilidade

- navegacao completa por teclado,
- labels e hierarquia semantica consistentes,
- foco visivel em elementos interativos,
- contraste minimo AA para texto principal,
- drawer mobile com fechamento previsivel e titulo acessivel,
- componentes clicaveis com area de toque adequada.

## Estrategia de testes

### Unitarios/componentes

- renderizacao da sidebar e destaque de rota ativa,
- `UploadPanel` com estados disabled/loading,
- `VideoCard` e `VideoStatusBadge` por status,
- estados de erro e vazio por pagina,
- `ShakaVideoPlayer` com mocks de lifecycle (`attach`, `load`, `destroy`).

### Integracao de front-end

- fluxo upload -> invalidacao -> atualizacao de lista,
- polling condicional por status,
- navegacao `/upload` -> `/catalog` -> `/videos/:id`.

### E2E

- validar menu lateral + drawer mobile,
- upload de video e transicao de status,
- abertura do player para item pronto,
- exibicao de fallback quando item nao pronto.

## Plano incremental recomendado

1. preparar base visual e infra (shadcn + tokens + shell + rotas),
2. migrar camada de dados para React Query,
3. implementar paginas finais (`/upload`, `/catalog`, `/videos/:id`),
4. substituir player por Shaka e ajustar UX de playback,
5. completar testes e validacao cross-device.

## Uso obrigatorio de skills na implementacao

Durante a execucao do plano, adotar explicitamente as seguintes skills para garantir consistencia tecnica e aderencia ao `STYLE.MD`:

- `frontend-design`: orientar composicao visual, layout e acabamento da experiencia,
- `shadcn`: setup e composicao correta dos componentes shadcn/ui,
- `context7-mcp`: consulta de documentacao atualizada das bibliotecas antes de implementar APIs e padroes,
- `vercel-react-best-practices`: boas praticas de performance, renderizacao e composicao React,
- `tanstack-start-best-practices`: aplicar somente praticas reutilizaveis para o ecossistema TanStack que sejam compativeis com Vite + React Router, sem migrar arquitetura.

Checklist minimo por etapa:

1. consultar docs atualizadas via Context7 para a biblioteca da etapa,
2. aplicar padroes de UX/visual com `frontend-design`,
3. implementar componentes UI com `shadcn`,
4. revisar componentes React com `vercel-react-best-practices`,
5. validar padroes de dados/polling/invalidação com referencias do ecossistema TanStack.

## Riscos e mitigacoes

- setup inicial do shadcn/tokens pode impactar estilos existentes: migrar por etapas e validar visual a cada tela.
- troca de player pode introduzir regressao de reproducao: cobrir lifecycle com testes e fallback de erro.
- polling agressivo pode aumentar carga da API: ajustar intervalos por pagina e condicao de status.

## Resultado esperado

Ao final desta melhoria, o front-end deve:

1. apresentar identidade visual moderna consistente,
2. disponibilizar as 3 paginas separadas com navegacao lateral,
3. usar React Query como camada padrao de dados assincros,
4. reproduzir MPEG-DASH via Shaka Player com melhor UX,
5. manter comportamento responsivo e acessivel em desktop e mobile.
