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
        async deleteUser(publicId) {
            await api.deleteUsers(publicId)
            this.updateUserList()
        },
        async updateUserList() {
            this.users = await api.getUsers()
        }
    },
    created() {
        this.updateUserList()
    },
    mounted() {
    }
};

const app = createApp(App)
app.mount('#app')