<template>
  <div class="app">
    <!-- Sidebar -->
    <Sidebar v-if="showSidebar" />

    <!-- Content -->
    <router-view />
  </div>


  <div class="center-form">
    <div class="container" :class="{ active: isSignup }">
      <div class="forms">
        <div class="form login">
          <span class="title">Onviz</span>
          <div class="input-field-country">
            <select v-model="selectedCountry" name="country" id="country-select">
              <option value="" disabled selected>Выберите страну</option>
              <option value="Russia">Россия</option>
              <option value="Australia">Австралия</option>
              <option value="Austria">Австрия</option>
              <option value="Azerbaijan">Азербайджан</option>
              <option value="Anguilla">Ангилья</option>
              <option value="Argentina">Аргентина</option>
              <option value="Armenia">Армения</option>
              <option value="Arulko">Арулько</option>
              <option value="Belarus">Беларусь</option>
              <option value="Belize">Белиз</option>
              <option value="Belgium">Бельгия</option>
              <option value="Bermuda">Бермуды</option>
              <option value="Bulgaria">Болгария</option>
              <option value="Brazil">Бразилия</option>
              <option value="UK">Великобритания</option>
              <option value="Hungary">Венгрия</option>
              <option value="Vietnam">Вьетнам</option>
              <option value="Haiti">Гаити</option>
              <option value="Guadeloupe">Гваделупа</option>
              <option value="Germany">Германия</option>
              <option value="Netherlands">Голландия</option>
              <option value="Honduras">Гондурас</option>
              <option value="HongKong">Гонконг</option>
              <option value="Greece">Греция</option>
              <option value="Georgia">Грузия</option>
              <option value="Denmark">Дания</option>
              <option value="DominicanRepublic">Доминиканская республика</option>
              <option value="Egypt">Египет</option>
              <option value="Israel">Израиль</option>
              <option value="India">Индия</option>
              <option value="Indonesia">Индонезия</option>
              <option value="Jordan">Иордания</option>
              <option value="Iraq">Ирак</option>
              <option value="Iran">Иран</option>
              <option value="Ireland">Ирландия</option>
              <option value="Spain">Испания</option>
              <option value="Italy">Италия</option>
              <option value="Kazakhstan">Казахстан</option>
              <option value="Cameroon">Камерун</option>
              <option value="Canada">Канада</option>
              <option value="Caribbean">Карибы</option>
              <option value="Cyprus">Кипр</option>
              <option value="Kyrgyzstan">Киргизстан</option>
              <option value="China">Китай</option>
              <option value="Korea">Корея</option>
              <option value="CostaRica">Коста-Рика</option>
              <option value="Cuba">Куба</option>
              <option value="Kuwait">Кувейт</option>
              <option value="Latvia">Латвия</option>
              <option value="Lebanon">Ливан</option>
              <option value="Libya">Ливан</option>
              <option value="Lithuania">Литва</option>
              <option value="Luxembourg">Люксембург</option>
              <option value="Macedonia">Македония</option>
              <option value="Malaysia">Малайзия</option>
              <option value="Malta">Мальта</option>
              <option value="Mexico">Мексика</option>
              <option value="Mozambique">Мозамбик</option>
              <option value="Moldova">Молдова</option>
              <option value="Monaco">Монако</option>
              <option value="Mongolia">Монголия</option>
              <option value="Morocco">Морокко</option>
              <option value="Netherlands">Нидерланды</option>
              <option value="NewZealand">Новая Зеландия</option>
              <option value="Norway">Норвегия</option>
              <option value="UAE">О.А.Э.</option>
              <option value="IsleOfMan">Остров Мэн</option>
              <option value="Pakistan">Пакистан</option>
              <option value="Peru">Перу</option>
              <option value="Poland">Польша</option>
              <option value="Portugal">Португалия</option>
              <option value="Reunion">Реюньон</option>
              <option value="Romania">Румыния</option>
              <option value="USA">США</option>
              <option value="ElSalvador">Сальвадор</option>
              <option value="Singapore">Сингапур</option>
              <option value="Syria">Сирия</option>
              <option value="Slovakia">Словакия</option>
              <option value="Slovenia">Словения</option>
              <option value="Suriname">Суринам</option>
              <option value="Tajikistan">Таджикистан</option>
              <option value="Taiwan">Тайвань</option>
              <option value="Thailand">Тайланд</option>
              <option value="Tunisia">Тунис</option>
              <option value="Turkmenistan">Туркменистан</option>
              <option value="Turkmenistan">Туркмения</option>
              <option value="TurksAndCaicos">Туркс и Кейкос</option>
              <option value="Turkey">Турция</option>
              <option value="Uganda">Уганда</option>
              <option value="Uzbekistan">Узбекистан</option>
              <option value="Ukraine">Украина</option>
              <option value="Finland">Финляндия</option>
              <option value="France">Франция</option>
              <option value="Croatia">Хорватия</option>
              <option value="CzechRepublic">Чехия</option>
              <option value="Chile">Чили</option>
              <option value="Switzerland">Швейцария</option>
              <option value="Sweden">Швеция</option>
              <option value="Ecuador">Эквадор</option>
              <option value="Estonia">Эстония</option>
              <option value="SouthAfrica">ЮАР</option>
              <option value="Yugoslavia">Югославия</option>
              <option value="SouthKorea">Южная Корея</option>
              <option value="Jamaica">Ямайка</option>
              <option value="Japan">Япония</option>
            </select>
          </div>

          <div class="input-field">
            <input v-model="loginEmail" type="text" placeholder="Enter your email" required>
            <i class="uil uil-envelope icon"></i>
          </div>
          <div class="input-field">
            <input v-model="loginPassword" type="password" class="password" placeholder="Enter your password" required>
            <i class="uil uil-lock icon"></i>
            <i class="uil uil-eye-slash showHidePw" @click="togglePasswordVisibility('loginPassword')"></i>
          </div>


          <div class="input-field button">
            <input type="button" value="Login" @click="signinUser">
          </div>


        </div>

        <!-- Registration Form -->
        <div class="form signup">
          <span class="title">Registration</span>


          <div class="input-field" name="username" id="username">
            <input v-model="signupName" type="text" placeholder="Enter your name" required>
            <i class="uil uil-user"></i>
          </div>
          <div class="input-field" name="email" id="email">
            <input v-model="signupEmail" type="text" placeholder="Enter your email" required>
            <i class="uil uil-envelope icon"></i>
          </div>
          <div class="input-field" name="password" id="password">
            <input v-model="signupPassword" type="password" class="password" placeholder="Create a password" required>
            <i class="uil uil-lock icon"></i>
          </div>
          <div class="input-field" name="passwordConfirm" id="passwordConfirm">
            <input v-model="signupPasswordConfirm" type="password" class="password" placeholder="Confirm a password" required>
            <i class="uil uil-lock icon"></i>
            <i class="uil uil-eye-slash showHidePw" @click="togglePasswordVisibility('signupPasswordConfirm')"></i>
          </div>

          <div class="checkbox-text">
            <div class="checkbox-content">
              <input type="checkbox" id="termCon">
              <label for="termCon" class="text">I accepted all terms and conditions</label>
            </div>
          </div>

          <div class="input-field button">
            <input type="button" value="Signup" @click="signupUser">
          </div>


          <div class="login-signup">
          <span class="text">Already a member?
            <a href="#" class="text login-link" @click="switchToLogin">Login Now</a>
          </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>


