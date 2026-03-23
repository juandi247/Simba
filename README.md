# Simba 🦁 

Simba is a **Distributed Coordinator**, inspired by Google Chubby, designed for **learning and for fun**.  

Currently, Simba is a **hand-implemented Raft consensus algorithm**, written entirely from scratch (no external libraries), paired with **Deterministic Simulation Testing (DST)** for verification and validation.  

## Deterministic Simulation Testing (DST)

DST allows simulating the behavior of nodes in a fully deterministic way. This means that every simulation run with the same seed and scenario produces exactly the same results, which is crucial for testing consensus algorithms. Simba reuses the same core logic both in real implementations and in simulations, making testing reliable and reproducible.  

## Inspiration

The project also draws inspiration from **TigerBeetle** and its **Tiger-style approach**, focusing on static allocation, fixed sizes, and determinism, even though Go is garbage collected. This allows having everything pre-allocated, deterministic, and predictable in the simulator.  

## Current Focus

Simba is currently focused on:
- Building the Raft core **entirely by hand**.  
- Developing the deterministic simulator for testing.  
- Reusing the core logic in both real and simulated environments.  

> **Status:** under construction — Raft and its simulator are the current priority.  
