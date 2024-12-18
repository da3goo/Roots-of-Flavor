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


// Логика для формы логина
document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

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
        const result = await response.json();

        if (response.ok) {
            console.log("Login successful for user:", result);
            userIsOnline = true;
            currentUser = result;
            document.getElementById('loginModal').style.display = 'none'; 
            document.getElementById('logInText').innerHTML = '<a href="profile.html">Profile</a>';
            document.getElementById('logOutText').style.display = 'inline'; 
        } else {
            
            console.log("Login failed:", result.error);
            if (result.error && result.error === "Invalid email or password") {
                alert('Incorrect email or password. Please try again.');
            } else {
                alert(result.error || 'Login failed');
            }
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Server error. Please try again later.');
    }
});




// Логика для кнопки "Logout"
document.getElementById('logOutText').addEventListener('click', async function() {
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
            document.getElementById('logInText').innerHTML = '<span>Login</span>';
            document.getElementById('logOutText').style.display = 'none'; 
        } else {
            
            const result = await response.json();
            console.error("Logout failed:", result);
            
        }
    } catch (error) {
        
        console.error('Logout error:', error);
        
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
            document.getElementById('logInText').innerHTML = '<a href="profile.html id="profileStyles">Profile</a>';
            document.getElementById('logOutText').style.display = 'inline'; 

            
        } else {
            userIsOnline = false;
            currentUser = null;
            document.getElementById('logInText').innerHTML = '<span>Login</span>';
            document.getElementById('logOutText').style.display = 'none'; 
        }
    } catch (error) {
        console.error('Error checking session:', error);
        

        
    }
}

checkUserOnLoad();








