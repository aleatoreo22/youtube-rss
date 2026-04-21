# YouTube RSS Aggregator

Uma aplicação que agrega conteúdo de canais do YouTube via RSS feeds, permitindo consumo offline e organização por usuário.

## 📋 Visão Geral

Este projeto consiste em duas partes:

1. **API Server (Go)** - Um servidor REST que:
   - Busca feeds RSS do YouTube
   - Armazena conteúdo em SQLite
   - Serve conteúdo via endpoints HTTP

2. **App Mobile (Flutter)** - Aplicativo que:
   - Consome a API para exibir vídeos
   - Permite adicionar/remover canais
   - Organiza conteúdo por data

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│                     youtube-rss-api (Go)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Service     │  │  Repository  │  │    Handler   │  │
│  │   Layer      │  │    Layer     │  │     Layer    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│  ┌──────────────┐                                      │
│  │   Database   │  (SQLite)                             │
│  └──────────────┘                                      │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                   youtube_rss_ui (Flutter)               │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Main App                             │   │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐    │   │
│  │  │ Videos │ │History │ │  User  │ │  Nav   │    │   │
│  │  └────────┘ └────────┘ └────────┘ └────────┘    │   │
│  └──────────────────────────────────────────────────┘   │
│                          ↓                                 │
│              API Client (YoutubeRSS service)              │
└─────────────────────────────────────────────────────────┘
```

## 🚀 Funcionalidades

- ✅ Agregação automática de conteúdo do YouTube
- ✅ Sincronização em background (a cada 5 minutos)
- ✅ Organização por usuários e canais
- ✅ Paginação de conteúdo
- ✅ Histórico de visualizações
- ✅ Tema escuro com Flutter
- ✅ Suporte a mobile e desktop

## 📦 Instalação

### Pré-requisitos

- [Go](https://golang.org/) (1.21+)
- [Flutter](https://flutter.dev/) (3.10+)
- [SQLite](https://www.sqlite.org/)

### API Server (Go)

```bash
cd youtube-rss-api
go mod download
go build -o api ./cmd/api
./api
```

A API estará rodando em `http://localhost:1234`

### Flutter App

```bash
cd youtube_rss_ui
flutter pub get
flutter run
```

## 📡 Endpoints da API

| Método | Path | Descrição |
|--------|------|-----------|
| GET | /ping | Health check |
| GET | /user/:id/content/:date | Conteúdo por data |
| GET | /user/:id/content?page=:p&limit=:l | Paginação |
| POST | /user/:id/channel | Adicionar canal |
| DELETE | /user/:id/channel | Remover canal |
| GET | /user/:id/channel | Listar canais do usuário |

## 🗄️ Banco de Dados

O projeto usa SQLite com as seguintes tabelas:

- `users` - Informações dos usuários
- `channels` - Canais do YouTube
- `user_channels` - Relação usuário-canais
- `content` - Conteúdo dos vídeos

## 📱 Uso do App

1. Adicione canais do YouTube
2. O app sincroniza automaticamente a cada 5 minutos
3. Navegue pelos vídeos em ordem cronológica
4. Use a barra de navegação inferior para alternar entre:
   - **Videos** - Feed principal
   - **History** - Histórico
   - **User** - Configurações do usuário

## 🛠️ Desenvolvimento

### Testes

```bash
# API
cd youtube-rss-api
go test ./...

# Flutter App
cd youtube_rss_ui
flutter test
```

### Linting

```bash
cd youtube_rss_ui
flutter analyze
```

## 📄 Licença

[MIT](LICENSE)

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, abra um issue ou PR.
