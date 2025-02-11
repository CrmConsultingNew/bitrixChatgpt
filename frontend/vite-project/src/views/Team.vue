<template>
  <Sidebar />

  <a href="http://localhost:5173/team"> Маркиза в коробе </a>

  <div id="calculator">
    <div id="title" name="title" class="title">
      <h3>Маркиза в коробе</h3>
    </div>

    <div class="display">
      <p>Высота Y</p>
      <input v-model="firstNumber" type="text" />
      <p>Ширина X</p>
      <input v-model="secondNumber" type="text" />

      <div class="blocks">
        <div class="info_block" id="info_block" name="info_block">
          <p>Короб</p>
          {{ secondNumber - 10 }} см
        </div>
        <div class="info_block" id="info_block" name="info_block">
          <p>Труба</p>
          {{ secondNumber - 9.5 }} см
        </div>
        <div class="info_block" id="info_block" name="info_block">
          <p>Ткань</p>
          {{ firstNumber - 19 + 50 }} см
        </div>
        <div class="info_block" id="info_block" name="info_block">
          <p>Нити</p>
          {{ Math.ceil((firstNumber - 30) / 45) }} шт
        </div>
      </div>
    </div>
    <div class="buttons">
      <!--      <button @click="add">Add</button>
      <button @click="subtract">Subtract</button>
      <button @click="multiply">Multiply</button>
      <button @click="divide">Divide</button>-->
    </div>
    <div class="result">
      <!--      Result: {{ result }}-->
    </div>
  </div>
</template>

<script>
import Sidebar from '@/components/Sidebar.vue'

export default {
  components: { Sidebar },
  data() {
    return {
      firstNumber: '',
      secondNumber: '',
      result: null
    }
  },
  methods: {
    addResultToFirstElem() {
      this.result = parseFloat(this.firstNumber) + parseFloat(this.secondNumber)
    },
    add() {
      this.result = parseFloat(this.firstNumber) + parseFloat(this.secondNumber)
    },
    subtract() {
      this.result = parseFloat(this.firstNumber) - parseFloat(this.secondNumber)
    },
    multiply() {
      this.result = parseFloat(this.firstNumber) * parseFloat(this.secondNumber)
    },
    divide() {
      const num2 = parseFloat(this.secondNumber)
      if (num2 !== 0) {
        this.result = parseFloat(this.firstNumber) / num2
      } else {
        this.result = 'Error: Division by zero'
      }
    }
  }
}
</script>

<style scoped>
.title {
  padding: 5px 20px 30px 10px;
}

.blocks {
  border-radius: 5px;
  border-color: black;
  display: flex;
  padding: 10px;
}

.info_block {
  width: 500px;
  height: 80px;
  margin: 10px; /* Add margin to create space around the info blocks */
  padding: 10px; /* Optionally, add padding for more spacing */
  border: 1px solid cornflowerblue; /* Add a border for visual separation */
  background-color: deepskyblue;
}
.info_block p {
  margin-bottom: 10px; /* Add margin to the bottom of <p> elements */
}

#calculator {
  width: 500px;
  height: 500px;
  margin: 0 auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  text-align: center;
}

.display {
  font-size: 17px;
  margin-bottom: 10px;
}

input {
  width: 100px;
  height: 50px;
  text-align: center;
}

button {
  width: 80px;
  height: 30px;
  margin: 5px;
  font-size: 16px;
  padding: 20px;
}

.result {
  font-size: 20px;
  margin-top: 10px;
}
</style>