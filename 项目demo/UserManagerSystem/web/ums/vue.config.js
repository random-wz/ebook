module.exports = {
    devServer: {
        proxy: {
            '/try/ajax/json_demo.json': {
                target: 'https://www.runoob.com',
                changeOrigin: true
            }
        }
    }
}