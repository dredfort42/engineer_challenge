# Interview Task: Golang Development for Event Handling and Docker Integration

## Objective

The task is designed to evaluate your ability to develop microservices using Golang and Docker, focusing on real-time data handling and querying from a time-series database.

---

## Task Overview

You are required to create a complete microservices setup using Golang and Docker. This setup will include:

-   A **daemon service** to generate JSON-formatted security events.
-   A **client service** to query and display critical events.
-   A **writer service** to write to InfluxDB.
-   A **reader service** to read data from InfluxDB.
-   **NATS** server for messaging between microservices.
-   **InfluxDB** for event storage.

---

## Specifications

### Daemon Service:

-   Develop a Golang daemon that runs continuously to generate random JSON events.
-   Each event must contain the following fields:
    -   `criticality` (integer)
    -   `timestamp` (ISO 8601 format string)
    -   `eventMessage` (string)
-   Publish these events on NATS using the subject `events`.

### Client Service:

-   Create a Golang client that queries, using NATS messages, to retrieve the **last 10 events with the criticality level higher than `x`**, where `x` is set in an environmet variable.
-   Display these events in a clear and concise format of your choice.
-   Package this client into a **Docker container**.

### Reader Service:

-   Create a Golang client that services NATS requests to query **InfluxDB** to retrieve the **last `x` events with the criticality higher than `y`**. where `x` and `y` are as specified in the request
-   respond via NATS with the requested events
-   Package this client into a **Docker container**.

### Writer Service:

-   Create a Golang client that subscribe to events published on NATS and write them on **InfluxDB**
-   Package this client into a **Docker container**.

### Docker Network Setup:

Define a **Docker Compose** file to manage a network comprising:

-   Your daemon service
-   Your client service
-   Your reader service
-   Your writer service
-   **NATS**
-   An **InfluxDB** container

---

## Requirements

1. **Security:** Implement secure-by-design principles in your code.
2. **Quality Assurance:** Utilize techniques to ensure code accuracy and quality.
3. **Best Practices:** Follow best coding practices and explain your choices.

---

## Deliverables

1. **Source code** for both the daemon and client services.
2. **Dockerfiles** and **Docker Compose** file for network setup.
3. A **README file** detailing:
    - Setup instructions
    - How to run the services
    - Any dependencies

---

## Presentation

You will have a **90-minute face-to-face interview** to:

1. Demonstrate your microservices setup and its functionality.
2. Discuss your design and implementation decisions.
3. Showcase your code and the outputs of both services.
4. Explain the techniques used for maintaining code quality and security.
5. Walk through the Docker integration and operations.

---

## Evaluation Criteria

1. **Functionality:** Completeness of the services and correctness of outputs.
2. **Code Quality:** Readability, structure, and adherence to best practices.
3. **Security:** Implementation of secure coding practices.
4. **Presentation Skills:** Clarity and depth of the explanation regarding design choices and implementations.

---

This task is an opportunity to showcase your **technical skills** as well as your ability to **communicate complex ideas effectively**. We look forward to your innovative solutions and a detailed discussion during your interview.
