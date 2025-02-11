<script>
import axios from "axios";

export default {
  name: 'WidgetData',
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
                  property732: item.PROPERTY_732
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

<template>

</template>

<style scoped>

</style>