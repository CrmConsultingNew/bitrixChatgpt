<template>
  <div class="container">
    <div class="page-container">
      <img src="../../public/harizma.jpg" alt="Image Description" class="logo">
    </div>
    <div class="thanks">
      <h1>Спасибо за оценку!</h1>
    </div>
  </div>


</template>


<style scoped>

.thanks{
  color: white;
  font-weight: normal;
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

.form-field label {
  display: block;
  width: 100%;
  color: white;
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


</style>
<script setup>
</script>