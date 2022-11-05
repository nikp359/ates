import { createApp } from 'vue'
import api from "./api.js";

const App = {
    data() {
        return {
            count: 0,
            message: "Hello dude",
            users: []
        }
    },
    methods: {
        edit(publicId) {
            console.log(publicId)
        }
    },
    async created() {
        this.users = await api.getUsers()
    },
    mounted() {
    }
};

const app = createApp(App)
app.mount('#app')