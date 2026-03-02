# SDD + MCP Analysis (CLI/OpenSpec)

## Estado actual del CLI/OpenSpec
- `openspec init/list/show/status/complete/archive` vive en `cmd/agentic-agent/openspec.go`; genera `proposal.md` y `tasks.md` con plantillas mínimas y registra el cambio en `.agentic/openspec/changes/`.
- `internal/openspec/manager.go` gestiona el ciclo draft → imported → implementing → implemented → archived; importa tareas de `tasks.md` y escribe `metadata.yaml`, pero no valida MCP, gates ni traceabilidad.
- `internal/openspec/parser.go` parsea listas y archivos `tasks/*.md` para descripción/AC/skills; no extrae `implements`, `context_pack`, contratos ni ADRs.
- Plantillas (`internal/openspec/templates/*.tmpl`) solo contienen secciones genéricas; no incluyen Tech Stack, Spec Graph ni campos de gobierno.
- `autopilot` (`internal/orchestrator/autopilot.go`) auto-importa tareas pero no verifica gates ni contexto antes de reclamar/ejecutar.

## Brecha vs Operando Modelo SDD+MCP (Enterprise)
- Falta relación con Platform Spec (SpecKit): no se declara `implements` ni se consumen fan-out tasks.
- No hay MCP Context Pack versionado ni Gate Check (Context, Domain, Integration, NFR, Ready).
- No se rastrean contratos, ADRs bloqueantes ni se actualiza un Spec Graph (`implements/dependsOn/affects/status`).
- No existe split de rutas para Feature/Change/Bug/Hotfix ni enforcement de políticas/roles.
- Todo ocurre en un solo repo de componente; no hay controles de gobernanza ni de rebase de contexto.

## Propuesta de incorporación (resumen)
1. **Metadatos obligatorios**: extender `TemplateData`, `metadata.yaml` y plantillas para incluir `implements`, `context_pack`, `contracts_referenced`, `blocked_by`, `status`, `outcome`, enlaces de Spec Graph; propagar a tareas.
2. **Gate checks en CLI**: comando `openspec check <id>` y bloqueo en `complete/autopilot` si falla algún gate (G1-G5).
3. **Ingesta de fan-out (SpecKit → OpenSpec)**: flag `--from-plan` que lea handoff tasks con `component_repo`, `platform_spec_id`, `context_pack_version`, `contract_change`, `blocked_by` y precargue metadatos.
4. **Enforcement de Context Pack MCP**: `--context-pack <id>` en `init`; comando para fetch/validar packs en `.agentic/context/packs/`; gate1 depende de esto.
5. **Spec Graph**: tras `openspec complete`, actualizar `.agentic/spec-graph.json` con nodos/aristas (`implements/dependsOn/affects/status`).
6. **Autopilot con gates**: antes de reclamar tarea, validar gates del cambio asociado; abortar con mensajes de faltantes.
7. **Flujos por tipo de cambio**: `--type {feature,change,bug,hotfix}` elige plantilla específica (incl. hotfix minimal) alineada al OPERATING-MODEL.
8. **Lint/CI**: `openspec lint` verifica metadatos, gates, context pack y links de Spec Graph para PR checks.
9. **Compatibilidad progresiva**: permitir flujo actual pero marcar estado `blocked` y avisar si faltan metadatos/gates en cambios heredados.

## Nuevos comandos de plataforma (SpecKit layer)
- `agentic-agent platform init`
  - Bootstraps repo de plataforma: `constitution/`, `initiatives/`, `platform-specs/`, `contracts/`, `adr/`, `spec-graph.json`, `.agentic/context/packs/`.
  - Configura `agnostic-agent.yaml` con rutas de plataforma y activa gates por defecto.

- `agentic-agent platform add-feature --initiative ECO-123 --context-pack cp-v3 --name "Cart Checkout"`
  - Crea un Platform Spec usando plantillas MCP-aware (con `implements` vacío, `context_pack` cp-v3, `status=Draft`).
  - Opcional: `--fanout cart,checkout,payments` genera archivo de handoff para OpenSpec (`fanout.yaml`).

- `agentic-agent platform change-feature --id SPEC-PLAT-123 --field status --value Approved`
  - Actualiza metadatos de un Platform Spec (status, context_pack, blocked_by ADRs, contract_change flag).
  - Revalida gates con `platform gate-check` y marca fan-outs afectados como "stale" si cambia `context_pack`.

- `agentic-agent platform change-priority --initiative ECO-123 --priority {Planned|Discovery|Draft|Approved|Implementing|Done|Paused|Blocked}`
  - Ajusta prioridad/estado de iniciativa y propaga etiqueta a specs hijos y fan-outs (señal para equipos de componente).

Integración: los comandos `platform *` producen/actualizan el fan-out consumido por `openspec init --from-plan` y sincronizan `spec-graph.json` en la capa de plataforma.
