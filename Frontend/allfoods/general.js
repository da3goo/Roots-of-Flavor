document.addEventListener('DOMContentLoaded', function() {
 
    const urlParams = new URLSearchParams(window.location.search);
    const foodName = urlParams.get('name');  // Получаем название еды из параметра URL

    
    if (!foodName) {
        console.error('Food name is missing in the URL');
        return;
    }


    const apiUrl = `http://localhost:8080/food?name=${foodName}`;


   
    function displayFoodData(data) {
        
        const foodNameElement = document.querySelector('#foodName');
        const originsContent = document.querySelector('#originsContent');
        const culturalContent = document.querySelector('#culturalContent');
        const healthContent = document.querySelector('#healthContent');
        const aroundWorldContent = document.querySelector('#aroundWorldContent');
        const foodImage = document.querySelector('#foodImage');
        const likeButton = document.getElementById('like-button');
        const likeIcon = document.getElementById('like-icon');
        const likeCount = document.getElementById('like-count');

        // Заполняем данные на странице
        foodNameElement.innerHTML = data.name;
        originsContent.innerHTML = data.description1;
        culturalContent.innerHTML = data.description2;
        healthContent.innerHTML = data.description3;
        aroundWorldContent.innerHTML = data.description4;
        foodImage.src = data.image_url;
        foodImage.alt = data.name;

        // Лайк кнопка
        let isLiked = false;
        let currentLikes = parseInt(likeCount.textContent);

        likeButton.addEventListener('click', () => {
            isLiked = !isLiked;

            if (isLiked) {
                likeIcon.classList.replace('far', 'fas');
            } else {
                likeIcon.classList.replace('fas', 'far');
            }
        });
    }

    //запросик
    fetch(apiUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error('Food not found or server error');
            }
            return response.json(); 
        })
        .then(data => {
            displayFoodData(data); // пон
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to fetch food data: ' + error.message);
        });
        

        

    // Обработчик для сабменю
    const foodsLink = document.getElementById('foodsLink');
    const submenu = document.getElementById('foodsSubmenu');

    if (foodsLink && submenu) {
        foodsLink.addEventListener('click', function(event) {
            event.preventDefault();

            const menuItem = this.parentElement;

            
            if (!menuItem.classList.contains('active')) {
                menuItem.classList.add('active');
            } else {
               
                menuItem.classList.remove('active');
            }
        });
    } else {
        console.error('Foods link or submenu not found');
    }

    // Темный режим
    const darkModeToggle = document.getElementById('dark-mode-toggle');
    const body = document.body;

    if (darkModeToggle) {
        if (localStorage.getItem('dark-mode') === 'enabled') {
            body.classList.add('dark-mode');
        }

        // Обработчик кнопки переключения
        darkModeToggle.addEventListener('click', () => {
            body.classList.toggle('dark-mode');
            if (body.classList.contains('dark-mode')) {
                localStorage.setItem('dark-mode', 'enabled');
            } else {
                localStorage.setItem('dark-mode', 'disabled');
            }
        });
    } else {
        console.error('Dark mode toggle button not found');
    }
});

document.getElementById('like-button').addEventListener('click', function () {
    var foodName = document.getElementById('foodName').textContent.trim();

    // Sending getting foodname by parameters
    fetch(`http://localhost:8080/addFoodToSaved?name=${encodeURIComponent(foodName)}`, {
        method: 'GET',
    })
    .then(response => response.text())
    .then(data => {
        alert(data);  
    })
    .catch(error => {
        alert('Error: ' + error.message);
    });
});
 