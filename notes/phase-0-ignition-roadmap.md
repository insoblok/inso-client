
# 🚀 Phase 0: `Ignition`

> _“We have fuel, we have wires. Now we spark the first signal.”_

## 🎯 Goals

- 🔥 Finalize component names and architecture (✅ done)
- 🧪 Continue with `TestClient` for contract interactions
- 🧱 Build WebSocket stream in `DevServer` for `ChainUI`
- 📦 Define JSON schema for event messages
- 🖥️ Setup `ChainUI` scaffold (Scala.js project)
- 🧪 Test socket connection end-to-end

---

## 📦 Components

| Name        | Role                             |
|-------------|----------------------------------|
| `DevNode`   | Local Ethereum dev chain (geth)  |
| `DevServer` | Go backend for wallet, WS relay  |
| `TestClient`| CLI-style test code in Go        |
| `ChainUI`   | Scala.js frontend for live data  |

---

## 🛸 Future Phase Names

| Phase | Codename    | Focus                              |
|-------|-------------|-------------------------------------|
| 0     | `Ignition`  | Bootstrap system + first socket    |
| 1     | `Liftoff`   | First decoded log shown in UI      |
| 2     | `Orbit`     | Contract filtering, UI layout      |
| 3     | `Telemetry` | Live stats, metrics, data panels   |
| 4     | `Docking`   | UI → TX interaction (send txs)     |
| 5     | `Expansion` | Scale to multi-node/cross-chain    |
