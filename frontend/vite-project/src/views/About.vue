<template>
  <main id="Home-page">
    <h1 class="title">Список компаний</h1>

    <!-- Date Filters -->
    <div class="filter-container">
      <div>
        <label for="date-from-filter">Дата от: </label>
        <input type="date" id="date-from-filter" v-model="dateFrom" />
      </div>
      <div>
        <label for="date-to-filter">Дата до: </label>
        <input type="date" id="date-to-filter" v-model="dateTo" />
      </div>
      <button @click="applyDateFilter" :class="{'apply-button': true, 'apply-button-active': applyButtonClicked}">
        Применить
      </button>
    </div>

    <!-- Text Filter -->
    <div class="filter-container">
      <label for="title-filter">Фильтр по названию:</label>
      <input type="text" id="title-filter" v-model="titleFilter" placeholder="Введите название..." />
    </div>

    <!-- Company List -->
    <div class="table-container">
      <ul class="table">
        <li v-for="company in filteredCompanies" :key="company.ID" class="list-item">
          <div class="item-details">
            <div class="company-header">
              <p @click="fetchAndToggleMenu(company.ID)" class="title-link">{{ company.TITLE }}</p>
              <button
                  @click="generateDoc(company.ID, company.TITLE)"
                  v-if="activeItem === company.ID"
                  class="generate-document-button"
              >
                Сформировать документ
              </button>
            </div>
            <div v-if="activeItem === company.ID" class="additional-info">
              <p>ID компании: {{ company.ID }}</p>
              <p>У компании есть email? {{ company.HAS_EMAIL }} (Y = да, N = нет)</p>
              <p v-if="company.Emails && company.Emails.length > 0">Email:</p>
              <ul v-if="company.Emails && company.Emails.length > 0">
                <li v-for="(email, index) in company.Emails" :key="index" class="email-item">{{ email.VALUE }}</li>
              </ul>
              <p v-if="!company.Emails || company.Emails.length === 0">Нет email-адресов.</p>
              <p>ID ответственного: {{ company.ASSIGNED_BY_ID }}</p>

              <!-- Суммы ДДС -->
              <p>Суммы ДДС:</p>
              <div v-if="processesData[company.ID]">
                <ul>
                  <li v-for="(item, index) in processesData[company.ID]" :key="index" class="dds-item">
                    {{ item.property628 }} от {{ item.property638 }}
                  </li>
                </ul>
              </div>

              <!-- Суммы актов -->
              <p>Суммы актов:</p>
              <div v-if="opportunityData[company.ID] && opportunityData[company.ID].length > 0">
                <ul>
                  <li v-for="(item, index) in opportunityData[company.ID]" :key="index" class="dds-item">
                    {{ item.opportunity }} от {{ item.ufCrm26_1712128088 }}
                  </li>
                </ul>
              </div>
              <p v-else>Нет данных по суммам актов.</p>


              <!-- Mail Button -->
              <div class="mail-container" v-if="activeItem === company.ID">
                <p class="mail-text">Отправить письмо клиенту</p>
                <div v-for="(email, index) in company.Emails" :key="index">
                  <button
                      class="mail-button"
                      @click="sendMail(email.VALUE, company.ID)"
                  >
                    Отправить на {{ email.VALUE }}
                  </button>
                </div>
              </div>
              <p
                  v-if="emailStatus[company.ID]"
                  :class="{
                  'email-success': emailStatus[company.ID] === 'success',
                  'email-failure': emailStatus[company.ID] === 'failure'
                }"
              >
                {{ emailStatus[company.ID] === 'success' ? 'Файл успешно отправлен' : 'Не удалось отправить файл' }}
              </p>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </main>
</template>

<script>
import axios from 'axios';

