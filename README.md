# Qoqos Backend Infrastructure

This repository contains the backend services and infrastructure for **Qoqos**, a tech repair & buyback platform. 

## 🏗 Architecture Overview

The backend is designed as a multi-service architecture running in Docker. It integrates traditional REST APIs with modern AI workflows:

- **API Service (Node.js)**: Core business logic and endpoints connecting to the Next.js frontend.
- **PostgreSQL**: Primary relational database for user data, orders, and application state.
- **n8n**: Workflow automation platform for orchestrating complex business processes.
- **Ollama**: Local execution of LLMs for natural language processing tasks and user interactions.
- **Qdrant**: Vector database for RAG (Retrieval-Augmented Generation) and semantic search functionality.
- **Telegram Bot**: Python-based integration for user notifications, escalations, and direct AI interaction.
- **Nginx**: Reverse proxy to route traffic securely between services.

## 🚀 Development Approach & Philosophy

This infrastructure was rapidly prototyped and iterated upon using **AI-assisted development tools** (LLMs, Cursor, etc.). 
I acted as the **System Architect**, leveraging these tools to drastically accelerate development while ensuring all components integrate seamlessly into a robust, containerized environment (Docker Compose). This project demonstrates the ability to conceptualize, design, and orchestrate complex microservice architectures with modern AI capabilities.

## ⚙️ Getting Started

### Prerequisites
- Docker & Docker Compose
- Node.js 18+ (for local API development)
- Python 3.10+ (for AI scripts & bots)

### Running Locally
1. Clone the repository and configure your `.env` file (ask the administrator for the `.env` template if not provided).
2. Start the infrastructure via Docker Compose:
   ```bash
   docker-compose up -d
   ```
3. Initialize the database using the provided `seed.sql` script if needed.

## 🧠 AI Integration (RAG)
The backend utilizes Qdrant and Ollama to power the `valera_brain.py` assistant. It processes unstructured data and provides automated, context-aware responses, which can be escalated to a human operator via the Telegram Bot integration.
