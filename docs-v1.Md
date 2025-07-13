Example WorkflowUser Setup:User installs the app on their phone and PC.
They sign in with the same account (cloud) or pair devices via QR code (local).

Device Discovery:Phone and PC detect each other via mDNS (local) or check in with the server (cloud).
The app shows a list of available devices.

Sharing a Link:On the phone, the user shares a link via the app’s share sheet.
The app sends the link to the PC (directly or via the server).
The PC receives the link, shows a notification, and opens it in the browser.

Reverse Sharing:On the PC, the user right-clicks a link in the browser and selects “Send to Phone.”
The phone receives the link and opens it.

1. Backend RequirementsThe backend will:Authenticate users: Allow users to sign in and associate devices with their accounts.
Register devices: Store device metadata (e.g., device ID, name) for discovery.
Handle device discovery: Provide a list of a user’s available devices.
Facilitate link transfer: Allow devices to send/receive links, with real-time delivery via WebSocket.
Ensure security: Use JWT for authentication and encrypt link payloads.

2. Tech StackLanguage: Go (v1.21 or later)
Framework: Gin (REST API handling)
Database: SQLite (lightweight, no external setup)
WebSocket: gorilla/websocket for real-time link delivery
Authentication: JWT (JSON Web Tokens) with golang-jwt/jwt
Encryption: AES for link payloads
Dependencies:github.com/gin-gonic/gin
github.com/mattn/go-sqlite3
github.com/gorilla/websocket
github.com/golang-jwt/jwt/v5

3. Project Structure

Step 2: Main Application (main.go)This file sets up the Gin server, initializes the database, and defines routes.

Documentation:Initializes SQLite database using db.InitDB().
Sets up a Gin router with a base path /api.
Defines endpoints for user registration (/register), login (/login), device management (/devices), and link sharing (/links).
WebSocket endpoint (/ws) for real-time link delivery.
Runs the server on port 8080.

Step 3: Database Setup (db/db.go)This file initializes SQLite and defines database operations.

Documentation:Opens a connection to a SQLite database (linkshare.db).
Creates three tables:users: Stores user email and hashed password.
devices: Stores user devices with a reference to user_id.
links: Stores shared links with references to user_id, from_device_id, and to_device_id.

Step 4: Models (models/*.go)Define data structures for users, devices, and links.

Documentation:User: Represents a user with an ID, email, and hashed password.
Device: Represents a device with an ID, user ID, name, and last seen timestamp.
Link: Represents a shared link with sender and receiver device IDs, URL, and creation time.

Step 5: Authentication Middleware (middleware/auth.go)This middleware validates JWT tokens for protected routes.

Documentation:Checks for a Bearer token in the Authorization header.
Validates the JWT using a secret key.
Extracts user_id from the token and stores it in the context.
Aborts with a 401 status if the token is missing or invalid.

Step 6: Authentication Handlers (handlers/auth.go)Handles user registration and login.

Documentation:Register: Creates a new user by hashing their password (SHA-256 for simplicity; use bcrypt in production) and storing it in the users table.
Login: Verifies credentials and returns a JWT with the user’s ID, valid for 24 hours.

Step 7: Device Handlers (handlers/device.go)Handles device registration and discovery.

Documentation:RegisterDevice: Registers a device with a name for the authenticated user.
ListDevices: Returns a list of devices associated with the user, used for device discovery.

Step 8: Link Handlers (handlers/link.go)Handles link sending and WebSocket communication.

Documentation:SendLink: Stores a link in the database and sends it to the target device via WebSocket if the device is online.
WebSocketHandler: Upgrades an HTTP connection to WebSocket, registers the device’s connection, and updates its last_seen timestamp. Keeps the connection open to receive links.
connections: A thread-safe map to track active WebSocket connections by user_id and device_id.

Step 9: Encryption Utility (utils/crypto.go)Encrypts link URLs for secure transfer (optional).

Documentation:Encrypt: Encrypts a string (e.g., URL) using AES-CFB with a provided key.
Decrypt: Decrypts an AES-encrypted string.
Note: In practice, use a secure key management system (e.g., store keys in environment variables).

5. How Devices Know Each Other’s AvailabilityDevice Registration: When a device starts the app, it sends a POST /devices request with its name, registering it in the devices table.
WebSocket Connection: The device connects to the /ws endpoint with its device_id, signaling it’s online. The backend updates the last_seen timestamp.
Discovery: The client app calls GET /devices to get a list of the user’s devices, showing only those with recent last_seen timestamps (e.g., within the last minute).
Link Transfer: When a link is sent (POST /links), the backend checks if the target device has an active WebSocket connection. If online, it sends the link immediately; if offline, it stores the link for later retrieval.

6. DeploymentBuild the Backend:bash