export default {
  name: 'AboutPage',
  data() {
    return {
      companies: [],
      activeItem: null,
      dateFrom: '', // Начальная дата
      dateTo: '',   // Конечная дата
      titleFilter: '', // Текстовый фильтр для поля TITLE
      processesData: {}, // Хранение дополнительных данных для каждой компании
      opportunityData: {}, // Хранение данных по суммам для каждой компании
      emailStatus: {}, // Статус отправки email для каждой компании
      applyButtonClicked: false, // Состояние кнопки Применить
    };
  },
  created() {
    this.fetchCompanies(); // Получение списка компаний при загрузке компонента
  },
  methods: {
    fetchCompanies() {
      let url = `https://crmconsulting-api.ru/api/companies`;

      axios
          .get(url)
          .then((response) => {
            console.log('Companies data received:', response.data);
            this.companies = response.data; // Сохранение полученных данных в companies
          })
          .catch((error) => {
            console.error('Error fetching companies data:', error);
          });
    },
    fetchAndToggleMenu(companyID) {
      // Переключение активности элемента
      this.activeItem = this.activeItem === companyID ? null : companyID;

      if (this.activeItem) {
        const url = `https://crmconsulting-api.ru/api/company?id=${companyID}`;
        console.log('Fetching company data for company ID:', companyID, 'from URL:', url);

        axios
            .get(url)
            .then((response) => {
              console.log('Company data received for company ID', companyID, ':', response.data);

              // Обновление данных компании
              this.companies = this.companies.map((company) =>
                  company.ID === companyID ? { ...company, Emails: response.data.EMAIL } : company
              );
            })
            .catch((error) => {
              console.error('Error fetching company data for company ID', companyID, ':', error);
            });

        // Получение данных по процессам
        const processesUrl = `https://crmconsulting-api.ru/api/processes?id=${companyID}`;
        console.log('Fetching processes for company ID:', companyID, 'from URL:', processesUrl);

        axios
            .get(processesUrl)
            .then((response) => {
              console.log('Processes data received for company ID', companyID, ':', response.data);
              // Обновление данных processesData
              const data = response.data ? response.data.map((item) => {
                return {
                  property628: item.PROPERTY_628,
                  property638: item.PROPERTY_638
                };
              }) : [];
              this.processesData = {
                ...this.processesData,
                [companyID]: data,
              };
              console.log('Updated processesData:', this.processesData);
            })
            .catch((error) => {
              console.error('Error fetching processes data for company ID', companyID, ':', error);
            });

        // Получение данных по суммам
        const opportunityUrl = `https://crmconsulting-api.ru/api/items?id=${companyID}`;
        axios
            .get(opportunityUrl)
            .then((response) => {
              console.log('Opportunity data received for company ID:', companyID, response.data);

              if (response.data && response.data.length > 0) {
                const data = response.data.map((item) => {
                  return {
                    opportunity: item.opportunity,
                    ufCrm26_1712128088: item.ufCrm26_1712128088
                  };
                });

                this.opportunityData = {
                  ...this.opportunityData,
                  [companyID]: data,
                };
                console.log('Updated opportunityData:', this.opportunityData);
              } else {
                console.log('No opportunity data found for company ID:', companyID);
                this.opportunityData = {
                  ...this.opportunityData,
                  [companyID]: [],
                };
              }
            })
            .catch((error) => {
              console.error('Error fetching opportunity data for company ID:', companyID, error);
            });
      }
    },
    generateDoc(companyID, title) {
      if (!this.dateFrom || !this.dateTo) {
        alert('Пожалуйста, выберите оба значения даты.');
        return;
      }
      console.log(`Generating document for Company ID: ${companyID}, Title: ${title}, Date From: ${this.dateFrom}, Date To: ${this.dateTo}`);
      axios({
        url: `https://crmconsulting-api.ru/api/generate_docx?id=${companyID}&title=${encodeURIComponent(title)}&date_from=${this.dateFrom}&date_to=${this.dateTo}`,
        method: 'GET',
        responseType: 'blob', // Важно для получения файла
      })
          .then((response) => {
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const currentDate = new Date().toLocaleDateString('ru-RU').replace(/\//g, '.');
            const fileName = `Акты сверок от CRM Consulting для ${title} (от ${currentDate}).docx`;
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', fileName);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
          })
          .catch((error) => {
            console.error('Error generating the document:', error);
          });
    },
    sendMail(email, companyID) {
      console.log('Sending mail to:', email);

      axios
          .post('https://crmconsulting-api.ru/api/send_email', { email })
          .then((response) => {
            console.log('Email sent successfully:', response.data);
            this.emailStatus[companyID] = 'success';
          })
          .catch((error) => {
            console.error('Error sending email:', error);
            this.emailStatus[companyID] = 'failure';
          });
    },
    applyDateFilter() {
      console.log('Applying date filter from:', this.dateFrom, 'to:', this.dateTo);
      this.applyButtonClicked = true;
      setTimeout(() => {
        this.applyButtonClicked = false;
      }, 1000); // Кнопка меняет цвет на 1 секунду
    },
  },
  computed: {
    filteredCompanies() {
      // Применение текстового фильтра по названию компании
      let filtered = this.companies;
      if (this.titleFilter) {
        filtered = filtered.filter((company) =>
            company.TITLE.toLowerCase().includes(this.titleFilter.toLowerCase())
        );
      }
      console.log('Filtered companies:', filtered);
      return filtered;
    },
  },
};
</script>

<style>
/* Общие стили для разделов */
.title {
  font-size: 2.5rem;
  color: #333;
  text-align: center;
  margin-bottom: 20px;
}

.filter-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.filter-container label {
  font-weight: bold;
  color: #555;
}

.filter-container input[type="date"],
.filter-container input[type="text"] {
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #ccc;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.filter-container input[type="date"]:focus,
.filter-container input[type="text"]:focus {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.3);
}

.apply-button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.8s ease;
}

