export default {
    async getUsers() {
        try {
            const response = await axios.get('/api/users');
            return response.data
        } catch (error) {
            console.error(error);
        }
    },
    async deleteUsers(publicId) {
        try {
            const response = await axios.delete('/api/users', {
                data: {
                    public_id: publicId
                }
            });
            return response.data
        } catch (error) {
            console.error(error);
        }
    }
}