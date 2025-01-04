


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







//server


// Login
document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    
    const blockedUntil = localStorage.getItem('blockedUntil');
    if (blockedUntil && new Date().getTime() < blockedUntil) {
        const timeLeft = Math.ceil((blockedUntil - new Date().getTime()) / 1000);
        alert(`You have to wait ${timeLeft} seconds before trying again.👿👿👿`);
        return;
    }

    const email = document.getElementById('emailInput').value;
    const password = document.getElementById('passwordInput').value;

    try {
        console.log("Sending login request for email:", email);
        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
            credentials: 'include'
        });

        console.log("Login response status:", response.status);

        if (response.ok) {
            const result = await response.json();
            console.log("Login successful for user:", result);
            userIsOnline = true;
            currentUser = result;
            document.getElementById('loginModal').style.display = 'none';
            document.getElementById('logInText').innerHTML = '<a href="/Frontend/profilePage/profile.html">Profile</a>';
        } else if (response.status === 429) {
            
            const resetTime = response.headers.get('X-RateLimit-Reset');
            const resetDate = new Date(parseInt(resetTime) * 1000);
            const timeToWait = Math.max(0, resetDate - new Date());
            const secondsToWait = Math.ceil(timeToWait / 1000);

            const blockedUntil = new Date().getTime() + 10000; //10 seconds
            localStorage.setItem('blockedUntil', blockedUntil);

            alert(`Rate limit exceeded. Please try again in ${secondsToWait} seconds👿👿👿.`);
            setTimeout(() => {
                alert("I said to try again after 10 seconds!!!👿👿👿");
            }, 10000);

        } else {
            const result = await response.json();
            console.log("Login failed:", result.error);

            if (result.error) {
                alert(result.error);
            } else {
                alert('Login failed. Please try again.');
            }
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Could not connect to the server. Please try again later.');
    }
});










// Логика для кнопки "Logout"
document.getElementById('logOutText').addEventListener('click', async function () {
    try {
        console.log("Logout button clicked, sending request to server...");

        const response = await fetch('http://localhost:8080/logout', {
            method: 'POST',
            credentials: 'include'
        });

        if (response.ok) {
            console.log("Logout successful.");
            userIsOnline = false;
            currentUser = null;
            window.location.href = '/Frontend/client/main_page.html';
        } else {
            const result = await response.json();
            console.error("Logout failed:", result);
            alert(result.message || 'Failed to log out. Please try again.');
        }
    } catch (error) {
        console.error('Logout error:', error);
        alert('Could not connect to the server. Please try again later.');
    }
});



async function checkUserOnLoad() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const result = await response.json();
            userIsOnline = true;
            currentUser = result;
            document.getElementById('loginModal').style.display = 'none';
            document.getElementById('logInText').innerHTML = '<a href="/Frontend/profilePage/profile.html" id="profileStyles">Profile</a>';
            document.getElementById('logOutText').style.display = 'inline';
        } else {
            userIsOnline = false;
            currentUser = null;
            document.getElementById('logInText').innerHTML = '<span>Login</span>';
            document.getElementById('logOutText').style.display = 'none';
        }
    } catch (error) {
        console.error('Error checking session:', error);
        alert('Could not connect to the server. Please try again later.');
    }
}


checkUserOnLoad();



//Reg
document.getElementById('registerForm').addEventListener('submit', async function (event) {
    event.preventDefault();

    const email = document.getElementById('emailRegisterInput').value;
    const password = document.getElementById('passwordRegisterInput').value;
    const fullname = document.getElementById('nameRegiserInput').value;

    try {
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password, fullname })
        });

        const result = await response.json();

        if (response.ok) {
            console.log("Registration successful:", result);
            alert(result.message || 'Registration successful!'); 
            window.location.href = '/Frontend/client/main_page.html';
        } else {
            console.error("Registration failed:", result);
            alert(result.error || "Registration failed!");
        }
    } catch (error) {
        console.error('Registration error:', error);
        alert('Could not connect to the server. Please try again later.');
    }
});









