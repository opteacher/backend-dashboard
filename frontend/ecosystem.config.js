module.exports = {
    apps : [{
        name   : "backend-dashboard",
        script : "app.js",
        watch  : true,
        env_production : {
            PORT : 4000
        }
    }]
};