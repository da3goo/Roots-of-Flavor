document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('logOutText').addEventListener('click', async function () {
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
});
