# Developer Preferences

This is a project-level configuration. Please strictly adhere to the following general principles regarding communication, code style, workflow, and prohibited actions.

## Language & Communication Style
- Respond entirely in **English**.
- Target vocabulary level: **3,000 words** (intermediate/CEFR B1 level).
- Whenever you use advanced, low-frequency, or complex words that exceed the 3,000-word vocabulary limit, attach the Chinese translation in parentheses right after the word, like this: `word (中文解释)`.
- Keep sentences clear, direct, and grammatically accurate to help me learn English while reading code/explanations.

## Coding Philosophy
- **Defensive Programming**: Always consider edge cases, null/undefined checks, and exception handling
- **KISS Principle**: Keep It Simple, Stupid. Avoid over-engineering; prioritize readability and maintainability over extreme performance optimization.

## Code Style & Linting
**Naming Conventions**:
- Pinyin naming is strictly forbidden! All variables, functions, classes, and filenames must use accurate English words.
- Frontend/Node.js: Use `camelCase` for variables and functions, `PascalCase` for classes and components, and `UPPER_CASE` for constants.
- Database/SQL: Use `snake_case` consistently for fields and table names.

**Go (Google Go Style)**:
- **Code Formatting**: Strictly follow `gofmt` (using **Tab** for indentation) and `goimports`.
- **Naming**: Use `mixedCaps` (i.e., `camelCase` or `PascalCase`). Acronyms must stay fully uppercase (e.g., `URL`, `API`, instead of `Url`, `Api`).
- **Error Handling**: Never ignore errors; errors must be returned as the final return value of a function. Prefer guard clauses using `if err != nil { return err }` to exit early.
- **Comments**: All exported names (starting with a capital letter) must have complete documentation comments.

## 5. Git Workflow
**Commit Message**: Must follow Conventional Commits (Angular style) in the format `<type>(scope): <description>`. For example:
- `feat: add user authentication layout`
- `fix(api): resolve token expiration memory leak`
**Atomic Commits**: Each commit or modification should address only a single issue. Never mix feature development with unrelated code formatting in a single commit.
**No Auto-Commit/Push**: Do not automatically perform `git commit` or `git push` unless explicitly requested.

## Red Lines & Dangerous Operations
The following operations require explicit user confirmation before execution, even in auto-accept mode:
- Deleting files, directories, or Git history.
- Modifying `.env` files, secrets, tokens, certificates, or CI/CD configurations.
- Executing `git push`, `git rebase`, `git reset --hard`, or force pushing.


> **WARNING**: If a project-level `CLAUDE.md` exists in the root directory of the current project and its rules conflict with this global file, **the project-level rules take precedence**.