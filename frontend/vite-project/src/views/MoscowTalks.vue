<template>
  <main id="WidgetFormPage">
    <h1 class="title">Говорит, Москва!</h1>
    <h2 class="title_second">Настройте приложение</h2>
    <div class="form-container">
      <form @submit.prevent="submitForm" class="form-content">

        <div class="form-group">
          <input
              type="text"
              v-model="form.Prompt"
              class="form-input"
              placeholder="Prompt here"
              required
          />
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.NaimenovanieKompanii"
              class="form-input"
              placeholder="Название вашей компании"
              @input="showDropdown = true"
          />
          <ul v-if="filteredCompanies.length && showDropdown" class="dropdown">
            <li
                v-for="company in filteredCompanies"
                :key="company.ID"
                @click="selectCompany(company)"
                class="dropdown-item"
            >
              {{ company.TITLE }}
            </li>
          </ul>
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.CompanyInfo"
              class="form-input"
              placeholder="Информация о компании"
              required
          />
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.CompanyContacts"
              class="form-input"
              placeholder="Контакты компании"
              required
          />
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.BotTargets"
              class="form-input"
              placeholder="Цель бота"
              required
          />
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.WhatDoNotDo"
              class="form-input"
              placeholder="Что запрещено делать"
              required
          />
        </div>

        <div class="form-group">
          <input
              type="text"
              v-model="form.AdditionalInfo"
              class="form-input"
              placeholder="Дополнительная информация"
              required
          />
        </div>

        <!-- Кнопка теперь вне form-group -->
        <button type="submit" class="submit-button" :class="{ 'submitted': isSubmitted }" :disabled="isSubmitted">
          {{ isSubmitted ? "Сохранено!" : "Сохранить данные" }}
        </button>

      </form>
    </div>
  </main>
</template>

<script>
import axios from 'axios';

export default {
  name: 'WidgetFormPage',
  data() {
    return {
      allCompanies: [],
      showDropdown: false,
      isSubmitted: false,
      fieldPositions: {
        Prompt: { x: 400, y: 400, z: 0 },
        CompanyInfo: { x: 400, y: 550, z: 0 },
        CompanyContacts: { x: 450, y: 700, z: 0 },
        BotTargets: { x: 500, y: 850, z: 0 },
        WhatDoNotDo: { x: 550, y: 1000, z: 0 },
        AdditionalInfo: { x: 600, y: 1150, z: 0 }
      },
      form: {
        Prompt :'',
        NaimenovanieKompanii: '',
        CompanyInfo: '',
        CompanyContacts: '',
        BotTargets: '',
        WhatDoNotDo: '',
        AdditionalInfo: ''
      },
      fieldLabels: {
        Prompt :'Prompt here',
        NaimenovanieKompanii: 'Наименование компании',
        CompanyInfo: 'Информация о компании',
        CompanyContacts: 'Контакты компании',
        BotTargets: 'Цель бота',
        WhatDoNotDo: 'Что запрещено делать',
        AdditionalInfo: 'Дополнительная информация'
      }
    };
  },
  computed: {
    filteredCompanies() {
      const searchTerm = this.form.NaimenovanieKompanii.toLowerCase();
      return this.allCompanies.filter(company =>
          company.TITLE.toLowerCase().includes(searchTerm)
      );
    }
  },
  created() {
    this.fetchDealData();
  },
  methods: {
    fetchDealData() {
      const url = 'https://crmconsulting-api.ru/api/widget_data';
      axios
          .get(url)
          .then(response => {
            this.allCompanies = response.data;
          })
          .catch(error => {
            console.error('Ошибка при получении данных о компаниях:', error);
          });
    },

    submitForm() {
      console.log('Отправляемые данные:', this.form);
      const url = 'https://crmconsulting-api.ru/api/message_event';

      axios.post(url, { ...this.form, from_frontend: true }, {
        headers: {
          'X-Frontend-Key': 'your-secure-key'
        }
      })
          .then(response => {
            console.log('Данные формы отправлены успешно:', response.data);
            this.isSubmitted = true;
            setTimeout(() => {
              this.isSubmitted = false;
            }, 2000);
          })
          .catch(error => {
            console.error('Ошибка при отправке данных формы:', error);
          });
    },
    selectCompany(company) {
      this.form.NaimenovanieKompanii = company.TITLE;
      this.showDropdown = false;
    },
    getFieldPosition(key) {
      const position = this.fieldPositions[key];
      return {
        position: 'absolute',
        left: `${position.x}px`,
        top: `${position.y}px`,
        zIndex: position.z
      };
    }
  }
};
</script>

<style>
body {
  background-color: #f4f4f9;
  font-family: 'Roboto', sans-serif;
}

.form-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px; /* Отступы между элементами */
}

.form-input {
  width: 200%; /* Растягиваем input на 80% ширины */
  height: 8rem; /* Увеличиваем высоту в 2 раза */
  max-width: 600px;
  margin-top: 20px; /* Сдвигаем вниз на 30px */
  font-size: 1.2rem; /* Немного увеличиваем шрифт */
  padding: 30px; /* Добавляем отступ внутри */
}

.submit-button {
  width: 80%;
  max-width: 300px;
  margin-top: 30px; /* Сдвигаем кнопку вверх, чтобы приблизить к инпутам */
}


.submit-button:hover {
  background-color: #45a049;
}

.submit-button.submitted {
  background-color: #81c784;
}

.dropdown {
  z-index: 1;
  position: absolute;
  left: calc(35% - 200px);
  top: 160px;
  background-color: white;
  border: 1px solid #ccc;
  max-height: 200px;
  overflow-y: auto;
  list-style: none;
  margin-top: 5px;
  padding: 0;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.dropdown-item {
  padding: 15px;
  font-size: 1.1rem;
  cursor: pointer;
}

.dropdown-item:hover {
  background-color: #e8f5e9;
}
</style>
