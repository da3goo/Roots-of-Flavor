const questions = [
    { 
      "question": "Where did tacos originate?", 
      "answers": ["Mexico", "USA", "Spain", "Colombia"], 
      "correct": 0 
    },
    { 
      "question": "What was the original meaning of the word 'taco' in the 18th century?", 
      "answers": ["A kind of fish", "A type of tortilla", "Gunpowder wrapped in paper", "A miner's tool"], 
      "correct": 2 
    },
    { 
      "question": "What is the key ingredient of a Mojito?", 
      "answers": ["Gin", "Rum", "Vodka", "Tequila"], 
      "correct": 1 
    },
    { 
      "question": "Which region is known for the production of Farine de Blé Noir de Bretagne?", 
      "answers": ["Provence", "Brittany", "Normandy", "Paris"], 
      "correct": 1 
    },
    { 
      "question": "What is the traditional meat used in Kabanosy?", 
      "answers": ["Pork", "Beef", "Lamb", "Chicken"], 
      "correct": 0 
    },
    { 
      "question": "What kind of milk is used to produce Kumis?", 
      "answers": ["Goat's milk", "Cow's milk", "Mare's milk", "Sheep's milk"], 
      "correct": 2 
    },
    { 
      "question": "Which of the following is a traditional filling in a meat pie from Australia and New Zealand?", 
      "answers": ["Fish", "Mushrooms", "Crocodile", "Cheese"], 
      "correct": 2 
    },
    { 
      "question": "What is the main ingredient of the dish Osso Bucco?", 
      "answers": ["Chicken", "Veal shank", "Pork belly", "Beef steak"], 
      "correct": 1 
    },
    { 
      "question": "Which drink is known as Japan's most famous alcoholic beverage?", 
      "answers": ["Sake", "Shochu", "Umeshu", "Sake-tini"], 
      "correct": 0 
    },
    { 
      "question": "What ingredient makes the intense red color in kimchi?", 
      "answers": ["Tomatoes", "Chili peppers", "Beetroot", "Paprika"], 
      "correct": 1 
    },
    { 
      "question": "Which country is the origin of Vodka?", 
      "answers": ["Poland", "Russia", "Finland", "Germany"], 
      "correct": 1 
    },
    { 
      "question": "What type of cheese is traditionally used in Pizza Margherita?", 
      "answers": ["Mozzarella", "Cheddar", "Parmesan", "Gorgonzola"], 
      "correct": 0 
    },
    { 
      "question": "What is the name of the famous French bread known for its crispy crust?", 
      "answers": ["Brioche", "Baguette", "Focaccia", "Pain de Campagne"], 
      "correct": 1 
    },
    { 
      "question": "Where did the Croissant originate?", 
      "answers": ["France", "Germany", "Austria", "Italy"], 
      "correct": 2 
    },
    { 
      "question": "Which of the following is a key component of a Sake brew?", 
      "answers": ["Rice", "Corn", "Potatoes", "Barley"], 
      "correct": 0 
    },
    { 
      "question": "Which type of meat is commonly used in Beshbarmak?", 
      "answers": ["Lamb", "Beef", "Chicken", "Horse"], 
      "correct": 3 
    },
    { 
      "question": "Which country is famous for its kimchi?", 
      "answers": ["Japan", "China", "Korea", "Thailand"], 
      "correct": 2 
    },
    { 
      "question": "What is the key flavor ingredient in Peking Duck?", 
      "answers": ["Ginger", "Soy sauce", "Orange peel", "Hoisin sauce"], 
      "correct": 3 
    },
    { 
      "question": "What vegetable is traditionally used in the base of a traditional Osso Bucco?", 
      "answers": ["Tomatoes", "Carrots", "Leeks", "Celery"], 
      "correct": 1 
    },
    { 
      "question": "What type of alcohol is used to prepare a Mojito?", 
      "answers": ["Vodka", "Rum", "Whiskey", "Gin"], 
      "correct": 1 
    },
    { 
      "question": "Which of the following is a characteristic of the traditional Russian Caravan tea?", 
      "answers": ["Floral", "Citrusy", "Smoky", "Minty"], 
      "correct": 2 
    },
    { 
      "question": "What is the traditional grain used to make Kumpir?", 
      "answers": ["Rice", "Barley", "Wheat", "Potato"], 
      "correct": 3 
    },
    { 
      "question": "What is the primary flavoring of Lapsang Souchong tea?", 
      "answers": ["Floral", "Smoky", "Fruity", "Nutty"], 
      "correct": 1 
    },
    { 
      "question": "What was the original ingredient of the French croissant?", 
      "answers": ["Brioche dough", "Puff pastry", "Doughnut dough", "Cake batter"], 
      "correct": 1 
    },
    { 
      "question": "In which European country is Osso Bucco considered a signature dish?", 
      "answers": ["Spain", "Italy", "France", "Germany"], 
      "correct": 1 
    },
    { 
      "question": "Which type of fish is commonly used in sushi?", 
      "answers": ["Salmon", "Tuna", "Mackerel", "All of the above"], 
      "correct": 3 
    },
    { 
      "question": "What does the word 'kabanosy' refer to in Polish cuisine?", 
      "answers": ["A type of bread", "A long sausage", "A type of cheese", "A soup"], 
      "correct": 1 
    },
    { 
      "question": "Which country’s food culture includes the famous dish Beshbarmak?", 
      "answers": ["Kyrgyzstan", "Kazakhstan", "Uzbekistan", "Tajikistan"], 
      "correct": 1 
    },
    { 
      "question": "What is the main characteristic of a true Peking Duck dish?", 
      "answers": ["Grilled meat", "Crispy skin", "Roasted in the ground", "Served with flatbreads"], 
      "correct": 1 
    },
    { 
      "question": "What ingredient is crucial to making Gyokuro tea?", 
      "answers": ["Matcha powder", "Unprocessed green tea leaves", "Shaded leaves", "Brewed with milk"], 
      "correct": 2 
    },
    { 
      "question": "What is the major ingredient in the preparation of Onigiri?", 
      "answers": ["Rice", "Seaweed", "Fish", "Soybeans"], 
      "correct": 0 
    },
    { 
      "question": "What is the key component of the traditional Margarita pizza?", 
      "answers": ["Cheddar", "Mozzarella", "Ricotta", "Feta"], 
      "correct": 1 
    },
    { 
      "question": "What is a key feature of the dough used in croissants?", 
      "answers": ["Yeast-leavened", "Baking powder-based", "Puff pastry", "Brioche dough"], 
      "correct": 2 
    },
    { 
      "question": "What does the word 'Nori' refer to?", 
      "answers": ["A type of fish", "A seaweed", "A tea", "A rice preparation"], 
      "correct": 1 
    },
    { 
      "question": "What makes the crust of a Baguette crispy?", 
      "answers": ["Baked at a high temperature", "Using egg wash", "Baked in a clay oven", "Sprinkling water on dough"], 
      "correct": 0 
    },
    { 
      "question": "What kind of meat is typically used for preparing Shawarma?", 
      "answers": ["Pork", "Lamb", "Beef", "All of the above"], 
      "correct": 3 
    },
    { 
      "question": "Where did the sandwich called 'Osso Bucco' originate?", 
      "answers": ["USA", "Germany", "France", "Italy"], 
      "correct": 3 
    },
    { 
      "question": "What color is the traditional Cheddar cheese?", 
      "answers": ["White", "Yellow", "Orange", "Red"], 
      "correct": 1 
    },
    { 
        "question": "Which region is famous for producing Kumpir?", 
        "answers": ["Turkey", "Greece", "Spain", "Italy"], 
        "correct": 0 
      },
      { 
        "question": "What is the origin of the name 'Pepperoni'?", 
        "answers": ["Italian", "French", "American", "Spanish"], 
        "correct": 0 
      },
      { 
        "question": "What is the characteristic color of Cheddar cheese?", 
        "answers": ["Yellow", "White", "Red", "Blue"], 
        "correct": 0 
      },
      { 
        "question": "What key ingredient distinguishes the taste of Gyokuro tea?", 
        "answers": ["Shaded leaves", "High heat", "Matcha powder", "Honey"], 
        "correct": 0 
      },
      { 
        "question": "What is the main ingredient in the dish Osso Bucco?", 
        "answers": ["Beef", "Veal", "Chicken", "Pork"], 
        "correct": 1 
      },
      { 
        "question": "Which famous tea blend uses a smoky flavor derived from dried leaves?", 
        "answers": ["Earl Grey", "Russian Caravan", "Green tea", "Oolong"], 
        "correct": 1 
      },
      { 
        "question": "What is the origin of Bagguette?", 
        "answers": ["Italy", "France", "Germany", "Spain"], 
        "correct": 1 
      },
      { 
        "question": "What is the primary seasoning in traditional Kimchi?", 
        "answers": ["Ginger", "Garlic", "Chili", "Salt"], 
        "correct": 2 
      },
      { 
        "question": "Which of these dishes is from Italy?", 
        "answers": ["Banh Mi", "Sushi", "Osso Bucco", "Tacos"], 
        "correct": 2 
      },
      { 
        "question": "What type of alcohol is used in a Mojito?", 
        "answers": ["Rum", "Tequila", "Vodka", "Whiskey"], 
        "correct": 0 
      },
      { 
        "question": "Where did the concept of Sushi originate?", 
        "answers": ["China", "Japan", "Korea", "Vietnam"], 
        "correct": 1 
      },
      { 
        "question": "What ingredient gives pepperoni its characteristic red color?", 
        "answers": ["Paprika", "Tomatoes", "Chili", "Cayenne pepper"], 
        "correct": 0 
      },
      { 
        "question": "What ingredient is added to make the crust of a Baguette crispy?", 
        "answers": ["Cornstarch", "Flour", "Water", "Egg wash"], 
        "correct": 2 
      },
      { 
        "question": "What is the term 'Sushi' actually referring to?", 
        "answers": ["The fish", "The rice", "The wrapping", "The seaweed"], 
        "correct": 1 
      },
      { 
        "question": "Which of the following is a traditional topping for a pizza Margherita?", 
        "answers": ["Basil", "Onion", "Olives", "Anchovies"], 
        "correct": 0 
      },
      { 
        "question": "Which country invented Croissant?", 
        "answers": ["Italy", "Austria", "France", "Germany"], 
        "correct": 1 
      },
      { 
        "question": "Which ingredient is the base for the dish 'Sushi'?", 
        "answers": ["Rice", "Seaweed", "Fish", "Soy sauce"], 
        "correct": 0 
      },
      { 
        "question": "What is the primary ingredient in a traditional Russian Caravan tea blend?", 
        "answers": ["Oolong", "Green tea", "Black tea", "Herbs"], 
        "correct": 2 
      },
      { 
        "question": "What traditional food is served as part of a French picnic?", 
        "answers": ["Baguette", "Sushi", "Tacos", "Kumpir"], 
        "correct": 0 
      },
      { 
        "question": "Which type of meat is traditionally used for Shawarma?", 
        "answers": ["Lamb", "Beef", "Chicken", "All of the above"], 
        "correct": 3 
      },
      { 
        "question": "What is a traditional accompaniment for Kumpir?", 
        "answers": ["Cheese", "Salad", "Mayonnaise", "All of the above"], 
        "correct": 3 
      },
      { 
        "question": "What type of cheese is used in Osso Bucco?", 
        "answers": ["Parmesan", "Mozzarella", "Ricotta", "Cheddar"], 
        "correct": 0 
      },
      { 
        "question": "Which country is most famous for producing Vodka?", 
        "answers": ["Poland", "Russia", "Sweden", "Germany"], 
        "correct": 1 
      },
      { 
        "question": "What is the primary ingredient used in making Kimchi?", 
        "answers": ["Cabbage", "Spinach", "Carrots", "Onions"], 
        "correct": 0 
      },
      { 
        "question": "What region is famous for the production of Farine de Blé Noir?", 
        "answers": ["Provence", "Brittany", "Normandy", "Paris"], 
        "correct": 1 
      },
      { 
        "question": "What is the key ingredient used to create the flavor of 'Tacos de Minero'?", 
        "answers": ["Chili", "Tortilla", "Meat", "Gunpowder"], 
        "correct": 3 
      },
      { 
        "question": "Which of these is a famous French bread?", 
        "answers": ["Brioche", "Baguette", "Focaccia", "Ciabatta"], 
        "correct": 1 
      },
      { 
        "question": "Which of the following is commonly used in Osso Bucco?", 
        "answers": ["Rice", "Pasta", "Polenta", "Bread"], 
        "correct": 2 
      },
      { 
        "question": "What is the basic ingredient in traditional Nori?", 
        "answers": ["Seaweed", "Rice", "Fish", "Vegetables"], 
        "correct": 0 
      },
      { 
        "question": "What is the name of the dried Japanese seaweed used in sushi?", 
        "answers": ["Kombu", "Nori", "Wakame", "Hijiki"], 
        "correct": 1 
      },
      { 
        "question": "What is the main protein used in a traditional Gyokuro tea ceremony?", 
        "answers": ["Matcha", "Tuna", "Rice", "Tea leaves"], 
        "correct": 3 
      },
      { 
        "question": "What does the word 'Mojito' translate to?", 
        "answers": ["Little Miracle", "Little Death", "Little Kiss", "Little Mint"], 
        "correct": 3 
      },
      { 
        "question": "What is the primary ingredient used in the preparation of pizza margherita?", 
        "answers": ["Basil", "Tomato", "Cheese", "Olives"], 
        "correct": 1 
      },
      { 
        "question": "Where did the origin of Sushi come from?", 
        "answers": ["China", "Japan", "Vietnam", "Korea"], 
        "correct": 1 
      },
      { 
        "question": "What is the main characteristic of a traditional croissant?", 
        "answers": ["Flaky", "Dense", "Chewy", "Spongy"], 
        "correct": 0 
      },
      { 
        "question": "Which country is famous for the dish Osso Bucco?", 
        "answers": ["Italy", "Spain", "France", "Greece"], 
        "correct": 0 
      },
      { 
        "question": "What is the origin of the food item Tacos?", 
        "answers": ["USA", "Mexico", "Italy", "Brazil"], 
        "correct": 1 
      },
      { 
        "question": "What is the key flavoring used in traditional Osso Bucco?", 
        "answers": ["Garlic", "Cilantro", "Rosemary", "Lemon zest"], 
        "correct": 3 
      },
      { 
        "question": "Which of these cheeses is traditionally used in pizza margherita?", 
        "answers": ["Mozzarella", "Cheddar", "Feta", "Brie"], 
        "correct": 0 
      },
      { 
        "question": "What meat is traditionally used in Beshbarmak?", 
        "answers": ["Lamb", "Pork", "Chicken", "Beef"], 
        "correct": 0 
      },
      { 
        "question": "Where did the dish Kumpir originate?", 
        "answers": ["Turkey", "Egypt", "Italy", "Greece"], 
        "correct": 0 
      },
      { 
        "question": "Which of the following is commonly served with pizza?", 
        "answers": ["French fries", "Garlic bread", "Salad", "Chips"], 
        "correct": 1 
      },
      { 
        "question": "What is a traditional base of the 'pizza margherita'?", 
        "answers": ["Pesto", "Tomato sauce", "Garlic butter", "Cheese sauce"], 
        "correct": 1 
      },
      { 
        "question": "What grain is primarily used in the production of sake?", 
        "answers": ["Barley", "Rice", "Wheat", "Oats"], 
        "correct": 1 
      },
      { 
        "question": "What is the key characteristic of Kabanosy sausages?", 
        "answers": ["Dry", "Spicy", "Fried", "Grilled"], 
        "correct": 0 
      },
      { 
        "question": "What is the country of origin for the sandwich 'Banh Mi'?", 
        "answers": ["Vietnam", "Thailand", "China", "Japan"], 
        "correct": 0 
      },
      { 
        "question": "Which country produces the most Vodka?", 
        "answers": ["Russia", "Poland", "USA", "Finland"], 
        "correct": 0 
      },
      { 
        "question": "Which tea variety is known for its smoke flavor?", 
        "answers": ["Darjeeling", "Lapsang Souchong", "Sencha", "Gyokuro"], 
        "correct": 1 
      },
      { 
        "question": "What is the traditional flavor of Osso Bucco?", 
        "answers": ["Sour", "Salty", "Bitter", "Savory"], 
        "correct": 3 
      },
      { 
        "question": "Which type of tea is a key component of the Russian Caravan blend?", 
        "answers": ["Black tea", "Green tea", "Oolong tea", "Herbal tea"], 
        "correct": 0 
      },
      { 
        "question": "What type of tea is traditionally used to make a matcha latte?", 
        "answers": ["Green tea", "Black tea", "Oolong tea", "White tea"], 
        "correct": 0 
      }
  ];

  let currentQuestionIndex = 0;
  let score = 0;
  let usedQuestions = [];
  let selectedButton = null;
  let timer;
  let timeLeft = 180; 
  
  // Timer update function
  function updateTimer() {
      const minutes = Math.floor(timeLeft / 60);
      const seconds = timeLeft % 60;
      document.getElementById("timer").innerText = `Time: ${minutes < 10 ? '0' + minutes : minutes}:${seconds < 10 ? '0' + seconds : seconds}`;
      
      if (timeLeft <= 0) {
          clearInterval(timer); 
          showResults(); 
      } else {
          timeLeft--;
      }
  }
  
  function loadQuestion() {
      if (usedQuestions.length === 10 || timeLeft <= 0) {
          showResults(); 
          return;
      }
  
      let randomIndex;
      do {
          randomIndex = Math.floor(Math.random() * questions.length);
      } while (usedQuestions.includes(randomIndex));
  
      usedQuestions.push(randomIndex);
      const question = questions[randomIndex];
      currentQuestionIndex = randomIndex;
  
      
      document.getElementById("question-counter").innerText = `Question ${usedQuestions.length} of 10`;
  
      document.getElementById("question").innerText = question.question;
      const answersDiv = document.getElementById("answers");
      answersDiv.innerHTML = "";
  
      selectedButton = null;
  
      question.answers.forEach((answer, index) => {
          const btn = document.createElement("button");
          btn.innerText = answer;
          btn.className = "answer-btn";
          btn.onclick = () => selectAnswer(btn, index);
          answersDiv.appendChild(btn);
      });
  }
  
  function selectAnswer(button, index) {
      if (selectedButton) {
          selectedButton.classList.remove("selected");
      }
      button.classList.add("selected");
      selectedButton = button;
  
      if (!button.dataset.checked) {
          button.dataset.checked = "true";
          if (index === questions[currentQuestionIndex].correct) {
              score++;
          }
      }
      document.getElementById("next-btn").style.display = "inline-block"; 
  }
  
  function nextQuestion() {
      if (!selectedButton) {
          
          return;
      }
      
      document.getElementById("next-btn").style.display = "none"; 
      loadQuestion();
  }
  
  function showResults() {
      clearInterval(timer); 
      const percentage = Math.round((score / 10) * 100); 
      document.getElementById("result").innerHTML = `Quiz Complete! Your score: ${percentage}%`;
      document.getElementById("quiz-card").style.display = "none";
      document.getElementById("next-btn").style.display = "none"; 
      document.getElementById("play-again-btn").style.display = "inline-block"; 
      document.getElementById("back-btn").style.display = "inline-block"; 
  }
  
  function playAgain() {
      
      window.location.reload();
  }
  
  function goBack() {
      
      if (usedQuestions.length > 1) {
          usedQuestions.pop(); 
          score--; 
          currentQuestionIndex = usedQuestions[usedQuestions.length - 1];
          loadQuestion();
      }
  }
  
  // Handler for the "Enter" key
  document.addEventListener("keydown", function(event) {
      if (event.key === "Enter" && selectedButton) {
          nextQuestion(); 
      }
  });
  

  timer = setInterval(updateTimer, 1000);
  
  loadQuestion();