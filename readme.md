# zbxedx

`zbxedx` is a small framework for writing **compiled Zabbix Agent2 plugins in Go**.

It aims to reduce boilerplate, make plugin behavior explicit, and simplify the development of plugins that require
custom metric collection or shared data snapshots.

The framework is intentionally minimal and focuses on **clarity and explicit lifecycle handling**, rather than
abstraction or automation.

---

## Motivation

Writing plugins for Zabbix Agent2 is possible, but in practice often more complex than expected. Several recurring
issues appear even in relatively small plugins.

### A. Non-obvious plugin lifecycle

Understanding how an Agent2 plugin actually works — especially built-in plugins — requires familiarity with Zabbix
internals such as watch-based metrics, item lifecycle, and configuration handling.

As a result, even simple plugins can be hard to reason about without reading existing source code or reverse-engineering
agent behavior.

---

### B. Excessive boilerplate per plugin

Each plugin typically reimplements the same patterns:

- configuration parsing
- metric registration
- item key handling
- lifecycle wiring

This boilerplate often obscures the actual purpose of the plugin and makes small plugins harder to write, review, and
maintain.

`zbxedx` aims to reduce this boilerplate to the minimum required for correctness, without hiding behavior behind
abstractions.

---

### C. Metrics derived from a shared data snapshot

Some plugins need to collect a single data snapshot and expose **multiple items** derived from it.

In a vanilla Zabbix setup, this usually leads to one of the following compromises:

- server- or proxy-side post-processing  
  (which may be inefficient or unnecessary for simple transformations)

- trapper items combined with `zabbix_sender`  
  (which are not supported directly by Agent2 and require additional tooling)

Neither option is always desirable, especially for frequently updated or locally available data.

`zbxedx` allows plugins to collect data once and publish multiple item values from that snapshot using watch-based
metrics, without additional agents, senders, or server-side post-processing.

---

### D. Code separation, reuse, and long-term maintenance

Traditionally, Agent2 plugins live inside the Zabbix codebase namespace and build environment. This tightly couples
plugin logic to Zabbix internals and adds friction to development and reuse.

In practice, this often leads to:

- plugins that cannot be built or tested outside the Zabbix source tree
- duplicated logic across plugins
- limited reuse of internal packages
- additional effort required for iteration and maintenance

`zbxedx` encourages a clearer separation between:

- Zabbix integration glue
- plugin-specific logic
- reusable internal packages

This makes plugins easier to develop independently, easier to reuse, and easier to maintain over time.

---

## Core Concepts

### Plugin

A plugin is registered once and defines:

- its name
- configuration options
- one or more metric sets

Registration typically happens in `init()`.

---

### MetricSet

A `MetricSet` describes a logical group of metrics.

It may include:

- a list of metrics (item keys)
- an optional asynchronous collector

Metric sets can operate independently from the Zabbix agent’s default collection cycle.

---

### Metrics

Each metric defines:

- an item key
- a description
- how its value is produced

Two models are supported:

1. **Exported metrics**  
   Collected by the Zabbix agent on demand via an `Export` function.

2. **Asynchronously collected metrics**  
   Produced by a custom collector running in a dedicated goroutine and
   published independently from the agent’s reporting cycle.

---

### Watch Filters

For asynchronously collected metrics, `zbxedx` uses **watch filters**.

A watch filter:

- is created by the agent when it needs to handle a specific item key
- receives published event payloads
- extracts a value for a concrete Zabbix item

Filters are created per item key (and parameters) and are not global.  
This mirrors how Zabbix handles item lifecycle and avoids shared mutable state.

---

## Bootstrap

The `bootstrap.sh` script downloads and prepares the Zabbix source tree used for local builds and integration testing.

## Quick start (Docker)

The easiest way to verify that the framework compiles inside Zabbix Agent2 and that
plugins are loaded at startup is to build and run the provided Docker image.

From the repository root:

```bash
docker build -t zbxedx-agent .
docker run --rm zbxedx-agent
```

The container uses example configuration files from conf/ and loads the dice
example plugin.

Note: the Docker image is intended for local testing and development. It is not a
production deployment recipe.

---

## Repository layout

- `edxagent/` — framework core (reusable plugin infrastructure)
- `plugins/` — example plugins (e.g., `plugins/dice`)
- `conf/` — example Zabbix Agent2 configuration files for local testing
- `zabbix_templates/` — exported Zabbix templates/host definitions for quick server-side setup
- `zabbix/` — Zabbix source tree used for local builds and integration testing. Created automatically during bootstrap
- `bootstrap.sh` — helper script to fetch/bootstrap the Zabbix source tree
- `Dockerfile` — local build image for smoke-testing plugin integration

---

## Example Plugin

The repository includes a fully documented example plugin: **`dice`**.

It demonstrates:

- plugin registration
- configuration parsing
- asynchronous metric collection
- watch filter usage
- coexistence of async and standard agent metrics

The example is intended to serve as **executable documentation**.

---

## Recommended Plugin Structure

A typical plugin can be split into a small set of files with clear responsibilities:

- `plugin.go` — plugin registration and configuration
- `metrics.go` — metric sets and collection logic
- `filters.go` — watch filters

This structure is recommended for clarity, but not enforced.

---

## Design Principles

- Explicit over implicit
- Minimal API surface
- No hidden goroutines
- No global registries beyond plugin scope
- Example code is part of the documentation

If something is not obvious from the code, it likely needs a comment.

---

## Status

`zbxedx` is intentionally small and opinionated.

It is not a general-purpose SDK and does not try to abstract Zabbix itself. The framework evolves together with
real plugins and real use cases.

---

## License

MIT License.