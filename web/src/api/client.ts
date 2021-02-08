import Axios from 'axios'

const httpBaseURL = process.env.VUE_APP_HTTP_BASE_URL
const webSocketBaseURL = process.env.VUE_APP_WEBSOCKET_BASE_URL

const client = Axios.create({ baseURL: httpBaseURL, withCredentials: true })

export default client
export { client, httpBaseURL, webSocketBaseURL }