.apply-button:hover {
  background-color: #0056b3;
}

/* Активное состояние кнопки "Применить" */
.apply-button-active {
  background-color: #ffa07a;
}

/* Стили для ссылок на заголовки */
.title-link {
  cursor: pointer;
  color: #007bff;
  font-weight: bold;
  text-decoration: none;
  font-size: 1.2rem;
  transition: color 0.3s ease;
}

.title-link:hover {
  color: #0056b3;
}

/* Стили для элемента информации */
.item-details {
  padding: 15px;
  border-radius: 10px;
  background-color: #f9f9f9;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.additional-info {
  margin-top: 15px;
  padding: 10px;
  background-color: #fff;
  border-radius: 8px;
  border-left: 4px solid #007bff;
}

.dds-item {
  color: #28a745;
  font-weight: bold;
  list-style-type: none;
  margin: 5px 0;
}

.email-item {
  color: #007bff;
  list-style-type: none;
  margin: 5px 0;
}

/* Стили для кнопки отправки почты */
.mail-container {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.mail-text {
  margin-right: 8px; /* Расстояние между текстом и иконкой */
  font-weight: bold;
  color: #555;
}

.mail-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2em;
  width: 24px;
  height: 24px;
  background-size: cover;
  background-position: center;
  background-image: url('../assets/mailicon.png'); /* Замените на ваш URL иконки */
  transition: transform 0.3s ease;
}

.mail-button:hover {
  transform: scale(1.1);
}

/* Убираем маркеры из списка */
.list-item {
  list-style-type: none;
  margin-left: 0;
  padding-left: 0;
  background-color: transparent;
  position: relative;
}

/* Расположение таблицы */
.table-container {
  text-align: left;
  width: 80%;
  margin: 0 auto;
}

.table {
  padding-left: 0;
}

.table li {
  margin-bottom: 15px;
}

/* Стили для заголовка компании */
.company-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Стили для кнопки создания документа */
.generate-document-button {
  background-color: #28a745;
  color: white;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  border-radius: 8px;
  font-weight: bold;
  transition: background-color 0.3s ease;
  position: relative;
}

.generate-document-button:hover {
  background-color: #218838;
}

.generate-document-button::after {
  content: "";
  position: absolute;
  top: 50%;
  right: -10px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 6px 0 6px 10px;
  border-color: transparent transparent transparent #28a745;
  transform: translateY(-50%);
}

/* Статус отправки email */
.email-success {
  color: #28a745;
  font-weight: bold;
  margin-top: 10px;
}

.email-failure {
  color: #dc3545;
  font-weight: bold;
  margin-top: 10px;
}
</style>
