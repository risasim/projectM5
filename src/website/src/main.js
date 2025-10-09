import router from './router'
import { createApp } from 'vue'
import App from './App.vue'
import IconEdit from './components/icons/IconEdit.vue'
import IconBack from './components/icons/IconBack.vue'


const app = createApp(App)
app.component('IconEdit', IconEdit)
app.component('IconBack', IconBack)
app.use(router)
app.mount('#app')