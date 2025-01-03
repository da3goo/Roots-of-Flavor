let users = [];     

async function fetchUsers(sortBy = '', emailFilter = '') {
    const url = new URL('http://localhost:8080/getUsers');
    const params = new URLSearchParams();

    if (sortBy) {
        params.append('sort', sortBy);
    }
    if (emailFilter) {
        params.append('email', emailFilter);
    }

    url.search = params.toString();

    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Ошибка при загрузке данных');
        }

        users = await response.json();
        
        renderTable(users); 
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

function renderTable(users) {
    const tableBody = document.getElementById('userTableBody');
    tableBody.innerHTML = ''; 

    users.forEach(user => {
        console.log(user.userstatus); 

        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.fullname}</td>
            <td>${user.email}</td>
            <td>${user.userstatus}</td>
            <td>${new Date(user.createdAt).toLocaleString()}</td>
        `;
        tableBody.appendChild(row);
    });
}


function applyFilters() {
    const sortValue = document.getElementById('sortSelect').value;
    const emailFilter = document.getElementById('emailFilter').value.toLowerCase();

    fetchUsers(sortValue, emailFilter);
}

window.onload = fetchUsers;
        const darkModeToggle = document.getElementById('darkModeBtn');
        const body = document.body;
        
        if (localStorage.getItem('dark-mode') === 'enabled') {
            body.classList.add('dark-mode');
        }

        
        darkModeToggle.addEventListener('click', () => {
        body.classList.toggle('dark-mode');

        
        if (body.classList.contains('dark-mode')) {
            localStorage.setItem('dark-mode', 'enabled');
        } else {
        localStorage.setItem('dark-mode', 'disabled');
        }
});


async function checkSession() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: "GET",
            credentials: "include", 
        });

        console.log("Response status:", response.status); 

        if (!response.ok) {
            alert("The session is invalid. Redirection to the main page.");
            window.location.href = "/Frontend/client/main_page.html";
            throw new Error("Unauthorized");
        }

        const data = await response.json();
        console.log("Recieved datt", data); я

        if (data.userStatus === "user") {
            alert("You do not have access to this section.");
            window.location.href = "/Frontend/client/main_page.html";
        } else if (data.userStatus === "admin") {
            console.log("User is admin. He has permission");
        }
    } catch (error) {
        console.error("Error checking session", error);
    }
}


document.addEventListener("DOMContentLoaded", checkSession);


