const Koa = require("koa")
const bodyparser = require("koa-bodyparser")
const json = require("koa-json")
const logger = require("koa-logger")
const cors = require("koa2-cors")
// [FRONTEND_ENV_BEG]
const path = require("path")
const statc = require("koa-static")
const views = require("koa-views")
// [FRONTEND_ENV_END]
const config = require("./config/server")
const router = require("./routes/index")

const app = new Koa()

// 跨域配置
app.use(cors())

// 路径解析
app.use(bodyparser())

// json解析
app.use(json())

// 日志输出
app.use(logger())

// [FRONTEND_ENV_BEG]
// 指定静态目录
app.use(statc(path.join(__dirname, "public")))

// 指定页面目录
app.use(views("./", {extension: "html"}))
// [FRONTEND_ENV_END]

// 路径分配
app.use(router.routes(), router.allowedMethods())

// 错误跳转
app.use(ctx => {
    ctx.status = 404
    ctx.body = "error"
})

app.listen(process.env.PORT || config.port)