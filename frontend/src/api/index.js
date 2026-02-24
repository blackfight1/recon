import axios from 'axios'

const api = axios.create({
    baseURL: '/api',
    timeout: 30000
})

// 请求拦截器
api.interceptors.request.use(
    config => config,
    error => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
    response => response.data,
    error => {
        console.error('API Error:', error)
        return Promise.reject(error)
    }
)

export default {
    // 统计信息
    getStats() {
        return api.get('/stats')
    },

    // 目标管理
    getTargets() {
        return api.get('/targets')
    },
    createTarget(data) {
        return api.post('/targets', data)
    },
    updateTarget(id, data) {
        return api.put(`/targets/${id}`, data)
    },
    deleteTarget(id) {
        return api.delete(`/targets/${id}`)
    },
    triggerScan(id) {
        return api.post(`/targets/${id}/scan`)
    },

    // 子域名
    getSubdomains(params) {
        return api.get('/subdomains', { params })
    },
    getSubdomainsByTarget(targetId) {
        return api.get(`/subdomains/target/${targetId}`)
    },

    // 变更日志
    getChanges(params) {
        return api.get('/changes', { params })
    },
    getChangesByTarget(targetId) {
        return api.get(`/changes/target/${targetId}`)
    },

    // 任务
    getTasks(params) {
        return api.get('/tasks', { params })
    }
}
