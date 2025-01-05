document.addEventListener('DOMContentLoaded', () => {
    checkSession();  
    document.getElementById("admin-btn").addEventListener("click", function () {
        window.location.href = "/Frontend/client/adminPage.html";
    });
    

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
    });
    document.getElementById('changeNameBtn').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModalName');
        const previousModal = document.getElementById('editProfileModal');
        previousModal.style.display = 'none';
        modal.style.display = 'block';
        loadUserData(); 
    });
    document.getElementById('changePasswordBtn').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModalPassword');
        const previousModal = document.getElementById('editProfileModal');
        previousModal.style.display = 'none';
        modal.style.display = 'block';
        loadUserData(); 
    });
    
    document.getElementById('close-btn-edit-profile').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModal');
        modal.style.display = 'none';
    });
    document.getElementById('close-btn-edit-profile-name').addEventListener('click', () => {
        const modal = document.getElementById('editProfileModalName');
        const prevModal = document.getElementById('editProfileModal');
        modal.style.display = 'none';
        prevModal.style.display = 'block';
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


    document.getElementById("submitChangingPassword").addEventListener("click", async (event) => {
        event.preventDefault();
    
        const oldPassword = document.getElementById("oldpassword").value.trim();
        const newPassword = document.getElementById("newpassword").value.trim();
        const newPasswordRetype = document.getElementById("newpasswordretype").value.trim();
    
        if (newPassword !== newPasswordRetype) {
            alert("New passwords do not match!");
            return;
        }
    
        try {
            const response = await fetch("http://localhost:8080/changePassword", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    oldPassword,
                    newPassword,
                    newPasswordRetype,
                }),
                credentials: "include",
            });
    
            const result = await response.json();
    
            if (!response.ok) {
                alert(result.message || "Failed to update password");
            } else {
                const updatedAtElement = document.getElementById("updatedPasswordAt");
                updatedAtElement.innerText = `Last updated: ${result.updated_at}`;
                alert(result.message);
            }
        } catch (error) {
            console.error("Error updating password:", error);
            alert("An error occurred. Please try again.");
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

            
            document.getElementById('profile-fullname').textContent = userData.fullname;
            document.getElementById('profile-email').textContent = userData.email;

            const createdAt = new Date(userData.createdAt);
            document.getElementById('profile-createdAt').textContent = isNaN(createdAt)
                ? 'Invalid account creation date'
                : `Account created on: ${createdAt.toLocaleString()}`;

            
            if (userData.userStatus === 'admin') {
                const adminBtn = document.getElementById('admin-btn');
                if (adminBtn) {
                    adminBtn.style.display = 'block'; 
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

            const updatedAtFullname = new Date(userData.updatedFullnameAt);
            document.getElementById('updatedAt').textContent = isNaN(updatedAtFullname)
                ? 'Invalid date format'
                : `Fullname last updated: ${updatedAtFullname.toLocaleString()}`;

            if (userData.updatedAt) {
                const updatedAtPassword = new Date(userData.updatedAt);
                document.getElementById('updatedPasswordAt').textContent = isNaN(updatedAtPassword)
                    ? 'Invalid date format'
                    : `Password last updated: ${updatedAtPassword.toLocaleString()}`;
            } else {
                document.getElementById('updatedPasswordAt').textContent = 'Password has not been updated yet.';
            }

        } else {
            alert('Session expired or not logged in. Please log in.');
            window.location.href = '/Frontend/client/main_page.html';
        }
    } catch (error) {
        console.error('Error loading user data:', error);
    }
}
