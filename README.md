
  <h1 align="center">📈 Trading Application</h1>

<p align="center">
  <strong>Modern Full-Stack Stock Trading & Portfolio Management Platform</strong>
</p>

<p align="center">
  <a href="#tech-stack"><img src="https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=for-the-badge&logo=vuedotjs&logoColor=white" alt="Vue 3" /></a>
  <a href="#tech-stack"><img src="https://img.shields.io/badge/Vuetify-3.0-1867C0?style=for-the-badge&logo=vuetify&logoColor=white" alt="Vuetify" /></a>
  <a href="#tech-stack"><img src="https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" /></a>
  <a href="#tech-stack"><img src="https://img.shields.io/badge/PostgreSQL-14-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" /></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License" /></a>
</p>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Project Status](#-project-status)
- [Features](#-features)
- [Security](#-security)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [Project Structure](#-project-structure)
- [API Reference](#-api-reference)
- [Contributing](#-contributing)
- [Disclaimer](#-disclaimer)
- [License](#-license)

---

## 📌 Overview

Trading Application is a full-stack stock trading platform developed using **Go (Golang)**, **Vue.js**, **Vuetify**, and **PostgreSQL**.

The application enables users to securely register, manage their profile, maintain wallets, place buy and sell orders, track portfolios, and manage investments through a modular and scalable architecture.

This project is being developed as a production-oriented learning and portfolio application following enterprise software engineering practices.

---

## 🚧 Project Status

This project is currently under active development.

### Core Modules Completed
- User Authentication
- Wallet Management
- Stock Trading
- Profile Management
- Admin Dashboard

### In Progress
- Live Stock Market API Integration
- Portfolio Analytics

### Planned
- Watchlist
- Notifications
- Real-Time Updates
- Mobile Optimization

---

## ✨ Features

### 🔐 Authentication & Security
- Secure user registration and login
- JWT-based authentication
- Comprehensive profile management

### 💰 Wallet & Orders
- **Wallet Management** — Maintain balances and transaction history
- **Order Placement** — Execute secure buy and sell orders
- **Real-Time Tracking** — Monitor stock prices and portfolio performance

### 📊 Portfolio Management
- **Dashboard** — Live overview of holdings and total balance
- **Investment Management** — Track and manage long-term investments

### 🚀 Enterprise Architecture
- **Modular & Scalable** — Designed following enterprise software engineering practices
- **High-Performance Backend** — Optimized concurrent processing with Go

---

## 🔒 Security

- JWT Authentication
- Password Hashing (bcrypt)
- Protected Routes
- Role-Based Authorization
- Session Management
- Backend Input Validation
- Frontend Input Validation
- Environment-Based Configuration

---

## 🛠️ Tech Stack

### Frontend
| Technology | Purpose |
|-----------|---------|
| **Vue 3** | Progressive JavaScript framework for building user interfaces |
| **Vuetify** | Material Design component framework for Vue |
| **Node.js** | Runtime environment for development |

### Backend
| Technology | Purpose |
|-----------|---------|
| **Go** | Core language for backend services |
| **Go Modules** | Dependency management |
| **RESTful API** | Communication between client and server |

### Infrastructure & Databases
| Technology | Purpose |
|-----------|---------|
| **PostgreSQL** | Primary robust relational database for users, wallets, and trades |

---

## 🏗️ Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                        CLIENT (Browser)                     │
│  ┌────────────┐  ┌──────────────┐  ┌─────────────────────┐  │
│  │   Vue 3    │  │   Vuetify    │  │    Trading UI       │  │
│  │ Components │  │   Styling    │  │     & Wallets       │  │
│  └─────┬──────┘  └──────────────┘  └──────────┬──────────┘  │
│        │                                       │            │
│        │  REST API                             │ HTTP       │
└────────┼───────────────────────────────────────┼────────────┘
         │                                       │
         ▼                                       ▼
┌─────────────────────┐              ┌────────────────────────┐
│    Go Backend       │              │      External APIs     │
│  ┌───────────────┐  │              │    (Market Data)       │
│  │  Controllers  │  │              └────────────────────────┘
│  ├───────────────┤  │
│  │  Services     │  │
│  ├───────────────┤  │
│  │  Repositories │  │
│  ├───────────────┤  │
│  │  Auth / JWT   │  │
│  └───────┬───────┘  │
│          │          │
└──────────┼──────────┘
           │
           ▼
┌───────────────────────────┐
│       Databases           │
│  ┌─────────────────────┐  │
│  │     PostgreSQL      │  │
│  │ (Profiles & Trades) │  │
│  └─────────────────────┘  │
└───────────────────────────┘
```

---

## 🚀 Getting Started

### Prerequisites

| Requirement   | Version |
|---------------|---------|
| **Node.js**   | ≥ 18.x  |
| **Go**        | ≥ 1.20  |
| **PostgreSQL**| Latest  |

### 1. Clone the Repository

```bash
git clone https://github.com/Bharani6/Trading_Application.git
cd Trading_Application
```

### 2. Backend Setup

Navigate to the `backend` directory, install dependencies, and start the Go server:

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 3. Frontend Setup

In a new terminal window, navigate to the `frontend` directory, install dependencies, and start the development server:

```bash
cd frontend
npm install
npm run dev
```

---

## 📁 Project Structure

```text
Trading_Application/
├── backend/                          # Go Backend (Modular Architecture)
│   ├── cmd/server/                   # Application entry point (main.go)
│   ├── configs/                      # Environment & app configurations
│   ├── internal/                     # Core Business Logic & Modules
│   │   ├── admin/                    # Admin panel controllers & services
│   │   ├── auth/                     # Authentication & JWT logic
│   │   ├── database/                 # DB connection & ORM setup
│   │   ├── market/                   # External market data integration
│   │   ├── middleware/               # HTTP middlewares (auth, logging)
│   │   ├── profile/                  # User profile & KYC management
│   │   ├── routes/                   # API route definitions
│   │   ├── trade/                    # Trading execution logic
│   │   ├── user/                     # User management & models
│   │   ├── wallet/                   # Fund management & ledgers
│   │   └── utils/                    # Shared utilities
│   ├── migrations/                   # Database migration files
│   ├── go.mod                        # Go dependencies
│   └── go.sum
│
├── frontend/                         # Vue 3 Frontend (Component-Based)
│   ├── public/                       # Static assets
│   ├── src/                          # Vue source code
│   │   ├── api/                      # Axios clients & API wrappers
│   │   ├── components/               # Reusable Vue/Vuetify components
│   │   ├── layouts/                  # Page layouts (Navbar, Footer, Sidebar)
│   │   ├── router/                   # Vue Router definitions
│   │   ├── services/                 # Business logic & external calls
│   │   ├── store/                    # State management (Vuex/Pinia)
│   │   ├── views/                    # Main page components
│   │   └── App.vue                   # Root component
│   ├── package.json                  # NPM dependencies
│   └── vite.config.js                # Build configuration
│
├── LICENSE                           # MIT License
└── README.md
```

---

## 📡 API Reference

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Authenticate & get JWT token |

### Trading & Wallets
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/wallet/balance` | Get current user's wallet balance |
| `GET` | `/api/v1/shares` | Get available shares to trade |
| `POST` | `/api/v1/trades/buy` | Execute a buy order |
| `POST` | `/api/v1/trades/sell` | Execute a sell order |
| `GET` | `/api/v1/trades/history` | Get user's trade history |

*(Note: API endpoints are examples based on standard trading architecture)*

---

## 🤝 Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

---

## ⚠️ Disclaimer

This application is developed for educational and portfolio purposes only.

It does not provide real financial advice or execute real stock market transactions.

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](./LICENSE) file for details.

---

## 👨💻 Author

**Bharanidharan R** — [@Bharani6](https://github.com/Bharani6)

---

<p align="center">
  <sub>Built with ❤️ for modern trading.</sub>
</p>
