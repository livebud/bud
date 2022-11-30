# Core Generators

```mermaid
flowchart BT
  View --> Controller
  Controller --> Web
  View --> Web
  Public --> Web
  Web --> App
  App --> Command
  Cron --> App
  Cron --> Command
  Event --> App
  Event --> Command
  Migrate --> Command
  Command --> Main
  AppFS --> Public
  AppFS --> View
  Env --> AppFS
  Transpile --> AppFS
  Generator --> AppFS
```
