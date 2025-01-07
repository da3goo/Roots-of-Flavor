document.getElementById("custom-aboutUsText").style.color = "#d13d64";
document.getElementById("custom-aboutUsText").style.textDecoration ="underline";
document.body.style.backgroundColor = "beige";
      
const darkModeToggle = document.getElementById('custom-darkModeBtn');
const body = document.body;
if (localStorage.getItem('dark-mode') === 'enabled') {
    body.classList.add('dark-mode');

}
// Toggle button handler
darkModeToggle.addEventListener('click', () => {
    body.classList.toggle('dark-mode');
    if (body.classList.contains('dark-mode')) {
        localStorage.setItem('dark-mode', 'enabled');
        document.body.style.backgroundColor = "#121212";
        document.getElementById('aboutUsText').style.color = '#e0e0e0';
        document.getElementById('aboutUsP1').style.color = '#e0e0e0';
        document.getElementById('aboutUsP2').style.color = '#e0e0e0';

    } else {
        localStorage.setItem('dark-mode', 'disabled');
        document.body.style.backgroundColor = "beige";
        document.getElementById('aboutUsText').style.color = '';
        document.getElementById('aboutUsP1').style.color = '#606060';
        document.getElementById('aboutUsP2').style.color = '#606060';


    }
});
document.getElementById('contactForm').addEventListener('submit', function (e) {
    e.preventDefault();

    var formData = new FormData(this);

    fetch('http://localhost:8080/send', {
      method: 'POST',
      body: formData
    })
    .then(response => response.text())
    .then(data => {
      document.getElementById('responseMessage').innerHTML = `<p style="color: green;">${data}</p>`;
      
      alert('The message was sent successfully!');

      document.getElementById('contactForm').reset();
    })
    .catch(error => {
      document.getElementById('responseMessage').innerHTML = `<p style="color: red;">Ошибка при отправке: ${error.message}</p>`;
      
      alert('Error in sending ' + error.message);
    });
});
      

        

        