# MobileDev — Tools & CLI Reference

## Task & Context Management

- `agentic-agent task list` — view backlog and in-progress tasks
- `agentic-agent task claim <ID>` — claim task (records branch + timestamp)
- `agentic-agent context build --task <ID>` — load context bundle (includes openspec + api_contracts + tech-stack)
- `agentic-agent validate` — run all quality validators before completing
- `agentic-agent task complete <ID>` — mark done (auto-captures commits)
- `agentic-agent openspec check <spec-id>` — verify spec readiness before claiming

## Flutter & Dart CLI

**Project Management:**
- `flutter create myapp` — create new Flutter project
- `flutter pub get` — fetch dependencies (same as `dart pub get`)
- `flutter pub upgrade` — upgrade dependencies
- `flutter clean` — clean build artifacts

**Development & Testing:**
- `flutter run` — run app on connected device/emulator
- `flutter run -d chrome` — run on web
- `flutter analyze` — check for errors/warnings (must be zero warnings)
- `flutter test` — run all unit tests
- `flutter test test/auth_test.dart` — run specific test file
- `flutter test --coverage` — generate coverage report
- `dart format lib/` — auto-format Dart code
- `dart fix --apply` — auto-fix analyzer issues

**Build & Release:**
- `flutter build ios` — build iOS app
- `flutter build android` — build Android app
- `flutter build web` — build web app
- `flutter build apk` — build Android APK
- `flutter build aab` — build Android App Bundle

## Dart/Flutter MCP Server Tools

**Setup:**
```bash
# Install and configure Dart MCP server (requires Dart 3.9+)
claude mcp add --transport stdio dart -- dart mcp-server
```

**Key Commands via MCP:**

Package Discovery:
- `pub_dev_search "package-name"` — search pub.dev for packages
- Returns: package name, description, stars, latest version, maintainers
- Use before adding any dependency to find most-maintained option

Widget Tree Inspection:
- `widget_tree_inspect` — get structure of rendered widget tree
- `widget_tree_validate <selector>` — check if widget matches expected layout
- Helps verify visual AC compliance without screenshots

Runtime Error Detection:
- `runtime_error_check` — catch runtime errors from running app
- `analyzer_check` — run Dart analyzer and return error list

Code & Formatting:
- `dart_format_check` — verify code is formatted
- `dart_format_apply` — auto-format all code

Test Execution:
- `test_run` — run all tests and return results
- `test_run "test/auth_test.dart"` — run specific test
- `coverage_report` — generate and parse coverage report

Dependency Management:
- `pubspec_check` — validate pubspec.yaml
- `dependency_update` — list outdated dependencies
- `dependency_add "package_name" "^1.0.0"` — add dependency to pubspec.yaml

## HTTP Client Library (for API calls)

**Dio (recommended):**
```dart
import 'package:dio/dio.dart';

final dio = Dio();
final response = await dio.post(
  '/api/v1/auth/login',
  data: {'email': 'user@example.com', 'password': 'secret'},
);
```

**Built-in http package:**
```dart
import 'package:http/http.dart' as http;

final response = await http.post(
  Uri.parse('http://localhost:8080/api/v1/auth/login'),
  body: jsonEncode({'email': 'user@example.com', 'password': 'secret'}),
);
```

Always read `api_contracts[].path` from context before writing HTTP code — never assume endpoint structure.

## State Management

**Riverpod (recommended):**
```dart
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});
```

**BLoC:**
```dart
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  // ...
}
```

**Provider (legacy):**
```dart
final authProvider = ChangeNotifierProvider((ref) => AuthProvider());
```

Check `tech-stack.md` for project's chosen pattern.

## Key Paths

- `.agentic/spec/` — all openspec proposals with acceptance criteria and platform scope
- `.agentic/contracts/` — API contracts (read before writing HTTP calls)
- `.agentic/context/` — global-context.md (design tokens, Material 3 theme), tech-stack.md (Flutter version, state management)
- `.agentic/coordination/` — announcements.yaml, reservations.yaml
- `lib/` — Dart source code
- `lib/src/` — internal implementation
- `test/` — unit tests
- `pubspec.yaml` — project dependencies and metadata
- `analysis_options.yaml` — Dart analyzer configuration

## API Contract Reference

Before writing any HTTP call with Dio or http:

1. Read the contract from `api_contracts[].path` in your context bundle
2. Contract is an OpenAPI spec (YAML or JSON) with:
   - `paths:` — endpoints (GET /api/v1/users, POST /api/v1/auth/login, etc.)
   - `components.schemas:` — data models (User, AuthToken, Error, etc.)
   - `securitySchemes:` — authentication (Bearer token, API key, OAuth, etc.)
3. Never assume additional fields or endpoints not in the contract
4. If contract is incomplete: report `contract-deviation` to TechLead

## Announcements Format

When announcing task completion to TechLead:

```yaml
- from_agent: mobile-dev
  to_agent: tech-lead
  project_id: proj-001
  status: complete
  summary: "Auth feature implemented on iOS and Android"
  data:
    task_id: TASK-044
    branch: feature/auth-mobile-v1
    platforms_tested: [ios, android, web]
    widgets_created:
      - lib/src/screens/auth/login_screen.dart
      - lib/src/screens/auth/register_screen.dart
      - lib/src/widgets/auth_form.dart
    test_results:
      unit_tests: 22 passed
      coverage: 89%
    ac_coverage:
      - "Login screen displays email/password fields (Material 3)" ✅
      - "Submit button disabled until both fields filled" ✅
      - "Login calls POST /auth/login and stores JWT token" ✅
    analyzer_warnings: 0
    commits: [abc123, def456, ghi789]

# Or if you find an API mismatch:
- from_agent: mobile-dev
  to_agent: tech-lead
  project_id: proj-001
  status: contract-deviation
  summary: "POST /auth/login returns different error structure than spec"
  data:
    task_id: TASK-044
    spec_ref: .agentic/contracts/auth-api.yaml
    endpoint: POST /api/v1/auth/login
    expected: '{"error": "Invalid credentials", "code": "AUTH_INVALID"}'
    actual: '{"message": "Login failed"}'
    reproduction: "Try login with invalid email"
    severity: blocking
```

## Git Workflow

- Create feature branch: `git checkout -b feature/task-id-short-name` (TechLead handles main)
- Commit messages: `feat: <screen/widget>: <AC description>` or `fix: <widget>: <bug description>`
- Push to feature branch: `git push origin feature/<branch>`
- Never force-push or merge directly to main
