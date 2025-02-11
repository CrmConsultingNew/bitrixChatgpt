// store.js

import { createStore } from 'vuex';

export default createStore({
  modules: {
    auth: {
      state: {
        loggedIn: false,
      },
      mutations: {
        setLoggedIn(state, loggedIn) {
          state.loggedIn = loggedIn;
        },
      },
      actions: {
        login({ commit }) {
          // Perform your login logic here
          // After successful login, commit 'setLoggedIn' mutation with true
          commit('setLoggedIn', true);
        },
        logout({ commit }) {
          // Perform logout logic here
          // After successful logout, commit 'setLoggedIn' mutation with false
          commit('setLoggedIn', false);
        },
      },
    },
  },
});
