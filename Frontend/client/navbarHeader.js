let userIsOnline = false;
let currentUser = null;

// Handler for showing the login modal or "Log out" button
document.getElementById('logInText').addEventListener('click', function() {
    if (!userIsOnline) {
        document.getElementById('loginModal').style.display = 'block';
    } else {
        document.getElementById('logOutText').style.display = 'inline';
    }
});

// Handler for the "Log out" button
document.getElementById('logOutText').addEventListener('click', async function() {
    try {
        const response = await fetch('http://localhost:3000/logout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: currentUser1.email })  
        });

        const result = await response.json();

        if (result.success) {
            currentUser1 = null;  
            document.getElementById('logInText').innerHTML = '<span>Login</span>';
            document.getElementById('logOutText').style.display = 'none';  
            userIsOnline = false;  
        } else {
            alert(result.message);
        }
    } catch (error) {
        console.error('Logout error:', error);
        alert('Server error. Please try again later.');
    }
});

// Closing modal windows when clicking outside the window
window.addEventListener('click', function(event) {
    if (event.target === document.getElementById('loginModal')) {
        document.getElementById('loginModal').style.display = 'none';
    } else if (event.target === document.getElementById('registerModal')) {
        document.getElementById('registerModal').style.display = 'none';
    }
});

// Sending login data
document.getElementById('loginForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const email = document.getElementById('emailInput').value;
    const password = document.getElementById('passwordInput').value;

    try {
        const response = await fetch('http://localhost:3000/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const result = await response.json();

        if (result.success) {
            currentUser1 = {
                username: result.username,
                email: email
            };

            localStorage.setItem('userEmail', email);

            document.getElementById('logInText').innerHTML = `<span>${currentUser1.username}</span>`;
            document.getElementById('loginModal').style.display = 'none';
            userIsOnline = true;
        } else {
            alert(result.message);
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Server error. Please try again later.');
    }
});

// Handler for sending registration data
document.getElementById('registerForm').addEventListener('submit', async function(event) {
    event.preventDefault();

    const username = document.getElementById('nameRegiserInput').value;
    const usersurname = document.getElementById('surnameRegisterInput').value;
    const email = document.getElementById('emailRegisterInput').value;
    const password = document.getElementById('passwordRegisterInput').value;
    const userstatus = 'user';

    try {
        const response = await fetch('http://localhost:3000/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, usersurname, email, password, userstatus })
        });

        if (!response.ok) {
            throw new Error('Ошибка на сервере. Статус: ' + response.status);
        }

        const result = await response.json();

        if (result.success) {
            currentUser1 = { username: result.username, email: email };
            
            document.getElementById('logInText').innerHTML = `<span>${currentUser1.username}</span>`;
            document.getElementById('registerModal').style.display = 'none';
            userIsOnline = true;
        } else {
            alert(result.message || 'Ошибка регистрации');
        }
    } catch (error) {
        console.error('Ошибка регистрации:', error);
        alert('Registration error. Please try again later.');
    }
});

// Opening a registration modal window when clicking "Register"
document.getElementById('registerText').addEventListener('click', function(event) {
    event.preventDefault();
    document.getElementById('loginModal').style.display = 'none'; 
    document.getElementById('registerModal').style.display = 'block'; 
});

// Close the registration modal window
document.querySelector('.closeBtnRegisterModal').addEventListener('click', function() {
    document.getElementById('registerModal').style.display = 'none'; 
});

// Click on username to show "Log out" button
document.getElementById('logInText').addEventListener('click', function() {
    if (userIsOnline) {
        document.getElementById('logOutText').style.display = 'inline';
    }
});
// Close the login modal window when clicking on the cross
document.querySelector('.closeBtnLoginModal').addEventListener('click', function() {
    document.getElementById('loginModal').style.display = 'none'; 
});


// Handler for submenu
document.getElementById('foodsLink').addEventListener('click', function(event) {
    event.preventDefault(); 
    
    const menuItem = this.parentElement; 
    const submenu = document.getElementById('foodsSubmenu');
    
    if (!menuItem.classList.contains('active')) {
        menuItem.classList.add('active');
    } else {
        menuItem.classList.remove('active');
    }
});

// Function to check user when page loads
async function checkUserOnLoad() {
    try {
        const email = localStorage.getItem('userEmail'); 

        if (email) {
            const response = await fetch('http://localhost:3000/check-user', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });

            const result = await response.json();

            if (result.success) {
                currentUser1 = {
                    username: result.username,
                    email: email,
                    userStatus: result.userStatus 
                };
                
                document.getElementById('logInText').innerHTML = `<span>${currentUser1.username}</span>`;
                document.getElementById('logOutText').style.display = 'inline'; 

                if (currentUser1.userStatus === 'admin') {
                    document.getElementById('offersLink').style.display = 'inline';
                }

                userIsOnline = true;
            } else {
                document.getElementById('logInText').innerHTML = '<span>Login</span>';
                document.getElementById('logOutText').style.display = 'none';
                userIsOnline = false;
            }
        } else {
            document.getElementById('logInText').innerHTML = '<span>Login</span>';
            document.getElementById('logOutText').style.display = 'none';
            userIsOnline = false;
        }
    } catch (error) {
        console.error('Check user error:', error);
        document.getElementById('logInText').innerHTML = '<span>Login</span>';
        document.getElementById('logOutText').style.display = 'none';
        userIsOnline = false;
    }
}

checkUserOnLoad();


