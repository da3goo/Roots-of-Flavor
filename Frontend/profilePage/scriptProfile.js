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
});


async function checkSession() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const userData = await response.json();
            console.log("User data:", userData);  // Выводим данные, чтобы понять, что приходит с сервера

            // Заполняем данные профиля
            document.getElementById('profile-fullname').textContent = userData.fullname;
            document.getElementById('profile-email').textContent = userData.email;

            // Обрабатываем дату
            console.log("Created at:", userData.createdAt);  // Смотрим, что приходит в поле createdAt
            const createdAt = new Date(userData.createdAt);

            if (isNaN(createdAt)) {
                document.getElementById('profile-createdAt').textContent = 'Invalid account creation date';
            } else {
                document.getElementById('profile-createdAt').textContent = `Account created on: ${createdAt.toLocaleString()}`;
            }
        } else {
            alert('Session expired or not logged in. Please log in.');
            window.location.href = '/Frontend/client/main_page.html'; // Перенаправляем на главную страницу
        }
    } catch (error) {
        console.error('Error checking session:', error);
        alert('Unable to verify session. Please try again later.');
        window.location.href = '/Frontend/client/main_page.html'; // Перенаправляем на главную страницу
    }
}