import Sidebar from '../components/Sidebar.vue';
import { ref } from 'vue';
import { onMounted } from 'vue';

const showSidebar = ref(false); // Initially hide the Sidebar

// Check if the user is signed in when the component is mounted
onMounted(() => {
  // Retrieve the login status from sessionStorage
  const isSignedIn = sessionStorage.getItem('isSignedIn');
  if (isSignedIn === 'true') {
    showSidebar.value = false;
  }
});

const selectedCountry = ref('');
const isSignup = ref(false);
const loginEmail = ref('');
const loginPassword = ref('');
const signupName = ref('');
const signupEmail = ref('');
const signupPassword = ref('');
const signupPasswordConfirm = ref('');
const tokenEndpoint = 'https://onviz-api.ru/api/yandex/token';
// Replace these with your actual client credentials and user credentials
const clientId = 'your-client-id';
const clientSecret = 'your-client-secret';
const username = 'user-username';
const password = 'user-password';

// Prepare the token request
const formData = new URLSearchParams();
formData.append('grant_type', 'password');
formData.append('client_id', clientId);
formData.append('client_secret', clientSecret);
formData.append('username', username);
formData.append('password', password);
fetch(tokenEndpoint, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  },
  body: formData,
})
  .then(response => response.json())
  .then(data => {
    console.log('Token response:', data);
    // Handle the token response as needed
  })
  .catch(error => {
    console.error('Token request failed:', error);
  });

