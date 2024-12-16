const foodData = {
    Kazakhstan: {
        dish: "Beshbarmak",
        image: "/pictures/beshbarmak.jpg",
        description: "A traditional Kazakh dish consisting of boiled horse meat or lamb with flat noodles and onion sauce. The name 'Beshbarmak' means 'five fingers' as it is traditionally eaten with hands."
    },
    France: {
        dish: "Croissant",
        image: "/pictures/croissant.jpg",
        description: "A small pastry, a crescent-shaped bun made of puff pastry. It's the popular product of French cuisine, served for breakfast with coffee for adults or with cocoa for children."
    },
    Italy: {
        dish: "Osso Buco",
        image: "/pictures/osso buco.jpg",
        description: "A Milanese specialty of cross-cut veal shanks braised with vegetables, white wine, and broth. It's traditionally served with gremolata and risotto alla Milanese."
    },
    China: {
        dish: "Peking Duck",
        image: "/pictures/peking duck.jpg",
        description: "A famous duck dish from Beijing that has been prepared since the imperial era. The duck is roasted until the skin is thin, crispy, and glossy brown, served with scallions, cucumber, and sweet bean sauce."
    },
    Australia: {
        dish: "Meat Pie",
        image: "/pictures/meat pie.jpg",
        description: "An iconic Australian savory pie containing minced meat and gravy, sometimes with onions, mushrooms, and cheese. It's a popular take-away food often enjoyed at sporting events."
    }
};

const pins = document.querySelectorAll('.pin');
const popup = document.querySelector('.food-popup');

function showPopup(event) {
    const pin = event.currentTarget;
    const country = pin.dataset.country;
    const data = foodData[country];
    
    // Update popup content
    popup.querySelector('img').src = data.image;
    popup.querySelector('img').alt = data.dish;
    popup.querySelector('h3').textContent = data.dish;
    popup.querySelector('p').textContent = data.description;

    // Position popup
    const pinRect = pin.getBoundingClientRect();
    const mapRect = document.querySelector('.world-map').getBoundingClientRect();
    
    let left = pinRect.left - mapRect.left;
    let top = pinRect.top - mapRect.top;

    // Adjust position to prevent popup from going off-screen
    if (left + 300 > mapRect.width) {
        left = left - 300 - 20;
    } else {
        left = left + 20;
    }

    if (top + 300 > mapRect.height) {
        top = top - 300;
    }

    popup.style.left = `${left}px`;
    popup.style.top = `${top}px`;
    
    // Show popup
    popup.classList.add('active');
}

function hidePopup() {
    popup.classList.remove('active');
}


pins.forEach(pin => {
    pin.addEventListener('mouseenter', showPopup);
    pin.addEventListener('mouseleave', hidePopup);
});
