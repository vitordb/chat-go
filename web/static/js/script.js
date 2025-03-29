document.addEventListener('DOMContentLoaded', () => {
    // DOM elements
    const authContainer = document.getElementById('auth-container');
    const chatContainer = document.getElementById('chat-container');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterLink = document.getElementById('show-register');
    const showLoginLink = document.getElementById('show-login');
    const logoutBtn = document.getElementById('logout-btn');
    const messageForm = document.getElementById('message-form');
    const messageInput = document.getElementById('message-input');
    const messagesContainer = document.getElementById('messages');
    const chatroomList = document.getElementById('chatroom-list');
    const createRoomBtn = document.getElementById('create-room-btn');
    const newRoomNameInput = document.getElementById('new-room-name');

    // State
    let currentUser = null;
    let currentChatroom = null;
    let socket = null;

    // Check if user is authenticated
    checkAuth();

    // Event listeners
    showRegisterLink.addEventListener('click', (e) => {
        e.preventDefault();
        document.querySelectorAll('.form-container').forEach(container => {
            container.style.display = container.style.display === 'none' ? 'block' : 'none';
        });
    });

    showLoginLink.addEventListener('click', (e) => {
        e.preventDefault();
        document.querySelectorAll('.form-container').forEach(container => {
            container.style.display = container.style.display === 'none' ? 'block' : 'none';
        });
    });

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                checkAuth();
            } else {
                const error = await response.text();
                alert(`Login failed: ${error}`);
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    });

    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('register-username').value;
        const password = document.getElementById('register-password').value;

        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                checkAuth();
            } else {
                const error = await response.text();
                alert(`Registration failed: ${error}`);
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    });

    logoutBtn.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/auth/logout', { method: 'POST' });
            if (response.ok) {
                disconnectSocket();
                currentUser = null;
                currentChatroom = null;
                authContainer.style.display = 'block';
                chatContainer.style.display = 'none';
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    });

    messageForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const content = messageInput.value.trim();
        if (content && socket) {
            socket.send(JSON.stringify({ content }));
            messageInput.value = '';
        }
    });

    createRoomBtn.addEventListener('click', async () => {
        const name = newRoomNameInput.value.trim();
        if (name) {
            try {
                const response = await fetch('/api/chatrooms', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name })
                });

                if (response.ok) {
                    newRoomNameInput.value = '';
                    fetchChatrooms();
                } else {
                    const error = await response.text();
                    alert(`Failed to create chatroom: ${error}`);
                }
            } catch (error) {
                alert(`Error: ${error.message}`);
            }
        }
    });

    // Functions
    async function checkAuth() {
        try {
            const response = await fetch('/api/auth/check');
            if (response.ok) {
                currentUser = await response.json();
                authContainer.style.display = 'none';
                chatContainer.style.display = 'block';
                fetchChatrooms();
            } else {
                authContainer.style.display = 'block';
                chatContainer.style.display = 'none';
            }
        } catch (error) {
            console.error(`Error checking auth: ${error.message}`);
            authContainer.style.display = 'block';
            chatContainer.style.display = 'none';
        }
    }

    async function fetchChatrooms() {
        try {
            const response = await fetch('/api/chatrooms');
            if (response.ok) {
                const chatrooms = await response.json();
                renderChatrooms(chatrooms);
                if (chatrooms.length > 0 && !currentChatroom) {
                    joinChatroom(chatrooms[0].id);
                }
            }
        } catch (error) {
            console.error(`Error fetching chatrooms: ${error.message}`);
        }
    }

    function renderChatrooms(chatrooms) {
        chatroomList.innerHTML = '';
        chatrooms.forEach(chatroom => {
            const li = document.createElement('li');
            li.textContent = chatroom.name;
            li.dataset.id = chatroom.id;
            if (currentChatroom && chatroom.id === currentChatroom.id) {
                li.classList.add('active');
            }
            li.addEventListener('click', () => joinChatroom(chatroom.id));
            chatroomList.appendChild(li);
        });
    }

    async function joinChatroom(chatroomId) {
        try {
            const response = await fetch(`/api/chatrooms/${chatroomId}`);
            if (response.ok) {
                const chatroom = await response.json();
                
                // Disconnect from previous chatroom
                disconnectSocket();
                
                // Clear messages
                messagesContainer.innerHTML = '';
                
                // Update current chatroom
                currentChatroom = chatroom;
                
                // Update UI
                document.querySelectorAll('#chatroom-list li').forEach(li => {
                    li.classList.toggle('active', li.dataset.id === chatroomId);
                });
                
                // Connect to new chatroom
                connectSocket(chatroomId);
            }
        } catch (error) {
            console.error(`Error joining chatroom: ${error.message}`);
        }
    }

    function connectSocket(chatroomId) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws/${chatroomId}`;
        
        socket = new WebSocket(wsUrl);
        
        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            renderMessage(message);
            
            // Scroll to bottom
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        };
        
        socket.onclose = () => {
            console.log('WebSocket connection closed');
        };
        
        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    function disconnectSocket() {
        if (socket) {
            socket.close();
            socket = null;
        }
    }

    function renderMessage(message) {
        const messageDiv = document.createElement('div');
        messageDiv.classList.add('message');
        
        // Add class based on message type
        messageDiv.classList.add(`message-${message.type}`);
        
        // Add class if the message is from the current user
        if (message.type === 'chat' && message.user_id === currentUser.id) {
            messageDiv.classList.add('own');
        }
        
        // Add message content
        if (message.type !== 'system') {
            const usernameDiv = document.createElement('div');
            usernameDiv.classList.add('username');
            usernameDiv.textContent = message.username;
            messageDiv.appendChild(usernameDiv);
        }
        
        const contentDiv = document.createElement('div');
        contentDiv.classList.add('content');
        contentDiv.textContent = message.content;
        messageDiv.appendChild(contentDiv);
        
        const timeDiv = document.createElement('div');
        timeDiv.classList.add('time');
        timeDiv.textContent = new Date(message.created_at).toLocaleTimeString();
        messageDiv.appendChild(timeDiv);
        
        messagesContainer.appendChild(messageDiv);
    }
}); 