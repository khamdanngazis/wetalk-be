# Chat Application System Architecture

## **Overview**
This document outlines the architecture of a chat application with the following features:
1. User registration with unique identifiers (`@username`).
2. User login system.
3. Real-time messaging between users.
4. Chat history retrieval (sorted by time in descending order).
5. Kafka-based notification service.
6. Message status tracking (`sent`, `delivered`, `read`).
7. Scalable design to handle high traffic and large user base.

---

## **System Architecture**

### **Components**
1. **Frontend (React.js)**
   - Handles user interactions and connects to the backend via:
     - REST APIs for registration, login, and chat history.
     - WebSocket for real-time messaging and status updates.
   - Displays:
     - Chat messages.
     - Notifications.
     - Message statuses (`sent`, `delivered`, `read`).

2. **Backend (Golang)**
   - Manages:
     - User authentication and session handling.
     - Real-time messaging via WebSocket.
     - Message statuses and chat history.
     - Kafka integration for decoupled notification processing.
   - Exposes REST APIs for frontend interactions.

3. **Database (PostgreSQL)**
   - Stores:
     - User information.
     - Chat messages.
     - Message statuses.

4. **Redis**
   - Caches:
     - Active WebSocket sessions.
     - Temporary unread notifications for offline users.

5. **Kafka**
   - **Topics**:
     - `message-events`: Tracks message lifecycle events (`sent`, `delivered`, `read`).
     - `unread-notifications`: Manages notifications for offline users.
   - Decouples the chat and notification services for scalability.

6. **Notification Service**
   - Kafka consumer that:
     - Processes events for unread notifications.
     - Pushes notifications to online users via WebSocket.
     - Stores unread notifications in Redis for offline users.

---

### **Architecture Diagram**

```plaintext
+-------------------+       +--------------------+       +--------------------+
|  React.js (Client)|<----->|    WebSocket       |<----->| Notification Service|
|                   |       |                    |       |   (Kafka Consumer)  |
+-------------------+       +--------------------+       +--------------------+
        |                              ^                         ^
        | REST API                     |                         |
        v                              |                         |
+-------------------+       +--------------------+       +--------------------+
| Backend (Golang)  |<----->|     Kafka          |<----->|   Redis (Cache)    |
|  REST & WebSocket |       |   (Message Bus)    |       |   Unread Messages  |
+-------------------+       +--------------------+       +--------------------+
        |                                                       |
        v                                                       |
+-------------------+                                           |
|  PostgreSQL       |<-----------------------------------------+
| Database          |
+-------------------+
```

## **Flow Details**

### **1. User Registration and Login**
- Users can register using their email address, which is then associated with a unique `@username`.
- After registration, users can log in via REST API. Upon successful login:
  - A session is created and cached in Redis to facilitate WebSocket management (for real-time communication).

### **2. Sending a Message**
- The sender sends a message via WebSocket to the backend.
- **Backend**:
  - Saves the message to the `Messages` table in the database with `status = sent`.
  - Publishes a `message-sent` event to the Kafka `message-events` topic.
  - If the receiver is online, the backend forwards the message to the receiver's WebSocket connection.

### **3. Delivering a Message**
- When the receiver's WebSocket receives the message, it is marked as `delivered` in the backend.
- **Backend**:
  - Updates the message status to `delivered` in the `Messages` table.
  - Publishes a `message-delivered` event to Kafka to notify the sender and the receiver.

### **4. Reading a Message**
- When the receiver reads a message, a `read-receipt` is sent back to the backend via WebSocket.
- **Backend**:
  - Updates the message status to `read` in the `Messages` table.
  - Publishes a `message-read` event to Kafka to notify both the sender and the receiver.

### **5. Offline Notifications**
- If the receiver is offline:
  - The `message-sent` event is published to Kafka, which the **Notification Service** consumes.
  - The Notification Service stores the notification in Redis for offline users.
  - When the receiver reconnects to the WebSocket, the notification is delivered in real-time.

### **6. Chat History Retrieval**
- The frontend can request the chat history between two users via a REST API endpoint.
  - The backend retrieves the chat messages from the `Messages` table.
  - Messages are sorted by `created_at DESC` to show the most recent messages first.

---

## **Database Schema**

### **Users Table**
| Column        | Type      | Constraints             |
|---------------|-----------|-------------------------|
| id            | UUID      | Primary Key             |
| username      | VARCHAR   | Unique, Not Null        |
| email         | VARCHAR   | Unique, Not Null        |
| password      | VARCHAR   | Hashed, Not Null        |
| created_at    | TIMESTAMP | Default: CURRENT_TIME   |
| deleted_at    | TIMESTAMP  | Default: null   |

### **Messages Table**
| Column        | Type       | Constraints             |
|---------------|------------|-------------------------|
| id            | UUID       | Primary Key             |
| sender_id     | UUID       | Foreign Key (users.id)  |
| receiver_id   | UUID       | Foreign Key (users.id)  |
| content       | TEXT       | Not Null                |
| status        | ENUM       | `sent`, `delivered`, `read` |
| created_at    | TIMESTAMP  | Default: CURRENT_TIME   |
| updated_at    | TIMESTAMP  | Default: CURRENT_TIME   |
| deleted_at    | TIMESTAMP  | Default: null   |

### **Kafka Topics**
1. **`message-events`**:
   - Tracks message lifecycle events (`sent`, `delivered`, `read`).
   - Partitions by `message_id`.

2. **`unread-notifications`**:
   - Handles offline notifications for unread messages.
   - Partitions by `receiver_id`.

---

## **Key Benefits**
1. **Scalability**: Kafka ensures the system can handle high traffic and a large volume of messages.
2. **Asynchronous Processing**: Decouples messaging and notification processing, enabling horizontal scalability.
3. **Real-Time Updates**: WebSocket connection ensures instant message delivery and message status updates.
4. **Offline Notifications**: Kafka and Redis provide reliable storage and delivery of notifications for offline users.
5. **Event-Driven Architecture**: Kafka allows for easy extension and future scalability, including analytics or additional features.

