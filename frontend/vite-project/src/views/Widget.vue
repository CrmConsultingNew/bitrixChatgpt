<template>
  <main id="WidgetFormPage">
    <h1 class="title">Аудит</h1>
    <h2 class="title_second">Заполните данные компании</h2>
    <div class="form-container">
      <form @submit.prevent="submitForm" class="form-content">
        <div class="form-group">
          <input
              type="text"
              v-model="form.NaimenovanieKompanii"
              id="NaimenovanieKompanii"
              class="form-input"
              placeholder="Введите название компании"
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

        <div v-for="key in filteredFormFields" :key="key" class="form-group">
          <template v-if="['Nisha', 'Geographia', 'KonechniyProduct', 'OsnovnieKanaliProdazh', 'Sezonnost', 'Zapros', 'CeliDoKontsaGoda', 'OzhidaniyaOtSotrudnichestva'].includes(key)">
            <textarea
                v-model="form[key]"
                :id="key"
                class="form-input"
                :style="getFieldPosition(key)"
                :placeholder="fieldLabels[key]"
                rows="2"
                required
            ></textarea>
          </template>
          <template v-else>
            <input
                type="text"
                v-model="form[key]"
                :id="key"
                class="form-input"
                :style="getFieldPosition(key)"
                :placeholder="fieldLabels[key]"
                @input="checkNumericInput(key)"
                required
            />
          </template>
        </div>

        <button type="submit" class="submit-button" :class="{ 'submitted': isSubmitted }">Отправить</button>
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
        INN: { x: 400, y: 290, z: 0 },
        Nisha: { x: 700, y: 200, z: 0 },
        Geographia: { x: 250, y: 380, z: 0 },
        KonechniyProduct: { x: 650, y: 370, z: 0 },
        OborotKompanii: { x: 750, y: 470, z: 0 },
        SrednemesyachnayaViryuchka: { x: 400, y: 510, z: 0 },
        KolichestvoSotrudnikov: { x: 50, y: 500, z: 0 },
        KolichestvoSotrudnikovOP: { x: 600, y: 720, z: 0 },
        EstROP: { x: 100, y: 710, z: 0 },
        EstHR: { x: 300, y: 620, z: 0 },
        OsnovnieKanaliProdazh: { x: 380, y: 820, z: 0 },
        Sezonnost: { x: 20, y: 850, z: 0 },
        Zapros: { x: 720, y: 840, z: 0 },
        CeliDoKontsaGoda: { x: 600, y: 960, z: 0 },
        OzhidaniyaOtSotrudnichestva: { x: 150, y: 1000, z: 0 }
      },
      form: {
        NaimenovanieKompanii: '',
        INN: '',
        Nisha: '',
        Geographia: '',
        KonechniyProduct: '',
        OborotKompanii: '',
        SrednemesyachnayaViryuchka: '',
        KolichestvoSotrudnikov: '',
        KolichestvoSotrudnikovOP: '',
        EstROP: '',
        EstHR: '',
        OsnovnieKanaliProdazh: '',
        Sezonnost: '',
        Zapros: '',
        CeliDoKontsaGoda: '',
        OzhidaniyaOtSotrudnichestva: ''
      },
      fieldLabels: {
        NaimenovanieKompanii: 'Наименование компании',
        INN: 'ИНН',
        Nisha: 'Ниша',
        Geographia: 'География',
        KonechniyProduct: 'Конечный продукт',
        OborotKompanii: 'Оборот компании',
        SrednemesyachnayaViryuchka: 'Среднемесячная выручка',
        KolichestvoSotrudnikov: 'Количество сотрудников',
        KolichestvoSotrudnikovOP: 'Кол-во сотрудников ОП',
        EstROP: 'Есть РОП',
        EstHR: 'Есть HR',
        OsnovnieKanaliProdazh: 'Основные каналы продаж',
        Sezonnost: 'Сезонность',
        Zapros: 'Запрос',
        CeliDoKontsaGoda: 'Цели до конца года',
        OzhidaniyaOtSotrudnichestva: 'Ожидания от сотрудничества'
      }
    };
  },
  computed: {
    filteredCompanies() {
      const searchTerm = this.form.NaimenovanieKompanii.toLowerCase();
      return this.allCompanies.filter(company =>
          company.TITLE.toLowerCase().includes(searchTerm)
      );
    },
    filteredFormFields() {
      return Object.keys(this.form).filter(key => key !== 'NaimenovanieKompanii');
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

    checkNumericInput(key) {
      const numericFields = ['INN', 'OborotKompanii', 'SrednemesyachnayaViryuchka', 'KolichestvoSotrudnikov', 'KolichestvoSotrudnikovOP'];
      if (numericFields.includes(key)) {
        const value = this.form[key];
        if (isNaN(value) && value.trim() !== '') {
          alert('Данное поле должно содержать только числовые значения');
        }
      }
    },

    submitForm() {
      console.log('Отправляемые данные:', this.form);
      const url = 'https://crmconsulting-api.ru/api/form_data';
      axios.post(url, this.form)
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
  background-image: url('/images/worldwide-connection_53876-90461.jpg');
  background-size: cover;
  background-position: center;
}


#NaimenovanieKompanii {
  width: 40%; /* Увеличенная ширина для поля Название компании */
  margin-left: -200px; /* Сдвигаем поле на 50 пикселей влево */
}

.form-input {
  width: 23%; /* Полная ширина для остальных полей */
  height: auto;
  min-height: 60px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #ddd;
  font-size: 1rem;
  box-shadow: 0 4px 12px dodgerblue;
  transition: all 0.3s ease;
  resize: none; /* Запретить изменение размера текстового поля */
}

.title {
  text-align: center;
  color: #2c3e50;
  font-size: 2.5rem;
  margin-bottom: 20px;
}
.title_second {
  text-align: center;
  color: dodgerblue;
  font-size: 1.5rem;
  margin-bottom: 20px;
}

#WidgetFormPage {
  margin-left: 50px; /* Сдвигаем все элементы на 50 пикселей влево */
}

.form-container {
  background-color: #ffffff;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  margin: 20px auto;
  width: 80%;
  min-height: 1200px; /* Минимальная высота контейнера */
}

.form-input {
  width: 23%; /* Полная ширина для остальных полей */
  height: auto;
  min-height: 60px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #ddd;
  font-size: 1rem;
  box-shadow: 0 4px 12px dodgerblue;
  transition: all 0.3s ease;
}

#NaimenovanieKompanii {
  width: 40%; /* Увеличенная ширина для поля Название компании */
}

.textarea {
  height: 100px;
}

.form-input:focus {
  border-color: dodgerblue;
  box-shadow: 0 0 8px dodgerblue; /* Тень при фокусе */
}

.submit-button {
  width: 100%;
  display: block; /* Положение кнопки */
  margin: 900px auto;
  background-color: dodgerblue;
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: bold;
  transition: background-color 0.3s ease;
}

.submit-button:hover {
  background-color: #1e90ff;
}

.submit-button.submitted {
  background-color: lightsalmon;
}

.dropdown {
  z-index: 1;
  position: absolute; /* Задаем абсолютное позиционирование */
  left: calc(35% - 200px); /* Смещение влево, чтобы совпадало с полем */
  top: 160px; /* Подберите нужное значение в зависимости от высоты input */
  background-color: white;
  border: 1px solid #ccc;
  max-height: 200px;
  overflow-y: auto;
  list-style: none;
  margin-top: 5px;
  padding: 0;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}


.dropdown-item {
  padding: 10px;
  cursor: pointer;
}

.dropdown-item:hover {
  background-color: #f0f0f0;
}
</style>
