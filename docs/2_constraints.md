# Constraints

## Database Implementation Constraints

The solution adopts SQLite for the POC phase with the following technical constraints:

- **GORM ORM** driver implementation requiring:
  - CGO-enabled builds due to SQLite dependencies

## CI/CD Pipeline Constraints

Current integration testing approach is constrained by:

- **Temporary pipeline simulation** using:
  - Shell script-based test execution (`*.sh` files)
  - Docker Compose for environment orchestration
- **Deferred automation** of:
  - Build pipeline implementation
  - QA stage tooling
