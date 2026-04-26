# Contexto

Essa aplicação se trata de um simples protótipo de plataforma de streaming, a ideia fundamental aqui é construir uma simples que possua duas grandes funcionalidades

# Funcionalidades

## Transcoding

- Receber vídeos MP4 e realizar o processo de transcoding de forma assíncrona para o formato MPEG-DASH.
- O back-end deve ser responsável por fazer esse processamento.
- O front-end deve fornecer uma interface amigável para fazer o upload desses vídeos e acompanhar os vídeos que estão sendo carregados e qual o seu status

## Streaming

- Servir o video em streaming para o front-end.
- O player de vídeo pode escolher entre as qualidades 480p, 720p e 1080p

## Listagem de vídeos

- Deve haver uma listagem simples de vídeos onde basta clicar no vídeo e ele começará a ser transmitido em outra página

# Fora do Escopo

- Lidar com login, cadastro ou qualquer tipo de autenticação
- Criação de legenda para os vídeos
- Qualquer implementação de cache

# Tecnologias a serem utilizadas

## Back-end

- Golang
- PostgreSQL
- Echo

## Front-end

- React.js (Vite)
- TypeScript
- shadcn
