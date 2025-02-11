<template>
  <Sidebar />

  <main id="about-page">
    <h1>About</h1>
    <p>This is the about page</p>

    <div>
      <button @click="toggleMenu">Toggle Menu</button>
      <ul v-if="showMenu">
        <!-- Menu items -->
        <li v-for="(item, index) in menuItems" :key="index">
          <button @click="toggleMenuItem(item)">{{ item.name }}</button>
          <ul v-if="item.showSubmenu">
            <!-- Submenu items -->
            <li>Subitem 1</li>
            <li>Subitem 2</li>
            <li>Subitem 3</li>
          </ul>
        </li>
      </ul>
    </div>
  </main>
</template>

<script>
import Sidebar from '@/components/Sidebar.vue'

export default {
  components: { Sidebar },
  data() {
    return {
      showMenu: false,
      menuItems: [
        { name: 'Item 1', showSubmenu: false },
        { name: 'Item 2', showSubmenu: false },
        { name: 'Item 3', showSubmenu: false },
        { name: 'Item 4', showSubmenu: false }
      ]
    }
  },
  methods: {
    toggleMenu() {
      this.showMenu = !this.showMenu
    },
    toggleMenuItem(item) {
      // Toggle the clicked item's submenu
      item.showSubmenu = !item.showSubmenu

      // Close submenus for other items
      this.menuItems.forEach((otherItem) => {
        if (otherItem !== item) {
          otherItem.showSubmenu = false
        }
      })
    }
  }
}
</script>