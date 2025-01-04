let users = [];
let totalPages = 1;
let currentPage = 1;

async function fetchUsers(sortBy = '', emailFilter = '', page = 1) {
    const url = new URL('http://localhost:8080/getUsers');
    const params = new URLSearchParams();

    if (sortBy) {
        params.append('sort', sortBy);
    }
    if (emailFilter) {
        params.append('email', emailFilter);
    }
    params.append('page', page);
    params.append('pageSize', 9);

    url.search = params.toString();

    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Ошибка при загрузке данных');
        }

        const data = await response.json();
        console.log("Received data:", data);

        if (Array.isArray(data.users)) {
            users = data.users;
            totalPages = data.totalPages;
            currentPage = data.currentPage;
            renderTable(users);
            renderPagination();
        } else {
            console.error("Received data is not in the expected format:", data);
        }

    } catch (error) {
        console.error('Ошибка:', error);
    }
}

function renderTable(users) {
    const tableBody = document.getElementById('userTableBody');
    tableBody.innerHTML = ''; 

    users.forEach(user => {
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

function renderPagination() {
    const paginationContainer = document.getElementById('paginationContainer');
    paginationContainer.innerHTML = '';

    const prevButton = document.createElement('button');
    prevButton.innerText = 'Previous';
    prevButton.disabled = currentPage <= 1;
    prevButton.onclick = () => fetchUsers(document.getElementById('sortSelect').value, document.getElementById('emailFilter').value, currentPage - 1);
    paginationContainer.appendChild(prevButton);

    for (let i = 1; i <= totalPages; i++) {
        const pageButton = document.createElement('button');
        pageButton.innerText = i;
        pageButton.classList.toggle('active', i === currentPage);
        pageButton.onclick = () => fetchUsers(document.getElementById('sortSelect').value, document.getElementById('emailFilter').value, i);
        paginationContainer.appendChild(pageButton);
    }

    const nextButton = document.createElement('button');
    nextButton.innerText = 'Next';
    nextButton.disabled = currentPage >= totalPages;
    nextButton.onclick = () => fetchUsers(document.getElementById('sortSelect').value, document.getElementById('emailFilter').value, currentPage + 1);
    paginationContainer.appendChild(nextButton);
}

function applyFilters() {
    const sortValue = document.getElementById('sortSelect').value;
    const emailFilter = document.getElementById('emailFilter').value.toLowerCase();

    fetchUsers(sortValue, emailFilter, 1);  // Always start from page 1 when filters are applied
}

async function checkSession() {
    try {
        const response = await fetch('http://localhost:8080/checksession', {
            method: "GET",
            credentials: "include", // Обязательно включаем cookie
        });

        console.log("Response status:", response.status);

        if (!response.ok) {
            alert("The session is invalid. Redirection to the main page.");
            window.location.href = "/Frontend/client/main_page.html";
            throw new Error("Unauthorized");
        }

        const data = await response.json();
        console.log("Received data", data);

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

document.addEventListener("DOMContentLoaded", async () => {
    await checkSession(); 
    fetchUsers(); 
});


















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
