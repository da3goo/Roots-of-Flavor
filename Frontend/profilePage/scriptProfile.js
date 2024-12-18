const editBtn = document.getElementById('edit-btn');
const editForm = document.getElementById('edit-form');
const profileName = document.getElementById('profile-name');
const profileInfo = document.getElementById('profile-info');
const editName = document.getElementById('edit-name');
const editAbout = document.getElementById('edit-about');

// Show form on edit button click
editBtn.addEventListener('click', () => {
    editForm.classList.toggle('hidden');
    profileInfo.classList.add('hidden');
    editName.value = profileName.innerText;
    editAbout.value = profileInfo.innerText || '';
});

// Handle form submission
editForm.addEventListener('submit', async (e) => {
e.preventDefault();
profileName.innerText = editName.value;
profileInfo.innerText = editAbout.value;
profileInfo.classList.remove('hidden');
editForm.classList.add('hidden');

//Send POST request
try {
    const response = await fetch('https://your-server-endpoint.com/update-profile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: editName.value,
            about: editAbout.value
        })
    });
    if (response.ok) {
        alert('Profile updated successfully!');
    } else {
        alert('Error updating profile');
    }
} catch (error) {
    console.error(error);
    alert('Failed to update profile. Please try again.');
}
});