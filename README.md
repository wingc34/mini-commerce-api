# mini-commerce-api

A production-minded RESTful API server for an e-commerce platform, built with Go. This is the backend counterpart to [mini-commerce](https://github.com/wingc34/mini-commerce).

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Authentication | Google OAuth2 + JWT |
| Payment | Stripe |
| Containerisation | Docker + docker-compose |

## Architecture

This project follows **Clean Architecture** with a three-layer separation:

```
HTTP Request
    ↓
Handler Layer    → Parse request, call service, return JSON response
    ↓
Service Layer    → Business logic (e.g. address default promotion, draft order creation)
    ↓
Repository Layer → Database operations only, no business logic
    ↓
PostgreSQL
```

Each layer depends on **interfaces**, not concrete implementations. This means:
- The service layer has no knowledge of GORM or SQL
- The handler layer has no knowledge of the database
- Any layer can be swapped or mocked independently for testing

### Project Structure

```
mini-commerce-api/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Environment config
│   ├── database/        # GORM connection
│   ├── handler/         # Gin HTTP handlers
│   ├── middleware/       # JWT auth middleware
│   ├── model/           # GORM model structs
│   ├── repository/      # Database queries
│   ├── router/          # Route registration
│   └── service/         # Business logic
└── pkg/
    ├── ctxutil/         # Gin context helpers
    ├── id/              # cuid2 ID generation
    ├── jwt/             # JWT sign / verify
    ├── oauth/           # Google OAuth2 config
    ├── response/        # Unified JSON response helpers
    └── stripe/          # Stripe client
```

## API Endpoints

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/auth/google` | No | Redirect to Google OAuth consent screen |
| `GET` | `/api/v1/auth/google/callback` | No | Exchange OAuth code, return JWT |

### Products
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/products` | No | Paginated product list (page size: 9) |
| `GET` | `/api/v1/products/recommended` | No | 4 random recommended products |
| `GET` | `/api/v1/products/:id` | No | Product detail with SKUs |
| `POST` | `/api/v1/products/:id/stock` | No | Check SKU stock by attributes |

### Users
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/v1/users/me` | ✅ | Get current user profile |
| `PATCH` | `/api/v1/users/me` | ✅ | Update name and phone number |
| `GET` | `/api/v1/users/me/addresses` | ✅ | List non-deleted addresses |
| `POST` | `/api/v1/users/me/addresses` | ✅ | Create address (auto-default if first) |
| `PUT` | `/api/v1/users/me/addresses/:id` | ✅ | Update address |
| `DELETE` | `/api/v1/users/me/addresses/:id` | ✅ | Soft delete, promote new default if needed |
| `POST` | `/api/v1/users/me/wishlist` | ✅ | Add product to wishlist |
| `DELETE` | `/api/v1/users/me/wishlist/:productId` | ✅ | Remove product from wishlist |

### Orders
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/orders/draft` | ✅ | Create draft order atomically with items |
| `GET` | `/api/v1/orders` | ✅ | Paginated order list (page size: 4) |
| `GET` | `/api/v1/orders/:id` | ✅ | Order detail (accepts order ID or draft order ID) |

### Payments
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/payments/intent` | ✅ | Create Stripe PaymentIntent, return clientSecret |
| `POST` | `/api/v1/webhooks/stripe` | Stripe sig | Handle payment_intent.succeeded, create Order from DraftOrder |

## Key Design Decisions

**Draft Order Pattern** — Orders are created in two steps to handle the async nature of Stripe payments. A `DraftOrder` is created first and persists until the Stripe webhook confirms payment success, at which point a real `Order` is created atomically.

**Soft Delete on Addresses** — Addresses are soft-deleted (`deleted_at`) rather than hard-deleted to preserve order history integrity. When the default address is deleted, another address is automatically promoted.

**JWT over Session** — Stateless JWT authentication with minimal payload (`user_id`, `email`, `exp`) keeps the server stateless and scales horizontally without shared session storage.

**Interface-driven Repository** — Every repository is defined as an interface, making it straightforward to inject mock implementations in tests without a real database.

## Getting Started

### Prerequisites

- Go 1.22+
- Docker and docker-compose
- A Google Cloud project with OAuth2 credentials
- A Stripe account

### Installation

```bash
git clone https://github.com/wingc34/mini-commerce-api
cd mini-commerce-api
go mod download
```

### Environment Variables

```bash
cp .env.example .env
```

| Variable | Description |
|----------|-------------|
| `PORT` | Server port (default: 8080) |
| `ENV` | Environment (`development` / `production`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `JWT_EXPIRY_HOURS` | JWT expiry duration in hours |
| `GOOGLE_CLIENT_ID` | Google OAuth2 client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth2 client secret |
| `GOOGLE_REDIRECT_URL` | OAuth2 callback URL |
| `STRIPE_SECRET_KEY` | Stripe secret key |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret |

### Running Locally

```bash
# Start PostgreSQL
docker-compose up -d

# Run database migrations
make migrate

# Start the server
make run
```

The server will be available at `http://localhost:8080`.

### Verify Setup

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Database Schema

The schema includes 9 tables across three domains:

**Users** — `users`, `addresses`, `wishlists`

**Products** — `products`, `skus`

**Orders** — `draft_orders`, `draft_order_items`, `orders`, `order_items`

See [`migrations/001_init.sql`](migrations/001_init.sql) for the full schema.

## Known Limitations

- No inventory locking during checkout — two users could theoretically purchase the last item simultaneously
- Expired draft orders are not automatically cleaned up (no background job)
- Google OAuth `state` parameter is a fixed string; production should use a random value per request to prevent CSRF

## Related

- **Frontend**: [mini-commerce](https://github.com/wingc34/mini-commerce) — Next.js + tRPC + Prisma