function togglePasswordVisibility(passwordField) {
  const field = eval(passwordField);
  field.type = field.type === 'password' ? 'text' : 'password';
}

function validateEmail(email) {
  // Basic email format validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

function validatePassword(password) {
  // Password should contain at least 6 characters, a symbol, and both upper and lower case letters
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&(){}])[A-Za-z\d!@#$%^&(){}]{6,}$/;
  return passwordRegex.test(password);
}


function signinUser() {
  // Create an object with the user's login data
  if (!selectedCountry.value) {
    alert('Please select a country');
    return;
  }

  const loginData = {
    email: loginEmail.value,
    password: loginPassword.value,
    country: selectedCountry.value, // Include selected country in login data
    clientId: clientId, // Add clientId to login data
    clientSecret: clientSecret, // Add clientSecret to login data
  };

  // Define the URL of your server where you want to send the login data
  const loginUrl = 'https://onviz-api.ru/api/login_page'; // Replace with your server URL

  // Send a POST request to the server
  fetch(loginUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      //'Authorization': 'Bearer'
    },
    body: JSON.stringify(loginData),
  })
    .then(response => {
      if (response.status === 200) {
        console.log(response.status)
        // Login successful, get the token
        //showSidebar.value = true;

         //window.location.href = "https://oauth.yandex.ru/authorize?client_id=4fed8408c435482b950afeb2d6e0f3cc&client_secret=dbb4420ab51f41fc86a2dedd37d2302b"

        window.location.href = "https://onviz-api.ru/api/token_oauth"
        //window.location.href = "https://onviz-api.ru/api/yandex/authorize"
        //window.location.href = "https://social.yandex.net/broker/redirect"

       // window.location.href = "https://onviz-api.ru/api/redirect"
        return response.text(); // Assuming the response is a SHA-256 token string

      } else if (response.status === 400) {
        //  return response.text();
      } else if (response.status === 401) {
        return response.text();
      }
      else {
        // Login failed, log response status and response data
        alert('console.error')
        console.error('Login failed. Status:', response.status);
        return response.text().then(text => {
          console.error('Response data:', text);
          throw new Error('Login failed');
        });
      }
    })
    .then(token => {
      // Handle the response token
      console.log('Received token:', token);
      // Store the token in session storage
      sessionStorage.setItem('token', token);

      // Set showSidebar to true to display the Sidebar

      // Redirect to the desired URL
      //window.location.href = 'http://localhost:5173/leads'; // Replace with the URL you want to redirect to
    })
    .catch(error => {
      // Handle errors, e.g., registration failure
      console.error('Error:', error);
      alert('Login failed...');
    });
}

