<template>
  <div class="container">
    <div class="page-container">
      <img src="../../public/harizma.jpg" alt="Image Description" class="logo">
    </div>
    <div class="feedback-form-container">
      <h2 class="form-heading">Оцени Харизму</h2>
      <form @submit.prevent="submitFeedback">
        <div class="form-field">
          <label for="rating">Оценка (1-10): *</label>
          <select v-model="rating" required class="rating-select">
            <option value="" disabled hidden>Поставь оценку</option>
            <option v-for="num in 10" :value="num" :key="num">{{ num }}</option>
          </select>
          <div class="mandatory-note" v-if="ratingError">Оценка обязательна</div>
        </div>
        <div class="form-field">
          <label for="comment">Чем вас удивить в следующий раз?</label>
          <textarea id="comment" v-model="comment"></textarea>
        </div>
        <div class="form-field">
          <input type="submit" value="Отправить">
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      code: '',  // Store the ID
      rating: '',
      comment: '',
      ratingError: false
    }
  },
  mounted() {
    const urlParams = new URLSearchParams(window.location.search);
    this.code = urlParams.get('code');
  },
  methods: {
    submitFeedback() {
      if (this.rating === '' || !this.code) {
        this.ratingError = true; // Display an error message if ID or rating is missing
        return;
      }
      this.ratingError = false;

      // Redirect immediately after validating the input
      window.location.href = 'https://b24-yeth0y.bitrix24site.ru/empty_jekf/';

      // Continue with the fetch request to send data
      const feedbackData = {
        code: this.code,  // Include the ID in the feedback data
        rating: this.rating.toString(),
        comment: this.comment
      };

      fetch('https://harizma-service.ru/api/user-form', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(feedbackData)
      })
          .then(response => {
            if (!response.ok) {
              throw new Error('Network response was not OK');
            }
            console.log('Feedback sent successfully');
          })
          .catch(error => {
            console.error('Error:', error);
          });
    }
  }
}
</script>



<style scoped>
.rating-select {
  height: 50px;
}
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}
.page-container {
  margin-bottom: 20px;
}
.logo {
  width: 100%;
  max-width: 500px;
  height: auto;
}
.feedback-form-container {
  width: 90%;
  max-width: 400px;
  box-sizing: border-box;
}
.form-heading {
  font-weight: normal;
  color: #FFB500;
  text-align: center;
  margin-bottom: 45px;
}
.form-field {
  margin-bottom: 30px;
}
.form-field label {
  display: block;
  width: 100%;
  color: white;
}
.mandatory-note {
  color: red;
}
.form-field select, .form-field textarea {
  display: block;
  width: 100%;
  color: white;
  background-color: #333;
  border: none;
}
.form-field input[type="submit"] {
  width: 100%;
  height: 40px;
  color: white;
  background-color: #333;
  border: 1px solid #D3D3D3;
}
.form-field input[type="submit"]:hover {
  background-color: #FF7F50;
}
#comment {
  width: 100%;
  height: 80px;
}
@media only screen and (min-width: 768px) {
  .feedback-form-container {
    width: 70%;
  }
}
@media only screen and (min-width: 1024px) {
  .feedback-form-container {
    width: 50%;
  }
}
</style>
