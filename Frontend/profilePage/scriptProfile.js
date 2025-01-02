document.addEventListener('DOMContentLoaded', () => {
    checkSession();  

    const logOutButton = document.getElementById('logOutText');
    if (logOutButton) {
        logOutButton.addEventListener('click', async function () {
            try {
                console.log("Logout button clicked, sending request to server...");
                const response = await fetch('http://localhost:8080/logout', {
                    method: 'POST',
                    credentials: 'include'
                });

                if (response.ok) {
                    console.log("Logout successful.");
                    window.location.href = '/Frontend/client/main_page.html'; // Укажите страницу, куда перенаправить
                } else {
                    const result = await response.json();
                    console.error("Logout failed:", result);
                    alert(result.message || 'Failed to log out. Please try again.');
                }
            } catch (error) {
                console.error('Logout error:', error);
                alert('Server error. Please try again later.');
            }
        });
    }
    document.getElementById('delete-btn').addEventListener('click', async function () {
        if (confirm('Are you sure you want to delete your profile? This action cannot be undone.')) {
            try {
                const response = await fetch('http://localhost:8080/deleteUser', {
                    method: 'DELETE',
                    credentials: 'include'
                });
    
                if (response.ok) {
                    alert('Your profile has been deleted successfully.');
                    window.location.href = '/Frontend/client/main_page.html'; 
                } else {
                    const result = await response.json();
                    alert(result.message || 'Failed to delete profile.');
                }
            } catch (error) {
                console.error('Error deleting profile:', error);
                alert('Unable to delete your profile. Please try again later.');
            }
        }
    });
    
    document.getElementById('edit-btn').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModal');
        modal.style.display = 'block';
        loadUserData(); 
    });
    
    document.getElementById('close-btn').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModal');
        modal.style.display = 'none';
    });
    document.getElementById('editProfileForm').addEventListener('submit', async function (event) {
        event.preventDefault(); 
    
        const updatedFullname = document.getElementById('fullname').value;
    
        
        try {
            const sessionResponse = await fetch('http://localhost:8080/checksession', {
                method: 'GET',
                credentials: 'include'
            });
    
            if (sessionResponse.ok) {
                const userData = await sessionResponse.json();
    
                
                if (userData.fullname === updatedFullname) {
                    alert('This name is already in use. No changes were made.');
                    window.location.href = 'profile.html'; 
                    return;
                }
            } else {
                alert('Session expired or not logged in. Please log in.');
                window.location.href = '/Frontend/client/main_page.html'; 
                return;
            }
        } catch (error) {
            console.error('Error checking session:', error);
            alert('Unable to verify session. Please try again later.');
            return;
        }
    
        
        const response = await fetch('http://localhost:8080/updateName', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ fullname: updatedFullname })
        });
    
        if (response.ok) {
            alert('Profile updated successfully');
            window.location.reload(); 
        } else {
            const result = await response.json();
            alert(result.message || 'Failed to update profile');
        }
    });


    
});


async function checkSession() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const userData = await response.json();
            console.log("User data:", userData);

            // Заполнение профиля пользователя
            document.getElementById('profile-fullname').textContent = userData.fullname;
            document.getElementById('profile-email').textContent = userData.email;

            const createdAt = new Date(userData.createdAt);
            document.getElementById('profile-createdAt').textContent = isNaN(createdAt)
                ? 'Invalid account creation date'
                : `Account created on: ${createdAt.toLocaleString()}`;

            // Проверка статуса пользователя
            if (userData.userStatus === 'admin') {
                const adminBtn = document.getElementById('admin-btn');
                if (adminBtn) {
                    adminBtn.style.display = 'block'; // Показываем кнопку
                }
            }
        } else {
            alert('Session expired or not logged in. Please log in.');
            window.location.href = '/Frontend/client/main_page.html';
        }
    } catch (error) {
        console.error('Error checking session:', error);
        alert('Unable to verify session. Please try again later.');
        window.location.href = '/Frontend/client/main_page.html';
    }
}


async function loadUserData() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const userData = await response.json();
            console.log("User data:", userData);

            
            document.getElementById('fullname').placeholder = userData.fullname;
            document.getElementById('fullname').value = userData.fullname;

            
            const updatedAt = new Date(userData.updatedFullnameAt); 
            document.getElementById('updatedAt').textContent = isNaN(updatedAt)
                ? 'Invalid date format'
                : `Last updated: ${updatedAt.toLocaleString()}`;
        } else {
            alert('Session expired or not logged in. Please log in.');
            window.location.href = '/Frontend/client/main_page.html';
        }
    } catch (error) {
        console.error('Error loading user data:', error);
    }
}