function signupUser() {
  // Client-side validation
  if (signupPassword.value !== signupPasswordConfirm.value) {
    alert('Password confirmation failed.');
    return;
  }

  if (
    signupName.value.length < 6 ||
    !validateEmail(signupEmail.value) ||
    !validatePassword(signupPassword.value)
  ) {
    alert('Invalid input. Please check your details.');
    return;
  }

  // Create an object with the user's registration data
  const userData = {
    username: signupName.value,
    email: signupEmail.value,
    password: signupPassword.value,
  };

  // Define the URL of your server where you want to send the registration data
  const registrationUrl = 'https://onviz-api.ru/api/auth_page'; // Replace with your server URL

  // Send a POST request to the server
  fetch(registrationUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  })
    .then(response => {
      if (response.status === 200) {
        // Registration successful
        showSidebar.value = false;

        return response.json();
      } else {
        // Log the error response for debugging
        console.error('Registration failed. Status:', response.status);
        return response.text().then(text => {
          console.error('Response data:', text);
          throw new Error('Registration failed');
        });
      }
    })
    .then(data => {
      // Handle the response data (if applicable)
      console.log('Server response:', data);
      if (data.message === 'Registration successful') {
        // Handle a successful registration
        alert('Registration successful!');
      } else {
        // Handle other responses from the server if needed
        if (data.message === 'User exist') {
          alert('User exists, please change the email address');
        } else {
          alert('Registration failed...');
        }
      }
    })
    .catch(error => {
      // Handle errors, e.g., registration failure
      console.error('Error:', error);
      alert('Registration failed...');
    });
}


function switchToSignup() {
  isSignup.value = true;
}

function switchToLogin() {
  isSignup.value = false;
}


</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap');

*{
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Poppins', sans-serif;
}

.input-field-country select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  outline: none;
  font-size: 14px;
  transition: border-color 0.2s ease;
  margin-top: 30px;
}

.input-field-country select:hover, .input-field select:focus {
  border-color: #4070f4;
}

body{
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #4070f4;
}

.container{
  position: relative;
  max-width: 430px;
  width: 100%;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin: 0 20px;
}

.container .forms{
  display: flex;
  align-items: center;
  height: 440px;
  width: 200%;
  transition: height 0.2s ease;
}


.container .form{
  width: 50%;
  padding: 30px;
  background-color: #fff;
  transition: margin-left 0.18s ease;
}

.container.active .login{
  margin-left: -50%;
  opacity: 0;
  transition: margin-left 0.18s ease, opacity 0.15s ease;
}

.container .signup{
  opacity: 0;
  transition: opacity 0.09s ease;
}
.container.active .signup{
  opacity: 1;
  transition: opacity 0.2s ease;
}

.container.active .forms{
  height: 600px;
}
.container .form .title{
  position: relative;
  font-size: 27px;
  font-weight: 600;
}

.form .title::before{
  content: '';
  position: absolute;
  left: 0;
  bottom: 0;
  height: 3px;
  width: 30px;
  background-color: #4070f4;
  border-radius: 25px;
}

.form .input-field{
  position: relative;
  height: 40px;
  width: 100%;
  margin-top: 30px;
}

.input-field input{
  position: absolute;
  height: 100%;
  width: 100%;
  padding: 0 35px;
  border: none;
  outline: none;
  font-size: 16px;
  border-bottom: 2px solid #ccc;
  border-top: 2px solid transparent;
  transition: all 0.2s ease;
}

.input-field input:is(:focus, :valid){
  border-bottom-color: #4070f4;
}

.input-field i{
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
  font-size: 23px;
  transition: all 0.2s ease;
}

.input-field input:is(:focus, :valid) ~ i{
  color: #4070f4;
}

.input-field i.icon{
  left: 0;
}
.input-field i.showHidePw{
  right: 0;
  cursor: pointer;
  padding: 10px;
}

.form .checkbox-text{
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
}

.checkbox-text .checkbox-content{
  display: flex;
  align-items: center;
}

.checkbox-content input{
  margin-right: 10px;
  accent-color: #4070f4;
}

.form .text{
  color: #333;
  font-size: 14px;
}

.form a.text{
  color: #4070f4;
  text-decoration: none;
}
.form a:hover{
  text-decoration: underline;
}

.form .button{
  margin-top: 35px;
}

.form .button input{
  border: none;
  color: #fff;
  font-size: 17px;
  font-weight: 500;
  letter-spacing: 1px;
  border-radius: 6px;
  background-color: #4070f4;
  cursor: pointer;
  transition: all 0.3s ease;
}

.button input:hover{
  background-color: #265df2;
}

.form .login-signup{
  margin-top: 30px;
  text-align: center;
}

.center-form {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

</style>

