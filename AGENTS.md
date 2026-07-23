# AI Engineering System

## Base Rules (Always Follow)
- **Language**: Communicate in Russian.
- **Scope**: Use 80/20 principle. Analyze requests critically. Do not perform unnecessary changes.
- **Safety**: Preserve current behavior during refactoring. Do not make commits or destructive git actions without explicit request.
- **Cross-Layer Access**: Если задача относится к одному слою (например, UI/UX), но требует изменения файлов другого слоя (например, БД или API), ОБЯЗАТЕЛЬНО спроси разрешения у пользователя перед поиском и открытием этих файлов.

## Context Rules (READ ONLY WHEN RELEVANT)
When starting a task, ONLY read the rules that directly apply to the files or stack you are modifying:

- `RULES/core.md` -> General coding standards (read when creating new code)
- `RULES/architecture.md` -> New feature design, patterns
- `RULES/frontend.md` -> UI/UX, React/Next, Vue/Nuxt, CSS, components
- `RULES/backend.md` -> API, server logic, routing
- `RULES/database.md` -> SQL, ORM, schemas, migrations
- `RULES/security.md` -> Auth, data validation, protection
- `RULES/testing.md` -> Unit and E2E tests
- `RULES/performance.md` -> Optimization, caching, memoization
- `RULES/devops.md` -> Docker, CI/CD, deployment